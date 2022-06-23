package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"strconv"
	"flag"
	"io/ioutil"
	"os"
	"os/user"
	"net/http"

	g8rv2 "github.com/g8rswimmer/go-twitter/v2"
)

var exitcode int = 0

type twidt int64
var twid_def twidt = 0
var next_max twidt = twid_def
var next_since twidt = twid_def
func print_id() {
	if uniqid != nil {
		err := uniqid.write()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
	fmt.Fprintf(os.Stderr, "--------\n-since_id=%d\n", next_since)
	fmt.Fprintf(os.Stderr,   "-max_id=%d\n", next_max)
}

//const onetimedefault = 10
const onetimemax = 100
const onetimemin_t = 5
const onetimemin_l = 1
const onetimemin_s = 10
var onetimemin = onetimemin_t
const sleepdot = 5

// // tweet/user hash for v2 api
// type tweetHash map[string] *gotwtr.Tweet
// type userHash map[string] *gotwtr.User

// TL type "enum"
type tltype int
const (
	tlnone tltype = iota
	tluser
	tlhome
	tlmention
	tlrtofme
	tllist
	tlsearch
	tlsearcha
)

type revtype bool
const (
	reverse revtype = true
	forward revtype = false
)

const (
	rsrecent string = "recent"
	rsall string = "all"
)

type idCheck map[string]bool
var uniqid idCheck = nil

func (c idCheck) checkID(id string) (exist bool) {
	if c[id] {
		return true
	} else {
		c[id] = true
		return false
	}
}

func (c idCheck) write() (err error) {
	bytes, _ := json.Marshal(c)
	err = ioutil.WriteFile("tempids.json", bytes, os.FileMode(0600))
	return err
}

func (c *idCheck) read() (err error) {
	*c = idCheck{}
	raw, err := ioutil.ReadFile("tempids.json")
	if err != nil {
		return err
	}
	if len(raw) == 0 {
		return nil
	}
	json.Unmarshal(raw, c)
	return nil
}

var twapi twSearchApi

func main(){
	var err error
	tLtypePtr := flag.String("get", "", "TLtype: user, mention, list, search")
	screennamePtr := flag.String("user", "", "twitter @ screenname")
	useridPtr := flag.String("userid", "0", "integer user Id")
	listnamePtr := flag.String("listname", "", "list name")
	listIDPtr := flag.String("listid", "0", "list ID")
	queryPtr := flag.String("query", "", "Query String")
	resulttypePtr := flag.String("restype", "", "result type: [recent]/all")
	countPtr := flag.Int("count", 0, "tweet count. 5-3200 ?")
	eachPtr := flag.Int("each", 0, "req count for each loop 5-100")
	max_idPtr := flag.Int64("max_id", 0, "starting tweet id")
	since_idPtr := flag.Int64("since_id", 0, "reverse start tweet id")
	reversePtr := flag.Bool("reverse", false, "reverse output. wait newest TL")
	loopsPtr := flag.Int("loops", 0, "get loop max")
	waitPtr := flag.Int64("wait", 0, "wait second for next loop")
	jsonPtr := flag.Bool("json", false, "dump json format")
	// nortPtr := flag.Bool("nort", false, "not include retweets")
	flag.Parse()
	var tLtype = *tLtypePtr
	var screenname = *screennamePtr
	var userid = *useridPtr
	var listname = *listnamePtr
	var listID = *listIDPtr
	var queryString = *queryPtr
	var resulttype = *resulttypePtr
	var count = *countPtr
	var eachcount = *eachPtr
	var max_id = twidt(*max_idPtr)
	var since_id = twidt(*since_idPtr)
	var reverseflag = *reversePtr
	var max_loop = *loopsPtr
	var waitsecond = *waitPtr
	var jsonp = *jsonPtr
	// var includeRTs = ! *nortPtr
	
	if flag.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "positional argument no need [%s]\n", flag.Arg(0))
		os.Exit(2)
	}

	var t tltype
	switch tLtype {
	case "user":    t = tluser
	//case "home":    t = tlhome
	case "mention": t = tlmention
	//case "rtofme":  t = tlrtofme
	case "list":    t = tllist
	case "search":  t = tlsearch
	case "":
		if listID != "0" || listname != "" {
			t = tllist
			tLtype = "list"
			fmt.Fprintln(os.Stderr, "assume -get=list")
		} else if userid != "0" || screenname != "" {
			t = tluser
			tLtype = "user"
			fmt.Fprintf(os.Stderr, "assume -get=%s\n", tLtype)
		} else if queryString != "" {
			t = tlsearch
			tLtype = "search"
			fmt.Fprintf(os.Stderr, "assume -get=%s\n", tLtype)
		} else if listID != "0" {
			t = tllist
			tLtype = "list"
			fmt.Fprintf(os.Stderr, "assume -get=%s\n", tLtype)
		} else {
			fmt.Fprintf(os.Stderr, "invalid type -get=%s\n", tLtype)
			os.Exit(2)
		}
	default:
		fmt.Fprintf(os.Stderr, "invalid type -get=%s\n", tLtype)
		os.Exit(2)
	}
	fmt.Fprintf(os.Stderr, "-get=%s\n", tLtype)
	
	twapi.client = connectTwitterApi()
	twapi.jsonp = jsonp
	
	switch t {
	case tluser: fallthrough
	case tlmention: fallthrough
	case tllist:
		if userid != "0" {
			fmt.Fprintf(os.Stderr, "user id=%s\n", userid)
			if (screenname != "") {
				fmt.Fprintln(os.Stderr, "screen name ignored.")
			}
		} else if screenname != "" {
			fmt.Printf("convert %s to userID\n", screenname)
			userid, _ = name2id(screenname)
			if userid == "" {
				os.Exit(2)
			}
			fmt.Fprintf(os.Stderr, "user id=%s %s\n", userid, screenname)
		} else if t != tllist {
			fmt.Fprintf(os.Stderr, "no user id\n")
			os.Exit(2)
		}
	default:
		if userid != "0" || screenname != "" {
			fmt.Fprintf(os.Stderr, "-get=%s no need userid/screenname\n", tLtype)
			os.Exit(2)
		}
	}

	switch t {
	case tllist:
		if max_id != 0 || since_id != 0 || reverseflag {
			fmt.Fprintln(os.Stderr, "-get=list can't handle -max_id, -since_id and -reverse")
			os.Exit(2)
		}
		if listID != "0" && listname != "" {
			fmt.Fprintln(os.Stderr, "list name ignored.")
			listname = ""
		}
		listID = listIDCheck(userid, listID, listname)
		fmt.Fprintf(os.Stderr, "listIDCheck returns '%s'\n", listID)
		if listID == "0" {
			fmt.Fprintln(os.Stderr, "-listid not specified")
			os.Exit(2)
		}
		onetimemin = onetimemin_l
		userid = listID
	default:
		if listID != "0" {
			fmt.Fprintln(os.Stderr, "-get=%s no need -listid", tLtype)
			os.Exit(2)
		}
	}
	
	switch t {
	case tlsearch:
		if queryString == "" {
			fmt.Fprintln(os.Stderr, "-query not specified")
			os.Exit(2)
		}
		switch {
		case strings.HasPrefix(rsrecent, resulttype):
			//default
		case strings.HasPrefix(rsall, resulttype):
			t = tlsearcha
		default:
			fmt.Fprintf(os.Stderr, "invalid -restype=%s\n", resulttype)
			os.Exit(2)
		}
		onetimemin = onetimemin_s
		err = uniqid.read()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		userid = queryString
	default:
		if queryString != "" {
			fmt.Fprintf(os.Stderr, "-get=%s no need -query\n", tLtype)
			os.Exit(2)
		}
		if resulttype != "" {
			fmt.Fprintf(os.Stderr, "-get=%s no need -restype=%s\n", tLtype, resulttype)
			os.Exit(2)
		}
	}		
	twapi.t = t

	fmt.Fprintf(os.Stderr, "count=%d\n", count)
	fmt.Fprintf(os.Stderr, "each=%d\n", eachcount)
	fmt.Fprintf(os.Stderr, "reverse=%v\n", reverseflag)
	fmt.Fprintf(os.Stderr, "loops=%d\n", max_loop)
	fmt.Fprintf(os.Stderr, "max_id=%d\n", max_id)
	fmt.Fprintf(os.Stderr, "since_id=%d\n", since_id)
	fmt.Fprintf(os.Stderr, "wait=%d\n", waitsecond)

	sgchn.sighandle()
	
	if reverseflag {
		if max_id != 0 {
			fmt.Fprintf(os.Stderr, "max id ignored when reverse\n")
		}
		if waitsecond <= 0 {
			waitsecond = 60
			fmt.Fprintf(os.Stderr, "wait default=%d (reverse)\n", waitsecond)
		}
		getReverseTLs(userid, count, max_loop, waitsecond, since_id)
	} else {
		if max_loop == 0 && since_id == 0 && count == 0 {
			count = onetimemin
			if count < 5 { count = 5 }
			fmt.Fprintf(os.Stderr, "set forward default count=%d\n", count)
		}
		if count != 0 && count < onetimemin {
			count = onetimemin
			fmt.Fprintf(os.Stderr, "set count=%d\n", count)
		}
		if eachcount != 0 && eachcount < onetimemin {
			eachcount = onetimemin
			fmt.Fprintf(os.Stderr, "set eachcount=%d\n", eachcount)
		} else if eachcount > onetimemax {
			eachcount = onetimemax
			fmt.Fprintf(os.Stderr, "set eachcount=%d\n", eachcount)
		}
			
		if max_id > 0 && max_id <= since_id {
			fmt.Fprintf(os.Stderr, "sincd id ignored when max<=since\n")
		}
		if waitsecond <= 0 {
			waitsecond = 5
			fmt.Fprintf(os.Stderr, "wait default=%d (forward)\n", waitsecond)
		}
		getFowardTLs(userid, count, eachcount, max_loop, waitsecond, max_id, since_id)
	}
	print_id()
	os.Exit(exitcode)
}

func getFowardTLs(userid string, count int, eachcount int, loops int, waitsecond int64, max twidt, since twidt) {
	totalc := 0
	var countlim bool = true
	if count <= 0 {
		countlim = false
	}
	if eachcount == 0 {
		if count > 0 {
			eachcount = count
			if eachcount > onetimemax {
				eachcount = onetimemax
			}
			fmt.Fprintf(os.Stderr, "-each=%d assumed\n", eachcount)
		}
	}
	if max > 0 {
		if max <= since {
			since = 0
		}
	}
	until := since
	if until > 0 {
		until -= 1
	}
	for i := 1; ; i++ {

		tweets, c, last, err := twapi.getTL(userid, eachcount, max, until)
		if err != nil {
			print_id()
			os.Exit(2)
		}
		// jsonTweets, _ := json.Marshal(tweets) //test
		// fmt.Println(string(jsonTweets))       //test
		totalc += c
		if c > 0 {
			firstid, lastid, nout := printTL(tweets, count, forward)
			// fmt.Println("printTL id:", firstid, "-", lastid)
			if next_since == twid_def {
				next_since = firstid
			}
			next_max = lastid

			if lastid <= since {
				break  // break by since_id
			}
			if countlim {
				count -= nout
				if count <= 0 { break }
			}
		}
		if loops > 0 && i >= loops {
			break
		}
		if last {
			if totalc == 0 {
				fmt.Fprintln(os.Stderr, "last record. break")
				exitcode = 1
			}
			break
		}
		sleep(waitsecond) //?
	}
	return
}

func getReverseTLs(userid string, count int, loops int, waitsecond int64, since twidt) {
	var tweets []*g8rv2.TweetDictionary
	var countlim bool = true
	if count <=  0 {
		countlim = false
	}
	var sinceid = since
	var delsince twidt = 0
	next_since = sinceid //default: same sinceid
	if sinceid <= 0 {
		fmt.Fprintf(os.Stderr, "since=%d. get %d tweet\n", sinceid, onetimemin)
		tweets, c, _, err := twapi.getTL(userid, onetimemin, 0, 0)
		if err != nil {
			print_id()
			os.Exit(2)
		}
		if c == 0 {
			fmt.Fprintln(os.Stderr, "Not 1 record available")
			sleep(waitsecond)
		} else {
			firstid, lastid, _ := printTL(tweets, 0, reverse)
			next_max = firstid
			next_since = lastid
			sinceid = lastid
			sleep(5)
		}
	} else {
		fmt.Fprintf(os.Stderr, "since=%d. start from this record.\n", sinceid)
	}
	for i:=1; ; i+=1 {
		tweets = getTLsince(userid, sinceid)
 
		c := len(tweets)
		if c > 0 {
			minid := str2twid(tweets[len(tweets) - 1].Tweet.ID)
			if minid <= sinceid {
				//指定ツイートまで取れたのでダブらないように削除する
				tweets = tweets[: len(tweets) - 1]
				c = len(tweets)
			} else {
				if delsince == 0 {
					//gap
					fmt.Fprintf(os.Stderr, "Gap exists. since_id=%d deleted?\n", sinceid)
					delsince = sinceid
				}
			}
			if c > 0 {
				firstid, lastid, nout := printTL(tweets, 0, reverse)
				if next_max == 0 {
					next_max = firstid
				}
				next_since = lastid
				sinceid = lastid
				delsince = 0
				if countlim {
					count -= nout
					if count <= 0 { break }
				}
			}
		} else {
			if delsince == 0 {
				//gap
				fmt.Fprintf(os.Stderr, "Gap exists. since_id=%d deleted?\n", sinceid)
				delsince = sinceid
			}
		}
		if loops > 0 && i >= loops {
			break
		}
		sleep(waitsecond)
	}
	return
}
 
func getTLsince(userid string, since twidt) (tweets []*g8rv2.TweetDictionary) {
	totalc := 0
	tweets = []*g8rv2.TweetDictionary{}
	var max_id twidt = 0
	until := since
	if until > 0 {
		until -= 1
	}
	twapi.rewindQuery()
	for i := 0; ; i++ {
 
		twts, c, last, err := twapi.getTL(userid, onetimemax, max_id, until)
		if err != nil{
			print_id()
			os.Exit(2)
		}
		totalc += c
		if c > 0 {

			lastid := str2twid(twts[c - 1].Tweet.ID)
			
			tweets = append(tweets, twts...)
			
			if since > 0 {
				if lastid <= since {
					break
				}
				max_id = lastid - 1
			} else {
				break
			}
		}
		if last {
			if totalc == 0 {
				//fmt.Fprintln(os.Stderr, "last record. break")
			}
			break
		}
		// 一度で取りきれなかった
		fmt.Fprintln(os.Stderr, "------continue")
 
		sleep(10) //??
	}
	return tweets
}

func str2twid(tweetID string) (twidt) {
	id64, _  := strconv.ParseInt(tweetID, 10, 64)
	return twidt(id64)
}

func printTL(tweets []*g8rv2.TweetDictionary, count int, revs revtype) (firstid twidt, lastid twidt, nout int) {

	firstid = twid_def
	lastid = twid_def
	imax := len(tweets)
	is := 0
	ip := 1
	if revs {
		is = imax - 1
		ip = -1
	}
	nout = 0
	for i := is; 0 <= i && i < imax; i += ip {
		tweet := tweets[i]
		id := str2twid(tweet.Tweet.ID)
		if i == is {
			firstid = id
			lastid = id
		}
		twtype, rt := ifRetweeted(tweet)
		//  RT > Reply > Mention > tweet
		var done bool
		if rt != nil {
			twtype2, _ := ifRetweeted(rt)
			done = printTweet(twtype, tweet, twtype2, rt)
		} else {
			done = printTweet("or", tweet, twtype, tweet)
		}
		if done {
			nout++
		}

		lastid = id
		
		if count > 0 && nout >= count {
			break
		}
	}
	return firstid, lastid, nout
}

func ifRetweeted(t *g8rv2.TweetDictionary) (twtype string, rt *g8rv2.TweetDictionary) {
	twtype = "tw"
	if t.InReplyUser != nil {
		twtype = "Mn"
	}
	rt = nil
	for _, r := range t.ReferencedTweets {
		switch r.Reference.Type {
		case "retweeted":
			rt = r.TweetDictionary
			twtype = "RT"
		case "replied_to":
			twtype = "Re"
		case "quoted":
		}
	}
	return twtype, rt
}

func printTweet(twtype1 string, tweet1 *g8rv2.TweetDictionary, twtype2 string, tweet2 *g8rv2.TweetDictionary) (bool) {
	tweetid := tweet1.Tweet.ID
	tweetuser := tweet1.Author.UserName

	origiid := tweet2.Tweet.ID
	origiuser := tweet2.Author.UserName
	
	firstp := true
	idst := "*Id:"
	if uniqid != nil {
		if uniqid.checkID(origiid) {
			firstp = false
			idst = "_id:"
		}
	}

	if tweetid == origiid {
		fmt.Fprintln(os.Stderr, idst, tweetid)
	} else {
		fmt.Fprintln(os.Stderr, idst, tweetid, origiid)
	}

	if firstp {
		fmt.Printf("%s\t@%s\t%s\t%s\t@%s\t%s\t\"%s\"\n",
			tweetid, tweetuser, twtype1,
			origiid, origiuser, twtype2, quoteText(tweet2.Tweet.Text))
	}
	return firstp
}

func quoteText(fulltext string) (qtext string) {
	quoted1 := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(fulltext, "\n", `\n`), "\r", `\r`), "\"", `\"`)
	qtext = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(quoted1, `&amp;`, `&`), `&lt;`, `<`), `&gt;`, `>`)
	return qtext
}


type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

func connectTwitterApi() (client *g8rv2.Client) {
	usr, _ := user.Current()
	raw, error := ioutil.ReadFile(usr.HomeDir + "/twitter/twitterBearerToken.json")
	if error != nil {
		fmt.Fprintln(os.Stderr, error.Error())
		os.Exit(2)
	}
	var twitterBearerToken TwitterBearerToken
	json.Unmarshal(raw, &twitterBearerToken)

	// raw, error = ioutil.ReadFile(usr.HomeDir + "/twitter/twitterAccount.json")
	// if error != nil {
	client =  &g8rv2.Client {
		Authorizer: authorize{
			Token: twitterBearerToken.BearerToken,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}
	// } else {
	//  	var twitterAccount TwitterAccount
	//  	json.Unmarshal(raw, &twitterAccount)
	//  
	//  	client =  g8rv2.New(twitterBearerToken.BearerToken,
	//  		g8rv2.WithConsumerKey(twitterAccount.ConsumerKey),
	//  		g8rv2.WithConsumerSecret(twitterAccount.ConsumerSecret))
	// }
	return client
}

type TwitterAccount struct {
	AccessToken string `json:"accessToken"`
	AccessTokenSecret string `json:"accessTokenSecret"`
	ConsumerKey string `json:"consumerKey"`
	ConsumerSecret string `json:"consumerSecret"`
}

type TwitterBearerToken struct {
	BearerToken string `json:"bearerToken"`
}

package main

import (
	"encoding/json"
	"context"
	"fmt"
	"strconv"
	"time"
	"os"
	//"errors"

	g8rv2 "github.com/g8rswimmer/go-twitter/v2"
)


func twid2str(twid twidt) (string) {
	return strconv.FormatInt(int64(twid), 10)
}


type twSearchApi struct {
	client *g8rv2.Client
	t  tltype
	tlopt g8rv2.UserTweetTimelineOpts
	tmopt g8rv2.UserMentionTimelineOpts
	lsopt g8rv2.ListTweetLookupOpts
	sropt g8rv2.TweetRecentSearchOpts
	saopt g8rv2.TweetSearchOpts
	nextToken string
	seq int
	jsonp bool
}

func (ta *twSearchApi) getTL(userID string, maxresult int, max twidt, since twidt) (tweets []*g8rv2.TweetDictionary, count int, last bool, err error) {
	count = 0
	last = true
	err = nil
	client := ta.client
	switch ta.t {
	case tluser:
		if ta.seq == 0 {
			ta.tlopt = g8rv2.UserTweetTimelineOpts {
				TweetFields: []g8rv2.TweetField{g8rv2.TweetFieldReferencedTweets, g8rv2.TweetFieldInReplyToUserID},
				Expansions:  []g8rv2.Expansion{g8rv2.ExpansionReferencedTweetsIDAuthorID},
				}
			if max > 0 {ta.tlopt.UntilID = twid2str(max)}
			if since > 0 {ta.tlopt.SinceID = twid2str(since)}
		} else {
			if ta.nextToken == "" {
				break
			}
			ta.tlopt.PaginationToken = ta.nextToken
		}
		if maxresult > 0 {ta.tlopt.MaxResults = maxresult}
		
		ta.seq++
		ta.nextToken = ""

		var res *g8rv2.UserTweetTimelineResponse
		res, err = client.UserTweetTimeline(context.Background(), userID, ta.tlopt)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		if res == nil {
			fmt.Fprintln(os.Stderr, "res == nil")
			break
		}
		if ta.jsonp {
			jsonraw, _ := json.MarshalIndent(res, "", "    ")
			fmt.Println(string(jsonraw))
		}
		if res.Raw != nil {
			for _, e := range res.Raw.Errors {
				fmt.Fprintf(os.Stderr, "%s: %s\n", e.Title, e.Detail)
				break
			}
			tweets = TweetDictionarySlice(res.Raw)
		}
		if res.Meta != nil {
			count = res.Meta.ResultCount
			if res.Meta.NextToken != "" {
				ta.nextToken = res.Meta.NextToken
				last = false
			}
		}
		fmt.Fprintf(os.Stderr, "%s get len: %d\n", time.Now().Format("15:04:05"), count)
		return tweets, count, last, err
	case tlhome:
	case tlmention:
		if ta.seq == 0 {
			ta.tmopt = g8rv2.UserMentionTimelineOpts {
				TweetFields: []g8rv2.TweetField{g8rv2.TweetFieldReferencedTweets, g8rv2.TweetFieldInReplyToUserID},
				Expansions:  []g8rv2.Expansion{g8rv2.ExpansionReferencedTweetsIDAuthorID},
				}
			if max > 0 {ta.tmopt.UntilID = twid2str(max)}
			if since > 0 {ta.tmopt.SinceID = twid2str(since)}
		} else {
			if ta.nextToken == "" {
				break
			}
			ta.tmopt.PaginationToken = ta.nextToken
		}
		if maxresult > 0 {ta.tmopt.MaxResults = maxresult}
		
		ta.seq++
		ta.nextToken = ""

		var res *g8rv2.UserMentionTimelineResponse
		res, err = client.UserMentionTimeline(context.Background(), userID, ta.tmopt)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		if res == nil {
			fmt.Fprintln(os.Stderr, "res == nil")
			break
		}
		if ta.jsonp {
			jsonraw, _ := json.MarshalIndent(res, "", "    ")
			fmt.Println(string(jsonraw))
		}
		if res.Raw != nil {
			for _, e := range res.Raw.Errors {
				fmt.Fprintf(os.Stderr, "%s: %s\n", e.Title, e.Detail)
				break
			}
			tweets = TweetDictionarySlice(res.Raw)
		}
		if res.Meta != nil {
			count = res.Meta.ResultCount
			if res.Meta.NextToken != "" {
				ta.nextToken = res.Meta.NextToken
				last = false
			}
		}
		fmt.Fprintf(os.Stderr, "%s get len: %d\n", time.Now().Format("15:04:05"), count)
		return tweets, count, last, err
	case tlrtofme:
	case tllist:
		if ta.seq == 0 {
			ta.lsopt = g8rv2.ListTweetLookupOpts {
				TweetFields: []g8rv2.TweetField{g8rv2.TweetFieldReferencedTweets, g8rv2.TweetFieldInReplyToUserID},
				Expansions:  []g8rv2.Expansion{g8rv2.ExpansionReferencedTweetsIDAuthorID},
				}
			//if max > 0 {ta.lsopt.UntilID = twid2str(max)}
			//if since > 0 {ta.lsopt.SinceID = twid2str(since)}
		} else {
			if ta.nextToken == "" {
				break
			}
			ta.lsopt.PaginationToken = ta.nextToken
		}
		if maxresult > 0 {ta.lsopt.MaxResults = maxresult}
		
		ta.seq++
		ta.nextToken = ""

		var res *g8rv2.ListTweetLookupResponse
		res, err = client.ListTweetLookup(context.Background(), userID, ta.lsopt)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		if res == nil {
			fmt.Fprintln(os.Stderr, "res == nil")
			break
		}
		if ta.jsonp {
			jsonraw, _ := json.MarshalIndent(res, "", "    ")
			fmt.Println(string(jsonraw))
		}
		if res.Raw != nil {
			for _, e := range res.Raw.Errors {
				fmt.Fprintf(os.Stderr, "%s: %s\n", e.Title, e.Detail)
				break
			}
			tweets = TweetDictionarySlice(res.Raw)
		}
		if res.Meta != nil {
			count = res.Meta.ResultCount
			if res.Meta.NextToken != "" {
				ta.nextToken = res.Meta.NextToken
				last = false
			}
		}
		fmt.Fprintf(os.Stderr, "%s get len: %d\n", time.Now().Format("15:04:05"), count)
		return tweets, count, last, err
	case tlsearch:
		if ta.seq == 0 {
			ta.sropt = g8rv2.TweetRecentSearchOpts {
				TweetFields: []g8rv2.TweetField{g8rv2.TweetFieldReferencedTweets, g8rv2.TweetFieldInReplyToUserID},
				Expansions:  []g8rv2.Expansion{g8rv2.ExpansionReferencedTweetsIDAuthorID},
				}
			if max > 0 {ta.sropt.UntilID = twid2str(max)}
			if since > 0 {ta.sropt.SinceID = twid2str(since)}
		} else {
			if ta.nextToken == "" {
				break
			}
			ta.sropt.NextToken = ta.nextToken
		}
		if maxresult > 0 {ta.sropt.MaxResults = maxresult}
		
		ta.seq++
		ta.nextToken = ""

		var res *g8rv2.TweetRecentSearchResponse
		res, err = client.TweetRecentSearch(context.Background(), userID, ta.sropt)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		if res == nil {
			fmt.Fprintln(os.Stderr, "res == nil")
			break
		}
		if ta.jsonp {
			jsonraw, _ := json.MarshalIndent(res, "", "    ")
			fmt.Println(string(jsonraw))
		}
		if res.Raw != nil {
			for _, e := range res.Raw.Errors {
				fmt.Fprintf(os.Stderr, "%s: %s\n", e.Title, e.Detail)
				break
			}
			tweets = TweetDictionarySlice(res.Raw)
		}
		if res.Meta != nil {
			count = res.Meta.ResultCount
			if res.Meta.NextToken != "" {
				ta.nextToken = res.Meta.NextToken
				last = false
			}
		}
		fmt.Fprintf(os.Stderr, "%s get len: %d\n", time.Now().Format("15:04:05"), count)
		return tweets, count, last, err
	case tlsearcha:
		if ta.seq == 0 {
			ta.saopt = g8rv2.TweetSearchOpts {
				TweetFields: []g8rv2.TweetField{g8rv2.TweetFieldReferencedTweets, g8rv2.TweetFieldInReplyToUserID},
				Expansions:  []g8rv2.Expansion{g8rv2.ExpansionReferencedTweetsIDAuthorID},
				}
			if max > 0 {ta.saopt.UntilID = twid2str(max)}
			if since > 0 {ta.saopt.SinceID = twid2str(since)}
		} else {
			if ta.nextToken == "" {
				break
			}
			ta.saopt.NextToken = ta.nextToken
		}
		if maxresult > 0 {ta.saopt.MaxResults = maxresult}
		
		ta.seq++
		ta.nextToken = ""

		var res *g8rv2.TweetSearchResponse
		res, err = client.TweetSearch(context.Background(), userID, ta.saopt)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
		if res == nil {
			fmt.Fprintln(os.Stderr, "res == nil")
			break
		}
		if ta.jsonp {
			jsonraw, _ := json.MarshalIndent(res, "", "    ")
			fmt.Println(string(jsonraw))
		}
		if res.Raw != nil {
			for _, e := range res.Raw.Errors {
				fmt.Fprintf(os.Stderr, "%s: %s\n", e.Title, e.Detail)
				break
			}
			tweets = TweetDictionarySlice(res.Raw)
		}
		if res.Meta != nil {
			count = res.Meta.ResultCount
			if res.Meta.NextToken != "" {
				ta.nextToken = res.Meta.NextToken
				last = false
			}
		}
		fmt.Fprintf(os.Stderr, "%s get len: %d\n", time.Now().Format("15:04:05"), count)
		return tweets, count, last, err
	}
	return nil, 0, true, err
}

func (ta *twSearchApi) rewindQuery() {
	ta.seq = 0
}

func TweetDictionarySlice(t *g8rv2.TweetRaw) (s []*g8rv2.TweetDictionary) {
	s = []*g8rv2.TweetDictionary{}
	for _, tweet := range t.Tweets {
		s = append(s, g8rv2.CreateTweetDictionary(*tweet, t.Includes))
	}
	return s
}

// // UserTweetTimelineOpts are the options for the user tweet timeline request
// type UserTweetTimelineOpts struct {
//  	Expansions      []Expansion
//  	MediaFields     []MediaField
//  	PlaceFields     []PlaceField
//  	PollFields      []PollField
//  	TweetFields     []TweetField
//  	UserFields      []UserField
//  	Excludes        []Exclude
//  	StartTime       time.Time
//  	EndTime         time.Time
//  	MaxResults      int
//  	PaginationToken string
//  	SinceID         string
//  	UntilID         string
// }

// // UserTweetTimeline will return the user tweet timeline
// func (c *Client) UserTweetTimeline(ctx context.Context, userID string, opts UserTweetTimelineOpts) (*UserTweetTimelineResponse, error) {

// // UserTweetTimelineResponse contains the information from the user tweet timeline callout
// type UserTweetTimelineResponse struct {
//  	Raw       *TweetRaw
//  	Meta      *UserTimelineMeta `json:"meta"`
//  	RateLimit *RateLimit
// }
// TweetRaw is the raw response from the tweet lookup endpoint

// type TweetRaw struct {
//  	Tweets       []*TweetObj       `json:"data"`
//  	Includes     *TweetRawIncludes `json:"includes,omitempty"`
//  	Errors       []*ErrorObj       `json:"errors,omitempty"`
//  	dictionaries map[string]*TweetDictionary
// }

// // UserTimelineMeta contains the meta data from the timeline callout
// type UserTimelineMeta struct {
//  	ResultCount   int    `json:"result_count"`
//  	NewestID      string `json:"newest_id"`
//  	OldestID      string `json:"oldest_id"`
//  	NextToken     string `json:"next_token"`
//  	PreviousToken string `json:"previous_token"`
// }

// // TweetDictionary is a struct of a tweet and all of the reference objects
// type TweetDictionary struct {
//  	Tweet            TweetObj
//  	Author           *UserObj
//  	InReplyUser      *UserObj
//  	Place            *PlaceObj
//  	AttachmentPolls  []*PollObj
//  	AttachmentMedia  []*MediaObj
//  	Mentions         []*TweetMention
//  	ReferencedTweets []*TweetReference
// }

// // TweetDictionaries create a map of tweet dictionaries from the raw tweet response
// func (t *TweetRaw) TweetDictionaries() map[string]*TweetDictionary {
//  	if t.dictionaries != nil {
//  		return t.dictionaries
//  	}
//  
//  	t.dictionaries = map[string]*TweetDictionary{}
//  	for _, tweet := range t.Tweets {
//  		t.dictionaries[tweet.ID] = CreateTweetDictionary(*tweet, t.Includes)
//  	}
//  	return t.dictionaries
// }

// // TweetDictionary is a struct of a tweet and all of the reference objects
// type TweetDictionary struct {
//  	Tweet            TweetObj
//  	Author           *UserObj
//  	InReplyUser      *UserObj
//  	Place            *PlaceObj
//  	AttachmentPolls  []*PollObj
//  	AttachmentMedia  []*MediaObj
//  	Mentions         []*TweetMention
//  	ReferencedTweets []*TweetReference
// }

// // TweetObj is the primary object on the tweets endpoints
// type TweetObj struct {
//  	ID                 string                       `json:"id"`
//  	Text               string                       `json:"text"`
//  	Attachments        *TweetAttachmentsObj         `json:"attachments,omitempty"`
//  	AuthorID           string                       `json:"author_id,omitempty"`
//  	ContextAnnotations []*TweetContextAnnotationObj `json:"context_annotations,omitempty"`
//  	ConversationID     string                       `json:"conversation_id,omitempty"`
//  	CreatedAt          string                       `json:"created_at,omitempty"`
//  	Entities           *EntitiesObj                 `json:"entities,omitempty"`
//  	Geo                *TweetGeoObj                 `json:"geo,omitempty"`
//  	InReplyToUserID    string                       `json:"in_reply_to_user_id,omitempty"`
//  	Language           string                       `json:"lang,omitempty"`
//  	NonPublicMetrics   *TweetMetricsObj             `json:"non_public_metrics,omitempty"`
//  	OrganicMetrics     *TweetMetricsObj             `json:"organic_metrics,omitempty"`
//  	PossiblySensitive  bool                         `json:"possibly_sensitive,omitempty"`
//  	PromotedMetrics    *TweetMetricsObj             `json:"promoted_metrics,omitempty"`
//  	PublicMetrics      *TweetMetricsObj             `json:"public_metrics,omitempty"`
//  	ReferencedTweets   []*TweetReferencedTweetObj   `json:"referenced_tweets,omitempty"`
//  	Source             string                       `json:"source,omitempty"`
//  	WithHeld           *WithHeldObj                 `json:"withheld,omitempty"`
// }

// // UserObj contains Twitter user account metadata describing the referenced user
// type UserObj struct {
//  	ID              string          `json:"id"`
//  	Name            string          `json:"name"`
//  	UserName        string          `json:"username"`
//  	CreatedAt       string          `json:"created_at,omitempty"`
//  	Description     string          `json:"description,omitempty"`
//  	Entities        *EntitiesObj    `json:"entities,omitempty"`
//  	Location        string          `json:"location,omitempty"`
//  	PinnedTweetID   string          `json:"pinned_tweet_id,omitempty"`
//  	ProfileImageURL string          `json:"profile_image_url,omitempty"`
//  	Protected       bool            `json:"protected,omitempty"`
//  	PublicMetrics   *UserMetricsObj `json:"public_metrics,omitempty"`
//  	URL             string          `json:"url,omitempty"`
//  	Verified        bool            `json:"verified,omitempty"`
//  	WithHeld        *WithHeldObj    `json:"withheld,omitempty"`
// }

// // TweetReference is the tweet referenced and it's dictionary
// type TweetReference struct {
//  	Reference       *TweetReferencedTweetObj
//  	TweetDictionary *TweetDictionary
// }

// // TweetReferencedTweetObj is a Tweet this Tweet refers to
// type TweetReferencedTweetObj struct {
//  	Type string `json:"type"`
//  	ID   string `json:"id"`
// }

// // ErrorObj is part of the partial errors in the response
// type ErrorObj struct {
//  	Title        string      `json:"title"`
//  	Detail       string      `json:"detail"`
//  	Type         string      `json:"type"`
//  	ResourceType string      `json:"resource_type"`
//  	Parameter    string      `json:"parameter"`
//  	Value        interface{} `json:"value"`
// }

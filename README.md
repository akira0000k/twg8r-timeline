# twg8r-timeline
Twitter timeline get command. API v2 with g8rswimmer/go-twitter library.

~~~
$ ./g8getl -help
Usage of ./g8getl:
-get string
	TLtype: user, mention, list and search

-user string
	twitter @ screenname
-userid string
	integer user Id (default "0")

-listid string
	list ID (default "0")
-listname string
	list name

-query string
	Query String
-restype string
	result type: [recent]/all

-reverse
	reverse output. wait newest TL

-max_id int
	starting tweet id
-since_id int
	reverse start tweet id

-count int
	tweet count. 5-3200 ?
-each int
	req count for each loop 5-100
-loops int
	get loop max
-wait int
	wait second for next loop
~~~


## parameter example
### ユーザーTL
    [-get=user]  -user=screenname / -userid=9999999
    -get=mention -user=screenname / -userid=9999999

### リストTL
    [-get=list] -listid=99999999
    [-get=list] -listname=リスト名 / -listid=99999999  -user=screenname / -userid=9999999
    -max_id -since_id -reverse は使えない

### 検索
	[-get=search] -query=検索文字列

### 取得方向
    -reverse  (逆。最新待ち受け取得)  順方向は過去へ

### 続き指示
順方向ではこの次から古いものをとる

    -max_id=1529278564566454273

逆方向ではこの次から待ち受ける

    -since_id=1529278731545882624 -reverse

### その他パラメタ
    -count=取得件数めやす　　(デフォルトは順10件, 逆のデフォルトは制限なし。全体件数の制御)
    -each=一回の取得件数　 　(順のみ、デフォルト20件, 最大100件)
    -loops=内部繰り返し数　　(-countで全体件数を制御するか、-eachと-loops で制御してもよい)
    -wait=秒             　(ループ間隔)  デフォルト 順10 逆60

## 認証
~/twitter/twitterBearerToken.json に所有者トークンを書き込み保存しておく。

~~~	
{
    "bearerToken": "????????????????????????????????????????????????????????????????????????????????????????????????????????????????"
}
~~~	

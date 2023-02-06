// https://qiita.com/ppco/items/8bf22a7bde9be13c22f1
package trainingTwitter

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/ozaki-physics/go-training-composition/trainingIo"
)

type creds struct {
	ConsumerKey       string `json:"consumer_key"`
	ConsumerSecret    string `json:"consumer_secret"`
	AccessToken       string `json:"access_token"`
	AccessTokenSecret string `json:"access_token_secret"`
}

func JsonToStruct(goStruct interface{}, path string) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	// JSON から struct にする
	if err := json.Unmarshal(bytes, &goStruct); err != nil {
		log.Fatalln(err)
	}
	// fmt.Println(goStruct)
}

// terminalArgs 実行時にコマンド引数から JSON ファイルのパスを受け取る
func terminalArgs() string {
	flag.Parse()
	args := flag.Args()

	var filepath string
	if len(args) > 0 {
		// コマンド引数の1個目を取得
		filepath = args[0]
		// ファイルが存在するか確認
		trainingIo.SearchFile(filepath)
	} else {
		// 標準入力で JSON の filapath を受け取る
		fmt.Println("JSON ファイルのパスが 実行時に渡されませんでした")
		log.Fatal("終了")
	}

	fmt.Printf("%s でアクセスします\n", filepath)
	return filepath
}

func newCreds() creds {
	var c creds
	// path := "./trainingTwitter/secret.json"
	path := terminalArgs()
	JsonToStruct(&c, path)
	return c
}

// 2.1 oauth_nonce
func createoauthNonce() string {
	key := make([]byte, 32)
	rand.Read(key)
	enc := base64.StdEncoding.EncodeToString(key)
	replaceStr := []string{"+", "/", "="}
	for _, str := range replaceStr {
		enc = strings.Replace(enc, str, "", -1)
	}
	return enc
}

// 2.2.1 oauth_signature以外のパラメータと機能ごとの追加パラメータを文字列順にソートし、パーセントエンコードを行い、&で結合
type sortedQuery struct {
	m    map[string]string
	keys []string
}

func mapMerge(m1, m2 map[string]string) map[string]string {
	m := map[string]string{}

	for k, v := range m1 {
		m[k] = v
	}
	for k, v := range m2 {
		m[k] = v
	}
	return m
}

func sortedQueryString(m map[string]string) string {
	sq := &sortedQuery{
		m:    m,
		keys: make([]string, len(m)),
	}
	var i int
	for key := range m {
		sq.keys[i] = key
		i++
	}
	sort.Strings(sq.keys)

	values := make([]string, len(sq.keys))
	for i, key := range sq.keys {
		values[i] = fmt.Sprintf("%s=%s", url.QueryEscape(key), url.QueryEscape(sq.m[key]))
	}
	return strings.Join(values, "&")
}

func calcHMACSHA1(base, key string) string {
	b := []byte(key)
	h := hmac.New(sha1.New, b)
	io.WriteString(h, base)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// まとめると
func manualOauthSettings(creds *creds, additionalParam map[string]string, httpMethod, uri string) string {
	m := map[string]string{}
	m["oauth_consumer_key"] = creds.ConsumerKey
	m["oauth_nonce"] = createoauthNonce()
	m["oauth_signature_method"] = "HMAC-SHA1"
	m["oauth_timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	m["oauth_token"] = creds.AccessToken
	m["oauth_version"] = "1.0"

	// 2.2.2 HTTPメソッドとAPIのエンドポイントをパーセントエンコードし、1で作成した文字列とさらに&で結合
	baseQueryString := sortedQueryString(mapMerge(m, additionalParam))

	base := []string{}
	base = append(base, url.QueryEscape(httpMethod))
	base = append(base, url.QueryEscape(uri))
	base = append(base, url.QueryEscape(baseQueryString))

	signatureBase := strings.Join(base, "&")

	// 2.2.3 発行されたConsumerSecretとAccessSecretをパーセントエンコードし、&で結合
	signatureKey := url.QueryEscape(creds.ConsumerSecret) + "&" + url.QueryEscape(creds.AccessTokenSecret)

	// 2.2.4 2および3で生成した文字列に対して、3をキーにして、HMACSHA1で計算したものに対し、BASE64エンコード
	m["oauth_signature"] = calcHMACSHA1(signatureBase, signatureKey)

	authHeader := fmt.Sprintf("OAuth oauth_consumer_key=\"%s\", oauth_nonce=\"%s\", oauth_signature=\"%s\", oauth_signature_method=\"%s\", oauth_timestamp=\"%s\", oauth_token=\"%s\", oauth_version=\"%s\"",
		url.QueryEscape(m["oauth_consumer_key"]),
		url.QueryEscape(m["oauth_nonce"]),
		url.QueryEscape(m["oauth_signature"]),
		url.QueryEscape(m["oauth_signature_method"]),
		url.QueryEscape(m["oauth_timestamp"]),
		url.QueryEscape(m["oauth_token"]),
		url.QueryEscape(m["oauth_version"]),
	)

	return authHeader
}

// 3.ツイート機能
func tweet(creds *creds, message string) (*http.Response, error) {
	addtionalParam := map[string]string{"status": message}
	url := "https://api.twitter.com/1.1/statuses/update.json"
	authHeader := manualOauthSettings(creds, addtionalParam, "POST", url)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", authHeader)
	req.URL.RawQuery = sortedQueryString(addtionalParam)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

// リツイートを取得
func GetRetweet() {
	c := newCreds()
	// message := "hello world !"
	// addtionalParam := map[string]string{"status": message}
	addtionalParam := map[string]string{}
	// url := "https://api.twitter.com/1.1/statuses/update.json"
	// url := "https://api.twitter.com/1.1/statuses/user_timeline.json"
	id := 1607485276577824773
	url := fmt.Sprintf("https://api.twitter.com/1.1/statuses/retweets/%d.json?count=4&trim_user=false", id)
	authHeader := manualOauthSettings(&c, addtionalParam, "POST", url)
	fmt.Println(authHeader)
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Authorization", authHeader)
	req.URL.RawQuery = sortedQueryString(addtionalParam)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	fmt.Println(string(b))
}

// go get github.com/ChimeraCoder/anaconda
func GetRTbyAnaconda() {
	c := newCreds()
	client := anaconda.NewTwitterApiWithCredentials(
		c.AccessToken,
		c.AccessTokenSecret,
		c.ConsumerKey,
		c.ConsumerSecret,
	)
	// var id int64 = 1604582274917351426
	// 818 件のリツイートがあるはずなのに 85 件しか取得できない...
	// var id int64 = 1504382036961865735
	// キャンペーン
	var id int64 = 1607979788102209536
	v := url.Values{}
	v.Set("count", "100")
	// v.Set("trim_user", "false")
	// v.Set("max_results", "100")
	tweets, err := client.GetRetweets(id, v)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(tweets))
	for _, tweet := range tweets {
		// fmt.Printf("%+v", tweet.User)
		// fmt.Printf("%+v\n", tweet.User.Name)
		// fmt.Printf("@%s\n", tweet.User.ScreenName)
		// fmt.Printf("%+v\n", tweet.User)
		fmt.Printf(
			"\"%s\",\"@%s\",\"%d\",%d,%d,\"%s\",\"https://twitter.com/%s\"\n",
			tweet.User.Name,           // 表示名
			tweet.User.ScreenName,     // ユーザー名
			tweet.User.Id,             // ユーザーID
			tweet.User.FriendsCount,   // フォロー
			tweet.User.FollowersCount, // フォロワー
			tweet.User.Description,    // プロフィール
			tweet.User.ScreenName,     // ユーザー名
		)
	}
}

func GetTimeLine() {
	c := newCreds()
	client := anaconda.NewTwitterApiWithCredentials(
		c.AccessToken,
		c.AccessTokenSecret,
		c.ConsumerKey,
		c.ConsumerSecret,
	)
	v := url.Values{}
	v.Set("count", "3")
	tweets, err := client.GetHomeTimeline(v)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(tweets))
	for _, tweet := range tweets {
		fmt.Println(tweet.User.Name)
		fmt.Println(tweet.Text)
	}
}

// https://zenn.dev/sivchari/articles/2b55ebfdb34621
const retweetsLookupURL = "https://api.twitter.com/2/tweets/%v/retweeted_by"

type oauth struct {
	tokenType   string
	accessToken string
}

func (o *oauth) unmarshalJSON(b io.ReadCloser) error {
	var d struct {
		TokenType   string `json:"token_type"`
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(b).Decode(&d); err != nil {
		return err
	}
	o.tokenType = d.TokenType
	o.accessToken = d.AccessToken
	return nil
}

// v2 で試す
func getBearerToken(ctx context.Context) string {
	c := newCreds()
	credentials := c.ConsumerKey + ":" + c.ConsumerSecret

	generateAppOnlyBearerTokenURL := "https://api.twitter.com/oauth2/token?grant_type=client_credentials"
	b64credentials := base64.StdEncoding.EncodeToString([]byte(credentials))
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, generateAppOnlyBearerTokenURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Basic "+b64credentials)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// b, _ := io.ReadAll(resp.Body)
	// log.Println(string(b))

	var o oauth
	err = o.unmarshalJSON(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return o.accessToken
}

type RetweetsResponse struct {
	Users []*User `json:"data"`
	// Includes *TweetIncludes      `json:"includes,omitempty"`
	// Errors   []*APIResponseError `json:"errors,omitempty"`
	// Meta     *RetweetsLookupMeta `json:"meta"`
	Title  string `json:"title,omitempty"`
	Detail string `json:"detail,omitempty"`
	Type   string `json:"type,omitempty"`
}
type User struct {
	ID              string      `json:"id"`
	Name            string      `json:"name"`
	UserName        string      `json:"username"`
	CreatedAt       string      `json:"created_at,omitempty"`
	Description     string      `json:"description,omitempty"`
	Entities        *UserEntity `json:"entities,omitempty"`
	Location        string      `json:"location,omitempty"`
	PinnedTweetID   string      `json:"pinned_tweet_id,omitempty"`
	ProfileImageURL string      `json:"profile_image_url,omitempty"`
	Protected       bool        `json:"protected,omitempty"`
	// PublicMetrics   *UserPublicMetrics `json:"public_metrics,omitempty"`
	URL      string `json:"url,omitempty"`
	Verified bool   `json:"verified,omitempty"`
	// Withheld        *UserWithheld      `json:"withheld,omitempty"`
}
type UserEntity struct {
	URL *UserURL `json:"url"`
	// Description *UserDescription `json:"description"`
}
type UserURL struct {
	URLs []*UserURLs `json:"urls"`
}
type UserURLs struct {
	Start       int    `json:"start"`
	End         int    `json:"end"`
	URL         string `json:"url"`
	ExpandedURL string `json:"expanded_url"`
	DisplayURL  string `json:"display_url"`
}

// v2 で試す
func getRetweetV2Access(ctx context.Context) *RetweetsResponse {
	tweetID := 1607485276577824773
	ep := fmt.Sprintf(retweetsLookupURL, tweetID)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	bearerToken := getBearerToken(ctx)
	fmt.Println(bearerToken)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	log.Println(string(b))

	var retweetsLookup RetweetsResponse
	if err := json.NewDecoder(resp.Body).Decode(&retweetsLookup); err != nil {
		log.Fatal(err)
	}

	return &retweetsLookup
}

// v2 で試す
func GetRetweetV2() {
	ctx := context.Background()
	retweets := getRetweetV2Access(ctx)
	fmt.Println(len(retweets.Users))
	for _, v := range retweets.Users {
		fmt.Println(v.ID)
		fmt.Println(v.Name)
	}
}

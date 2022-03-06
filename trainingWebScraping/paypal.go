package trainingWebScraping

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// PAYPAL_SANDBOX_URL テスト環境の基本となる URL
const PAYPAL_SANDBOX_URL = "https://api-m.sandbox.paypal.com"

// Key paypal の key が格納された JSON のための struct
type Key struct {
	Paypal struct {
		ClientID string `json:"clientID"`
		Secret   string `json:"secret"`
	} `json:"paypal"`
}

// GetPaypalAccessToken アクセストークンを取得する
func GetPaypalAccessToken(keypath string) string {
	urlString := PAYPAL_SANDBOX_URL + "/v1/oauth2/token"

	// JSON から値を取り出して Go に格納する
	bytes, err := os.ReadFile(keypath)
	if err != nil {
		panic(err)
	}
	var keyJSON Key
	if err := json.Unmarshal(bytes, &keyJSON); err != nil {
		panic(err)
	}

	// POST で送信する中身のデータ
	values := url.Values{}
	values.Add("grant_type", "client_credentials")

	// 10秒でタイムアウトするクライアントを作成
	client := &http.Client{Timeout: time.Duration(10) * time.Second}
	// Request を生成
	req, err := http.NewRequest("POST", urlString, strings.NewReader(values.Encode()))
	if err != nil {
		panic(err)
	}
	// リクエストヘッダ の追加
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Language", "en_US")
	// BASIC 認証(curl では -u に該当する)
	req.SetBasicAuth(keyJSON.Paypal.ClientID, keyJSON.Paypal.Secret)

	// リクエストする
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// レスポンスから body を取り出す
	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("ステータスコード: " + strconv.Itoa(resp.StatusCode))
	fmt.Println("レスポンスボディ: " + string(bodyByte))

	// レスポンスの body(JSON 形式)の中でも使いたい値だけ取り出すため
	type ResponseJSON struct {
		AccessToken string `json:"access_token"`
	}
	var responseJSON ResponseJSON
	if err := json.Unmarshal(bodyByte, &responseJSON); err != nil {
		panic(err)
	}

	fmt.Println("アクセストークン: " + responseJSON.AccessToken)
	return responseJSON.AccessToken
}

// GetPaypalClientToken クライアントトークンを取得する
func GetPaypalClientToken(accessToken string) {
	urlString := PAYPAL_SANDBOX_URL + "/v1/identity/generate-token"
	// 10秒でタイムアウトするクライアントを作成
	client := &http.Client{Timeout: time.Duration(10) * time.Second}
	// Request を生成
	req, err := http.NewRequest("POST", urlString, nil)
	if err != nil {
		panic(err)
	}
	// リクエストヘッダ の追加
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Accept-Language", "en_US")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	// リクエストする
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// レスポンスから body を取り出す
	bodyByte, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("ステータスコード: " + strconv.Itoa(resp.StatusCode))
	fmt.Println("レスポンスボディ: " + string(bodyByte))

	// レスポンスの body(JSON 形式)の中でも使いたい値だけ取り出すため
	type ResponseJSON struct {
		ClientToken string `json:"client_token"`
	}
	var responseJSON ResponseJSON
	if err := json.Unmarshal(bodyByte, &responseJSON); err != nil {
		panic(err)
	}

	fmt.Println("クライアントトークン: " + responseJSON.ClientToken)
}

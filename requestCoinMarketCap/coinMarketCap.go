package requestCoinMarketCap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

// CoinMarketCap_sample ドキュメントに書いてあったサンプル
func CoinMarketCap_sample() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://sandbox-api.coinmarketcap.com/v1/cryptocurrency/listings/latest", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := url.Values{}
	q.Add("start", "1")
	q.Add("limit", "5000")
	q.Add("convert", "USD")

	req.Header.Set("Accepts", "application/json")
	req.Header.Add("X-CMC_PRO_API_KEY", "b54bcf4d-1bca-4e8e-9a24-22ff2c3d462c")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to server")
		os.Exit(1)
	}
	fmt.Println(resp.Status)
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respBody))
}

const (
	SandboxBaseEndpoint = "https://sandbox-api.coinmarketcap.com"
	SandboxAPIKey       = "b54bcf4d-1bca-4e8e-9a24-22ff2c3d462c"
	LiveBaseEndpoint    = "https://pro-api.coinmarketcap.com"
)

// CoinMarketCap クレデンシャル(key と URL を保持)
type CoinMarketCap struct {
	Service struct {
		Key     string `json:"apiKey"`
		baseURL string
	} `json:"CoinMarketCap"`
}

// ResponseStatus レスポンスの共通部分
type ResponseStatus struct {
	Timestamp    string `json:"timestamp"`
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Elapsed      int    `json:"elapsed"`
	CreditCount  int    `json:"credit_count"`
	Notice       string `json:"notice"`
}

// Platform CMCIDMap と Metadata の共通部分
type Platform struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Symbol       string `json:"symbol"`
	Slug         string `json:"slug"`
	TokenAddress string `json:"token_address"`
}

// getCredential JSON から値を取り出して struct に格納する
func GetCredential(keypath string, isLive bool) CoinMarketCap {
	// Sandbox モード
	if isLive == false {
		var sandboxCredential CoinMarketCap
		sandboxCredential.Service.Key = SandboxAPIKey
		sandboxCredential.Service.baseURL = SandboxBaseEndpoint
		return sandboxCredential
	}

	// Live モード
	bytes, err := os.ReadFile(keypath)
	if err != nil {
		log.Fatal(err)
	}
	var keyJSON CoinMarketCap
	if err := json.Unmarshal(bytes, &keyJSON); err != nil {
		log.Fatal(err)
	}
	keyJSON.Service.baseURL = LiveBaseEndpoint
	return keyJSON
}

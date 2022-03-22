package requestCoinMarketCap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type KeyInfo struct {
	Plan struct {
		CreditLimitDaily                 int    `json:"credit_limit_daily"`
		CreditLimitDailyReset            string `json:"credit_limit_daily_reset"`
		CreditLimitDailyResetTimestamp   string `json:"credit_limit_daily_reset_timestamp"`
		CreditLimitMonthly               int    `json:"credit_limit_monthly"`
		CreditLimitMonthlyReset          string `json:"credit_limit_monthly_reset"`
		CreditLimitMonthlyResetTimestamp string `json:"credit_limit_monthly_reset_timestamp"`
		RateLimitMinute                  int    `json:"rate_limit_minute"`
	} `json:"plan"`
	Usage struct {
		CurrentMinute struct {
			RequestsMade int `json:"requests_made"`
			RequestsLeft int `json:"requests_left"`
		} `json:"current_minute"`
		CurrentDay struct {
			CreditsUsed int `json:"credits_used"`
			CreditsLeft int `json:"credits_left"`
		} `json:"current_day"`
		CurrentMonth struct {
			CreditsUsed int `json:"credits_used"`
			CreditsLeft int `json:"credits_left"`
		} `json:"current_month"`
	} `json:"usage"`
}

type ResponseKeyInfo struct {
	Data   KeyInfo        `json:"data"`
	Status ResponseStatus `json:"status"`
}

func (c *CoinMarketCap) GetKeyInfo() {
	const entryURL = "/v1/key/info"

	// クライアントの作成
	client := &http.Client{}
	req, err := http.NewRequest("GET", c.Service.baseURL+entryURL, nil)
	if err != nil {
		log.Println(err)
	}
	// リクエストヘッダーの作成
	req.Header.Set("Accepts", "application/json")
	// req.Header.Set("Accept-Encoding", "deflate, gzip")
	req.Header.Set("X-CMC_PRO_API_KEY", c.Service.Key)

	// リクエストする
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	// ステータスコードを確認
	log.Println(resp.Status)

	// ResponseBody を取り出す
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	// とりあえず出力
	// fmt.Println(string(respBody))
	// インデントを整えて出力
	// var buf bytes.Buffer
	// json.Indent(&buf, respBody, "", "  ")
	// fmt.Println(buf.String())

	// JSON にして出力
	var responseJSON ResponseKeyInfo
	if err := json.Unmarshal(respBody, &responseJSON); err != nil {
		log.Println(err)
	}
	fmt.Println("struct に入っているか?")
	fmt.Printf("%+v\n", responseJSON)
}

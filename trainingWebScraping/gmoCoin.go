package trainingWebScraping

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// enum みたいなモノを作ってみる
type Symbol int

const (
	_ Symbol = iota
	ALL
	BTC
	ETH
	BCH
	LTC
	XRP
	XEM
	XLM
	XYM
	MONA
	BTC_JPY
	ETH_JPY
	BCH_JPY
	LTC_JPY
	XRP_JPY
)

// enum の 値を定義する
func (s Symbol) String() string {
	switch {
	case s == BTC:
		return fmt.Sprintf("%s", "BTC")
	case s == ETH:
		return fmt.Sprintf("%s", "ETH")
	case s == BCH:
		return fmt.Sprintf("%s", "BCH")
	case s == LTC:
		return fmt.Sprintf("%s", "LTC")
	case s == XRP:
		return fmt.Sprintf("%s", "XRP")
	case s == XEM:
		return fmt.Sprintf("%s", "XEM")
	case s == XLM:
		return fmt.Sprintf("%s", "XLM")
	case s == XYM:
		return fmt.Sprintf("%s", "XYM")
	case s == MONA:
		return fmt.Sprintf("%s", "MONA")
	case s == BTC_JPY:
		return fmt.Sprintf("%s", "BTC_JPY")
	case s == ETH_JPY:
		return fmt.Sprintf("%s", "ETH_JPY")
	case s == BCH_JPY:
		return fmt.Sprintf("%s", "BCH_JPY")
	case s == LTC_JPY:
		return fmt.Sprintf("%s", "LTC_JPY")
	case s == XRP_JPY:
		return fmt.Sprintf("%s", "XRP_JPY")
	default:
		return ""
	}
}

const endPoint = "https://api.coin.z.com/public"

func GetGMOCoin() {
	var s Symbol
	s = ETH
	getTickerRate(&s)
}

// GetGMOCoin GMO コインの public API を呼んで暗号資産のレートを取得する
func getTickerRate(s *Symbol) {
	path := "/v1/ticker"
	accessURL := endPoint + path

	// URL パラメータを付与
	values := url.Values{}
	if *s != ALL {
		values.Add("symbol", s.String())
		accessURL += "?" + values.Encode()
	}

	response, err := http.Get(accessURL)
	if err != nil {
		log.Println(err)
	}
	body, err := ioutil.ReadAll(response.Body)
	if response.StatusCode > 299 {
		log.Printf("ステータスコードが正常ではない:\n  ステータスコード:\n  %d\n  body:\n%s", response.StatusCode, body)
		return
	}
	if err != nil {
		log.Println(err)
	}
	var buf bytes.Buffer
	// インデントを整える
	json.Indent(&buf, body, "", "  ")
	fmt.Println(buf.String())
}

package requestCoinMarketCap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type CMCID struct {
	ID                  int      `json:"id"`
	Name                string   `json:"name"`
	Symbol              string   `json:"symbol"`
	Slug                string   `json:"slug"`
	Rank                int      `json:"rank"`
	IsActive            int      `json:"is_active"`
	FirstHistoricalData string   `json:"first_historical_data"`
	LastHistoricalData  string   `json:"last_historical_data"`
	Platform            Platform `json:"platform"`
}

type ResponseCryptocurrencyMap struct {
	Status ResponseStatus `json:"status"`
	Data   []CMCID        `json:"data"`
}

type UseCMCID struct {
	Symbol string
	CMCID  int
	Name   string
	Slug   string
}

func ExampleSearchCMCID(c *CoinMarketCap) {
	var symbols []string
	symbols = append(symbols, "USDT")
	symbols = append(symbols, "ETH")
	tmp := c.SearchCMCID(symbols)
	fmt.Printf("%+v\n", tmp)
}

// searchCMCID CoinMarketCap での ID を調べる
// Sandbox 環境では クエリ文字列パラメータ が一部 正常に機能していないと思われる
func (c *CoinMarketCap) SearchCMCID(symbols []string) []UseCMCID {
	const entryURL = "/v1/cryptocurrency/map"

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

	// クエリ文字列パラメータの作成
	q := url.Values{}
	if symbol := makeQueryParmSymbol(&symbols); symbol != "" {
		q.Add("symbol", symbol)
	}
	q.Add("limit", "5")       // もし Symbol がブランクだったときのため
	q.Add("sort", "cmc_rank") // CMC の上位順でソートして取得
	// q.Add("aux", "status")
	req.URL.RawQuery = q.Encode()

	// リクエストする
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	// ステータスコードを確認
	fmt.Println(resp.Status)

	// ResponseBody を取り出す
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	// とりあえず出力する
	// fmt.Println(string(respBody))
	// インデントを整えて出力してみる
	// var buf bytes.Buffer
	// json.Indent(&buf, respBody, "", "  ")
	// fmt.Println(buf.String())
	// JSON にして出力してみる
	var responseJSON ResponseCryptocurrencyMap
	if err := json.Unmarshal(respBody, &responseJSON); err != nil {
		log.Println(err)
	}
	fmt.Println("struct に入っているか?")
	fmt.Printf("%+v\n", responseJSON)
	return generateSymbolAndCMCID(&responseJSON.Data)
}

// makeQueryParmSymbol Symbol のスライスをカンマ区切りの1個の string にする
// makeQueryParmCMCID とアルゴリズムは同じ
func makeQueryParmSymbol(symbols *[]string) string {
	var queryParmSymbol string
	if len(*symbols) == 0 {
		return queryParmSymbol
	}

	for i, symbol := range *symbols {
		if i != 0 {
			queryParmSymbol += ","
		}
		queryParmSymbol += symbol
	}
	return queryParmSymbol
}

// generateSymbolAndCMCID レスポンスから よく使いそうな値たちだけの構造体に変換する
func generateSymbolAndCMCID(respData *[]CMCID) []UseCMCID {
	var cmcIDMap []UseCMCID
	for _, data := range *respData {
		tmp := UseCMCID{
			data.Symbol,
			data.ID,
			data.Name,
			data.Slug,
		}
		cmcIDMap = append(cmcIDMap, tmp)
	}
	return cmcIDMap
}

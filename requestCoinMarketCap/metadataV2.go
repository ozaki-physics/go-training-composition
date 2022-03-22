package requestCoinMarketCap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type ContractAddr struct {
	ContractAddress string `json:"contract_address"`
	Platform        struct {
		Name string `json:"name"`
		Coin struct {
			ID     string `json:"id"`
			Name   string `json:"name"`
			Symbol string `json:"symbol"`
			Slug   string `json:"slug"`
		} `json:"coin"`
	} `json:"platform"`
}

type CryptocurrencyInfo struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Symbol      string   `json:"symbol"`
	Category    string   `json:"category"`
	Description string   `json:"description"`
	Slug        string   `json:"slug"`
	Logo        string   `json:"logo"`
	Subreddit   string   `json:"subreddit"`
	Notice      string   `json:"notice"`
	Tags        []string `json:"tags"`
	TagNames    []string `json:"tag-names"`
	TagGroups   []string `json:"tag-groups"`
	Urls        struct {
		Website      []string `json:"website"`
		Twitter      []string `json:"twitter"`
		MessageBoard []string `json:"message_board"`
		Chat         []string `json:"chat"`
		Facebook     []string `json:"facebook"`
		Explorer     []string `json:"explorer"`
		Reddit       []string `json:"reddit"`
		TechnicalDoc []string `json:"technical_doc"`
		SourceCode   []string `json:"source_code"`
		Announcement []string `json:"announcement"`
	} `json:"urls"`
	Platform                      Platform       `json:"platform"`
	DateAdded                     string         `json:"date_added"`
	TwitterUsername               string         `json:"twitter_username"`
	IsHidden                      int            `json:"is_hidden"`
	DateLaunched                  string         `json:"date_launched"`
	ContractAddress               []ContractAddr `json:"contract_address"`
	SelfReportedCirculatingSupply int            `json:"self_reported_circulating_supply"`
	SelfReportedTags              []string       `json:"self_reported_tags"`
	SelfReportedMarketCap         int            `json:"self_reported_market_cap"`
}

type ResponseCryptocurrencyInfo struct {
	Status ResponseStatus `json:"status"`
	// key が動的だから 単純に struct に変換できない map[string]interface{} を使う
	Data map[string]interface{} `json:"data"`
}

type UseMetaData struct {
	Symbol   string
	CMCID    int
	Name     string
	WebSite  []string
	Logo     string
	Explorer []string
}

func ExampleGetMetadata(c *CoinMarketCap) {
	var CMCIDs []int
	CMCIDs = append(CMCIDs, 1)
	CMCIDs = append(CMCIDs, 1027)
	tmp := c.GetMetadata(CMCIDs)
	fmt.Printf("%+v\n", tmp)
}

// GetMetadata CoinMarketCap で メタデータを取得する
func (c *CoinMarketCap) GetMetadata(CMCIDs []int) []UseMetaData {
	const entryURL = "/v2/cryptocurrency/info"

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
	id := makeQueryParmCMCID(&CMCIDs)
	if id == "" {
		log.Fatalln("必須パラメータが無い")
	}
	q.Add("id", id)
	// q.Add("aux", "")
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
	// JSON にして出力してみる
	// 動的な key と status を削ったレスポンス
	var ccInfo []CryptocurrencyInfo
	// key が動的だから 単純に struct に変換できない map[string]interface{} を使う
	var responseJSON ResponseCryptocurrencyInfo
	if err := json.Unmarshal(respBody, &responseJSON); err != nil {
		log.Println(err)
	}
	// 動的な部分を扱う 第1戻り値が 動的な key
	for _, value := range responseJSON.Data {
		// 動的な key の中身 struct にするために 一度 []byte にし直す
		// 通常だと value は map[string]interface{} 型になってしまう
		byteValue, err := json.Marshal(value)
		if err != nil {
			log.Panicln(err)
		}
		// 一度 []byte にした部分を struct にする
		var cci CryptocurrencyInfo
		if err := json.Unmarshal(byteValue, &cci); err != nil {
			log.Println(err)
		}
		ccInfo = append(ccInfo, cci)
	}

	// fmt.Println("struct に入っているか?")
	// fmt.Printf("%+v\n", responseJSON)
	return generateSymbolAndMetaData(&ccInfo)
}

// makeQueryParmCMCID CMCID のスライスをカンマ区切りの1個の string にする
// makeQueryParmSymbol とアルゴリズムは同じ
func makeQueryParmCMCID(CMCIDs *[]int) string {
	var queryParmID string
	if len(*CMCIDs) == 0 {
		return queryParmID
	}

	for i, id := range *CMCIDs {
		if i != 0 {
			queryParmID += ","
		}
		queryParmID += strconv.Itoa(id)
	}
	return queryParmID
}

func generateSymbolAndMetaData(CCInfo *[]CryptocurrencyInfo) []UseMetaData {
	var metaDataMap []UseMetaData
	for _, info := range *CCInfo {
		tmp := UseMetaData{
			info.Symbol,
			info.ID,
			info.Name,
			info.Urls.Website,
			info.Logo,
			info.Urls.Explorer,
		}
		metaDataMap = append(metaDataMap, tmp)
	}
	return metaDataMap
}

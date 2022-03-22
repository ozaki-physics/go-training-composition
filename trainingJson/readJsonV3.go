package trainingJson

import (
	"encoding/json"
	"fmt"
)

// https://medium.com/@taka.abc.hiko/goで-jsonをパースしたい時のやつ-d54bee39064d
func ExampleDynamicJSONParse() {
	myJSON := `{
		"hello": {
			"a": "aaa",
			"b": "bbb"
		},
		"world": "none"
	 }`

	var result map[string]interface{}
	json.Unmarshal([]byte(myJSON), &result)
	fish := result["hello"].(map[string]interface{})
	for key, value := range fish {
		fmt.Println(key, value.(string))
	}
}

func ExampleDynamicJSONParse02() {
	type nono struct {
		ID     float64 `json:"id"`
		Name   string  `json:"name"`
		Symbol string  `json:"symbol"`
	}
	myJSON := `{
		"status":{
			"timestamp": "2022-03-13T14:17:16.630Z",
			"error_code":0
		},
		"data": {
			"1": {
				"id": 1,
				"name": "Bitcoin",
				"symbol": "BTC"
			},
			"1027": {
				"id": 1027,
				"name": "Ethereum",
				"symbol": "ETH"
			}
		}
	}`
	byteJSON := []byte(myJSON)
	var result map[string]interface{}
	json.Unmarshal(byteJSON, &result)
	// JSON から 一部取り出す
	aaa := result["data"]
	// 型変換 map[string]interface{} へ
	no := aaa.(map[string]interface{})
	for key, value := range no {
		fmt.Println(key)
		// fmt.Println(value)
		// 型変換 map[string]interface{} へ
		bbb := value.(map[string]interface{})
		ccc := nono{
			bbb["id"].(float64),
			bbb["name"].(string),
			bbb["symbol"].(string),
		}
		fmt.Printf("%+v\n", ccc)
	}
}

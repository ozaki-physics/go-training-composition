package assetCalc

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type coin struct {
	Symbol      string `json:"symbol"`
	Transaction []struct {
		Price     float64 `json:"price"`
		Size      float64 `json:"size"`
		Timestamp string  `json:"timestamp"`
	} `json:"transaction"`
}

type data struct {
	Coin map[string]interface{} `json:"data"`
}

type Position struct {
	Symbol       string
	AveragePrice float64
	Size         float64
}

func Calc(symbol string) (map[string]Position, error) {
	// 取引履歴を読み込む
	bytes, err := os.ReadFile("./assetCalc/crypto.json")
	if err != nil {
		log.Fatalln(err)
	}
	// JSON から struct にする
	var d data
	if err := json.Unmarshal(bytes, &d); err != nil {
		log.Fatalln(err)
	}

	var pMap = make(map[string]Position)
	// JSON の key が動的だから 再度 byte 化して struct 化する
	for key, value := range d.Coin {
		byteValue, err := json.Marshal(value)
		if err != nil {
			log.Fatalln(err)
		}
		var c coin
		if err := json.Unmarshal(byteValue, &c); err != nil {
			log.Fatalln(err)
		}

		var p Position
		if key != c.Symbol {
			log.Printf("key(%s) または 取引履歴のシンボル(%s)が間違っている", key, c.Symbol)
			continue
		}
		p.Symbol = c.Symbol
		// 取引履歴より 平均値と合計保有数を計算
		for _, t := range c.Transaction {
			p.AveragePrice = (p.AveragePrice*p.Size + t.Price*t.Size) / (p.Size + t.Size)
			p.Size += t.Size
		}

		if symbol == "all" {
			fmt.Printf("%s : %g : %g\n", p.Symbol, p.AveragePrice, p.Size)
		}
		pMap[key] = p
	}
	if symbol != "" && symbol != "all" {
		p := pMap[symbol]
		fmt.Printf("%s : %g : %g\n", p.Symbol, p.AveragePrice, p.Size)
	}
	return pMap, nil
}

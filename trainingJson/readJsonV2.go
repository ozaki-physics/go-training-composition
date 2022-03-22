package trainingJson

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// ExampleTags 構造体のフィールド名と JSON の名前が異なっていても動作する
// JSON の 順番と 構造体のフィールド順番 が一致していなくても JSON のメタデータに応じて 格納してくれる
// だから eins が Threethree に格納される
// また メタタグのオプション
// ハイフンを入れたから JSON から値を取ってこなくなっている and JSON にしたときに存在しなくなる
// omitempty は JSON にしたとき(エンコード)でフィールドを省略する
func ExampleTags() {
	var j = []byte(`[
	{"no": 1, "example": {"eins": 111, "zwei": "1個目のツヴァイ", "drei": 333, "vier": "1個目のフィーア"}},
	{"no": 2, "example": {"eins": 222, "zwei": "2個目のツヴァイ", "drei": 444, "vier": "2個目のフィーア"}},
	{"no": 3, "example": {"eins": 555, "zwei": "3個目のツヴァイ", "vier": "3個目のフィーア"}}
]`)

	type aaa struct {
		Num int `json:"no"`
		Ex  struct {
			Oneone     int    `json:"drei,omitempty"` // 順番が入れ替わってる
			Twotwo     string `json:"zwei"`
			Threethree int    `json:"eins"` // 順番が入れ替わってる
			Fourfour   string `json:"-"`    // JSON が格納されなくなる
		} `json:"example"`
	}
	var many []aaa
	json.Unmarshal(j, &many)
	fmt.Printf("%+v\n", many)

	v, _ := json.Marshal(many)
	fmt.Println(string(v)) // omitempty をつけたから JSON に存在しなくなっている
}

func ExampleRawMessage() {
	type aaa struct {
		Num int `json:"no"`
	}
	var a01 aaa
	var a02 aaa

	// RawMessage を定義する?
	a := json.RawMessage(`{"no":1}`)
	// []byte になっている JSON を取り出す?
	b, _ := a.MarshalJSON()
	fmt.Println("bだよ: " + string(b))

	// フツーに JSON をデコード
	json.Unmarshal(a, &a01)
	fmt.Printf("%+v\n", a01)

	// []byte になっている JSON を書き換える?
	a.UnmarshalJSON([]byte(`{"nono":2}`))
	// []byte になっている JSON を取り出す?
	c, _ := a.MarshalJSON()
	fmt.Println("cだよ: " + string(c))
	// フツーに JSON をデコード
	json.Unmarshal(a, &a02)
	fmt.Printf("%+v\n", a02)
	fmt.Println("メタタグは no なのに JSON は nono だから取り出せてない")
}

func ExampleRawMessageMarshal() {
	h := json.RawMessage(`{"precomputed": true}`)

	c := struct {
		Header *json.RawMessage `json:"header"`
		Body   string           `json:"body"`
	}{Header: &h, Body: "Hello Gophers!"}

	b, err := json.MarshalIndent(&c, "", "\t")
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
}

func ExampleRawMessageUnmarshal() {
	type Color struct {
		Space string
		Point json.RawMessage // delay parsing until we know the color space
	}
	type RGB struct {
		R uint8
		G uint8
		B uint8
	}
	type YCbCr struct {
		Y  uint8
		Cb int8
		Cr int8
	}

	var j = []byte(`[
	{"Space": "YCbCr", "Point": {"Y": 255, "Cb": 0, "Cr": -10}},
	{"Space": "RGB",   "Point": {"R": 98, "G": 218, "B": 255}}
]`)
	var colors []Color
	err := json.Unmarshal(j, &colors)
	if err != nil {
		log.Fatalln("error:", err)
	}

	for _, c := range colors {
		// JSON の中の値に応じて 動的に struct を切り替えられる
		var dst interface{}
		switch c.Space {
		case "RGB":
			dst = new(RGB)
		case "YCbCr":
			dst = new(YCbCr)
		}
		err := json.Unmarshal(c.Point, dst)
		if err != nil {
			log.Fatalln("error:", err)
		}
		fmt.Println(c.Space, dst)
	}
}

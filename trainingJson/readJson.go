package trainingJson

import (
	"encoding/json"
	"fmt"
	"github.com/ozaki-physics/go-training-composition/utils"
	"os"
)

// Example 公式ドキュメントにあったサンプル
// See: https://pkg.go.dev/encoding/json@go1.17#example-Unmarshal
func Example() {
	var jsonBlob = []byte(`[
	{"Name": "Platypus", "Order": "Monotremata"},
	{"Name": "Quoll",    "Order": "Dasyuromorphia"}
]`)
	// struct の定義
	type Animal struct {
		Name  string
		Order string
	}
	// JSON が配列で始まるから 配列で初期化する
	var animals []Animal
	err := json.Unmarshal(jsonBlob, &animals)
	utils.ErrCheck(err)
	fmt.Printf("%+v\n", animals)
}

// ReadJson01 JSON を読み込む
func ReadJson01() {
	bytes, err01 := os.ReadFile("trainingJson/filepath01.json")
	utils.ErrCheck(err01)

	// struct の定義
	// メソッド内で struct を定義するなら メソッド内で使われる変数の順番に気をつけること
	// メソッド外で package の struct とするなら 順番は気にしなくていい

	// struct の type は 小文字で書いて private でもいいが
	// 要素は 大文字で書いて public にしないといけない
	// 要素名と JSON の key 名は同じにしないといけない(大文字小文字は関係ない)

	// struct02 は struct01 で使うから先に定義する
	type struct02 struct {
		In01 string
		In02 int
		In03 []string
	}

	type struct01 struct {
		Out01 struct02
		Out02 struct02
	}

	// JSON が 配列じゃなくてオブジェクトだから [] をつけない
	var outStruct struct01
	// struct のポインタを渡す
	err02 := json.Unmarshal(bytes, &outStruct)
	utils.ErrCheck(err02)
	// JSON(実体は struct) の全体を出力
	fmt.Printf("%+v\n", outStruct)
	// JSON(実体は struct) の 一部を出力
	fmt.Println(outStruct.Out01.In02)
}

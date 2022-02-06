package main

import (
	"fmt"
)

func output(i I) {
	fmt.Printf("値: %v, 型: %T\n", i, i)
}

// インタフェース宣言
type I interface {
	M()
}

// struct 宣言 T 型は string を持つ
type T struct {
	S string
}

// レシーバとして *T 型を指定して メソッド M() を実装
func (t *T) M() {
	// nil だったときの処理
	if t == nil {
		fmt.Println("<nil>")
		return
	}
	fmt.Println(t.S)
}

func main() {
	// インタフェース型の変数宣言
	var i I
	// struct のポインタを宣言
	var t *T
	// インタフェースに格納
	i = t
	// この時点ではまだ struct の実体?がない
	output(i)
	// 出力 値: <nil>, 型: *main.T
	i.M()
	// 出力 <nil>

	i = &T{"hello"}
	output(i)
	// 出力 値: &{hello}, 型: *main.T
	i.M()
	// 出力 hello
}

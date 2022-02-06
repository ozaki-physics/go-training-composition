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
	fmt.Println(t.S)
}

// struct 宣言 F 型は float64 を持つ
type F float64

// レシーバとして F 型を指定して メソッド M() を実装
func (f F) M() {
	fmt.Println(f)
}

func main() {
	// インタフェース型の変数宣言
	var i I
	// T 型のポインタを i に格納
	i = &T{"Hello"}
	// (value, type) を確認
	output(i)
	// 出力 値: &{Hello}, 型: *main.T
	i.M()
	// 出力 Hello

	// F 型を i に格納
	i = F(3.14)
	output(i)
	// 出力 値: 3.14, 型: main.F
	i.M()
	// 出力 3.14
}

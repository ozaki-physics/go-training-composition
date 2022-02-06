package main

import (
	"fmt"
)

// 型 switch と型アサーション
func do(i interface{}) {
	// i.(type) でとりあえず実体の型を受け取って case で判定する
	switch v := i.(type) {
	case int:
		fmt.Printf("型: %T, 値: %v, 2乗: %v\n", v, v, v*2)
	case string:
		fmt.Printf("型: %T, 値: %v, %q is %v bytes long\n", v, v, v, len(v))
	default:
		fmt.Printf("%T型なんて知らん!\n", v)
	}
}

func main() {
	// 実体は string 型の 任意の型を保持できる空インタフェースを宣言
	var i interface{} = "hello"

	// 型アサーション
	s, ok := i.(string)
	fmt.Println(s, ok)
	// 出力 hello true

	// 型アサーション
	f, ok := i.(float64)
	fmt.Println(f, ok)
	// 出力 0 false 実体は string だから

	do(3)
	// 出力 型: int, 値: 3, 2乗: 6
	do("hello")
	// 出力 型: string, 値: hello, "hello" is 5 bytes long
	do(true)
	// 出力 bool型なんて知らん!
}

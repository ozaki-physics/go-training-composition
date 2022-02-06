package main

import (
	"fmt"
)

func output(i interface{}) {
	fmt.Printf("値: %v, 型: %T\n", i, i)
}

func main() {
	var i interface{}
	output(i)
	// 出力 値: <nil>, 型: <nil>

	i = 42
	output(i)
	// 出力 値: 42, 型: int

	i = "hello"
	output(i)
	// 出力 値: hello, 型: string
}

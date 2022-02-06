package main

import (
	"fmt"
)

func adder() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

func main() {
	// pos には sum += x の関数が入っている
	pos := adder()
	for i := 0; i < 3; i++ {
		// pos(i) は sum
		// i は x
		fmt.Println(pos(i))
	}
	// 出力
	// 0
	// 1
	// 3

	fmt.Println(pos(10))
	// 出力 13 なぜなら今保持してる3が足されるから
}

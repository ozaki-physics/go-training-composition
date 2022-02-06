package main

import (
	"fmt"
)

// 渡された引数の順序を入れ替える
func swap(x, y string) (string, string) {
	return y, x
}

func main() {
	// 入れ替えたものを格納する
	a, b := swap("hello", "world")
	fmt.Println(a, b)
}

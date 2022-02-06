package main

import (
	"fmt"
)

// 戻り値を関数定義の時点で宣言する
func split(in int) (out01, out02 int) {
	out01 = in * 4
	out02 = in + 4
	return
}

func main() {
	a, b := split(3)
	fmt.Println(a, b)
}

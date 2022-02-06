package goTour

import (
	"fmt"
)

func outputEx25(i interface{}) {
	fmt.Printf("値: %v, 型: %T\n", i, i)
}

func ex25() {
	var i interface{}
	outputEx25(i)
	// 出力 値: <nil>, 型: <nil>

	i = 42
	outputEx25(i)
	// 出力 値: 42, 型: int

	i = "hello"
	outputEx25(i)
	// 出力 値: hello, 型: string
}

package goTour

import (
	"fmt"
)

func ex06() {
	var a = `
	hello
	world
	`
	fmt.Println(a)
	// 以下が出力
	//
	// hello
	// world
	//

	b := false
	fmt.Printf("Type: %T Value: %v\n", b, b)
}

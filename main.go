package main

import (
	"fmt"
	"github.com/ozaki-physics/go-training-composition/package01"
	"github.com/ozaki-physics/go-training-composition/package02"
)

func main() {
	fmt.Println("Hello World!")
	fmt.Println(package01.Const_big)
	// これはエラー
	// fmt.Println(package01.const_small)
	fmt.Println(package02.Const_big)
	// package02.Sample_server()
}

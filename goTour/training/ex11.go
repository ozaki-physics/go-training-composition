package main

import (
	"fmt"
)

func aaa() {
	defer fmt.Println("world")
	fmt.Println("hello")
	// hello worldと出力される
}

func bbb() {
	for i := 0; i < 10; i++ {
		defer fmt.Println(i)
	}
	// 結果は9,8,7,6,5,4,3,2,1,0
}

func main() {
	aaa()
	bbb()
}

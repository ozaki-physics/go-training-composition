package goTour

import (
	"fmt"
)

func aaaEx11() {
	defer fmt.Println("world")
	fmt.Println("hello")
	// hello worldと出力される
}

func bbbEx11() {
	for i := 0; i < 10; i++ {
		defer fmt.Println(i)
	}
	// 結果は9,8,7,6,5,4,3,2,1,0
}

func ex11() {
	aaaEx11()
	bbbEx11()
}

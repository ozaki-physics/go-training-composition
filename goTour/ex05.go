package goTour

import (
	"fmt"
	"math/cmplx"
)

func ex05() {
	var a int
	fmt.Println(a)
	a = 1
	fmt.Println(a)

	var i, j = 1, 2
	fmt.Println(i, j)

	b := 3
	fmt.Println(b)

	var c = 4
	fmt.Println(c)

	var z complex128 = cmplx.Sqrt(-5 + 12i)
	fmt.Println(z)
	// 2+3i が出力される

	var (
		ToBe   bool   = false
		MaxInt uint64 = 1<<64 - 1
	)
	fmt.Println(ToBe, MaxInt)
}

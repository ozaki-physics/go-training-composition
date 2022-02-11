package goTour

import (
	"fmt"
)

func ex14() {
	var a [3]int
	fmt.Println(a[0], a[1])
	// 出力 0 0
	primes := [6]int{1, 2, 3, 4, 5, 6}
	fmt.Println(primes[1], primes[3])
	// 出力 2 4

	var s []int = primes[1:4]
	fmt.Println(s)
	// 出力 [2 3 4]

	names := [3]string{
		"Alice",
		"Bob",
		"Carol",
	}
	// 改行して配列を作った場合は最後にもコンマがいる
	fmt.Println(names)
	// 出力 [Alice Bob Carol]
	b := names[1:3]
	fmt.Println(b)
	// 出力 [Bob Carol]
	b[0] = "XXX"
	fmt.Println(names)
	// 出力 [Alice XXX Carol]

	// 配列リテラルを作成し、同時に配列リテラルを参照するスライスを作成する
	c := []struct {
		i int
		b bool
	}{
		{2, true},
		{3, true},
		{5, true},
		{7, false},
		{11, false},
		{13, true},
	}
	fmt.Println(c)
	// 出力 [{2 true} {3 true} {5 true} {7 false} {11 false} {13 true}]
}

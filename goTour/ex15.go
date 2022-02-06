package goTour

import (
	"fmt"
	"strings"
)

func ex15() {
	a := make([]int, 5)
	fmt.Println(a, len(a), cap(a))
	// 出力 [0 0 0 0 0] 5 5
	b := make([]int, 0, 5)
	fmt.Println(b, len(b), cap(b))
	// 出力 [] 0 5

	board := [][]string{
		[]string{"-", "-", "-"},
		[]string{"-", "-", "-"},
		[]string{"-", "-", "-"},
	}
	for i := 0; i < len(board); i++ {
		fmt.Printf("%s\n", strings.Join(board[i], " "))
	}
	// 出力
	// - - -
	// - - -
	// - - -

	var c []int
	fmt.Println(c, len(c), cap(c))
	// 出力 [] 0 0
	c = append(c, 0, 1, 2)
	fmt.Println(c, len(c), cap(c))
	// 出力 [0 1 2] 3 4
	// 追加するリストの容量が 元のリストの容量より大きい場合 メモリ上はリストを割り当て直すみたい
	// その時に 容量が変化するのかも

	d := make([]int, 3)
	d[1] = 10
	for i, value := range d {
		fmt.Printf("%d: %d\n", i, value)
	}
	// 出力
	// 0: 0
	// 1: 10
	// 2: 0
}

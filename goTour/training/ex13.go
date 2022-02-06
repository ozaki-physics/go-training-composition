package main

import (
	"fmt"
)

type Vertex struct {
	X int
	Y int
}

type Vertex2 struct {
	X int
	Y int
	Z int
}

func main() {
	fmt.Println(Vertex{1, 2})
	// 出力 {1 2}

	v := Vertex{1, 2}
	v.X = 4
	fmt.Println(v.X)
	// 出力 4

	w := Vertex{3, 5}
	p := &w
	fmt.Println((*p).X)

	v2 := Vertex{X: 6, Y: 7}
	fmt.Println(v2)
	// 出力 {6 7}
	v3 := Vertex{X: 8}
	fmt.Println(v3)
	// 出力 {8 0}

	v4 := Vertex2{Z: 9, Y: 10}
	fmt.Println(v4)
	// 出力 {0 9 10}

	v5 := &Vertex{2, 1}
	fmt.Println(v5)
	// 出力 &{2 1}
}

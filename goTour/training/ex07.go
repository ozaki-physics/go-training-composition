package main

import (
	"fmt"
)

func main() {
	var a int
	b := a
	fmt.Printf("%T\n", b)
	// int
	i := 42
	j := 3.142
	k := 0.867 + 0.5i
	fmt.Printf("%T, %T, %T\n", i, j, k)
	// int, float64, complex128

	const Pi = 3.14
	fmt.Println(Pi)
}

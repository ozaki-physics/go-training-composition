package main

import (
	"fmt"
)

func add(x int, y int) int {
	return x + y
}

func add02(x, y int) int {
	return x + y
}

func main() {
	ans := add(2, 3)
	fmt.Println(ans)
	ans02 := add02(3, 4)
	fmt.Println(ans02)
}

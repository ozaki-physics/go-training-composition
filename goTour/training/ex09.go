package main

import (
	"fmt"
	"math"
)

func sqrt(x float64) string {
	if x < 0 {
		return sqrt(-x) + "i"
	}
	return fmt.Sprint(math.Sqrt(x))
}

func aaa(x, n, lim float64) float64 {
	if v := math.Pow(x, n); lim < v {
		return v
	}
	return 10
}

func bbb(x, n, lim float64) float64 {
	if v := math.Pow(x, n); lim < v {
		return v
	} else if v == lim {
		return v * 2
	}
	return 10
}

func main() {
	fmt.Println(sqrt(2))
	fmt.Println(sqrt(-2))

	a := -2
	if a < 0 {
		fmt.Println("負")
	} else {
		fmt.Println("正")
	}

	fmt.Println(aaa(2, 3, 7))
	fmt.Println(aaa(2, 3, 9))
	fmt.Printf("else if を通って%g\n", bbb(2, 3, 8))

	fmt.Println(
		math.Pow(3, 2),
		math.Pow(3, 3),
	)
}

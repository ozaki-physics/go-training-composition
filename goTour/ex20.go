package goTour

import (
	"fmt"
)

func adderEx20() func(int) int {
	sum := 5
	fmt.Println("-1-", sum)
	return func(x int) int {
		fmt.Println("-2-", sum)
		sum += x
		fmt.Println("-3-", x)
		fmt.Println("-4-", sum)
		return sum
	}
}

func ex20() {
	fmt.Println("when 01")
	pos := adderEx20()
	fmt.Println("when 02")
	fmt.Println(pos(10))
	fmt.Println("when 03")
	for i := 0; i < 3; i++ {
		fmt.Println("when 04")
		fmt.Println(pos(i))
	}
	// 出力
	// when 01
	// -1- 5
	// when 02
	// -2- 5
	// -3- 10
	// -4- 15
	// 15
	// when 03
	// when 04
	// -2- 15
	// -3- 0
	// -4- 15
	// 15
	// when 04
	// -2- 15
	// -3- 1
	// -4- 16
	// 16
	// when 04
	// -2- 16
	// -3- 2
	// -4- 18
	// 18
}

package goTour

import (
	"fmt"
)

func ex08() {
	var sum int
	for i := 0; i < 10; i++ {
		sum += i
	}
	fmt.Println(sum)

	var sum02 = 1
	for sum02 < 10 {
		sum02++
	}
	fmt.Println(sum02)

	var sum03 = 1
	for sum03 < 10 {
		sum03++
	}
	fmt.Println(sum03)
}

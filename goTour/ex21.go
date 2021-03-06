package goTour

import (
	"fmt"
)

func fibonacci() func() int {
	a := 0
	b := 1
	return func() int {
		out := a
		tmp := a + b
		a = b
		b = tmp
		return out
	}
}

func ex21() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}

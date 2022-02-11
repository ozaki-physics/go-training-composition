package goTour

import (
	"fmt"
)

func ex31() {
	ch := make(chan int, 3)
	ch <- 1
	ch <- 2
	ch <- 3
	fmt.Println(<-ch)
	// 先に ch から受け取ってからじゃないとエラーになる
	ch <- 4
	fmt.Println(<-ch)
	fmt.Println(<-ch)
	fmt.Println(<-ch)
}

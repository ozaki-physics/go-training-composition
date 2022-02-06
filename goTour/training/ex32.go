package main

import (
	"fmt"
)

// チャネルを使った自作フィボナッチ数列
// close の扱いが微妙な気がする
func fibonacci(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
	close(c)
}
func main() {
	// バッファを変えると n 個目の数列まで計算する
	c := make(chan int, 3)
	go fibonacci(cap(c), c)
	for i := range c {
		fmt.Println(i)
	}
}

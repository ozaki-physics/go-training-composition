package main

import (
	"fmt"
)

func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		// どれかの case が実行できるようになるまで goroutine スレッドを待機させる
		select {
		case c <- x: // x を c チャネルに送信したら実行する
			x, y = y, x+y
		case <-quit: // quit チャネルを値を受け取ったら実行する
			fmt.Println("quit")
			return
		}
	}
}

func main() {
	// 2個のチャネルを生成
	c := make(chan int)
	quit := make(chan int)
	// for文中は c チャネルに値を渡して 終わったら quit チャネルに値を渡す という goroutine スレッドを作る
	go func() {
		for i := 0; i < 4; i++ {
			fmt.Println(<-c)
		}
		quit <- 0
	}()

	fibonacci(c, quit)
	// 出力
	// 0
	// 1
	// 1
	// 2
	// quit
}

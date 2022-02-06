package main

import (
	"fmt"
)

// int配列 を受け取るのと 返す先のチャネルを指定されて 指定されたチャネルへ合計値を返す
func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum // send sum to c
}

func main() {
	s := []int{7, 2, 8, -9, 4, 0}

	// c というチャネルを作成する
	c := make(chan int)
	// int配列を2等分して c チャネルへ送る
	go sum(s[:len(s)/2], c)
	go sum(s[len(s)/2:], c)
	// c チャネルから受け取る
	x, y := <-c, <-c // receive from c
	// 受け取った値を出力する
	fmt.Println(x, y, x+y)
	// 出力 -5 17 12
}

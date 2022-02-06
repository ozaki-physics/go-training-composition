package goTour

import (
	"fmt"
	"sync"
	"time"
)

// SafeCounter 型は排他制御ができて key の数を保持する
type SafeCounter struct {
	mu sync.Mutex
	v  map[string]int
}

// Inc メソッドは指定された key のカウンタを増やす
func (c *SafeCounter) Inc(key string) {
	// 一度に1つの goroutine しか map にアクセスできないようにロックする
	c.mu.Lock()
	c.v[key]++
	c.mu.Unlock()
}

// Value メソッドは指定された key のカウンタ値を返す
func (c *SafeCounter) Value(key string) int {
	// 一度に1つの goroutine しか map にアクセスできないようにロックする
	c.mu.Lock()
	// mutex が Unlock されることを保証するために defer を使うこともできる
	defer c.mu.Unlock()
	return c.v[key]
}

func ex35() {
	c := SafeCounter{
		v: make(map[string]int),
	}
	for i := 0; i < 10; i++ {
		go c.Inc("key")
		fmt.Println("今は", i, "番目")
	}
	// main とは 別の goroutine の処理によって結果が変わるため 数秒待つ
	fmt.Println(time.Second)
	// 出力 1s
	time.Sleep(time.Second)
	fmt.Println("keyの数は", c.Value("key"))
	// 出力 keyの数は 10
}

package main

import (
	"fmt"
	"time"
)

// 自作で select の default を実験した
func main() {
	tick := time.Tick(100 * time.Millisecond)
	boom := time.After(500 * time.Millisecond)
	for {
		select {
		case <-tick:
			fmt.Println("チック.")
		case <-boom:
			fmt.Println("ボーン!")
			return
		default:
			fmt.Println("    .")
			time.Sleep(50 * time.Millisecond)
		}
	}
}

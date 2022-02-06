package main

import (
	"fmt"
	"strings"
)

// 単語数を数えるメソッド
func WordCount(s string) map[string]int {
	a := strings.Fields(s)
	b := make(map[string]int)
	for i := 0; i < len(a); i++ {
		count := 0
		for j := 0; j < len(a); j++ {
			if a[j] == a[i] {
				count += 1
			}
		}
		b[string(a[i])] = count
	}
	return b
}

func main() {
	str := "Hello World Hello"
	fmt.Println(WordCount(str))
	// 出力 map[Hello:2 World:1]
}

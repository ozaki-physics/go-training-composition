package main

import (
	"fmt"
)

type Vertex struct {
	x, y float64
}

func main() {
	// key は string, value は Vertex で宣言
	// まだ nil な map で使えない
	var m map[string]Vertex
	// make() で使える map にする
	m = make(map[string]Vertex)
	m["key01"] = Vertex{
		3.14, -2.71,
	}
	m["key02"] = Vertex{1, 3}
	fmt.Println(m["key01"])
	// 出力 {3.14 -2.71}
	fmt.Println(m)
	// 出力 map[key01:{3.14 -2.71} key02:{1 3}]

	// make() を使わない方法
	var p = map[string]Vertex{
		"key01": Vertex{1, 2},
		"key02": Vertex{3, 4},
	}
	fmt.Println(p)
	// 出力 map[key01:{1 2} key02:{3 4}]

	// 型を省略した書き方
	var r = map[string]Vertex{
		"key01": {1.41, 1.73},
		"key02": {5, 6},
	}
	fmt.Println(r)
	// 出力 map[key01:{1.41 1.73} key02:{5 6}]

	a := make(map[string]int)
	a["Answer"] = 42
	fmt.Println(a)
	// 出力 map[Answer:42]
	delete(a, "Answer")
	v, ok := a["Answer"]
	fmt.Println(v, ok)
	// 出力 0 false
}

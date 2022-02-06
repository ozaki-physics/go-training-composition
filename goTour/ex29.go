package goTour

import (
	"fmt"
	"io"
	"strings"
)

func ex29() {
	r := strings.NewReader("Hello, Reader!")

	// ゼロ値が8個のバイト型の配列
	b := make([]byte, 8)
	fmt.Println(b, len(b), cap(b))
	// 出力 [0 0 0 0 0 0 0 0] 8 8
	for {
		// 8個のバイトごとに取り出す感じ
		n, err := r.Read(b)
		fmt.Printf("n: %v, err: %v\n", n, err)
		fmt.Printf("b[:n]: %q\n", b[:n])
		if err == io.EOF {
			break
		}
	}
	// 出力
	// n: 8, err: <nil>
	// b[:n]: "Hello, R"
	// n: 6, err: <nil>
	// b[:n]: "eader!"
	// n: 0, err: EOF
	// b[:n]: ""
}

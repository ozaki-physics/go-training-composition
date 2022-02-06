package main

import (
	"fmt"
)

func main() {
	i, j := 6, 12

	// 変数iのポインタをpに渡す
	// pはiのポインタ
	p := &i
	fmt.Println(*p)
	// 出力 6 fmt.Println(i) と同義
	fmt.Println(p)
	// fmt.Println(*i) // エラー
	// 出力 0xc000018038 fmt.Println(&i) と同義
	fmt.Println(i)
	// 出力 6
	fmt.Println("-1--")

	// 21をポインタへ代入する
	// ポインタ元のiも21になる
	// 同じiの場所だからポインタも変化しない
	*p = 16
	fmt.Println(*p)
	// 出力 16 fmt.Println(i) と同義
	fmt.Println(p)
	// 出力 0xc000018038 fmt.Println(&i) と同義
	fmt.Println(i)
	// 出力 16
	fmt.Println("-2--")

	// 変数jのポインタをpに渡す
	// pはjのポインタ
	// 異なるjの場所だからポインタが変わる
	p = &j
	fmt.Println(*p)
	// 出力 12 fmt.Println(j) と同義
	fmt.Println(p)
	// 出力 0xc00010c008 fmt.Println(&j) と同義
	fmt.Println("-3--")

	// ポインタをポインタに再代入
	*p = *p / 2
	// j = j / 2 と同義
	fmt.Println(j)
	// 出力 6

	fmt.Println(&j)
}

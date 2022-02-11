package goTour

import (
	"fmt"
)

func outputEx24(i I_Ex24) {
	fmt.Printf("値: %v, 型: %T\n", i, i)
}

// インタフェース宣言
type I_Ex24 interface {
	M()
}

// struct 宣言 T_Ex24 型は string を持つ
type T_Ex24 struct {
	S string
}

// レシーバとして *T_Ex24 型を指定して メソッド M() を実装
func (t *T_Ex24) M() {
	// nil だったときの処理
	if t == nil {
		fmt.Println("<nil>")
		return
	}
	fmt.Println(t.S)
}

func ex24() {
	// インタフェース型の変数宣言
	var i I_Ex24
	// struct のポインタを宣言
	var t *T_Ex24
	// インタフェースに格納
	i = t
	// この時点ではまだ struct の実体?がない
	outputEx24(i)
	// 出力 値: <nil>, 型: *main.T_Ex24
	i.M()
	// 出力 <nil>

	i = &T_Ex24{"hello"}
	outputEx24(i)
	// 出力 値: &{hello}, 型: *main.T_Ex24
	i.M()
	// 出力 hello
}

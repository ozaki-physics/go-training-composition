package goTour

import (
	"fmt"
	"math"
)

type VertexEx22 struct {
	X, Y float64
}

// メソッドの作成
// Abs01 メソッドは v という名前の VertexEx22 型のレシーバを持つ。型は任意。
func (v VertexEx22) Abs01() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// 関数の作成
func Abs02(v VertexEx22) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// レシーバをポインタした書き方
func (v *VertexEx22) Scale(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

func ex22() {
	v := VertexEx22{3, 4}
	fmt.Println(v.Abs01())
	// 出力 5
	fmt.Println(Abs02(v))
	// 出力 5
	v.Scale(10)
	fmt.Println(v)
	// 出力 {30 40}
}

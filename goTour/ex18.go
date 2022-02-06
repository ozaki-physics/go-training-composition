package goTour

import (
	"fmt"
	"math"
)

func compute(fn func(float64, float64) float64) float64 {
	return fn(3, 4)
}

func ex18() {
	hypot := func(x, y float64) float64 {
		return math.Sqrt(x*x + y*y)
	}

	fmt.Println(hypot(3, 4))
	// 出力 5 なぜなら3^2+4^2=5^2 の平方根
	fmt.Println(compute(hypot))
	// 出力 5 なぜなら渡す値はcomputeで既に決まっているから
	fmt.Println(compute(math.Pow))
	// 出力 81 なぜなら3^4=81
}

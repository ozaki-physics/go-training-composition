package trainingEmbedding

import (
	"fmt"
)

func Example01() {
	fmt.Println("mainEmbedding")

	p := Parent{"Universe", 138}
	pProfile := p.WhoParent()
	fmt.Println(pProfile)
	// 親の who メソッド
	// 名前: Universe, 年齢: 138

	fmt.Println("---")

	c := Child{Parent{"Earth", 45}}
	cProfile := c.WhoParent()
	fmt.Println(cProfile)
	cProfileChild := c.WhoChild()
	fmt.Println(cProfileChild)
	// 親の who メソッド
	// 名前: Earth, 年齢: 45
	// 子の who メソッド
	// 名前: Earth, 年齢: 45

	fmt.Println("---")

	c02 := Child02{"Jupiter", Parent{"Earth", 46}}
	c02Profile := c02.WhoParent()
	fmt.Println(c02Profile)
	c02ProfileChild := c02.WhoChild02()
	fmt.Println(c02ProfileChild)
	// 親の who メソッド
	// 名前: Earth, 年齢: 46
	// 子の who メソッド
	// 名前: Jupiter, 年齢: 46
}

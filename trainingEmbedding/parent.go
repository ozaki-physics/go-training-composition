package trainingEmbedding

import (
	"fmt"
	"strconv"
)

type Parent struct {
	Name string
	Age  int
}

func (p *Parent) WhoParent() string {
	profile := "名前: " + p.Name + ", 年齢: " + strconv.Itoa(p.Age)
	fmt.Println("親の who メソッド")
	return profile
}

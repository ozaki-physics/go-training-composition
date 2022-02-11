package trainingEmbedding

import (
	"fmt"
	"strconv"
)

type Child struct {
	Parent
}

func (c *Child) WhoChild() string {
	profile := "名前: " + c.Name + ", 年齢: " + strconv.Itoa(c.Age)
	fmt.Println("子の who メソッド")
	return profile
}

type Child02 struct {
	Name string
	Parent
}

func (c *Child02) WhoChild02() string {
	profile := "名前: " + c.Name + ", 年齢: " + strconv.Itoa(c.Age)
	fmt.Println("子の who メソッド")
	return profile
}

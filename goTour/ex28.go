package goTour

import (
	"fmt"
	"time"
)

type MyError struct {
	When time.Time
	What string
}

func (e *MyError) Error() string {
	return fmt.Sprintf("いつ: %v, 何が: %s", e.When, e.What)
}

func run() error {
	return &MyError{
		time.Now(),
		"実行できません",
	}
}

func ex28() {
	// err が nil じゃない つまり エラーが起こっている
	if err := run(); err != nil {
		fmt.Println(err)
		// 出力 いつ: 2021-01-16 00:01:35.7540284 +0000 UTC m=+0.000046701, 何が: 実行できません
	}
}

package pkg03_test

import (
	"fmt"
	// package 名を別にしているから ドットを付けないと クラス名書かないといけなくなる
	. "github.com/ozaki-physics/go-training-composition/pkg03"
	"testing"
)

// TestMain を書かない場合
func TestSubPublicTrue(t *testing.T) {
	patterns := []struct {
		a        int
		b        int
		expected int
	}{
		{4, 2, 2},
		{1, 2, -1},
		{10, -2, 12},
		{-10, -2, -8},
	}

	for idx, pattern := range patterns {
		t.Run(fmt.Sprintf("いけた?"), func(t *testing.T) {
			actual := SubPublicTrue(pattern.a, pattern.b)
			// 本当は err の戻り値も付けたいが自作メソッドに実装していない
			// if err != nil {
			// 	t.Fatal("エラーだよ! エラーが無(nil)じゃないから")
			// }
			if actual != pattern.expected {
				t.Errorf("インデックス %d で期待は %d なのに実際は %d だった", idx, pattern.expected, actual)
			}
		})
	}
}

func TestSubPublicFalse(t *testing.T) {
	patterns := []struct {
		a        int
		b        int
		expected int
	}{
		{4, 2, 2},
		{1, 2, -1},
		{10, -2, 12},
		{-10, -2, -8},
	}

	for idx, pattern := range patterns {
		t.Run(fmt.Sprintf("いけた?"), func(t *testing.T) {
			actual := SubPublicFalse(pattern.a, pattern.b)
			// 本当は err の戻り値も付けたいが自作メソッドに実装していない
			// if err != nil {
			// 	t.Fatal("エラーだよ! エラーが無(nil)じゃないから")
			// }
			if actual != pattern.expected {
				t.Errorf("インデックス %d で期待は %d なのに実際は %d だった", idx, pattern.expected, actual)
			}
		})
	}
}

// go-training-composition にいる状態で go test -v ./pkg03/

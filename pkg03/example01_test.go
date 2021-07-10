package pkg03_test

import (
	"fmt"
	// package 名を別にしているから ドットを付けないと クラス名書かないといけなくなる
	. "github.com/ozaki-physics/go-training-composition/pkg03"
	"os"
	"testing"
)

// TestMain を書く場合
func TestMain(m *testing.M) {
	fmt.Println("test の前処理")
	code := m.Run()
	fmt.Println("test の後処理")
	// Exit を書くのは慣例
	os.Exit(code)
}

func TestAddPublicTrue(t *testing.T) {
	patterns := []struct {
		a        int
		b        int
		expected int
	}{
		{1, 2, 3},
		{10, -2, 8},
		{-10, -2, -12},
	}

	for idx, pattern := range patterns {
		actual := AddPublicTrue(pattern.a, pattern.b)
		// 結果が一致するか確かめる
		if pattern.expected != actual {
			t.Errorf("インデックス %d で期待は %d なのに実際は %d だった", idx, pattern.expected, actual)
		}
	}
}

func TestAddPublicFalse(t *testing.T) {
	patterns := []struct {
		a        int
		b        int
		expected int
	}{
		{1, 2, 3},
		{10, -2, 8},
		{-10, -2, -12},
	}

	for idx, pattern := range patterns {
		actual := AddPublicFalse(pattern.a, pattern.b)
		// 結果が一致するか確かめる
		if pattern.expected != actual {
			t.Errorf("インデックス %d で期待は %d なのに実際は %d だった", idx, pattern.expected, actual)
		}
	}
}

// go-training-composition にいる状態で go test -v ./pkg03/

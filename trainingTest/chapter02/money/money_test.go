package money_test

import (
	"testing"

	. "github.com/ozaki-physics/go-training-composition/trainingTest/chapter02/money"
)

func TestCreateDollar(t *testing.T) {
	t.Run("コンストラクタ", func(t *testing.T) {
		five := CreateDollar(5)
		expected := 5
		if five.Amout != expected {
			t.Errorf("期待値 %d で実際は %d", expected, five.Amout)
		}
	})
}

func TestTimes(t *testing.T) {
	t.Run("もとの$5は変わらず 何度でもかけ算ができる", func(t *testing.T) {
		// 型を明示しても書いてみる
		var five Dollar = CreateDollar(5)
		var product Dollar = five.Times(2)

		expected := 10
		if product.Amout != expected {
			t.Errorf("期待値 %d で実際は %d", expected, product.Amout)
		}

		product = five.Times(3)
		expected = 15
		if product.Amout != expected {
			t.Errorf("期待値 %d で実際は %d", expected, product.Amout)
		}
	})
}

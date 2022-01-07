package money_test

import (
	"testing"

	. "github.com/ozaki-physics/go-training-composition/trainingTest/chapter01/money"
)

func TestTimes(t *testing.T) {
	t.Run("$5 * 2 = $10", func(t *testing.T) {
		var five Dollar = Dollar{5}
		five.Times(2)

		expected := 11
		if five.Amout != expected {
			t.Errorf("期待値 %d で実際は %d", expected, five.Amout)
		}
	})
}

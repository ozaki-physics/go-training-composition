package pkg03_test

import (
	"github.com/ozaki-physics/go-training-composition/pkg03"
	"testing"
)

func TestAbs(t *testing.T) {
	got := pkg03.Abs(-1)
	t.Run("all ok", func(t *testing.T) {
		if got != 1 {
			t.Errorf("Abs(-1) = %d; want 1", got)
		}
	})

	t.Run("本当に?", func(t *testing.T) {
		if got != 1 {
			t.Errorf("Abs(-1) = %d; want 1", got)
		}
	})
}

func TestReverseAbss(t *testing.T) {
	got := pkg03.ReverseAbs(-1)
	if got != -1 {
		t.Errorf("Abs(-1) = %d; want 1", got)
	}
}

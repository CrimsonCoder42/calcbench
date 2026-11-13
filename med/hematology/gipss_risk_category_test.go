package hematology

import (
	"errors"
	"testing"
)

// Tefferi 2018: 0 Low, 1 Intermediate-1, 2 Intermediate-2, ≥ 3 High.
func TestGIPSSRiskCategory_example(t *testing.T) {
	want := []GIPSSCategory{GIPSSLow, GIPSSIntermediate1, GIPSSIntermediate2, GIPSSHigh, GIPSSHigh, GIPSSHigh, GIPSSHigh}
	for s, w := range want {
		if got, err := GIPSSRiskCategory(s); err != nil || got != w {
			t.Errorf("GIPSSRiskCategory(%d) = (%q, %v), want %q", s, got, err, w)
		}
	}
	for _, s := range []int{-1, 7} {
		if _, err := GIPSSRiskCategory(s); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("GIPSSRiskCategory(%d) error = %v, want ErrOutOfRange", s, err)
		}
	}
}

func TestGIPSSRiskCategory_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { g, _ := GIPSSRiskCategory(2); sinkString = string(g) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkGIPSSRiskCategory(b *testing.B) {
	b.ReportAllocs()
	s := 2
	for b.Loop() {
		g, _ := GIPSSRiskCategory(s)
		sinkString = string(g)
	}
}

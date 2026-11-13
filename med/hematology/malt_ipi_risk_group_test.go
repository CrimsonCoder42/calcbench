package hematology

import (
	"errors"
	"testing"
)

// Thieblemont 2017 abstract: 0 low, 1 intermediate, ≥ 2 high.
func TestMALTIPIRiskGroup_example(t *testing.T) {
	want := []MALTIPIGroup{MALTIPILow, MALTIPIIntermediate, MALTIPIHigh, MALTIPIHigh}
	for s, w := range want {
		if got, err := MALTIPIRiskGroup(s); err != nil || got != w {
			t.Errorf("MALTIPIRiskGroup(%d) = (%q, %v), want %q", s, got, err, w)
		}
	}
	for _, s := range []int{-1, 4} {
		if _, err := MALTIPIRiskGroup(s); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("MALTIPIRiskGroup(%d) error = %v, want ErrOutOfRange", s, err)
		}
	}
}

func TestMALTIPIRiskGroup_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { g, _ := MALTIPIRiskGroup(1); sinkString = string(g) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkMALTIPIRiskGroup(b *testing.B) {
	b.ReportAllocs()
	s := 1
	for b.Loop() {
		g, _ := MALTIPIRiskGroup(s)
		sinkString = string(g)
	}
}

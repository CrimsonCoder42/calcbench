package hematology

import (
	"errors"
	"testing"
)

// Zhou 2014 Table 3 groups, each edge on both sides.
func TestNCCNIPIRiskGroup_example(t *testing.T) {
	want := []NCCNIPIGroup{NCCNIPILow, NCCNIPILow, NCCNIPILowIntermediate, NCCNIPILowIntermediate,
		NCCNIPIHighIntermediate, NCCNIPIHighIntermediate, NCCNIPIHigh, NCCNIPIHigh, NCCNIPIHigh}
	for s, w := range want {
		if got, err := NCCNIPIRiskGroup(s); err != nil || got != w {
			t.Errorf("NCCNIPIRiskGroup(%d) = (%q, %v), want %q", s, got, err, w)
		}
	}
	for _, s := range []int{-1, 9} {
		if _, err := NCCNIPIRiskGroup(s); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("NCCNIPIRiskGroup(%d) error = %v, want ErrOutOfRange", s, err)
		}
	}
}

func TestNCCNIPIRiskGroup_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { g, _ := NCCNIPIRiskGroup(5); sinkString = string(g) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkNCCNIPIRiskGroup(b *testing.B) {
	b.ReportAllocs()
	s := 5
	for b.Loop() {
		g, _ := NCCNIPIRiskGroup(s)
		sinkString = string(g)
	}
}

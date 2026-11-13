package neuro

import (
	"errors"
	"testing"
)

// Tokuhashi 2005 abstract groups, each edge on both sides.
func TestTokuhashiRevisedGroup_example(t *testing.T) {
	tests := []struct {
		total int
		want  TokuhashiGroup
	}{
		{0, TokuhashiUnder6Months}, {8, TokuhashiUnder6Months}, {9, Tokuhashi6MonthsOrMore},
		{11, Tokuhashi6MonthsOrMore}, {12, Tokuhashi1YearOrMore}, {15, Tokuhashi1YearOrMore},
	}
	for _, tc := range tests {
		if got, err := TokuhashiRevisedGroup(tc.total); err != nil || got != tc.want {
			t.Errorf("TokuhashiRevisedGroup(%d) = (%q, %v), want %q", tc.total, got, err, tc.want)
		}
	}
	for _, bad := range []int{-1, 16} {
		if _, err := TokuhashiRevisedGroup(bad); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("TokuhashiRevisedGroup(%d) error = %v, want ErrOutOfRange", bad, err)
		}
	}
}

func TestTokuhashiRevisedGroup_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { g, _ := TokuhashiRevisedGroup(10); sinkString = string(g) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkTokuhashiRevisedGroup(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		g, _ := TokuhashiRevisedGroup(10)
		sinkString = string(g)
	}
}

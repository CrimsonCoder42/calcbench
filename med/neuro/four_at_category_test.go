package neuro

import (
	"errors"
	"testing"
)

// Bellelli 2014 bands (0; 1–3; ≥ 4), each edge on both sides.
func TestFourATCategoryOf_example(t *testing.T) {
	tests := []struct {
		total int
		want  FourATCategory
	}{
		{0, FourATUnlikely}, {1, FourATPossibleCognitiveImpairment}, {3, FourATPossibleCognitiveImpairment},
		{4, FourATPossibleDelirium}, {12, FourATPossibleDelirium},
	}
	for _, tc := range tests {
		if got, err := FourATCategoryOf(tc.total); err != nil || got != tc.want {
			t.Errorf("FourATCategoryOf(%d) = (%q, %v), want %q", tc.total, got, err, tc.want)
		}
	}
	for _, bad := range []int{-1, 13} {
		if _, err := FourATCategoryOf(bad); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("FourATCategoryOf(%d) error = %v, want ErrOutOfRange", bad, err)
		}
	}
}

func TestFourATCategoryOf_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { c, _ := FourATCategoryOf(4); sinkString = string(c) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkFourATCategoryOf(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		c, _ := FourATCategoryOf(4)
		sinkString = string(c)
	}
}

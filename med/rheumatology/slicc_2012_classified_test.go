package rheumatology

import (
	"errors"
	"testing"
)

// Petri 2012 rule: ≥ 4 criteria with ≥ 1 clinical and ≥ 1 immunologic, or biopsy nephritis with ANA/anti-dsDNA.
func TestSLICC2012Classified_example(t *testing.T) {
	tests := []struct {
		c, i   int
		biopsy bool
		want   bool
	}{
		{3, 1, false, true}, {1, 3, false, true}, {2, 1, false, false}, {4, 0, false, false}, {0, 4, false, false},
		{0, 6, false, false}, {11, 6, false, true}, {0, 0, true, true}, {1, 1, false, false},
	}
	for _, tc := range tests {
		if got, err := SLICC2012Classified(tc.c, tc.i, tc.biopsy); err != nil || got != tc.want {
			t.Errorf("SLICC2012Classified(%d, %d, %v) = (%v, %v), want %v", tc.c, tc.i, tc.biopsy, got, err, tc.want)
		}
	}
}

func TestSLICC2012Classified_errors(t *testing.T) {
	for _, tc := range [][2]int{{-1, 0}, {12, 0}, {0, -1}, {0, 7}} {
		if _, err := SLICC2012Classified(tc[0], tc[1], false); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("SLICC2012Classified(%d, %d) error = %v, want ErrOutOfRange", tc[0], tc[1], err)
		}
	}
}

func TestSLICC2012Classified_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkBool, _ = SLICC2012Classified(3, 1, false) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkSLICC2012Classified(b *testing.B) {
	b.ReportAllocs()
	c, i := 3, 1
	for b.Loop() {
		sinkBool, _ = SLICC2012Classified(c, i, false)
	}
}

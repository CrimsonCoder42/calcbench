package pediatrics

import (
	"errors"
	"math"
	"testing"
)

var (
	sinkInt    int
	sinkFloat  float64
	sinkString string
)

// Kocher 1999 abstract: fever history, non-weight-bearing, ESR "at least forty" mm/h, WBC "more than 12,000"
// cells/mm³. ESR 40 counts, 39.9 does not; WBC 12,000 does not, 12,001 does.
func TestKocherCriteriaCount_example(t *testing.T) {
	tests := []struct {
		fever, nwb bool
		esr, wbc   float64
		want       int
	}{
		{false, false, 39.9, 12000, 0}, {false, false, 40, 12000, 1}, {false, false, 39.9, 12001, 1},
		{true, false, 10, 8000, 1}, {false, true, 10, 8000, 1}, {true, true, 40, 12001, 4}, {true, true, 20, 9000, 2},
	}
	for _, tc := range tests {
		if got, err := KocherCriteriaCount(tc.fever, tc.nwb, tc.esr, tc.wbc); err != nil || got != tc.want {
			t.Errorf("KocherCriteriaCount(%+v) = (%d, %v), want %d", tc, got, err, tc.want)
		}
	}
}

func TestKocherCriteriaCount_errors(t *testing.T) {
	cases := []struct {
		esr, wbc float64
		want     error
	}{
		{math.NaN(), 1000, ErrNotFinite}, {20, math.Inf(1), ErrNotFinite}, {math.Inf(-1), 1000, ErrNotFinite},
		{-1, 1000, ErrOutOfRange}, {201, 1000, ErrOutOfRange}, {20, -1, ErrOutOfRange}, {20, 500001, ErrOutOfRange},
	}
	for _, c := range cases {
		if _, err := KocherCriteriaCount(false, false, c.esr, c.wbc); !errors.Is(err, c.want) {
			t.Errorf("KocherCriteriaCount(esr %v, wbc %v) error = %v, want %v", c.esr, c.wbc, err, c.want)
		}
	}
	for _, c := range [][2]float64{{0, 0}, {200, 500000}} {
		if _, err := KocherCriteriaCount(false, false, c[0], c[1]); err != nil {
			t.Errorf("KocherCriteriaCount(%v, %v) error = %v, want nil (range edge)", c[0], c[1], err)
		}
	}
}

func TestKocherCriteriaCount_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = KocherCriteriaCount(true, false, 45, 13000) }); n != 0 {
		t.Errorf("allocates %v times, want 0", n)
	}
}

func BenchmarkKocherCriteriaCount(b *testing.B) {
	b.ReportAllocs()
	esr, wbc := 45.0, 13000.0
	for b.Loop() {
		sinkInt, _ = KocherCriteriaCount(true, false, esr, wbc)
	}
}

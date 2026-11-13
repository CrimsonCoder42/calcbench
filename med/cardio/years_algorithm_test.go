package cardio

import (
	"math"
	"testing"
)

// van der Hulle 2017 abstract: 0 items → excluded below 1000 ng/mL; ≥ 1 item → excluded below 500 ng/mL.
func TestYEARSPulmonaryEmbolismExcluded(t *testing.T) {
	tests := []struct {
		items int
		dd    float64
		want  bool
	}{
		{0, 999.9, true}, {0, 1000, false}, {0, 0, true},
		{1, 499.9, true}, {1, 500, false}, {3, 499, true}, {3, 900, false},
		{0, 900, true}, {1, 600, false}, // doc Example
	}
	for _, tc := range tests {
		if got, err := YEARSPulmonaryEmbolismExcluded(tc.items, tc.dd); err != nil || got != tc.want {
			t.Errorf("(%d, %v) = (%v, %v), want %v", tc.items, tc.dd, got, err, tc.want)
		}
	}
}

// parity: ledger years_algorithm — no source test file or JSON example; source case 1 is the abstract's rule.
func TestYEARSPulmonaryEmbolismExcluded_parity(t *testing.T) {
	if got, _ := YEARSPulmonaryEmbolismExcluded(2, 450); !got { // parity: ledger years_algorithm, source case 1
		t.Error("2 items, 450 ng/mL: want true")
	}
}

func TestYEARSPulmonaryEmbolismExcluded_errors(t *testing.T) {
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		_, err := YEARSPulmonaryEmbolismExcluded(0, v)
		checkErr(t, "non-finite", err, ErrNotFinite)
	}
	for _, c := range []struct {
		items int
		dd    float64
	}{{-1, 100}, {4, 100}, {0, -0.1}, {0, 100000.1}} {
		_, err := YEARSPulmonaryEmbolismExcluded(c.items, c.dd)
		checkErr(t, "range", err, ErrOutOfRange)
	}
}

func TestYEARSPulmonaryEmbolismExcluded_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkBool, _ = YEARSPulmonaryEmbolismExcluded(1, 400) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkYEARSPulmonaryEmbolismExcluded(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkBool, _ = YEARSPulmonaryEmbolismExcluded(1, 400)
	}
}

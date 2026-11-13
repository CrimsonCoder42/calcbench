package critical

import (
	"errors"
	"math"
	"testing"
)

var (
	sinkB6Int   int
	sinkB6Bool  bool
	sinkB6Float float64
)

var nonFiniteB6 = []float64{math.NaN(), math.Inf(1), math.Inf(-1)}

// wantErrB6 fails unless err wraps want.
func wantErrB6(t *testing.T, name string, err, want error) {
	t.Helper()
	if !errors.Is(err, want) {
		t.Errorf("%s: err = %v, want %v", name, err, want)
	}
}

// Miller 2013 table (PMC3618004): each item's bands, both sides of each printed edge. HAND-DERIVED.
func TestATLASScore_bands(t *testing.T) {
	type in struct {
		age      float64
		abx      bool
		wbc, alb float64
		cr       float64
		want     int
	}
	tests := []in{
		{59.9, false, 15999, 36, 120, 0},
		{60, false, 15999, 36, 120, 1}, {79.9, false, 15999, 36, 120, 1}, {80, false, 15999, 36, 120, 2},
		{40, true, 15999, 36, 120, 2},
		{40, false, 16000, 36, 120, 1}, {40, false, 25000, 36, 120, 1}, {40, false, 25001, 36, 120, 2},
		{40, false, 15999, 35.1, 120, 0}, {40, false, 15999, 35, 120, 1}, {40, false, 15999, 26, 120, 1},
		{40, false, 15999, 25.5, 120, 1}, {40, false, 15999, 25, 120, 2},
		{40, false, 15999, 36, 120.5, 1}, {40, false, 15999, 36, 121, 1}, {40, false, 15999, 36, 179, 1},
		{40, false, 15999, 36, 180, 2},
		{82, true, 26000, 25, 180, 10},
	}
	for i, tc := range tests {
		if got, err := ATLASScore(tc.age, tc.abx, tc.wbc, tc.alb, tc.cr); err != nil || got != tc.want {
			t.Errorf("case %d %+v: (%d, %v), want %d", i, tc, got, err, tc.want)
		}
	}
}

// parity: ledger atlas_score, source case 1. The app sums category pickers with antibiotics = 2 points, as the
// table does; age_cat 2, abx, wbc_cat 1, albumin_cat 1, creatinine_cat 0 → 2+2+1+1+0 = 6.
func TestATLASScore_parity(t *testing.T) {
	if got, _ := ATLASScore(85, true, 20000, 30, 100); got != 6 {
		t.Errorf("got %d, want 6", got)
	}
}

func TestATLASScore_errors(t *testing.T) {
	for _, v := range nonFiniteB6 {
		for i := 0; i < 4; i++ {
			a := [4]float64{40, 10000, 40, 80}
			a[i] = v
			_, err := ATLASScore(a[0], false, a[1], a[2], a[3])
			wantErrB6(t, "non-finite", err, ErrNotFinite)
		}
	}
	for _, a := range [][4]float64{{-1, 10000, 40, 80}, {131, 10000, 40, 80}, {40, -1, 40, 80}, {40, 10000, -1, 80}, {40, 10000, 40, -1}} {
		_, err := ATLASScore(a[0], false, a[1], a[2], a[3])
		wantErrB6(t, "range", err, ErrOutOfRange)
	}
}

func TestATLASScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkB6Int, _ = ATLASScore(70, true, 20000, 30, 150) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkATLASScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkB6Int, _ = ATLASScore(70, true, 20000, 30, 150)
	}
}

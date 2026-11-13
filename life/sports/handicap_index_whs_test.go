package sports

import (
	"errors"
	"math"
	"testing"
)

func seq(from, to float64) []float64 {
	var s []float64
	for v := from; v <= to+1e-9; v++ {
		s = append(s, v)
	}
	return s
}

// Cases HAND-DERIVED from Rules of Handicapping Rule 5.2a's table (read 2026-11-13), one per table row:
// differentials 10, 11, 12, … (ascending), so "lowest k" is the first k.
func TestHandicapIndexWHS_table52a(t *testing.T) {
	tests := []struct {
		n    int
		want float64
	}{
		{3, 8.0},   // lowest 1 = 10.0, −2.0
		{4, 9.0},   // lowest 1, −1.0
		{5, 10.0},  // lowest 1
		{6, 9.5},   // mean(10, 11) = 10.5, −1.0
		{7, 10.5},  // lowest 2
		{8, 10.5},  // lowest 2
		{9, 11.0},  // lowest 3: 10, 11, 12
		{11, 11.0}, // lowest 3
		{12, 11.5}, // lowest 4
		{14, 11.5}, // lowest 4
		{15, 12.0}, // lowest 5
		{16, 12.0}, // lowest 5
		{17, 12.5}, // lowest 6
		{18, 12.5}, // lowest 6
		{19, 13.0}, // lowest 7
		{20, 13.5}, // lowest 8: mean 10..17
	}
	for _, tc := range tests {
		d := seq(10, 10+float64(tc.n-1))
		if got, err := HandicapIndexWHS(d); err != nil || got != tc.want {
			t.Errorf("n=%d: (%v, %v), want %v", tc.n, got, err, tc.want)
		}
	}
}

// Only the most recent 20 count: 25 differentials whose first 5 (oldest) are 0.0 give the same result as the last 20.
// The lowest are found regardless of order, and the input is not modified.
func TestHandicapIndexWHS_recentAndOrder(t *testing.T) {
	d := append([]float64{0, 0, 0, 0, 0}, seq(1, 20)...)
	if got, err := HandicapIndexWHS(d); err != nil || got != 4.5 {
		t.Errorf("25 scores: (%v, %v), want 4.5 (mean of 1..8)", got, err)
	}
	r := []float64{20, 3, 17, 1, 9, 12, 5, 8, 2, 15, 4, 19, 7, 11, 6, 14, 10, 18, 13, 16}
	orig := append([]float64(nil), r...)
	if got, err := HandicapIndexWHS(r); err != nil || got != 4.5 {
		t.Errorf("shuffled 1..20: (%v, %v), want 4.5", got, err)
	}
	for i := range r {
		if r[i] != orig[i] {
			t.Fatalf("input modified at %d", i)
		}
	}
}

// Rounding to the nearest tenth, ties half to even (house default; the Rules give no tie rule):
// 7 scores, lowest 2 = 10.1 and 10.2 → mean 10.15 → 10.2 (even tenth 102); 10.2 and 10.3 → 10.25 → 10.2.
// Negative: −1.1 and −1.2 → −1.15 → −1.2. Lowest 3 = 10.0, 10.0, 10.1 → 10.0333 → 10.0.
func TestHandicapIndexWHS_rounding(t *testing.T) {
	pad := []float64{30, 30, 30, 30, 30}
	tests := []struct {
		low  []float64
		want float64
	}{
		{[]float64{10.1, 10.2}, 10.2},
		{[]float64{10.2, 10.3}, 10.2},
		{[]float64{-1.1, -1.2}, -1.2},
	}
	for _, tc := range tests {
		d := append(append([]float64(nil), pad...), tc.low...)
		if got, err := HandicapIndexWHS(d); err != nil || got != tc.want {
			t.Errorf("%v: (%v, %v), want %v", tc.low, got, err, tc.want)
		}
	}
	d := []float64{10.0, 10.0, 10.1, 30, 30, 30, 30, 30, 30}
	if got, err := HandicapIndexWHS(d); err != nil || got != 10.0 {
		t.Errorf("thirds: (%v, %v), want 10.0", got, err)
	}
}

// Rule 5.3: the result is at most 54.0.
func TestHandicapIndexWHS_cap(t *testing.T) {
	if got, err := HandicapIndexWHS([]float64{60, 61, 62, 63, 64}); err != nil || got != 54.0 {
		t.Errorf("(%v, %v), want 54.0", got, err)
	}
}

func TestHandicapIndexWHS_errors(t *testing.T) {
	cases := []struct {
		name string
		d    []float64
		want error
	}{
		{"nil", nil, ErrOutOfRange},
		{"two", []float64{1, 2}, ErrOutOfRange},
		{"NaN", []float64{1, math.NaN(), 3}, ErrNotFinite},
		{"+Inf", []float64{1, 2, math.Inf(1)}, ErrNotFinite},
		{"too large", []float64{1, 2, 1e300}, ErrOutOfRange},
	}
	for _, tc := range cases {
		if _, err := HandicapIndexWHS(tc.d); !errors.Is(err, tc.want) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.want)
		}
	}
}

func TestHandicapIndexWHS_noAllocation(t *testing.T) {
	d := seq(1, 20)
	var s float64
	if n := testing.AllocsPerRun(100, func() { s, _ = HandicapIndexWHS(d) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
	_ = s
}

func BenchmarkHandicapIndexWHS(b *testing.B) {
	for _, n := range []int{20, 1000} {
		d := seq(1, float64(n))
		b.Run(map[int]string{20: "20", 1000: "1000"}[n], func(b *testing.B) {
			b.ReportAllocs()
			var s float64
			for b.Loop() {
				s, _ = HandicapIndexWHS(d)
			}
			_ = s
		})
	}
}

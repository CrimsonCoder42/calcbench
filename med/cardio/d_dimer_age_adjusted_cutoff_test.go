package cardio

import (
	"math"
	"testing"
)

// Schouten 2013: for a person aged 78 years, "the D-dimer concentration would be considered normal below 780 μg/L".
// Douma 2010: age × 10 above 50, 500 otherwise; edge 50 on both sides. Tolerance 1e-9 absolute.
func TestDDimerAgeAdjustedCutoffUgL_values(t *testing.T) {
	tests := []struct{ age, want float64 }{{78, 780}, {40, 500}, {50, 500}, {50.5, 505}, {120, 1200}, {18, 500}}
	for _, tc := range tests {
		if got, err := DDimerAgeAdjustedCutoffUgL(tc.age); err != nil || math.Abs(got-tc.want) > 1e-9 {
			t.Errorf("age %v: (%v, %v), want %v", tc.age, got, err, tc.want)
		}
	}
	_, err := DDimerAgeAdjustedCutoffUgL(math.NaN())
	checkErr(t, "NaN", err, ErrNotFinite)
	_, err = DDimerAgeAdjustedCutoffUgL(math.Inf(1))
	checkErr(t, "+Inf", err, ErrNotFinite)
	_, err = DDimerAgeAdjustedCutoffUgL(120.1)
	checkErr(t, "120.1", err, ErrOutOfRange)
}

func TestDDimerAgeAdjustedCutoffUgL_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = DDimerAgeAdjustedCutoffUgL(78) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkDDimerAgeAdjustedCutoffUgL(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = DDimerAgeAdjustedCutoffUgL(78)
	}
}

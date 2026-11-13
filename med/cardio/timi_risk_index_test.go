package cardio

import (
	"math"
	"testing"
)

// HAND-DERIVED from Morrow 2001's formula (heart rate × (age/10)²) / SBP. Tolerance 1e-12 relative (a few float ops).
func TestTIMIRiskIndex_values(t *testing.T) {
	tests := []struct{ hr, age, sbp, want float64 }{
		{80, 60, 120, 24},
		{100, 18, 100, 3.24},
		{60, 120, 200, 43.2},
	}
	for _, tc := range tests {
		got, err := TIMIRiskIndex(tc.hr, tc.age, tc.sbp)
		if err != nil || math.Abs(got-tc.want) > 1e-12*tc.want {
			t.Errorf("TIMIRiskIndex(%v, %v, %v) = (%v, %v), want %v", tc.hr, tc.age, tc.sbp, got, err, tc.want)
		}
	}
}

func TestTIMIRiskIndex_errors(t *testing.T) {
	nan, inf := math.NaN(), math.Inf(1)
	cases := []struct {
		name         string
		hr, age, sbp float64
		want         error
	}{
		{"hr NaN", nan, 60, 120, ErrNotFinite},
		{"age +Inf", 80, inf, 120, ErrNotFinite},
		{"sbp -Inf", 80, 60, -inf, ErrNotFinite},
		{"hr 0", 0, 60, 120, ErrOutOfRange},
		{"sbp 0", 80, 60, 0, ErrOutOfRange},
		{"age 17.9", 80, 17.9, 120, ErrOutOfRange},
		{"age 120.1", 80, 120.1, 120, ErrOutOfRange},
	}
	for _, tc := range cases {
		_, err := TIMIRiskIndex(tc.hr, tc.age, tc.sbp)
		checkErr(t, tc.name, err, tc.want)
	}
}

func TestTIMIRiskIndex_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = TIMIRiskIndex(80, 60, 120) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkTIMIRiskIndex(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = TIMIRiskIndex(80, 60, 120)
	}
}

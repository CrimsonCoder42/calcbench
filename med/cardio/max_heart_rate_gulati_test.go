package cardio

import (
	"math"
	"testing"
)

// parity: ledger max_hr, source cases 1–6, "gulati" output (printed to 0.1; tolerance 0.05).
func TestMaxHeartRateGulati_parity(t *testing.T) {
	tests := []struct{ age, want float64 }{
		{35, 175.2}, {40, 170.8}, {40, 170.8}, {65, 148.8}, {20, 188.4}, {90, 126.8},
	}
	for i, tc := range tests {
		if got, err := MaxHeartRateGulati(tc.age); err != nil || math.Abs(got-tc.want) > 0.05 {
			t.Errorf("source case %d: (%v, %v), want %v", i+1, got, err, tc.want)
		}
	}
}

func TestMaxHeartRateGulati_edges(t *testing.T) {
	// 206 − 0.88 × 18 = 190.16; 206 − 0.88 × 100 = 118. Tolerance 1e-12 absolute.
	for _, tc := range []struct{ age, want float64 }{{18, 190.16}, {100, 118}} {
		if got, err := MaxHeartRateGulati(tc.age); err != nil || math.Abs(got-tc.want) > 1e-12 {
			t.Errorf("age %v: (%v, %v), want %v", tc.age, got, err, tc.want)
		}
	}
	for _, a := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		_, err := MaxHeartRateGulati(a)
		checkErr(t, "non-finite", err, ErrNotFinite)
	}
	for _, a := range []float64{17.99, 100.01} {
		_, err := MaxHeartRateGulati(a)
		checkErr(t, "range", err, ErrOutOfRange)
	}
}

func TestMaxHeartRateGulati_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = MaxHeartRateGulati(40) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkMaxHeartRateGulati(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = MaxHeartRateGulati(40)
	}
}

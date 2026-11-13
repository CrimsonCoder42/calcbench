package cardio

import (
	"math"
	"testing"
)

// parity: ledger max_hr, source cases 1–6, "hunt" output (printed to 0.1; tolerance 0.05).
func TestMaxHeartRateNes_parity(t *testing.T) {
	tests := []struct{ age, want float64 }{
		{35, 188.6}, {40, 185.4}, {40, 185.4}, {65, 169.4}, {20, 198.2}, {90, 153.4},
	}
	for i, tc := range tests {
		if got, err := MaxHeartRateNes(tc.age); err != nil || math.Abs(got-tc.want) > 0.05 {
			t.Errorf("source case %d: (%v, %v), want %v", i+1, got, err, tc.want)
		}
	}
}

func TestMaxHeartRateNes_edges(t *testing.T) {
	// 211 − 0.64 × 18 = 199.48; 211 − 64 = 147. Tolerance 1e-12 absolute.
	for _, tc := range []struct{ age, want float64 }{{18, 199.48}, {100, 147}} {
		if got, err := MaxHeartRateNes(tc.age); err != nil || math.Abs(got-tc.want) > 1e-12 {
			t.Errorf("age %v: (%v, %v), want %v", tc.age, got, err, tc.want)
		}
	}
	for _, a := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		_, err := MaxHeartRateNes(a)
		checkErr(t, "non-finite", err, ErrNotFinite)
	}
	for _, a := range []float64{17.99, 100.01} {
		_, err := MaxHeartRateNes(a)
		checkErr(t, "range", err, ErrOutOfRange)
	}
}

func TestMaxHeartRateNes_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = MaxHeartRateNes(40) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkMaxHeartRateNes(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = MaxHeartRateNes(40)
	}
}

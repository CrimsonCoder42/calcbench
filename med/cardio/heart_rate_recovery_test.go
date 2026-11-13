package cardio

import (
	"math"
	"testing"
)

// parity: ledger hrr1, source cases 1–7 (case 7, one-minute rate above peak, is refused by the source too).
func TestHeartRateRecovery1Min_parity(t *testing.T) {
	tests := []struct{ peak, one, want float64 }{
		{180, 150, 30}, // source case 1
		{170, 148, 22}, // source case 2
		{160, 142, 18}, // source case 3
		{150, 135, 15}, // source case 4
		{150, 137, 13}, // source case 5
		{140, 130, 10}, // source case 6
	}
	for i, tc := range tests {
		if got, err := HeartRateRecovery1Min(tc.peak, tc.one); err != nil || got != tc.want {
			t.Errorf("case %d: (%v, %v), want %v", i+1, got, err, tc.want)
		}
	}
	_, err := HeartRateRecovery1Min(100, 120) // source case 7
	checkErr(t, "case 7", err, ErrDomain)
}

func TestHeartRateRecovery1Min_edges(t *testing.T) {
	if got, err := HeartRateRecovery1Min(120, 120); err != nil || got != 0 {
		t.Errorf("equal rates: (%v, %v), want 0", got, err)
	}
	_, err := HeartRateRecovery1Min(math.NaN(), 100)
	checkErr(t, "NaN", err, ErrNotFinite)
	_, err = HeartRateRecovery1Min(150, math.Inf(1))
	checkErr(t, "+Inf", err, ErrNotFinite)
	_, err = HeartRateRecovery1Min(0, 0)
	checkErr(t, "zero", err, ErrOutOfRange)
	_, err = HeartRateRecovery1Min(150, -1)
	checkErr(t, "negative", err, ErrOutOfRange)
}

func TestHeartRateRecovery1Min_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = HeartRateRecovery1Min(150, 135) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkHeartRateRecovery1Min(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = HeartRateRecovery1Min(150, 135)
	}
}

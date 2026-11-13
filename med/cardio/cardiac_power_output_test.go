package cardio

import (
	"math"
	"testing"
)

// HAND-DERIVED from the printed definition CPO = MAP × CO / 451 (Fincke 2004). Tolerance 1e-12 relative.
func TestCardiacPowerOutputW_values(t *testing.T) {
	tests := []struct {
		call func() (float64, error)
		want float64
	}{
		{func() (float64, error) { return CardiacPowerOutputW(70, 4) }, 280.0 / 451},
		{func() (float64, error) { return CardiacPowerOutputW(100, 5) }, 500.0 / 451},
	}
	for i, tc := range tests {
		got, err := tc.call()
		if err != nil || math.Abs(got-tc.want) > 1e-12*math.Abs(tc.want) {
			t.Errorf("case %d: (%v, %v), want %v", i, got, err, tc.want)
		}
	}
}

func TestCardiacPowerOutputW_errors(t *testing.T) {
	var err error
	_, err = CardiacPowerOutputW(math.NaN(), 4)
	checkErr(t, "math.NaN(), 4", err, ErrNotFinite)
	_, err = CardiacPowerOutputW(70, math.Inf(-1))
	checkErr(t, "70, math.Inf(-1)", err, ErrNotFinite)
	_, err = CardiacPowerOutputW(0, 4)
	checkErr(t, "0, 4", err, ErrOutOfRange)
	_, err = CardiacPowerOutputW(70, -1)
	checkErr(t, "70, -1", err, ErrOutOfRange)
}

func TestCardiacPowerOutputW_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = CardiacPowerOutputW(70, 4) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkCardiacPowerOutputW(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = CardiacPowerOutputW(70, 4)
	}
}

package electrical

import (
	"errors"
	"math"
	"testing"
)

// Hand computations: 10 A × 0.0012 Ω/m × 2 × 152.4 m = 3.6576 V; three-phase √3 × 10 × 0.0012 × 152.4 =
// 3.167574... V (√3 = 1.7320508075688772). Tolerance 1e-12 relative.
func TestVoltageDropTwoWire_cases(t *testing.T) {
	got, err := VoltageDropTwoWire(10, 0.0012, 152.4)
	if err != nil || !near(got, 3.6576) {
		t.Errorf("got %v %v", got, err)
	}
	if got, err := VoltageDropTwoWire(0, 0.0012, 152.4); err != nil || got != 0 {
		t.Errorf("zero current: %v %v", got, err)
	}
	for _, c := range []struct {
		i, r, l float64
		want    error
	}{
		{-1, 1, 1, ErrNegative},
		{1, -1, 1, ErrNegative},
		{1, 1, math.NaN(), ErrNotFinite},
		{math.MaxFloat64, 2, 2, ErrOverflow},
	} {
		if _, err := VoltageDropTwoWire(c.i, c.r, c.l); !errors.Is(err, c.want) {
			t.Errorf("(%v, %v, %v): %v, want %v", c.i, c.r, c.l, err, c.want)
		}
	}
}

func TestVoltageDropThreePhase_cases(t *testing.T) {
	got, err := VoltageDropThreePhase(10, 0.0012, 152.4)
	if err != nil || !near(got, math.Sqrt(3)*1.8288) {
		t.Errorf("got %v %v", got, err)
	}
	if _, err := VoltageDropThreePhase(1, math.Inf(1), 1); !errors.Is(err, ErrNotFinite) {
		t.Errorf("Inf: %v", err)
	}
}

func BenchmarkVoltageDropTwoWire(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = VoltageDropTwoWire(10, 0.0012, 152.4)
	}
}

func BenchmarkVoltageDropThreePhase(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = VoltageDropThreePhase(10, 0.0012, 152.4)
	}
}

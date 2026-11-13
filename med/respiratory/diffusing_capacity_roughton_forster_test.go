package respiratory

import (
	"errors"
	"math"
	"testing"
)

// Hand-derived from 1/DL = 1/DM + 1/(θ·VC): DM 50, θ·VC = 1.0 × 75 = 75 → 1/DL = 1/50 + 1/75 = 5/150, DL = 30.
func TestDiffusingCapacityRoughtonForster_handExample(t *testing.T) {
	got, err := DiffusingCapacityRoughtonForster(50, 1.0, 75)
	if err != nil || math.Abs(got-30) > 1e-12 {
		t.Errorf("got (%v, %v), want (30, nil)", got, err)
	}
}

// Oracle: the reciprocal form written out, on a grid of inputs; agreement within 4 ulps-ish (1e-13 relative).
func TestDiffusingCapacityRoughtonForster_reciprocalOracle(t *testing.T) {
	for _, dm := range []float64{5, 20, 50, 120} {
		for _, theta := range []float64{0.3, 0.6, 1.2} {
			for _, vc := range []float64{20, 60, 150} {
				want := 1 / (1/dm + 1/(theta*vc))
				got, err := DiffusingCapacityRoughtonForster(dm, theta, vc)
				if err != nil || math.Abs(got-want) > 1e-13*want {
					t.Errorf("(%v, %v, %v) = (%v, %v), want %v", dm, theta, vc, got, err, want)
				}
			}
		}
	}
}

// The result is below both components (a series conductance).
func TestDiffusingCapacityRoughtonForster_belowBothComponents(t *testing.T) {
	got, _ := DiffusingCapacityRoughtonForster(40, 0.8, 90)
	if !(got < 40 && got < 0.8*90) {
		t.Errorf("DL = %v, want < DM 40 and < θVC 72", got)
	}
}

func TestDiffusingCapacityRoughtonForster_errors(t *testing.T) {
	tests := []struct {
		name          string
		dm, theta, vc float64
		wantErr       error
	}{
		{"zero DM", 0, 1, 75, ErrOutOfRange},
		{"negative theta", 50, -1, 75, ErrOutOfRange},
		{"zero VC", 50, 1, 0, ErrOutOfRange},
		{"NaN DM", math.NaN(), 1, 75, ErrNotFinite},
		{"+Inf theta", 50, math.Inf(1), 75, ErrNotFinite},
		{"-Inf VC", 50, 1, math.Inf(-1), ErrNotFinite},
		{"overflowing θ·VC", 50, math.MaxFloat64, 2, ErrDomain},
	}
	for _, tc := range tests {
		_, err := DiffusingCapacityRoughtonForster(tc.dm, tc.theta, tc.vc)
		if !errors.Is(err, tc.wantErr) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.wantErr)
		}
	}
}

var sinkFloat float64

func TestDiffusingCapacityRoughtonForster_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = DiffusingCapacityRoughtonForster(50, 1, 75) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkDiffusingCapacityRoughtonForster(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = DiffusingCapacityRoughtonForster(50, 1, 75)
	}
}

package fitness

import (
	"errors"
	"math"
	"testing"
)

// Worked examples, HAND-DERIVED from Herrmann et al. 2024 §2.5 (1 MET ≈ 1.0 kcal/kg/h):
// 6.0 METs (the paper's code 01084 example value), 70 kg, 0.5 h → 210 kcal; 3.5 METs, 68.0388555 kg, 0.75 h →
// 178.60199568750 kcal; 1 MET, 1 kg, 1 h → 1 kcal; any MET for 0 h → 0 kcal.
// Tolerance 1e-15 relative: two multiplies.
func TestEnergyExpenditureMET_example(t *testing.T) {
	tests := []struct{ met, kg, h, want float64 }{
		{6.0, 70, 0.5, 210},
		{3.5, 68.0388555, 0.75, 178.6019956875},
		{1, 1, 1, 1},
		{8, 80, 0, 0},
	}
	for _, tc := range tests {
		got, err := EnergyExpenditureMET(tc.met, tc.kg, tc.h)
		if err != nil || math.Abs(got-tc.want) > 1e-15*math.Max(math.Abs(tc.want), 1) {
			t.Errorf("EnergyExpenditureMET(%v, %v, %v) = (%v, %v), want %v", tc.met, tc.kg, tc.h, got, err, tc.want)
		}
	}
}

func TestEnergyExpenditureMET_errors(t *testing.T) {
	nan, inf := math.NaN(), math.Inf(1)
	cases := []struct {
		name       string
		met, kg, h float64
		want       error
	}{
		{"NaN met", nan, 70, 1, ErrNotFinite},
		{"+Inf mass", 6, inf, 1, ErrNotFinite},
		{"-Inf hours", 6, 70, -inf, ErrNotFinite},
		{"zero met", 0, 70, 1, ErrOutOfRange},
		{"negative met", -1, 70, 1, ErrOutOfRange},
		{"zero mass", 6, 0, 1, ErrOutOfRange},
		{"negative hours", 6, 70, -0.1, ErrOutOfRange},
		{"overflow", math.MaxFloat64, 70, 1, ErrDomain},
	}
	for _, tc := range cases {
		if _, err := EnergyExpenditureMET(tc.met, tc.kg, tc.h); !errors.Is(err, tc.want) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.want)
		}
	}
}

func TestEnergyExpenditureMET_noAllocation(t *testing.T) {
	var s float64
	if n := testing.AllocsPerRun(100, func() { s, _ = EnergyExpenditureMET(6, 70, 0.5) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
	_ = s
}

func BenchmarkEnergyExpenditureMET(b *testing.B) {
	b.ReportAllocs()
	var s float64
	for b.Loop() {
		s, _ = EnergyExpenditureMET(6, 70, 0.5)
	}
	_ = s
}

package cardio

import (
	"errors"
	"math"
	"testing"
)

// Hand-derived from Korteweg 1878 eq. (10): a₁ 0.001 m, E₁ 4e5 Pa, ρ 1050 kg/m³, R₁ 0.01 m →
// a₁E₁ = 400, 2ρR₁ = 21, √(400/21) = 4.364357804719847 m/s.
func TestPulseWaveVelocityMoensKorteweg_handExample(t *testing.T) {
	got, err := PulseWaveVelocityMoensKorteweg(0.001, 4e5, 1050, 0.01)
	if err != nil || math.Abs(got-4.364357804719847) > 1e-12 {
		t.Errorf("got (%v, %v), want (4.364357804719847, nil)", got, err)
	}
}

// Scaling check from eq. (10): quadrupling E₁ doubles the speed; quadrupling R₁ halves it.
func TestPulseWaveVelocityMoensKorteweg_scaling(t *testing.T) {
	base, _ := PulseWaveVelocityMoensKorteweg(0.001, 4e5, 1050, 0.01)
	e4, _ := PulseWaveVelocityMoensKorteweg(0.001, 16e5, 1050, 0.01)
	r4, _ := PulseWaveVelocityMoensKorteweg(0.001, 4e5, 1050, 0.04)
	if math.Abs(e4-2*base) > 1e-12 || math.Abs(r4-base/2) > 1e-12 {
		t.Errorf("scaling: base %v, 4E %v (want %v), 4R %v (want %v)", base, e4, 2*base, r4, base/2)
	}
}

func TestPulseWaveVelocityMoensKorteweg_errors(t *testing.T) {
	tests := []struct {
		name       string
		a, e, p, r float64
		wantErr    error
	}{
		{"zero thickness", 0, 4e5, 1050, 0.01, ErrOutOfRange},
		{"negative modulus", 0.001, -1, 1050, 0.01, ErrOutOfRange},
		{"zero density", 0.001, 4e5, 0, 0.01, ErrOutOfRange},
		{"negative radius", 0.001, 4e5, 1050, -0.01, ErrOutOfRange},
		{"NaN thickness", math.NaN(), 4e5, 1050, 0.01, ErrNotFinite},
		{"+Inf modulus", 0.001, math.Inf(1), 1050, 0.01, ErrNotFinite},
		{"-Inf radius", 0.001, 4e5, 1050, math.Inf(-1), ErrNotFinite},
		{"overflowing product", math.MaxFloat64, math.MaxFloat64, 1, 1, ErrDomain},
		{"underflow to zero", math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64, 1, 1, ErrDomain},
	}
	for _, tc := range tests {
		_, err := PulseWaveVelocityMoensKorteweg(tc.a, tc.e, tc.p, tc.r)
		if !errors.Is(err, tc.wantErr) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.wantErr)
		}
	}
}

func TestPulseWaveVelocityMoensKorteweg_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = PulseWaveVelocityMoensKorteweg(0.001, 4e5, 1050, 0.01) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkPulseWaveVelocityMoensKorteweg(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = PulseWaveVelocityMoensKorteweg(0.001, 4e5, 1050, 0.01)
	}
}

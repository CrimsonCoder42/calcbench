package physics

import (
	"errors"
	"math"
	"testing"
)

// Hand computations from ρ = m/V (OpenStax University Physics Vol. 1, §14.1). Tolerance 1e-12 relative:
// one multiplication or division.
func TestDensity_cases(t *testing.T) {
	cases := []struct{ m, v, want float64 }{
		{8900, 1, 8900},
		{1, 0.001, 1000},
		{0, 2, 0},
	}
	for _, c := range cases {
		got, err := Density(c.m, c.v)
		if err != nil || math.Abs(got-c.want) > 1e-12*math.Max(1, c.want) {
			t.Errorf("Density(%v, %v) = %v %v, want %v", c.m, c.v, got, err, c.want)
		}
	}
}

func TestMassFromDensity_cases(t *testing.T) {
	if got, err := MassFromDensity(8900, 0.5); err != nil || got != 4450 {
		t.Errorf("got %v %v", got, err)
	}
	if _, err := MassFromDensity(math.MaxFloat64, 2); !errors.Is(err, ErrOverflow) {
		t.Errorf("overflow: %v", err)
	}
}

func TestVolumeFromDensity_cases(t *testing.T) {
	if got, err := VolumeFromDensity(8900, 8900); err != nil || got != 1 {
		t.Errorf("got %v %v", got, err)
	}
	if _, err := VolumeFromDensity(1, 0); !errors.Is(err, ErrZero) {
		t.Errorf("zero: %v", err)
	}
}

func TestDensity_errors(t *testing.T) {
	cases := []struct {
		m, v float64
		want error
	}{
		{1, 0, ErrZero},
		{-1, 1, ErrNegative},
		{1, -1, ErrNegative},
		{math.NaN(), 1, ErrNotFinite},
		{1, math.Inf(1), ErrNotFinite},
		{math.MaxFloat64, 0.5, ErrOverflow},
	}
	for _, c := range cases {
		if _, err := Density(c.m, c.v); !errors.Is(err, c.want) {
			t.Errorf("Density(%v, %v): err %v, want %v", c.m, c.v, err, c.want)
		}
	}
}

func BenchmarkDensity(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Density(8900, 1)
	}
}

func BenchmarkMassFromDensity(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = MassFromDensity(8900, 0.5)
	}
}

func BenchmarkVolumeFromDensity(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = VolumeFromDensity(8900, 8900)
	}
}

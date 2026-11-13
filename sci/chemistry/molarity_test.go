package chemistry

import (
	"errors"
	"math"
	"testing"
)

// Hand computations from M = mol solute / L solution (OpenStax Chemistry 2e, §3.3):
// 0.133 mol / 0.355 L = 0.3746478873239437 M; 100 g NaCl / 58.44 g/mol / 2 L = 0.8555783709787816 M;
// 0.5 M × 2 L × 58.44 g/mol = 58.44 g. Tolerance 1e-12 relative: two or three operations.
func TestMolarity_cases(t *testing.T) {
	got, err := Molarity(0.133, 0.355)
	if err != nil || math.Abs(got-0.3746478873239437) > 1e-12 {
		t.Errorf("got %v %v", got, err)
	}
	for _, c := range []struct {
		n, v float64
		want error
	}{
		{1, 0, ErrZero},
		{-1, 1, ErrNegative},
		{math.NaN(), 1, ErrNotFinite},
		{math.MaxFloat64, 0.5, ErrOverflow},
	} {
		if _, err := Molarity(c.n, c.v); !errors.Is(err, c.want) {
			t.Errorf("Molarity(%v, %v): %v, want %v", c.n, c.v, err, c.want)
		}
	}
}

func TestMolarityFromMass_cases(t *testing.T) {
	got, err := MolarityFromMass(100, 58.44, 2)
	if err != nil || math.Abs(got-0.8555783709787816) > 1e-12 {
		t.Errorf("got %v %v", got, err)
	}
	if _, err := MolarityFromMass(100, 0, 2); !errors.Is(err, ErrZero) {
		t.Errorf("zero molar mass: %v", err)
	}
	if _, err := MolarityFromMass(100, 58.44, 0); !errors.Is(err, ErrZero) {
		t.Errorf("zero volume: %v", err)
	}
}

func TestMassForMolarity_cases(t *testing.T) {
	got, err := MassForMolarity(0.5, 2, 58.44)
	if err != nil || math.Abs(got-58.44) > 1e-12*58.44 {
		t.Errorf("got %v %v", got, err)
	}
	for _, c := range []struct {
		c, v, mm float64
		want     error
	}{
		{-1, 1, 1, ErrNegative},
		{1, math.NaN(), 1, ErrNotFinite},
		{1, 1, math.Inf(-1), ErrNotFinite},
		{math.MaxFloat64, 2, 2, ErrOverflow},
	} {
		if _, err := MassForMolarity(c.c, c.v, c.mm); !errors.Is(err, c.want) {
			t.Errorf("MassForMolarity(%v, %v, %v): %v, want %v", c.c, c.v, c.mm, err, c.want)
		}
	}
}

func BenchmarkMolarity(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = Molarity(0.133, 0.355)
	}
}

func BenchmarkMolarityFromMass(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = MolarityFromMass(100, 58.44, 2)
	}
}

func BenchmarkMassForMolarity(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = MassForMolarity(0.5, 2, 58.44)
	}
}

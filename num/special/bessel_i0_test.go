package special

import (
	"errors"
	"math"
	"testing"
)

func TestBesselI0_TableValues(t *testing.T) {
	// A&S Table 9.8 (e^(−x) I₀(x) at x = 1: 0.4657596077); values to 16 digits from DLMF 10.25.2.
	cases := []struct{ x, want float64 }{
		{0, 1},
		{1, 1.2660658777520082},
		{-1, 1.2660658777520082},
		{5, 27.239871823604442},
		{10, 2815.716628466254},
	}
	for _, c := range cases {
		got, err := BesselI0(c.x)
		if err != nil || math.Abs(got-c.want) > 1e-14*c.want {
			t.Errorf("I0(%v) = %v, %v; want %v", c.x, got, err, c.want)
		}
	}
	got, _ := BesselI0(1)
	// e⁻¹·I₀(1) = 0.46575960759364043… (A&S Table 9.8 lists 0.4657596076; full digits from the
	// DLMF §10.25.2 series Σ (1/4)ᵏ/(k!)², I₀(1) = 1.2660658777520083…). Tolerance 1e-15 (a few ulp).
	if math.Abs(math.Exp(-1)*got-0.46575960759364043) > 1e-15 {
		t.Errorf("e^-1 I0(1) = %v; want 0.46575960759364043 (A&S Table 9.8, DLMF §10.25.2)", math.Exp(-1)*got)
	}
}

// The asymptotic branch (|x| > 20) must agree with the full positive series (oracle, no cancellation).
func TestBesselI0_AsymptoticMatchesSeries(t *testing.T) {
	for _, x := range []float64{20.0000001, 25, 40, 80} {
		q, term, sum := x*x/4, 1.0, 1.0
		for k := 1.0; k < 400; k++ {
			term *= q / (k * k)
			sum += term
		}
		got, err := BesselI0(x)
		if err != nil || math.Abs(got-sum) > 1e-13*sum {
			t.Errorf("I0(%v) = %v, %v; series %v", x, got, err, sum)
		}
	}
}

func TestBesselI0_Errors(t *testing.T) {
	for _, x := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := BesselI0(x); !errors.Is(err, ErrNotFinite) {
			t.Errorf("I0(%v): err = %v; want ErrNotFinite", x, err)
		}
	}
	if _, err := BesselI0(720); !errors.Is(err, ErrOverflow) {
		t.Errorf("I0(720): err = %v; want ErrOverflow", err)
	}
	if v, err := BesselI0(-700); err != nil || math.IsInf(v, 0) {
		t.Errorf("I0(-700) = %v, %v; want finite", v, err)
	}
}

func BenchmarkBesselI0(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkF, _ = BesselI0(8.6)
	}
}

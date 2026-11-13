package special

import (
	"errors"
	"math"
	"testing"
)

// Closed forms from DLMF Table 18.3.1 / A&S Tables 22.3–22.10; tolerance 1e-13 relative (a few ulps
// per recurrence step, n ≤ 3 here).
func TestOrthogonalPolynomials_ClosedForms(t *testing.T) {
	type poly func(int, float64) (float64, error)
	cases := []struct {
		name string
		f    poly
		n    int
		ref  func(x float64) float64
	}{
		{"LegendreP0", LegendreP, 0, func(x float64) float64 { return 1 }},
		{"LegendreP1", LegendreP, 1, func(x float64) float64 { return x }},
		{"LegendreP2", LegendreP, 2, func(x float64) float64 { return (3*x*x - 1) / 2 }},
		{"LegendreP3", LegendreP, 3, func(x float64) float64 { return (5*x*x*x - 3*x) / 2 }},
		{"ChebyshevT3", ChebyshevT, 3, func(x float64) float64 { return 4*x*x*x - 3*x }},
		{"ChebyshevU2", ChebyshevU, 2, func(x float64) float64 { return 4*x*x - 1 }},
		{"ChebyshevU3", ChebyshevU, 3, func(x float64) float64 { return 8*x*x*x - 4*x }},
		{"HermiteH2", HermiteH, 2, func(x float64) float64 { return 4*x*x - 2 }},
		{"HermiteH3", HermiteH, 3, func(x float64) float64 { return 8*x*x*x - 12*x }},
		{"HermiteHe3", HermiteHe, 3, func(x float64) float64 { return x*x*x - 3*x }},
		{"HermiteHe4", HermiteHe, 4, func(x float64) float64 { return x*x*x*x - 6*x*x + 3 }},
		{"LaguerreL1", LaguerreL, 1, func(x float64) float64 { return 1 - x }},
		{"LaguerreL2", LaguerreL, 2, func(x float64) float64 { return (x*x - 4*x + 2) / 2 }},
		{"LaguerreL3", LaguerreL, 3, func(x float64) float64 { return (-x*x*x + 9*x*x - 18*x + 6) / 6 }},
	}
	for _, c := range cases {
		for _, x := range []float64{-2, -0.7, 0, 0.5, 1, 3.25} {
			got, err := c.f(c.n, x)
			if want := c.ref(x); err != nil || !within(got, want, 1e-13) {
				t.Errorf("%s(%v) = %v, %v; want %v", c.name, x, got, err, want)
			}
		}
	}
}

func TestChebyshev_Trigonometric(t *testing.T) {
	// DLMF 18.5.1–18.5.2: T_n(cos θ) = cos nθ, U_n(cos θ) = sin((n+1)θ)/sin θ.
	for n := 0; n <= 50; n++ {
		for _, th := range []float64{0.3, 1.1, 2.9} {
			x := math.Cos(th)
			tv, _ := ChebyshevT(n, x)
			if math.Abs(tv-math.Cos(float64(n)*th)) > 1e-12 {
				t.Errorf("T_%d(cos %v) = %v; want %v", n, th, tv, math.Cos(float64(n)*th))
			}
			uv, _ := ChebyshevU(n, x)
			want := math.Sin(float64(n+1)*th) / math.Sin(th)
			if math.Abs(uv-want) > 1e-11 {
				t.Errorf("U_%d(cos %v) = %v; want %v", n, th, uv, want)
			}
		}
	}
}

func TestLegendreP_AtOne(t *testing.T) {
	// DLMF Table 18.6.1: P_n(1) = 1, P_n(−1) = (−1)^n.
	for n := 0; n <= 100; n++ {
		p, _ := LegendreP(n, 1)
		m, _ := LegendreP(n, -1)
		if math.Abs(p-1) > 1e-13 || math.Abs(m-math.Pow(-1, float64(n))) > 1e-13 {
			t.Errorf("P_%d(±1) = %v, %v", n, p, m)
		}
	}
}

func TestOrthogonalPolynomials_Errors(t *testing.T) {
	fs := map[string]func(int, float64) (float64, error){
		"LegendreP": LegendreP, "ChebyshevT": ChebyshevT, "ChebyshevU": ChebyshevU,
		"HermiteH": HermiteH, "HermiteHe": HermiteHe, "LaguerreL": LaguerreL,
	}
	for name, f := range fs {
		if _, err := f(-1, 0.5); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("%s(-1): err = %v; want ErrOutOfRange", name, err)
		}
		if _, err := f(3, math.NaN()); !errors.Is(err, ErrNotFinite) {
			t.Errorf("%s(NaN): err = %v; want ErrNotFinite", name, err)
		}
		if _, err := f(3, math.Inf(-1)); !errors.Is(err, ErrNotFinite) {
			t.Errorf("%s(-Inf): err = %v; want ErrNotFinite", name, err)
		}
		// L₁₀₀₀(1000) ≈ −2.594e215 is finite, so LaguerreL gets a larger x: L_n(x) ~ xⁿ/n! and
		// 1e6¹⁰⁰⁰/1000! ≈ 1e6000/4e2567 ≈ 1e3432 > MaxFloat64 (≈1.8e308).
		x := 1e3
		if name == "LaguerreL" {
			x = 1e6
		}
		if _, err := f(1000, x); !errors.Is(err, ErrOverflow) {
			t.Errorf("%s(1000, %v): err = %v; want ErrOverflow", name, x, err)
		}
	}
}

func BenchmarkLegendreP(b *testing.B) {
	for _, n := range []int{10, 1000} {
		b.Run(itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkF, _ = LegendreP(n, 0.37)
			}
		})
	}
}

func BenchmarkHermiteH(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkF, _ = HermiteH(20, 1.3)
	}
}

func BenchmarkLaguerreL(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkF, _ = LaguerreL(20, 1.3)
	}
}

func BenchmarkChebyshevT(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkF, _ = ChebyshevT(20, 0.3)
	}
}

func BenchmarkChebyshevU(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkF, _ = ChebyshevU(20, 0.3)
	}
}

func BenchmarkHermiteHe(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkF, _ = HermiteHe(20, 1.3)
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for ; n > 0; n /= 10 {
		i--
		buf[i] = byte('0' + n%10)
	}
	return string(buf[i:])
}

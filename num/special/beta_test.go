package special

import (
	"errors"
	"math"
	"testing"
)

func TestBeta_Values(t *testing.T) {
	cases := []struct{ a, b, want float64 }{
		{2, 3, 1.0 / 12},                 // DLMF 5.12.1: 1!·2!/4!
		{0.5, 0.5, math.Pi},              // Γ(½)² / Γ(1) = π (DLMF 5.4.6)
		{1, 4, 0.25},                     // B(1, b) = 1/b
		{10, 10, 1.0825088224469026e-06}, // 9!·9!/19!
		{0.2, 3, 0},                      // oracle below
	}
	cases[3].want = math.Exp(lgammaOf(10) + lgammaOf(10) - lgammaOf(20))
	cases[4].want = math.Gamma(0.2) * math.Gamma(3) / math.Gamma(3.2)
	for _, c := range cases {
		got, err := Beta(c.a, c.b)
		if err != nil || math.Abs(got-c.want) > 1e-13*c.want {
			t.Errorf("Beta(%v, %v) = %v, %v; want %v", c.a, c.b, got, err, c.want)
		}
		sym, _ := Beta(c.b, c.a)
		if math.Abs(sym-got) > 1e-15*got {
			t.Errorf("Beta not symmetric at (%v, %v)", c.a, c.b)
		}
	}
}

func lgammaOf(x float64) float64 { v, _ := math.Lgamma(x); return v }

func TestLogBeta_Values(t *testing.T) {
	got, err := LogBeta(2, 3)
	if err != nil || math.Abs(got+math.Log(12)) > 1e-14 {
		t.Errorf("LogBeta(2, 3) = %v, %v; want %v", got, err, -math.Log(12))
	}
	// Tiny a: ln B(a, 1) = −ln a exactly.
	got, _ = LogBeta(1e-200, 1)
	if math.Abs(got-200*math.Ln10) > 1e-12*200*math.Ln10 {
		t.Errorf("LogBeta(1e-200, 1) = %v; want %v", got, 200*math.Ln10)
	}
}

func TestBeta_Errors(t *testing.T) {
	cases := []struct {
		name string
		a, b float64
		want error
	}{
		{"a zero", 0, 1, ErrOutOfRange},
		{"b negative", 1, -2, ErrOutOfRange},
		{"a NaN", math.NaN(), 1, ErrNotFinite},
		{"b Inf", 1, math.Inf(1), ErrNotFinite},
	}
	for _, c := range cases {
		if _, err := Beta(c.a, c.b); !errors.Is(err, c.want) {
			t.Errorf("Beta %s: err = %v; want %v", c.name, err, c.want)
		}
		if _, err := LogBeta(c.a, c.b); !errors.Is(err, c.want) {
			t.Errorf("LogBeta %s: err = %v; want %v", c.name, err, c.want)
		}
	}
	// B(a, 1) = 1/a overflows for a = 1e-310 (1/a > MaxFloat64).
	if _, err := Beta(1e-310, 1); !errors.Is(err, ErrOverflow) {
		t.Errorf("Beta(1e-310, 1): err = %v; want ErrOverflow", err)
	}
}

func BenchmarkBeta(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkF, _ = Beta(2.5, 3.5)
	}
}

func BenchmarkLogBeta(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkF, _ = LogBeta(250, 3.5)
	}
}

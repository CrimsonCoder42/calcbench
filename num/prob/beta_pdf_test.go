package prob

import (
	"math"
	"testing"
)

// JKB vol. 2 ch. 25 eq. 25.2: α = 2, β = 3, x = 0.5 → 12·0.5·0.25 = 1.5 (B(2, 3) = 1/12, hand-checked).
func TestBetaPDF_sourceExample(t *testing.T) {
	if got, err := BetaPDF(0.5, 2, 3); err != nil || !within(got, 1.5, 1e-14) {
		t.Errorf("BetaPDF(0.5, 2, 3) = (%v, %v), want 1.5", got, err)
	}
}

func TestBetaPDF_cases(t *testing.T) {
	cases := []struct {
		name          string
		x, a, b, want float64
	}{
		{"outside below", -0.1, 2, 3, 0},
		{"outside above", 1.1, 2, 3, 0},
		{"x = 0, α > 1", 0, 2, 3, 0},
		{"x = 1, β = 1", 1, 2, 1, 2},
		{"uniform", 0.3, 1, 1, 1},
	}
	for _, c := range cases {
		if got, err := BetaPDF(c.x, c.a, c.b); err != nil || !within(got, c.want, 1e-14) {
			t.Errorf("%s: (%v, %v), want %v", c.name, got, err, c.want)
		}
	}
}

func TestBetaPDF_integratesToCDF(t *testing.T) {
	f := func(x float64) float64 { d, _ := BetaPDF(x, 2.5, 4); return d }
	ca, _ := BetaCDF(0.1, 2.5, 4)
	cb, _ := BetaCDF(0.7, 2.5, 4)
	if got := simpson(f, 0.1, 0.7, 2000); math.Abs(got-(cb-ca)) > 1e-10 {
		t.Errorf("∫f = %v, CDF diff = %v", got, cb-ca)
	}
}

func TestBetaPDF_errors(t *testing.T) {
	_, err := BetaPDF(0, 0.5, 2)
	requireErr(t, "x = 0, α < 1", err, ErrOutOfRange)
	_, err = BetaPDF(1, 2, 0.5)
	requireErr(t, "x = 1, β < 1", err, ErrOutOfRange)
	_, err = BetaPDF(0.5, 0, 1)
	requireErr(t, "α = 0", err, ErrOutOfRange)
	_, err = BetaPDF(nan, 1, 1)
	requireErr(t, "x NaN", err, ErrNotFinite)
	_, err = BetaPDF(5e-324, 0.01, 1)
	requireErr(t, "overflow", err, ErrOverflow)
}

func BenchmarkBetaPDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = BetaPDF(0.3, 2, 3)
	}
}

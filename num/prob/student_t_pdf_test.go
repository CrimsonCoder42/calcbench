package prob

import (
	"math"
	"testing"
)

// A&S 26.7.1 / JKB ch. 28 eq. 28.2: ν = 1 is Cauchy, f(0) = 1/π; ν = 2, f(0) = 1/(2√2).
func TestStudentTPDF_sourceExample(t *testing.T) {
	cases := []struct{ x, df, want float64 }{
		{0, 1, 1 / math.Pi},
		{1, 1, 1 / (2 * math.Pi)},
		{0, 2, 1 / (2 * math.Sqrt2)},
		{1, 2, math.Pow(1.5, -1.5) / (2 * math.Sqrt2)},
	}
	for _, c := range cases {
		if got, err := StudentTPDF(c.x, c.df); err != nil || !within(got, c.want, 1e-14) {
			t.Errorf("StudentTPDF(%v, %v) = (%v, %v), want %v", c.x, c.df, got, err, c.want)
		}
	}
}

// Large ν tends to the standard normal density; symmetric; integrates to the CDF difference.
func TestStudentTPDF_properties(t *testing.T) {
	g, _ := StudentTPDF(1, 1e6)
	if !within(g, math.Exp(-0.5)/math.Sqrt(2*math.Pi), 1e-5) {
		t.Errorf("ν → ∞: %v", g)
	}
	a, _ := StudentTPDF(-1.3, 5)
	b, _ := StudentTPDF(1.3, 5)
	if a != b {
		t.Errorf("asymmetric: %v %v", a, b)
	}
	f := func(x float64) float64 { d, _ := StudentTPDF(x, 5); return d }
	ca, _ := StudentTCDF(-2, 5)
	cb, _ := StudentTCDF(1, 5)
	if got := simpson(f, -2, 1, 2000); math.Abs(got-(cb-ca)) > 1e-11 {
		t.Errorf("∫f = %v, CDF diff = %v", got, cb-ca)
	}
}

func TestStudentTPDF_errors(t *testing.T) {
	_, err := StudentTPDF(0, 0)
	requireErr(t, "df 0", err, ErrOutOfRange)
	_, err = StudentTPDF(nan, 1)
	requireErr(t, "t NaN", err, ErrNotFinite)
	if got, err := StudentTPDF(1e200, 3); err != nil || got != 0 {
		t.Errorf("far tail: (%v, %v)", got, err)
	}
}

func BenchmarkStudentTPDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = StudentTPDF(1.5, 5)
	}
}

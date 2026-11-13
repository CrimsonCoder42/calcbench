package prob

import (
	"math"
	"testing"
)

// Ross §5.3 (F(x) = (x − a)/(b − a)): a = 0, b = 10, x = 3 → 0.3. Hand-checked.
func TestUniformCDF_sourceExample(t *testing.T) {
	if got, err := UniformCDF(3, 0, 10); err != nil || got != 0.3 {
		t.Errorf("UniformCDF(3, 0, 10) = (%v, %v), want (0.3, nil)", got, err)
	}
}

func TestUniformCDF_cases(t *testing.T) {
	for _, c := range []struct{ x, a, b, want float64 }{{-1, 0, 10, 0}, {0, 0, 10, 0}, {10, 0, 10, 1}, {1e300, 0, 10, 1}, {0, -1, 1, 0.5}} {
		if got, err := UniformCDF(c.x, c.a, c.b); err != nil || got != c.want {
			t.Errorf("UniformCDF(%v, %v, %v) = (%v, %v), want %v", c.x, c.a, c.b, got, err, c.want)
		}
	}
}

// b − a overflows float64 but (x − a)/(b − a) is representable: computed from the halved terms. Expected
// values hand-checked (x = 0 is the midpoint; x = 5e307 is three quarters of the way; x = −Max/2 one quarter).
func TestUniformCDF_wideInterval(t *testing.T) {
	for _, c := range []struct{ x, a, b, want float64 }{
		{0, -1e308, 1e308, 0.5},
		{5e307, -1e308, 1e308, 0.75},
		{0, -math.MaxFloat64, math.MaxFloat64, 0.5},
		{-math.MaxFloat64 / 2, -math.MaxFloat64, math.MaxFloat64, 0.25},
	} {
		if got, err := UniformCDF(c.x, c.a, c.b); err != nil || !within(got, c.want, 2.3e-16) {
			t.Errorf("UniformCDF(%v, %v, %v) = (%v, %v), want (%v, nil)", c.x, c.a, c.b, got, err, c.want)
		}
	}
}

func TestUniformCDF_errors(t *testing.T) {
	_, err := UniformCDF(0, 1, 0)
	requireErr(t, "a > b", err, ErrOutOfRange)
	_, err = UniformCDF(inf, 0, 1)
	requireErr(t, "x Inf", err, ErrNotFinite)
}

func TestUniformCDF_noAlloc(t *testing.T) {
	requireNoAlloc(t, "UniformCDF", func() { sinkFloat, _ = UniformCDF(3, 0, 10) })
}

func BenchmarkUniformCDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = UniformCDF(3, 0, 10)
	}
}

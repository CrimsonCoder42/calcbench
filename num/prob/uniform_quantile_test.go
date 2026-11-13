package prob

import (
	"math"
	"testing"
)

// Inverting the uniform CDF of NIST §1.3.6.6.2 / Ross §5.3: a = 0, b = 10, u = 0.3 → 3. Hand-checked
// (FMA: one rounding of 0.3·10, which is 3 exactly).
func TestUniformQuantile_sourceExample(t *testing.T) {
	if got, err := UniformQuantile(0.3, 0, 10); err != nil || got != 3 {
		t.Errorf("UniformQuantile(0.3, 0, 10) = (%v, %v), want (3, nil)", got, err)
	}
}

func TestUniformQuantile_cases(t *testing.T) {
	for _, c := range []struct{ u, a, b, want float64 }{{0, -2, 2, -2}, {1, -2, 2, 2}, {0.5, -2, 2, 0}, {1, 0, math.MaxFloat64, math.MaxFloat64}} {
		if got, err := UniformQuantile(c.u, c.a, c.b); err != nil || got != c.want {
			t.Errorf("UniformQuantile(%v, %v, %v) = (%v, %v), want %v", c.u, c.a, c.b, got, err, c.want)
		}
	}
}

// b − a overflows float64 but a + u·(b − a) is representable: computed as 2·(a/2 + u·(b/2 − a/2)). Expected
// values hand-checked (exact in binary: the halves of ±1e308 and ±Max are exact).
func TestUniformQuantile_wideInterval(t *testing.T) {
	for _, c := range []struct{ u, a, b, want float64 }{
		{0.5, -1e308, 1e308, 0},
		{0.5, -math.MaxFloat64, math.MaxFloat64, 0},
		{0, -1e308, 1e308, -1e308},
		{1, -1e308, 1e308, 1e308},
		{0.75, -1e308, 1e308, 5e307},
		{0.25, -math.MaxFloat64, math.MaxFloat64, -math.MaxFloat64 / 2},
	} {
		if got, err := UniformQuantile(c.u, c.a, c.b); err != nil || got != c.want {
			t.Errorf("UniformQuantile(%v, %v, %v) = (%v, %v), want (%v, nil)", c.u, c.a, c.b, got, err, c.want)
		}
	}
}

func TestUniformQuantile_errors(t *testing.T) {
	_, err := UniformQuantile(1.5, 0, 1)
	requireErr(t, "u > 1", err, ErrProbability)
	_, err = UniformQuantile(0.5, 1, 1)
	requireErr(t, "a = b", err, ErrOutOfRange)
	_, err = UniformQuantile(nan, 0, 1)
	requireErr(t, "NaN", err, ErrNotFinite)
}

// Oracle: UniformCDF(UniformQuantile(u)) = u within 2 ulps of 1.
func TestUniformQuantile_oracle(t *testing.T) {
	rng := newTestRand()
	for trial := range 500 {
		a := (rng.Float64() - 0.5) * 100
		b := a + 0.1 + 100*rng.Float64()
		u := rng.Float64()
		x, err := UniformQuantile(u, a, b)
		if err != nil {
			t.Fatal(err)
		}
		if back, _ := UniformCDF(x, a, b); absDiff(back, u) > 1e-13 {
			t.Fatalf("trial %d: UniformCDF(UniformQuantile(%v)) = %v", trial, u, back)
		}
	}
}

func TestUniformQuantile_noAlloc(t *testing.T) {
	requireNoAlloc(t, "UniformQuantile", func() { sinkFloat, _ = UniformQuantile(0.3, 0, 10) })
}

func BenchmarkUniformQuantile(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = UniformQuantile(0.3, 0, 10)
	}
}

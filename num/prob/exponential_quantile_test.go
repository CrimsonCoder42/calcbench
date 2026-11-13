package prob

import (
	"math"
	"testing"
)

// NIST §1.3.6.6.7 (PPF G(p) = −β ln(1 − p)): λ = 2, u = 0.5 → ln 2/2 = 0.34657359027997264. Hand-checked;
// tolerance 1 ulp.
func TestExponentialQuantile_sourceExample(t *testing.T) {
	if got, err := ExponentialQuantile(0.5, 2); err != nil || !within(got, 0.34657359027997264, 2.5e-16) {
		t.Errorf("ExponentialQuantile(0.5, 2) = (%v, %v), want (0.34657359027997264, nil)", got, err)
	}
}

func TestExponentialQuantile_cases(t *testing.T) {
	for _, c := range []struct{ u, rate, want float64 }{{0, 2, 0}, {1e-20, 1, 1e-20}, {1 - 0x1p-53, 1, 53 * math.Ln2}} {
		if got, err := ExponentialQuantile(c.u, c.rate); err != nil || !within(got, c.want, 5e-16) {
			t.Errorf("ExponentialQuantile(%v, %v) = (%v, %v), want %v", c.u, c.rate, got, err, c.want)
		}
	}
}

func TestExponentialQuantile_errors(t *testing.T) {
	_, err := ExponentialQuantile(1, 2)
	requireErr(t, "u = 1", err, ErrOutOfRange)
	_, err = ExponentialQuantile(0.5, 0)
	requireErr(t, "rate = 0", err, ErrOutOfRange)
	_, err = ExponentialQuantile(0.5, 5e-324)
	requireErr(t, "quantile overflows", err, ErrOverflow)
	_, err = ExponentialQuantile(-0.5, 1)
	requireErr(t, "u < 0", err, ErrProbability)
	_, err = ExponentialQuantile(nan, 1)
	requireErr(t, "NaN", err, ErrNotFinite)
}

// Oracle: ExponentialCDF(ExponentialQuantile(u)) = u. Tolerance 4 ulps relative.
func TestExponentialQuantile_oracle(t *testing.T) {
	rng := newTestRand()
	for trial := range 500 {
		u, rate := rng.Float64(), 0.01+100*rng.Float64()
		x, err := ExponentialQuantile(u, rate)
		if err != nil {
			t.Fatal(err)
		}
		if back, _ := ExponentialCDF(x, rate); !within(back, u, 1e-15) {
			t.Fatalf("trial %d: ExponentialCDF(ExponentialQuantile(%v, %v)) = %v", trial, u, rate, back)
		}
	}
}

func TestExponentialQuantile_noAlloc(t *testing.T) {
	requireNoAlloc(t, "ExponentialQuantile", func() { sinkFloat, _ = ExponentialQuantile(0.5, 2) })
}

func BenchmarkExponentialQuantile(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = ExponentialQuantile(0.5, 2)
	}
}

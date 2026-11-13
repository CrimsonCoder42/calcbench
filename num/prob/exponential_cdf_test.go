package prob

import "testing"

// NIST §1.3.6.6.7 (F(x) = 1 − e^(−x/β)): λ = 2, x = 1 → 1 − e⁻² = 0.8646647167633873. Hand-checked;
// tolerance 2 ulps.
func TestExponentialCDF_sourceExample(t *testing.T) {
	if got, err := ExponentialCDF(1, 2); err != nil || !within(got, 0.8646647167633873, 5e-16) {
		t.Errorf("ExponentialCDF(1, 2) = (%v, %v), want (0.8646647167633873, nil)", got, err)
	}
}

func TestExponentialCDF_cases(t *testing.T) {
	for _, c := range []struct{ x, rate, want float64 }{{0, 2, 0}, {-1, 2, 0}, {1e300, 1, 1}, {1e-20, 1, 1e-20}} {
		if got, err := ExponentialCDF(c.x, c.rate); err != nil || !within(got, c.want, 5e-16) {
			t.Errorf("ExponentialCDF(%v, %v) = (%v, %v), want %v", c.x, c.rate, got, err, c.want)
		}
	}
}

func TestExponentialCDF_errors(t *testing.T) {
	_, err := ExponentialCDF(1, -2)
	requireErr(t, "rate < 0", err, ErrOutOfRange)
	_, err = ExponentialCDF(inf, 1)
	requireErr(t, "x Inf", err, ErrNotFinite)
}

func TestExponentialCDF_noAlloc(t *testing.T) {
	requireNoAlloc(t, "ExponentialCDF", func() { sinkFloat, _ = ExponentialCDF(1, 2) })
}

func BenchmarkExponentialCDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = ExponentialCDF(1, 2)
	}
}

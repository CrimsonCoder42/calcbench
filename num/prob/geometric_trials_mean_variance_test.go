package prob

import "testing"

// Casella–Berger §3.2 (geometric: EX = 1/p, Var X = (1−p)/p²): p = 0.2 → (5, 20). Hand-checked;
// tolerance 2 ulps.
func TestGeometricTrialsMeanVariance_sourceExample(t *testing.T) {
	m, v, err := GeometricTrialsMeanVariance(0.2)
	if err != nil || !within(m, 5, 5e-16) || !within(v, 20, 5e-16) {
		t.Errorf("GeometricTrialsMeanVariance(0.2) = (%v, %v, %v), want (5, 20, nil)", m, v, err)
	}
}

func TestGeometricTrialsMeanVariance_cases(t *testing.T) {
	for _, c := range []struct{ p, m, v float64 }{{1, 1, 0}, {0.5, 2, 2}, {0.25, 4, 12}} {
		if m, v, err := GeometricTrialsMeanVariance(c.p); err != nil || m != c.m || v != c.v {
			t.Errorf("GeometricTrialsMeanVariance(%v) = (%v, %v, %v), want (%v, %v)", c.p, m, v, err, c.m, c.v)
		}
	}
}

func TestGeometricTrialsMeanVariance_errors(t *testing.T) {
	_, _, err := GeometricTrialsMeanVariance(0)
	requireErr(t, "p = 0", err, ErrOutOfRange)
	_, _, err = GeometricTrialsMeanVariance(1e-200)
	requireErr(t, "variance overflows", err, ErrOverflow)
	_, _, err = GeometricTrialsMeanVariance(nan)
	requireErr(t, "NaN", err, ErrNotFinite)
}

func TestGeometricTrialsMeanVariance_noAlloc(t *testing.T) {
	requireNoAlloc(t, "GeometricTrialsMeanVariance", func() { sinkFloat, _, _ = GeometricTrialsMeanVariance(0.2) })
}

func BenchmarkGeometricTrialsMeanVariance(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _, _ = GeometricTrialsMeanVariance(0.2)
	}
}

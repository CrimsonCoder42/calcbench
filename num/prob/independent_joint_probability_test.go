package prob

import "testing"

// Ross §3.4 (independent events): two fair coins both heads, 0.5·0.5 = 0.25. Hand-checked.
func TestIndependentJointProbability_sourceExample(t *testing.T) {
	if got, err := IndependentJointProbability(0.5, 0.5); err != nil || got != 0.25 {
		t.Errorf("IndependentJointProbability(0.5, 0.5) = (%v, %v), want (0.25, nil)", got, err)
	}
}

func TestIndependentJointProbability_cases(t *testing.T) {
	for _, c := range []struct{ a, b, want float64 }{{0, 1, 0}, {1, 1, 1}, {0.2, 0.5, 0.1}, {1, 0.3, 0.3}} {
		if got, err := IndependentJointProbability(c.a, c.b); err != nil || got != c.want {
			t.Errorf("IndependentJointProbability(%v, %v) = (%v, %v), want %v", c.a, c.b, got, err, c.want)
		}
	}
}

func TestIndependentJointProbability_errors(t *testing.T) {
	_, err := IndependentJointProbability(nan, 0.5)
	requireErr(t, "pA NaN", err, ErrNotFinite)
	_, err = IndependentJointProbability(0.5, inf)
	requireErr(t, "pB Inf", err, ErrNotFinite)
	_, err = IndependentJointProbability(1.5, 0.5)
	requireErr(t, "pA > 1", err, ErrProbability)
	_, err = IndependentJointProbability(0.5, -0.5)
	requireErr(t, "pB < 0", err, ErrProbability)
}

func TestIndependentJointProbability_noAlloc(t *testing.T) {
	requireNoAlloc(t, "IndependentJointProbability", func() { sinkFloat, _ = IndependentJointProbability(0.5, 0.5) })
}

func BenchmarkIndependentJointProbability(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = IndependentJointProbability(0.5, 0.5)
	}
}

package prob

import "testing"

// Ross §3.3 / Casella–Berger §1.3: P(B|A) = 0.9, P(A) = 0.01, P(B) = 0.0585 → 0.009/0.0585 = 2/13.
// Hand-checked with exact fractions; tolerance 4 ulps for the rounded decimal inputs.
func TestBayesTheorem_sourceExample(t *testing.T) {
	if got, err := BayesTheorem(0.9, 0.01, 0.0585); err != nil || !within(got, 2.0/13, 1e-15) {
		t.Errorf("BayesTheorem(0.9, 0.01, 0.0585) = (%v, %v), want (2/13, nil)", got, err)
	}
}

func TestBayesTheorem_cases(t *testing.T) {
	for _, c := range []struct{ ba, a, b, want float64 }{{0, 0.5, 0.5, 0}, {1, 0.5, 0.5, 1}, {0.5, 1, 0.5, 1}, {0.5, 0.5, 1, 0.25}} {
		if got, err := BayesTheorem(c.ba, c.a, c.b); err != nil || got != c.want {
			t.Errorf("BayesTheorem(%v, %v, %v) = (%v, %v), want %v", c.ba, c.a, c.b, got, err, c.want)
		}
	}
}

func TestBayesTheorem_errors(t *testing.T) {
	_, err := BayesTheorem(0.5, 0.5, 0)
	requireErr(t, "pB = 0", err, ErrOutOfRange)
	_, err = BayesTheorem(0.9, 0.9, 0.5)
	requireErr(t, "posterior > 1", err, ErrInconsistent)
	_, err = BayesTheorem(nan, 0.5, 0.5)
	requireErr(t, "NaN", err, ErrNotFinite)
	_, err = BayesTheorem(0.5, -0.5, 0.5)
	requireErr(t, "pA < 0", err, ErrProbability)
}

func TestBayesTheorem_noAlloc(t *testing.T) {
	requireNoAlloc(t, "BayesTheorem", func() { sinkFloat, _ = BayesTheorem(0.9, 0.01, 0.0585) })
}

func BenchmarkBayesTheorem(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = BayesTheorem(0.9, 0.01, 0.0585)
	}
}

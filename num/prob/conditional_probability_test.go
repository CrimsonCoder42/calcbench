package prob

import "testing"

// Ross §3.2 / Casella–Berger §1.3: one die, A = {6}, B = {even}: (1/6)/(1/2) = 1/3.
// Hand-checked; tolerance 2 ulps.
func TestConditionalProbability_sourceExample(t *testing.T) {
	if got, err := ConditionalProbability(1.0/6, 0.5); err != nil || !within(got, 1.0/3, 5e-16) {
		t.Errorf("ConditionalProbability(1/6, 1/2) = (%v, %v), want (1/3, nil)", got, err)
	}
}

func TestConditionalProbability_cases(t *testing.T) {
	for _, c := range []struct{ ab, b, want float64 }{{0, 0.5, 0}, {0.5, 0.5, 1}, {0.25, 1, 0.25}, {5e-324, 5e-324, 1}} {
		if got, err := ConditionalProbability(c.ab, c.b); err != nil || got != c.want {
			t.Errorf("ConditionalProbability(%v, %v) = (%v, %v), want %v", c.ab, c.b, got, err, c.want)
		}
	}
}

func TestConditionalProbability_errors(t *testing.T) {
	_, err := ConditionalProbability(0, 0)
	requireErr(t, "pB = 0 (never NaN)", err, ErrOutOfRange)
	_, err = ConditionalProbability(0.6, 0.5)
	requireErr(t, "pAAndB > pB", err, ErrInconsistent)
	_, err = ConditionalProbability(nan, 0.5)
	requireErr(t, "NaN", err, ErrNotFinite)
	_, err = ConditionalProbability(0.1, 1.5)
	requireErr(t, "pB > 1", err, ErrProbability)
}

func TestConditionalProbability_noAlloc(t *testing.T) {
	requireNoAlloc(t, "ConditionalProbability", func() { sinkFloat, _ = ConditionalProbability(0.25, 0.5) })
}

func BenchmarkConditionalProbability(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = ConditionalProbability(0.25, 0.5)
	}
}

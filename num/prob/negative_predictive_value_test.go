package prob

import "testing"

// Ross §3.3 (Bayes's formula) on the same 2×2 table: NPV = TN/(TN + FN) = 54/81 = 2/3. Hand-checked with
// exact fractions; tolerance 8 ulps for the rounded inputs.
func TestNegativePredictiveValue_sourceExample(t *testing.T) {
	got, err := NegativePredictiveValue(231.0/258, 54.0/86, 258.0/344)
	if err != nil || !within(got, 2.0/3, 2e-15) {
		t.Errorf("NegativePredictiveValue(2×2 table) = (%v, %v), want (2/3, nil)", got, err)
	}
}

func TestNegativePredictiveValue_cases(t *testing.T) {
	for _, c := range []struct{ se, sp, pi, want float64 }{{1, 0.9, 0.3, 1}, {0.9, 0.9, 1, 0}, {0.5, 0.5, 0.2, 0.8}} {
		if got, err := NegativePredictiveValue(c.se, c.sp, c.pi); err != nil || !within(got, c.want, 1e-15) {
			t.Errorf("NPV(%v, %v, %v) = (%v, %v), want %v", c.se, c.sp, c.pi, got, err, c.want)
		}
	}
}

func TestNegativePredictiveValue_errors(t *testing.T) {
	_, err := NegativePredictiveValue(1, 0, 0.5)
	requireErr(t, "no negative possible", err, ErrOutOfRange)
	_, err = NegativePredictiveValue(0.5, inf, 0.5)
	requireErr(t, "Inf", err, ErrNotFinite)
	_, err = NegativePredictiveValue(1.1, 0.5, 0.5)
	requireErr(t, "sensitivity > 1", err, ErrProbability)
}

func TestNegativePredictiveValue_noAlloc(t *testing.T) {
	requireNoAlloc(t, "NegativePredictiveValue", func() { sinkFloat, _ = NegativePredictiveValue(0.9, 0.6, 0.75) })
}

func BenchmarkNegativePredictiveValue(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = NegativePredictiveValue(0.9, 0.6, 0.75)
	}
}

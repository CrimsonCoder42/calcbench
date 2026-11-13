package prob

import "testing"

// Ross §8.2 (Markov's inequality; the weekly-production illustration, example number not confirmed in
// print): mean 50, threshold 75 → 2/3. Hand-checked; tolerance 1 ulp.
func TestMarkovBound_sourceExample(t *testing.T) {
	if got, err := MarkovBound(50, 75); err != nil || !within(got, 2.0/3, 2.5e-16) {
		t.Errorf("MarkovBound(50, 75) = (%v, %v), want (2/3, nil)", got, err)
	}
}

func TestMarkovBound_cases(t *testing.T) {
	for _, c := range []struct{ m, a, want float64 }{{0, 1, 0}, {2, 1, 1}, {1, 4, 0.25}, {1e300, 1e-10, 1}} {
		if got, err := MarkovBound(c.m, c.a); err != nil || got != c.want {
			t.Errorf("MarkovBound(%v, %v) = (%v, %v), want %v", c.m, c.a, got, err, c.want)
		}
	}
}

func TestMarkovBound_errors(t *testing.T) {
	_, err := MarkovBound(-1, 1)
	requireErr(t, "mean < 0", err, ErrOutOfRange)
	_, err = MarkovBound(1, 0)
	requireErr(t, "threshold = 0", err, ErrOutOfRange)
	_, err = MarkovBound(nan, 1)
	requireErr(t, "NaN", err, ErrNotFinite)
}

func TestMarkovBound_noAlloc(t *testing.T) {
	requireNoAlloc(t, "MarkovBound", func() { sinkFloat, _ = MarkovBound(50, 75) })
}

func BenchmarkMarkovBound(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = MarkovBound(50, 75)
	}
}

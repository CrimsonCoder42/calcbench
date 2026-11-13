package prob

import "testing"

// Ross §3.2 (multiplication rule): two cards without replacement, both aces: (3/51)·(4/52) = 1/221.
// Hand-checked; tolerance 2 ulps (two rounded quotients and one product).
func TestJointProbability_sourceExample(t *testing.T) {
	if got, err := JointProbability(3.0/51, 4.0/52); err != nil || !within(got, 1.0/221, 5e-16) {
		t.Errorf("JointProbability(3/51, 4/52) = (%v, %v), want (%v, nil)", got, err, 1.0/221)
	}
}

func TestJointProbability_cases(t *testing.T) {
	for _, c := range []struct{ a, b, want float64 }{{0, 0.7, 0}, {1, 0.7, 0.7}, {0.5, 0.5, 0.25}} {
		if got, err := JointProbability(c.a, c.b); err != nil || got != c.want {
			t.Errorf("JointProbability(%v, %v) = (%v, %v), want %v", c.a, c.b, got, err, c.want)
		}
	}
}

func TestJointProbability_errors(t *testing.T) {
	_, err := JointProbability(nan, 0.5)
	requireErr(t, "pAGivenB NaN", err, ErrNotFinite)
	_, err = JointProbability(0.5, 2)
	requireErr(t, "pB > 1", err, ErrProbability)
	_, err = JointProbability(-1, 0.5)
	requireErr(t, "pAGivenB < 0", err, ErrProbability)
}

func TestJointProbability_noAlloc(t *testing.T) {
	requireNoAlloc(t, "JointProbability", func() { sinkFloat, _ = JointProbability(0.5, 0.5) })
}

func BenchmarkJointProbability(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = JointProbability(0.5, 0.5)
	}
}

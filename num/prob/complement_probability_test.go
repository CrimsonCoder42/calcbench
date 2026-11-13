package prob

import "testing"

// Ross §2.4 / Casella–Berger §1.2: P(Aᶜ) = 1 − P(A); 0.25 → 0.75. Hand-checked.
func TestComplementProbability_sourceExample(t *testing.T) {
	if got, err := ComplementProbability(0.25); err != nil || got != 0.75 {
		t.Errorf("ComplementProbability(0.25) = (%v, %v), want (0.75, nil)", got, err)
	}
}

func TestComplementProbability_cases(t *testing.T) {
	for _, c := range []struct{ p, want float64 }{{0, 1}, {1, 0}, {0.5, 0.5}, {0.75, 0.25}, {5e-324, 1}} {
		if got, err := ComplementProbability(c.p); err != nil || got != c.want {
			t.Errorf("ComplementProbability(%v) = (%v, %v), want (%v, nil)", c.p, got, err, c.want)
		}
	}
}

func TestComplementProbability_errors(t *testing.T) {
	for _, p := range []float64{nan, inf, -inf} {
		_, err := ComplementProbability(p)
		requireErr(t, "ComplementProbability", err, ErrNotFinite)
	}
	for _, p := range []float64{-0.1, 1.0000000000000002} {
		_, err := ComplementProbability(p)
		requireErr(t, "ComplementProbability", err, ErrProbability)
	}
}

func TestComplementProbability_noAlloc(t *testing.T) {
	requireNoAlloc(t, "ComplementProbability", func() { sinkFloat, _ = ComplementProbability(0.25) })
}

func BenchmarkComplementProbability(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = ComplementProbability(0.25)
	}
}

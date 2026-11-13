package prob

import (
	"math"
	"testing"
)

// Egidi et al. 2018 §2 (inverse odds): 2.5 → 0.4; 1.9 → 1/1.9 = 0.5263157894736842 (correctly rounded). Hand-checked.
func TestImpliedProbability_sourceExample(t *testing.T) {
	for _, c := range []struct{ d, want float64 }{{2.5, 0.4}, {1.9, 0.5263157894736842}} {
		if got, err := ImpliedProbability(c.d); err != nil || got != c.want {
			t.Errorf("ImpliedProbability(%v) = (%v, %v), want (%v, nil)", c.d, got, err, c.want)
		}
	}
}

// Wheatcroft 2021, preprint §2.1: O₁ = 3, O₂ = 1.4 → r₁ = 1/3 ≈ 0.33, r₂ = 1/1.4 ≈ 0.71 (the source's rounded
// values; checked here against the correctly rounded quotients and to the source's 2 decimals).
func TestImpliedProbability_wheatcroftExample(t *testing.T) {
	for _, c := range []struct{ d, want, rounded float64 }{{3, 1.0 / 3, 0.33}, {1.4, 1 / 1.4, 0.71}} {
		got, err := ImpliedProbability(c.d)
		if err != nil || got != c.want || math.Round(got*100)/100 != c.rounded {
			t.Errorf("ImpliedProbability(%v) = (%v, %v), want (%v ≈ %v, nil)", c.d, got, err, c.want, c.rounded)
		}
	}
}

func TestImpliedProbability_cases(t *testing.T) {
	for _, c := range []struct{ d, want float64 }{{1, 1}, {2, 0.5}, {1e300, 1e-300}} {
		if got, err := ImpliedProbability(c.d); err != nil || !within(got, c.want, 5e-16) {
			t.Errorf("ImpliedProbability(%v) = (%v, %v), want %v", c.d, got, err, c.want)
		}
	}
}

func TestImpliedProbability_errors(t *testing.T) {
	_, err := ImpliedProbability(0.99)
	requireErr(t, "d < 1", err, ErrOutOfRange)
	_, err = ImpliedProbability(inf)
	requireErr(t, "Inf", err, ErrNotFinite)
}

func TestImpliedProbability_noAlloc(t *testing.T) {
	requireNoAlloc(t, "ImpliedProbability", func() { sinkFloat, _ = ImpliedProbability(1.9) })
}

func BenchmarkImpliedProbability(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = ImpliedProbability(1.9)
	}
}

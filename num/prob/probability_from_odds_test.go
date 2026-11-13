package prob

import (
	"math"
	"math/rand/v2"
	"testing"
)

// newTestRand returns the package's fixed-seed generator for oracle tests.
func newTestRand() *rand.Rand {
	return rand.New(rand.NewPCG(1, 2))
}

// Bland, Altman, BMJ 2000;320:1468 (odds = p/(1 − p), so p = o/(1 + o)): 3 → 0.75. Hand-checked.
func TestProbabilityFromOdds_sourceExample(t *testing.T) {
	if got, err := ProbabilityFromOdds(3); err != nil || got != 0.75 {
		t.Errorf("ProbabilityFromOdds(3) = (%v, %v), want (0.75, nil)", got, err)
	}
}

func TestProbabilityFromOdds_cases(t *testing.T) {
	for _, c := range []struct{ o, want float64 }{{0, 0}, {1, 0.5}, {0.25, 0.2}, {1e300, 1}, {math.MaxFloat64, 1}, {5e-324, 5e-324}} {
		if got, err := ProbabilityFromOdds(c.o); err != nil || !within(got, c.want, 5e-16) {
			t.Errorf("ProbabilityFromOdds(%v) = (%v, %v), want %v", c.o, got, err, c.want)
		}
	}
}

func TestProbabilityFromOdds_errors(t *testing.T) {
	_, err := ProbabilityFromOdds(-1)
	requireErr(t, "odds < 0", err, ErrOutOfRange)
	_, err = ProbabilityFromOdds(inf)
	requireErr(t, "Inf", err, ErrNotFinite)
}

func TestProbabilityFromOdds_noAlloc(t *testing.T) {
	requireNoAlloc(t, "ProbabilityFromOdds", func() { sinkFloat, _ = ProbabilityFromOdds(3) })
}

func BenchmarkProbabilityFromOdds(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = ProbabilityFromOdds(3)
	}
}

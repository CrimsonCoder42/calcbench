package prob

import (
	"math"
	"math/rand/v2"
	"testing"
)

// Deeks & Altman, BMJ 2004;329:168: pre-test probability 0.1 → odds 0.11; × LR 20.43 → 2.27; → 0.69.
// Exactly: 0.1·20.43/(0.1·20.43 + 0.9) = 2.043/2.943. Tolerance 1e-14 relative.
func TestPostTestProbability_sourceExample(t *testing.T) {
	got, err := PostTestProbability(0.1, 20.43)
	if want := 2.043 / 2.943; err != nil || !within(got, want, 1e-14) {
		t.Errorf("PostTestProbability(0.1, 20.43) = (%v, %v), want %v", got, err, want)
	}
	if math.Round(got*100)/100 != 0.69 {
		t.Errorf("PostTestProbability(0.1, 20.43) = %v, want 0.69 as printed", got)
	}
}

func TestPostTestProbability_cases(t *testing.T) {
	tests := []struct {
		name  string
		p, lr float64
		want  float64
	}{
		{"LR 1 leaves p unchanged", 0.37, 1, 0.37},
		{"LR 0 rules out", 0.5, 0, 0},
		{"p 0 stays 0", 0, 50, 0},
		{"p 1 stays 1", 1, 0.2, 1},
		{"p 1, large LR", 1, 1e300, 1},
		{"even odds, LR 3", 0.5, 3, 0.75},
		{"even odds, LR 1/3", 0.5, 1.0 / 3, 0.25},
		{"huge LR", 0.5, math.MaxFloat64, 1},
	}
	for _, tc := range tests {
		got, err := PostTestProbability(tc.p, tc.lr)
		if err != nil || !within(got, tc.want, 1e-15) {
			t.Errorf("%s: PostTestProbability(%v, %v) = (%v, %v), want %v", tc.name, tc.p, tc.lr, got, err, tc.want)
		}
	}
}

// Oracle: the odds route of the source (OddsFromProbability, × LR, ProbabilityFromOdds) on 300 random
// cases with p in (0, 1) and LR in [1e-3, 1e3], fixed seed. Tolerance 1e-13 relative.
func TestPostTestProbability_oddsOracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		p := rng.Float64()*0.998 + 0.001
		lr := math.Exp((rng.Float64()*2 - 1) * math.Log(1e3))
		odds, _ := OddsFromProbability(p)
		want, _ := ProbabilityFromOdds(odds * lr)
		got, err := PostTestProbability(p, lr)
		if err != nil || !within(got, want, 1e-13) {
			t.Fatalf("PostTestProbability(%v, %v) = (%v, %v), want %v", p, lr, got, err, want)
		}
	}
}

func TestPostTestProbability_errors(t *testing.T) {
	_, err := PostTestProbability(nan, 2)
	requireErr(t, "PostTestProbability(NaN, 2)", err, ErrNotFinite)
	_, err = PostTestProbability(0.5, inf)
	requireErr(t, "PostTestProbability(0.5, Inf)", err, ErrNotFinite)
	_, err = PostTestProbability(1.5, 2)
	requireErr(t, "PostTestProbability(1.5, 2)", err, ErrProbability)
	_, err = PostTestProbability(0.5, -1)
	requireErr(t, "PostTestProbability(0.5, -1)", err, ErrOutOfRange)
	_, err = PostTestProbability(1, 0)
	requireErr(t, "PostTestProbability(1, 0)", err, ErrOutOfRange)
}

func TestPostTestProbability_noAlloc(t *testing.T) {
	requireNoAlloc(t, "PostTestProbability", func() { sinkFloat, _ = PostTestProbability(0.2, 4) })
}

func BenchmarkPostTestProbability(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = PostTestProbability(0.2, 4)
	}
}

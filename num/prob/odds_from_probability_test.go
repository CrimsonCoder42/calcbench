package prob

import "testing"

// Bland, Altman, BMJ 2000;320:1468: odds = P(event)/P(no event); p = 0.75 → 3 (3 to 1). Hand-checked.
func TestOddsFromProbability_sourceExample(t *testing.T) {
	if got, err := OddsFromProbability(0.75); err != nil || got != 3 {
		t.Errorf("OddsFromProbability(0.75) = (%v, %v), want (3, nil)", got, err)
	}
}

func TestOddsFromProbability_cases(t *testing.T) {
	for _, c := range []struct{ p, want float64 }{{0, 0}, {0.5, 1}, {0.2, 0.25}, {0.9375, 15}} {
		if got, err := OddsFromProbability(c.p); err != nil || !within(got, c.want, 5e-16) {
			t.Errorf("OddsFromProbability(%v) = (%v, %v), want %v", c.p, got, err, c.want)
		}
	}
}

func TestOddsFromProbability_errors(t *testing.T) {
	_, err := OddsFromProbability(1)
	requireErr(t, "p = 1", err, ErrOutOfRange)
	_, err = OddsFromProbability(nan)
	requireErr(t, "NaN", err, ErrNotFinite)
	_, err = OddsFromProbability(-0.5)
	requireErr(t, "p < 0", err, ErrProbability)
}

// Round trip with ProbabilityFromOdds over a fixed-seed sample; tolerance 4 ulps.
func TestOddsFromProbability_roundTrip(t *testing.T) {
	rng := newTestRand()
	for range 500 {
		p := rng.Float64()
		o, err := OddsFromProbability(p)
		if err != nil {
			t.Fatal(err)
		}
		if back, err := ProbabilityFromOdds(o); err != nil || !within(back, p, 1e-15) {
			t.Fatalf("ProbabilityFromOdds(OddsFromProbability(%v)) = (%v, %v)", p, back, err)
		}
	}
}

func TestOddsFromProbability_noAlloc(t *testing.T) {
	requireNoAlloc(t, "OddsFromProbability", func() { sinkFloat, _ = OddsFromProbability(0.75) })
}

func BenchmarkOddsFromProbability(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = OddsFromProbability(0.75)
	}
}

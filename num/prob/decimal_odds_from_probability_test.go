package prob

import "testing"

// Wheatcroft 2021 §2.1 (rᵢ = 1/Oᵢ) and Wunderlich–Memmert 2018 (no margin: inverse odds = probability):
// p = 0.4 → 2.5. Hand-checked.
func TestDecimalOddsFromProbability_sourceExample(t *testing.T) {
	if got, err := DecimalOddsFromProbability(0.4); err != nil || got != 2.5 {
		t.Errorf("DecimalOddsFromProbability(0.4) = (%v, %v), want (2.5, nil)", got, err)
	}
}

func TestDecimalOddsFromProbability_cases(t *testing.T) {
	for _, c := range []struct{ p, want float64 }{{1, 1}, {0.5, 2}, {0.25, 4}, {1e-300, 1e300}} {
		if got, err := DecimalOddsFromProbability(c.p); err != nil || !within(got, c.want, 5e-16) {
			t.Errorf("DecimalOddsFromProbability(%v) = (%v, %v), want %v", c.p, got, err, c.want)
		}
	}
}

func TestDecimalOddsFromProbability_errors(t *testing.T) {
	_, err := DecimalOddsFromProbability(0)
	requireErr(t, "p = 0", err, ErrOutOfRange)
	_, err = DecimalOddsFromProbability(5e-324)
	requireErr(t, "subnormal p", err, ErrOverflow)
	_, err = DecimalOddsFromProbability(nan)
	requireErr(t, "NaN", err, ErrNotFinite)
	_, err = DecimalOddsFromProbability(1.5)
	requireErr(t, "p > 1", err, ErrProbability)
}

func TestDecimalOddsFromProbability_noAlloc(t *testing.T) {
	requireNoAlloc(t, "DecimalOddsFromProbability", func() { sinkFloat, _ = DecimalOddsFromProbability(0.4) })
}

func BenchmarkDecimalOddsFromProbability(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = DecimalOddsFromProbability(0.4)
	}
}

package prob

import "testing"

// I_0.5(2, 3) = P(Binomial(4, 0.5) ≥ 2) = 11/16 (DLMF 8.17.5 relation, hand-checked); I_x(a, 1) = x^a.
func TestBetaCDF_sourceExample(t *testing.T) {
	cases := []struct{ x, a, b, want float64 }{
		{0.5, 2, 3, 0.6875},
		{0.3, 3, 1, 0.027},
		{0.3, 1, 1, 0.3},
	}
	for _, c := range cases {
		if got, err := BetaCDF(c.x, c.a, c.b); err != nil || !within(got, c.want, 1e-13) {
			t.Errorf("BetaCDF(%v, %v, %v) = (%v, %v), want %v", c.x, c.a, c.b, got, err, c.want)
		}
	}
}

func TestBetaCDF_cases(t *testing.T) {
	if got, _ := BetaCDF(-1, 2, 3); got != 0 {
		t.Errorf("below: %v", got)
	}
	if got, _ := BetaCDF(1, 2, 3); got != 1 {
		t.Errorf("at 1: %v", got)
	}
}

func TestBetaCDF_errors(t *testing.T) {
	_, err := BetaCDF(0.5, 2, -1)
	requireErr(t, "β < 0", err, ErrOutOfRange)
	_, err = BetaCDF(inf, 2, 3)
	requireErr(t, "x Inf", err, ErrNotFinite)
}

func BenchmarkBetaCDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = BetaCDF(0.3, 2, 3)
	}
}

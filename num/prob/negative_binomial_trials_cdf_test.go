package prob

import "testing"

// Casella–Berger §3.2 (negative binomial, trials): r = 2, p = 0.5, x = 4 → 0.25 + 0.25 + 0.1875 = 0.6875.
// Hand-checked; tolerance 1e-14.
func TestNegativeBinomialTrialsCDF_sourceExample(t *testing.T) {
	if got, err := NegativeBinomialTrialsCDF(4, 2, 0.5); err != nil || !within(got, 0.6875, 1e-14) {
		t.Errorf("NegativeBinomialTrialsCDF(4, 2, 0.5) = (%v, %v), want (0.6875, nil)", got, err)
	}
}

func TestNegativeBinomialTrialsCDF_cases(t *testing.T) {
	tests := []struct {
		name string
		x, r int
		p    float64
		want float64
	}{
		{"x < r", 1, 2, 0.5, 0},
		{"x = r", 2, 2, 0.5, 0.25},
		{"r = 1 is geometric", 3, 1, 0.2, 0.488},
		{"p = 1", 5, 3, 1, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := NegativeBinomialTrialsCDF(tt.x, tt.r, tt.p); err != nil || !within(got, tt.want, 1e-14) {
				t.Errorf("NegativeBinomialTrialsCDF(%d, %d, %v) = (%v, %v), want %v", tt.x, tt.r, tt.p, got, err, tt.want)
			}
		})
	}
}

func TestNegativeBinomialTrialsCDF_errors(t *testing.T) {
	_, err := NegativeBinomialTrialsCDF(4, 0, 0.5)
	requireErr(t, "r < 1", err, ErrOutOfRange)
	_, err = NegativeBinomialTrialsCDF(4, 2, 0)
	requireErr(t, "p = 0", err, ErrOutOfRange)
	_, err = NegativeBinomialTrialsCDF(4, 2, inf)
	requireErr(t, "Inf", err, ErrNotFinite)
}

// Oracle: the CDF equals the compensated sum of the PMF over r…x. Tolerance 2e-12 relative plus 1e-15.
func TestNegativeBinomialTrialsCDF_oracle(t *testing.T) {
	rng := newTestRand()
	for trial := range 300 {
		r := 1 + rng.IntN(20)
		x := r + rng.IntN(150)
		p := 0.05 + 0.95*rng.Float64()
		var s compensatedSum
		for i := r; i <= x; i++ {
			v, _ := NegativeBinomialTrialsPMF(i, r, p)
			s.add(v)
		}
		want := s.total()
		if got, err := NegativeBinomialTrialsCDF(x, r, p); err != nil || absDiff(got, want) > 2e-12*want+1e-15 {
			t.Fatalf("trial %d: NegativeBinomialTrialsCDF(%d, %d, %v) = (%v, %v), want %v", trial, x, r, p, got, err, want)
		}
	}
}

func TestNegativeBinomialTrialsCDF_noAlloc(t *testing.T) {
	requireNoAlloc(t, "NegativeBinomialTrialsCDF", func() { sinkFloat, _ = NegativeBinomialTrialsCDF(40, 10, 0.3) })
}

func BenchmarkNegativeBinomialTrialsCDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = NegativeBinomialTrialsCDF(40, 10, 0.3)
	}
}

package prob

import "testing"

// Deeks & Altman, BMJ 2004;329:168, table: ≥ 40 pack-years in 42 of 148 with the disease and 2 of 144
// without, so sensitivity 42/148, specificity 142/144 and LR+ = (42/148)/(2/144) = 20.43 (printed 20.4).
// Tolerance 1e-14 relative (three roundings).
func TestLikelihoodRatioPositive_sourceExample(t *testing.T) {
	got, err := LikelihoodRatioPositive(42.0/148, 142.0/144)
	if want := (42.0 / 148) / (2.0 / 144); err != nil || !within(got, want, 1e-14) {
		t.Errorf("LikelihoodRatioPositive = (%v, %v), want %v", got, err, want)
	}
	if !within(got, 20.4, 0.003) {
		t.Errorf("LikelihoodRatioPositive = %v, want 20.4 as printed", got)
	}
}

func TestLikelihoodRatioPositive_cases(t *testing.T) {
	tests := []struct {
		name   string
		se, sp float64
		want   float64
	}{
		{"uninformative test", 0.7, 0.3, 1},
		{"zero sensitivity", 0, 0.9, 0},
		{"zero specificity", 0.8, 0, 0.8},
		{"perfect sensitivity", 1, 0.95, 20},
	}
	for _, tc := range tests {
		got, err := LikelihoodRatioPositive(tc.se, tc.sp)
		if err != nil || !within(got, tc.want, 1e-14) {
			t.Errorf("%s: LikelihoodRatioPositive(%v, %v) = (%v, %v), want %v", tc.name, tc.se, tc.sp, got, err, tc.want)
		}
	}
}

func TestLikelihoodRatioPositive_errors(t *testing.T) {
	_, err := LikelihoodRatioPositive(nan, 0.5)
	requireErr(t, "LR+(NaN, 0.5)", err, ErrNotFinite)
	_, err = LikelihoodRatioPositive(0.5, inf)
	requireErr(t, "LR+(0.5, Inf)", err, ErrNotFinite)
	_, err = LikelihoodRatioPositive(1.2, 0.5)
	requireErr(t, "LR+(1.2, 0.5)", err, ErrProbability)
	_, err = LikelihoodRatioPositive(0.5, -0.1)
	requireErr(t, "LR+(0.5, -0.1)", err, ErrProbability)
	_, err = LikelihoodRatioPositive(0.9, 1)
	requireErr(t, "LR+(0.9, 1)", err, ErrOutOfRange)
}

func TestLikelihoodRatioPositive_noAlloc(t *testing.T) {
	requireNoAlloc(t, "LikelihoodRatioPositive", func() { sinkFloat, _ = LikelihoodRatioPositive(0.9, 0.8) })
}

func BenchmarkLikelihoodRatioPositive(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = LikelihoodRatioPositive(0.9, 0.8)
	}
}

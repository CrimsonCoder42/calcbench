package prob

import "testing"

// Deeks & Altman, BMJ 2004;329:168, same table dichotomised at 40 pack-years: LR− = (1 − 42/148)/(142/144)
// = (106/148)/(142/144) ≈ 0.7263 (hand-checked). Tolerance 1e-14 relative.
func TestLikelihoodRatioNegative_sourceExample(t *testing.T) {
	got, err := LikelihoodRatioNegative(42.0/148, 142.0/144)
	if want := (106.0 / 148) / (142.0 / 144); err != nil || !within(got, want, 1e-14) {
		t.Errorf("LikelihoodRatioNegative = (%v, %v), want %v", got, err, want)
	}
}

func TestLikelihoodRatioNegative_cases(t *testing.T) {
	tests := []struct {
		name   string
		se, sp float64
		want   float64
	}{
		{"uninformative test", 0.4, 0.6, 1},
		{"perfect sensitivity", 1, 0.7, 0},
		{"perfect specificity", 0.8, 1, 0.2},
		{"zero sensitivity", 0, 0.5, 2},
	}
	for _, tc := range tests {
		got, err := LikelihoodRatioNegative(tc.se, tc.sp)
		if err != nil || !within(got, tc.want, 1e-14) {
			t.Errorf("%s: LikelihoodRatioNegative(%v, %v) = (%v, %v), want %v", tc.name, tc.se, tc.sp, got, err, tc.want)
		}
	}
}

func TestLikelihoodRatioNegative_errors(t *testing.T) {
	_, err := LikelihoodRatioNegative(nan, 0.5)
	requireErr(t, "LR−(NaN, 0.5)", err, ErrNotFinite)
	_, err = LikelihoodRatioNegative(0.5, 1.5)
	requireErr(t, "LR−(0.5, 1.5)", err, ErrProbability)
	_, err = LikelihoodRatioNegative(0.5, 0)
	requireErr(t, "LR−(0.5, 0)", err, ErrOutOfRange)
}

func TestLikelihoodRatioNegative_noAlloc(t *testing.T) {
	requireNoAlloc(t, "LikelihoodRatioNegative", func() { sinkFloat, _ = LikelihoodRatioNegative(0.9, 0.8) })
}

func BenchmarkLikelihoodRatioNegative(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = LikelihoodRatioNegative(0.9, 0.8)
	}
}

package prob

import "testing"

// Ross §3.3 (Bayes's formula) on a 2×2 table TP 231, FP 32, FN 27, TN 54 (n = 344):
// Se = 231/258, Sp = 54/86, prevalence 258/344 → PPV = TP/(TP + FP) = 231/263 ≈ 0.878.
// The exact 231/263 is hand-checked with fractions; tolerance 8 ulps for the rounded inputs.
func TestPositivePredictiveValue_sourceExample(t *testing.T) {
	got, err := PositivePredictiveValue(231.0/258, 54.0/86, 258.0/344)
	if err != nil || !within(got, 231.0/263, 2e-15) {
		t.Errorf("PositivePredictiveValue(2×2 table) = (%v, %v), want (%v, nil)", got, err, 231.0/263)
	}
}

func TestPositivePredictiveValue_cases(t *testing.T) {
	tests := []struct {
		name             string
		se, sp, pi, want float64
	}{
		{"perfect specificity", 0.9, 1, 0.1, 1},
		{"zero prevalence", 0.9, 0.9, 0, 0},
		{"coin-flip test", 0.5, 0.5, 0.2, 0.2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := PositivePredictiveValue(tt.se, tt.sp, tt.pi); err != nil || !within(got, tt.want, 1e-15) {
				t.Errorf("PPV(%v, %v, %v) = (%v, %v), want %v", tt.se, tt.sp, tt.pi, got, err, tt.want)
			}
		})
	}
}

func TestPositivePredictiveValue_errors(t *testing.T) {
	_, err := PositivePredictiveValue(0, 1, 0.5)
	requireErr(t, "no positive possible", err, ErrOutOfRange)
	_, err = PositivePredictiveValue(nan, 0.5, 0.5)
	requireErr(t, "NaN", err, ErrNotFinite)
	_, err = PositivePredictiveValue(0.5, 1.5, 0.5)
	requireErr(t, "specificity > 1", err, ErrProbability)
	_, err = PositivePredictiveValue(0.5, 0.5, -0.1)
	requireErr(t, "prevalence < 0", err, ErrProbability)
}

func TestPositivePredictiveValue_noAlloc(t *testing.T) {
	requireNoAlloc(t, "PositivePredictiveValue", func() { sinkFloat, _ = PositivePredictiveValue(0.9, 0.6, 0.75) })
}

func BenchmarkPositivePredictiveValue(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = PositivePredictiveValue(0.9, 0.6, 0.75)
	}
}

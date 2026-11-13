package prob

import (
	"math"
	"math/big"
	"testing"
)

var (
	dieFaces = []float64{1, 2, 3, 4, 5, 6}
	dieProbs = []float64{1.0 / 6, 1.0 / 6, 1.0 / 6, 1.0 / 6, 1.0 / 6, 1.0 / 6}
)

// Ross §4.3 (expected value): a fair die → 7/2 = 3.5. Hand-checked; tolerance 2 ulps (1/6 is rounded).
func TestExpectedValue_sourceExample(t *testing.T) {
	if got, err := ExpectedValue(dieFaces, dieProbs); err != nil || !within(got, 3.5, 5e-16) {
		t.Errorf("ExpectedValue(die) = (%v, %v), want (3.5, nil)", got, err)
	}
}

func TestExpectedValue_cases(t *testing.T) {
	tests := []struct {
		name   string
		xs, ps []float64
		want   float64
	}{
		{"one point", []float64{-7}, []float64{1}, -7},
		{"bet: win 10 with 0.25, lose 5", []float64{10, -5}, []float64{0.25, 0.75}, -1.25},
		{"zero probability outcome", []float64{1e300, 2}, []float64{0, 1}, 2},
		{"near overflow, scaled retry", []float64{math.MaxFloat64, math.MaxFloat64}, []float64{0.5, 0.5}, math.MaxFloat64},
		{"opposite extremes", []float64{-math.MaxFloat64, math.MaxFloat64}, []float64{0.5, 0.5}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := ExpectedValue(tt.xs, tt.ps); err != nil || !within(got, tt.want, 5e-16) {
				t.Errorf("ExpectedValue(%v, %v) = (%v, %v), want %v", tt.xs, tt.ps, got, err, tt.want)
			}
		})
	}
}

func TestExpectedValue_errors(t *testing.T) {
	_, err := ExpectedValue(nil, nil)
	requireErr(t, "empty", err, ErrOutOfRange)
	_, err = ExpectedValue([]float64{1, 2}, []float64{1})
	requireErr(t, "length mismatch", err, ErrOutOfRange)
	_, err = ExpectedValue([]float64{inf}, []float64{1})
	requireErr(t, "Inf outcome", err, ErrNotFinite)
	_, err = ExpectedValue([]float64{1}, []float64{nan})
	requireErr(t, "NaN probability", err, ErrNotFinite)
	_, err = ExpectedValue([]float64{1, 2}, []float64{-0.5, 1.5})
	requireErr(t, "probability < 0", err, ErrProbability)
	_, err = ExpectedValue([]float64{1, 2}, []float64{0.5, 0.4})
	requireErr(t, "sum 0.9", err, ErrNotDistribution)
}

// Oracle: the exact rational Σ xᵢ·pᵢ of the float64 inputs. Tolerance 2 ulps of max|xᵢ| (compensated sum
// of exact-error products).
func TestExpectedValue_oracle(t *testing.T) {
	rng := newTestRand()
	for trial := range 300 {
		n := 1 + rng.IntN(40)
		xs, ps := make([]float64, n), make([]float64, n)
		exact := new(big.Rat)
		maxAbs := 0.0
		for i := range n {
			xs[i] = (rng.Float64() - 0.5) * 1e6
			ps[i] = 1 / float64(n)
			maxAbs = math.Max(maxAbs, math.Abs(xs[i]))
			exact.Add(exact, new(big.Rat).Mul(new(big.Rat).SetFloat64(xs[i]), new(big.Rat).SetFloat64(ps[i])))
		}
		want, _ := exact.Float64()
		if got, err := ExpectedValue(xs, ps); err != nil || absDiff(got, want) > 5e-16*maxAbs {
			t.Fatalf("trial %d: ExpectedValue = (%v, %v), want %v", trial, got, err, want)
		}
	}
}

func TestExpectedValue_noAlloc(t *testing.T) {
	requireNoAlloc(t, "ExpectedValue", func() { sinkFloat, _ = ExpectedValue(dieFaces, dieProbs) })
}

func BenchmarkExpectedValue(b *testing.B) {
	for _, size := range []struct {
		name string
		n    int
	}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}} {
		xs, ps := make([]float64, size.n), make([]float64, size.n)
		for i := range size.n {
			xs[i], ps[i] = float64(i), 1/float64(size.n)
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, _ = ExpectedValue(xs, ps)
			}
		})
	}
}

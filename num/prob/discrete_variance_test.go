package prob

import (
	"math"
	"math/big"
	"testing"
)

// Ross §4.5 (variance): a fair die → 35/12 ≈ 2.9166666666666665. Hand-checked; tolerance 4 ulps.
func TestDiscreteVariance_sourceExample(t *testing.T) {
	if got, err := DiscreteVariance(dieFaces, dieProbs); err != nil || !within(got, 35.0/12, 1e-15) {
		t.Errorf("DiscreteVariance(die) = (%v, %v), want (35/12, nil)", got, err)
	}
}

func TestDiscreteVariance_cases(t *testing.T) {
	tests := []struct {
		name   string
		xs, ps []float64
		want   float64
	}{
		{"one point", []float64{5}, []float64{1}, 0},
		{"bernoulli 0.3", []float64{0, 1}, []float64{0.7, 0.3}, 0.21},
		{"large offset, no cancellation", []float64{1e9 + 1, 1e9 + 3}, []float64{0.5, 0.5}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := DiscreteVariance(tt.xs, tt.ps); err != nil || !within(got, tt.want, 1e-15) {
				t.Errorf("DiscreteVariance(%v, %v) = (%v, %v), want %v", tt.xs, tt.ps, got, err, tt.want)
			}
		})
	}
}

// Chan, Golub, LeVeque 1983 eq. (1.7): the correction term removes the error of the rounded mean. Two
// outcomes 1e9 + 0.1 and 1e9 + 0.3 (as float64), p = ½ each: the exact variance, from big.Rat on the float64
// inputs, is ((x₂ − x₁)/2)² = 0.0099999928474439…; without the correction the result was
// 0.0099999928474475 (3.6e-13 relative). Tolerance 2 ulps relative.
func TestDiscreteVariance_correctionTerm(t *testing.T) {
	xs, ps := []float64{1e9 + 0.1, 1e9 + 0.3}, []float64{0.5, 0.5}
	half := new(big.Rat).SetFrac64(1, 2)
	mean := new(big.Rat)
	for _, x := range xs {
		mean.Add(mean, new(big.Rat).Mul(new(big.Rat).SetFloat64(x), half))
	}
	exact := new(big.Rat)
	for _, x := range xs {
		d := new(big.Rat).Sub(new(big.Rat).SetFloat64(x), mean)
		exact.Add(exact, d.Mul(d.Mul(d, d), half))
	}
	want, _ := exact.Float64()
	if got, err := DiscreteVariance(xs, ps); err != nil || !within(got, want, 4.5e-16) {
		t.Errorf("DiscreteVariance(1e9+0.1, 1e9+0.3) = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestDiscreteVariance_errors(t *testing.T) {
	_, err := DiscreteVariance([]float64{-math.MaxFloat64, math.MaxFloat64}, []float64{0.5, 0.5})
	requireErr(t, "variance overflows", err, ErrOverflow)
	_, err = DiscreteVariance(nil, nil)
	requireErr(t, "empty", err, ErrOutOfRange)
	_, err = DiscreteVariance([]float64{nan}, []float64{1})
	requireErr(t, "NaN", err, ErrNotFinite)
	_, err = DiscreteVariance([]float64{1, 2}, []float64{0.5, 0.6})
	requireErr(t, "sum 1.1", err, ErrNotDistribution)
}

// Oracle: the exact rational Σ pᵢ(xᵢ − μ)² with μ the exact mean. Tolerance 1e-13 relative (two passes,
// the second over rounded deviations).
func TestDiscreteVariance_oracle(t *testing.T) {
	rng := newTestRand()
	for trial := range 300 {
		n := 2 + rng.IntN(40)
		xs, ps := make([]float64, n), make([]float64, n)
		mean := new(big.Rat)
		for i := range n {
			xs[i] = 1000 + rng.Float64()
			ps[i] = 1 / float64(n)
			mean.Add(mean, new(big.Rat).Mul(new(big.Rat).SetFloat64(xs[i]), new(big.Rat).SetFloat64(ps[i])))
		}
		exact := new(big.Rat)
		for i := range n {
			d := new(big.Rat).Sub(new(big.Rat).SetFloat64(xs[i]), mean)
			exact.Add(exact, d.Mul(d.Mul(d, d), new(big.Rat).SetFloat64(ps[i])))
		}
		want, _ := exact.Float64()
		if got, err := DiscreteVariance(xs, ps); err != nil || !within(got, want, 1e-13) {
			t.Fatalf("trial %d: DiscreteVariance = (%v, %v), want %v", trial, got, err, want)
		}
	}
}

func TestDiscreteVariance_noAlloc(t *testing.T) {
	requireNoAlloc(t, "DiscreteVariance", func() { sinkFloat, _ = DiscreteVariance(dieFaces, dieProbs) })
}

func BenchmarkDiscreteVariance(b *testing.B) {
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
				sinkFloat, _ = DiscreteVariance(xs, ps)
			}
		})
	}
}

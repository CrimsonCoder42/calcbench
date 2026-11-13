package prob

import (
	"math"
	"testing"
)

func requireSliceClose(t *testing.T, name string, got, want []float64, tol float64) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s: len %d, want %d", name, len(got), len(want))
	}
	for i := range want {
		if math.Abs(got[i]-want[i]) > tol*math.Max(1, math.Abs(want[i])) {
			t.Errorf("%s[%d] = %v, want %v", name, i, got[i], want[i])
		}
	}
}

// JKB 1997 ch. 35 §2: E = n·p, var = n·p(1 − p), cov = −n·pᵢpⱼ.
func TestMultinomialMeanVariance(t *testing.T) {
	mean, cov, err := MultinomialMeanVariance(nil, nil, 10, []float64{0.2, 0.3, 0.5})
	if err != nil {
		t.Fatal(err)
	}
	requireSliceClose(t, "mean", mean, []float64{2, 3, 5}, 1e-15)
	requireSliceClose(t, "cov", cov, []float64{1.6, -0.6, -1, -0.6, 2.1, -1.5, -1, -1.5, 2.5}, 1e-15)
	// Each row of the covariance sums to 0 (Σ Xᵢ = n is constant).
	for i := range 3 {
		if s := cov[3*i] + cov[3*i+1] + cov[3*i+2]; math.Abs(s) > 1e-14 {
			t.Errorf("row %d sums to %v", i, s)
		}
	}
	mean, cov, _ = MultinomialMeanVariance(mean, cov, 0, []float64{0.5, 0.5})
	requireSliceClose(t, "n = 0 mean", mean, []float64{0, 0}, 0)
	requireSliceClose(t, "n = 0 cov", cov, []float64{0, 0, 0, 0}, 0)
	_, _, err = MultinomialMeanVariance(nil, nil, -1, []float64{1})
	requireErr(t, "n < 0", err, ErrOutOfRange)
	_, _, err = MultinomialMeanVariance(nil, nil, 3, nil)
	requireErr(t, "empty", err, ErrOutOfRange)
	_, _, err = MultinomialMeanVariance(nil, nil, 3, []float64{0.5, 0.6})
	requireErr(t, "not a distribution", err, ErrNotDistribution)
	md, cd := make([]float64, 3), make([]float64, 9)
	p := []float64{0.2, 0.3, 0.5}
	requireNoAlloc(t, "MultinomialMeanVariance", func() { sinkSlice, _, _ = MultinomialMeanVariance(md, cd, 10, p) })
}

// KBJ vol. 1 ch. 49 §2: α = (1, 2, 2), α₀ = 5: mean (0.2, 0.4, 0.4), var α̃(1 − α̃)/6, cov −α̃ᵢα̃ⱼ/6.
func TestDirichletMeanVariance(t *testing.T) {
	mean, cov, err := DirichletMeanVariance(nil, nil, []float64{1, 2, 2})
	if err != nil {
		t.Fatal(err)
	}
	requireSliceClose(t, "mean", mean, []float64{0.2, 0.4, 0.4}, 1e-15)
	requireSliceClose(t, "cov", cov, []float64{
		0.16 / 6, -0.08 / 6, -0.08 / 6,
		-0.08 / 6, 0.24 / 6, -0.16 / 6,
		-0.08 / 6, -0.16 / 6, 0.24 / 6,
	}, 1e-14)
	// Beta(2, 3) as a 2-component Dirichlet: var = αβ/((α + β)²(α + β + 1)) = 6/150 = 0.04.
	_, cov, _ = DirichletMeanVariance(nil, nil, []float64{2, 3})
	if !within(cov[0], 0.04, 1e-14) {
		t.Errorf("Beta(2, 3) variance %v, want 0.04", cov[0])
	}
	_, _, err = DirichletMeanVariance(nil, nil, []float64{1})
	requireErr(t, "m = 1", err, ErrOutOfRange)
	_, _, err = DirichletMeanVariance(nil, nil, []float64{1, 0})
	requireErr(t, "α = 0", err, ErrOutOfRange)
	_, _, err = DirichletMeanVariance(nil, nil, []float64{1, math.NaN()})
	requireErr(t, "α NaN", err, ErrNotFinite)
	_, _, err = DirichletMeanVariance(nil, nil, []float64{math.MaxFloat64, math.MaxFloat64})
	requireErr(t, "α₀ overflow", err, ErrOverflow)
	md, cd := make([]float64, 3), make([]float64, 9)
	a := []float64{1, 2, 2}
	requireNoAlloc(t, "DirichletMeanVariance", func() { sinkSlice, _, _ = DirichletMeanVariance(md, cd, a) })
}

func BenchmarkMultinomialMeanVariance(b *testing.B) {
	b.ReportAllocs()
	md, cd := make([]float64, 4), make([]float64, 16)
	p := []float64{0.1, 0.2, 0.3, 0.4}
	for b.Loop() {
		sinkSlice, _, _ = MultinomialMeanVariance(md, cd, 100, p)
	}
}

func BenchmarkDirichletMeanVariance(b *testing.B) {
	b.ReportAllocs()
	md, cd := make([]float64, 4), make([]float64, 16)
	a := []float64{1, 2, 3, 4}
	for b.Loop() {
		sinkSlice, _, _ = DirichletMeanVariance(md, cd, a)
	}
}

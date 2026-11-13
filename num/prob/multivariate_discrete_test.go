package prob

import "testing"

// JKB 1997 ch. 35 (multinomial PMF), hand-computed: 4!/(1!2!1!)·0.2·0.5²·0.3 = 0.18.
func TestMultinomialPMF(t *testing.T) {
	if got, err := MultinomialPMF([]int{1, 2, 1}, []float64{0.2, 0.5, 0.3}); err != nil || !within(got, 0.18, 1e-14) {
		t.Errorf("source example = (%v, %v)", got, err)
	}
	// Two categories reduce to the binomial.
	for k := 0; k <= 10; k++ {
		got, _ := MultinomialPMF([]int{k, 10 - k}, []float64{0.3, 0.7})
		want, _ := BinomialPMF(k, 10, 0.3)
		if !within(got, want, 1e-13) {
			t.Errorf("k = %d: %v, binomial %v", k, got, want)
		}
	}
	// Sum over all outcomes of 4 trials in 3 categories is 1.
	sum := 0.0
	for a := 0; a <= 4; a++ {
		for b := 0; a+b <= 4; b++ {
			v, _ := MultinomialPMF([]int{a, b, 4 - a - b}, []float64{0.2, 0.5, 0.3})
			sum += v
		}
	}
	if !within(sum, 1, 1e-14) {
		t.Errorf("Σ = %v", sum)
	}
	if got, _ := MultinomialPMF([]int{1, 0}, []float64{0, 1}); got != 0 {
		t.Errorf("p = 0 with k > 0 = %v", got)
	}
	if got, _ := MultinomialPMF([]int{0, 0}, []float64{0.5, 0.5}); got != 1 {
		t.Errorf("n = 0 = %v", got)
	}
	if got, _ := MultinomialPMF([]int{1000000, 1000000}, []float64{0.5, 0.5}); !(got > 0 && got < 1e-3) {
		t.Errorf("large n = %v", got)
	}
	_, err := MultinomialPMF([]int{1}, []float64{0.5, 0.5})
	requireErr(t, "lengths", err, ErrOutOfRange)
	_, err = MultinomialPMF(nil, nil)
	requireErr(t, "empty", err, ErrOutOfRange)
	_, err = MultinomialPMF([]int{-1, 1}, []float64{0.5, 0.5})
	requireErr(t, "negative count", err, ErrOutOfRange)
	_, err = MultinomialPMF([]int{1 << 53, 1}, []float64{0.5, 0.5})
	requireErr(t, "n too large", err, ErrOutOfRange)
	_, err = MultinomialPMF([]int{1, 1}, []float64{0.5, 0.6})
	requireErr(t, "not distribution", err, ErrNotDistribution)
	_, err = MultinomialPMF([]int{1, 1}, []float64{nan, 0.5})
	requireErr(t, "NaN", err, ErrNotFinite)
	_, err = MultinomialPMF([]int{1, 1}, []float64{-0.5, 1.5})
	requireErr(t, "p outside", err, ErrProbability)
	c, p := []int{1, 2, 1}, []float64{0.2, 0.5, 0.3}
	requireNoAlloc(t, "MultinomialPMF", func() { sinkFloat, _ = MultinomialPMF(c, p) })
}

// KBJ vol. 1 ch. 49 (Dirichlet density), hand-computed: Γ(6)/Γ(2)³·0.2·0.3·0.5 = 3.6.
func TestDirichletPDF(t *testing.T) {
	if got, err := DirichletPDF([]float64{0.2, 0.3, 0.5}, []float64{2, 2, 2}); err != nil || !within(got, 3.6, 1e-14) {
		t.Errorf("source example = (%v, %v)", got, err)
	}
	// Two components reduce to the beta density.
	for _, x := range []float64{0.1, 0.4, 0.75} {
		got, _ := DirichletPDF([]float64{x, 1 - x}, []float64{2.5, 4})
		want, _ := BetaPDF(x, 2.5, 4)
		if !within(got, want, 1e-13) {
			t.Errorf("x = %v: %v, beta %v", x, got, want)
		}
	}
	// α = (1, 1, 1) is uniform on the simplex: density Γ(3) = 2 everywhere, faces included.
	if got, _ := DirichletPDF([]float64{0, 0.25, 0.75}, []float64{1, 1, 1}); !within(got, 2, 1e-15) {
		t.Errorf("uniform = %v", got)
	}
	if got, _ := DirichletPDF([]float64{0, 0.5, 0.5}, []float64{2, 1, 1}); got != 0 {
		t.Errorf("face with α > 1 = %v", got)
	}
	_, err := DirichletPDF([]float64{0, 1}, []float64{0.5, 1})
	requireErr(t, "face with α < 1", err, ErrOverflow)
	_, err = DirichletPDF([]float64{1}, []float64{1})
	requireErr(t, "dimension 1", err, ErrOutOfRange)
	_, err = DirichletPDF([]float64{0.5, 0.5}, []float64{1})
	requireErr(t, "lengths", err, ErrOutOfRange)
	_, err = DirichletPDF([]float64{0.5, 0.5}, []float64{1, 0})
	requireErr(t, "α = 0", err, ErrOutOfRange)
	_, err = DirichletPDF([]float64{0.5, 0.5}, []float64{1, inf})
	requireErr(t, "α Inf", err, ErrNotFinite)
	_, err = DirichletPDF([]float64{0.5, 0.6}, []float64{1, 1})
	requireErr(t, "not on simplex", err, ErrNotDistribution)
	_, err = DirichletPDF([]float64{1.5, -0.5}, []float64{1, 1})
	requireErr(t, "x outside", err, ErrProbability)
	_, err = DirichletPDF([]float64{5e-324, 1}, []float64{0.01, 1})
	requireErr(t, "overflow", err, ErrOverflow)
	x, a := []float64{0.2, 0.3, 0.5}, []float64{2, 2, 2}
	requireNoAlloc(t, "DirichletPDF", func() { sinkFloat, _ = DirichletPDF(x, a) })
}

// JKB 1997 ch. 39 (multivariate hypergeometric), hand-computed: C(5,2)C(10,2)C(15,2)/C(30,6) = 47250/593775.
func TestMultivariateHypergeometricPMF(t *testing.T) {
	if got, err := MultivariateHypergeometricPMF([]int{2, 2, 2}, []int{5, 10, 15}); err != nil || !within(got, 47250.0/593775, 1e-13) {
		t.Errorf("source example = (%v, %v)", got, err)
	}
	// Two colours reduce to the hypergeometric PMF.
	for k := 0; k <= 5; k++ {
		got, _ := MultivariateHypergeometricPMF([]int{k, 5 - k}, []int{8, 12})
		want, _ := HypergeometricPMF(k, 20, 8, 5)
		if !within(got, want, 1e-12) {
			t.Errorf("k = %d: %v, hypergeometric %v", k, got, want)
		}
	}
	if got, _ := MultivariateHypergeometricPMF([]int{3, 0}, []int{2, 5}); got != 0 {
		t.Errorf("k > K = %v", got)
	}
	_, err := MultivariateHypergeometricPMF([]int{1}, []int{1, 2})
	requireErr(t, "lengths", err, ErrOutOfRange)
	_, err = MultivariateHypergeometricPMF([]int{-1, 1}, []int{1, 2})
	requireErr(t, "negative", err, ErrOutOfRange)
	_, err = MultivariateHypergeometricPMF([]int{1, 1}, []int{1 << 53, 1})
	requireErr(t, "too many items", err, ErrOutOfRange)
	d, c := []int{2, 2, 2}, []int{5, 10, 15}
	requireNoAlloc(t, "MultivariateHypergeometricPMF", func() { sinkFloat, _ = MultivariateHypergeometricPMF(d, c) })
}

func BenchmarkMultinomialPMF(b *testing.B) {
	c, p := []int{1, 2, 1}, []float64{0.2, 0.5, 0.3}
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = MultinomialPMF(c, p)
	}
}

func BenchmarkDirichletPDF(b *testing.B) {
	x, a := []float64{0.2, 0.3, 0.5}, []float64{2, 2, 2}
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = DirichletPDF(x, a)
	}
}

func BenchmarkMultivariateHypergeometricPMF(b *testing.B) {
	d, c := []int{2, 2, 2}, []int{5, 10, 15}
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = MultivariateHypergeometricPMF(d, c)
	}
}

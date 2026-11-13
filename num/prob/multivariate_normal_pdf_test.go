package prob

import (
	"math"
	"testing"
)

// KBJ ch. 45 §1 density, hand-evaluated: Σ = [[4,2],[2,3]], det 8, Σ⁻¹ = [[3,−2],[−2,4]]/8.
var mvnL2 = []float64{2, 0, 1, math.Sqrt2}

func TestMultivariateNormalPDF_sourceExample(t *testing.T) {
	cases := []struct {
		x, mu []float64
		want  float64
	}{
		{[]float64{1, 1}, []float64{0, 0}, math.Exp(-3.0/16) / (2 * math.Pi * math.Sqrt(8))},
		{[]float64{0, 0}, []float64{0, 0}, 1 / (2 * math.Pi * math.Sqrt(8))},
		{[]float64{3, -1}, []float64{1, 1}, math.Exp(-0.5*(3*4-2*2*2*(-2)+4*4)/8) / (2 * math.Pi * math.Sqrt(8))},
	}
	for _, c := range cases {
		got, err := MultivariateNormalPDF(c.x, c.mu, mvnL2)
		if err != nil || !within(got, c.want, 1e-14) {
			t.Errorf("PDF(%v; %v) = (%v, %v), want %v", c.x, c.mu, got, err, c.want)
		}
		lg, err := MultivariateNormalLogPDF(c.x, c.mu, mvnL2)
		if err != nil || !within(lg, math.Log(c.want), 1e-14) {
			t.Errorf("LogPDF(%v) = (%v, %v), want %v", c.x, lg, err, math.Log(c.want))
		}
	}
}

// n = 1 is NormalPDF; a diagonal L is the product of independent normal densities.
func TestMultivariateNormalPDF_reduces(t *testing.T) {
	got, _ := MultivariateNormalPDF([]float64{1.3}, []float64{0.2}, []float64{1.7})
	want, _ := NormalPDF(1.3, 0.2, 1.7)
	if !within(got, want, 1e-14) {
		t.Errorf("n = 1: %v, want %v", got, want)
	}
	L := []float64{2, 0, 0, 0, 0.5, 0, 0, 0, 3}
	got, _ = MultivariateNormalPDF([]float64{1, 2, 3}, []float64{0, 1, -1}, L)
	a, _ := NormalPDF(1, 0, 2)
	b, _ := NormalPDF(2, 1, 0.5)
	c, _ := NormalPDF(3, -1, 3)
	if !within(got, a*b*c, 1e-14) {
		t.Errorf("diagonal: %v, want %v", got, a*b*c)
	}
	// Upper triangle is ignored.
	got2, _ := MultivariateNormalPDF([]float64{1, 1}, []float64{0, 0}, []float64{2, 99, 1, math.Sqrt2})
	got1, _ := MultivariateNormalPDF([]float64{1, 1}, []float64{0, 0}, mvnL2)
	if got1 != got2 {
		t.Errorf("upper triangle read: %v vs %v", got1, got2)
	}
	// Far tail: PDF underflows to 0, LogPDF stays finite.
	if d, err := MultivariateNormalPDF([]float64{100}, []float64{0}, []float64{1}); err != nil || d != 0 {
		t.Errorf("tail PDF = (%v, %v)", d, err)
	}
	if lg, err := MultivariateNormalLogPDF([]float64{100}, []float64{0}, []float64{1}); err != nil || !within(lg, -5000-0.5*math.Log(2*math.Pi), 1e-15) {
		t.Errorf("tail LogPDF = (%v, %v)", lg, err)
	}
	// Above the stack scratch size.
	n := 20
	big := make([]float64, n*n)
	x := make([]float64, n)
	mu := make([]float64, n)
	for i := range n {
		big[i*n+i] = 1
		x[i] = 0.1
	}
	lg, _ := MultivariateNormalLogPDF(x, mu, big)
	if want := -0.5 * (float64(n)*math.Log(2*math.Pi) + float64(n)*0.01); !within(lg, want, 1e-14) {
		t.Errorf("n = 20: %v, want %v", lg, want)
	}
}

func TestMultivariateNormalPDF_errors(t *testing.T) {
	_, err := MultivariateNormalPDF(nil, nil, nil)
	requireErr(t, "empty", err, ErrDimension)
	_, err = MultivariateNormalPDF([]float64{1, 1}, []float64{0}, mvnL2)
	requireErr(t, "mean length", err, ErrDimension)
	_, err = MultivariateNormalLogPDF([]float64{1, 1}, []float64{0, 0}, []float64{1, 0, 1})
	requireErr(t, "L length", err, ErrDimension)
	_, err = MultivariateNormalPDF([]float64{nan, 1}, []float64{0, 0}, mvnL2)
	requireErr(t, "x NaN", err, ErrNotFinite)
	_, err = MultivariateNormalPDF([]float64{1, 1}, []float64{0, inf}, mvnL2)
	requireErr(t, "mean Inf", err, ErrNotFinite)
	_, err = MultivariateNormalPDF([]float64{1, 1}, []float64{0, 0}, []float64{2, 0, nan, 1})
	requireErr(t, "L NaN", err, ErrNotFinite)
	_, err = MultivariateNormalPDF([]float64{1, 1}, []float64{0, 0}, []float64{2, 0, 1, 0})
	requireErr(t, "L diag 0", err, ErrOutOfRange)
	_, err = MultivariateNormalLogPDF([]float64{1, 1}, []float64{0, 0}, []float64{-2, 0, 1, 1})
	requireErr(t, "L diag < 0", err, ErrOutOfRange)
	_, err = MultivariateNormalPDF([]float64{1e300, 0}, []float64{-1e300, 0}, []float64{1e-300, 0, 0, 1})
	requireErr(t, "quad overflow", err, ErrOverflow)
	_, err = MultivariateNormalPDF([]float64{0}, []float64{0}, []float64{1e-310})
	requireErr(t, "density overflow", err, ErrOverflow)
}

func TestMultivariateNormalPDF_noAlloc(t *testing.T) {
	x, mu := []float64{1, 1}, []float64{0, 0}
	requireNoAlloc(t, "MultivariateNormalPDF", func() { sinkFloat, _ = MultivariateNormalPDF(x, mu, mvnL2) })
	requireNoAlloc(t, "MultivariateNormalLogPDF", func() { sinkFloat, _ = MultivariateNormalLogPDF(x, mu, mvnL2) })
}

func BenchmarkMultivariateNormalPDF(b *testing.B) {
	L := []float64{2, 0, 0, 1, 1.5, 0, 0.5, -0.3, 1.2}
	x, mu := []float64{1, 2, 3}, []float64{0, 1, 2}
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = MultivariateNormalPDF(x, mu, L)
	}
}

func BenchmarkMultivariateNormalLogPDF(b *testing.B) {
	L := []float64{2, 0, 0, 1, 1.5, 0, 0.5, -0.3, 1.2}
	x, mu := []float64{1, 2, 3}, []float64{0, 1, 2}
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = MultivariateNormalLogPDF(x, mu, L)
	}
}

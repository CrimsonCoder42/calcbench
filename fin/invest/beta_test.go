package invest

import (
	"fmt"
	"math"
	"testing"
)

func TestBeta_handWorked(t *testing.T) {
	m := []float64{0.01, 0.02, -0.01, 0.03}
	a := make([]float64, len(m))
	for i, x := range m {
		a[i] = 2*x + 0.001
	}
	got, err := Beta(a, m)
	if err != nil || !near(got, 2, 1e-13) {
		t.Errorf("Beta = %v, %v; want 2", got, err)
	}
}

// Hand-worked: market 1, 2, 3 (s² = 1), asset 2, 1, 6 (s_xy = Σ(x − 2)(y − 3)/2 = (1 + 0 + 3)/2 = 2) → 2.
func TestBeta_cases(t *testing.T) {
	got, err := Beta([]float64{2, 1, 6}, []float64{1, 2, 3})
	if err != nil || !near(got, 2, 1e-15) {
		t.Errorf("Beta = %v, %v; want 2", got, err)
	}
	got, err = Beta([]float64{0.05, 0.05, 0.05}, []float64{0.01, 0.02, 0.03})
	if err != nil || !near(got, 0, 1e-15) {
		t.Errorf("constant asset: %v, %v; want 0", got, err)
	}
	got, err = Beta([]float64{-0.01, -0.02, -0.03}, []float64{0.01, 0.02, 0.03})
	if err != nil || !near(got, -1, 1e-15) {
		t.Errorf("mirror: %v, %v; want -1", got, err)
	}
}

func TestBeta_errors(t *testing.T) {
	_, err := Beta([]float64{0.01}, []float64{0.01})
	wantErr(t, "one pair", err, ErrTooFew)
	_, err = Beta([]float64{0.01, 0.02}, []float64{0.01})
	wantErr(t, "lengths", err, ErrLengthMismatch)
	_, err = Beta([]float64{0.01, math.NaN()}, []float64{0.01, 0.02})
	wantErr(t, "NaN", err, ErrNotFinite)
	_, err = Beta([]float64{0.01, 0.02}, []float64{0.03, 0.03})
	wantErr(t, "flat market", err, ErrZeroDenominator)
	_, err = Beta([]float64{-1.7e308, 1.7e308}, []float64{-1e-10, 1e-10})
	wantErr(t, "overflow", err, ErrOverflow)
}

func TestBeta_noAlloc(t *testing.T) {
	a, m := []float64{2, 1, 6}, []float64{1, 2, 3}
	if n := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = Beta(a, m) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkBeta(b *testing.B) {
	for _, n := range benchSizes {
		m := genReturns(n)
		a := genReturns(n + 1)[1:]
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, sinkErr = Beta(a, m)
			}
		})
	}
}

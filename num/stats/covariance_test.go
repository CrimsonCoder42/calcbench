package stats

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"testing"
)

// nistCov is the e-Handbook §6.5.4.1 example: five observations of three variables.
var nistCov = [][]float64{
	{4.0, 4.2, 3.9, 4.3, 4.1},
	{2.0, 2.1, 2.0, 2.1, 2.2},
	{0.60, 0.59, 0.58, 0.62, 0.63},
}

// Printed S (e-Handbook §6.5.4.1), checked with exact rational arithmetic.
var nistS = []float64{
	0.025, 0.0075, 0.00175,
	0.0075, 0.007, 0.00135,
	0.00175, 0.00135, 0.00043,
}

func TestSampleCovariance_nist(t *testing.T) {
	// Tolerance 1e-13 relative: the inputs are decimal fractions not exact in binary.
	if got, err := SampleCovariance(nistCov[0], nistCov[1]); err != nil || !near(got, 0.0075, 1e-13) {
		t.Errorf("SampleCovariance = %v, %v; want 0.0075", got, err)
	}
}

func TestSampleCovariance_cases(t *testing.T) {
	m := math.MaxFloat64
	for _, tt := range []struct {
		name string
		x, y []float64
		want float64
	}{
		{"two", []float64{0, 1}, []float64{5, 7}, 1},
		{"hand-checked", []float64{1, 2, 3, 4}, []float64{2, 1, 4, 3}, 1},
		{"negative", []float64{1, 2, 3}, []float64{3, 2, 1}, -1},
		{"constant y", []float64{1, 2, 3}, []float64{7, 7, 7}, 0},
		{"huge values", []float64{-m / 2, m / 2}, []float64{-1, 1}, m},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := SampleCovariance(tt.x, tt.y); err != nil || !near(got, tt.want, 1e-15) {
				t.Errorf("SampleCovariance(%v, %v) = %v, %v; want %v", tt.x, tt.y, got, err, tt.want)
			}
		})
	}
}

func TestPopulationCovariance_cases(t *testing.T) {
	for _, tt := range []struct {
		name string
		x, y []float64
		want float64
	}{
		{"one pair", []float64{3}, []float64{4}, 0},
		{"hand-checked", []float64{1, 2, 3, 4}, []float64{2, 1, 4, 3}, 0.75},
		{"self is population variance", []float64{2, 4, 4, 4, 5, 5, 7, 9}, []float64{2, 4, 4, 4, 5, 5, 7, 9}, 4},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := PopulationCovariance(tt.x, tt.y); err != nil || !near(got, tt.want, 1e-15) {
				t.Errorf("PopulationCovariance(%v, %v) = %v, %v; want %v", tt.x, tt.y, got, err, tt.want)
			}
		})
	}
}

// Oracle: the covariance of x with itself is the variance.
func TestSampleCovariance_selfIsVariance(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for range 200 {
		x := randomFloats(r, 2+r.IntN(50))
		c, err1 := SampleCovariance(x, x)
		v, err2 := SampleVariance(x)
		// Tolerance 1e-14 relative: same centred sums, different code paths.
		if err1 != nil || err2 != nil || !near(c, v, 1e-14) {
			t.Fatalf("SampleCovariance(x, x) = %v (%v), SampleVariance = %v (%v)", c, err1, v, err2)
		}
	}
}

func TestCovariance_errors(t *testing.T) {
	m := math.MaxFloat64
	for _, tt := range []struct {
		name string
		x, y []float64
		want error
		pop  bool
	}{
		{"empty", nil, nil, ErrEmpty, false},
		{"empty pop", nil, nil, ErrEmpty, true},
		{"mismatch", []float64{1, 2}, []float64{1}, ErrLengthMismatch, false},
		{"mismatch pop", []float64{1}, []float64{1, 2}, ErrLengthMismatch, true},
		{"one pair", []float64{1}, []float64{1}, ErrTooFew, false},
		{"NaN", []float64{1, math.NaN()}, []float64{1, 2}, ErrNotFinite, false},
		{"Inf pop one", []float64{math.Inf(1)}, []float64{1}, ErrNotFinite, true},
		{"Inf y pop", []float64{1, 2}, []float64{1, math.Inf(-1)}, ErrNotFinite, true},
		{"overflow", []float64{-m, m}, []float64{-m, m}, ErrOverflow, false},
		{"overflow pop", []float64{-m, m}, []float64{-m, m}, ErrOverflow, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := SampleCovariance
			if tt.pop {
				f = PopulationCovariance
			}
			if _, err := f(tt.x, tt.y); !errors.Is(err, tt.want) {
				t.Errorf("err = %v; want %v", err, tt.want)
			}
		})
	}
}

func TestCovarianceMatrix_nist(t *testing.T) {
	got, err := CovarianceMatrix(nistCov)
	if err != nil {
		t.Fatal(err)
	}
	for i := range nistS {
		// Tolerance 1e-12 relative: decimal inputs; the smallest entry is 4.3e-4.
		if !near(got[i], nistS[i], 1e-12) {
			t.Errorf("S[%d] = %v; want %v", i, got[i], nistS[i])
		}
	}
}

func TestCovarianceMatrix_errors(t *testing.T) {
	m := math.MaxFloat64
	for _, tt := range []struct {
		name string
		in   [][]float64
		want error
	}{
		{"no variables", nil, ErrEmpty},
		{"no observations", [][]float64{{}}, ErrEmpty},
		{"mismatch", [][]float64{{1, 2}, {1}}, ErrLengthMismatch},
		{"one observation", [][]float64{{1}, {2}}, ErrTooFew},
		{"NaN", [][]float64{{1, 2}, {1, math.NaN()}}, ErrNotFinite},
		{"overflow", [][]float64{{-m, m}}, ErrOverflow},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := CovarianceMatrix(tt.in); !errors.Is(err, tt.want) {
				t.Errorf("err = %v; want %v", err, tt.want)
			}
		})
	}
}

func BenchmarkSampleCovariance(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		in := benchInputs(n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = SampleCovariance(in[0], in[1])
			}
		})
	}
}

func BenchmarkPopulationCovariance(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		in := benchInputs(n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = PopulationCovariance(in[0], in[1])
			}
		})
	}
}

func BenchmarkCovarianceMatrix(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		in := benchInputs(n)
		b.Run(fmt.Sprintf("k=8/%d", n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = CovarianceMatrix(in[:])
			}
		})
	}
}

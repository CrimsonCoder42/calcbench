package stats

import (
	"errors"
	"fmt"
	"math"
	"testing"
)

// Hand-checked: r_xy = 0.6 (as in the PearsonCorrelation case), z reverses x.
func TestCorrelationMatrix_cases(t *testing.T) {
	in := [][]float64{{1, 2, 3, 4}, {2, 1, 4, 3}, {4, 3, 2, 1}}
	want := []float64{1, 0.6, -1, 0.6, 1, -0.6, -1, -0.6, 1}
	got, err := CorrelationMatrix(in)
	if err != nil {
		t.Fatal(err)
	}
	for i := range want {
		if !near(got[i], want[i], 1e-15) {
			t.Errorf("R[%d] = %v; want %v", i, got[i], want[i])
		}
	}
}

// Agrees with PearsonCorrelation pairwise on the e-Handbook §6.5.4.1 data.
func TestCorrelationMatrix_matchesPearson(t *testing.T) {
	got, err := CorrelationMatrix(nistCov)
	if err != nil {
		t.Fatal(err)
	}
	for i := range 3 {
		for j := range 3 {
			want := 1.0
			if i != j {
				want, _ = PearsonCorrelation(nistCov[i], nistCov[j])
			}
			if got[i*3+j] != want {
				t.Errorf("R[%d,%d] = %v; want %v", i, j, got[i*3+j], want)
			}
		}
	}
}

func TestCorrelationMatrix_errors(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   [][]float64
		want error
	}{
		{"no variables", nil, ErrEmpty},
		{"mismatch", [][]float64{{1, 2}, {1}}, ErrLengthMismatch},
		{"one observation", [][]float64{{1}}, ErrTooFew},
		{"Inf", [][]float64{{1, math.Inf(1)}}, ErrNotFinite},
		{"constant variable", [][]float64{{1, 2, 3}, {5, 5, 5}}, ErrZeroVariance},
		{"constant first", [][]float64{{5, 5}}, ErrZeroVariance},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := CorrelationMatrix(tt.in); !errors.Is(err, tt.want) {
				t.Errorf("err = %v; want %v", err, tt.want)
			}
		})
	}
}

func BenchmarkCorrelationMatrix(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		in := benchInputs(n)
		b.Run(fmt.Sprintf("k=8/%d", n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = CorrelationMatrix(in[:])
			}
		})
	}
}

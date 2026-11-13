package stats

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Worked example (Scott 2015 §3.4 bin counts ν_kl, hand-checked).
func TestFrequencyDistribution2D_worked(t *testing.T) {
	x := []float64{0.5, 1.5, 1.5, 2}
	y := []float64{0.5, 0.5, 1.5, 2}
	e := []float64{0, 1, 2}
	got, err := FrequencyDistribution2D(x, y, e, e)
	if want := []int{1, 0, 1, 2}; err != nil || !slices.Equal(got, want) {
		t.Errorf("FrequencyDistribution2D = %v, %v; want %v", got, err, want)
	}
}

// Oracle: the row sums equal FrequencyDistribution of x, the column sums that of y.
func TestFrequencyDistribution2D_margins(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	xe, ye := []float64{0, 0.3, 0.5, 1}, []float64{0, 0.25, 0.5, 0.75, 1}
	for range 100 {
		n := 1 + r.IntN(50)
		x, y := make([]float64, n), make([]float64, n)
		for i := range x {
			x[i], y[i] = r.Float64(), r.Float64()
		}
		got, err := FrequencyDistribution2D(x, y, xe, ye)
		fx, _ := FrequencyDistribution(x, xe)
		fy, _ := FrequencyDistribution(y, ye)
		if err != nil {
			t.Fatal(err)
		}
		bx, by := len(xe)-1, len(ye)-1
		for i := range bx {
			s := 0
			for j := range by {
				s += got[i*by+j]
			}
			if s != fx[i] {
				t.Fatalf("row %d sum %d; want %d", i, s, fx[i])
			}
		}
		for j := range by {
			s := 0
			for i := range bx {
				s += got[i*by+j]
			}
			if s != fy[j] {
				t.Fatalf("column %d sum %d; want %d", j, s, fy[j])
			}
		}
	}
}

func TestFrequencyDistribution2D_errors(t *testing.T) {
	e := []float64{0, 1}
	for _, tt := range []struct {
		name         string
		x, y, xe, ye []float64
		want         error
	}{
		{"empty", nil, nil, e, e, ErrEmpty},
		{"mismatch", []float64{0.5}, nil, e, e, ErrLengthMismatch},
		{"NaN x", []float64{math.NaN()}, []float64{0.5}, e, e, ErrNotFinite},
		{"Inf y", []float64{0.5}, []float64{math.Inf(1)}, e, e, ErrNotFinite},
		{"one x edge", []float64{0.5}, []float64{0.5}, []float64{0}, e, ErrTooFew},
		{"y edges descending", []float64{0.5}, []float64{0.5}, e, []float64{1, 0}, ErrNotAscending},
		{"x outside", []float64{2}, []float64{0.5}, e, e, ErrOutOfRange},
		{"y outside", []float64{0.5}, []float64{-1}, e, e, ErrOutOfRange},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := FrequencyDistribution2D(tt.x, tt.y, tt.xe, tt.ye); !errors.Is(err, tt.want) {
				t.Errorf("err = %v; want %v", err, tt.want)
			}
		})
	}
}

func BenchmarkFrequencyDistribution2D(b *testing.B) {
	edges := make([]float64, 21)
	for i := range edges {
		edges[i] = float64(i) * 10 // 0 … 200
	}
	for _, n := range []int{1_000, 10_000, 100_000} {
		in := benchInputs(n)
		for _, v := range in[:2] {
			for i := range v {
				v[i] = min(max(v[i], 0), 200)
			}
		}
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = FrequencyDistribution2D(in[0], in[1], edges, edges)
			}
		})
	}
}

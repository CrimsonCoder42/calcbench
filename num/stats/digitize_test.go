package stats

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Worked example (Knuth §6.2.1 insertion point, hand-checked).
func TestDigitize_worked(t *testing.T) {
	got, err := Digitize([]float64{-1, 0, 0.5, 1, 3}, []float64{0, 1, 2})
	if want := []int{0, 1, 1, 2, 3}; err != nil || !slices.Equal(got, want) {
		t.Errorf("Digitize = %v, %v; want %v", got, err, want)
	}
	if got, err := Digitize(nil, []float64{0}); err != nil || len(got) != 0 {
		t.Errorf("Digitize(nil) = %v, %v; want empty", got, err)
	}
	if got, err := Digitize([]float64{-1, 0, 1}, []float64{0}); err != nil || !slices.Equal(got, []int{0, 1, 1}) {
		t.Errorf("one edge = %v, %v", got, err)
	}
}

// Oracle: a linear count of edges ≤ v.
func TestDigitize_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	edges := []float64{-2, -0.5, 0, 0.25, 1, 3}
	for range 300 {
		v := r.Float64()*8 - 4
		want := 0
		for _, e := range edges {
			if e <= v {
				want++
			}
		}
		if got, _ := Digitize([]float64{v}, edges); got[0] != want {
			t.Fatalf("Digitize(%v) = %d; want %d", v, got[0], want)
		}
	}
}

func TestDigitize_errors(t *testing.T) {
	for _, tt := range []struct {
		name          string
		values, edges []float64
		want          error
	}{
		{"no edges", []float64{1}, nil, ErrTooFew},
		{"NaN value", []float64{math.NaN()}, []float64{0}, ErrNotFinite},
		{"Inf edge", []float64{1}, []float64{0, math.Inf(1)}, ErrNotFinite},
		{"repeated edge", []float64{1}, []float64{0, 0}, ErrNotAscending},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Digitize(tt.values, tt.edges); !errors.Is(err, tt.want) {
				t.Errorf("err = %v; want %v", err, tt.want)
			}
		})
	}
}

func BenchmarkDigitize(b *testing.B) {
	edges := make([]float64, 41)
	for i := range edges {
		edges[i] = float64(i) * 5
	}
	for _, n := range []int{1_000, 10_000, 100_000} {
		in := benchInputs(n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = Digitize(in[0], edges)
			}
		})
	}
}

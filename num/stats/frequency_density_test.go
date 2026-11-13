package stats

import (
	"errors"
	"fmt"
	"math"
	"testing"
)

// Worked example (Scott 2015 §3.1 eq. 3.1, hand-checked): counts 1, 3, 1 of n = 5 in widths 2 → 0.1, 0.3, 0.1.
func TestFrequencyDensity_worked(t *testing.T) {
	got, err := FrequencyDensity([]float64{1, 2, 2, 3, 5}, []float64{0, 2, 4, 6})
	want := []float64{0.1, 0.3, 0.1}
	for i := range want {
		if err != nil || !near(got[i], want[i], 1e-15) {
			t.Fatalf("FrequencyDensity = %v, %v; want %v", got, err, want)
		}
	}
}

// Property: Σ density·width = 1 with unequal widths.
func TestFrequencyDensity_integratesToOne(t *testing.T) {
	values := []float64{0.1, 0.2, 0.7, 1.5, 2.5, 3.9, 4, 9}
	edges := []float64{0, 0.5, 1, 4, 10}
	got, err := FrequencyDensity(values, edges)
	if err != nil {
		t.Fatal(err)
	}
	var area float64
	for j, d := range got {
		area += d * (edges[j+1] - edges[j])
	}
	if !near(area, 1, 1e-15) {
		t.Errorf("area = %v; want 1", area)
	}
	m := math.MaxFloat64
	if got, err := FrequencyDensity([]float64{0}, []float64{-m, m}); err != nil || !near(got[0], 0.5/m, 1e-15) {
		t.Errorf("full-range class = %v, %v; want %v", got, err, 0.5/m)
	}
}

func TestFrequencyDensity_errors(t *testing.T) {
	tiny := math.SmallestNonzeroFloat64
	for _, tt := range []struct {
		name          string
		values, edges []float64
		want          error
	}{
		{"empty", nil, []float64{0, 1}, ErrEmpty},
		{"one edge", []float64{1}, []float64{0}, ErrTooFew},
		{"NaN", []float64{math.NaN()}, []float64{0, 1}, ErrNotFinite},
		{"not ascending", []float64{1}, []float64{1, 0}, ErrNotAscending},
		{"outside", []float64{2}, []float64{0, 1}, ErrOutOfRange},
		{"overflow", []float64{0}, []float64{0, tiny}, ErrOverflow},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := FrequencyDensity(tt.values, tt.edges); !errors.Is(err, tt.want) {
				t.Errorf("err = %v; want %v", err, tt.want)
			}
		})
	}
}

func BenchmarkFrequencyDensity(b *testing.B) {
	edges := make([]float64, 41)
	for i := range edges {
		edges[i] = float64(i) * 5 // 0 … 200
	}
	for _, n := range []int{1_000, 10_000, 100_000} {
		in := benchInputs(n)
		v := in[0]
		for i := range v {
			v[i] = min(max(v[i], 0), 200)
		}
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = FrequencyDensity(v, edges)
			}
		})
	}
}

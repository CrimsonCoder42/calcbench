package stats

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Worked example: e-Handbook §1.3.3.14 counts the points in each class. Values 1, 2, 2, 3, 5 with edges
// 0, 2, 4, 6 → [0,2) holds 1, [2,4) holds 2, 2, 3, [4,6] holds 5 → 1, 3, 1 (hand-checked).
func TestFrequencyDistribution_worked(t *testing.T) {
	if got, err := FrequencyDistribution([]float64{1, 2, 2, 3, 5}, []float64{0, 2, 4, 6}); err != nil || !slices.Equal(got, []int{1, 3, 1}) {
		t.Errorf("FrequencyDistribution = %v, %v; want [1 3 1]", got, err)
	}
}

func TestFrequencyDistribution_cases(t *testing.T) {
	for _, tt := range []struct {
		name          string
		values, edges []float64
		want          []int
	}{
		{"one class", []float64{0, 1, 0.5}, []float64{0, 1}, []int{3}},
		{"value on the last edge joins the last class", []float64{6}, []float64{0, 2, 4, 6}, []int{0, 0, 1}},
		{"value on an inner edge goes right", []float64{2, 4}, []float64{0, 2, 4, 6}, []int{0, 1, 1}},
		{"value on the first edge", []float64{0}, []float64{0, 2}, []int{1}},
		{"negative edges", []float64{-1.5, -0.5, -0.5}, []float64{-2, -1, 0}, []int{1, 2}},
		{"empty classes", []float64{9}, []float64{0, 1, 2, 10}, []int{0, 0, 1}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := FrequencyDistribution(tt.values, tt.edges); err != nil || !slices.Equal(got, tt.want) {
				t.Errorf("FrequencyDistribution(%v, %v) = %v, %v; want %v", tt.values, tt.edges, got, err, tt.want)
			}
		})
	}
}

func TestFrequencyDistribution_errors(t *testing.T) {
	for _, tt := range []struct {
		name          string
		values, edges []float64
		want          error
	}{
		{"no values", nil, []float64{0, 1}, ErrEmpty},
		{"one edge", []float64{1}, []float64{0}, ErrTooFew},
		{"NaN value", []float64{math.NaN()}, []float64{0, 1}, ErrNotFinite},
		{"Inf edge", []float64{1}, []float64{0, math.Inf(1)}, ErrNotFinite},
		{"equal edges", []float64{1}, []float64{0, 1, 1}, ErrNotAscending},
		{"descending edges", []float64{1}, []float64{2, 1}, ErrNotAscending},
		{"value below", []float64{-0.1}, []float64{0, 1}, ErrOutOfRange},
		{"value above", []float64{1.1}, []float64{0, 1}, ErrOutOfRange},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := FrequencyDistribution(tt.values, tt.edges); !errors.Is(err, tt.want) || got != nil {
				t.Errorf("FrequencyDistribution = %v, %v; want nil, %v", got, err, tt.want)
			}
		})
	}
}

func TestFrequencyDistribution_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		values := randomFloats(rng, 1+rng.IntN(300))
		b := 1 + rng.IntN(20)
		lo, hi := slices.Min(values), slices.Max(values)
		edges := make([]float64, b+1)
		for j := range edges {
			edges[j] = lo + (hi-lo+1)*float64(j)/float64(b)
		}
		edges[0] = lo
		want := make([]int, b)
		for _, v := range values { // linear scan oracle
			for j := range b {
				if v >= edges[j] && (v < edges[j+1] || j == b-1) {
					want[j]++
					break
				}
			}
		}
		if got, err := FrequencyDistribution(values, edges); err != nil || !slices.Equal(got, want) {
			t.Fatalf("trial %d: FrequencyDistribution = %v, %v; want %v", trial, got, err, want)
		}
	}
}

func TestFrequencyDistribution_allocations(t *testing.T) {
	edges := []float64{299, 299.8, 299.9, 301}
	if a := testing.AllocsPerRun(100, func() { _, _ = FrequencyDistribution(michelso, edges) }); a != 1 {
		t.Errorf("FrequencyDistribution allocates %v times per call, want 1 (the result)", a)
	}
}

func BenchmarkFrequencyDistribution(b *testing.B) {
	edges := make([]float64, 41) // 40 classes covering mean ± 10 SD
	for j := range edges {
		edges[j] = 5 * float64(j)
	}
	for _, size := range benchSizes {
		inputs := benchInputs(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; b.Loop(); i++ {
				counts, _ := FrequencyDistribution(inputs[i&7], edges)
				benchSink = float64(counts[20])
			}
		})
	}
}

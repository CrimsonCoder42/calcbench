package stats

import (
	"errors"
	"math"
	"testing"
)

// Worked example: OpenStax Introductory Statistics 2e §2.5, mean of a frequency table = Σfm/Σf with m
// the class midpoint. Classes [0,10) [10,20) [20,30) with counts 2, 5, 3 → (2·5 + 5·15 + 3·25)/10 = 16
// (hand-checked).
func TestGroupedMean_worked(t *testing.T) {
	if got, err := GroupedMean([]float64{0, 10, 20, 30}, []int{2, 5, 3}); err != nil || got != 16 {
		t.Errorf("GroupedMean = %v, %v; want 16", got, err)
	}
}

func TestGroupedMean_cases(t *testing.T) {
	m := math.MaxFloat64
	for _, tt := range []struct {
		name   string
		edges  []float64
		counts []int
		want   float64
	}{
		{"one class", []float64{2, 4}, []int{7}, 3},
		{"zero counts ignored", []float64{0, 1, 2, 3}, []int{0, 4, 0}, 1.5},
		{"unequal widths", []float64{0, 2, 10}, []int{1, 1}, 3.5},
		{"negative edges", []float64{-4, -2, 0}, []int{1, 3}, -1.5},
		{"products overflow", []float64{0, m}, []int{4}, m / 2},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// Tolerance 1e-15: a few roundings.
			if got, err := GroupedMean(tt.edges, tt.counts); err != nil || !near(got, tt.want, 1e-15) {
				t.Errorf("GroupedMean(%v, %v) = %v, %v; want %v", tt.edges, tt.counts, got, err, tt.want)
			}
		})
	}
}

// groupedErrorCases are the input errors shared by GroupedMean and GroupedMedian.
var groupedErrorCases = []struct {
	name   string
	edges  []float64
	counts []int
	want   error
}{
	{"no counts", []float64{0}, nil, ErrEmpty},
	{"edges too short", []float64{0, 1}, []int{1, 1}, ErrLengthMismatch},
	{"edges too long", []float64{0, 1, 2}, []int{1}, ErrLengthMismatch},
	{"NaN edge", []float64{0, math.NaN()}, []int{1}, ErrNotFinite},
	{"edges not ascending", []float64{0, 2, 1}, []int{1, 1}, ErrNotAscending},
	{"negative count", []float64{0, 1, 2}, []int{1, -1}, ErrNegative},
	{"all counts zero", []float64{0, 1}, []int{0}, ErrZeroTotal},
	{"total overflows", []float64{0, 1, 2}, []int{math.MaxInt, 1}, ErrOverflow},
}

func TestGroupedMean_errors(t *testing.T) {
	for _, tt := range groupedErrorCases {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := GroupedMean(tt.edges, tt.counts); !errors.Is(err, tt.want) || got != 0 {
				t.Errorf("GroupedMean(%v, %v) = %v, %v; want 0, %v", tt.edges, tt.counts, got, err, tt.want)
			}
		})
	}
}

// The oracle expands each class into f copies of its midpoint and takes ArithmeticMean.
func TestGroupedMean_oracle(t *testing.T) {
	edges := []float64{0, 1.5, 4, 4.25, 9, 20}
	counts := []int{3, 0, 7, 2, 5}
	var expanded []float64
	for j, f := range counts {
		for range f {
			expanded = append(expanded, (edges[j]+edges[j+1])/2)
		}
	}
	want, _ := ArithmeticMean(expanded)
	// Tolerance 1e-15: both compensated.
	if got, err := GroupedMean(edges, counts); err != nil || !near(got, want, 1e-15) {
		t.Errorf("GroupedMean = %v, %v; want %v", got, err, want)
	}
}

func TestGroupedMean_noAllocation(t *testing.T) {
	edges, counts := []float64{0, 10, 20, 30}, []int{2, 5, 3}
	if a := testing.AllocsPerRun(100, func() { benchSink, _ = GroupedMean(edges, counts) }); a != 0 {
		t.Errorf("GroupedMean allocates %v times per call, want 0", a)
	}
}

// groupedBenchInput returns b classes of width 1 with counts in 0..96.
func groupedBenchInput(b int) ([]float64, []int) {
	edges := make([]float64, b+1)
	counts := make([]int, b)
	for j := range edges {
		edges[j] = float64(j)
	}
	for j := range counts {
		counts[j] = (j * 37) % 97
	}
	return edges, counts
}

func BenchmarkGroupedMean(b *testing.B) {
	for _, size := range benchSizes {
		edges, counts := groupedBenchInput(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchSink, _ = GroupedMean(edges, counts)
			}
		})
	}
}

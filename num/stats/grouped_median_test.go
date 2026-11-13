package stats

import (
	"errors"
	"testing"
)

// Worked example: Dataplot "GROUPED PERCENTILE", p = l + ((p%·N/100) − F)/f·C at p% = 50. Classes
// [0,10) [10,20) [20,30) with counts 2, 5, 3 (N = 10): the median class is [10,20) (cumulative 2 < 5 ≤
// 7), so 10 + (5 − 2)/5·10 = 16 (hand-checked).
func TestGroupedMedian_worked(t *testing.T) {
	if got, err := GroupedMedian([]float64{0, 10, 20, 30}, []int{2, 5, 3}); err != nil || got != 16 {
		t.Errorf("GroupedMedian = %v, %v; want 16", got, err)
	}
}

func TestGroupedMedian_cases(t *testing.T) {
	for _, tt := range []struct {
		name   string
		edges  []float64
		counts []int
		want   float64
	}{
		{"one class", []float64{2, 4}, []int{8}, 3},
		{"half reached exactly at a class end", []float64{0, 1, 2}, []int{3, 3}, 1},
		{"leading empty class", []float64{0, 1, 2}, []int{0, 4}, 1.5},
		{"odd total", []float64{0, 10, 20}, []int{1, 2}, 10 + 0.5/2*10},
		{"negative edges", []float64{-10, -5, 0}, []int{4, 4}, -5},
		{"last class", []float64{0, 1, 2}, []int{1, 9}, 1 + 4.0/9},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// Tolerance 1e-15: one interpolation.
			if got, err := GroupedMedian(tt.edges, tt.counts); err != nil || !near(got, tt.want, 1e-15) {
				t.Errorf("GroupedMedian(%v, %v) = %v, %v; want %v", tt.edges, tt.counts, got, err, tt.want)
			}
		})
	}
}

func TestGroupedMedian_errors(t *testing.T) {
	for _, tt := range groupedErrorCases {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := GroupedMedian(tt.edges, tt.counts); !errors.Is(err, tt.want) || got != 0 {
				t.Errorf("GroupedMedian(%v, %v) = %v, %v; want 0, %v", tt.edges, tt.counts, got, err, tt.want)
			}
		})
	}
}

func TestGroupedMedian_noAllocation(t *testing.T) {
	edges, counts := []float64{0, 10, 20, 30}, []int{2, 5, 3}
	if a := testing.AllocsPerRun(100, func() { benchSink, _ = GroupedMedian(edges, counts) }); a != 0 {
		t.Errorf("GroupedMedian allocates %v times per call, want 0", a)
	}
}

func BenchmarkGroupedMedian(b *testing.B) {
	for _, size := range benchSizes {
		edges, counts := groupedBenchInput(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchSink, _ = GroupedMedian(edges, counts)
			}
		})
	}
}

package interval

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Worked example written for this repo (the one in the doc comment): 0 1 2 4 5 7 → [0,2] [4,5] [7,7].
func TestSummaryRanges_workedExample(t *testing.T) {
	got, err := SummaryRanges([]int{0, 1, 2, 4, 5, 7})
	if want := []Interval[int]{iv(0, 2), iv(4, 5), iv(7, 7)}; err != nil || !slices.Equal(got, want) {
		t.Errorf("SummaryRanges = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestSummaryRanges_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   []Interval[int]
	}{
		{"empty", nil, []Interval[int]{}},
		{"one", []int{5}, []Interval[int]{iv(5, 5)}},
		{"two consecutive", []int{5, 6}, []Interval[int]{iv(5, 6)}},
		{"two apart", []int{5, 7}, []Interval[int]{iv(5, 5), iv(7, 7)}},
		{"duplicates join", []int{1, 1, 2, 2, 4, 4}, []Interval[int]{iv(1, 2), iv(4, 4)}},
		{"negative through zero", []int{-2, -1, 0, 1}, []Interval[int]{iv(-2, 1)}},
		{"extremes", []int{math.MinInt, math.MinInt + 1, math.MaxInt - 1, math.MaxInt}, []Interval[int]{iv(math.MinInt, math.MinInt+1), iv(math.MaxInt-1, math.MaxInt)}},
		{"far apart, no wrap", []int{math.MinInt, math.MaxInt}, []Interval[int]{iv(math.MinInt, math.MinInt), iv(math.MaxInt, math.MaxInt)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SummaryRanges(tt.values)
			if err != nil || got == nil || !slices.Equal(got, tt.want) {
				t.Errorf("SummaryRanges(%v) = (%#v, %v), want (%v, nil)", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestSummaryRanges_smallTypes(t *testing.T) {
	if got, _ := SummaryRanges([]uint8{0, 254, 255}); !slices.Equal(got, []Interval[uint8]{{0, 0}, {254, 255}}) {
		t.Errorf("uint8 = %v", got)
	}
	if got, _ := SummaryRanges([]int8{-128, 127}); !slices.Equal(got, []Interval[int8]{{-128, -128}, {127, 127}}) {
		t.Errorf("int8 = %v", got)
	}
}

// Problem-site shape, in our own words: page numbers to print, compressed into ranges.
func TestSummaryRanges_problemSiteShape(t *testing.T) {
	got, _ := SummaryRanges([]int{1, 2, 3, 8, 10, 11})
	if want := []Interval[int]{iv(1, 3), iv(8, 8), iv(10, 11)}; !slices.Equal(got, want) {
		t.Errorf("pages = %v, want %v", got, want)
	}
}

func TestSummaryRanges_errors(t *testing.T) {
	for _, values := range [][]int{{2, 1}, {1, 2, 3, 3, 0}} {
		if got, err := SummaryRanges(values); !errors.Is(err, ErrNotSorted) || got != nil {
			t.Errorf("SummaryRanges(%v) = (%v, %v), want ErrNotSorted", values, got, err)
		}
	}
}

// Oracle: SummaryRanges is MergeOverlapping of the unit intervals [v, v+1], with 1 taken off each End
// (the closed form of joining consecutive integers).
func TestSummaryRanges_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		values := make([]int, rng.IntN(40))
		for i := range values {
			values[i] = rng.IntN(60) - 20
		}
		slices.Sort(values)
		units := make([]Interval[int], len(values))
		for i, v := range values {
			units[i] = iv(v, v+1)
		}
		want, _ := MergeOverlapping(units)
		for i := range want {
			want[i].End--
		}
		got, err := SummaryRanges(values)
		if err != nil || !slices.Equal(got, want) {
			t.Fatalf("trial %d: SummaryRanges(%v) = (%v, %v), want %v", trial, values, got, err, want)
		}
	}
}

func BenchmarkSummaryRanges(b *testing.B) {
	for _, size := range benchSizes {
		values := make([]int, size.n)
		for i := range values {
			values[i] = i + i/4 // runs of four
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = SummaryRanges(values)
			}
		})
	}
}

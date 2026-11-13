package interval

import (
	"errors"
	"math"
	"math/bits"
	"math/rand/v2"
	"slices"
	"testing"
)

// Cormen et al., Introduction to Algorithms, 3rd ed., §16.1, Figure 16.1: the 11 activities
// (s, f) = (1,4) (3,5) (0,6) (5,7) (3,9) (5,9) (6,10) (8,11) (8,12) (2,14) (12,16). The book's greedy
// answer is {a1, a4, a8, a11}; with closed intervals the same set is chosen (each next start is strictly
// after the previous finish), 0-based indices 0 3 7 10.
func TestIntervalSchedulingMaxNonOverlap_clrsFigure(t *testing.T) {
	in := []Interval[int]{iv(1, 4), iv(3, 5), iv(0, 6), iv(5, 7), iv(3, 9), iv(5, 9), iv(6, 10), iv(8, 11), iv(8, 12), iv(2, 14), iv(12, 16)}
	got, err := IntervalSchedulingMaxNonOverlap(in)
	if want := []int{0, 3, 7, 10}; err != nil || !slices.Equal(got, want) {
		t.Errorf("CLRS 16.1 = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestIntervalSchedulingMaxNonOverlap_cases(t *testing.T) {
	tests := []struct {
		name string
		in   []Interval[int]
		want []int
	}{
		{"empty", nil, []int{}},
		{"one", []Interval[int]{iv(1, 2)}, []int{0}},
		{"two apart", []Interval[int]{iv(5, 6), iv(1, 2)}, []int{1, 0}},
		{"shared endpoint conflicts (closed)", []Interval[int]{iv(1, 3), iv(3, 5)}, []int{0}},
		{"long one loses to two short", []Interval[int]{iv(0, 10), iv(1, 2), iv(4, 5)}, []int{1, 2}},
		{"duplicates keep the lowest index", []Interval[int]{iv(2, 4), iv(2, 4), iv(2, 4)}, []int{0}},
		{"points", []Interval[int]{iv(3, 3), iv(3, 3), iv(4, 4)}, []int{0, 2}},
		{"negative", []Interval[int]{iv(-5, -3), iv(-2, 0)}, []int{0, 1}},
		{"extremes", []Interval[int]{iv(math.MinInt, 0), iv(1, math.MaxInt)}, []int{0, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IntervalSchedulingMaxNonOverlap(tt.in)
			if err != nil || got == nil || !slices.Equal(got, tt.want) {
				t.Errorf("(%v) = (%#v, %v), want (%v, nil)", tt.in, got, err, tt.want)
			}
		})
	}
}

// Problem-site shape, in our own words: the fewest intervals to delete so the rest do not overlap is
// n minus the size of the greedy set.
func TestIntervalSchedulingMaxNonOverlap_problemSiteShape(t *testing.T) {
	in := []Interval[int]{iv(1, 2), iv(2, 3), iv(3, 4), iv(1, 3)}
	got, _ := IntervalSchedulingMaxNonOverlap(in)
	// Closed: [1,2] and [3,4] are disjoint; [2,3] touches both. Keep 2, delete 2.
	if len(in)-len(got) != 2 {
		t.Errorf("kept %v, want 2 kept", got)
	}
}

func TestIntervalSchedulingMaxNonOverlap_errors(t *testing.T) {
	if got, err := IntervalSchedulingMaxNonOverlap([]Interval[int]{iv(1, 2), iv(7, 3)}); !errors.Is(err, ErrInverted) || got != nil {
		t.Errorf("inverted: (%v, %v), want ErrInverted", got, err)
	}
}

// Oracle: the largest pairwise-disjoint subset by trying every subset (n <= 12); the result must be
// pairwise disjoint and that large.
func TestIntervalSchedulingMaxNonOverlap_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		in := make([]Interval[int], rng.IntN(13))
		for i := range in {
			start := rng.IntN(30)
			in[i] = iv(start, start+rng.IntN(8))
		}
		disjoint := func(idx []int) bool {
			for x := range idx {
				for y := x + 1; y < len(idx); y++ {
					p, q := in[idx[x]], in[idx[y]]
					if p.Start <= q.End && q.Start <= p.End {
						return false
					}
				}
			}
			return true
		}
		best := 0
		for mask := uint(0); mask < 1<<len(in); mask++ {
			var idx []int
			for i := range in {
				if mask&(1<<i) != 0 {
					idx = append(idx, i)
				}
			}
			if bits.OnesCount(mask) > best && disjoint(idx) {
				best = bits.OnesCount(mask)
			}
		}
		got, err := IntervalSchedulingMaxNonOverlap(in)
		if err != nil || len(got) != best || !disjoint(got) {
			t.Fatalf("trial %d: (%v) = (%v, %v), want %d disjoint", trial, in, got, err, best)
		}
	}
}

func BenchmarkIntervalSchedulingMaxNonOverlap(b *testing.B) {
	for _, size := range benchSizes {
		rng := rand.New(rand.NewPCG(1, 2))
		in := make([]Interval[int], size.n)
		for i := range in {
			start := rng.IntN(10 * size.n)
			in[i] = iv(start, start+rng.IntN(50))
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = IntervalSchedulingMaxNonOverlap(in)
			}
		})
	}
}

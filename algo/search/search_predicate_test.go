package search

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// Knuth, TAOCP Vol. 3, 2nd ed., §6.2.1 binary search, phrased as a predicate over indexes: the first
// index of knuthSorted whose key is >= 653 is 10 (the 11th key).
func TestSearchPredicate_knuthExample(t *testing.T) {
	got, err := SearchPredicate(0, len(knuthSorted), func(i int) bool { return knuthSorted[i] >= 653 })
	if err != nil || got != 10 {
		t.Errorf("SearchPredicate(knuth >= 653) = (%d, %v), want (10, nil)", got, err)
	}
}

// Binary search on the answer, our own examples: the smallest x with x·x >= 50 is 8; the smallest
// daily capacity that ships loads 3 2 2 4 1 4 within 3 days is 6.
func TestSearchPredicate_onTheAnswer(t *testing.T) {
	if got, _ := SearchPredicate(0, 100, func(x int) bool { return x*x >= 50 }); got != 8 {
		t.Errorf("smallest x with x² >= 50 = %d, want 8", got)
	}
	loads := []int{3, 2, 2, 4, 1, 4}
	daysNeeded := func(capacity int) int {
		days, carried := 1, 0
		for _, w := range loads {
			if carried+w > capacity {
				days++
				carried = 0
			}
			carried += w
		}
		return days
	}
	if got, _ := SearchPredicate(4, 17, func(c int) bool { return daysNeeded(c) <= 3 }); got != 6 {
		t.Errorf("smallest capacity for 3 days = %d, want 6", got)
	}
}

func TestSearchPredicate_cases(t *testing.T) {
	tests := []struct {
		name      string
		lo, hi    int
		threshold int // pred(x) = x >= threshold
		want      int
	}{
		{"empty range", 5, 5, 0, 5},
		{"one value, true", 3, 4, 3, 3},
		{"one value, false", 3, 4, 4, 4},
		{"true everywhere", 0, 10, -5, 0},
		{"false everywhere", 0, 10, 20, 10},
		{"negative range", -10, -2, -6, -6},
		{"whole int range", math.MinInt, math.MaxInt, 0, 0},
		{"whole int range, threshold at the top", math.MinInt, math.MaxInt, math.MaxInt - 1, math.MaxInt - 1},
		{"whole int range, threshold at the bottom", math.MinInt, math.MaxInt, math.MinInt, math.MinInt},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			got, err := SearchPredicate(tt.lo, tt.hi, func(x int) bool {
				calls++
				if x < tt.lo || x >= tt.hi {
					t.Fatalf("pred called with %d outside [%d, %d)", x, tt.lo, tt.hi)
				}
				return x >= tt.threshold
			})
			if err != nil || got != tt.want {
				t.Errorf("SearchPredicate(%d, %d, >= %d) = (%d, %v), want %d", tt.lo, tt.hi, tt.threshold, got, err, tt.want)
			}
			if calls > 65 {
				t.Errorf("pred called %d times, want at most 65", calls)
			}
		})
	}
}

func TestSearchPredicate_errors(t *testing.T) {
	if _, err := SearchPredicate(5, 4, func(int) bool { return true }); !errors.Is(err, ErrBounds) {
		t.Errorf("lo > hi: err = %v, want ErrBounds", err)
	}
	if _, err := SearchPredicate(math.MaxInt, math.MinInt, func(int) bool { return true }); !errors.Is(err, ErrBounds) {
		t.Errorf("MaxInt > MinInt: err = %v, want ErrBounds", err)
	}
	if _, err := SearchPredicate(0, 4, nil); !errors.Is(err, ErrNilPredicate) {
		t.Errorf("nil pred: err = %v, want ErrNilPredicate", err)
	}
}

func TestSearchPredicate_nonMonotoneDoesNotPanic(t *testing.T) {
	got, err := SearchPredicate(0, 10, func(x int) bool { return x%3 == 0 })
	if err != nil || got < 0 || got > 10 {
		t.Errorf("non-monotone: (%d, %v), want a value in [0, 10]", got, err)
	}
}

func TestSearchPredicate_noAllocation(t *testing.T) {
	pred := func(i int) bool { return knuthSorted[i] >= 653 }
	if allocs := testing.AllocsPerRun(100, func() { benchSink, _ = SearchPredicate(0, 16, pred) }); allocs != 0 {
		t.Errorf("SearchPredicate allocates %v times per call, want 0", allocs)
	}
}

func TestSearchPredicate_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 500 {
		lo := rng.IntN(200) - 100
		hi := lo + rng.IntN(100)
		threshold := lo - 5 + rng.IntN(hi-lo+10)
		want := hi
		for x := lo; x < hi; x++ {
			if x >= threshold {
				want = x
				break
			}
		}
		if got, err := SearchPredicate(lo, hi, func(x int) bool { return x >= threshold }); err != nil || got != want {
			t.Fatalf("trial %d: SearchPredicate(%d, %d, >= %d) = (%d, %v), want %d", trial, lo, hi, threshold, got, err, want)
		}
	}
}

func BenchmarkSearchPredicate(b *testing.B) {
	for _, size := range benchSizes {
		sorted, targets := benchData(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			sink := 0
			for i := 0; b.Loop(); i++ {
				t := targets[i&1023]
				x, _ := SearchPredicate(0, len(sorted), func(j int) bool { return sorted[j] >= t })
				sink += x
			}
			benchSink = sink
		})
	}
}

package interval

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// randomSortedDisjoint returns a sorted, disjoint list (MergeOverlapping of random intervals in [0, maxValue]).
func randomSortedDisjoint(rng *rand.Rand, count, maxValue int) []Interval[int] {
	in := make([]Interval[int], count)
	for i := range in {
		start := rng.IntN(maxValue)
		in[i] = iv(start, start+rng.IntN(min(6, maxValue-start+1)))
	}
	out, _ := MergeOverlapping(in)
	return out
}

// Worked example written for this repo: [1,2] [5,7] [9,9] [12,15] with [6,10] added. By hand: [6,10]
// overlaps [5,7] and contains [9,9], giving [5,10]; [1,2] and [12,15] are untouched.
func TestInsertInterval_workedExample(t *testing.T) {
	sorted := []Interval[int]{iv(1, 2), iv(5, 7), iv(9, 9), iv(12, 15)}
	got, err := InsertInterval(sorted, iv(6, 10))
	if want := []Interval[int]{iv(1, 2), iv(5, 10), iv(12, 15)}; err != nil || !slices.Equal(got, want) {
		t.Errorf("InsertInterval = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestInsertInterval_cases(t *testing.T) {
	tests := []struct {
		name   string
		sorted []Interval[int]
		added  Interval[int]
		want   []Interval[int]
	}{
		{"into empty", nil, iv(3, 4), []Interval[int]{iv(3, 4)}},
		{"before all", []Interval[int]{iv(5, 6)}, iv(1, 2), []Interval[int]{iv(1, 2), iv(5, 6)}},
		{"after all", []Interval[int]{iv(5, 6)}, iv(8, 9), []Interval[int]{iv(5, 6), iv(8, 9)}},
		{"between", []Interval[int]{iv(1, 2), iv(8, 9)}, iv(4, 5), []Interval[int]{iv(1, 2), iv(4, 5), iv(8, 9)}},
		{"touches left end", []Interval[int]{iv(1, 2), iv(4, 5)}, iv(2, 3), []Interval[int]{iv(1, 3), iv(4, 5)}},
		{"touches both", []Interval[int]{iv(1, 2), iv(4, 5)}, iv(2, 4), []Interval[int]{iv(1, 5)}},
		{"integer neighbour stays apart", []Interval[int]{iv(1, 2)}, iv(3, 3), []Interval[int]{iv(1, 2), iv(3, 3)}},
		{"inside one", []Interval[int]{iv(1, 10)}, iv(3, 4), []Interval[int]{iv(1, 10)}},
		{"covers all", []Interval[int]{iv(2, 3), iv(5, 6)}, iv(0, 9), []Interval[int]{iv(0, 9)}},
		{"negative", []Interval[int]{iv(-9, -7), iv(-2, 0)}, iv(-6, -3), []Interval[int]{iv(-9, -7), iv(-6, -3), iv(-2, 0)}},
		{"extremes", []Interval[int]{iv(math.MinInt, -1)}, iv(-1, math.MaxInt), []Interval[int]{iv(math.MinInt, math.MaxInt)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := InsertInterval(tt.sorted, tt.added)
			if err != nil || !slices.Equal(got, tt.want) {
				t.Errorf("InsertInterval(%v, %v) = (%v, %v), want (%v, nil)", tt.sorted, tt.added, got, err, tt.want)
			}
		})
	}
}

// Problem-site shapes, in our own words: one booking slots into a gap; one booking swallows a run of
// three neighbours; the list starts empty.
func TestInsertInterval_problemSiteShapes(t *testing.T) {
	got, _ := InsertInterval([]Interval[int]{iv(1, 3), iv(7, 9)}, iv(4, 5))
	if want := []Interval[int]{iv(1, 3), iv(4, 5), iv(7, 9)}; !slices.Equal(got, want) {
		t.Errorf("gap = %v, want %v", got, want)
	}
	got, _ = InsertInterval([]Interval[int]{iv(0, 1), iv(3, 4), iv(6, 7), iv(9, 11), iv(14, 15)}, iv(4, 10))
	if want := []Interval[int]{iv(0, 1), iv(3, 11), iv(14, 15)}; !slices.Equal(got, want) {
		t.Errorf("swallow = %v, want %v", got, want)
	}
}

func TestInsertInterval_errors(t *testing.T) {
	if _, err := InsertInterval([]Interval[int]{iv(1, 2)}, iv(5, 4)); !errors.Is(err, ErrInverted) {
		t.Errorf("inverted added: err = %v, want ErrInverted", err)
	}
	if _, err := InsertInterval([]Interval[int]{iv(3, 1)}, iv(5, 6)); !errors.Is(err, ErrInverted) {
		t.Errorf("inverted in list: err = %v, want ErrInverted", err)
	}
	for _, bad := range [][]Interval[int]{{iv(5, 6), iv(1, 2)}, {iv(1, 3), iv(3, 4)}, {iv(1, 5), iv(2, 3)}} {
		if got, err := InsertInterval(bad, iv(10, 11)); !errors.Is(err, ErrNotSortedDisjoint) || got != nil {
			t.Errorf("InsertInterval(%v) = (%v, %v), want ErrNotSortedDisjoint", bad, got, err)
		}
	}
	if _, err := InsertInterval([]Interval[float64]{{1, 2}}, Interval[float64]{1, math.NaN()}); !errors.Is(err, ErrInverted) {
		t.Errorf("[1, NaN]: err = %v, want ErrInverted (NaN sorts first, D-002)", err)
	}
}

// Oracle: InsertInterval equals MergeOverlapping of the list with the new interval appended.
func TestInsertInterval_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		sorted := randomSortedDisjoint(rng, rng.IntN(15), 60)
		before := slices.Clone(sorted)
		start := rng.IntN(60)
		added := iv(start, start+rng.IntN(12))
		want, _ := MergeOverlapping(append(slices.Clone(sorted), added))
		got, err := InsertInterval(sorted, added)
		if err != nil || !slices.Equal(got, want) {
			t.Fatalf("trial %d: InsertInterval(%v, %v) = (%v, %v), want %v", trial, sorted, added, got, err, want)
		}
		if !slices.Equal(sorted, before) {
			t.Fatalf("trial %d: input modified", trial)
		}
	}
}

func BenchmarkInsertInterval(b *testing.B) {
	for _, size := range benchSizes {
		sorted := make([]Interval[int], size.n)
		for i := range sorted {
			sorted[i] = iv(10*i, 10*i+5)
		}
		added := iv(5*size.n, 5*size.n+50) // lands in the middle and merges a few
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = InsertInterval(sorted, added)
			}
		})
	}
}

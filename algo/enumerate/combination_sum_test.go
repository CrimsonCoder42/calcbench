package enumerate

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

func mustCombinationSum(t *testing.T, candidates []int, target int) [][]int {
	t.Helper()
	seq, err := CombinationSum(candidates, target)
	if err != nil {
		t.Fatalf("CombinationSum(%v, %d): %v", candidates, target, err)
	}
	return collect(seq)
}

// Backtracking (Knuth §7.2.2 Algorithm B; Skiena §7.1), worked by hand: parts from {2, 3, 5} summing to
// 8, in non-decreasing order, are 2+2+2+2, 2+3+3, 3+5.
func TestCombinationSum_workedExample(t *testing.T) {
	got := mustCombinationSum(t, []int{5, 3, 2}, 8)
	want := [][]int{{2, 2, 2, 2}, {2, 3, 3}, {3, 5}}
	if !slices.EqualFunc(got, want, slices.Equal) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestCombinationSum_cases(t *testing.T) {
	tests := []struct {
		name       string
		candidates []int
		target     int
		want       [][]int
	}{
		{"target zero", []int{1, 2}, 0, [][]int{{}}},
		{"no candidates", nil, 3, nil},
		{"no candidates, target zero", nil, 0, [][]int{{}}},
		{"unreachable", []int{4, 6}, 7, nil},
		{"one candidate", []int{3}, 9, [][]int{{3, 3, 3}}},
		{"duplicates merged", []int{2, 2, 3}, 5, [][]int{{2, 3}}},
		{"candidate equals target", []int{7}, 7, [][]int{{7}}},
		{"coins 1 and 2 make 4", []int{1, 2}, 4, [][]int{{1, 1, 1, 1}, {1, 1, 2}, {2, 2}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mustCombinationSum(t, tt.candidates, tt.target)
			if len(got) != len(tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
			for i := range got {
				if !slices.Equal(got[i], tt.want[i]) {
					t.Errorf("[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestCombinationSum_errors(t *testing.T) {
	tests := []struct {
		name       string
		candidates []int
		target     int
		want       error
	}{
		{"negative target", []int{1}, -1, ErrNegative},
		{"zero candidate", []int{1, 0}, 3, ErrNonPositive},
		{"negative candidate", []int{-2, 3}, 3, ErrNonPositive},
		{"min int candidate", []int{math.MinInt}, 3, ErrNonPositive},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := CombinationSum(tt.candidates, tt.target); !errors.Is(err, tt.want) {
				t.Errorf("err = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestCombinationSum_hugeTargetEarlyBreak(t *testing.T) {
	seq, err := CombinationSum([]int{math.MaxInt / 2, math.MaxInt}, math.MaxInt)
	if err != nil {
		t.Fatal(err)
	}
	got := collect(seq)
	if len(got) != 1 || !slices.Equal(got[0], []int{math.MaxInt}) {
		t.Errorf("got %v, want [[MaxInt]]", got)
	}
	seq, _ = CombinationSum([]int{1, 2, 3}, 1000)
	count := 0
	for range seq {
		count++
		if count == 3 {
			break
		}
	}
	if count != 3 {
		t.Errorf("count = %d, want 3", count)
	}
}

func TestCombinationSum_inputNotModified(t *testing.T) {
	in := []int{5, 3, 3, 2}
	seq, _ := CombinationSum(in, 10)
	for range seq {
	}
	if !slices.Equal(in, []int{5, 3, 3, 2}) {
		t.Errorf("input modified: %v", in)
	}
}

// countWays is the oracle: the number of multisets of cands (distinct, ascending) summing to target, by
// the classic coin-counting table.
func countWays(cands []int, target int) int {
	ways := make([]int, target+1)
	ways[0] = 1
	for _, c := range cands {
		for s := c; s <= target; s++ {
			ways[s] += ways[s-c]
		}
	}
	return ways[target]
}

func TestCombinationSum_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		cands := make([]int, 1+rng.IntN(5))
		for i := range cands {
			cands[i] = 1 + rng.IntN(9)
		}
		target := rng.IntN(25)
		got := mustCombinationSum(t, cands, target)
		uniq := slices.Compact(slices.Sorted(slices.Values(cands)))
		if len(got) != countWays(uniq, target) {
			t.Fatalf("trial %d: %v → %d: %d combinations, want %d", trial, cands, target, len(got), countWays(uniq, target))
		}
		for i, c := range got {
			sum := 0
			for _, v := range c {
				sum += v
			}
			if sum != target || !slices.IsSorted(c) {
				t.Fatalf("trial %d: %v sums to %d or is not sorted", trial, c, sum)
			}
			if i > 0 && slices.Compare(got[i-1], c) >= 0 {
				t.Fatalf("trial %d: %v does not follow %v", trial, c, got[i-1])
			}
		}
	}
}

func BenchmarkCombinationSum(b *testing.B) {
	for _, target := range []int{20, 40, 60} {
		cands := []int{2, 3, 5, 7}
		b.Run(fmt.Sprint(target), func(b *testing.B) {
			b.ReportAllocs()
			sink := 0
			for b.Loop() {
				seq, _ := CombinationSum(cands, target)
				for c := range seq {
					sink += len(c)
				}
			}
			benchSink = sink
		})
	}
}

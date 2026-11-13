package twopointer

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// The definition of Manacher (1975): a sequence equal to its reversal. Hand-made data.
func TestIsPalindrome_worked(t *testing.T) {
	if !IsPalindrome([]int{3, 1, 4, 1, 3}) {
		t.Error("IsPalindrome(3 1 4 1 3) = false, want true")
	}
	if IsPalindrome([]int{3, 1, 4, 1, 5}) {
		t.Error("IsPalindrome(3 1 4 1 5) = true, want false")
	}
}

func TestIsPalindrome_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   bool
	}{
		{"empty", nil, true},
		{"one", []int{7}, true},
		{"two equal", []int{7, 7}, true},
		{"two different", []int{7, 8}, false},
		{"even length", []int{1, 2, 2, 1}, true},
		{"differs only in the middle pair", []int{1, 2, 3, 1}, false},
		{"sorted", []int{1, 2, 3}, false},
		{"negative and zero", []int{-1, 0, -1}, true},
		{"extremes", []int{math.MinInt, math.MaxInt, math.MinInt}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPalindrome(tt.values); got != tt.want {
				t.Errorf("IsPalindrome(%v) = %v, want %v", tt.values, got, tt.want)
			}
		})
	}
	if !IsPalindrome([]float64{1, math.NaN(), 1}) {
		t.Error("NaN at the centre is never compared: want true")
	}
	if IsPalindrome([]float64{math.NaN(), math.NaN()}) {
		t.Error("NaN != NaN: want false")
	}
	if !IsPalindrome([]string{"ab", "c", "ab"}) {
		t.Error("IsPalindrome(ab c ab) = false, want true")
	}
}

func TestIsPalindrome_noAllocation(t *testing.T) {
	values := make([]int, 1000)
	allocs := testing.AllocsPerRun(20, func() { _ = IsPalindrome(values) })
	if allocs != 0 {
		t.Errorf("IsPalindrome allocates %v times per call, want 0", allocs)
	}
}

// Oracle: compare with a reversed copy.
func TestIsPalindrome_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		half := make([]int, rng.IntN(6))
		for i := range half {
			half[i] = rng.IntN(3)
		}
		values := slices.Clone(half)
		if rng.IntN(2) == 0 {
			values = append(values, rng.IntN(3))
		}
		rev := slices.Clone(half)
		slices.Reverse(rev)
		values = append(values, rev...)
		if rng.IntN(2) == 0 && len(values) > 0 {
			values[rng.IntN(len(values))] = rng.IntN(3)
		}
		before := slices.Clone(values)
		reversed := slices.Clone(values)
		slices.Reverse(reversed)
		if got, want := IsPalindrome(values), slices.Equal(values, reversed); got != want {
			t.Fatalf("trial %d: IsPalindrome(%v) = %v, want %v", trial, values, got, want)
		}
		if !slices.Equal(values, before) {
			t.Fatalf("trial %d: input modified", trial)
		}
	}
}

func BenchmarkIsPalindrome(b *testing.B) {
	for _, size := range benchSizes {
		values := make([]int64, size.n) // all zeros: a palindrome, the full scan
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if IsPalindrome(values) {
					benchSink++
				}
			}
		})
	}
}

package window

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// naiveMinLength is the oracle: every run, sums in int, shortest first then earliest start.
func naiveMinLength(values []uint8, target int) (int, int) {
	for length := 1; length <= len(values); length++ {
		for start := 0; start+length <= len(values); start++ {
			sum := 0
			for _, v := range values[start : start+length] {
				sum += int(v)
			}
			if sum >= target {
				return start, length
			}
		}
	}
	return 0, 0
}

// On the array of Laaksonen, Competitive Programmer's Handbook (2018), §8.1 "Subarray sum", p. 77
// (1 3 2 5 1 1 2 3, target 8): no two neighbours reach 8 and 3 + 2 + 5 = 10 is the first run of three
// that does, so the answer is (1, 3). Hand-derived.
func TestMinLengthSubarraySumAtLeast_laaksonenExample(t *testing.T) {
	start, length, err := MinLengthSubarraySumAtLeast([]int{1, 3, 2, 5, 1, 1, 2, 3}, 8)
	if err != nil || start != 1 || length != 3 {
		t.Errorf("MinLengthSubarraySumAtLeast(laaksonen, 8) = (%d, %d, %v), want (1, 3, nil)", start, length, err)
	}
}

func TestMinLengthSubarraySumAtLeast_cases(t *testing.T) {
	tests := []struct {
		name          string
		values        []int64
		target        int64
		start, length int
	}{
		{"empty", nil, 5, 0, 0},
		{"empty, target zero", nil, 0, 0, 0},
		{"one, reached", []int64{5}, 5, 0, 1},
		{"one, not reached", []int64{4}, 5, 0, 0},
		{"target zero: one value", []int64{3, 4}, 0, 0, 1},
		{"negative target: one value", []int64{3, 4}, -9, 0, 1},
		{"needs everything", []int64{1, 1, 1}, 3, 0, 3},
		{"never reached", []int64{1, 1, 1}, 4, 0, 0},
		{"single large value later", []int64{1, 1, 9, 1}, 9, 2, 1},
		{"tie: earliest start", []int64{2, 3, 3, 2}, 5, 0, 2},
		{"duplicates", []int64{2, 2, 2, 2}, 6, 0, 3},
		{"sorted", []int64{1, 2, 3, 4, 5}, 9, 3, 2},
		{"reversed", []int64{5, 4, 3, 2, 1}, 9, 0, 2},
		{"sum beyond int64 without overflow", []int64{math.MaxInt64, math.MaxInt64, 1}, math.MaxInt64, 0, 1},
		{"target at the maximum", []int64{math.MaxInt64 - 1, 1, 1}, math.MaxInt64, 0, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, length, err := MinLengthSubarraySumAtLeast(tt.values, tt.target)
			if err != nil || start != tt.start || length != tt.length {
				t.Errorf("MinLengthSubarraySumAtLeast(%v, %d) = (%d, %d, %v), want (%d, %d, nil)", tt.values, tt.target, start, length, err, tt.start, tt.length)
			}
		})
	}
	start, length, err := MinLengthSubarraySumAtLeast([]uint8{200, 100, 250}, 255)
	if err != nil || start != 0 || length != 2 {
		t.Errorf("uint8 200 100 250, 255 = (%d, %d, %v), want (0, 2, nil)", start, length, err)
	}
}

func TestMinLengthSubarraySumAtLeast_errors(t *testing.T) {
	for _, values := range [][]int{{1, 0, 2}, {-1}, {3, math.MinInt}} {
		if _, _, err := MinLengthSubarraySumAtLeast(values, 1); !errors.Is(err, ErrNotPositive) {
			t.Errorf("MinLengthSubarraySumAtLeast(%v) error = %v, want ErrNotPositive", values, err)
		}
	}
	if _, _, err := MinLengthSubarraySumAtLeast([]uint{4, 0}, 1); !errors.Is(err, ErrNotPositive) {
		t.Errorf("uint zero: error = %v, want ErrNotPositive", err)
	}
}

func TestMinLengthSubarraySumAtLeast_noAllocation(t *testing.T) {
	values := benchPositive(1000)
	allocs := testing.AllocsPerRun(20, func() { _, _, _ = MinLengthSubarraySumAtLeast(values, 5000) })
	if allocs != 0 {
		t.Errorf("MinLengthSubarraySumAtLeast allocates %v times per call, want 0", allocs)
	}
}

// Oracle: random positive uint8 data and uint8 targets (window sums pass the type's range, which the
// function must never compute), against the O(n³) search.
func TestMinLengthSubarraySumAtLeast_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		values := make([]uint8, rng.IntN(20))
		for i := range values {
			values[i] = uint8(1 + rng.IntN(255))
		}
		before := slices.Clone(values)
		target := uint8(rng.IntN(256))
		wantStart, wantLength := naiveMinLength(values, int(target))
		if target == 0 && len(values) > 0 {
			wantStart, wantLength = 0, 1
		}
		start, length, err := MinLengthSubarraySumAtLeast(values, target)
		if err != nil || start != wantStart || length != wantLength {
			t.Fatalf("trial %d: MinLengthSubarraySumAtLeast(%v, %d) = (%d, %d, %v), want (%d, %d)", trial, values, target, start, length, err, wantStart, wantLength)
		}
		if !slices.Equal(values, before) {
			t.Fatalf("trial %d: input modified", trial)
		}
	}
}

// benchPositive returns n fixed-seed values in [1, 100].
func benchPositive(n int) []int64 {
	rng := rand.New(rand.NewPCG(1, 2))
	s := make([]int64, n)
	for i := range s {
		s[i] = int64(1 + rng.IntN(100))
	}
	return s
}

func BenchmarkMinLengthSubarraySumAtLeast(b *testing.B) {
	for _, size := range benchSizes {
		values := benchPositive(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, length, _ := MinLengthSubarraySumAtLeast(values, 2500)
				benchSinkDeque = length
			}
		})
	}
}

package prefix

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// naiveCountSubarrays is the oracle: every subarray summed directly in int, O(n²).
func naiveCountSubarrays(values []int8, target int) int {
	count := 0
	for lo := range values {
		sum := 0
		for _, v := range values[lo:] {
			sum += int(v)
			if sum == target {
				count++
			}
		}
	}
	return count
}

// Laaksonen, Competitive Programmer's Handbook (2018), §8.1, p. 77: the array 1 3 2 5 1 1 2 3 holds a
// subarray with sum 8 (2 + 5 + 1). Prefix sums 0 1 4 6 11 12 13 15 18 differ by 8 only for 12 − 4, so
// the count is 1 (hand-derived).
func TestCountSubarraysSumEquals_laaksonenExample(t *testing.T) {
	if got := CountSubarraysSumEquals([]int{1, 3, 2, 5, 1, 1, 2, 3}, 8); got != 1 {
		t.Errorf("CountSubarraysSumEquals(laaksonen, 8) = %d, want 1", got)
	}
}

func TestCountSubarraysSumEquals_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []int64
		target int64
		want   int
	}{
		{"empty", nil, 0, 0},
		{"one, equal", []int64{4}, 4, 1},
		{"one, not equal", []int64{4}, 3, 0},
		{"all zeros, target zero", []int64{0, 0, 0}, 0, 6},
		{"alternating signs, target zero", []int64{1, -1, 1, -1}, 0, 4},
		{"duplicates", []int64{2, 2, 2}, 4, 2},
		{"negative target", []int64{-3, 1, -1}, -3, 2},
		{"no match", []int64{1, 2, 3}, 7, 0},
		{"extremes, exact in 128 bits", []int64{math.MaxInt64, math.MaxInt64, math.MinInt64}, math.MaxInt64, 2},
		{"sum beyond int64 is never equal", []int64{math.MaxInt64, math.MaxInt64}, -2, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CountSubarraysSumEquals(tt.values, tt.target); got != tt.want {
				t.Errorf("CountSubarraysSumEquals(%v, %d) = %d, want %d", tt.values, tt.target, got, tt.want)
			}
		})
	}
	if got := CountSubarraysSumEquals([]uint8{255, 1, 255}, 256-256); got != 0 {
		t.Errorf("uint8: count = %d, want 0 (no wrap-around)", got)
	}
}

// Oracle: random int8 data (sums beyond int8 included) against the O(n²) count.
func TestCountSubarraysSumEquals_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		values := randomInt8s(rng, rng.IntN(25))
		before := slices.Clone(values)
		target := int8(rng.IntN(41) - 20)
		if got, want := CountSubarraysSumEquals(values, target), naiveCountSubarrays(values, int(target)); got != want {
			t.Fatalf("trial %d: CountSubarraysSumEquals(%v, %d) = %d, want %d", trial, values, target, got, want)
		}
		if !slices.Equal(values, before) {
			t.Fatalf("trial %d: input modified", trial)
		}
	}
}

func BenchmarkCountSubarraysSumEquals(b *testing.B) {
	for _, size := range benchSizes {
		values := benchInts(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchSink = CountSubarraysSumEquals(values, 100)
			}
		})
	}
}

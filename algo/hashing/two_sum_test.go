package hashing

import (
	"math"
	"math/rand/v2"
	"testing"
)

// Hand-derived: in 3 8 11 −2 7 the first pair summing to 9 (by smallest j) is 11 + (−2), indices 2 and 3.
func TestTwoSum_example(t *testing.T) {
	if i, j := TwoSum([]int{3, 8, 11, -2, 7}, 9); i != 2 || j != 3 {
		t.Errorf("TwoSum = (%d, %d), want (2, 3)", i, j)
	}
}

// Our own wording of the usual exercise cases.
func TestTwoSum_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		target int
		wi, wj int
	}{
		{"empty", nil, 0, -1, -1},
		{"one", []int{5}, 10, -1, -1},
		{"two", []int{4, 6}, 10, 0, 1},
		{"same value twice", []int{5, 1, 5}, 10, 0, 2},
		{"a value is not its own partner", []int{5}, 10, -1, -1},
		{"smallest i for the j", []int{2, 2, 7}, 9, 0, 2},
		{"smallest j wins", []int{1, 9, 4, 6}, 10, 0, 1},
		{"negative and zero", []int{-4, 0, 4}, 0, 0, 2},
		{"no pair", []int{1, 2, 3}, 100, -1, -1},
		{"extremes sum", []int{math.MaxInt, math.MinInt}, -1, 0, 1},
		{"overflow never matches", []int{math.MaxInt, 2, math.MinInt + 1}, math.MinInt + 1, -1, -1},
		{"underflow never matches", []int{math.MinInt, -2}, math.MaxInt - 1, -1, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if i, j := TwoSum(tt.values, tt.target); i != tt.wi || j != tt.wj {
				t.Errorf("TwoSum(%v, %d) = (%d, %d), want (%d, %d)", tt.values, tt.target, i, j, tt.wi, tt.wj)
			}
		})
	}
}

func TestTwoSum_unsigned(t *testing.T) {
	// 250 + 10 wraps to 4 in uint8; the true sum is 260, so no pair sums to 4.
	if i, j := TwoSum([]uint8{250, 10}, 4); i != -1 || j != -1 {
		t.Errorf("TwoSum(uint8 wrap) = (%d, %d), want (-1, -1)", i, j)
	}
	if i, j := TwoSum([]uint8{200, 55}, 255); i != 0 || j != 1 {
		t.Errorf("TwoSum(uint8) = (%d, %d), want (0, 1)", i, j)
	}
}

// Oracle: the O(n²) double loop, scanning j then i ascending; int8 values make wrapping common.
func TestTwoSum_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 500 {
		n := rng.IntN(40)
		values := make([]int8, n)
		for i := range values {
			values[i] = int8(rng.IntN(256) - 128)
		}
		target := int8(rng.IntN(256) - 128)
		wi, wj := -1, -1
	search:
		for j := range values {
			for i := range j {
				if int(values[i])+int(values[j]) == int(target) {
					wi, wj = i, j
					break search
				}
			}
		}
		if i, j := TwoSum(values, target); i != wi || j != wj {
			t.Fatalf("trial %d: TwoSum(%v, %d) = (%d, %d), want (%d, %d)", trial, values, target, i, j, wi, wj)
		}
	}
}

func BenchmarkTwoSum(b *testing.B) {
	for _, size := range benchSizes {
		// Distinct values with no pair summing to -1: the scan runs to the end.
		values := make([]int, size.n)
		for i := range values {
			values[i] = i
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = TwoSum(values, -1)
			}
		})
	}
}

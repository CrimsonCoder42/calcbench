package hashing

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// Hand-derived: in 5 1 7 5 the two 5s are 3 apart.
func TestContainsNearbyDuplicate_example(t *testing.T) {
	values := []int{5, 1, 7, 5}
	if got, err := ContainsNearbyDuplicate(values, 3); err != nil || !got {
		t.Errorf("ContainsNearbyDuplicate(gap 3) = (%v, %v), want (true, nil)", got, err)
	}
	if got, err := ContainsNearbyDuplicate(values, 2); err != nil || got {
		t.Errorf("ContainsNearbyDuplicate(gap 2) = (%v, %v), want (false, nil)", got, err)
	}
}

func TestContainsNearbyDuplicate_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		maxGap int
		want   bool
	}{
		{"empty", nil, 3, false},
		{"one", []int{1}, 3, false},
		{"adjacent, gap 0", []int{1, 1}, 0, false},
		{"adjacent, gap 1", []int{1, 1}, 1, true},
		{"gap larger than n", []int{2, 3, 2}, math.MaxInt, true},
		{"all distinct", []int{1, 2, 3, 4}, 3, false},
		{"far duplicate only", []int{9, 1, 2, 3, 9}, 3, false},
		{"negative and zero", []int{0, -1, 0}, 2, true},
		{"extremes", []int{math.MinInt, math.MaxInt, math.MinInt}, 2, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := ContainsNearbyDuplicate(tt.values, tt.maxGap); err != nil || got != tt.want {
				t.Errorf("ContainsNearbyDuplicate(%v, %d) = (%v, %v), want (%v, nil)", tt.values, tt.maxGap, got, err, tt.want)
			}
		})
	}
}

func TestContainsNearbyDuplicate_nan(t *testing.T) {
	if got, _ := ContainsNearbyDuplicate([]float64{math.NaN(), math.NaN()}, 1); got {
		t.Error("ContainsNearbyDuplicate(NaN NaN) = true, want false: a NaN equals nothing")
	}
}

func TestContainsNearbyDuplicate_errors(t *testing.T) {
	for _, gap := range []int{-1, math.MinInt} {
		if got, err := ContainsNearbyDuplicate([]int{1, 1}, gap); !errors.Is(err, ErrNegativeGap) || got {
			t.Errorf("ContainsNearbyDuplicate(gap %d) = (%v, %v), want (false, ErrNegativeGap)", gap, got, err)
		}
	}
}

// Oracle: the O(n²) pair check.
func TestContainsNearbyDuplicate_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 500 {
		n := rng.IntN(60)
		values := randomInts(rng, n, 1+rng.IntN(2*n+1))
		gap := rng.IntN(n + 2)
		want := false
		for j := range values {
			for i := max(0, j-gap); i < j; i++ {
				want = want || values[i] == values[j]
			}
		}
		if got, err := ContainsNearbyDuplicate(values, gap); err != nil || got != want {
			t.Fatalf("trial %d: ContainsNearbyDuplicate(%v, %d) = (%v, %v), want %v", trial, values, gap, got, err, want)
		}
	}
}

func BenchmarkContainsNearbyDuplicate(b *testing.B) {
	for _, size := range benchSizes {
		// Distinct values: the scan runs to the end with a full window of 100.
		values := make([]int, size.n)
		for i := range values {
			values[i] = i
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = ContainsNearbyDuplicate(values, 100)
			}
		})
	}
}

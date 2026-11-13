package hashing

import (
	"math"
	"math/rand/v2"
	"testing"
)

// Hand-derived: the bytes of "statistics" count s 3, t 3, a 1, i 2, c 1; 'a' at index 2 is first unique.
func TestFirstUnique_statistics(t *testing.T) {
	if got := FirstUnique([]byte("statistics")); got != 2 {
		t.Errorf("FirstUnique(statistics) = %d, want 2", got)
	}
}

func TestFirstUnique_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   int
	}{
		{"empty", nil, -1},
		{"one", []int{4}, 0},
		{"two equal", []int{4, 4}, -1},
		{"two distinct", []int{4, 5}, 0},
		{"unique at end", []int{1, 2, 1, 2, 3}, 4},
		{"all repeat", []int{3, 1, 3, 1}, -1},
		{"negative and zero", []int{-1, 0, -1}, 1},
		{"extremes", []int{math.MaxInt, math.MinInt, math.MaxInt}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FirstUnique(tt.values); got != tt.want {
				t.Errorf("FirstUnique(%v) = %d, want %d", tt.values, got, tt.want)
			}
		})
	}
}

func TestFirstUnique_nan(t *testing.T) {
	if got := FirstUnique([]float64{1, 1, math.NaN(), math.NaN()}); got != 2 {
		t.Errorf("FirstUnique(1 1 NaN NaN) = %d, want 2 (a NaN equals nothing)", got)
	}
}

// Oracle: the O(n²) pairwise count.
func TestFirstUnique_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		n := rng.IntN(60)
		values := randomInts(rng, n, 1+rng.IntN(n+1))
		want := -1
		for i, v := range values {
			c := 0
			for _, w := range values {
				if w == v {
					c++
				}
			}
			if c == 1 {
				want = i
				break
			}
		}
		if got := FirstUnique(values); got != want {
			t.Fatalf("trial %d: FirstUnique(%v) = %d, want %d", trial, values, got, want)
		}
	}
}

func BenchmarkFirstUnique(b *testing.B) {
	for _, size := range benchSizes {
		// Every value repeats except the last: the second scan runs to the end.
		values := make([]int, size.n)
		for i := range values {
			values[i] = i / 2
		}
		values[size.n-1] = -1
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = FirstUnique(values)
			}
		})
	}
}

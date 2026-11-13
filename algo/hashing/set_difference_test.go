package hashing

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Hand-derived: a = 1 2 2 3 4, b = 4 3 3 5 → a ∖ b = 1 2.
func TestSetDifference_example(t *testing.T) {
	if got := SetDifference([]int{1, 2, 2, 3, 4}, []int{4, 3, 3, 5}); !slices.Equal(got, []int{1, 2}) {
		t.Errorf("SetDifference = %v, want [1 2]", got)
	}
}

func TestSetDifference_cases(t *testing.T) {
	tests := []struct {
		name string
		a, b []int
		want []int
	}{
		{"both empty", nil, nil, []int{}},
		{"a empty", nil, []int{1}, []int{}},
		{"b empty", []int{3, 3, 1}, nil, []int{3, 1}},
		{"b covers a", []int{1, 2}, []int{2, 1, 0}, []int{}},
		{"negative and zero", []int{0, -1, -2}, []int{-1}, []int{0, -2}},
		{"extremes", []int{math.MinInt, math.MaxInt}, []int{math.MinInt}, []int{math.MaxInt}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SetDifference(tt.a, tt.b)
			if got == nil || !slices.Equal(got, tt.want) {
				t.Errorf("SetDifference(%v, %v) = %#v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestSetDifference_nan(t *testing.T) {
	nan := math.NaN()
	got := SetDifference([]float64{nan, 1, nan}, []float64{nan})
	if len(got) != 3 || !math.IsNaN(got[0]) || got[1] != 1 || !math.IsNaN(got[2]) {
		t.Errorf("SetDifference with NaN = %v, want [NaN 1 NaN]", got)
	}
}

func TestSetDifference_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		a := randomInts(rng, rng.IntN(80), 50)
		b := randomInts(rng, rng.IntN(80), 50)
		if got, want := SetDifference(a, b), setOracle(a, b, "difference"); !slices.Equal(got, want) {
			t.Fatalf("trial %d: SetDifference = %v, want %v", trial, got, want)
		}
	}
}

func BenchmarkSetDifference(b *testing.B) {
	for _, size := range benchSizes {
		rng := rand.New(rand.NewPCG(1, 2))
		x, y := randomInts(rng, size.n, 2*size.n), randomInts(rng, size.n, 2*size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = SetDifference(x, y)
			}
		})
	}
}

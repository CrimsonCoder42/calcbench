package hashing

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Hand-derived: a = 1 2 2 3 4, b = 4 3 3 5 → a ∪ b = 1 2 3 4 5 (a's order, then b's new values).
func TestSetUnion_example(t *testing.T) {
	if got := SetUnion([]int{1, 2, 2, 3, 4}, []int{4, 3, 3, 5}); !slices.Equal(got, []int{1, 2, 3, 4, 5}) {
		t.Errorf("SetUnion = %v, want [1 2 3 4 5]", got)
	}
}

func TestSetUnion_cases(t *testing.T) {
	tests := []struct {
		name string
		a, b []int
		want []int
	}{
		{"both empty", nil, nil, []int{}},
		{"a empty", nil, []int{2, 2, 1}, []int{2, 1}},
		{"b empty", []int{1, 1}, nil, []int{1}},
		{"disjoint", []int{5, 1}, []int{0, -3}, []int{5, 1, 0, -3}},
		{"equal sets", []int{1, 2}, []int{2, 1}, []int{1, 2}},
		{"extremes", []int{math.MaxInt}, []int{math.MinInt, math.MaxInt}, []int{math.MaxInt, math.MinInt}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SetUnion(tt.a, tt.b)
			if got == nil || !slices.Equal(got, tt.want) {
				t.Errorf("SetUnion(%v, %v) = %#v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestSetUnion_nan(t *testing.T) {
	nan := math.NaN()
	got := SetUnion([]float64{nan, 1}, []float64{nan, 1})
	if len(got) != 3 || !math.IsNaN(got[0]) || got[1] != 1 || !math.IsNaN(got[2]) {
		t.Errorf("SetUnion with NaN = %v, want [NaN 1 NaN]", got)
	}
}

func TestSetUnion_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		a := randomInts(rng, rng.IntN(80), 50)
		b := randomInts(rng, rng.IntN(80), 50)
		if got, want := SetUnion(a, b), setOracle(a, b, "union"); !slices.Equal(got, want) {
			t.Fatalf("trial %d: SetUnion = %v, want %v", trial, got, want)
		}
	}
}

func BenchmarkSetUnion(b *testing.B) {
	for _, size := range benchSizes {
		rng := rand.New(rand.NewPCG(1, 2))
		x, y := randomInts(rng, size.n, 2*size.n), randomInts(rng, size.n, 2*size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = SetUnion(x, y)
			}
		})
	}
}

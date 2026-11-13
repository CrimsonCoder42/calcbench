package hashing

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// setOracle builds a set operation the slow way: scan a (then b for union) in order, keep a value when
// keep says so and it is not yet in the output. Shared by the three set-operation tests.
func setOracle(a, b []int, op string) []int {
	out := []int{}
	add := func(v int) {
		if !slices.Contains(out, v) {
			out = append(out, v)
		}
	}
	for _, v := range a {
		inB := slices.Contains(b, v)
		if op == "union" || (op == "intersection" && inB) || (op == "difference" && !inB) {
			add(v)
		}
	}
	if op == "union" {
		for _, v := range b {
			add(v)
		}
	}
	return out
}

// Hand-derived: a = 1 2 2 3 4, b = 4 3 3 5 → a ∩ b = 3 4 (a's order).
func TestSetIntersection_example(t *testing.T) {
	if got := SetIntersection([]int{1, 2, 2, 3, 4}, []int{4, 3, 3, 5}); !slices.Equal(got, []int{3, 4}) {
		t.Errorf("SetIntersection = %v, want [3 4]", got)
	}
}

func TestSetIntersection_cases(t *testing.T) {
	tests := []struct {
		name string
		a, b []int
		want []int
	}{
		{"both empty", nil, nil, []int{}},
		{"a empty", nil, []int{1}, []int{}},
		{"b empty", []int{1}, nil, []int{}},
		{"disjoint", []int{1, 2}, []int{3}, []int{}},
		{"duplicates in both", []int{2, 2, 1, 2}, []int{2, 2}, []int{2}},
		{"order of a", []int{9, -1, 5, 0}, []int{0, 5, 9}, []int{9, 5, 0}},
		{"extremes", []int{math.MinInt, math.MaxInt}, []int{math.MaxInt}, []int{math.MaxInt}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SetIntersection(tt.a, tt.b)
			if got == nil || !slices.Equal(got, tt.want) {
				t.Errorf("SetIntersection(%v, %v) = %#v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestSetIntersection_nan(t *testing.T) {
	nan := math.NaN()
	if got := SetIntersection([]float64{nan, 1}, []float64{nan, 1}); !slices.Equal(got, []float64{1}) {
		t.Errorf("SetIntersection with NaN = %v, want [1]", got)
	}
}

func TestSetIntersection_doesNotModifyInput(t *testing.T) {
	a, b := []int{3, 1, 3}, []int{1, 3}
	_ = SetIntersection(a, b)
	if !slices.Equal(a, []int{3, 1, 3}) || !slices.Equal(b, []int{1, 3}) {
		t.Errorf("SetIntersection modified its input: %v %v", a, b)
	}
}

func TestSetIntersection_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		a := randomInts(rng, rng.IntN(80), 50)
		b := randomInts(rng, rng.IntN(80), 50)
		if got, want := SetIntersection(a, b), setOracle(a, b, "intersection"); !slices.Equal(got, want) {
			t.Fatalf("trial %d: SetIntersection = %v, want %v", trial, got, want)
		}
	}
}

func BenchmarkSetIntersection(b *testing.B) {
	for _, size := range benchSizes {
		rng := rand.New(rand.NewPCG(1, 2))
		x, y := randomInts(rng, size.n, 2*size.n), randomInts(rng, size.n, 2*size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = SetIntersection(x, y)
			}
		})
	}
}

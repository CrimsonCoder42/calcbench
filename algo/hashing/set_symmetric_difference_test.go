package hashing

import (
	"math"
	"math/rand/v2"
	"slices"
	"strconv"
	"testing"
)

// Our example (Knuth vol. 3 §6.4): {1 2 3 4} xor {3 4 5 6} = 1 2 (from a) then 5 6 (from b).
func TestSetSymmetricDifference_example(t *testing.T) {
	got := SetSymmetricDifference([]int{1, 2, 3, 4}, []int{3, 4, 5, 6})
	if !slices.Equal(got, []int{1, 2, 5, 6}) {
		t.Errorf("SetSymmetricDifference = %v", got)
	}
}

func TestSetSymmetricDifference_cases(t *testing.T) {
	tests := []struct {
		name string
		a, b []int
		want []int
	}{
		{"both empty", nil, nil, []int{}},
		{"a empty", nil, []int{2, 1, 2}, []int{2, 1}},
		{"b empty", []int{3, 3, 1}, nil, []int{3, 1}},
		{"equal", []int{1, 2}, []int{2, 1}, []int{}},
		{"disjoint", []int{1}, []int{2}, []int{1, 2}},
		{"duplicates", []int{5, 1, 5}, []int{1, 7, 7}, []int{5, 7}},
	}
	for _, tc := range tests {
		got := SetSymmetricDifference(tc.a, tc.b)
		if got == nil || !slices.Equal(got, tc.want) {
			t.Errorf("%s: = %v, want %v", tc.name, got, tc.want)
		}
	}
	got := SetSymmetricDifference([]float64{math.NaN(), 0}, []float64{math.NaN(), math.Copysign(0, -1)})
	if len(got) != 2 || !math.IsNaN(got[0]) || !math.IsNaN(got[1]) {
		t.Errorf("NaN and signed zero: %v", got)
	}
}

// Oracle: the result as a set equals (A \ B) ∪ (B \ A) built from SetDifference, on 300 random inputs.
func TestSetSymmetricDifference_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		a, b := make([]int, r.IntN(40)), make([]int, r.IntN(40))
		for i := range a {
			a[i] = r.IntN(25)
		}
		for i := range b {
			b[i] = r.IntN(25)
		}
		got := SetSymmetricDifference(a, b)
		want := append(SetDifference(a, b), SetDifference(b, a)...)
		if !slices.Equal(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func BenchmarkSetSymmetricDifference(b *testing.B) {
	for _, n := range []int{1e3, 1e4, 1e5} {
		r := rand.New(rand.NewPCG(3, 4))
		x, y := make([]int, n), make([]int, n)
		for i := range x {
			x[i], y[i] = r.IntN(2*n), r.IntN(2*n)
		}
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = SetSymmetricDifference(x, y)
			}
		})
	}
}

package hashing

import (
	"math"
	"math/rand/v2"
	"slices"
	"strconv"
	"testing"
)

// Our example (Knuth vol. 3 §6.4, membership by hashing): which of 1 2 3 4 5 are in {2, 4, 9}?
func TestIsIn_example(t *testing.T) {
	got := IsIn([]int{1, 2, 3, 4, 5}, []int{2, 4, 9})
	if !slices.Equal(got, []bool{false, true, false, true, false}) {
		t.Errorf("IsIn = %v", got)
	}
}

func TestIsIn_cases(t *testing.T) {
	if got := IsIn([]int{}, []int{1}); got == nil || len(got) != 0 {
		t.Errorf("empty values: %v", got)
	}
	if got := IsIn([]int{1, 2}, nil); !slices.Equal(got, []bool{false, false}) {
		t.Errorf("empty test: %v", got)
	}
	if got := IsIn([]int{7, 7}, []int{7, 7}); !slices.Equal(got, []bool{true, true}) {
		t.Errorf("duplicates: %v", got)
	}
	nan := math.NaN()
	if got := IsIn([]float64{nan, math.Copysign(0, -1)}, []float64{nan, 0}); !slices.Equal(got, []bool{false, true}) {
		t.Errorf("NaN and signed zero: %v", got)
	}
}

// Oracle: agrees with slices.Contains on 300 random inputs.
func TestIsIn_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		values, test := make([]int, r.IntN(50)), make([]int, r.IntN(50))
		for i := range values {
			values[i] = r.IntN(30)
		}
		for i := range test {
			test[i] = r.IntN(30)
		}
		got := IsIn(values, test)
		for i, v := range values {
			if got[i] != slices.Contains(test, v) {
				t.Fatalf("IsIn[%d] = %v for %d", i, got[i], v)
			}
		}
	}
}

func BenchmarkIsIn(b *testing.B) {
	for _, n := range []int{1e3, 1e4, 1e5} {
		r := rand.New(rand.NewPCG(3, 4))
		values, test := make([]int, n), make([]int, n)
		for i := range values {
			values[i], test[i] = r.IntN(2*n), r.IntN(2*n)
		}
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = IsIn(values, test)
			}
		})
	}
}

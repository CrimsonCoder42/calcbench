package selection

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Worked example on the data of Knuth, TAOCP Vol. 3, 2nd ed., ch. 5 (the 16 keys): the 8th and 9th of
// the sorted keys are 509 and 512. With 703 dropped, 15 keys remain and the 8th is 509. Hand-derived.
func TestMedian_knuthExample(t *testing.T) {
	lower, upper, err := Median(knuth16)
	if err != nil || lower != 509 || upper != 512 {
		t.Errorf("Median(knuth16) = (%d, %d, %v), want (509, 512, nil)", lower, upper, err)
	}
	lower, upper, err = Median(knuth16[:15])
	if err != nil || lower != 509 || upper != 509 {
		t.Errorf("Median(knuth16[:15]) = (%d, %d, %v), want (509, 509, nil)", lower, upper, err)
	}
}

func TestMedian_cases(t *testing.T) {
	tests := []struct {
		name         string
		values       []int
		lower, upper int
	}{
		{"one element", []int{42}, 42, 42},
		{"two elements", []int{9, 3}, 3, 9},
		{"three elements", []int{9, 3, 5}, 5, 5},
		{"duplicates straddle the middle", []int{2, 1, 2, 1}, 1, 2},
		{"all equal", []int{7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7}, 7, 7},
		{"already sorted, even", []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}, 9, 10},
		{"reversed, even", []int{19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}, 9, 10},
		{"negative and zero", []int{-3, 0, -7, 5}, -3, 0},
		{"extremes of the type", []int{math.MaxInt, math.MinInt}, math.MinInt, math.MaxInt},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lower, upper, err := Median(tt.values)
			if err != nil || lower != tt.lower || upper != tt.upper {
				t.Errorf("Median(%v) = (%d, %d, %v), want (%d, %d, nil)", tt.values, lower, upper, err, tt.lower, tt.upper)
			}
		})
	}
}

func TestMedian_floats(t *testing.T) {
	// cmp.Compare order: NaN, -Inf, 1, +Inf, so the middles are -Inf and 1.
	lower, upper, err := Median([]float64{math.Inf(1), 1, math.NaN(), math.Inf(-1)})
	if err != nil || !math.IsInf(lower, -1) || upper != 1 {
		t.Errorf("Median(floats) = (%v, %v, %v), want (-Inf, 1, nil)", lower, upper, err)
	}
}

func TestMedian_errors(t *testing.T) {
	for _, values := range [][]int{nil, {}} {
		if lower, upper, err := Median(values); !errors.Is(err, ErrEmpty) || lower != 0 || upper != 0 {
			t.Errorf("Median(%v) = (%d, %d, %v), want (0, 0, ErrEmpty)", values, lower, upper, err)
		}
	}
}

func TestMedian_doesNotModifyInput(t *testing.T) {
	values := slices.Clone(knuth16)
	if _, _, err := Median(values); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(values, knuth16) {
		t.Errorf("Median modified its input: %v", values)
	}
}

func TestMedian_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		n := 1 + rng.IntN(300)
		values := randomInts(rng, n, 1+rng.IntN(2*n))
		sorted := slices.Sorted(slices.Values(values))
		lower, upper, err := Median(values)
		if err != nil || lower != sorted[(n-1)/2] || upper != sorted[n/2] {
			t.Fatalf("trial %d: Median(n=%d) = (%d, %d, %v), want (%d, %d)",
				trial, n, lower, upper, err, sorted[(n-1)/2], sorted[n/2])
		}
	}
}

func BenchmarkMedian(b *testing.B) {
	for _, size := range benchSizes {
		inputs := benchInputs(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; b.Loop(); i++ {
				lower, _, _ := Median(inputs[i&7])
				benchSink = lower
			}
		})
	}
}

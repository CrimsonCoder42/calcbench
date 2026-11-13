package selection

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Worked example on the 16 keys of Knuth, TAOCP Vol. 3, 2nd ed., ch. 5. The k-th largest is read off the
// sorted list from the right by hand.
func TestKthLargest_knuthExample(t *testing.T) {
	for k, want := range map[int]int{1: 908, 3: 765, 8: 512, 16: 61} {
		got, err := KthLargest(knuth16, k)
		if err != nil || got != want {
			t.Errorf("KthLargest(knuth16, %d) = (%d, %v), want (%d, nil)", k, got, err, want)
		}
	}
}

// Our own data, answers by hand: the 2nd largest of six scores (descending 23 17 14 11 8 5), and the
// 4th largest of a list with repeats, which count separately (descending 12 12 9 9 7 4 4 2).
func TestKthLargest_exercises(t *testing.T) {
	if got, _ := KthLargest([]int{14, 8, 23, 5, 17, 11}, 2); got != 17 {
		t.Errorf("2nd largest = %d, want 17", got)
	}
	if got, _ := KthLargest([]int{9, 4, 9, 12, 4, 7, 12, 2}, 4); got != 9 {
		t.Errorf("4th largest = %d, want 9", got)
	}
}

func TestKthLargest_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		k      int
		want   int
	}{
		{"one element", []int{42}, 1, 42},
		{"two elements, first", []int{9, 3}, 1, 9},
		{"two elements, second", []int{9, 3}, 2, 3},
		{"duplicates count separately", []int{4, 1, 4, 1, 4}, 4, 1},
		{"all equal", []int{7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7}, 11, 7},
		{"already sorted", []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}, 3, 17},
		{"reversed", []int{19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0}, 3, 17},
		{"negative and zero", []int{-3, 0, -7, 5, -1}, 2, 0},
		{"maximum of the type", []int{math.MaxInt, math.MinInt, 0}, 1, math.MaxInt},
		{"minimum of the type", []int{math.MaxInt, math.MinInt, 0}, 3, math.MinInt},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := KthLargest(tt.values, tt.k)
			if err != nil || got != tt.want {
				t.Errorf("KthLargest(%v, %d) = (%d, %v), want (%d, nil)", tt.values, tt.k, got, err, tt.want)
			}
		})
	}
}

func TestKthLargest_floats(t *testing.T) {
	values := []float64{1, math.NaN(), math.Inf(1), math.Copysign(0, -1), math.Inf(-1), 0.5}
	// cmp.Compare order, largest first: +Inf, 1, 0.5, -0, -Inf, NaN.
	if got, err := KthLargest(values, 1); err != nil || got != math.Inf(1) {
		t.Errorf("KthLargest(floats, 1) = (%v, %v), want +Inf", got, err)
	}
	if got, err := KthLargest(values, 6); err != nil || !math.IsNaN(got) {
		t.Errorf("KthLargest(floats, 6) = (%v, %v), want NaN", got, err)
	}
}

func TestKthLargest_errors(t *testing.T) {
	tests := []struct {
		name    string
		values  []int
		k       int
		wantErr error
	}{
		{"nil input", nil, 1, ErrEmpty},
		{"k zero", []int{1, 2, 3}, 0, ErrRank},
		{"k negative", []int{1, 2, 3}, -1, ErrRank},
		{"k above n", []int{1, 2, 3}, 4, ErrRank},
		{"k at the minimum of int", []int{1, 2, 3}, math.MinInt, ErrRank},
		{"k at the maximum of int", []int{1, 2, 3}, math.MaxInt, ErrRank},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := KthLargest(tt.values, tt.k)
			if !errors.Is(err, tt.wantErr) || got != 0 {
				t.Errorf("KthLargest(%v, %d) = (%d, %v), want (0, %v)", tt.values, tt.k, got, err, tt.wantErr)
			}
		})
	}
}

func TestKthLargest_doesNotModifyInput(t *testing.T) {
	values := slices.Clone(knuth16)
	if _, err := KthLargest(values, 4); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(values, knuth16) {
		t.Errorf("KthLargest modified its input: %v", values)
	}
}

func TestKthLargest_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		n := 1 + rng.IntN(300)
		values := randomInts(rng, n, 1+rng.IntN(2*n))
		sorted := slices.Sorted(slices.Values(values))
		k := 1 + rng.IntN(n)
		got, err := KthLargest(values, k)
		if err != nil || got != sorted[n-k] {
			t.Fatalf("trial %d: KthLargest(n=%d, k=%d) = (%d, %v), want %d", trial, n, k, got, err, sorted[n-k])
		}
	}
}

func BenchmarkKthLargest(b *testing.B) {
	for _, size := range benchSizes {
		inputs := benchInputs(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; b.Loop(); i++ {
				v, _ := KthLargest(inputs[i&7], size.n/10)
				benchSink = v
			}
		})
	}
}

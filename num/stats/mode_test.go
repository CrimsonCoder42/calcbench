package stats

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Worked example: e-Handbook §1.3.5.1, the mode "occurs with the greatest frequency. It is not
// necessarily unique." 2, 3, 3, 5, 5, 5 → [5] ×3; 1, 1, 2, 2, 3 → [1 2] ×2 (hand-checked).
func TestMode_worked(t *testing.T) {
	if modes, count, err := Mode([]float64{2, 3, 3, 5, 5, 5}); err != nil || count != 3 || !slices.Equal(modes, []float64{5}) {
		t.Errorf("Mode(2 3 3 5 5 5) = %v, %d, %v; want [5], 3", modes, count, err)
	}
	if modes, count, err := Mode([]int{1, 1, 2, 2, 3}); err != nil || count != 2 || !slices.Equal(modes, []int{1, 2}) {
		t.Errorf("Mode(1 1 2 2 3) = %v, %d, %v; want [1 2], 2", modes, count, err)
	}
}

func TestMode_cases(t *testing.T) {
	for _, tt := range []struct {
		name   string
		values []int
		modes  []int
		count  int
	}{
		{"one element", []int{7}, []int{7}, 1},
		{"all distinct", []int{3, 1, 2}, []int{1, 2, 3}, 1},
		{"all equal", []int{4, 4, 4}, []int{4}, 3},
		{"unsorted ties", []int{9, -1, 9, -1, 0}, []int{-1, 9}, 2},
		{"extremes of the type", []int{math.MaxInt, math.MinInt, math.MaxInt}, []int{math.MaxInt}, 2},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if modes, count, err := Mode(tt.values); err != nil || count != tt.count || !slices.Equal(modes, tt.modes) {
				t.Errorf("Mode(%v) = %v, %d, %v; want %v, %d", tt.values, modes, count, err, tt.modes, tt.count)
			}
		})
	}
}

func TestMode_floatsAndStrings(t *testing.T) {
	// cmp.Compare: NaN equals NaN and sorts first; −0 equals +0.
	modes, count, err := Mode([]float64{math.NaN(), 1, math.NaN(), 1})
	if err != nil || count != 2 || len(modes) != 2 || !math.IsNaN(modes[0]) || modes[1] != 1 {
		t.Errorf("Mode(NaN 1 NaN 1) = %v, %d, %v; want [NaN 1], 2", modes, count, err)
	}
	if modes, count, err := Mode([]float64{math.Copysign(0, -1), 0, 2}); err != nil || count != 2 || len(modes) != 1 || modes[0] != 0 {
		t.Errorf("Mode(-0 0 2) = %v, %d, %v; want [0], 2", modes, count, err)
	}
	if modes, count, err := Mode([]string{"b", "a", "b"}); err != nil || count != 2 || !slices.Equal(modes, []string{"b"}) {
		t.Errorf("Mode(b a b) = %v, %d, %v; want [b], 2", modes, count, err)
	}
}

func TestMode_errors(t *testing.T) {
	if modes, count, err := Mode([]int{}); !errors.Is(err, ErrEmpty) || modes != nil || count != 0 {
		t.Errorf("Mode([]) = %v, %d, %v; want nil, 0, ErrEmpty", modes, count, err)
	}
}

func TestMode_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		n := 1 + rng.IntN(300)
		values := make([]int, n)
		for i := range values {
			values[i] = rng.IntN(1 + n/3)
		}
		freq := map[int]int{}
		best := 0
		for _, v := range values {
			freq[v]++
			best = max(best, freq[v])
		}
		var want []int
		for v, c := range freq {
			if c == best {
				want = append(want, v)
			}
		}
		slices.Sort(want)
		orig := slices.Clone(values)
		modes, count, err := Mode(values)
		if err != nil || count != best || !slices.Equal(modes, want) || !slices.Equal(values, orig) {
			t.Fatalf("trial %d: Mode = %v, %d, %v; want %v, %d (input unchanged)", trial, modes, count, err, want, best)
		}
	}
}

func BenchmarkMode(b *testing.B) {
	for _, size := range benchSizes {
		inputs := benchInputs(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; b.Loop(); i++ {
				modes, _, _ := Mode(inputs[i&7])
				benchSink = modes[0]
			}
		})
	}
}

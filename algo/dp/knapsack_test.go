package dp

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// knapsack01Slow tries every subset (n ≤ 14).
func knapsack01Slow(weights []int, values []int64, capacity int) int64 {
	var best int64
	for mask := 0; mask < 1<<len(weights); mask++ {
		w, v := 0, int64(0)
		for i := range weights {
			if mask>>i&1 == 1 {
				w += weights[i]
				v += values[i]
			}
		}
		if w <= capacity {
			best = max(best, v)
		}
	}
	return best
}

// unboundedSlow tries every vector of copies item by item (small inputs only): a different shape from
// the capacity-indexed table under test.
func unboundedSlow(weights []int, values []int64, capacity int) int64 {
	if len(weights) == 0 {
		return 0
	}
	var best int64
	for copies := 0; copies*weights[0] <= capacity; copies++ {
		rest := unboundedSlow(weights[1:], values[1:], capacity-copies*weights[0])
		best = max(best, int64(copies)*values[0]+rest)
	}
	return best
}

// checkChosen verifies that items is ascending, distinct, within capacity and worth best.
func checkChosen(t *testing.T, weights []int, values []int64, capacity int, best int64, items []int) {
	t.Helper()
	w, v := 0, int64(0)
	for k, i := range items {
		if i < 0 || i >= len(weights) || (k > 0 && i <= items[k-1]) {
			t.Fatalf("items %v not ascending distinct indices", items)
		}
		w += weights[i]
		v += values[i]
	}
	if w > capacity || v != best {
		t.Fatalf("items %v weigh %d (cap %d) and are worth %d, want %d", items, w, capacity, v, best)
	}
}

// CLRS 3rd ed., §16.2, Figure 16.2: items (w, v) = (10, 60), (20, 100), (30, 120), W = 50; the 0-1
// optimum takes items 2 and 3 (0-based 1 and 2) for 220.
func TestKnapsack01_clrsFigure16_2(t *testing.T) {
	best, items, err := Knapsack01([]int{10, 20, 30}, []int64{60, 100, 120}, 50)
	if err != nil || best != 220 || !slices.Equal(items, []int{1, 2}) {
		t.Errorf("Knapsack01(CLRS 16.2) = (%d, %v, %v), want (220, [1 2], nil)", best, items, err)
	}
}

func TestKnapsack01_cases(t *testing.T) {
	tests := []struct {
		name     string
		weights  []int
		values   []int64
		capacity int
		best     int64
		items    []int
	}{
		{"no items", nil, nil, 10, 0, []int{}},
		{"zero capacity", []int{1}, []int64{5}, 0, 0, []int{}},
		{"one fits", []int{3}, []int64{7}, 3, 7, []int{0}},
		{"one too heavy", []int{4}, []int64{7}, 3, 0, []int{}},
		{"two, take both", []int{1, 2}, []int64{3, 4}, 3, 7, []int{0, 1}},
		{"zero-weight item is free", []int{0, 5}, []int64{9, 1}, 0, 9, []int{0}},
		{"zero-value item is not taken", []int{1}, []int64{0}, 5, 0, []int{}},
		{"duplicates", []int{2, 2, 2}, []int64{5, 5, 5}, 4, 10, []int{0, 1}},
		{"capacity far above total", []int{1, 1}, []int64{1, 2}, math.MaxInt, 3, []int{0, 1}},
		{"extreme value", []int{1}, []int64{math.MaxInt64}, 1, math.MaxInt64, []int{0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			best, items, err := Knapsack01(tt.weights, tt.values, tt.capacity)
			if err != nil || best != tt.best || items == nil || !slices.Equal(items, tt.items) {
				t.Errorf("Knapsack01 = (%d, %#v, %v), want (%d, %v, nil)", best, items, err, tt.best, tt.items)
			}
		})
	}
}

func TestKnapsack01_errors(t *testing.T) {
	tests := []struct {
		name     string
		weights  []int
		values   []int64
		capacity int
		want     error
	}{
		{"length mismatch", []int{1, 2}, []int64{1}, 5, ErrLengthMismatch},
		{"negative capacity", []int{1}, []int64{1}, -1, ErrNegative},
		{"negative weight", []int{-1}, []int64{1}, 5, ErrNegative},
		{"negative value", []int{1}, []int64{-1}, 5, ErrNegative},
		{"overflow", []int{1, 1}, []int64{math.MaxInt64, 1}, 2, ErrOverflow},
		{"table too large", []int{1 << 40}, []int64{1}, 1 << 40, ErrTooLarge},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, _, err := Knapsack01(tt.weights, tt.values, tt.capacity); !errors.Is(err, tt.want) {
				t.Errorf("Knapsack01 error = %v, want %v", err, tt.want)
			}
			if tt.want == ErrOverflow || tt.want == ErrTooLarge {
				return
			}
			if _, _, err := UnboundedKnapsack(tt.weights, tt.values, tt.capacity); !errors.Is(err, tt.want) {
				t.Errorf("UnboundedKnapsack error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestKnapsack01_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		n := rng.IntN(13)
		weights, values := make([]int, n), make([]int64, n)
		for i := range n {
			weights[i] = rng.IntN(20)
			values[i] = rng.Int64N(50)
		}
		capacity := rng.IntN(60)
		best, items, err := Knapsack01(weights, values, capacity)
		if err != nil || best != knapsack01Slow(weights, values, capacity) {
			t.Fatalf("trial %d: Knapsack01(%v, %v, %d) = (%d, %v), want %d", trial, weights, values, capacity, best, err, knapsack01Slow(weights, values, capacity))
		}
		checkChosen(t, weights, values, capacity, best, items)
	}
}

// CLRS 3rd ed., §15.1, Figure 15.1: with one item per rod length 1…10 priced 1 5 8 9 10 17 17 20 24 30,
// the unbounded knapsack of capacity 4 is the rod-cutting revenue r4 = 10 (two pieces of length 2).
func TestUnboundedKnapsack_clrsRodPrices(t *testing.T) {
	weights := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	prices := []int64{1, 5, 8, 9, 10, 17, 17, 20, 24, 30}
	best, counts, err := UnboundedKnapsack(weights, prices, 4)
	if err != nil || best != 10 || !slices.Equal(counts, []int{0, 2, 0, 0, 0, 0, 0, 0, 0, 0}) {
		t.Errorf("UnboundedKnapsack(rod, 4) = (%d, %v, %v), want (10, two of length 2, nil)", best, counts, err)
	}
}

func TestUnboundedKnapsack_cases(t *testing.T) {
	tests := []struct {
		name     string
		weights  []int
		values   []int64
		capacity int
		best     int64
	}{
		{"no items", nil, nil, 5, 0},
		{"zero capacity", []int{1}, []int64{5}, 0, 0},
		{"repeat one item", []int{2}, []int64{3}, 7, 9},
		{"mix", []int{2, 3}, []int64{3, 5}, 7, 11},
		{"zero weight zero value ignored", []int{0, 2}, []int64{0, 1}, 4, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			best, counts, err := UnboundedKnapsack(tt.weights, tt.values, tt.capacity)
			if err != nil || best != tt.best || len(counts) != len(tt.weights) {
				t.Fatalf("UnboundedKnapsack = (%d, %v, %v), want %d", best, counts, err, tt.best)
			}
			checkCounts(t, tt.weights, tt.values, tt.capacity, best, counts)
		})
	}
}

func checkCounts(t *testing.T, weights []int, values []int64, capacity int, best int64, counts []int) {
	t.Helper()
	w, v := 0, int64(0)
	for i, c := range counts {
		if c < 0 {
			t.Fatalf("negative count %v", counts)
		}
		w += c * weights[i]
		v += int64(c) * values[i]
	}
	if w > capacity || v != best {
		t.Fatalf("counts %v weigh %d (cap %d) and are worth %d, want %d", counts, w, capacity, v, best)
	}
}

func TestUnboundedKnapsack_errors(t *testing.T) {
	if _, _, err := UnboundedKnapsack([]int{0}, []int64{1}, 5); !errors.Is(err, ErrUnbounded) {
		t.Errorf("zero weight, positive value: err = %v, want ErrUnbounded", err)
	}
	if _, _, err := UnboundedKnapsack([]int{1}, []int64{math.MaxInt64 / 2}, 3); !errors.Is(err, ErrOverflow) {
		t.Errorf("overflow: err = %v, want ErrOverflow", err)
	}
	for _, capacity := range []int{1 << 40, math.MaxInt} {
		if _, _, err := UnboundedKnapsack([]int{1}, []int64{1}, capacity); !errors.Is(err, ErrTooLarge) {
			t.Errorf("capacity %d: err = %v, want ErrTooLarge", capacity, err)
		}
	}
}

func TestUnboundedKnapsack_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		n := rng.IntN(5)
		weights, values := make([]int, n), make([]int64, n)
		for i := range n {
			weights[i] = 1 + rng.IntN(8)
			values[i] = rng.Int64N(30)
		}
		capacity := rng.IntN(25)
		best, counts, err := UnboundedKnapsack(weights, values, capacity)
		if want := unboundedSlow(weights, values, capacity); err != nil || best != want {
			t.Fatalf("trial %d: UnboundedKnapsack(%v, %v, %d) = (%d, %v), want %d", trial, weights, values, capacity, best, err, want)
		}
		checkCounts(t, weights, values, capacity, best, counts)
	}
}

// 100 items; the capacity grows, so time grows as n·W.
func BenchmarkKnapsack01(b *testing.B) {
	for _, size := range linearSizes {
		weights, values := benchItems(100, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _, _ = Knapsack01(weights, values, size.n)
			}
		})
	}
}

func BenchmarkUnboundedKnapsack(b *testing.B) {
	for _, size := range linearSizes {
		weights, values := benchItems(100, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _, _ = UnboundedKnapsack(weights, values, size.n)
			}
		})
	}
}

func benchItems(n, capacity int) ([]int, []int64) {
	rng := rand.New(rand.NewPCG(1, 2))
	weights, values := make([]int, n), make([]int64, n)
	for i := range n {
		weights[i] = 1 + rng.IntN(capacity/10+1)
		values[i] = rng.Int64N(1000)
	}
	return weights, values
}

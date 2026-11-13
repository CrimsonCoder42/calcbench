package twopointer

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Lomuto's one-pass partition (Cormen et al., 3rd ed., §7.1) with the test "not zero": the kept
// elements reach the front in their original order. Hand-made data.
func TestMoveZeroesInPlace_worked(t *testing.T) {
	values := []int{4, 0, 0, -2, 7, 0}
	n := MoveZeroesInPlace(values)
	if want := []int{4, -2, 7, 0, 0, 0}; n != 3 || !slices.Equal(values, want) {
		t.Errorf("MoveZeroesInPlace = %d, %v, want 3, %v", n, values, want)
	}
}

func TestMoveZeroesInPlace_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   []int
		n      int
	}{
		{"empty", nil, nil, 0},
		{"one zero", []int{0}, []int{0}, 0},
		{"one non-zero", []int{5}, []int{5}, 1},
		{"two, zero first", []int{0, 5}, []int{5, 0}, 1},
		{"all zeros", []int{0, 0, 0}, []int{0, 0, 0}, 0},
		{"no zeros", []int{3, 1, 2}, []int{3, 1, 2}, 3},
		{"duplicates keep order", []int{2, 0, 2, 0, 1}, []int{2, 2, 1, 0, 0}, 3},
		{"zeros already last", []int{1, 2, 0}, []int{1, 2, 0}, 2},
		{"extremes", []int{0, math.MinInt, 0, math.MaxInt}, []int{math.MinInt, math.MaxInt, 0, 0}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values := slices.Clone(tt.values)
			n := MoveZeroesInPlace(values)
			if n != tt.n || !slices.Equal(values, tt.want) {
				t.Errorf("MoveZeroesInPlace(%v) = %d, %v, want %d, %v", tt.values, n, values, tt.n, tt.want)
			}
		})
	}
	words := []string{"", "a", "", "b"}
	if n := MoveZeroesInPlace(words); n != 2 || !slices.Equal(words, []string{"a", "b", "", ""}) {
		t.Errorf("MoveZeroesInPlace(strings) = %d, %q", n, words)
	}
}

// -0 == 0 counts as zero and is swapped, not rewritten, so its sign bit survives; NaN is not zero.
func TestMoveZeroesInPlace_floats(t *testing.T) {
	negZero := math.Copysign(0, -1)
	values := []float64{negZero, math.NaN(), 1}
	n := MoveZeroesInPlace(values)
	if n != 2 || !math.IsNaN(values[0]) || values[1] != 1 || values[2] != 0 || !math.Signbit(values[2]) {
		t.Errorf("MoveZeroesInPlace(-0 NaN 1) = %d, %v, want 2, [NaN 1 -0]", n, values)
	}
}

func TestMoveZeroesInPlace_noAllocation(t *testing.T) {
	src := make([]int, 1000)
	for i := range src {
		src[i] = i % 3
	}
	values := make([]int, len(src))
	allocs := testing.AllocsPerRun(20, func() {
		copy(values, src)
		benchSink = MoveZeroesInPlace(values)
	})
	if allocs != 0 {
		t.Errorf("MoveZeroesInPlace allocates %v times per call, want 0", allocs)
	}
}

// Oracle: the non-zeros filtered in order, then as many zeros as there were.
func TestMoveZeroesInPlace_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		values := make([]int, rng.IntN(25))
		for i := range values {
			values[i] = rng.IntN(4) - 1
		}
		var want []int
		for _, v := range values {
			if v != 0 {
				want = append(want, v)
			}
		}
		nonZero := len(want)
		want = append(want, make([]int, len(values)-nonZero)...)
		original := slices.Clone(values)
		if n := MoveZeroesInPlace(values); n != nonZero || !slices.Equal(values, want) {
			t.Fatalf("trial %d: MoveZeroesInPlace(%v) = %d, %v, want %d, %v", trial, original, n, values, nonZero, want)
		}
	}
}

func BenchmarkMoveZeroesInPlace(b *testing.B) {
	for _, size := range benchSizes {
		rng := rand.New(rand.NewPCG(1, 2))
		src := make([]int64, size.n)
		for i := range src {
			src[i] = int64(rng.IntN(3))
		}
		values := make([]int64, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				copy(values, src)
				benchSink = MoveZeroesInPlace(values)
			}
		})
	}
}

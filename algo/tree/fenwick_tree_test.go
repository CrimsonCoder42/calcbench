package tree

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Laaksonen (2018) §9.2 pp. 86–87: the array 1 3 4 8 6 1 4 2 has the tree 1 4 4 16 6 7 4 29, and the
// sum of its first 7 values is 16 + 7 + 4 = 27.
func TestNewFenwickTree_laaksonenExample(t *testing.T) {
	f, err := NewFenwickTree([]int{1, 3, 4, 8, 6, 1, 4, 2})
	if err != nil {
		t.Fatal(err)
	}
	if want := []int{1, 4, 4, 16, 6, 7, 4, 29}; !slices.Equal(f.t, want) {
		t.Errorf("tree = %v, want %v", f.t, want)
	}
	if got, err := f.PrefixSum(7); err != nil || got != 27 {
		t.Errorf("PrefixSum(7) = (%d, %v), want 27", got, err)
	}
}

// Laaksonen §9.2 p. 88: changing the value at (1-based) position 3 changes tree positions 3, 4 and 8.
func TestFenwickTreeAdd_laaksonenExample(t *testing.T) {
	f, _ := NewFenwickTree([]int{1, 3, 4, 8, 6, 1, 4, 2})
	if err := f.Add(2, 5); err != nil {
		t.Fatal(err)
	}
	if want := []int{1, 4, 9, 21, 6, 7, 4, 34}; !slices.Equal(f.t, want) {
		t.Errorf("tree after Add(2, 5) = %v, want %v", f.t, want)
	}
}

func TestFenwickTreeRangeSum_cases(t *testing.T) {
	values := []int{1, 3, 4, 8, 6, 1, 4, 2}
	f, _ := NewFenwickTree(values)
	for lo := 0; lo <= len(values); lo++ {
		for hi := lo; hi <= len(values); hi++ {
			want := 0
			for _, v := range values[lo:hi] {
				want += v
			}
			if got, err := f.RangeSum(lo, hi); err != nil || got != want {
				t.Errorf("RangeSum(%d, %d) = (%d, %v), want %d", lo, hi, got, err, want)
			}
		}
	}
}

func TestFenwickTree_emptyAndOne(t *testing.T) {
	f, err := NewFenwickTree[int](nil)
	if err != nil || f.Len() != 0 {
		t.Fatalf("empty: (%v, %v)", f, err)
	}
	if s, err := f.PrefixSum(0); err != nil || s != 0 {
		t.Errorf("empty PrefixSum(0) = (%d, %v)", s, err)
	}
	g, _ := NewFenwickTree([]int{-7})
	if s, _ := g.RangeSum(0, 1); s != -7 {
		t.Errorf("one: RangeSum = %d", s)
	}
	var zero FenwickTree[int]
	var nilTree *FenwickTree[int]
	if zero.Len() != 0 || nilTree.Len() != 0 {
		t.Error("zero or nil Len != 0")
	}
	if _, err := nilTree.PrefixSum(0); err != nil {
		t.Errorf("nil PrefixSum(0) = %v, want nil", err)
	}
	if err := nilTree.Add(0, 1); !errors.Is(err, ErrRange) {
		t.Errorf("nil Add: %v, want ErrRange", err)
	}
}

func TestFenwickTree_errors(t *testing.T) {
	f, _ := NewFenwickTree([]int{1, 2, 3})
	for _, i := range []int{-1, 3, math.MaxInt} {
		if err := f.Add(i, 1); !errors.Is(err, ErrRange) {
			t.Errorf("Add(%d): %v, want ErrRange", i, err)
		}
	}
	for _, r := range [][2]int{{-1, 2}, {2, 1}, {0, 4}} {
		if _, err := f.RangeSum(r[0], r[1]); !errors.Is(err, ErrRange) {
			t.Errorf("RangeSum(%d, %d): %v, want ErrRange", r[0], r[1], err)
		}
	}
	if _, err := f.PrefixSum(4); !errors.Is(err, ErrRange) {
		t.Errorf("PrefixSum(4): %v, want ErrRange", err)
	}
	if _, err := NewFenwickTree([]int8{100, 100}); !errors.Is(err, ErrOverflow) {
		t.Errorf("int8 node 200: %v, want ErrOverflow", err)
	}
	if _, err := NewFenwickTree([]float64{1, math.NaN()}); !errors.Is(err, ErrNotFinite) {
		t.Errorf("NaN value: %v, want ErrNotFinite", err)
	}
	g, _ := NewFenwickTree([]float64{1, 2})
	if err := g.Add(0, math.Inf(1)); !errors.Is(err, ErrNotFinite) {
		t.Errorf("Inf delta: %v, want ErrNotFinite", err)
	}
}

// An overflowing Add leaves every node as it was.
func TestFenwickTreeAdd_atomicOnOverflow(t *testing.T) {
	f, _ := NewFenwickTree([]int8{1, 120, 1, 1})
	before := slices.Clone(f.t)
	if err := f.Add(0, 10); !errors.Is(err, ErrOverflow) { // node [0, 2) would be 131
		t.Fatalf("Add: %v, want ErrOverflow", err)
	}
	if !slices.Equal(f.t, before) {
		t.Errorf("tree changed on error: %v, was %v", f.t, before)
	}
}

func TestFenwickTree_extremes(t *testing.T) {
	f, err := NewFenwickTree([]int64{math.MaxInt64, math.MinInt64, 5})
	if err != nil {
		t.Fatal(err)
	}
	if s, err := f.RangeSum(0, 3); err != nil || s != 4 {
		t.Errorf("RangeSum(0, 3) = (%d, %v), want 4", s, err)
	}
	if s, err := f.RangeSum(0, 1); err != nil || s != math.MaxInt64 {
		t.Errorf("RangeSum(0, 1) = (%d, %v)", s, err)
	}
}

func TestFenwickTree_floats(t *testing.T) {
	// Halves and quarters are exact in binary floating point, so the sums compare exactly.
	f, _ := NewFenwickTree([]float64{0.5, -0.25, 2, 0})
	if s, err := f.RangeSum(1, 3); err != nil || s != 1.75 {
		t.Errorf("RangeSum(1, 3) = (%v, %v), want 1.75", s, err)
	}
}

func TestFenwickTree_doesNotModifyInput(t *testing.T) {
	values := []int{1, 2, 3, 4}
	f, _ := NewFenwickTree(values)
	_ = f.Add(0, 10)
	if !slices.Equal(values, []int{1, 2, 3, 4}) {
		t.Errorf("input modified: %v", values)
	}
}

// Oracle: random adds and range sums against a plain array.
func TestFenwickTree_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		n := rng.IntN(70)
		model := make([]int, n)
		for i := range model {
			model[i] = rng.IntN(201) - 100
		}
		f, err := NewFenwickTree(model)
		if err != nil {
			t.Fatal(err)
		}
		for op := range 100 {
			if n > 0 && rng.IntN(2) == 0 {
				i, d := rng.IntN(n), rng.IntN(201)-100
				if err := f.Add(i, d); err != nil {
					t.Fatal(err)
				}
				model[i] += d
				continue
			}
			lo := rng.IntN(n + 1)
			hi := lo + rng.IntN(n-lo+1)
			want := 0
			for _, v := range model[lo:hi] {
				want += v
			}
			if got, err := f.RangeSum(lo, hi); err != nil || got != want {
				t.Fatalf("trial %d op %d: RangeSum(%d, %d) = (%d, %v), want %d", trial, op, lo, hi, got, err, want)
			}
		}
	}
}

func TestFenwickTree_noAllocation(t *testing.T) {
	f, _ := NewFenwickTree(make([]int, 1000))
	if a := testing.AllocsPerRun(100, func() {
		_ = f.Add(17, 3)
		_, _ = f.PrefixSum(900)
		_, _ = f.RangeSum(100, 700)
		_ = f.Len()
	}); a != 0 {
		t.Errorf("allocations per run = %v, want 0", a)
	}
}

func fenwickOf(n int) *FenwickTree[int] {
	values := make([]int, n)
	for i := range values {
		values[i] = i % 97
	}
	f, _ := NewFenwickTree(values)
	return f
}

func BenchmarkNewFenwickTree(b *testing.B) {
	for _, size := range benchSizes {
		values := make([]int, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = NewFenwickTree(values)
			}
		})
	}
}

func BenchmarkFenwickTreeAdd(b *testing.B) {
	for _, size := range benchSizes {
		f := fenwickOf(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			i := 0
			for b.Loop() {
				_ = f.Add(i%size.n, 1)
				i += 7919
			}
		})
	}
}

func BenchmarkFenwickTreePrefixSum(b *testing.B) {
	for _, size := range benchSizes {
		f := fenwickOf(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = f.PrefixSum(size.n - 1)
			}
		})
	}
}

func BenchmarkFenwickTreeRangeSum(b *testing.B) {
	for _, size := range benchSizes {
		f := fenwickOf(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = f.RangeSum(1, size.n-1)
			}
		})
	}
}

func BenchmarkFenwickTreeLen(b *testing.B) {
	f := fenwickOf(10)
	b.ReportAllocs()
	for b.Loop() {
		_ = f.Len()
	}
}

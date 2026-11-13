package tree

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Laaksonen (2018) §28.1 pp. 258–259: the 16-value array sums to 72 (root 72/0); increasing the values
// from position 5 to 13 (0-based, inclusive: [5, 14) here) by 2 makes the root 90 and both halves 45; the
// query that follows pushes the pending update down without changing any sum.
func TestNewRangeAddSumTree_laaksonenExample(t *testing.T) {
	s, err := NewRangeAddSumTree(laaksonen16)
	if err != nil {
		t.Fatal(err)
	}
	if got, _ := s.Query(0, 16); got != 72 {
		t.Fatalf("Query(0, 16) = %d, want 72", got)
	}
	if err := s.RangeUpdate(5, 14, 2); err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct{ lo, hi, want int }{{0, 16, 90}, {0, 8, 45}, {8, 16, 45}, {5, 14, 43 + 18}, {12, 14, 8 + 4}} {
		if got, err := s.Query(tt.lo, tt.hi); err != nil || got != tt.want {
			t.Errorf("Query(%d, %d) = (%d, %v), want %d", tt.lo, tt.hi, got, err, tt.want)
		}
	}
}

func TestNewRangeAddMinTree_cases(t *testing.T) {
	s, _ := NewRangeAddMinTree(laaksonen16)
	_ = s.RangeUpdate(9, 10, 10) // the 1 at index 9 becomes 11
	_ = s.RangeUpdate(0, 8, -3)  // 2 5 3 0 −1 4 −1 3
	for _, tt := range []struct{ lo, hi, want int }{{0, 16, -1}, {8, 16, 2}, {0, 4, 0}, {9, 10, 11}} {
		if got, err := s.Query(tt.lo, tt.hi); err != nil || got != tt.want {
			t.Errorf("min Query(%d, %d) = (%d, %v), want %d", tt.lo, tt.hi, got, err, tt.want)
		}
	}
	if _, err := s.Query(2, 2); !errors.Is(err, ErrEmptyRange) {
		t.Errorf("empty: %v, want ErrEmptyRange", err)
	}
}

func TestNewRangeAddMaxTree_cases(t *testing.T) {
	s, _ := NewRangeAddMaxTree(laaksonen16)
	_ = s.RangeUpdate(0, 16, 1)
	_ = s.RangeUpdate(1, 2, -5) // the 8 (now 9) becomes 4
	for _, tt := range []struct{ lo, hi, want int }{{0, 16, 8}, {0, 3, 7}, {15, 16, 3}} {
		if got, err := s.Query(tt.lo, tt.hi); err != nil || got != tt.want {
			t.Errorf("max Query(%d, %d) = (%d, %v), want %d", tt.lo, tt.hi, got, err, tt.want)
		}
	}
	if _, err := s.Query(0, 0); !errors.Is(err, ErrEmptyRange) {
		t.Errorf("empty: %v, want ErrEmptyRange", err)
	}
}

// affine is the update x → mul·x + add, applied to (sum, count) aggregates: a non-trivial Compose.
type affine struct{ mul, add int }

// A generic lazy tree with range assign-or-scale updates (affine maps) on sums, checked by a model.
func TestNewLazySegmentTree_affine(t *testing.T) {
	ops := LazyOps[int, affine]{
		Combine:  func(l, r int) int { return l + r },
		Identity: 0,
		Apply:    func(u affine, sum, length int) int { return u.mul*sum + u.add*length },
		Compose: func(newer, older affine) affine {
			return affine{newer.mul * older.mul, newer.mul*older.add + newer.add}
		},
		NoUpdate: affine{1, 0},
	}
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 200 {
		n := 1 + rng.IntN(40)
		model := make([]int, n)
		for i := range model {
			model[i] = rng.IntN(11) - 5
		}
		s, err := NewLazySegmentTree(model, ops)
		if err != nil {
			t.Fatal(err)
		}
		for op := range 40 {
			lo := rng.IntN(n + 1)
			hi := lo + rng.IntN(n-lo+1)
			if rng.IntN(2) == 0 {
				u := affine{rng.IntN(3) - 1, rng.IntN(7) - 3} // includes assignment (mul 0)
				if err := s.RangeUpdate(lo, hi, u); err != nil {
					t.Fatal(err)
				}
				for i := lo; i < hi; i++ {
					model[i] = u.mul*model[i] + u.add
				}
				continue
			}
			want := 0
			for _, v := range model[lo:hi] {
				want += v
			}
			if got, err := s.Query(lo, hi); err != nil || got != want {
				t.Fatalf("trial %d op %d: Query(%d, %d) = (%d, %v), want %d", trial, op, lo, hi, got, err, want)
			}
		}
	}
}

func TestLazySegmentTree_emptyAndOne(t *testing.T) {
	s, err := NewRangeAddSumTree[int](nil)
	if err != nil || s.Len() != 0 {
		t.Fatalf("empty: (%v, %v)", s, err)
	}
	if v, err := s.Query(0, 0); err != nil || v != 0 {
		t.Errorf("empty Query = (%d, %v)", v, err)
	}
	if err := s.RangeUpdate(0, 0, 5); err != nil {
		t.Errorf("empty RangeUpdate: %v", err)
	}
	one, _ := NewRangeAddMinTree([]int{3})
	_ = one.RangeUpdate(0, 1, -10)
	if v, _ := one.Query(0, 1); v != -7 {
		t.Errorf("one: %d, want −7", v)
	}
	var nilTree *LazySegmentTree[int, int]
	if nilTree.Len() != 0 || nilTree.RangeUpdate(0, 0, 1) != nil {
		t.Error("nil: Len or empty RangeUpdate")
	}
	if _, err := nilTree.Query(0, 0); !errors.Is(err, ErrEmptyRange) {
		t.Errorf("nil Query(0, 0): %v, want ErrEmptyRange", err)
	}
}

func TestLazySegmentTree_errors(t *testing.T) {
	good := LazyOps[int, int]{
		Combine: func(l, r int) int { return l + r },
		Apply:   func(u, s, n int) int { return s + u*n },
		Compose: func(a, b int) int { return a + b },
	}
	for _, mutate := range []func(*LazyOps[int, int]){
		func(o *LazyOps[int, int]) { o.Combine = nil },
		func(o *LazyOps[int, int]) { o.Apply = nil },
		func(o *LazyOps[int, int]) { o.Compose = nil },
	} {
		ops := good
		mutate(&ops)
		if s, err := NewLazySegmentTree([]int{1}, ops); !errors.Is(err, ErrNilFunc) || s != nil {
			t.Errorf("nil func: (%v, %v), want ErrNilFunc", s, err)
		}
	}
	s, _ := NewRangeAddSumTree([]int{1, 2, 3})
	for _, r := range [][2]int{{-1, 1}, {2, 1}, {0, 4}} {
		if _, err := s.Query(r[0], r[1]); !errors.Is(err, ErrRange) {
			t.Errorf("Query(%d, %d): %v, want ErrRange", r[0], r[1], err)
		}
		if err := s.RangeUpdate(r[0], r[1], 1); !errors.Is(err, ErrRange) {
			t.Errorf("RangeUpdate(%d, %d): %v, want ErrRange", r[0], r[1], err)
		}
	}
	for _, mk := range []func([]float64) (*LazySegmentTree[float64, float64], error){
		NewRangeAddSumTree[float64], NewRangeAddMinTree[float64], NewRangeAddMaxTree[float64],
	} {
		if _, err := mk([]float64{1, math.NaN()}); !errors.Is(err, ErrNotFinite) {
			t.Errorf("NaN value: %v, want ErrNotFinite", err)
		}
		f, _ := mk([]float64{1, 2})
		if err := f.RangeUpdate(0, 2, math.Inf(1)); !errors.Is(err, ErrNotFinite) {
			t.Errorf("Inf update: %v, want ErrNotFinite", err)
		}
	}
	if _, err := NewRangeAddSumTree([]int8{100, 100}); !errors.Is(err, ErrOverflow) {
		t.Errorf("int8 200: %v, want ErrOverflow", err)
	}
}

// An update that overflows partway (after pushes and writes) leaves every node and pending update as it was.
func TestLazySegmentTreeRangeUpdate_atomicOnOverflow(t *testing.T) {
	s, _ := NewRangeAddSumTree([]int8{1, 2, 3, 4, 5, 6, 7, 8})
	if err := s.RangeUpdate(0, 8, 1); err != nil { // pending +1 at the root: 2 … 9, sum 44
		t.Fatal(err)
	}
	nodes, lazy := slices.Clone(s.node), slices.Clone(s.lazy)
	// +12 on [1, 8) makes the root 44 + 84 = 128, one past int8: found after pushes and writes below.
	if err := s.RangeUpdate(1, 8, 12); !errors.Is(err, ErrOverflow) { // 44 + 84 = 128
		t.Fatalf("RangeUpdate: %v, want ErrOverflow", err)
	}
	if !slices.Equal(s.node, nodes) || !slices.Equal(s.lazy, lazy) {
		t.Errorf("tree changed on error:\n nodes %v (was %v)\n lazy %v (was %v)", s.node, nodes, s.lazy, lazy)
	}
	if v, _ := s.Query(0, 8); v != 44 {
		t.Errorf("Query after failed update = %d, want 44", v)
	}
	// delta·length itself overflowing T.
	m, _ := NewRangeAddSumTree([]int8{0, 0, 0, 0})
	if err := m.RangeUpdate(0, 4, 40); !errors.Is(err, ErrOverflow) { // 40·4 = 160
		t.Errorf("40·4 in int8: %v, want ErrOverflow", err)
	}
}

// Query leaves the tree untouched (read-only), even when pending updates are on its path.
func TestLazySegmentTreeQuery_readOnly(t *testing.T) {
	s, _ := NewRangeAddSumTree(laaksonen16)
	_ = s.RangeUpdate(5, 14, 2)
	nodes, lazy := slices.Clone(s.node), slices.Clone(s.lazy)
	for lo := range 16 {
		for hi := lo + 1; hi <= 16; hi++ {
			_, _ = s.Query(lo, hi)
		}
	}
	if !slices.Equal(s.node, nodes) || !slices.Equal(s.lazy, lazy) {
		t.Error("Query modified the tree")
	}
}

func TestLazySegmentTree_floats(t *testing.T) {
	// Quarters are exact in binary floating point: exact comparison is safe.
	s, _ := NewRangeAddSumTree([]float64{0.25, 0.5, -1})
	_ = s.RangeUpdate(0, 3, 0.25)
	if v, _ := s.Query(0, 3); v != 0.5 {
		t.Errorf("float sum = %v, want 0.5", v)
	}
}

// Oracle: random range adds and queries on all three ready-made trees against a plain array.
func TestLazySegmentTree_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		n := 1 + rng.IntN(70)
		model := make([]int, n)
		for i := range model {
			model[i] = rng.IntN(201) - 100
		}
		sum, _ := NewRangeAddSumTree(model)
		mn, _ := NewRangeAddMinTree(model)
		mx, _ := NewRangeAddMaxTree(model)
		for op := range 80 {
			lo := rng.IntN(n)
			hi := lo + 1 + rng.IntN(n-lo)
			if rng.IntN(2) == 0 {
				d := rng.IntN(41) - 20
				for _, s := range []*LazySegmentTree[int, int]{sum, mn, mx} {
					if err := s.RangeUpdate(lo, hi, d); err != nil {
						t.Fatal(err)
					}
				}
				for i := lo; i < hi; i++ {
					model[i] += d
				}
				continue
			}
			want := 0
			for _, v := range model[lo:hi] {
				want += v
			}
			got, err := sum.Query(lo, hi)
			gotMin, _ := mn.Query(lo, hi)
			gotMax, _ := mx.Query(lo, hi)
			if err != nil || got != want || gotMin != slices.Min(model[lo:hi]) || gotMax != slices.Max(model[lo:hi]) {
				t.Fatalf("trial %d op %d: [%d, %d) sum/min/max = %d/%d/%d, want %d/%d/%d", trial, op, lo, hi,
					got, gotMin, gotMax, want, slices.Min(model[lo:hi]), slices.Max(model[lo:hi]))
			}
		}
	}
}

func TestLazySegmentTree_noAllocation(t *testing.T) {
	s, _ := NewRangeAddSumTree(make([]int, 1000))
	if a := testing.AllocsPerRun(100, func() {
		_ = s.RangeUpdate(17, 903, 3)
		_, _ = s.Query(100, 900)
		_ = s.Len()
	}); a != 0 {
		t.Errorf("allocations per run = %v, want 0", a)
	}
}

func BenchmarkNewLazySegmentTree(b *testing.B) {
	ops := LazyOps[int, int]{
		Combine: func(l, r int) int { return max(l, r) },
		Apply:   func(u, agg, _ int) int { return agg + u },
		Compose: func(a, c int) int { return a + c },
	}
	for _, size := range benchSizes {
		values := segValues(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = NewLazySegmentTree(values, ops)
			}
		})
	}
}

func benchmarkLazyCtor(b *testing.B, mk func([]int) (*LazySegmentTree[int, int], error)) {
	for _, size := range benchSizes {
		values := segValues(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = mk(values)
			}
		})
	}
}

func BenchmarkNewRangeAddSumTree(b *testing.B) { benchmarkLazyCtor(b, NewRangeAddSumTree[int]) }
func BenchmarkNewRangeAddMinTree(b *testing.B) { benchmarkLazyCtor(b, NewRangeAddMinTree[int]) }
func BenchmarkNewRangeAddMaxTree(b *testing.B) { benchmarkLazyCtor(b, NewRangeAddMaxTree[int]) }

func BenchmarkLazySegmentTreeRangeUpdate(b *testing.B) {
	for _, size := range benchSizes {
		s, _ := NewRangeAddSumTree(segValues(size.n))
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = s.RangeUpdate(1, size.n-1, 1)
				_ = s.RangeUpdate(1, size.n-1, -1)
			}
		})
	}
}

func BenchmarkLazySegmentTreeQuery(b *testing.B) {
	for _, size := range benchSizes {
		s, _ := NewRangeAddSumTree(segValues(size.n))
		_ = s.RangeUpdate(size.n/3, size.n/2, 5)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = s.Query(1, size.n-1)
			}
		})
	}
}

func BenchmarkLazySegmentTreeLen(b *testing.B) {
	s, _ := NewRangeAddSumTree(segValues(10))
	b.ReportAllocs()
	for b.Loop() {
		_ = s.Len()
	}
}

// Fix M3: the min tree also keeps maximums (and the max tree minimums), so an add that pushes an element
// that is not the node's aggregate out of T is rejected and nothing changes.
func TestNewRangeAddMinTree_rejectsOverflowOfNonExtreme(t *testing.T) {
	for _, mk := range []func([]int8) (*LazySegmentTree[int8, int8], error){NewRangeAddMinTree[int8], NewRangeAddMaxTree[int8]} {
		s, _ := mk([]int8{100, -100})
		node, lazy, opp := slices.Clone(s.node), slices.Clone(s.lazy), slices.Clone(s.opp)
		if err := s.RangeUpdate(0, 2, 50); !errors.Is(err, ErrOverflow) { // 100 + 50 = 150
			t.Errorf("+50 on [100, −100]: %v, want ErrOverflow", err)
		}
		if err := s.RangeUpdate(0, 2, -50); !errors.Is(err, ErrOverflow) { // −100 − 50 = −150
			t.Errorf("−50 on [100, −100]: %v, want ErrOverflow", err)
		}
		if !slices.Equal(s.node, node) || !slices.Equal(s.lazy, lazy) || !slices.Equal(s.opp, opp) {
			t.Error("tree changed on error")
		}
		if v, err := s.Query(0, 1); err != nil || v != 100 {
			t.Errorf("Query(0, 1) = (%d, %v), want 100", v, err)
		}
		if err := s.RangeUpdate(0, 2, 27); err != nil { // 127 and −73 fit
			t.Fatalf("+27: %v", err)
		}
		if v, _ := s.Query(0, 1); v != 127 {
			t.Errorf("after +27: Query(0, 1) = %d, want 127", v)
		}
		if v, _ := s.Query(1, 2); v != -73 {
			t.Errorf("after +27: Query(1, 2) = %d, want −73", v)
		}
	}
}

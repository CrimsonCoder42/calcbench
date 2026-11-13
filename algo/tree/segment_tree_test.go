package tree

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"strings"
	"testing"
)

// laaksonen16 is the 16-element array of Laaksonen (2018) §9.3 p. 89 and §28.1 p. 258 (total 72).
var laaksonen16 = []int{5, 8, 6, 3, 2, 7, 2, 6, 7, 1, 7, 5, 6, 2, 3, 2}

// Laaksonen §9.3 p. 89: in the array 5 8 6 3 2 7 2 6, sum(2, 7) = 6 + 3 + 2 + 7 + 2 + 6 = 26 (0-based,
// inclusive); here the half-open range [2, 8).
func TestNewSumSegmentTree_laaksonenExample(t *testing.T) {
	s, err := NewSumSegmentTree(laaksonen16[:8])
	if err != nil {
		t.Fatal(err)
	}
	if got, err := s.Query(2, 8); err != nil || got != 26 {
		t.Errorf("Query(2, 8) = (%d, %v), want 26", got, err)
	}
	full, _ := NewSumSegmentTree(laaksonen16)
	if got, _ := full.Query(0, 16); got != 72 {
		t.Errorf("Query(0, 16) = %d, want 72", got)
	}
}

func TestSegmentTreeSet_updatesQueries(t *testing.T) {
	s, _ := NewSumSegmentTree(laaksonen16[:8])
	if err := s.Set(5, 10); err != nil { // 7 → 10 inside [2, 8)
		t.Fatal(err)
	}
	if got, _ := s.Query(2, 8); got != 29 {
		t.Errorf("Query(2, 8) after Set(5, 10) = %d, want 29", got)
	}
}

func TestNewMinSegmentTree_cases(t *testing.T) {
	s, _ := NewMinSegmentTree(laaksonen16)
	for _, tt := range []struct{ lo, hi, want int }{{0, 16, 1}, {0, 9, 2}, {0, 1, 5}, {10, 16, 2}, {1, 3, 6}} {
		if got, err := s.Query(tt.lo, tt.hi); err != nil || got != tt.want {
			t.Errorf("min Query(%d, %d) = (%d, %v), want %d", tt.lo, tt.hi, got, err, tt.want)
		}
	}
	if _, err := s.Query(3, 3); !errors.Is(err, ErrEmptyRange) {
		t.Errorf("empty min range: %v, want ErrEmptyRange", err)
	}
}

func TestNewMaxSegmentTree_cases(t *testing.T) {
	s, _ := NewMaxSegmentTree(laaksonen16)
	for _, tt := range []struct{ lo, hi, want int }{{0, 16, 8}, {2, 8, 7}, {13, 16, 3}, {15, 16, 2}} {
		if got, err := s.Query(tt.lo, tt.hi); err != nil || got != tt.want {
			t.Errorf("max Query(%d, %d) = (%d, %v), want %d", tt.lo, tt.hi, got, err, tt.want)
		}
	}
	if _, err := s.Query(0, 0); !errors.Is(err, ErrEmptyRange) {
		t.Errorf("empty max range: %v, want ErrEmptyRange", err)
	}
}

func TestNewMinSegmentTree_floatsAndStrings(t *testing.T) {
	f, _ := NewMinSegmentTree([]float64{3, math.NaN(), -1})
	if got, _ := f.Query(0, 3); !math.IsNaN(got) {
		t.Errorf("min with NaN = %v, want NaN (least under cmp.Compare)", got)
	}
	g, _ := NewMaxSegmentTree([]string{"pear", "fig", "plum"})
	if got, _ := g.Query(0, 2); got != "pear" {
		t.Errorf("max string = %q, want pear", got)
	}
}

// A non-commutative monoid (string concatenation) proves the left-to-right combine order.
func TestNewSegmentTree_nonCommutative(t *testing.T) {
	words := []string{"a", "b", "c", "d", "e", "f", "g"}
	s, err := NewSegmentTree(words, func(l, r string) string { return l + r }, "")
	if err != nil {
		t.Fatal(err)
	}
	for lo := 0; lo <= len(words); lo++ {
		for hi := lo; hi <= len(words); hi++ {
			if got, err := s.Query(lo, hi); err != nil || got != strings.Join(words[lo:hi], "") {
				t.Errorf("Query(%d, %d) = (%q, %v)", lo, hi, got, err)
			}
		}
	}
	_ = s.Set(3, "X")
	if got, _ := s.Query(0, 7); got != "abcXefg" {
		t.Errorf("after Set: %q", got)
	}
}

func TestSegmentTree_emptyAndOne(t *testing.T) {
	s, err := NewSumSegmentTree[int](nil)
	if err != nil || s.Len() != 0 {
		t.Fatalf("empty: (%v, %v)", s, err)
	}
	if v, err := s.Query(0, 0); err != nil || v != 0 {
		t.Errorf("empty Query(0, 0) = (%d, %v), want 0", v, err)
	}
	one, _ := NewMinSegmentTree([]int{-4})
	if v, _ := one.Query(0, 1); v != -4 {
		t.Errorf("one: %d", v)
	}
	if err := one.Set(0, 9); err != nil {
		t.Fatal(err)
	}
	if v, _ := one.Query(0, 1); v != 9 {
		t.Errorf("one after Set: %d", v)
	}
	var nilTree *SegmentTree[int]
	var zero SegmentTree[int]
	if nilTree.Len() != 0 || zero.Len() != 0 {
		t.Error("nil or zero Len != 0")
	}
	if _, err := nilTree.Query(0, 0); !errors.Is(err, ErrEmptyRange) {
		t.Errorf("nil Query(0, 0): %v, want ErrEmptyRange (no identity)", err)
	}
	if err := nilTree.Set(0, 1); !errors.Is(err, ErrRange) {
		t.Errorf("nil Set: %v, want ErrRange", err)
	}
}

func TestSegmentTree_errors(t *testing.T) {
	if _, err := NewSegmentTree([]int{1}, nil, 0); !errors.Is(err, ErrNilFunc) {
		t.Errorf("nil combine: %v", err)
	}
	s, _ := NewSumSegmentTree([]int{1, 2, 3})
	for _, r := range [][2]int{{-1, 1}, {2, 1}, {0, 4}} {
		if _, err := s.Query(r[0], r[1]); !errors.Is(err, ErrRange) {
			t.Errorf("Query(%d, %d): %v, want ErrRange", r[0], r[1], err)
		}
	}
	for _, i := range []int{-1, 3} {
		if err := s.Set(i, 0); !errors.Is(err, ErrRange) {
			t.Errorf("Set(%d): %v, want ErrRange", i, err)
		}
	}
	if _, err := NewSumSegmentTree([]int8{100, 100}); !errors.Is(err, ErrOverflow) {
		t.Errorf("int8 200: %v, want ErrOverflow", err)
	}
	if _, err := NewSumSegmentTree([]float64{math.NaN()}); !errors.Is(err, ErrNotFinite) {
		t.Errorf("NaN: %v, want ErrNotFinite", err)
	}
	f, _ := NewSumSegmentTree([]float64{1, 2})
	if err := f.Set(0, math.Inf(1)); !errors.Is(err, ErrNotFinite) {
		t.Errorf("Set(Inf): %v, want ErrNotFinite", err)
	}
	if _, err := NewSumSegmentTree([]float64{math.MaxFloat64, math.MaxFloat64}); !errors.Is(err, ErrOverflow) {
		t.Errorf("float sum to Inf: %v, want ErrOverflow", err)
	}
}

// A Set whose path overflows leaves the tree as it was; queries that overflow report it.
func TestSegmentTreeSet_atomicOnOverflow(t *testing.T) {
	s, _ := NewSumSegmentTree([]int8{1, 2, 3, 120})
	before := slices.Clone(s.node)
	if err := s.Set(0, 10); !errors.Is(err, ErrOverflow) { // the root would be 10+2+3+120 = 135
		t.Fatalf("Set: %v, want ErrOverflow", err)
	}
	if !slices.Equal(s.node, before) {
		t.Errorf("nodes changed on error: %v, was %v", s.node, before)
	}
	// Nodes [0,2) = 127, [2,4) = −97 and the root 30 fit, but [1,3) is formed as 126 + 3.
	t2, err := NewSumSegmentTree([]int8{1, 126, 3, -100})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := t2.Query(1, 3); !errors.Is(err, ErrOverflow) {
		t.Errorf("Query(1, 3) of 126 + 3 in int8: %v, want ErrOverflow", err)
	}
}

func TestSegmentTree_extremes(t *testing.T) {
	s, err := NewSumSegmentTree([]int64{math.MaxInt64, math.MinInt64, 0})
	if err != nil {
		t.Fatal(err)
	}
	if v, err := s.Query(0, 3); err != nil || v != -1 {
		t.Errorf("Query(0, 3) = (%d, %v), want −1", v, err)
	}
	u, err := NewSumSegmentTree([]uint8{200, 55})
	if err != nil {
		t.Fatal(err)
	}
	if v, _ := u.Query(0, 2); v != 255 {
		t.Errorf("uint8 sum = %d, want 255", v)
	}
	if err := u.Set(1, 56); !errors.Is(err, ErrOverflow) {
		t.Errorf("uint8 256: %v, want ErrOverflow", err)
	}
}

func TestSegmentTree_doesNotModifyInput(t *testing.T) {
	values := slices.Clone(laaksonen16)
	s, _ := NewSumSegmentTree(values)
	_ = s.Set(0, 100)
	if !slices.Equal(values, laaksonen16) {
		t.Error("input modified")
	}
}

// Oracle: random sets and queries for sum, min and max against a plain array.
func TestSegmentTree_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		n := 1 + rng.IntN(70)
		model := make([]int, n)
		for i := range model {
			model[i] = rng.IntN(201) - 100
		}
		sum, _ := NewSumSegmentTree(model)
		mn, _ := NewMinSegmentTree(model)
		mx, _ := NewMaxSegmentTree(model)
		for op := range 80 {
			if rng.IntN(3) == 0 {
				i, v := rng.IntN(n), rng.IntN(201)-100
				for _, s := range []*SegmentTree[int]{sum, mn, mx} {
					if err := s.Set(i, v); err != nil {
						t.Fatal(err)
					}
				}
				model[i] = v
				continue
			}
			lo := rng.IntN(n)
			hi := lo + 1 + rng.IntN(n-lo)
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

func TestSegmentTree_noAllocation(t *testing.T) {
	s, _ := NewSumSegmentTree(make([]int, 1000))
	if a := testing.AllocsPerRun(100, func() {
		_ = s.Set(17, 3)
		_, _ = s.Query(100, 900)
		_ = s.Len()
	}); a != 0 {
		t.Errorf("allocations per run = %v, want 0", a)
	}
}

func segValues(n int) []int {
	values := make([]int, n)
	for i := range values {
		values[i] = i % 97
	}
	return values
}

func BenchmarkNewSegmentTree(b *testing.B) {
	for _, size := range benchSizes {
		values := segValues(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = NewSegmentTree(values, func(l, r int) int { return l ^ r }, 0)
			}
		})
	}
}

func BenchmarkNewSumSegmentTree(b *testing.B) {
	for _, size := range benchSizes {
		values := segValues(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = NewSumSegmentTree(values)
			}
		})
	}
}

func BenchmarkNewMinSegmentTree(b *testing.B) {
	for _, size := range benchSizes {
		values := segValues(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = NewMinSegmentTree(values)
			}
		})
	}
}

func BenchmarkNewMaxSegmentTree(b *testing.B) {
	for _, size := range benchSizes {
		values := segValues(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = NewMaxSegmentTree(values)
			}
		})
	}
}

func BenchmarkSegmentTreeSet(b *testing.B) {
	for _, size := range benchSizes {
		s, _ := NewSumSegmentTree(segValues(size.n))
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			i := 0
			for b.Loop() {
				_ = s.Set(i%size.n, i&63)
				i += 7919
			}
		})
	}
}

func BenchmarkSegmentTreeQuery(b *testing.B) {
	for _, size := range benchSizes {
		s, _ := NewSumSegmentTree(segValues(size.n))
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = s.Query(1, size.n-1)
			}
		})
	}
}

func BenchmarkSegmentTreeLen(b *testing.B) {
	s, _ := NewSumSegmentTree(segValues(10))
	b.ReportAllocs()
	for b.Loop() {
		_ = s.Len()
	}
}

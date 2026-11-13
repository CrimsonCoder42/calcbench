package tree

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

func allValues[T signed](t *testing.T, f *FenwickRangeUpdate[T]) []T {
	t.Helper()
	out := make([]T, f.Len())
	for i := range out {
		v, err := f.Get(i)
		if err != nil {
			t.Fatalf("Get(%d): %v", i, err)
		}
		out[i] = v
	}
	return out
}

// Laaksonen (2018) §9.4 pp. 93–94: the array 3 3 1 1 1 5 2 2; increasing positions 1…4 by 5 gives
// 3 8 6 6 6 5 2 2 (its difference array 3 5 −2 0 0 −1 −3 0, as printed).
func TestFenwickRangeUpdate_laaksonenExample(t *testing.T) {
	f, err := NewFenwickRangeUpdate([]int{3, 3, 1, 1, 1, 5, 2, 2})
	if err != nil {
		t.Fatal(err)
	}
	if err := f.RangeAdd(1, 5, 5); err != nil {
		t.Fatal(err)
	}
	if got, want := allValues(t, f), []int{3, 8, 6, 6, 6, 5, 2, 2}; !slices.Equal(got, want) {
		t.Errorf("values = %v, want %v", got, want)
	}
	// The value 2 at position 6 is the sum of differences 3 − 2 + 4 − 3 (before the update; p. 93).
	g, _ := NewFenwickRangeUpdate([]int{3, 3, 1, 1, 1, 5, 2, 2})
	if v, _ := g.Get(6); v != 2 {
		t.Errorf("Get(6) = %d, want 2", v)
	}
}

func TestFenwickRangeUpdateRangeAdd_cases(t *testing.T) {
	tests := []struct {
		name   string
		lo, hi int
		delta  int
		want   []int
	}{
		{"empty range", 2, 2, 9, []int{0, 0, 0, 0, 0}},
		{"zero delta", 0, 5, 0, []int{0, 0, 0, 0, 0}},
		{"whole array", 0, 5, -3, []int{-3, -3, -3, -3, -3}},
		{"one element", 3, 4, 7, []int{0, 0, 0, 7, 0}},
		{"to the end", 2, 5, 1, []int{0, 0, 1, 1, 1}},
		{"from the start", 0, 2, 4, []int{4, 4, 0, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, _ := NewFenwickRangeUpdate(make([]int, 5))
			if err := f.RangeAdd(tt.lo, tt.hi, tt.delta); err != nil {
				t.Fatal(err)
			}
			if got := allValues(t, f); !slices.Equal(got, tt.want) {
				t.Errorf("values = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFenwickRangeUpdate_errors(t *testing.T) {
	f, _ := NewFenwickRangeUpdate([]int{1, 2, 3})
	for _, r := range [][2]int{{-1, 1}, {2, 1}, {0, 4}} {
		if err := f.RangeAdd(r[0], r[1], 1); !errors.Is(err, ErrRange) {
			t.Errorf("RangeAdd(%d, %d): %v, want ErrRange", r[0], r[1], err)
		}
	}
	for _, i := range []int{-1, 3} {
		if _, err := f.Get(i); !errors.Is(err, ErrRange) {
			t.Errorf("Get(%d): %v, want ErrRange", i, err)
		}
	}
	if _, err := NewFenwickRangeUpdate([]int8{-100, 100}); !errors.Is(err, ErrOverflow) {
		t.Errorf("difference 200 in int8: %v, want ErrOverflow", err)
	}
	if _, err := NewFenwickRangeUpdate([]float64{math.Inf(-1)}); !errors.Is(err, ErrNotFinite) {
		t.Errorf("Inf value: %v, want ErrNotFinite", err)
	}
	g, _ := NewFenwickRangeUpdate([]float64{1})
	if err := g.RangeAdd(0, 1, math.NaN()); !errors.Is(err, ErrNotFinite) {
		t.Errorf("NaN delta: %v, want ErrNotFinite", err)
	}
	var nilF *FenwickRangeUpdate[int]
	if nilF.Len() != 0 || nilF.RangeAdd(0, 0, 1) != nil {
		t.Error("nil: Len != 0 or empty RangeAdd failed")
	}
	if _, err := nilF.Get(0); !errors.Is(err, ErrRange) {
		t.Errorf("nil Get: %v, want ErrRange", err)
	}
}

// A RangeAdd whose subtraction would overflow leaves the additions undone too.
func TestFenwickRangeUpdateRangeAdd_atomicOnOverflow(t *testing.T) {
	f, _ := NewFenwickRangeUpdate([]int8{0, 0, 0, 0, -100, 0, 0, 0}) // d[4] = −100
	before := slices.Clone(f.diff.t)
	// Nodes 1, 2, 4 would gain 100, then node 5 (d[4]) would become −200: nothing may change.
	if err := f.RangeAdd(0, 4, 100); !errors.Is(err, ErrOverflow) {
		t.Fatalf("RangeAdd: %v, want ErrOverflow", err)
	}
	if !slices.Equal(f.diff.t, before) {
		t.Errorf("tree changed on error: %v, was %v", f.diff.t, before)
	}
}

// A node covering both ends of the range is left alone, so no spurious overflow there.
func TestFenwickRangeUpdateRangeAdd_sharedNodeUntouched(t *testing.T) {
	f, _ := NewFenwickRangeUpdate([]int8{0, 0, 0, 0, 0, 0, 0, 127}) // node 8 holds a[7] − a[−1] = 127
	if err := f.RangeAdd(1, 3, 100); err != nil {
		t.Fatalf("RangeAdd: %v (nodes 4 and 8 see +100 − 100)", err)
	}
	if got, want := allValues(t, f), []int8{0, 100, 100, 0, 0, 0, 0, 127}; !slices.Equal(got, want) {
		t.Errorf("values = %v, want %v", got, want)
	}
}

func TestFenwickRangeUpdate_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		n := 1 + rng.IntN(60)
		model := make([]int, n)
		for i := range model {
			model[i] = rng.IntN(201) - 100
		}
		f, err := NewFenwickRangeUpdate(model)
		if err != nil {
			t.Fatal(err)
		}
		for op := range 60 {
			lo := rng.IntN(n + 1)
			hi := lo + rng.IntN(n-lo+1)
			d := rng.IntN(201) - 100
			if err := f.RangeAdd(lo, hi, d); err != nil {
				t.Fatal(err)
			}
			for i := lo; i < hi; i++ {
				model[i] += d
			}
			i := rng.IntN(n)
			if got, err := f.Get(i); err != nil || got != model[i] {
				t.Fatalf("trial %d op %d: Get(%d) = (%d, %v), want %d", trial, op, i, got, err, model[i])
			}
		}
	}
}

func TestFenwickRangeUpdate_noAllocation(t *testing.T) {
	f, _ := NewFenwickRangeUpdate(make([]int, 1000))
	if a := testing.AllocsPerRun(100, func() {
		_ = f.RangeAdd(10, 900, 2)
		_, _ = f.Get(500)
		_ = f.Len()
	}); a != 0 {
		t.Errorf("allocations per run = %v, want 0", a)
	}
}

func BenchmarkNewFenwickRangeUpdate(b *testing.B) {
	for _, size := range benchSizes {
		values := make([]int, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = NewFenwickRangeUpdate(values)
			}
		})
	}
}

func BenchmarkFenwickRangeUpdateRangeAdd(b *testing.B) {
	for _, size := range benchSizes {
		f, _ := NewFenwickRangeUpdate(make([]int, size.n))
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = f.RangeAdd(1, size.n-1, 1)
				_ = f.RangeAdd(1, size.n-1, -1)
			}
		})
	}
}

func BenchmarkFenwickRangeUpdateGet(b *testing.B) {
	for _, size := range benchSizes {
		f, _ := NewFenwickRangeUpdate(make([]int, size.n))
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = f.Get(size.n - 1)
			}
		})
	}
}

func BenchmarkFenwickRangeUpdateLen(b *testing.B) {
	f, _ := NewFenwickRangeUpdate(make([]int, 10))
	b.ReportAllocs()
	for b.Loop() {
		_ = f.Len()
	}
}

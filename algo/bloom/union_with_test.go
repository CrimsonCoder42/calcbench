package bloom

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"slices"
	"testing"
)

// Broder–Mitzenmacher §2.4: the OR of two filters with the same bits and hash functions is the filter of
// the union. Hand-built: "apple" in one filter, "banana" in the other (m = 64, k = 3, bits as in
// TestBloomFilter_example), union = bits 13 35 41 48 56 58.
func TestUnionWith_example(t *testing.T) {
	a, b := mustFilter(t, 64, 3), mustFilter(t, 64, 3)
	a.Add([]byte("apple"))
	b.Add([]byte("banana"))
	if err := a.UnionWith(b); err != nil {
		t.Fatal(err)
	}
	if got, want := setBitList(a), []uint64{13, 35, 41, 48, 56, 58}; !slices.Equal(got, want) || a.SetBits() != 6 {
		t.Errorf("union bits = %v (SetBits %d), want %v", got, a.SetBits(), want)
	}
	if got := setBitList(b); !slices.Equal(got, []uint64{13, 35, 58}) {
		t.Errorf("other modified: %v", got)
	}
}

// Oracle: the union equals one filter fed both item sets, word for word, with the same SetBits.
func TestUnionWith_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		m := 1 + rng.Uint64N(3000)
		k := 1 + rng.IntN(10)
		a, b, both := mustFilter(t, m, k), mustFilter(t, m, k), mustFilter(t, m, k)
		for i := range rng.IntN(60) {
			item := fmt.Appendf(nil, "a-%d-%d", trial, i)
			a.Add(item)
			both.Add(item)
		}
		for i := range rng.IntN(60) {
			item := fmt.Appendf(nil, "b-%d-%d", trial, i)
			b.Add(item)
			both.Add(item)
		}
		if err := a.UnionWith(b); err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(a.words, both.words) || a.SetBits() != both.SetBits() {
			t.Fatalf("trial %d (m=%d, k=%d): union differs from the filter of both sets (SetBits %d vs %d)", trial, m, k, a.SetBits(), both.SetBits())
		}
	}
}

func TestUnionWith_cases(t *testing.T) {
	t.Run("with itself", func(t *testing.T) {
		f := mustFilter(t, 100, 3)
		f.Add([]byte("x"))
		before, n := slices.Clone(f.words), f.SetBits()
		if err := f.UnionWith(f); err != nil || !slices.Equal(before, f.words) || f.SetBits() != n {
			t.Errorf("self-union changed the filter or failed: %v", err)
		}
	})
	t.Run("with an empty filter", func(t *testing.T) {
		f, empty := mustFilter(t, 100, 3), mustFilter(t, 100, 3)
		f.Add([]byte("x"))
		n := f.SetBits()
		if err := f.UnionWith(empty); err != nil || f.SetBits() != n || !f.MayContain([]byte("x")) {
			t.Errorf("union with an empty filter: SetBits %d → %d, err %v", n, f.SetBits(), err)
		}
	})
	t.Run("zero values", func(t *testing.T) {
		var f, g BloomFilter
		if err := f.UnionWith(&g); err != nil || f.SetBits() != 0 {
			t.Errorf("zero-value union: %v", err)
		}
	})
}

func TestUnionWith_errors(t *testing.T) {
	f := mustFilter(t, 100, 3)
	f.Add([]byte("x"))
	before := slices.Clone(f.words)
	for name, other := range map[string]*BloomFilter{
		"nil":              nil,
		"other bit count":  mustFilter(t, 101, 3),
		"other hash count": mustFilter(t, 100, 4),
		"zero-value other": {},
	} {
		t.Run(name, func(t *testing.T) {
			if err := f.UnionWith(other); !errors.Is(err, ErrIncompatible) {
				t.Errorf("UnionWith(%s) = %v, want ErrIncompatible", name, err)
			}
			if !slices.Equal(before, f.words) {
				t.Errorf("UnionWith(%s) modified the filter", name)
			}
		})
	}
}

func TestUnionWith_noAllocation(t *testing.T) {
	a, b := mustFilter(t, 10_000, 3), mustFilter(t, 10_000, 3)
	if allocs := testing.AllocsPerRun(20, func() { _ = a.UnionWith(b) }); allocs != 0 {
		t.Errorf("UnionWith allocates %v times per call, want 0", allocs)
	}
}

// The size is the number of 64-bit words (m = 64·size bits): O(m/64).
func BenchmarkUnionWith(b *testing.B) {
	for _, size := range benchSizes {
		m := uint64(size.n) * 64
		x, y := mustFilter(b, m, 3), mustFilter(b, m, 3)
		for i := range size.n {
			y.Add(fmt.Appendf(nil, "key-%d", i))
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = x.UnionWith(y)
			}
		})
	}
}

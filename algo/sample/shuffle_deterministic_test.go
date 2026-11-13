package sample

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Knuth, TAOCP Vol. 2, §3.4.2, Algorithm P with a fixed seed: the output is a permutation of 0…n−1 and
// the same seed reproduces it (no golden values: math/rand/v2 does not promise identical draws across
// Go releases).
func TestShuffleDeterministic_sameSeedSamePermutation(t *testing.T) {
	a, err := ShuffleDeterministic(nil, 20, 42)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := ShuffleDeterministic(nil, 20, 42)
	if !slices.Equal(a, b) {
		t.Errorf("same seed, different permutations: %v vs %v", a, b)
	}
	if !slices.Equal(slices.Sorted(slices.Values(a)), func() []int {
		s := make([]int, 20)
		for i := range s {
			s[i] = i
		}
		return s
	}()) {
		t.Errorf("%v is not a permutation of 0…19", a)
	}
	c, _ := ShuffleDeterministic(nil, 20, 43)
	if slices.Equal(a, c) {
		t.Errorf("seeds 42 and 43 gave the same permutation of 20 (probability 1/20!)")
	}
}

func TestShuffleDeterministic_cases(t *testing.T) {
	for _, n := range []int{0, 1, 2, 100} {
		got, err := ShuffleDeterministic(nil, n, 7)
		if err != nil || len(got) != n {
			t.Fatalf("n=%d: len %d, %v", n, len(got), err)
		}
		seen := make([]bool, n)
		for _, v := range got {
			if v < 0 || v >= n || seen[v] {
				t.Fatalf("n=%d: %v is not a permutation", n, got)
			}
			seen[v] = true
		}
	}
	got, _ := ShuffleDeterministic([]int{-1}, 3, 7)
	if len(got) != 4 || got[0] != -1 {
		t.Errorf("append to dst: %v", got)
	}
}

func TestShuffleDeterministic_errors(t *testing.T) {
	dst := []int{5}
	if got, err := ShuffleDeterministic(dst, -1, 1); !errors.Is(err, ErrNegative) || !slices.Equal(got, dst) {
		t.Errorf("negative n: %v, %v", got, err)
	}
	for _, n := range []int{math.MaxInt, maxAllocBytes / 8} {
		if got, err := ShuffleDeterministic(dst, n, 1); !errors.Is(err, ErrTooLarge) || !slices.Equal(got, dst) {
			t.Errorf("n=%d: %v, %v; want ErrTooLarge, no panic", n, got, err)
		}
	}
}

// TestShuffleDeterministic_oracle: it equals ShuffleInPlace on 0…n−1 with the documented generator.
func TestShuffleDeterministic_oracle(t *testing.T) {
	for seed := range uint64(300) {
		n := int(seed % 40)
		got, _ := ShuffleDeterministic(nil, n, seed)
		want := make([]int, n)
		for i := range want {
			want[i] = i
		}
		_ = ShuffleInPlace(newPCGRand(seed), want)
		if !slices.Equal(got, want) {
			t.Fatalf("seed %d: %v, want %v", seed, got, want)
		}
	}
}

func BenchmarkShuffleDeterministic(b *testing.B) {
	for _, size := range benchSizes {
		dst := make([]int, 0, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				dst, _ = ShuffleDeterministic(dst[:0], size.n, 1)
			}
			benchSink = dst[0]
		})
	}
}

// newPCGRand is the generator ShuffleDeterministic documents: PCG seeded (seed, pcgStream).
func newPCGRand(seed uint64) *rand.Rand { return rand.New(rand.NewPCG(seed, pcgStream)) }

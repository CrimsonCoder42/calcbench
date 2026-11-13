package sample

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// benchSink keeps benchmark results observable. Shared by every benchmark in this package.
var benchSink int

// benchSizes are the collection sizes of .claude/rules/go-complexity.md.
var benchSizes = []struct {
	name string
	n    int
}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}}

// within5Sigma reports whether a binomial count of trials with success probability p is within five
// standard deviations of its mean. With fixed seeds the tests are deterministic; 5σ only guards the
// choice of seed against a borderline draw.
func within5Sigma(count, trials int, p float64) bool {
	mean := float64(trials) * p
	sigma := math.Sqrt(float64(trials) * p * (1 - p))
	return math.Abs(float64(count)-mean) <= 5*sigma
}

// Knuth, TAOCP Vol. 2, §3.4.2, Algorithm P: every one of the n! orders is equally likely. Checked for
// n = 3: 60 000 shuffles, each of the 6 orders within 5σ of 10 000.
func TestShuffleInPlace_uniformThree(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	counts := map[string]int{}
	const trials = 60_000
	for range trials {
		v := []string{"a", "b", "c"}
		if err := ShuffleInPlace(rng, v); err != nil {
			t.Fatal(err)
		}
		counts[fmt.Sprint(v)]++
	}
	if len(counts) != 6 {
		t.Fatalf("saw %d orders, want 6", len(counts))
	}
	for k, c := range counts {
		if !within5Sigma(c, trials, 1.0/6) {
			t.Errorf("order %s: %d of %d", k, c, trials)
		}
	}
}

func TestShuffleInPlace_cases(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for _, in := range [][]int{nil, {}, {5}, {1, 2}, {3, 3, 3}, {-1, 0, math.MaxInt, math.MinInt}} {
		v := slices.Clone(in)
		if err := ShuffleInPlace(rng, v); err != nil {
			t.Fatal(err)
		}
		a, b := slices.Sorted(slices.Values(v)), slices.Sorted(slices.Values(in))
		if !slices.Equal(a, b) {
			t.Errorf("ShuffleInPlace(%v) = %v: not a permutation", in, v)
		}
	}
}

func TestShuffleInPlace_errors(t *testing.T) {
	v := []int{1, 2, 3}
	if err := ShuffleInPlace(nil, v); !errors.Is(err, ErrNilRand) {
		t.Errorf("err = %v, want ErrNilRand", err)
	}
	if !slices.Equal(v, []int{1, 2, 3}) {
		t.Errorf("input changed on error: %v", v)
	}
}

func TestShuffleInPlace_noAllocation(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	v := make([]int, 64)
	if allocs := testing.AllocsPerRun(100, func() { _ = ShuffleInPlace(rng, v) }); allocs != 0 {
		t.Errorf("ShuffleInPlace allocates %v times per call, want 0", allocs)
	}
}

// TestShuffleInPlace_oracle: each value lands in each position with probability 1/n (n = 5, 50 000 runs),
// and the result is always a permutation.
func TestShuffleInPlace_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	const n, trials = 5, 50_000
	var at [n][n]int
	for range trials {
		v := []int{0, 1, 2, 3, 4}
		_ = ShuffleInPlace(rng, v)
		for pos, val := range v {
			at[val][pos]++
		}
	}
	for val := range n {
		for pos := range n {
			if !within5Sigma(at[val][pos], trials, 1.0/n) {
				t.Errorf("value %d at position %d: %d of %d", val, pos, at[val][pos], trials)
			}
		}
	}
}

func BenchmarkShuffleInPlace(b *testing.B) {
	for _, size := range benchSizes {
		v := make([]int, size.n)
		rng := rand.New(rand.NewPCG(1, 2))
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = ShuffleInPlace(rng, v)
			}
			benchSink = v[0]
		})
	}
}

package hashing

import (
	"maps"
	"math"
	"math/rand/v2"
	"strings"
	"testing"
)

// benchSizes are the collection sizes of .claude/rules/go-complexity.md. Shared by the benchmarks.
var benchSizes = []struct {
	name string
	n    int
}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}}

// randomInts returns n values drawn from [0, maxValue) with a fixed seed. Shared by the tests and benchmarks.
func randomInts(rng *rand.Rand, n, maxValue int) []int {
	s := make([]int, n)
	for i := range s {
		s[i] = rng.IntN(maxValue)
	}
	return s
}

// Sedgewick & Wayne, Algorithms, 4th ed., §3.1, client FrequencyCounter (count each word of a text),
// applied to our own six words; counts read off by hand.
func TestFrequencyCount_sedgewickWayneClient(t *testing.T) {
	got := FrequencyCount(strings.Fields("to be or not to be"))
	want := map[string]int{"to": 2, "be": 2, "or": 1, "not": 1}
	if !maps.Equal(got, want) {
		t.Errorf("FrequencyCount = %v, want %v", got, want)
	}
}

func TestFrequencyCount_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   map[int]int
	}{
		{"empty", nil, map[int]int{}},
		{"one", []int{7}, map[int]int{7: 1}},
		{"two distinct", []int{7, 8}, map[int]int{7: 1, 8: 1}},
		{"all equal", []int{4, 4, 4}, map[int]int{4: 3}},
		{"negative and zero", []int{-1, 0, -1}, map[int]int{-1: 2, 0: 1}},
		{"extremes", []int{math.MinInt, math.MaxInt, math.MinInt}, map[int]int{math.MinInt: 2, math.MaxInt: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FrequencyCount(tt.values)
			if got == nil || !maps.Equal(got, tt.want) {
				t.Errorf("FrequencyCount(%v) = %#v, want %v", tt.values, got, tt.want)
			}
		})
	}
}

func TestFrequencyCount_floats(t *testing.T) {
	got := FrequencyCount([]float64{math.NaN(), math.NaN(), math.Copysign(0, -1), 0, 1})
	// Each NaN is its own key; -0 and +0 are one key.
	if len(got) != 4 || got[0] != 2 || got[1] != 1 {
		t.Errorf("FrequencyCount(floats) = %v, want 2 NaN keys, 0:2, 1:1", got)
	}
}

// Oracle: count with a nested loop over the distinct values.
func TestFrequencyCount_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		n := rng.IntN(200)
		values := randomInts(rng, n, 1+rng.IntN(n+1))
		got := FrequencyCount(values)
		total := 0
		for v, c := range got {
			want := 0
			for _, w := range values {
				if w == v {
					want++
				}
			}
			if c != want {
				t.Fatalf("trial %d: count of %d = %d, want %d", trial, v, c, want)
			}
			total += c
		}
		if total != n {
			t.Fatalf("trial %d: counts sum to %d, want %d", trial, total, n)
		}
	}
}

func BenchmarkFrequencyCount(b *testing.B) {
	for _, size := range benchSizes {
		values := randomInts(rand.New(rand.NewPCG(1, 2)), size.n, size.n/4)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_ = FrequencyCount(values)
			}
		})
	}
}

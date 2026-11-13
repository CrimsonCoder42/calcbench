package hashing

import (
	"cmp"
	"errors"
	"maps"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// rankAll is the slow oracle: every entry, sorted by count descending then value ascending.
func rankAll[T cmp.Ordered](counts map[T]int) []Count[T] {
	all := make([]Count[T], 0, len(counts))
	for v, c := range counts {
		all = append(all, Count[T]{Value: v, Count: c})
	}
	slices.SortFunc(all, func(a, b Count[T]) int {
		if a.Count != b.Count {
			return cmp.Compare(b.Count, a.Count)
		}
		return cmp.Compare(a.Value, b.Value)
	})
	return all
}

// The ranking rule of the consumer (the consumer backend knowledgebase/series/rank.go:19-37): count descending,
// ties by key ascending. Counts below are our own.
func TestTopKByCount_consumerOrder(t *testing.T) {
	counts := map[string]int{"cough": 4, "fever": 7, "rash": 4, "ache": 2, "chill": 7}
	got, err := TopKByCount(counts, 3)
	want := []Count[string]{{"chill", 7}, {"fever", 7}, {"cough", 4}}
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("TopKByCount = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestTopKByCount_cases(t *testing.T) {
	tests := []struct {
		name   string
		counts map[int]int
		k      int
		want   []Count[int]
	}{
		{"empty, k 0", nil, 0, []Count[int]{}},
		{"empty, k 3", map[int]int{}, 3, []Count[int]{}},
		{"k 0", map[int]int{1: 1}, 0, []Count[int]{}},
		{"one", map[int]int{5: 2}, 1, []Count[int]{{5, 2}}},
		{"k above m", map[int]int{5: 2, 6: 9}, 10, []Count[int]{{6, 9}, {5, 2}}},
		{"all tied", map[int]int{3: 1, 1: 1, 2: 1}, 2, []Count[int]{{1, 1}, {2, 1}}},
		{"zero and negative counts", map[int]int{1: 0, 2: -3, 3: 0}, 3, []Count[int]{{1, 0}, {3, 0}, {2, -3}}},
		{"extreme keys", map[int]int{math.MinInt: 1, math.MaxInt: 1}, 2, []Count[int]{{math.MinInt, 1}, {math.MaxInt, 1}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TopKByCount(tt.counts, tt.k)
			if err != nil || got == nil || !slices.Equal(got, tt.want) {
				t.Errorf("TopKByCount(%v, %d) = (%#v, %v), want (%v, nil)", tt.counts, tt.k, got, err, tt.want)
			}
		})
	}
}

func TestTopKByCount_errors(t *testing.T) {
	for _, k := range []int{-1, math.MinInt} {
		if got, err := TopKByCount(map[int]int{1: 1}, k); !errors.Is(err, ErrNegativeK) || got != nil {
			t.Errorf("TopKByCount(k=%d) = (%v, %v), want (nil, ErrNegativeK)", k, got, err)
		}
	}
}

func TestTopKByCount_doesNotModifyInput(t *testing.T) {
	counts := map[int]int{1: 3, 2: 5, 3: 1}
	before := maps.Clone(counts)
	if _, err := TopKByCount(counts, 2); err != nil {
		t.Fatal(err)
	}
	if !maps.Equal(counts, before) {
		t.Errorf("TopKByCount modified its input: %v", counts)
	}
}

// Oracle: rank every entry and take the first k. Many ties make map order matter if it could leak.
func TestTopKByCount_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		m := rng.IntN(120)
		counts := make(map[int]int, m)
		for range m {
			counts[rng.IntN(500)-250] = rng.IntN(6)
		}
		k := rng.IntN(len(counts) + 3)
		want := rankAll(counts)[:min(k, len(counts))]
		got, err := TopKByCount(counts, k)
		if err != nil || !slices.Equal(got, want) {
			t.Fatalf("trial %d: TopKByCount(m=%d, k=%d) = (%v, %v), want %v", trial, len(counts), k, got, err, want)
		}
	}
}

func BenchmarkTopKByCount(b *testing.B) {
	for _, size := range benchSizes {
		counts := FrequencyCount(randomInts(rand.New(rand.NewPCG(1, 2)), 4*size.n, size.n))
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = TopKByCount(counts, 100)
			}
		})
	}
}

package hashing

import (
	"errors"
	"math/rand/v2"
	"slices"
	"strings"
	"testing"
)

// Sedgewick & Wayne §3.1 FrequencyCounter then §2.4 TopM, on our own words: "be" and "to" both occur
// twice; equal counts rank by value ascending, so "be" comes first. Hand-derived.
func TestTopKFrequent_wordExample(t *testing.T) {
	got, err := TopKFrequent(strings.Fields("to be or not to be"), 2)
	want := []Count[string]{{"be", 2}, {"to", 2}}
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("TopKFrequent = (%v, %v), want (%v, nil)", got, err, want)
	}
}

// Our own wording of the usual exercise: the two most frequent of a short list.
func TestTopKFrequent_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		k      int
		want   []Count[int]
	}{
		{"empty", nil, 2, []Count[int]{}},
		{"k 0", []int{1, 1}, 0, []Count[int]{}},
		{"one value", []int{9}, 1, []Count[int]{{9, 1}}},
		{"two most frequent", []int{4, 4, 4, 6, 6, 8}, 2, []Count[int]{{4, 3}, {6, 2}}},
		{"k above distinct", []int{2, 1, 2}, 5, []Count[int]{{2, 2}, {1, 1}}},
		{"ties by value", []int{9, 3, 7, 3, 9, 7}, 2, []Count[int]{{3, 2}, {7, 2}}},
		{"negative values", []int{-5, -5, 0, -1}, 1, []Count[int]{{-5, 2}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TopKFrequent(tt.values, tt.k)
			if err != nil || got == nil || !slices.Equal(got, tt.want) {
				t.Errorf("TopKFrequent(%v, %d) = (%#v, %v), want (%v, nil)", tt.values, tt.k, got, err, tt.want)
			}
		})
	}
}

func TestTopKFrequent_errors(t *testing.T) {
	if got, err := TopKFrequent([]int{1}, -1); !errors.Is(err, ErrNegativeK) || got != nil {
		t.Errorf("TopKFrequent(k=-1) = (%v, %v), want (nil, ErrNegativeK)", got, err)
	}
}

func TestTopKFrequent_doesNotModifyInput(t *testing.T) {
	values := []int{3, 1, 3, 2}
	if _, err := TopKFrequent(values, 2); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(values, []int{3, 1, 3, 2}) {
		t.Errorf("TopKFrequent modified its input: %v", values)
	}
}

// Oracle: count by FrequencyCount and rank every entry.
func TestTopKFrequent_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		n := rng.IntN(300)
		values := randomInts(rng, n, 1+rng.IntN(n/3+1))
		k := rng.IntN(20)
		all := rankAll(FrequencyCount(values))
		want := all[:min(k, len(all))]
		got, err := TopKFrequent(values, k)
		if err != nil || !slices.Equal(got, want) {
			t.Fatalf("trial %d: TopKFrequent(n=%d, k=%d) = (%v, %v), want %v", trial, n, k, got, err, want)
		}
	}
}

func BenchmarkTopKFrequent(b *testing.B) {
	for _, size := range benchSizes {
		values := randomInts(rand.New(rand.NewPCG(1, 2)), size.n, size.n/4)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = TopKFrequent(values, 10)
			}
		})
	}
}

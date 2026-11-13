package enumerate

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"testing"
)

// Knuth, TAOCP Vol. 4A, §7.2.1.1, Algorithm M with radices (2, 3) for the multiset {1, 2, 2}: the tuples
// (0,0) (1,0) (0,1) (1,1) (0,2) (1,2) give {} 1 2 12 22 122.
func TestSubsetsUnique_knuthExample(t *testing.T) {
	got := collect(SubsetsUnique([]int{2, 1, 2}))
	want := [][]int{{}, {1}, {2}, {1, 2}, {2, 2}, {1, 2, 2}}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if !slices.Equal(got[i], want[i]) {
			t.Errorf("[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}

func TestSubsetsUnique_cases(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want int
	}{
		{"empty", nil, 1},
		{"one", []int{4}, 2},
		{"all equal", []int{4, 4, 4}, 4},
		{"distinct", []int{3, 1, 2}, 8},
		{"negative and zero", []int{0, -1, 0}, 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := slices.Clone(tt.in)
			if got := len(collect(SubsetsUnique(in))); got != tt.want {
				t.Errorf("SubsetsUnique(%v) yields %d, want %d", tt.in, got, tt.want)
			}
			if !slices.Equal(in, tt.in) {
				t.Errorf("input modified: %v", in)
			}
		})
	}
}

func TestSubsetsUnique_earlyBreak(t *testing.T) {
	count := 0
	for range SubsetsUnique(make([]int, 1000)) {
		count++
		if count == 2 {
			break
		}
	}
	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}
}

// TestSubsetsUnique_oracle: the slow version takes every subset, sorts it and deduplicates by key.
func TestSubsetsUnique_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 200 {
		n := rng.IntN(9)
		vals := make([]int, n)
		for i := range vals {
			vals[i] = rng.IntN(4)
		}
		want := map[string]bool{}
		for s := range Subsets(vals) {
			c := slices.Clone(s)
			slices.Sort(c)
			want[fmt.Sprint(c)] = true
		}
		got := map[string]bool{}
		for s := range SubsetsUnique(vals) {
			if !slices.IsSorted(s) {
				t.Fatalf("trial %d: %v not ascending", trial, s)
			}
			key := fmt.Sprint(s)
			if got[key] {
				t.Fatalf("trial %d: %v repeated", trial, s)
			}
			got[key] = true
		}
		if len(got) != len(want) {
			t.Fatalf("trial %d: %d sub-multisets, want %d", trial, len(got), len(want))
		}
		for k := range want {
			if !got[k] {
				t.Fatalf("trial %d: missing %s", trial, k)
			}
		}
	}
}

func BenchmarkSubsetsUnique(b *testing.B) {
	for _, n := range []int{10, 15, 20} {
		vals := make([]int, n)
		for i := range vals {
			vals[i] = i / 2
		}
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			sink := 0
			for b.Loop() {
				for s := range SubsetsUnique(vals) {
					sink += len(s)
				}
			}
			benchSink = sink
		})
	}
}

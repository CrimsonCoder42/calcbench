package enumerate

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"testing"
)

// Knuth, TAOCP Vol. 4A, §7.2.1.1, Algorithm M with radices 2: binary counting 000, 001, …, 111 read with
// values[0] as the least significant digit gives {} a b ab c ac bc abc.
func TestSubsets_knuthExample(t *testing.T) {
	got := collect(Subsets([]string{"a", "b", "c"}))
	want := [][]string{{}, {"a"}, {"b"}, {"a", "b"}, {"c"}, {"a", "c"}, {"b", "c"}, {"a", "b", "c"}}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if !slices.Equal(got[i], want[i]) {
			t.Errorf("[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}

func TestSubsets_cases(t *testing.T) {
	for n, want := range map[int]int{0: 1, 1: 2, 2: 4, 5: 32} {
		if got := len(collect(Subsets(make([]int, n)))); got != want {
			t.Errorf("Subsets of %d values yields %d, want %d", n, got, want)
		}
	}
	// Duplicates are not merged.
	if got := len(collect(Subsets([]int{3, 3}))); got != 4 {
		t.Errorf("duplicates: %d, want 4", got)
	}
}

func TestSubsets_earlyBreakLargeN(t *testing.T) {
	count := 0
	for range Subsets(make([]int, 100)) { // 2^100: beyond any mask type
		count++
		if count == 10 {
			break
		}
	}
	if count != 10 {
		t.Errorf("count = %d, want 10", count)
	}
}

// TestSubsets_oracle: the k-th subset is the set bits of k.
func TestSubsets_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 200 {
		n := rng.IntN(11)
		vals := rng.Perm(50)[:n]
		k := 0
		for s := range Subsets(vals) {
			var want []int
			for i := range n {
				if k>>i&1 == 1 {
					want = append(want, vals[i])
				}
			}
			if !slices.Equal(s, want) && !(len(s) == 0 && len(want) == 0) {
				t.Fatalf("trial %d: subset %d = %v, want %v", trial, k, s, want)
			}
			k++
		}
		if k != 1<<n {
			t.Fatalf("trial %d: %d subsets, want %d", trial, k, 1<<n)
		}
	}
}

func BenchmarkSubsets(b *testing.B) {
	for _, n := range []int{10, 15, 20} {
		vals := make([]int, n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			sink := 0
			for b.Loop() {
				for s := range Subsets(vals) {
					sink += len(s)
				}
			}
			benchSink = sink
		})
	}
}

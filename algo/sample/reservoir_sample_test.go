package sample

import (
	"errors"
	"fmt"
	"iter"
	"math/rand/v2"
	"slices"
	"testing"
)

// Vitter 1985, Algorithm R: every k-subset is equally likely. Checked for 2 of 4 items: 120 000 samples,
// each of the C(4,2) = 6 pairs within 5σ of 20 000.
func TestReservoirSample_vitterUniformPairs(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	counts := map[string]int{}
	const trials = 120_000
	var buf []int
	for range trials {
		var err error
		buf, err = ReservoirSample(buf[:0], rng, slices.Values([]int{1, 2, 3, 4}), 2)
		if err != nil {
			t.Fatal(err)
		}
		pair := slices.Sorted(slices.Values(buf))
		counts[fmt.Sprint(pair)]++
	}
	if len(counts) != 6 {
		t.Fatalf("saw %d pairs, want 6", len(counts))
	}
	for k, c := range counts {
		if !within5Sigma(c, trials, 1.0/6) {
			t.Errorf("pair %s: %d of %d", k, c, trials)
		}
	}
}

func TestReservoirSample_cases(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	tests := []struct {
		name string
		in   []int
		k    int
		want int
	}{
		{"empty sequence", nil, 3, 0},
		{"k zero", []int{1, 2, 3}, 0, 0},
		{"fewer than k", []int{1, 2}, 5, 2},
		{"exactly k", []int{1, 2, 3}, 3, 3},
		{"one of many", []int{1, 2, 3, 4, 5, 6}, 1, 1},
		{"huge k does not preallocate", []int{1, 2}, 1 << 62, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReservoirSample(nil, rng, slices.Values(tt.in), tt.k)
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != tt.want {
				t.Errorf("len = %d, want %d", len(got), tt.want)
			}
			for _, v := range got {
				if !slices.Contains(tt.in, v) {
					t.Errorf("%d is not an input item", v)
				}
			}
		})
	}
}

func TestReservoirSample_appendsToDst(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	got, err := ReservoirSample([]int{-7}, rng, slices.Values([]int{1, 2, 3}), 2)
	if err != nil || len(got) != 3 || got[0] != -7 {
		t.Errorf("got %v, %v; want -7 then two items", got, err)
	}
}

func TestReservoirSample_kZeroDoesNotRead(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	read := false
	seq := iter.Seq[int](func(yield func(int) bool) { read = true })
	if _, err := ReservoirSample(nil, rng, seq, 0); err != nil || read {
		t.Errorf("err = %v, read = %v; want nil, false", err, read)
	}
}

func TestReservoirSample_errors(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	dst := []int{9}
	tests := []struct {
		name string
		rng  *rand.Rand
		seq  iter.Seq[int]
		k    int
		want error
	}{
		{"nil rng", nil, slices.Values([]int{1}), 1, ErrNilRand},
		{"nil seq", rng, nil, 1, ErrNilSeq},
		{"negative k", rng, slices.Values([]int{1}), -1, ErrNegative},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReservoirSample(dst, tt.rng, tt.seq, tt.k)
			if !errors.Is(err, tt.want) || !slices.Equal(got, dst) {
				t.Errorf("got %v, %v; want dst unchanged, %v", got, err, tt.want)
			}
		})
	}
}

// TestReservoirSample_oracle: each of n items is included with probability k/n (n = 10, k = 3, 50 000
// runs), and a sample never repeats a position.
func TestReservoirSample_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	const n, k, trials = 10, 3, 50_000
	items := make([]int, n)
	for i := range items {
		items[i] = i
	}
	var counts [n]int
	var buf []int
	for range trials {
		buf, _ = ReservoirSample(buf[:0], rng, slices.Values(items), k)
		seen := map[int]bool{}
		for _, v := range buf {
			if seen[v] {
				t.Fatalf("item %d sampled twice: %v", v, buf)
			}
			seen[v] = true
			counts[v]++
		}
	}
	for i, c := range counts {
		if !within5Sigma(c, trials, float64(k)/n) {
			t.Errorf("item %d: %d of %d", i, c, trials)
		}
	}
}

func BenchmarkReservoirSample(b *testing.B) {
	for _, size := range benchSizes {
		items := make([]int, size.n)
		rng := rand.New(rand.NewPCG(1, 2))
		dst := make([]int, 0, 16)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				dst, _ = ReservoirSample(dst[:0], rng, slices.Values(items), 16)
			}
			benchSink = len(dst)
		})
	}
}

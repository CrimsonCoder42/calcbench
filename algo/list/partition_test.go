package list

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// partitionModel is the oracle: a stable split of a slice copy.
func partitionModel(in []int, pivot int) []int {
	var less, rest []int
	for _, v := range in {
		if v < pivot {
			less = append(less, v)
		} else {
			rest = append(rest, v)
		}
	}
	return append(less, rest...)
}

// Worked example (hand-derived, catalog/list.md ALG-LIST-014); also the problem-site pattern, in our words.
func TestPartitionInPlace_workedExample(t *testing.T) {
	head, err := PartitionInPlace(FromSlice([]int{1, 4, 3, 2, 5, 2}), 3)
	if got := values(t, head); err != nil || !slices.Equal(got, []int{1, 2, 2, 4, 3, 5}) {
		t.Errorf("got (%v, %v)", got, err)
	}
}

func TestPartitionInPlace_cases(t *testing.T) {
	tests := []struct {
		name  string
		in    []int
		pivot int
	}{
		{"empty", nil, 0},
		{"one less", []int{1}, 2},
		{"one not less", []int{2}, 2},
		{"all less", []int{3, 1, 2}, 9},
		{"none less", []int{3, 1, 2}, -9},
		{"duplicates of pivot", []int{2, 1, 2, 1}, 2},
		{"extremes", []int{math.MaxInt, math.MinInt, 0}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			head, err := PartitionInPlace(FromSlice(tt.in), tt.pivot)
			if got, want := values(t, head), partitionModel(tt.in, tt.pivot); err != nil || !slices.Equal(got, want) {
				t.Errorf("got (%v, %v), want %v", got, err, want)
			}
		})
	}
}

func TestPartitionInPlace_nanPivot(t *testing.T) {
	in := []float64{2, math.Inf(-1), 1}
	head, err := PartitionInPlace(FromSlice(in), math.NaN())
	if got := values(t, head); err != nil || !slices.Equal(got, in) {
		t.Errorf("NaN pivot: (%v, %v), want unchanged %v", got, err, in)
	}
}

func TestPartitionInPlace_cycle(t *testing.T) {
	cyc, _ := cyclic([]int{5, 1, 4}, 0)
	if got, err := PartitionInPlace(cyc, 3); !errors.Is(err, ErrCycle) || got != cyc || !DetectCycle(cyc) {
		t.Errorf("err = %v", err)
	}
}

func TestPartitionInPlace_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		in := randInts(rng, rng.IntN(40), 10)
		pivot := rng.IntN(12) - 1
		head, err := PartitionInPlace(FromSlice(in), pivot)
		if got, want := values(t, head), partitionModel(in, pivot); err != nil || !slices.Equal(got, want) {
			t.Fatalf("(%v, %d) = %v, want %v", in, pivot, got, want)
		}
	}
}

func TestPartitionInPlace_noAlloc(t *testing.T) {
	head := FromSlice(countUpTo(100))
	if a := testing.AllocsPerRun(100, func() { head, _ = PartitionInPlace(head, 50) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkPartitionInPlace(b *testing.B) {
	for _, s := range benchSizes {
		head := FromSlice(randInts(rand.New(rand.NewPCG(1, 2)), s.n, 1000))
		b.Run(s.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				head, _ = PartitionInPlace(head, 500)
			}
		})
	}
}

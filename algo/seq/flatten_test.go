package seq

import (
	"iter"
	"math/rand/v2"
	"slices"
	"testing"
)

// seqsOf turns slices into a sequence of sequences; a nil slice becomes a nil inner iterator when
// nilAsNil is set.
func seqsOf(groups [][]int, nilAsNil bool) iter.Seq[iter.Seq[int]] {
	return func(yield func(iter.Seq[int]) bool) {
		for _, g := range groups {
			var s iter.Seq[int]
			if g != nil || !nilAsNil {
				s = slices.Values(g)
			}
			if !yield(s) {
				return
			}
		}
	}
}

// Hand-derived: [1 2] [] [3] flattens to 1 2 3.
func TestFlatten_example(t *testing.T) {
	got := slices.Collect(Flatten(seqsOf([][]int{{1, 2}, {}, {3}}, false)))
	if !slices.Equal(got, []int{1, 2, 3}) {
		t.Errorf("Flatten = %v, want [1 2 3]", got)
	}
}

func TestFlatten_cases(t *testing.T) {
	tests := []struct {
		name   string
		groups [][]int
		want   []int
	}{
		{"no inner", nil, nil},
		{"all inner empty", [][]int{{}, {}}, nil},
		{"one", [][]int{{4}}, []int{4}},
		{"nil inner skipped", [][]int{nil, {1}, nil, {2, 2}}, []int{1, 2, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slices.Collect(Flatten(seqsOf(tt.groups, true)))
			if !slices.Equal(got, tt.want) {
				t.Errorf("Flatten(%v) = %v, want %v", tt.groups, got, tt.want)
			}
		})
	}
	if got := slices.Collect(Flatten[int](nil)); len(got) != 0 {
		t.Errorf("Flatten(nil) = %v, want empty", got)
	}
}

func TestFlatten_stopsInnerAndOuterOnBreak(t *testing.T) {
	var inner, later probe
	var outerStopped bool
	outerProduced := 0
	outer := func(yield func(iter.Seq[int]) bool) {
		for _, s := range []iter.Seq[int]{inner.ints(-1), later.ints(-1)} {
			outerProduced++
			if !yield(s) {
				outerStopped = true
				return
			}
		}
	}
	for v := range Flatten(outer) {
		if v == 2 {
			break
		}
	}
	if !inner.stopped || inner.produced != 3 || !outerStopped || outerProduced != 1 || later.produced != 0 {
		t.Errorf("inner produced %d stopped %v, outer produced %d stopped %v, later produced %d",
			inner.produced, inner.stopped, outerProduced, outerStopped, later.produced)
	}
}

func TestFlatten_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		groups := make([][]int, rng.IntN(6))
		var want []int
		for i := range groups {
			groups[i] = randomInts(rng, rng.IntN(5))
			want = append(want, groups[i]...)
		}
		if got := slices.Collect(Flatten(seqsOf(groups, false))); !slices.Equal(got, want) {
			t.Fatalf("Flatten(%v) = %v, want %v", groups, got, want)
		}
	}
}

func TestFlatten_noAllocationPerElement(t *testing.T) {
	run := func(n int) float64 {
		s := Flatten(slices.Values([]iter.Seq[int]{count(n), count(n)}))
		return testing.AllocsPerRun(20, func() {
			for range s {
			}
		})
	}
	if small, large := run(10), run(10_000); small != large {
		t.Errorf("Flatten allocates %v times over 20 elements and %v over 20000; want equal", small, large)
	}
}

func BenchmarkFlatten(b *testing.B) {
	for _, size := range benchSizes {
		b.Run(size.name, func(b *testing.B) {
			inner := make([]iter.Seq[int], 10)
			for i := range inner {
				inner[i] = count(size.n / 10)
			}
			s := Flatten(slices.Values(inner))
			b.ReportAllocs()
			for b.Loop() {
				sum := 0
				for v := range s {
					sum += v
				}
				_ = sum
			}
		})
	}
}

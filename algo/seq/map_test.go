package seq

import (
	"iter"
	"math/rand/v2"
	"slices"
	"strconv"
	"testing"
)

// probe is an upstream iterator that records how it was consumed. Shared by every test file in this
// package. The runtime itself fails a range-over-func loop whose iterator calls yield after it returned
// false, so a probe that keeps going after a stop would make the test panic.
type probe struct {
	produced int  // elements handed to yield
	stopped  bool // yield returned false
	finished bool // ran to its natural end
}

// ints yields 0, 1, …, n-1, or forever when n < 0.
func (p *probe) ints(n int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := 0; n < 0 || i < n; i++ {
			p.produced++
			if !yield(i) {
				p.stopped = true
				return
			}
		}
		p.finished = true
	}
}

// benchSizes are the collection sizes of .claude/rules/go-complexity.md. Shared by the benchmarks.
var benchSizes = []struct {
	name string
	n    int
}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}}

// count yields 0 … n-1 with no bookkeeping, for benchmarks and allocation tests.
func count(n int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := range n {
			if !yield(i) {
				return
			}
		}
	}
}

// randomInts returns n fixed-seed values in [0, 10), shared by the oracle tests.
func randomInts(rng *rand.Rand, n int) []int {
	s := make([]int, n)
	for i := range s {
		s[i] = rng.IntN(10)
	}
	return s
}

// Hand-derived: squaring 1 2 3 gives 1 4 9 (the spec's range-over-func rule, applied element by element).
func TestMap_example(t *testing.T) {
	got := slices.Collect(Map(slices.Values([]int{1, 2, 3}), func(v int) int { return v * v }))
	if !slices.Equal(got, []int{1, 4, 9}) {
		t.Errorf("Map(1 2 3, square) = %v, want [1 4 9]", got)
	}
}

func TestMap_changesType(t *testing.T) {
	got := slices.Collect(Map(slices.Values([]int{-1, 0, 7}), strconv.Itoa))
	if !slices.Equal(got, []string{"-1", "0", "7"}) {
		t.Errorf("Map(-1 0 7, Itoa) = %q", got)
	}
}

func TestMap_empty(t *testing.T) {
	double := func(v int) int { return 2 * v }
	if got := slices.Collect(Map(slices.Values([]int(nil)), double)); len(got) != 0 {
		t.Errorf("Map(empty) = %v, want empty", got)
	}
	if got := slices.Collect(Map(nil, double)); len(got) != 0 {
		t.Errorf("Map(nil seq) = %v, want empty", got)
	}
	var p probe
	if got := slices.Collect(Map[int, int](p.ints(3), nil)); len(got) != 0 || p.produced != 0 {
		t.Errorf("Map(nil fn) = %v, produced %d, want empty and 0", got, p.produced)
	}
}

func TestMap_stopsUpstreamOnBreak(t *testing.T) {
	var p probe
	calls := 0
	for v := range Map(p.ints(-1), func(v int) int { calls++; return v + 1 }) {
		if v == 3 {
			break
		}
	}
	if p.produced != 3 || !p.stopped || calls != 3 {
		t.Errorf("after break at the 3rd element: produced %d, stopped %v, fn calls %d; want 3, true, 3", p.produced, p.stopped, calls)
	}
}

func TestMap_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		in := randomInts(rng, rng.IntN(50))
		want := make([]int, len(in))
		for i, v := range in {
			want[i] = 3*v - 1
		}
		got := slices.Collect(Map(slices.Values(in), func(v int) int { return 3*v - 1 }))
		if !slices.Equal(got, want) {
			t.Fatalf("Map(%v) = %v, want %v", in, got, want)
		}
	}
}

// The allocations of one range loop do not grow with the number of elements.
func TestMap_noAllocationPerElement(t *testing.T) {
	run := func(n int) float64 {
		s := Map(count(n), func(v int) int { return v + 1 })
		return testing.AllocsPerRun(20, func() {
			for range s {
			}
		})
	}
	if small, large := run(10), run(10_000); small != large {
		t.Errorf("Map allocates %v times over 10 elements and %v over 10000; want equal", small, large)
	}
}

func BenchmarkMap(b *testing.B) {
	for _, size := range benchSizes {
		b.Run(size.name, func(b *testing.B) {
			s := Map(count(size.n), func(v int) int { return v + 1 })
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

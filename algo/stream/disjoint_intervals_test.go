package stream

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/CrimsonCoder42/calcbench/algo/interval"
)

type iv = interval.Interval[int]

// r is the closed range [lo, hi].
func r(lo, hi int) iv { return iv{Start: lo, End: hi} }

// Hand-derived: adding 1, 3, 7, 2, 6 gives [1, 3] (2 fills the gap) and [6, 7].
func TestDisjointIntervals_example(t *testing.T) {
	var d DisjointIntervals[int]
	steps := []struct {
		v    int
		want []iv
	}{
		{1, []iv{r(1, 1)}},
		{3, []iv{r(1, 1), r(3, 3)}},
		{7, []iv{r(1, 1), r(3, 3), r(7, 7)}},
		{2, []iv{r(1, 3), r(7, 7)}},
		{6, []iv{r(1, 3), r(6, 7)}},
	}
	for _, s := range steps {
		d.Add(s.v)
		if got := d.Intervals(nil); !slices.Equal(got, s.want) {
			t.Errorf("after Add(%d): %v, want %v", s.v, got, s.want)
		}
	}
}

func TestDisjointIntervals_cases(t *testing.T) {
	tests := []struct {
		name string
		adds []int
		want []iv
	}{
		{"empty", nil, nil},
		{"one", []int{5}, []iv{r(5, 5)}},
		{"duplicates", []int{5, 5, 5}, []iv{r(5, 5)}},
		{"inside a range", []int{1, 2, 3, 2}, []iv{r(1, 3)}},
		{"ascending", []int{1, 2, 3, 4}, []iv{r(1, 4)}},
		{"descending", []int{4, 3, 2, 1}, []iv{r(1, 4)}},
		{"joins left", []int{1, 2}, []iv{r(1, 2)}},
		{"joins right", []int{2, 1}, []iv{r(1, 2)}},
		{"joins both", []int{1, 3, 2}, []iv{r(1, 3)}},
		{"negative and zero", []int{-1, 1, 0, -5}, []iv{r(-5, -5), r(-1, 1)}},
		{"extremes", []int{math.MaxInt, math.MinInt, math.MaxInt - 1, math.MinInt + 1}, []iv{r(math.MinInt, math.MinInt+1), r(math.MaxInt-1, math.MaxInt)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d DisjointIntervals[int]
			for _, v := range tt.adds {
				d.Add(v)
			}
			if got := d.Intervals(nil); !slices.Equal(got, tt.want) || d.Len() != len(tt.want) {
				t.Errorf("Add %v: %v (Len %d), want %v", tt.adds, got, d.Len(), tt.want)
			}
		})
	}
}

func TestDisjointIntervals_smallTypes(t *testing.T) {
	var s DisjointIntervals[int8]
	for _, v := range []int8{127, -128, 126, -127, 0} {
		s.Add(v)
	}
	want := []interval.Interval[int8]{{Start: -128, End: -127}, {Start: 0, End: 0}, {Start: 126, End: 127}}
	if got := s.Intervals(nil); !slices.Equal(got, want) {
		t.Errorf("int8: %v, want %v", got, want)
	}
	var u DisjointIntervals[uint8]
	for _, v := range []uint8{255, 0, 1, 254} {
		u.Add(v)
	}
	wantU := []interval.Interval[uint8]{{Start: 0, End: 1}, {Start: 254, End: 255}}
	if got := u.Intervals(nil); !slices.Equal(got, wantU) {
		t.Errorf("uint8: %v, want %v", got, wantU)
	}
}

func TestDisjointIntervals_contains(t *testing.T) {
	var d DisjointIntervals[int]
	for _, v := range []int{1, 2, 3, 7} {
		d.Add(v)
	}
	for v, want := range map[int]bool{0: false, 1: true, 2: true, 3: true, 4: false, 6: false, 7: true, 8: false} {
		if got := d.Contains(v); got != want {
			t.Errorf("Contains(%d) = %v, want %v", v, got, want)
		}
	}
	var empty DisjointIntervals[int]
	if empty.Contains(0) {
		t.Error("the empty set contains 0")
	}
}

func TestDisjointIntervals_appendsToDst(t *testing.T) {
	var d DisjointIntervals[int]
	d.Add(4)
	got := d.Intervals([]iv{r(0, 0)})
	if !slices.Equal(got, []iv{r(0, 0), r(4, 4)}) {
		t.Errorf("Intervals([[0 0]]) = %v", got)
	}
}

// slowRanges is the oracle: a presence table over [0, span), read as runs.
func slowRanges(present []bool) []iv {
	var out []iv
	for i := 0; i < len(present); i++ {
		if !present[i] {
			continue
		}
		j := i
		for j+1 < len(present) && present[j+1] {
			j++
		}
		out = append(out, r(i, j))
		i = j
	}
	return out
}

func TestDisjointIntervals_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		const span = 60
		present := make([]bool, span)
		var d DisjointIntervals[int]
		for range rng.IntN(80) {
			v := rng.IntN(span)
			present[v] = true
			d.Add(v)
			if !d.Contains(v) {
				t.Fatalf("Contains(%d) false right after Add", v)
			}
		}
		if got, want := d.Intervals(nil), slowRanges(present); !slices.Equal(got, want) {
			t.Fatalf("Intervals = %v, want %v", got, want)
		}
	}
}

func TestDisjointIntervals_noAllocation(t *testing.T) {
	var d DisjointIntervals[int]
	for v := 0; v < 100; v += 2 {
		d.Add(v)
	}
	buf := make([]iv, 0, 64)
	allocs := testing.AllocsPerRun(20, func() {
		d.Add(10) // already covered
		_ = d.Contains(11)
		buf = d.Intervals(buf[:0])
	})
	if allocs != 0 {
		t.Errorf("Add (covered), Contains and Intervals into a buffer allocate %v times, want 0", allocs)
	}
}

// Adds n fixed-seed values from [0, 2n): about n/4 ranges at the end, so Add's memmove is O(k).
func BenchmarkDisjointIntervals(b *testing.B) {
	for _, size := range benchSizes {
		b.Run(size.name, func(b *testing.B) {
			rng := rand.New(rand.NewPCG(1, 2))
			values := make([]int, size.n)
			for i := range values {
				values[i] = rng.IntN(2 * size.n)
			}
			b.ReportAllocs()
			for b.Loop() {
				var d DisjointIntervals[int]
				for _, v := range values {
					d.Add(v)
				}
			}
		})
	}
}

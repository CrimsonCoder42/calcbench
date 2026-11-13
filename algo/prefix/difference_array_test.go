package prefix

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Laaksonen, Competitive Programmer's Handbook (2018), §9.4 "Range updates", p. 93: the array
// 3 3 1 1 1 5 2 2 has the difference array 3 0 −2 0 0 4 −3 0; increasing positions 1…4 (closed) by 5
// changes only two differences. The materialized result 3 8 6 6 6 5 2 2 is hand-derived.
func TestDifferenceArray_laaksonenExample(t *testing.T) {
	initial := []int{3, 3, 1, 1, 1, 5, 2, 2}
	d := NewDifferenceArray(initial)
	wantDiff := []int{3, 0, -2, 0, 0, 4, -3, 0}
	for i, w := range wantDiff {
		if d.diff[i] != widen(w) {
			t.Fatalf("diff[%d] = %v, want %d", i, d.diff[i], w)
		}
	}
	if err := d.RangeAdd(1, 5, 5); err != nil {
		t.Fatal(err)
	}
	got, err := d.Materialize(nil)
	if want := []int{3, 8, 6, 6, 6, 5, 2, 2}; err != nil || !slices.Equal(got, want) {
		t.Errorf("after RangeAdd(1, 5, 5): Materialize = (%v, %v), want (%v, nil)", got, err, want)
	}
	if !slices.Equal(initial, []int{3, 3, 1, 1, 1, 5, 2, 2}) {
		t.Errorf("NewDifferenceArray modified its input: %v", initial)
	}
}

func TestDifferenceArray_cases(t *testing.T) {
	d := NewDifferenceArray(make([]int, 4))
	steps := []struct {
		lo, hi int
		delta  int
		want   []int
	}{
		{0, 0, 9, []int{0, 0, 0, 0}},      // empty range: no-op
		{4, 4, 9, []int{0, 0, 0, 0}},      // empty range at the end
		{0, 4, -1, []int{-1, -1, -1, -1}}, // whole array, negative
		{3, 4, 5, []int{-1, -1, -1, 4}},   // last element only
		{1, 3, 0, []int{-1, -1, -1, 4}},   // zero delta
		{0, 2, 2, []int{1, 1, -1, 4}},     // prefix
		{1, 2, 3, []int{1, 4, -1, 4}},     // one element
	}
	for i, s := range steps {
		if err := d.RangeAdd(s.lo, s.hi, s.delta); err != nil {
			t.Fatalf("step %d: RangeAdd error %v", i, err)
		}
		got, err := d.Materialize(nil)
		if err != nil || !slices.Equal(got, s.want) {
			t.Fatalf("step %d: Materialize = (%v, %v), want %v", i, got, err, s.want)
		}
	}
	if d.Len() != 4 {
		t.Errorf("Len = %d, want 4", d.Len())
	}
}

func TestDifferenceArray_zeroAndNil(t *testing.T) {
	var zero DifferenceArray[int]
	if zero.Len() != 0 {
		t.Errorf("zero Len = %d", zero.Len())
	}
	if err := zero.RangeAdd(0, 0, 1); err != nil {
		t.Errorf("zero RangeAdd(0, 0) error %v", err)
	}
	if got, err := zero.Materialize(nil); err != nil || len(got) != 0 {
		t.Errorf("zero Materialize = (%v, %v)", got, err)
	}
	var nilArray *DifferenceArray[int]
	if nilArray.Len() != 0 {
		t.Errorf("nil Len = %d", nilArray.Len())
	}
	if err := nilArray.RangeAdd(0, 1, 1); !errors.Is(err, ErrRange) {
		t.Errorf("nil RangeAdd(0, 1) error = %v, want ErrRange", err)
	}
	empty := NewDifferenceArray[int](nil)
	if got, err := empty.Materialize(nil); err != nil || len(got) != 0 {
		t.Errorf("empty Materialize = (%v, %v)", got, err)
	}
}

func TestDifferenceArray_errors(t *testing.T) {
	d := NewDifferenceArray([]int{1, 2, 3})
	for _, r := range [][2]int{{-1, 1}, {2, 1}, {0, 4}, {math.MinInt, 0}} {
		if err := d.RangeAdd(r[0], r[1], 1); !errors.Is(err, ErrRange) {
			t.Errorf("RangeAdd(%d, %d) error = %v, want ErrRange", r[0], r[1], err)
		}
	}
	if got, _ := d.Materialize(nil); !slices.Equal(got, []int{1, 2, 3}) {
		t.Errorf("a failed RangeAdd changed the array: %v", got)
	}
	// Above the int8 range, then back: only the final values matter.
	d8 := NewDifferenceArray([]int8{127, -128})
	_ = d8.RangeAdd(0, 1, 1)
	dst := []int8{7}
	if got, err := d8.Materialize(dst); !errors.Is(err, ErrOverflow) || !slices.Equal(got, []int8{7}) {
		t.Errorf("Materialize(127+1) = (%v, %v), want ([7], ErrOverflow)", got, err)
	}
	_ = d8.RangeAdd(0, 2, -1)
	if got, err := d8.Materialize(nil); !errors.Is(err, ErrOverflow) {
		t.Errorf("Materialize(-128-1) = (%v, %v), want ErrOverflow", got, err)
	}
	_ = d8.RangeAdd(1, 2, 1)
	if got, err := d8.Materialize(nil); err != nil || !slices.Equal(got, []int8{127, -128}) {
		t.Errorf("Materialize after returning in range = (%v, %v), want ([127 -128], nil)", got, err)
	}
	// Many adds of the extreme int64 values: the 128-bit differences stay exact.
	d64 := NewDifferenceArray([]int64{0, 0})
	for range 4 {
		_ = d64.RangeAdd(0, 1, math.MaxInt64)
	}
	for range 4 {
		_ = d64.RangeAdd(0, 1, -math.MaxInt64)
	}
	if got, err := d64.Materialize(nil); err != nil || !slices.Equal(got, []int64{0, 0}) {
		t.Errorf("int64 extremes = (%v, %v), want ([0 0], nil)", got, err)
	}
}

func TestDifferenceArray_noAllocation(t *testing.T) {
	d := NewDifferenceArray(make([]int64, 1000))
	buf := make([]int64, 0, 1000)
	allocs := testing.AllocsPerRun(20, func() {
		_ = d.RangeAdd(10, 500, 3)
		buf, _ = d.Materialize(buf[:0])
	})
	if allocs != 0 {
		t.Errorf("RangeAdd + Materialize with a reused buffer allocate %v times, want 0", allocs)
	}
}

// Oracle: random range adds against a plain array updated element by element.
func TestDifferenceArray_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		n := rng.IntN(20)
		naive := make([]int64, n)
		for i := range naive {
			naive[i] = int64(rng.IntN(200) - 100)
		}
		d := NewDifferenceArray(naive)
		for range rng.IntN(30) {
			lo := rng.IntN(n + 1)
			hi := lo + rng.IntN(n-lo+1)
			delta := int64(rng.IntN(2001) - 1000)
			if err := d.RangeAdd(lo, hi, delta); err != nil {
				t.Fatal(err)
			}
			for i := lo; i < hi; i++ {
				naive[i] += delta
			}
		}
		got, err := d.Materialize(nil)
		if err != nil || !slices.Equal(got, naive) {
			t.Fatalf("trial %d: Materialize = (%v, %v), want %v", trial, got, err, naive)
		}
	}
}

func BenchmarkNewDifferenceArray(b *testing.B) {
	for _, size := range benchSizes {
		initial := benchInts(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchSink = NewDifferenceArray(initial).Len()
			}
		})
	}
}

func BenchmarkRangeAdd(b *testing.B) {
	for _, size := range benchSizes {
		d := NewDifferenceArray(benchInts(size.n))
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				for i := 0; i < size.n; i++ {
					_ = d.RangeAdd(i/2, i, 1)
				}
			}
		})
	}
}

func BenchmarkMaterialize(b *testing.B) {
	for _, size := range benchSizes {
		d := NewDifferenceArray(benchInts(size.n))
		buf := make([]int64, 0, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				buf, _ = d.Materialize(buf[:0])
			}
		})
	}
}

func BenchmarkDifferenceArrayLen(b *testing.B) {
	d := NewDifferenceArray(make([]int64, 1000))
	b.ReportAllocs()
	for b.Loop() {
		benchSink = d.Len()
	}
}

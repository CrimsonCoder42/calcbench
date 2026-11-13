package window

import (
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// dequeContents returns the elements held, front to back.
func dequeContents[T float64 | int](d *MonotonicDeque[T]) []T {
	out := make([]T, 0, d.size)
	for i := range d.size {
		out = append(out, d.ring[(d.head+i)&(len(d.ring)-1)])
	}
	return out
}

// slidingExtremes drives a deque over fixed windows of length w the way callers do, returning the front
// of each full window.
func slidingExtremes(values []int, w int, wantMax bool) []int {
	d := NewMonotonicDeque[int](wantMax)
	var out []int
	for i, v := range values {
		if i >= w {
			d.PopFrontIfEqual(values[i-w])
		}
		d.PushBack(v)
		if i >= w-1 {
			f, _ := d.Front()
			out = append(out, f)
		}
	}
	return out
}

// Laaksonen, Competitive Programmer's Handbook (2018), §8.3 "Sliding window minimum", pp. 81–82: the
// array 2 1 4 5 3 4 1 2 with window 4. The handbook shows the queue after each move as 1 4 5, 1 3,
// 3 4, 1, 1 2, and the minima 1, 1, 3, 1, 1.
func TestMonotonicDeque_laaksonenExample(t *testing.T) {
	values := []int{2, 1, 4, 5, 3, 4, 1, 2}
	const w = 4
	wantQueues := [][]int{{1, 4, 5}, {1, 3}, {3, 4}, {1}, {1, 2}}
	var d MonotonicDeque[int] // the zero value is a min-deque
	var minima []int
	for i, v := range values {
		if i >= w {
			d.PopFrontIfEqual(values[i-w])
		}
		d.PushBack(v)
		if i >= w-1 {
			if got := dequeContents(&d); !slices.Equal(got, wantQueues[i-w+1]) {
				t.Errorf("window %d: queue %v, want %v", i-w+1, got, wantQueues[i-w+1])
			}
			f, _ := d.Front()
			minima = append(minima, f)
		}
	}
	if want := []int{1, 1, 3, 1, 1}; !slices.Equal(minima, want) {
		t.Errorf("minima %v, want %v", minima, want)
	}
}

func TestMonotonicDeque_operations(t *testing.T) {
	d := NewMonotonicDeque[int](false)
	if _, ok := d.Front(); ok || d.Len() != 0 {
		t.Fatal("new deque not empty")
	}
	if _, ok := d.PopFront(); ok {
		t.Error("PopFront on empty: ok = true")
	}
	if d.PopFrontIfEqual(0) {
		t.Error("PopFrontIfEqual on empty: true")
	}
	steps := []struct {
		push    int
		evicted int
		want    []int
	}{
		{5, 0, []int{5}},
		{7, 0, []int{5, 7}},
		{7, 0, []int{5, 7, 7}}, // equal elements are kept
		{6, 2, []int{5, 6}},
		{math.MinInt, 2, []int{math.MinInt}},
		{math.MaxInt, 0, []int{math.MinInt, math.MaxInt}},
	}
	for i, s := range steps {
		if ev := d.PushBack(s.push); ev != s.evicted {
			t.Errorf("step %d: PushBack(%d) evicted %d, want %d", i, s.push, ev, s.evicted)
		}
		if got := dequeContents(d); !slices.Equal(got, s.want) {
			t.Fatalf("step %d: contents %v, want %v", i, got, s.want)
		}
	}
	if d.PopFrontIfEqual(0) {
		t.Error("PopFrontIfEqual(0) with front MinInt: true")
	}
	if v, ok := d.PopFront(); !ok || v != math.MinInt || d.Len() != 1 {
		t.Errorf("PopFront = (%d, %v), Len %d; want (MinInt, true), 1", v, ok, d.Len())
	}
	if !d.PopFrontIfEqual(math.MaxInt) || d.Len() != 0 {
		t.Errorf("PopFrontIfEqual(MaxInt) failed, Len %d", d.Len())
	}
}

func TestMonotonicDeque_max(t *testing.T) {
	d := NewMonotonicDeque[int](true)
	for _, v := range []int{3, 1, 2, 2, 5, 4} {
		d.PushBack(v)
	}
	if got, want := dequeContents(d), []int{5, 4}; !slices.Equal(got, want) {
		t.Errorf("max-deque contents %v, want %v", got, want)
	}
	if f, ok := d.Front(); !ok || f != 5 {
		t.Errorf("max-deque Front = (%d, %v), want (5, true)", f, ok)
	}
}

func TestMonotonicDeque_nilAndNaN(t *testing.T) {
	var nilDeque *MonotonicDeque[float64]
	if _, ok := nilDeque.Front(); ok || nilDeque.Len() != 0 || nilDeque.PopFrontIfEqual(1) {
		t.Error("nil deque is not empty")
	}
	if _, ok := nilDeque.PopFront(); ok {
		t.Error("nil PopFront ok = true")
	}
	// cmp.Compare order: NaN is the least value, and equal to itself.
	var d MonotonicDeque[float64]
	d.PushBack(2)
	d.PushBack(math.NaN())
	d.PushBack(1)
	if f, _ := d.Front(); !math.IsNaN(f) || d.Len() != 2 {
		t.Errorf("min-deque after 2 NaN 1: front %v, len %d; want NaN, 2", f, d.Len())
	}
	if !d.PopFrontIfEqual(math.NaN()) {
		t.Error("PopFrontIfEqual(NaN) with a NaN front: false")
	}
}

// Growth across the ring boundary: many pushes without eviction, interleaved pops.
func TestMonotonicDeque_wrapAndGrow(t *testing.T) {
	var d MonotonicDeque[int]
	next, popped := 0, 0
	for round := range 50 {
		for range round % 7 {
			d.PushBack(next)
			next++
		}
		for range round % 3 {
			if v, ok := d.PopFront(); ok {
				if v != popped {
					t.Fatalf("PopFront = %d, want %d", v, popped)
				}
				popped++
			}
		}
		if d.Len() != next-popped {
			t.Fatalf("Len = %d, want %d", d.Len(), next-popped)
		}
	}
}

func TestMonotonicDeque_noAllocationOnceGrown(t *testing.T) {
	values := make([]int, 1000)
	for i := range values {
		values[i] = (i * 7919) % 1000
	}
	d := NewMonotonicDeque[int](false)
	for _, v := range values {
		d.PushBack(v)
	}
	allocs := testing.AllocsPerRun(20, func() {
		for d.Len() > 0 {
			d.PopFront()
		}
		for i, v := range values {
			if i >= benchWindow {
				d.PopFrontIfEqual(values[i-benchWindow])
			}
			d.PushBack(v)
			_, _ = d.Front()
		}
	})
	if allocs != 0 {
		t.Errorf("a grown deque allocates %v times per pass, want 0", allocs)
	}
}

// Oracle: a deque driven over fixed windows agrees with RollingMin and RollingMax (and so with the
// naive rescan their tests use), with many ties.
func TestMonotonicDeque_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		n := 1 + rng.IntN(200)
		values := make([]int, n)
		for i := range values {
			values[i] = rng.IntN(1 + rng.IntN(20))
		}
		w := 1 + rng.IntN(n)
		if got, want := slidingExtremes(values, w, false), naiveWindowExtremes(values, w, false); !slices.Equal(got, want) {
			t.Fatalf("trial %d: sliding min %v, want %v", trial, got, want)
		}
		if got, want := slidingExtremes(values, w, true), naiveWindowExtremes(values, w, true); !slices.Equal(got, want) {
			t.Fatalf("trial %d: sliding max %v, want %v", trial, got, want)
		}
	}
}

func BenchmarkNewMonotonicDeque(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		d := NewMonotonicDeque[float64](true)
		benchSinkDeque = d.Len()
	}
}

// benchSinkDeque keeps the deque benchmarks' results observable.
var benchSinkDeque int

// BenchmarkMonotonicDequePushBack drives a sliding minimum of window benchWindow over n values: n
// pushes, n front reads and up to n evictions in total.
func BenchmarkMonotonicDequePushBack(b *testing.B) {
	for _, size := range benchSizes {
		values := benchFloats(size.n)
		d := NewMonotonicDeque[float64](false)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				for i, v := range values {
					if i >= benchWindow {
						d.PopFrontIfEqual(values[i-benchWindow])
					}
					d.PushBack(v)
				}
			}
			benchSinkDeque = d.Len()
		})
	}
}

func BenchmarkMonotonicDequeFront(b *testing.B) {
	var d MonotonicDeque[float64]
	d.PushBack(1)
	b.ReportAllocs()
	for b.Loop() {
		f, _ := d.Front()
		benchSinkDeque = int(f)
	}
}

// BenchmarkMonotonicDequePopFront pushes n ascending values (no evictions), then pops them all.
func BenchmarkMonotonicDequePopFront(b *testing.B) {
	for _, size := range benchSizes {
		d := NewMonotonicDeque[int](false)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				for i := range size.n {
					d.PushBack(i)
				}
				for d.Len() > 0 {
					d.PopFront()
				}
			}
		})
	}
}

func BenchmarkMonotonicDequePopFrontIfEqual(b *testing.B) {
	for _, size := range benchSizes {
		d := NewMonotonicDeque[int](false)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				for i := range size.n {
					d.PushBack(i)
				}
				for i := range size.n {
					d.PopFrontIfEqual(i)
				}
			}
		})
	}
}

func BenchmarkMonotonicDequeLen(b *testing.B) {
	var d MonotonicDeque[int]
	d.PushBack(1)
	b.ReportAllocs()
	for b.Loop() {
		benchSinkDeque = d.Len()
	}
}

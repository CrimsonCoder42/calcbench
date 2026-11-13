package heap

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

func intLess(a, b int) bool    { return a < b }
func intGreater(a, b int) bool { return a > b }

// drain pops every element of h in order.
func drain[T any](h *Heap[T]) []T {
	out := []T{}
	for {
		v, ok := h.Pop()
		if !ok {
			return out
		}
		out = append(out, v)
	}
}

// Cormen et al., Introduction to Algorithms, 3rd ed., §6.3, Figure 6.3: BUILD-MAX-HEAP on
// A = 4 1 3 2 16 9 10 14 8 7 leaves the array 16 14 10 8 7 9 3 2 4 1 (also reproduced by hand).
func TestHeapify_clrsFigure(t *testing.T) {
	in := []int{4, 1, 3, 2, 16, 9, 10, 14, 8, 7}
	h, err := Heapify(intGreater, in)
	if want := []int{16, 14, 10, 8, 7, 9, 3, 2, 4, 1}; err != nil || !slices.Equal(h.items, want) {
		t.Fatalf("Heapify(CLRS 6.3) layout = %v (%v), want %v", h.items, err, want)
	}
	if got, want := drain(h), []int{16, 14, 10, 9, 8, 7, 4, 3, 2, 1}; !slices.Equal(got, want) {
		t.Errorf("drain = %v, want %v", got, want)
	}
}

// Williams 1964: heapsort by repeated OUTHEAP. On Knuth's 16 keys a min-heap pops them in sorted order.
func TestHeapPop_knuthSorted(t *testing.T) {
	h, err := NewHeap(intLess, 0)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range knuth16 {
		if !h.Push(v) {
			t.Fatalf("Push(%d) = false", v)
		}
	}
	if h.Len() != 16 {
		t.Fatalf("Len = %d, want 16", h.Len())
	}
	want := []int{61, 87, 154, 170, 275, 426, 503, 509, 512, 612, 653, 677, 703, 765, 897, 908}
	if got := drain(h); !slices.Equal(got, want) {
		t.Errorf("drain = %v, want %v", got, want)
	}
}

func TestHeap_cases(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want []int
	}{
		{"empty", nil, []int{}},
		{"one", []int{7}, []int{7}},
		{"two", []int{9, 2}, []int{2, 9}},
		{"duplicates", []int{3, 1, 3, 1, 3}, []int{1, 1, 3, 3, 3}},
		{"sorted", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"reversed", []int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
		{"negative and zero", []int{0, -4, 2, -1}, []int{-4, -1, 0, 2}},
		{"extremes", []int{0, math.MaxInt, math.MinInt}, []int{math.MinInt, 0, math.MaxInt}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viaHeapify, err := Heapify(intLess, tt.in)
			if err != nil {
				t.Fatal(err)
			}
			if got := drain(viaHeapify); !slices.Equal(got, tt.want) {
				t.Errorf("Heapify then drain = %v, want %v", got, tt.want)
			}
			viaPush, _ := NewHeap(intLess, len(tt.in))
			for _, v := range tt.in {
				viaPush.Push(v)
			}
			if got := drain(viaPush); !slices.Equal(got, tt.want) {
				t.Errorf("Push then drain = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHeap_floatsMaxHeap(t *testing.T) {
	h, _ := Heapify(func(a, b float64) bool { return a > b }, []float64{1.5, math.Inf(-1), math.Inf(1), -0.5})
	if got, want := drain(h), []float64{math.Inf(1), 1.5, -0.5, math.Inf(-1)}; !slices.Equal(got, want) {
		t.Errorf("max-heap drain = %v, want %v", got, want)
	}
}

func TestHeap_structElements(t *testing.T) {
	type job struct {
		name     string
		priority int
	}
	h, _ := NewHeap(func(a, b job) bool { return a.priority < b.priority }, 4)
	h.Push(job{"write", 3})
	h.Push(job{"read", 1})
	h.Push(job{"lint", 2})
	if top, ok := h.Peek(); !ok || top.name != "read" {
		t.Errorf("Peek = (%v, %v), want read", top, ok)
	}
	if h.Len() != 3 {
		t.Errorf("Peek changed Len to %d", h.Len())
	}
}

func TestHeap_emptyAndZero(t *testing.T) {
	h, _ := NewHeap(intLess, 0)
	if v, ok := h.Pop(); ok || v != 0 {
		t.Errorf("Pop on empty = (%d, %v), want (0, false)", v, ok)
	}
	if v, ok := h.Peek(); ok || v != 0 {
		t.Errorf("Peek on empty = (%d, %v), want (0, false)", v, ok)
	}
	var zero Heap[int]
	if zero.Push(1) || zero.Len() != 0 {
		t.Errorf("zero Heap: Push stored a value (Len %d)", zero.Len())
	}
	if _, ok := zero.Pop(); ok {
		t.Error("zero Heap: Pop ok")
	}
	var nilHeap *Heap[int]
	if nilHeap.Push(1) || nilHeap.Len() != 0 {
		t.Error("nil Heap: Push or Len wrong")
	}
	if _, ok := nilHeap.Pop(); ok {
		t.Error("nil Heap: Pop ok")
	}
	if _, ok := nilHeap.Peek(); ok {
		t.Error("nil Heap: Peek ok")
	}
}

func TestNewHeap_errors(t *testing.T) {
	if h, err := NewHeap[int](nil, 4); !errors.Is(err, ErrNilLess) || h != nil {
		t.Errorf("NewHeap(nil) = (%v, %v), want (nil, ErrNilLess)", h, err)
	}
	h, err := NewHeap(intLess, -5)
	if err != nil || h.Len() != 0 {
		t.Errorf("NewHeap(less, -5) = (%v, %v), want an empty heap", h, err)
	}
}

// Rule L41: a capacity hint whose buffer exceeds one allocation's limit is an error, never a panic.
func TestNewHeap_tooLarge(t *testing.T) {
	for _, c := range []int{1 << 60, math.MaxInt} {
		if h, err := NewHeap(intLess, c); !errors.Is(err, ErrTooLarge) || h != nil {
			t.Errorf("NewHeap(less, %d) = (%v, %v), want (nil, ErrTooLarge)", c, h, err)
		}
	}
	// Zero-size elements take no memory, so any hint is allowed.
	if h, err := NewHeap(func(a, b struct{}) bool { return false }, 1<<60); err != nil || h.Len() != 0 {
		t.Errorf("NewHeap(struct{}, 1<<60) = (%v, %v), want an empty heap", h, err)
	}
}

func TestHeapify_errors(t *testing.T) {
	if h, err := Heapify[int](nil, []int{1}); !errors.Is(err, ErrNilLess) || h != nil {
		t.Errorf("Heapify(nil) = (%v, %v), want (nil, ErrNilLess)", h, err)
	}
}

func TestHeapify_doesNotModifyInput(t *testing.T) {
	in := slices.Clone(knuth16)
	h, _ := Heapify(intLess, in)
	h.Pop()
	h.Push(-1)
	if !slices.Equal(in, knuth16) {
		t.Errorf("Heapify modified its input: %v", in)
	}
}

// Oracle: a random mix of Push and Pop against a sorted slice that models the heap.
func TestHeap_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		start := make([]int, rng.IntN(40))
		for i := range start {
			start[i] = rng.IntN(50) - 25
		}
		h, _ := Heapify(intLess, start)
		model := slices.Sorted(slices.Values(start))
		for op := range 200 {
			if rng.IntN(3) == 0 {
				v, ok := h.Pop()
				if len(model) == 0 {
					if ok {
						t.Fatalf("trial %d op %d: Pop on empty gave %d", trial, op, v)
					}
					continue
				}
				if !ok || v != model[0] {
					t.Fatalf("trial %d op %d: Pop = (%d, %v), want %d", trial, op, v, ok, model[0])
				}
				model = model[1:]
			} else {
				v := rng.IntN(50) - 25
				h.Push(v)
				i, _ := slices.BinarySearch(model, v)
				model = slices.Insert(model, i, v)
			}
			if h.Len() != len(model) {
				t.Fatalf("trial %d op %d: Len = %d, want %d", trial, op, h.Len(), len(model))
			}
			if top, ok := h.Peek(); ok != (len(model) > 0) || (ok && top != model[0]) {
				t.Fatalf("trial %d op %d: Peek = (%d, %v)", trial, op, top, ok)
			}
		}
	}
}

func TestHeap_noAllocWithinCapacity(t *testing.T) {
	h, _ := NewHeap(intLess, 64)
	allocs := testing.AllocsPerRun(100, func() {
		for i := range 64 {
			h.Push(64 - i)
		}
		for h.Len() > 0 {
			h.Pop()
		}
		h.Peek()
	})
	if allocs != 0 {
		t.Errorf("Push/Pop/Peek within capacity allocated %v times, want 0", allocs)
	}
}

func randomInts(n int) []int {
	rng := rand.New(rand.NewPCG(1, 2))
	values := make([]int, n)
	for i := range values {
		values[i] = rng.IntN(n)
	}
	return values
}

func BenchmarkNewHeap(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		h, _ := NewHeap(intLess, 16)
		_ = h
	}
}

func BenchmarkHeapify(b *testing.B) {
	for _, size := range benchSizes {
		values := randomInts(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				h, _ := Heapify(intLess, values)
				_ = h
			}
		})
	}
}

// n pushes into a heap with room for them: O(n log n) worst case, O(n) average on random input.
func BenchmarkHeapPush(b *testing.B) {
	for _, size := range benchSizes {
		values := randomInts(size.n)
		h, _ := NewHeap(intLess, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				h.items = h.items[:0]
				for _, v := range values {
					h.Push(v)
				}
			}
		})
	}
}

// n pops from a full heap (restored by an O(n) copy each round): O(n log n).
func BenchmarkHeapPop(b *testing.B) {
	for _, size := range benchSizes {
		full, _ := Heapify(intLess, randomInts(size.n))
		h, _ := Heapify(intLess, full.items)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				h.items = h.items[:size.n]
				copy(h.items, full.items)
				for h.Len() > 0 {
					h.Pop()
				}
			}
		})
	}
}

func BenchmarkHeapPeek(b *testing.B) {
	h, _ := Heapify(intLess, randomInts(1_000))
	b.ReportAllocs()
	for b.Loop() {
		_, _ = h.Peek()
	}
}

func BenchmarkHeapLen(b *testing.B) {
	h, _ := Heapify(intLess, randomInts(1_000))
	b.ReportAllocs()
	for b.Loop() {
		_ = h.Len()
	}
}

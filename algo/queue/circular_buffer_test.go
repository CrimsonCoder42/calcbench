package queue

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// contents reads the buffer oldest first through At.
func contents[T any](b *CircularBuffer[T]) []T {
	out := []T{}
	for i := range b.Len() {
		v, _ := b.At(i)
		out = append(out, v)
	}
	return out
}

// Hand-derived (catalog): capacity 3; Overwrite keeps 2 3 4, Reject keeps 1 2 3 and refuses 4.
func TestCircularBuffer_example(t *testing.T) {
	over, err := NewCircularBuffer[int](3, Overwrite)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range []int{1, 2, 3, 4} {
		overwritten, err := over.Push(v)
		if err != nil || overwritten != (v == 4) {
			t.Fatalf("Overwrite Push(%d) = (%v, %v)", v, overwritten, err)
		}
	}
	if got := contents(over); !slices.Equal(got, []int{2, 3, 4}) {
		t.Errorf("Overwrite contents = %v, want [2 3 4]", got)
	}

	rej, err := NewCircularBuffer[int](3, Reject)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range []int{1, 2, 3} {
		if _, err := rej.Push(v); err != nil {
			t.Fatalf("Reject Push(%d) = %v", v, err)
		}
	}
	if overwritten, err := rej.Push(4); !errors.Is(err, ErrFull) || overwritten {
		t.Errorf("Reject Push(4) = (%v, %v), want (false, ErrFull)", overwritten, err)
	}
	if got := contents(rej); !slices.Equal(got, []int{1, 2, 3}) {
		t.Errorf("Reject contents = %v, want [1 2 3]", got)
	}
}

func TestCircularBuffer_accessors(t *testing.T) {
	b, _ := NewCircularBuffer[string](2, Overwrite)
	if _, ok := b.Pop(); ok {
		t.Error("Pop on empty returned ok")
	}
	if _, ok := b.Front(); ok {
		t.Error("Front on empty returned ok")
	}
	if _, ok := b.Back(); ok {
		t.Error("Back on empty returned ok")
	}
	b.Push("a")
	b.Push("b")
	b.Push("c") // overwrites "a"; the ring has wrapped
	if f, _ := b.Front(); f != "b" {
		t.Errorf("Front = %q, want b", f)
	}
	if l, _ := b.Back(); l != "c" {
		t.Errorf("Back = %q, want c", l)
	}
	for _, i := range []int{-1, 2, math.MaxInt} {
		if _, ok := b.At(i); ok {
			t.Errorf("At(%d) returned ok", i)
		}
	}
	if b.Len() != 2 || b.Cap() != 2 {
		t.Errorf("Len, Cap = %d, %d; want 2, 2", b.Len(), b.Cap())
	}
	if v, ok := b.Pop(); !ok || v != "b" {
		t.Errorf("Pop = (%q, %v), want (b, true)", v, ok)
	}
}

func TestCircularBuffer_capacityOne(t *testing.T) {
	b, _ := NewCircularBuffer[int](1, Overwrite)
	for v := range 5 {
		b.Push(v)
	}
	if got := contents(b); !slices.Equal(got, []int{4}) {
		t.Errorf("capacity 1 contents = %v, want [4]", got)
	}
}

func TestCircularBuffer_zeroValue(t *testing.T) {
	var b CircularBuffer[int]
	if _, err := b.Push(1); !errors.Is(err, ErrFull) {
		t.Errorf("zero-value Push = %v, want ErrFull", err)
	}
	if _, ok := b.Pop(); ok || b.Len() != 0 || b.Cap() != 0 {
		t.Error("zero value is not an empty buffer of capacity 0")
	}
}

func TestCircularBuffer_errors(t *testing.T) {
	for _, c := range []int{0, -1} {
		if b, err := NewCircularBuffer[int](c, Reject); !errors.Is(err, ErrCapacity) || b != nil {
			t.Errorf("NewCircularBuffer(%d) = (%v, %v), want (nil, ErrCapacity)", c, b, err)
		}
	}
	for _, c := range []int{math.MaxInt, maxAllocBytes/8 + 1} { // just past the limit, on 32- and 64-bit
		if b, err := NewCircularBuffer[int64](c, Reject); !errors.Is(err, ErrTooLarge) || b != nil {
			t.Errorf("NewCircularBuffer[int64](%d) = (%v, %v), want (nil, ErrTooLarge)", c, b, err)
		}
	}
	if b, err := NewCircularBuffer[[1 << 20]byte](maxAllocBytes>>20+1, Overwrite); !errors.Is(err, ErrTooLarge) || b != nil {
		t.Errorf("NewCircularBuffer[1 MiB](limit/1MiB + 1) = (%v, %v), want (nil, ErrTooLarge)", b, err)
	}
	// Zero-size elements need no memory; only the index bound MaxInt/2 applies.
	if b, err := NewCircularBuffer[struct{}](math.MaxInt/2+1, Reject); !errors.Is(err, ErrTooLarge) || b != nil {
		t.Errorf("NewCircularBuffer[struct{}](MaxInt/2+1) = (%v, %v), want (nil, ErrTooLarge)", b, err)
	}
	if b, err := NewCircularBuffer[struct{}](math.MaxInt/2, Reject); err != nil || b.Cap() != math.MaxInt/2 {
		t.Errorf("NewCircularBuffer[struct{}](MaxInt/2) = (%v, %v), want a buffer", b, err)
	}
	for _, p := range []FullPolicy{-1, 2} {
		if b, err := NewCircularBuffer[int](3, p); !errors.Is(err, ErrPolicy) || b != nil {
			t.Errorf("NewCircularBuffer(policy %d) = (%v, %v), want (nil, ErrPolicy)", p, b, err)
		}
	}
}

func TestCircularBuffer_noAllocAfterConstruction(t *testing.T) {
	b, _ := NewCircularBuffer[int](8, Overwrite)
	allocs := testing.AllocsPerRun(100, func() {
		for v := range 20 {
			b.Push(v)
		}
		for b.Len() > 0 {
			b.Pop()
		}
	})
	if allocs != 0 {
		t.Errorf("Push/Pop allocated %v times, want 0", allocs)
	}
}

// Oracle: a plain slice that drops its oldest element (Overwrite) or refuses (Reject) at capacity.
func TestCircularBuffer_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 200 {
		capacity := 1 + rng.IntN(6)
		policy := FullPolicy(rng.IntN(2))
		b, _ := NewCircularBuffer[int](capacity, policy)
		var model []int
		for step := range 200 {
			if rng.IntN(3) == 0 {
				v, ok := b.Pop()
				if ok != (len(model) > 0) || (ok && v != model[0]) {
					t.Fatalf("trial %d step %d: Pop = (%d, %v), model %v", trial, step, v, ok, model)
				}
				if ok {
					model = model[1:]
				}
			} else {
				v := rng.IntN(100)
				overwritten, err := b.Push(v)
				full := len(model) == capacity
				switch {
				case full && policy == Reject:
					if !errors.Is(err, ErrFull) || overwritten {
						t.Fatalf("trial %d step %d: Push = (%v, %v), want ErrFull", trial, step, overwritten, err)
					}
				case full:
					model = append(model[1:], v)
					if err != nil || !overwritten {
						t.Fatalf("trial %d step %d: Push = (%v, %v), want overwritten", trial, step, overwritten, err)
					}
				default:
					model = append(model, v)
					if err != nil || overwritten {
						t.Fatalf("trial %d step %d: Push = (%v, %v), want plain push", trial, step, overwritten, err)
					}
				}
			}
			if got := contents(b); !slices.Equal(got, model) {
				t.Fatalf("trial %d step %d: contents %v, want %v", trial, step, got, model)
			}
		}
	}
}

// n pushes into a buffer of capacity 64 under Overwrite, with a Pop after every third: O(1) per operation.
func BenchmarkCircularBuffer(b *testing.B) {
	for _, size := range benchSizes {
		buf, _ := NewCircularBuffer[int](64, Overwrite)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				for i := range size.n {
					buf.Push(i)
					if i%3 == 0 {
						buf.Pop()
					}
				}
			}
		})
	}
}

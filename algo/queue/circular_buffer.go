package queue

import (
	"errors"
	"fmt"
	"math"
	"unsafe"
)

// Errors of CircularBuffer. ErrCapacity and ErrPolicy are wrapped with the offending value; test with errors.Is.
var (
	// ErrCapacity: the capacity is less than 1.
	ErrCapacity = errors.New("queue: capacity must be at least 1")
	// ErrPolicy: the full-buffer policy is not Reject or Overwrite.
	ErrPolicy = errors.New("queue: unknown full-buffer policy")
	// ErrFull: Push on a full buffer whose policy is Reject.
	ErrFull = errors.New("queue: buffer is full")
	// ErrTooLarge: the slots would need more than maxAllocBytes of memory.
	ErrTooLarge = errors.New("queue: capacity needs more memory than one allocation may hold")
)

// maxAllocBytes is the largest single allocation a constructor will request: 2⁴⁷ bytes on 64-bit and
// 2³¹ on 32-bit, half the Go runtime's maxAlloc (runtime/malloc.go), above which make panics.
const maxAllocBytes = 1 << (31 + 16*(^uint(0)>>63))

// FullPolicy says what Push does when a CircularBuffer is full.
type FullPolicy int

const (
	// Reject makes Push on a full buffer return ErrFull and leave the buffer unchanged.
	Reject FullPolicy = iota
	// Overwrite makes Push on a full buffer replace the oldest element.
	Overwrite
)

// CircularBuffer is a first-in, first-out ring buffer of fixed capacity. Its slots are allocated once,
// by NewCircularBuffer; no operation allocates after that. When full, Push follows the policy chosen at
// construction: Reject or Overwrite. Not safe for concurrent use; the zero value has capacity 0 and
// rejects every Push with ErrFull.
//
// Methods: Push(v) adds v at the back. Pop removes and returns the oldest element. Front and Back return
// the oldest and newest without removing them. At(i) returns the i-th oldest (At(0) is Front). Pop,
// Front, Back and At return (zero, false) when there is no such element. Len is the element count, Cap
// the capacity.
//
// ID:      ALG-QUEUE-002
// Time:    O(1) worst case per operation; NewCircularBuffer O(capacity)
// Space:   O(capacity), one allocation in NewCircularBuffer
// Source:  Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.2 (a queue in a circular array, pointers advanced modulo M)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §10.1 (ENQUEUE and DEQUEUE with wrap-around)
// Tags:    queue, fifo, ring-buffer, circular-buffer, bounded, overwrite, data-structure
type CircularBuffer[T any] struct {
	slots  []T
	head   int // index of the oldest element
	count  int
	policy FullPolicy
}

// NewCircularBuffer returns an empty CircularBuffer with room for capacity elements and the given
// policy for Push on a full buffer.
// Errors: ErrCapacity when capacity < 1; ErrTooLarge when capacity × the size of T exceeds
// maxAllocBytes (2⁴⁷ bytes on 64-bit, 2³¹ on 32-bit) or capacity > math.MaxInt/2; ErrPolicy when policy is not Reject or Overwrite.
func NewCircularBuffer[T any](capacity int, policy FullPolicy) (*CircularBuffer[T], error) {
	if capacity < 1 {
		return nil, fmt.Errorf("%w: capacity = %d", ErrCapacity, capacity)
	}
	if capacity > math.MaxInt/2 {
		// head + count must stay below 2·capacity without overflowing int (see wrap).
		return nil, fmt.Errorf("%w: capacity = %d exceeds MaxInt/2", ErrTooLarge, capacity)
	}
	var zero T
	// Dividing the limit, not multiplying capacity, so the check itself cannot overflow.
	if size := uint64(unsafe.Sizeof(zero)); size > 0 && uint64(capacity) > maxAllocBytes/size {
		return nil, fmt.Errorf("%w: capacity = %d of %d-byte elements", ErrTooLarge, capacity, size)
	}
	if policy != Reject && policy != Overwrite {
		return nil, fmt.Errorf("%w: policy = %d", ErrPolicy, policy)
	}
	return &CircularBuffer[T]{slots: make([]T, capacity), policy: policy}, nil
}

// Push adds v at the back. On a full buffer, under Overwrite it replaces the oldest element and returns
// overwritten = true; under Reject it returns ErrFull and changes nothing. O(1).
func (b *CircularBuffer[T]) Push(v T) (overwritten bool, err error) {
	capacity := len(b.slots)
	if b.count == capacity {
		if b.policy != Overwrite || capacity == 0 {
			return false, ErrFull
		}
		b.slots[b.head] = v // the oldest slot becomes the newest
		b.head = b.wrap(b.head + 1)
		return true, nil
	}
	b.slots[b.wrap(b.head+b.count)] = v
	b.count++
	return false, nil
}

// Pop removes and returns the oldest element, or (zero, false) when the buffer is empty. O(1).
func (b *CircularBuffer[T]) Pop() (T, bool) {
	var zero T
	if b.count == 0 {
		return zero, false
	}
	v := b.slots[b.head]
	b.slots[b.head] = zero // release what v refers to
	b.head = b.wrap(b.head + 1)
	b.count--
	return v, true
}

// Front returns the oldest element without removing it, or (zero, false) when the buffer is empty. O(1).
func (b *CircularBuffer[T]) Front() (T, bool) {
	return b.At(0)
}

// Back returns the newest element without removing it, or (zero, false) when the buffer is empty. O(1).
func (b *CircularBuffer[T]) Back() (T, bool) {
	return b.At(b.count - 1)
}

// At returns the i-th oldest element (At(0) is the oldest), or (zero, false) when i is outside
// [0, Len()). O(1).
func (b *CircularBuffer[T]) At(i int) (T, bool) {
	if i < 0 || i >= b.count {
		var zero T
		return zero, false
	}
	return b.slots[b.wrap(b.head+i)], true
}

// Len returns the number of elements. O(1).
func (b *CircularBuffer[T]) Len() int {
	return b.count
}

// Cap returns the capacity. O(1).
func (b *CircularBuffer[T]) Cap() int {
	return len(b.slots)
}

// wrap reduces an index in [0, 2·capacity) to [0, capacity) without a division.
func (b *CircularBuffer[T]) wrap(i int) int {
	if i >= len(b.slots) {
		i -= len(b.slots)
	}
	return i
}

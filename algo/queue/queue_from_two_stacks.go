package queue

// QueueFromTwoStacks is a first-in, first-out queue held as two stacks: new elements are pushed on an
// input stack; elements leave from an output stack, which is refilled, by moving the whole input stack
// over in reverse, only when it is empty. Each element is moved once, so every operation is O(1)
// amortized. The zero value is an empty queue, ready to use. Not safe for concurrent use.
//
// Methods: Enqueue(v) adds v at the back. Dequeue removes and returns the front. Front returns the front
// without removing it. Dequeue and Front return (zero, false) on an empty queue. Len is the element count.
//
// ID:      ALG-QUEUE-001
// Time:    O(1) amortized per Enqueue, Dequeue and Front (one Dequeue or Front can cost O(n), n = Len); O(1) Len
// Space:   O(n) for n elements
// Source:  Okasaki, "Simple and efficient purely functional queues and deques", Journal of Functional Programming, 1995, 5(4):583–592, doi:10.1017/S0956796800001489 (§2, the batched two-list queue)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, Exercise 10.1-6 (a queue from two stacks)
// Tags:    queue, fifo, two-stacks, amortized, data-structure
type QueueFromTwoStacks[T any] struct {
	in, out []T // out's top (its last element) is the front of the queue
}

// Enqueue adds v at the back of the queue. O(1) amortized.
func (q *QueueFromTwoStacks[T]) Enqueue(v T) {
	q.in = append(q.in, v)
}

// Dequeue removes and returns the front element, or (zero, false) when the queue is empty.
// O(1) amortized.
func (q *QueueFromTwoStacks[T]) Dequeue() (T, bool) {
	if !q.refill() {
		var zero T
		return zero, false
	}
	n := len(q.out)
	v := q.out[n-1]
	var zero T
	q.out[n-1] = zero // release what v refers to
	q.out = q.out[:n-1]
	return v, true
}

// Front returns the front element without removing it, or (zero, false) when the queue is empty.
// O(1) amortized.
func (q *QueueFromTwoStacks[T]) Front() (T, bool) {
	if !q.refill() {
		var zero T
		return zero, false
	}
	return q.out[len(q.out)-1], true
}

// Len returns the number of elements. O(1).
func (q *QueueFromTwoStacks[T]) Len() int {
	return len(q.in) + len(q.out)
}

// refill moves the input stack, reversed, onto the empty output stack, and reports whether the queue
// holds any element.
func (q *QueueFromTwoStacks[T]) refill() bool {
	if len(q.out) > 0 {
		return true
	}
	if len(q.in) == 0 {
		return false
	}
	for i := len(q.in) - 1; i >= 0; i-- {
		q.out = append(q.out, q.in[i])
	}
	clear(q.in) // release references; the backing array is reused by later Enqueues
	q.in = q.in[:0]
	return true
}

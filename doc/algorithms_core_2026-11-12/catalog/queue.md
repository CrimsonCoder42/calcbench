# algo/queue (SHORT = QUEUE)

Builder W1-B3. Not yet in `doc/taxonomy.md` (lead to add: `ALG | algo/queue | QUEUE | FIFO queues: two-stack
queue, ring buffer`). The monotonic deque has one home, `algo/window` (`window.MonotonicDeque`, another
builder); it is not built here.

| ID | Name | Signature |
|---|---|---|
| ALG-QUEUE-001 | QueueFromTwoStacks | `type QueueFromTwoStacks[T any]` with `Enqueue(T)`, `Dequeue() (T, bool)`, `Front() (T, bool)`, `Len() int` |
| ALG-QUEUE-002 | CircularBuffer | `func NewCircularBuffer[T any](capacity int, policy FullPolicy) (*CircularBuffer[T], error)` with `Push(T) (overwritten bool, err error)`, `Pop() (T, bool)`, `Front() (T, bool)`, `Back() (T, bool)`, `At(i int) (T, bool)`, `Len() int`, `Cap() int` |

### ALG-QUEUE-001 QueueFromTwoStacks
- **Source**: Okasaki, "Simple and efficient purely functional queues and deques", Journal of Functional
  Programming 5(4):583–592, 1995, doi:10.1017/S0956796800001489 (verified 2026-11-12; §2 recalls the
  two-list "batched" queue: a front list and a reversed rear list, reversed when the front empties, O(1)
  amortized). Okasaki, *Purely Functional Data Structures*, Cambridge University Press, 1998, §5.2.
  Cormen et al., *Introduction to Algorithms*, 3rd ed., 2009, Exercise 10.1-6 (a queue from two stacks).
  Hood and Melville, "Real-time queue operations in pure LISP", Information Processing Letters 13(2), 1981,
  give the worst-case O(1) version (citation not re-verified; not used).
- **Best known**: O(1) amortized per operation (each element is pushed and popped on each stack once); O(1)
  worst case needs Hood–Melville's incremental reversal (or a ring buffer, ALG-QUEUE-002, when the
  capacity is bounded). O(n) space.
- **Approach**: two Go slices as stacks; Dequeue/Front move the whole input stack, reversed, to the output
  stack only when the output is empty. **Lost**: Hood–Melville real-time queue (worst-case O(1), several
  times the bookkeeping, no gain for a non-persistent Go value); a slice with re-slicing from the front
  (`q = q[1:]`) leaks the underlying array's head until growth copies.
- **Errors**: none; Dequeue/Front on an empty queue return (zero, false). Zero value is an empty queue.
  Not safe for concurrent use.
- **Worked example**: enqueue 1 2 3, dequeue → 1, enqueue 4, dequeue all → 2 3 4. Hand-derived.
- **Tags**: LeetCode 232.

### ALG-QUEUE-002 CircularBuffer
- **Source**: Knuth, TAOCP Vol. 1, 3rd ed., 1997, §2.2.2 (sequential allocation; a queue in a circular array
  of M slots, the rear and front pointers advancing modulo M); Cormen et al., 3rd ed., §10.1 (ENQUEUE and
  DEQUEUE on an array with wrap-around).
- **Best known**: O(1) worst case per operation, O(capacity) space, allocated once at construction.
- **Approach**: a slice of fixed length, a head index and a count; the tail is (head + count) mod capacity.
  Policy chosen by the caller: `Reject` (Push on a full buffer returns `ErrFull` and changes nothing) or
  `Overwrite` (Push on a full buffer replaces the oldest element and reports overwritten = true).
  **Lost**: head and tail indices with one slot left empty to tell full from empty (wastes a slot);
  power-of-two capacity with a mask (faster modulo but forces the capacity).
- **Errors**: `ErrCapacity` (capacity < 1), `ErrTooLarge` (capacity × element size > 2⁴⁷ bytes on 64-bit,
  2³¹ on 32-bit, rule L41; or capacity > MaxInt/2, so index sums cannot overflow), `ErrPolicy` (an unknown policy), `ErrFull` (Push under Reject
  when full). Pop/Front/Back/At on missing elements return (zero, false).
- **Worked example**: capacity 3, Overwrite: push 1 2 3 4 → 2 3 4 (1 overwritten); Reject: push 1 2 3 4 →
  ErrFull on 4, holds 1 2 3. Hand-derived.
- **Tags**: LeetCode 622, 641.

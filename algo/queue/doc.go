// Package queue holds first-in, first-out queues: a queue built from two stacks and a fixed-capacity
// circular buffer.
//
// ID short form: QUEUE (IDs are ALG-QUEUE-NNN).
//
// Every queue's zero value or constructor result is ready to use; none is safe for concurrent use.
// Removing or reading a missing element returns (zero, false), never a panic.
// The monotonic deque (sliding-window minimum and maximum) is not here: see algo/window.MonotonicDeque.
//
// Functions:
//
//	ALG-QUEUE-001  QueueFromTwoStacks  FIFO queue from two stacks, O(1) amortized per operation
//	ALG-QUEUE-002  CircularBuffer      fixed-capacity ring buffer that rejects or overwrites when full
package queue

// Package window computes sliding-window sums, means, minima and maxima over a sequence.
//
// ID short form: WINDOW (IDs are ALG-WINDOW-NNN).
//
// Windows: for windowLen = w and n values, window i (i = w-1 … n-1) is values[i-w+1 : i+1], so the
// output has n-w+1 elements and each includes the current value. A window longer than the data gives
// no windows (doc/ASSUMPTIONS.md A-004).
//
// Output: every function appends its results to dst and returns the extended slice, like
// strconv.AppendInt (doc/DECISIONS.md D-003). Pass nil to allocate once, sized up front, or buf[:0]
// to reuse a buffer. On error, dst is returned with its length unchanged.
//
// Numbers: RollingSum and RollingMean take float64 and reject NaN and ±Inf. RollingMin and RollingMax
// take any ordered type and use cmp.Compare's total order, in which NaN is the least value (D-002).
// No function modifies values. Results are in the unit of values.
//
// Variable windows: MinLengthSubarraySumAtLeast, LongestWithAtMostKDistinct and LongestWithoutRepeat
// find one run values[start : start+length] by two pointers that move right only, and return
// (start, length); ties go to the run that starts first. ExponentialMovingAverage weights all past
// values, not a fixed window. MonotonicDeque is the monotonic queue of RollingMin and RollingMax (Lemire 2006),
// exported for windows the caller drives (by time, by count, or streaming); the zero value is a
// min-deque. It is the one home of the monotonic queue in calcbench.
//
// Functions:
//
//	ALG-WINDOW-001  RollingSum   compensated sum of each window
//	ALG-WINDOW-002  RollingMean  mean of each window (the single moving average)
//	ALG-WINDOW-003  RollingMin   least element of each window
//	ALG-WINDOW-004  RollingMax   greatest element of each window
//	ALG-WINDOW-005  NewMonotonicDeque               an empty min- or max-deque
//	ALG-WINDOW-006  MonotonicDeque.PushBack         evict beaten elements from the back, then append, O(1) amortized
//	ALG-WINDOW-007  MonotonicDeque.Front            the least (or greatest) element held
//	ALG-WINDOW-008  MonotonicDeque.PopFront         remove the front element
//	ALG-WINDOW-009  MonotonicDeque.PopFrontIfEqual  remove the front element if it is the value leaving the window
//	ALG-WINDOW-010  MonotonicDeque.Len              the number of elements held
//	ALG-WINDOW-011  MinLengthSubarraySumAtLeast     shortest run of positive values with sum at least a target
//	ALG-WINDOW-012  LongestWithAtMostKDistinct      longest run with at most k distinct values
//	ALG-WINDOW-013  LongestWithoutRepeat            longest run with no repeated value
//	ALG-WINDOW-014  ExponentialMovingAverage        single exponential smoothing (NIST §6.4.3.1)
//	ALG-WINDOW-015  MovingMedian                    median of each window (mean of the two middles when even), two heaps
//	ALG-WINDOW-016  CoOccurrenceWithinWindow  per ordered label pair, events of A with an event of B within the window
package window

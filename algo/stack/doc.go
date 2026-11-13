// Package stack holds stack algorithms: monotonic stacks (nearest greater and smaller elements, the
// largest rectangle in a histogram or a grid), a stack with O(1) minimum, and stack evaluators of
// reverse Polish expressions and bracket balance.
//
// ID short form: STACK (IDs are ALG-STACK-NNN).
//
// Ordering: values are compared with cmp.Compare, a total order in which a floating-point NaN sorts
// before every other value and -0 equals +0 (doc/DECISIONS.md D-002). Comparisons are strict: an equal
// value is neither greater nor smaller.
//
// Output: functions that give one index per element append to dst and return the extended slice, like
// strconv.AppendInt (D-003). Pass nil to allocate once, sized up front, or buf[:0] to reuse a buffer.
// An index is -1 where no element qualifies. No function modifies its input.
//
// Unix path simplification is not here: it is the standard library's path.Clean on a rooted path.
//
// Functions:
//
//	ALG-STACK-001  NextGreater                index of the next strictly greater element
//	ALG-STACK-002  NextSmaller                index of the next strictly smaller element
//	ALG-STACK-003  PreviousGreater            index of the previous strictly greater element
//	ALG-STACK-004  PreviousSmaller            index of the previous strictly smaller element
//	ALG-STACK-005  LargestRectangleHistogram  area of the largest rectangle under a histogram
//	ALG-STACK-006  MaximalRectangle           area of the largest all-true rectangle in a grid
//	ALG-STACK-007  DailyTemperatures          positions to wait for a strictly greater reading
//	ALG-STACK-008  MinStack                   stack with Push, Pop, Top and Min in O(1)
//	ALG-STACK-009  EvaluateRPN                value of an int64 reverse Polish expression
//	ALG-STACK-010  BalancedBrackets           whether brackets of the given pairs are balanced
package stack

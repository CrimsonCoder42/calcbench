// Package prefix computes prefix sums and what they answer: range sums in one and two dimensions,
// range updates through a difference array, subarray counts and maxima, and prefix products.
//
// ID short form: PREFIX (IDs are ALG-PREFIX-NNN).
//
// Ranges are half-open, Go style: [lo, hi) is values[lo] … values[hi-1], and lo == hi is empty.
// PrefixSums returns n+1 sums for n values, the first being 0, so RangeSum(prefix, lo, hi) is
// prefix[hi] − prefix[lo].
//
// Numbers: every function takes any built-in integer type (signed or unsigned, 8 to 64 bits). Floats
// are not accepted: a float prefix sum loses the small values to cancellation (see algo/window's
// compensated RollingSum instead). Overflow is never silent. Results that must fit the element type
// return ErrOverflow when they do not; intermediate values that need not fit it are carried exactly in
// 128-bit arithmetic (CountSubarraysSumEquals, EquilibriumIndex, RectangleSum, DifferenceArray), so
// those report an error only when the answer itself is out of range.
//
// Output: functions that build a slice append to dst and return the extended slice (doc/DECISIONS.md
// D-003). Pass nil to allocate once, sized up front. On error dst is returned with its length
// unchanged. No function modifies its input.
//
// Functions:
//
//	ALG-PREFIX-001  PrefixSums                   the n+1 prefix sums of a slice
//	ALG-PREFIX-002  RangeSum                     the sum of values[lo:hi] from its prefix sums, O(1)
//	ALG-PREFIX-003  PrefixSums2D                 the summed-area table of a grid
//	ALG-PREFIX-004  SummedAreaTable.RectangleSum the sum of a rectangle of the grid, O(1)
//	ALG-PREFIX-005  NewDifferenceArray           a difference array holding a copy of a slice
//	ALG-PREFIX-006  DifferenceArray.RangeAdd     add a constant to a range, O(1)
//	ALG-PREFIX-007  DifferenceArray.Materialize  the current elements, O(n)
//	ALG-PREFIX-008  CountSubarraysSumEquals      the number of subarrays with a given sum
//	ALG-PREFIX-009  MaxSubarray                  the largest subarray sum and its range (Kadane)
//	ALG-PREFIX-010  ProductExceptSelf            for each i, the product of all values but values[i], no division
//	ALG-PREFIX-011  EquilibriumIndex             the first index whose left and right sums are equal
//	ALG-PREFIX-012  DifferenceArray.Len          the number of elements
//	ALG-PREFIX-013  SummedAreaTable.Rows         the number of grid rows
//	ALG-PREFIX-014  SummedAreaTable.Cols         the number of grid columns
package prefix

package stack

import (
	"errors"
	"fmt"
	"math"
	"math/bits"
)

// Errors of the rectangle functions. They are wrapped with the offending input; test with errors.Is.
var (
	// ErrNegativeHeight: a histogram bar is lower than 0.
	ErrNegativeHeight = errors.New("stack: height must not be negative")
	// ErrOverflow: a result does not fit in int64.
	ErrOverflow = errors.New("stack: result exceeds the int64 range")
)

// LargestRectangleHistogram returns the area of the largest axis-aligned rectangle that fits under the
// histogram whose bars, each one unit wide, have the given heights. Empty input gives 0.
// Units: heights in any unit u; the area is in u × bar widths. heights is read, never modified.
// Errors: ErrNegativeHeight when a height is below 0; ErrOverflow when height × width exceeds int64.
//
// ID:      ALG-STACK-005
// Time:    O(n)   n = len(heights); each index is pushed and popped once
// Space:   O(n) extra for the index stack (one allocation)
// Source:  Berkman, Schieber, Vishkin, "Optimal doubly logarithmic parallel algorithms based on finding all nearest smaller values", Journal of Algorithms, 1993, 14(3):344–370, doi:10.1006/jagm.1993.1018
// Source:  Gabow, Bentley, Tarjan, "Scaling and related techniques for geometry problems", Proc. 16th ACM STOC, 1984, pp. 135–143, doi:10.1145/800057.808675
// Tags:    stack, monotonic-stack, histogram, largest-rectangle, area, ansv
func LargestRectangleHistogram(heights []int64) (int64, error) {
	for i, h := range heights {
		if h < 0 {
			return 0, fmt.Errorf("%w: heights[%d] = %d", ErrNegativeHeight, i, h)
		}
	}
	area, _, err := largestRectangle(heights, make([]int, 0, len(heights)+1))
	return area, err
}

// largestRectangle is the one-pass monotonic-stack scan over non-negative heights. stack is scratch
// space (len 0, any capacity); it is returned so a caller can reuse it.
func largestRectangle(heights []int64, stack []int) (int64, []int, error) {
	n := len(heights)
	var best int64
	for i := 0; i <= n; i++ {
		h := int64(-1) // past the end: a sentinel lower than every bar, which closes every rectangle
		if i < n {
			h = heights[i]
		}
		// Popping on >= lets the last of a run of equal bars measure the run's full width.
		for len(stack) > 0 && heights[stack[len(stack)-1]] >= h {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			left := -1 // the previous smaller bar, or the left edge
			if len(stack) > 0 {
				left = stack[len(stack)-1]
			}
			width := i - left - 1
			hi, lo := bits.Mul64(uint64(heights[top]), uint64(width))
			if hi != 0 || lo > math.MaxInt64 {
				return 0, stack, fmt.Errorf("%w: height %d × width %d", ErrOverflow, heights[top], width)
			}
			best = max(best, int64(lo))
		}
		stack = append(stack, i)
	}
	return best, stack[:0], nil
}

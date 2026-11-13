package series

import (
	"fmt"
	"time"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// WindowTotal is the sum and the number of the values in one time window.
type WindowTotal struct {
	Sum   float64
	Count int
}

// TrailingTimeWindowSum returns, for each point i, the sum and count of the values of the points j ≤ i
// whose time lies in the trailing window (points[i].At − window, points[i].At]: a sliding window measured
// in time, not in a number of points, so irregular sampling is handled. Two pointers: the window's left
// edge only moves forward. Sums are compensated (Neumaier) so subtracting the values that leave the
// window does not accumulate error. points is read, never modified.
// Errors: ErrOutOfRange when window ≤ 0; ErrUnsorted when a point's time is before the previous one's;
// ErrNotFinite when a value is NaN or ±Inf; ErrOverflow when a window's sum is beyond the float64 range.
//
// ID:      NUM-SERIES-003
// Time:    O(n)   n = len(points); each point enters and leaves the window once (amortized two pointers)
// Space:   O(1) extra; O(n) for the result (one allocation)
// Source:  Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 17 "Amortized Analysis" (the aggregate method bounds the two-pointer pass)
// Source:  Neumaier A, "Rundungsfehleranalyse einiger Verfahren zur Summation endlicher Summen", ZAMM, 1974, 54(1):39–51, doi:10.1002/zamm.19740540106 (compensated summation)
// Units:   values in any one unit; window a time.Duration; Sum in the unit of the values; Count points
// Example: values 30, 45, 20 at day 0 07:00, day 1 18:00, day 1 19:10, window 7 days → (30, 1), (75, 2), (95, 3) (hand-checked)
// Tags:    rolling-sum, sliding-window, time-window, trailing, two-pointers, irregular-sampling, time-series
func TrailingTimeWindowSum(points []Point, window time.Duration) ([]WindowTotal, error) {
	if window <= 0 {
		return nil, fmt.Errorf("%w: window = %v, want > 0", ErrOutOfRange, window)
	}
	if err := checkPoints(points); err != nil {
		return nil, err
	}
	out := make([]WindowTotal, len(points))
	lo := 0
	var sum, comp float64
	for hi, p := range points {
		sum, comp = reduce.NeumaierAdd(sum, comp, p.Value)
		edge := p.At.Add(-window)
		for !points[lo].At.After(edge) {
			sum, comp = reduce.NeumaierAdd(sum, comp, -points[lo].Value)
			lo++
		}
		total := sum + comp
		if !reduce.AllFinite(total) || !reduce.AllFinite(sum) {
			return nil, fmt.Errorf("%w: window ending at points[%d]", ErrOverflow, hi)
		}
		out[hi] = WindowTotal{Sum: total, Count: hi - lo + 1}
	}
	return out, nil
}

// checkPoints returns ErrNotFinite for a NaN or ±Inf value and ErrUnsorted for a time before the
// previous point's.
func checkPoints(points []Point) error {
	for i, p := range points {
		if !reduce.AllFinite(p.Value) {
			return fmt.Errorf("%w: points[%d].Value = %v", ErrNotFinite, i, p.Value)
		}
		if i > 0 && p.At.Before(points[i-1].At) {
			return fmt.Errorf("%w: points[%d].At = %v is before points[%d].At = %v", ErrUnsorted, i, p.At, i-1, points[i-1].At)
		}
	}
	return nil
}

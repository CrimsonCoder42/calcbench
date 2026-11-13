package series

import (
	"fmt"
	"math"
	"math/bits"
	"time"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// LocalHourWindowMean returns the arithmetic mean and the count of the values of the points whose local
// clock hour in loc lies in the window [fromHour, toHour): fromHour ≤ h < toHour, or, when
// fromHour > toHour, a window that wraps past midnight (h ≥ fromHour or h < toHour). For example
// [5, 12) is 05:00–11:59 and [22, 6) is 22:00–05:59. The hour is that of the local wall clock on the
// day of the point (package time with the IANA rules of loc), so a DST change moves which instants fall
// in the window. Two windows' means and their difference (a morning–evening comparison, say) are two
// calls and a subtraction; any threshold on them is the caller's. Points need not be sorted. The sum is
// compensated. points is read, never modified.
// Errors: ErrNilLocation when loc is nil; ErrOutOfRange unless 0 ≤ fromHour ≤ 23, 0 ≤ toHour ≤ 24 and
// fromHour ≠ toHour; ErrNotFinite when a value in the window is NaN or ±Inf; ErrTooFew when no point
// falls in the window.
//
// ID:      NUM-SERIES-011
// Time:    O(n)   n = len(points); one pass (two only when the sum overflows)
// Space:   O(1) extra, no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §1.3.5.1 "Measures of Location", https://www.itl.nist.gov/div898/handbook/eda/section3/eda351.htm (read 2026-11-12) (mean Ȳ = Σ Yᵢ/N)
// Source:  ISO 8601-1:2019, Date and time — Representations for information interchange — Part 1: Basic rules (local time of day: hour 00–23)
// Source:  Neumaier A, "Rundungsfehleranalyse einiger Verfahren zur Summation endlicher Summen", ZAMM, 1974, 54(1):39–51, doi:10.1002/zamm.19740540106 (compensated summation)
// Tags:    time-of-day, hour-window, diurnal, mean, local-time, time-series
// Units:   values in any one unit; hours 0–24 of the local clock; mean in the unit of the values
// Example: 140 at 06:30, 150 at 11:59, 120 at 18:00 local, window [5, 12) → mean 145, count 2; window [17, 24) → 120, count 1 (hand-checked)
func LocalHourWindowMean(points []Point, loc *time.Location, fromHour, toHour int) (mean float64, count int, err error) {
	if loc == nil {
		return 0, 0, fmt.Errorf("%w: loc", ErrNilLocation)
	}
	if fromHour < 0 || fromHour > 23 || toHour < 0 || toHour > 24 || fromHour == toHour {
		return 0, 0, fmt.Errorf("%w: hours [%d, %d), want 0 ≤ from ≤ 23, 0 ≤ to ≤ 24, from ≠ to", ErrOutOfRange, fromHour, toHour)
	}
	in := func(t time.Time) bool {
		h := t.In(loc).Hour()
		if fromHour < toHour {
			return h >= fromHour && h < toHour
		}
		return h >= fromHour || h < toHour
	}
	var sum, comp float64
	for i, p := range points {
		if !in(p.At) {
			continue
		}
		if !reduce.AllFinite(p.Value) {
			return 0, 0, fmt.Errorf("%w: points[%d].Value = %v", ErrNotFinite, i, p.Value)
		}
		sum, comp = reduce.NeumaierAdd(sum, comp, p.Value)
		count++
	}
	if count == 0 {
		return 0, 0, fmt.Errorf("%w: no point in hours [%d, %d)", ErrTooFew, fromHour, toHour)
	}
	if t := sum + comp; reduce.AllFinite(t) {
		return t / float64(count), count, nil
	}
	// The sum overflowed: redo on values scaled by 2^-k, 2^k > count (exact for normal numbers); the mean
	// of finite values lies between their extremes, so it is finite.
	k := bits.Len(uint(count))
	sum, comp = 0, 0
	for _, p := range points {
		if in(p.At) {
			sum, comp = reduce.NeumaierAdd(sum, comp, math.Ldexp(p.Value, -k))
		}
	}
	return math.Ldexp((sum+comp)/float64(count), k), count, nil
}

package series

import (
	"fmt"
	"math/rand/v2"
	"time"
)

// secondsPerDay is the length of a day of elapsed time, 86 400 s (SI: d = 24 h = 86 400 s; BIPM, The
// International System of Units, 9th ed., 2019, Table 8).
const secondsPerDay = 86400

// TheilSenSlopePerDay returns the Theil–Sen slope (NUM-SERIES-001) of the points' values against elapsed
// time in days: x is the time since the first point, in days of 86 400 s of elapsed (not calendar) time,
// so a DST day does not stretch the axis. Points need not be sorted; points at the same instant form no
// slope with each other. points is read, never modified.
// Errors: as NUM-SERIES-001 (ErrNilRandom, ErrTooFew when fewer than 2 points or all at one instant,
// ErrNotFinite for a NaN or ±Inf value, ErrOverflow).
//
// ID:      NUM-SERIES-002
// Time:    O(n log n) expected   n = len(points)
// Space:   O(max(n, 2²²)) extra (the x and y arrays plus NUM-SERIES-001's buffers); result is a scalar
// Source:  Sen PK, "Estimates of the regression coefficient based on Kendall's tau", Journal of the American Statistical Association, 1968, 63(324):1379–1389, doi:10.1080/01621459.1968.10480934
// Source:  Dillencourt MB, Mount DM, Netanyahu NS, "A randomized algorithm for slope selection", International Journal of Computational Geometry & Applications, 1992, 2(1):1–27, doi:10.1142/S0218195992000020
// Units:   point values in unit Y; result in unit Y per day (86 400 s)
// Example: values 10, 12, 11 at 0, 1 and 2 days → slopes 2, 0.5, −1 → 0.5 per day (hand-checked)
// Tags:    theil-sen, sen-slope, robust-regression, trend, slope-per-day, time-series
func TheilSenSlopePerDay(points []Point, rng *rand.Rand) (float64, error) {
	if rng == nil {
		return 0, fmt.Errorf("%w: rng", ErrNilRandom)
	}
	if len(points) < 2 {
		return 0, fmt.Errorf("%w: len(points) = %d, want ≥ 2", ErrTooFew, len(points))
	}
	x := make([]float64, 2*len(points))
	x, y := x[:len(points)], x[len(points):]
	origin := points[0].At
	for i, p := range points {
		x[i] = elapsedDays(origin, p.At)
		y[i] = p.Value
	}
	return TheilSenSlope(x, y, rng)
}

// elapsedDays returns the elapsed time from origin to t in days of 86 400 s, without the ±292-year
// saturation of time.Time.Sub.
func elapsedDays(origin, t time.Time) float64 {
	s := float64(t.Unix() - origin.Unix())
	ns := float64(t.Nanosecond() - origin.Nanosecond())
	return (s + ns/1e9) / secondsPerDay
}

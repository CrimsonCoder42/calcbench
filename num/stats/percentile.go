package stats

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/selection"
)

// QuantileMethod selects one of the nine sample-quantile definitions of Hyndman and Fan (1996), by their
// type number. Every definition is Q(p) = (1 − γ)·X(j) + γ·X(j+1), with X(1) ≤ … ≤ X(n) the sorted
// values, j = ⌊np + m⌋ and g = np + m − j; the types differ in m and γ:
//
//	type 1  m = 0          γ = 0 if g = 0, else 1           (inverse of the empirical distribution function)
//	type 2  m = 0          γ = 1/2 if g = 0, else 1         (as type 1, averaging at discontinuities)
//	type 3  m = −1/2       γ = 0 if g = 0 and j even, else 1 (nearest even order statistic)
//	type 4  m = 0          γ = g
//	type 5  m = 1/2        γ = g
//	type 6  m = p          γ = g                             (NIST e-Handbook §7.2.6.2)
//	type 7  m = 1 − p      γ = g                             (the default of R and of spreadsheets' PERCENTILE)
//	type 8  m = (p + 1)/3  γ = g                             (recommended by Hyndman and Fan)
//	type 9  m = p/4 + 3/8  γ = g
//
// An index below 1 uses X(1) and one above n uses X(n).
type QuantileMethod int

// The nine Hyndman–Fan sample-quantile types.
const (
	HyndmanFan1 QuantileMethod = 1 + iota
	HyndmanFan2
	HyndmanFan3
	HyndmanFan4
	HyndmanFan5
	HyndmanFan6
	HyndmanFan7
	HyndmanFan8
	HyndmanFan9
)

// DefaultQuantileMethod is the documented default, Hyndman–Fan type 7 (doc/DECISIONS.md D-010). Go has no
// default arguments: pass it explicitly, so the method is visible at every call site.
const DefaultQuantileMethod = HyndmanFan7

// Percentile returns the p-quantile of values (p = 0.9 for the 90th percentile) by the Hyndman–Fan
// definition method. It selects the one or two order statistics it needs in linear time, without sorting.
// A fractional part g within 4ε·max(1, |np + m|) of an integer (ε = 2^-52) is taken as that integer, so
// rounding in np + m cannot move a discontinuous type to the next order statistic.
// values is read, never modified: a copy is reordered.
// Errors: ErrEmpty when values is empty; ErrNotFinite when a value is NaN or ±Inf; ErrOutOfRange when p
// is not in [0, 1] (NaN included) or method is not HyndmanFan1 … HyndmanFan9.
//
// ID:      NUM-STATS-015
// Time:    O(n) worst case   n = len(values); one linear-time selection and one pass
// Space:   O(n) extra (one allocation, the selection's copy) plus O(log n) stack
// Source:  Hyndman, Fan, "Sample quantiles in statistical packages", The American Statistician, 1996, 50(4):361–365, doi:10.1080/00031305.1996.10473566 (eq. 1; definitions 1–3 §2, 4–9 §3 and Table 3)
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §7.2.6.2 "Percentiles", https://www.itl.nist.gov/div898/handbook/prc/section2/prc262.htm (read 2026-11-12)
// Tags:    percentile, quantile, order-statistic, hyndman-fan, descriptive
// Units:   values in any one unit; p dimensionless in [0, 1]; result in the unit of values
// Example: e-Handbook §7.2.6.2, the 12 resistivity values, p = 0.90: type 6 → 95.1981, type 7 → 95.1957, type 8 → 95.1972 (as printed)
func Percentile(values []float64, p float64, method QuantileMethod) (float64, error) {
	if err := checkSample(values, 1); err != nil {
		return 0, err
	}
	if !(p >= 0 && p <= 1) {
		return 0, fmt.Errorf("%w: p = %v, want 0..1", ErrOutOfRange, p)
	}
	if method < HyndmanFan1 || method > HyndmanFan9 {
		return 0, fmt.Errorf("%w: method = %d, want HyndmanFan1..HyndmanFan9", ErrOutOfRange, int(method))
	}
	n := len(values)
	j, gamma := quantilePosition(n, p, method)
	lo := min(max(j, 1), n)   // rank of X(j), clamped to 1..n
	hi := min(max(j+1, 1), n) // rank of X(j+1), clamped to 1..n
	if gamma == 0 || lo == hi {
		if gamma == 1 {
			lo = hi
		}
		x, _ := selection.KthSmallest(values, lo) // cannot fail: 1 ≤ lo ≤ n
		return x, nil
	}
	a, _ := selection.KthSmallest(values, lo) // cannot fail: 1 ≤ lo ≤ n
	b := nextOrderStatistic(values, a, lo)
	return lerp(a, b, gamma), nil
}

// quantilePosition returns j and γ of Hyndman–Fan definition method for n values and probability p.
func quantilePosition(n int, p float64, method QuantileMethod) (j int, gamma float64) {
	var m float64
	switch method {
	case HyndmanFan3:
		m = -0.5
	case HyndmanFan5:
		m = 0.5
	case HyndmanFan6:
		m = p
	case HyndmanFan7:
		m = 1 - p
	case HyndmanFan8:
		m = (p + 1) / 3
	case HyndmanFan9:
		m = p/4 + 3.0/8
	}
	x := float64(n)*p + m
	fl := math.Floor(x)
	g := x - fl
	tol := 4 * 0x1p-52 * max(1, math.Abs(x))
	if g < tol {
		g = 0
	} else if g > 1-tol {
		fl++
		g = 0
	}
	j = int(fl)
	switch method {
	case HyndmanFan1:
		gamma = 1
		if g == 0 {
			gamma = 0
		}
	case HyndmanFan2:
		gamma = 1
		if g == 0 {
			gamma = 0.5
		}
	case HyndmanFan3:
		gamma = 1
		if g == 0 && j%2 == 0 {
			gamma = 0
		}
	default:
		gamma = g
	}
	return j, gamma
}

// nextOrderStatistic returns X(k+1) given a = X(k), 1 ≤ k < n: a itself when more than k values are ≤ a
// (a tie straddles ranks k and k+1), otherwise the least value greater than a. One pass, no allocation.
func nextOrderStatistic(values []float64, a float64, k int) float64 {
	atMost := 0
	next := math.Inf(1)
	for _, v := range values {
		if v <= a {
			atMost++
		} else if v < next {
			next = v
		}
	}
	if atMost > k {
		return a
	}
	return next
}

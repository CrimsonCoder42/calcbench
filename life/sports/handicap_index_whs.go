package sports

import (
	"errors"
	"fmt"
	"math"
)

var (
	// ErrNotFinite is returned when an input is NaN or ±Inf.
	ErrNotFinite = errors.New("sports: value is NaN or infinite")
	// ErrOutOfRange is returned when an input is outside the range the source defines.
	ErrOutOfRange = errors.New("sports: input outside the range the source defines")
)

// whsTable52a is Rule 5.2a's table, indexed by the number of Score Differentials (3–20): how many of the
// lowest are averaged, and the adjustment added.
var whsTable52a = [21]struct {
	lowest     int
	adjustment float64
}{
	3: {1, -2.0}, 4: {1, -1.0}, 5: {1, 0}, 6: {2, -1.0}, 7: {2, 0}, 8: {2, 0},
	9: {3, 0}, 10: {3, 0}, 11: {3, 0}, 12: {4, 0}, 13: {4, 0}, 14: {4, 0},
	15: {5, 0}, 16: {5, 0}, 17: {6, 0}, 18: {6, 0}, 19: {7, 0}, 20: {8, 0},
}

// whsMaxHandicapIndex is Rule 5.3's maximum Handicap Index.
const whsMaxHandicapIndex = 54.0

// HandicapIndexWHS returns the World Handicap System Handicap Index from a scoring record of Score
// Differentials, oldest first: of the most recent 20 (all of them when there are fewer), the average of the
// lowest k plus an adjustment, by Rule 5.2a's table (3 → lowest 1 − 2.0; 4 → lowest 1 − 1.0; 5 → lowest 1;
// 6 → lowest 2 − 1.0; 7–8 → 2; 9–11 → 3; 12–14 → 4; 15–16 → 5; 17–18 → 6; 19 → 7; 20 → lowest 8), rounded to
// the nearest tenth (Rule 5.2b), and at most 54.0 (Rule 5.3). The source does not state a tie rule for the
// rounding; ties (an exact .x5) round half to even, the house default. Score Differentials themselves (Rule
// 5.1) and exceptional-score or Committee adjustments are the caller's. The input is not modified.
// Each differential is first taken to the nearest tenth (Rule 5.1 states them in tenths).
// Errors: ErrNotFinite for a NaN or ±Inf differential; ErrOutOfRange for fewer than 3 differentials or a
// differential beyond ±1000 (an overflow guard, not a Rules limit).
//
// ID:      LIF-SPORT-001
// Time:    O(n)   n = len(differentials); the lowest-k selection is over at most 20 values (constant)
// Space:   O(1), no allocation (a fixed 20-element array)
// Source:  The R&A and USGA. Rules of Handicapping (World Handicap System), effective 1 January 2024. Rule 5.2a Table "Number of Score Differentials in scoring record / Score Differential(s) to be used / Adjustment"; Rule 5.2b "Average the lowest 8 of the most recent 20 Score Differentials ... and round to the nearest tenth"; Rule 5.3 "The maximum Handicap Index a player can obtain is 54.0" (randa.org/roh/the-rules-of-handicapping/rule-5, read 2026-11-13)
// Tags:    golf, handicap, whs, world-handicap-system, sports
// Units:   differentials strokes (Score Differentials); result strokes (Handicap Index, tenths)
// Range:   len(differentials) ≥ 3; only the last 20 are used; |differential| ≤ 1000 (overflow guard)
// Example: 3 differentials 10.0, 12.0, 14.0 → 10.0 − 2.0 = 8.0; 20 differentials 1.0 … 20.0 → mean of 1 … 8 = 4.5 (HAND-DERIVED from Table 5.2a)
func HandicapIndexWHS(differentials []float64) (float64, error) {
	n := len(differentials)
	if n < 3 {
		return 0, fmt.Errorf("%w: %d Score Differentials, want ≥ 3 (Rule 5.2a)", ErrOutOfRange, n)
	}
	recent := differentials
	if n > 20 {
		recent = differentials[n-20:]
	}
	// Work in whole tenths (Rule 5.1 Score Differentials are tenths), so the average and its rounding are exact.
	var buf [20]int64
	m := len(recent)
	for i, d := range recent {
		if math.IsNaN(d) || math.IsInf(d, 0) {
			return 0, fmt.Errorf("%w: differential %d = %v", ErrNotFinite, i, d)
		}
		if math.Abs(d) > 1000 {
			return 0, fmt.Errorf("%w: differential %d = %v, want |d| ≤ 1000", ErrOutOfRange, i, d)
		}
		buf[i] = int64(math.Round(d * 10))
	}
	// Insertion sort of at most 20 values.
	for i := 1; i < m; i++ {
		v := buf[i]
		j := i - 1
		for j >= 0 && buf[j] > v {
			buf[j+1] = buf[j]
			j--
		}
		buf[j+1] = v
	}
	row := whsTable52a[m]
	k := int64(row.lowest)
	sum := int64(0)
	for i := 0; i < row.lowest; i++ {
		sum += buf[i]
	}
	// Mean in tenths = sum / k, rounded half to even; then the adjustment (whole tenths).
	q, r := sum/k, sum%k
	if r < 0 {
		q, r = q-1, r+k
	}
	if 2*r > k || (2*r == k && q%2 != 0) {
		q++
	}
	hi := float64(q+int64(math.Round(row.adjustment*10))) / 10
	if hi > whsMaxHandicapIndex {
		hi = whsMaxHandicapIndex
	}
	return hi, nil
}

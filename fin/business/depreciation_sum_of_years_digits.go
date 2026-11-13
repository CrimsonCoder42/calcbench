package business

import "math/bits"

// DepreciationSumOfYearsDigits returns the per-period depreciation expense of an asset by the sum of the
// years-digits method, as a schedule of lifePeriods amounts in cents. Period k (1-based) gets
// (cost − salvage) × (n − k + 1) ÷ S with S = n(n+1)/2, computed in exact 128-bit integer arithmetic and
// rounded half-even to the cent; the last period takes the rounding remainder, so the schedule sums to
// cost − salvage exactly; an amount is cut to what is left if rounding would overrun it, so no amount is
// negative. The result is a new slice.
// Errors: ErrOutOfRange when a money input is outside ±2⁵³; ErrNegative when a money input is negative;
// ErrDomain when salvageCents > costCents; ErrPeriods when lifePeriods < 1 or > 1 000 000.
//
// ID:      FIN-BIZ-051
// Time:    O(n)   n = lifePeriods
// Space:   O(n) for the result (one allocation), O(1) extra
// Source:  26 CFR 1.167(b)-3(a) "Sum of the years-digits method" (numerator = remaining useful life including the current year, denominator = sum of the years digits of the useful life, applied to cost less salvage), Example (1), https://www.law.cornell.edu/cfr/text/26/1.167(b)-3 (read 2026-11-13)
// Tags:    depreciation, sum-of-years-digits, syd, accelerated-depreciation, schedule, fixed-assets, accounting
// Units:   costCents, salvageCents and every schedule amount in cents; lifePeriods a count of periods
// Example: cost $1,750, salvage $250, 5 years → $500, $400, $300, $200, $100 (26 CFR 1.167(b)-3(a), Example (1))
func DepreciationSumOfYearsDigits(costCents, salvageCents int64, lifePeriods int) ([]int64, error) {
	if err := checkDepreciationInputs(costCents, salvageCents, lifePeriods); err != nil {
		return nil, err
	}
	base := uint64(costCents - salvageCents) // ≤ 2⁵³
	n := uint64(lifePeriods)
	sum := n * (n + 1) / 2 // ≤ 5·10¹¹
	schedule := make([]int64, lifePeriods)
	var used uint64
	for k := range schedule[:lifePeriods-1] {
		hi, lo := bits.Mul64(base, n-uint64(k)) // < 2⁵³ · 10⁶, fits 128 bits; hi < sum so Div64 is safe
		q, r := bits.Div64(hi, lo, sum)
		if 2*r > sum || (2*r == sum && q%2 == 1) {
			q++
		}
		q = min(q, base-used) // never past cost − salvage, so the last amount is ≥ 0
		schedule[k] = int64(q)
		used += q
	}
	schedule[lifePeriods-1] = int64(base - used)
	return schedule, nil
}

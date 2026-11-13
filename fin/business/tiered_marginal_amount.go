package business

import (
	"fmt"
	"math"
)

// TieredMarginalAmount applies a marginal rate schedule to amountCents: tier k covers the part of the amount
// from tierStartsCents[k] up to tierStartsCents[k+1] (the last tier is open-ended) and is charged
// ratesPerTier[k] on that part only. It is the shape of a progressive tax schedule ("$5,535, plus 28 % of the
// excess over $36,900") and of a graduated commission. The caller supplies the schedule; none is built in.
// The per-tier products are summed in float64 and the total is rounded half-even to the cent once; it is
// within 1 cent of the exact value for totals ≤ 2⁴³ cents.
// Errors: ErrOutOfRange when amountCents or a tier start is outside ±2⁵³, or the two slices differ in length
// or are empty; ErrNegative when amountCents is negative; ErrDomain when tierStartsCents[0] ≠ 0 or the
// starts are not strictly increasing; ErrNotFinite when a rate is NaN or ±Inf; ErrRate when a rate is
// outside 0…1.
//
// ID:      FIN-BIZ-055
// Time:    O(t)   t = len(tierStartsCents)
// Space:   O(1), no allocation
// Source:  26 U.S.C. § 1(a) rate table ("Over $36,900 but not over $89,150: $5,535, plus 28% of the excess over $36,900"; each row's base is the full lower tiers at their own rates), https://www.law.cornell.edu/uscode/text/26/1 (read 2026-11-13)
// Tags:    tiered-rate, marginal-rate, progressive, bracket, graduated-commission, schedule
// Units:   amountCents, tierStartsCents and the result in cents; ratesPerTier decimals
// Example: § 1(a) schedule (0: 15 %, 36,900: 28 %, 89,150: 31 %), $100,000 → $20,165 + 31 % × $10,850 = $23,528.50 (hand computation from the table)
func TieredMarginalAmount(amountCents int64, tierStartsCents []int64, ratesPerTier []float64) (int64, error) {
	if err := checkNonNegativeCents("amountCents", amountCents); err != nil {
		return 0, err
	}
	if len(tierStartsCents) == 0 || len(tierStartsCents) != len(ratesPerTier) {
		return 0, fmt.Errorf("%w: %d tier starts and %d rates, want the same non-zero count", ErrOutOfRange, len(tierStartsCents), len(ratesPerTier))
	}
	if tierStartsCents[0] != 0 {
		return 0, fmt.Errorf("%w: tierStartsCents[0] = %d, want 0", ErrDomain, tierStartsCents[0])
	}
	for k, r := range ratesPerTier {
		if math.IsNaN(r) || math.IsInf(r, 0) {
			return 0, fmt.Errorf("%w: ratesPerTier[%d] = %v", ErrNotFinite, k, r)
		}
		if r < 0 || r > 1 {
			return 0, fmt.Errorf("%w: ratesPerTier[%d] = %v, want 0…1", ErrRate, k, r)
		}
		if k > 0 {
			if err := checkCents("tierStartsCents", tierStartsCents[k]); err != nil {
				return 0, err
			}
			if tierStartsCents[k] <= tierStartsCents[k-1] {
				return 0, fmt.Errorf("%w: tierStartsCents[%d] = %d is not above the previous start", ErrDomain, k, tierStartsCents[k])
			}
		}
	}
	var total float64
	for k, start := range tierStartsCents {
		if amountCents <= start {
			break
		}
		top := amountCents
		if k+1 < len(tierStartsCents) && tierStartsCents[k+1] < top {
			top = tierStartsCents[k+1]
		}
		total += ratesPerTier[k] * float64(top-start)
	}
	return int64(math.RoundToEven(total)), nil
}

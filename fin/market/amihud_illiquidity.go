package market

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// AmihudIlliquidity returns Amihud's illiquidity ratio ILLIQ = (1/D)·Σ_d |R_d| / VOLD_d, the average
// absolute return per unit of traded value over D days. Amihud averages over days with trading; a day
// with zero dollar volume is rejected, so the caller drops such days first. The source reports ILLIQ
// × 10⁶; any such scale is the caller's.
// Errors: ErrEmpty if returns is empty; ErrLengthMismatch if the lengths differ; ErrNotFinite if a value
// is NaN or ±Inf; ErrNegative if a dollar volume is negative; ErrDomain if a dollar volume is 0;
// ErrOverflow if the result exceeds float64.
//
// ID:      FIN-MKT-015
// Time:    O(D)   D = number of days
// Space:   O(1), no allocation
// Source:  Amihud, "Illiquidity and stock returns: cross-section and time-series effects", Journal of Financial Markets 5(1):31–56, 2002, doi:10.1016/S1386-4181(01)00024-6, §2, eq. (1) (ILLIQ_iy = (1/D_iy)·Σ_d |R_iyd|/VOLD_iyd)
// Tags:    market, microstructure, amihud, illiquidity, price-impact, liquidity
// Units:   returns as fractions (0.01 = 1 %); dollarVolumes in currency; result in fraction per currency unit
// Example: hand-checked: returns 0.02, −0.01, volumes 1e6, 5e5 → (2e-8 + 2e-8)/2 = 2e-8
func AmihudIlliquidity(returns, dollarVolumes []float64) (float64, error) {
	d := len(returns)
	if d == 0 {
		return 0, fmt.Errorf("%w: returns", ErrEmpty)
	}
	if len(dollarVolumes) != d {
		return 0, fmt.Errorf("%w: len(returns) = %d, len(dollarVolumes) = %d", ErrLengthMismatch, d, len(dollarVolumes))
	}
	var sum, comp float64
	for i := range d {
		r := returns[i]
		if !reduce.AllFinite(r) {
			return 0, fmt.Errorf("%w: returns[%d] = %v", ErrNotFinite, i, r)
		}
		if err := checkSize("dollarVolumes", dollarVolumes[i]); err != nil {
			return 0, fmt.Errorf("day %d: %w", i, err)
		}
		if dollarVolumes[i] == 0 {
			return 0, fmt.Errorf("%w: dollarVolumes[%d] = 0", ErrDomain, i)
		}
		if r < 0 {
			r = -r
		}
		sum, comp = reduce.NeumaierAdd(sum, comp, r/dollarVolumes[i])
	}
	v := (sum + comp) / float64(d)
	if !reduce.AllFinite(v) {
		return 0, fmt.Errorf("%w: ILLIQ", ErrOverflow)
	}
	return v, nil
}

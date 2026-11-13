package interest

import "math"

// EffectiveAnnualRateContinuous returns the effective annual rate of a nominal annual rate
// compounded continuously, e^r − 1, evaluated with math.Expm1 (no cancellation for small r). The
// result is a decimal (0.083287 for r = 0.08), unrounded.
// Errors: ErrRate if nominalRatePerYear is NaN, ±Inf or ≤ −1; ErrOverflow if the result exceeds float64.
//
// ID:      FIN-INT-030
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  R. A. DeFusco, D. W. McLeavey, J. E. Pinto and D. E. Runkle, Quantitative Investment Analysis, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ch. 1 "The Time Value of Money", §3.3, Eq. 6 EAR = e^(r_s) − 1
// Tags:    interest, rate-conversion, effective-annual-rate, ear, continuous-compounding
// Units:   nominalRatePerYear decimal per year (continuously compounded); result decimal per year
// Example: 8 % continuous → 0.0833 (QIA ch. 1 §3.3, Eq. 6)
func EffectiveAnnualRateContinuous(nominalRatePerYear float64) (float64, error) {
	if err := checkRate("nominalRatePerYear", nominalRatePerYear); err != nil {
		return 0, err
	}
	return factorResult(math.Expm1(nominalRatePerYear))
}

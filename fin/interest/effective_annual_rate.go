package interest

import "math"

// EffectiveAnnualRate returns the effective annual rate of a nominal (stated) annual rate compounded
// periodsPerYear times a year: (1 + r/m)^m − 1, evaluated as expm1(m·log1p(r/m)) so that small rates
// keep their precision (Goldberg 1991). The result is a decimal (0.0816 = 8.16 %), unrounded.
// Errors: ErrNotFinite if nominalRatePerYear is NaN or ±Inf; ErrPeriods if periodsPerYear < 1;
// ErrRate if nominalRatePerYear/periodsPerYear ≤ −1; ErrOverflow if the result exceeds float64.
//
// ID:      FIN-INT-029
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NCEES, FE Reference Handbook 10.0.1, NCEES, 2020, Engineering Economics, p. 231 "Non-Annual Compounding", i_e = (1 + r/m)^m − 1
// Source:  R. A. DeFusco, D. W. McLeavey, J. E. Pinto and D. E. Runkle, Quantitative Investment Analysis, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ch. 1 "The Time Value of Money", §3.3, Eq. 5 EAR = (1 + Periodic interest rate)^m − 1
// Tags:    interest, rate-conversion, effective-annual-rate, ear, nominal-rate, compounding-frequency
// Units:   nominalRatePerYear decimal per year; periodsPerYear compounding periods per year; result decimal per year
// Example: 8 % semiannual → 0.0816 (QIA ch. 1 §3.3); 18 % monthly → 0.19562 (Principles of Finance §8.4, Eq. 8.38)
func EffectiveAnnualRate(nominalRatePerYear float64, periodsPerYear int) (float64, error) {
	if err := checkFinite("nominalRatePerYear", nominalRatePerYear); err != nil {
		return 0, err
	}
	if err := checkPeriods("periodsPerYear", periodsPerYear, 1); err != nil {
		return 0, err
	}
	m := float64(periodsPerYear)
	if err := checkRate("nominalRatePerYear/periodsPerYear", nominalRatePerYear/m); err != nil {
		return 0, err
	}
	return factorResult(math.Expm1(m * math.Log1p(nominalRatePerYear/m)))
}

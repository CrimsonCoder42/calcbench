package interest

import "math"

// PeriodicRateFromEffectiveAnnual returns the rate per period that, compounded periodsPerYear times,
// gives effectiveRatePerYear: (1 + i_e)^(1/m) − 1, evaluated as expm1(log1p(i_e)/m) so that small rates
// keep their precision. It is the inverse of EffectiveAnnualRate (which returns i_e from the nominal
// rate r = m × the periodic rate). Use it to pay monthly into an account quoted at an effective annual
// rate. The result is a decimal per period, unrounded.
// Errors: ErrRate if effectiveRatePerYear is NaN, ±Inf or ≤ −1; ErrPeriods if periodsPerYear < 1.
//
// ID:      FIN-INT-044
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NCEES, FE Reference Handbook 10.0.1, NCEES, 2020, Engineering Economics, p. 231 "Non-Annual Compounding", i_e = (1 + r/m)^m − 1, solved for the periodic rate r/m
// Source:  R. A. DeFusco, D. W. McLeavey, J. E. Pinto and D. E. Runkle, Quantitative Investment Analysis, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ch. 1 "The Time Value of Money", §3.3, Eq. 5 EAR = (1 + Periodic interest rate)^m − 1
// Tags:    interest, rate-conversion, effective-annual-rate, periodic-rate, equivalent-rate, compounding-frequency
// Units:   effectiveRatePerYear decimal per year; periodsPerYear periods per year; result decimal per period
// Example: 8.16 % effective, semiannual → 0.04 per half-year (QIA ch. 1 §3.3 read backwards)
func PeriodicRateFromEffectiveAnnual(effectiveRatePerYear float64, periodsPerYear int) (float64, error) {
	if err := checkRate("effectiveRatePerYear", effectiveRatePerYear); err != nil {
		return 0, err
	}
	if err := checkPeriods("periodsPerYear", periodsPerYear, 1); err != nil {
		return 0, err
	}
	return math.Expm1(math.Log1p(effectiveRatePerYear) / float64(periodsPerYear)), nil
}

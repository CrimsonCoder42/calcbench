package stats

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/num/prob"
)

func checkIntervalInputs(mean, sd float64, n int) error {
	if err := checkScalar("mean", mean); err != nil {
		return err
	}
	if err := checkScalar("sd", sd); err != nil {
		return err
	}
	if sd < 0 {
		return fmt.Errorf("%w: sd = %v", ErrNegative, sd)
	}
	if n < 1 {
		return fmt.Errorf("%w: n = %d, want ≥ 1", ErrOutOfRange, n)
	}
	return nil
}

// MeanConfidenceIntervalZ returns the two-sided confidence interval for a population mean when the
// standard deviation σ is known (or n is large enough to treat s as σ): x̄ ± z·σ/√n, z the normal
// quantile at 1/2 + confidence/2. It returns the lower and upper limits and the margin z·σ/√n.
// Errors: ErrNotFinite for NaN or ±Inf; ErrNegative when sd < 0; ErrOutOfRange unless n ≥ 1 and
// 0 < confidence < 1.
//
// ID:      NUM-STATS-050
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Montgomery and Runger, Applied Statistics and Probability for Engineers, 6th ed., Wiley, 2014, §8.1 "Confidence Interval on the Mean of a Normal Distribution, Variance Known" (x̄ − z_{α/2}σ/√n ≤ μ ≤ x̄ + z_{α/2}σ/√n)
// Tags:    confidence-interval, mean, z-interval, margin-of-error, inference
// Units:   mean, sd and the limits in one unit; n a count; confidence a fraction (0.95)
// Example: x̄ 20.6, σ 3.2, n 50, 95 % → 20.6 ± 1.959964·3.2/√50 = 20.6 ± 0.886981 (hand-checked)
func MeanConfidenceIntervalZ(mean, sd float64, n int, confidence float64) (lo, hi, margin float64, err error) {
	if err := checkIntervalInputs(mean, sd, n); err != nil {
		return 0, 0, 0, err
	}
	z, err := twoSidedZ(confidence)
	if err != nil {
		return 0, 0, 0, err
	}
	margin = z * (sd / math.Sqrt(float64(n)))
	lo, hi = mean-margin, mean+margin
	if math.IsInf(lo, 0) || math.IsInf(hi, 0) {
		return 0, 0, 0, fmt.Errorf("%w: interval limits", ErrOverflow)
	}
	return lo, hi, margin, nil
}

// MeanConfidenceIntervalT returns the two-sided confidence interval for a normal population mean from a
// sample of n with sample standard deviation s: x̄ ± t·s/√n, t the Student t quantile at
// 1/2 + confidence/2 with n − 1 degrees of freedom. It returns the lower and upper limits and the margin.
// Errors: ErrNotFinite for NaN or ±Inf; ErrNegative when sd < 0; ErrOutOfRange unless n ≥ 2 and
// 0 < confidence < 1.
//
// ID:      NUM-STATS-051
// Time:    O(B)   B = the bisection steps of NUM-PROB-079 StudentTQuantile (bounded there)
// Space:   O(1), no allocation
// Source:  NIST/SEMATECH, e-Handbook of Statistical Methods, §7.2.2.1 "Confidence limits for the mean" (Ȳ ± t_{1−α/2, N−1}·s/√N), https://www.itl.nist.gov/div898/handbook/prc/section2/prc221.htm
// Source:  Montgomery and Runger, Applied Statistics and Probability for Engineers, 6th ed., Wiley, 2014, §8.2 "Confidence Interval on the Mean of a Normal Distribution, Variance Unknown"
// Tags:    confidence-interval, mean, t-interval, student-t, margin-of-error, inference
// Units:   mean, sd and the limits in one unit; n a count; confidence a fraction (0.95)
// Example: x̄ 20.6, s 3.2, n 50, 95 % → t₀.₉₇₅,₄₉ = 2.009575, margin 0.909433 (hand-checked)
func MeanConfidenceIntervalT(mean, sd float64, n int, confidence float64) (lo, hi, margin float64, err error) {
	if err := checkIntervalInputs(mean, sd, n); err != nil {
		return 0, 0, 0, err
	}
	if n < 2 {
		return 0, 0, 0, fmt.Errorf("%w: n = %d, want ≥ 2", ErrOutOfRange, n)
	}
	if _, err := twoSidedZ(confidence); err != nil {
		return 0, 0, 0, err
	}
	t, err := prob.StudentTQuantile(0.5+confidence/2, float64(n-1))
	if err != nil {
		return 0, 0, 0, err
	}
	margin = t * (sd / math.Sqrt(float64(n)))
	lo, hi = mean-margin, mean+margin
	if math.IsInf(lo, 0) || math.IsInf(hi, 0) || math.IsInf(margin, 0) {
		return 0, 0, 0, fmt.Errorf("%w: interval limits", ErrOverflow)
	}
	return lo, hi, margin, nil
}

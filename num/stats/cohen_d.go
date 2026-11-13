package stats

import "fmt"

// CohenD returns Cohen's d, the standardized difference between two means, d = (meanA − meanB)/sd, with
// sd the standard deviation common to both populations (Cohen's σ). For two samples the usual estimate
// of sd is PooledStandardDeviation. Its sign follows meanA − meanB. It states the number only; labels for
// ranges of d are interpretation and are not given here.
// Errors: ErrNotFinite when meanA, meanB or sd is NaN or ±Inf; ErrNegative when sd < 0; ErrZeroVariance
// when sd = 0; ErrOverflow when d exceeds the float64 range.
//
// ID:      NUM-STATS-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Cohen, Statistical Power Analysis for the Behavioral Sciences, 2nd ed., Lawrence Erlbaum, 1988, §2.2, eq. 2.2.1, p. 20
// Source:  Lakens, "Calculating and reporting effect sizes to facilitate cumulative science: a practical primer for t-tests and ANOVAs", Frontiers in Psychology, 2013, 4:863, doi:10.3389/fpsyg.2013.00863 (eq. 1, d_s with the pooled standard deviation)
// Source:  Goulet-Pelletier, Cousineau, "A review of effect sizes and their confidence intervals, Part I: The Cohen's d family", The Quantitative Methods for Psychology, 2018, 14(4):242–265, doi:10.20982/tqmp.14.4.p242 (eq. 1)
// Tags:    cohen-d, effect-size, standardized-mean-difference, two-sample, descriptive
// Units:   meanA, meanB and sd in one unit; result dimensionless
// Example: meanA = 30.48101, meanB = 20.14458, sd = 6.34260 (e-Handbook §1.3.5.3 AUTO83B means and pooled s) → 1.62968 (hand-checked)
func CohenD(meanA, meanB, sd float64) (float64, error) {
	if err := checkScalar("meanA", meanA); err != nil {
		return 0, err
	}
	if err := checkScalar("meanB", meanB); err != nil {
		return 0, err
	}
	if err := checkScalar("sd", sd); err != nil {
		return 0, err
	}
	if sd < 0 {
		return 0, fmt.Errorf("%w: sd = %v", ErrNegative, sd)
	}
	if sd == 0 {
		return 0, fmt.Errorf("%w: sd = 0", ErrZeroVariance)
	}
	return standardizedDifference(meanA, meanB, sd, "d")
}

package prob

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/num/special"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// StudentTPDF returns the density of Student's t distribution with ν = df degrees of freedom,
// f(t) = (1 + t²/ν)^(−(ν+1)/2)/(√ν·B(ν/2, ½)), evaluated in log space with log1p and num/special LogBeta.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when df ≤ 0.
//
// ID:      NUM-PROB-077
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 28 "t-Distributions", §2 (eq. 28.2)
// Source:  Abramowitz, Stegun, Handbook of Mathematical Functions, NBS, 1964, §26.7.1
// Tags:    probability, distribution, student-t, t-distribution, pdf, continuous
// Units:   dimensionless
// Example: ν = 1, t = 0 → 1/π ≈ 0.3183098861837907 (the Cauchy case); ν = 2, t = 0 → 1/(2√2) ≈ 0.3535533905932738
func StudentTPDF(t, df float64) (float64, error) {
	if err := checkPositive("df", df); err != nil {
		return 0, err
	}
	if err := checkFinite("t", t); err != nil {
		return 0, err
	}
	lb, err := special.LogBeta(df/2, 0.5)
	if err != nil {
		return 0, specialErr("log beta", err)
	}
	d := math.Exp(-(df+1)/2*math.Log1p(t*t/df) - 0.5*math.Log(df) - lb)
	if !reduce.AllFinite(d) {
		return 0, fmt.Errorf("%w: density at t = %v, df = %v", ErrOverflow, t, df)
	}
	return d, nil
}

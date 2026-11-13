package prob

// studentTTail returns P(T ≤ −|t|) and its complement P(T > −|t|) for validated df. With r = ν/t²,
// x = ν/(ν + t²) = r/(1 + r) and 1 − x = 1/(1 + r) are both formed without cancellation, and
// P(T ≤ −|t|) = ½·I_x(ν/2, ½) (A&S 26.7.1 with 26.5.27).
func studentTTail(t, df float64) (tail, rest float64, err error) {
	if t == 0 {
		return 0.5, 0.5, nil
	}
	r := df / (t * t) // t² may overflow to +Inf: r = 0, the tail is 0
	i, _, err := betaSplit(r/(1+r), 1/(1+r), df/2, 0.5)
	if err != nil {
		return 0, 0, err
	}
	tail = i / 2
	return tail, 1 - tail, nil
}

// StudentTCDF returns P(T ≤ t) for Student's t distribution with ν = df degrees of freedom:
// ½·I_(ν/(ν+t²))(ν/2, ½) for t < 0 and one minus that for t > 0, from the regularized incomplete beta
// (num/special NUM-SPEC-005). The lower tail keeps its relative accuracy far into the tail.
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrOutOfRange when df ≤ 0, or when the incomplete
// beta does not converge (wrapping special.ErrNoConvergence).
//
// ID:      NUM-PROB-078
// Time:    O(√ν + log(1/ε))   ν = df, ε = 2⁻⁵²
// Space:   O(1), no allocation
// Source:  Johnson, Kotz, Balakrishnan, Continuous Univariate Distributions, vol. 2, 2nd ed., Wiley, 1995, ch. 28 "t-Distributions", §2
// Source:  Abramowitz, Stegun, Handbook of Mathematical Functions, NBS, 1964, §26.7.1 and §26.5.27 (A(t|ν) via I_x)
// Tags:    probability, distribution, student-t, t-distribution, cdf, p-value, t-test
// Units:   dimensionless; result a probability in [0, 1]
// Example: ν = 1, t = 1 → 0.75 (the Cauchy case, ½ + atan(1)/π); ν = 2, t = 2 → ½ + 1/√6 ≈ 0.9082482904638631
func StudentTCDF(t, df float64) (float64, error) {
	if err := checkPositive("df", df); err != nil {
		return 0, err
	}
	if err := checkFinite("t", t); err != nil {
		return 0, err
	}
	tail, rest, err := studentTTail(t, df)
	if err != nil {
		return 0, err
	}
	if t < 0 {
		return tail, nil
	}
	return rest, nil
}

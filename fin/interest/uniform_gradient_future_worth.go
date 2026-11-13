package interest

import "math"

// expm1MinusXOverX2 returns (eˣ − 1 − x)/x², with its limit ½ at x = 0. For |x| < ½ it sums the
// Taylor series Σ xᵏ/(k+2)! (no cancellation); otherwise it evaluates the expression directly, where
// the subtraction loses at most about 3 bits.
func expm1MinusXOverX2(x float64) float64 {
	if math.Abs(x) < 0.5 {
		// Terms xᵏ/(k+2)!; 20 terms bring |x|²⁰/22! below 1e-27 for |x| < ½.
		sum, term := 0.0, 0.5
		for k := 0; k < 20; k++ {
			sum += term
			term *= x / float64(k+3)
		}
		return sum
	}
	return (math.Expm1(x) - x) / (x * x)
}

// log1pMinusXOverX2 returns (ln(1+x) − x)/x², with its limit −½ at x = 0. For |x| < 0.1 it sums
// Σ_{k≥2} (−1)^(k+1) x^(k−2)/k (no cancellation); otherwise it evaluates the expression directly,
// where the subtraction loses at most about 5 bits.
func log1pMinusXOverX2(x float64) float64 {
	if math.Abs(x) < 0.1 {
		// 18 terms bring 0.1¹⁸/20 below 1e-19.
		sum, pow := 0.0, 1.0
		sign := -1.0
		for k := 2; k < 20; k++ {
			sum += sign * pow / float64(k)
			pow *= x
			sign = -sign
		}
		return sum
	}
	return (math.Log1p(x) - x) / (x * x)
}

// uniformGradientFutureWorth returns ((1+i)ⁿ − 1 − n·i)/i² for checked inputs and n ≥ 2. With
// L = ln(1+i) and x = n·L the numerator is (eˣ − 1 − x) + n·(L − i), so the value is
// (eˣ−1−x)/x² · (n·L/i)² + n·(L − i)/i², two terms free of cancellation that agree in sign for
// n ≥ 2 only up to a factor n/(n−1). At i = 0 it is n(n−1)/2.
func uniformGradientFutureWorth(i float64, n int) float64 {
	nf := float64(n)
	if i == 0 {
		return nf * (nf - 1) / 2
	}
	x := nf * math.Log1p(i)
	r := nf * log1pOverX(i)
	return expm1MinusXOverX2(x)*r*r + nf*log1pMinusXOverX2(i)
}

// UniformGradientFutureWorth returns the factor (F/G,i,n) = ((1+i)ⁿ − 1)/i² − n/i, the value at the
// end of period n of the gradient series 0, G, 2G, …, (n−1)G paid at the ends of periods 1…n, per
// unit of G. It equals ((F/A,i,n) − n)/i, evaluated without the cancellation that form has for small
// i; i = 0 gives exactly n(n−1)/2, and periods ≤ 1 gives 0.
// Errors: ErrRate if ratePerPeriod is NaN, ±Inf or ≤ −1; ErrPeriods if periods < 0; ErrOverflow if
// the factor exceeds float64.
//
// ID:      FIN-INT-012
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NCEES, FE Reference Handbook 10.0.1, NCEES, 2020, Engineering Economics, factor table p. 230, Uniform Gradient Future Worth (F/G, i%, n) = [(1 + i)^n − 1]/i² − n/i, footnote F/G = (F/A − n)/i
// Tags:    interest, engineering-economy, factor, gradient, future-worth, arithmetic-gradient
// Units:   ratePerPeriod decimal per period; periods count; result dimensionless (per unit of G)
// Example: i = 1 %, n = 10 → 46.2213 = (10.4622 − 10)/0.01 by the handbook's footnote (hand-checked from p. 233)
func UniformGradientFutureWorth(ratePerPeriod float64, periods int) (float64, error) {
	if err := checkRate("ratePerPeriod", ratePerPeriod); err != nil {
		return 0, err
	}
	if err := checkPeriods("periods", periods, 0); err != nil {
		return 0, err
	}
	if periods <= 1 {
		return 0, nil
	}
	return factorResult(uniformGradientFutureWorth(ratePerPeriod, periods))
}

package business

import (
	"fmt"
	"math"
)

// maxConvertible is 2⁶²: a float64 below it in magnitude truncates to an int64 without overflow, and the
// int64 sums below cannot overflow either.
const maxConvertible = 1 << 62

// roundHalfEvenPlus returns x + k rounded half-even to an integer, exactly (x is taken as the exact value
// it holds). |x| must be < 2⁶² and |k| ≤ 2⁶¹.
func roundHalfEvenPlus(x float64, k int64) int64 {
	return roundHalfEvenExact(x, 0, k)
}

// roundHalfEvenExact returns the exact real p + e + k rounded half-even to an integer, where p is a
// float64 product fl(a·b) and e = math.FMA(a, b, −p) its exact rounding error (so p + e = a·b exactly), or
// e = 0. |p| must be < 2⁶², |e| ≤ ½·ulp(p), and |k| ≤ 2⁶¹.
//
// Why it is exact: p and e are split into integer parts (Trunc, exact) and fractions fp, fe (exact
// subtractions); the integer parts are summed with k in int64 as t. If |p| ≥ 2⁵², p is an integer, so
// fp = 0. Otherwise |e| ≤ ¼, so fe = e. The exact fraction f = fp + fe has |f| < 1.5 and is compared with
// ±½ as fe against ±½ − fp. When |fp| ≥ ¼ that difference is exact (Sterbenz). When |fp| < ¼ the
// computed ½ − fp is still ≥ ¼, while |fe| ≤ ⅛ unless |p| ≥ 2⁵¹ — and then p is a multiple of ½, so
// fp ∈ {0, ±½} and the difference is exact again. Ties go to the even total.
// Method: Goldberg, "What Every Computer Scientist Should Know About Floating-Point Arithmetic", ACM
// Computing Surveys 23(1):5–48, 1991, doi:10.1145/103162.103163, Theorem 11 (Sterbenz); FMA gives the
// product error exactly (it is representable).
func roundHalfEvenExact(p, e float64, k int64) int64 {
	wp, we := math.Trunc(p), math.Trunc(e)
	fp, fe := p-wp, e-we
	t := int64(wp) + int64(we) + k
	up, down := 0.5-fp, -0.5-fp
	switch {
	case fe > up:
		return t + 1
	case fe < down:
		return t - 1
	case fe == up:
		if t%2 != 0 {
			return t + 1
		}
	case fe == down:
		if t%2 != 0 {
			return t - 1
		}
	}
	return t
}

// EconomicValueAdded returns after-tax operating income minus the capital charge, nopat −
// investedCapital × costOfCapital, in cents, rounded half-even once on the total. The result is the
// correctly rounded exact value (costOfCapital taken as the float64 it holds): the product's rounding
// error is recovered with math.FMA and the half-cent decision is made on the exact sum.
// Errors: ErrOutOfRange when a money input is outside ±2⁵³; ErrNegative when investedCapitalCents is
// negative; ErrNotFinite when costOfCapital is NaN or ±Inf; ErrRate when costOfCapital ≤ −1; ErrOverflow
// when the result is outside ±2⁵³ cents.
//
// ID:      FIN-BIZ-034
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Franklin, Graybeal and Cooper, Principles of Accounting, Vol. 2: Managerial Accounting, OpenStax (Rice University), 2019, §12.3 "Evaluate an Operating Segment or a Project Using Return on Investment, Residual Income, and Economic Value Added" (EVA = after-tax income − invested capital × weighted average cost of capital), https://openstax.org/books/principles-managerial-accounting/pages/12-3-evaluate-an-operating-segment-or-a-project-using-return-on-investment-residual-income-and-economic-value-added (read 2026-11-12)
// Tags:    value, economic-value-added, eva, residual-income, cost-of-capital
// Units:   nopatCents, investedCapitalCents and the result in cents; costOfCapital a decimal per period (0.09 = 9 %)
// Example: new machinery: $150,000 − $1,500,000 × 0.09 = $15,000 (Principles of Accounting Vol. 2 §12.3)
func EconomicValueAdded(nopatCents, investedCapitalCents int64, costOfCapital float64) (int64, error) {
	if err := checkCents("nopatCents", nopatCents); err != nil {
		return 0, err
	}
	if err := checkNonNegativeCents("investedCapitalCents", investedCapitalCents); err != nil {
		return 0, err
	}
	if math.IsNaN(costOfCapital) || math.IsInf(costOfCapital, 0) {
		return 0, fmt.Errorf("%w: costOfCapital = %v", ErrNotFinite, costOfCapital)
	}
	if costOfCapital <= -1 {
		return 0, fmt.Errorf("%w: costOfCapital = %v, want > −1", ErrRate, costOfCapital)
	}
	charge := float64(float64(investedCapitalCents) * costOfCapital) // explicit conversion: rounded, never fused (Go spec, Floating-point operators)
	if math.Abs(charge) >= maxConvertible {
		return 0, fmt.Errorf("%w: investedCapitalCents × costOfCapital", ErrOverflow)
	}
	chargeErr := math.FMA(float64(investedCapitalCents), costOfCapital, -charge)
	eva := roundHalfEvenExact(-charge, -chargeErr, nopatCents)
	if err := checkResultCents("economic value added", eva); err != nil {
		return 0, err
	}
	return eva, nil
}

package interest

import (
	"errors"
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// The package's sentinel errors. They live in the file of the first function that uses them
// (doc/ASSUMPTIONS.md A-008). Every error a function returns wraps one of these with the name of the
// offending parameter; test with errors.Is.
var (
	// ErrNotFinite: a float64 input is NaN or ±Inf. A non-finite rate also wraps ErrRate.
	ErrNotFinite = errors.New("interest: value is NaN or infinite")
	// ErrOutOfRange: a money input is outside ±2⁵³ cents, or another input is outside its stated range.
	ErrOutOfRange = errors.New("interest: input out of range")
	// ErrOverflow: a money result is outside ±2⁵³ cents, or a float64 result is beyond float64.
	ErrOverflow = errors.New("interest: result overflows")
	// ErrRate: a rate is NaN, ±Inf, or not greater than −1.
	ErrRate = errors.New("interest: rate must be finite and greater than -1")
	// ErrPeriods: a period, payment or day count is below its minimum.
	ErrPeriods = errors.New("interest: period count out of range")
	// ErrNegative: an input that must not be negative is negative.
	ErrNegative = errors.New("interest: negative input")
	// ErrDomain: the inputs are valid one by one but the formula is undefined for them together.
	ErrDomain = errors.New("interest: inputs outside the formula's domain")
	// ErrNoConvergence: the root finder did not meet its tolerance within its iteration limit.
	ErrNoConvergence = errors.New("interest: root finder did not converge")
	// ErrNoSolution: no rate solves the equation (no sign change was found).
	ErrNoSolution = errors.New("interest: no solution")
)

// maxExactCents is the largest money magnitude, in cents, that float64 holds exactly (2⁵³).
const maxExactCents = 1 << 53

// checkCents returns ErrOutOfRange when |cents| > 2⁵³.
func checkCents(name string, cents int64) error {
	if cents > maxExactCents || cents < -maxExactCents {
		return fmt.Errorf("%w: %s = %d is outside ±2^53 cents", ErrOutOfRange, name, cents)
	}
	return nil
}

// checkRate returns ErrRate (also wrapping ErrNotFinite for NaN and ±Inf) unless rate is finite and > −1.
func checkRate(name string, rate float64) error {
	if !reduce.AllFinite(rate) {
		return fmt.Errorf("%w: %w: %s = %v", ErrRate, ErrNotFinite, name, rate)
	}
	if rate <= -1 {
		return fmt.Errorf("%w: %s = %v", ErrRate, name, rate)
	}
	return nil
}

// checkPeriods returns ErrPeriods when n < minimum.
func checkPeriods(name string, n, minimum int) error {
	if n < minimum {
		return fmt.Errorf("%w: %s = %d, want ≥ %d", ErrPeriods, name, n, minimum)
	}
	return nil
}

// checkFinite returns ErrNotFinite for NaN or ±Inf.
func checkFinite(name string, v float64) error {
	if !reduce.AllFinite(v) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, v)
	}
	return nil
}

// roundCents rounds a float64 cent amount half-even (math.RoundToEven) to int64 cents.
// A NaN, an infinity or a magnitude above 2⁵³ after rounding is ErrOverflow.
func roundCents(v float64) (int64, error) {
	r := math.RoundToEven(v)
	if math.IsNaN(r) || r > maxExactCents || r < -maxExactCents {
		return 0, fmt.Errorf("%w: result %v cents is outside ±2^53", ErrOverflow, v)
	}
	return int64(r), nil
}

// compound returns (1+i)ⁿ as exp(n·log1p(i)), the form that keeps the precision of small i
// (Goldberg 1991, §1.4, Theorem 4). It may return +Inf or 0 at the ends of float64.
func compound(i float64, n int) float64 {
	return math.Exp(float64(n) * math.Log1p(i))
}

// FutureValue returns the value after periods compounding periods of presentValueCents at
// ratePerPeriod, PV·(1+i)ⁿ, rounded half-even to the cent. The rate is applied once per period and
// interest is left to compound; a negative amount gives the mirror-image result. periods = 0 returns
// the amount unchanged.
// The result is the float64 value rounded half-even; it equals the correctly rounded exact value
// except when that value lies within about n·ε·|x| of a half cent, and is within 1 cent of it for
// |result| ≤ 2⁴³ cents.
// Errors: ErrOutOfRange if |presentValueCents| > 2⁵³; ErrRate if ratePerPeriod is NaN, ±Inf or ≤ −1;
// ErrPeriods if periods < 0; ErrOverflow if the result exceeds ±2⁵³ cents.
//
// ID:      FIN-INT-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  J. Dahlquist and R. Knight, Principles of Finance, OpenStax, 2022, §7.2 "Time Value of Money (TVM) Basics", FV = PV × (1 + r)^n, https://openstax.org/books/principles-finance/pages/7-2-time-value-of-money-tvm-basics
// Source:  R. A. DeFusco, D. W. McLeavey, J. E. Pinto and D. E. Runkle, Quantitative Investment Analysis, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ch. 1 "The Time Value of Money", §3, Eq. 2
// Tags:    interest, time-value, future-value, compound, single-payment
// Units:   presentValueCents and the result in cents; ratePerPeriod decimal per period; periods count
// Example: PV $1,000, 3 %, 4 periods → $1,125.51 (Principles of Finance §7.3, Example 1)
func FutureValue(presentValueCents int64, ratePerPeriod float64, periods int) (int64, error) {
	if err := checkCents("presentValueCents", presentValueCents); err != nil {
		return 0, err
	}
	if err := checkRate("ratePerPeriod", ratePerPeriod); err != nil {
		return 0, err
	}
	if err := checkPeriods("periods", periods, 0); err != nil {
		return 0, err
	}
	if presentValueCents == 0 || periods == 0 {
		return presentValueCents, nil
	}
	return roundCents(float64(presentValueCents) * compound(ratePerPeriod, periods))
}

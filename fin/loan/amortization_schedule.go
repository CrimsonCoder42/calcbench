package loan

import (
	"errors"
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
	"github.com/CrimsonCoder42/calcbench/fin/interest"
)

// The package's sentinel errors. They live in the file of the first function that uses them
// (doc/ASSUMPTIONS.md A-008). Every error a function returns wraps one of these with the name of the
// offending parameter; test with errors.Is. An error that comes from fin/interest wraps both the
// matching sentinel here and the fin/interest one.
var (
	// ErrNotFinite: a float64 input is NaN or ±Inf. A non-finite rate also wraps ErrRate.
	ErrNotFinite = errors.New("loan: value is NaN or infinite")
	// ErrOutOfRange: a money input is outside ±2⁵³ cents, or another input is outside its stated range.
	ErrOutOfRange = errors.New("loan: input out of range")
	// ErrOverflow: a money result or intermediate balance is outside ±2⁵³ cents.
	ErrOverflow = errors.New("loan: result overflows")
	// ErrRate: a rate is NaN, ±Inf, or not greater than −1.
	ErrRate = errors.New("loan: rate must be finite and greater than -1")
	// ErrPeriods: a period or payment count is outside its range.
	ErrPeriods = errors.New("loan: period count out of range")
	// ErrNegative: an input that must not be negative is negative.
	ErrNegative = errors.New("loan: negative input")
	// ErrZeroDenominator: the input used as the denominator is zero.
	ErrZeroDenominator = errors.New("loan: zero denominator")
	// ErrDomain: the inputs are valid one by one but the formula is undefined for them together.
	ErrDomain = errors.New("loan: inputs outside the formula's domain")
)

// maxExactCents is the largest money magnitude, in cents, that float64 holds exactly (2⁵³).
const maxExactCents = 1 << 53

// maxSchedulePeriods bounds the length of a schedule so that a bad input cannot request a
// multi-gigabyte allocation (10⁷ rows of 40 bytes is 400 MB).
const maxSchedulePeriods = 10_000_000

// checkCents returns ErrOutOfRange when |cents| > 2⁵³.
func checkCents(name string, cents int64) error {
	if cents > maxExactCents || cents < -maxExactCents {
		return fmt.Errorf("%w: %s = %d is outside ±2^53 cents", ErrOutOfRange, name, cents)
	}
	return nil
}

// checkNonNegativeCents returns ErrOutOfRange when |cents| > 2⁵³ and ErrNegative when cents < 0.
func checkNonNegativeCents(name string, cents int64) error {
	if err := checkCents(name, cents); err != nil {
		return err
	}
	if cents < 0 {
		return fmt.Errorf("%w: %s = %d", ErrNegative, name, cents)
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

// checkPeriods returns ErrPeriods unless minimum ≤ n ≤ maximum.
func checkPeriods(name string, n, minimum, maximum int) error {
	if n < minimum || n > maximum {
		return fmt.Errorf("%w: %s = %d, want %d…%d", ErrPeriods, name, n, minimum, maximum)
	}
	return nil
}

// roundCents rounds a float64 cent amount half-even to int64 cents; NaN, ±Inf or a magnitude above 2⁵³
// is ErrOverflow.
func roundCents(name string, v float64) (int64, error) {
	r := math.RoundToEven(v)
	if math.IsNaN(r) || r > maxExactCents || r < -maxExactCents {
		return 0, fmt.Errorf("%w: %s = %v cents is outside ±2^53", ErrOverflow, name, v)
	}
	return int64(r), nil
}

// checkResultCents returns ErrOverflow when an integer money result leaves ±2⁵³.
func checkResultCents(name string, cents int64) error {
	if cents > maxExactCents || cents < -maxExactCents {
		return fmt.Errorf("%w: %s = %d cents is outside ±2^53", ErrOverflow, name, cents)
	}
	return nil
}

// interestSentinels pairs each fin/interest sentinel with this package's.
var interestSentinels = [...][2]error{
	{interest.ErrOverflow, ErrOverflow},
	{interest.ErrRate, ErrRate},
	{interest.ErrNotFinite, ErrNotFinite},
	{interest.ErrOutOfRange, ErrOutOfRange},
	{interest.ErrPeriods, ErrPeriods},
	{interest.ErrNegative, ErrNegative},
	{interest.ErrDomain, ErrDomain},
}

// fromInterest wraps an error from fin/interest with the matching sentinel of this package, so callers
// can test with loan.ErrX (and still with interest.ErrX).
func fromInterest(err error) error {
	for _, p := range interestSentinels {
		if errors.Is(err, p[0]) {
			return fmt.Errorf("%w: %w", p[1], err)
		}
	}
	return err
}

// periodInterest returns balance·rate rounded half-even to the cent.
func periodInterest(balanceCents int64, ratePerPeriod float64) (int64, error) {
	return roundCents("interest", float64(balanceCents)*ratePerPeriod)
}

// AmortizationRow is one period of an amortization schedule. BalanceCents is the balance after the
// period's payment.
type AmortizationRow struct {
	Period         int   // 1-based period number
	PaymentCents   int64 // amount paid this period = InterestCents + PrincipalCents
	InterestCents  int64 // balance before the payment × ratePerPeriod, rounded half-even
	PrincipalCents int64 // part of the payment that reduces the balance
	BalanceCents   int64 // balance after the payment
}

// amortizeRow computes one schedule row from the balance before it. The row pays the loan off (pays
// the whole balance plus its interest, leaving exactly 0) when it is the last period or when the
// balance plus interest does not exceed the level payment; payoff reports which.
func amortizeRow(balanceCents, paymentCents int64, ratePerPeriod float64, last bool) (row AmortizationRow, payoff bool, err error) {
	in, err := periodInterest(balanceCents, ratePerPeriod)
	if err != nil {
		return AmortizationRow{}, false, err
	}
	row = AmortizationRow{PaymentCents: paymentCents, InterestCents: in}
	payoff = last || balanceCents+in <= paymentCents
	if payoff {
		row.PrincipalCents = balanceCents
		row.PaymentCents = in + balanceCents
	} else {
		row.PrincipalCents = paymentCents - in
	}
	row.BalanceCents = balanceCents - row.PrincipalCents
	if err := checkResultCents("PaymentCents", row.PaymentCents); err != nil {
		return AmortizationRow{}, false, err
	}
	if err := checkResultCents("BalanceCents", row.BalanceCents); err != nil {
		return AmortizationRow{}, false, err
	}
	return row, payoff, nil
}

// loanPayment validates a fully amortizing loan's inputs and returns its level payment.
func loanPayment(principalCents int64, ratePerPeriod float64, periods, maxPeriods int) (int64, error) {
	if err := checkNonNegativeCents("principalCents", principalCents); err != nil {
		return 0, err
	}
	if err := checkRate("ratePerPeriod", ratePerPeriod); err != nil {
		return 0, err
	}
	if err := checkPeriods("periods", periods, 1, maxPeriods); err != nil {
		return 0, err
	}
	p, err := interest.AnnuityPaymentFromPresentValue(principalCents, ratePerPeriod, periods)
	if err != nil {
		return 0, fromInterest(err)
	}
	return p, nil
}

// AmortizationSchedule returns the amortization schedule of a fully amortizing loan of principalCents
// repaid in periods level end-of-period payments at ratePerPeriod. The level payment is
// interest.AnnuityPaymentFromPresentValue (P·i/(1 − (1+i)⁻ⁿ), rounded half-even). In each row the
// interest is the opening balance × ratePerPeriod rounded half-even and the principal is the payment
// less the interest; this per-row half-even rounding is the definition of every row (not an
// approximation of a closed form). The payoff row pays the remaining balance plus its interest, so the
// principal column sums to principalCents, no balance is negative and the final balance is exactly 0.
// The payoff row is row periods, or an earlier row when the rounded level payment would repay the loan
// sooner (a tiny loan whose payment rounds up: AmortizationSchedule(100, 0.005, 360) pays 1 cent a
// period and returns 100 rows); then fewer than periods rows are returned. For ratePerPeriod ≥ 0 the
// payoff payment differs from the level payment by at most (F/A,i,n) cents: the payment's rounding
// (≤ ½ cent) and each row's interest rounding (≤ ½ cent) compound to the end (about 1,005 cents at
// i = 0.005, n = 360; typically about half that). principalCents = 0 gives an empty, non-nil slice. The
// slice is allocated once, with capacity periods.
// Errors: ErrOutOfRange if |principalCents| > 2⁵³; ErrNegative if principalCents < 0; ErrRate if
// ratePerPeriod is NaN, ±Inf or ≤ −1; ErrPeriods if periods < 1 or > 10⁷; ErrOverflow if a payment or
// interest amount exceeds ±2⁵³ cents.
//
// ID:      FIN-LOAN-001
// Time:    O(n)   n = periods
// Space:   O(n) for the result (one allocation); O(1) extra
// Source:  Dahlquist & Knight, Principles of Finance, OpenStax (Rice University), 2022, §8.3 "Loan Amortization", Table 8.4 and the 20-year mortgage table, https://openstax.org/books/principles-finance/pages/8-3-loan-amortization (read 2026-11-12)
// Source:  Dahlquist & Knight, Principles of Finance, OpenStax, 2022, §8.2 "Annuities", eq. 8.14 (the level payment), https://openstax.org/books/principles-finance/pages/8-2-annuities (read 2026-11-12)
// Tags:    loan, amortization, schedule, mortgage, payment, interest, principal, balance
// Units:   principalCents and every row amount in cents; ratePerPeriod decimal per payment period; periods count
// Example: OpenStax §8.3 Table 8.4: $32,000, 0.5 %/month, 36 → payment 973.50; rows 1–3 interest 160.00, 155.93, 151.84, balances 31,186.50, 30,368.93, 29,547.27
func AmortizationSchedule(principalCents int64, ratePerPeriod float64, periods int) ([]AmortizationRow, error) {
	payment, err := loanPayment(principalCents, ratePerPeriod, periods, maxSchedulePeriods)
	if err != nil {
		return nil, err
	}
	if principalCents == 0 {
		return []AmortizationRow{}, nil
	}
	rows := make([]AmortizationRow, periods)
	balance := principalCents
	for k := range rows {
		row, payoff, err := amortizeRow(balance, payment, ratePerPeriod, k == periods-1)
		if err != nil {
			return nil, fmt.Errorf("period %d: %w", k+1, err)
		}
		row.Period = k + 1
		rows[k] = row
		balance = row.BalanceCents
		if payoff {
			return rows[:k+1], nil
		}
	}
	return rows, nil
}

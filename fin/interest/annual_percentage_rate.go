package interest

import "fmt"

// AnnualPercentageRate returns the annual percentage rate of a closed-end loan with a single advance
// and payments equal level payments, by the actuarial method of Regulation Z, 12 CFR 1026 Appendix
// J: the unit-period rate i solving amountFinanced = payment·(1 − (1+i)⁻ⁿ)/i, times the number of
// unit-periods in a year ((b)(1)). It assumes the unit-period is the payment period and the first
// period is a full regular period (t = 1, f = 0 in (c)(1)); odd first periods, irregular payments and
// multiple advances are not covered. The result is a decimal (0.0969 = 9.69 %), unrounded. The rate
// is AnnuityRate's with no final payment, found by the same bracketed Newton–bisection solver.
// Errors: ErrOutOfRange if |amountFinancedCents| or |paymentCents| > 2⁵³; ErrDomain if
// amountFinancedCents ≤ 0 or paymentCents ≤ 0; ErrPeriods if payments < 1 or periodsPerYear < 1;
// ErrNoSolution if the rate lies outside (−1 + 1e-12, 1e6]. (The solver's ErrNoConvergence cannot
// occur: its widest bracket needs at most 41 bisections.)
//
// ID:      FIN-INT-032
// Time:    O(k)   k ≤ 72 scan steps + 200 refinement steps
// Space:   O(1), no allocation
// Source:  12 CFR 1026 (Regulation Z) Appendix J (b)(1), (b)(8), (c)(1), read at https://www.consumerfinance.gov/rules-policy/regulations/1026/j/ (CFPB eRegulations) 2026-11-12; equal payments, regular first period (c)(1) with t = 1, f = 0
// Source:  W. H. Press, S. A. Teukolsky, W. T. Vetterling and B. P. Flannery, Numerical Recipes, 3rd ed., Cambridge University Press, 2007, §9.4 "Newton-Raphson Method Using Derivative" (rtsafe), p. 456 (solver)
// Tags:    interest, loan, apr, annual-percentage-rate, regulation-z, truth-in-lending, actuarial-method
// Units:   amountFinancedCents and paymentCents in cents; payments count; periodsPerYear payments per year; result decimal per year
// Example: $5,000 advanced, 24 monthly payments of $230 → 9.69 % (Appendix J (c)(1), Example (i))
func AnnualPercentageRate(amountFinancedCents, paymentCents int64, payments, periodsPerYear int) (float64, error) {
	if err := checkCents("amountFinancedCents", amountFinancedCents); err != nil {
		return 0, err
	}
	if err := checkCents("paymentCents", paymentCents); err != nil {
		return 0, err
	}
	if amountFinancedCents <= 0 {
		return 0, fmt.Errorf("%w: amountFinancedCents = %d, want > 0", ErrDomain, amountFinancedCents)
	}
	if paymentCents <= 0 {
		return 0, fmt.Errorf("%w: paymentCents = %d, want > 0", ErrDomain, paymentCents)
	}
	if err := checkPeriods("payments", payments, 1); err != nil {
		return 0, err
	}
	if err := checkPeriods("periodsPerYear", periodsPerYear, 1); err != nil {
		return 0, err
	}
	i, err := solveAnnuityRate(float64(amountFinancedCents), float64(paymentCents), 0, payments)
	if err != nil {
		return 0, err
	}
	return i * float64(periodsPerYear), nil
}

// Package loan holds loan and mortgage calculators: the amortization schedule and the balances it implies,
// balloon loans, prepayment, and the ratios lenders underwrite with (LTV, DTI, DSCR, PITI).
//
// ID short form: LOAN (IDs are FIN-LOAN-NNN).
//
// Money is int64 cents; every money result is rounded half-even (math.RoundToEven) to the cent. Rates are
// float64 decimals per payment period. The level payment comes from fin/interest
// (AnnuityPaymentFromPresentValue); every per-period interest amount is rounded half-even, and a schedule's
// last payment absorbs the rounding so the balance ends at exactly 0. RemainingBalance reproduces a
// schedule row exactly.
//
// Functions:
//
//	FIN-LOAN-001  AmortizationSchedule      per-period payment, interest, principal and balance
//	FIN-LOAN-002  RemainingBalance          balance after k payments, equal to schedule row k
//	FIN-LOAN-003  BalloonBalance            balance after k payments of a given amount
//	FIN-LOAN-004  BalloonLoanPayment        level payment that leaves a stated balloon
//	FIN-LOAN-005  ExtraPaymentPayoff        periods to payoff and interest saved with an extra payment
//	FIN-LOAN-006  LoanToValueRatio          loan / property value
//	FIN-LOAN-007  PITIPayment               principal and interest + taxes + insurance (+ mortgage insurance)
//	FIN-LOAN-008  DebtServiceCoverageRatio  net operating income / debt service
//	FIN-LOAN-009  DebtToIncomeRatio         monthly debt payments / gross monthly income
//	FIN-LOAN-010  AffordableLoanAmount      largest loan whose housing payment fits a housing ratio
package loan

// Package interest holds the time value of money: single-sum and annuity values, the nine
// engineering-economy factors, growing annuities and perpetuities, rate conversions (effective annual
// rate, Regulation DD APY, Regulation Z APR), day-count year fractions, net present value and the
// internal rate of return. Other fin packages build on it; it imports only the standard library.
//
// ID short form: INT (IDs are FIN-INT-NNN).
//
// Money is int64 cents. Every money input must satisfy |x| ≤ 2⁵³ so that its float64 working value is
// exact (ErrOutOfRange otherwise); money results are computed in float64 and rounded half-even
// (math.RoundToEven) to the cent, and a result beyond ±2⁵³ cents is ErrOverflow. Accuracy contract
// (finance decision F-5): a money result equals the correctly rounded exact value except when that value
// lies within about n·ε·|x| of a half cent, and is within 1 cent of it for |result| ≤ 2⁴³ cents.
// Rates are float64 decimals per stated period (0.05 = 5 %) and must be finite and > −1 (ErrRate). (1+i)ⁿ is evaluated
// as exp(n·log1p(i)) and (1+i)ⁿ − 1 as expm1(n·log1p(i)), which keeps the precision that
// pow(1+i, n) − 1 loses for small i (Goldberg 1991, "Rounding Error", "Cancellation"); every i = 0
// limit is exact. Factors are returned unrounded as float64. Scalar functions allocate nothing when
// they succeed (an error allocates its wrapped message). Errors wrap the package sentinels with the
// offending parameter; test with errors.Is.
//
// The only root finder in fin is here, behind InternalRateOfReturn, AnnuityRate and
// AnnualPercentageRate: a scan outward from rate 0, then Newton steps guarded by bisection
// (Press et al. 2007, §9.4).
//
// Functions:
//
//	FIN-INT-001  FutureValue                            PV·(1+i)ⁿ, cents
//	FIN-INT-002  PresentValue                           FV·(1+i)⁻ⁿ, cents
//	FIN-INT-003  FutureValueContinuous                  PV·e^(rt), cents
//	FIN-INT-004  PresentValueContinuous                 FV·e^(−rt), cents
//	FIN-INT-005  SinglePaymentCompoundAmount            (F/P,i,n)
//	FIN-INT-006  SinglePaymentPresentWorth              (P/F,i,n)
//	FIN-INT-007  UniformSeriesCompoundAmount            (F/A,i,n)
//	FIN-INT-008  UniformSeriesSinkingFund               (A/F,i,n)
//	FIN-INT-009  UniformSeriesPresentWorth              (P/A,i,n)
//	FIN-INT-010  CapitalRecovery                        (A/P,i,n)
//	FIN-INT-011  UniformGradientPresentWorth            (P/G,i,n)
//	FIN-INT-012  UniformGradientFutureWorth             (F/G,i,n)
//	FIN-INT-013  UniformGradientUniformSeries           (A/G,i,n)
//	FIN-INT-014  FutureValueAnnuity                     FV of an ordinary annuity, cents
//	FIN-INT-015  FutureValueAnnuityDue                  FV of an annuity due, cents
//	FIN-INT-016  PresentValueAnnuity                    PV of an ordinary annuity, cents
//	FIN-INT-017  PresentValueAnnuityDue                 PV of an annuity due, cents
//	FIN-INT-018  AnnuityPaymentFromPresentValue         level end-of-period payment for a PV (loan payment)
//	FIN-INT-019  AnnuityPaymentFromFutureValue          level end-of-period payment for an FV (sinking fund)
//	FIN-INT-020  AnnuityDuePaymentFromPresentValue      level start-of-period payment for a PV
//	FIN-INT-021  AnnuityDuePaymentFromFutureValue       level start-of-period payment for an FV
//	FIN-INT-022  PresentValueGrowingAnnuity             PV of a growing ordinary annuity, cents
//	FIN-INT-023  FutureValueGrowingAnnuity              FV of a growing ordinary annuity, cents
//	FIN-INT-024  GrowingAnnuityPaymentFromPresentValue  first payment of a growing annuity for a PV
//	FIN-INT-025  GrowingAnnuityPaymentFromFutureValue   first payment of a growing annuity for an FV
//	FIN-INT-026  PresentValuePerpetuity                 C/r, cents
//	FIN-INT-027  PresentValueGrowingPerpetuity          C/(r−g), cents; also terminal value
//	FIN-INT-029  EffectiveAnnualRate                    (1 + r/m)^m − 1
//	FIN-INT-030  EffectiveAnnualRateContinuous          e^r − 1
//	FIN-INT-031  AnnualPercentageYield                  Regulation DD Appendix A APY, decimal
//	FIN-INT-032  AnnualPercentageRate                   Regulation Z Appendix J APR, equal payments, decimal
//	FIN-INT-033  RuleOf72                               approximate periods to double, 72/(100·i)
//	FIN-INT-034  DoublingTime                           exact periods to double, ln 2/ln(1+i)
//	FIN-INT-035  SimpleInterest                         P·r·t, cents
//	FIN-INT-036  YearFraction30360                      MSRB Rule G-33 30/360 year fraction
//	FIN-INT-037  YearFractionActual365                  actual days/365 year fraction
//	FIN-INT-038  NetPresentValue                        Σ cₜ/(1+i)ᵗ, cents
//	FIN-INT-039  InternalRateOfReturn                   rate where NPV = 0
//	FIN-INT-040  GrowthRatePerPeriod                    (FV/PV)^(1/n) − 1 (CAGR)
//	FIN-INT-042  AnnuityRate                            rate solving PV = PMT·(P/A) + FV·(P/F) (bond yield)
//	FIN-INT-043  AnnuityPeriods                         number of level payments that repay a PV, N = −ln(1 − i·PV/PMT)/ln(1+i)
//	FIN-INT-044  PeriodicRateFromEffectiveAnnual        (1 + i_e)^(1/m) − 1, the inverse of EffectiveAnnualRate
//
// IDs held back, not built (no confirmed source; see doc/finance_2026-11-12/research/interest.md):
// FIN-INT-028 FutureValueAnnuityContinuous, FIN-INT-041 AddOnInterestPayment.
package interest

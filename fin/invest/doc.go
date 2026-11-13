// Package invest holds investment calculators: bond prices, yields and duration, Treasury bill
// rates, stock returns and dividend-discount values, CAPM, Black-Scholes option values, capital
// budgeting measures (MIRR, profitability index, payback), the Sharpe
// ratio, WACC, tax-equivalent yield, fund NAV and the bid-ask spread; and return-series statistics:
// simple and log returns, compound and annualised returns, historical, realized, EWMA and range-based
// (Parkinson, Garman–Klass, Rogers–Satchell) volatility, beta, tracking error, information and Sortino
// ratios, maximum drawdown, the Calmar ratio, historical and parametric VaR, expected shortfall and the
// Kelly fraction.
//
// ID short form: INVEST (IDs are FIN-INVEST-NNN).
//
// Money is int64 cents (|x| ≤ 2⁵³); money results are rounded half-even. Rates, yields and ratios
// are float64 decimals (0.05 = 5 %). Return series are []float64 of decimal returns per period; every
// annualisation factor, decay factor and confidence level is a parameter. Errors wrap the package
// sentinels (ErrRate, ErrDomain, ErrTooFew, …).
// Time-value factors and the root finder come from package fin/interest.
//
// Functions:
//
//	FIN-INVEST-001  CurrentYield                    annual coupon / price
//	FIN-INVEST-002  BondPrice                       price of a level-coupon bond on a coupon date
//	FIN-INVEST-003  YieldToMaturity                 nominal annual yield to maturity from price
//	FIN-INVEST-004  MacaulayDuration                PV-weighted average time of a bond's flows, years
//	FIN-INVEST-005  ModifiedDuration                −(1/P)·dP/dy, years
//	FIN-INVEST-006  Convexity                       (1/P)·d²P/dy², years²
//	FIN-INVEST-007  BankDiscountYield               Treasury bill discount rate, 360-day basis
//	FIN-INVEST-009  BondEquivalentYield             Treasury bill investment rate (≤ half a year)
//	FIN-INVEST-010  HoldingPeriodYield              total return over a holding period
//	FIN-INVEST-011  AnnualizedHoldingPeriodYield    effective annual rate of a holding-period return
//	FIN-INVEST-012  CapitalGainsYield               price change / beginning price
//	FIN-INVEST-013  DividendYield                   annual dividend / price
//	FIN-INVEST-014  GordonGrowthPrice               D₁/(r − g)
//	FIN-INVEST-015  ZeroGrowthStockPrice            D/r (also preferred stock)
//	FIN-INVEST-016  CAPMExpectedReturn              Rf + β(Rm − Rf)
//	FIN-INVEST-017  RiskPremium                     expected return − risk-free rate
//	FIN-INVEST-018  BlackScholesCall                European call value
//	FIN-INVEST-019  BlackScholesPut                 European put value
//	FIN-INVEST-020  ModifiedInternalRateOfReturn    MIRR with finance and reinvestment rates
//	FIN-INVEST-021  ProfitabilityIndex              PV(inflows)/PV(outflows)
//	FIN-INVEST-022  PaybackPeriod                   periods to recover the investment
//	FIN-INVEST-023  DiscountedPaybackPeriod         payback on discounted flows
//	FIN-INVEST-024  REMOVED → FIN-INT-018 (equivalent annual annuity = annuity payment of the NPV)
//	FIN-INVEST-025  SharpeRatio                     (Rp − Rf)/σp
//	FIN-INVEST-026  WeightedAverageCostOfCapital    WACC with debt and equity
//	FIN-INVEST-027  TaxEquivalentYield              tax-free yield/(1 − T)
//	FIN-INVEST-028  NetAssetValuePerShare           (assets − liabilities)/shares
//	FIN-INVEST-030  BidAskSpread                    relative quoted spread
//	FIN-INVEST-032  SimpleReturns                   Pₜ/Pₜ₋₁ − 1 for a price series (dst)
//	FIN-INVEST-033  LogReturns                      ln(Pₜ/Pₜ₋₁) for a price series (dst)
//	FIN-INVEST-034  CumulativeReturn                Π(1 + Rₜ) − 1
//	FIN-INVEST-035  AnnualizedReturn                [Π(1 + Rₜ)]^(periodsPerYear/n) − 1
//	FIN-INVEST-036  AnnualizedVolatility            sample s · √periodsPerYear
//	FIN-INVEST-037  RealizedVariance                Σ rⱼ² (ABDL 2003)
//	FIN-INVEST-038  RealizedVolatility              √(Σ rⱼ²)
//	FIN-INVEST-039  EWMAVariance                    σ² ← λσ² + (1 − λ)r² (RiskMetrics)
//	FIN-INVEST-040  ParkinsonVolatility             high-low range estimator
//	FIN-INVEST-041  GarmanKlassVolatility           open-high-low-close estimator
//	FIN-INVEST-042  RogersSatchellVolatility        drift-independent OHLC estimator
//	FIN-INVEST-043  Beta                            Cov(Rᵢ, R_M)/Var(R_M)
//	FIN-INVEST-044  TrackingError                   sample s of active returns
//	FIN-INVEST-045  InformationRatio                mean active return / tracking error
//	FIN-INVEST-046  SortinoRatio                    (mean − MAR)/downside deviation
//	FIN-INVEST-047  MaxDrawdown                     largest peak-to-trough decline, with indices
//	FIN-INVEST-048  CalmarRatio                     annualised return / maximum drawdown
//	FIN-INVEST-049  HistoricalVaR                   −Q(1 − c) of observed returns
//	FIN-INVEST-050  ParametricVaR                   σ·z_c − μ (normal)
//	FIN-INVEST-051  ExpectedShortfall               empirical ES with the quantile-atom correction
//	FIN-INVEST-052  KellyFraction                   p − q/b
//	FIN-INVEST-053  CapitalizationRate              annual NOI / price (cap rate)
//
// Not built (no confirmed source): FIN-INVEST-008 MoneyMarketYield, FIN-INVEST-029 LeveredBetaHamada,
// FIN-INVEST-031 ImpliedRepoRate (doc/finance_2026-11-12/research/invest.md; 031 stays reserved for it).
package invest

# Finance build brief (for the four package agents; the lead owns everything else)

Repo: /Users/devinanderson/Desktop/calcbench. Module github.com/CrimsonCoder42/calcbench (go 1.27).
Read first: `.claude/CLAUDE.md`, `.claude/rules/go-functions.md`, `go-testing.md`, `go-complexity.md`,
`numeric.md`, `sources.md`. Other subject agents are live in this tree; follow CLAUDE.md "Several subjects at once".

## What you own (and nothing else)

| Agent | Package(s) | Research file you write | IDs |
|---|---|---|---|
| A | `fin/interest/` | `doc/finance_2026-11-12/research/interest.md` | FIN-INT-001… |
| B | `fin/invest/` | `doc/finance_2026-11-12/research/invest.md` | FIN-INVEST-001… |
| C | `fin/business/` | `doc/finance_2026-11-12/research/business.md` | FIN-BIZ-001… |
| D | `fin/loan/`, `fin/economics/` | `doc/finance_2026-11-12/research/loan.md`, `economics.md` | FIN-LOAN-001…, FIN-ECON-001… |

Never edit INDEX.md, doc/catalog/*, doc/DECISIONS.md or any other shared file. Anything outside your paths
goes into your research file under "FOR THE LEAD". No git writes. No `rm` (use `mv -n` to `_delete/`). No
new dependency: standard library only. Do NOT import `num/...`. Do NOT fetch vcalc.com.

## Your job, in order

1. **Research** every function in your list. Confirm each citation by opening a page that shows it, and write one entry
   per function in the exact catalog shape below into your research file, in ID order.
2. **Build** each confirmed function: code, doc comment, `_test.go`, benchmark, and the package `doc.go`.
   After each file run only `gofmt -l`, `go build ./fin/<pkg>/` and `go vet ./fin/<pkg>/`.
   **You never run `go test`.** The lead runs the tests once, at the end. Do not run benchmarks either.
3. Finish with a short hand-back: the functions built (ID, name, file:line), the functions held back and why,
   and anything for the lead.

A function whose citation cannot be confirmed (D-015) is NOT built. Its entry says
`SOURCE NOT CONFIRMED: <what was tried>`, and it is listed in your hand-back. A function with no acceptable
source at all says `SOURCE NOT FOUND`. Never invent a formula, constant, example or page number.

## Sources (decided by the lead; logged as finance decisions F-1 and F-2)

Acceptable as the primary source:
- Brealey-Myers-Allen; Ross-Westerfield-Jaffe; Bodie-Kane-Marcus; Fabozzi; Hull (Black-Scholes); Newnan
  (factor notation); CFA Institute readings.
- Regulations: 12 CFR 1026 (Regulation Z) Appendix J (APR); 12 CFR 1030 Appendix A (APY). Read them on ecfr.gov.
- Official bodies' pages: Federal Reserve, BLS, US Treasury and TreasuryDirect, SEC, IRS (Pub. 946 depreciation),
  FFIEC (UBPR user guide), DOJ/FTC Merger Guidelines (HHI), MSRB Rule G-33 (day counts), NCEES FE Reference
  Handbook (engineering-economy factors), CFPB.
- The original paper, with its DOI confirmed on Crossref (https://api.crossref.org/works?query.bibliographic=…).
  Examples: Black & Scholes 1973, Sharpe 1964, Gordon & Shapiro 1956, Macaulay 1938 (NBER), Hamada 1972,
  Harris 1913 (EOQ; reprint Operations Research 1990), Fisher 1930 *The Theory of Interest*.
- **F-1: the OpenStax textbooks (Rice University, CC BY 4.0, free on openstax.org).** These are *Principles of
  Finance*, *Principles of Accounting* Vol. 1 and 2, and *Principles of Economics* 3e. They count as a recognised textbook
  and as the confirmed second source when the named paid textbook cannot be opened (D-015 allows a confirmed second source).
  Cite them as author(s), title, edition, OpenStax, year, chapter.section, and the URL of the section page.
  Their "Example"/"Solution" boxes are the preferred worked examples, with the section as the locator.
- Numerics (for the method, not the formula): Goldberg, "What Every Computer Scientist Should Know About
  Floating-Point Arithmetic", ACM Computing Surveys 23(1):5–48, 1991, doi:10.1145/103162.103163. Press et al.,
  *Numerical Recipes* 3rd ed., Cambridge University Press, 2007, §9.4 (`rtsafe`).
Not acceptable as primary: vCalc, Wikipedia (a pointer only), Investopedia and other blogs, Q&A sites, library code.

## Numeric policy (Devin decided; do not change it)

- **Money** is `int64` minor units (cents). Every money parameter and result name ends in `Cents`. Never a float
  for a currency amount, in or out (prices per share and option values are money too).
- Money inputs must satisfy |x| ≤ 2⁵³ (`maxExactCents = 1 << 53`) so the float64 working value is exact;
  otherwise ErrOutOfRange. Money results are computed in float64, then rounded **half-even**
  (`math.RoundToEven`) to int64 cents; |result| > 2⁵³ gives ErrOverflow. Half-up only where a regulation says
  so: cite it. The doc comment states which rounding applies.
- **Rates** are `float64`, decimal per stated period (0.05 = 5 %), named `ratePerPeriod`, `ratePerYear`,
  `growthPerPeriod`, etc. A rate must be finite and > −1 (ErrRate). Ratios and yields are returned as decimals
  (0.05), never as percentages, except where the defining source states a percentage (HHI, CPI): say so.
- **Period counts** are `int` named `periods`, `periodsPerYear`, `days`, … (ErrPeriods when out of range).
  Dates are `time.Time` passed in; never `time.Now`.
- **Compounding**: (1+i)ⁿ = `math.Exp(float64(n)*math.Log1p(i))` and (1+i)ⁿ − 1 = `math.Expm1(float64(n)*math.Log1p(i))`.
  This avoids the cancellation of pow(1+i, n) − 1 for small i (Goldberg 1991, §1.4 / Theorem 4, which uses
  exactly the compound-interest example; cite the section you confirm). The i = 0 limits are handled exactly
  (n, n(n−1)/2, …), never by dividing by zero.
- **Schedules** return slices sized up front: one allocation, no append growth.
- **IRR/YTM**: only `fin/interest` has the root finder; it is exported through `InternalRateOfReturn`. It works in
  two phases. First, bracket from rate 0 outward, upward by doubling to 1e6 and downward toward −1 + 1e-12. Then
  run Newton steps guarded by bisection (Numerical Recipes §9.4 `rtsafe`) on NPV(rate), with tolerance
  |Δrate| ≤ 1e-12·max(1,|rate|) and at most 200 iterations. ErrNoConvergence if not met; ErrNoSolution if no sign change is found.
  Document that a series with several sign changes can have several IRRs, and which one is returned (the root
  bracketed nearest 0).
- **Black-Scholes** uses `math.Erf` for the normal CDF: N(x) = ½·erfc(−x/√2) (`math.Erfc`, accurate in the tails).
- **Ratios**: inputs in cents (int64), result float64. A zero denominator gives ErrZeroDenominator. A negative
  denominator is allowed where it means something (negative equity), and the doc says so.
- **Errors**: package sentinels, wrapped with the parameter name (`fmt.Errorf("%w: ratePerPeriod = %v", ErrRate, r)`).
  Use these names where they apply: ErrNotFinite, ErrOutOfRange, ErrOverflow, ErrRate, ErrPeriods,
  ErrZeroDenominator, ErrNegative, ErrDomain, ErrNoConvergence, ErrNoSolution. The `var (...)` block lives in
  the file of the first function that uses them (doc/ASSUMPTIONS.md A-008). The doc comment lists every error
  the function can return.
- **No allocation** for scalar functions (test it with `testing.AllocsPerRun` == 0).

## Doc comment (exact shape; see `.claude/rules/go-functions.md`)

```go
// FutureValue returns the value after periods compounding periods of presentValueCents at ratePerPeriod,
// rounded half-even to the cent. presentValueCents is read only. …what it assumes…
// Errors: ErrOutOfRange …; ErrRate …; ErrPeriods …; ErrOverflow ….
//
// ID:      FIN-INT-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  <primary: author(s), title, edition, publisher, year, §/eq./table, URL or DOI>
// Source:  <second source if any, same detail>
// Tags:    interest, time-value, future-value, compound
// Units:   presentValueCents and the result in cents; ratePerPeriod decimal per period; periods count
// Example: <the source's worked numbers, with the locator> (or "hand-checked")
func FutureValue(presentValueCents int64, ratePerPeriod float64, periods int) (int64, error) {
```

## Research entry shape (one per function, in your research file)

```
### <Name>      (ID: FIN-XXX-NNN)
Formula: <define every symbol>
Inputs and units: …
Output: …
Errors: …
Primary source: <full citation, cited to the line; CONFIRMED <how: URL opened 2026-11-12>>
Worked example: <inputs → expected, with the source's locator; or HAND-CHECKED with the arithmetic shown>
Variants: …
Best approach: … (time, space, allocations)
Notes: … (what it must NOT do)
Covers rows: <the vCalc rows it covers, from the lead's list>
```

## Tests (`.claude/rules/go-testing.md`)

Each test file contains, for each function:
- The source's worked example, cited in a comment, compared exactly in cents (a rate result to a stated tolerance, with the reason).
- Table cases: a zero rate, zero periods, one period, a negative amount where meaningful, the extremes
  (2⁵³ ok, 2⁵³+1 → ErrOutOfRange), NaN and ±Inf rates.
- Every documented error.
- An oracle where one exists:
  - money results against `math/big` (`big.Float` with 256-bit precision, or `big.Rat` for integer powers);
  - identities between factors (F/A = (F/P − 1)/i, A/P = A/F + i, …);
  - an amortization schedule that reconciles to the cent (Σprincipal = loan, final balance = 0);
  - NPV(IRR) ≈ 0.
- A benchmark with `b.ReportAllocs()`: one case for scalars; 1e3/1e4/1e5 for slice inputs and schedules.
- An `AllocsPerRun` == 0 test for each scalar. Test names follow `Test<Function>_<case>`; use fixed seeds only.

## Function lists (IDs fixed by the lead, in catalog order; signatures are the contract)

### A: fin/interest (all others may import these; build 018, 010, 038, 039 FIRST, others depend on them)
001 FutureValue(presentValueCents int64, ratePerPeriod float64, periods int) (int64, error)
002 PresentValue(futureValueCents int64, ratePerPeriod float64, periods int) (int64, error)
003 FutureValueContinuous(presentValueCents int64, ratePerYear, years float64) (int64, error)
004 PresentValueContinuous(futureValueCents int64, ratePerYear, years float64) (int64, error)
005 SinglePaymentCompoundAmount(ratePerPeriod float64, periods int) (float64, error)   // (F/P,i,n)
006 SinglePaymentPresentWorth(ratePerPeriod float64, periods int) (float64, error)     // (P/F,i,n)
007 UniformSeriesCompoundAmount(ratePerPeriod float64, periods int) (float64, error)   // (F/A,i,n)
008 UniformSeriesSinkingFund(ratePerPeriod float64, periods int) (float64, error)      // (A/F,i,n)
009 UniformSeriesPresentWorth(ratePerPeriod float64, periods int) (float64, error)     // (P/A,i,n)
010 CapitalRecovery(ratePerPeriod float64, periods int) (float64, error)               // (A/P,i,n)
011 UniformGradientPresentWorth(ratePerPeriod float64, periods int) (float64, error)   // (P/G,i,n)
012 UniformGradientFutureWorth(ratePerPeriod float64, periods int) (float64, error)    // (F/G,i,n)
013 UniformGradientUniformSeries(ratePerPeriod float64, periods int) (float64, error)  // (A/G,i,n)
014 FutureValueAnnuity(paymentCents int64, ratePerPeriod float64, periods int) (int64, error)
015 FutureValueAnnuityDue(paymentCents int64, ratePerPeriod float64, periods int) (int64, error)
016 PresentValueAnnuity(paymentCents int64, ratePerPeriod float64, periods int) (int64, error)
017 PresentValueAnnuityDue(paymentCents int64, ratePerPeriod float64, periods int) (int64, error)
018 AnnuityPaymentFromPresentValue(presentValueCents int64, ratePerPeriod float64, periods int) (int64, error)  // = level loan payment
019 AnnuityPaymentFromFutureValue(futureValueCents int64, ratePerPeriod float64, periods int) (int64, error)
020 AnnuityDuePaymentFromPresentValue(presentValueCents int64, ratePerPeriod float64, periods int) (int64, error)
021 AnnuityDuePaymentFromFutureValue(futureValueCents int64, ratePerPeriod float64, periods int) (int64, error)
022 PresentValueGrowingAnnuity(firstPaymentCents int64, ratePerPeriod, growthPerPeriod float64, periods int) (int64, error)
023 FutureValueGrowingAnnuity(firstPaymentCents int64, ratePerPeriod, growthPerPeriod float64, periods int) (int64, error)
024 GrowingAnnuityPaymentFromPresentValue(presentValueCents int64, ratePerPeriod, growthPerPeriod float64, periods int) (int64, error) // first payment
025 GrowingAnnuityPaymentFromFutureValue(futureValueCents int64, ratePerPeriod, growthPerPeriod float64, periods int) (int64, error)
026 PresentValuePerpetuity(paymentCents int64, ratePerPeriod float64) (int64, error)
027 PresentValueGrowingPerpetuity(firstPaymentCents int64, ratePerPeriod, growthPerPeriod float64) (int64, error)  // also terminal value
028 FutureValueAnnuityContinuous(paymentCents int64, ratePerYear float64, years int) (int64, error)  // end-of-year payments, continuous compounding; build only if sourced
029 EffectiveAnnualRate(nominalRatePerYear float64, periodsPerYear int) (float64, error)
030 EffectiveAnnualRateContinuous(nominalRatePerYear float64) (float64, error)
031 AnnualPercentageYield(principalCents, interestCents int64, daysInTerm int) (float64, error)  // 12 CFR 1030 App. A; decimal, unrounded
032 AnnualPercentageRate(amountFinancedCents, paymentCents int64, payments, periodsPerYear int) (float64, error)  // Reg Z App. J, equal payments, unit period = payment period
033 RuleOf72(ratePerPeriod float64) (float64, error)      // approximate periods to double
034 DoublingTime(ratePerPeriod float64) (float64, error)  // exact ln 2 / ln(1+i)
035 SimpleInterest(principalCents int64, ratePerYear, yearFraction float64) (int64, error)
036 YearFraction30360(start, end time.Time) (float64, error)    // MSRB Rule G-33 30/360
037 YearFractionActual365(start, end time.Time) (float64, error)
038 NetPresentValue(ratePerPeriod float64, cashFlowsCents []int64) (int64, error)  // cashFlowsCents[0] at t = 0
039 InternalRateOfReturn(cashFlowsCents []int64) (float64, error)
040 GrowthRatePerPeriod(presentValueCents, futureValueCents int64, periods int) (float64, error)  // CAGR; zero-coupon yield; "investment return rate"
041 AddOnInterestPayment(principalCents int64, ratePerYear, years float64, payments int) (int64, error)

### B: fin/invest (imports fin/interest for 003 [via InternalRateOfReturn], 020, 021, 023, 024)
001 CurrentYield(annualCouponCents, priceCents int64) (float64, error)
002 BondPrice(faceCents int64, couponRatePerYear, yieldPerYear float64, periods, periodsPerYear int) (int64, error)  // on a coupon date
003 YieldToMaturity(priceCents, faceCents int64, couponRatePerYear float64, periods, periodsPerYear int) (float64, error)  // nominal annual = periodsPerYear × periodic
004 MacaulayDuration(faceCents int64, couponRatePerYear, yieldPerYear float64, periods, periodsPerYear int) (float64, error)  // years
005 ModifiedDuration(same) (float64, error)
006 Convexity(same) (float64, error)  // years²
007 BankDiscountYield(faceCents, priceCents int64, daysToMaturity int) (float64, error)
008 MoneyMarketYield(faceCents, priceCents int64, daysToMaturity int) (float64, error)
009 BondEquivalentYield(faceCents, priceCents int64, daysToMaturity, daysInYear int) (float64, error)  // Treasury investment rate, ≤ ½ year
010 HoldingPeriodYield(beginningValueCents, endingValueCents, incomeCents int64) (float64, error)  // also total stock return
011 AnnualizedHoldingPeriodYield(holdingPeriodYield float64, days int) (float64, error)
012 CapitalGainsYield(beginningPriceCents, endingPriceCents int64) (float64, error)
013 DividendYield(annualDividendPerShareCents, pricePerShareCents int64) (float64, error)
014 GordonGrowthPrice(nextDividendCents int64, requiredReturn, growthRate float64) (int64, error)
015 ZeroGrowthStockPrice(dividendCents int64, requiredReturn float64) (int64, error)  // also preferred stock
016 CAPMExpectedReturn(riskFreeRate, beta, marketReturn float64) (float64, error)
017 RiskPremium(expectedReturn, riskFreeRate float64) (float64, error)
018 BlackScholesCall(spotCents, strikeCents int64, riskFreeRatePerYear, volatilityPerYear, years float64) (int64, error)
019 BlackScholesPut(spotCents, strikeCents int64, riskFreeRatePerYear, volatilityPerYear, years float64) (int64, error)
020 ModifiedInternalRateOfReturn(cashFlowsCents []int64, financeRatePerPeriod, reinvestRatePerPeriod float64) (float64, error)
021 ProfitabilityIndex(ratePerPeriod float64, cashFlowsCents []int64) (float64, error)
022 PaybackPeriod(cashFlowsCents []int64) (float64, error)
023 DiscountedPaybackPeriod(ratePerPeriod float64, cashFlowsCents []int64) (float64, error)
024 EquivalentAnnualAnnuity(netPresentValueCents int64, ratePerPeriod float64, periods int) (int64, error)
025 SharpeRatio(portfolioReturn, riskFreeRate, returnStdDev float64) (float64, error)
026 WeightedAverageCostOfCapital(equityValueCents, debtValueCents int64, costOfEquity, costOfDebt, taxRate float64) (float64, error)
027 TaxEquivalentYield(taxFreeYield, marginalTaxRate float64) (float64, error)
028 NetAssetValuePerShare(totalAssetsCents, totalLiabilitiesCents, sharesOutstanding int64) (int64, error)
029 LeveredBetaHamada(unleveredBeta, taxRate, debtToEquity float64) (float64, error)
030 BidAskSpread(bidCents, askCents int64) (float64, error)  // relative spread; define per source
031 ImpliedRepoRate(spotPriceCents, futuresPriceCents int64, days int) (float64, error)  // build only if sourced

### C: fin/business
001 CurrentRatio(currentAssetsCents, currentLiabilitiesCents int64) (float64, error)
002 QuickRatio(cashCents, marketableSecuritiesCents, receivablesCents, currentLiabilitiesCents int64) (float64, error)  // acid test
003 DebtRatio(totalLiabilitiesCents, totalAssetsCents int64) (float64, error)
004 DebtToEquityRatio(totalLiabilitiesCents, totalEquityCents int64) (float64, error)
005 EquityMultiplier(totalAssetsCents, totalEquityCents int64) (float64, error)
006 EquityToDebtRatio(totalEquityCents, totalLiabilitiesCents int64) (float64, error)  // owners' to creditors' equity
007 TimesInterestEarned(ebitCents, interestExpenseCents int64) (float64, error)  // interest coverage
008 FixedChargeCoverageRatio(ebitCents, fixedChargesCents, interestExpenseCents int64) (float64, error)
009 CashFlowAdequacyRatio(operatingCashFlowCents, capitalExpendituresCents, debtRepaymentsCents, dividendsPaidCents int64) (float64, error)
010 CashFlowToCurrentLiabilities(operatingCashFlowCents, averageCurrentLiabilitiesCents int64) (float64, error)  // CFOCL
011 CashFlowToCapitalExpenditures(operatingCashFlowCents, capitalExpendituresCents int64) (float64, error)  // CFOCE
012 AssetTurnover(netSalesCents, averageTotalAssetsCents int64) (float64, error)
013 AssetToSalesRatio(totalAssetsCents, netSalesCents int64) (float64, error)
014 ReceivablesTurnover(netCreditSalesCents, averageReceivablesCents int64) (float64, error)
015 AverageCollectionPeriod(averageReceivablesCents, netCreditSalesCents int64, daysInPeriod int) (float64, error)  // days
016 InventoryTurnover(costOfGoodsSoldCents, averageInventoryCents int64) (float64, error)
017 DaysInInventory(averageInventoryCents, costOfGoodsSoldCents int64, daysInPeriod int) (float64, error)
018 GrossProfitMargin(netSalesCents, costOfGoodsSoldCents int64) (float64, error)
019 OperatingProfitMargin(operatingIncomeCents, netSalesCents int64) (float64, error)
020 NetProfitMargin(netIncomeCents, netSalesCents int64) (float64, error)
021 EarningsPerShare(netIncomeCents, preferredDividendsCents, weightedAverageShares int64) (int64, error)  // cents per share, half-even
022 DilutedEarningsPerShare(netIncomeCents, preferredDividendsCents, convertibleInterestAfterTaxCents, weightedAverageShares, dilutiveShares int64) (int64, error)
023 PriceToEarningsRatio(pricePerShareCents, earningsPerShareCents int64) (float64, error)
024 PriceToBookRatio(pricePerShareCents, bookValuePerShareCents int64) (float64, error)
025 PriceToSalesRatio(marketCapitalizationCents, revenueCents int64) (float64, error)
026 BookValuePerShare(totalEquityCents, preferredEquityCents, sharesOutstanding int64) (int64, error)
027 DividendsPerShare(totalDividendsCents, sharesOutstanding int64) (int64, error)
028 DividendPayoutRatio(dividendsCents, netIncomeCents int64) (float64, error)
029 RetentionRatio(dividendsCents, netIncomeCents int64) (float64, error)
030 ReturnOnAssets(netIncomeCents, averageTotalAssetsCents int64) (float64, error)
031 ReturnOnEquity(netIncomeCents, averageEquityCents int64) (float64, error)
032 ReturnOnCommonEquity(netIncomeCents, preferredDividendsCents, averageCommonEquityCents int64) (float64, error)  // RCSE
033 ReturnOnInvestment(finalValueCents, costCents int64) (float64, error)  // (final − cost)/cost; also advertising ROI
034 EconomicValueAdded(nopatCents, investedCapitalCents int64, costOfCapital float64) (int64, error)
035 FreeCashFlowToFirm(ebitCents int64, taxRate float64, depreciationCents, capitalExpendituresCents, changeInWorkingCapitalCents int64) (int64, error)
036 FreeCashFlowToEquity(fcffCents, interestExpenseCents int64, taxRate float64, netBorrowingCents int64) (int64, error)
037 NetWorkingCapital(currentAssetsCents, currentLiabilitiesCents int64) (int64, error)
038 DegreeOfOperatingLeverage(contributionMarginCents, operatingIncomeCents int64) (float64, error)
039 DegreeOfFinancialLeverage(ebitCents, interestExpenseCents int64) (float64, error)
040 ContributionMargin(salesCents, variableCostsCents int64) (int64, error)
041 ContributionMarginRatio(salesCents, variableCostsCents int64) (float64, error)
042 BreakEvenUnits(fixedCostsCents, pricePerUnitCents, variableCostPerUnitCents int64) (float64, error)
043 EconomicOrderQuantity(annualDemandUnits float64, orderCostCents, holdingCostPerUnitPerYearCents int64) (float64, error)  // units
044 OptimalOrderCycle(annualDemandUnits float64, orderCostCents, holdingCostPerUnitPerYearCents int64) (float64, error)  // years between orders = EOQ/D
045 ReorderPoint(dailyDemandUnits, leadTimeDays, safetyStockUnits float64) (float64, error)  // units
046 DepreciationStraightLine(costCents, salvageCents int64, lifePeriods int) ([]int64, error)  // per-period schedule; remainder in the last period
047 DepreciationDecliningBalance(costCents, salvageCents int64, lifePeriods int, rateMultiple float64) ([]int64, error)  // 2 = double declining; never below salvage
048 GoodwillFromAcquisition(purchasePriceCents, fairValueOfNetAssetsCents int64) (int64, error)  // negative = bargain purchase, documented
049 BankEfficiencyRatio(noninterestExpenseCents, netInterestIncomeCents, noninterestIncomeCents int64) (float64, error)
050 LoanToDepositRatio(totalLoansCents, totalDepositsCents int64) (float64, error)

### D1: fin/loan (imports fin/interest for 018 and 010)
001 AmortizationSchedule(principalCents int64, ratePerPeriod float64, periods int) ([]AmortizationRow, error)
    `type AmortizationRow struct{ Period int; PaymentCents, InterestCents, PrincipalCents, BalanceCents int64 }`.
    payment = interest.AnnuityPaymentFromPresentValue. Each row's interest is rounded half-even. The final
    payment is adjusted so the balance ends at exactly 0.
002 RemainingBalance(principalCents int64, ratePerPeriod float64, periods, paymentsMade int) (int64, error)  // equals schedule row paymentsMade exactly (iterate O(k), no allocation)
003 BalloonBalance(principalCents, paymentCents int64, ratePerPeriod float64, paymentsMade int) (int64, error)  // balance after k level payments, same per-period rounding
004 BalloonLoanPayment(principalCents, balloonCents int64, ratePerPeriod float64, periods int) (int64, error)
005 ExtraPaymentPayoff(principalCents int64, ratePerPeriod float64, periods int, extraPerPeriodCents int64) (periodsToPayoff int, interestSavedCents int64, err error)
006 LoanToValueRatio(loanCents, propertyValueCents int64) (float64, error)
007 PITIPayment(principalAndInterestCents, annualPropertyTaxCents, annualInsuranceCents, monthlyMortgageInsuranceCents int64) (int64, error)
008 DebtServiceCoverageRatio(netOperatingIncomeCents, debtServiceCents int64) (float64, error)
009 DebtToIncomeRatio(monthlyDebtPaymentsCents, grossMonthlyIncomeCents int64) (float64, error)
010 AffordableLoanAmount(grossMonthlyIncomeCents int64, maxHousingRatio float64, monthlyTaxesAndInsuranceCents int64, ratePerPeriod float64, periods int) (int64, error)

### D2: fin/economics
001 ConsumerPriceIndex(basketCostCurrentCents, basketCostBaseCents int64) (float64, error)  // index, base = 100 (BLS)
002 InflationRate(previousIndex, currentIndex float64) (float64, error)
003 RealInterestRate(nominalRate, inflationRate float64) (float64, error)  // exact Fisher (1+i)/(1+π) − 1; also inflation-adjusted return
004 HerfindahlHirschmanIndex(marketSharesPercent []float64) (float64, error)  // DOJ/FTC: sum of squared percent shares, 0–10 000

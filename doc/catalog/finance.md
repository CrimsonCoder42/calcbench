# Catalog: finance

Status: BUILT 2026-11-12 (catalog and build in one run, per Devin's instruction). Work folder: `../finance_2026-11-12/`.
Input: the names in vCalc's Finance category, pasted by Devin. vCalc pages were not fetched, and nothing but the names was taken from it.
Decisions logged in `../finance_2026-11-12/PENDING_SHARED.md` (F-1…F-n), each with the alternative that lost.
Packages: `fin/interest` (FIN-INT), `fin/loan` (FIN-LOAN), `fin/invest` (FIN-INVEST), `fin/business` (FIN-BIZ),
`fin/economics` (FIN-ECON), plus a proposed `fin/personal` for later work (salary, retirement).

## Step 1: Clean

One row per input name, in the order given. A name marked "(x2)" in the input is two rows; the second
is MERGEd into the first. KEEP and MERGE rows name the function ID that covers them.

| # | Name as given | Decision | Home | Reason |
|---|---|---|---|---|
| | **Calculators and collections** | | | |
| c1 | Annuity Calculator | DROP | — | A calculator (its formulas are FIN-INT-014…025) |
| c2 | FINANCE 21J01ABT033 | DROP | — | A collection with a course code for a name |
| c3 | Finance 20s01abt023 | DROP | — | A collection |
| c4 | Finance Calc | DROP | — | A calculator |
| c5 | Finance Calculator | DROP | — | A calculator |
| c6 | Finance Calculator | DROP | — | A calculator (duplicate name) |
| c7 | Finance Monitoring Calculator | DROP | — | A calculator |
| c8 | Future Value (FV) Calculator | DROP | — | A calculator (formula is FIN-INT-001) |
| c9 | Goodwill Calculator | DROP | — | A calculator (formula is FIN-BIZ-048) |
| c10 | Investment & Loan Calculator | DROP | — | A calculator |
| c11 | Investment Help Calculator | DROP | — | A calculator |
| c12 | Loan and Mortage | DROP | — | A collection (formulas are in `fin/loan`) |
| c13 | Microeconomics Calculator | DROP | — | A calculator; its formulas are not named |
| c14 | Salary Calculator | MOVE | `fin/personal` (proposed) | Salary row (F-4) |
| | **Equations** | | | |
| 1 | 2018 Federal Tax Rate Schedules | DROP | — | DATA: a year's tax schedule is a table, not a function |
| 2 | Accounting Equation (Capital) | DROP | — | One subtraction (equity = assets − liabilities); no rule beyond it (D-008 by analogy) |
| 3 | Accounting Equation [Data set] | DROP | — | A data set |
| 4 | Acid Test Ratio | KEEP | `fin/business` FIN-BIZ-002 | Quick ratio |
| 5 | Advertising ROI | MERGE | → FIN-BIZ-033 | ROI with a domain label |
| 6 | Annual Percentage Yield | KEEP | `fin/interest` FIN-INT-031 | APY (12 CFR 1030) |
| 7 | Annual Percentage Yield | MERGE | → FIN-INT-031 | Duplicate |
| 8 | Annualized Holding Period Yield | KEEP | `fin/invest` FIN-INVEST-011 | |
| 9 | Annuity - Future Value with Continuous Compounding | KEEP | `fin/interest` FIN-INT-028 | Built only if sourced |
| 10 | Annuity Due-Payment (FV) | KEEP | `fin/interest` FIN-INT-021 | |
| 11 | Annuity Due-Payment (PV) | KEEP | `fin/interest` FIN-INT-020 | |
| 12 | Annuity-Payment (FV) | KEEP | `fin/interest` FIN-INT-019 | |
| 13 | Annuity-Payment (PV) | KEEP | `fin/interest` FIN-INT-018 | Also the level loan payment |
| 14 | APY (FDIC) | MERGE | → FIN-INT-031 | The same 12 CFR 1030 formula |
| 15 | Arbitrary Goodwill Value Agreement | DROP | — | An agreed value, not a formula |
| 16 | Asset to Sales Ratio | KEEP | `fin/business` FIN-BIZ-013 | |
| 17 | Asset Turnover | KEEP | `fin/business` FIN-BIZ-012 | |
| 18 | Asset Turnover Ratio | MERGE | → FIN-BIZ-012 | Duplicate |
| 19 | Average Collection Period | KEEP | `fin/business` FIN-BIZ-015 | |
| 20 | Balloon Loan Payments | KEEP | `fin/loan` FIN-LOAN-004 | |
| 21 | Bank Discount Yield | KEEP | `fin/invest` FIN-INVEST-007 | |
| 22 | Bank Efficiency Ratio | KEEP | `fin/business` FIN-BIZ-049 | |
| 23 | Benefits of Extra Mortgage Payments | KEEP | `fin/loan` FIN-LOAN-005 | |
| 24 | Bid Ask Spread | KEEP | `fin/invest` FIN-INVEST-030 | |
| 25 | Bid Ask Spread_Copy | MERGE | → FIN-INVEST-030 | A copy |
| 26 | Black-Scholes | KEEP | `fin/invest` FIN-INVEST-018, -019 | Call and put |
| 27 | Bond Equivalent Yield | KEEP | `fin/invest` FIN-INVEST-009 | |
| 28 | Book Value per Share | KEEP | `fin/business` FIN-BIZ-026 | |
| 29 | Book Value per Share | MERGE | → FIN-BIZ-026 | Duplicate |
| 30 | Capital Asset Pricing Model | KEEP | `fin/invest` FIN-INVEST-016 | |
| 31 | Capital Asset Pricing Model | MERGE | → FIN-INVEST-016 | Duplicate |
| 32 | Capital Gains Yield | KEEP | `fin/invest` FIN-INVEST-012 | |
| 33 | Cash Flow Adequacy | KEEP | `fin/business` FIN-BIZ-009 | |
| 34 | CFOCE | KEEP | `fin/business` FIN-BIZ-011 | Operating cash flow to capital expenditures |
| 35 | CFOCL | KEEP | `fin/business` FIN-BIZ-010 | Operating cash flow to current liabilities |
| 36 | Compound Interest Future Value | KEEP | `fin/interest` FIN-INT-001 | |
| 37 | Consumer Price Index | KEEP | `fin/economics` FIN-ECON-001 | |
| 38 | Contribution Margin | KEEP | `fin/business` FIN-BIZ-040, -041 | Amount and ratio |
| 39 | Cost of Goods Available for Sale | DROP | — | One addition (beginning inventory + purchases) |
| 40 | Current Ratio | KEEP | `fin/business` FIN-BIZ-001 | |
| 41 | Current Yield | KEEP | `fin/invest` FIN-INVEST-001 | |
| 42 | Daily Salary | MOVE | `fin/personal` (proposed) | Salary row (F-4) |
| 43 | Days in Inventory | KEEP | `fin/business` FIN-BIZ-017 | |
| 44 | Debt Coverage Ratio | MERGE | → FIN-LOAN-008 | The same quantity as DSCR |
| 45 | Debt Ratio | KEEP | `fin/business` FIN-BIZ-003 | |
| 46 | Debt Service Coverage Ratio | KEEP | `fin/loan` FIN-LOAN-008 | |
| 47 | Debt to Equity Ratio | KEEP | `fin/business` FIN-BIZ-004 | |
| 48 | Diluted EPS | KEEP | `fin/business` FIN-BIZ-022 | |
| 49 | Discount | DROP | — | Ambiguous name (price discount or discount factor); the discount factor is FIN-INT-006 |
| 50 | Dividend Payout Ratio | KEEP | `fin/business` FIN-BIZ-028 | |
| 51 | Dividend Payout Ratio | MERGE | → FIN-BIZ-028 | Duplicate |
| 52 | Dividend Yield (Stock) | KEEP | `fin/invest` FIN-INVEST-013 | |
| 53 | Dividend Yield Ratio | MERGE | → FIN-INVEST-013 | Duplicate |
| 54 | Dividends Per Share | KEEP | `fin/business` FIN-BIZ-027 | |
| 55 | Earnings per Share | KEEP | `fin/business` FIN-BIZ-021 | |
| 56 | Earnings per Share | MERGE | → FIN-BIZ-021 | Duplicate |
| 57 | Economic Order Quantity | KEEP | `fin/business` FIN-BIZ-043 | |
| 58 | Economic Value Added | KEEP | `fin/business` FIN-BIZ-034 | |
| 59 | Effective Annual Rate | KEEP | `fin/interest` FIN-INT-029 | |
| 60 | End-of-Year Wealth | MERGE | → FIN-INT-001 | Wealth after one period of return (n = 1) |
| 61 | Equity Multiplier | KEEP | `fin/business` FIN-BIZ-005 | |
| 62 | Equivalent Annual Annuity | KEEP | `fin/invest` FIN-INVEST-024 | |
| 63 | Estimated Earnings | DROP | — | Ambiguous name; no formula identifiable from the name alone |
| 64 | FCFE | KEEP | `fin/business` FIN-BIZ-036 | |
| 65 | FCFF | KEEP | `fin/business` FIN-BIZ-035 | |
| 66 | Financial Leverage | KEEP | `fin/business` FIN-BIZ-039 | Degree of financial leverage |
| 67 | Fisher Equation | KEEP | `fin/economics` FIN-ECON-003 | Exact form |
| 68 | Fixed Charge Coverage Ratio | KEEP | `fin/business` FIN-BIZ-008 | |
| 69 | Future Value | MERGE | → FIN-INT-001 | Duplicate |
| 70 | Future Value | MERGE | → FIN-INT-001 | Duplicate |
| 71 | Future Value - (decimal input) | MERGE | → FIN-INT-001 | The same formula; rates are always decimal here |
| 72 | Future Value Factor | MERGE | → FIN-INT-005 | (F/P,i,n) |
| 73 | Future Value Factor | MERGE | → FIN-INT-005 | Duplicate |
| 74 | Future Value of Annuity | KEEP | `fin/interest` FIN-INT-014 | |
| 75 | Future Value [r,n,t] | MERGE | → FIN-INT-001 | Rate r/n per period over n·t periods: the caller's conversion |
| 76 | Future Value-Annuity Due | KEEP | `fin/interest` FIN-INT-015 | |
| 77 | FV of 401k Contributions | MOVE | `fin/personal` (proposed) | Retirement row (F-4); it reduces to FIN-INT-014 |
| 78 | Goodwill Multiple | DROP | — | No formula identifiable from the name; the goodwill computation is FIN-BIZ-048 |
| 79 | Goodwill MultipleEE | MERGE | → #78 | Duplicate of a dropped row |
| 80 | Gordon Growth Model | KEEP | `fin/invest` FIN-INVEST-014 | |
| 81 | Gross Profit Margin | KEEP | `fin/business` FIN-BIZ-018 | |
| 82 | Growing Annuity-Payment (FV) | KEEP | `fin/interest` FIN-INT-025 | |
| 83 | Growing Annuity-Payment (PV) | KEEP | `fin/interest` FIN-INT-024 | |
| 84 | Hamada's equation | KEEP | `fin/invest` FIN-INVEST-029 | |
| 85 | Herfindahl-Hirschman Index | KEEP | `fin/economics` FIN-ECON-004 | Market concentration (F-6) |
| 86 | Holding Cost per Day | DROP | — | Annual holding cost / days: a unit conversion of FIN-BIZ-043's input, not a function |
| 87 | Holding Period Yield | KEEP | `fin/invest` FIN-INVEST-010 | |
| 88 | How Much House Can I Afford | KEEP | `fin/loan` FIN-LOAN-010 | |
| 89 | Implied Repo Rate | KEEP | `fin/invest` FIN-INVEST-031 | Built only if sourced |
| 90 | Inflation-Adjusted Return | MERGE | → FIN-ECON-003 | The same exact Fisher relation |
| 91 | Interest Coverage Ratio | MERGE | → FIN-BIZ-007 | The same quantity as times interest earned |
| 92 | Investment Return Rate | MERGE | → FIN-INT-040 | Compound growth rate per period |
| 93 | Liquidity Ratio | MERGE | → FIN-BIZ-001 | Read as the current ratio, the usual liquidity ratio |
| 94 | Loan - Balloon Balance | KEEP | `fin/loan` FIN-LOAN-003 | |
| 95 | Loan to Deposit Ratio | KEEP | `fin/business` FIN-BIZ-050 | |
| 96 | Loan to Value (LTV) Ratio | KEEP | `fin/loan` FIN-LOAN-006 | |
| 97 | Money Market Yield | KEEP | `fin/invest` FIN-INVEST-008 | |
| 98 | Monthly Salary | MOVE | `fin/personal` (proposed) | Salary row (F-4) |
| 99 | Mortgage Loan Balance After n Payments | KEEP | `fin/loan` FIN-LOAN-002 | |
| 100 | Net Asset Value (NAV) | KEEP | `fin/invest` FIN-INVEST-028 | Per share |
| 101 | Net Income | DROP | — | One subtraction (revenue − expenses) |
| 102 | Net Income to Net Sales | MERGE | → FIN-BIZ-020 | The same quantity as net profit margin |
| 103 | Net Profit Margin | KEEP | `fin/business` FIN-BIZ-020 | |
| 104 | Net Working Capital | KEEP | `fin/business` FIN-BIZ-037 | Named by Devin; checked int64 subtraction |
| 105 | Operating Leverage | KEEP | `fin/business` FIN-BIZ-038 | Degree of operating leverage |
| 106 | Operating Profit Margin | KEEP | `fin/business` FIN-BIZ-019 | |
| 107 | Optimal stock cycle length | KEEP | `fin/business` FIN-BIZ-044 | EOQ / demand |
| 108 | Optimal Value of Variable Cost per Unit | DROP | — | No formula identifiable from the name |
| 109 | Owners' Equity to Creditors' Equity Ratio | KEEP | `fin/business` FIN-BIZ-006 | |
| 110 | Payback Period | KEEP | `fin/invest` FIN-INVEST-022 | |
| 111 | Perpetuity | KEEP | `fin/interest` FIN-INT-026 | |
| 112 | Personal Pension Calculator | MOVE | `fin/personal` (proposed) | Retirement row (F-4) |
| 113 | PITI Monthly Payment | KEEP | `fin/loan` FIN-LOAN-007 | |
| 114 | PITI Monthly Payment | MERGE | → FIN-LOAN-007 | Duplicate |
| 115 | PITI Plus Mortgage Insurance Monthly Payment | MERGE | → FIN-LOAN-007 | The mortgage-insurance term is a parameter (0 when absent) |
| 116 | Preferred Stock Value | MERGE | → FIN-INVEST-015 | The zero-growth formula D / r |
| 117 | Present Discounted Value | KEEP | `fin/interest` FIN-INT-002 | |
| 118 | Present Value Annuity Due | KEEP | `fin/interest` FIN-INT-017 | |
| 119 | Present Value Constant Growth Perpetual Annuity | MERGE | → FIN-INT-027 | Growing perpetuity |
| 120 | Present Value Factor | MERGE | → FIN-INT-006 | (P/F,i,n) |
| 121 | Present Value Factor | MERGE | → FIN-INT-006 | Duplicate |
| 122 | Present Value of Stock - Constant Growth | MERGE | → FIN-INVEST-014 | Gordon growth |
| 123 | Present Value of Stock - Zero Growth | KEEP | `fin/invest` FIN-INVEST-015 | |
| 124 | Present Value- Growing Annuity | KEEP | `fin/interest` FIN-INT-022 | |
| 125 | Present Value- With Continuous Compounding | KEEP | `fin/interest` FIN-INT-004 | |
| 126 | Present Value-Annuity | KEEP | `fin/interest` FIN-INT-016 | |
| 127 | Present Value-Annuity Due | MERGE | → FIN-INT-017 | Duplicate of #118 |
| 128 | Present Value-Annuity Factor | MERGE | → FIN-INT-009 | (P/A,i,n) |
| 129 | Present Value-Growing Perpetuity | KEEP | `fin/interest` FIN-INT-027 | |
| 130 | Price to Book Value | KEEP | `fin/business` FIN-BIZ-024 | |
| 131 | Price to Earnings Ratio (P/E) | KEEP | `fin/business` FIN-BIZ-023 | |
| 132 | Price to Sales Ratio | KEEP | `fin/business` FIN-BIZ-025 | |
| 133 | Price/Earnings Ratio | MERGE | → FIN-BIZ-023 | Duplicate |
| 134 | Rate of Inflation | KEEP | `fin/economics` FIN-ECON-002 | |
| 135 | RCSE | KEEP | `fin/business` FIN-BIZ-032 | Return on common stockholders' equity |
| 136 | Receivables Turnover Ratio | KEEP | `fin/business` FIN-BIZ-014 | |
| 137 | Remaining Balance-Amortization | MERGE | → FIN-LOAN-002 | The same quantity |
| 138 | Reorder Level | KEEP | `fin/business` FIN-BIZ-045 | Reorder point |
| 139 | Reorder Level | MERGE | → FIN-BIZ-045 | Duplicate |
| 140 | Retention Ratio | KEEP | `fin/business` FIN-BIZ-029 | |
| 141 | Retirement Savings | MOVE | `fin/personal` (proposed) | Retirement row (F-4) |
| 142 | Return On Assets | KEEP | `fin/business` FIN-BIZ-030 | |
| 143 | Return on Equity (ROE) | KEEP | `fin/business` FIN-BIZ-031 | |
| 144 | Return on Investment (ROI) | KEEP | `fin/business` FIN-BIZ-033 | |
| 145 | Risk Premium | KEEP | `fin/invest` FIN-INVEST-017 | |
| 146 | ROI | MERGE | → FIN-BIZ-033 | Duplicate |
| 147 | Rule of 72 | KEEP | `fin/interest` FIN-INT-033 | Plus the exact FIN-INT-034 |
| 148 | Single Payment Compound Amount (SPCA) Factor | KEEP | `fin/interest` FIN-INT-005 | |
| 149 | Single Payment Present Worth (SPPW) Factor | KEEP | `fin/interest` FIN-INT-006 | |
| 150 | Tax Equivalent Yield | KEEP | `fin/invest` FIN-INVEST-027 | |
| 151 | Terminal Value | MERGE | → FIN-INT-027 | The Gordon terminal value is a growing perpetuity of the next cash flow |
| 152 | Times Interest Earned Ratio | KEEP | `fin/business` FIN-BIZ-007 | |
| 153 | Total Stock Return | MERGE | → FIN-INVEST-010 | (P₁ − P₀ + D)/P₀, the holding period yield |
| 154 | Travel Distance on Gas Budget | MOVE | `life/automotive` | Not finance |
| 155 | Uniform Gradient Future Worth (UGFW) | KEEP | `fin/interest` FIN-INT-012 | |
| 156 | Uniform Gradient Present Worth (UGPW) | KEEP | `fin/interest` FIN-INT-011 | |
| 157 | Uniform Gradient Uniform Series (UGUS) | KEEP | `fin/interest` FIN-INT-013 | |
| 158 | Uniform Series Compound Amount (USCA) | KEEP | `fin/interest` FIN-INT-007 | |
| 159 | Uniform Series Present Worth (USPW) | KEEP | `fin/interest` FIN-INT-009 | |
| 160 | Uniform Series Sinking Fund (USSF) | KEEP | `fin/interest` FIN-INT-008 | |
| 161 | Variable Cost Per Unit Time | DROP | — | No formula identifiable from the name |
| 162 | Yield to Maturity | KEEP | `fin/invest` FIN-INVEST-003 | |
| 163 | Zero Coupon Bond Value | MERGE | → FIN-INT-002 | The price is the present value of the face amount |
| 164 | Zero Coupon Bond Yield | MERGE | → FIN-INT-040 | (F/P)^(1/n) − 1 |
| | **Data items, data sets, WikiClips** | | | |
| d1 | 7/6 SOFR ARM | DROP | — | DATA: a market rate |
| d2 | Current 15 year fixed interest rate | DROP | — | DATA: a market rate |
| d3 | Current 30 year fixed interest rate | DROP | — | DATA: a market rate |
| d4 | Mileage Rate - Charitable (2014) | DROP | — | DATA: an IRS rate for one year |
| d5 | work hours per week | DROP | — | DATA: a convention, not a function |
| d6 | work hours per year | DROP | — | DATA |
| d7 | work weeks per year | DROP | — | DATA |
| d8 | Assets and Liabilities | DROP | — | A data set |
| d9 | Current Mortgage Rates | DROP | — | A data set |
| d10 | Basic Investor Calculators | DROP | — | WikiClip |
| d11 | Financial Calculators | DROP | — | WikiClip |
| d12 | US Canada Trade Tools | DROP | — | WikiClip |

## Step 2: The standard set (not on the list)

Added from a corporate-finance formula sheet:
- **Interest:** FIN-INT-003 FutureValueContinuous, -010 CapitalRecovery, -023 FutureValueGrowingAnnuity,
  -030 EffectiveAnnualRateContinuous, -032 AnnualPercentageRate, -034 DoublingTime, -035 SimpleInterest,
  -036 YearFraction30360, -037 YearFractionActual365, -038 NetPresentValue, -039 InternalRateOfReturn,
  -041 AddOnInterestPayment.
- **Loan:** FIN-LOAN-001 AmortizationSchedule, -009 DebtToIncomeRatio.
- **Invest:** FIN-INVEST-002 BondPrice, -004 MacaulayDuration, -005 ModifiedDuration, -006 Convexity,
  -020 ModifiedInternalRateOfReturn, -021 ProfitabilityIndex, -023 DiscountedPaybackPeriod, -025 SharpeRatio,
  -026 WeightedAverageCostOfCapital.
- **Business:** FIN-BIZ-016 InventoryTurnover, -042 BreakEvenUnits, -046 DepreciationStraightLine,
  -047 DepreciationDecliningBalance.

## Step 3: Research entries

The entries below are the per-package research files, written by the package agents and checked by the lead.

## Build status (2026-11-12)

Built: **122 functions**. Final `go test ./fin/...` passes in all five packages (`../finance_2026-11-12/out/go_test_final.txt`).
Count per package: interest 40 · loan 10 · invest 27 · business 41 · economics 4.

Held back (D-015: the citation could not be confirmed, so the function was not built). Each research entry below says what was tried.
- **FIN-INT-028** FutureValueAnnuityContinuous: no source with discrete payments and continuous compounding.
- **FIN-INT-041** AddOnInterestPayment: no acceptable source opened.
- **FIN-INVEST-008** MoneyMarketYield: found only on prep sites and blogs.
- **FIN-INVEST-029** LeveredBetaHamada: the DOI is confirmed, but there was no readable text and no second source.
- **FIN-INVEST-031** ImpliedRepoRate: no source found.
- **FIN-BIZ-006** EquityToDebtRatio, **013** AssetToSalesRatio: no acceptable page defines them.
- **FIN-BIZ-009** CashFlowAdequacyRatio, **010** CashFlowToCurrentLiabilities, **011** CashFlowToCapitalExpenditures: only summaries and blogs, no formula.
- **FIN-BIZ-039** DegreeOfFinancialLeverage: defined in words only.
- **FIN-BIZ-044** OptimalOrderCycle: Harris's T is a lead time, not EOQ/D, and no source states EOQ/D.
- **FIN-BIZ-045** ReorderPoint: no formula in any source opened.

Removed before release, as duplicates. The IDs are never reused:
- **FIN-INVEST-024** EquivalentAnnualAnnuity → FIN-INT-018. The equivalent annual annuity is the annuity payment of the NPV.
- **FIN-BIZ-033** ReturnOnInvestment → FIN-INVEST-012. CapitalGainsYield computes the same (final − cost)/cost.

Added during verify: **FIN-INT-042** AnnuityRate. It solves PV = A·(P/A,i,n) + F·(P/F,i,n) for i, and YieldToMaturity and AnnualPercentageRate use it.

The vCalc rows that point at a held-back ID (#9, 33–35, 16, 66, 89, 97, 107, 138–139, 109, 84) are not covered until a source is found.

## Research: fin/interest


Agent A, 2026-11-12. Contract: `doc/finance_2026-11-12/BRIEF.md`. Code: `fin/interest/`.
Every citation below was confirmed on 2026-11-12 by opening a page that shows it. Some pages were opened by a
research sub-agent working for me; I re-checked the quoted worked numbers in Python.
Numerics: D. Goldberg, "What Every Computer Scientist Should Know About Floating-Point Arithmetic", ACM Computing Surveys 23(1):5–48, 1991, doi:10.1145/103162.103163 (DOI confirmed on Crossref, api.crossref.org/works/10.1145/103162.103163); the compound-interest example and Theorem 4 are under "Rounding Error" > "Cancellation" in the Oracle reprint https://docs.oracle.com/cd/E19957-01/806-3568/ncg_goldberg.html, which has no section numbers, so the brief's "§1.4" was not seen.
Root finder: W. H. Press, S. A. Teukolsky, W. T. Vetterling, B. P. Flannery, *Numerical Recipes*, 3rd ed., Cambridge University Press, 2007, ISBN 978-0-521-88068-8, §9.4 "Newton-Raphson Method Using Derivative", p. 456 (TOC at https://numerical.recipes/book.html; Crossref has no DOI for the book).

Built: 40 functions (001–027, 029–040 and 042, added at the lead's request). Held back: 028 and 041 (SOURCE NOT CONFIRMED for both).

### FutureValue      (ID: FIN-INT-001)
Formula: FV = PV·(1+i)ⁿ, where PV is the present value, i the rate per period and n the number of periods.
Inputs and units: money in int64 cents (|x| ≤ 2⁵³); rate decimal per period (finite, > −1); periods int
Output: FV in cents, rounded half-even to the cent
Errors: ErrOutOfRange, ErrRate (+ErrNotFinite for NaN/Inf), ErrPeriods (n < 0), ErrOverflow
Primary source: J. Dahlquist and R. Knight, *Principles of Finance*, OpenStax (Rice University, CC BY 4.0), 2022, §7.2 "Time Value of Money (TVM) Basics" (FV = PV × (1 + r)^n), https://openstax.org/books/principles-finance/pages/7-2-time-value-of-money-tvm-basics; second: R. A. DeFusco, D. W. McLeavey, J. E. Pinto, D. E. Runkle, *Quantitative Investment Analysis*, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ISBN 978-1-119-10422-3, ch. 1 "The Time Value of Money", §3 Eq. 2. CONFIRMED: page opened 2026-11-12.
Worked example: Principles of Finance §7.3 Example 1: PV $1,000, 3 %, N 4 → $1,125.51. §7.2: $1,000 at 4 % for 10 years → $1,480.24; $400,000 at 4 % for 6 years → $506,127.61. All three reproduce exactly.
Variants: Continuous form (003); factor form (005).
Best approach: exp(n·log1p(i)) beat pow(1+i, n) because it keeps the low bits of a small i (Goldberg). O(1) time, O(1) space, 0 allocations
Notes: No float money in or out. periods = 0 and PV = 0 return the input exactly.
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### PresentValue      (ID: FIN-INT-002)
Formula: PV = FV·(1+i)⁻ⁿ.
Inputs and units: money in int64 cents (|x| ≤ 2⁵³); rate decimal per period (finite, > −1); periods int
Output: PV in cents, rounded half-even to the cent
Errors: ErrOutOfRange, ErrRate, ErrPeriods, ErrOverflow (negative rate)
Primary source: J. Dahlquist and R. Knight, *Principles of Finance*, OpenStax (Rice University, CC BY 4.0), 2022, §7.4 "Applications of TVM in Finance", Eq. 7.32, https://openstax.org/books/principles-finance/pages/7-4-applications-of-tvm-in-finance; second: R. A. DeFusco, D. W. McLeavey, J. E. Pinto, D. E. Runkle, *Quantitative Investment Analysis*, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ISBN 978-1-119-10422-3, ch. 1 "The Time Value of Money", §5 Eq. 8. CONFIRMED: page opened 2026-11-12.
Worked example: Principles of Finance §7.3 Example 2: FV $50,000, 5 %, 5 → $39,176.31. §7.4: $1,000/(1.05)^30 → $231.38.
Variants: Continuous form (004); factor form (006).
Best approach: exp(−n·log1p(i)). O(1) time, O(1) space, 0 allocations
Notes: A long horizon rounds to 0 cents; that is not an error.
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### FutureValueContinuous      (ID: FIN-INT-003)
Formula: FV = PV·e^(r·t), where r is the nominal annual rate compounded continuously and t the time in years.
Inputs and units: PV in cents; ratePerYear decimal (finite, > −1); years float64 ≥ 0
Output: FV in cents, rounded half-even to the cent
Errors: ErrOutOfRange, ErrRate, ErrNotFinite (years), ErrNegative (years < 0), ErrOverflow
Primary source: R. A. DeFusco, D. W. McLeavey, J. E. Pinto, D. E. Runkle, *Quantitative Investment Analysis*, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ISBN 978-1-119-10422-3, ch. 1 "The Time Value of Money", §3.2 "Continuous Compounding", Eq. 4 FV_N = PV·e^(r_s N), publisher excerpt https://catalogimages.wiley.com/images/db/pdf/9781119104223.excerpt.pdf (chapter 1 in full); second: G. Strang and E. Herman, *Calculus Volume 1*, OpenStax, 2016, §6.8 "Exponential Growth and Decay", https://openstax.org/books/calculus-volume-1/pages/6-8-exponential-growth-and-decay ("Balance = Pe^(rt)"). CONFIRMED: page opened 2026-11-12. Principles of Finance does not have it (its index has no entry for continuous compounding).
Worked example: QIA Example 6: $10,000 at 8 % for 2 years → $11,735.11 (exact). QIA Table 1: $1 → $1.083287.
Variants: —
Best approach: math.Exp(r·t). O(1) time, O(1) space, 0 allocations
Notes: The rate must be > −1 to match the package policy, although e^(rt) is defined for any r.
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### PresentValueContinuous      (ID: FIN-INT-004)
Formula: PV = FV·e^(−r·t).
Inputs and units: FV in cents; ratePerYear; years ≥ 0
Output: PV in cents, rounded half-even to the cent
Errors: ErrOutOfRange, ErrRate, ErrNotFinite, ErrNegative, ErrOverflow
Primary source: G. Strang and E. Herman, *Calculus Volume 1*, OpenStax, 2016, §6.8 "Exponential Growth and Decay", https://openstax.org/books/calculus-volume-1/pages/6-8-exponential-growth-and-decay, Example 6.43; second: R. A. DeFusco, D. W. McLeavey, J. E. Pinto, D. E. Runkle, *Quantitative Investment Analysis*, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ISBN 978-1-119-10422-3, ch. 1 "The Time Value of Money" §3.2 Eq. 4, solved for PV. CONFIRMED: page opened 2026-11-12.
Worked example: Calculus Vol. 1 Example 6.43: $1,000,000 due in 40 years at 5 % → $135,335.28; at 6 % → $90,717.95 (both exact).
Variants: —
Best approach: math.Exp(−r·t). O(1) time, O(1) space, 0 allocations
Notes: —
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### SinglePaymentCompoundAmount      (ID: FIN-INT-005)
Formula: (F/P,i,n) = (1+i)ⁿ, where i is the rate per period and n the number of periods.
Inputs and units: ratePerPeriod decimal (finite, > −1); periods int
Output: float64 factor, unrounded
Errors: ErrRate, ErrPeriods (n < 0, or n < 1 for A/F, A/P and A/G), ErrOverflow (factor beyond float64; not for A/F, A/P or A/G, which cannot overflow: see the review round)
Primary source: NCEES, *FE Reference Handbook* 10.0.1 (2nd printing May 2020, ISBN 978-1-947801-11-0), Engineering Economics, factor table p. 230; copy hosted by Mississippi State University https://www.cee.msstate.edu/wp-content/uploads/fe-handbook-10-0-1.pdf, row "Single Payment Compound Amount". CONFIRMED: page opened 2026-11-12. The handbook PDF was read page by page. NCEES forbids reproducing the tables, so only the formulas are cited and one table value is used as a check.
Worked example: FE table p. 233, i = 1.00 %, n = 10 → 1.1046 (4 decimals; the test allows 5e-5).
Variants: —
Best approach: exp(n·log1p(i)); n = 0 gives exactly 1. O(1) time, O(1) space, 0 allocations. The tests check against a 256-bit math/big oracle and against the handbook identities.
Notes: Factors are float64 because they are dimensionless multipliers, not money.
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### SinglePaymentPresentWorth      (ID: FIN-INT-006)
Formula: (P/F,i,n) = (1+i)⁻ⁿ, where i is the rate per period and n the number of periods.
Inputs and units: ratePerPeriod decimal (finite, > −1); periods int
Output: float64 factor, unrounded
Errors: ErrRate, ErrPeriods (n < 0, or n < 1 for A/F, A/P and A/G), ErrOverflow (factor beyond float64; not for A/F, A/P or A/G, which cannot overflow: see the review round)
Primary source: NCEES, *FE Reference Handbook* 10.0.1 (2nd printing May 2020, ISBN 978-1-947801-11-0), Engineering Economics, factor table p. 230; copy hosted by Mississippi State University https://www.cee.msstate.edu/wp-content/uploads/fe-handbook-10-0-1.pdf, row "Single Payment Present Worth". CONFIRMED: page opened 2026-11-12. The handbook PDF was read page by page. NCEES forbids reproducing the tables, so only the formulas are cited and one table value is used as a check.
Worked example: FE table p. 233, i = 1.00 %, n = 10 → 0.9053 (4 decimals; the test allows 5e-5).
Variants: —
Best approach: exp(−n·log1p(i)); underflow to 0 is not an error. O(1) time, O(1) space, 0 allocations. The tests check against a 256-bit math/big oracle and against the handbook identities.
Notes: Factors are float64 because they are dimensionless multipliers, not money.
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### UniformSeriesCompoundAmount      (ID: FIN-INT-007)
Formula: (F/A,i,n) = ((1+i)ⁿ − 1)/i, where i is the rate per period and n the number of periods.
Inputs and units: ratePerPeriod decimal (finite, > −1); periods int
Output: float64 factor, unrounded
Errors: ErrRate, ErrPeriods (n < 0, or n < 1 for A/F, A/P and A/G), ErrOverflow (factor beyond float64; not for A/F, A/P or A/G, which cannot overflow: see the review round)
Primary source: NCEES, *FE Reference Handbook* 10.0.1 (2nd printing May 2020, ISBN 978-1-947801-11-0), Engineering Economics, factor table p. 230; copy hosted by Mississippi State University https://www.cee.msstate.edu/wp-content/uploads/fe-handbook-10-0-1.pdf, row "Uniform Series Compound Amount". CONFIRMED: page opened 2026-11-12. The handbook PDF was read page by page. NCEES forbids reproducing the tables, so only the formulas are cited and one table value is used as a check.
Worked example: FE table p. 233, i = 1.00 %, n = 10 → 10.4622 (4 decimals; the test allows 5e-5).
Variants: —
Best approach: expm1(x)/x · n · log1p(i)/i with x = n·log1p(i); no cancellation, and i = 0 gives exactly n. O(1) time, O(1) space, 0 allocations. The tests check against a 256-bit math/big oracle and against the handbook identities.
Notes: Factors are float64 because they are dimensionless multipliers, not money.
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### UniformSeriesSinkingFund      (ID: FIN-INT-008)
Formula: (A/F,i,n) = i/((1+i)ⁿ − 1), where i is the rate per period and n the number of periods.
Inputs and units: ratePerPeriod decimal (finite, > −1); periods int
Output: float64 factor, unrounded
Errors: ErrRate, ErrPeriods (n < 0, or n < 1 for A/F, A/P and A/G), ErrOverflow (factor beyond float64; not for A/F, A/P or A/G, which cannot overflow: see the review round)
Primary source: NCEES, *FE Reference Handbook* 10.0.1 (2nd printing May 2020, ISBN 978-1-947801-11-0), Engineering Economics, factor table p. 230; copy hosted by Mississippi State University https://www.cee.msstate.edu/wp-content/uploads/fe-handbook-10-0-1.pdf, row "Uniform Series Sinking Fund". CONFIRMED: page opened 2026-11-12. The handbook PDF was read page by page. NCEES forbids reproducing the tables, so only the formulas are cited and one table value is used as a check.
Worked example: FE table p. 233, i = 1.00 %, n = 10 → 0.0956 (4 decimals; the test allows 5e-5).
Variants: —
Best approach: 1/(F/A); n ≥ 1; i = 0 gives 1/n. O(1) time, O(1) space, 0 allocations. The tests check against a 256-bit math/big oracle and against the handbook identities.
Notes: Factors are float64 because they are dimensionless multipliers, not money.
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### UniformSeriesPresentWorth      (ID: FIN-INT-009)
Formula: (P/A,i,n) = ((1+i)ⁿ − 1)/(i(1+i)ⁿ), where i is the rate per period and n the number of periods.
Inputs and units: ratePerPeriod decimal (finite, > −1); periods int
Output: float64 factor, unrounded
Errors: ErrRate, ErrPeriods (n < 0, or n < 1 for A/F, A/P and A/G), ErrOverflow (factor beyond float64; not for A/F, A/P or A/G, which cannot overflow: see the review round)
Primary source: NCEES, *FE Reference Handbook* 10.0.1 (2nd printing May 2020, ISBN 978-1-947801-11-0), Engineering Economics, factor table p. 230; copy hosted by Mississippi State University https://www.cee.msstate.edu/wp-content/uploads/fe-handbook-10-0-1.pdf, row "Uniform Series Present Worth". CONFIRMED: page opened 2026-11-12. The handbook PDF was read page by page. NCEES forbids reproducing the tables, so only the formulas are cited and one table value is used as a check.
Worked example: FE table p. 233, i = 1.00 %, n = 10 → 9.4713 (4 decimals; the test allows 5e-5).
Variants: —
Best approach: expm1(−x)/(−x) · n · log1p(i)/i; i = 0 gives n. O(1) time, O(1) space, 0 allocations. The tests check against a 256-bit math/big oracle and against the handbook identities.
Notes: Factors are float64 because they are dimensionless multipliers, not money.
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### CapitalRecovery      (ID: FIN-INT-010)
Formula: (A/P,i,n) = i(1+i)ⁿ/((1+i)ⁿ − 1), where i is the rate per period and n the number of periods.
Inputs and units: ratePerPeriod decimal (finite, > −1); periods int
Output: float64 factor, unrounded
Errors: ErrRate, ErrPeriods (n < 0, or n < 1 for A/F, A/P and A/G), ErrOverflow (factor beyond float64; not for A/F, A/P or A/G, which cannot overflow: see the review round)
Primary source: NCEES, *FE Reference Handbook* 10.0.1 (2nd printing May 2020, ISBN 978-1-947801-11-0), Engineering Economics, factor table p. 230; copy hosted by Mississippi State University https://www.cee.msstate.edu/wp-content/uploads/fe-handbook-10-0-1.pdf, row "Capital Recovery". CONFIRMED: page opened 2026-11-12. The handbook PDF was read page by page. NCEES forbids reproducing the tables, so only the formulas are cited and one table value is used as a check.
Worked example: FE table p. 233, i = 1.00 %, n = 10 → 0.1056 (4 decimals; the test allows 5e-5).
Variants: —
Best approach: 1/(P/A); n ≥ 1; i = 0 gives 1/n. O(1) time, O(1) space, 0 allocations. The tests check against a 256-bit math/big oracle and against the handbook identities.
Notes: Factors are float64 because they are dimensionless multipliers, not money.
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### UniformGradientPresentWorth      (ID: FIN-INT-011)
Formula: (P/G,i,n) = ((1+i)ⁿ − 1)/(i²(1+i)ⁿ) − n/(i(1+i)ⁿ), where i is the rate per period and n the number of periods.
Inputs and units: ratePerPeriod decimal (finite, > −1); periods int
Output: float64 factor, unrounded
Errors: ErrRate, ErrPeriods (n < 0, or n < 1 for A/F, A/P and A/G), ErrOverflow (factor beyond float64; not for A/F, A/P or A/G, which cannot overflow: see the review round)
Primary source: NCEES, *FE Reference Handbook* 10.0.1 (2nd printing May 2020, ISBN 978-1-947801-11-0), Engineering Economics, factor table p. 230; copy hosted by Mississippi State University https://www.cee.msstate.edu/wp-content/uploads/fe-handbook-10-0-1.pdf, row "Uniform Gradient Present Worth". CONFIRMED: page opened 2026-11-12. The handbook PDF was read page by page. NCEES forbids reproducing the tables, so only the formulas are cited and one table value is used as a check.
Worked example: FE table p. 233, i = 1.00 %, n = 10 → 41.8435 (4 decimals; the test allows 5e-5).
Variants: —
Best approach: (F/G)·(1+i)⁻ⁿ when n·ln(1+i) ≤ 1, ((P/A) − n(1+i)⁻ⁿ)/i above that; i = 0 gives n(n−1)/2. O(1) time, O(1) space, 0 allocations. The tests check against a 256-bit math/big oracle and against the handbook identities.
Notes: Factors are float64 because they are dimensionless multipliers, not money.
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### UniformGradientFutureWorth      (ID: FIN-INT-012)
Formula: (F/G,i,n) = ((1+i)ⁿ − 1)/i² − n/i, where i is the rate per period and n the number of periods.
Inputs and units: ratePerPeriod decimal (finite, > −1); periods int
Output: float64 factor, unrounded
Errors: ErrRate, ErrPeriods (n < 0, or n < 1 for A/F, A/P and A/G), ErrOverflow (factor beyond float64; not for A/F, A/P or A/G, which cannot overflow: see the review round)
Primary source: NCEES, *FE Reference Handbook* 10.0.1 (2nd printing May 2020, ISBN 978-1-947801-11-0), Engineering Economics, factor table p. 230; copy hosted by Mississippi State University https://www.cee.msstate.edu/wp-content/uploads/fe-handbook-10-0-1.pdf, row "Uniform Gradient Future Worth". CONFIRMED: page opened 2026-11-12. The handbook PDF was read page by page. NCEES forbids reproducing the tables, so only the formulas are cited and one table value is used as a check.
Worked example: FE table p. 233, i = 1.00 %, n = 10 → 46.2213 (the table has no F/G column; this is (10.4622 − 10)/0.01 from the handbook footnote F/G = (F/A − n)/i) (4 decimals; the test allows 5e-5).
Variants: —
Best approach: ((eˣ−1−x)/x²)·(nL/i)² + n(L − i)/i², with L = log1p(i) and x = nL; both bracketed functions use a Taylor series near 0, so there is no cancellation. O(1) time, O(1) space, 0 allocations. The tests check against a 256-bit math/big oracle and against the handbook identities.
Notes: Factors are float64 because they are dimensionless multipliers, not money.
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### UniformGradientUniformSeries      (ID: FIN-INT-013)
Formula: (A/G,i,n) = 1/i − n/((1+i)ⁿ − 1), where i is the rate per period and n the number of periods.
Inputs and units: ratePerPeriod decimal (finite, > −1); periods int
Output: float64 factor, unrounded
Errors: ErrRate, ErrPeriods (n < 0, or n < 1 for A/F, A/P and A/G), ErrOverflow (factor beyond float64; not for A/F, A/P or A/G, which cannot overflow: see the review round)
Primary source: NCEES, *FE Reference Handbook* 10.0.1 (2nd printing May 2020, ISBN 978-1-947801-11-0), Engineering Economics, factor table p. 230; copy hosted by Mississippi State University https://www.cee.msstate.edu/wp-content/uploads/fe-handbook-10-0-1.pdf, row "Uniform Gradient Uniform Series". CONFIRMED: page opened 2026-11-12. The handbook PDF was read page by page. NCEES forbids reproducing the tables, so only the formulas are cited and one table value is used as a check.
Worked example: FE table p. 233, i = 1.00 %, n = 10 → 4.4179 (4 decimals; the test allows 5e-5).
Variants: —
Best approach: (F/G)/(F/A) when x ≤ 1, the direct form above that; i = 0 gives (n−1)/2. O(1) time, O(1) space, 0 allocations. The tests check against a 256-bit math/big oracle and against the handbook identities.
Notes: Factors are float64 because they are dimensionless multipliers, not money.
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### FutureValueAnnuity      (ID: FIN-INT-014)
Formula: FV = A·((1+i)ⁿ − 1)/i (ordinary annuity: payments at the end of each period), where A is the level payment, i the rate per period and n the number of payments.
Inputs and units: money in int64 cents (|x| ≤ 2⁵³); rate decimal per period (finite, > −1); periods int
Output: cents, rounded half-even to the cent
Errors: ErrOutOfRange, ErrRate, ErrPeriods (n < 0), ErrOverflow
Primary source: R. A. DeFusco, D. W. McLeavey, J. E. Pinto, D. E. Runkle, *Quantitative Investment Analysis*, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ISBN 978-1-119-10422-3, ch. 1 "The Time Value of Money", §4.1 Eq. 7; second: J. Dahlquist and R. Knight, *Principles of Finance*, OpenStax (Rice University, CC BY 4.0), 2022, §8.2 "Annuities" Eq. 8.23, https://openstax.org/books/principles-finance/pages/8-2-annuities. CONFIRMED: page opened 2026-11-12.
Worked example: QIA Example 7: €20,000 a year, 30 years, 9 % → €2,726,150.77 (exact). QIA Figure 3: $1,000, 5 years, 5 % → $5,525.63 (exact). Principles of Finance §8.2 shows $16,248.98 for $3,000, 5 years, 4 %; the exact value is $16,248.9677, which rounds to $16,248.97. The book rounds its factor first, and the test documents this.
Variants: Ordinary annuity vs annuity due (×(1+i)).
Best approach: Multiply or divide the amount by the cancellation-free factor helper (no pow − 1). O(1) time, O(1) space, 0 allocations.
Notes: —
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### FutureValueAnnuityDue      (ID: FIN-INT-015)
Formula: FV = A·((1+i)ⁿ − 1)/i·(1+i), where A is the level payment, i the rate per period and n the number of payments.
Inputs and units: money in int64 cents (|x| ≤ 2⁵³); rate decimal per period (finite, > −1); periods int
Output: cents, rounded half-even to the cent
Errors: ErrOutOfRange, ErrRate, ErrPeriods (n < 0), ErrOverflow
Primary source: J. Dahlquist and R. Knight, *Principles of Finance*, OpenStax (Rice University, CC BY 4.0), 2022, §8.2 Eq. 8.25 ("multiply by one additional period of interest, (1 + i)"), https://openstax.org/books/principles-finance/pages/8-2-annuities. CONFIRMED: page opened 2026-11-12.
Worked example: HAND-CHECKED: $1,000, 5 %, 5 → $5,525.63125 × 1.05 = $5,801.91.
Variants: Ordinary annuity vs annuity due (×(1+i)).
Best approach: Multiply or divide the amount by the cancellation-free factor helper (no pow − 1). O(1) time, O(1) space, 0 allocations.
Notes: —
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### PresentValueAnnuity      (ID: FIN-INT-016)
Formula: PV = A·(1 − (1+i)⁻ⁿ)/i, where A is the level payment, i the rate per period and n the number of payments.
Inputs and units: money in int64 cents (|x| ≤ 2⁵³); rate decimal per period (finite, > −1); periods int
Output: cents, rounded half-even to the cent
Errors: ErrOutOfRange, ErrRate, ErrPeriods (n < 0), ErrOverflow
Primary source: R. A. DeFusco, D. W. McLeavey, J. E. Pinto, D. E. Runkle, *Quantitative Investment Analysis*, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ISBN 978-1-119-10422-3, ch. 1 "The Time Value of Money", §6.1 Eq. 11; second: J. Dahlquist and R. Knight, *Principles of Finance*, OpenStax (Rice University, CC BY 4.0), 2022, §8.2 Eq. 8.14. CONFIRMED: page opened 2026-11-12.
Worked example: QIA Example 11: €1,000 a year, 5 years, 12 % → €3,604.78 (exact). Not used: Principles of Finance's $99,817.81 for $25,000, 5 years, 8 %, which is exactly $99,817.75.
Variants: Ordinary annuity vs annuity due (×(1+i)).
Best approach: Multiply or divide the amount by the cancellation-free factor helper (no pow − 1). O(1) time, O(1) space, 0 allocations.
Notes: —
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### PresentValueAnnuityDue      (ID: FIN-INT-017)
Formula: PV = A·(1 − (1+i)⁻ⁿ)/i·(1+i), where A is the level payment, i the rate per period and n the number of payments.
Inputs and units: money in int64 cents (|x| ≤ 2⁵³); rate decimal per period (finite, > −1); periods int
Output: cents, rounded half-even to the cent
Errors: ErrOutOfRange, ErrRate, ErrPeriods (n < 0), ErrOverflow
Primary source: R. A. DeFusco, D. W. McLeavey, J. E. Pinto, D. E. Runkle, *Quantitative Investment Analysis*, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ISBN 978-1-119-10422-3, ch. 1 "The Time Value of Money", §6.1 footnote 7; second: J. Dahlquist and R. Knight, *Principles of Finance*, OpenStax (Rice University, CC BY 4.0), 2022, §8.2 Eq. 8.17. CONFIRMED: page opened 2026-11-12.
Worked example: QIA §6.1: four $100 payments, the first today, at 12 % → $340.18 (exact).
Variants: Ordinary annuity vs annuity due (×(1+i)).
Best approach: Multiply or divide the amount by the cancellation-free factor helper (no pow − 1). O(1) time, O(1) space, 0 allocations.
Notes: —
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### AnnuityPaymentFromPresentValue      (ID: FIN-INT-018)
Formula: A = PV·i/(1 − (1+i)⁻ⁿ) = PV·(A/P,i,n), where A is the level payment, i the rate per period and n the number of payments.
Inputs and units: money in int64 cents (|x| ≤ 2⁵³); rate decimal per period (finite, > −1); periods int
Output: cents, rounded half-even to the cent
Errors: ErrOutOfRange, ErrRate, ErrPeriods (n < 1), ErrOverflow
Primary source: R. A. DeFusco, D. W. McLeavey, J. E. Pinto, D. E. Runkle, *Quantitative Investment Analysis*, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ISBN 978-1-119-10422-3, ch. 1 "The Time Value of Money", §7.3 "Solving for the Size of Annuity Payments" (A = PV / PV annuity factor); second: J. Dahlquist and R. Knight, *Principles of Finance*, OpenStax (Rice University, CC BY 4.0), 2022, §8.2 Eq. 8.30 and §8.3 "Loan Amortization" Eq. 8.31, https://openstax.org/books/principles-finance/pages/8-3-loan-amortization. CONFIRMED: page opened 2026-11-12.
Worked example: QIA Example 20: $100,000, 8 %/12, 360 → $733.76. Principles of Finance §8.3: $32,000, 6 %/12, 36 → $973.50; $140,000, 3.6 %/12, 240 → $819.16. All three are exact. §8.2's $16,376.60 comes from a rounded factor; the exact value $16,376.5929 rounds to $16,376.59.
Variants: Ordinary annuity vs annuity due (×(1+i)).
Best approach: Multiply or divide the amount by the cancellation-free factor helper (no pow − 1). O(1) time, O(1) space, 0 allocations.
Notes: The payment is rounded once, at the end. Schedules in fin/loan absorb the leftover cents in the last payment.
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### AnnuityPaymentFromFutureValue      (ID: FIN-INT-019)
Formula: A = FV·i/((1+i)ⁿ − 1) = FV·(A/F,i,n), where A is the level payment, i the rate per period and n the number of payments.
Inputs and units: money in int64 cents (|x| ≤ 2⁵³); rate decimal per period (finite, > −1); periods int
Output: cents, rounded half-even to the cent
Errors: ErrOutOfRange, ErrRate, ErrPeriods (n < 1)
Primary source: NCEES, *FE Reference Handbook* 10.0.1 (2nd printing May 2020, ISBN 978-1-947801-11-0), Engineering Economics, factor table p. 230; copy hosted by Mississippi State University https://www.cee.msstate.edu/wp-content/uploads/fe-handbook-10-0-1.pdf, Uniform Series Sinking Fund; second: R. A. DeFusco, D. W. McLeavey, J. E. Pinto, D. E. Runkle, *Quantitative Investment Analysis*, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ISBN 978-1-119-10422-3, ch. 1 "The Time Value of Money", §4.1 Eq. 7 solved for A. CONFIRMED: page opened 2026-11-12.
Worked example: HAND-CHECKED, the inverse of QIA Figure 3: $5,525.63, 5 %, 5 → $1,000.00.
Variants: Ordinary annuity vs annuity due (×(1+i)).
Best approach: Multiply or divide the amount by the cancellation-free factor helper (no pow − 1). O(1) time, O(1) space, 0 allocations.
Notes: |payment| ≤ |amount|, because the factor's sum starts with a term of 1, so ErrOverflow cannot occur.
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### AnnuityDuePaymentFromPresentValue      (ID: FIN-INT-020)
Formula: A = PV·i/((1 − (1+i)⁻ⁿ)(1+i)), where A is the level payment, i the rate per period and n the number of payments.
Inputs and units: money in int64 cents (|x| ≤ 2⁵³); rate decimal per period (finite, > −1); periods int
Output: cents, rounded half-even to the cent
Errors: ErrOutOfRange, ErrRate, ErrPeriods (n < 1)
Primary source: R. A. DeFusco, D. W. McLeavey, J. E. Pinto, D. E. Runkle, *Quantitative Investment Analysis*, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ISBN 978-1-119-10422-3, ch. 1 "The Time Value of Money", §6.1 footnote 7 solved for A; second: J. Dahlquist and R. Knight, *Principles of Finance*, OpenStax (Rice University, CC BY 4.0), 2022, §8.2 Eq. 8.17. CONFIRMED: page opened 2026-11-12.
Worked example: HAND-CHECKED, the inverse of QIA §6.1: $340.18, 12 %, 4 → $100.00.
Variants: Ordinary annuity vs annuity due (×(1+i)).
Best approach: Multiply or divide the amount by the cancellation-free factor helper (no pow − 1). O(1) time, O(1) space, 0 allocations.
Notes: |payment| ≤ |amount|, because the factor's sum starts with a term of 1, so ErrOverflow cannot occur.
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### AnnuityDuePaymentFromFutureValue      (ID: FIN-INT-021)
Formula: A = FV·i/(((1+i)ⁿ − 1)(1+i)), where A is the level payment, i the rate per period and n the number of payments.
Inputs and units: money in int64 cents (|x| ≤ 2⁵³); rate decimal per period (finite, > −1); periods int
Output: cents, rounded half-even to the cent
Errors: ErrOutOfRange, ErrRate, ErrPeriods (n < 1), ErrOverflow
Primary source: J. Dahlquist and R. Knight, *Principles of Finance*, OpenStax (Rice University, CC BY 4.0), 2022, §8.2 Eq. 8.25 solved for PYMT; second: R. A. DeFusco, D. W. McLeavey, J. E. Pinto, D. E. Runkle, *Quantitative Investment Analysis*, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ISBN 978-1-119-10422-3, ch. 1 "The Time Value of Money", §4.1 Eq. 7. CONFIRMED: page opened 2026-11-12.
Worked example: HAND-CHECKED: $5,801.91, 5 %, 5 → $1,000.00.
Variants: Ordinary annuity vs annuity due (×(1+i)).
Best approach: Multiply or divide the amount by the cancellation-free factor helper (no pow − 1). O(1) time, O(1) space, 0 allocations.
Notes: —
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### PresentValueGrowingAnnuity      (ID: FIN-INT-022)
Formula: PV = C·(1 − ((1+g)/(1+r))ⁿ)/(r − g), and C·n/(1+r) when r = g. Damodaran writes C = A(1+g), where A is the current cash flow.
Inputs and units: money in int64 cents (|x| ≤ 2⁵³); rate decimal per period (finite, > −1); periods int; g is the growth per period of the payments; the first payment C is at the end of period 1; the r = g limit is handled exactly
Output: cents, rounded half-even to the cent
Errors: ErrOutOfRange, ErrRate (r or g), ErrPeriods, ErrOverflow
Primary source: A. Damodaran, *Applied Corporate Finance*, 3rd ed., Wiley, 2010 (ISBN 978-0-470-38464-0), Appendix 3 "Time Value of Money", author-hosted https://pages.stern.nyu.edu/~adamodar/pdfiles/acf3E/book/app3.pdf, "Growing Annuities" (the formula and the r = g case "nA"). CONFIRMED: page opened 2026-11-12: the author-hosted PDF of the book appendix was read in full.
Worked example: Damodaran's gold mine: 5,000 oz × $300 growing 3 %, discounted at 10 %, for 20 years → $16,145,980 to the dollar. The exact value is $16,145,979.98, and the test uses cents.
Variants: FV (023); payment (024, 025).
Best approach: With u = (r−g)/(1+g) and ℓ = log1p(u): C·n·expm1(−nℓ)/(−nℓ)·(ℓ/u)/(1+g), which has no cancellation near r = g. O(1) time, O(1) space, 0 allocations
Notes: Not in Principles of Finance (its index has no growing annuity) and not in QIA ch. 1. Damodaran is a recognised Wiley textbook but is NOT on the brief's named list; see FOR THE LEAD.
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### FutureValueGrowingAnnuity      (ID: FIN-INT-023)
Formula: FV = C·((1+r)ⁿ − (1+g)ⁿ)/(r − g) = PV_growing·(1+r)ⁿ.
Inputs and units: money in int64 cents (|x| ≤ 2⁵³); rate decimal per period (finite, > −1); periods int; g is the growth per period of the payments; the first payment C is at the end of period 1; the r = g limit is handled exactly
Output: cents, rounded half-even to the cent
Errors: ErrOutOfRange, ErrRate, ErrPeriods, ErrOverflow
Primary source: A. Damodaran, *Applied Corporate Finance*, 3rd ed., Wiley, 2010 (ISBN 978-0-470-38464-0), Appendix 3 "Time Value of Money", author-hosted https://pages.stern.nyu.edu/~adamodar/pdfiles/acf3E/book/app3.pdf for the PV, carried forward by FV = PV(1+r)^n (R. A. DeFusco, D. W. McLeavey, J. E. Pinto, D. E. Runkle, *Quantitative Investment Analysis*, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ISBN 978-1-119-10422-3, ch. 1 "The Time Value of Money" §3 Eq. 2). CONFIRMED: page opened 2026-11-12.
Worked example: HAND-CHECKED in exact rational arithmetic: the gold mine → $108,622,079.49.
Variants: —
Best approach: PV helper × exp(n·log1p(r)). O(1) time, O(1) space, 0 allocations
Notes: Same flag as 022.
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### GrowingAnnuityPaymentFromPresentValue      (ID: FIN-INT-024)
Formula: C = PV·(r − g)/(1 − ((1+g)/(1+r))ⁿ).
Inputs and units: money in int64 cents (|x| ≤ 2⁵³); rate decimal per period (finite, > −1); periods int; g is the growth per period of the payments; the first payment C is at the end of period 1; the r = g limit is handled exactly
Output: first payment in cents, rounded half-even to the cent
Errors: ErrOutOfRange, ErrRate, ErrPeriods (n < 1), ErrOverflow
Primary source: A. Damodaran, *Applied Corporate Finance*, 3rd ed., Wiley, 2010 (ISBN 978-0-470-38464-0), Appendix 3 "Time Value of Money", author-hosted https://pages.stern.nyu.edu/~adamodar/pdfiles/acf3E/book/app3.pdf, the PV formula solved for the first payment. CONFIRMED: page opened 2026-11-12.
Worked example: HAND-CHECKED: $16,145,979.98, 10 %, 3 %, 20 → $1,545,000.00.
Variants: —
Best approach: PV / helper. O(1) time, O(1) space, 0 allocations
Notes: Same flag as 022.
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### GrowingAnnuityPaymentFromFutureValue      (ID: FIN-INT-025)
Formula: C = FV·(r − g)/((1+r)ⁿ − (1+g)ⁿ).
Inputs and units: money in int64 cents (|x| ≤ 2⁵³); rate decimal per period (finite, > −1); periods int; g is the growth per period of the payments; the first payment C is at the end of period 1; the r = g limit is handled exactly
Output: first payment in cents, rounded half-even to the cent
Errors: ErrOutOfRange, ErrRate, ErrPeriods (n < 1), ErrOverflow
Primary source: A. Damodaran, *Applied Corporate Finance*, 3rd ed., Wiley, 2010 (ISBN 978-0-470-38464-0), Appendix 3 "Time Value of Money", author-hosted https://pages.stern.nyu.edu/~adamodar/pdfiles/acf3E/book/app3.pdf, the PV formula × (1+r)^n, solved for the first payment. CONFIRMED: page opened 2026-11-12.
Worked example: HAND-CHECKED: $108,622,079.49 → $1,545,000.00.
Variants: —
Best approach: FV / (helper·exp(n·log1p(r))). O(1) time, O(1) space, 0 allocations
Notes: Same flag as 022.
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### PresentValuePerpetuity      (ID: FIN-INT-026)
Formula: PV = C/r, with the first payment one period from now; it converges only for r > 0.
Inputs and units: payment in cents; rate decimal
Output: cents, rounded half-even to the cent
Errors: ErrOutOfRange, ErrRate, ErrDomain (r ≤ 0), ErrOverflow
Primary source: J. Dahlquist and R. Knight, *Principles of Finance*, OpenStax (Rice University, CC BY 4.0), 2022, §8.1 "Perpetuities" Eq. 8.1 PV = C/Rs, https://openstax.org/books/principles-finance/pages/8-1-perpetuities; second: R. A. DeFusco, D. W. McLeavey, J. E. Pinto, D. E. Runkle, *Quantitative Investment Analysis*, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ISBN 978-1-119-10422-3, ch. 1 "The Time Value of Money", §6.2 Eq. 13. CONFIRMED: page opened 2026-11-12.
Worked example: §8.1 Eq. 8.2: $2.00/0.07 = $28.57. Endowment: $75,000/0.035 = $2,142,857(.14). QIA Example 14: £100/0.05 = £2,000. All exact.
Variants: Growing perpetuity (027).
Best approach: One division. O(1) time, O(1) space, 0 allocations
Notes: —
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### PresentValueGrowingPerpetuity      (ID: FIN-INT-027)
Formula: PV = C/(r − g); requires r > g. This is also the terminal value and the Gordon growth price.
Inputs and units: first payment in cents; r and g decimals
Output: cents, rounded half-even to the cent
Errors: ErrOutOfRange, ErrRate, ErrDomain (r ≤ g), ErrOverflow
Primary source: J. Dahlquist and R. Knight, *Principles of Finance*, OpenStax (Rice University, CC BY 4.0), 2022, §8.1 Eq. 8.3 PV = C/(Rs − G), https://openstax.org/books/principles-finance/pages/8-1-perpetuities; second: A. Damodaran, *Applied Corporate Finance*, 3rd ed., Wiley, 2010 (ISBN 978-0-470-38464-0), Appendix 3 "Time Value of Money", author-hosted https://pages.stern.nyu.edu/~adamodar/pdfiles/acf3E/book/app3.pdf, "Growing Perpetuities". CONFIRMED: page opened 2026-11-12.
Worked example: §8.1 Rooney: $1.75/(0.058 − 0.012) = $38.04. Savo: $4.00/(0.03 − 0.002) = $142.86. With zero growth, $1.75/0.058 = $30.17.
Variants: —
Best approach: One division. O(1) time, O(1) space, 0 allocations
Notes: fin/invest GordonGrowthPrice (FIN-INVEST-014) can call this.
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### FutureValueAnnuityContinuous      (ID: FIN-INT-028)
Formula: (proposed) FV = C·(e^(r·n) − 1)/(e^r − 1), for end-of-year payments compounded continuously.
Inputs and units: payment in cents; ratePerYear decimal; years int.
Output: cents.
Errors: —
Primary source: SOURCE NOT CONFIRMED. Tried and not found:
- Principles of Finance: its index has no continuous compounding.
- QIA ch. 1: Eq. 4 covers only the continuous FV of a lump sum.
- NCEES FE Handbook 10.0.1 pp. 230–232: no continuous-compounding factors.
- Calculus Vol. 1 §6.8: lump sum only.
- Hull ch. 4: the O'Reilly page returned 403, and copies of the slides are not acceptable.
- Newnan: the publisher's pages could not be reached.
The brief says "build only if sourced".
Worked example: —
Variants: —
Best approach: it would be C·expm1(r·n)/expm1(r), which has no cancellation.
Notes: NOT BUILT. Summing a geometric series from Eq. 4 would derive it, but the result would be a formula with no citable line.
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### EffectiveAnnualRate      (ID: FIN-INT-029)
Formula: EAR = (1 + r/m)^m − 1, where r is the nominal (stated) annual rate and m the number of compounding periods per year.
Inputs and units: nominalRatePerYear decimal (finite); periodsPerYear int ≥ 1; r/m > −1
Output: decimal per year, unrounded
Errors: ErrNotFinite, ErrPeriods (m < 1), ErrRate (r/m ≤ −1), ErrOverflow
Primary source: NCEES, *FE Reference Handbook* 10.0.1 (2nd printing May 2020, ISBN 978-1-947801-11-0), Engineering Economics, p. 231 "Non-Annual Compounding" i_e = (1 + r/m)^m − 1; copy hosted by Mississippi State University https://www.cee.msstate.edu/wp-content/uploads/fe-handbook-10-0-1.pdf; second: R. A. DeFusco, D. W. McLeavey, J. E. Pinto, D. E. Runkle, *Quantitative Investment Analysis*, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ISBN 978-1-119-10422-3, ch. 1 "The Time Value of Money", §3.3 Eq. 5. CONFIRMED: page opened 2026-11-12.
Worked example: Principles of Finance §8.4 Eq. 8.38: 18 % compounded monthly gives (1.015)^12 − 1 = 0.19562. QIA §3.3: 8 % semiannual → 0.0816. QIA Table 1, daily → 0.083278.
Variants: Continuous form (030). APY (031) is the deposit-regulation form.
Best approach: expm1(m·log1p(r/m)), which is Goldberg's own compound-interest example. O(1) time, O(1) space, 0 allocations
Notes: Principles of Finance §8.4 has only the worked form and no general formula, so the formula is cited from NCEES and QIA.
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### EffectiveAnnualRateContinuous      (ID: FIN-INT-030)
Formula: EAR = e^r − 1.
Inputs and units: nominalRatePerYear decimal (finite, > −1)
Output: decimal per year
Errors: ErrRate, ErrOverflow
Primary source: R. A. DeFusco, D. W. McLeavey, J. E. Pinto, D. E. Runkle, *Quantitative Investment Analysis*, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ISBN 978-1-119-10422-3, ch. 1 "The Time Value of Money", §3.3 Eq. 6 EAR = e^(r_s) − 1. CONFIRMED: page opened 2026-11-12.
Worked example: QIA §3.3: 8 % → 8.33 % (0.083287 in Table 1).
Variants: —
Best approach: math.Expm1. O(1) time, O(1) space, 0 allocations
Notes: —
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### AnnualPercentageYield      (ID: FIN-INT-031)
Formula: APY = (1 + Interest/Principal)^(365/Days in term) − 1. The regulation writes 100[…] to give percent; this function returns a decimal.
Inputs and units: principalCents > 0; interestCents > −principal; daysInTerm ≥ 1
Output: decimal, unrounded. 12 CFR 1030.3(f)(1) rounds the disclosed APY to 0.01 percentage point; that rounding is left to the caller.
Errors: ErrOutOfRange, ErrDomain (principal ≤ 0, or interest ≤ −principal), ErrPeriods (days < 1), ErrOverflow
Primary source: 12 CFR Part 1030 (Regulation DD) Appendix A Part I.A, read on the CFPB's official eRegulations copy https://www.consumerfinance.gov/rules-policy/regulations/1030/a/ ; the rounding rule is §1030.3(f)(1), at .../1030/3/. CONFIRMED: page opened 2026-11-12. ecfr.gov, the site the brief names, redirects automated readers to unblock.federalregister.gov, so the CFPB copy of the same CFR text was used; see FOR THE LEAD.
Worked example: Part I.A: $61.68 on $1,000 over 365 days → 6.17 %; $30.37 on $1,000 over 182 days → 6.18 %. Part I.B: $26.68 on $1,000 over 183 days → 5.39 %. All three reproduce after rounding.
Variants: Part I.E and Part II (APY earned) are not covered.
Best approach: expm1((365/d)·log1p(I/P)). O(1) time, O(1) space, 0 allocations
Notes: —
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### AnnualPercentageRate      (ID: FIN-INT-032)
Formula: The unit-period rate i solves A = P·(1 − (1+i)⁻ⁿ)/i, and APR = i·w, where w is the number of unit-periods per year. This follows App. J (b)(1) and (c)(1) with t = 1, f = 0.
Inputs and units: amountFinancedCents > 0; paymentCents > 0; payments n ≥ 1; periodsPerYear w ≥ 1
Output: decimal per year, unrounded
Errors: ErrOutOfRange, ErrDomain, ErrPeriods, ErrNoSolution
Primary source: 12 CFR Part 1026 (Regulation Z) Appendix J (b)(1), (b)(4), the (b)(8) general equation, and (c)(1), read on https://www.consumerfinance.gov/rules-policy/regulations/1026/j/ . The equations are images; the research agent read the (c)(1) image at https://files.consumerfinance.gov/eregulations/ER22DE11.064-1.png. CONFIRMED: page opened 2026-11-12. Solver: W. H. Press, S. A. Teukolsky, W. T. Vetterling, B. P. Flannery, *Numerical Recipes*, 3rd ed., Cambridge University Press, 2007, ISBN 978-0-521-88068-8, §9.4 "Newton-Raphson Method Using Derivative", p. 456 (TOC at https://numerical.recipes/book.html; Crossref has no DOI for the book).
Worked example: (c)(1) Example (i): A $5,000, P $230, n 24, monthly, t = 1, f = 0 → 9.69 % (computed: 9.6857 %). Not used: the (b)(9) iteration example as summarised ($1,000 repaid by six payments of $176.51 → 12.83 %). It does NOT reproduce (it gives 19.97 %), so the summary was garbled.
Variants: Not covered: an odd first period (f ≠ 0), irregular payments, multiple advances.
Best approach: The shared bracket + rtsafe solver on f(i) = P·(P/A) − A, with an analytic derivative. O(1) per evaluation, ≤ 272 evaluations, 0 allocations: the closure does not escape (checked with go build -gcflags=-m).
Notes: The regulatory tolerance (§1026.22) is not applied; the caller rounds.
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### RuleOf72      (ID: FIN-INT-033)
Formula: N ≈ 72/(100·i).
Inputs and units: ratePerPeriod decimal > 0
Output: periods (approximate)
Errors: ErrRate, ErrDomain (i ≤ 0), ErrOverflow
Primary source: J. Dahlquist and R. Knight, *Principles of Finance*, OpenStax (Rice University, CC BY 4.0), 2022, §7.4, subheading "The Rule of 72", Eq. 7.37, https://openstax.org/books/principles-finance/pages/7-4-applications-of-tvm-in-finance; second: R. A. DeFusco, D. W. McLeavey, J. E. Pinto, D. E. Runkle, *Quantitative Investment Analysis*, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ISBN 978-1-119-10422-3, ch. 1 "The Time Value of Money", §7.2 footnote 9. CONFIRMED: page opened 2026-11-12.
Worked example: §7.4: 9 % → 8 years; 6 % → 12; 4 % → 18. QIA fn. 9: 7 % → 10.3.
Variants: The exact form is 034.
Best approach: One division. O(1) time, O(1) space, 0 allocations
Notes: An approximation by definition. The test checks it against DoublingTime (within 2 % for rates of 4–11 %).
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### DoublingTime      (ID: FIN-INT-034)
Formula: N = ln 2/ln(1+i).
Inputs and units: ratePerPeriod decimal > 0
Output: periods, fractional
Errors: ErrRate, ErrDomain, ErrOverflow
Primary source: R. A. DeFusco, D. W. McLeavey, J. E. Pinto, D. E. Runkle, *Quantitative Investment Analysis*, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ISBN 978-1-119-10422-3, ch. 1 "The Time Value of Money", §7.2 "Solving for the Number of Periods", Example 19 (N = ln(FV/PV)/ln(1 + r) with FV/PV = 2). CONFIRMED: page opened 2026-11-12.
Worked example: QIA Example 19: 7 % → 10.24 years.
Variants: —
Best approach: math.Ln2/math.Log1p(i). O(1) time, O(1) space, 0 allocations
Notes: Not rounded up to a whole period.
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### SimpleInterest      (ID: FIN-INT-035)
Formula: I = P·r·t.
Inputs and units: principalCents; ratePerYear; yearFraction ≥ 0
Output: cents, half-even
Errors: ErrOutOfRange, ErrRate, ErrNotFinite, ErrNegative, ErrOverflow
Primary source: D. Kirk, *Contemporary Mathematics*, OpenStax, 2023, §6.3 "Simple Interest", https://openstax.org/books/contemporary-mathematics/pages/6-3-simple-interest. CONFIRMED: page opened 2026-11-12. Principles of Finance does not have it (not in its index). MSRB G-33(a) has an accrued-interest formula, but it is an image that was not read, so it is not cited.
Worked example: Example 6.27: $4,000, 5.5 %, 4 years → $880.00. Example 6.28: $11,995, 14.9 %, 2 years → $3,574.51. Example 6.30: $5,600, 11.75 %, 45/365 → $81.1233; half-even gives $81.12, and the book rounds up to $81.13 (interest is "almost uniformly rounded up").
Variants: The day counts 036 and 037 supply t.
Best approach: One product. O(1) time, O(1) space, 0 allocations
Notes: Contemporary Mathematics is an OpenStax book but is not on F-1's named list; see FOR THE LEAD. The book rounds up, which differs from the package's half-even policy; the test documents the difference.
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### YearFraction30360      (ID: FIN-INT-036)
Formula: Days = (Y2−Y1)·360 + (M2−M1)·30 + (D2−D1). If D2 = 31 and D1 is 30 or 31, D2 becomes 30; then if D1 = 31, D1 becomes 30. t = Days/360.
Inputs and units: start and end time.Time, each read in its own Location; time of day is ignored
Output: years (decimal)
Errors: ErrDomain (end before start)
Primary source: MSRB Rule G-33(e)(i) "Day Count Basis" and (e)(ii) "Day Count Formula", https://www.msrb.org/Rules-and-Interpretations/MSRB-Rules/General/Rule-G-33 (amended per Release 34-77316, 2016). CONFIRMED: page opened 2026-11-12: the formula and both rules for the 31st were quoted. The brief's locator "(b)(i)(B)" is wrong; the day count is in (e).
Worked example: HAND-CHECKED: 2026-01-31 → 2026-03-31 gives 60/360. The test also has a table of 31st-of-month and February cases.
Variants: No end-of-February rule. The rule text has none. The SIA "30/360 US" February rule was not confirmed on msrb.org and is not applied.
Best approach: Integer arithmetic. O(1) time, O(1) space, 0 allocations
Notes: —
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### YearFractionActual365      (ID: FIN-INT-037)
Formula: t = (actual calendar days from start to end)/365.
Inputs and units: start and end time.Time
Output: years (decimal)
Errors: ErrDomain (end before start)
Primary source: D. Kirk, *Contemporary Mathematics*, OpenStax, 2023, §6.3 "Simple Interest", https://openstax.org/books/contemporary-mathematics/pages/6-3-simple-interest ("This may be referred to as Actual/365"; Example 6.30); second: Regulation DD Appendix A Part I.A, which divides the actual days in the term by 365. CONFIRMED: page opened 2026-11-12.
Worked example: Example 6.30: 45 days → 45/365, which gives interest of $81.12 half-even.
Variants: —
Best approach: Calendar dates are converted to UTC-midnight Unix days, so daylight saving cannot shift the count. O(1) time, O(1) space, 0 allocations
Notes: Same F-1 flag as 035.
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### NetPresentValue      (ID: FIN-INT-038)
Formula: NPV = Σₜ cₜ/(1+i)ᵗ for t = 0…n−1; c₀ is not discounted.
Inputs and units: ratePerPeriod; cashFlowsCents []int64 (each |c| ≤ 2⁵³)
Output: cents, rounded half-even to the cent
Errors: ErrRate, ErrDomain (empty), ErrOutOfRange, ErrOverflow
Primary source: J. Dahlquist and R. Knight, *Principles of Finance*, OpenStax (Rice University, CC BY 4.0), 2022, §16.2 "Net Present Value (NPV) Method" Eq. 16.1, https://openstax.org/books/principles-finance/pages/16-2-net-present-value-npv-method; second: R. A. DeFusco, D. W. McLeavey, J. E. Pinto, D. E. Runkle, *Quantitative Investment Analysis*, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ISBN 978-1-119-10422-3, ch. 1 "The Time Value of Money", §6.4. CONFIRMED: page opened 2026-11-12.
Worked example: §16.2 Table 16.5: −30,000, then 6,000 … 18,000, at 9 % → $26,946.90 (exact). Table 16.3 shows $2,835.63 because it adds rounded rows; the unrounded $2,835.6194 rounds to $2,835.62, and the test documents this.
Variants: —
Best approach: Horner's rule in v = 1/(1+i): O(n) time, O(1) space, 0 allocations, which meets the Ω(n) lower bound for reading the input. Rejected: an exp per term, the same order but about 10× slower.
Notes: The relative error grows about n·ε; the tests allow that much on cancelling sums.
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### InternalRateOfReturn      (ID: FIN-INT-039)
Formula: r such that Σ cₜ/(1+r)ᵗ = 0.
Inputs and units: cashFlowsCents []int64; n ≥ 2; at least one inflow and one outflow
Output: decimal per period
Errors: ErrDomain, ErrOutOfRange, ErrNoSolution, ErrNoConvergence
Primary source: J. Dahlquist and R. Knight, *Principles of Finance*, OpenStax (Rice University, CC BY 4.0), 2022, §16.3 "Internal Rate of Return (IRR) Method" (the definition; Eq. 16.3; multiple IRRs in Table 16.8), https://openstax.org/books/principles-finance/pages/16-3-internal-rate-of-return-irr-method. Method: W. H. Press, S. A. Teukolsky, W. T. Vetterling, B. P. Flannery, *Numerical Recipes*, 3rd ed., Cambridge University Press, 2007, ISBN 978-0-521-88068-8, §9.4 "Newton-Raphson Method Using Derivative", p. 456 (TOC at https://numerical.recipes/book.html; Crossref has no DOI for the book). CONFIRMED: page opened 2026-11-12.
Worked example: §16.3 Table 16.7: the embroidery machine's flows → 14.09 % (computed: 14.0855 %).
Variants: Flows with several sign changes can have several IRRs. This returns the root bracketed nearest 0, upward first when both directions bracket one at the same scan step.
Best approach: Scan outward from 0: upward, doubling from 0.01 to 1e6; downward to −0.64, then halving the gap to −1 + 1e-12. Then rtsafe. Horner computes NPV and dNPV/dr together. O(n·k) with k ≤ 72 + 200; 0 allocations.
Notes: ErrNoConvergence cannot be reached through the public brackets (at most 41 bisections are ever needed), so it is tested on the solver directly. The book does not state the Olympic example's two IRRs; the test computes them from the quadratic.
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### GrowthRatePerPeriod      (ID: FIN-INT-040)
Formula: g = (FV/PV)^(1/N) − 1.
Inputs and units: presentValueCents and futureValueCents (non-zero, same sign); periods ≥ 1
Output: decimal per period
Errors: ErrOutOfRange, ErrDomain, ErrPeriods
Primary source: R. A. DeFusco, D. W. McLeavey, J. E. Pinto, D. E. Runkle, *Quantitative Investment Analysis*, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ISBN 978-1-119-10422-3, ch. 1 "The Time Value of Money", §7.1 "Solving for Interest Rates and Growth Rates", Eq. 14. CONFIRMED: page opened 2026-11-12.
Worked example: QIA Example 17: 10,503.0 → 14,146.4 over 4 periods → 0.077291. Example 18: 8.52 → 7.35 over 5 → −0.029111. Principles of Finance §7.3: $1,000 → $1,125.51 in 4 → 3 %.
Variants: Also serves as CAGR, zero-coupon yield and implied return.
Best approach: expm1(log(FV/PV)/N). O(1) time, O(1) space, 0 allocations
Notes: —
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### AddOnInterestPayment      (ID: FIN-INT-041)
Formula: (proposed) payment = (P + P·r·t)/payments.
Inputs and units: principal in cents; ratePerYear; years; payments.
Output: cents.
Errors: —
Primary source: SOURCE NOT CONFIRMED. Searched and not found:
- Principles of Finance (its index);
- Contemporary Mathematics §6.3 and §6.4;
- Principles of Accounting (search);
- any federalreserve.gov or consumerfinance.gov page.
The only textbook hit was COMAP, *For All Practical Purposes* 10e, ch. 22, at digfir-published.macmillanusa.com; the site did not resolve (DNS failure). Blogs (CFI, Bankrate) are not acceptable sources.
Worked example: —
Variants: —
Best approach: one product and one division, O(1).
Notes: NOT BUILT.
Covers rows: the lead's vCalc row list is not in the tree (no doc/catalog finance file on 2026-11-12); see FOR THE LEAD.

### AnnuityRate      (ID: FIN-INT-042)
Formula: Find i with PV = PMT·(1 − (1+i)⁻ⁿ)/i + FV·(1+i)⁻ⁿ, where all amounts are ≥ 0 and payments fall at the end of each period. This is the ordinary-annuity PV (QIA Eq. 11) plus the single-sum PV (QIA Eq. 8), solved for the rate.
Inputs and units: presentValueCents > 0; paymentCents ≥ 0 and futureValueCents ≥ 0, not both 0 (all |x| ≤ 2⁵³); periods ≥ 1, any int.
Output: rate per period as a decimal, unrounded.
Errors: ErrOutOfRange, ErrNegative, ErrDomain (PV = 0, or PMT = FV = 0), ErrPeriods (n < 1), ErrNoSolution (root outside (−1 + 1e-12, 1e6]). ErrNoConvergence cannot occur through this function: the widest bracket needs at most 41 bisections.
Primary source: R. A. DeFusco, D. W. McLeavey, J. E. Pinto, D. E. Runkle, *Quantitative Investment Analysis*, 3rd ed., Wiley (CFA Institute Investment Series), 2015, ch. 1, §6.1 Eq. 11 and §5 Eq. 8 for the equation, and §7.1 "Solving for Interest Rates and Growth Rates" for solving it for r. Second: J. Dahlquist and R. Knight, *Principles of Finance*, OpenStax, 2022, §7.3 "Methods for Solving Time Value of Money Problems" (solving for I/Y), https://openstax.org/books/principles-finance/pages/7-3-methods-for-solving-time-value-of-money-problems. Solver: Press et al., *Numerical Recipes* 3rd ed., §9.4. CONFIRMED: pages opened 2026-11-12. Principles of Finance §7.3 solves for the rate with a calculator only (its "Solving for Interest Rate" case); it has no algebraic method, which is why the solver is cited separately.
Worked example: Principles of Finance §7.3: PV $1,000, FV $1,125.51, N 4 → 3 %. QIA Example 20 reversed: $100,000 repaid by 360 payments of $733.76 → 0.08/12. HAND-CHECKED: a par bond (PV = FV = $1,000, PMT $25, n = 20) → 2.5 %.
Variants: AnnualPercentageRate (032) is this with FV = 0, times the number of periods per year, and now calls the same internal solver. fin/invest YieldToMaturity can call it as price = coupon·(P/A) + face·(P/F).
Best approach: the shared bracket-plus-rtsafe solver on the closed forms. Each evaluation is O(1) whatever n is (exp/log1p/expm1, no loop over periods); there are at most 272 evaluations and 0 allocations (the closure does not escape; checked with -gcflags=-m). Rejected: building the cash-flow slice and calling InternalRateOfReturn, which is O(n) per step, allocates n + 1 values, and cannot handle huge n.
Notes: With non-negative amounts the right side falls strictly as i rises, so the root is unique. A term whose coefficient is 0 is skipped, so 0·Inf never makes a NaN. The tests cover periods = math.MaxInt.
Covers rows: the lead's vCalc row list is not in the tree; see FOR THE LEAD.


## Review round (lead rulings, 2026-11-12)

- F-1 covers every OpenStax title, and Damodaran *Applied Corporate Finance* 3e is accepted, so the sources
  of 004, 022–025, 035 and 037 stand.
- F-5 accuracy contract: every cents result's doc comment (001–004, 014–027, 035, 038) now says: "The
  result is the float64 value rounded half-even; it equals the correctly rounded exact value except when
  that value lies within about n·ε·|x| of a half cent, and is within 1 cent of it for |result| ≤ 2⁴³ cents."
- Oracle tests: magnitudes are drawn ≤ 2⁴³. A 1-cent difference is accepted only when the exact value
  lies within max(1e-6 cent, 1e-14·|x|) of a half cent; everywhere else equality is exact. The 1e-14·|x|
  term is the closed forms' own float error (≈ (2|n·ln(1+i)| + 3)·ε ≤ 3.4e-15 relative here). With a
  bare 1e-6 cent a correct implementation would fail on values above about 10⁹ cents.
- 008 UniformSeriesSinkingFund: ErrOverflow is removed. (F/A) ≥ 1, so (A/F) ≤ 1.
- 013 UniformGradientUniformSeries: ErrOverflow is removed. (A/G) is a weighted mean of 0 … n−1.
- 010 CapitalRecovery: ErrOverflow is removed. The factor is now evaluated as i + (A/F) for i ≥ 0 (the
  handbook identity) and as 1/(P/A) for i < 0, so it is at most 1 + i and stays finite even at
  math.MaxFloat64.
- 032 AnnualPercentageRate: ErrNoConvergence is removed from the doc comment, because it cannot be
  reached. The solver-level test stays.

## FOR THE LEAD

1. Some functions cite sources that are not on the brief's named list, because nothing on the list covers
   their formula. sources.md accepts any recognised textbook, but the brief's list may be meant as closed;
   that is your call. The sources:
   - OpenStax *Contemporary Mathematics* (2023) §6.3, for 035 SimpleInterest and 037 YearFractionActual365;
   - OpenStax *Calculus Volume 1* (2016) §6.8, the primary source for 004 and the second source for 003;
   - Damodaran, *Applied Corporate Finance* 3e (Wiley 2010) App. 3, for 022–025 (growing annuities).
   If F-1 does not extend to them, 022–025, 035 and 037 lose their primary source and should be marked
   unconfirmed. 003 and 004 still stand on QIA (CFA Institute) Eq. 4.
2. DeFusco et al., *Quantitative Investment Analysis* 3e (CFA Institute Investment Series) is the CFA Institute
   reading used as the primary source for 003, 014, 016–018, 030, 034 and 040. The brief lists "CFA
   Institute readings" as acceptable.
3. ecfr.gov blocks automated reading (it answers with a 302 to unblock.federalregister.gov). 12 CFR 1030
   App. A and 12 CFR 1026 App. J were read on the CFPB's official eRegulations copies instead.
4. Two locators in the brief need correcting:
   - The MSRB 30/360 day count is in G-33(e), not (b)(i)(B).
   - Goldberg's compound-interest example and Theorem 4 sit under "Rounding Error" > "Cancellation"; no
     "§1.4" is visible in the reprint.
5. Some of the books' worked numbers are a cent off because the books round intermediate values. The tests
   use the exact half-even cent and say so:
   - Principles of Finance §8.2: $16,248.98 (exact .97);
   - Principles of Finance §8.2 Eq. 8.30: $16,376.60 (exact .59);
   - Principles of Finance §16.2 Table 16.3: $2,835.63 (exact .62);
   - Contemporary Mathematics rounds interest up: $81.13 against $81.12 half-even.
6. These are yours:
   - INDEX.md rows;
   - doc/sources/ entries (OpenStax Principles of Finance, QIA, NCEES FE 10.0.1, Damodaran, Contemporary
     Mathematics, Calculus Vol. 1, Reg DD, Reg Z, MSRB G-33, Goldberg, Numerical Recipes);
   - doc/catalog coverage. No finance "Covers rows" list existed in the tree, so that line is unfilled
     in every entry.
7. The tests and benchmarks are written but were never run, per the brief.
8. Functions the other agents will call:
   - 018 AnnuityPaymentFromPresentValue
   - 010 CapitalRecovery
   - 038 NetPresentValue
   - 039 InternalRateOfReturn
   - 027 PresentValueGrowingPerpetuity (Gordon)
   - 040 GrowthRatePerPeriod

## Research: fin/loan


Date read: 2026-11-12. Every URL below was opened on that date and the quoted text read on the page.
Common source for the schedule mechanics (S-AM):
Dahlquist, J. & Knight, R., *Principles of Finance*, OpenStax (Rice University), 2022, §8.3 "Loan
Amortization", Table 8.4 (auto loan) and the 20-year mortgage table,
https://openstax.org/books/principles-finance/pages/8-3-loan-amortization — CONFIRMED (opened 2026-11-12).
Auto loan: $32,000, 6 %/yr monthly (i = 0.005), 36 payments, "approximately $973.50 per month"; rows
1: 160.00 / 813.50 / 31,186.50; 2: 155.93 / 817.57 / 30,368.93; 3: 151.84 / 821.66 / 29,547.27;
36: 4.84 interest. Mortgage: $140,000, 3.6 %/yr (i = 0.003), 240 payments, "$819.16"; row 239: 4.89 / 814.27;
row 240: 2.45 interest; "After three years (month 36), you receive a bonus of $2,000" applied to principal.
Annuity PV (S-PV): same book §8.2 "Annuities", eq. 8.14 "PVa = PYMT × [1 − 1/(1 + r)^n] / r",
https://openstax.org/books/principles-finance/pages/8-2-annuities — CONFIRMED.
Lump-sum + annuity (S-BOND): same book §10.2 "Bond Valuation", eq. 10.9 "Bond Price = Par Value × 1/(1+r)^n +
Coupon Payment Amount × (1−1/(1+r)^n)/r", https://openstax.org/books/principles-finance/pages/10-2-bond-valuation
— CONFIRMED.

### AmortizationSchedule      (ID: FIN-LOAN-001)
Formula: A = round_he(P·i/(1−(1+i)⁻ⁿ)) via interest.AnnuityPaymentFromPresentValue (FIN-INT-018). Row k:
I_k = round_he(B_{k−1}·i); principal_k = A − I_k; B_k = B_{k−1} − principal_k. Row n: principal_n = B_{n−1},
payment_n = I_n + B_{n−1}, B_n = 0. Review ruling: a row whose B + I ≤ A is the payoff row and the schedule
ends there (fewer rows; capacity n; never a negative balance). Bound for i ≥ 0: |last − A| ≤ (F/A,i,n) cents.
principal 0 → empty non-nil slice.
Inputs and units: principalCents ≥ 0; ratePerPeriod decimal > −1; periods 1…maxSchedulePeriods (10⁷).
Output: []AmortizationRow, len = periods, one allocation.
Errors: ErrOutOfRange, ErrNegative, ErrRate (+ErrNotFinite), ErrPeriods, ErrOverflow.
Primary source: S-AM (Table 8.4 rows). Payment: S-PV. CONFIRMED.
Worked example: (3200000, 0.005, 36): payment 97350; rows 1–3 exactly as the book (16000/81350/3118650,
15593/81757/3036893, 15184/82166/2954727); row 36 interest 484 (book 4.84). Our row 36 pays 97360 and
principal 96876 because the rounded payment's shortfall is absorbed in the last payment; the book's 968.66
carries the unrounded payment. Mortgage (14000000, 0.003, 240): payment 81916; row 239 489/81427 = book;
row 240 interest 245 = book. HAND-CHECKED in Python with the same float operations.
Oracle: Σprincipal = P, B_n = 0, Σpayments = P + Σinterest, RemainingBalance(k) = row k.
Best approach: one pass O(n), one allocation of n rows. Lower bound Ω(n) (the output).
Notes: half-even per row (brief numeric policy); never append.
Covers rows: "Amortization schedule", "Loan amortization", "Mortgage amortization".

### RemainingBalance      (ID: FIN-LOAN-002)
Formula: B_k by the same row recursion as FIN-LOAN-001 (k = paymentsMade; k = n gives 0).
Inputs: principalCents, ratePerPeriod, periods (1…10⁷), paymentsMade (0…periods); 0 after an early payoff.
Output: int64 cents. Errors: as 001, plus ErrPeriods when paymentsMade ∉ [0, periods].
Primary source: S-AM (balance column of Table 8.4). CONFIRMED.
Worked example: (3200000, 0.005, 36, 3) → 2954727 (book 29,547.27).
Best approach: O(k) iteration, no allocation; the closed form B_k = P(1+i)^k − A·((1+i)^k−1)/i is O(1) but
does not reproduce the per-row rounding, so it would disagree with the schedule by cents (rejected).
Covers rows: "Remaining balance", "Loan balance".

### BalloonBalance      (ID: FIN-LOAN-003)
Formula: B_k = B_{k−1} + round_he(B_{k−1}·i) − paymentCents, k = 1…paymentsMade; no last-row adjustment.
Inputs: principalCents ≥ 0, paymentCents ≥ 0, ratePerPeriod, paymentsMade ≥ 0.
Output: int64 cents; 0 once a payment covers balance + interest (never negative; review ruling 2026-11-12).
Errors: ErrOutOfRange, ErrNegative, ErrRate, ErrPeriods, ErrOverflow (balance leaves ±2⁵³ under negative amortization).
Primary source: CFPB, "What is a balloon payment? When is one allowed?",
https://www.consumerfinance.gov/ask-cfpb/what-is-a-balloon-payment-when-is-one-allowed-en-104/ — "A balloon
payment on a mortgage is a large, one-time payment at the end of the loan term." CONFIRMED. Mechanics: S-AM.
Worked example: 30-year $100,000 at 0.5 %/month, payment 59955, 60 payments → 9305437 (HAND-CHECKED).
Best approach: O(k), no allocation (closed form rejected for the same reason as 002).
Covers rows: "Balloon balance".

### BalloonLoanPayment      (ID: FIN-LOAN-004)
Formula: P = A·(P/A,i,n) + B·(P/F,i,n) (S-BOND eq. 10.9 with coupon → payment, par → balloon), so
A = (P − B·(1+i)⁻ⁿ)·(A/P,i,n), rounded half-even once. Uses interest.SinglePaymentPresentWorth (FIN-INT-006)
and interest.CapitalRecovery (FIN-INT-010).
Errors: ErrOutOfRange, ErrNegative, ErrRate, ErrPeriods, ErrDomain (balloon > P(1+i)ⁿ: payment would be
negative), ErrOverflow.
Primary source: S-BOND + S-PV (CONFIRMED); balloon definition: CFPB page above (CONFIRMED).
Worked example: (10000000, 5000000, 0.005, 60) → 121664; balloon 0 → equals AnnuityPaymentFromPresentValue;
balloon = principal → interest only P·i. Oracle: BalloonBalance(P, A, i, n) within n cents of B. HAND-CHECKED.
Covers rows: "Balloon loan payment".

### ExtraPaymentPayoff      (ID: FIN-LOAN-005)
Formula: scheduled A as 001; each period I_k = round_he(B·i); pay A + extra; when B + I_k ≤ A + extra (or k = n)
the loan is paid off with B + I_k. interestSavedCents = Σinterest(schedule 001) − Σinterest(with extra).
Errors: ErrOutOfRange, ErrNegative, ErrRate, ErrPeriods, ErrOverflow.
Primary source: S-AM (the book's extra $2,000 applied to principal pays the mortgage off early). CONFIRMED.
Worked example: (14000000, 0.003, 240, 10000) → 204 periods, 925723 cents saved; extra 0 → periods, 0 saved.
HAND-CHECKED in Python.
Best approach: O(n) time (baseline pass + extra pass), O(1) space, no allocation.
Covers rows: "Extra payment calculator", "Mortgage payoff".

### LoanToValueRatio      (ID: FIN-LOAN-006)
Formula: LTV = loan / property value (decimal).
Errors: ErrOutOfRange, ErrNegative, ErrZeroDenominator.
Primary source: Fannie Mae *Selling Guide* B2-1.2-01 "Loan-to-Value (LTV) Ratios" (06/01/2022),
https://selling-guide.fanniemae.com/sel/b2-1.2-01/loan-value-ltv-ratios — "Divide the original loan amount by
the property value." CONFIRMED. Second: HUD Handbook 4000.1 Glossary (last revised 8/12/2026): "The
Loan-to-Value (LTV) is computed as the Base Loan Amount divided by the Adjusted Value."
https://www.hud.gov/sites/default/files/Housing/documents/40001-hsgh-Update-18-Glossary.pdf CONFIRMED.
CFPB "What is a loan-to-value ratio…" (en-121) states the concept only.
Worked example: (36000000, 40000000) → 0.9 (HAND-CHECKED).
Notes: Fannie Mae's delivery rule (truncate to 2 decimals, round up to whole percent) is lender policy, not
applied; the result is unrounded.
Covers rows: "Loan to value".

### PITIPayment      (ID: FIN-LOAN-007)
Formula: total = PI + MI + round_he((annualTax + annualInsurance)/12), computed exactly in integers.
Errors: ErrOutOfRange, ErrNegative, ErrOverflow.
Primary source: CFPB, "On a mortgage, what's the difference between my principal and interest payment and my
total monthly payment?" (served at https://www.consumerfinance.gov/ask-cfpb/what-is-a-piti-en-1941/):
"Principal + interest + mortgage insurance (if applicable) + escrow (homeowners insurance and tax) = total
monthly payment". CONFIRMED. HUD 4000.1 Glossary: "PITI - Principal, Interest, Taxes, and Insurance".
Worked example: (150000, 360000, 120000, 5000) → 150000 + 5000 + 40000 = 195000 (HAND-CHECKED).
Covers rows: "PITI", "Mortgage payment with taxes and insurance".

### DebtServiceCoverageRatio      (ID: FIN-LOAN-008)
Formula: DSCR = NOI / debt service (same period).
Errors: ErrOutOfRange, ErrNegative (debt service < 0), ErrZeroDenominator. NOI may be negative.
Primary source: OCC, *Comptroller's Handbook: Commercial Real Estate Lending*, Version 2.0, March 2022, p. 43
"Debt-Service Coverage Ratio": "The DSCR, calculated by dividing the NOI by the annual debt service
requirements, measures the borrower's ability to service its debt."; Glossary p. 137.
https://www.occ.treas.gov/publications-and-resources/publications/comptrollers-handbook/files/commercial-real-estate-lending/pub-ch-commercial-real-estate.pdf
CONFIRMED (PDF text read 2026-11-12). No numeric example in the handbook.
Worked example: (12500000, 10000000) → 1.25 (HAND-CHECKED).
Covers rows: "DSCR".

### DebtToIncomeRatio      (ID: FIN-LOAN-009)
Formula: DTI = monthly debt payments / gross monthly income.
Errors: ErrOutOfRange, ErrNegative, ErrZeroDenominator.
Primary source: CFPB, "What is a debt-to-income ratio?",
https://www.consumerfinance.gov/ask-cfpb/what-is-a-debt-to-income-ratio-en-1791/ — "all your monthly debt
payments divided by your gross monthly income"; example $1,500 + $100 + $400 = $2,000 on $6,000 → 33 %.
CONFIRMED. Second: Fannie Mae Selling Guide B3-6-02 (04/02/2025). CONFIRMED.
Worked example: (200000, 600000) → 0.3333… (CFPB "33 percent").
Covers rows: "Debt to income".

### AffordableLoanAmount      (ID: FIN-LOAN-010)
Formula: allowance = income·ratio − taxesAndInsurance (monthly); loan = round_he(allowance·(P/A,i,n)); 0 when
the allowance ≤ 0. Uses interest.UniformSeriesPresentWorth (FIN-INT-009).
Inputs: grossMonthlyIncomeCents ≥ 0, maxHousingRatio ∈ [0, 1] (caller's front-end limit; no constant encoded),
monthlyTaxesAndInsuranceCents ≥ 0, ratePerPeriod, periods ≥ 1.
Errors: ErrOutOfRange, ErrNegative, ErrNotFinite, ErrRate, ErrPeriods, ErrOverflow.
Sources (a composition of confirmed definitions, no single page states it end to end):
HUD 4000.1 Glossary "PTI - Total Mortgage Payment to Effective Income Ratio" (the housing ratio); CFPB total
monthly payment formula (above); S-PV eq. 8.14 (loan = PV of the P&I payments). All CONFIRMED.
Worked example: (600000, 0.28, 40000, 0.005, 360) → allowance 128000 → 21349327 (HAND-CHECKED).
Covers rows: "How much house can I afford", "Mortgage qualification".

## FOR THE LEAD
- INDEX.md rows FIN-LOAN-001…010; doc/sources entries for OpenStax Principles of Finance, CFPB pages, Fannie Mae
  Selling Guide, HUD 4000.1 Glossary, OCC CRE handbook.
- ASSUMPTIONS rows needed: (1) AmortizationSchedule caps periods at 10⁷ (ErrPeriods) so a bad input cannot
  request a multi-GB allocation; (2) the last payment absorbs the payment's rounding (brief), so the book's
  printed last row differs by cents; (3) superseded by the review ruling: BalloonBalance returns 0 after payoff; (4) AffordableLoanAmount returns 0 (no error) when taxes and insurance exceed the allowance;
  (5) fin/loan maps fin/interest errors onto its own sentinels (errors.Is works with loan.ErrX and interest.ErrX).
- fin/interest FIN-INT-018 example: OpenStax §8.2 $75,000, 3 %, 5 → book ≈ $16,376.60; half-even gives 1637659.

## Research: fin/invest


All pages below were opened on 2026-11-12. "OpenStax PoF" = Julie Dahlquist, Rainford Knight, et al., *Principles of
Finance*, 1st ed., OpenStax (Rice University), 2022, CC BY 4.0, https://openstax.org/books/principles-finance/
(a 2e dated 2026-06-24 also exists at /books/principles-finance-2e/; section numbers checked for 10.4 only).
Crossref DOIs were confirmed with https://api.crossref.org/works?query.bibliographic=… on 2026-11-12.

Covers rows: the lead's vCalc row list was not handed to agent B; the lead maps rows to these IDs.

Sources confirmed once and reused below:
- [MAC] Frederick R. Macaulay, *Some Theoretical Problems Suggested by the Movements of Interest Rates, Bond Yields and
  Stock Prices in the United States since 1856*, NBER, 1938, ch. 2 "The Concept of Long Term Interest Rates", pp. 44–51,
  formula for D on p. 48, closed form p. 49, table "Duration in years" p. 51.
  https://www.nber.org/system/files/chapters/c6342/c6342.pdf (scanned pages read as images). No DOI (NBER book, 1938);
  RePEc record https://ideas.repec.org/b/nbr/nberbk/maca38-1.html.
- [SOA] Society of Actuaries, "Notation and terminology used for Exam FM",
  https://www.soa.org/globalassets/assets/Files/Edu/2019/exam-fm-notation-terminology2.pdf: "Modified duration = −P′(i)/P(i)
  and Modified convexity = P″(i)/P(i)", Macaulay duration Σ t·R_t(1+i)^−t / Σ R_t(1+i)^−t.
  NOT on the F-2 list (see FOR THE LEAD).
- [TD-BILL] US Treasury (TreasuryDirect), "Price, Yield and Rate Calculations for a Treasury Bill" (2004),
  https://www.treasurydirect.gov/instit/annceresult/press/preanre/2004/ofcalc6decbill.pdf, pp. 1–3.
- [BS] Fischer Black, Myron Scholes, "The Pricing of Options and Corporate Liabilities", Journal of Political Economy
  81(3):637–654, 1973, doi:10.1086/260062 (Crossref confirmed). Readable copy:
  https://www.cs.princeton.edu/courses/archive/fall09/cos323/papers/black_scholes73.pdf; eq. (13) p. 644 (call),
  eqs. (25)–(27) p. 647 (put), ideal conditions a)–g) p. 640.

### CurrentYield      (ID: FIN-INVEST-001)
Formula: CY = C / P; C = annual coupon interest in dollars, P = price paid.
Inputs and units: annualCouponCents ≥ 0, priceCents > 0 (cents). Output: decimal (0.05 = 5 %).
Errors: ErrOutOfRange (|x| > 2⁵³), ErrNegative (coupon < 0 or price < 0), ErrZeroDenominator (price = 0).
Primary source: MSRB, "Municipal Bond Basics", https://www.msrb.org/node/13886: "Current Yield is the ratio of the annual
dollar amount of interest paid on a bond to its purchase price." CONFIRMED (URL opened 2026-11-12).
Worked example: HAND-CHECKED: coupon $50.00 (5000), price $952.06 (95206) → 5000/95206 = 0.0525177…
Variants: dividend yield (FIN-INVEST-013) is the same ratio for a stock.
Best approach: one division, O(1), no allocation.
Notes: no rounding of the ratio; never a percentage.

### BondPrice      (ID: FIN-INVEST-002)
Formula: P = c·[1 − (1+i)⁻ⁿ]/i + F·(1+i)⁻ⁿ; F face, c = F·couponRatePerYear/m coupon per period (unrounded),
i = yieldPerYear/m, n periods, m periodsPerYear. i = 0: P = c·n + F.
Inputs and units: faceCents ≥ 0; couponRatePerYear ≥ 0 decimal; yieldPerYear decimal nominal annual; periods ≥ 1;
periodsPerYear ≥ 1. Output: cents, half-even. Valid on a coupon date (no accrued interest).
Errors: ErrOutOfRange, ErrNegative, ErrRate (NaN/Inf or i ≤ −1), ErrPeriods, ErrOverflow.
Primary source: OpenStax PoF §10.2 "Bond Valuation", https://openstax.org/books/principles-finance/pages/10-2-bond-valuation,
"Bond Price = Par Value×1/(1+r)^n + Coupon Payment Amount×[1−1/(1+r)^n]/r". CONFIRMED (opened 2026-11-12).
Second: [MAC] p. 48 footnote 22, P = I/R + … + I/Rⁿ + F/Rⁿ. CONFIRMED.
Worked example: OpenStax §10.2 3M bond: F $1,000, 2.25 %, annual, 15 years, YTM 1.24 % → $1,137.47 (computed
113747.096 cents → 113747). §10.2 calculator example: 4 % semiannual, 30 periods, 5 % → $895.35 (89534.85 → 89535).
Variants: annuity factor from fin/interest is in cents per payment; the coupon here is fractional cents, so the
formula is evaluated directly.
Best approach: O(1) with (1+i)⁻ⁿ = exp(−n·log1p(i)) and 1 − (1+i)⁻ⁿ = −expm1(−n·log1p(i)) (Goldberg 1991 §1.4).
Notes: no accrued interest / dirty price; no day counts.

### YieldToMaturity      (ID: FIN-INVEST-003)
Formula: the periodic i solving P = Σₖ cₖ(1+i)⁻ᵏ, cₖ = coupon (+ face at k = n); result = m·i (nominal annual).
Inputs and units: priceCents > 0, faceCents ≥ 0, couponRatePerYear ≥ 0, periods ≥ 1, periodsPerYear ≥ 1.
The coupon per period is rounded half-even to the cent (a paid coupon is whole cents) because the root finder
(interest.InternalRateOfReturn) takes cents.
Output: decimal nominal annual rate.
Errors: ErrOutOfRange, ErrNegative, ErrDomain (price = 0), ErrRate, ErrPeriods, ErrNoSolution, ErrNoConvergence.
Primary source: OpenStax PoF §10.2 (YTM is the discount rate that makes PV of the bond's flows = price). CONFIRMED.
Worked example: OpenStax §10.2 Coca-Cola: price $952.06, par $1,000, $5.00 per half-year, 14 periods → semiannual
0.8651 %, annual 1.73 % (hand: 0.0086506 → 0.0173012).
Best approach: interest.AnnuityRate (FIN-INT-042) on (price, coupon cents, face, n): closed-form factors, O(k)
evaluations independent of n, O(1) space, no allocation (revised after review; the n+1-flow IRR version panicked at
periods = MaxInt and allocated).

### MacaulayDuration      (ID: FIN-INVEST-004)
Formula: D = Σₖ k·cₖ·vᵏ / Σₖ cₖ·vᵏ (in periods), v = 1/(1+i), i = y/m; result in years = D/m.
Inputs: as BondPrice. Output: years.
Errors: as BondPrice plus ErrDomain (price 0: face and coupon both 0).
Primary source: [MAC] p. 48 (definition, "D = the 'duration' of the bond—in half years"), p. 49 (closed form). CONFIRMED.
Second: [SOA] Macaulay duration. CONFIRMED (list caveat).
Worked example: [MAC] p. 51 table "Duration in years" (semiannual, bonds selling at par): 4 % 1 yr .990, 4 % 10 yr
8.339, 5 % 10 yr 7.989, 6 % 10 yr 7.662, 6 % 3 yr 2.790, 4 % 100 yr 25.014 (13 of 15 entries reproduced to 3 dp; 5 % 1 yr .987 and 6 % 25 yr 13.254 do not, see FOR THE LEAD).
Best approach: closed form (p. 49) is O(1) but divides by R − 1 and cancels as i → 0; chosen: one Horner pass
for p(v), p′(v), p″(v): O(n) time, O(1) space, no allocation (n ≤ 1200 for a 100-year monthly bond).
Notes: no effective duration, no option-adjusted.

### ModifiedDuration      (ID: FIN-INVEST-005)
Formula: −(1/P)·dP/dy = D_years/(1 + y/m), with P written as a function of the nominal annual yield y.
Primary source: [SOA] "Modified duration = −P′(i)/P(i)". CONFIRMED (SOA not on F-2 list).
Worked example: HAND-CHECKED from [MAC] 6 % 10 yr at par: 7.661900/1.03 = 7.438737; oracle: central difference of
BondPrice.
Best approach: same Horner pass as 004.

### Convexity      (ID: FIN-INVEST-006)
Formula: (1/P)·d²P/dy² = Σ k(k+1)·cₖ·vᵏ⁺² / (P·m²), years².
Primary source: [SOA] "Modified convexity = P″(i)/P(i)". CONFIRMED (SOA not on F-2 list).
Worked example: HAND-CHECKED 6 % 10 yr semiannual at par: 68.774822 years²; oracle: second difference of price.

### BankDiscountYield      (ID: FIN-INVEST-007)
Formula: d = ((F − P)/F)·(360/r); r days to maturity.
Primary source: [TD-BILL] p. 1 "d = ((100 − P)/100) * (360/r)". CONFIRMED. Also TreasuryDirect "Understanding pricing",
https://treasurydirect.gov/marketable-securities/understanding-pricing/ (Price = Face(1 − d·t/360), $1,000, 0.145 %, 182 days → $999.27).
Worked example: [TD-BILL] p. 1: P = 99.937778 per 100, r = 28 → d = 0.007999971428571 (before display rounding 0.800 %).
In cents: face 100000000, price 99937778.
Errors: ErrOutOfRange, ErrNegative, ErrZeroDenominator (face 0), ErrDomain (price 0), ErrPeriods (days < 1).

### MoneyMarketYield      (ID: FIN-INVEST-008)
SOURCE NOT CONFIRMED: (F − P)/P × 360/t appears only in CFA prep sites and blogs (analystprep, xplaind, nasdaq,
CFI). Tried: TreasuryDirect bill-calculation PDF and "Understanding pricing" (no money market yield), Fed Trading
and Capital-Markets Activities Manual (no formula found in text search), searches of federalreserve.gov, newyorkfed.org,
sec.gov, treasury.gov. CFA reading is behind the learning portal. NOT BUILT.

### BondEquivalentYield      (ID: FIN-INVEST-009)
Formula: i = ((F − P)/P)·(y/r); y = days in the year (365 or 366), r days to maturity, only r ≤ half a year.
Primary source: [TD-BILL] p. 2 "Calculate Coupon Equivalent Yield For bills of not more than one half-year to maturity:
i = ((100 − P)/P) * (y/r)". CONFIRMED. The > half-year quadratic (p. 2–3) is not built (signature has no place for it).
Worked example: [TD-BILL] p. 2: P = 99.937778, y = 366, r = 28 → i = 0.008138368141143 (0.814 %).
Errors: ErrOutOfRange, ErrNegative, ErrZeroDenominator (price 0), ErrPeriods (days < 1),
ErrDomain (daysInYear ∉ {365, 366} or 2·days > daysInYear).

### HoldingPeriodYield      (ID: FIN-INVEST-010)
Formula: (P₁ − P₀ + D)/P₀ = capital gain yield + dividend yield.
Primary source: OpenStax PoF §15.1 "Risk and Return to an Individual Asset",
https://openstax.org/books/principles-finance/pages/15-1-risk-and-return-to-an-individual-asset
("Total Percent Return = Dividend Yield + Capital Gain Yield"). CONFIRMED.
Worked example: §15.1 Target: bought $128.74, sold $176.53, dividends $2.68 → $50.47, 39.20 % (5047/12874 = 0.392030…).
3M: $175 → $169.72, dividends $5.91 → 0.36 % (63/17500 = 0.0036 per share, $6.30 for 10 shares).
Errors: ErrOutOfRange, ErrNegative (beginning or ending value < 0), ErrZeroDenominator (beginning 0).
Best approach: numerator in exact int64 (|sum| ≤ 3·2⁵³), one float division.

### AnnualizedHoldingPeriodYield      (ID: FIN-INVEST-011)
Formula: EAR = (1 + HPR)^m − 1, m = holding periods per year = 365/days.
Primary source: OpenStax PoF §15.1: "EAR = [1 + Holding Period Percentage Return]^m − 1 … m is the number of holding
periods in a year". CONFIRMED. The 365-day year that turns days into m is NOT stated by OpenStax (FOR THE LEAD).
Worked example: §15.1 three-month example m = 4: (1.1462)⁴ − 1 = 0.7260 (not reachable with integer days, since
365/4 = 91.25); test uses days = 73 (m = 5) hand-checked, and days = 365 (m = 1) → HPR.
Best approach: expm1(m·log1p(HPR)). Errors: ErrRate (NaN/Inf, ≤ −1), ErrPeriods (days < 1), ErrOverflow.

### CapitalGainsYield      (ID: FIN-INVEST-012)
Formula: (P₁ − P₀)/P₀. Primary source: OpenStax PoF §15.1 ("The capital gain yield is the change in the stock price
divided by the initial stock price"). CONFIRMED.
Variants: return on investment, (final − cost)/cost: absorbs FIN-BIZ-033 (removed as its duplicate); tags roi,
return-on-investment.
Worked example: Target 4779/12874 = 0.371213…; 3M (169.72 − 175)/175 = −0.030171… (source −3.02 %).

### DividendYield      (ID: FIN-INVEST-013)
Formula: D/P₀. Primary source: OpenStax PoF §15.1 ("dividing the dividends you received by the initial stock price").
CONFIRMED. Worked example: 3M 5.91/175 = 0.033771 (3.38 %); Target 2.68/128.74 = 0.020817.

### GordonGrowthPrice      (ID: FIN-INVEST-014)
Formula: P₀ = D₁/(r − g), D₁ next dividend, r required return, g constant growth, r > g.
Primary source: OpenStax PoF §11.2 "Dividend Discount Models (DDMs)",
https://openstax.org/books/principles-finance/pages/11-2-dividend-discount-models-ddms ("Stock Value = D1/(r − g)",
named the Gordon growth model). CONFIRMED. Original: M. J. Gordon, E. Shapiro, "Capital Equipment Analysis: The
Required Rate of Profit", Management Science 3(1):102–110, 1956, doi:10.1287/mnsc.3.1.102 (Crossref CONFIRMED; the
paper's own text was not opened).
Worked example: §11.2: D₀ $5.00, g 4 %, r 8 % → $130.00 (D₁ = 520 cents → 13000).
Errors: ErrOutOfRange, ErrNegative (dividend), ErrRate, ErrDomain (r ≤ g), ErrOverflow.

### ZeroGrowthStockPrice      (ID: FIN-INVEST-015)
Formula: P = D/r (also preferred stock). Primary source: OpenStax PoF §11.2 ("PV = D/r"). CONFIRMED.
Worked example: §11.2 $2.00, 6 % → $33.33 (3333.33 → 3333). Errors: ErrDomain (r ≤ 0), ErrRate, ErrNegative.

### CAPMExpectedReturn      (ID: FIN-INVEST-016)
Formula: Re = Rf + β(Rm − Rf). Primary source: OpenStax PoF §15.3 "The Capital Asset Pricing Model (CAPM)",
https://openstax.org/books/principles-finance/pages/15-3-the-capital-asset-pricing-model-capm. CONFIRMED.
Original: W. F. Sharpe, "Capital Asset Prices: A Theory of Market Equilibrium under Conditions of Risk", Journal of
Finance 19(3):425–442, 1964, doi:10.1111/j.1540-6261.1964.tb02865.x (Crossref CONFIRMED; text not opened).
Worked example: §15.3 DAL: Rf 3.36 %, Rm 11.64 %, β 1.39 → 14.87 % (0.148692).
Errors: ErrRate (rates NaN/Inf/≤ −1), ErrNotFinite (beta), ErrOverflow.

### RiskPremium      (ID: FIN-INVEST-017)
Formula: E[R] − Rf. Primary source: OpenStax PoF §15.3 ("Risk Premium = S&P 500 Avg Return − US T-Bill Avg Return =
11.64 % − 3.36 % = 8.28 %"). CONFIRMED. Worked example: that line (0.0828).

### BlackScholesCall      (ID: FIN-INVEST-018)
Formula: [BS] eq. (13): w = x·N(d₁) − c·e^{r(t−t*)}·N(d₂); d₁ = [ln(x/c) + (r + ½v²)(t* − t)]/(v√(t* − t)),
d₂ = d₁ − v√(t* − t). x spot, c strike, r riskless rate (continuous, per year), v volatility per year, t* − t years.
N(x) = ½·erfc(−x/√2) (brief). Assumptions [BS] p. 640: European, no dividends, constant r and v.
Primary source: [BS] CONFIRMED (Crossref + scanned page 644 read).
Worked example: the paper has no numeric example. HAND-CHECKED: S $42, K $40, r 10 %, σ 20 %, T 0.5 →
475.942 cents → 476; put 80.860 → 81 (Python math.erfc).
Limits: v√T = 0 → boundary (8)/discounted intrinsic max(S − K·e^{−rT}, 0).
Errors: ErrOutOfRange, ErrNegative, ErrDomain (spot or strike 0), ErrRate, ErrNotFinite (vol, years), ErrOverflow.

### BlackScholesPut      (ID: FIN-INVEST-019)
Formula: [BS] eq. (27) p. 647: u = −x·N(−d₁) + c·e^{−r(t*−t)}·N(−d₂); parity eq. (25) w − u = x − c·e^{r(t−t*)}.
Worked example: as 018 (put 81 cents); test checks parity (25) to the cent.

### ModifiedInternalRateOfReturn      (ID: FIN-INVEST-020)
Formula: MIRR = (TV/PV_out)^{1/n} − 1; PV_out = Σ_{cₜ<0} |cₜ|(1+f)^{−t}; TV = Σ_{cₜ>0} cₜ(1+r)^{n−t}, n = last period.
Primary source: OpenStax PoF §16.4 "Alternative Methods",
https://openstax.org/books/principles-finance/pages/16-4-alternative-methods (three steps: PV of outflows at the cost
of capital, terminal value of inflows at the cost of capital, solve). CONFIRMED. The source uses one rate for both;
the API's two rates (finance, reinvest) reduce to it when equal.
Worked example: §16.4 embroidery machine, 9 %, −16,000; 2,000; 4,000; 5,000 ×4 → TV $31,595.22 (the listed terms sum to
$31,589.23: source arithmetic slip, recorded), MIRR 12 % (hand 0.1200476).
Errors: ErrRate, ErrDomain (< 2 flows), ErrOutOfRange, ErrNoSolution (no inflow or no outflow), ErrOverflow.

### ProfitabilityIndex      (ID: FIN-INVEST-021)
Formula: PI = PV(inflows)/PV(outflows). Primary source: OpenStax PoF §16.4 ("PI = PV(Cash Inflows)/PV(Cash Outflows)").
CONFIRMED. Worked example: §16.4 $18,836/$16,000 = 1.18 (hand 1.177226).
Errors: ErrRate, ErrDomain (empty), ErrOutOfRange, ErrZeroDenominator (no outflow).

### PaybackPeriod      (ID: FIN-INVEST-022)
Formula: smallest t with cumulative Σ_{0..t} cₖ ≥ 0; fractional: (t − 1) + (−cum_{t−1})/cₜ.
Primary source: OpenStax PoF §16.1 "Payback Period Method" (definition, Sam's Sporting Goods example 4 years) and §16.4
(fraction rule "5 + 145.72/2,981.34 = 5.05 years"). CONFIRMED.
Worked example: §16.1: −16,000; 2,000; 4,000; 5,000; 5,000; 5,000 → 4 years.
Errors: ErrDomain (empty), ErrOutOfRange, ErrNoSolution (never recovered). Cumulative sum exact in int64.

### DiscountedPaybackPeriod      (ID: FIN-INVEST-023)
Formula: as 022 on cₜ(1+i)^{−t}. Primary source: OpenStax PoF §16.4. CONFIRMED.
Worked example: §16.4 9 % table → 5 + 145.72/2,981.34 = 5.05 years (hand 5.048876).

### EquivalentAnnualAnnuity      (ID: FIN-INVEST-024)
REMOVED (never released; lead decision 2026-11-12): duplicate of FIN-INT-018 AnnuityPaymentFromPresentValue
(equivalent annual annuity = annuity payment of the NPV). Files moved to _delete/fin/invest/. Research kept below.
Formula: EAA = NPV·(A/P,i,n) = NPV·i/(1 − (1+i)^{−n}). Primary source: OpenStax PoF §16.5 "Choosing between Projects",
https://openstax.org/books/principles-finance/pages/16-5-choosing-between-projects ("Find the annuity that has the same
present value as the NPV and the same number of periods as the project"). CONFIRMED.
Worked example: §16.5 older truck NPV $49,474, 10 %, 3 years → $19,894 (hand 1989422.78 → 1989423 cents); newer truck
$80,658, 6 years → $18,520 (1851967.21 → 1851967).
Best approach: interest.CapitalRecovery (FIN-INT-010), one multiply, half-even.

### SharpeRatio      (ID: FIN-INVEST-025)
Formula: (Rp − Rf)/σp. Primary source: OpenStax PoF §15.4 "Applications in Performance Measurement", eq. 15.14,
https://openstax.org/books/principles-finance/pages/15-4-applications-in-performance-measurement. CONFIRMED.
Originals: W. F. Sharpe, "Mutual Fund Performance", Journal of Business 39(S1):119–138, 1966, doi:10.1086/294846;
"The Sharpe Ratio", Journal of Portfolio Management 21(1):49–58, 1994, doi:10.3905/jpm.1994.409501 (Crossref CONFIRMED).
Worked example: §15.4 McKinley: 20 %, 3 %, σ 26 % → 0.65 (0.653846).
Errors: ErrNotFinite, ErrRate, ErrZeroDenominator (σ = 0), ErrNegative (σ < 0).

### WeightedAverageCostOfCapital      (ID: FIN-INVEST-026)
Formula: WACC = E/(D+E)·re + D/(D+E)·rd·(1 − T) (no preferred; signature has none).
Primary source: OpenStax PoF §17.3 "Calculating the Weighted Average Cost of Capital",
https://openstax.org/books/principles-finance/pages/17-3-calculating-the-weighted-average-cost-of-capital. CONFIRMED.
Worked example: Bluebonnet: D = 0.97 × $5,000,000 = $4,850,000 (§17.1), E = 1,000,000 × $15 (§17.1), rd 6.312 %,
T 21 % (§17.2), re 13.4 % (CAPM, §17.3) → 11.33 % in source (weights rounded to 24 %/75.6 %); exact 0.1134430.
Errors: ErrOutOfRange, ErrNegative (values), ErrZeroDenominator (D + E = 0), ErrRate, ErrDomain (T ∉ [0, 1]).

### TaxEquivalentYield      (ID: FIN-INVEST-027)
Formula: tax-exempt yield / (1 − tax bracket). Primary source: MSRB, "Understanding Taxable Municipal Bonds",
https://msrb.org/node/13519. CONFIRMED. Worked example: 3.00 %/(1 − 0.37) = 4.76 % (0.0476190).
Errors: ErrRate, ErrDomain (T ∉ [0, 1)).

### NetAssetValuePerShare      (ID: FIN-INVEST-028)
Formula: (total assets − total liabilities)/shares outstanding.
Primary source: US SEC, Fast Answers "Net Asset Value", http://www.sec.gov/answers/nav.htm (now redirects to
investor.gov glossary, which returned 403 to the fetcher; read from the archived copy
https://webarchive.library.unt.edu/web/20130214105729mp_/http://www.sec.gov/answers/nav.htm). CONFIRMED.
Worked example: SEC: assets $100 m, liabilities $10 m → NAV $90 m; NAV $100 m / 10,000,000 shares → $10.
Best approach: exact int64 quotient and remainder with half-even rounding; no float.
Errors: ErrOutOfRange, ErrNegative (assets, liabilities or shares < 0), ErrZeroDenominator (shares 0).

### LeveredBetaHamada      (ID: FIN-INVEST-029)
SOURCE NOT CONFIRMED: Hamada, "The Effect of the Firm's Capital Structure on the Systematic Risk of Common Stocks",
Journal of Finance 27(2):435–452, 1972, doi:10.1111/j.1540-6261.1972.tb00971.x is confirmed on Crossref, but no
readable copy of the paper was found, and βL = βU[1 + (1 − T)D/E] appears only on Wikipedia, consultancies (Oxera)
and calculator sites. OpenStax PoF has no levered/unlevered beta. NOT BUILT.

### BidAskSpread      (ID: FIN-INVEST-030)
Formula: relative quoted spread = (ask − bid)/((ask + bid)/2).
Primary source: Charles Collver, "A characterization of market quality for small capitalization US equities", US SEC,
Division of Trading and Markets / DERA staff paper, September 2014,
https://www.sec.gov/marketstructure/research/small_cap_liquidity.pdf, p. 2: "Quoted spread is the difference between
the best ask price and best bid price. Relative quoted spread is the quoted spread scaled by the bid-ask midpoint."
Locked or crossed quotes (spread ≤ 0) are excluded there. CONFIRMED (staff paper; disclaims Commission views).
Worked example: HAND-CHECKED bid $24.19, ask $24.21 → 2/2420 = 0.000826446.
Errors: ErrOutOfRange, ErrNegative (bid < 0), ErrDomain (ask ≤ bid: locked or crossed).

### ImpliedRepoRate      (ID: FIN-INVEST-031)
SOURCE NOT FOUND within the time box: no OpenStax, Treasury, Fed or SEC page with the implied repo formula was
found ("build only if sourced"). NOT BUILT.

## FOR THE LEAD
Rulings received 2026-11-12: SOA Exam FM material accepted (F-3) for 005/006; 011's 365-day year logged as an
assumption; 024 removed (→ FIN-INT-018); duration/convexity periods capped at 10 000 000.
1. FIN-INVEST-005 ModifiedDuration and 006 Convexity cite the Society of Actuaries Exam FM notation sheet, which is
   not on the F-2 list (Macaulay 1938 defines only duration). Built and flagged: accept SOA as an official exam body,
   or hold both back (move the two files and tests to _delete/).
2. FIN-INVEST-011: the 365-day year in m = 365/days is not in OpenStax §15.1 (which defines m as holding periods per
   year). Needs an assumption row (A-xxx) or a source; Treasury's bill formula [TD-BILL] uses y = 365/366.
3. Duration/convexity are O(n) by choice (closed form cancels as i → 0); please benchmark at n = 1e3/1e4/1e5.
4. YieldToMaturity rounds the per-period coupon to the cent before calling interest.InternalRateOfReturn (cents API).
5. ProfitabilityIndex, MIRR and discounted payback discount directly (Horner) and do not call fin/interest; only 003
   (InternalRateOfReturn) and 024 (CapitalRecovery) import it.
7. Macaulay 1938 p. 51 table: 5 % 1-year (.987) and 6 % 25-year (13.254) disagree with the formula on p. 48
   (0.98780, 13.25083); the other 13 entries agree to 3 dp. Left out of the test, flagged here.
6. OpenStax §16.4 MIRR terminal value $31,595.22 does not equal its own column ($31,589.23).

## Research: fin/business


Every page below was opened on 2026-11-12. Short names used in the entries:

- **PoF** = Julie Dahlquist and Rainford Knight, *Principles of Finance*, 1st ed., OpenStax (Rice University), 2022.
  Section pages under https://openstax.org/books/principles-finance/pages/
- **PA1** = Mitchell Franklin, Patty Graybeal and Dixon Cooper, *Principles of Accounting, Volume 1: Financial
  Accounting*, OpenStax (Rice University), 2019. Pages under https://openstax.org/books/principles-financial-accounting/pages/
- **PA2** = Franklin, Graybeal and Cooper, *Principles of Accounting, Volume 2: Managerial Accounting*, OpenStax
  (Rice University), 2019. Pages under https://openstax.org/books/principles-managerial-accounting/pages/
- **SEC** = U.S. Securities and Exchange Commission, "Beginners' Guide to Financial Statements",
  https://www.sec.gov/about/reports-publications/beginners-guide-financial-statements
- **CFA-FCF** = CFA Institute, "Free Cash Flow Valuation", 2026 Curriculum CFA Program Level II Equity Valuation,
  refresher-reading page, Summary, https://www.cfainstitute.org/insights/professional-learning/refresher-readings/2025/free-cash-flow-valuation
- **UBPR-3** = FFIEC, *UBPR User's Guide*, "Noninterest Income and Expenses — Page 3", 2026-08-20, item 32
  "Efficiency Ratio" (UBPRE088), https://www.cdr.ffiec.gov/CDRDownload/CDR/UserGuide/v186/FFIEC%20UBPR%20User%20Guide%20Noninterest%20Income%20and%20Expenses--Page%203_2026-08-20.PDF
- **UBPR-1** = FFIEC, *UBPR User's Guide*, "Summary Ratios — Page 1", 2026-08-20, item 32 "Net Loans & Leases to
  Deposits" (UBPRE600), https://www.cdr.ffiec.gov/CDRDownload/CDR/UserGuide/v186/FFIEC%20UBPR%20User%20Guide%20Summary%20Ratios--Page%201_2026-08-20.PDF

Conventions for every entry (numeric policy of the brief): money is int64 cents with |x| ≤ 2⁵³ (else
ErrOutOfRange); ratios return float64 decimals; a zero denominator is ErrZeroDenominator; money results are
rounded half-even; |money result| > 2⁵³ is ErrOverflow. "Clear Lake" = Clear Lake Sporting Goods (PoF ch. 6);
"Banyan" = Banyan Goods (PA1 Appendix A); both use the same numbers.

---

### CurrentRatio      (ID: FIN-BIZ-001)
Formula: current assets ÷ current liabilities.
Inputs and units: currentAssetsCents, currentLiabilitiesCents (cents, ≥ 0).
Output: ratio, decimal (2 means 2:1).
Errors: ErrOutOfRange, ErrNegative, ErrZeroDenominator.
Primary source: PoF §6.3 "Liquidity Ratios", https://openstax.org/books/principles-finance/pages/6-3-liquidity-ratios; CONFIRMED (URL opened 2026-11-12).
Second source: PA1 Appendix A "Financial Statement Analysis", https://openstax.org/books/principles-financial-accounting/pages/a-financial-statement-analysis; CONFIRMED.
Worked example: Clear Lake current year $200,000 ÷ $100,000 = 2 (PoF §6.3; PA1 App. A, Banyan).
Variants: cash ratio (not in list C).
Best approach: one division; O(1), no allocation.
Notes: must not round the ratio.
Covers rows: current ratio.

### QuickRatio      (ID: FIN-BIZ-002)
Formula: (cash + short-term investments + accounts receivable) ÷ current liabilities.
Inputs and units: cashCents, marketableSecuritiesCents (= short-term investments), receivablesCents, currentLiabilitiesCents (cents, ≥ 0).
Output: ratio, decimal.
Errors: ErrOutOfRange, ErrNegative, ErrZeroDenominator.
Primary source: PoF §6.3; CONFIRMED (opened 2026-11-12). Second: PA1 App. A; CONFIRMED.
Worked example: ($110,000 + $20,000 + $30,000) ÷ $100,000 = 1.6 (PoF §6.3, Clear Lake).
Variants: acid-test (same ratio).
Best approach: exact int64 sum (three values ≤ 2⁵³ cannot overflow int64), one division. O(1), no allocation.
Notes: inventory is NOT in the numerator.
Covers rows: quick ratio, acid test.

### DebtRatio      (ID: FIN-BIZ-003)
Formula: (current liabilities + long-term liabilities) ÷ total assets = total liabilities ÷ total assets.
Inputs and units: totalLiabilitiesCents, totalAssetsCents (≥ 0).
Output: decimal.
Errors: ErrOutOfRange, ErrNegative, ErrZeroDenominator.
Primary source: PoF §6.4 "Solvency Ratios" (debt-to-assets), https://openstax.org/books/principles-finance/pages/6-4-solvency-ratios; CONFIRMED.
Worked example: ($100,000 + $50,000) ÷ $250,000 = 0.6 (PoF §6.4, Clear Lake).
Best approach: O(1), no allocation.
Covers rows: debt ratio, debt-to-assets.

### DebtToEquityRatio      (ID: FIN-BIZ-004)
Formula: total liabilities ÷ total stockholders' equity.
Inputs and units: totalLiabilitiesCents (≥ 0), totalEquityCents (any sign; negative equity gives a negative ratio, documented).
Output: decimal.
Errors: ErrOutOfRange, ErrNegative (liabilities), ErrZeroDenominator.
Primary source: PoF §6.4; CONFIRMED. Second: PA1 App. A; SEC Beginners' Guide ("Total Liabilities / Shareholders' Equity"); CONFIRMED.
Worked example: $150,000 ÷ $100,000 = 1.5 (PoF §6.4, Clear Lake).
Covers rows: debt to equity.

### EquityMultiplier      (ID: FIN-BIZ-005)
Formula: average total assets ÷ average stockholders' equity.
Inputs and units: totalAssetsCents (≥ 0), totalEquityCents (any sign, ≠ 0). The caller passes averages when following the source.
Output: decimal.
Errors: ErrOutOfRange, ErrNegative (assets), ErrZeroDenominator.
Primary source: PoF §6.6 "Profitability Ratios and the DuPont Method", https://openstax.org/books/principles-finance/pages/6-6-profitability-ratios-and-the-dupont-method; CONFIRMED.
Worked example: DuPont factor 2.37 for Clear Lake = $225,000 ÷ $95,000 = 2.3684… (PoF §6.6; inputs are the average assets and average equity used in the same section).
Covers rows: equity multiplier.

### EquityToDebtRatio      (ID: FIN-BIZ-006)
SOURCE NOT CONFIRMED: searched PoF ch. 6 and its index, PA1 App. A and its index, SEC Beginners' Guide, CFA
"Financial Analysis Techniques" page (summary only, no formulas). No acceptable page names an equity-to-debt
(owners' to creditors' equity) ratio. Not built. (It is 1 ÷ FIN-BIZ-004, but the name and definition need a source.)
Retry (agent C2 retry, 2026-11-12): CFA Level II Financial Ratio List (CFA Institute, "Financial Ratio List", CFA Program Level II, https://www.cfainstitute.org/sites/default/files/-/media/documents/support/programs/cfa/cfa_program_level_ii_financial_ratio_list.pdf; ratios 1–41 read) has debt-to-equity (31) and
debt-to-capital (32) but no equity-to-debt ratio; PoF index (https://openstax.org/books/principles-finance/pages/index)
has no "equity ratio" entry. Still not built.

### TimesInterestEarned      (ID: FIN-BIZ-007)
Formula: EBIT ÷ interest expense.
Inputs and units: ebitCents (any sign), interestExpenseCents (≥ 0).
Output: decimal ("times").
Errors: ErrOutOfRange, ErrNegative, ErrZeroDenominator.
Primary source: PoF §6.4; CONFIRMED. Second: PA1 App. A; CONFIRMED.
Worked example: $43,000 ÷ $2,000 = 21.5 (PoF §6.4, Clear Lake).
Covers rows: times interest earned, interest coverage.

### FixedChargeCoverageRatio      (ID: FIN-BIZ-008)
Formula: (EBIT + lease payments) ÷ (interest payments + lease payments). fixedChargesCents = the fixed charge other than interest (lease payments; in S-K terms the interest portion of rental expense).
Inputs and units: ebitCents (any sign), fixedChargesCents (≥ 0), interestExpenseCents (≥ 0), cents.
Output: decimal ("times").
Errors: ErrOutOfRange, ErrNegative (fixed charges, interest), ErrZeroDenominator (interest + lease = 0).
Primary source: CFA Institute, "Financial Ratio List", CFA Program Level II, https://www.cfainstitute.org/sites/default/files/-/media/documents/support/programs/cfa/cfa_program_level_ii_financial_ratio_list.pdf, ratio 34: "Fixed charge coverage ratio = (Earnings before interest and taxes + Lease payments) ÷ (Interest payments + Lease payments)"; CONFIRMED (PDF opened 2026-11-12, page ii).
Second source: SEC Regulation S-K, 17 CFR 229.503(d) as in force 2018-01-01 (later removed from S-K), https://www.ecfr.gov/api/renderer/v1/content/enhanced/2018-01-01/title-17?part=229&section=229.503: "If you register debt securities, show a ratio of earnings to fixed charges"; Instruction 1: fixed charges = interest expensed and capitalized, amortized premiums/discounts, estimated interest within rental expense, preference dividends of consolidated subsidiaries; earnings = pre-tax income from continuing operations plus fixed charges (and listed adjustments); CONFIRMED (opened 2026-11-12). Same structure: (pre-tax + interest + rent interest) ÷ (interest + rent interest) = (EBIT + F) ÷ (I + F); the S-K adjustments (capitalized interest, equity investees, preference dividends) are not modelled.
Worked example: HAND-CHECKED (no numbers in either source): EBIT $50,000, lease $10,000, interest $20,000 → 60,000 ÷ 30,000 = 2.
Variants: S-K "ratio of combined fixed charges and preference dividends to earnings" (not built).
Best approach: exact int64 sums (≤ 2⁵⁴), one division; ≤ 2 ulp when a sum exceeds 2⁵³. O(1), no allocation.
Notes: FOR THE LEAD: the signature's fixedChargesCents is read as lease payments (fixed charges other than interest), per the CFA list.
Covers rows: fixed charge coverage.

### CashFlowAdequacyRatio      (ID: FIN-BIZ-009)
SOURCE NOT CONFIRMED: PA1 §16.5 (cash-flow ratios) defines free cash flow, cash flows to sales and to assets only;
CFA "Analyzing Statements of Cash Flows II" page is summary only. Not built.
Retry (agent C2 retry, 2026-11-12): CFA Level II Financial Ratio List (ratios 1–41) has no cash-flow adequacy ratio; the CFA page
https://www.cfainstitute.org/insights/professional-learning/refresher-readings/2026/analyzing-statements-of-cash-flows-ii
shows learning outcomes only; web search restricted to .gov/openstax.org/cfainstitute.org found only blogs
(fathomhq, wallstreetprep, accountingtools: not acceptable; they also disagree on numerator/denominator order). Still not built.

### CashFlowToCurrentLiabilities      (ID: FIN-BIZ-010)
SOURCE NOT CONFIRMED: same search as 009 (PA1 §16.5, CFA pages summary only). Not built.
Retry (agent C2 retry, 2026-11-12): same pages as 009 retry; no acceptable page defines CFO ÷ average current liabilities. Still not built.

### CashFlowToCapitalExpenditures      (ID: FIN-BIZ-011)
SOURCE NOT CONFIRMED: same search as 009. Not built.
Retry (agent C2 retry, 2026-11-12): same pages as 009 retry (the CFA reading's "reinvestment ratio" is named in outside summaries
but its formula is not visible on any free CFA page). Still not built.

### AssetTurnover      (ID: FIN-BIZ-012)
Formula: net sales ÷ average total assets.
Inputs and units: netSalesCents (≥ 0), averageTotalAssetsCents (≥ 0).
Output: decimal ("times").
Errors: ErrOutOfRange, ErrNegative, ErrZeroDenominator.
Primary source: PoF §6.2 "Operating Efficiency Ratios", https://openstax.org/books/principles-finance/pages/6-2-operating-efficiency-ratios; CONFIRMED. Second: PA1 App. A; CONFIRMED.
Worked example: $120,000 ÷ $225,000 = 0.53 (PoF §6.2, Clear Lake).
Covers rows: total asset turnover.

### AssetToSalesRatio      (ID: FIN-BIZ-013)
SOURCE NOT CONFIRMED: no acceptable page defines an asset-to-sales ratio (PoF and PA1 indexes searched). Not built.
Retry (agent C2 retry, 2026-11-12): CFA Level II Financial Ratio List has total asset turnover (14) but no asset-to-sales ratio; PoF index
has no "capital intensity" or "asset to sales" entry. Still not built.

### ReceivablesTurnover      (ID: FIN-BIZ-014)
Formula: net credit sales ÷ average accounts receivable; average = (beginning + ending) ÷ 2.
Inputs and units: netCreditSalesCents (≥ 0), averageReceivablesCents (≥ 0).
Errors: ErrOutOfRange, ErrNegative, ErrZeroDenominator.
Primary source: PA1 §9.3 "Determine the Efficiency of Receivables Management Using Financial Ratios", https://openstax.org/books/principles-financial-accounting/pages/9-3-determine-the-efficiency-of-receivables-management-using-financial-ratios; CONFIRMED. Second: PoF §6.2; CONFIRMED.
Worked example: BWW 2017: $400,000 ÷ $80,000 = 5 (PA1 §9.3); Clear Lake $100,000 ÷ $25,000 = 4 (PoF §6.2).
Covers rows: receivables turnover.

### AverageCollectionPeriod      (ID: FIN-BIZ-015)
Formula: days in period ÷ receivables turnover = days × average receivables ÷ net credit sales.
Inputs and units: averageReceivablesCents, netCreditSalesCents (cents, ≥ 0), daysInPeriod (int > 0; the source uses 365).
Output: days.
Errors: ErrOutOfRange, ErrNegative, ErrZeroDenominator, ErrPeriods.
Primary source: PA1 §9.3 ("number of days' sales in receivables" = 365 ÷ ART); CONFIRMED. Second: PoF §19.1 (receivables collection period = 365 ÷ AR turnover), https://openstax.org/books/principles-finance/pages/19-1-what-is-working-capital; CONFIRMED.
Worked example: BWW 2017: 365 ÷ 5 = 73 days (PA1 §9.3) → (8,000,000¢, 40,000,000¢, 365) = 73.
Best approach: days × receivables ÷ sales in one expression (one rounding fewer than 365 ÷ (sales ÷ receivables)).
Covers rows: average collection period, days sales outstanding.

### InventoryTurnover      (ID: FIN-BIZ-016)
Formula: cost of goods sold ÷ average inventory.
Errors: ErrOutOfRange, ErrNegative, ErrZeroDenominator.
Primary source: PoF §6.2; CONFIRMED. Second: PA1 App. A; SEC guide ("Cost of Sales / Average Inventory for the Period"); CONFIRMED.
Worked example: $60,000 ÷ $37,500 = 1.6 (PoF §6.2, Clear Lake).
Covers rows: inventory turnover.

### DaysInInventory      (ID: FIN-BIZ-017)
Formula: days × inventory ÷ cost of goods sold (= days ÷ inventory turnover).
Inputs and units: averageInventoryCents, costOfGoodsSoldCents (≥ 0), daysInPeriod (> 0; sources use 365).
Output: days.
Errors: ErrOutOfRange, ErrNegative, ErrZeroDenominator, ErrPeriods.
Primary source: PoF §19.1 ("Inventory Conversion Period = 365/Inventory Turnover", average inventory); CONFIRMED.
Second: PoF §6.2 and PA1 App. A ("Days' Sales in Inventory = Ending Inventory / COGS × 365"); CONFIRMED.
Worked example: $40,000 ÷ $60,000 × 365 = 243.33 → "243 days" (PoF §6.2, Clear Lake; ending inventory passed in the inventory parameter).
Notes: the two sources differ only in which inventory figure (ending vs average) they put in; the function takes whichever the caller passes; the parameter name follows the brief. FOR THE LEAD.
Covers rows: days in inventory, inventory conversion period.

### GrossProfitMargin      (ID: FIN-BIZ-018)
Formula: (net sales − cost of goods sold) ÷ net sales.
Errors: ErrOutOfRange, ErrNegative, ErrZeroDenominator.
Primary source: PA1 §6.6 "Describe and Prepare Multi-Step and Simple Income Statements for Merchandising Companies", https://openstax.org/books/principles-financial-accounting/pages/6-6-describe-and-prepare-multi-step-and-simple-income-statements-for-merchandising-companies; CONFIRMED.
Worked example: CBS ($293,500 − $180,000) ÷ $293,500 = 0.3867 → "0.39" (PA1 §6.6).
Best approach: exact int64 subtraction, one division.
Covers rows: gross margin.

### OperatingProfitMargin      (ID: FIN-BIZ-019)
Formula: income from operations ÷ net revenues.
Errors: ErrOutOfRange, ErrNegative (sales), ErrZeroDenominator.
Primary source: SEC Beginners' Guide, "Operating Margin"; CONFIRMED (opened 2026-11-12).
Worked example: HAND-CHECKED: 2,500,000¢ ÷ 10,000,000¢ = 0.25 (no numeric example in the source).
Covers rows: operating margin.

### NetProfitMargin      (ID: FIN-BIZ-020)
Formula: net income ÷ net sales.
Primary source: PoF §6.6 ("Profit Margin"); CONFIRMED. Second: PA1 App. A; CONFIRMED.
Worked example: $35,000 ÷ $120,000 = 0.2917 → "0.29" (PoF §6.6, Clear Lake).
Covers rows: profit margin.

### EarningsPerShare      (ID: FIN-BIZ-021)
Formula: (net income − preferred dividends) ÷ weighted average common shares outstanding.
Inputs and units: netIncomeCents (any sign), preferredDividendsCents (≥ 0), weightedAverageShares (> 0, ≤ 2⁵³).
Output: cents per share, rounded half-even.
Errors: ErrOutOfRange, ErrNegative, ErrZeroDenominator, ErrOverflow.
Primary source: PoF §6.5 "Market Value Ratios", https://openstax.org/books/principles-finance/pages/6-5-market-value-ratios; CONFIRMED. Second: PA1 §14.5, https://openstax.org/books/principles-financial-accounting/pages/14-5-discuss-the-applicability-of-earnings-per-share-as-a-method-to-measure-performance; CONFIRMED.
Worked example: ($35,000 − $5,000) ÷ 8,000 = $3.75 (PoF §6.5); ($50,000 − $10,000) ÷ 5,000 = $8.00 (PA1 §14.5, Sanaron).
Best approach: exact integer division with half-even remainder test (no float), O(1), no allocation.
Covers rows: EPS.

### DilutedEarningsPerShare      (ID: FIN-BIZ-022)
Formula: (net income − preferred dividends + after-tax interest on convertible debt) ÷ (weighted average ordinary shares + dilutive potential ordinary shares). Rounded half-even to the cent.
Inputs and units: netIncomeCents (any sign), preferredDividendsCents (≥ 0; dividends on preferred shares NOT assumed converted), convertibleInterestAfterTaxCents (≥ 0), weightedAverageShares (> 0, ≤ 2⁵³), dilutiveShares (≥ 0, ≤ 2⁵³).
Output: cents per share, half-even (exact integer division).
Errors: ErrOutOfRange, ErrNegative, ErrZeroDenominator (weightedAverageShares = 0), ErrOverflow.
Primary source: IAS 33 "Earnings per Share" as adopted in the Annex to Commission Regulation (EU) 2023/1803 of 13 August 2023 (OJ L 237, 26.9.2023), https://eur-lex.europa.eu/legal-content/EN/TXT/HTML/?uri=CELEX:32023R1803; CONFIRMED (page opened in a browser 2026-11-12, IAS 33 text read). ¶31: "an entity shall adjust profit or loss attributable to ordinary equity holders of the parent entity, and the weighted average number of shares outstanding, for the effects of all dilutive potential ordinary shares". ¶33 adjusts profit by the after-tax effect of (a) dividends on dilutive potential ordinary shares, (b) "any interest recognised in the period related to dilutive potential ordinary shares", (c) other consequential changes. ¶36: shares = weighted average (¶19, ¶26) "plus the weighted average number of ordinary shares that would be issued on the conversion of all the dilutive potential ordinary shares". ¶12: basic profit is after preference dividends.
Second source: CFA Institute, "CFA Program Level I Sample Questions", Q5 Sublyme Corporation, https://www.cfainstitute.org/programs/cfa-program/cfa-program-level-i-sample-questions; CONFIRMED (opened 2026-11-12): "Diluted EPS = (Net income)/(Weighted average number of shares outstanding + New common shares that would have been issued at conversion)= USD200,000,000/[50,000,000 + (2,000,000 × 2)] = USD3.70".
Worked example: Sublyme (CFA Q5): (20,000,000,000¢, 0, 0, 50,000,000, 4,000,000) → 370.37… → 370¢. HAND-CHECKED with interest: (100,000,000 − 10,000,000 + 6,000,000)¢ ÷ (1,000,000 + 200,000) = 80¢.
Variants: treasury-stock method for options (¶45) — the caller passes the incremental shares in dilutiveShares.
Best approach: exact int64 numerator (|.| ≤ 3·2⁵³) and denominator (≤ 2⁵⁴), integer half-even division (divRoundHalfEven). O(1), no allocation.
Notes: must NOT decide dilution: the antidilution test (¶41–44) and ¶33(c) consequential changes are the caller's. FOR THE LEAD: the signature has no separate "convertible preferred dividends" input; the caller leaves those dividends out of preferredDividendsCents (equivalent to ¶33(a)).
Covers rows: diluted EPS.

### PriceToEarningsRatio      (ID: FIN-BIZ-023)
Formula: market price per share ÷ earnings per share.
Inputs and units: pricePerShareCents (≥ 0), earningsPerShareCents (any sign ≠ 0; negative EPS gives a negative ratio, documented).
Primary source: PoF §6.5; CONFIRMED. Second: PoF §11.1, SEC guide; CONFIRMED.
Worked example: $69.41 ÷ $3.75 = 18.51 (PoF §6.5); $24.00 ÷ $4.00 = 6 (PoF §11.1, Cameo Corp.); $20 ÷ $2 = 10 (SEC).
Covers rows: P/E.

### PriceToBookRatio      (ID: FIN-BIZ-024)
Formula: market price per share ÷ book value per share.
Primary source: PoF §11.1 "Multiple Approaches to Stock Valuation", https://openstax.org/books/principles-finance/pages/11-1-multiple-approaches-to-stock-valuation; CONFIRMED.
Worked example: HAND-CHECKED from PoF §6.5 numbers: $69.41 ÷ $10.00 = 6.941 (§11.1 has no numeric P/B example).
Covers rows: price to book, market to book.

### PriceToSalesRatio      (ID: FIN-BIZ-025)
Formula: PoF §11.1: share price ÷ sales per share. With N shares, (price × N) ÷ (sales) = market capitalization ÷ revenue: the brief's signature uses the whole-firm form, identical for any N > 0.
Primary source: PoF §11.1; CONFIRMED.
Worked example: HAND-CHECKED: 5,000,000,000¢ ÷ 2,000,000,000¢ = 2.5.
Notes: FOR THE LEAD: the source's form is per share; the signature's whole-firm form is the same ratio.
Covers rows: price to sales.

### BookValuePerShare      (ID: FIN-BIZ-026)
Formula: (total equity − preferred equity) ÷ total shares outstanding.
Output: cents per share, half-even.
Primary source: PoF §6.5; CONFIRMED.
Worked example: ($100,000 − $20,000) ÷ 8,000 = $10 (PoF §6.5).
Covers rows: book value per share.

### DividendsPerShare      (ID: FIN-BIZ-027)
Formula: total cash dividend ÷ shares outstanding (PA1 §14.3: total = per-share amount × outstanding shares, outstanding = issued − treasury).
Output: cents per share, half-even.
Primary source: PA1 §14.3, https://openstax.org/books/principles-financial-accounting/pages/14-3-record-transactions-and-the-effects-on-financial-statements-for-cash-dividends-property-dividends-stock-dividends-and-stock-splits; CONFIRMED.
Worked example: La Cantina 9,200 shares × $0.50 = $4,600, so $4,600 ÷ 9,200 = $0.50 (PA1 §14.3).
Covers rows: dividends per share.

### DividendPayoutRatio      (ID: FIN-BIZ-028)
Formula: common share dividends ÷ net income attributable to common shares.
Inputs and units: dividendsCents (≥ 0), netIncomeCents (> 0; net income attributable to common shares), cents.
Output: decimal (may exceed 1).
Errors: ErrOutOfRange, ErrNegative (dividends; negative net income: a payout share of a loss is rejected), ErrZeroDenominator.
Primary source: CFA Institute, "Financial Ratio List", CFA Program Level II, https://www.cfainstitute.org/sites/default/files/-/media/documents/support/programs/cfa/cfa_program_level_ii_financial_ratio_list.pdf, ratio 35: "Dividend payout ratio = Common share dividends ÷ Net income attributable to common shares"; CONFIRMED (PDF opened 2026-11-12, page ii).
Worked example: HAND-CHECKED (no numbers in the source): 25,000,000¢ ÷ 100,000,000¢ = 0.25. (PoF §17.6 https://openstax.org/books/principles-finance/pages/17-6-alternative-sources-of-funds has the numbers "$1 million in net income … pays $250,000 in dividends and retains $750,000" but does not name the ratio.)
Best approach: one division. O(1), no allocation.
Notes: FOR THE LEAD: rejecting net income ≤ 0 is a numeric-policy choice (ErrNegative / ErrZeroDenominator).
Covers rows: dividend payout ratio.

### RetentionRatio      (ID: FIN-BIZ-029)
Formula: (net income attributable to common shares − common share dividends) ÷ net income attributable to common shares = 1 − payout ratio.
Inputs and units: as 028.
Output: decimal (negative when dividends exceed earnings).
Errors: as 028.
Primary source: CFA Institute, "Financial Ratio List", CFA Program Level II, https://www.cfainstitute.org/sites/default/files/-/media/documents/support/programs/cfa/cfa_program_level_ii_financial_ratio_list.pdf, ratio 36: "Retention rate = (Net income attributable to common shares – Common share dividends) ÷ Net income attributable to common shares = 1 – Payout ratio"; CONFIRMED (PDF opened 2026-11-12). Ratio 37 uses it: "Sustainable growth rate = Retention rate × Return on equity".
Worked example: HAND-CHECKED: (100,000,000 − 25,000,000)¢ ÷ 100,000,000¢ = 0.75.
Variants: plowback ratio (same).
Best approach: exact int64 difference (within ±2⁵³), one division (correctly rounded, better than 1 − payout). O(1), no allocation.
Covers rows: retention ratio, plowback ratio.

### ReturnOnAssets      (ID: FIN-BIZ-030)
Formula: net income ÷ average total assets.
Primary source: PoF §6.6; CONFIRMED. Second: PA1 App. A; CONFIRMED.
Worked example: $35,000 ÷ $225,000 = 0.1556 → "0.16" (PoF §6.6).
Covers rows: return on assets.

### ReturnOnEquity      (ID: FIN-BIZ-031)
Formula: net income ÷ average stockholders' equity.
Inputs: averageEquityCents any sign ≠ 0 (negative equity documented).
Primary source: PA1 App. A; CONFIRMED.
Worked example: $35,000 ÷ $95,000 = 0.368 → "0.37" (PA1 App. A, Banyan).
Covers rows: return on equity.

### ReturnOnCommonEquity      (ID: FIN-BIZ-032)
Formula: (net income − preferred dividends) ÷ average common stockholders' equity.
Primary source: PoF §6.6, whose ROE worked example deducts preferred dividends: "($35,000 − $5,000) / $95,000 = 0.32"; CONFIRMED.
Worked example: (3,500,000¢ − 500,000¢) ÷ 9,500,000¢ = 0.3158 (PoF §6.6).
Notes: PoF labels it ROE; the brief's name (RCSE) separates it from FIN-BIZ-031. FOR THE LEAD.
Covers rows: return on common stockholders' equity.

### ReturnOnInvestment      (ID: FIN-BIZ-033)
REMOVED → FIN-INVEST-012 (lead's ruling, 2026-11-12): duplicate of CapitalGainsYield, the same (final − cost) ÷ cost. Never released; code moved to _delete/fin/business/. The research below is kept for the record.
Formula: income ÷ invested capital, with income = final value − cost and invested capital = cost.
Primary source: PA2 §12.3 "Evaluate an Operating Segment or a Project Using Return on Investment, Residual Income, and Economic Value Added" ("ROI equals Income divided by Average Capital Assets"), https://openstax.org/books/principles-managerial-accounting/pages/12-3-evaluate-an-operating-segment-or-a-project-using-return-on-investment-residual-income-and-economic-value-added; CONFIRMED.
Worked example: donut division income $1,000,000 on $2,850,000 → 0.35 (PA2 §12.3); as (final, cost) = (385,000,000¢, 285,000,000¢) → 0.3509.
Notes: FOR THE LEAD: the signature's (final − cost)/cost is the source's income ÷ investment with income = gain.
Covers rows: ROI, advertising ROI.

### EconomicValueAdded      (ID: FIN-BIZ-034)
Formula: after-tax income − invested capital × weighted average cost of capital.
Inputs: nopatCents (any sign), investedCapitalCents (≥ 0), costOfCapital (decimal, finite, > −1).
Output: cents, half-even (one rounding).
Errors: ErrOutOfRange, ErrNegative, ErrNotFinite, ErrRate, ErrOverflow.
Primary source: PA2 §12.3 ("EVA equals After-Tax Income minus (Invested Capital times Weighted Average Cost of Capital)"); CONFIRMED.
Worked example: $150,000 − $1,500,000 × 0.09 = $15,000 (PA2 §12.3).
Covers rows: EVA.

### FreeCashFlowToFirm      (ID: FIN-BIZ-035)
Formula: FCFF = EBIT(1 − tax rate) + Dep − FCInv − WCInv.
Inputs: ebitCents, depreciationCents, capitalExpendituresCents (= FCInv), changeInWorkingCapitalCents (= WCInv), taxRate in [0, 1].
Output: cents, half-even (one rounding of the only fractional term, exact half-even on the total).
Errors: ErrOutOfRange, ErrNotFinite, ErrRate, ErrOverflow.
Primary source: CFA-FCF, Summary ("FCFF = EBIT(1 – Tax rate) + Dep – FCInv – WCInv"); CONFIRMED (page opened 2026-11-12, formula visible).
Second: PoF §5.6 (FCF = operating cash flow − capital expenditures), https://openstax.org/books/principles-finance/pages/5-6-operating-cash-flow-and-free-cash-flow-to-the-firm-fcff; CONFIRMED (same structure, different starting line).
Worked example: HAND-CHECKED: EBIT $100,000, t = 0.25, Dep $10,000, FCInv $20,000, WCInv $5,000 → 75,000 + 10,000 − 20,000 − 5,000 = $60,000.
Covers rows: FCFF.

### FreeCashFlowToEquity      (ID: FIN-BIZ-036)
Formula: FCFE = FCFF − Int(1 − tax rate) + net borrowing.
Primary source: CFA-FCF, Summary ("FCFE = FCFF – Int(1 – Tax rate) + Net borrowing"); CONFIRMED.
Worked example: HAND-CHECKED: FCFF $60,000, interest $4,000, t = 0.25, net borrowing $2,000 → 60,000 − 3,000 + 2,000 = $59,000.
Covers rows: FCFE.

### NetWorkingCapital      (ID: FIN-BIZ-037)
Formula: current assets − current liabilities.
Output: cents (exact).
Primary source: PA1 App. A (working capital); CONFIRMED. Second: PoF §19.1; SEC guide; CONFIRMED.
Worked example: $200,000 − $100,000 = $100,000 (PA1 App. A, Banyan).
Covers rows: working capital.

### DegreeOfOperatingLeverage      (ID: FIN-BIZ-038)
Formula: contribution margin ÷ net operating income.
Primary source: PA2 §3.5 "Calculate and Interpret a Company's Margin of Safety and Operating Leverage", https://openstax.org/books/principles-managerial-accounting/pages/3-5-calculate-and-interpret-a-companys-margin-of-safety-and-operating-leverage; CONFIRMED.
Worked example: Company A $148,000 ÷ $85,000 = 1.74; Company B $210,000 ÷ $85,000 = 2.47 (PA2 §3.5).
Covers rows: degree of operating leverage.

### DegreeOfFinancialLeverage      (ID: FIN-BIZ-039)
SOURCE NOT CONFIRMED: PoF §17.4, §17 key terms and index have no DFL formula; only non-acceptable pages (LibreTexts
course not derived from the listed OpenStax books, blogs) show EBIT ÷ (EBIT − interest). Not built.
Retry (agent C2 retry, 2026-11-12): PoF 1e index (https://openstax.org/books/principles-finance/pages/index) has only "financial leverage"
(§17.1, no DFL); PoF 2e §17.4 (https://openstax.org/books/principles-finance-2e/pages/17-4-capital-structure-choices) no
formula; CFA "Company Analysis: Past and Present" refresher page
(https://www.cfainstitute.org/insights/professional-learning/refresher-readings/2026/company-analysis-past-and-present)
defines DFL in words only ("sensitivity of net income to changes in operating income"), no formula; CFA Level II
Financial Ratio List has no DFL. Still not built.

### ContributionMargin      (ID: FIN-BIZ-040)
Formula: sales − variable costs.
Output: cents (exact).
Primary source: PA2 §3.1 "Explain Contribution Margin and Calculate Contribution Margin per Unit, Contribution Margin Ratio, and Total Contribution Margin", https://openstax.org/books/principles-managerial-accounting/pages/3-1-explain-contribution-margin-and-calculate-contribution-margin-per-unit-contribution-margin-ratio-and-total-contribution-margin; CONFIRMED.
Worked example: Hicks April $50,000 − $10,000 = $40,000; per unit $100 − $20 = $80 (PA2 §3.1).
Covers rows: contribution margin.

### ContributionMarginRatio      (ID: FIN-BIZ-041)
Formula: contribution margin ÷ sales.
Primary source: PA2 §3.1; CONFIRMED.
Worked example: $80 ÷ $100 = 0.80 (PA2 §3.1, Blue Jay).
Covers rows: contribution margin ratio.

### BreakEvenUnits      (ID: FIN-BIZ-042)
Formula: total fixed costs ÷ contribution margin per unit = F ÷ (p − v).
Errors: ErrOutOfRange, ErrNegative, ErrDomain (p − v ≤ 0: no break-even).
Primary source: PA2 §3.2 "Calculate a Break-Even Point in Units and Dollars", https://openstax.org/books/principles-managerial-accounting/pages/3-2-calculate-a-break-even-point-in-units-and-dollars; CONFIRMED.
Worked example: $18,000 ÷ $80 = 225 units (Hicks); $14,000 ÷ $250 = 56 (Marshall & Hirito); $100,000 ÷ $200 = 500 (College Creations) (PA2 §3.2).
Notes: not rounded up to whole units (the caller decides).
Covers rows: break-even units.

### EconomicOrderQuantity      (ID: FIN-BIZ-043)
Formula: EOQ = √(2·D·S ÷ H); D = annual demand (units/year), S = set-up/order cost per order, H = carrying cost per unit per year. Harris's form: X = √(240·M·S ÷ C), M = units used per month, C = unit cost, with a 10 % yearly interest-and-depreciation charge on the average stock value (CX + S)/2: identical with D = 12M, H = C/10 (240 = 2 × 12 ÷ 0.10). The S in the stock value is a constant and does not move the optimum.
Inputs and units: annualDemandUnits (float64, finite, ≥ 0, units/year), orderCostCents (≥ 0), holdingCostPerUnitPerYearCents (> 0).
Output: units per order, float64, not rounded.
Errors: ErrNotFinite, ErrNegative, ErrOutOfRange, ErrZeroDenominator (H = 0).
Primary source: F. W. Harris, "What Quantity to Make at Once", ch. IV in Operation and Costs (Library of Factory Management vol. V; Shaw Factory Management Series), A. W. Shaw Co., Chicago, © 1915, pp. 47–52; HathiTrust full text uc2.ark:/13960/t02z13n82, read via https://distantreader.org/stacks/trust/uc2/uc2_ark+=13960=t02z13n82.txt; CONFIRMED (opened 2026-11-12; the HathiTrust page babel.hathitrust.org returned 403). Read: p. 49 symbols ("M equals the number of units used per month (movement)", "C equals … the unit cost", "S equals the set-up cost of an order in dollars", "X equals the unknown size of order, or lot size, which is most economical"); p. 50 "(CX+S)÷240M equals I" (10 % charge); p. 51 X "reduces to the square root of (240 MS divided by C)".
Second source (bibliographic): D. Erlenkotter, "Ford Whitman Harris's Economical Lot Size Model", UCLA Anderson School, May 30, 2013, https://escholarship.org/content/qt4s8482p2/qt4s8482p2_noSplash_5dfa07046f9562ca39dee86651320d0b.pdf; CONFIRMED (opened 2026-11-12): identifies Harris 1915 (Library of Factory Management vol. V, 47–52) as Harris's EOQ model and Harris 1913, Factory 10(2):135–136, 152, reprinted Operations Research 38(6):947–950, 1990. Crossref (https://api.crossref.org/works?query.bibliographic=…) CONFIRMED doi:10.1287/opre.38.6.937 (Erlenkotter 1990, OR 38(6):937–946) and doi:10.1287/opre.38.6.947 (Harris reprint); the 1913/1990 article text itself was not opened (INFORMS 403).
Worked example: Harris p. 51: "an article has a movement, M, of 1,000 units per month with a set-up cost of two dollars and a unit cost of ten cents … the most economical size of lot is 2,190 units" → (12,000, 200¢, 1¢) = √4,800,000 = 2,190.89… (Harris truncates).
Best approach: √(2·S/H)·√D: S/H one rounding, no overflow for any finite D (vs √(2DS/H), which overflows for D near MaxFloat64). ≤ 4 ulp. O(1), no allocation.
Notes: must NOT round to whole units. FOR THE LEAD: the H-form is an algebraic restatement of Harris's 10 %-charge formula; the cited 1915 chapter is the book version of the 1913 paper named in the brief.
Covers rows: economic order quantity.

### OptimalOrderCycle      (ID: FIN-BIZ-044)
SOURCE NOT CONFIRMED: depends on 043. Not built.
Retry (agent C2 retry, 2026-11-12): Harris 1915 (above) defines T as the "manufacturing interval" (time to make and deliver an order, i.e. a lead time), not the time between orders, and states no X ÷ M cycle formula; PoF 2e §19.5 (https://openstax.org/books/principles-finance-2e/pages/19-5-inventory-management) names EOQ with no formula. EOQ ÷ D is not stated by any acceptable page opened. Still not built.

### ReorderPoint      (ID: FIN-BIZ-045)
SOURCE NOT CONFIRMED: PoF §19.5 has no reorder-point formula; no acceptable page found. Not built.
Retry (agent C2 retry, 2026-11-12): Harris 1915 pp. 49–50 describes the "safe stock minimum, or smallest quantity the storekeeper may allow his stock to fall to before he must enter an order for more" but gives no formula (no M × T + reserve); PoF 2e §19.5 no formula; site:openstax.org search for "reorder point" found only a COBOL example (Introduction to Computer Science §7.1); other hits were blogs/Wikipedia. Still not built.

### DepreciationStraightLine      (ID: FIN-BIZ-046)
Formula: (cost − salvage value) ÷ useful life, each period.
Output: []int64 of length lifePeriods; each = ⌊(cost − salvage) ÷ n⌋ (exact integer), the last period also takes the remainder, so Σ = cost − salvage exactly.
Errors: ErrOutOfRange, ErrNegative, ErrDomain (salvage > cost), ErrPeriods (n < 1 or n > 1 000 000).
Primary source: PA1 §11.3 "Explain and Apply Depreciation Methods to Allocate Capitalized Costs", https://openstax.org/books/principles-financial-accounting/pages/11-3-explain-and-apply-depreciation-methods-to-allocate-capitalized-costs; CONFIRMED.
Worked example: Kenzie, cost $58,000, salvage $10,000, 5 years → $9,600 each year, Σ $48,000 (PA1 §11.3).
Best approach: one allocation sized n, O(n).
Notes: IRS Pub. 946 not used as a second source: it describes MACRS (conventions, switch to straight line), not book depreciation.
Covers rows: straight-line depreciation.

### DepreciationDecliningBalance      (ID: FIN-BIZ-047)
Formula: expense_t = book_{t−1} × (rateMultiple ÷ life), salvage not subtracted first; "depreciation expense is not permitted to take the book value below the estimated salvage value".
Output: []int64 length n; each expense rounded half-even, capped at book − salvage.
Errors: ErrOutOfRange, ErrNegative, ErrDomain (salvage > cost, rateMultiple ≤ 0), ErrNotFinite, ErrPeriods.
Primary source: PA1 §11.3 (double-declining-balance); CONFIRMED.
Worked example: Kenzie $58,000, salvage $10,000, 5 years, 40 %: 23,200; 13,920; 8,352; 2,528; 0 (PA1 §11.3).
Notes: the source gives no rule for a balance that has not reached salvage at the end of the life (no switch to
straight line, no final plug), so none is applied: the total is then below cost − salvage. FOR THE LEAD.
Covers rows: double declining balance, declining balance.

### GoodwillFromAcquisition      (ID: FIN-BIZ-048)
Formula: purchase price − (fair) value of net assets acquired.
Output: cents (exact); negative when the price is below the net assets (the source covers only the positive case; the function does not label the negative case).
Primary source: PA1 §11.4 "Describe Accounting for Intangible Assets and Record Related Transactions", https://openstax.org/books/principles-financial-accounting/pages/11-4-describe-accounting-for-intangible-assets-and-record-related-transactions; CONFIRMED.
Worked example: London Hoops $10 million − $7 million = $3 million (PA1 §11.4).
Covers rows: goodwill.

### BankEfficiencyRatio      (ID: FIN-BIZ-049)
Formula: total overhead (noninterest) expense ÷ (net interest income (TE) + noninterest income). UBPR shows it as a percent (PCTOF); returned here as a decimal.
Primary source: UBPR-3, item 32 "Efficiency Ratio" (UBPRE088); CONFIRMED (PDF opened 2026-11-12).
Worked example: HAND-CHECKED: 6,000,000¢ ÷ (8,000,000¢ + 2,000,000¢) = 0.6.
Covers rows: bank efficiency ratio.

### LoanToDepositRatio      (ID: FIN-BIZ-050)
Formula: net loans and leases ÷ total deposits (UBPR shows a percent).
Primary source: UBPR-1, item 32 "Net Loans & Leases to Deposits" (UBPRE600); CONFIRMED.
Worked example: HAND-CHECKED: 80,000,000¢ ÷ 100,000,000¢ = 0.8.
Covers rows: loan to deposit.

---

## FOR THE LEAD

1. Held back (SOURCE NOT CONFIRMED): 006, 009, 010, 011, 013, 039, 044, 045. Built on retry by agent C2: 008, 022, 028, 029, 043
   (new sources for doc/sources/: cfa-financial-ratio-list.md, ias-33-eu-2023-1803.md, sec-reg-sk-503d.md, harris-1915-eoq.md).
2. Ratios from UBPR (049, 050) are percents in the source; returned as decimals per the numeric policy.
3. 047: no final-period plug when the declining balance does not reach salvage (source silent). Decide.
4. 046/047 cap lifePeriods at 1 000 000 (ErrPeriods) so a huge n cannot exhaust memory: a numeric-policy
   choice to log in doc/ASSUMPTIONS.md.
5. 017 sources disagree on ending vs average inventory; 025 per-share vs whole-firm form; 032 PoF calls it ROE;
   033 ROI as gain ÷ cost.
6. Sources for doc/sources/: openstax-principles-finance.md, openstax-principles-accounting.md, sec-beginners-guide.md,
   cfa-free-cash-flow.md, ffiec-ubpr.md (rows to write from this file).

## Research: fin/economics


Date read: 2026-11-12. Every URL below was opened on that date; the quoted text was read on the page.

### ConsumerPriceIndex      (ID: FIN-ECON-001)
Formula: index = 100 · C_t / C_b, where C_t is the cost of the fixed basket in the current period and C_b its
cost in the base period. The base period has index 100 by definition.
Inputs and units: basketCostCurrentCents, basketCostBaseCents (int64 cents).
Output: float64 index points (base = 100). A percentage-style number, because the defining sources state the
index on a base of 100 (brief: "except where the defining source states a percentage (HHI, CPI)").
Errors: ErrOutOfRange (|cents| > 2⁵³); ErrNegative (either cost < 0); ErrZeroDenominator (base cost = 0).
Primary source: Greenlaw, Shapiro & MacDonald, *Principles of Economics* 3e, OpenStax (Rice University), 2022,
§22.1 "Tracking Inflation", subsection on index numbers: "The base year, by definition, has an index number
equal to 100." https://openstax.org/books/principles-economics-3e/pages/22-1-tracking-inflation
CONFIRMED: URL opened 2026-11-12.
Second source: U.S. Bureau of Labor Statistics, *Handbook of Methods*, "Consumer Price Index: Calculation",
https://www.bls.gov/opub/hom/cpi/calculation.htm — "A modified Laspeyres price index is used to aggregate basic
indexes…"; "These quantity weights remain fixed for a 1-year period…"; "most published indexes have an index
base period of φ = 1982–1984". CONFIRMED: URL opened 2026-11-12. (BLS's operational formula is the
modified-Laspeyres aggregation of basic indexes; this function is the textbook fixed-basket ratio, which is the
concept the BLS index implements. The BLS weighting/aggregation itself is not implemented.)
Worked example: OpenStax §22.1, base = Period 3 basket $107: Period 4 $117.50 → 117.50/1.07 = 109.8;
Period 1 $100 → 93.4; Period 3 → 100.0. In cents: (11750, 10700) → 109.813084…, (10000, 10700) → 93.457943…,
(10700, 10700) → 100 exactly.
Variants: chained (C-CPI-U) and BLS modified Laspeyres aggregation (not built).
Best approach: O(1), no allocation; float64(current)/float64(base)·100.
Notes: must NOT round the index (BLS publishes 3 decimals; rounding is the caller's).
Covers rows: "Consumer Price Index", "CPI" (lead's list).

### InflationRate      (ID: FIN-ECON-002)
Formula: π = I_t / I_{t−1} − 1 (decimal). OpenStax writes it as (level new − level previous)/level previous.
Inputs and units: previousIndex, currentIndex (float64 index points, any common base).
Output: float64 decimal (0.098 = 9.8 %).
Errors: ErrNotFinite (NaN/Inf); ErrNegative (currentIndex < 0); ErrZeroDenominator (previousIndex = 0);
ErrNegative (previousIndex < 0); ErrOverflow (non-finite result).
Primary source: Greenlaw, Shapiro & MacDonald, *Principles of Economics* 3e, OpenStax, 2022, §22.1:
"Level in new year – Level in previous year / Level in previous year x 100 = Percentage change"; worked:
"109.8 – 100/100 = 0.098 = 9.8%". URL as above. CONFIRMED: opened 2026-11-12.
Worked example: (100, 109.8) → 0.098 (tolerance 1e-15 relative: 109.8 is not exact in binary).
Variants: annualised monthly change (not built).
Best approach: O(1); currentIndex/previousIndex − 1 (one division, one subtraction; no cancellation beyond the
inherent one when the indexes are equal).
Notes: returns a decimal, not a percent.
Covers rows: "Inflation rate", "Inflation (CPI)".

### RealInterestRate      (ID: FIN-ECON-003)
Formula: exact Fisher relation 1 + i = (1 + r)(1 + π), so r = (1 + i)/(1 + π) − 1 = (i − π)/(1 + π).
i nominal rate per period, π inflation per the same period, r real rate.
Inputs and units: nominalRate, inflationRate (float64 decimal per the same period).
Output: float64 decimal per period.
Errors: ErrRate (either rate NaN, ±Inf or ≤ −1; NaN/Inf also wrap ErrNotFinite); ErrOverflow (non-finite result).
Primary source: Fisher, I., *Appreciation and Interest*, Publications of the American Economic Association
11(4), Macmillan for the AEA, 1896, ch. II §3, eq. (2) "1 + j = (1 + a)(1 + i)" and eq. (3) "j = i + a + ia",
p. 9 [351]: "The rate of interest in the (relatively) depreciating standard is equal to the sum of three terms,
viz., the rate of interest in the appreciating standard, the rate of appreciation itself and the product of
these two elements." Read on FRASER (Federal Reserve Bank of St. Louis):
https://fraser.stlouisfed.org/files/docs/meltzer/appreciationinte00fish.pdf CONFIRMED: PDF opened and the text
layer read 2026-11-12. Mapping: money is the depreciating standard (j = nominal i), goods the appreciating
standard (i = real r), a = inflation π.
No Crossref DOI exists for the 1896 monograph (Crossref query returned only reviews of it; the retrospective
Dimand & Gómez Betancourt, J. Econ. Perspectives 26(4):185–196, 2012, doi:10.1257/jep.26.4.185, is a pointer).
Not used: Fisher 1930 *The Theory of Interest* ch. II (econlib.org, opened 2026-11-12) states only the additive
approximation ("one point would be subtracted from, or added to"); its exact formula is not on the OLL pages.
OpenStax *Principles of Finance* §3.4 (opened 2026-11-12) gives only "Real = Nominal − Inflation".
Worked example: OpenStax *Principles of Finance* §3.4 narrative: $1,000 at 6 % gives $1,060; the TV rises 2 %
to $1,020: the exact real gain is 1060/1020 − 1 = 0.0392157 (the book quotes the approximation 4 %).
RealInterestRate(0.06, 0.02) = 0.04/1.02 = 0.0392156862745098 (HAND-CHECKED; tolerance 1e-15).
Fisher round trip: (1 + r)(1 + π) − 1 = i to 1e-15.
Variants: approximate r ≈ i − π (not built; the exact form is the decimal-safe one).
Best approach: (i − π)/(1 + π): one subtraction, one division, no cancellation of 1 + i against 1 + π.
Notes: also the "inflation-adjusted return" row.
Covers rows: "Real interest rate", "Inflation-adjusted return", "Fisher equation".

### HerfindahlHirschmanIndex      (ID: FIN-ECON-004)
Formula: HHI = Σ s_k², s_k the market share of firm k in percent (0–100). Range 0–10 000.
Inputs and units: marketSharesPercent []float64, percent points.
Output: float64 HHI points (squared percent).
Errors: ErrDomain (empty slice); ErrNotFinite (NaN/Inf share); ErrOutOfRange (share < 0 or > 100).
Primary source: U.S. Department of Justice and Federal Trade Commission, *Merger Guidelines*, December 18, 2023,
§2.1 "Guideline 1", p. 5: "The HHI is defined as the sum of the squares of the market shares; it is small when
there are many small firms and grows larger as the market becomes more concentrated, reaching 10,000 in a
market with a single firm." Footnote 13: "the HHI for a market of five equal firms is 2,000 (5 x 20² = 2,000)
and for six equal firms is 1,667 (6 x 16.67² = 1667)." Footnote on the change: "The change in HHI from a merger
of firms with shares a and b is equal to 2ab … 2 x 20 x 5 = 200."
https://www.justice.gov/d9/2023-12/2023%20Merger%20Guidelines.pdf CONFIRMED: PDF opened, text read 2026-11-12.
Worked example: [20,20,20,20,20] → 2000 exactly; six × 16.67 → 1667.3334 (the Guidelines print 1,667);
[100] → 10000; merger [20,5,75] → [25,75] changes HHI by 200.
Variants: ΔHHI = 2ab (not built; one subtraction of two HHI calls).
Best approach: one pass, Neumaier-compensated sum of squares: O(n) time, O(1) space, no allocation. Ω(n) bound
(every share is read).
Notes: the sum of shares is not checked against 100 (rounded published shares rarely total exactly 100); the
thresholds (1,800 / 100) are policy and are not encoded.
Covers rows: "Herfindahl-Hirschman Index", "HHI".

## FOR THE LEAD
- INDEX.md rows for FIN-ECON-001…004 and a `doc/sources/` entry for the 2023 Merger Guidelines, the BLS HOM CPI
  page, OpenStax Principles of Economics 3e and Fisher 1896 (FRASER copy).
- `doc/taxonomy.md`: confirm the short form ECON for fin/economics.

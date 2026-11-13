# financial: notes for the lead (A-FIN)

## Driver gaps
- DRIVER: amortization-calculator: the per-period schedule table is not captured in result.tables (only the summary table).

- DRIVER: real-estate-calculator, auto-loan-calculator, finance-calculator and loan-calculator: `fields` returns only calcSearchForm (status 200); the calculator inputs are not captured (perhaps not inside a <form>). Rows 7, 17, 22 and 50 held as BLOCKED until fixed.
- DRIVER: investment-calculator and payment-calculator: only the default tab's fields are captured; the solve-for tabs are not reachable.
- DRIVER: repayment-calculator: `run` fill of cpaybackwayamount resolves to `<table id="cpaybackwayamount">` (the `[name=..], [id=..]` locator takes the table first) and throws; prefer `[name=...]` over `[id=...]`.
- DRIVER (session 2): college-cost-calculator: result block holds only "Results:" after submit (2 runs).
- DRIVER (session 2): auto-loan and loan calculators: `run` with the form name that `fields` reports (calform, calform2, calform3) sends nothing, every field MISSING; form '' works but clicks only the first Calculate button. Rows 7 (real-estate) still shows no fields after driver update 2.

## Handoffs

## INDEX / taxonomy notes
- New: FIN-INT-043 | AnnuityPeriods | fin/interest | number of level payments that repay a PV, N = −ln(1 − i·PV/PMT)/ln(1+i) | O(1) | O(1), no allocation | NCEES, FE Reference Handbook 10.0.1, NCEES, 2020, Engineering Economics, factor table p. 230, (P/A,i,n) solved for N | interest, time-value, annuity, number-of-periods, nper, loan-term, payoff
  (second source still wanted: OpenStax Principles of Finance §8.2 does not solve for N.)
- New: FIN-INVEST-053 | CapitalizationRate | fin/invest | annual NOI / price (cap rate) | O(1) | O(1), no allocation | Office of the Comptroller of the Currency, Comptroller's Handbook: Commercial Real Estate Lending, Version 2.0, March 2022, Glossary p. 138 "Capitalization rate"; p. 43; p. 140 NOI | invest, real-estate, cap-rate, capitalization-rate, noi, rental-property, valuation
- doc/sources: the OCC CRE handbook entry (already used by FIN-LOAN-008) should add FIN-INVEST-053 (Glossary p. 138).
- Cash-on-cash return (equity dividend rate): NEEDS SOURCE; Appraisal Institute, The Appraisal of Real Estate, is the likely primary text (not opened).
- Lease payment (rows 19, 69): NEEDS SOURCE for rent charge = (adj. cap cost + residual) × money factor and MF = APR/2400. Reg M 12 CFR 1013.4(f) has the disclosure lines only. Candidate: Federal Reserve Board, "Keys to Vehicle Leasing" (federalreserve.gov/pubs/leasing, 403 from here) — Devin may be able to open it.
- New: FIN-INT-044 | PeriodicRateFromEffectiveAnnual | fin/interest | (1 + i_e)^(1/m) − 1, the inverse of EffectiveAnnualRate | O(1) | O(1), no allocation | NCEES, FE Reference Handbook 10.0.1, NCEES, 2020, Engineering Economics, p. 231 "Non-Annual Compounding", solved for r/m; QIA 3rd ed. ch. 1 §3.3 Eq. 5 | interest, rate-conversion, effective-annual-rate, periodic-rate, equivalent-rate, compounding-frequency
- Mutual fund costs (row 29): NEEDS SOURCE; SEC Form N-1A Item 3 Example instructions (expense ratio vs 5 % return, loads). sec.gov/files/formn-1a.pdf returned 404.
- Gap (rows 30, 31 likely): IRR on dated cash flows (actual/365, XIRR-style money-weighted return) is not built; candidate source CFA Institute GIPS 2020. Would live in fin/interest beside FIN-INT-039.
- New package fin/personal (taxonomy row PERSON, planned → built): doc.go header for the lead to check.
- New: FIN-PERSON-001 | WeeklyPayEquivalent | fin/personal | pay × periods a year ÷ 52, half-even | O(1) | O(1), no allocation | 29 CFR 778.113(b) "Salary for periods other than workweek", https://www.law.cornell.edu/cfr/text/29/778.113 (read 2026-11-13) | personal, salary, pay-period, weekly-wage, payroll, flsa
- New: FIN-PERSON-002 | HourlyRateFromWeeklyPay | fin/personal | weekly pay ÷ hours it covers, half-even | O(1) | O(1), no allocation | 29 CFR 778.113(a) "Weekly salary", same URL | personal, salary, hourly-rate, regular-rate, payroll, flsa
- doc/sources: new entry for 29 CFR 778.113 (FIN-PERSON-001, -002).
- Multi-debt payoff (rows 56, 57): NEEDS SOURCE for the avalanche/snowball allocation and minimum-payment roll-over; candidate CFPB consumer-education pages. Would build in fin/loan on FIN-LOAN-003's recursion.
- New: FIN-BIZ-051 | DepreciationSumOfYearsDigits | fin/business | sum-of-the-years-digits schedule, remainder in the last period | O(n)   n = lifePeriods | O(n) for the result (one allocation), O(1) extra | 26 CFR 1.167(b)-3(a) "Sum of the years-digits method", Example (1), https://www.law.cornell.edu/cfr/text/26/1.167(b)-3 (read 2026-11-13) | depreciation, sum-of-years-digits, syd, accelerated-depreciation, schedule, fixed-assets, accounting

- New: FIN-BIZ-052 | MarkupRate | fin/business | (selling price − cost) ÷ cost | O(1) | O(1), no allocation | D. Kirk, Contemporary Mathematics, OpenStax, 2023, §6.2 "Discounts, Markups, and Sales Tax", Example 6.16, https://openstax.org/books/contemporary-mathematics/pages/6-2-discounts-markups-and-sales-tax (read 2026-11-13) | markup, pricing, retail-price, cost, margin
- New: FIN-BIZ-053 | RetailPriceFromMarkup | fin/business | cost × (1 + markup rate), cents | O(1) | O(1), no allocation | same §6.2, Example 6.16 | markup, pricing, retail-price, cost
- New: FIN-BIZ-054 | SalePriceAfterDiscount | fin/business | original × (1 − discount rate), cents | O(1) | O(1), no allocation | same §6.2, Example 6.10 | discount, pricing, sale-price, percent-off
- New: FIN-BIZ-055 | TieredMarginalAmount | fin/business | marginal rate schedule on a supplied tier list, cents | O(t)   t = len(tierStartsCents) | O(1), no allocation | 26 U.S.C. § 1(a) rate table, https://www.law.cornell.edu/uscode/text/26/1 (read 2026-11-13) | tiered-rate, marginal-rate, progressive, bracket, graduated-commission, schedule
  (also usable by the SKIPPED tax rows 45/47/48 with a caller-supplied bracket slice.)
- New: FIN-ECON-005 | GDPExpenditure | fin/economics | C + I + G + (X − M), cents | O(1) | O(1), no allocation | Greenlaw, Shapiro and MacDonald, Principles of Macroeconomics 3e, OpenStax, 2022, §6.1, Table 6.1, https://openstax.org/books/principles-macroeconomics-3e/pages/6-1-measuring-the-size-of-the-economy-gross-domestic-product (read 2026-11-13) | economics, gdp, expenditure-approach, national-accounts, macroeconomics
- Handoffs received: rows 72 gdp-calculator (BUILT FIN-ECON-005) and 73 electricity-calculator (MATCH by hand). Proposal for 73: a generic `CostAtUnitPrice(quantity float64, priceCentsPerUnit float64) (int64, error)` in fin/business (quantity × price, half-even) if the lead wants a fin home for energy/fuel/utility cost results; not built (no formula beyond the product, API decision).

## _delete moves (path | replaced by | date)

## FOUND, NOT FIXED
- INCIDENT (A-FIN caused it): the session scratchpad is shared between agents. From row 3 on, my helper
  `scratchpad/setrow.py` had been replaced by A-OTHER's script of the same name, which writes to
  `other/PROGRESS.md`. My row closes for financial rows 3, 4, 5, 6, 7, 8, 9, 10, 13, 14, 15, 16, 17, 18, 19, 20, 21,
  22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 38, 40, 41, 46, 50, 51, 53, 55, 56, 59 (and a debug
  "TESTX" on row 3) were therefore written into the SAME-NUMBERED rows of other/PROGRESS.md (status + IDs columns).
  Rows there that still carry FIN- IDs, "TESTX" or my status texts are mine and should go back to their prior state
  (most likely OPEN). A restore script was refused by the permission classifier (outside my owned paths), so it is
  left for the lead / A-OTHER. financial/PROGRESS.md has since been rewritten correctly from a private script.
  Suggest: each agent keeps helpers under its own sub-folder of the scratchpad (I now use scratchpad/afin/).

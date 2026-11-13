# financial: progress (A-FIN)

One row per calculator. status: OPEN until closed as MATCH / FIXED (what, source) / BUILT <IDs> / SKIPPED (reason) / NEEDS SOURCE / HANDOFF <agent> / BLOCKED / UNRESOLVED. A closed row is skipped on relaunch.
first-pass match is from doc/calculator_net_map_2026-11-13 and must be re-checked by grep.

| # | name | slug | url | first-pass match | status | IDs | check file |
|---|---|---|---|---|---|---|---|
| 1 | Mortgage Calculator | mortgage-calculator | https://www.calculator.net/mortgage-calculator.html | HAVE FIN-LOAN-001, FIN-LOAN-007, FIN-INT-018 | MATCH (payment, PITI; total-of-payments THEIRS DIFFERS, see DISCREPANCIES) | FIN-INT-018, FIN-LOAN-001, FIN-LOAN-007 | checks/mortgage-calculator.md |
| 2 | Amortization Calculator | amortization-calculator | https://www.calculator.net/amortization-calculator.html | HAVE FIN-LOAN-001 | MATCH (payment, payoff periods; cent totals THEIRS DIFFERS by convention; DRIVER: schedule table not captured) | FIN-LOAN-001, FIN-INT-018, FIN-LOAN-005 | checks/amortization-calculator.md |
| 3 | Mortgage Payoff Calculator | mortgage-payoff-calculator | https://www.calculator.net/mortgage-payoff-calculator.html | HAVE FIN-LOAN-005, FIN-LOAN-002 | BUILT FIN-INT-043; MATCH (payoff months) | FIN-INT-043, FIN-LOAN-002, FIN-LOAN-005 | checks/mortgage-payoff-calculator.md |
| 4 | House Affordability Calculator | house-affordability-calculator | https://www.calculator.net/house-affordability-calculator.html | PARTIAL FIN-LOAN-010, FIN-LOAN-009 (DTI limits are lender policy, not a formula) | MATCH (loan from housing ratio; ratio limits are lender policy, %-of-price wrapper not built) | FIN-LOAN-010 | checks/house-affordability-calculator.md |
| 5 | Rent Calculator | rent-calculator | https://www.calculator.net/rent-calculator.html | PARTIAL FIN-LOAN-009 (rent-to-income rule of thumb only) | MATCH (DTI arithmetic; 28/36 thresholds are lender policy, not built) | FIN-LOAN-009 | checks/rent-calculator.md |
| 6 | Debt-to-Income Ratio Calculator | debt-ratio-calculator | https://www.calculator.net/debt-ratio-calculator.html | HAVE FIN-LOAN-009 | MATCH | FIN-LOAN-009 | checks/debt-ratio-calculator.md |
| 7 | Real Estate Calculator | real-estate-calculator | https://www.calculator.net/real-estate-calculator.html | PARTIAL FIN-LOAN-001, FIN-INT-018 (points and closing costs not modelled) | BLOCKED (no calculator fields after driver update 2; page shows only the search box) | PARTIAL FIN-LOAN-001, FIN-INT-018, FIN-INT-032 | checks/real-estate-calculator.md |
| 8 | Refinance Calculator | refinance-calculator | https://www.calculator.net/refinance-calculator.html | PARTIAL FIN-LOAN-002, FIN-INT-018 (no refinance break-even comparison) | MATCH (payment, APR with points; remaining term THEIRS DIFFERS; break-even not built, no definition/source) | FIN-INT-043, FIN-INT-018, FIN-INT-032, FIN-LOAN-003 | checks/refinance-calculator.md |
| 9 | Rental Property Calculator | rental-property-calculator | https://www.calculator.net/rental-property-calculator.html | PARTIAL FIN-INT-039, FIN-INT-038 (no cap rate or cash-on-cash return) | BUILT FIN-INVEST-053; MATCH (cap rate, IRR); cash-on-cash NEEDS SOURCE | FIN-INVEST-053, FIN-INT-039 | checks/rental-property-calculator.md |
| 10 | APR Calculator | apr-calculator | https://www.calculator.net/apr-calculator.html | HAVE FIN-INT-032 | MATCH (2 of 3 exact; 1 THEIRS DIFFERS by 0.0007 pp) | FIN-INT-032, FIN-INT-018 | checks/apr-calculator.md |
| 11 | FHA Loan Calculator | fha-loan-calculator | https://www.calculator.net/fha-loan-calculator.html | OUT OF SCOPE (US FHA mortgage-insurance rate tables) | SKIPPED (FHA upfront and annual MIP rates are HUD tables that change by year) |  | checks/fha-loan-calculator.md |
| 12 | VA Mortgage Calculator | va-mortgage-calculator | https://www.calculator.net/va-mortgage-calculator.html | OUT OF SCOPE (US VA funding-fee tables) | SKIPPED (VA funding-fee table changes by year and service history) |  | checks/va-mortgage-calculator.md |
| 13 | Home Equity Loan Calculator | home-equity-loan-calculator | https://www.calculator.net/home-equity-loan-calculator.html | PARTIAL FIN-LOAN-006, FIN-INT-018 (CLTV cap is lender policy) | MATCH (payment, APR, LTV/CLTV; CLTV cap is lender policy) | FIN-INT-018, FIN-INT-032, FIN-LOAN-006 | checks/home-equity-loan-calculator.md |
| 14 | HELOC Calculator | heloc-calculator | https://www.calculator.net/heloc-calculator.html | PARTIAL FIN-INT-018 (draw/interest-only then repayment phases not modelled) | MATCH (composed: FIN-INT-035 draw + FIN-INT-018 repay) | FIN-INT-035, FIN-INT-018 | checks/heloc-calculator.md |
| 15 | Down Payment Calculator | down-payment-calculator | https://www.calculator.net/down-payment-calculator.html | PARTIAL FIN-LOAN-006, FIN-LOAN-010 | MATCH (composed: FIN-LOAN-006, FIN-INT-018) | FIN-LOAN-006, FIN-INT-018 | checks/down-payment-calculator.md |
| 16 | Rent vs. Buy Calculator | rent-vs-buy-calculator | https://www.calculator.net/rent-vs-buy-calculator.html | NONE | SKIPPED (needs the yearly federal standard deduction by filing status; the rest is a scenario of FIN-INT-018/LOAN-002/INT-001) |  | checks/rent-vs-buy-calculator.md |
| 17 | Auto Loan Calculator | auto-loan-calculator | https://www.calculator.net/auto-loan-calculator.html | PARTIAL FIN-INT-018, FIN-LOAN-001 (trade-in, fees, sales tax not modelled) | MATCH (composed: FIN-INT-018; trade-in tax treatment is a state rule) | FIN-INT-018 | checks/auto-loan-calculator.md |
| 18 | Cash Back or Low Interest Calculator | cash-back-or-low-interest-calculator | https://www.calculator.net/cash-back-or-low-interest-calculator.html | PARTIAL FIN-INT-018 (offer comparison not modelled) | MATCH (composed: FIN-INT-018, FIN-LOAN-001) | FIN-INT-018, FIN-LOAN-001 | checks/cash-back-or-low-interest-calculator.md |
| 19 | Auto Lease Calculator | auto-lease-calculator | https://www.calculator.net/auto-lease-calculator.html | NONE (money factor / residual lease payment) | NEEDS SOURCE (money-factor rent charge: not in Reg M 12 CFR 1013.4(f); Fed leasing guide unreachable) |  | checks/auto-lease-calculator.md |
| 20 | Interest Calculator | interest-calculator | https://www.calculator.net/interest-calculator.html | HAVE FIN-INT-001, FIN-INT-014 | MATCH | FIN-INT-001, FIN-INT-014, FIN-INT-015, FIN-INT-002, FIN-INT-029 | checks/interest-calculator.md |
| 21 | Investment Calculator | investment-calculator | https://www.calculator.net/investment-calculator.html | PARTIAL FIN-INT-001, FIN-INT-014, FIN-INT-040 (solve-for-any-variable wrapper not built) | BUILT FIN-INT-044; MATCH (end amount; other tabs not reachable: DRIVER) | FIN-INT-044, FIN-INT-001, FIN-INT-014, FIN-INT-015 | checks/investment-calculator.md |
| 22 | Finance Calculator | finance-calculator | https://www.calculator.net/finance-calculator.html | PARTIAL FIN-INT-001, FIN-INT-002, FIN-INT-016, FIN-INT-018, FIN-INT-042 (no number-of-periods solver) | MATCH (FV tab); other tabs DRIVER | FIN-INT-001, FIN-INT-014, FIN-INT-015 | checks/finance-calculator.md |
| 23 | Compound Interest Calculator | compound-interest-calculator | https://www.calculator.net/compound-interest-calculator.html | HAVE FIN-INT-001, FIN-INT-029, FIN-INT-030 | MATCH | FIN-INT-029, FIN-INT-044 | checks/compound-interest-calculator.md |
| 24 | Interest Rate Calculator | interest-rate-calculator | https://www.calculator.net/interest-rate-calculator.html | HAVE FIN-INT-042 | MATCH | FIN-INT-042 | checks/interest-rate-calculator.md |
| 25 | Savings Calculator | savings-calculator | https://www.calculator.net/savings-calculator.html | HAVE FIN-INT-001, FIN-INT-014, FIN-INT-023 | MATCH | FIN-INT-001, FIN-INT-023 | checks/savings-calculator.md |
| 26 | Simple Interest Calculator | simple-interest-calculator | https://www.calculator.net/simple-interest-calculator.html | HAVE FIN-INT-035 | MATCH | FIN-INT-035 | checks/simple-interest-calculator.md |
| 27 | CD Calculator | cd-calculator | https://www.calculator.net/cd-calculator.html | HAVE FIN-INT-001, FIN-INT-031 | MATCH | FIN-INT-001, FIN-INT-003, FIN-INT-029, FIN-INT-030 | checks/cd-calculator.md |
| 28 | Bond Calculator | bond-calculator | https://www.calculator.net/bond-calculator.html | HAVE FIN-INVEST-002, FIN-INVEST-003 | MATCH (price; YTM THEIRS DIFFERS by 0.0001 pp; dated form: accrued interest not built) | FIN-INVEST-002, FIN-INVEST-003 | checks/bond-calculator.md |
| 29 | Mutual Fund Calculator | mutual-fund-calculator | https://www.calculator.net/mutual-fund-calculator.html | PARTIAL FIN-INT-001, FIN-INT-014 (sales loads and expense ratio not modelled) | NEEDS SOURCE (expense-ratio and load conventions: SEC Form N-1A Item 3 Example not reachable) | PARTIAL FIN-INT-001, FIN-INT-014 | checks/mutual-fund-calculator.md |
| 30 | Average Return Calculator | average-return-calculator | https://www.calculator.net/average-return-calculator.html | HAVE FIN-INVEST-035, FIN-INVEST-034 | MATCH (form 2, composed); form 1 dated money-weighted return not built | FIN-INVEST-034, FIN-INT-040, FIN-INT-044, FIN-INT-001, FIN-INT-029 | checks/average-return-calculator.md |
| 31 | IRR Calculator | irr-calculator | https://www.calculator.net/irr-calculator.html | HAVE FIN-INT-039 | MATCH | FIN-INT-039 | checks/irr-calculator.md |
| 32 | ROI Calculator | roi-calculator | https://www.calculator.net/roi-calculator.html | HAVE FIN-BIZ-033, FIN-INVEST-035 | MATCH | FIN-INVEST-010, FIN-INVEST-011, FIN-INT-040 | checks/roi-calculator.md |
| 33 | Payback Period Calculator | payback-period-calculator | https://www.calculator.net/payback-period-calculator.html | HAVE FIN-INVEST-022, FIN-INVEST-023 | MATCH (listed flows, IRR); growing-flow form THEIRS DIFFERS | FIN-INVEST-022, FIN-INVEST-023, FIN-INT-039 | checks/payback-period-calculator.md |
| 34 | Present Value Calculator | present-value-calculator | https://www.calculator.net/present-value-calculator.html | HAVE FIN-INT-002, FIN-INT-016 | MATCH | FIN-INT-002, FIN-INT-017 | checks/present-value-calculator.md |
| 35 | Future Value Calculator | future-value-calculator | https://www.calculator.net/future-value-calculator.html | HAVE FIN-INT-001, FIN-INT-014 | MATCH | FIN-INT-001, FIN-INT-014, FIN-INT-015 | checks/future-value-calculator.md |
| 36 | Retirement Calculator | retirement-calculator | https://www.calculator.net/retirement-calculator.html | PARTIAL FIN-INT-001, FIN-INT-014, FIN-INT-016, FIN-INT-024 (no full retirement-plan model) | MATCH (form 4 composed); plan forms not built (D-021 fin/personal, no single source) | FIN-INT-043, FIN-INT-044, FIN-INT-018 | checks/retirement-calculator.md |
| 37 | 401K Calculator | 401k-calculator | https://www.calculator.net/401k-calculator.html | OUT OF SCOPE (IRS contribution limits and match rules) | SKIPPED (IRS contribution limits and employer-match rules; growth maths is FIN-INT-001/014) |  | checks/401k-calculator.md |
| 38 | Pension Calculator | pension-calculator | https://www.calculator.net/pension-calculator.html | PARTIAL FIN-INT-002, FIN-INT-016 (lump-sum vs annuity comparison not built) | MATCH (break-even age, COLA 0); COLA/joint/early-late forms not built (D-021) | FIN-INT-043, FIN-INT-044 | checks/pension-calculator.md |
| 39 | Social Security Calculator | social-security-calculator | https://www.calculator.net/social-security-calculator.html | OUT OF SCOPE (SSA benefit tables) | SKIPPED (SSA bend points and benefit tables change by year) |  | checks/social-security-calculator.md |
| 40 | Annuity Calculator | annuity-calculator | https://www.calculator.net/annuity-calculator.html | HAVE FIN-INT-014, FIN-INT-015 | MATCH | FIN-INT-001, FIN-INT-014, FIN-INT-015, FIN-INT-044 | checks/annuity-calculator.md |
| 41 | Annuity Payout Calculator | annuity-payout-calculator | https://www.calculator.net/annuity-payout-calculator.html | HAVE FIN-INT-018, FIN-INT-020 | MATCH | FIN-INT-018, FIN-INT-044 | checks/annuity-payout-calculator.md |
| 42 | Roth IRA Calculator | roth-ira-calculator | https://www.calculator.net/roth-ira-calculator.html | OUT OF SCOPE (IRS limits and tax rules) | SKIPPED (IRS contribution limits and tax rules) |  | checks/roth-ira-calculator.md |
| 43 | IRA Calculator | ira-calculator | https://www.calculator.net/ira-calculator.html | OUT OF SCOPE (IRS limits and tax rules) | SKIPPED (IRS contribution limits and tax rules) |  | checks/ira-calculator.md |
| 44 | RMD Calculator | rmd-calculator | https://www.calculator.net/rmd-calculator.html | OUT OF SCOPE (IRS life-expectancy tables) | SKIPPED (IRS Uniform Lifetime Table divisors) |  | checks/rmd-calculator.md |
| 45 | Income Tax Calculator | tax-calculator | https://www.calculator.net/tax-calculator.html | OUT OF SCOPE (tax tables) | SKIPPED (federal and state tax brackets change by year) |  | checks/tax-calculator.md |
| 46 | Salary Calculator | salary-calculator | https://www.calculator.net/salary-calculator.html | NONE (pay-period conversion) | BUILT FIN-PERSON-001, FIN-PERSON-002; MATCH (unadjusted); days-off adjustment NEEDS SOURCE | FIN-PERSON-001, FIN-PERSON-002 | checks/salary-calculator.md |
| 47 | Marriage Tax Calculator | marriage-calculator | https://www.calculator.net/marriage-calculator.html | OUT OF SCOPE (tax tables) | SKIPPED (tax brackets by filing status change by year) |  | checks/marriage-calculator.md |
| 48 | Estate Tax Calculator | estate-tax-calculator | https://www.calculator.net/estate-tax-calculator.html | OUT OF SCOPE (tax tables) | SKIPPED (estate exemption and rate schedule change by year) |  | checks/estate-tax-calculator.md |
| 49 | Take-Home-Paycheck Calculator | take-home-pay-calculator | https://www.calculator.net/take-home-pay-calculator.html | OUT OF SCOPE (tax tables) | SKIPPED (federal, state and FICA withholding tables change by year) |  | checks/take-home-pay-calculator.md |
| 50 | Loan Calculator | loan-calculator | https://www.calculator.net/loan-calculator.html | HAVE FIN-INT-018, FIN-INT-019, FIN-LOAN-003, FIN-LOAN-004 | MATCH | FIN-INT-018, FIN-INT-001, FIN-INT-002, FIN-INT-044 | checks/loan-calculator.md |
| 51 | Payment Calculator | payment-calculator | https://www.calculator.net/payment-calculator.html | PARTIAL FIN-INT-018 (fixed-payment term solve not built) | MATCH (fixed term); fixed-payment tab DRIVER gap, ours FIN-INT-043 | FIN-INT-018, FIN-INT-043 | checks/payment-calculator.md |
| 52 | Currency Calculator | currency-calculator | https://www.calculator.net/currency-calculator.html | OUT OF SCOPE (live exchange rates) | SKIPPED (live exchange rates) |  | checks/currency-calculator.md |
| 53 | Inflation Calculator | inflation-calculator | https://www.calculator.net/inflation-calculator.html | PARTIAL FIN-ECON-002 (theirs uses the BLS CPI history table) | MATCH (flat-rate forms); CPI-history form needs the BLS table (data), formula FIN-ECON-002 | FIN-INT-001, FIN-INT-002, FIN-ECON-002 | checks/inflation-calculator.md |
| 54 | Sales Tax Calculator | sales-tax-calculator | https://www.calculator.net/sales-tax-calculator.html | OUT OF SCOPE (jurisdiction rates) | SKIPPED (jurisdiction rate tables; with a supplied rate it is one percentage, no formula to cite) |  | checks/sales-tax-calculator.md |
| 55 | Credit Card Calculator | credit-card-calculator | https://www.calculator.net/credit-card-calculator.html | PARTIAL FIN-INT-018, FIN-LOAN-005 (minimum-payment rules not modelled) | MATCH (months, payment; interest cents by convention) | FIN-INT-043, FIN-LOAN-003, FIN-INT-018 | checks/credit-card-calculator.md |
| 56 | Credit Cards Payoff Calculator | credit-card-payoff-calculator | https://www.calculator.net/credit-card-payoff-calculator.html | NONE (multi-card avalanche/snowball) | NEEDS SOURCE (avalanche/snowball allocation rule; per-card recursion is FIN-LOAN-003) |  | checks/credit-card-payoff-calculator.md |
| 57 | Debt Payoff Calculator | debt-payoff-calculator | https://www.calculator.net/debt-payoff-calculator.html | NONE (avalanche/snowball schedule) | NEEDS SOURCE (multi-debt allocation rule, as row 56; single debt is FIN-INT-043) |  | checks/debt-payoff-calculator.md |
| 58 | Debt Consolidation Calculator | debt-consolidation-calculator | https://www.calculator.net/debt-consolidation-calculator.html | PARTIAL FIN-INT-018, FIN-INT-032 (comparison not built) | MATCH (loan APR, combined APR); payment roll-over rule NEEDS SOURCE | FIN-INT-032, FIN-INT-039, FIN-INT-043, FIN-LOAN-003 | checks/debt-consolidation-calculator.md |
| 59 | Repayment Calculator | repayment-calculator | https://www.calculator.net/repayment-calculator.html | HAVE FIN-INT-018, FIN-LOAN-001 | MATCH (fixed term); fixed-amount mode DRIVER gap (id collision), ours FIN-INT-043 | FIN-INT-018, FIN-INT-043 | checks/repayment-calculator.md |
| 60 | Student Loan Calculator | student-loan-calculator | https://www.calculator.net/student-loan-calculator.html | PARTIAL FIN-INT-018 (US repayment-plan rules not modelled) | MATCH (payment, payoff months, grace accrual); in-school capitalization rules not built | FIN-INT-018, FIN-INT-043, FIN-INT-001 | checks/student-loan-calculator.md |
| 61 | College Cost Calculator | college-cost-calculator | https://www.calculator.net/college-cost-calculator.html | NONE | BLOCKED (DRIVER: result block empty after submit); composes from FIN-INT-001/019/023 | PARTIAL FIN-INT-001, FIN-INT-019, FIN-INT-023 | checks/college-cost-calculator.md |
| 62 | VAT Calculator | vat-calculator | https://www.calculator.net/vat-calculator.html | OUT OF SCOPE (jurisdiction rates) | SKIPPED (jurisdiction rate tables; with a supplied rate it is one percentage, no formula to cite) |  | checks/vat-calculator.md |
| 63 | Depreciation Calculator | depreciation-calculator | https://www.calculator.net/depreciation-calculator.html | PARTIAL FIN-BIZ-046, FIN-BIZ-047 (sum-of-years and MACRS not built) | BUILT FIN-BIZ-051; MATCH (SYD, DDB); MACRS and partial-year conventions not built | FIN-BIZ-051, FIN-BIZ-046, FIN-BIZ-047 | checks/depreciation-calculator.md |
| 64 | Margin Calculator | margin-calculator | https://www.calculator.net/margin-calculator.html | PARTIAL FIN-BIZ-018 (markup not built) | BUILT FIN-BIZ-052, FIN-BIZ-053; MATCH | FIN-BIZ-052, FIN-BIZ-053, FIN-BIZ-018 | checks/margin-calculator.md |
| 65 | Discount Calculator | discount-calculator | https://www.calculator.net/discount-calculator.html | NONE | BUILT FIN-BIZ-054; MATCH | FIN-BIZ-054 | checks/discount-calculator.md |
| 66 | Business Loan Calculator | business-loan-calculator | https://www.calculator.net/business-loan-calculator.html | PARTIAL FIN-INT-018, FIN-INT-032 (fee handling not built) | MATCH (payment; APR 1 of 2 THEIRS DIFFERS by 0.0016 pp) | FIN-INT-018, FIN-INT-032 | checks/business-loan-calculator.md |
| 67 | Personal Loan Calculator | personal-loan-calculator | https://www.calculator.net/personal-loan-calculator.html | PARTIAL FIN-INT-018, FIN-LOAN-001 (origination fee not modelled) | MATCH (payment; APR 1 of 2 THEIRS DIFFERS by 0.0024 pp) | FIN-INT-018, FIN-INT-032 | checks/personal-loan-calculator.md |
| 68 | Boat Loan Calculator | boat-loan-calculator | https://www.calculator.net/boat-loan-calculator.html | PARTIAL FIN-INT-018 (fees and sales tax not modelled) | MATCH (composed: FIN-INT-018) | FIN-INT-018 | checks/boat-loan-calculator.md |
| 69 | Lease Calculator | lease-calculator | https://www.calculator.net/lease-calculator.html | NONE (residual-value lease payment) | MATCH | FIN-LOAN-004 | checks/lease-calculator.md |
| 70 | Budget Calculator | budget-calculator | https://www.calculator.net/budget-calculator.html | NONE | SKIPPED (line-item sums; suggested category shares are unsourced rules of thumb) |  | checks/budget-calculator.md |
| 71 | Commission Calculator | commission-calculator | https://www.calculator.net/commission-calculator.html | NONE | BUILT FIN-BIZ-055; MATCH | FIN-BIZ-055 | checks/commission-calculator.md |
| 72 | GDP Calculator | gdp-calculator | https://www.calculator.net/gdp-calculator.html | HANDOFF from A-OTHER (expenditure identity, fin/economics) | BUILT FIN-ECON-005; MATCH | FIN-ECON-005 | checks/gdp-calculator.md |
| 73 | Electricity Calculator | electricity-calculator | https://www.calculator.net/electricity-calculator.html | HANDOFF from A-OTHER (energy cost in cents) | MATCH (hand: kWh × price; no fin function built, proposal in LEAD.md) |  | checks/electricity-calculator.md |

## Site log

Page loads and submissions made (count per calculator), any 403/429/block.
- mortgage-calculator: 3 submissions (+1 lead proof run), no block.
- amortization-calculator: 1 page load (fields) + 3 submissions, no block.
- mortgage-payoff-calculator: 1 page load + 2 submissions, no block.
- house-affordability-calculator: 1 page load + 3 submissions, no block.
- rent-calculator: 1 page load + 1 submission, no block.
- debt-ratio-calculator: 1 page load + 2 submissions, no block.
- real-estate-calculator: 2 page loads (fields twice, same result), 0 submissions; DRIVER gap.
- refinance-calculator: 1 page load + 1 submission, no block.
- rental-property-calculator: 1 page load + 2 submissions, no block.
- apr-calculator: 1 page load + 3 submissions, no block.
- home-equity-loan-calculator: 1 page load + 2 submissions, no block.
- heloc-calculator: 1 page load + 2 submissions, no block.
- down-payment-calculator: 1 page load + 1 submission, no block.
- rent-vs-buy-calculator: 1 page load, 0 submissions.
- auto-loan-calculator: 2 page loads (fields), 0 submissions; DRIVER gap.
- cash-back-or-low-interest-calculator: 1 page load + 2 submissions, no block.
- auto-lease-calculator: 1 page load, 0 submissions.
- interest-calculator: 1 page load + 2 submissions, no block.
- investment-calculator: 1 page load + 2 submissions, no block.
- finance-calculator: 1 page load (fields, empty), 0 submissions; DRIVER gap.
- compound-interest-calculator: 1 page load + 2 submissions, no block.
- interest-rate-calculator: 1 page load + 2 submissions, no block.
- savings-calculator: 1 page load + 2 submissions, no block.
- simple-interest-calculator: 1 page load + 2 submissions, no block.
- cd-calculator: 1 page load + 3 submissions, no block.
- bond-calculator: 1 page load + 2 submissions, no block.
- mutual-fund-calculator: 1 page load + 1 submission, no block.
- average-return-calculator: 1 page load + 1 submission, no block.
- irr-calculator: 1 page load + 2 submissions, no block.
- roi-calculator: 1 page load + 2 submissions, no block.
- payback-period-calculator: 1 page load + 2 submissions, no block.
- present-value-calculator: 1 page load + 3 submissions, no block.
- future-value-calculator: 1 page load + 2 submissions, no block.
- retirement-calculator: 1 page load + 1 submission, no block.
- pension-calculator: 1 page load + 1 submission, no block.
- annuity-calculator: 1 page load + 2 submissions, no block.
- annuity-payout-calculator: 1 page load + 2 submissions, no block.
- salary-calculator: 1 page load + 2 submissions, no block.
- loan-calculator: 1 page load (fields, empty), 0 submissions; DRIVER gap.
- payment-calculator: 1 page load + 1 submission, no block.
- inflation-calculator: 1 page load + 2 submissions, no block.
- repayment-calculator: 1 page load + 2 runs (second failed in driver: id collision), no block.
- credit-card-calculator: 1 page load + 2 submissions, no block.
- credit-card-payoff-calculator: 1 page load, 0 submissions.
- Session 1 (A-FIN, 2026-11-13 09:08–09:46): no 403/429/block. About 45 page loads and 70 submissions in all, at most 3 per calculator.
- Session 2: real-estate-calculator 1 page load (still no fields); auto-loan-calculator 1 page load + 3 runs.
- finance-calculator: session 2, 1 page load + 2 submissions.
- loan-calculator: session 2, 1 page load + 4 runs (3 sent nothing: named-form DRIVER gap).
- debt-payoff-calculator: 1 page load, 0 submissions.
- debt-consolidation-calculator: 2 page loads + 1 submission.
- student-loan-calculator: 1 page load + 1 submission.
- college-cost-calculator: 1 page load + 2 submissions (results not captured).
- depreciation-calculator: 1 page load + 2 submissions.
- margin-calculator: 1 page load + 2 submissions.
- discount-calculator: 1 page load + 2 submissions.
- business-loan-calculator: 1 page load + 2 submissions.
- personal-loan-calculator: 1 page load + 2 submissions.
- boat-loan-calculator: 1 page load + 1 submission.
- lease-calculator: 1 page load + 2 submissions.
- budget-calculator: 1 page load, 0 submissions.
- commission-calculator: 1 page load + 2 submissions.
- gdp-calculator (row 72, handoff): 1 page load + 2 submissions.
- electricity-calculator (row 73, handoff): 1 page load + 2 submissions.

# Research: fin/loan (agent D, D1)

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

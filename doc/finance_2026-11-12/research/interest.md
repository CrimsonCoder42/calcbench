# Research: fin/interest (FIN-INT-001…041)

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

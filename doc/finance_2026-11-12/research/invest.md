# Research: fin/invest (agent B, 2026-11-12)

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

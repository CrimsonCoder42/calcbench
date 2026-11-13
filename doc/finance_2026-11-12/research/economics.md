# Research: fin/economics (agent D, D2)

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

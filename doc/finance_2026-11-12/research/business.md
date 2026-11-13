# Research: fin/business (agent C), 2026-11-12

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

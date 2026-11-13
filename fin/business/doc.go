// Package business computes financial-statement ratios, per-share amounts, cash-flow and value measures,
// cost-volume-profit figures, the economic order quantity, depreciation schedules and two bank ratios from
// the FFIEC UBPR.
//
// ID short form: BIZ (IDs are FIN-BIZ-NNN).
//
// Conventions: money is int64 cents with |x| ≤ 2⁵³ (ErrOutOfRange otherwise); money results are rounded
// half-even to the cent (per-share amounts with exact integer arithmetic) and |result| > 2⁵³ is
// ErrOverflow. Ratios return float64 decimals (0.25 = 25 %, 2 = 2:1), never percentages, and are not
// rounded. A zero denominator is ErrZeroDenominator; a negative denominator is accepted only where it
// means something (negative equity, negative earnings), as each doc comment says. Rates and tax rates are
// float64 decimals; NaN or ±Inf is ErrNotFinite. Errors are wrapped with the parameter name; test with
// errors.Is. Every scalar function is O(1) with no allocation; the depreciation schedules are O(n) with one
// allocation.
//
// Functions:
//
//	FIN-BIZ-001  CurrentRatio                   current assets ÷ current liabilities
//	FIN-BIZ-002  QuickRatio                     (cash + marketable securities + receivables) ÷ current liabilities
//	FIN-BIZ-003  DebtRatio                      total liabilities ÷ total assets
//	FIN-BIZ-004  DebtToEquityRatio              total liabilities ÷ total equity
//	FIN-BIZ-005  EquityMultiplier               total assets ÷ total equity
//	FIN-BIZ-007  TimesInterestEarned            EBIT ÷ interest expense
//	FIN-BIZ-008  FixedChargeCoverageRatio       (EBIT + lease payments) ÷ (interest + lease payments)
//	FIN-BIZ-012  AssetTurnover                  net sales ÷ average total assets
//	FIN-BIZ-014  ReceivablesTurnover            net credit sales ÷ average receivables
//	FIN-BIZ-015  AverageCollectionPeriod        days × average receivables ÷ net credit sales
//	FIN-BIZ-016  InventoryTurnover              cost of goods sold ÷ average inventory
//	FIN-BIZ-017  DaysInInventory                days × inventory ÷ cost of goods sold
//	FIN-BIZ-018  GrossProfitMargin              (net sales − cost of goods sold) ÷ net sales
//	FIN-BIZ-019  OperatingProfitMargin          operating income ÷ net sales
//	FIN-BIZ-020  NetProfitMargin                net income ÷ net sales
//	FIN-BIZ-021  EarningsPerShare               (net income − preferred dividends) ÷ weighted average shares, cents
//	FIN-BIZ-022  DilutedEarningsPerShare        (net income − preferred + convertible interest) ÷ diluted shares, cents
//	FIN-BIZ-023  PriceToEarningsRatio           price per share ÷ earnings per share
//	FIN-BIZ-024  PriceToBookRatio               price per share ÷ book value per share
//	FIN-BIZ-025  PriceToSalesRatio              market capitalization ÷ revenue
//	FIN-BIZ-026  BookValuePerShare              (total equity − preferred equity) ÷ shares, cents
//	FIN-BIZ-027  DividendsPerShare              total dividends ÷ shares outstanding, cents
//	FIN-BIZ-028  DividendPayoutRatio            dividends ÷ net income
//	FIN-BIZ-029  RetentionRatio                 (net income − dividends) ÷ net income
//	FIN-BIZ-030  ReturnOnAssets                 net income ÷ average total assets
//	FIN-BIZ-031  ReturnOnEquity                 net income ÷ average equity
//	FIN-BIZ-032  ReturnOnCommonEquity           (net income − preferred dividends) ÷ average common equity
//	FIN-BIZ-033  REMOVED → FIN-INVEST-012       (CapitalGainsYield computes the same (final − cost) ÷ cost)
//	FIN-BIZ-034  EconomicValueAdded             after-tax income − invested capital × cost of capital, cents
//	FIN-BIZ-035  FreeCashFlowToFirm             EBIT(1 − t) + depreciation − capex − Δ working capital, cents
//	FIN-BIZ-036  FreeCashFlowToEquity           FCFF − interest(1 − t) + net borrowing, cents
//	FIN-BIZ-037  NetWorkingCapital              current assets − current liabilities, cents
//	FIN-BIZ-038  DegreeOfOperatingLeverage      contribution margin ÷ operating income
//	FIN-BIZ-040  ContributionMargin             sales − variable costs, cents
//	FIN-BIZ-041  ContributionMarginRatio        (sales − variable costs) ÷ sales
//	FIN-BIZ-042  BreakEvenUnits                 fixed costs ÷ (price − variable cost) per unit
//	FIN-BIZ-043  EconomicOrderQuantity          √(2 × annual demand × order cost ÷ holding cost), units
//	FIN-BIZ-046  DepreciationStraightLine       straight-line schedule, remainder in the last period
//	FIN-BIZ-047  DepreciationDecliningBalance   declining-balance schedule, never below salvage
//	FIN-BIZ-048  GoodwillFromAcquisition        purchase price − fair value of net assets, cents
//	FIN-BIZ-049  BankEfficiencyRatio            noninterest expense ÷ (net interest income + noninterest income)
//	FIN-BIZ-050  LoanToDepositRatio             net loans and leases ÷ total deposits
//	FIN-BIZ-051  DepreciationSumOfYearsDigits   sum-of-the-years-digits schedule, remainder in the last period
//	FIN-BIZ-052  MarkupRate                     (selling price − cost) ÷ cost
//	FIN-BIZ-053  RetailPriceFromMarkup          cost × (1 + markup rate), cents
//	FIN-BIZ-054  SalePriceAfterDiscount         original × (1 − discount rate), cents
//	FIN-BIZ-055  TieredMarginalAmount           marginal rate schedule on a supplied tier list, cents
//
// IDs 006, 009–011, 013, 039, 044 and 045 are assigned but not built: their sources could
// not be confirmed (doc/finance_2026-11-12/research/business.md).
package business

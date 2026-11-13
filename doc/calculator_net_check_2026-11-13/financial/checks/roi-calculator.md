# roi-calculator

FIELDS (form `calform`): beginbalance amount invested (1,000), endbalance amount returned (2,000), investmenttime
(date|length), investmentlength (2.5 y), beginbalanceday / endbalanceday (dates).

MATCH: HAVE FIN-INVEST-010 HoldingPeriodYield (ROI = gain / amount invested; FIN-BIZ-033 was REMOVED → FIN-INVEST-012),
FIN-INVEST-011 AnnualizedHoldingPeriodYield (whole days, 365-day year), FIN-INT-040 GrowthRatePerPeriod (CAGR over
whole periods). A fractional number of years is reached by choosing the period (half-years for 2.5 y).

COMPARE (scratch: checks/scratch/roi-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 1,000 → 2,000 over 2.5 y | ROI 100.00 %; annualized via 5 half-years 31.9508 %; via 913 days 31.9308 % | 100.00 %; 31.95 % | 0 (half-years); −0.02 pp (913 days, a day count not 2.5 × 365 = 912.5) | MATCH |
| 25,000 → 21,000 over 3 y | ROI −16.00 %; CAGR −5.6461 % | −16.00 %; −5.65 % | 0 | MATCH |

Submissions: 2.

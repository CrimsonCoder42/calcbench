# depreciation-calculator

FIELDS (form ''): cmethod (straight|declining|sum), ccost (11,000), csalvage (1,000), cyears (5), cfactor (2),
cround (yes|no), cpartical partial first year (no) + cparticalmethod convention (day|halfmonth|fullmonth|halfquarter|
fullquarter|halfyear), startday, caymonth + cayday (start of the accounting year).

MATCH: PARTIAL FIN-BIZ-046 (straight line), FIN-BIZ-047 (declining balance). Sum-of-years-digits was missing.
Partial-year conventions and MACRS (IRS tables, out of scope) not built.
RESEARCH: 26 CFR 1.167(b)-3(a) "Sum of the years-digits method", Example (1) (read 2026-11-13 at
law.cornell.edu/cfr/text/26/1.167(b)-3): fraction = remaining life including the current year ÷ sum of the years'
digits, applied to cost − salvage; $1,750, $250 salvage, 5 years → $500, $400, $300, $200, $100.
BUILD: FIN-BIZ-051 DepreciationSumOfYearsDigits, fin/business/depreciation_sum_of_years_digits.go:21.

COMPARE (scratch: checks/scratch/depreciation-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| SYD: 11,000, salvage 1,000, 5 y | 3,333.33; 2,666.67; 2,000.00; 1,333.33; 666.67 | same | 0 | MATCH |
| DDB factor 2: same asset | 4,400.00; 2,640.00; 1,584.00; 950.40; 425.60 | same | 0 | MATCH |

Submissions: 2.

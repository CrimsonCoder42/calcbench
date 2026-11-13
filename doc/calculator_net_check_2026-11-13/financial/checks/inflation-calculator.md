# inflation-calculator

FIELDS: form `calform` (CPI history): cstartingamount1 (100), cinmonth1 (1…12, 13 = annual average) + cinyear1,
coutmonth1 + coutyear1. Form `calform2` (forward flat rate): cstartingamount2 (100), cinrate2 (3 %), cinyear2 (10).
Form `calform3` (backward flat rate): cstartingamount3 (100), cinrate3 (3 %), cinyear3 (10 years ago).

MATCH: form 1 is OUT OF SCOPE for data (it needs the BLS CPI-U history table); the formula, amount × CPI_out/CPI_in,
is FIN-ECON-002 InflationRate applied to two index levels the caller supplies. Forms 2 and 3 HAVE FIN-INT-001
FutureValue and FIN-INT-002 PresentValue at the inflation rate.

COMPARE (scratch: checks/scratch/inflation-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 2,500, 3.5 %/yr, 12 years forward | 3,777.67 | 3,777.67 | 0 | MATCH |
| 2,500, 3.5 %/yr, 12 years back | 1,654.46 | 1,654.46 | 0 | MATCH |

Submissions: 2.

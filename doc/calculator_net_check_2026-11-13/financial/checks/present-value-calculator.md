# present-value-calculator

FIELDS: form `calform` (lump sum): c1futurevalue (1,000), c1yearsv N (10), c1interestratev I/Y (6 %).
Form `calform2` (annuity): cyearsv N (10), cinterestratev (6 %), ccontributeamountv PMT (100/period), ciadditionat1
(beginning|end).

MATCH: HAVE FIN-INT-002 PresentValue, FIN-INT-016 PresentValueAnnuity, FIN-INT-017 PresentValueAnnuityDue.

COMPARE (scratch: checks/scratch/present-value-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| FV 1,000, 10 periods, 6 % | 558.39 | 558.39 | 0 | MATCH |
| PMT 100 at beginning, 10, 6 % | 780.17 | 780.17 | 0 | MATCH |
| (ordinary annuity, same) | 736.01 | not submitted | — | — |

Submissions: 3 (form 2 repeated once because the first read cut the output short).

# annuity-calculator

FIELDS (form `calform`): cstartingprinciple (20,000), cannualaddition (10,000), cmonthlyaddition (0), cadditionat1
(beginning|end), cinterestrate annual growth (6 %), cyears (10).

MATCH: HAVE FIN-INT-001, FIN-INT-014, FIN-INT-015; monthly additions use FIN-INT-044 (monthly rate equivalent to the
annual growth rate).

COMPARE (scratch: checks/scratch/annuity-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 20k + 10k/yr at beginning, 6 %, 10 y | 35,816.95 + 139,716.43 = 175,533.38 | 175,533.38 | 0 | MATCH |
| 300/month at end, 7 %/yr, 20 y | 152,260.91 (effective-rate convention; 156,278.00 at nominal 7 %/12) | 152,260.91 | 0 | MATCH |

Submissions: 2.

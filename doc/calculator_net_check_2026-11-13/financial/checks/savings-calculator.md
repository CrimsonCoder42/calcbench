# savings-calculator

FIELDS (form `calform`): cstartingprinciple (20,000), cannualaddition (5,000) + cannualadditionincrease (3 %/yr),
cmonthlyaddition (0) + cmonthlyadditionincrease (0), cinterestrate (3 %), ccompound (annually…biweekly), cyears (10),
ctaxtrate (0).

MATCH: HAVE FIN-INT-001 FutureValue + FIN-INT-023 FutureValueGrowingAnnuity (end-of-year deposits growing g a year;
r = g handled). Monthly growing deposits with annual growth steps and tax on interest are not covered.

COMPARE (scratch: checks/scratch/savings-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 20k, +5k/yr growing 3 %, 3 % annual, 10 y (r = g) | 26,878.33 + 65,238.66 = 92,116.99 | 92,116.99 | 0 | MATCH |
| 10k, +2k/yr growing 5 %, 7 % annual, 20 y | 38,696.84 + 121,638.68 = 160,335.52 | 160,335.52 | 0 | MATCH |

Submissions: 2.

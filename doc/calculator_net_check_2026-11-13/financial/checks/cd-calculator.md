# cd-calculator

FIELDS (form `calform`): cstartingprinciple (10,000), cinterestrate (5 %), ccompound (annually|semiannually|quarterly|
monthly|continuously), cyears (3) + cmonths (0), ctaxtrate (0).

MATCH: HAVE FIN-INT-001 FutureValue, FIN-INT-003 FutureValueContinuous, FIN-INT-029/030 effective annual rate.
(FIN-INT-031 Reg DD APY takes interest earned and days; their page shows the EAR, not a Reg DD APY.)

COMPARE (scratch: checks/scratch/cd-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 10,000, 5 % annual, 3 y | 11,576.25 | 11,576.25 | 0 | MATCH |
| 25,000, 4.5 % monthly, 18 mo | 26,742.38; EAR 4.594 % | 26,742.38; 4.594 % | 0 | MATCH |
| 10,000, 4 % continuous, 2 y | 10,832.87; EAR 4.081 % | 10,832.87; 4.081 % | 0 | MATCH |

Submissions: 3.

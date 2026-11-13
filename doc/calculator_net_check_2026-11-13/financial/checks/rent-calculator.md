# rent-calculator

FIELDS (form `calform`): income (80,000) + incomeunit (y|m), monthlydebt (0).

MATCH: PARTIAL FIN-LOAN-009 DebtToIncomeRatio. Their two figures are (ratio × monthly income − monthly debt) at
36 % and 28 % (the mortgage 28/36 rule applied to rent + debt). The ratios are lender rules of thumb with no
primary source for rent; HUD's own measure is a 30 % "cost burden" threshold. Nothing built: the inverse is one
multiply and a subtraction, and the thresholds are policy the caller supplies.

COMPARE (scratch: checks/scratch/rent-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 80,000/yr, debt 500/mo | DTI((1,900+500), 6,666.67) = 0.3600 | max rent 1,900 | 0 | MATCH (36 %) |
| same | DTI((1,366.67+500), 6,666.67) = 0.2800 | recommended 1,367 | 0 | MATCH (28 % of rent + debt) |

Submissions: 1.

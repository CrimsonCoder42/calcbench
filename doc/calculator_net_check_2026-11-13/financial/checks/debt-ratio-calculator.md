# debt-ratio-calculator

FIELDS (form `calform`): income lines salary (60,000/y), pension, investment, otherincome; housing lines rental
(1,200/m), mortgage, propertytax (/y), hoa, insurance (/y); other debts creditcards (200/m), studentloan, autoloan
(250/m), otherloan. Each has a unit select (m|y).

MATCH: HAVE FIN-LOAN-009 DebtToIncomeRatio (monthly debts / gross monthly income). Front-end = housing only,
back-end = all debts; the caller sums the lines. Their "rating" words and the house-value estimate (fixed 7.263 %,
20 % down, 2 % costs) are not formulas we ship.

COMPARE (scratch: checks/scratch/debt-ratio-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 60,000/y; rent 1,200; cards 200; auto 250 | back 0.3300, front 0.2400 | 33 %, 24 % | 0 | MATCH |
| 95,000/y + 500/m pension; mortgage 2,100; tax 4,800/y; ins 1,200/y; auto 400; student 300 | back 0.3921, front 0.3089 | 39 %, 31 % | 0 | MATCH |

Submissions: 2.

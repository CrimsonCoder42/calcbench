# credit-card-calculator

FIELDS (form `calform`): balance (8,000), rate (18 %/yr), payoffoption (1 = fixed payment | 0 = pay off in),
fixedpaymentamount (200/month), year (2) + month (0).

MATCH: HAVE FIN-INT-043 AnnuityPeriods (months to pay off at a fixed payment; built this session), FIN-LOAN-003
BalloonBalance (balance before the last, partial payment), FIN-INT-018 (payment to clear in a set time). Their page
has no minimum-payment rule (% of balance), so none is needed here.

COMPARE (scratch: checks/scratch/credit-card-calculator.go; 1.5 %/month)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 8,000, 18 %, 200/month | N = 61.543 → 62 months = 5 y 2 mo; interest = 61 × 200 + 107.34 × 1.015 − 8,000 ≈ 4,308.95 | 5 y 2 mo; 4,308.61 | months 0; interest +0.34 | months MATCH; interest THEIRS DIFFERS (cent-rounding convention, see mortgage-calculator) |
| 8,000, 18 %, pay off in 2 y | 399.39/month | 399.39 | 0 | MATCH |

Submissions: 2.

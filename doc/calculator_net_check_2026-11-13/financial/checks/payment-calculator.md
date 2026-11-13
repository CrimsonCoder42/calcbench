# payment-calculator

FIELDS (form `calform`): cloanamount (200,000), cloanterm (15 y), cmonthlypay (2,000), cinterestrate (6 %). The page
has a fixed-term tab and a fixed-payment tab; the driver submits the fixed-term one only.

MATCH: fixed term HAVE FIN-INT-018; fixed payment (term solve) HAVE FIN-INT-043 AnnuityPeriods (built this session,
row 3).

COMPARE (scratch: checks/scratch/payment-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 200,000, 15 y, 6 % | 1,687.71 | 1,687.71 | 0 | MATCH |
| 200,000, 2,000/month, 6 % | N = 138.976 → 139 months | not reachable | — | DRIVER |

DRIVER: the fixed-payment tab is not selectable through `run` (no mode field in the form).
Submissions: 1.

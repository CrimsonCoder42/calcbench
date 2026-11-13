# annuity-payout-calculator

FIELDS (form `calform`): cstartingprinciple (500,000), cinterestrate (6 %), cyearstopayout (10) or camounttopayout
(5,000; solve for length), cpayfrequency (annually|semiannually|quarterly|monthly|semimonthly|biweekly).

MATCH: HAVE FIN-INT-018 AnnuityPaymentFromPresentValue with FIN-INT-044 (per-period rate equivalent to the annual
rate); the length tab maps to FIN-INT-043.

COMPARE (scratch: checks/scratch/annuity-payout-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 500,000, 6 %/yr, 10 y, monthly | 5,511.20 (end, effective-rate) | 5,511.20 | 0 | MATCH |
| 250,000, 4 %/yr, 20 y, quarterly | 4,531.44 | 4,531.44 | 0 | MATCH |

Totals (661,344.16; 362,515.45) follow the unrounded-payment × n convention (see mortgage-calculator).
Submissions: 2.

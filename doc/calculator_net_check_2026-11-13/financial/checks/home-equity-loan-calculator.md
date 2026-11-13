# home-equity-loan-calculator

FIELDS: form `calform`: cloanamount (150,000), cinterestrate (8 %), cyears (15), caddoptional closing costs (false),
corigamount (7,500) + corigisa (percentage|amount), corigpaid (deduct|upfront).
Form `calform2` (how much can I borrow): chousevalue (600,000), cmbalance (250,000), ltv select (90|85|80|75|70|60; 80).

MATCH: HAVE FIN-INT-018 (payment), FIN-INT-032 (APR with the deducted closing cost as a prepaid finance charge),
FIN-LOAN-006 LoanToValueRatio (LTV; CLTV is the same ratio on the sum of liens). The CLTV cap itself is lender
policy (caller-supplied); the maximum loan cap × value − balance is arithmetic, not built.

COMPARE (scratch: checks/scratch/home-equity-loan-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 150k, 8 %, 15 y | pmt 1,433.48 | 1,433.48 | 0 | MATCH |
| same, 7,500 closing deducted | APR 8.85993 % | 8.860 % | 0 | MATCH |
| value 600k, balance 250k | LTV 0.4167 | 41.7 % | 0 | MATCH |
| same, CLTV cap 85 % | CLTV(250k + 260k, 600k) = 0.8500 | max loan 260,000 | 0 | MATCH |

Total of payments (their 258,026.06) follows the same unrounded-payment × n convention as mortgage-calculator.
Submissions: 2.

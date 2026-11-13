# debt-consolidation-calculator

FIELDS (form ''): per debt k: cnm_k, cbal_k, cmpy_k, cint_k (defaults 10,000 / 260 / 17.99 %; 7,500 / 190 / 19.99 %;
6,500 / 180 / 18.99 %); consolidation loan cloanamount (25,000), cyears (5) + cmonths (0), cloancost (5) + unit (p|d).
The loan's own rate field was not listed by `fields`; the page default gives a 543.44 payment.

MATCH: PARTIAL → composed. Loan APR with the fee: FIN-INT-032. Existing debts' combined APR: FIN-INT-039 IRR on the
summed monthly flows (each debt at its own payment; FIN-INT-043 + FIN-LOAN-003 for the last partial payment).
Their payoff time assumes the 630 total keeps being paid with freed payments rolled over to the remaining debts
(stated on the page): the same unsourced allocation rule as rows 56–57, not built.

COMPARE (scratch: checks/scratch/debt-consolidation-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| loan 25,000, 60 mo, payment 543.44, fee 5 % | rate 10.99 %; APR 13.2508 % | APR 13.25 % | 0 | MATCH |
| three debts, own payments | combined APR (IRR) 18.9218 % | 18.92 % | 0 | MATCH |
| same, 630/mo pooled with roll-over | pooled at 18.9218 %: N = 58.68 → 59 months, ≈36,967 paid | 59 months; 36,963.17 | months 0 | consistent; roll-over rule NEEDS SOURCE |

Submissions: 1.

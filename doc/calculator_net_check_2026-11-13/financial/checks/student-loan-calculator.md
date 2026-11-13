# student-loan-calculator

FIELDS: `calform` (payment): cloanamount1 (30,000), cloanterm1 (10 y), cinterestrate1 (6.8 %), cmonthlypay1 (blank,
solve-for). `calform2` (payoff with extra): cloanamount3 (30,000), cmonthlypay3 (350), cinterestrate3 (6.8 %),
cpayoffoption (together|extra|original), cadditionalmonth (150), cadditionalyear (0), cadditionalonetime (0).
`calform3` (projection while in school): cgraduate (2 y), cannualamount (10,000/yr), cloanamount (20,000), cloanterm
(10), cgraceperiod (6 months), cinterestrate (6.8 %), cpayinterest (yes|no).

MATCH: PARTIAL. Forms 1–2 HAVE FIN-INT-018, FIN-INT-043. Form 3: interest accrues monthly through the 6-month grace
(FIN-INT-001) and the payment is FIN-INT-018; how disbursements during school accrue and capitalize is a federal
program rule (US repayment plans out of scope), not built.

COMPARE (scratch: checks/scratch/student-loan-calculator.go; 6.8 %/12)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 30,000, 10 y | 345.24 | 345.24 | 0 | MATCH |
| 30,000, 350/mo | N = 117.68 → 118 mo = 9 y 10 mo | 9 y 10 mo | 0 | MATCH |
| same + 150/mo | N = 73.53 → 74 mo = 6 y 2 mo | 6 y 2 mo | 0 | MATCH |
| form 3: their 44,263.99 after graduation, 6-month grace | 45,790.45 | 45,790.44 | +0.01 | MATCH (their balance is given to the cent and carries hidden cents) |
| form 3: 45,790.44 over 10 y | 526.96 | 526.96 | 0 | MATCH |

Submissions: 1 (the page recalculated all three forms).

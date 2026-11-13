# lease-calculator

FIELDS (form ''): cloanamount asset value (20,000), cresidualvalue (8,000), cloantermy (3) + cloantermm (0),
cinterestrate (6 %), cmonthlypay (400; solve-for tab).

MATCH: HAVE FIN-LOAN-004 BalloonLoanPayment: this page's lease is a level end-of-month payment that leaves the
residual as the balance (PV = PMT·(P/A) + residual·(P/F)), with no money factor or taxes. (A payment in advance
would be 403.05 in case 1; theirs is in arrears.) The money-factor vehicle lease (row 19) stays NEEDS SOURCE.

COMPARE (scratch: checks/scratch/lease-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 20,000, residual 8,000, 36 mo, 6 % | 405.06 | 405.06 | 0 | MATCH |
| 45,000, residual 25,000, 51 mo, 4.5 % | 525.33 | 525.33 | 0 | MATCH |

Submissions: 2.

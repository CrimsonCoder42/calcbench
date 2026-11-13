# loan-calculator

FIELDS (after driver update 2): form `calform` (amortized): cloanamount (100,000), cloanterm (10 y) + cloantermmonth,
cinterestrate (6 %), ccompound (annually…biweekly; monthly), cpayback (daily…year; month). `calform2` (deferred, one
lump sum at maturity): c2loanamount, c2loanterm + month, c2interestrate, c2compound (annually). `calform3` (bond, a
predetermined amount due at maturity): c3loanamount, c3loanterm + month, c3interestrate, c3compound (annually).

MATCH: HAVE FIN-INT-018 (amortized), FIN-INT-001 (amount due), FIN-INT-002 (amount received); compounding different
from the payment period uses FIN-INT-044.

COMPARE (scratch: checks/scratch/loan-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| amortized 100k, 10 y, 6 % monthly | 1,110.21 | 1,110.21 | 0 | MATCH |
| same, compounded annually, paid monthly | 1,102.24 | 1,102.24 | 0 | MATCH |
| deferred 100k, 10 y, 6 % annual | 179,084.77 | 179,084.77 | 0 | MATCH |
| bond 100k due, 10 y, 6 % annual | 55,839.48 | 55,839.48 | 0 | MATCH |

DRIVER: `run` with form `calform`, `calform2` or `calform3` sends nothing (all MISSING), as on auto-loan; with
form '' only the first form recalculates, so forms 2 and 3 were compared at their defaults.
Submissions: 4 (3 sent nothing).

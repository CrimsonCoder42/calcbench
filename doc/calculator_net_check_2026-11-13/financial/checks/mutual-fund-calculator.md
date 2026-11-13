# mutual-fund-calculator

FIELDS (form `calform`): cinvestamt (20,000), cannualadd (0), cmonthlyadd (1,000), creturnrate (5 %), cyears (5) +
cmonths (0), csalecharge (2 %, front load), cdefercharge (0 %, back load), coperatecost (0.5 %/yr expense ratio).

MATCH: PARTIAL FIN-INT-001, FIN-INT-014. Theirs treats the expense ratio as subtracted from the return
(5 % − 1 % → 4 % a year), shown below; whether loads come off each contribution or the balance is a convention.
RESEARCH: the primary source for fund-cost illustration is SEC Form N-1A, Item 3 "Example" and its instructions
(5 % assumed return, expense ratio, loads). The form PDF was not reachable (sec.gov/files/formn-1a.pdf 404), so the
expense and load conventions are unconfirmed: NEEDS SOURCE. Nothing built.

COMPARE (scratch: checks/scratch/mutual-fund-calculator.go), as a check on the convention only

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 10,000, 5 %, 10 y, expense 1 %, no loads | FutureValue(10,000, 4 %, 10) = 14,802.44 | 14,802.44; net IRR 4.000 % | 0 | convention identified (r − ER), not a sourced MATCH |

Submissions: 1.

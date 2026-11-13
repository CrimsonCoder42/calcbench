# interest-rate-calculator

FIELDS (form `calform`): cloanamount (32,000), cloanterm (3 y) + cloantermmonth (0), cmonthlypay (960).

MATCH: HAVE FIN-INT-042 AnnuityRate (rate per period with FV = 0; nominal annual = 12 × monthly).

COMPARE (scratch: checks/scratch/interest-rate-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 32,000, 36 mo, 960/mo | 0.00422067/mo → 5.06480 %/yr | 5.065 % | 0 | MATCH |
| 250,000, 360 mo, 1,600/mo | 0.00551691/mo → 6.62030 %/yr | 6.620 % | 0 | MATCH |

Submissions: 2.

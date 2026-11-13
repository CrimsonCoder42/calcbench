# future-value-calculator

FIELDS (form `calform`): cyearsv N (10), cstartingprinciplev PV (1,000), cinterestratev I/Y (6 %), ccontributeamountv
PMT (100/period), ciadditionat1 (beginning|end).

MATCH: HAVE FIN-INT-001 FutureValue, FIN-INT-014 FutureValueAnnuity, FIN-INT-015 FutureValueAnnuityDue.

COMPARE (scratch: checks/scratch/future-value-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| PV 1,000, PMT 100 at end, 6 %, 10 | 1,790.85 + 1,318.08 = 3,108.93 | 3,108.93 | 0 | MATCH |
| PV 0, PMT 250 at beginning, 0.5 %, 40 | 11,094.91 | 11,094.91 | 0 | MATCH |

Submissions: 2.

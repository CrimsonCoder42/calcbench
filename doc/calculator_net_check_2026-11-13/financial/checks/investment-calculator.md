# investment-calculator

FIELDS (form `calform`; only the "end amount" tab's fields are captured, the other tabs solve for contribution,
rate, starting amount, length): ctargetamountv (1,000,000), cstartingprinciplev (20,000), cyearsv (10),
cinterestratev (6 %), ccompound (annually…biweekly), ccontributeamountv (1,000), cadditionat1 (beginning|end),
ciadditionat1 (monthly|annually).

MATCH: PARTIAL. End amount: FIN-INT-001 + FIN-INT-014/015. Monthly contributions into an account compounding
annually need the monthly rate equivalent to the annual one, (1 + i)^(1/12) − 1, which we did not have. Other tabs
map to FIN-INT-019 (contribution), FIN-INT-042 (rate), FIN-INT-002 (starting amount), FIN-INT-043 (length).
RESEARCH: NCEES FE Reference Handbook 10.0.1 (2020) p. 231, i_e = (1 + r/m)^m − 1, solved for r/m; QIA 3rd ed. ch. 1
§3.3 Eq. 5.
BUILD: FIN-INT-044 PeriodicRateFromEffectiveAnnual, fin/interest/periodic_rate_from_effective_annual.go:20.

COMPARE (scratch: checks/scratch/investment-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 20k, 10 y, 6 % annual, +1,000/month at end | j = 0.00486755; 35,816.95 + 162,473.44 = 198,290.39 | 198,290.40 | −0.01 | MATCH (two separately rounded parts; exact 198,290.396) |
| 5k, 15 y, 8 % annual, +3,000/yr at start | 15,860.85 + 87,972.85 = 103,833.70 | 103,833.69 | +0.01 | MATCH (same cause) |

DRIVER: the other tabs (solve for contribution / rate / start / length) are not exposed as separate forms, so only
the end-amount tab was compared.
Submissions: 2.

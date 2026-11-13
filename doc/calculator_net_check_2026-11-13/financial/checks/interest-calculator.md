# interest-calculator

FIELDS (form `calform`): cstartingprinciple (20,000), cannualaddition (5,000), cmonthlyaddition (0), cadditionat1
(beginning|end), cinterestrate (5 %), ccompound (annually|semiannually|quarterly|monthly|semimonthly|biweekly),
cyears (5) + cmonths (0), ctaxtrate (0), cinflationrate (3 %).

MATCH: HAVE FIN-INT-001 FutureValue, FIN-INT-014 FutureValueAnnuity, FIN-INT-015 FutureValueAnnuityDue,
FIN-INT-002 PresentValue (inflation-adjusted buying power), FIN-INT-029 EffectiveAnnualRate. Not covered: a
contribution frequency different from the compounding frequency (caller converts the rate), tax on interest.

COMPARE (scratch: checks/scratch/interest-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 20k, +5k/yr at start, 5 %, annual, 5 y | 25,525.63 + 29,009.56 = 54,535.19 | 54,535.20; contribution interest 4,009.56 | −0.01 | MATCH (each part rounded half-even before the sum; the exact sum 54,535.195 rounds to .20) |
| same, inflation 3 % | PV(54,535.19, 3 %, 5) = 47,042.53 | 47,042.54 | −0.01 | MATCH (same cause) |
| 10k, +200/mo at end, 6 %, monthly, 10 y | 18,193.97 + 32,775.87 = 50,969.84 | 50,969.84 | 0 | MATCH |
| same | EAR 6.168 % | 6.168 % | 0 | MATCH |

The 1-cent gaps come from adding two separately rounded results in the scratch, not from either function.
Submissions: 2.

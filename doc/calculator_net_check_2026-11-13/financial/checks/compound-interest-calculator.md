# compound-interest-calculator

FIELDS (form `calform`): cinterestrate (6 %), cincompound (annually|semiannually|quarterly|monthly|semimonthly|biweekly;
monthly), coutcompound (same options; annually). It converts a nominal rate between compounding frequencies.

MATCH: HAVE FIN-INT-029 EffectiveAnnualRate; the conversion to another frequency uses FIN-INT-044
PeriodicRateFromEffectiveAnnual (built this session for row 21).

COMPARE (scratch: checks/scratch/compound-interest-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 6 % monthly → annually | 6.16778 % | 6.16778 % | 0 | MATCH |
| 9 % quarterly → biweekly (26/yr) | 8.91549 % nominal; 0.34290 % per 2 weeks | 8.91549 %; 0.34290 % | 0 | MATCH |

Submissions: 2.

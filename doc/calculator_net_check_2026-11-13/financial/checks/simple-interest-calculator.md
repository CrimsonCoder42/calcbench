# simple-interest-calculator

FIELDS (form `calform`): balance (30,000; solve-for mode), principal (20,000), rate (3) + ratebase (y|m),
term (10) + termbase (y|m).

MATCH: HAVE FIN-INT-035 SimpleInterest (P·r·t, half-even to the cent).

COMPARE (scratch: checks/scratch/simple-interest-calculator.go)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 20,000, 3 %/yr, 10 y | interest 6,000.00; balance 26,000.00 | 6,000.00; 26,000.00 | 0 | MATCH |
| 1,234.56, 7.25 %/yr, 18 months | 134.26; 1,368.82 | 134.26; 1,368.82 | 0 | MATCH |

Submissions: 2.

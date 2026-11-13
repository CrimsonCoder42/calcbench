# exponent-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope.
FIELDS (1 load): calform base (default 2), exponent (default 5), fresult (solve for any one of the three).
MATCH: PARTIAL ALG-INTMTH-001 PowerBySquaring (exact int64 powers, overflow checked). Real/fractional bases and exponents and the inverse solves (root, log) are stdlib math.Pow / math.Log; no wrapper built (no maths to add).
COMPARE (scratch/exponent-calculator.go; 4 submissions)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| 2^5 | 32 | 32 | 0 | match |
| (−3)^39 | −4052555153018976267 | −4052555153018976267 | 0 | match |
| 7^22 | 3909821048582988049 | 3909821048582988049 | 0 | match |
| 10^19 | ErrOverflow (> int64) | 1.0E+19 | n/a | by design: int64 range |

Outcome: MATCH on the integer range we cover.

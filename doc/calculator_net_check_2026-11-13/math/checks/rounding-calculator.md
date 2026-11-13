# rounding-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope.
FIELDS (1 load): calform cnum (empty), cpre select {1..6 decimals, 0 ones, −1 tens, −2, −3, −6, and fractions 1/2…1/64}, cpren (number, 2; custom precision when sp = 1), cmode {nearest, up (half up), down (half down), ceil, floor, even, odd, away, toward}, sp radio {0 select, 1 input}.
MATCH: PARTIAL NUM-ALGEB-022 RoundDecimals (ties to even only).
RESEARCH: IEEE Std 754-2019 §4.3.1 (roundTiesToEven, roundTiesToAway), §4.3.2 (toward +∞, −∞, 0), §9.5 (roundTiesToZero).
BUILD: NUM-ALGEB-037 RoundDecimalsMode with type RoundMode (six IEEE modes), num/algebra/round_decimals_mode.go; tests in round_decimals_mode_test.go (all modes on exact binary ties, agreement with NUM-ALGEB-022, errors, benchmark).
NOT BUILT (NEEDS SOURCE): their "half up" (ties toward +∞), "half down" (ties toward −∞) and "half to odd" (ties to the odd neighbour) are not IEEE 754 attributes; no primary source found in the time spent (IEEE 754-2019, Muller et al. Handbook §2.2 checked from memory only). Fraction precisions (nearest 1/2…1/64) not built: RoundDecimalsMode on x·2ᵏ would cover them; left for LEAD.

COMPARE (scratch/rounding-calculator.go; 7 submissions)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| −2.5, ones, ties to even | −2 | −2 | 0 | match |
| −2.5, ones, ties away | −3 | −3 | 0 | match |
| −2.5, ones, ties toward zero | −2 | −2 | 0 | match |
| 1.125, hundredths, ties to even | 1.12 | 1.12 | 0 | match |
| −1.21, tenths, floor | −1.3 | −1.3 | 0 | match |
| 2.675, hundredths, ties away | 2.68 | 2.68 | 0 | match (2.675·100 rounds to 267.5 in binary, so both see a tie) |
| 1250, hundreds, ceiling | 1300 | 1300 | 0 | match |

# half-life-calculator (A-MATH, 2026-11-13)

CLASSIFY: in scope. The pure mathematics (exponential decay A₀·2^(−t/T½)) lives in num/logs; physics-specific wrappers (activity, nuclide tables) would be sci/** for A-OTHER, none needed here.
FIELDS (1 load): calform nt (10), n0 (100), t (50), t12 (empty) — any three give the fourth; calform2 chalflife, cmeanlifetime, cdecayconstant — any one gives the other two.
MATCH: NONE (only FIN-INT-034 DoublingTime and NUM-SERIES-019 MeanReversionHalfLife, different models).
RESEARCH: Abramson, College Algebra 2e, OpenStax 2021, §6.7 "Half-Life" (A = A₀e^{kt}, t½ = −ln 2/k), Examples 2 and 3; Ling, Sanny & Moebs, University Physics Vol. 3, OpenStax 2016, §10.3, eqs. 10.15 (λ = 0.693/T½) and 10.16 (lifetime = 1/λ).
BUILD: NUM-LOGS-005..010 in num/logs/exponential_decay.go; tests num/logs/exponential_decay_test.go (Abramson Ex. 3 bone, 13 304.65 exact vs book's 13 301 from a rounded k). New sentinels ErrOutOfRange, ErrOverflow in num/logs.

COMPARE (scratch/half-life-calculator.go; 4 submissions)

| inputs | ours | theirs | diff | verdict |
|---|---|---|---|---|
| N0 100, Nt 10, t 50 → T½ | 15.051499783199057 | 15.051499783199 | 0 at 12 dp | match |
| same → τ | 21.71472409516259 (= 50/ln 10) | 21.714724093908 | 1.3e-9 | THEIRS DIFFERS |
| N0 1, Nt 0.2, T½ 5730 → t (Abramson Ex. 3) | 13304.647983704584 | 13304.647983705 | 0 at 9 dp | match |
| N0 250, t 1000, T½ 5730 → Nt | 221.51555738584966 | 221.51555738585 | 0 at 11 dp | match |
| T½ 5730 → τ | 8266.64258429376 | 8266.6425838161 | 4.8e-7 | THEIRS DIFFERS |
| T½ 5730 → λ | 1.2096809433855938e-4 | 1.2096809434555e-4 | 7e-15 (rel 6e-11) | THEIRS DIFFERS |

DISCREPANCY: THEIRS DIFFERS. Their λ·T½ = 0.6931471806000 and τ/T½ = 1.4426950408056: they use ln 2 rounded to 10 places (0.6931471806) in the λ/τ conversions, while their T½, t and Nt use full precision. ln 2 = 0.693147180559945309… (Abramson §6.7 uses the exact ln 2; Ling eq. 10.15 shows 0.693 only as a rounded display). Ours uses math.Ln2 and stays.

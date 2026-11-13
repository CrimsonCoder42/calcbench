# Fix pass: num/trig (2026-11-12)

Fixer for three reviewer findings. Owned paths only: num/trig/**, doc/specs/s6_trig.md, this file. Public API unchanged.

## New unexported helper
- `num/trig/third_angle.go:6` `piLo = 1.2246467991473532e-16` (π − math.Pi); `third_angle.go:14` `thirdAngle(x, y) (z, sinZ)`:
  x + y formed exactly by TwoSum (s + e), z = ((math.Pi − s) + piLo) − e (math.Pi − s exact by Sterbenz when s ≥ π/2, the
  only case where z is small); sin z = sin s when s < π/2 (z near π), else sin z; z ≤ 0 → (0, 0).

## Finding 1 (MEDIUM) SolveTriangleSAS NUM-TRIG-023
- `solve_triangle_sas.go:29-46`: angle opposite the shorter side (b ≤ a) is now `atan2(b·sin C, (a − b) + 2b·sin²(C/2))`
  (line 38), from c sin B = b sin C (DLMF 4.42.4) and c cos B = a − b cos C (DLMF 4.42.6, projection formula, confirmed on
  the page) with 1 − cos C = 2 sin²(C/2) (DLMF 4.21.28, confirmed). No term cancels; the rounded c is not used. Sides beyond
  2^±500 scaled exactly (sideShift). The other angle is `thirdAngle(C, B)` (line 39); when a = b both angles are B.
  The Kahan §3 angle core is no longer used here; c is still Kahan §7 (LawOfCosinesSide).
- Doc comment rewritten (lines 3-9), Source lines changed (below).

## Finding 2 (MEDIUM) SSA NUM-TRIG-026, AAS NUM-TRIG-025, ASA NUM-TRIG-024
- `solve_triangle_ssa.go:34`, `solve_triangle_aas.go:28`, `solve_triangle_asa.go:28`: C and sin C from `thirdAngle`.
  The side stays the law of sines (c = a sin C / sin A), now well conditioned because C and sin C are accurate.
  The projection rule c = a cos B + b cos A was NOT used: it cancels when A is obtuse and a ≈ b + c (e.g. a = 2,
  b = 1.999999). Source lines of these three are unchanged (DLMF 4.42.4, A&S 4.3.148 / Kahan §8-9 for SSA).
- Doc comments updated (AAS lines 5-8, ASA lines 5-8, SSA lines 6-9).

## Finding 3 (LOW) CentralAngle NUM-TRIG-031
- `central_angle.go`: the Sinnott 1984 `Source:` line removed; lines 14-16 mention Sinnott in prose as not seen/not relied on.
  Remaining Source lines (DLMF 4.42.8 + 4.21.28; A&S 4.3.147) unchanged.

## Tests added (oracles are 256-bit big.Float, sin/cos/atan from Maclaurin series DLMF 4.19.1, 4.19.2, 4.24.3; never the code)
- `third_angle_test.go:91` TestBigOracles (oracle self-check: sin π/6, cos π/3, atan 1, sin(math.Pi) = piLo).
- `third_angle_test.go:107` TestThirdAngle: 24 cases incl. δ = 1e-3…1e-12 both orders, 4 ulps.
- `solve_triangle_sas_test.go:84` TestSolveTriangleSAS_nearPi: (2,1), (1,2), (1,1), (1e300,3e299) × δ = 1e-3…1e-12,
  all six parts relative ≤ 16 ulps (old code: 100% error at δ = 1e-10).
- `solve_triangle_ssa_test.go:110` TestSolveTriangleSSA_nearPi: a = 2, b = 1, A = π − δ, δ = 1e-3…1e-12, 16 ulps.
- `solve_triangle_aas_test.go:75`, `solve_triangle_asa_test.go:82` …_nearPi: (π − δ, δ/3) and (δ/3, π − δ), 16 ulps.
- All existing tests (Kahan Tables 1/2, needle, huge, oracles, no-allocation) pass unchanged.

## Checks
- gofmt -l num/trig: clean; go vet ./num/trig/: clean.
- `go test -count=1 ./num/trig/`: `ok  	github.com/CrimsonCoder42/calcbench/num/trig	0.291s`

## Benchmarks (out/bench_trig_fix.txt, -benchmem -benchtime=200ms, Apple M1 Pro)
| Function | ns/op | B/op | allocs |
|---|---|---|---|
| SolveTriangleAAS | 30.44 | 0 | 0 |
| SolveTriangleASA | 29.28 | 0 | 0 |
| SolveTriangleSAS | 53.37 | 0 | 0 |
| SolveTriangleSSA | 90.00 | 0 | 0 |
Time/Space stay O(1), no allocation.

## Changed Source lines, verbatim (for INDEX.md / doc/sources regeneration)
NUM-TRIG-023 SolveTriangleSAS, now (solve_triangle_sas.go:18-20):
```
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.42.5, https://dlmf.nist.gov/4.42 (read 2026-11-12)
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eqs. 4.42.4 (law of sines) and 4.42.6 (projection formula), which give tan B = b sin C / (a − b cos C), https://dlmf.nist.gov/4.42 (read 2026-11-12); eq. 4.21.28 (cos 2z = 1 − 2 sin² z), https://dlmf.nist.gov/4.21 (read 2026-11-12)
// Source:  Kahan, "Miscalculating Area and Angles of a Needle-like Triangle", lecture notes, UC Berkeley, 2014-09-04, §7 (side c), p. 6, https://people.eecs.berkeley.edu/~wkahan/Triangle.pdf (read 2026-11-12)
```
(The Kahan line lost "and §3 (angles), p. 3": §3 is no longer used by SAS.)

NUM-TRIG-031 CentralAngle: line removed:
```
// Source:  Sinnott, "Virtues of the Haversine", Sky and Telescope, 1984, 68(2):159 (original not seen; formula and locator as reproduced in the comp.infosystems.gis FAQ Q5.1, R. G. Chamberlain, 1996, and in Wikipedia "Haversine formula", both read 2026-11-12)
```
NUM-TRIG-024/025/026: Source lines unchanged.

## FOUND, NOT FIXED (outside owned paths)
- doc/sources/sinnott_1984.md and doc/sources/README.md:69 still list NUM-TRIG-031 as using Sinnott 1984; doc/sources/kahan_2014.md
  (if it lists §3 for NUM-TRIG-023) and dlmf entries need 4.42.6 / 4.21.28 for NUM-TRIG-023. notes_trig.md:97,114 still describe
  the Sinnott Source line. INDEX.md Source column for 023 and 031 to regenerate.

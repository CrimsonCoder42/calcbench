# Kahan, Miscalculating Area and Angles of a Needle-like Triangle

**Citation**: W. Kahan, "Miscalculating Area and Angles of a Needle-like Triangle (from Lecture Notes for Introductory Numerical Analysis Classes)", University of California, Berkeley, 2014-09-04, https://people.eecs.berkeley.edu/~wkahan/Triangle.pdf (read 2026-11-12).

| Taken | Locator | Used by | How checked |
|---|---|---|---|
| How to compute ∆: sort a ≥ b ≥ c, reject when c − (a − b) < 0, the parenthesised formula | §2, p. 3 | NUM-GEOM-014, -016 | `TestTriangleAreaHeron_kahanTable1`; exact big.Float oracle |
| Angle C with μ, the a ≥ b swap, "do not remove parentheses" | §3, p. 3 | NUM-TRIG-021, -022 | Table 1 |
| Table 1 "Area ∆ and Angle C" (worked examples) | §5, p. 4 | NUM-GEOM-016; NUM-TRIG-021 | Trig: all 10 numeric rows plus the error rows in law_of_cosines_angle_test.go; geometry: rows reproduced except row 7 |
| c := √((a−b)² + 4ab sin²(C/2)) | §7, p. 6 | NUM-TRIG-020, -023 (side c only) | 512-bit oracle |
| B_S, B_T, B_C and the selection B(a, A, b) with the 0.756 threshold | §8, p. 6 | NUM-TRIG-019, -026 | Table 2 |
| Table 2 (worked examples for B) | §9, p. 8 | NUM-TRIG-019, -026 | 15 rows in law_of_sines_angle_test.go |

Status: PDF read page by page 2026-11-12 (both the trig and geometry builders). Where Kahan's decimal data round to float64 the angle moves (Table 1 rows 4, 5, 7): those rows carry a wider stated tolerance and are also checked against a 512-bit oracle. Table 1 row 7 (5278.64055, 94721.35941, 99999.99996 → 0) cannot be reproduced in binary and is omitted from the geometry table test. Correction to the catalog: Table 2 is in §9 (p. 8), not §8. Kahan's B_T is used in an equal factored form beyond 2⁵⁰⁰ (the literal (a/tan A)² overflows); literal form everywhere else. Since the fix pass of 2026-11-12 (fix_trig.md, finding 1), NUM-TRIG-023 SolveTriangleSAS no longer uses §3: its angles come from DLMF 4.42.4, 4.42.6 and 4.21.28 (see nist_dlmf.md); §7 is still used for side c.

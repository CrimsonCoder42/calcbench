# Fix: num/algebra (reviewer findings, 2026-11-12)

Owned paths only: `num/algebra/**`, `doc/specs/s5_algebra.md`, this file. Public API unchanged (names, parameters, results).
Final: `gofmt -l num/algebra` clean · `go vet ./num/algebra/` ok · `go test -count=1 ./num/algebra/` →
`ok  github.com/CrimsonCoder42/calcbench/num/algebra 0.562s`.

## 1. HIGH: CubicRoots, widely spread roots (NUM-ALGEB-009)

Coefficient order: `CubicRoots(a, b, c, d)` is a·x³ + b·x² + c·x + d (the four parameters, highest power first;
D-006's ascending order applies to `coeffs []float64`). The reviewer's tuples read as x³ + 1e10x² + x + 1e-10,
1e-300x³ + x² + x + 1, and (x − 1)(x² − 1e8x + 1); all three reproduce their reported roots.

Causes found: (a) forward deflation by the largest root (old :100, the test chose forward when |y|³ ≥ |C|, the
wrong way round); (b) Viète's small roots taken directly (cancellation); (c) a third: the scaled C underflows to 0
for 1e-300x³ + …, so the scaled cubic lost the small pair; (d) a fourth: `r*r < q*q*q` and
`math.Sqrt(r*r-q*q*q)` were fused into FMAs differently on arm64, so the Cardano branch took √(negative) = NaN for
(x − 1)(x² − 1e8x + 1).

Changed (`num/algebra/cubic_roots.go`):
- :90–91: R² − Q³ computed once as `float64(r*r) - q3` (explicit conversions stop FMA fusion); branch and √ use it.
- :94–103: Viète branch keeps only the root of largest modulus (no cancellation).
- :124–135: original coefficients normalised by a power of two (largest in [½, 1)); x₁ re-polished on them.
- :136–150: deflation to a·x² + f1·x + f0: backward (f0 = −d/x₁, f1 = (f0 − c)/x₁) when |a|·|x₁|³ ≥ |d| (x₁ at least
  the geometric mean of the other two), forward otherwise. ErrOverflow if f0/f1 not finite.
- :151–162: QuadraticRoots (small root by Vieta), then Newton polish of each root on the original coefficients:
  real roots by polishCubicRoot, a conjugate pair by the new polishComplexCubicRoot (:190) and re-conjugated.
- Doc comment :16–33 rewritten. Time O(1), Space O(1) no allocation: unchanged.

Tests (`cubic_roots_test.go`): `bigCubicRoots` (:90, 8192-bit bisection + exact deflation + quadratic formula,
independent of the code) and `TestCubicRoots_widelySpreadRoots` (:160): the three cases, each root within 8u of
its own modulus, plus case 3 against its exact factorisation 1, 5e7 ± √(25e14 − 1). Results: −1e10,
−5e-11 ± 8.660254037844386e-11i; −9.999999999999999e299, −0.5 ± 0.8660254037844387i; 1.0000000000000002e-08, 1,
9.999999999999999e7. Extra check, not kept (moved to `_delete/num/algebra/zz_stress_test.go`): 3000 random cubics
with coefficients over 2^±60: worst error 2.0·u·max(1, cond) per root.

Peters & Wilkinson 1971 NOT cited: the bibliographic record is confirmed (J. Inst. Maths Applics 8:16–35, via NAG
C02 intro), but the forward/backward rule could not be read in the paper or in a confirmed secondary text. The
direction rule is stated in the doc comment with its one-line cancellation argument instead. No Source line changed.

## 2. MEDIUM: PolynomialRoots on Wilkinson-20 (NUM-ALGEB-010)

Measured cause: plain complex Horner cannot locate roots 11–17 of the float Wilkinson-20 better than ~1e-2
(rounding noise in p), and Bini's residual rule declared points 0.5 away converged. Fix
(`num/algebra/polynomial_roots.go`):
- :245 `compensatedHorner`: p and p′ by the compensated Horner scheme in complex arithmetic (Graillat &
  Ménissier-Morain 2012, Alg. 2.1, 3.1, 3.4, 5.2, 5.4), FMA TwoProduct (:279), TwoSum (:296). Our own: the error terms
  are summed in plain floating point (not AccSum), and p′'s recurrence is compensated the same way (needed: with
  plain p′ the fourfold root (x − 1)⁴ stalled at 1e-6 and did not converge).
- :211 `newtonRatio` uses it; `atNoise` = |p| ≤ (6n·u)²·Σ|aᵢ||z|ⁱ (:231, Theorem 5.5 eq. (5.17), γ̃₂ₙ ≈ 4√2·n·u ≤ 6n·u).
- :113–150: convergence test (our own, documented as such): |Δ| ≤ 4u·|z| (step applied), or |Δ| ≥ previous |Δ|
  while atNoise (step not applied); sweep cap 100 + n and ErrNoConvergence kept. A separate Newton polish was not
  added: the final Aberth sweeps already run on the original coefficients with the compensated evaluation, and a
  plain Newton step can collapse two approximations of a cluster onto one root.
- :321–330 real snap now `|p(Re z)| ≤ |p(z)|` (compensated, `logAbsValue` :306) instead of the residual rule.
- Wilkinson-20 result: every root within 1.74e-4 of k (= the float polynomial's own roots), real.
- Allocations 3 → 4 (lastStep slice). Doc Space line updated; Time unchanged O(n²·k).

Tests (`polynomial_roots_test.go`): `TestPolynomialRoots_wilkinson20` (:134): references by 8192-bit Newton from k
(asserted within 2e-4 of k); tolerance per root 10·(γ̃²·Σ|aᵢ||r|ⁱ/|p′(r)| + 4u|r|) < 1e-11, justified in the comment.
Allocation test renamed `TestPolynomialRoots_atMostFourAllocations` (:293). Stale "Bini's rule" tolerance comments
rewritten. All existing tests pass unchanged in value.

## 3. MEDIUM: citation of 0.7 and 3.8

Bini 1996 full text not readable (Springer 303 to login; Aberth 1973 AMS PDF 403). Abstract (search) confirms only
"starting approximations ... equispaced on circles of suitable radii" (Rouché) and a backward-error stop condition.
So: 3.8 and Bini's stopping rule REMOVED (replaced by the convergence test above, documented as our own, and the
compensated-Horner bound from a source read in full). 0.7 KEPT as the angle offset but relabelled as our own
choice (:14–18, doc :32): it only selects start points; any σ off the conjugation-symmetric set works and the
returned roots do not depend on it beyond the tolerance. DECISION FOR DEVIN: accept σ = 0.7 as an unsourced
engineering constant, or I replace it (no source for any specific value was found).
Bini Source line narrowed to what the abstract confirms. Spec updated (`doc/specs/s5_algebra.md`, 009 and 010).

## 4. LOW: PolynomialRoots([1e-300, 1e300])

Now returns `[0], nil`: a group whose Newton-polygon start radius underflows to 0 is returned as 0 (underflow, not
an error), documented in the doc comment (:44–47). `aberthStartPoints` (:165, :192–196) returns the count; those
roots are marked converged. ErrOverflow kept only for a radius above the range; message now "a root's modulus is
above the float64 range". Overflow guard is now n·Σ|coeffs| (was (1 + 3.8n)·Σ). Tests: the old error row removed;
`TestPolynomialRoots_underflowingRoot` (:184) checks [0] and (x − 1)(1e300x + 1e-300) → [0, 1].

## Source lines (verbatim, for INDEX.md and doc/sources)

NUM-ALGEB-009: unchanged.
NUM-ALGEB-010, changed:
`Source:  Bini, "Numerical computation of polynomial zeros by means of Aberth's method", Numerical Algorithms, 1996, 13(2):179–200, doi:10.1007/BF02207694 (start points equispaced on circles of radii from the Newton polygon)`
NUM-ALGEB-010, new:
`Source:  Graillat, Ménissier-Morain, "Accurate summation, dot product and polynomial evaluation in complex floating point arithmetic", Information and Computation, 2012, 216:57–71, doi:10.1016/j.ic.2011.09.003, Algorithms 2.1, 3.1, 3.4, 5.2, 5.4 and Theorem 5.5, eq. (5.17) (compensated Horner evaluation and its error bound)`
(confirmed: full paper PDF read, pp. 57–68, convexoptimization.com/TOOLS/Graillat.pdf, 2026-11-12). Aberth and
Ehrlich lines unchanged. New Space for 010: `O(n) for the result (one allocation), O(n) extra (hull indices, last corrections and converged flags, at most three allocations)`.
doc/sources: bini_1996.md "Taken" should drop "angle offset 0.7; stopping rule ... (1 + 3.8i)"; new
graillat_menissier_morain_2012.md (taken: Alg. 2.1, 3.1, 3.4, 5.2, 5.4, Thm 5.5; used: NUM-ALGEB-010; checked:
Wilkinson-20 and the DLMF quartic).

## Benchmarks (`out/bench_algebra_fix.txt`, M1 Pro, before → after)

| Benchmark | before | after |
|---|---|---|
| CubicRoots/scalar | 148.5 ns, 0 allocs | 213.4 ns, 0 allocs |
| PolynomialRoots/1e1 | 4.2 µs, 2 allocs | 11.7 µs, 3 allocs |
| PolynomialRoots/1e2 | 0.59 ms, 3 allocs | 1.44 ms, 4 allocs |
| PolynomialRoots/1e3 | 48.5 ms, 3 allocs | 126 ms, 4 allocs |

PolynomialRoots is about 2.6× slower: compensated evaluation (~3× the flops of Horner) plus running to |Δ| ≤ 4u|z|.
Growth 1e2 → 1e3 is 88× (was 82×), consistent with O(n²·k). Kept for correctness (go-complexity rule 7).

## MOVED TO _delete/
`_delete/num/algebra/zz_debug_test.go`, `_delete/num/algebra/zz_stress_test.go` (scratch tests).

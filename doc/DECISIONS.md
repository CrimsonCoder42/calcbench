# Decisions

Each decision: what, who, the alternative that lost and why. IDs are never reused.

## D-001 Selection is introselect (2026-11-12, Devin)
`algo/selection` uses quickselect with a median-of-3 pivot. When progress is poor, it switches to a
median-of-medians pivot. The switch rule is a 6n work budget, measured in verify (ASSUMPTIONS A-010).
Result: O(n) worst case, no randomness, and quickselect's speed on ordinary data.
Lost: median-of-medians only (also O(n) worst case, but measured 2.0x slower at 1e5); randomized quickselect
(O(n²) worst case, and it would need a `*rand.Rand` parameter to stay pure).

## D-002 NaN policy for ordered algorithms (2026-11-12, Devin)
Search, selection, heap, merge, interval, RollingMin and RollingMax order values with `cmp.Compare`. This is a
total order: NaN sorts before every other value, and -0 equals +0. It is the same order `slices.Sort` uses, so
a slice sorted by `slices.Sort` is valid input to `algo/search`, and it costs nothing extra. RollingSum and
RollingMean are arithmetic: they return `ErrNotFinite` on NaN or ±Inf (`.claude/rules/numeric.md`).
Lost: returning `ErrNaN` from every O(n) function. Search cannot scan for NaN without losing O(log n), so the
package would behave two ways, and every element would pay an extra branch. Also lost: propagating NaN as
Go's builtin `min`/`max` do, which would make the rolling functions disagree with selection on the same data.

## D-003 Window functions append to a caller-supplied dst (2026-11-12, Devin)
Signature shape: `RollingX(dst, values []T, windowLen int) ([]T, error)`. Results are appended to dst, like
`strconv.AppendInt`. Passing nil allocates once, sized up front. Passing `buf[:0]` reuses the buffer.
Lost: returning a fresh slice (one allocation on every call, and a dst variant later would cost 4 more IDs);
shipping both forms (8 functions for one convenience).

## D-004 Median returns both middle elements (2026-11-12, Devin)
`Median[T cmp.Ordered](values []T) (lower, upper T, err error)`. When n is odd, lower == upper. The function
stays generic and O(n). The arithmetic median, (lower+upper)/2, belongs to the caller or to a later
`num/stats` function.
Lost: a float64-only average (not generic, and it is `num/stats` work); the lower median only (an even-length
caller would need a second call).

## D-005 Number theory gets its own package, `num/integer` (2026-11-12, Devin)
Primes, factorisation, gcd, lcm and modular arithmetic live in `num/integer` (SHORT `INTEG`).
Lost: keeping SmallestPrimeFactor in `num/algebra`. Algebra's taxonomy line is roots, polynomials and basic
operations; number theory has its own sources (Knuth TAOCP Vol. 2 §4.5) and will grow.

## D-006 Polynomials take ascending coefficients; "no real roots" is a result (2026-11-12, Devin)
`coeffs[i]` multiplies xⁱ. A root finder that finds no real root returns a count of 0, not an error.
Lost: descending order, as written on paper and in NumPy's `polyval`. Index ≠ power invites off-by-one errors,
and appending a higher term would shift every index. Lost: an error for no real roots (it is a valid answer).

## D-007 The general polynomial root finder is Aberth–Ehrlich (2026-11-12, Devin)
NUM-ALGEB-010 uses simultaneous Aberth–Ehrlich iteration with Bini's start points: O(n²) a sweep, no
deflation, no dependency. Lost: Laguerre with deflation (accumulated deflation error, needs polishing);
companion-matrix eigenvalues (O(n³), and BLOCKED ON num/linalg); Jenkins–Traub (three stages, hard to verify,
no gain at our degrees).

## D-008 Standard-library one-liners are dropped, not wrapped (2026-11-12, Devin)
10ˣ, eˣ, eˣ−1 (`math.Expm1`) and scalb (`math.Ldexp`) get no calcbench function anywhere. Percent change and
ratio are not built now (a later Basic Operations pass, IDs from NUM-ALGEB-011).
Lost: thin wrappers for discoverability: two names for one thing.

## D-009 Imports between `num` packages follow the build order (2026-11-12, Devin)
Inside `num`, a package may import only packages earlier in this order: logs, integer, combinatorics,
algebra, trig, geometry, stats, prob. In particular `num/geometry` may import `num/algebra` and `num/trig`,
and `num/prob` imports `num/combinatorics`. Never the reverse; no cycles.
Lost: forbidding imports inside `num` (the tangent line and the law of cosines would be duplicated or forced
into the wrong package).

## D-010 Statistics conventions (2026-11-12, Devin)
- Percentile takes an explicit Hyndman–Fan method argument (types 1–9); the documented default is type 7.
  Lost: one fixed method (callers matching R, Excel or NIST would get silently different numbers).
- Population and sample forms are separate functions (PopulationVariance / SampleVariance, and the SDs).
  Lost: a `ddof` or boolean parameter (a bare flag reads ambiguously at the call site).
- Inputs are `float64`; Mode is generic and returns every tied mode, ascending. Lost: generic numeric inputs
  everywhere (float arithmetic is needed anyway); Mode returning only the first (hides ties).
- The NIST Dataplot Reference Manual is an accepted source; ISO 13528 is accepted as a second source.
- %RSD and the median absolute deviation get their own IDs. Gamma goes to combinatorics (where it is
  `math.Gamma`/`Lgamma` plus LogFactorial/LogBinomial); rank order centroid is LATER; "Sum Squared" is DROPPED;
  WeightedGeometricMean is out unless a source turns up.
- Sources fixed by Devin: population variance and SD cite Hogg, McKean, Craig 7th ed. §1.9 or Casella and
  Berger 2nd ed. §2.3; CohenD cites Cohen 1988 2nd ed. §2.2 eq. 2.2.1; MatthewsCorrelation cites Matthews 1975,
  Biochim. Biophys. Acta 405(2):442–451; grouped mean and mode cite Spiegel, Schaum's Statistics. A paywalled
  primary is citable when two secondary sources reproduce the formula; both are recorded.

## D-011 Geometry conventions (2026-11-12, Devin)
- Radians only; every angle name says `Rad`. Lost: degree twins (doubles the API; `num/trig` converts).
- A package-local `Point{X, Y float64}`. Lost: a shared Vec2 in `num/linalg` (not yet built; would block).
- Lengths are plain `float64` in one consistent unit; areas and volumes come out in that unit squared and
  cubed. Lost: `unit/length` types (`unit` may import `num`, so the reverse would be a cycle).
- "Interface Angle" is the angle between two lines (NUM-GEOM-006). "Bell-shaped" gets no ID.
- The golden ratio goes to `num/logs` as a constant (NUM-LOGS-001), with nothing else in that package yet.
- The law of cosines belongs to `num/trig`: NUM-GEOM-018 and -019 are REMOVED rows pointing at NUM-TRIG-020
  and -021. Lost: owning it in geometry (one function, one home; trig also needs it for triangle solving).
- All 82 cataloged IDs are built (80 functions after the two REMOVED rows).

## D-012 Trigonometry conventions (2026-11-12, Devin)
- Radians only; DegToRad and RadToDeg are the conversions. Lost: Deg twins of every function (double API).
- NormalizeAngle returns [0, 2π); NormalizeAngleSigned returns (−π, π]; Acot returns (−π/2, π/2].
- SolveTriangleSSA returns the number of solutions (0, 1 or 2) and two solution slots. Lost: an error for
  "no triangle" and a slice of solutions (an allocation, and the ambiguous case is a result, not a failure).
- CentralAngle (NUM-TRIG-031) lives in trig; Earth-surface distance is a later `sci/geoscience` function.

## D-013 Probability conventions (2026-11-12, Devin)
- The geometric distribution counts TRIALS (support 1, 2, …), and its function names say so. Lost: the
  failures convention (support 0, 1, …), which is X − 1.
- Plain functions per distribution (PMF, CDF, MeanVariance, Quantile), no interface. Lost: a Distribution
  interface (an allocation and indirection per call for no consumer that needs polymorphism).
- Odds conversions live in `num/prob`. Cauchy stays only if sourced.
- The log binomial coefficient comes from `num/combinatorics` (NUM-COMBI-004), built in the same run; the
  interim `math.Lgamma` route is not needed.

## D-014 Combinatorics conventions (2026-11-12, Devin)
- Exact counts are `uint64` with `ErrOverflow`, plus a `*big.Int` twin for the counts that overflow early
  (NUM-COMBI-017…022). Lost: `uint64` only (Bell(26), Catalan(37) unreachable); `big.Int` only (allocates on
  every call for the common small case).
- Enumeration (generating permutations, combinations, subsets) is a later `algo` package, not built now.
- Stirling numbers of both kinds are in, with their sources.

## D-015 Unconfirmed citations hold back only their own function (2026-11-12, Devin)
A citation marked unconfirmed in a catalog is confirmed during the build. One that cannot be confirmed holds
back that single function (listed for Devin); everything else ships. A second confirmed source can carry a
function whose first source failed. Lost: holding a whole package until every citation is confirmed in print.

## D-016 Planned homes for computer-science catalog rows (2026-11-12, Devin)
Data size goes to `unit/datasize`; performance laws and software metrics to `eng/computing`; check digits to
`algo/checksum`. (No compsci catalog existed at build time; nothing was built under these.)

## D-017 Algorithms core, S10 (2026-11-12, lead under Devin's delegation "every decision is the lead's")
The lead's decisions L1–L41, each with the alternative that lost, are in
`algorithms_core_2026-11-12/README.md`. The ones that change the map: packages `algo/sorting`, `algo/hashing`
and `algo/seq` instead of `sort`, `hash`, `iter` (standard packages a caller imports beside ours); gcd, lcm,
modular arithmetic, the sieve and primality stay in `num/integer` (D-005), requested there, not built in
`algo/intmath`; the monotonic queue's one home is `algo/window`; interview books and problem sites are not
sources (17 functions HELD); every size argument is guarded against makeslice panics (L41: n·elemSize ≤ 2⁴⁷
bytes on 64-bit, else ErrTooLarge); per-structure caps of 128 MiB where a table is built from user data
(L13, L35); one ID per type when its methods are accessors, one per method when the method is an algorithm
(L23). Lost: the requested names and homes, for the reasons in the README.


## Cross-reference: run-local decision numbers (merged 2026-11-12 by the lockdown pass)
finance F-1…F-13 = D-018…D-030 · linear algebra (unnumbered table, 15 rows, in its order) = D-031…D-045 · medicine D-M1…D-M9 = D-046…D-054 · port D-P1…D-P15 = D-055…D-069 (A-P1 = A-057) · algorithms core L1…L41 stay inside D-017 · lockdown = D-070 onward. The catalog passes' open questions (algebra Q1–Q6, geometry, trig, prob Q1–Q7, combinatorics Q1) were answered by Devin as D-005…D-014.
## D-018 OpenStax textbooks are acceptable finance sources (2026-11-12, finance lead; finance F-1)
OpenStax *Principles of Finance*, *Principles of Accounting* Vol. 1–2, *Principles of Economics* 3e, *Contemporary
Mathematics* and *Calculus* Vol. 1 (Rice University, CC BY 4.0) count as recognised textbooks. Each is a confirmed
second source when the named paid textbook (Brealey-Myers-Allen, Ross-Westerfield-Jaffe, Bodie-Kane-Marcus,
Fabozzi, Newnan) cannot be opened (D-015). Lost: citing the paid books unread (that breaks "confirm every
citation") or holding back about 100 functions.

## D-019 The finance source list is extended to official bodies and named works (2026-11-12, finance lead; finance F-2)
Accepted, each confirmed by opening it:
- regulations and agencies: Reg Z App. J and Reg DD App. A, both read on the CFPB eRegulations copy because
  ecfr.gov blocks automated reading; Treasury/TreasuryDirect, SEC, MSRB, FFIEC UBPR, OCC, HUD 4000.1, CFPB,
  Fannie Mae Selling Guide, BLS, DOJ/FTC;
- professional material: NCEES FE Reference Handbook; CFA Institute's Quantitative Investment Analysis 3e,
  ratio list, refresher readings and sample questions;
- named works: IAS 33 as adopted by Regulation (EU) 2023/1803; Damodaran, *Applied Corporate Finance* 3e;
  Harris 1915, read as public-domain full text through a mirror; original papers with Crossref-confirmed DOIs.
Lost: the brief's narrower list, which would hold back about 40 more functions.

## D-020 Society of Actuaries Exam FM material is accepted for duration and convexity (2026-11-12, finance lead; finance F-3)
A professional body's official exam notation, on the same footing as the CFA readings. Macaulay 1938 is
cited for the duration itself. Lost: holding back FIN-INVEST-005 and -006.

## D-021 A later `fin/personal` package (2026-11-12, finance lead; finance F-4)
Salary conversions and retirement or pension rows (vCalc rows c14, 42, 77, 98, 112, 141) belong in a
personal-finance package, SHORT `PERSON`. They are not built now. Lost: putting them in `fin/interest` (they
are household planning conventions, not time-value formulas) or dropping them.

## D-022 Accuracy contract for money results (2026-11-12, finance lead; finance F-5)
Money is int64 cents. A result computed in float64 (powers, exponentials, several products) is rounded half-even
once. It equals the correctly rounded exact value except within about n·ε·|x| of a half cent, and is within
1 cent for |result| ≤ 2⁴³ cents (about $88 billion). Formulas made of integer sums plus one cents × rate
product are exact over the whole ±2⁵³ range (D-026). Schedules round each row half-even by definition.
Lost: math/big everywhere (an allocation and roughly 50× the time on every call, against Devin's
zero-allocation rule for scalars); a claim of "exact to the cent" at 2⁵³, which float64 cannot keep.

## D-023 HHI lives in `fin/economics` (2026-11-12, finance lead; finance F-6)
Market concentration is defined by the DOJ/FTC Merger Guidelines, an economics and antitrust measure. Lost:
`fin/business`, as in the brief's list.

## D-024 The rate solver, IRR, NPV and APR live in `fin/interest`; import order inside `fin` (2026-11-12, finance lead; finance F-7)
`fin/interest` imports nothing from `fin`. `fin/loan` and `fin/invest` may import `fin/interest`.
`fin/business` and `fin/economics` import nothing from `fin`. The bracketed Newton/bisection solver
(Numerical Recipes 3e §9.4) is reached through InternalRateOfReturn and AnnuityRate (FIN-INT-042).
Lost: IRR in `fin/invest`, which would force `fin/interest` (APR) to import invest or to keep a second solver.

## D-025 Duplicates removed before release (2026-11-12, finance lead; finance F-8)
FIN-INVEST-024 EquivalentAnnualAnnuity becomes REMOVED → FIN-INT-018 (the same formula). FIN-BIZ-033
ReturnOnInvestment becomes REMOVED → FIN-INVEST-012 (the same (b − a)/a). The files are in `_delete/fin/`.
Lost: two IDs for one formula (CLAUDE.md: one function, one home).

## D-026 Exact linear money formulas (2026-11-12, finance lead; finance F-9)
EconomicValueAdded, FreeCashFlowToFirm and FreeCashFlowToEquity recover the product's rounding error with
math.FMA and round half-even on the exact sum. Each product is written `float64(float64(x) * r)`; the Go spec
says an explicit conversion rounds and blocks FMA fusion. FCFF uses ebit − ebit·t, never ebit·(1 − t).
Lost: plain float64 with the D-022 contract (cheap exactness given up).

## D-027 Bounds on periods (2026-11-12, finance lead; finance F-10)
- Loan schedules and loops: periods and paymentsMade ≤ 10⁷.
- Duration and convexity: ≤ 10⁷.
- Depreciation schedules: ≤ 10⁶.
Above the bound: ErrPeriods. Lost: no bound, which lets an input allocate gigabytes or loop for hours.

## D-028 Rounded level payments and early payoff (2026-11-12, finance lead; finance F-11)
AmortizationSchedule uses the half-even payment from FIN-INT-018. Each row's interest is rounded half-even.
The last payment absorbs the residue, with |last − level| ≤ (F/A,i,n) cents. When the rounded payment pays a
tiny loan off early, the schedule ends at the payoff row, and no balance goes negative. Lost: carrying the
unrounded payment, as OpenStax Table 8.4 does (its rows then disagree with any real statement by cents).

## D-029 Ratios are decimals (2026-11-12, finance lead; finance F-12)
Ratios and yields are returned as decimals (0.05), including the FFIEC UBPR ratios that the guide prints as
percents. The exceptions are the index values whose source defines them as numbers (CPI base 100, HHI
0–10 000). Lost: per-function percent outputs (callers would have to know which is which).

## D-030 Pure add/subtract rows are dropped unless Devin named them (2026-11-12, finance lead; finance F-13)
Net income, the accounting equation and cost of goods available for sale are dropped (one addition or
subtraction, D-008 by analogy). Net working capital and contribution margin are kept because Devin's
instruction names them. They are checked int64 arithmetic. Lost: a function per accounting identity.

## D-031 Angles between vectors by Kahan's θ = 2·atan2(‖â − b̂‖, ‖â + b̂‖) in 2D, 3D and n-D (2026-11-12, linear algebra lead; linalg table row 1)
Angles between vectors by Kahan's θ = 2·atan2(‖â − b̂‖, ‖â + b̂‖) in 2D, 3D and n-D. Lost: atan2(‖a × b‖, a·b) as the request suggested, because no primary source analysing it was found; Kahan 2006 §12 analyses his formula, shows arccos and arcsin failing, and his works in any dimension.

## D-032 Rotations are ACTIVE on column vectors (2026-11-12, linear algebra lead; linalg table row 2)
Rotations are ACTIVE on column vectors (v′ = Rv, v′ = qvq⁻¹). Lost: Diebel's frame (passive) convention, because callers turning a vector expect v′ = Rv; Diebel's matrices are cited as the transposes and checked so in tests.

## D-033 Quat is a struct {W, X, Y, Z} (2026-11-12, linear algebra lead; linalg table row 3)
Quat is a struct {W, X, Y, Z}. Lost: [4]float64, because array order (scalar first or last) is the classic source of quaternion bugs; named fields remove it.

## D-034 Every function returns an error, even Vec3Add (2026-11-12, linear algebra lead; linalg table row 4)
Every function returns an error, even Vec3Add. Lost: error-free fast arithmetic, because `.claude/rules/numeric.md`: NaN/Inf in must be an error, never a silent NaN; benchmarks show 2–10 ns, 0 allocations.

## D-035 Closed-form inverses (2026-11-12, linear algebra lead; linalg table row 5)
Closed-form inverses (2×2, 3×3, 4×4) scale each row by an exact power of two first. Lost: scale by the largest element only (first build), because a badly scaled, well-conditioned matrix (diag(1, 1e−300, 1e−300)) gave a determinant that underflowed and a false ErrSingular; found by a test.

## D-036 4×4 inverse by cofactors (2026-11-12, linear algebra lead; linalg table row 6)
4×4 inverse by cofactors (Laplace expansion with 2×2 minors). Lost: LU on the stack, because the request asked for closed forms for small sizes; cofactors cost no allocation; LU (MatInverse) is documented as the choice for ill-conditioned input.

## D-037 MatRank default tolerance max(rows, cols)·ε·max\|aᵢⱼ\| when tol < 0 (2026-11-12, linear algebra lead; linalg table row 7)
MatRank default tolerance max(rows, cols)·ε·max\|aᵢⱼ\| when tol < 0. Lost: no default (tol required), because the numeric policy asks for a documented default; elimination's backward error is of that order.

## D-038 QuatToEulerZYX takes lockTol (2026-11-12, linear algebra lead; linalg table row 8)
QuatToEulerZYX takes lockTol (default 2⁻⁵⁰ when < 0). Lost: a hidden constant; no lock handling, because near pitch ±π/2 the generic atan2 formulas split yaw and roll arbitrarily; at the lock the function returns roll 0 and the exact yaw, and the threshold is the caller's.

## D-039 QuatSlerp limits t to [0, 1] and takes the shorter arc (2026-11-12, linear algebra lead; linalg table row 9)
QuatSlerp limits t to [0, 1] and takes the shorter arc. Lost: extrapolation; the literal Shoemake formula without the sign flip, because extrapolated slerp is rarely meant; q and −q are the same rotation, and the long arc is almost never wanted.

## D-040 Slerp weights as s·sinc(sθ)/sinc(θ) (2026-11-12, linear algebra lead; linalg table row 10)
Slerp weights as s·sinc(sθ)/sinc(θ). Lost: a threshold switching to lerp for small θ, because no hidden tolerance is needed; the weights are exact (1 − t, t) at θ = 0.

## D-041 isFinite is \|x\| ≤ MaxFloat64; twoSum/twoProduct/det2 use explicit float64 conversions (2026-11-12, linear algebra lead; linalg table row 11)
isFinite is \|x\| ≤ MaxFloat64; twoSum/twoProduct/det2 use explicit float64 conversions. Lost: x − x == 0 and unconverted products, because Go may fuse an inlined product into a following add/sub (FMA): x − x became the product's rounding error and finite values were reported as overflow; the conversions keep the error-free transformations exact (Go spec: explicit conversion prevents fusion).

## D-042 Mat3SymmetricEigenvalues by Smith's closed form with its √ε limit documented (2026-11-12, linear algebra lead; linalg table row 12)
Mat3SymmetricEigenvalues by Smith's closed form with its √ε limit documented. Lost: an iterative Jacobi method now, because the request asked for closed forms; the limit for repeated eigenvalues (found by the oracle test) is stated on the function, and an iterative solver is LATER.

## D-043 Mat4Det/Mat4Inverse by the generalised Laplace expansion along two rows (2026-11-12, linear algebra lead; linalg table row 13)
Mat4Det/Mat4Inverse by the generalised Laplace expansion along two rows (six complementary 2×2 minor pairs), from Strang §5.2's cofactors. Lost: citing Eberly's "The Laplace Expansion Theorem" note (Geometric Tools, 2008), because the reviewer noted the s/c minor naming resembles that widely copied note; the code was derived from the expansion theorem, but the resemblance is recorded here; Eberly is not on the acceptable-source list, so it is not cited — Devin may rename or add it.

## D-044 Determinant sums (2026-11-12, linear algebra lead; linalg table row 14)
Determinant sums (Mat3Det, Mat4Det) by Dot2, with a row-rescaled retry when an intermediate overflows. Lost: plain FMA sums, because review found false ErrOverflow for finite determinants and x² − fl(x²) residues for exactly singular matrices.

## D-045 2D distance, polar coordinates, arc segments stay in num/geometry; angle normalisation… (2026-11-12, linear algebra lead; linalg table row 15)
2D distance, polar coordinates, arc segments stay in num/geometry; angle normalisation in num/trig. Lost: Vec2 versions here, because one function, one home (the request).

## D-046 Medicine packages and SHORTs (2026-11-12, lead under Devin's no-stop instruction; medicine D-M1)
Nine `med/<specialty>` packages: anthropometry ANTHRO, cardio CARDIO, respiratory RESP, renal RENAL, chemistry CHEM,
fluids FLUID, hepatic HEPAT, scales SCALE, pharmacology PHARM. Seven exist after this run (fluids and pharmacology
hold only HELD rows, so no folder was created: a package is created with its first function).
`med/chemistry` is new to the request's list: derived clinical-chemistry values (corrected calcium, osmolality,
eAG, LDL, anion gap) are one subject, not renal.
Lost: putting the lab-derived values in `med/renal` (corrected calcium, eAG and LDL are not kidney measures; one
function, one home by subject) and a `med/labs` grab-bag (CLAUDE.md forbids grab-bags).

## D-047 Citation bar for medicine: the formula must be SEEN in the primary (2026-11-12, medicine lead; medicine D-M2)
A function ships only when its formula, coefficients and labels were read in the primary source (abstract, full
text, scan or official document). A confirmed citation whose formula was seen only in a secondary source is HELD
(LEFT FOR DEVIN), never shipped.
Lost: shipping "well-known" formulas (Du Bois, Mosteller, MAP one-third rule, Friedewald, MDRD, Parkland) on
secondary pointers. That would break `.claude/rules/sources.md` exactly where it was written to hold.

## D-048 One function per unit system (2026-11-12, medicine lead; medicine D-M3)
Where clinicians use both systems, there are twin functions (`…SI`, `…Imperial`, `…MmolL`, `…GDL`). Each
converts once at entry with a named constant whose source is on the line above, and the tests assert that the
twins agree. Conversion constants come from the source itself where it gives one (ADAG 1.5944/2.594, the Martin
2013 caption factors 0.0259/0.0113, Dorwart's own /18 and /2.8). Otherwise they come from KDIGO 2024 p. S126
(creatinine 88.4, calcium 0.2495).
Lost: a unit parameter (`unit Unit`) on one function. It hides the unit from the name, and a wrong enum value is a
silent ×88 error.

## D-049 Categories are typed strings carrying the source's labels verbatim (2026-11-12, medicine lead; medicine D-M4)
`type XCategory string`, whose constant values are the source's own words (CDC "Healthy Weight", KDIGO "Normal or
high", Rimel "severe"). KDIGO uses the code as the value plus a `Term()` method.
Lost: integer enums with a separate label table. The label is the contract, so it is the value.

## D-050 Held back by policy: every output that is a drug dose or an amount to administer (2026-11-12, Devin's rule; medicine D-M5)
This covers the dosage family, pediatric acetaminophen, both dewormer rows and IV dose, plus fluid volumes to
administer: Parkland/ABA burn resuscitation and Holliday–Segar maintenance. The fluid volumes are amounts to
administer, so the same rule applies.
Lost: treating resuscitation volumes as "not a drug". The rule's purpose (Devin chooses the source before anything
tells a clinician how much to give) covers fluids.

## D-051 Corrected calcium is Payne 1973's form only (2026-11-12, medicine lead; medicine D-M6)
Adjusted Ca = Ca − albumin + 4.0 (mg/dL, g/dL), from the Payne abstract. The common 0.8 mg/dL-per-g/dL form is not
built.
Lost: the 0.8 form. No primary source was found; a pointer attributes it to a "modified Orrell" not seen.

## D-052 MELD: the UNOS/OPTN form as quoted in Kim 2021, plus MELDNa and MELD 3.0 (2026-11-12, medicine lead; medicine D-M7)
Lost: Kamath 2001's own coefficient form (0.957/0.378/1.120/0.643), which was not read in the paper.

## D-053 Epworth: score and published bands only, no item text (2026-11-12, medicine lead; medicine D-M8)
The ESS is © Murray W Johns and the official site says a licence is needed to use it. calcbench ships the
arithmetic (the sum of eight 0–3 ratings) and the five band labels; the licensed caller supplies the item wording.
Lost: embedding the item text (copyright), and leaving Epworth out entirely (the request names it).
Flag for Devin: confirm whether the consumer's use needs a licence.

## D-054 Moves (2026-11-12, medicine lead; medicine D-M9)
- Hagen–Poiseuille and Banana Equivalent Dose go to `sci/physics`.
- Diluting stock solutions (C₁V₁ = C₂V₂) goes to `sci/chemistry`.
- Direct Care FTEs goes to `fin/business`.
- Bayes for disease testing goes to `num/prob` (already NUM-PROB-010/-011).
- Starling stays in `med/cardio` (HELD).
- Moens–Korteweg stays in `med/cardio`, because arterial pulse-wave velocity is its use.
- Time of death and Weber fraction are LATER.
Lost: keeping Hagen–Poiseuille in cardio. It is general fluid physics, and cardio would call it.

## D-055 The medicine run's packages are additive-only in this run (2026-11-12, port lead; port D-P1)
The medicine agent finished: no change in med/ after 14:19, and its report was delivered. Its seven packages take new
files and new IDs (after the last), but none of its 44 files is edited, and the lead places the doc.go lines.
Lost: treating them as fully off-limits and pushing every cardio/renal/respiratory/hepatic port into a follow-up wave.
That would idle two of the four W1 builders, while the agent the rule protects is no longer writing.

## D-056 The source path lives only in the ledger; parity tests cite the ledger row (2026-11-12, port lead; port D-P2)
Parity tests carry `// parity: ledger <entry id>, source case N`, never a file path or a product name.
Lost: the source file:line in the test comment. The request also says nothing in calcbench may name the product,
and the ledger is the only place the path may appear; the ledger row has the file:line.

## D-057 Routing by home package, not by source folder (2026-11-12, port lead; port D-P3)
The source folders (general, cardiac…) mix specialties. Each entry is routed to the builder that owns its calcbench
home package, so no two builders share a package within a wave. The waves keep the request's order.
Lost: one builder per source folder writing into whatever package fits. Two builders would write into med/gastro
(Rome criteria live in general/ and metabolic/), which breaks "builders never touch another builder's package".

## D-058 New packages and SHORTs (2026-11-12, port lead; port D-P4)
med/gastro GASTRO · med/endocrine ENDO · med/psychiatry PSYCH · med/critical CRIT · med/hematology HEMAT ·
med/oncology ONCO · med/rheumatology RHEUM · med/infectious INFECT · med/pediatrics PEDS ·
med/thermoregulation THERMO · med/emergency EMERG · med/neuro NEURO · med/trauma TRAUMA · med/obstetric OB.
Also, from doc/taxonomy.md: unit/temperature TEMP and life/fitness FIT. life/nutrition NUTR is new.
Lost: a single med/scales holding every scored instrument. Packages are by subject (CLAUDE.md), and a scale is a
shape, not a subject.

## D-059 Citation bar carried over (D-047, medicine D-M2) and policy rows (2026-11-12, port lead; port D-P5)
A formula is ported only when it is SEEN in the primary. Dose and amount-to-administer rows, decision pathways and
administrative rules are NOT PORTED, with the reason recorded.
Lost: porting on the reference's citation. The request calls it a hypothesis, and the medicine run found several
wrong ones (PEF, Liebermeister, Harris–Benedict).

## D-060 VTE risk scores live in med/cardio (2026-11-12, port lead; port D-P6)
IMPEDE VTE (myeloma) and TRiP(cast) (lower-limb immobilisation) are venous-thromboembolism risk scores, and VTE
(Wells, Geneva, PERC, D-dimer) is homed in med/cardio. Lost: med/hematology for IMPEDE and med/trauma for TRiP(cast);
that would split the VTE family across three packages by the population each one studied, when the subject is the same.

## D-061 MED-CARDIO-009 withdrawn before release (2026-11-12, port lead; port D-P7)
The Canadian Syncope Risk Score point table obtained (sum −3…9) contradicts the paper's stated range (−3…11), so the
builder withdrew the function. The ID is not reused. Its INDEX row reads "MED-CARDIO-009 | REMOVED before release |
Canadian Syncope Risk Score: point table not confirmed".

## D-062 Series kernel: num/series is ours; num/stats and algo are reused, never written (2026-11-12, port lead; port D-P8)
The reference's series/ math goes to num/series, a new package. taxonomy.md plans it as "sequences and series,
trend and smoothing", and no subject has claimed it. num/stats belongs to the statistics subject (29 functions
built; another subject is still live in num/), and algo/* belongs to S1. A series function that duplicates one of
theirs (median, MAD, percentile, TopK, interval merge, rolling sums) is REUSED. A statistic they lack is a
PENDING_SHARED request to their owner, and is never written here.
Lost: writing the robust statistics into num/stats ourselves. That edits another subject's package while it may
be live.

## D-063 Fitness loads take one unit, any (`…KgOrLb`), no kg/lb twins (2026-11-12, port lead; port D-P9)
1RM, load-at-percent and volume-load equations are homogeneous in load: the result is in the unit of the input.
Lost: kg and lb twins. They would be identical code with two IDs; the FENa precedent (MED-RENAL-009) already takes
"one unit, any" for a ratio.

## D-064 LIF-FIT-004/005 (Mayhew 1992) withdrawn before release (2026-11-12, port lead; port D-P10)
The only text of the abstract seen was the reference's own copy of the publisher page, and the live page returns
402. The reference's copy is the hypothesis, not the primary (D-059). Files moved to _delete/life/fitness/. IDs
retired, not reused.
Lost: accepting the copied abstract. It would make the reference its own witness.

## D-065 The identity word the request excludes is allowed only inside a verbatim title in a Source: line (2026-11-12, port lead; port D-P11)
A citation's title is quoted exactly. Everywhere else (identifiers, prose, tests, doc.go) the identity rule holds.
Lost: rewording titles (a corrupted citation) or dropping titles (an incomplete citation, against sources.md).

## D-066 %DV rounding is half-even (2026-11-12, port lead; port D-P12)
21 CFR 101.9(d)(7)(ii) says "nearest whole percent" and states no tie rule, so .claude/rules/numeric.md's default
(half-even) applies. The one parity case (2.5 g of 20 g DV: source 13, ours 12) is a recorded disagreement for
Devin.
Lost: half-up. It matches the reference, but no cited rule supports it.

## D-067 A definition from the SI Brochure is a source for pace (2026-11-12, port lead; port D-P13)
Pace is the reciprocal of speed. The BIPM SI Brochure (9th ed., 2019) defines speed, so LIF-FIT-009 cites it,
as LIF-FIT-007 cites the definition of "% of 1RM".
Lost: NOT PORTED for want of a paper. A unit definition is the primary for a unit relation.

## D-068 The NEWS2 SpO₂ bands stay in two packages, guarded by a cross-check test (2026-11-12, port lead; port D-P14)
MED-CRIT-001 NEWS2 keeps its own SpO₂ switch, and MED-RESP-009/010 export the same bands. med/* packages do not
import each other in production code (medicine CONVENTIONS). An external test,
med/critical/news2_spo2_crosscheck_test.go, compares the two over every SpO₂ value 70–100 on both scales, on air and
on oxygen, so a drift fails the build.
Lost: having NEWS2 call med/respiratory, which would be the first production import between med packages.

## D-069 Wording-only edits to three medicine-run files (exception to D-055) (2026-11-12, port lead; port D-P15)
The request's identity rule applies to all of calcbench, so the excluded identity word describing study populations outside
quotation marks was changed to "persons" or the population's own words. Files: med/hepatic/fib4.go (2), the
Cockcroft–Gault file (1) and the Epworth score file (1). No code changed. The medicine agent had finished, so the
reason for D-055 does not apply to a comment fix.
Lost: leaving them. That breaks the identity rule, which the reviewer checks across the tree.

## D-070 INDEX.md is generated from the code (2026-11-12, lockdown pass)
Every INDEX row is built from the doc comment (Time, Space, Source lines joined by "; ", Tags) and the package
doc.go line (Does); a literal "|" in a cell is written "\|". REMOVED rows are kept by hand. The nine
PENDING_SHARED INDEX sections were checked against the generated rows: finance 124/124 identical, linalg 113/115,
medicine 42/44 and port 209/212 (the rest changed by later review fixes); the five catalog passes' rows (algebra,
combinatorics, geometry, prob, trig) predate the build and were superseded by the math build's rows. testing/index_test.go
keeps the two in step. Lost: pasting each pending section as written (stale cells, and 490 functions had no row).

## D-071 What carries an ID: functions, ID'd types and constants (2026-11-12, lockdown pass)
An ID sits on an exported function, method, type or constant. An exported method needs no ID of its own when its
receiver type carries an ID or is the result type of an ID'd function in the same package (accessors and labels,
D-017 L23); a constructor New<T> needs none when T carries an ID. Everything else exported needs an ID. Lost: one
ID per method (89 new IDs for Len, Push, Term and the like, against L23).

## D-072 doc/taxonomy.md lists one package per row with a Status (2026-11-12, lockdown pass)
Each package has its own row and a Status of built or planned, so a test can check the map against the tree. Rows
added for packages created without a taxonomy row: the seven medicine-run packages (D-046) and the thirteen port-run
med packages plus life/nutrition (D-058). Lost: the grouped rows (one row for eleven unit packages), which a test
cannot read.

## D-073 %DV rounds an exact half up (2026-11-12, lockdown pass; supersedes D-066)
FDA, Guidance for Industry: A Food Labeling Guide (January 2013), Appendix H, p. 129: "%DV values [that] fall
exactly halfway between two whole numbers or higher (e.g., 2.5 to 2.99) … round up". LIF-NUTR-006 now uses
math.Round (inputs are ≥ 0, so half away from zero is half up); the parity case 2.5 g of 20 g gives 13, as the
reference did. Name, parameters and result type are unchanged; only the tie behaviour moved to the agency's
written rule. Lost: half-even (D-066), chosen only because 21 CFR 101.9(d)(7)(ii) states no tie rule; the
lockdown request kept it "unless the FDA guidance says otherwise", and it does.

## D-074 Verbatim quotations keep the excluded words; prose never does (2026-11-12, lockdown pass; extends D-065)
The lockdown grep found the identity word the request excludes in 3 kinds of place. (1) Verbatim citation titles in
Source: lines, source notes and catalog "Primary source:" lines, and the instrument name "Patient Health
Questionnaire". (2) Sentences quoted from a source inside a doc comment: the Teasdale 2014 GCS quote, the MASCC
item wording, the Shapiro rule. (3) The source's own item name "outpatient status" (field MASCCFactors.Outpatient).
All three are kept, because changing a quotation corrupts the citation. Prose describing study populations was
rewritten ("adults", "persons"). The consumer's product and repository names were removed from every file outside
the S14 port work folder, and two consumer-specific files moved into it. The S14 work folder's own path keeps its
name wherever a link needs it. Lost: rewriting quotations (a false citation) or dropping them (an incomplete one).

## D-075 Braden and Mentzer withdrawn before release; fracture index homed in med/endocrine (2026-11-12, Devin)
MED-SCALE-010 BradenScore and MED-HEMAT-018 MentzerIndex were built in port wave 2 as PORTED-SECONDARY, but in each
case one of the two secondary sources was not seen citing the primary, so the two-source bar is not met. Both move to
NEEDS PDF (Bergstrom 1987 PMID 3299278; Mentzer 1973 PMID 4123424); their files are in _delete/ and the two IDs are
retired, never reused (REMOVED rows). The fracture index (Black 2001, PMID 11527048), which fits no package by name,
belongs in med/endocrine when its paper arrives. Lost: keeping the two functions on a weak second source, and
med/rheumatology for the fracture index.

## D-076 Secondary ports that fail the two-source bar are withdrawn before release (2026-11-12, port lead, wave 2; applies D-075)
The wave 2 reviewer and two verifiers re-checked every PORTED-SECONDARY function against the bar Devin applied in
D-075: two independent authoritative sources that both visibly cite the primary. Independent means not by each
other's authors and not by the primary's authors. Fourteen entries failed: bp_classify, hospital_score, Framingham
general CVD, DIPSS, SCORMA, Burton mean body temperature, O'Connor 1RM, AIR, PUCAI bands, DECAF, CASPAR, modified
Aldrete, RCVS2 and mNIHSS (out/verify_A.md, out/verify_B.md). Their 18 functions moved to _delete/ and the IDs are
retired, never reused (REMOVED rows): MED-CARDIO-061, 066, 067, 068 · MED-HEMAT-014, 015, 021 · MED-THERMO-003, 004 ·
LIF-FIT-013 · MED-GASTRO-012, 015 · MED-RESP-044 · MED-RHEUM-007, 008 · MED-SCALE-009 · MED-NEURO-029, 031. Every
entry is now NEEDS PDF. Eight entries passed only after a failing secondary was replaced by one that meets the bar,
and their Source and Verified lines now name it: CAM-ICU, Silverman–Andersen, BODE, GDS-15, CAS, Alvarado, PUCAI
(score), CAT. BODE's replacement cites Celli 2004 with wrong page numbers, so it is on CONFIRM. The 0–2 item check that
FLACC and PAINAD borrowed from the Aldrete file now lives in flacc_score.go. Lost: keeping the functions on a weak
second source; asking Devin about each (D-075 already decided the ground).

## D-077 The IOTA simple rules live in med/obstetric (2026-11-12, port lead, wave 2)
MED-OB-002 IOTASimpleRules classifies an adnexal mass, a gynaecology rule, and doc/taxonomy.md has no gynaecology
package. It sits in med/obstetric, the closest planned home for the reproductive organs, and the package comment
says so. Lost: med/oncology (the rule is used before any cancer diagnosis) and a new med/gynecology package (one
function; a new package is Devin's call, and a later move would be a new ID under the additive-API rule).

## D-078 Matrix gaps: package edges and scope (2026-11-12, matrix run, under "every decision is yours")
- num/linalg comes after num/prob in D-009's build order, so it may import num/algebra (CubicRoots for the 3×3
  characteristics' eigenvalues). Lost: a second cubic solver inside linalg (two homes for one formula).
- "is square" is not a function: rows == cols is one comparison (D-008 by analogy). Lost: a one-line wrapper.
- The covariance ⇄ correlation bridges live in num/linalg: they are the diagonal scalings D^(−1/2)·Σ·D^(−1/2) of a
  matrix, and num/stats has no matrix functions. Lost: num/stats (it works on samples, not matrices).
- General non-symmetric eigenvalues for n > 2 are LATER: they need Hessenberg reduction and the shifted QR algorithm.
  Power iteration covers the dominant eigenpair. Lost: building the QR algorithm now (a separate slice of its own).
- Golub & Van Loan 4th ed. is cited only by the sections whose titles were confirmed (§2.3, 3.1, 4.2, 5.2, 8.5, 8.6),
  with Strang 5th ed. (whole table of contents confirmed) as the second source. Lost: section numbers from memory.

## D-079 Fixed-size matrix functions keep their own small rank routine (2026-11-12, matrix run)
Mat2Characteristics and Mat3Characteristics compute the rank with rankFixed: MatRank's complete pivoting and
default tolerance, on a stack array, so the fixed-size path allocates nothing. MatRank itself is unchanged, because
released code is reused and never rewritten. Lost: calling MatRank, which allocates its n² work copy (reviewer
finding); refactoring MatRank, which would edit a released function.

## D-080 Matrix factorisations scale by powers of two (2026-11-12, matrix run; reviewer finding)
QRHouseholder, LeastSquaresQR and Mat2Eigen scale their input by an exact power of two, as SVDJacobi and
MatSymmetricEigenJacobi do. Each Householder vector is normalised by its largest entry, so vᵀv ≥ 1. Before this,
inputs near 1e±170 silently skipped a reflection or falsely overflowed. Mat2Eigen forms its discriminant as
((a − d)/2)² + bc rather than trace² − 4·det, which kept only √u accuracy for nearly equal eigenvalues. Lost: the
unscaled forms, which were wrong without an error at the extremes.

## D-081 Four new packages and the shared reduction API (2026-11-13, numpy build, DECIDED D1/D5)
algo/reduce (RED), num/signal (SIG), num/special (SPEC) and num/float (FLT) are added to doc/taxonomy.md.
algo/reduce's API was fixed by the lead before the build, so the other builders could code against it. It adds
two functions to D5's list: NeumaierAdd (one step, for running and streaming sums) and TwoSum. These are the
"summation helpers" D5 consolidates. AllFinite is variadic, so `AllFinite(a, b)` checks scalars and
`AllFinite(xs...)` checks a slice, with no allocation. Lost: Sum alone. The streaming sums in algo/window,
algo/stream and num/series accumulate one value at a time and could not have called a slice Sum, so their private
kernels would have stayed.

## D-082 Polynomial arithmetic is split by the import order (2026-11-13, numpy build)
num/linalg imports num/algebra (D-078), so num/algebra can never import num/linalg. Coefficient arithmetic
(add, sub, div, pow, derivative, integral, fromroots, companion matrix, trim, 2-D evaluation) goes in num/algebra.
Roots by companion-matrix eigenvalues (PolynomialRootsCompanion) and fitting (PolynomialFit,
PolynomialFitCovariance, Vandermonde) go in num/linalg, because they need the eigen-solver and LeastSquaresQR.
Lost: all of it in num/algebra as D1 says, which would create an import cycle.

## D-083 The multivariate normal takes the caller's Cholesky factor (2026-11-13, numpy build)
num/prob comes before num/linalg in D-009's order, so it may not import NUM-LINALG-152 Cholesky.
MultivariateNormalPDF and MultivariateNormalLogPDF take the lower-triangular factor L of Σ. The log-determinant is
2·Σ log Lᵢᵢ and the quadratic form is computed by forward substitution inside the density. Lost: a private
Cholesky in num/prob (two homes for one algorithm); putting the density in num/linalg (D9 names num/prob);
reordering D-009.

## D-084 What the consolidation replaces, and how replaced code is kept (2026-11-13, numpy build)
- **Scope of "NaN/Inf check":** a named private helper. Its predicate becomes reduce.AllFinite. A wrapper whose
  only job is to format the package's own sentinel error stays, with its body calling AllFinite. Inline
  `math.IsNaN || math.IsInf` expressions are not helpers and are left (LATER).
- **Kept private, with reasons:**
  - algo/tree's generic `checkFinite[T number]`: AllFinite takes float64 only.
  - algo/geometry2d binaryGCD: algo may not import num.
  - num/stats' two-pass, scaled kernels: their summation step now calls NeumaierAdd.
  - IRR, MIRR, profitability index and Macaulay duration keep their own Horner loops. They need the sign of an
    overflowed value, a second derivative, or an O(n) buffer, which NUM-ALGEB-001/002 cannot give with O(1) space.
    The numbers are unchanged.
- **How replaced code is kept:** a file that is edited rather than removed is first copied with `cp -n` to
  `_delete/<same path>`, with one line in `_delete/LEDGER.md`. Lost: `mv` (the file stays in use) and deleting
  the helper with no copy (the delete policy).

## D-085 API choices left open by DECIDED (2026-11-13, numpy build)
- **Errors and refusals:**
  - Interp (NUM-CALC-004) returns ErrOutOfRange outside the table; numpy clamps. This follows numeric.md:
    enforce the range.
  - MatLogDet returns ErrSingular; numpy returns (0, −Inf).
  - Sign(NaN) is ErrNotFinite.
  - FFTShift and IFFTShift are generic and return no error.
- **Argument and result shapes:**
  - LexSortIndices takes keys[0] as the primary key; numpy uses the last key.
  - BinCount takes the bin count from the caller, so its allocation is bounded.
  - CovarianceMatrix, CorrelationMatrix and FrequencyDistribution2D return row-major flat slices, as num/linalg does.
- **Dropped by the D1 list:**
  - hfft and ihfft.
  - The per-family coefficient arithmetic of the orthogonal polynomials (add, mul, der, int, fit, vander,
    companion, polyutils). D1 keeps only evaluation and roots.
- **Lost:** copying numpy's choices. They return NaN or clamp silently where calcbench's rules require an error.

## D-086 HFT and execution maths: placement and numeric limits (2026-11-13, S19, DECIDED D1–D5)
- **New package fin/market (MKT).** It imports no other fin package (the taxonomy's fin rule).
- **Domain helpers in algo/stream.** The tick and basis-point helpers, StreamingVWAP and WindowedVWAP sit in
  algo/stream because the request's owner list puts them there (B4). That conflicts with CLAUDE.md's "algo: no
  domain words", and the conflict is left for Devin (REPORT). Lost: moving them to fin/market without asking.
- **FIN-INVEST-031 is skipped.** fin/invest/doc.go already reserves it for ImpliedRepoRate (not built). IDs are never
  reused, so S19 starts at 032.
- **WindowedVariance accuracy.** It uses West's sliding update, O(1) per Add, plus a caller-invoked O(w) Resync.
  It does no periodic recompute: a hidden O(w) step would break "O(1) per Add". The downdate's rounding scales with
  the data that has passed through, not with the current window's variance. So "1e-9 relative at every step" holds
  for w ≥ 7 on the test series, and for w ≤ 3 the test uses a scale-relative bound (reviewer judged this a real
  numerical limit; Chan, Golub & LeVeque 1983). Lost: re-summing the window (D3 forbids it).
- **NEEDS SOURCE, not built:** DepthWeightedImbalance and TWAP (no primary source defines either down to the
  equation), and the continuous Kelly form (μ − r)/σ² (Kelly 1956 gives only the binary bet).
- **WeightedMid** is the imbalance-weighted mid only; Stoikov's full Markov-chain micro-price is not built.

## D-087 toolbench requests R-1…R-3: homes and shapes (2026-11-13, S20, DECIDED D1–D3)
- **algo/backoff (BACKOF), new.** Exponential, FullJitter, EqualJitter, DecorrelatedJitter on integer nanoseconds;
  the random draw u ∈ [0, 1) is a parameter, scaled exactly as ⌊d·⌊u·2⁵³⌋ / 2⁵³⌋ in 128 bits, so every function is
  pure. Lost: an RNG inside (impure, untestable), and float64 duration arithmetic (inexact above 2⁵³ ns).
- **num/integer CeilDiv (NUM-INTEG-013).** Truncated quotient plus one when the remainder has the divisor's sign; no
  a + b − 1 overflow. b = 0 is the package's existing ErrZeroDenominator (num/integer has no ErrDivByZero; a second
  sentinel for the same condition would split callers' errors.Is checks). Lost: a new ErrDivByZero.
- **algo/cron (CRON), new.** POSIX five-field syntax only (IEEE Std 1003.1-2024 crontab): no steps, names or
  macros, which POSIX does not define (no source, no ship). Next searches civil fields with carries over bitsets,
  never minute by minute; DST rules as the request states (gap → the transition instant, overlap → first occurrence
  only). Lost: a time/calendar package (not in the taxonomy) and the Vixie extensions (no primary standard).

## D-088 toolbench requests R-4, R-5: two methods on existing types (2026-11-13, S21, DECIDED D1–D2)
- **TokenBucket.Restore (ALG-RATELM-005).** Gives back n tokens the caller took earlier: refill to now, then add n
  tokens capped at the burst; n ≤ 0 and the zero value are no-ops. The caller counts the tokens to return (as
  x/time/rate's Reservation.CancelAt does); Allow is unchanged. Lost: a Reserve that may go into debt (a second,
  negative-level state that every existing method would have to learn, and a change to released semantics).
- **LRUCache.Oldest (ALG-LIST-018).** The entry Put would evict next, read without a use. Newest is not added: it
  would need its own ID, doc comment and tests, which is more than the one line D2 allowed. Lost: Newest beside it.

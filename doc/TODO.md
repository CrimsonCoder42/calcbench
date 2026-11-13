# TODO

Open work, newest first. Each item names its slice or ID.

- Merged at the 2026-11-12 lockdown from the finance, linear algebra, medicine and port PENDING_SHARED files
  (requests between subjects; satisfied ones are noted with their ID and are not open):
  - SATISFIED: med → num/prob Bayes for disease testing = NUM-PROB-008/010/011; linalg → num/geometry polar rows =
    NUM-GEOM-007/008; algebra → num/linalg vector addition = NUM-LINALG-001 Vec2Add (and 011 Vec3Add);
    port → algo/heap "TopK with a caller comparison (count with tie-break)" = ALG-HASH-003 TopKByCount.
  - num/stats: TukeyFences (Tukey 1977; hinges vs type-7 quartiles to decide); chi-square goodness of fit;
    frequency-weighted SD; weighted median (port requests).
  - algo/stream or num/stats: Welford with remove (push/pop) — ALG-STREAM-004 RunningStats has add and merge, no
    removal (port request, Welford 1962).
  - algo/prefix: float64 prefix sums with compensated summation — ALG-PREFIX-001 is integer-only (port request).
  - algo/window: co-occurring pairs within a window (port triage).
  - num/prob: likelihood-ratio odds chain posterior; chi-square quantile (port triage).
  - num/series: average real variability (Mena 2005, PMID 15716690; formula not in the abstract).
  - unit/mass, unit/length: kg⇄lb, cm⇄in (NIST SP 811) (port triage). When they exist, the conversion constants in
    med/* that are the sources' own (CDC 703, KDIGO 88.4/0.2495, Martin 0.0259/0.0113) stay; any other cm/in or
    kg/lb conversion goes through unit/* (medicine).
  - fin/personal (planned, D-021): salary conversions, retirement and pension rows (vCalc c14, 42, 77, 98, 112, 141).
  - fin/business: Direct Care FTEs (medicine move). sci/chemistry: dilution C₁V₁ = C₂V₂ (medicine move).
    sci/physics: Hagen–Poiseuille; Banana Equivalent Dose (no peer-reviewed origin found) (medicine move).
  - sci/physics: force and velocity vector rows and the 29 LM collections (linalg move; offer Vec2/Vec3 functions).
    sci/astronomy: ephemeris, orbit position and velocity, slant range (linalg move). sci/geoscience: LLA ⇄ ECEF on
    WGS 84 (NIMA TR8350.2 3rd ed.; EPSG:4979 → 4978), river flow compensation (linalg move).
  - medicine: release the HELD-SOURCE rows by reading the primaries in catalog/medicine.md §5; Devin to supply the
    sources for the dose family (§4) before any is built (D-050).
- Math build (2026-11-12), catalog text stale vs the code (catalogs are not edited by the build; owners to update):
  - All six catalogs (algebra, combinatorics, geometry, trig, stats, prob) still say "awaiting approval"/"DRAFT" and
    carry their Q1–Q7 as open; D-005…D-014 answered them.
  - trig.md: Kahan 2014 Table 2 is §9 (p. 8), not §8; ±Inf limits superseded (A-016); A&S 4.4.13 and Kahan 1987 dropped.
  - prob.md: NUM-PROB-041 λ = 1000 example is 0.0126146113487215 (not 0.012614611348719664); NUM-PROB-026 saddle-point
    sign is +δ(n) − δ(k) − δ(n−k) (Loader eq. 5); "machine precision for every n" holds near the mode only.
  - combinatorics.md and its PENDING_SHARED: NUM-COMBI-011 Derangements limit is n ≤ 20, not 21 (A-012); NUM-COMBI-009/010
    cite OEIS A111492/A144151, not Brualdi §2.2; STEP 1 rows 45/46 point at "63"/"64".
  - algebra.md and its PENDING_SHARED: Goldberg's cancellation-free formula is eq. (5), not (4) (NUM-ALGEB-006); Stewart,
    NR 3rd ed. and Knuth §4.6.4 eq. (2) are superseded by the built Source lines.
  - stats.md: NUM-STATS-026 names weigtmea.htm (the weighted TRIMMED mean) for the plain weighted mean; it is
    ch2/weigmean.pdf. NUM-STATS-030/031 are not in the catalog. D-010 names Hogg 7th ed.; the code cites the 8th (A-026).
  - geometry.md and its PENDING_SHARED: CRC 32nd ed. (†) and Bronshtein (‡) → CRC 31st ed.; NUM-GEOM-033 label
    "12 + 2√8" → 10 + 4√2; NUM-GEOM-046 threshold 2⁻²⁶ → 2⁻³⁰; NUM-GEOM-070 DLMF 19.21.10 → 19.21.11; NUM-GEOM-057
    "degree ≤ 3" → quadratic; NUM-GEOM-082 "no cancellation" is wrong; SCOPE row says 82 functions (80 built).
  - Every subject's PENDING_SHARED §1 INDEX rows are superseded by the rows in doc/math_build_2026-11-12/notes_*.md.
- NUM-ALGEB-009 CubicRoots: state the 2^±1000 coefficient-ratio underflow (A-014) in the doc comment.
- num/stats: deviations below ~1e-154 underflow when squared (ErrZeroVariance or 0 for tiny spreads, e.g.
  NUM-STATS-014); WeightedMean (NUM-STATS-026) loses underflowed w·x terms. A scale-up retry would fix both.
- algo/selection: an exported in-place select would save NUM-STATS-031 MAD one O(n) allocation (and enable the
  ALG-HEAP-001 hybrid below). New ID; Devin to decide.
- num/integer: Montgomery multiplication in place of bits.Rem64 in the internal mulMod (faster where 128/64 division
  is slow, arm64). Internal swap, no API change (NUM-INTEG-001…003).
- Elliptic integrals home (geometry PENDING_SHARED 3b): R_F, R_D, R_G and the AGM are unexported in num/geometry
  (NUM-GEOM-046, -070). When a second subject needs them: `num/calculus` or a new `num/special` (Devin).
- Regularized incomplete beta/gamma home (prob Q6: `num/special`, SHORT SPECL, or `num/calculus`). Lifts the 2²⁸-term
  CDF cap of A-029 (NUM-PROB discrete CDFs).
- Cauchy distribution: sourced (NIST e-Handbook §1.3.6.6.3) but LATER per D-013; no ID yet.
- Odds formats (decimal ↔ fractional ↔ American): held, no ID; no primary source found (only betting/industry sites).
- NUM-STATS-027 WeightedGeometricMean: not built, no confirmed source (Bullen, Handbook of Means, ch. II page
  unconfirmed). Out per D-010 unless a source turns up.
- NUM-STATS-007 GroupedMode: held. Needs one library check of Spiegel & Stephens, Schaum's Statistics 4th ed., ch. 3
  (L + Δ1/(Δ1+Δ2)·c); the code waits in `_delete/num/stats/grouped_mode.go`.
- Random sampling (prob PENDING_SHARED §5): a sampler, if ever wanted (DiceRoll, Monte Carlo), takes a caller's
  `rand.Source`, gets a new ID; no package-level generator.
- Basic Operations pass (D-008): percent change and ratio helpers in num/algebra, IDs from NUM-ALGEB-011.
- num/integer: Prime Number Theorem (MOVE from combinatorics row 21).
- num/series: golden ratio sequence (powers of φ, Fibonacci ratio convergence) (MOVE from geometry row 24; uses NUM-LOGS-001).
- num/stats: chi-square (prob request; the LATER inferential slice).
- unit/mass and unit/speed: mass and velocity unit conversions (MOVE from geometry rows 33, 44).
- sci/geoscience: Earth great-circle distance = cited Earth radius × NUM-TRIG-031 CentralAngle (D-012; trig request).
- sci/astronomy: sea-level horizon distance (calls NUM-GEOM-042) and horizon central angle arccos(R/(R+h))
  (geometry rows 18, 15); Drake and Seager equations (prob #24, #99).
- sci/physics: angular momentum in planetary motion (geometry row 4); grazing angle, torque at an angle, transverse
  force on a string, optics collections c2, c4 (trig #21, #36, #37); LM 33.x collections (prob #2–#5).
- eng/mechanical: conveyor capacity, mass capacity, volume output (algebra); beam-section centroid and second moment
  of area (geometry rows 16, 35). eng/electrical: impedance collection c1 (trig).
- life/games: Guess 1 to 31, Powerball (C(69,5)·26), five-card poker hands, dice, Magic (combinatorics and prob
  offers: NUM-COMBI-003, NUM-PROB-035…040, -001, -018). life/sports: NCAA bracket (NUM-PROB-012…017).
- life/home: quilting; spiral garland length (calls NUM-GEOM-081) (geometry rows 6, 31).
- Not planned, ask Devin (prob catalog): Friedman equation, Happiness probability, Probability of fatality, Sod's law.
- S1: Devin to confirm the page references for Knuth Vol. 3 §6.2.1 (the search for 653) and Bentley,
  Programming Pearls 2nd ed., Column 9.3, against a physical copy (see `sources/knuth_taocp_vol3.md`,
  `sources/bentley_programming_pearls.md`).
- ALG-HEAP-001: a TopK with the worst-case bound O(n + k log k) (heap first, batched selection once too many
  values enter the heap) needs an exported in-place select. Devin to decide: new ID(s) or leave it (A-011).
- ALG-WINDOW-001/-002: integer rolling sums with an overflow policy (spec, out of scope).

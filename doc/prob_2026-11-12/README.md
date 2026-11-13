# Work: probability catalog pass, 2026-11-12

Request: clean the vCalc "Probability" name list, catalog the standard set a probability package needs,
research every KEEP from primary sources, order the build, and write shared-file rows. DOCS ONLY.
Owned paths: `doc/catalog/prob.md` and this folder. `doc/catalog/stats.md` was read, not edited.

| File | What |
|---|---|
| [../catalog/prob.md](../catalog/prob.md) | The deliverable: Step 1 (139 rows), Step 2 standard set, Step 3 (62 entries), Step 4 order |
| [PENDING_SHARED.md](PENDING_SHARED.md) | 62 INDEX rows, SCOPE row, taxonomy option, requests to other subjects |
| [PROGRESS.md](PROGRESS.md) | Step log |

## Questions for Devin (recommendation · alternative · why it lost)

- **Q1 geometric convention.** TRIALS (X = trial of the first success, support 1, 2, …; mean 1/p), with the
  negative binomial as Pascal trials to match. Alternative: failures (support 0, …; R's dgeom). Lost: Ross
  and Casella–Berger and vCalc's "Pascal" rows all use trials; one convention across both distributions.
- **Q2 API.** Plain functions per distribution (`BinomialPMF(k, n, p)`), mean and variance as one pair
  `<Dist>MeanVariance`, parameters named for their parameterisation (exponential `rate`, lognormal `mu`,
  `sigma` of ln X). Alternative: a `Distribution` interface or struct methods (`Binomial{N,P}.PMF(k)`). Lost:
  IDs, INDEX rows and doc comments are per function in calcbench; methods still validate per call or need a
  fallible constructor; dynamic dispatch buys nothing a caller asked for.
- **Q3 odds home.** `num/prob`. Alternative: `life/sports`. Lost: life/games needs the same conversions and
  domain groups cannot import each other; the maths is probability.
- **Q4 Cauchy.** LATER: no consumer, no mean or variance, so it does not fit the set; closed forms make it
  cheap to add when sci/physics asks (Lorentzian line shape). Alternative: keep PDF/CDF/quantile now. Lost: scope.
- **Q5 "sums to 1" check** (TotalProbability, BayesPosterior, ExpectedValue). Reject when |Σp − 1| > n·u
  (u = 2⁻⁵³; each input carries at most one rounding). Alternative: a caller-supplied tolerance. Lost: every
  caller would invent a number. Second alternative, silent renormalising, lost: hides caller bugs.
- **Q6 special functions home** (regularized incomplete beta/gamma for O(1) binomial/Poisson CDFs now, and
  t/χ²/F later). Recommend a new `num/special` (SHORT SPECL). Alternative: `num/combinatorics` beside gamma.
  Lost: incomplete beta/gamma are analysis, not counting, and the inference slice needs them too.
- **Q7 American and fractional odds: SOURCE NOT FOUND.** Only betting sites define them. Held without ID or
  formula until you name an acceptable source (a regulator's rule book or a paper). Alternative: ship from
  the converter pages. Lost: violates `.claude/rules/sources.md`.

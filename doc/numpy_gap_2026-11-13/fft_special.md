# numpy gap: FFT, windows, random distributions (maths only), special functions (agent A4)

Pages read (WebFetch, 2026-11-13): routines.fft.html · routines.window.html · random/generator.html ·
routines.math.html ("Other special functions" section). Pages that failed: none.
Inventory: `calcbench_functions.md` (grepped by meaning), `num/prob/*.go` for unexported helpers, `doc/taxonomy.md`.

**Counts (58 rows): HAVE 7 · PARTIAL 6 · MISSING, IN SCOPE 40 · OUT OF SCOPE 5** (fftn/ifftn/rfftn/irfftn share one row)

Proposed new homes (not in `doc/taxonomy.md`, so a decision for Devin):
- `num/signal` NEW, SHORT `SIGNAL`: DFT/FFT, frequency bins, shifts, window functions. No FFT, DFT, window
  or `math/cmplx` spectral code exists anywhere in the tree (grep of the inventory and of `*.go`).
  Alternative: put windows in `num/series` (smoothing lives there); lost because windows are FFT tapers and
  belong beside the transform that uses them.
- `num/special` NEW, SHORT `SPECL`: i0, sinc, and the incomplete gamma/beta that several distributions need.
  Today `regularizedGamma` exists only unexported in `num/prob/chi_square_cdf.go:45` (DLMF 8.7.1 / 8.9.2).
  Alternative: keep them unexported in `num/prob`; lost because `num/signal` (Kaiser needs I0) would then
  have to import `num/prob` for a non-probability function.
- Distributions stay in `num/prob` (its existing home: PMF/PDF, CDF, Survival, Quantile, MeanVariance per family).

## numpy.fft

| numpy name | what it computes | calcbench match | verdict | proposed home + primary source |
|---|---|---|---|---|
| fft | 1-D complex DFT, O(n log n) | none | MISSING, IN SCOPE | num/signal; Cooley & Tukey 1965, Math. Comp. 19:297–301 (radix-2); Bluestein 1970, IEEE Trans. Audio Electroacoust. AU-18(4):451–455 (any n); Oppenheim & Schafer, Discrete-Time Signal Processing 3rd ed., §8.1, ch. 9 |
| ifft | inverse 1-D DFT | none | MISSING, IN SCOPE | num/signal; Oppenheim & Schafer 3rd ed. §8.2 (eq. 8.67–8.68) |
| rfft | DFT of a real input, n/2+1 bins | none | MISSING, IN SCOPE | num/signal; Press et al., Numerical Recipes 3rd ed. §12.3.1 (FFT of real functions) |
| irfft | inverse of rfft, real output | none | MISSING, IN SCOPE | num/signal; Numerical Recipes 3rd ed. §12.3.1 |
| hfft | FFT of a Hermitian-symmetric signal (real spectrum) | none | MISSING, IN SCOPE | num/signal; low priority, follows from irfft by conjugation; Oppenheim & Schafer 3rd ed. §8.6 (symmetry) |
| ihfft | inverse of hfft | none | MISSING, IN SCOPE | num/signal; low priority, as hfft |
| fftfreq | sample frequencies of the DFT bins, k/(n·d) | none | MISSING, IN SCOPE | num/signal; Oppenheim & Schafer 3rd ed. §8.1 (bin k ↔ ω = 2πk/N) |
| rfftfreq | non-negative bin frequencies for rfft | none | MISSING, IN SCOPE | num/signal; as fftfreq |
| fftshift | rotate zero frequency to the centre | none (algo has no slice rotate in inventory) | MISSING, IN SCOPE | num/signal; a cyclic shift by ⌊n/2⌋; source as fftfreq |
| ifftshift | inverse of fftshift | none | MISSING, IN SCOPE | num/signal; rotation by ⌈n/2⌉ |
| fft2 | 2-D DFT | none | OUT OF SCOPE | beyond the brief's FFT set (1-D, real-valued, radix-2/Bluestein); row-column composition of fft is the caller's |
| ifft2 | inverse 2-D DFT | none | OUT OF SCOPE | as fft2 |
| rfft2 | 2-D real DFT | none | OUT OF SCOPE | as fft2 |
| irfft2 | inverse 2-D real DFT | none | OUT OF SCOPE | as fft2 |
| fftn / ifftn / rfftn / irfftn | N-D DFTs | none | OUT OF SCOPE | ndarray machinery (axes over N dimensions) |

## Window functions

| numpy name | what it computes | calcbench match | verdict | proposed home + primary source |
|---|---|---|---|---|
| bartlett | triangular window, M points | none | MISSING, IN SCOPE | num/signal; Harris 1978, Proc. IEEE 66(1):51–83, §V.B, eq. (20) |
| blackman | 3-term Blackman window (0.42, 0.5, 0.08) | none | MISSING, IN SCOPE | num/signal; Harris 1978 §V.E, eq. (32) |
| hamming | raised cosine 0.54 − 0.46 cos | none | MISSING, IN SCOPE | num/signal; Harris 1978 §V.D, eq. (30a) |
| hanning | Hann window 0.5 − 0.5 cos | none | MISSING, IN SCOPE | num/signal; Harris 1978 §V.C, eq. (27a) |
| kaiser | Kaiser–Bessel window I0(β√(1−x²))/I0(β) | none | MISSING, IN SCOPE | num/signal (needs i0 below); Kaiser 1974, Proc. IEEE ISCAS, 20–23; Harris 1978 §V.H, eq. (46) |

Note: `HammingDistance` (ALG-BITS-013) and `IsHammingNumber` (ALG-INTMTH-009) are unrelated name collisions.

## Special functions (routines.math "Other special functions")

| numpy name | what it computes | calcbench match | verdict | proposed home + primary source |
|---|---|---|---|---|
| i0 | modified Bessel function of the first kind, order 0 | none; Go stdlib has math.J0/Y0 but no I0 | MISSING, IN SCOPE | num/special; DLMF §10.25.2 (series), §10.40.1 (large x); Abramowitz & Stegun 9.8.1–9.8.2 (polynomial approximations) |
| sinc | normalised sinc sin(πx)/(πx), 1 at x = 0 | none | MISSING, IN SCOPE | num/special; DLMF §4.14 for sin; definition: Oppenheim & Schafer 3rd ed. §2.9 / §4.3; Taylor branch near 0 to avoid 0/0 |

## numpy.random distributions (maths only: PMF/PDF, CDF, quantile; generators out of scope)

| numpy name | what it computes | calcbench match | verdict | proposed home + primary source |
|---|---|---|---|---|
| normal | normal distribution | NUM-PROB-047 NormalPDF, 048 NormalCDF, 049 NormalSurvival, 061 NormalQuantile | HAVE | |
| standard_normal | N(0,1) | NUM-PROB-047/048/061 with mean 0, sd 1 | HAVE | |
| lognormal | lognormal | NUM-PROB-051 PDF, 052 CDF, 053 MeanVariance, 062 Quantile | HAVE | |
| exponential | exponential (numpy uses scale β = 1/λ; calcbench rate λ) | NUM-PROB-054/055/056/060 | HAVE | |
| standard_exponential | Exp(1) | NUM-PROB-054/055/060 with rate 1 | HAVE | |
| uniform | continuous uniform | NUM-PROB-044/045/046/059 | HAVE | |
| geometric | geometric, trials convention (same as numpy) | NUM-PROB-029 PMF, 030 CDF, 031 MeanVariance, 057 Quantile | HAVE | |
| binomial | binomial | NUM-PROB-026 PMF, 027 CDF, 028 MeanVariance | PARTIAL: no quantile | num/prob BinomialQuantile; search on BinomialCDF; Johnson, Kemp & Kotz, Univariate Discrete Distributions 3rd ed., ch. 3 |
| poisson | Poisson | NUM-PROB-041 PMF, 042 CDF, 043 MeanVariance | PARTIAL: no quantile | num/prob PoissonQuantile; Johnson, Kemp & Kotz 3rd ed., ch. 4 |
| hypergeometric | hypergeometric | NUM-PROB-035 PMF, 036 CDF, 037 MeanVariance | PARTIAL: no quantile | num/prob; Johnson, Kemp & Kotz 3rd ed., ch. 6 |
| negative_binomial | numpy counts FAILURES before the n-th success | NUM-PROB-032/033/034 (TRIALS convention, X_trials = X_failures + r) | PARTIAL: failures convention only by a caller shift; no quantile | num/prob NegativeBinomialFailures* or a doc note; Johnson, Kemp & Kotz 3rd ed., ch. 5 |
| chisquare | chi-square | NUM-PROB-063 CDF, 064 Survival, 065 Quantile | PARTIAL: no PDF, no MeanVariance | num/prob ChiSquarePDF; DLMF §8.2 / Johnson, Kotz & Balakrishnan, Continuous Univariate Distributions vol. 1, 2nd ed., ch. 18 |
| multinomial | multinomial PMF n!/Πxᵢ! Πpᵢ^xᵢ | NUM-COMBI-008 Multinomial (coefficient only) | PARTIAL: coefficient exists, PMF missing | num/prob MultinomialPMF (log-space via math.Lgamma); Johnson, Kotz & Balakrishnan, Discrete Multivariate Distributions 1997, ch. 35 |
| gamma | gamma (shape k, scale θ) | none exported; `regularizedGamma` unexported at num/prob/chi_square_cdf.go:45 | MISSING, IN SCOPE | num/prob Gamma PDF/CDF/Survival/Quantile (CDF = P(k, x/θ)); DLMF §8.2.4, 8.7.1, 8.9.2 |
| standard_gamma | gamma with θ = 1 | as gamma | MISSING, IN SCOPE | covered by gamma with scale 1 |
| beta | beta distribution | none (no incomplete beta) | MISSING, IN SCOPE | num/special RegularizedIncompleteBeta (DLMF §8.17.22 continued fraction; Numerical Recipes 3rd ed. §6.4) + num/prob Beta PDF/CDF/Quantile |
| standard_t | Student's t | none | MISSING, IN SCOPE | num/prob; CDF via incomplete beta, DLMF §8.17 / A&S 26.7.1; Johnson, Kotz & Balakrishnan vol. 2, 2nd ed., ch. 28; Hill 1970 CACM Alg. 396 (quantile) |
| f | Snedecor F | none | MISSING, IN SCOPE | num/prob; CDF = I_{d1x/(d1x+d2)}(d1/2, d2/2), A&S 26.6.2; JKB vol. 2 ch. 27 |
| dirichlet | Dirichlet density on the simplex | none | MISSING, IN SCOPE | num/prob DirichletPDF (log-space via math.Lgamma); Kotz, Balakrishnan & Johnson, Continuous Multivariate Distributions vol. 1, 2nd ed., ch. 49 |
| gumbel | Gumbel (extreme value I) | none | MISSING, IN SCOPE | num/prob; closed-form PDF/CDF/quantile; JKB vol. 2, 2nd ed., ch. 22 |
| laplace | Laplace (double exponential) | none | MISSING, IN SCOPE | num/prob; closed form; JKB vol. 2, 2nd ed., ch. 24 |
| logistic | logistic | none | MISSING, IN SCOPE | num/prob; closed form; JKB vol. 2, 2nd ed., ch. 23 |
| logseries | logarithmic series PMF −pᵏ/(k ln(1−p)) | none | MISSING, IN SCOPE | num/prob; Johnson, Kemp & Kotz 3rd ed., ch. 7 |
| multivariate_hypergeometric | multivariate hypergeometric PMF | none | MISSING, IN SCOPE | num/prob; product of Binomial coefficients (NUM-COMBI-003, log form); Johnson, Kotz & Balakrishnan 1997, ch. 39 |
| multivariate_normal | MVN density | none (NUM-LINALG-152 Cholesky exists for the factor) | MISSING, IN SCOPE | num/prob MultivariateNormalPDF via Cholesky + log-det; Kotz, Balakrishnan & Johnson vol. 1, ch. 45; Golub & Van Loan 4th ed. §4.2 |
| noncentral_chisquare | noncentral χ² | none | MISSING, IN SCOPE | num/prob; Poisson mixture of central χ², Ding 1992, Appl. Stat. 41(2):478–482 (Alg. AS 275); low priority |
| noncentral_f | noncentral F | none | MISSING, IN SCOPE | num/prob; Poisson mixture of incomplete betas, Lenth 1987, Appl. Stat. 36(2):241–244 (AS 226); low priority |
| pareto | Lomax / Pareto II (numpy's form: shifted by 1) | none | MISSING, IN SCOPE | num/prob; closed form; JKB vol. 1, 2nd ed., ch. 20 |
| power | power-function distribution x^(a−1)·a on [0,1] | none | MISSING, IN SCOPE | num/prob; closed form; JKB vol. 2, 2nd ed., ch. 25 (beta(a,1) special case) |
| rayleigh | Rayleigh | none | MISSING, IN SCOPE | num/prob; closed form; JKB vol. 1, 2nd ed., ch. 18 §10 |
| standard_cauchy | Cauchy (0,1) | none | MISSING, IN SCOPE | num/prob; closed form incl. quantile tan(π(u−½)); JKB vol. 1, 2nd ed., ch. 16 |
| triangular | triangular on [a, c, b] | none | MISSING, IN SCOPE | num/prob; closed form; Johnson, Kotz & Balakrishnan vol. 2 ch. 26 §9 |
| vonmises | von Mises (circular) | none | MISSING, IN SCOPE | num/prob PDF needs I0 (num/special); CDF series Hill 1977, ACM TOMS 3(3):279–284 (Alg. 518) |
| wald | inverse Gaussian | none | MISSING, IN SCOPE | num/prob; CDF in terms of Φ (NUM-PROB-048); Chhikara & Folks, The Inverse Gaussian Distribution 1989, ch. 2 |
| weibull | Weibull (numpy: scale 1) | none | MISSING, IN SCOPE | num/prob; closed form; JKB vol. 1, 2nd ed., ch. 21 |
| zipf | zeta distribution k^(−a)/ζ(a), infinite support | none (no Riemann zeta) | MISSING, IN SCOPE | num/special RiemannZeta (DLMF §25.2, Euler–Maclaurin §25.2.9 or Borwein 2000) + num/prob ZipfPMF |

## Out of scope, listed once

Random bit generators and the sampling methods themselves (Generator.integers, random, choice, bytes, shuffle,
permutation, permuted, and the draw side of every distribution above). The nearest calcbench maths with a
caller-supplied source is already in algo/sample (ALG-SAMPLE-001 ShuffleInPlace, 002 ReservoirSample, 004
WeightedPicker.Pick, 005 RejectionSample). N-D FFTs (fftn family) and 2-D FFTs (fft2 family): ndarray/axes
machinery beyond the brief's 1-D FFT set. Complex dtype handling, `norm=` modes and `out=` arrays are API
machinery; a calcbench fft would fix one normalisation (unscaled forward, 1/n inverse) in its doc comment.

## Hand loops in calcbench

- `fin/invest/black_scholes_call.go:9` defines its own `normalCDF` (0.5·Erfc(−x/√2)) instead of importing
  NUM-PROB-048 NormalCDF (fin may import num). Same pattern likely in black_scholes_put.go. Not a numpy gap,
  but a duplicated distribution; FOUND, NOT FIXED.
- `num/prob/chi_square_cdf.go:31,45`: `gammaPrefactor` and `regularizedGamma` are the whole of the gamma CDF,
  unexported. Exporting them (num/special) unlocks gamma, standard_gamma and the Poisson/χ² quantile paths
  with almost no new maths.
- No FFT-shaped hand loop (direct O(n²) DFT) was found in sci/ or eng/ (no `math/cmplx` use outside
  num/linalg/eigen.go and num/algebra roots).

# numpy gap: mathematical functions, logic, isclose family, floating-point helpers (agent A1)

Pages read (WebFetch, 2026-11-13, NumPy v2.5 reference):
- https://numpy.org/doc/stable/reference/routines.math.html (read, complete)
- https://numpy.org/doc/stable/reference/routines.logic.html (read, complete)
- isclose/allclose: on routines.logic (Comparison). Floating helpers: routines.math "Floating point routines" plus
  nan_to_num (Miscellaneous); finfo/iinfo added as one row each from the brief.
Pages that failed: none.

Matching: by meaning against `calcbench_functions.md` (grep only; no package opened). "stdlib scalar" means Go's
`math` (or a builtin / `slices` / `math/cmplx`) already does the scalar; per D-012 calcbench never wraps it, and the
slice form is the vec-layer question. Proposed IDs are not assigned here.

Counts: HAVE 68 · PARTIAL 9 · MISSING IN SCOPE 31 · OUT OF SCOPE 16 (124 rows; some rows cover aliases, e.g. arcsin / asin)

| numpy name | what it computes | calcbench match | verdict | proposed home + primary source |
|---|---|---|---|---|
| sin | sine | stdlib: math.Sin (stdlib scalar) | HAVE | |
| cos | cosine | stdlib: math.Cos (stdlib scalar) | HAVE | |
| tan | tangent | stdlib: math.Tan (stdlib scalar) | HAVE | |
| arcsin / asin | inverse sine | stdlib: math.Asin (stdlib scalar) | HAVE | |
| arccos / acos | inverse cosine | stdlib: math.Acos (stdlib scalar) | HAVE | |
| arctan / atan | inverse tangent | stdlib: math.Atan (stdlib scalar) | HAVE | |
| hypot | √(x²+y²) without overflow | stdlib: math.Hypot (stdlib scalar); n-D: NUM-LINALG-041 VecNorm | HAVE | |
| arctan2 / atan2 | quadrant-correct atan(y/x) | stdlib: math.Atan2 (stdlib scalar) | HAVE | |
| degrees | radians to degrees | NUM-TRIG-002 RadToDeg | HAVE | |
| radians | degrees to radians | NUM-TRIG-001 DegToRad | HAVE | |
| deg2rad | degrees to radians | NUM-TRIG-001 DegToRad | HAVE | |
| rad2deg | radians to degrees | NUM-TRIG-002 RadToDeg | HAVE | |
| unwrap | remove 2π jumps from a phase sequence (cumulative wrapped differences) | NUM-TRIG-005 AngleDifference (the per-step kernel), NUM-TRIG-004 | PARTIAL: no sequence form, no custom period | num/trig `UnwrapAngles`; Itoh, "Analysis of the phase unwrapping algorithm", Applied Optics 21(14):2470, 1982 |
| sinh | hyperbolic sine | stdlib: math.Sinh (stdlib scalar) | HAVE | |
| cosh | hyperbolic cosine | stdlib: math.Cosh (stdlib scalar) | HAVE | |
| tanh | hyperbolic tangent | stdlib: math.Tanh (stdlib scalar) | HAVE | |
| arcsinh / asinh | inverse hyperbolic sine | stdlib: math.Asinh (stdlib scalar) | HAVE | |
| arccosh / acosh | inverse hyperbolic cosine | stdlib: math.Acosh (stdlib scalar) | HAVE | |
| arctanh / atanh | inverse hyperbolic tangent | stdlib: math.Atanh (stdlib scalar) | HAVE | |
| round | round to d decimals, ties to even | stdlib: math.RoundToEven (d = 0 only) | PARTIAL: no decimals argument; a correctly rounded decimal round is real numerics (x·10^d is inexact) | num/algebra `RoundDecimals`; IEEE 754-2019 §5.3.1 (roundToIntegralTiesToEven) + Muller et al., Handbook of Floating-Point Arithmetic 2nd ed. (2018) §4.9 (conversions) |
| around | alias of round | same as round | PARTIAL: as round | as round |
| rint | round to nearest integer, ties to even | stdlib: math.RoundToEven (stdlib scalar) | HAVE | |
| fix | round toward zero | stdlib: math.Trunc (stdlib scalar) | HAVE | |
| floor | floor | stdlib: math.Floor (stdlib scalar) | HAVE | |
| ceil | ceiling | stdlib: math.Ceil (stdlib scalar) | HAVE | |
| trunc | truncate | stdlib: math.Trunc (stdlib scalar) | HAVE | |
| prod | product of the elements | none (ALG-PREFIX-010 ProductExceptSelf is a different routine) | MISSING, IN SCOPE | num/linalg `VecProduct` with overflow/underflow scaling; Higham, Accuracy and Stability of Numerical Algorithms 2nd ed. (2002) §3.1 (products) |
| sum | sum of the elements | ALG-PREFIX-001 PrefixSums (last element), ALG-PREFIX-002 RangeSum; compensated sums only inside VecNorm1/MatTrace/RollingSum | PARTIAL: no standalone (compensated) Sum of a slice | num/linalg `VecSum` (Neumaier compensated); Neumaier, ZAMM 54:39–51, 1974; Higham 2002 §4.3 |
| nanprod | product treating NaN as 1 | none | MISSING, IN SCOPE (policy conflict: numeric.md says NaN in → error; needs a D-id) | num/linalg `VecProductSkipNaN`; Higham 2002 §3.1 |
| nansum | sum treating NaN as 0 | none | MISSING, IN SCOPE (same NaN-policy D-id) | num/linalg `VecSumSkipNaN`; Neumaier 1974 |
| cumulative_sum | running sum, optional leading 0 | ALG-PREFIX-001 PrefixSums (n+1 sums with the leading 0 = include_initial=True) | HAVE | |
| cumsum | running sum, n values | ALG-PREFIX-001 PrefixSums (drop element 0) | HAVE (offset by one; uncompensated per Does text) | |
| cumulative_prod | running product | none | MISSING, IN SCOPE | algo/prefix `PrefixProducts`; Knuth TAOCP vol 1 §1.2.3 (products) |
| cumprod | running product | none | MISSING, IN SCOPE | as cumulative_prod |
| nancumprod | running product, NaN as 1 | none | MISSING, IN SCOPE (NaN-policy D-id) | as cumulative_prod |
| nancumsum | running sum, NaN as 0 | none | MISSING, IN SCOPE (NaN-policy D-id) | algo/prefix; as cumsum |
| diff | n-th forward difference a[i+1]−a[i] | none exported (ALG-PREFIX-005 DifferenceArray holds differences but exposes no slice; NUM-SERIES-009 uses them inside) | MISSING, IN SCOPE | algo/prefix `Differences(values, order)`; Abramowitz & Stegun (1964) §25.1 (forward differences Δⁿ) |
| ediff1d | first differences of a flattened array, optional ends | none | MISSING, IN SCOPE (covered by `Differences` order 1) | as diff |
| gradient | second-order central differences, one-sided at ends, non-uniform spacing | match?: NUM-ALGEB-002 (polynomial derivative only) | MISSING, IN SCOPE | num/calculus (planned) `Gradient`; Fornberg, "Generation of finite difference formulas on arbitrarily spaced grids", Math. Comp. 51(184):699–706, 1988 |
| cross | cross product of 3-vectors (and 2-D scalar) | NUM-LINALG-016 Vec3Cross, NUM-LINALG-006 Vec2Cross | HAVE (no batch over rows) | |
| trapezoid | composite trapezoidal integral of samples, with x or dx | none (NUM-GEOM-025 is the shape's area) | MISSING, IN SCOPE | num/calculus (planned) `TrapezoidIntegral`; Abramowitz & Stegun §25.4.2; Press et al., Numerical Recipes 3rd ed. §4.1 |
| exp | eˣ | stdlib: math.Exp (stdlib scalar) | HAVE | |
| expm1 | eˣ − 1 accurately | stdlib: math.Expm1 (stdlib scalar) | HAVE | |
| exp2 | 2ˣ | stdlib: math.Exp2 (stdlib scalar) | HAVE | |
| log | ln x | stdlib: math.Log (stdlib scalar) | HAVE | |
| log10 | log₁₀ x | stdlib: math.Log10 (stdlib scalar) | HAVE | |
| log2 | log₂ x | stdlib: math.Log2 (stdlib scalar) | HAVE | |
| log1p | ln(1+x) accurately | stdlib: math.Log1p (stdlib scalar) | HAVE | |
| logaddexp | ln(eˣ + eʸ) without overflow | none | MISSING, IN SCOPE | num/logs `LogAddExp` (+ slice `LogSumExp`); Blanchard, Higham & Higham, "Accurately computing the log-sum-exp and softmax functions", IMA J. Numer. Anal. 41(4):2311–2330, 2021, doi:10.1093/imanum/draa038 |
| logaddexp2 | log₂(2ˣ + 2ʸ) | none | MISSING, IN SCOPE | num/logs `LogAddExp2`; same source |
| i0 | modified Bessel I₀ | none | MISSING, IN SCOPE (A4 owns the row in fft_special.md) | NEW num/special (SPECL) or per A4; DLMF §10.25, A&S §9.8.1–9.8.2 |
| sinc | sin(πx)/(πx) | none | MISSING, IN SCOPE (A4 owns the row) | per A4; Oppenheim & Schafer, Discrete-Time Signal Processing 3rd ed. §4.3 |
| signbit | whether the sign bit is set | stdlib: math.Signbit (stdlib scalar) | HAVE | |
| copysign | magnitude of x, sign of y | stdlib: math.Copysign (stdlib scalar) | HAVE | |
| frexp | mantissa and binary exponent | stdlib: math.Frexp (stdlib scalar) | HAVE | |
| ldexp | x·2ⁿ | stdlib: math.Ldexp (stdlib scalar) | HAVE | |
| nextafter | next float toward y | stdlib: math.Nextafter (stdlib scalar) | HAVE | |
| spacing | ulp(x): gap to the next float away from zero | none (composable from math.Nextafter, but sign/Inf/NaN/MaxFloat edge rules are real) | MISSING, IN SCOPE | NEW num/float (FLOAT) `Ulp`; Muller et al., Handbook of Floating-Point Arithmetic 2nd ed. (2018) §2.6.1 (ulp definitions) |
| lcm | least common multiple of \|a\|, \|b\| | none (math/big has only GCD for big.Int) | MISSING, IN SCOPE | num/integer (taxonomy already names gcd/lcm there, D-005) `LCM`, overflow checked; Knuth TAOCP vol 2 §4.5.2 |
| gcd | greatest common divisor of \|a\|, \|b\| | none in calcbench; stdlib only big.Int.GCD | MISSING, IN SCOPE | num/integer `GCD` (binary or Euclid); Knuth TAOCP vol 2 §4.5.2 Algorithms A and B |
| add | elementwise a + b | NUM-LINALG-035 VecAdd (slice form); scalar: operator | HAVE | |
| subtract | elementwise a − b | NUM-LINALG-036 VecSub | HAVE | |
| multiply | elementwise a·b | NUM-LINALG-037 VecScale (scalar × vector only) | PARTIAL: no elementwise (Hadamard) product of two slices | num/linalg `VecHadamard`; Horn & Johnson, Topics in Matrix Analysis (1991) §5.0 |
| divide / true_divide | elementwise a / b | NUM-LINALG-038 VecDiv (vector / scalar only) | PARTIAL: no elementwise slice ÷ slice with zero checks | num/linalg `VecDivElementwise`; operator definition, no special source |
| reciprocal | 1/x | operator (stdlib scalar) | HAVE | |
| positive | +x | operator | HAVE | |
| negative | −x | operator; slice: NUM-LINALG-037 VecScale(−1) | HAVE | |
| power / pow | xʸ | stdlib: math.Pow (stdlib scalar) | HAVE | |
| float_power | xʸ in float64 | stdlib: math.Pow; integer: ALG-INTMTH-001 PowerBySquaring | HAVE | |
| floor_divide | ⌊a/b⌋ | stdlib composition math.Floor(a/b) (inexact near integers; numpy uses (a − fmod)/b) | MISSING, IN SCOPE (with mod, as one floor-division pair) | num/algebra `FloorDivMod`; Knuth TAOCP vol 1 §1.2.4 eq. (1) x mod y = x − y⌊x/y⌋ |
| fmod | remainder, sign of the dividend (C fmod) | stdlib: math.Mod (stdlib scalar) | HAVE | |
| mod | remainder, sign of the divisor (floored) | none (math.Mod is truncated, math.Remainder is IEEE nearest; NUM-TRIG-003 NormalizeAngle is the 2π special case) | MISSING, IN SCOPE | num/algebra `FloorMod` (float) and num/integer (int64); Knuth TAOCP vol 1 §1.2.4 |
| modf | fractional and integral parts | stdlib: math.Modf (stdlib scalar) | HAVE | |
| remainder | alias of mod (floored) | none | MISSING, IN SCOPE (same function as mod) | as mod |
| divmod | floored quotient and remainder together | ALG-INTMTH-003 LongDivision (unsigned integer, truncated) | MISSING, IN SCOPE (float floored pair) | num/algebra `FloorDivMod`; Knuth TAOCP vol 1 §1.2.4 |
| angle | argument of a complex number | stdlib: cmplx.Phase (stdlib scalar) | HAVE | |
| real | real part | builtin real() | HAVE | |
| imag | imaginary part | builtin imag() | HAVE | |
| conj / conjugate | complex conjugate | stdlib: cmplx.Conj | HAVE | |
| maximum | elementwise max, NaN propagates | NUM-LINALG-051 VecMax, NUM-LINALG-034 Vec3Max; scalar: builtin max | HAVE | |
| max / amax | largest element | stdlib: slices.Max (panics on empty; propagates NaN) | HAVE (stdlib; panics on empty input, so a calcbench caller needs a length check) | |
| fmax | elementwise max, NaN ignored | none (math.Max and builtin max propagate NaN) | MISSING, IN SCOPE (NaN-policy D-id) | NEW num/float `FMax`; IEEE 754-2008 §5.3.1 maxNum |
| nanmax | largest element ignoring NaN | none | MISSING, IN SCOPE (NaN-policy D-id) | num/linalg `VecMaxSkipNaN`; IEEE 754-2008 §5.3.1 maxNum |
| minimum | elementwise min, NaN propagates | NUM-LINALG-050 VecMin, NUM-LINALG-033 Vec3Min; scalar: builtin min | HAVE | |
| min / amin | smallest element | stdlib: slices.Min (panics on empty; propagates NaN) | HAVE (as max) | |
| fmin | elementwise min, NaN ignored | none | MISSING, IN SCOPE (NaN-policy D-id) | NEW num/float `FMin`; IEEE 754-2008 §5.3.1 minNum |
| nanmin | smallest element ignoring NaN | none | MISSING, IN SCOPE (NaN-policy D-id) | num/linalg `VecMinSkipNaN`; IEEE 754-2008 §5.3.1 |
| convolve | discrete linear convolution of two 1-D sequences (full/same/valid) | NUM-ALGEB-003 PolynomialMultiply (identical maths to mode='full') | PARTIAL: no 'same'/'valid' modes, no FFT path for long inputs | num/series or num/algebra `Convolve(a, v, mode)`; Oppenheim & Schafer DTSP 3rd ed. §2.3 |
| clip | limit values to [lo, hi] | builtins min/max (Go 1.21): min(max(x, lo), hi) (stdlib scalar) | HAVE (scalar; slice form is vec layer; no lo > hi check) | |
| sqrt | square root | stdlib: math.Sqrt (stdlib scalar) | HAVE | |
| cbrt | cube root | stdlib: math.Cbrt (stdlib scalar) | HAVE | |
| square | x² | operator x*x | HAVE | |
| absolute | \|x\| | stdlib: math.Abs (stdlib scalar) | HAVE | |
| fabs | \|x\| for floats | stdlib: math.Abs | HAVE | |
| sign | −1, 0, +1 (NaN for NaN) | none (Go has math.Signbit only) | MISSING, IN SCOPE (small) | num/algebra `Sign`; DLMF §1.16(iv) sign function / A&S §3.1 |
| heaviside | step H(x) with chosen H(0) | none | MISSING, IN SCOPE (small) | num/algebra `Heaviside(x, atZero)`; DLMF §1.16.13 |
| nan_to_num | replace NaN/±Inf with finite values | none | MISSING, IN SCOPE (NaN-policy D-id; arguably a data-cleaning helper) | NEW num/float `ReplaceNonFinite`; IEEE 754-2019 §3.4 (special values); no formula source needed beyond the caller's replacements |
| real_if_close | drop imaginary parts within tolerance | none | OUT OF SCOPE (dtype conversion of complex arrays) | |
| interp | piecewise-linear interpolation through sorted (xp, fp), with end values | NUM-LINALG-049 VecLerp, NUM-LINALG-032 Vec3Lerp (two-point lerp only) | PARTIAL: no table interpolation (search + lerp), no left/right/period | NEW num/interp (INTERP) or num/calculus `LinearInterpolate`; Press et al., Numerical Recipes 3rd ed. §3.1–3.3 |
| all | whether every element is true | none; Go: loop / !slices.ContainsFunc | OUT OF SCOPE (boolean-array machinery) | |
| any | whether some element is true | Go: slices.ContainsFunc | OUT OF SCOPE (boolean-array machinery) | |
| isfinite | not Inf and not NaN | stdlib: !math.IsInf(x,0) && !math.IsNaN(x) (stdlib scalar) | HAVE | |
| isinf | ±Inf | stdlib: math.IsInf(x, 0) | HAVE | |
| isnan | NaN | stdlib: math.IsNaN | HAVE | |
| isnat | Not-a-Time | none | OUT OF SCOPE (datetime64) | |
| isneginf | −Inf | stdlib: math.IsInf(x, −1) | HAVE | |
| isposinf | +Inf | stdlib: math.IsInf(x, 1) | HAVE | |
| iscomplex | elementwise has nonzero imaginary part | none | OUT OF SCOPE (dtype testing) | |
| iscomplexobj | complex dtype check | none | OUT OF SCOPE (dtype) | |
| isfortran | Fortran-contiguous layout | none | OUT OF SCOPE (strides) | |
| isreal | elementwise zero imaginary part | none | OUT OF SCOPE (dtype testing) | |
| isrealobj | non-complex dtype check | none | OUT OF SCOPE (dtype) | |
| isscalar | scalar type check | none | OUT OF SCOPE (dtype; Go is statically typed) | |
| logical_and | elementwise AND | Go operator && | OUT OF SCOPE (ufunc over bool arrays) | |
| logical_or | elementwise OR | Go operator \|\| | OUT OF SCOPE (ufunc) | |
| logical_not | elementwise NOT | Go operator ! | OUT OF SCOPE (ufunc) | |
| logical_xor | elementwise XOR | Go a != b | OUT OF SCOPE (ufunc) | |
| allclose | every pair within atol + rtol·\|b\| | none | MISSING, IN SCOPE | NEW num/float (FLOAT) `AllClose(a, b []float64, relTol, absTol)`; Knuth TAOCP vol 2 3rd ed. §4.2.2 eqs. (21)–(24) (approximately/essentially equal); note numpy's form is asymmetric in b, Knuth's is symmetric: a D-id picks one |
| isclose | \|a − b\| ≤ atol + rtol·\|b\| per element | none (tolerances appear only inside MatRank/MatRowEchelon) | MISSING, IN SCOPE | NEW num/float `IsClose(a, b, relTol, absTol float64)`; Knuth TAOCP vol 2 §4.2.2 |
| array_equal | same length and elements | stdlib: slices.Equal (NaN ≠ NaN, as numpy's default) | HAVE | |
| array_equiv | equal after broadcasting | none | OUT OF SCOPE (broadcasting) | |
| greater / greater_equal / less / less_equal / equal / not_equal | elementwise comparisons | Go operators | OUT OF SCOPE (ufunc over arrays; listed as one row) | |
| finfo (eps, tiny, max, resolution) | float64 machine constants | stdlib: math.MaxFloat64, math.SmallestNonzeroFloat64; no Epsilon constant (0x1p-52 literal) | PARTIAL: no named machine epsilon / smallest normal | NEW num/float constants `MachineEpsilon`, `SmallestNormal`; IEEE 754-2019 §3.3 Table 3.5 (binary64 parameters) |
| iinfo | integer type limits | stdlib: math.MaxInt64, math.MinInt64, etc. | HAVE | |

Out of scope, listed once: ndarray, dtypes and dtype tests (iscomplexobj, isrealobj, isscalar, isfortran,
real_if_close), broadcasting (array_equiv), ufunc machinery (comparison and logical_* ufuncs, all/any over bool
arrays, out=/where=/casting= arguments of every ufunc above), datetime64 (isnat), masked arrays.

Decision for Devin raised by this page (not made here): seven MISSING rows (nansum, nanprod, nancumsum,
nancumprod, nanmax/nanmin, fmax/fmin, nan_to_num) contradict `.claude/rules/numeric.md` ("NaN or Inf in means an
error out"). Recommendation: one D-id allowing explicitly named `...SkipNaN` variants in a NEW `num/float` package;
the alternative (never ship NaN-skipping, callers filter first) loses because every caller then writes the same
filter loop. Also needs a D-id: isclose symmetric (Knuth) vs numpy's asymmetric form.

Hand loops in calcbench (from Does text only, code not opened; each sums, means, min/maxes or differences its own
slice where numpy has one routine):
- NUM-STATS-001 ArithmeticMean, NUM-STATS-025 SumOfSquares, NUM-STATS-026 WeightedMean: own sum loops (would use a
  compensated `VecSum`).
- NUM-STATS-013 Midrange, NUM-STATS-018 Range, NUM-SERIES-006 FixedStrideBuckets: own min/max scans.
- NUM-SERIES-009 AverageRealVariability: mean of |diff|, i.e. mean(abs(diff(x))).
- NUM-SERIES-003 TrailingTimeWindowSum, NUM-SERIES-008 WeeklySums: windowed sums (ALG-PREFIX-001 / ALG-WINDOW-001
  could serve).
- NUM-LINALG-042 VecNorm1, NUM-LINALG-093 MatTrace, ALG-WINDOW-001 RollingSum: each carries its own compensated sum
  (three copies of the kernel a `VecSum` would hold).
- FIN-ECON-004 HerfindahlHirschmanIndex: sum of squares of a slice.

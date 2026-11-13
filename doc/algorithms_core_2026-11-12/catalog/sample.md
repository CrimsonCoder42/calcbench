# Catalog: algo/sample (SHORT = SAMPLE), new package

Builder W3-B12, 2026-11-12. NOT in doc/taxonomy.md: created on the lead's brief (W3-B12); the taxonomy row
`ALG | algo/sample | SAMPLE | random shuffles and samples from a caller-supplied source` is for the lead/Devin.
Every function takes a *rand.Rand (math/rand/v2) from the caller; there is no global source, no clock,
no goroutine. Nil rng is ErrNilRand. Uniformity is exact given a uniform source: math/rand/v2's IntN and
Uint64N are unbiased (Lemire's multiply-and-reject, per the standard library documentation).

### ALG-SAMPLE-001 ShuffleInPlace
`func ShuffleInPlace[T any](rng *rand.Rand, values []T) error`
- **Source**: Knuth, TAOCP Vol. 2: Seminumerical Algorithms, 3rd ed., Addison-Wesley, 1997, §3.4.2,
  Algorithm P (shuffling); Durstenfeld, "Algorithm 235: Random permutation", Communications of the ACM
  7(7):420, 1964, doi:10.1145/364520.364540 (confirmed on Crossref 2026-11-12). Fisher & Yates (1938) is
  the pencil-and-paper original.
- **Best known**: O(n) time, n−1 random draws, O(1) extra; Ω(n) (every element may move).
- **Approach**: for j = n−1 down to 1, swap values[j] with values[rng.IntN(j+1)]. **Lost**: sort by random
  keys (O(n log n), ties bias it); rand.Rand.Shuffle (same algorithm, but takes a swap callback; ours is
  generic, documents the source and returns ErrNilRand instead of panicking).
- **Caveat**: a source with s bits of state reaches at most 2^s of the n! orders (PCG: 2^128 < 35!).
- **Worked example**: chi-square-free check: over 60 000 shuffles of [a b c] each of the 6 orders appears
  10 000 ± 5σ times. **Errors**: ErrNilRand. **Modifies input**: yes (name says so).
- **Tags**: shuffle, fisher-yates, knuth-shuffle, random-permutation. LeetCode 384.

### ALG-SAMPLE-002 ReservoirSample
`func ReservoirSample[T any](dst []T, rng *rand.Rand, seq iter.Seq[T], k int) ([]T, error)`
- **Source**: Vitter, "Random sampling with a reservoir", ACM Transactions on Mathematical Software
  11(1):37–57, 1985, doi:10.1145/3147.3165, Algorithm R (credited there to Waterman); Knuth Vol. 2, 3rd ed.,
  §3.4.2, Algorithm R (reservoir sampling).
- **Best known**: Algorithm R: O(N) time, N−k draws, O(k) space; Vitter's Algorithm Z: O(k(1+log(N/k)))
  draws expected, same output distribution; Li, "Reservoir-sampling algorithms of time complexity
  O(n(1 + log(N/n)))", ACM TOMS 20(4):481–493, 1994, doi:10.1145/198429.198435 (DOI not confirmed),
  Algorithm L: the same O(k(1+log(N/k))) draws with simpler arithmetic.
- **Approach**: Algorithm R: keep the first k; the t-th item (t > k, 1-based) replaces a uniform slot with
  probability k/t. **Lost**: Algorithm Z (fewer random draws, but a sequence must still be read item by
  item through iter.Seq, so the O(N) read dominates and Z's skip arithmetic gains little); Li 1994
  Algorithm L (O(k(1+log(N/k))) random draws instead of N−k, geometric skips; lost for the same reason:
  the skipped items must still be read from the iter.Seq, so only the draw count falls, not the O(N) time;
  a later variant for costly sources).
- **Order**: the reservoir's slot order, not random; shuffle if order matters.
- **Worked example**: sample 2 of 4 items 120 000 times: each of the 6 pairs appears 20 000 ± 5σ times.
- **Errors**: ErrNilRand, ErrNilSeq, ErrNegative (k < 0). k = 0 returns dst unchanged without reading seq.
  Fewer than k items: all of them. **Tags**: reservoir-sampling, streaming, random-sample. LeetCode 382, 398.

### ALG-SAMPLE-003 NewWeightedPicker · ALG-SAMPLE-004 WeightedPicker.Pick
`func NewWeightedPicker[W integer](weights []W) (*WeightedPicker[W], error)` · `func (p *WeightedPicker[W]) Pick(rng *rand.Rand) (int, error)`
- **Source**: the inversion method on a discrete distribution: Devroye, *Non-Uniform Random Variate
  Generation*, Springer, 1986, ch. III, §2 (inversion by sequential or binary search of the cumulative
  table; subsection not confirmed first-hand); Knuth Vol. 2, 3rd ed., §3.4.1 (discrete distributions). Prefix sums: ALG-PREFIX-001; search:
  ALG-SEARCH-002 UpperBound — both called, not copied.
- **Best known**: O(n) build, O(log n) per pick by binary search; O(1) per pick by Walker's alias method
  (Walker, "An efficient method for generating discrete random variables with general distributions", ACM
  TOMS 3(3):253–256, 1977, doi:10.1145/355744.355749; O(n) construction by Vose, "A linear algorithm for
  generating random numbers with a given distribution", IEEE Trans. Software Engineering 17(9):972–975,
  1991, doi:10.1109/32.92917 — DOI not confirmed). The best known pick is therefore O(1), not O(log n).
  **Lost**: the alias method (Walker 1977 / Vose 1991) — O(1) picks, but its
  table is built in floating point (or rationals) and a pick needs two draws; with integer weights the
  inversion is exact and O(log n) is fast in practice. Kept as a later variant.
- **Approach**: prefix[i] = w₀+…+w_{i−1}; draw r uniform in [0, total); answer UpperBound(prefix, r) − 1.
  Zero weights are never picked.
- **Worked example**: weights [1 3] → index 1 three times as often as 0 (r ∈ {0} → 0, r ∈ {1,2,3} → 1).
- **Errors**: ErrNegativeWeight, ErrZeroTotal (empty or all zero), ErrOverflow (total does not fit W);
  Pick: ErrNilRand, ErrNilPicker. Weights are copied (prefix sums) and never modified.
- **Tags**: weighted-random, inversion, prefix-sum, binary-search, discrete-distribution. LeetCode 528.

### ALG-SAMPLE-005 RejectionSample
`func RejectionSample[T any](rng *rand.Rand, propose func(*rand.Rand) T, accept func(candidate T, u float64) bool, maxTries int) (T, error)`
- **Source**: von Neumann, "Various techniques used in connection with random digits", in *Monte Carlo
  Method*, National Bureau of Standards Applied Mathematics Series 12, 1951, pp. 36–38; Devroye 1986,
  ch. II, §3 (the rejection method; the number of tries is geometric with mean c).
- **Best known**: expected c tries, c = sup f/(c·g) bound of the proposal; the cap turns an unbounded
  loop into an error.
- **Approach**: repeat: x ← propose(rng), u ← rng.Float64() in [0, 1); return x if accept(x, u). The caller
  writes accept as u < f(x)/(c·g(x)), or as a plain membership test ignoring u.
- **Worked example**: points uniform in the unit disc from the square [−1,1]²: acceptance rate π/4 (checked
  to 4σ over 100 000 tries). **Errors**: ErrNilRand, ErrNilFunc, ErrMaxTries (maxTries < 1 is
  ErrNonPositive). **Tags**: rejection-sampling, accept-reject, monte-carlo. LeetCode 470.

### ALG-SAMPLE-006 RandomElement
`func RandomElement[T any](rng *rand.Rand, seq iter.Seq[T]) (T, error)`
- **Name**: the brief's "LinkedListRandomNode". It takes any iter.Seq (a list's values, a channel drained by
  the caller, a file's lines), so "LinkedList" in the name would mislead; the problem is "a uniform element
  of a sequence of unknown length". Recorded for the lead.
- **Source**: Vitter 1985, Algorithm R with k = 1; calls ReservoirSample (ALG-SAMPLE-002).
- **Best known**: O(N) time (one pass), O(1) space. **Errors**: ErrNilRand, ErrNilSeq, ErrEmpty.
- **Worked example**: 5 items, 50 000 draws, each 10 000 ± 5σ. **Tags**: reservoir-sampling, streaming, random-element. LeetCode 382.

### ALG-SAMPLE-007 ShuffleDeterministic
`func ShuffleDeterministic(dst []int, n int, seed uint64) ([]int, error)`
- **Source**: ShuffleInPlace (ALG-SAMPLE-001) on 0…n−1 with a PCG source (O'Neill, "PCG: A Family of
  Simple Fast Space-Efficient Statistically Good Algorithms for Random Number Generation", Harvey Mudd
  College HMC-CS-2014-0905, 2014; math/rand/v2 NewPCG) seeded (seed, fixed stream constant).
- **Not cryptographic**: anyone who knows the seed knows the permutation; never use it for secrets.
  Reproducible within one Go toolchain; math/rand/v2 does not promise identical IntN output across Go
  releases, so tests check same-seed equality and the permutation property, not golden values.
- **Errors**: ErrNegative (n < 0), ErrTooLarge (len(dst)+n ints exceed a slice). Appends to dst. **Tags**: shuffle, seeded, deterministic, testing.

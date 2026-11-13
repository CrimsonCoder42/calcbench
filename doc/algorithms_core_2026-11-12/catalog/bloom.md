# algo/bloom (SHORT = BLOOM)

Builder W4-B14. Not yet in `doc/taxonomy.md` (lead to add: `ALG | algo/bloom | BLOOM | Bloom filters:
approximate set membership with false positives, parameter formulas`).

| ID | Name | Signature |
|---|---|---|
| ALG-BLOOM-001 | BloomFilter | `func NewBloomFilter(bitCount uint64, hashCount int) (*BloomFilter, error)` with `Add(item []byte)`, `MayContain(item []byte) bool`, `EstimatedFPR() float64`, `BitCount() uint64`, `HashCount() int`, `SetBits() uint64` |
| ALG-BLOOM-002 | FalsePositiveRate | `func FalsePositiveRate(bitCount, items uint64, hashCount int) (float64, error)` |
| ALG-BLOOM-003 | OptimalParams | `func OptimalParams(expectedItems uint64, falsePositiveRate float64) (bitCount uint64, hashCount int, err error)` |
| ALG-BLOOM-004 | BloomFilter.UnionWith | `func (f *BloomFilter) UnionWith(other *BloomFilter) error` |
| — | Intersection | **Not built**: Broder–Mitzenmacher §2.4 uses two filters only to *estimate* the size of an intersection; a bitwise-AND "intersection filter" is not in the sources read. Optional in the brief. |

## Decisions (for the lead's table)

- **B1 Hash function: SHA-256 (`crypto/sha256.Sum256`, FIPS 180-4)**, h₁ = digest bytes 0–7, h₂ = bytes 8–15
  (big-endian). Deterministic across processes, machines and Go versions, so a filter can be persisted,
  exchanged and unioned between processes and tests pin exact bit positions; its output is close to the
  "fully random" hashes the analysis assumes; `Sum256` keeps no state, so MayContain only reads and
  concurrent MayContain calls are safe (Add and UnionWith still need the caller's lock); no allocation.
  **Lost (measured)**: FNV-1a-128 (`hash/fnv`), the first choice: its two 64-bit halves are correlated
  for short keys (the prime is 2⁸⁸ + 0x13B), and a Python model of this exact construction gave a measured
  FPR of 0.00093 and a fill-based estimate of 0.0055 against the predicted 0.0100 (m = 9586, k = 7,
  1000 keys "member-i", 100 000 probes "other-i"); with SHA-256 the same run measured 0.01004 (scratch
  script, 2026-11-12). **Lost**: `hash/maphash` with a caller's `maphash.Seed` (faster and keyed, but a Seed
  cannot be serialised, so a filter means nothing outside its process, and tests could only check
  statistics). Cost: SHA-256 is slower per byte than FNV or maphash (hardware-accelerated on amd64/arm64);
  the benchmark records it. The hash is public, so an adversary choosing items can search offline for
  ones that hit set bits; a keyed variant would be a new function.
- **B2 Index arithmetic (review fix M7)**: enhanced double hashing, Dillinger and Manolios, "Bloom filters in
  probabilistic verification", FMCAD 2004, LNCS 3312, pp. 367–381, doi:10.1007/978-3-540-30494-4_26 (Crossref
  confirmed; §5.1–5.2 read): x, y = h₁ mod m, h₂ mod m; per probe x ← x + y, y ← y + i (mod m), i.e.
  gᵢ = h₁ + i·h₂ + (i³ − i)/6 mod m (eq. 6). Plain double hashing (the first build) gives one index when
  h₂ mod m = 0 and min(k, m/gcd) indices when h₂ shares a factor with m (DM §5.1); the cubic term removes
  both without restricting h₂ or making m prime. Kirsch–Mitzenmacher Theorem 5.2 does NOT assume prime m
  (§5.2: h₁, h₂ uniform on [m], any m; §3 alone uses a prime p) and covers any added f(i), so the asymptotic
  FPR is unchanged. **Lost**: step = 1 + h₂ mod (m−1) (no zero step, but still a shared factor with a
  composite m); prime m (DM's recommendation for plain double hashing; changes the caller's m). Exact,
  no overflow: every sum < 2m ≤ 2³¹ (k ≤ m). Python model with this scheme: measured FPR 0.01038,
  EstimatedFPR 0.01008, predicted 0.01003; worked example indices "apple" 41 48 56, "banana" 13 35 58.
- **B3 Size cap** (L13/L26; review fix L4: 1 ≤ k ≤ min(m, `MaxHashCount` = 1074), else `ErrHashCount` — at the optimum f = (1/2)ᵏ (survey §2.1), so k > 1074 aims below the smallest float64; OptimalParams clamps to it): bitCount ≤ `MaxBits` = 1<<30 bits (128 MiB), exported, else `ErrTooLarge`; OptimalParams
  returns `ErrTooLarge` when the m it computes exceeds the cap.
- **B4 Items are `[]byte`**. A string caller converts (`[]byte(s)`). **Lost**: a generic `T comparable` API
  (needs maphash.Comparable, see B1).

### ALG-BLOOM-001 BloomFilter
- **Source**: Bloom, "Space/time trade-offs in hash coding with allowable errors", Communications of the ACM
  13(7):422–426, 1970, doi:10.1145/362686.362692 (Crossref metadata confirmed 2026-11-12; full text not read: the
  origin, credited by both papers below). The statement read and implemented: Kirsch and Mitzenmacher, "Less hashing, same performance: building a better Bloom filter",
  Random Structures & Algorithms 33(2):187–218, 2008, doi:10.1002/rsa.20208, §2 "Standard Bloom filters"
  (read: m bits initially 0, k hash functions, set bits hᵢ(x), test whether all hᵢ(y) are 1; given the
  fraction ρ of 0 bits, the false-positive probability is (1 − ρ)ᵏ) and §5.2, Theorem 5.2 (read: double
  hashing, location h₁(u) + i·h₂(u) mod m for i ∈ [k] with h₁, h₂ independent uniform on [m], has the same
  asymptotic false-positive probability (1 − e^{−k/c})ᵏ as k independent hashes, c = m/n).
- **Best known**: Add and MayContain O(L + k) time (L = item length in bytes, hashed once; then k index steps),
  O(1) extra space; O(m) bits total. Space lower bound for any set representation with false-positive rate ε:
  n·log₂(1/ε) bits; a Bloom filter with optimal k uses log₂e ≈ 1.44 times that (Broder–Mitzenmacher §2.2).
- **Approach**: bits packed in `[]uint64`; one SHA-256 digest per call, k indices by double hashing (B2).
  A running count of set bits (updated when Add flips a 0) gives EstimatedFPR = (setBits/m)ᵏ in O(1)
  (KM §2: conditioned on ρ, FPR = (1−ρ)ᵏ), which, unlike a count of Adds, ignores duplicates and stays right
  after UnionWith. **Lost**: k independent hash functions (k digests per call; KM Theorem 5.2 shows no loss
  in the asymptotic FPR); the partitioned variant (KM §5.1; same asymptotics, harder to union with an
  unpartitioned filter); a counter of Adds for the FPR (overcounts duplicates).
- **Worked example (computed independently with Python hashlib, 2026-11-12; the test also recomputes gᵢ with
  math/big)**: m = 64, k = 3: "apple" → bits 41, 48, 56; "banana" → bits 13, 35, 58. After both Adds,
  SetBits = 6, both MayContain, EstimatedFPR = (6/64)³.
- **Errors**: `ErrBitCount` (bitCount < 1), `ErrHashCount` (hashCount < 1), `ErrTooLarge`
  (bitCount > 1<<30). Items are read, never modified or retained.
- **Tags**: bloom-filter, probabilistic, set-membership, approximate, double-hashing, sketch.

### ALG-BLOOM-002 FalsePositiveRate
- **Source**: Broder and Mitzenmacher, "Network applications of Bloom filters: a survey", Internet
  Mathematics 1(4):485–509, 2004, doi:10.1080/15427951.2004.10129096, §2.1 (read 2026-11-12): after n
  items, a given bit is still 0 with probability p′ = (1 − 1/m)^{kn} ≈ e^{−kn/m}; FPR f = (1 − e^{−kn/m})ᵏ.
  KM 2008 §2 states the same.
- **Best known**: O(1) time and space (one exp, one pow).
- **Approach**: `-math.Expm1(-k·n/m)` for 1 − e^{−kn/m} (accurate when kn/m is small), then `math.Pow(·, k)`.
  **Lost**: the exact (1 − (1 − 1/m)^{kn})ᵏ (differs by O(1/m); the brief and both sources use the
  exponential form).
- **Worked example (from survey §2.1)**: k = ln 2 · m/n gives f = (1/2)ᵏ; m/n = 10, k = 7 gives
  (1 − e^{−0.7})⁷ ≈ 0.008194 (hand-computed, checked against the closed form within 1e-12 relative).
- **Errors**: `ErrBitCount` (bitCount < 1), `ErrHashCount` (hashCount < 1). items = 0 gives 0.
- **Tags**: bloom-filter, false-positive-rate, probability.

### ALG-BLOOM-003 OptimalParams
- **Source**: Broder–Mitzenmacher 2004, §2.1 (k = ln 2 · (m/n) minimises f) and §2.2 (f ≤ ε requires
  m ≥ n·log₂(1/ε)/ln 2 = −n·ln ε/(ln 2)²) (read).
- **Best known**: O(1).
- **Approach**: m = ⌈−n·ln p/(ln 2)²⌉, k = max(1, round(m/n · ln 2)). Rounding k to the nearest integer is the
  survey's "in practice k must be an integer" (§2.1); the test checks the returned (m, k) meet p by
  FalsePositiveRate within a 1% relative margin (rounding k moves f by less than that for m ≥ 64).
  **Lost**: floor(k) (always cheaper but can exceed p by more).
- **Worked example**: n = 1000, p = 0.01 → m = ⌈9585.06⌉ = 9586, k = round(6.6446) = 7 (hand-computed).
- **Errors**: `ErrItems` (expectedItems < 1), `ErrProbability` (p not in (0, 1), or NaN), `ErrTooLarge`
  (m > 1<<30).
- **Tags**: bloom-filter, sizing, optimal-parameters.

### ALG-BLOOM-004 BloomFilter.UnionWith
- **Source**: Broder–Mitzenmacher 2004, §2.4 "Standard Bloom filter tricks" (read): two filters with the same
  number of bits and the same hash functions represent the union of their sets by the OR of their bit vectors.
- **Best known**: O(m/64) time, O(1) extra space (every word must be read).
- **Approach**: word-wise OR into the receiver, then recount set bits with `math/bits.OnesCount64`.
- **Errors**: `ErrIncompatible` when other is nil or has a different bitCount or hashCount (the hash
  function is fixed package-wide, so equal (m, k) is compatible). The receiver is modified (method on the
  filter it updates); other is read only.
- **Tags**: bloom-filter, union, set-operations, merge.

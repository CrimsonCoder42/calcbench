# Catalog: algo/bits (SHORT = BITS), new package

Builder W3-B10, 2026-11-12. Package name `bits` is the one approved in doc/taxonomy.md; the standard
library's math/bits is imported beside it under an alias (same reasoning as L24 for algo/list). Inputs are
never modified; nothing allocates except PopCountTable when dst lacks capacity.

Sources used throughout (confirmed by table of contents / search, 2026-11-12):
- Knuth, *The Art of Computer Programming*, Vol. 4A, Addison-Wesley, 2011, §7.1.3 "Bitwise tricks and
  techniques" (rightmost bits, x & −x, sideways addition ν, subset enumeration x ← (x − χ) & χ,
  bit interleaving); §7.2.1.1 "Generating all n-tuples" (Gray binary code g(k) = k ⊕ ⌊k/2⌋ and its inverse).
- Warren, *Hacker's Delight*, 2nd ed., Addison-Wesley, 2013: §2-1 "Manipulating Rightmost Bits", §3-2
  "Rounding Up/Down to the Next Power of 2" (clp2, flp2), ch. 5 "Counting Bits", §7-1 "Reversing Bits and
  Bytes", §7-2 "Shuffling Bits", ch. 13 "Gray Code".
- Gray, F., "Pulse code communication", U.S. Patent 2,632,058 (filed 1947, issued 17 Mar 1953).
- Hamming, R. W., "Error detecting and error correcting codes", Bell System Technical Journal 29(2):147–160,
  1950, doi:10.1002/j.1538-7305.1950.tb00463.x.
- Morton, G. M., "A computer oriented geodetic data base and a new technique in file sequencing", IBM Ltd.,
  Ottawa, 1966 (WorldCat OCLC 52181262; confirmed 2026-11-12).
- Laaksonen, *Competitive Programmer's Handbook* (L18), ch. 10 "Bit manipulation" (submask iteration).

## Dropped (exactly one math/bits call; logged)
- **HammingWeight / popcount** = `math/bits.OnesCount64`.
- **LowestSetBit index** = `math/bits.TrailingZeros64`; **HighestSetBit index** = `math/bits.Len64(x) − 1`.
  The VALUE forms (x & −x; the greatest power of two ≤ x) are built (BITS-007, -008): they are not one call.
- **Full 64-bit ReverseBits** = `math/bits.Reverse64`; the field form with a width is built (BITS-004).

### ALG-BITS-001 IsPowerOfTwo
`func IsPowerOfTwo[T integer](x T) bool`
- **Source**: Warren §2-1 (x & (x−1) turns off the rightmost 1); Knuth §7.1.3. **Best known**: O(1).
- **Approach**: x > 0 && x&(x−1) == 0. **Lost**: OnesCount64(x) == 1 (needs the unsigned conversion and a
  sign test anyway; same O(1)). Signed minimum is negative → false.
- **Worked example**: 2^k for k = 0…63 → true; 0, 6, −4 → false. **Errors**: none. **Tags**: bits, power-of-two. LeetCode 231.

### ALG-BITS-002 IsPowerOfFour
`func IsPowerOfFour[T integer](x T) bool`
- **Source**: as 001; 4^k = 2^(2k) has its one bit at an even index (mask 0x5555…). **Best known**: O(1).
- **Lost**: loop dividing by 4 (O(log x)); x % 3 == 1 test (correct, but a division instead of an AND).
- **Worked example**: 1, 4, 16, 2^62 → true; 2, 8, 2^63 → false. **Errors**: none. **Tags**: bits, power-of-four. LeetCode 342.

### ALG-BITS-003 NextPowerOfTwo
`func NextPowerOfTwo(x uint64) (uint64, error)`
- **Source**: Warren §3-2, clp2(x) = 2^(64 − nlz(x − 1)). **Best known**: O(1).
- **Lost**: the smear-right sequence x |= x>>1 … >>32 (six steps; LeadingZeros64 is one instruction).
- **Worked example**: 0 → 1, 1 → 1, 5 → 8, 1000 → 1024, 2^63 → 2^63.
- **Errors**: ErrOverflow for x > 2^63. **Tags**: bits, power-of-two, round-up.

### ALG-BITS-004 ReverseBits
`func ReverseBits(x uint64, width int) (uint64, error)`
- **Source**: Warren §7-1 (reverse the word, shift right by w − width); Knuth §7.1.3. **Best known**: O(1).
- **Lost**: bit-by-bit loop (O(width)); a 256-entry byte table (a global table and 8 lookups; math/bits is
  one instruction on arm64, a few on amd64).
- **Worked example**: (0b0011, 4) → 0b1100; (1, 64) → 2^63.
- **Errors**: ErrWidth (width ∉ 1…64); ErrOutOfRange (x has bits at or above width). **Tags**: bits, reverse. LeetCode 190.

### ALG-BITS-005 GrayEncode · ALG-BITS-006 GrayDecode
`func GrayEncode(x uint64) uint64` · `func GrayDecode(g uint64) uint64`
- **Source**: Gray 1953 patent; Knuth §7.2.1.1 (g(k) = k ⊕ ⌊k/2⌋, inverse k = g ⊕ ⌊g/2⌋ ⊕ ⌊g/4⌋ ⊕ …); Warren ch. 13.
- **Best known**: O(1) encode; decode O(log w) by prefix XOR (6 steps for w = 64). **Lost** for decode:
  the bit-serial loop (O(w)).
- **Worked example** (Knuth §7.2.1.1): 0…7 → 000 001 011 010 110 111 101 100, and back.
- **Errors**: none. **Tags**: bits, gray-code. LeetCode 89 (sequence = GrayEncode(0…2^n−1)).

### ALG-BITS-007 LowestSetBit · ALG-BITS-008 HighestSetBit
`func LowestSetBit(x uint64) uint64` · `func HighestSetBit(x uint64) uint64`
- **Source**: Knuth §7.1.3 "rightmost bits" (x & −x); Warren §2-1; Warren §3-2 flp2 (greatest power of 2 ≤ x).
- **Best known**: O(1). **Lost** for 008: smear-right then x − (x >> 1) (six steps vs one nlz).
- **Worked example**: 0b10110 → 0b00010 and 0b10000; 0 → 0 for both. **Errors**: none.
- **Tags**: bits, lowest-set-bit, highest-set-bit, fenwick.

### ALG-BITS-009 PopCountTable
`func PopCountTable(dst []uint8, n int) ([]uint8, error)`
- **Source**: Knuth §7.1.3 sideways addition with ν(x) = ν(⌊x/2⌋) + (x mod 2); Warren ch. 5.
- **Best known**: O(n) (the output has n+1 entries: Ω(n)). **Approach**: table recurrence, one add per entry.
  **Lost**: OnesCount per entry (also O(n) with POPCNT, but not on every target; the recurrence is
  branch-free and uses only the table).
- **Worked example**: n = 8 → 0 1 1 2 1 2 2 3 1. **Errors**: ErrNegative (n < 0); ErrTooLarge (n+1
  overflows int; beyond that, a table too large for memory is the runtime's, per L26).
- **Input modified**: no; appends to dst. **Tags**: bits, popcount, counting-bits, dp. LeetCode 338.

### ALG-BITS-010 LoneElementAmongPairs
`func LoneElementAmongPairs[T integer](values []T) T`
- **Source**: Knuth §7.1.3 (⊕ associative, commutative, x ⊕ x = 0). **Best known**: O(n) time, O(1) space
  (every value must be read). **Lost**: hash-map count (O(n) space); sort (O(n log n)).
- **Worked example**: [3 8 5 3 5] → 8. **Errors**: none; empty → 0; precondition unchecked (XOR of all).
- **Tags**: bits, xor, single-number. LeetCode 136.

### ALG-BITS-011 LoneElementAmongTriples
`func LoneElementAmongTriples[T integer](values []T) T`
- **Source**: Knuth §7.1.3 (bitwise operations act on all positions in parallel) plus the proof below; no
  textbook names this exact problem. FLAGGED for Devin as a cited-step composition (like L15).
- **Proof**: fix a bit position. If every value but one occurs 3 times, the number of values with that bit
  set is 3m + b, where b is the lone value's bit. Keep a two-bit counter (twos, ones) of that number modulo
  3. The updates ones ← (ones ⊕ v) ∧ ¬twos, then twos ← (twos ⊕ v) ∧ ¬ones map, for a set bit of v:
  (0,0) → (0,1) → (1,0) → (0,0), and leave the state unchanged for a clear bit (checked case by case: with
  v = 0, ones ← ones ∧ ¬twos = ones because the state (1,1) never occurs, and likewise twos). So after all
  values the counter is (0, b) and ones = b at every position, i.e. ones is the lone value. Two's
  complement makes this hold for signed types too.
- **Best known**: O(n) time, O(1) space. **Lost**: per-bit column count modulo 3 (O(w·n), w = 64);
  hash-map count (O(n) space).
- **Worked example**: [4 9 4 1 9 4 9] → 1. **Errors**: none; empty → 0; precondition unchecked.
- **Tags**: bits, single-number, modular-counting. LeetCode 137.

### ALG-BITS-012 BitwiseAndRange
`func BitwiseAndRange(lo, hi uint64) (uint64, error)`
- **Source**: Knuth §7.1.3 (AND; the leftmost bit of x ⊕ y is the first differing position) plus the proof
  below. FLAGGED for Devin as a cited-step composition (like L15).
- **Proof**: if lo = hi the AND is lo. Otherwise let d be the highest bit where lo and hi differ; then lo
  has 0 there and hi has 1, and they agree above d. Every v in [lo, hi] agrees with them above d (they are
  the bounds of an interval inside one aligned block of size 2^(d+1)). The value m = (hi with bits < d
  cleared) satisfies lo < m ≤ hi, and m − 1 ≥ lo; m has bits 0…d−1 clear and m − 1 has bit d clear, so
  every bit ≤ d is 0 in the AND. Result: lo with bits 0…d cleared = lo >> (d+1) << (d+1).
- **Best known**: O(1). **Lost**: shift both right until equal (O(w)); Kernighan clearing hi & (hi−1)
  while hi > lo (O(w)).
- **Worked example**: [5, 7] → 4; [12, 15] → 12; [8, 16] → 0. **Errors**: ErrRange (lo > hi).
- **Tags**: bits, and, range. LeetCode 201.

### ALG-BITS-013 HammingDistance
`func HammingDistance(x, y uint64) int`
- **Source**: Hamming 1950 (distance between code points); Knuth §7.1.3 ν(x ⊕ y). **Best known**: O(1).
- **Worked example**: 0b1011101 vs 0b1001001 → 2. **Errors**: none. **Tags**: bits, hamming-distance. LeetCode 461, 477.
- Hamming weight dropped (= math/bits.OnesCount64).

### ALG-BITS-014 Submasks
`func Submasks(mask uint64) iter.Seq[uint64]`
- **Source**: Knuth §7.1.3 (running through subsets with x ← (x − χ) & χ); Laaksonen ch. 10 (subset iteration).
- **Best known**: O(2^k) for k set bits (the output size), O(1) per submask. **Lost**: iterate 0…mask and
  test s & ^mask (O(mask), exponential in the mask's top bit, not its popcount).
- **Worked example**: 0b1101 → 1101 1100 1001 1000 0101 0100 0001 0000. **Errors**: none; early break supported.
- **Tags**: bits, submask, subset-enumeration, iterator, bitmask-dp.

### ALG-BITS-015 MortonEncode2D · ALG-BITS-016 MortonDecode2D
`func MortonEncode2D(x, y uint32) uint64` · `func MortonDecode2D(code uint64) (x, y uint32)`
- **Source**: Morton 1966 (the interleaved file sequence); Warren §7-2 "Shuffling Bits" (outer perfect shuffle
  by masks and shifts in log2 w steps); Knuth §7.1.3 (bit interleaving / zipper). Knuth names it in 7.1.3;
  equation numbers not confirmed (section-level citation).
- **Best known**: O(log w) with mask-and-shift (5 steps per coordinate), O(1) with BMI2 PDEP/PEXT (not
  reachable from portable Go). **Lost**: bit-by-bit loop (O(w)); 256-entry lookup tables (globals).
- **Worked example**: (0,0) (1,0) (0,1) (1,1) → 0 1 2 3; (2,3) → 14; round trip for every input.
- **Errors**: none. **Tags**: bits, morton, z-order, spatial-index.

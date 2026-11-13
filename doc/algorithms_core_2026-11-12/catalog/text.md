# Catalog: algo/text (SHORT = TEXT), new package

Builder W3-B11, 2026-11-12. Oracles are capped (≤ 80 bytes, ≤ 2000 trials; the exponential wildcard
oracle ≤ 9 runes) per the go-testing.md 60 s rule. New package, numbers from 001. `algo/text` is in doc/taxonomy.md
("string matching and distance").

**Conventions (stated once in doc.go).** Byte semantics by default: the matching, hashing and
palindrome functions are generic over `S ~string | ~[]byte`, index bytes, and report byte offsets
(one implementation serves both string and []byte callers without conversion). Rune semantics only
where the answer changes meaning on UTF-8 text: IsAnagram, WildcardMatch ('?' is one rune) and
RegexMatch ('.' is one rune); invalid UTF-8 bytes read as U+FFFD there (the Go conversion rule, as
ALG-HASH-010). ReverseWords and StringToInt use ASCII white space, which is UTF-8 safe (no multi-byte
rune contains an ASCII byte). Searches allocate only their pattern tables (KMP: the prefix function,
m ints); Rabin–Karp allocates nothing beyond dst. Large outputs append to dst (D-003 style).
Lost alternative: separate string and []byte functions (doubles the API) / rune-level everywhere
(the classical algorithms are defined on symbol arrays; bytes are the symbols of a Go string, and
rune decoding would make every index ambiguous).

**Confirmations (2026-11-12)**: Crossref confirmed KMP 1977 (10.1137/0206024), Karp–Rabin 1987
(10.1147/rd.312.0249), Manacher 1975 (10.1145/321892.321896), Apostolico–Breslauer–Galil 1995
(10.1016/0304-3975(94)00083-U), Foster 1947 (10.2307/3029479). Read first hand: SemVer 2.0.0 (§2, §9–§11,
BNF), RFC 8259 §6 ABNF, ISO/IEC 9899:2011 draft N1570 §7.22.1.4 (strtol), TIFF 6.0 Section 9 (p. 42).
Thompson 1968 (10.1145/363347.363387): Crossref returned 429; citation from memory, reviewer to check.
CLRS, Gusfield, Bentley, Knuth, Aho–Lam–Sethi–Ullman, Laaksonen: textbook sections cited from memory.

### ALG-TEXT-001 PrefixFunction
`func PrefixFunction[S ~string | ~[]byte](dst []int, s S) []int`
- **Source**: Knuth, Morris, Pratt, "Fast Pattern Matching in Strings", SIAM J. Comput. 6(2):323–350, 1977,
  doi:10.1137/0206024; CLRS 3rd ed. §32.4, COMPUTE-PREFIX-FUNCTION.
- **Best known**: O(m) time (amortised: k rises at most m times), Ω(m) to read s. O(1) extra.
- **Approach**: CLRS's loop, π[q] = length of the longest proper prefix of s[:q+1] that is also its
  suffix (0-based). **Lost**: KMP's optimised `next` table (same bound; the plain π is what callers reuse
  for periods and borders).
- **Worked example**: CLRS Fig. 32.10: ababaca → 0 0 1 2 3 0 1. **Errors**: none. **Tags**: LeetCode 1392, 459.

### ALG-TEXT-002 KMPSearch
`func KMPSearch[S ~string | ~[]byte](dst []int, text, pattern S) ([]int, error)`
- **Source**: as 001 (KMP 1977; CLRS §32.4 KMP-MATCHER).
- **Best known**: O(n + m) worst case, O(m) extra (the π table); Ω(n + m).
- **Approach**: KMP matcher appending every start (overlapping). **Lost**: Boyer–Moore (sublinear on
  average, but O(nm) for all occurrences in its basic form); naive O(nm).
- **Worked example**: (hand) text abababacaba, pattern ababaca (CLRS Fig. 32.10's pattern) → [2].
- **Errors**: ErrEmptyPattern. **Tags**: LeetCode 28 (StrStr: DROPPED — first match only is exactly
  strings.Index / bytes.Index, which the standard library already does in O(n + m)).

### ALG-TEXT-003 ZFunction
`func ZFunction[S ~string | ~[]byte](dst []int, s S) []int`
- **Source**: Gusfield, *Algorithms on Strings, Trees, and Sequences*, Cambridge UP, 1997, §1.3 (Z values)
  and §1.4 (the Z algorithm, linear time).
- **Best known**: O(n), O(1) extra. **Approach**: the Z-box [l, r) algorithm. **Lost**: naive O(n²).
- **Convention**: z[0] = n (Gusfield leaves Z₁ undefined). **Worked example**: Gusfield §1.3, aabcaabxaaz:
  Z₅ = 3, Z₆ = 1, Z₇ = Z₈ = 0, Z₉ = 2 (1-based) → z = 11 1 0 0 3 1 0 0 2 1 0. **Errors**: none.

### ALG-TEXT-004 RollingHash (NewRollingHash; methods Append, Roll, Sum, Len, Reset — one ID, L23)
`func NewRollingHash(base, modulus uint64, windowLen int) (RollingHash, error)`
- **Source**: Karp, Rabin, "Efficient randomized pattern-matching algorithms", IBM J. Res. Dev.
  31(2):249–260, 1987, doi:10.1147/rd.312.0249; CLRS §32.2 (Horner's rule and the rolling update mod q).
- **Cost**: O(1) per Append/Roll (one 128-bit product and remainder, math/bits). Zero value inert (M2). No collision-probability claim (M6: not stated in Karp–Rabin or CLRS §32.2). **Lost**: Mersenne-only
  2⁶¹−1 reduction (faster, but cannot reproduce CLRS's q = 13 example); a fixed modulus 2⁶⁴ (no prime,
  easy to break with Thue–Morse strings).
- **Worked example**: CLRS Fig. 32.5 (base 10, q = 13, window 5): windows 31415 and 67399 hash alike.
- **Errors**: ErrHashParams (modulus < 2, base < 1 or ≥ modulus), ErrWindowLen (< 1).

### ALG-TEXT-005 RabinKarpSearch
`func RabinKarpSearch[S ~string | ~[]byte](dst []int, text, pattern S, base, modulus uint64) ([]int, error)`
- **Source**: as 004. Every fingerprint hit is verified byte by byte, so there are no false positives.
- **Best known**: O(n + m) expected for a random base; worst O(nm) (every window a hit). O(1) extra.
- **Lost**: KMP (002) — deterministic, but Rabin–Karp extends to many patterns and 2-D. Both kept.
- **Worked example**: CLRS Fig. 32.5: 2359023141526739921, 31415, q = 13 → [6]; spurious hit at 12 rejected.
- **Errors**: ErrEmptyPattern, ErrHashParams.

### ALG-TEXT-006 RepeatedSubstrings
`func RepeatedSubstrings[S ~string | ~[]byte](dst []int, text S, length int) ([]int, error)`
- **Source**: Karp–Rabin 1987 (fingerprints of all windows); every hash match verified.
- **Cost**: O(n + L·h) expected, h = windows whose fingerprint equals a stored one (O(nL) worst, e.g. aⁿ);
  O(n) extra for the table. **Lost**: suffix array + LCP (deterministic O(n log n)); much more code.
- **Output**: the start of the first occurrence of each distinct length-L substring that occurs ≥ 2 times
  (overlaps count), ascending. **Worked example** (hand): xyzxyzxy, L = 3 → [0, 1, 2].
- **Errors**: ErrLength (L < 1). **Tags**: LeetCode 187, 1044 (fixed-L step).

### ALG-TEXT-007 LongestCommonPrefix
`func LongestCommonPrefix[S ~string | ~[]byte](strs []S) int`
- **Source**: Knuth, TAOCP Vol. 3, 2nd ed., §6.3 (digital searching: the common prefix of a key set is the
  trie path above the first branching node). Walking that path column by column needs no trie.
- **Best known**: O(k·p + k) where p = result, k = len(strs); Ω the same (every string must be read to p).
- **Lost**: build a trie (ALG-TRIE) O(Σ|s|) and memory; sort then compare first/last O(Σ|s| log k).
- **Worked example** (hand): carpet carpool cartoon → 3. **Errors**: none; no strings → 0.

### ALG-TEXT-008 PalindromeRadii (Manacher — the one home)
`func PalindromeRadii[S ~string | ~[]byte](dst []int, s S) []int`
- **Source**: Manacher, JACM 22(3):346–351, 1975, doi:10.1145/321892.321896; Apostolico, Breslauer, Galil,
  "Parallel detection of all palindromes in a string", TCS 141:163–173, 1995 (Manacher's method yields all
  maximal palindromes, odd and even centres).
- **Best known**: O(n), Ω(n). **Approach**: Manacher over the implicit separator-interleaved string
  (no copy). **Lost**: expand around centre O(n²) (that one is ALG-DP's).
- **Output**: 2n−1 lengths; entry c is the longest palindrome centred at byte c/2 (c even) or between
  bytes (c−1)/2 and (c+1)/2 (c odd). **Worked example** (hand): abaaba → 1 0 3 0 1 6 1 0 3 0 1.

### ALG-TEXT-009 LongestPalindromicSubstring
`func LongestPalindromicSubstring[S ~string | ~[]byte](s S) (start, length int)` — via 008; leftmost on ties.
- **Worked example** (hand): xyracecarz → (2, 7). **Tags**: LeetCode 5, 647.

### ALG-TEXT-010 IsAnagram
`func IsAnagram(a, b string) bool` — rune multisets equal.
- **Source**: Bentley, *Programming Pearls*, 2nd ed., Column 2 (anagram signatures); CLRS §8.2 (counting).
- **Best known**: O(n) expected; ASCII fast path O(1) extra (a 128-entry table), else O(d) map.
- **Lost**: sort the runes (Bentley's signature; O(n log n)). **Tags**: LeetCode 242.

### ALG-TEXT-011 RunLengthEncode / ALG-TEXT-012 RunLengthDecode (PackBits)
`func RunLengthEncode(dst, src []byte) []byte` · `func RunLengthDecode(dst, src []byte, maxLen int) ([]byte, error)`
- **Format**: Apple PackBits as specified in TIFF Revision 6.0, Adobe, 3 June 1992, Section 9 "PackBits
  Compression", p. 42: header n (int8) 0…127 → n+1 literal bytes follow; −127…−1 → the next byte
  repeated 1−n times; −128 → no-op. Binary safe and unambiguous (a text "3a2b" form is ambiguous for
  digit symbols and has no standard).
- **Encoder policy** (TIFF's advice): runs of ≥ 3 → replicate; a run of 2 joins the literal in progress
  when it is followed by a single byte, otherwise replicate; runs longer than 128 split. Output ≤
  n + ⌊n/3⌋ + 1 bytes (each non-final, non-full literal header is followed by a replicate run, so it
  covers ≥ 3 input bytes; TIFF's "1 extra byte per 128" is the toolbox encoder's bound, not this policy's). **Decode**: ErrTruncated (header without its data), ErrTooLong (output would exceed
  maxLen), ErrMaxLen (maxLen < 0). Counts are at most 128 per header, so no count can overflow; maxLen is
  the cap against decompression bombs (expansion ≤ 64×).
- **Worked example** (hand, from the rules): AA AA AA 80 00 2A AA AA AA AA 80 00 2A 22 AA×10 →
  FE AA 02 80 00 2A FD AA 03 80 00 2A 22 F7 AA. **Time** O(n + output), O(1) extra.

### ALG-TEXT-013 ReverseWords
`func ReverseWords(s string) string`
- **Source**: Bentley, *Programming Pearls*, 2nd ed., Column 2.3 (rotation by reversals: reverse the
  whole, then each piece). Words = maximal runs of non-ASCII-white-space bytes; output joined by one space.
- **Best known**: O(n), one allocation. **Lost**: strings.Fields + slices.Reverse + strings.Join (same
  output for ASCII white space, but k+2 allocations; Fields also splits on Unicode spaces).

### ALG-TEXT-014 MinimumWindowSubstring
`func MinimumWindowSubstring(text, target string) (start, length int, found bool)` — byte multiset.
- **Source**: Laaksonen, *Competitive Programmer's Handbook*, 2018, §8.1 "Two pointers method" (L18):
  right pointer extends until the window covers target; left pointer advances while it still covers.
  Each pointer moves n times → O(n + m), 256-entry table O(1) extra. **Lost**: binary search on length
  O(n log n). Ties: leftmost. Empty target → (0, 0, true). **Tags**: LeetCode 76.

### ALG-TEXT-015 StringToInt
`func StringToInt(s string) (value int64, n int, err error)`
- **Source**: ISO/IEC 9899:2011 §7.22.1.4 paragraphs 2–5 (strtol, base 10), read in draft N1570;
  white space per §7.4.1.10 ("C" locale: space, \f \n \r \t \v).
- **Rule**: skip white space, optional sign, longest digit run; n = bytes consumed (0 when there are no
  digits: ErrNoDigits). Overflow: ErrOverflow with value 0 (not strtol's clamp). **Not strconv.ParseInt**:
  ParseInt rejects leading white space and trailing text and does not report how much it consumed.
- **Tags**: atoi, strtol. LeetCode 8 (its clamping rule is not ours).

### ALG-TEXT-016 IsJSONNumber
`func IsJSONNumber(s string) bool`
- **Source**: RFC 8259 (STD 90), 2017, §6 ABNF `number = [ minus ] int [ frac ] [ exp ]`.
- **Approach**: a 9-state DFA, O(n), no allocation. **Differs from strconv.ParseFloat**, which also
  accepts +1, .5, 5., 01, Inf, NaN, hex floats and underscores-free forms, and rejects out-of-range values
  that JSON's grammar accepts (1e999). json.Valid also accepts strings, objects and surrounding white space.

### ALG-TEXT-017 ColumnLabel / ALG-TEXT-018 ColumnNumber
`func ColumnLabel(number int) (string, error)` · `func ColumnNumber(label string) (int, error)`
- **Source**: Foster, "A Number System without a Zero-Symbol", Mathematics Magazine 21(1):39–41, 1947,
  doi:10.2307/3029479 (bijective base-k numeration; here k = 26, digits A…Z = 1…26).
- **Worked example** (hand): 1 A · 26 Z · 27 AA · 52 AZ · 53 BA · 702 ZZ · 703 AAA · 16384 XFD.
- **Errors**: ErrColumnNumber (< 1), ErrColumnLabel (empty or not A–Z), ErrOverflow (> MaxInt).
  Names avoid the product name (spreadsheet columns; tag excel). **Tags**: LeetCode 168, 171.

### ALG-TEXT-019 CompareSemVer
`func CompareSemVer(a, b string) (int, error)`
- **Source**: Semantic Versioning 2.0.0, https://semver.org/spec/v2.0.0.html (read 2026-11-12): §2 (X.Y.Z,
  no leading zeros), §9 (pre-release), §10 (build metadata ignored), §11 (precedence), the BNF.
- **Approach**: validate and compare in one scan without allocation; numeric parts compared as digit
  strings (no size limit, no overflow). **Lost**: plain dotted integers ("1.01" = "1.1") — no citable
  rule; HELD as CompareVersions. **Worked example**: §11's chain 1.0.0-alpha < … < 1.0.0.
- **Errors**: ErrSemVer.

### ALG-TEXT-020 WildcardMatch
`func WildcardMatch(text, pattern string) bool` — '?' one rune, '*' any run, whole text.
- **Source**: semantics: IEEE Std 1003.1-2017 (POSIX), XCU §2.13.1 (? and *, without brackets);
  method: NFA simulation, Thompson 1968 CACM 11(6):419–422, doi:10.1145/363347.363387; Aho, Lam, Sethi,
  Ullman, *Compilers*, 2nd ed., 2006, §3.7.2 Algorithm 3.22. Recurrence D[i][j] (text i runes, pattern j):
  '*' → D[i][j−1] ∨ D[i−1][j]; '?' or equal → D[i−1][j−1].
- **Cost**: O(n·m), O(m) extra (two rows). **Lost**: greedy with backtracking (O(1) space, no acceptable
  source found). **Not path.Match**: no brackets, no escapes, '*' crosses '/'. **Tags**: glob. LeetCode 44.

### ALG-TEXT-021 RegexMatch
`func RegexMatch(text, pattern string) (bool, error)` — '.' any rune, 'x*' zero or more x, whole text.
- **Source**: as 020 (Thompson 1968; ALSU §3.3 regular expressions, §3.7.2 simulation). Recurrence: x* →
  D[i][j−1] ∨ (x matches tᵢ ∧ D[i−1][j]); x → x matches tᵢ ∧ D[i−1][j−1].
- **Cost**: O(n·m), O(m) extra. **Errors**: ErrPattern ('*' first or after '*'). **Tags**: LeetCode 10.

### HELD / DROPPED
- **IntegerToRoman / RomanToInteger — HELD (no source)**: no primary reference for the strict subtractive
  grammar could be read (encyclopediaofmath.org returned 502; only secondary pages). Candidates for Devin:
  Ifrah, *The Universal History of Numbers*, Wiley 2000, ch. 16; Menninger, *Number Words and Number Symbols*.
- **StrStr — DROPPED**: identical to strings.Index / bytes.Index.
- **CompareVersions (plain dotted integers) — HELD**: no citable rule; CompareSemVer built instead.

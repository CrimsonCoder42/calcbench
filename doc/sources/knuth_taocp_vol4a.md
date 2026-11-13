# Knuth

**Citations** (as written in the doc comments):
- Knuth, The Art of Computer Programming, Vol. 4A, Addison-Wesley, 2011, §7.1.3 "Bitwise tricks and techniques" (bit interleaving, the "zipper" of two words, by δ-shifts)
- Knuth, The Art of Computer Programming, Vol. 4A, Addison-Wesley, 2011, §7.1.3 "Bitwise tricks and techniques" (bitwise AND; the leftmost bit of x ⊕ y marks the first position where x and y differ); proof of the range identity in catalog/bits.md
- Knuth, The Art of Computer Programming, Vol. 4A, Addison-Wesley, 2011, §7.1.3 "Bitwise tricks and techniques" (bitwise operations act on all bit positions in parallel); the per-position counter modulo 3 and the proof are in doc/algorithms_core_2026-11-12/catalog/bits.md (ALG-BITS-011)
- Knuth, The Art of Computer Programming, Vol. 4A, Addison-Wesley, 2011, §7.1.3 "Bitwise tricks and techniques" (running through all subsets of a set represented by a mask with x ← (x − χ) & χ)
- Knuth, The Art of Computer Programming, Vol. 4A, Addison-Wesley, 2011, §7.1.3 "Bitwise tricks and techniques" (x ⊕ x = 0, x ⊕ 0 = x; ⊕ is associative and commutative)
- Knuth, The Art of Computer Programming, Vol. 4A, Addison-Wesley, 2011, §7.1.3 "Bitwise tricks and techniques", "Working with the rightmost bits"
- Knuth, The Art of Computer Programming, Vol. 4A, Addison-Wesley, 2011, §7.1.3 "Bitwise tricks and techniques", "Working with the rightmost bits" (x & −x extracts the rightmost 1)
- Knuth, The Art of Computer Programming, Vol. 4A, Addison-Wesley, 2011, §7.1.3 "Bitwise tricks and techniques", bit reversal by swaps
- Knuth, The Art of Computer Programming, Vol. 4A, Addison-Wesley, 2011, §7.1.3 "Bitwise tricks and techniques", sideways addition ν(x ⊕ y)
- Knuth, The Art of Computer Programming, Vol. 4A, Addison-Wesley, 2011, §7.1.3 "Bitwise tricks and techniques", sideways addition ν(x) and ν(x) = ν(⌊x/2⌋) + (x mod 2)
- Knuth, The Art of Computer Programming, Vol. 4A, Addison-Wesley, 2011, §7.1.3 "Bitwise tricks and techniques", the leftmost bit 2^⌊lg x⌋
- Knuth, The Art of Computer Programming, Vol. 4A, Addison-Wesley, 2011, §7.2.1.1 "Generating all n-tuples" (inverse Gray code: k = g ⊕ ⌊g/2⌋ ⊕ ⌊g/4⌋ ⊕ …)
- Knuth, The Art of Computer Programming, Vol. 4A, Addison-Wesley, 2011, §7.2.1.1 "Generating all n-tuples", Gray binary code g(k) = k ⊕ ⌊k/2⌋
- Knuth, The Art of Computer Programming, Vol. 4A: Combinatorial Algorithms, Part 1, Addison-Wesley, 2011, §7.2.1.1 (the binary reflected Gray code, g(k) = k ⊕ ⌊k/2⌋)
- Knuth, The Art of Computer Programming, Vol. 4A: Combinatorial Algorithms, Part 1, Addison-Wesley, 2011, §7.2.1.1, Algorithm M (mixed-radix generation) with every radix 2
- Knuth, The Art of Computer Programming, Vol. 4A: Combinatorial Algorithms, Part 1, Addison-Wesley, 2011, §7.2.1.1, Algorithm M (mixed-radix generation, radix mⱼ+1 for the j-th distinct value)
- Knuth, The Art of Computer Programming, Vol. 4A: Combinatorial Algorithms, Part 1, Addison-Wesley, 2011, §7.2.1.2, Algorithm L (lexicographic permutation generation, steps L2–L4)
- Knuth, The Art of Computer Programming, Vol. 4A: Combinatorial Algorithms, Part 1, Addison-Wesley, 2011, §7.2.1.2, Algorithm L (stated for a multiset a₁ ≤ … ≤ aₙ); calls NextPermutationInPlace (ALG-ENUM-001)
- Knuth, The Art of Computer Programming, Vol. 4A: Combinatorial Algorithms, Part 1, Addison-Wesley, 2011, §7.2.1.2, Algorithm L, run on the index array 0…n−1 with every interchange mirrored on the values
- Knuth, The Art of Computer Programming, Vol. 4A: Combinatorial Algorithms, Part 1, Addison-Wesley, 2011, §7.2.1.3 (generating all combinations)
- Knuth, The Art of Computer Programming, Vol. 4A: Combinatorial Algorithms, Part 1, Addison-Wesley, 2011, §7.2.1.6, Algorithm P (nested parentheses in lexicographic order)

**Used by**: ALG-BITS-001, ALG-BITS-004, ALG-BITS-005, ALG-BITS-006, ALG-BITS-007, ALG-BITS-008, ALG-BITS-009, ALG-BITS-010, ALG-BITS-011, ALG-BITS-012, ALG-BITS-013, ALG-BITS-014, ALG-BITS-015, ALG-ENUM-001, ALG-ENUM-002, ALG-ENUM-003, ALG-ENUM-004, ALG-ENUM-005, ALG-ENUM-006, ALG-ENUM-008, ALG-ENUM-009.

**Checked**: the builder's catalog section in doc/algorithms_core_2026-11-12/catalog/ states how each citation was confirmed (Crossref DOI, table of contents, or the text read); citations not read first-hand are listed in doc/algorithms_core_2026-11-12/REPORT.md. The worked examples are reproduced in the functions' tests.

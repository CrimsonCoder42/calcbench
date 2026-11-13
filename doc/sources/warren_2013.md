# Warren

**Citations** (as written in the doc comments):
- Warren, Hacker's Delight, 2nd ed., Addison-Wesley, 2013, ch. 13 "Gray Code", §13-1 (conversion from Gray code by a prefix XOR in log2 w steps)
- Warren, Hacker's Delight, 2nd ed., Addison-Wesley, 2013, ch. 5 "Counting Bits", §5-1 "Counting 1-Bits"
- Warren, Hacker's Delight, 2nd ed., Addison-Wesley, 2013, ch. 9 "Integer Division" (shift-and-subtract restoring division; truncating signed division from unsigned)
- Warren, Hacker's Delight, 2nd ed., Addison-Wesley, 2013, §2-1 "Manipulating Rightmost Bits" (power-of-two test); 4^k = 2^(2k), so the single bit sits at an even index (mask test)
- Warren, Hacker's Delight, 2nd ed., Addison-Wesley, 2013, §2-1 "Manipulating Rightmost Bits" (x & (x−1) turns off the rightmost 1-bit; "can be used to determine if an unsigned integer is a power of 2")
- Warren, Hacker's Delight, 2nd ed., Addison-Wesley, 2013, §2-1 "Manipulating Rightmost Bits" (x & (−x) isolates the rightmost 1-bit)
- Warren, Hacker's Delight, 2nd ed., Addison-Wesley, 2013, §3-2 "Rounding Up/Down to the Next Power of 2" (clp2: 2^(w − nlz(x − 1)))
- Warren, Hacker's Delight, 2nd ed., Addison-Wesley, 2013, §3-2 "Rounding Up/Down to the Next Power of 2" (flp2: the greatest power of 2 ≤ x, 0 for x = 0)
- Warren, Hacker's Delight, 2nd ed., Addison-Wesley, 2013, §7-1 "Reversing Bits and Bytes" (reverse the word, then shift right by w − width to reverse a field)
- Warren, Hacker's Delight, 2nd ed., Addison-Wesley, 2013, §7-2 "Shuffling Bits" (inverse perfect shuffle)
- Warren, Hacker's Delight, 2nd ed., Addison-Wesley, 2013, §7-2 "Shuffling Bits" (outer perfect shuffle by masks and shifts in log2 w steps)

**Used by**: ALG-BITS-001, ALG-BITS-002, ALG-BITS-003, ALG-BITS-004, ALG-BITS-006, ALG-BITS-007, ALG-BITS-008, ALG-BITS-009, ALG-BITS-015, ALG-BITS-016, ALG-INTMTH-003.

**Checked**: the builder's catalog section in doc/algorithms_core_2026-11-12/catalog/ states how each citation was confirmed (Crossref DOI, table of contents, or the text read); citations not read first-hand are listed in doc/algorithms_core_2026-11-12/REPORT.md. The worked examples are reproduced in the functions' tests.

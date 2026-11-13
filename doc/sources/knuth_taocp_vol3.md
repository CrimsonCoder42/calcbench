# Knuth, The Art of Computer Programming, Vol. 3

**Citation**: Knuth, *The Art of Computer Programming*, Vol. 3: *Sorting and Searching*, 2nd ed., Addison-Wesley, 1998.

| Taken | Locator | Used by | How checked |
|---|---|---|---|
| Binary search (Algorithm B) and the ⌈log2(n+1)⌉ comparison bound | §6.2.1 | ALG-SEARCH-001…004 | Oracle tests against a linear scan, 500 fixed-seed inputs per function |
| The 16 running-example keys 503 087 512 061 908 170 897 275 653 426 154 509 612 677 765 703 | ch. 5 (e.g. §5.2) and §6.2.1 | test data for ALG-SEARCH-*, ALG-SELECT-*, ALG-HEAP-001, ALG-MERGE-001 | Expected results hand-derived from the sorted keys |
| Search for K = 653 (the 11th key) | §6.2.1 | ALG-SEARCH-* worked-example tests | NOT verified against a copy here (no copy available); see TODO.md |
| Multiway merging | §5.4.1 | ALG-MERGE-001 (alternative noted: tree of losers) | Oracle test against concatenate-and-sort |


## Used by S10 (algorithms core, 2026-11-12)

IDs: ALG-HASH-006, ALG-HASH-007, ALG-HASH-008, ALG-INTVL-002, ALG-INTVL-003, ALG-INTVL-005, ALG-LIST-011, ALG-MATRIX-006, ALG-SAMPLE-004, ALG-SEARCH-005, ALG-SEARCH-006, ALG-SEARCH-008, ALG-SEARCH-009, ALG-SORT-001, ALG-SORT-002, ALG-SORT-003, ALG-SORT-004, ALG-SORT-006, ALG-SORT-009, ALG-STREAM-003, ALG-TEXT-007, ALG-TREE-011, ALG-TRIE-001, ALG-TRIE-002, ALG-TRIE-003, ALG-TRIE-004, ALG-TRIE-006, ALG-TRIE-008, ALG-TRIE-009, ALG-TRIE-010, ALG-TRIE-012, ALG-TWOPTR-004.

Citations as written in the doc comments:
- Knuth, The Art of Computer Programming, Vol. 3, 2nd ed., Addison-Wesley, 1998, §5.2.4 (one-pass merge of sorted sequences; here a list and one element)
- Knuth, The Art of Computer Programming, Vol. 3, 2nd ed., Addison-Wesley, 1998, §5.2.4 (one-pass merge of two sorted sequences: the range and values)
- Knuth, The Art of Computer Programming, Vol. 3, 2nd ed., Addison-Wesley, 1998, §5.2.4, Algorithm L (list merge sort; its merge relinks two sorted linked lists)
- Knuth, The Art of Computer Programming, Vol. 3, 2nd ed., Addison-Wesley, 1998, §5.2.4, Algorithm M (two-way merge: advance the list whose current item ends first)
- Knuth, The Art of Computer Programming, Vol. 3, 2nd ed., Addison-Wesley, 1998, §6.4 (hashing)
- Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, ch. 5 introduction (sorted order K_p(1) ≤ K_p(2) ≤ … ≤ K_p(N))
- Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §5.2, Algorithm D (distribution counting)
- Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §5.2.1, Algorithm S (straight insertion)
- Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §5.2.4 "Sorting by Merging", Algorithm M (two-way merge), p. 158 (run here from the high ends)
- Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §5.2.4 (merge sort); prefix sums via PrefixSums (ALG-PREFIX-001)
- Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §5.2.4 (sorting by merging)
- Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §5.2.5 (radix sorting)
- Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.2.1, Algorithm B (binary search in an ordered table), via ALG-SEARCH-002
- Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.2.1, Algorithm B (binary search)
- Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.2.1, Algorithm B (binary search), applied twice
- Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.2.1, Algorithm B (binary search), over the rows and then within one row
- Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.2.3 (balanced trees: Adel'son-Vel'skii and Landis's condition, subtree heights differ by at most one)
- Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.3 "Digital Searching"
- Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.3 "Digital Searching" (Algorithm T, trie search)
- Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.3 "Digital Searching" (a trie's subtree holds exactly the keys with that prefix, in lexicographic order)
- Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.3 "Digital Searching" (all keys with a prefix lie in one subtree)
- Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.3 "Digital Searching" (deletion from a trie)
- Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.3 "Digital Searching" (the search path passes every stored prefix of the argument, shortest first)
- Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.3 "Digital Searching" (trie insertion)
- Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.3 "Digital Searching" (trie search, one branch per key position)
- Knuth, The Art of Computer Programming, Vol. 3: Sorting and Searching, 2nd ed., Addison-Wesley, 1998, §6.3 (digital searching: the prefix shared by a set of keys is the trie path above its first branching node, walked here column by column without building the trie)

Checked: see the builder's catalog section in doc/algorithms_core_2026-11-12/catalog/ (confirmation status per citation; not first-hand items are listed in PROGRESS.md and REPORT.md).

# Kleinberg

**Citations** (as written in the doc comments):
- Kleinberg, Tardos, Algorithm Design, Addison-Wesley, 2006, ch. 5, Solved Exercise 1 (peak entry of a unimodal array with O(log n) probes)
- Kleinberg, Tardos, Algorithm Design, Addison-Wesley, 2006, ch. 6, Exercise 5 (segmenting a string of letters into words by dynamic programming over prefixes)
- Kleinberg, Tardos, Algorithm Design, Addison-Wesley, 2006, §4.1 (interval scheduling: earliest finish time, "the greedy algorithm stays ahead")
- Kleinberg, Tardos, Algorithm Design, Addison-Wesley, 2006, §4.1, "A Related Problem: Scheduling All Intervals" (interval partitioning; rooms needed = depth)
- Kleinberg, Tardos, Algorithm Design, Addison-Wesley, 2006, §5.3 "Counting Inversions" (Merge-and-Count: the pairs i < j split across sorted halves, examined with two pointers while merging)
- Kleinberg, Tardos, Algorithm Design, Addison-Wesley, 2006, §6.1 (the same recurrence OPT(j) = max(v_j + OPT(p(j)), OPT(j−1)); on a path, p(j) = j−2)

**Used by**: ALG-DP-011, ALG-DP-015, ALG-INTVL-004, ALG-INTVL-008, ALG-MATRIX-006, ALG-SEARCH-007.

**Checked**: the builder's catalog section in doc/algorithms_core_2026-11-12/catalog/ states how each citation was confirmed (Crossref DOI, table of contents, or the text read); citations not read first-hand are listed in doc/algorithms_core_2026-11-12/REPORT.md. The worked examples are reproduced in the functions' tests.

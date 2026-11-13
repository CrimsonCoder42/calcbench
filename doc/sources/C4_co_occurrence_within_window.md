# C4 (pending correlation.go:172) → ALG-WINDOW-016 CoOccurrenceWithinWindow (PORTED)
- Pattern, no formula: sort by time, two pointers bracketing [t − w, t + w], a set of labels present with O(1) insert/delete. Laaksonen, Competitive Programmer's Handbook, draft 3 July 2018, §8.1 "Two pointers method", p. 77 (as cited across algo/window); CLRS 3rd ed. ch. 11, 17.
- Definition matched to the source: Count(A, B) = number of events labelled A having ≥ 1 event labelled B with |Δt| ≤ w (the source's JointCount for names A < B; its `break` counts each A once). The source's rate JointCount/min(countA, countB) is not ported (it can exceed 1: CONSUMER DEFECTS).
- Complexity: source O(names²·n²) → O(n log n + P).
- Checks: O(n²) oracle on 300 random sets; source-shaped case JointCount 2. Read 2026-11-12.

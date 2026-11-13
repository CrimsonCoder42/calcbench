# Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms

**Citation**: Cormen, Leiserson, Rivest, Stein, *Introduction to Algorithms*, 3rd ed., MIT Press, 2009.

| Taken | Locator | Used by | How checked |
|---|---|---|---|
| Lower and upper median as order statistics | ch. 9 introduction | ALG-SELECT-002 | Oracle test against sort-and-index |
| Randomized and worst-case linear selection; the 3n/10 − 6 median-of-medians bound | §9.2–9.3 | ALG-SELECT-001 | `TestMedianOfMedians_guarantee` checks the bound on 200 inputs |
| k-way merge with a min-heap | §6.5, Exercise 6.5-9 | ALG-MERGE-001 | Oracle test |
| Product of polynomials in coefficient form, c_j = Σ a_k b_{j−k}, Θ(n²) | §30.1, eqs. (30.1)–(30.2) | NUM-ALGEB-003 | Hand-checked example; oracle vs exact convolution. Locator via walkccc.me/CLRS/Chap30/30.1 (Exercise 30.1-1 "using equations (30.1) and (30.2)") and a search summary of the CLRS ch. 30 PDF |


## Used by S10 (algorithms core, 2026-11-12)

IDs: ALG-DP-002, ALG-DP-003, ALG-DP-006, ALG-DP-007, ALG-DP-008, ALG-DP-013, ALG-DP-014, ALG-DP-018, ALG-DP-019, ALG-DP-020, ALG-DP-021, ALG-DP-024, ALG-GEO2D-001, ALG-GEO2D-002, ALG-GEO2D-003, ALG-GEO2D-005, ALG-GRAPH-001, ALG-GRAPH-002, ALG-GRAPH-003, ALG-GRAPH-004, ALG-GRAPH-005, ALG-GRAPH-006, ALG-GRAPH-007, ALG-GRAPH-008, ALG-GRAPH-009, ALG-GRAPH-010, ALG-GRAPH-011, ALG-GRAPH-012, ALG-GRAPH-013, ALG-GRAPH-014, ALG-GRAPH-015, ALG-GRAPH-016, ALG-GRAPH-017, ALG-GRAPH-018, ALG-GRAPH-019, ALG-GRAPH-020, ALG-GRAPH-021, ALG-GRAPH-022, ALG-GRAPH-023, ALG-GRAPH-024, ALG-GRAPH-025, ALG-GRAPH-026, ALG-HASH-001, ALG-HASH-002, ALG-HASH-004, ALG-HASH-005, ALG-HASH-006, ALG-HASH-007, ALG-HASH-008, ALG-HASH-009, ALG-HASH-010, ALG-HASH-011, ALG-HASH-012, ALG-HASH-013, ALG-HEAP-002, ALG-HEAP-003, ALG-HEAP-004, ALG-HEAP-005, ALG-HEAP-006, ALG-HEAP-007, ALG-HEAP-008, ALG-HEAP-009, ALG-HEAP-012, ALG-INTMTH-001, ALG-INTVL-004, ALG-INTVL-008, ALG-LIST-001, ALG-LIST-005, ALG-LIST-006, ALG-LIST-011, ALG-LIST-014, ALG-LIST-015, ALG-LIST-016, ALG-MATRIX-002, ALG-PREFIX-008, ALG-QUEUE-001, ALG-QUEUE-002, ALG-SELECT-003, ALG-SELECT-004, ALG-SORT-001, ALG-SORT-002, ALG-SORT-003, ALG-SORT-004, ALG-SORT-005, ALG-STREAM-003, ALG-TEXT-001, ALG-TEXT-002, ALG-TEXT-004, ALG-TEXT-005, ALG-TEXT-006, ALG-TEXT-010, ALG-TREE-001, ALG-TREE-002, ALG-TREE-003, ALG-TREE-004, ALG-TREE-005, ALG-TREE-009, ALG-TREE-010, ALG-TREE-012, ALG-TREE-016, ALG-TREE-017, ALG-TREE-018, ALG-TREE-022, ALG-TRIE-007, ALG-TWOPTR-004, ALG-TWOPTR-006, ALG-UNION-001, ALG-UNION-002, ALG-UNION-003, ALG-UNION-004, ALG-UNION-005, ALG-UNION-006, ALG-UNION-007, ALG-UNION-008, ALG-WINDOW-012, ALG-WINDOW-013.

Citations as written in the doc comments:
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, Appendix B.4 (directed and undirected graphs)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, Appendix B.5.3 (height of a rooted tree: edges on the longest simple downward path to a leaf)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, Appendix D.1 (the transpose, (Aᵀ)ᵢⱼ = aⱼᵢ)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, Exercise 10.1-6 (a queue from two stacks)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, Exercise 10.2-7 (Θ(n) nonrecursive reversal in constant storage)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, Exercise 16.2-2 (an O(nW) dynamic program for the 0-1 knapsack problem)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, Exercise 2.3-7 (the pair-sum problem) and ch. 11 (hash tables)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, Exercise 22.4-2 (the number of paths between two vertices of a DAG in linear time: a vertex's count is the sum over its predecessors, in topological order; here row-major order)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, Problem 16-1(d) (an O(nk) dynamic program for coin changing with k denominations; CLRS assumes a 1-cent coin, which the recurrence here does not need)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 11 "Hash Tables" (O(1) expected search and insert)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 11 "Hash Tables" (O(1) expected search, insert, delete)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 11 (counting in a hash table)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 11 (dictionary INSERT and SEARCH)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 11 (dictionary INSERT, SEARCH, DELETE)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 11 (dictionary operations)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 11 (expected O(1) SEARCH in a hash table)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 11 (grouping by key in a hash table)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 11 (hash tables; Theorem 11.2)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 12 (binary search trees)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 13 (red-black trees: the O(log k) Add that lost, see the catalog)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 9 introduction (lower and upper median)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, ch. 9 introduction (order statistics) and §9.2–9.3
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §10.1 (ENQUEUE and DEQUEUE with wrap-around)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §10.2 (LIST-INSERT, LIST-DELETE with a sentinel, Figure 10.4)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §10.2 (doubly linked list, O(1) delete) and ch. 11 (hash tables)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §10.2 (singly linked lists)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §12.1 (T.root)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §12.1 (an inorder walk of a BST lists its keys in sorted order)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §12.1 (binary-search-tree property; Theorem 12.1, the inorder walk of a BST is sorted)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §12.1, Theorem 12.1 (inorder walk of a BST yields its keys in sorted order)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §12.2 (ITERATIVE-TREE-SEARCH)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §12.2 (search path; the LCA is the first node on it with p <= key <= q, proved in doc/algorithms_core_2026-11-12/catalog/tree.md)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §12.3 (TRANSPLANT, TREE-DELETE)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §12.3 (TREE-INSERT)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §15.1 (the same recurrence with one item per length: rod cutting)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §15.1, eq. (15.2) r_n = max_{1≤i≤n}(p_i + r_{n−i}), BOTTOM-UP-CUT-ROD and EXTENDED-BOTTOM-UP-CUT-ROD, Figure 15.1 (the price table)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §15.2, eq. (15.7), MATRIX-CHAIN-ORDER and PRINT-OPTIMAL-PARENS, Figure 15.5 (the six-matrix example)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §15.4, eq. (15.9) (the LCS recurrence)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §15.4, eq. (15.9) (the LCS recurrence) and Exercise 15.4-4 (one row of min(m, n) entries)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §16.1 (activity selection; Theorem 16.1, GREEDY-ACTIVITY-SELECTOR)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §16.1, Exercise 16.1-4 (lecture halls; interval-graph coloring)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §2.1 (INSERTION-SORT)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §2.3.1 (MERGE)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §2.3.1, procedure MERGE
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §21.1 (CONNECTED-COMPONENTS with disjoint sets)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §21.1 (SAME-COMPONENT)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §21.1 (connected components of a graph whose edges arrive dynamically, maintained with disjoint sets)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §21.1 (each UNION reduces the number of sets by one)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §21.1 (n = the number of MAKE-SET operations)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §21.3 (FIND-SET with path compression), §21.4 (the α bound)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §21.3 (MAKE-SET for each element)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §21.3 (MAKE-SET)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §21.3 (UNION, LINK with union by rank)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.1
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.1 (Adj[u])
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.1 (adjacency lists; the copy duplicates every list)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.1 (adjacency-list representation)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.1 (adjacency-list representation; an undirected edge appears in both lists)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.1 (|E|)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.1 (|V|)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.2 (BFS from each undiscovered vertex counts the connected components of the implicit grid graph)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.2 (BFS on the implicit grid graph)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.2 (BFS shortest paths)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.2 (BFS)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.2 (BFS, here on the implicit grid graph)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.2 (BFS; seeding the queue with every border cell is a BFS from one virtual source joined to them)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.2 (Theorem 22.5: BFS finds shortest-path distances; PRINT-PATH)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.2 (breadth-first search by distance from the root)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.2 (breadth-first search gives fewest-edge paths; on this implicit graph each BFS layer is a contiguous run of indices, so a layer is two numbers and the search is one pass; proof in catalog/dp.md)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.2 (breadth-first search: vertices in order of distance from the source)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.2 (breadth-first search; on the implicit graph i → i+1 … i+maxJump[i] the reachable set from 0 is a prefix [0, f], so reachability reduces to the farthest index f; proof in catalog/dp.md)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.2, Exercise 22.2-7 (split into two kinds with every rivalry across, by BFS, O(n + r))
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.3 (DFS-VISIT)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.3 (DFS-VISIT, here on the implicit grid graph)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §22.4 (topological sort; Exercise 22.4-5 gives this in-degree method)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §23.2 (MST-KRUSKAL)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §23.2 (MST-PRIM)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §24.1 (BELLMAN-FORD with the negative-cycle check)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §24.2 (shortest paths in a directed acyclic graph by relaxing in topological order; row-major order is topological for right and down moves, with each cell's cost on its incoming edges)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §24.2 (shortest paths in a directed acyclic graph in reverse topological order: each entry's best path to the base is its value plus the better of its two children's)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §24.3 (DIJKSTRA; binary-heap bound)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §25.2 (FLOYD-WARSHALL)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §31.6 "Powers of an element" (repeated squaring)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §32.2 (Horner's rule, the rolling update, Figure 32.5)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §32.2 (rolling fingerprints) and ch. 11 (hash tables)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §32.2, RABIN-KARP-MATCHER and Figure 32.5
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §32.4, COMPUTE-PREFIX-FUNCTION and Figure 32.10
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §32.4, KMP-MATCHER
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §33.1 (SEGMENTS-INTERSECT, ON-SEGMENT)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §33.1 (cross products, DIRECTION)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §33.3 (convex hull; lower bound)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §33.4 (divide and conquer, presorting, the strip)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §6.1 (binary heap), §6.5 (priority queues)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §6.1 (heap-size)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §6.3 (BUILD-MAX-HEAP, O(n) bound)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §6.5 (HEAP-EXTRACT-MAX)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §6.5 (HEAP-MAXIMUM)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §6.5 (MAX-HEAP-INSERT)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §6.5 (a bounded priority queue keeps the k best of a stream)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §6.5, Exercise 6.5-9 (k-way merge of the rows with a min-heap, stopped at the k-th output)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §6.5, Exercise 6.5-9 (k-way merge with a min-heap; row i is the ascending list a[i]+b[0], a[i]+b[1], …)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §7.1 (PARTITION around a pivot; the linked form keeps order in O(1) space)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §7.1, procedure PARTITION (Lomuto's scheme: one pass, a boundary index, swap each kept element to the boundary; here the test is "not zero" and there is no pivot)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §8.2 (COUNTING-SORT)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §8.2 (counting occurrences of each key) and ch. 11 (hash tables)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §8.2 (counting) and ch. 11 (hash tables)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §8.3 (RADIX-SORT, Lemma 8.4: b-bit keys in r-bit digits, Θ((b/r)(n + 2^r)))
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §8.4 (BUCKET-SORT) and Exercise 8.4-2 (worst case)
- Cormen, Leiserson, Rivest, Stein, Introduction to Algorithms, 3rd ed., MIT Press, 2009, §9.2–9.3

Checked: see the builder's catalog section in doc/algorithms_core_2026-11-12/catalog/ (confirmation status per citation; not first-hand items are listed in PROGRESS.md and REPORT.md).

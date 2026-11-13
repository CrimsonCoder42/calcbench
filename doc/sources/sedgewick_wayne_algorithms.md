# Sedgewick & Wayne, Algorithms

**Citation**: Sedgewick, Wayne, *Algorithms*, 4th ed., Addison-Wesley, 2011, §2.4 "Priority Queues", client TopM; https://algs4.cs.princeton.edu/24pq/ read 2026-11-12.

Taken: keep the M largest with a min-priority-queue of size at most M+1 and print them in descending order (quoted from the TopM javadoc). Used by ALG-HEAP-001. Checked: oracle against sort-descending-and-take.


## Used by S10 (algorithms core, 2026-11-12)

IDs: ALG-GRAPH-017, ALG-HASH-001, ALG-HASH-002, ALG-HASH-003, ALG-HASH-004, ALG-HEAP-010, ALG-HEAP-011, ALG-HEAP-012, ALG-HEAP-013, ALG-SORT-001, ALG-SORT-002, ALG-SORT-007, ALG-SORT-009, ALG-SORT-010, ALG-STACK-009, ALG-STACK-010, ALG-TWOPTR-001, ALG-TWOPTR-002, ALG-TWOPTR-003, ALG-WINDOW-015.

Citations as written in the doc comments:
- Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §1.3 (Dijkstra's two-stack evaluator; Exercise 1.3.11 EvaluatePostfix)
- Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §1.3, Exercise 1.3.4 (Parentheses)
- Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §1.4, Exercise 1.4.15 "Faster 3-sum" (a linear two-sum on a sorted array)
- Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §1.4, Exercise 1.4.15 "Faster 3-sum" (a quadratic 3-sum from the linear two-sum on a sorted array)
- Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §2.2 (bottom-up mergesort; between 1/2 N lg N and N lg N compares)
- Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §2.3 (3-way partitioning: a[lo..lt-1] < v, a[lt..i-1] = v, a[gt+1..hi] > v)
- Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §2.3 (insertion sort for tiny arrays, cutoff 5 to 15)
- Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §2.4 (client TopM: a min-heap of the M largest)
- Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §2.4 Creative Problems, "Dynamic median-finding" (max-heap below the median, min-heap above)
- Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §2.4, Exercise 2.4.30 "Dynamic median-finding"
- Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §2.4, Exercise 2.4.30 "Dynamic median-finding" (a max-heap below the median, a min-heap above; algs4.cs.princeton.edu/24pq, read 2026-11-12)
- Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §2.4, Exercise 2.4.30 "Dynamic median-finding" (find the median in constant time)
- Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §2.4, Exercise 2.4.30 "Dynamic median-finding" (insert in logarithmic time)
- Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §2.5 "Applications", "Duplicates" (sort, then one pass noting equal values that are now adjacent), https://algs4.cs.princeton.edu/25applications/ (read 2026-11-12)
- Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §2.5 "Duplicates" (sort, then one pass noting equal values that appear consecutively)
- Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §2.5 (pointer sorting: sort references, not the data) and §2.2 (mergesort, stable)
- Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §3.1 (client FrequencyCounter) and §3.4
- Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §3.4
- Sedgewick, Wayne, Algorithms, 4th ed., Addison-Wesley, 2011, §4.1 (two-colorability: is a graph bipartite?)

Checked: see the builder's catalog section in doc/algorithms_core_2026-11-12/catalog/ (confirmation status per citation; not first-hand items are listed in PROGRESS.md and REPORT.md).

# S1 report: reference algorithms (2026-11-12)

**Built** (ID | name | file:line | Time | Space | source). Full citations are in INDEX.md and doc/sources/.
| ALG-SEARCH-001…004 | LowerBound, UpperBound, BinarySearchFirst, BinarySearchLast | algo/search/lower_bound.go:15, upper_bound.go:15, binary_search_first.go:16, binary_search_last.go:16 | O(log n) | O(1), 0 alloc | Knuth Vol. 3 §6.2.1; Bentley Col. 9.3 |
| ALG-SELECT-001/002 | KthSmallest, Median | algo/selection/kth_smallest.go:31, median.go:21 | O(n) worst | O(n) copy + O(log n) stack | Musser 1997; Blum et al. 1973; Hoare 1961; CLRS ch. 9 |
| ALG-WINDOW-001/002 | RollingSum, RollingMean | algo/window/rolling_sum.go:36, rolling_mean.go:18 | O(n) | O(1) extra | NIST e-Handbook §6.4.2.1; Neumaier 1974; Higham §4.3 |
| ALG-WINDOW-003/004 | RollingMin, RollingMax | algo/window/rolling_min.go:16, rolling_max.go:16 | O(n) | O(w), 0 alloc for w ≤ 64 | Lemire 2006 |
| ALG-HEAP-001 | TopK | algo/heap/top_k.go:24 | O(n log k) | O(1) extra + O(k) result | Sedgewick & Wayne §2.4; Floyd 1964; Williams 1964 |
| ALG-MERGE-001 | MergeKSorted | algo/merge/merge_k_sorted.go:17 | O(n log k + m) | O(k) extra | CLRS §6.5, Ex. 6.5-9; Knuth §5.4.1; Floyd 1964 |
| ALG-INTVL-001 | MergeOverlapping | algo/interval/merge_overlapping.go:31 | O(n log n) | O(log n) extra | Klee 1977; Fredman & Weide 1978 |

**go test** (final, out/go_test_final.txt): heap ok · interval ok · merge ok · search ok · selection ok · window ok. 86 test functions; the reviewer measured 100% statement coverage.

**Benchmarks** (ns/op at 1e3 / 1e4 / 1e5 · B/op · allocs; Apple M-series, out/bench_1.txt, bench_selection_3.txt, bench_window_2.txt):
| LowerBound 33 / 49 / 75 · 0 · 0 | log ✓ | UpperBound 15 / 24 / 54 · 0 · 0 | log ✓ |
| BinarySearchFirst 32 / 47 / 76 · 0 · 0 | log ✓ | BinarySearchLast 24 / 40 / 57 · 0 · 0 | log ✓ |
| KthSmallest 8.9k / 103k / 988k · 8n · 1 | linear ✓ | Median 9.4k / 106k / 1.04M · 8n · 1 | linear ✓ |
| RollingSum 5.1k / 61k / 797k · 0 · 0 | linear ✓ (×13 at 1e5: memory) | RollingMean 5.1k / 59k / 799k · 0 · 0 | linear ✓ |
| RollingMin 7.1k / 117k / 1.28M · 0 · 0 | linear ✓ | RollingMax 6.9k / 113k / 1.22M · 0 · 0 | linear ✓ |
| TopK k=100 random 4.4k / 11k / 61k · 896 · 1 | sub-linear ✓ | TopK ascending 17k / 277k / 2.9M · 896 · 1 | n log k ✓ |
| MergeKSorted k=16 23k / 348k / 3.57M · 8n · 2 | linear at fixed k ✓ | MergeOverlapping 49k / 808k / 9.96M · 16n · 1 | n log n ✓ |

**Probes** (moved to _delete/; raw output in out/):
- Selection fallback (out/probe_fallback.txt, 200 inputs per n): "halve every 2 rounds" switched to median
  of medians on 59–80% of random inputs. A 6n work budget switched on 0 of 600. Median-of-3 averages 2.75n.
  Median of medians alone takes 2.02 ms at 1e5, against 0.99 ms for introselect. Fixed, A-010.
- TopK heap against batched selection at n = 1e5 (out/probe_topk_batch.txt): random input, k = 10: 51 µs
  vs 718 µs. Ascending input, k = 10⁴: 5.28 ms vs 0.32 ms. Neither dominates, so the heap is kept (A-011).

**Reviewer** (fresh context, out/bench_reviewer.txt), with the fix for each finding:
1. RollingSum and RollingMean leaked rounding from values that had left the window: an error of 9e-4 at
   n = 1e6, and a wrong sign in a 4-window repro. Fixed: the sum restarts every w windows. The tests'
   tolerance is now local to the values since the restart, with the repro and a 20,000-value history test.
2. A false ErrOverflow fired when a partial window sum overflowed. Now documented in both functions.
3–7. Introselect's worst case, the deque ring, merge stability, the worked examples and no-panic were all
   confirmed. Doc drift fixed: SCOPE count, spec A-range, Floyd's title, "at most" comparisons, Floyd now
   cited by MergeKSorted, INDEX "Does" now taken from doc.go, this REPORT written.

**Decisions logged**: D-001…D-004 (Devin). Assumptions A-001…A-011.
**FOUND, NOT FIXED**: none in code. The Knuth §6.2.1 "653" example and Bentley Col. 9.3 were not checked
against a physical copy.
**MOVED TO _delete/**: algo/selection/probe_fallback_test.go, algo/heap/probe_batch_select_test.go (probes);
INDEX.md.new (scratch).
**LEFT FOR DEVIN**: TODO.md: the book page checks, the TopK worst-case hybrid (new IDs?), integer rolling sums.

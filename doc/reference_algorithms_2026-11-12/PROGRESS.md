# Progress

- [x] Orient: rules, taxonomy and skills read; repo had no code and no docs skeleton.
- [x] Decisions asked once; Devin picked the recommendation on all four (D-001…D-004).
- [x] Step 1: skeleton docs, INDEX.md header, work folder, `.gitignore` rule for `doc/*/out/`.
- [x] Step 2: spec `doc/specs/s1_reference_algorithms.md`. The NIST §6.4.2.1 example and the algs4 §2.4 TopM
      client were checked on the web.
- [x] Step 3: build of 6 packages and 13 functions. Each file passed gofmt, build and vet. The TopK k = 0
      panic was found while re-reading and fixed before any test ran.
- [x] Step 4a: `go test ./...`, all 6 ok (`out/go_test_1.txt`). Benchmarks in `out/bench_1.txt`.
- [x] Step 4b: KthSmallest grew 22x from 1e3 to 1e4. A probe (`out/probe_fallback.txt`) showed the
      "halve every 2 rounds" fallback firing on 59–80% of random inputs. It was replaced by a 6n work
      budget (A-010). The selection benchmarks now rotate 8 inputs, because one repeated 1e3 input let the
      branch predictor learn it (`out/bench_selection_3.txt`).
- [x] Step 4c: TopK heap against batched selection, probe in `out/probe_topk_batch.txt`. The heap is kept
      (A-011); both probes were moved to `_delete/`.
- [x] Step 4d: the rules were tightened mid-run (full citations, DOIs). Source lines were rewritten and the
      DOIs checked on Crossref (the Williams 1964 DOI guess was wrong and was corrected). INDEX.md was
      regenerated from the doc comments. `doc/sources/` has one file per source.
- [x] Step 4e: the reviewer found 1 real defect (window sums leaking history), fixed and re-tested; doc drift fixed; final go test all ok; ledgers, cleanup counts (13|13|13|13) and REPORT.md done.

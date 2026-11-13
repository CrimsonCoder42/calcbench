# Progress

- [x] Orient: CLAUDE.md ("Several subjects at once"), go-functions / go-testing / go-complexity / numeric /
      sources rules, INDEX.md, geometry.md (Point, distance, polar, arcs) and trig.md (angle normalisation) read.
- [x] Research: DLMF §1.5 read; Kahan "Mindless", Diebel 2006, Shoemake 1985 (OCR), Ogita–Rump–Oishi 2005 and Strang
      §10.6 read from PDFs; DOIs on Crossref; ISO 80000-2 unreadable (403), not cited.
- [x] Build: 115 functions in 26 files (vectors 001–051, coordinates 052–055, matrices 056–100, quaternions
      101–115), doc.go generated from the doc comments; gofmt and vet clean throughout.
- [x] Tests: worked examples, error tables (every documented error), oracle identities with fixed seeds,
      AllocsPerRun = 0 for all 86 fixed-size functions and for n-D functions with a reused dst; 115 benchmarks.
- [x] Verify 1: `go test ./num/linalg/...` FAILED (out/go_test_1.txt): FMA-fusion false overflow, Dot2 fusion,
      test tolerance below one ulp, slerp endpoint expectation, repeated-eigenvalue accuracy. A false
      ErrSingular for badly scaled inverses and a subnormal scaling overflow were found while writing the tests.
- [x] Fixes (README "Defects"), then verify 2: `go test ./num/linalg/...` ok (out/go_test_2.txt).
- [x] Benchmarks: out/bench_1.txt; growth matches the documented complexity; allocations as documented.
- [x] Catalog doc/catalog/linalg.md (136 rows, 115 entries), PENDING_SHARED.md, 13 new doc/sources files.
- [x] Reviewer: 2 HIGH, 3 MEDIUM, 7 LOW (README table); all fixed with regression tests; go test ok
      (out/go_test_4.txt); benchmarks re-run (out/bench_2.txt), fixed-size paths still 0 allocs.
- [x] Catalog, doc.go and INDEX rows regenerated from the doc comments after the fixes.
- [x] Cleanup: counts 115 | 115 (pending) | 115; scratch zz_debug_test.go moved to _delete/; final go test ok (out/go_test_final.txt).

# S21 report: toolbench requests R-4, R-5 (2026-11-13)

Built: ALG-RATELM-005 TokenBucket.Restore `algo/ratelimit/token_bucket_restore.go:20` · ALG-LIST-018
LRUCache.Oldest `algo/list/lru_cache_oldest.go:14`. Newest not built (A-073).

Tests: `go test ./... -count=1` 91 ok, 0 FAIL (out/test.txt); `-race` 91 ok, 0 FAIL (out/race.txt). After the
benchmark fix: algo/ratelimit and testing rerun with -race, ok (out/race_rerun.txt).

Benchmarks (-benchtime 1s -count 1, all 0 B/op 0 allocs/op): Oldest 2.19 ns (out/bench.txt) · Restore
25.2 µs per 1e3 Allow+Restore pairs, linear to 1e5 (out/bench_rerun.txt).

Index: 1663 → 1665 live functions (1697 → 1699 rows, 34 REMOVED); INDEX == code == doc.go (out/index.txt).

Reviewer: 2 low findings. The Restore benchmark only ran the full-bucket path: fixed by B1 (take 2, restore 1).
RFC 2697 cited for the cap only: the lead opened §3 ("The maximum size of the token bucket C is CBS"; "If Tc is
less than CBS, Tc is incremented by one"), and the citation stands. x/time/rate CancelAt was read in module source v0.16.0.

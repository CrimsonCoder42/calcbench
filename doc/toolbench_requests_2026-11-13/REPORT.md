# S20 report: toolbench requests R-1…R-3 (2026-11-13)

Built: ALG-BACKOF-001 Exponential `algo/backoff/exponential.go:37` · -002 FullJitter `algo/backoff/full_jitter.go:23` ·
-003 EqualJitter `algo/backoff/equal_jitter.go:20` · -004 DecorrelatedJitter `algo/backoff/decorrelated_jitter.go:24` ·
NUM-INTEG-013 CeilDiv `num/integer/ceil_div.go:22` · ALG-CRON-001 ParseCron `algo/cron/parse_cron.go:62` ·
ALG-CRON-002 CronSpec.Next `algo/cron/next.go:38`.

Tests: `go test ./... -count=1` 91 ok, 0 FAIL (out/test.txt); `-race` 91 ok, 0 FAIL (out/race.txt). After the review's
comment-only fixes: backoff, cron, integer and testing reran ok (out/test_rerun.txt, out/race_rerun.txt). Fuzz 30 s
each, all pass (out/fuzz_*.txt).

Benchmarks (out/bench.txt, -benchtime 1s -count 1, all 0 B/op 0 allocs/op): Exponential 4.2 ns · FullJitter 5.6 ·
EqualJitter 5.8 · DecorrelatedJitter 7.6 · CeilDiv 2.1 · ParseCron 29 (stars) / 84 (lists) · Next 50 (every minute) /
112 (weekdays) / 202 (Feb 29, 2097→2104) / 177 (NY gap) / 128 (NY overlap).

Index: 1656 → 1663 live functions (1690 → 1697 rows, 34 REMOVED); INDEX == code == doc.go (testing passes).

Reviewer: 5 findings, all fixed (cron day rule described as POSIX → A-068 plus doc wording; worked numbers attributed
to Brooker; dates read on web citations; Concrete Mathematics locator → §3.1 eq. (3.4), Ex. 3.12; Go spec quoted).
Web citations were opened by the lead (POSIX crontab, Brooker, gRPC); Brooker's formulas are images, so they are cited by figure.
The reviewer brute-forced cron against 5598 specs in 14 zones near real transitions: 0 mismatches.

FOUND, NOT FIXED: gofmt/parse errors in `_delete/algo/reduce/errors.go`, `_delete/num/prob/weibull_test_helpers_test.go`,
`doc/calculator_net_check_2026-11-13/math/checks/scratch/confidence-interval-calculator.go`.

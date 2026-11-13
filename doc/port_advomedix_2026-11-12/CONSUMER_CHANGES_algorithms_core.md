# Consumer changes (lists only; the consumers were read, never edited)

Each row: a hand-rolled copy found in a consumer, file:line, and the ALG ID that replaces it. Paths are relative
to /Users/devinanderson/Desktop.

| # | Consumer | file:line | Hand-rolled code | Replace with | Notes |
|---|---|---|---|---|---|
| 1 | backendAdvo | advo_production/advomedix-v2/backendAdvo/knowledgebase/series/rank.go:19-37 | countHeap + TopK over map counts | ALG-HASH-003 TopKByCount | field `Key` → `Value`; k < 0 now errors (catalog/hashing.md) |
| 2 | backendAdvo | knowledgebase/series/rank.go:70 | MergeIntervals with a gap (half-open, time.Time) | ALG-INTVL-001 MergeOverlapping | convert to int64 ns, add gap to End, merge, subtract gap (catalog/interval.md) |
| 3 | backendAdvo | api/handlers/patient/feed.go:725-732 | sortItems: sort.Slice with an ID tiebreak because sort is not stable | ALG-SORT-001 MergeSortStable or ALG-SORT-009 SortIndices | only matters if input order is meaningful; the tiebreak already fixes the order |
| 4 | backendAdvo | api/middleware/rate_limit.go:273-288 | per-IP golang.org/x/time/rate limiter | ALG-RATELM-001 TokenBucket (per-key arithmetic only) | the map, sweep and locking stay theirs; TokenBucket does not lock: add a mutex per entry |
| 5 | backendAdvo | knowledgebase/series/moments.go:119-130 | float prefix sums and range sums | ALG-PREFIX-001 PrefixSums / -002 RangeSum (integer types) | our prefix sums are integer-only; a float pattern stays theirs unless a float variant is added (LATER) |
| 6 | red_lite_cache | red_lite_cache/internal/adminauth/limiter.go:52-89 | map of token buckets via x/time/rate | ALG-RATELM-001 TokenBucket (per-key arithmetic) | as row 4 |
| 7 | red_lite_cache | internal/cache/ordered.go:45 | sort.Search for the first key greater than k | ALG-SEARCH-002 UpperBound (S1) | |
| 8 | red_lite_cache | internal/cache/ordered.go:64, 118, 201 | sort.SearchStrings | ALG-SEARCH-001 LowerBound (S1) | |

Not replaced (reasons in the catalog): red_lite_cache ordered.go:209-231 appendPrefix and sheet/sweeplog.go:17-24
covers() (trie functions do not fit a B+tree leaf walk or a log of prefix/suffix ops); cache/wheel.go:17 expiry
wheel (no timing-wheel function in this slice); backendAdvo series days.go:62 BucketFixed, glucose.go:11
TimeInRange, fitness.go:269 RollingUsualRange; no hand-rolled LRU exists in red_lite_cache (its store does not
evict).

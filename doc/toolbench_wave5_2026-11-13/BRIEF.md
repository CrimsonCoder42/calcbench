# S21 brief: toolbench requests R-4, R-5 (2026-11-13)

Request: `toolbench/doc/wave5_2026-11-13/calcbench_requests.md` (read-only). Decision D-088. Additive only: no
released signature or behaviour changes.

## B1 · algo/ratelimit · ALG-RATELM-005 · `func (b *TokenBucket) Restore(n int64, now time.Time)`
- Gives back n tokens taken earlier: refill to now first (same clock rule as Allow: a now before the latest
  seen counts as the latest), then level += n·intervalNs, capped at capacity (burst tokens); never negative.
- n <= 0: no-op (nothing changes, not even the refill). Zero value (burst 0): no-op. n >= burst: level =
  capacity directly, so n·intervalNs never overflows.
- A bucket never used yet starts full at the refill, so Restore leaves it full.
- Restore does not undo the passage of time and does not know which tokens a later Allow took: the caller
  passes how many to return (x/time/rate Reservation.Cancel computes that count; toolbench does the same).
- Sources: RFC 2697 §3 (bucket C capped at CBS; Tc incremented by one up to CBS); The Go Authors,
  golang.org/x/time/rate, `func (r *Reservation) CancelAt` (restores tokens not consumed by later
  reservations; the bucket is capped at burst), https://pkg.go.dev/golang.org/x/time/rate (read 2026-11-13).
- Example (oracle): burst 5, take 3 at t0 → 2 left; Restore(3, t0) → 5; Restore(10, t0) → 5 (capped).
- Tests: the example; capped restore; restore after a later take (take 3, take 2 → 0, Restore(3) → 3);
  n = 0 and n < 0 no-op; zero value; unused bucket; huge n (math.MaxInt64) no overflow; restore then refill
  still capped; now earlier than latest; Allow semantics unchanged after Restore; 0 allocs (AllocsPerRun);
  BenchmarkTokenBucketRestore. Existing benchmarks untouched.

## B2 · algo/list · ALG-LIST-018 · `func (c *LRUCache[K, V]) Oldest() (key K, value V, ok bool)`
- The least recently used entry (the one Put would evict next), not removed, recency unchanged (not a use);
  ok false on an empty cache, the zero value or nil. O(1): the list's Back.
- Newest() is not added (DECIDED D2: only if it costs one line; a method with its own ID, doc and tests does not).
- Source: as ALG-LIST-015 (Mattson et al. 1970; CLRS §10.2).
- Example (oracle): Put a, b, c; Get a → Oldest() = b; Oldest() again = b (unchanged).
- Tests: the example; empty, zero value, nil; after Get; after Put of an existing key (Put a,b,c; Put a again →
  Oldest = b); after Peek (unchanged); after Delete of the oldest; Oldest equals the next Put's eviction;
  capacity 1; 0 allocs; BenchmarkLRUCacheOldest. Existing benchmarks untouched.

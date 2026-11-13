# algo/ratelimit (SHORT = RATELM)

Builder W4-B14. Not yet in `doc/taxonomy.md` (lead to add: `ALG | algo/ratelimit | RATELM | rate meters with
the clock passed in: token bucket, GCRA leaky bucket, fixed (jumping) window, sliding (moving) window log`).

| ID | Name | Signature |
|---|---|---|
| ALG-RATELM-001 | TokenBucket | `func NewTokenBucket(tokensPerInterval int64, interval time.Duration, burst int64) (*TokenBucket, error)` with `Allow(now time.Time, n int64) (ok bool, retryAfter time.Duration)`, `Available(now time.Time) int64` |
| ALG-RATELM-002 | LeakyBucket | `func NewLeakyBucket(emissionInterval, tolerance time.Duration) (*LeakyBucket, error)` with `Allow(now time.Time) (ok bool, retryAfter time.Duration)` |
| ALG-RATELM-003 | FixedWindow | `func NewFixedWindow(limit int64, window time.Duration) (*FixedWindow, error)` with `Allow(now time.Time, n int64) (ok bool, retryAfter time.Duration)` |
| ALG-RATELM-004 | SlidingWindowLog | `func NewSlidingWindowLog(limit int64, window time.Duration) (*SlidingWindowLog, error)` with `Allow(now time.Time, n int64) (ok bool, retryAfter time.Duration)` |
| — | SlidingWindowCounter | **HELD (no source)**: the "weighted previous window" estimate is described only by vendor blogs and library docs (searched 2026-11-12); no paper, RFC or standard found. No ID given. |

## Conventions shared by the package (builder decisions, for the lead's table)

- **Clock as a parameter.** Every call takes `now time.Time`; no type reads the clock, starts a goroutine or
  holds a lock. Not safe for concurrent use: the caller's mutex guards a meter (one meter per key, the
  caller's map). Times should come from one clock (`time.Now` readings carry the monotonic clock, which
  `Sub`/`Before` use when both operands have it).
- **Clock going backwards** (decision R1): a `now` earlier than the latest `now` a meter has seen is treated
  as that latest time: no elapsed time, no refill, no window change; the call is still decided. **Lost**:
  returning an error (a wall clock steps back under NTP; the caller can do nothing better than "no time
  passed", and an error on the hot path would be ignored or turned into an outage).
- **Time arithmetic**: only `Time.Sub` (saturates at ±292 years) and `Time.Add` of a non-negative duration
  never larger than an already-measured gap, or of the emission interval with a check that the result moved
  forward. Scaled integer products are checked with `math/bits.Mul64` at construction.
- **Requests that can never conform** (`n` greater than the burst or the limit): `(false, Never)` where
  `Never = time.Duration(math.MaxInt64)`. `n <= 0` is `(true, 0)` and changes nothing. A zero-value meter
  refuses everything with `Never` (its burst/limit is 0). **Lost**: a third error result (the signature the
  lead asked for is `(ok, retryAfter)`; the two documented values cover both cases without one).
- **Size cap** (L13/L26): SlidingWindowLog keeps at most `limit` entries of 32 bytes, allocated once; limit is
  capped at 1<<22 (128 MiB) and larger returns `ErrTooLarge`.

### ALG-RATELM-001 TokenBucket
- **Source**: Heinanen and Guérin, "A Single Rate Three Color Marker", RFC 2697, IETF, September 1999, §3
  "Metering" (read 2026-11-12, https://www.rfc-editor.org/rfc/rfc2697.txt): the committed bucket C starts
  full, Tc(0) = CBS; Tc is incremented CIR times per second up to CBS; a packet of size B conforms (green) when
  Tc − B ≥ 0 and then Tc is decremented by B, otherwise nothing is decremented. Only the C bucket is used
  (no excess bucket E: EBS = 0 makes the yellow branch empty). Turner, "New directions in communications (or
  which way to the information age?)", IEEE Communications Magazine 24(10):8–15, 1986,
  doi:10.1109/MCOM.1986.1092946 is the historical origin of the term (Crossref metadata confirmed; full text not
  read; not used for any rule).
- **Best known**: O(1) time and O(1) space per call; this is the lower bound (each call must read the state).
- **Approach**: exact integer fluid model. One token = `interval` nanoseconds of "credit units"; each elapsed
  nanosecond adds `tokensPerInterval` units, so the level is exact at every nanosecond and never drifts
  (capacity = burst·interval units, checked with `bits.Mul64`). Refill is clamped to the capacity before
  multiplying, so `elapsed·tokensPerInterval` never overflows. `retryAfter` = ⌈deficit / tokensPerInterval⌉ ns.
  This credits whole tokens at the same average rate as RFC 2697's "CIR times per second", in finer steps.
  **Lost**: float64 tokens (as in golang.org/x/time/rate; rounding accumulates and equality tests become
  tolerances); RFC 2697's discrete one-token ticks (needs a tick phase the caller would have to choose).
- **Equivalence used as an oracle**: with `tokensPerInterval = 1`, `interval = T`, `burst = b`, TokenBucket
  makes the same decisions and returns the same retryAfter as LeakyBucket(T, (b−1)·T), by I.371 Annex A.2:
  tokens = (b·T − max(X′, 0))/T where X′ is the continuous-state bucket content. Proof in the test comment.
- **Worked example (hand-derived from RFC 2697 §3)**: 2 tokens per second, burst 3, at t=0: Allow(n=3) ok;
  Allow(1) refused, retryAfter 500 ms; at t=500 ms Allow(1) ok; at t=10 s Available = 3 (capped at burst).
- **Errors**: `ErrRate` (tokensPerInterval < 1), `ErrInterval` (interval ≤ 0), `ErrBurst` (burst < 1),
  `ErrTooLarge` (burst·interval > MaxInt64). Input never modified.
- **Consumer evidence**: red_lite_cache `internal/adminauth/limiter.go:52-89` (a map of
  golang.org/x/time/rate limiters with a sweep and a cap) and the consumer backend
  `api/middleware/rate_limit.go:273-288` (`checkIPLimit`: `rate.NewLimiter(RequestsPerSecond, BurstSize)` in a
  `sync.Map`, `Allow()`). **What TokenBucket replaces**: the per-key bucket arithmetic only, i.e. the
  `*rate.Limiter` value and its `Allow`/`ReserveN` call (limiter.go `bucket.l`; rate_limit.go:279/287).
  Their call becomes `ok, retry := b.Allow(now, 1)` and `retry` feeds Retry-After directly (limiter.go
  computes it from a Reservation). **What stays theirs**: the key → bucket map, the idle sweep and
  MaxEntries cap, the clock (`now func() time.Time` / `time.Now()`), and the locking. Note for the consumer backend:
  `rate.Limiter` locks internally and is used from a `sync.Map` without an outer lock; TokenBucket does not
  lock, so the entry needs its own `sync.Mutex` around `Allow` (or a striped lock). Removes the
  golang.org/x/time/rate dependency from both. Consumers are listed, never edited.
- **Tags**: rate-limiting, token-bucket, traffic-policing, throttling, retry-after. Pattern pointers only:
  code-heim go_82 (not a source).

### ALG-RATELM-002 LeakyBucket
- **Source**: ITU-T Recommendation I.371 (03/2004), "Traffic control and congestion control in B-ISDN",
  Annex A "Generic Cell Rate Algorithm GCRA(T,τ)", §A.1 Virtual Scheduling Algorithm and §A.2
  Continuous-state leaky bucket algorithm, pp. 86–87 (read 2026-11-12 from the ITU free PDF): TAT is
  initialised to ta(1); a cell with ta(k) < TAT − τ is non-conforming and TAT is unchanged; if
  TAT − τ ≤ ta(k) < TAT, TAT += T; if ta(k) > TAT, TAT = ta(k) + T. The NOTE states A.1 and A.2 decide the
  same cells. Second source: Berger, Eckberg, Hou, Lucantoni, "Performance characterizations of traffic
  monitoring, and associated control, mechanisms for broadband 'packet' networks", Proc. IEEE GLOBECOM '90,
  pp. 350–354, 1990, doi:10.1109/GLOCOM.1990.116535, §4 p. 351 (the leaky-bucket monitor: a counter
  incremented per arrival and drained at a constant rate, excess over a threshold non-conforming; read).
- **Best known**: O(1) time, O(1) space per call (lower bound).
- **Approach**: the virtual scheduling form (A.1): one `time.Time` (TAT) of state, comparisons by
  `TAT.Sub(now)` (saturating). A cell conforms iff TAT − now ≤ τ; retryAfter = (TAT − now) − τ. The next TAT
  is checked to have moved forward (Time.Add saturates at the end of time; then `(false, Never)`).
  **Lost**: the continuous-state form (A.2) with X and LCT (two fields for the same decisions; kept as the
  test oracle); a queue-based "leaky bucket as a queue" (that shapes/delays traffic, needs a scheduler or
  goroutine; this package meters).
- **Weight**: one cell per call, as the source defines it. A caller wanting n units per request uses
  TokenBucket.
- **Worked example (hand-derived from I.371 A.1)**: T = 100 ms, τ = 200 ms, all at t=0: three calls conform
  (TAT 100, 200, 300 ms), the fourth is refused with retryAfter 100 ms; at t=100 ms one conforms (TAT 400 ms).
- **Errors**: `ErrInterval` (emissionInterval ≤ 0), `ErrTolerance` (tolerance < 0).
- **Tags**: rate-limiting, leaky-bucket, gcra, virtual-scheduling, traffic-policing, atm.

### ALG-RATELM-003 FixedWindow
- **Source**: Berger, Eckberg, Hou, Lucantoni, GLOBECOM '90 (above), §4 p. 351, the jumping window scheme:
  "the counter starts at zero, at the beginning of a fixed time interval (window). The counter is incremented
  upon cell arrivals within the window, and resets to zero at the end of the window. The next window starts
  immediately after the preceding window" (read). Rathgeb, "Modeling and performance comparison of policing
  mechanisms for ATM networks", IEEE J. Sel. Areas Commun. 9(3):325–334, 1991, doi:10.1109/49.76630
  compares the same jumping and moving windows (Crossref and abstract confirmed; full text not read).
- **Best known**: O(1) time, O(1) space per call.
- **Approach**: the first call opens window 0 at its `now`; windows are [start + jW, start + (j+1)W). On a
  call past the window, start jumps by ⌊elapsed/W⌋·W (never past now, so no overflow) and the count resets.
  A request conforms when count + n ≤ limit (checked as n ≤ limit − count). retryAfter is the time to the
  window's end. **Lost**: windows aligned to the Unix epoch (`now.Truncate(W)`; useful across machines but
  the source anchors at a window beginning, and the caller can align by passing aligned times); the triggered
  jumping window (Rathgeb's TJW: window starts at the first arrival after the last expired) — not requested.
- **Known property (tested)**: up to 2·limit can pass in a span shorter than W across a window edge; the
  sliding log (004) does not allow this.
- **Worked example (hand-derived)**: limit 3, W = 1 s: at 0, 0.2, 0.9 s ok; at 0.95 s refused, retryAfter
  50 ms; at 1.0 s ok (new window).
- **Errors**: `ErrLimit` (limit < 1), `ErrWindow` (window ≤ 0).
- **Tags**: rate-limiting, fixed-window, jumping-window, counter, traffic-policing.

### ALG-RATELM-004 SlidingWindowLog
- **Source**: Berger, Eckberg, Hou, Lucantoni, GLOBECOM '90, §4 p. 351, the moving window scheme: "the maximum
  number of cell arrivals within a window of fixed time interval (W time slots) is limited as in the jumping
  window scheme. The difference is that adjacent windows overlap by W−1 slots … it requires that the arrival
  time of a cell be stored for W slots" (read). Rathgeb 1991 (above) as the second, metadata-confirmed source.
- **Best known**: O(1) amortized time per accepted call (each entry is pushed and evicted once), O(limit)
  space: the log must remember every accepted arrival still inside the window, and up to `limit` of them
  can be (lower bound for an exact moving window). A refused call costs O(min(n, limit)) to compute
  retryAfter (O(1) for n = 1).
- **Approach**: a ring buffer of (time, n) entries, `queue.CircularBuffer` (ALG-QUEUE-002, reused by call per
  L11), capacity `limit` (each entry holds ≥ 1 unit, so it never fills), and a running total. Entry t counts
  for now ∈ [t, t+W). retryAfter = the time until enough of the oldest entries expire for n to fit.
  **Lost**: a slice of timestamps re-sliced from the front (leaks the head until growth copies);
  per-unit entries (n entries per request; O(n) per call).
- **Worked example (hand-derived from the moving-window rule)**: limit 3, W = 1 s: at 0, 0.2, 0.9 s ok; at
  0.95 s refused, retryAfter 50 ms (the 0 s entry expires at 1.0 s); at 1.0 s ok (the 0 s entry has left);
  at 1.1 s refused, retryAfter 100 ms (the 0.2 s entry expires at 1.2 s). Same arrivals as FixedWindow's
  example, so the two tests show where the schemes differ.
- **Errors**: `ErrLimit`, `ErrWindow`, `ErrTooLarge` (limit > 1<<22).
- **Tags**: rate-limiting, sliding-window, moving-window, sliding-log, traffic-policing. Pattern pointers
  only: code-heim go_69/go_77 (not sources).

### SlidingWindowCounter — HELD (no source)
The estimate count ≈ previous·(1 − elapsed/W) + current is described by vendor engineering blogs and library
documentation only (search 2026-11-12: limits.readthedocs.io, Kong, Vercel, Appwrite, others); no paper, RFC,
ITU recommendation or textbook found. Berger et al.'s "exponential smoothing scheme" (§4) is a different,
cited mechanism and could be offered instead if Devin wants a fixed-memory sliding estimate.

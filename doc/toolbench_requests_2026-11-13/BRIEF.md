# BRIEF: toolbench requests R-1, R-2, R-3 (2026-11-13), slice S20

The request: `/Users/devinanderson/Desktop/toolbench/doc/wave1_2026-11-13/calcbench_requests.md` (read-only).
This brief is the slice spec. DECIDED by Devin (D-087): D1 algo/backoff (SHORT BACKOF), D2 num/integer CeilDiv,
D3 algo/cron (SHORT CRON). Stdlib only. Builders: gofmt -l, go build, go vet on their own package only; write
tests, fuzz and benchmarks; never run them. Every command bounded with `perl -e 'alarm N; exec @ARGV' --`.

## B1 algo/backoff (new) — ALG-BACKOF-001…004

| ID | Function | File |
|---|---|---|
| ALG-BACKOF-001 | `Exponential(attempt int, base, max time.Duration) (time.Duration, error)` | exponential.go |
| ALG-BACKOF-002 | `FullJitter(attempt int, base, max time.Duration, u float64) (time.Duration, error)` | full_jitter.go |
| ALG-BACKOF-003 | `EqualJitter(attempt int, base, max time.Duration, u float64) (time.Duration, error)` | equal_jitter.go |
| ALG-BACKOF-004 | `DecorrelatedJitter(prev, base, max time.Duration, u float64) (time.Duration, error)` | decorrelated_jitter.go |

- Exponential: sleep = min(max, base·2^attempt) in integer nanoseconds, saturating: base<<attempt is formed only
  when attempt < bits.LeadingZeros64(uint64(base)) (the result then fits in 63 bits, the ALG-BITS-003 pattern),
  otherwise the result is max. attempt ≥ 63 → max. No float.
- The draw u ∈ [0, 1) is scaled exactly in integers: k = uint64(u·2⁵³) (exact, u < 1), scale(d, u) =
  ⌊d·k / 2⁵³⌋ via bits.Mul64 (128-bit product), so scale(d, u) ∈ [0, d) for d > 0 and 0 for u = 0. Unexported helper
  in full_jitter.go.
- FullJitter = scale(sleep, u) ∈ [0, sleep). EqualJitter = (sleep − sleep/2) + scale(sleep/2, u) ∈ [⌈sleep/2⌉, sleep]
  (Brooker's temp/2 + random_between(0, temp/2); the ceiling half keeps an odd nanosecond, result ≤ sleep).
- DecorrelatedJitter = min(max, base + scale(hi − base, u)), hi = 3·prev saturating at MaxInt64 (Brooker
  random_between(base, sleep·3), then min(cap, ·)). prev < base → ErrPrev (the previous sleep is never below base).
- Errors (exported sentinels in exponential.go, wrapped with %w and the value): ErrAttempt attempt < 0;
  ErrBase base ≤ 0; ErrMax max < base; ErrDraw u NaN, < 0 or ≥ 1; ErrPrev prev < base. Check order: attempt/prev,
  base, max, u.
- Sources (both, on every function as relevant): Brooker, "Exponential Backoff And Jitter", AWS Architecture Blog,
  2015-03-04, https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/ ; gRPC, "Connection Backoff
  Protocol", https://github.com/grpc/grpc/blob/master/doc/connection-backoff.md (the cap MAX_BACKOFF, min(current·
  multiplier, MAX_BACKOFF)); saturating shift: Warren, Hacker's Delight, 2nd ed., 2013, §5-3 / ALG-BITS-003 (nlz).
- Worked (Brooker numbers per request): base 100 ms, max 10 s: attempt 3 → 800 ms; FullJitter u 0.5 → 400 ms;
  EqualJitter u 0.5 → 600 ms; attempt 10 → 10 s; Decorrelated(prev 300 ms, base 100 ms, max 10 s, u 0.5) → 500 ms.
- Tests: worked examples; attempt 0 (→ base), 62, 63, 64, MaxInt; max < base; max == base; base = 1 ns and
  MaxInt64; u = 0 and math.Nextafter(1, 0) (result < sleep, ≥ sleep − 1 for FullJitter); u = NaN, −0.0 ok?,
  −tiny, 1, +Inf; Decorrelated prev = MaxInt64 (saturation); Fuzz each (result within its documented bounds,
  never error on valid input), run with -fuzztime; AllocsPerRun == 0; Benchmark each.

## B2 num/integer — NUM-INTEG-013 CeilDiv

- `func CeilDiv(a, b int64) (int64, error)` in ceil_div.go. q, r := a/b, a%b; if r != 0 && (r > 0) == (b > 0) { q++ }.
- b = 0 → ErrZeroDenominator (num/integer's existing division-by-zero sentinel, fraction.go:11; no ErrDivByZero
  exists). a = MinInt64, b = −1 → ErrOverflow (lcm.go:10). Exact for every other int64 pair.
- Source: Graham, Knuth, Patashnik, Concrete Mathematics, 2nd ed., Addison-Wesley, 1994, §3.1–3.2 (⌈x⌉ = −⌊−x⌋;
  ⌈n/m⌉ = ⌊(n + m − 1)/m⌋ for m > 0); The Go Programming Language Specification, "Integer operators"
  (truncated division, q·b + r = a, sign of r = sign of a; MinInt64 / −1 overflows).
- Worked: (10,3)=4, (9,3)=3, (0,5)=0, (−7,2)=−3, (MaxInt64,1)=MaxInt64, (2⁵³+1,1)=2⁵³+1. Plus all four sign
  combinations, (MinInt64,1), (MinInt64,2), (MaxInt64,−1), (MinInt64,MinInt64), (MaxInt64,MaxInt64), (1,MinInt64)=0
  (⌈1/−2⁶³⌉ = 0), errors. Oracle: math/big Int (Div with Euclid adjust → ceil) on a few hundred seeded pairs;
  Fuzz against math/big. AllocsPerRun 0. Benchmark.
- Add the line `NUM-INTEG-013  CeilDiv  ⌈a/b⌉ for int64, exact, overflow-checked` to num/integer/doc.go.

## B3 algo/cron (new) — ALG-CRON-001 ParseCron, ALG-CRON-002 CronSpec.Next

- `func ParseCron(spec string) (CronSpec, error)`; `func (c CronSpec) Next(after time.Time) (time.Time, error)`.
  CronSpec: a value type of bitsets (minute uint64, hour uint32, dom uint32, month uint16, dow uint8) plus two bools
  (dom literally `*`, dow literally `*`). Zero value never fires.
- Syntax: POSIX only. Exactly five blank-separated fields: minute 0–59, hour 0–23, day of month 1–31, month 1–12,
  day of week 0–6 (0 = Sunday). Each field `*` or a comma list of elements, an element a number or `lo-hi` (lo ≤ hi).
  No steps (`/`), names, `7`, `?`, `L`, `@hourly`, seconds: ErrSyntax. Number outside the field range or lo > hi:
  ErrRange. Errors wrapped with %w naming the field.
- Day matching: month always restricts. If dom and dow are both not `*`, a day matches when either matches; if
  only one is restricted, it alone decides; both `*` every day. (The request's rule and common cron; POSIX's literal
  text also ORs when the MONTH is restricted and dom is `*` — this is an assumption, A-row.)
- ErrNeverFires from ParseCron when no day can ever match: dow is `*`, dom restricted, and no allowed month has
  an allowed day ≤ its longest length (Feb = 29). Next on a zero/never spec → ErrNeverFires; Next also stops with
  ErrNeverFires after scanning a bounded window (8 years cover Feb 29 across 2100; use e.g. 9 years) — no loop is
  unbounded.
- Next: the earliest instant strictly after `after` whose wall clock in after.Location() matches; seconds and
  nanoseconds are zero. Search civil fields year → month → day → hour → minute with carries, each next allowed
  value found with bits.TrailingZeros on the masked bitset; never minute by minute.
- DST, written into the doc comment exactly: a matching wall time that does not exist (spring-forward gap) fires
  at the first instant after the gap (the transition instant, e.g. 03:00 EDT); several matching times inside one
  gap fire once. A matching wall time that occurs twice (fall-back overlap) fires once, at its first (earlier-offset)
  occurrence; the second occurrence never fires. Use time.Time.ZoneBounds to find the transition; do not trust
  time.Date's unspecified choice in a gap or overlap.
- Source: IEEE Std 1003.1-2024 (POSIX.1-2024), Shell & Utilities volume, "crontab" utility, INPUT FILES,
  https://pubs.opengroup.org/onlinepubs/9799919799/utilities/crontab.html ; Go time package (Date, ZoneBounds) for
  zone rules.
- Time: Next O(1): at most 9 years × 12 months × 31 days candidate days, one trailing-zero count per hour/minute;
  state it as such. ParseCron O(n), n = len(spec).
- Worked: "30 2 * * *" in America/New_York after 2026-03-07T03:00-05:00 → 2026-03-08T03:00-04:00; "0 10 * * 1-5"
  after 2026-11-13T10:00Z → 2026-11-16T10:00Z; "0 0 31 2 *" → ErrNeverFires.
- Tests: worked examples; fall-back 2026-11-01 New_York "30 1 * * *" after 00:00 EDT → 01:30 EDT, after 01:30 EDT →
  2026-11-02 01:30 EST, after 01:10 EST → next day; "* 2 * * *" on spring-forward day after 01:59 → 03:00 EDT once,
  and "* 2 * * *" after 03:00 EDT → next day 02:00; "0 0 31 * *" (skips 30-day months); "0 0 29 2 *"
  from 2097 → 2104-02-29; "0 0 13 * 5" both restricted (OR); Feb 30/31 never; month only "0 0 * 2 *"; UTC and a
  fixed zone; after with seconds; after = a match (strictly after); year end carry; parse errors (4 and 6 fields,
  "*/5", "60", "24", "0", "32", "13", "7", "5-3", "", "1,,2", "-1", "a"). Oracle: a minute-stepping brute force in UTC
  over a capped window (few hundred random specs from a fixed seed, ≤ 2 years each) — UTC only, no DST. Fuzz
  ParseCron (no panic; round-trip Next result matches) with -fuzztime. Benchmark ParseCron and Next.
  time.LoadLocation("America/New_York") needs zoneinfo: use it; if it errors, t.Skip with the reason (stdlib
  time/tzdata embed is allowed in the test file only if the reviewer agrees — default: import _ "time/tzdata" in
  the _test.go so the test never depends on the host).

## Shared (lead only)

INDEX.md (regenerated, D-070), doc/taxonomy.md rows (algo/backoff BACKOF built; algo/cron CRON built),
doc/scope/SCOPE.md S20, DECISIONS D-087, ASSUMPTIONS rows, CHANGELOG, doc/requests/toolbench_reply_2026-11-13.md.

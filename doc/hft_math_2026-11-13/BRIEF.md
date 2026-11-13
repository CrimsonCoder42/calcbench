# Brief: HFT and execution maths (2026-11-13), S19

EDITS CODE. Risk NORMAL. Builders read this file and do not re-orient.

Rules that apply to every line: .claude/CLAUDE.md and .claude/rules/{go-functions,go-testing,go-complexity,numeric,
sources}.md. Read them once; they are short.

## DECIDED (verbatim; do not ask)

D1 Prices and sizes: price helpers take int64 ticks (price = ticks × tickSize); statistics take float64. Returns are
   float64 fractions (0.01 = 1%), never percent. Money stays int64 cents where a function already uses it. Basis
   points: 1 bp = 1e-4, helpers convert both ways. No hidden constants: tick size, annualisation factor (252, 390,
   seconds per day...), EWMA lambda, confidence level and risk-free rate are always parameters.
D2 New package fin/market, SHORT MKT (add the row to doc/taxonomy.md). Additions to fin/invest, num/series,
   num/signal and algo/stream continue those prefixes' ID sequences.
D3 Streaming types follow ALG-STREAM-004 RunningStats: a struct with Add, the readers, Reset, and Merge where the
   mathematics allows it; O(1) per Add, O(w) memory for windowed ones using the existing CircularBuffer / MonotonicDeque.
   Windowed variance uses the sliding Welford update (West 1979), never a re-sum of the window.
D4 Every function follows the doc-comment template (ID, Time, Space, Source to chapter or equation, Tags, Units,
   Example). "No source, no ship": a row whose primary source you cannot cite to the equation is left in the work
   folder as NEEDS SOURCE, not built. Industry blogs, Investopedia, Medium, StackExchange and Reddit are never the
   citation; the paper or textbook they point to is.
D5 Tests use oracle values worked by hand or from the cited source on small fixed inputs; each test under 60 s; no
   fuzz without -fuzztime. Builders write tests and never run them.

## Lead's notes

- **Reuse, never duplicate.** Grep INDEX.md and the package for the same mathematics under another name. A match is
  a REUSED row, not a new function.
  - Already exist: ALG-WINDOW-002 RollingMean, ALG-WINDOW-014 ExponentialMovingAverage, ALG-STREAM-004 RunningStats,
    FIN-INVEST-025 SharpeRatio, FIN-INVEST-030 BidAskSpread, NUM-STATS-014 PearsonCorrelation, NUM-STATS-009
    LinearRegression, NUM-STATS-015 Percentile, NUM-STATS-028 ZScore, NUM-STATS-038/039 Sample/PopulationCovariance,
    NUM-STATS-021 SampleVariance, ALG-RED-001 Sum and the rest of algo/reduce (CumSum, Diff, ArgMax...).
  - Also available: num/prob NormalQuantile/NormalCDF (for parametric VaR), ALG-QUEUE-002 CircularBuffer and
    ALG-WINDOW-005 MonotonicDeque.
- **Imports.**
  - algo/stream may import other algo packages (algo/queue, algo/window, algo/reduce), never num/**.
  - fin/market and fin/invest may import algo/** and num/**. fin/market imports NO other fin package
    (doc/taxonomy.md: inside fin, only loan and invest may import interest). If fin/market needs a return helper,
    it works from its own inputs.
  - num/series and num/signal may import algo/** and earlier num packages (D-009 order: ... stats, prob; linalg after).
- **Shared private helpers** (finite checks, compensated sums) are never re-implemented: call reduce.AllFinite,
  reduce.Sum and reduce.NeumaierAdd (S18, D-081).
- **Errors.** Each package keeps its own sentinels; new sentinels go in the file of the concept that defines them.
  A literal such as 252, 0.94, 1.645 or 2.326 never appears in code. In tests and doc examples it appears only as a
  stated parameter value with its source.
- **Streaming types** (D3, D-071): the type carries the ID, its methods need none, and New<T> needs none.
  - Struct with Add, readers, Reset, and Merge where the mathematics allows it.
  - O(1) per Add, checked by reading the code: no loop over the window inside Add.

## Owners (exclusive)

| Builder | Owns | Next free ID | Row file |
|---|---|---|---|
| B1 | fin/market (new; write doc.go with SHORT MKT) | FIN-MKT-001 | doc/hft_math_2026-11-13/B1.md |
| B2 | fin/invest | FIN-INVEST-031 | doc/hft_math_2026-11-13/B2.md |
| B3 | num/series, num/signal | NUM-SERIES-012, NUM-SIG-016 | doc/hft_math_2026-11-13/B3.md |
| B4 | algo/stream | ALG-STREAM-005 | doc/hft_math_2026-11-13/B4.md |

Row format in your file: `name | ID | source | BUILT` or `name | — | — | REUSED <ID>` or `name | — | what is
missing | NEEDS SOURCE`.

## Every builder

- Build, write doc.go entries and _test.go files with oracle values (D5), and benchmarks with b.ReportAllocs().
- Run ONLY: `gofmt -l <your dirs>`, `go build ./...`, `go vet ./<your pkg>/...`. **Never go test.** Never use
  `timeout`: it does not exist on this Mac.
- Never git, rm, or anything in the background. Never edit INDEX.md, taxonomy, DECISIONS, CHANGELOG, SCOPE or
  another builder's paths.
- If you must replace code in an existing file: first `mkdir -p _delete/<dir> && cp -n <path> _delete/<same path>`
  (add the .2026-11-13 suffix if it exists), then append a line to _delete/LEDGER.md. Existing exported signatures
  never change.
- Work within 25 minutes. If not finished, leave the code building and hand back the remaining list.
- Hand back at most 12 lines.

## Baseline (2026-11-13, after S18)

- `go build ./...` and `go vet ./...` are clean.
- `go test ./... -count=1 -run XXX -timeout 120s`: 76 ok.
- `go test -list`: 4,939 Test functions.
- INDEX.md: 1,533 rows = 1,505 live + 28 REMOVED.

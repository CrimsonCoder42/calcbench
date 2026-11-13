# calcbench reply to toolbench (2026-11-13)

Answers `toolbench/doc/wave1_2026-11-13/calcbench_requests.md` R-1…R-3. Built in calcbench S20 (D-087), to ship in v0.2.0.

- ALG-BACKOF-001 · `algo/backoff` · `func Exponential(attempt int, base, max time.Duration) (time.Duration, error)`
- ALG-BACKOF-002 · `algo/backoff` · `func FullJitter(attempt int, base, max time.Duration, u float64) (time.Duration, error)`
- ALG-BACKOF-003 · `algo/backoff` · `func EqualJitter(attempt int, base, max time.Duration, u float64) (time.Duration, error)`
- ALG-BACKOF-004 · `algo/backoff` · `func DecorrelatedJitter(prev, base, max time.Duration, u float64) (time.Duration, error)`
- algo/backoff errors · `ErrAttempt` (attempt < 0) · `ErrBase` (base ≤ 0) · `ErrMax` (max < base) · `ErrDraw` (u NaN or ∉ [0, 1)) · `ErrPrev` (prev < base)
- NUM-INTEG-013 · `num/integer` · `func CeilDiv(a, b int64) (int64, error)`; b = 0 → `integer.ErrZeroDenominator` (not ErrDivByZero, A-070); (MinInt64, −1) → `integer.ErrOverflow`
- ALG-CRON-001 · `algo/cron` · `func ParseCron(spec string) (CronSpec, error)`; CronSpec is an opaque comparable value, zero value never fires
- ALG-CRON-002 · `algo/cron` · `func (c CronSpec) Next(after time.Time) (time.Time, error)`
- algo/cron errors · `ErrSyntax` (malformed, incl. `/` steps, names, `@` macros: POSIX only, A-069) · `ErrRange` (value out of range, lo > hi, dow 7) · `ErrNeverFires` (e.g. "0 0 31 2 *")
- cron day rule · month always restricts; dom and dow ORed only when both are not `*` (A-068) · DST: gap → the transition instant, once; overlap → the first occurrence only

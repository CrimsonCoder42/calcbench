package backoff

import "time"

// EqualJitter returns Brooker's "Equal Jitter" delay: half of sleep kept, plus a draw in [0, sleep/2) chosen by
// the caller's u, where sleep = Exponential(attempt, base, max). The kept half is sleep − ⌊sleep/2⌋ = ⌈sleep/2⌉
// (an odd nanosecond stays), so the result is in [⌈sleep/2⌉, sleep]: ⌈sleep/2⌉ for u = 0, and below sleep
// except when sleep = 1 ns.
// Errors: ErrAttempt when attempt < 0; ErrBase when base ≤ 0; ErrMax when max < base; ErrDraw when u is
// NaN, < 0 or ≥ 1 (checked in that order; −0 is accepted as 0).
//
// ID:      ALG-BACKOF-003
// Time:    O(1)   one leading-zero count, one shift and one 64×64→128-bit multiply
// Space:   O(1), no allocation
// Source:  Brooker, "Exponential Backoff And Jitter", AWS Architecture Blog, 2015-03-04, https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/, figure 10 (Equal Jitter) (read 2026-11-13)
// Source:  grpc/grpc doc/connection-backoff.md (master), "Min(current_backoff * MULTIPLIER, MAX_BACKOFF)", https://github.com/grpc/grpc/blob/master/doc/connection-backoff.md (read 2026-11-13); saturating shift: Warren, Hacker's Delight, 2nd ed., Addison-Wesley, 2013, §5-3 nlz (the ALG-BITS-003 pattern)
// Tags:    backoff, retry, jitter, equal-jitter, exponential, duration
// Units:   attempt count; base, max and the result time.Duration (integer nanoseconds); u dimensionless in [0, 1)
// Example: base 100 ms, max 10 s, attempt 3 (sleep 800 ms), u 0.5 → 400 ms + 200 ms = 600 ms
func EqualJitter(attempt int, base, max time.Duration, u float64) (time.Duration, error) {
	if err := checkSchedule(attempt, base, max); err != nil {
		return 0, err
	}
	if err := checkDraw(u); err != nil {
		return 0, err
	}
	sleep := exponential(attempt, base, max)
	half := sleep / 2
	return (sleep - half) + scale(half, u), nil
}

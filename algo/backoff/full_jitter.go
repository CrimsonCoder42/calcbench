package backoff

import (
	"fmt"
	"math/bits"
	"time"
)

// FullJitter returns Brooker's "Full Jitter" delay: a value in [0, sleep) chosen by the caller's draw u,
// where sleep = Exponential(attempt, base, max). u = 0 gives 0; the result is ⌊sleep·⌊u·2⁵³⌋ / 2⁵³⌋, exact
// in integers, so it is always < sleep (for sleep ≤ 2⁵³ ns, u = the largest float64 below 1 gives sleep − 1).
// Errors: ErrAttempt when attempt < 0; ErrBase when base ≤ 0; ErrMax when max < base; ErrDraw when u is
// NaN, < 0 or ≥ 1 (checked in that order; −0 is accepted as 0).
//
// ID:      ALG-BACKOF-002
// Time:    O(1)   one leading-zero count, one shift and one 64×64→128-bit multiply
// Space:   O(1), no allocation
// Source:  Brooker, "Exponential Backoff And Jitter", AWS Architecture Blog, 2015-03-04, https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/, figure 6 (Full Jitter) (read 2026-11-13)
// Source:  grpc/grpc doc/connection-backoff.md (master), "Min(current_backoff * MULTIPLIER, MAX_BACKOFF)", https://github.com/grpc/grpc/blob/master/doc/connection-backoff.md (read 2026-11-13); saturating shift: Warren, Hacker's Delight, 2nd ed., Addison-Wesley, 2013, §5-3 nlz (the ALG-BITS-003 pattern)
// Tags:    backoff, retry, jitter, full-jitter, exponential, duration
// Units:   attempt count; base, max and the result time.Duration (integer nanoseconds); u dimensionless in [0, 1)
// Example: base 100 ms, max 10 s, attempt 3 (sleep 800 ms), u 0.5 → 400 ms
func FullJitter(attempt int, base, max time.Duration, u float64) (time.Duration, error) {
	if err := checkSchedule(attempt, base, max); err != nil {
		return 0, err
	}
	if err := checkDraw(u); err != nil {
		return 0, err
	}
	return scale(exponential(attempt, base, max), u), nil
}

// checkDraw validates u ∈ [0, 1); the negated test also rejects NaN.
func checkDraw(u float64) error {
	if !(u >= 0 && u < 1) {
		return fmt.Errorf("%w: u = %v", ErrDraw, u)
	}
	return nil
}

// scale returns ⌊d·k / 2⁵³⌋ with k = ⌊u·2⁵³⌋, for d ≥ 0 and u ∈ [0, 1). u·2⁵³ is an exact float (a power-of-two
// multiply) below 2⁵³, so k < 2⁵³; the product d·k < 2¹¹⁶ is formed exactly by bits.Mul64 and shifted right
// 53 bits. The result is in [0, d) for d > 0 and is 0 for u = 0 or d = 0.
func scale(d time.Duration, u float64) time.Duration {
	k := uint64(u * (1 << 53))
	hi, lo := bits.Mul64(uint64(d), k)
	return time.Duration(hi<<11 | lo>>53)
}

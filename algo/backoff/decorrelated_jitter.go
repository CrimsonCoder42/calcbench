package backoff

import (
	"fmt"
	"math"
	"time"
)

// DecorrelatedJitter returns Brooker's "Decorrelated Jitter" delay: min(max, a draw in [base, 3·prev)) chosen by
// the caller's u, where prev is the previous sleep (pass base before the first retry). The draw is
// base + ⌊(hi − base)·⌊u·2⁵³⌋ / 2⁵³⌋ with hi = 3·prev saturating at math.MaxInt64, so it never overflows; the
// result is in [base, max], and below hi except when base = math.MaxInt64 (then hi = base and the result is base).
// Errors: ErrPrev when prev < base; ErrBase when base ≤ 0; ErrMax when max < base; ErrDraw when u is NaN, < 0
// or ≥ 1 (checked in that order; −0 is accepted as 0).
//
// ID:      ALG-BACKOF-004
// Time:    O(1)   one 64×64→128-bit multiply
// Space:   O(1), no allocation
// Source:  Brooker, "Exponential Backoff And Jitter", AWS Architecture Blog, 2015-03-04, https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/, figure 11 (Decorrelated Jitter) (read 2026-11-13)
// Source:  grpc/grpc doc/connection-backoff.md (master), "Min(current_backoff * MULTIPLIER, MAX_BACKOFF)", https://github.com/grpc/grpc/blob/master/doc/connection-backoff.md (read 2026-11-13)
// Tags:    backoff, retry, jitter, decorrelated-jitter, duration
// Units:   prev, base, max and the result time.Duration (integer nanoseconds); u dimensionless in [0, 1)
// Example: prev 300 ms, base 100 ms, max 10 s, u 0.5 → 100 ms + 800 ms/2 = 500 ms
func DecorrelatedJitter(prev, base, max time.Duration, u float64) (time.Duration, error) {
	if prev < base {
		return 0, fmt.Errorf("%w: prev = %v, base = %v", ErrPrev, prev, base)
	}
	if err := checkBaseMax(base, max); err != nil {
		return 0, err
	}
	if err := checkDraw(u); err != nil {
		return 0, err
	}
	hi := time.Duration(math.MaxInt64)
	if prev <= math.MaxInt64/3 {
		hi = 3 * prev
	}
	// hi ≥ prev ≥ base > 0, so hi − base ≥ 0 and base + scale(hi − base, u) < hi ≤ MaxInt64.
	if sleep := base + scale(hi-base, u); sleep < max {
		return sleep, nil
	}
	return max, nil
}

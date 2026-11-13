package backoff

import (
	"errors"
	"fmt"
	"math/bits"
	"time"
)

// Errors returned by this package. They are wrapped with the offending value; test with errors.Is.
var (
	// ErrAttempt: the attempt number is negative.
	ErrAttempt = errors.New("backoff: attempt is negative")
	// ErrBase: the base delay is zero or negative.
	ErrBase = errors.New("backoff: base must be positive")
	// ErrMax: the cap is below the base delay.
	ErrMax = errors.New("backoff: max is below base")
	// ErrDraw: the draw u is NaN or outside [0, 1).
	ErrDraw = errors.New("backoff: u must be in [0, 1)")
	// ErrPrev: the previous sleep is below the base delay.
	ErrPrev = errors.New("backoff: prev is below base")
)

// Exponential returns the capped exponential backoff delay min(max, base·2^attempt): base for attempt 0,
// doubling each attempt until it reaches max. The product is formed only when it fits in 63 bits;
// otherwise the result is max, so it never overflows.
// Errors: ErrAttempt when attempt < 0; ErrBase when base ≤ 0; ErrMax when max < base (checked in that order).
//
// ID:      ALG-BACKOF-001
// Time:    O(1)   one leading-zero count and one shift
// Space:   O(1), no allocation
// Source:  Brooker, "Exponential Backoff And Jitter", AWS Architecture Blog, 2015-03-04, https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/, figure 3 (capped exponential) (read 2026-11-13)
// Source:  grpc/grpc doc/connection-backoff.md (master), "Min(current_backoff * MULTIPLIER, MAX_BACKOFF)", https://github.com/grpc/grpc/blob/master/doc/connection-backoff.md (read 2026-11-13); saturating shift: Warren, Hacker's Delight, 2nd ed., Addison-Wesley, 2013, §5-3 nlz (the ALG-BITS-003 pattern)
// Tags:    backoff, retry, exponential, cap, saturating, duration
// Units:   attempt count (0 = first retry); base, max and the result time.Duration (integer nanoseconds)
// Example: base 100 ms, max 10 s: attempt 3 → 800 ms; attempt 10 → 10 s (102.4 s capped)
func Exponential(attempt int, base, max time.Duration) (time.Duration, error) {
	if err := checkSchedule(attempt, base, max); err != nil {
		return 0, err
	}
	return exponential(attempt, base, max), nil
}

// checkSchedule validates attempt, base and max in the package's check order.
func checkSchedule(attempt int, base, max time.Duration) error {
	if attempt < 0 {
		return fmt.Errorf("%w: attempt = %d", ErrAttempt, attempt)
	}
	return checkBaseMax(base, max)
}

// checkBaseMax validates base and max.
func checkBaseMax(base, max time.Duration) error {
	if base <= 0 {
		return fmt.Errorf("%w: base = %v", ErrBase, base)
	}
	if max < base {
		return fmt.Errorf("%w: max = %v, base = %v", ErrMax, max, base)
	}
	return nil
}

// exponential is Exponential on validated input (attempt ≥ 0, 0 < base ≤ max). base<<attempt is formed only
// when attempt < nlz(base): base < 2^(64−nlz), so the shifted value is < 2^63 and fits in an int64.
func exponential(attempt int, base, max time.Duration) time.Duration {
	if attempt >= bits.LeadingZeros64(uint64(base)) {
		return max
	}
	if sleep := base << uint(attempt); sleep < max {
		return sleep
	}
	return max
}

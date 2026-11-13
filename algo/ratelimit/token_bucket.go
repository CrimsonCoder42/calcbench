package ratelimit

import (
	"errors"
	"fmt"
	"math"
	"math/bits"
	"time"
)

// Errors returned by this package. They are wrapped with the offending parameter; test with errors.Is.
var (
	// ErrRate: tokensPerInterval is less than 1.
	ErrRate = errors.New("ratelimit: tokens per interval must be at least 1")
	// ErrInterval: an interval is zero or negative.
	ErrInterval = errors.New("ratelimit: interval must be positive")
	// ErrBurst: burst is less than 1.
	ErrBurst = errors.New("ratelimit: burst must be at least 1")
	// ErrTolerance: a tolerance is negative.
	ErrTolerance = errors.New("ratelimit: tolerance must not be negative")
	// ErrLimit: limit is less than 1.
	ErrLimit = errors.New("ratelimit: limit must be at least 1")
	// ErrWindow: window is zero or negative.
	ErrWindow = errors.New("ratelimit: window must be positive")
	// ErrTooLarge: the meter's state would not fit (an int64 product overflows, or the log exceeds its cap).
	ErrTooLarge = errors.New("ratelimit: parameters too large")
)

// Never is the retryAfter of a request that can never conform: n is greater than the meter's burst or
// limit, the meter is a zero value, or its schedule has reached the end of representable time.
const Never time.Duration = math.MaxInt64

// TokenBucket is a token bucket meter: it holds at most burst tokens, gains tokensPerInterval tokens per
// interval, starts full at the first call, and a request for n tokens conforms when n tokens are
// present, which it then removes (RFC 2697 §3, the committed bucket alone). Refill is exact: the level
// is kept in integer units of 1/interval-nanoseconds of a token, so nothing drifts however long it runs.
// A refused request removes nothing. Not safe for concurrent use; the caller's mutex guards it.
//
// Methods: Allow(now, n) takes n tokens if present and returns (true, 0), else returns (false, wait)
// where wait is the time until n tokens will be present, or Never when n > burst. Available(now) is the
// whole number of tokens present at now.
//
// ID:      ALG-RATELM-001
// Time:    O(1) per call
// Space:   O(1)
// Source:  Heinanen and Guérin, "A Single Rate Three Color Marker", RFC 2697, IETF, 1999, §3 Metering (bucket C starts full at CBS, refills at CIR up to CBS; a packet of B conforms when Tc − B ≥ 0 and Tc is then decremented by B)
// Source:  ITU-T Recommendation I.371 (03/2004), Annex A.2, continuous-state leaky bucket (the equivalent GCRA form; tokens = burst − X′/T)
// Tags:    rate-limiting, token-bucket, traffic-policing, throttling, retry-after
type TokenBucket struct {
	tokensPerInterval int64     // units added per elapsed nanosecond
	intervalNs        int64     // units per token
	burst             int64     // tokens
	capacity          int64     // burst·intervalNs units
	level             int64     // units present, 0 ≤ level ≤ capacity
	last              time.Time // latest now seen
	started           bool
}

// NewTokenBucket returns a full TokenBucket that gains tokensPerInterval tokens every interval and
// holds at most burst tokens. For a rate of r per second use (r, time.Second, burst).
// Errors: ErrRate when tokensPerInterval < 1; ErrInterval when interval <= 0; ErrBurst when burst < 1;
// ErrTooLarge when burst × interval in nanoseconds exceeds math.MaxInt64.
func NewTokenBucket(tokensPerInterval int64, interval time.Duration, burst int64) (*TokenBucket, error) {
	if tokensPerInterval < 1 {
		return nil, fmt.Errorf("%w: tokensPerInterval = %d", ErrRate, tokensPerInterval)
	}
	if interval <= 0 {
		return nil, fmt.Errorf("%w: interval = %v", ErrInterval, interval)
	}
	if burst < 1 {
		return nil, fmt.Errorf("%w: burst = %d", ErrBurst, burst)
	}
	hi, lo := bits.Mul64(uint64(burst), uint64(interval))
	if hi != 0 || lo > math.MaxInt64 {
		return nil, fmt.Errorf("%w: burst = %d × interval = %v", ErrTooLarge, burst, interval)
	}
	return &TokenBucket{
		tokensPerInterval: tokensPerInterval,
		intervalNs:        int64(interval),
		burst:             burst,
		capacity:          int64(lo),
		level:             int64(lo),
	}, nil
}

// Allow takes n tokens at time now if they are present and returns (true, 0); otherwise it takes
// nothing and returns (false, the time until n tokens are present). n > burst returns (false, Never);
// n <= 0 returns (true, 0). A now before the latest now seen counts as that latest now.
func (b *TokenBucket) Allow(now time.Time, n int64) (ok bool, retryAfter time.Duration) {
	if n <= 0 {
		return true, 0
	}
	if n > b.burst {
		return false, Never
	}
	b.refill(now)
	need := n * b.intervalNs // n ≤ burst, so need ≤ capacity: no overflow
	if b.level >= need {
		b.level -= need
		return true, 0
	}
	deficit := need - b.level
	wait := deficit / b.tokensPerInterval
	if deficit%b.tokensPerInterval != 0 {
		wait++
	}
	return false, time.Duration(wait)
}

// Available returns the whole number of tokens present at time now (0 for a zero value).
func (b *TokenBucket) Available(now time.Time) int64 {
	if b.burst == 0 {
		return 0
	}
	b.refill(now)
	return b.level / b.intervalNs
}

// refill credits the time since the latest now seen. The first call starts the bucket full; a now that
// is not after the latest leaves everything as it is.
func (b *TokenBucket) refill(now time.Time) {
	if !b.started {
		b.started, b.last, b.level = true, now, b.capacity
		return
	}
	elapsed := int64(now.Sub(b.last)) // saturates at ±292 years
	if elapsed <= 0 {
		return
	}
	b.last = now
	room := b.capacity - b.level
	// Nanoseconds that would fill the room, rounded up; comparing first keeps the product below room.
	fill := room / b.tokensPerInterval
	if room%b.tokensPerInterval != 0 {
		fill++
	}
	if elapsed >= fill {
		b.level = b.capacity
		return
	}
	b.level += elapsed * b.tokensPerInterval // elapsed ≤ fill−1 < room/tokensPerInterval, so the product < room
}

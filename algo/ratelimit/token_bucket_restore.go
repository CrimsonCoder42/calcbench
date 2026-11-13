package ratelimit

import "time"

// Restore gives back n tokens taken earlier, at time now: the bucket first refills to now (a now before the
// latest now seen counts as that latest now, as in Allow), then gains n tokens, capped at burst. A bucket
// never used yet starts full at that refill, so Restore leaves it full. n <= 0 and the zero value change
// nothing, not even the refill.
// Restore does not undo the passage of time and does not know which tokens a later Allow took: the caller
// passes how many to return (golang.org/x/time/rate's Reservation.CancelAt computes that count the same way).
//
// ID:      ALG-RATELM-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Heinanen and Guérin, "A Single Rate Three Color Marker", RFC 2697, IETF, 1999, §3 Metering (bucket C is capped at CBS; Tc is incremented by one up to CBS)
// Source:  The Go Authors, golang.org/x/time/rate, func (*Reservation) CancelAt (tokens not consumed by later reservations are restored; the bucket is capped at burst), https://pkg.go.dev/golang.org/x/time/rate (read 2026-11-13)
// Tags:    rate-limiting, token-bucket, refund, cancel, throttling
// Units:   n in tokens; now is an instant
// Example: burst 5, Allow(t0, 3) → 2 tokens left; Restore(3, t0) → 5; Restore(10, t0) → 5 (capped at burst)
func (b *TokenBucket) Restore(n int64, now time.Time) {
	if n <= 0 || b.burst == 0 {
		return
	}
	b.refill(now)
	if n >= b.burst {
		b.level = b.capacity
		return
	}
	add := n * b.intervalNs // n < burst, so add < capacity: no overflow
	if add >= b.capacity-b.level {
		b.level = b.capacity
		return
	}
	b.level += add
}

package ratelimit

import (
	"fmt"
	"time"
)

// LeakyBucket is the leaky bucket of the Generic Cell Rate Algorithm GCRA(T, τ) in its virtual
// scheduling form: one cell per call, conforming when it arrives no earlier than the theoretical arrival
// time TAT minus the tolerance τ. The first cell sets TAT to its arrival time; a conforming cell moves TAT
// to max(TAT, now) + T; a non-conforming cell changes nothing. Equivalently (I.371 A.2), a bucket draining
// one unit per unit time, filled by T per conforming cell, that admits a cell while its content is at most
// τ: T is the emission interval (1/rate) and τ = (b−1)·T lets b cells pass back to back.
// Not safe for concurrent use; the caller's mutex guards it.
//
// Methods: Allow(now) returns (true, 0) and schedules the cell if it conforms, else (false, wait) where
// wait = TAT − τ − now is the time until it would conform.
//
// ID:      ALG-RATELM-002
// Time:    O(1) per call
// Space:   O(1)
// Source:  ITU-T Recommendation I.371 (03/2004), "Traffic control and congestion control in B-ISDN", Annex A, GCRA(T,τ), §A.1 Virtual Scheduling Algorithm and §A.2 continuous-state leaky bucket, pp. 86–87
// Source:  Berger, Eckberg, Hou, Lucantoni, "Performance characterizations of traffic monitoring, and associated control, mechanisms for broadband 'packet' networks", Proc. IEEE GLOBECOM '90, 1990, pp. 350–354, §4 (the leaky-bucket scheme), doi:10.1109/GLOCOM.1990.116535
// Tags:    rate-limiting, leaky-bucket, gcra, virtual-scheduling, traffic-policing, atm
type LeakyBucket struct {
	emission  time.Duration // T
	tolerance time.Duration // τ
	tat       time.Time     // theoretical arrival time of the next cell
	last      time.Time     // latest now seen
	started   bool
}

// NewLeakyBucket returns a LeakyBucket for GCRA(emissionInterval, tolerance).
// Errors: ErrInterval when emissionInterval <= 0; ErrTolerance when tolerance < 0.
func NewLeakyBucket(emissionInterval, tolerance time.Duration) (*LeakyBucket, error) {
	if emissionInterval <= 0 {
		return nil, fmt.Errorf("%w: emissionInterval = %v", ErrInterval, emissionInterval)
	}
	if tolerance < 0 {
		return nil, fmt.Errorf("%w: tolerance = %v", ErrTolerance, tolerance)
	}
	return &LeakyBucket{emission: emissionInterval, tolerance: tolerance}, nil
}

// Allow decides the cell arriving at now: (true, 0) when it conforms (TAT − now ≤ τ), after which TAT is
// max(TAT, now) + T; otherwise (false, TAT − now − τ) and nothing changes. A zero value, or a TAT that
// would pass the end of representable time, returns (false, Never). A now before the latest now seen
// counts as that latest now.
func (b *LeakyBucket) Allow(now time.Time) (ok bool, retryAfter time.Duration) {
	if b.emission <= 0 {
		return false, Never
	}
	if !b.started {
		b.started, b.last, b.tat = true, now, now
	} else if now.Before(b.last) {
		now = b.last
	} else {
		b.last = now
	}
	early := b.tat.Sub(now) // saturating
	if early > b.tolerance {
		return false, early - b.tolerance
	}
	base := b.tat
	if now.After(base) {
		base = now
	}
	next := base.Add(b.emission) // Time.Add saturates at the end of time
	if !next.After(base) {
		return false, Never
	}
	b.tat = next
	return true, 0
}

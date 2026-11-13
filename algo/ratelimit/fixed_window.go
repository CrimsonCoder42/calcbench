package ratelimit

import (
	"fmt"
	"time"
)

// FixedWindow is the jumping window meter: time is cut into consecutive windows of length window, the
// first starting at the first call; a counter starts at zero at the beginning of each window, and a
// request for n units conforms while count + n ≤ limit. Up to 2·limit units can pass within less than one
// window across a window edge; SlidingWindowLog does not allow that. Not safe for concurrent use; the
// caller's mutex guards it.
//
// Methods: Allow(now, n) counts n and returns (true, 0) if they fit in the current window, else
// (false, the time until the window ends), or (false, Never) when n > limit.
//
// ID:      ALG-RATELM-003
// Time:    O(1) per call
// Space:   O(1)
// Source:  Berger, Eckberg, Hou, Lucantoni, "Performance characterizations of traffic monitoring, and associated control, mechanisms for broadband 'packet' networks", Proc. IEEE GLOBECOM '90, 1990, pp. 350–354, §4 p. 351 (the jumping window scheme), doi:10.1109/GLOCOM.1990.116535
// Source:  Rathgeb, "Modeling and performance comparison of policing mechanisms for ATM networks", IEEE Journal on Selected Areas in Communications 9(3):325–334, 1991, doi:10.1109/49.76630
// Tags:    rate-limiting, fixed-window, jumping-window, counter, traffic-policing
type FixedWindow struct {
	limit   int64
	window  time.Duration
	start   time.Time // beginning of the current window
	count   int64     // units counted in the current window, 0 ≤ count ≤ limit
	last    time.Time // latest now seen, never before start
	started bool
}

// NewFixedWindow returns a FixedWindow admitting at most limit units per window.
// Errors: ErrLimit when limit < 1; ErrWindow when window <= 0.
func NewFixedWindow(limit int64, window time.Duration) (*FixedWindow, error) {
	if limit < 1 {
		return nil, fmt.Errorf("%w: limit = %d", ErrLimit, limit)
	}
	if window <= 0 {
		return nil, fmt.Errorf("%w: window = %v", ErrWindow, window)
	}
	return &FixedWindow{limit: limit, window: window}, nil
}

// Allow counts n units at time now and returns (true, 0) when they fit in the window containing now;
// otherwise it counts nothing and returns (false, the time until that window ends). n > limit returns
// (false, Never); n <= 0 returns (true, 0). A now before the latest now seen counts as that latest now.
func (w *FixedWindow) Allow(now time.Time, n int64) (ok bool, retryAfter time.Duration) {
	if n <= 0 {
		return true, 0
	}
	if n > w.limit {
		return false, Never
	}
	if !w.started {
		w.started, w.start, w.count = true, now, 0
	} else if now.Before(w.last) {
		now = w.last
	}
	w.last = now
	elapsed := now.Sub(w.start) // saturating; ≥ 0 since start ≤ last ≤ now
	if elapsed >= w.window {
		jump := elapsed / w.window * w.window // ≤ elapsed: the new start is not after now
		w.start = w.start.Add(jump)
		w.count = 0
		elapsed -= jump
	}
	if n <= w.limit-w.count {
		w.count += n
		return true, 0
	}
	return false, w.window - elapsed
}

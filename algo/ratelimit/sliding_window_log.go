package ratelimit

import (
	"fmt"
	"time"

	"github.com/CrimsonCoder42/calcbench/algo/queue"
)

// maxLogEntries caps SlidingWindowLog's limit: 1<<22 entries of 32 bytes is 128 MiB, the largest
// allocation a library call makes silently in this run (lead decisions L13, L26).
const maxLogEntries = 1 << 22

// logEntry is one accepted request: its time and its units.
type logEntry struct {
	at time.Time
	n  int64
}

// SlidingWindowLog is the moving window meter: it remembers the time of every accepted request still
// inside the window, and a request for n units at now conforms while the units accepted in
// (now − window, now] plus n are at most limit. A request accepted at t counts until t + window,
// exclusive. Not safe for concurrent use; the caller's mutex guards it.
//
// Methods: Allow(now, n) records n and returns (true, 0) if they fit, else (false, the time until enough
// of the oldest requests leave the window for n to fit), or (false, Never) when n > limit.
//
// ID:      ALG-RATELM-004
// Time:    O(1) amortized per accepted call (each entry added and evicted once); a refused call O(min(n, e)), e = entries held; NewSlidingWindowLog O(limit)
// Space:   O(limit), one allocation in NewSlidingWindowLog (32 bytes per entry)
// Source:  Berger, Eckberg, Hou, Lucantoni, "Performance characterizations of traffic monitoring, and associated control, mechanisms for broadband 'packet' networks", Proc. IEEE GLOBECOM '90, 1990, pp. 350–354, §4 p. 351 (the moving window scheme), doi:10.1109/GLOCOM.1990.116535
// Source:  Rathgeb, "Modeling and performance comparison of policing mechanisms for ATM networks", IEEE Journal on Selected Areas in Communications 9(3):325–334, 1991, doi:10.1109/49.76630
// Tags:    rate-limiting, sliding-window, moving-window, sliding-log, traffic-policing
type SlidingWindowLog struct {
	limit   int64
	window  time.Duration
	log     *queue.CircularBuffer[logEntry] // accepted requests, oldest first, times nondecreasing
	total   int64                           // units in log, 0 ≤ total ≤ limit
	last    time.Time                       // latest now seen
	started bool
}

// NewSlidingWindowLog returns a SlidingWindowLog admitting at most limit units in any span of length
// window. It allocates room for limit entries once.
// Errors: ErrLimit when limit < 1; ErrWindow when window <= 0; ErrTooLarge when limit > 1<<22.
func NewSlidingWindowLog(limit int64, window time.Duration) (*SlidingWindowLog, error) {
	if limit < 1 {
		return nil, fmt.Errorf("%w: limit = %d", ErrLimit, limit)
	}
	if window <= 0 {
		return nil, fmt.Errorf("%w: window = %v", ErrWindow, window)
	}
	if limit > maxLogEntries {
		return nil, fmt.Errorf("%w: limit = %d > %d", ErrTooLarge, limit, maxLogEntries)
	}
	// Each entry holds at least one unit and total ≤ limit, so limit slots never fill.
	log, err := queue.NewCircularBuffer[logEntry](int(limit), queue.Reject)
	if err != nil {
		return nil, err
	}
	return &SlidingWindowLog{limit: limit, window: window, log: log}, nil
}

// Allow records n units at time now and returns (true, 0) when the units accepted in the window ending at
// now plus n are at most limit; otherwise it records nothing and returns (false, the shortest wait after
// which they would be). n > limit returns (false, Never); n <= 0 returns (true, 0). A now before the
// latest now seen counts as that latest now.
func (s *SlidingWindowLog) Allow(now time.Time, n int64) (ok bool, retryAfter time.Duration) {
	if n <= 0 {
		return true, 0
	}
	if n > s.limit || s.log == nil {
		return false, Never
	}
	if s.started && now.Before(s.last) {
		now = s.last
	}
	s.started, s.last = true, now
	for {
		oldest, found := s.log.Front()
		if !found || now.Sub(oldest.at) < s.window {
			break
		}
		s.log.Pop()
		s.total -= oldest.n
	}
	if n <= s.limit-s.total {
		if _, err := s.log.Push(logEntry{at: now, n: n}); err != nil {
			return false, Never // unreachable: see NewSlidingWindowLog
		}
		s.total += n
		return true, 0
	}
	// Free the oldest entries until n fits; the one that makes room decides the wait.
	excess := s.total + n - s.limit // 1 ≤ excess ≤ n, no overflow: total, n ≤ limit
	var freed int64
	for i := range s.log.Len() {
		e, _ := s.log.At(i)
		freed += e.n
		if freed >= excess {
			return false, s.window - now.Sub(e.at) // 0 ≤ now − e.at < window
		}
	}
	return false, Never // unreachable: freeing every entry frees total ≥ excess
}

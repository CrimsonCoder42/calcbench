package stream

import (
	"errors"
	"fmt"
	"math"
	"time"
)

// Errors of HitCounter, wrapped with the offending value; test with errors.Is.
var (
	// ErrSpan: a window or resolution is not positive, or the resolution does not divide the window.
	ErrSpan = errors.New("stream: window and resolution must be positive and the resolution must divide the window")
	// ErrBuckets: window/resolution exceeds maxBuckets.
	ErrBuckets = errors.New("stream: window/resolution exceeds 1<<23 buckets")
	// ErrHits: a hit count is negative.
	ErrHits = errors.New("stream: hit count must not be negative")
	// ErrTimeRange: a time is outside the years 1678–2262 that time.Time.UnixNano can express.
	ErrTimeRange = errors.New("stream: time outside the range of UnixNano")
)

// maxBuckets caps a HitCounter at 1<<23 buckets (two int64 each: 128 MiB), the run's bound on what one
// library call allocates silently (lead decisions L13, L35).
const maxBuckets = 1 << 23

// HitCounter counts hits over the last window of time, at a resolution the caller chooses: time is
// cut into ticks of length resolution (tick = ⌊UnixNano / resolution⌋), and Count(now) sums the hits whose
// tick is one of the window/resolution ticks ending with now's tick. A hit at time x is therefore counted
// iff tick(now) − B < tick(x) ≤ tick(now), B = window/resolution: the effective window is between
// window − resolution and window long, never more (the resolution is the error bound). With window 300 s
// and resolution 1 s, Count at second 301 includes the hits of seconds 2…301 and drops second 1.
// The hits of one tick share one bucket of a ring of B buckets, slot = tick mod B, each stamped with its
// tick so a bucket left from an older lap of the ring reads as empty. Hits in the future of now are not
// counted. Time may go backward: a hit older than the window of the newest tick already in its slot is
// dropped (no Count at or after that tick could include it). The caller passes now on every call and owns
// locking. The zero value is not usable (its methods return ErrUninitialized).
//
// Methods: Hit(now) records one hit, Add(now, hits) records several, Count(now) returns the hits in the
// window ending at now.
//
// ID:      ALG-STREAM-002
// Time:    Hit, Add O(1); Count O(B)   B = window/resolution buckets; NewHitCounter O(B)
// Space:   O(B), two allocations in NewHitCounter; no method allocates
// Source:  Datar, Gionis, Indyk, Motwani, "Maintaining stream statistics over sliding windows", SIAM Journal on Computing, 2002, 31(6):1794–1813, doi:10.1137/S0097539701398363 (the sliding-window model; an exact count over N time units needs Θ(N) bits, here B buckets at the chosen resolution)
// Source:  Knuth, The Art of Computer Programming, Vol. 1, 3rd ed., Addison-Wesley, 1997, §2.2.2 (a circular array indexed modulo M)
// Tags:    stream, hit-counter, rate-limit, sliding-window, time-window, ring-buffer, buckets, counter
type HitCounter struct {
	resolutionNanos int64
	ticks           []int64 // ticks[slot] is the tick whose hits counts[slot] holds
	counts          []int64
}

// NewHitCounter returns an empty HitCounter over the last window of time, in buckets of resolution.
// Errors: ErrSpan when window or resolution is not positive or window is not a multiple of resolution;
// ErrBuckets when window/resolution > 1<<23.
func NewHitCounter(window, resolution time.Duration) (*HitCounter, error) {
	if window <= 0 || resolution <= 0 || window%resolution != 0 {
		return nil, fmt.Errorf("%w: window = %v, resolution = %v", ErrSpan, window, resolution)
	}
	buckets := int64(window / resolution)
	if buckets > maxBuckets {
		return nil, fmt.Errorf("%w: window/resolution = %d", ErrBuckets, buckets)
	}
	h := &HitCounter{
		resolutionNanos: int64(resolution),
		ticks:           make([]int64, buckets),
		counts:          make([]int64, buckets),
	}
	for i := range h.ticks {
		h.ticks[i] = math.MinInt64 // older than every real tick: the first hit in a slot resets it
	}
	return h, nil
}

// Hit records one hit at now. Same as Add(now, 1). O(1).
// Errors: ErrTimeRange, ErrOverflow, ErrUninitialized, as Add.
func (h *HitCounter) Hit(now time.Time) error {
	return h.Add(now, 1)
}

// Add records hits hits at now. O(1).
// Errors: ErrHits when hits < 0; ErrTimeRange when now is outside the UnixNano range; ErrOverflow when
// the bucket's count would exceed math.MaxInt64; ErrUninitialized on the zero value. On error nothing
// changes.
func (h *HitCounter) Add(now time.Time, hits int64) error {
	if h.counts == nil {
		return ErrUninitialized
	}
	if hits < 0 {
		return fmt.Errorf("%w: hits = %d", ErrHits, hits)
	}
	t, err := h.tick(now)
	if err != nil {
		return err
	}
	slot := h.slot(t)
	switch stamp := h.ticks[slot]; {
	case stamp == t:
		if hits > math.MaxInt64-h.counts[slot] {
			return fmt.Errorf("%w: bucket count %d + %d hits", ErrOverflow, h.counts[slot], hits)
		}
		h.counts[slot] += hits
	case stamp < t: // an older lap of the ring, or never used
		h.ticks[slot] = t
		h.counts[slot] = hits
	default: // a newer tick holds the slot; this hit is outside its window
	}
	return nil
}

// Count returns the number of hits whose tick lies in the window of B ticks ending with now's tick. O(B).
// Errors: ErrTimeRange when now is outside the UnixNano range; ErrOverflow when the total would exceed
// math.MaxInt64; ErrUninitialized on the zero value.
func (h *HitCounter) Count(now time.Time) (int64, error) {
	if h.counts == nil {
		return 0, ErrUninitialized
	}
	nowTick, err := h.tick(now)
	if err != nil {
		return 0, err
	}
	buckets := uint64(len(h.ticks))
	var total int64
	for i, stamp := range h.ticks {
		c := h.counts[i]
		if c == 0 || stamp > nowTick {
			continue
		}
		// nowTick − stamp ≥ 0 computed in uint64 is exact even when the int64 difference overflows.
		if uint64(nowTick)-uint64(stamp) >= buckets {
			continue
		}
		if c > math.MaxInt64-total {
			return 0, fmt.Errorf("%w: total hits", ErrOverflow)
		}
		total += c
	}
	return total, nil
}

// tick returns ⌊now.UnixNano() / resolution⌋ (floor, also for times before 1970).
func (h *HitCounter) tick(now time.Time) (int64, error) {
	if now.Before(time.Unix(0, math.MinInt64)) || now.After(time.Unix(0, math.MaxInt64)) {
		return 0, fmt.Errorf("%w: now = %v", ErrTimeRange, now)
	}
	ns := now.UnixNano()
	t := ns / h.resolutionNanos
	if ns%h.resolutionNanos != 0 && ns < 0 {
		t--
	}
	return t, nil
}

// slot returns tick mod B in [0, B).
func (h *HitCounter) slot(tick int64) int {
	b := int64(len(h.ticks))
	s := tick % b
	if s < 0 {
		s += b
	}
	return int(s)
}

package stream

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
	"time"
)

// t0 is an arbitrary fixed instant; tests count seconds from it.
var t0 = time.Unix(1_700_000_000, 0).UTC()

func at(seconds float64) time.Time {
	return t0.Add(time.Duration(seconds * float64(time.Second)))
}

func newHC(t *testing.T, window, resolution time.Duration) *HitCounter {
	t.Helper()
	h, err := NewHitCounter(window, resolution)
	if err != nil {
		t.Fatal(err)
	}
	return h
}

func mustCount(t *testing.T, h *HitCounter, now time.Time) int64 {
	t.Helper()
	c, err := h.Count(now)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

// Hand-derived, window 300 s at 1 s resolution: hits at seconds 1, 2, 3 → 3 counted at second 4;
// three more at second 300 → 6 at second 300; at second 301, second 1 has left → 5.
func TestHitCounter_example(t *testing.T) {
	h := newHC(t, 300*time.Second, time.Second)
	for _, s := range []float64{1, 2, 3} {
		if err := h.Hit(at(s)); err != nil {
			t.Fatal(err)
		}
	}
	if got := mustCount(t, h, at(4)); got != 3 {
		t.Errorf("Count(4) = %d, want 3", got)
	}
	if err := h.Add(at(300), 3); err != nil {
		t.Fatal(err)
	}
	if got := mustCount(t, h, at(300)); got != 6 {
		t.Errorf("Count(300) = %d, want 6", got)
	}
	if got := mustCount(t, h, at(301)); got != 5 {
		t.Errorf("Count(301) = %d, want 5", got)
	}
	if got := mustCount(t, h, at(1000)); got != 0 {
		t.Errorf("Count(1000) = %d, want 0", got)
	}
}

// With 5 s buckets over a 10 s window, a hit at 4.9 s shares bucket 0 (0–5 s): it is counted at 9.9 s
// (bucket 1) and gone at 10.0 s (bucket 2), 5.1 s later. The resolution bounds the window's error.
func TestHitCounter_resolution(t *testing.T) {
	h := newHC(t, 10*time.Second, 5*time.Second)
	base := time.Unix(0, 0) // tick boundaries at multiples of 5 s since the epoch
	_ = h.Hit(base.Add(4900 * time.Millisecond))
	if got := mustCount(t, h, base.Add(9900*time.Millisecond)); got != 1 {
		t.Errorf("Count(9.9 s) = %d, want 1", got)
	}
	if got := mustCount(t, h, base.Add(10*time.Second)); got != 0 {
		t.Errorf("Count(10 s) = %d, want 0", got)
	}
}

func TestHitCounter_beforeEpoch(t *testing.T) {
	h := newHC(t, 2*time.Second, time.Second)
	_ = h.Hit(time.Unix(0, -500_000_000)) // tick −1 (floor), not 0
	if got := mustCount(t, h, time.Unix(0, 200_000_000)); got != 1 {
		t.Errorf("Count(0.2 s) = %d, want 1", got)
	}
	if got := mustCount(t, h, time.Unix(1, 0)); got != 0 {
		t.Errorf("Count(1 s) = %d, want 0 (tick −1 is two ticks back)", got)
	}
}

func TestHitCounter_futureAndStaleHits(t *testing.T) {
	h := newHC(t, 3*time.Second, time.Second)
	_ = h.Hit(at(10))
	if got := mustCount(t, h, at(9)); got != 0 {
		t.Errorf("a hit after now was counted: %d", got)
	}
	_ = h.Hit(at(7)) // same slot as 10, older lap: dropped
	_ = h.Hit(at(8)) // an empty slot, inside the window of 10
	if got := mustCount(t, h, at(10)); got != 2 {
		t.Errorf("Count(10) = %d, want 2 (10 and 8; 7 dropped)", got)
	}
}

func TestHitCounter_zeroHits(t *testing.T) {
	h := newHC(t, 3*time.Second, time.Second)
	if err := h.Add(at(1), 0); err != nil {
		t.Fatal(err)
	}
	if got := mustCount(t, h, at(1)); got != 0 {
		t.Errorf("Count = %d, want 0", got)
	}
}

func TestHitCounter_errors(t *testing.T) {
	ctor := []struct {
		window, resolution time.Duration
		want               error
	}{
		{0, time.Second, ErrSpan},
		{-time.Second, time.Second, ErrSpan},
		{time.Second, 0, ErrSpan},
		{time.Second, -time.Nanosecond, ErrSpan},
		{10 * time.Second, 3 * time.Second, ErrSpan},
		{time.Second, 2 * time.Second, ErrSpan},
		{(maxBuckets + 1) * time.Nanosecond, time.Nanosecond, ErrBuckets},
		{math.MaxInt64, 1, ErrBuckets},
	}
	for _, c := range ctor {
		if h, err := NewHitCounter(c.window, c.resolution); !errors.Is(err, c.want) || h != nil {
			t.Errorf("NewHitCounter(%v, %v) = (%v, %v), want (nil, %v)", c.window, c.resolution, h, err, c.want)
		}
	}
	h := newHC(t, 3*time.Second, time.Second)
	if err := h.Add(at(1), -1); !errors.Is(err, ErrHits) {
		t.Errorf("Add(-1) error = %v, want ErrHits", err)
	}
	far := time.Unix(1<<40, 0)
	if err := h.Hit(far); !errors.Is(err, ErrTimeRange) {
		t.Errorf("Hit(year far beyond 2262) error = %v, want ErrTimeRange", err)
	}
	if _, err := h.Count(far); !errors.Is(err, ErrTimeRange) {
		t.Errorf("Count(year far beyond 2262) error = %v, want ErrTimeRange", err)
	}
	if err := h.Hit(time.Unix(-(1 << 40), 0)); !errors.Is(err, ErrTimeRange) {
		t.Errorf("Hit(year far before 1678) error = %v, want ErrTimeRange", err)
	}
	if err := h.Add(at(1), math.MaxInt64); err != nil {
		t.Fatal(err)
	}
	if err := h.Hit(at(1)); !errors.Is(err, ErrOverflow) {
		t.Errorf("bucket beyond MaxInt64: error = %v, want ErrOverflow", err)
	}
	if got := mustCount(t, h, at(1)); got != math.MaxInt64 {
		t.Errorf("a rejected Add changed the count: %d", got)
	}
	_ = h.Hit(at(2))
	if _, err := h.Count(at(2)); !errors.Is(err, ErrOverflow) {
		t.Errorf("total beyond MaxInt64: error = %v, want ErrOverflow", err)
	}
	var zero HitCounter
	if err := zero.Hit(at(1)); !errors.Is(err, ErrUninitialized) {
		t.Errorf("zero value Hit error = %v, want ErrUninitialized", err)
	}
	if _, err := zero.Count(at(1)); !errors.Is(err, ErrUninitialized) {
		t.Errorf("zero value Count error = %v, want ErrUninitialized", err)
	}
}

// The extremes of UnixNano are accepted at the finest resolution.
func TestHitCounter_extremeTimes(t *testing.T) {
	h := newHC(t, 4*time.Nanosecond, time.Nanosecond)
	lo, hi := time.Unix(0, math.MinInt64), time.Unix(0, math.MaxInt64)
	if err := h.Hit(lo); err != nil {
		t.Fatal(err)
	}
	if err := h.Hit(hi); err != nil {
		t.Fatal(err)
	}
	if got := mustCount(t, h, hi); got != 1 {
		t.Errorf("Count(max) = %d, want 1", got)
	}
	if got := mustCount(t, h, lo); got != 1 {
		t.Errorf("Count(min) = %d, want 1", got)
	}
}

// Oracle: a list of every hit time, counted directly, over non-decreasing random timelines.
func TestHitCounter_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		resolution := time.Duration(1+rng.IntN(5)) * 100 * time.Millisecond
		buckets := 1 + rng.IntN(20)
		window := time.Duration(buckets) * resolution
		h := newHC(t, window, resolution)
		var hits []int64 // ticks, one per hit
		now := t0
		for range 200 {
			now = now.Add(time.Duration(rng.IntN(400)) * time.Millisecond)
			tick := now.UnixNano() / int64(resolution)
			if rng.IntN(3) > 0 {
				n := int64(rng.IntN(3))
				if err := h.Add(now, n); err != nil {
					t.Fatal(err)
				}
				for range n {
					hits = append(hits, tick)
				}
				continue
			}
			var want int64
			for _, x := range hits {
				if x <= tick && tick-x < int64(buckets) {
					want++
				}
			}
			if got := mustCount(t, h, now); got != want {
				t.Fatalf("window %v, resolution %v: Count = %d, want %d", window, resolution, got, want)
			}
		}
	}
}

func TestHitCounter_noAllocation(t *testing.T) {
	h := newHC(t, 300*time.Second, time.Second)
	allocs := testing.AllocsPerRun(20, func() {
		_ = h.Hit(at(5))
		_, _ = h.Count(at(6))
	})
	if allocs != 0 {
		t.Errorf("Hit and Count allocate %v times, want 0", allocs)
	}
}

// Sizes are bucket counts: Count is O(B), Hit O(1).
func BenchmarkHitCounter(b *testing.B) {
	for _, size := range benchSizes {
		b.Run(size.name, func(b *testing.B) {
			h, _ := NewHitCounter(time.Duration(size.n)*time.Millisecond, time.Millisecond)
			now := t0
			b.ReportAllocs()
			for b.Loop() {
				now = now.Add(time.Millisecond)
				_ = h.Hit(now)
				_, _ = h.Count(now)
			}
		})
	}
}

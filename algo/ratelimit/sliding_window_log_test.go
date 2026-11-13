package ratelimit

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
	"time"
)

func mustSlidingWindowLog(t testing.TB, limit int64, window time.Duration) *SlidingWindowLog {
	t.Helper()
	s, err := NewSlidingWindowLog(limit, window)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

// Hand-derived from the moving window scheme of Berger, Eckberg, Hou, Lucantoni, GLOBECOM '90, §4 p. 351
// (at most limit arrivals in any window of length W; each arrival time is stored for W): limit 3 per 1 s,
// the same arrivals as TestFixedWindow_jumpingWindowExample, then one more at 1.1 s.
func TestSlidingWindowLog_movingWindowExample(t *testing.T) {
	s := mustSlidingWindowLog(t, 3, time.Second)
	runSteps(t, func(now time.Time, n int64) (bool, time.Duration) { return s.Allow(now, n) }, []step{
		{0, 1, decision{true, 0}},
		{200 * time.Millisecond, 1, decision{true, 0}},
		{900 * time.Millisecond, 1, decision{true, 0}},
		{950 * time.Millisecond, 1, decision{false, 50 * time.Millisecond}},   // the 0 s arrival leaves at 1 s
		{time.Second, 1, decision{true, 0}},                                   // it has left (windows are half-open)
		{1100 * time.Millisecond, 1, decision{false, 100 * time.Millisecond}}, // the 0.2 s arrival leaves at 1.2 s
	})
}

func TestSlidingWindowLog_cases(t *testing.T) {
	allowOf := func(s *SlidingWindowLog) func(time.Time, int64) (bool, time.Duration) {
		return func(now time.Time, n int64) (bool, time.Duration) { return s.Allow(now, n) }
	}
	t.Run("no burst across a window edge", func(t *testing.T) {
		// The arrivals that FixedWindow lets through as 5 units in 10 ms.
		s := mustSlidingWindowLog(t, 3, time.Second)
		runSteps(t, allowOf(s), []step{
			{0, 1, decision{true, 0}},
			{990 * time.Millisecond, 2, decision{true, 0}},
			{time.Second, 3, decision{false, 990 * time.Millisecond}}, // fits once the 0.99 s pair leaves
		})
	})
	t.Run("a large request waits for several entries", func(t *testing.T) {
		s := mustSlidingWindowLog(t, 5, time.Second)
		runSteps(t, allowOf(s), []step{
			{0, 2, decision{true, 0}},
			{100 * time.Millisecond, 2, decision{true, 0}},
			{200 * time.Millisecond, 1, decision{true, 0}},
			{500 * time.Millisecond, 3, decision{false, 600 * time.Millisecond}}, // needs both of the first two
			{1100 * time.Millisecond, 3, decision{true, 0}},
		})
	})
	t.Run("duplicates at one instant", func(t *testing.T) {
		s := mustSlidingWindowLog(t, 2, time.Second)
		runSteps(t, allowOf(s), []step{
			{0, 1, decision{true, 0}},
			{0, 1, decision{true, 0}},
			{0, 1, decision{false, time.Second}},
			{time.Second, 2, decision{true, 0}},
		})
	})
	t.Run("clock going back counts as the latest time", func(t *testing.T) {
		s := mustSlidingWindowLog(t, 1, time.Second)
		runSteps(t, allowOf(s), []step{
			{2 * time.Second, 1, decision{true, 0}},
			{500 * time.Millisecond, 1, decision{false, time.Second}}, // as at 2 s
			{2500 * time.Millisecond, 1, decision{false, 500 * time.Millisecond}},
			{3 * time.Second, 1, decision{true, 0}},
		})
	})
	t.Run("n zero, negative, above the limit", func(t *testing.T) {
		s := mustSlidingWindowLog(t, 2, time.Second)
		runSteps(t, allowOf(s), []step{
			{0, 0, decision{true, 0}},
			{0, -5, decision{true, 0}},
			{0, 3, decision{false, Never}},
			{0, math.MaxInt64, decision{false, Never}},
			{0, 2, decision{true, 0}},
		})
	})
	t.Run("zero value refuses", func(t *testing.T) {
		var s SlidingWindowLog
		if ok, retry := s.Allow(t0, 1); ok || retry != Never {
			t.Errorf("zero value Allow = (%v, %v), want (false, Never)", ok, retry)
		}
	})
	t.Run("from the zero time to the end of time", func(t *testing.T) {
		s := mustSlidingWindowLog(t, 1, time.Nanosecond)
		if ok, _ := s.Allow(time.Time{}, 1); !ok {
			t.Fatal("first call refused")
		}
		if ok, _ := s.Allow(endOfTime, 1); !ok {
			t.Error("call billions of years later refused")
		}
	})
}

func TestSlidingWindowLog_errors(t *testing.T) {
	tests := []struct {
		name    string
		limit   int64
		window  time.Duration
		wantErr error
	}{
		{"limit zero", 0, time.Second, ErrLimit},
		{"limit negative", -1, time.Second, ErrLimit},
		{"window zero", 1, 0, ErrWindow},
		{"window negative", 1, math.MinInt64, ErrWindow},
		{"limit above the cap", maxLogEntries + 1, time.Second, ErrTooLarge},
		{"limit at the maximum of int64", math.MaxInt64, time.Second, ErrTooLarge},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := NewSlidingWindowLog(tt.limit, tt.window)
			if s != nil || !errors.Is(err, tt.wantErr) {
				t.Errorf("NewSlidingWindowLog(%d, %v) = (%v, %v), want (nil, %v)", tt.limit, tt.window, s, err, tt.wantErr)
			}
		})
	}
}

// Oracle: every accepted request kept in a plain slice; the units in (t − W, t] are summed from scratch at
// each call (t clamped to the latest time), and the wait is found by trying each accepted request's
// leaving time in order.
func TestSlidingWindowLog_oracle(t *testing.T) {
	type accepted struct {
		at time.Duration
		n  int64
	}
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		limit := 1 + rng.Int64N(8)
		window := time.Duration(1 + rng.Int64N(1000))
		s := mustSlidingWindowLog(t, limit, window)
		var log []accepted
		started := false
		var now, latest time.Duration
		for stepNo := range 300 {
			now += time.Duration(rng.Int64N(int64(window)/3+1)) - window/20 // sometimes backwards
			n := rng.Int64N(limit+3) - 1                                    // -1 … limit+1
			ok, retry := s.Allow(at(now), n)

			want := decision{true, 0}
			switch {
			case n <= 0:
			case n > limit:
				want = decision{false, Never}
			default:
				if !started {
					started, latest = true, now
				}
				latest = max(latest, now)
				// Forget requests that can never count again (keeps the oracle O(limit²) per call).
				kept := log[:0]
				for _, a := range log {
					if latest-a.at < window {
						kept = append(kept, a)
					}
				}
				log = kept
				inWindow := func(t time.Duration) int64 {
					var sum int64
					for _, a := range log {
						if t-a.at < window && t >= a.at {
							sum += a.n
						}
					}
					return sum
				}
				if inWindow(latest)+n <= limit {
					log = append(log, accepted{latest, n})
				} else {
					want = decision{false, Never}
					for _, a := range log {
						if d := a.at + window - latest; d > 0 && inWindow(latest+d)+n <= limit {
							want = decision{false, d}
							break
						}
					}
				}
			}
			if (decision{ok, retry}) != want {
				t.Fatalf("trial %d step %d (limit %d, W=%v): Allow(t=%v, n=%d) = (%v, %v), want %v", trial, stepNo, limit, window, now, n, ok, retry, want)
			}
		}
	}
}

func TestSlidingWindowLog_noAllocation(t *testing.T) {
	s := mustSlidingWindowLog(t, 10, time.Second)
	i := 0
	allocs := testing.AllocsPerRun(100, func() {
		i++
		s.Allow(at(time.Duration(i)*30*time.Millisecond), 1)
	})
	if allocs != 0 {
		t.Errorf("Allow allocates %v times per call, want 0", allocs)
	}
}

// The size is the limit (entries held); each iteration makes 1000 calls at a rate slightly above the
// limit, so the log stays full and every call evicts. Flat ns/op across sizes shows O(1) amortized.
func BenchmarkSlidingWindowLog(b *testing.B) {
	for _, size := range benchSizes {
		s := mustSlidingWindowLog(b, int64(size.n), time.Second)
		gap := time.Second / time.Duration(size.n) * 9 / 10
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			now := t0
			for b.Loop() {
				for range 1000 {
					now = now.Add(gap)
					s.Allow(now, 1)
				}
			}
		})
	}
}

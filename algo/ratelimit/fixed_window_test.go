package ratelimit

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
	"time"
)

func mustFixedWindow(t testing.TB, limit int64, window time.Duration) *FixedWindow {
	t.Helper()
	w, err := NewFixedWindow(limit, window)
	if err != nil {
		t.Fatal(err)
	}
	return w
}

type step struct {
	now  time.Duration
	n    int64
	want decision
}

// Hand-derived from the jumping window scheme of Berger, Eckberg, Hou, Lucantoni, GLOBECOM '90, §4 p. 351
// (the counter starts at zero at the beginning of each window and the next window starts immediately
// after the preceding one): limit 3 per 1 s window, the first window opening at the first call.
func TestFixedWindow_jumpingWindowExample(t *testing.T) {
	w := mustFixedWindow(t, 3, time.Second)
	steps := []step{
		{0, 1, decision{true, 0}},
		{200 * time.Millisecond, 1, decision{true, 0}},
		{900 * time.Millisecond, 1, decision{true, 0}},
		{950 * time.Millisecond, 1, decision{false, 50 * time.Millisecond}}, // window full until 1 s
		{time.Second, 1, decision{true, 0}},                                 // a new window
	}
	runSteps(t, func(now time.Time, n int64) (bool, time.Duration) { return w.Allow(now, n) }, steps)
}

func runSteps(t *testing.T, allow func(time.Time, int64) (bool, time.Duration), steps []step) {
	t.Helper()
	for i, s := range steps {
		ok, retry := allow(at(s.now), s.n)
		if (decision{ok, retry}) != s.want {
			t.Fatalf("step %d: Allow(t=%v, n=%d) = (%v, %v), want %v", i, s.now, s.n, ok, retry, s.want)
		}
	}
}

func TestFixedWindow_cases(t *testing.T) {
	t.Run("up to twice the limit across a window edge", func(t *testing.T) {
		w := mustFixedWindow(t, 3, time.Second)
		runSteps(t, func(now time.Time, n int64) (bool, time.Duration) { return w.Allow(now, n) }, []step{
			{0, 1, decision{true, 0}},
			{990 * time.Millisecond, 2, decision{true, 0}},
			{time.Second, 3, decision{true, 0}}, // 5 units within 10 ms, limit 3: the scheme allows it
		})
	})
	t.Run("several empty windows are skipped", func(t *testing.T) {
		w := mustFixedWindow(t, 3, time.Second)
		runSteps(t, func(now time.Time, n int64) (bool, time.Duration) { return w.Allow(now, n) }, []step{
			{0, 3, decision{true, 0}},
			{3500 * time.Millisecond, 1, decision{true, 0}}, // window [3 s, 4 s)
			{3500 * time.Millisecond, 2, decision{true, 0}},
			{3600 * time.Millisecond, 1, decision{false, 400 * time.Millisecond}},
		})
	})
	t.Run("clock going back counts as the latest time", func(t *testing.T) {
		w := mustFixedWindow(t, 1, time.Second)
		runSteps(t, func(now time.Time, n int64) (bool, time.Duration) { return w.Allow(now, n) }, []step{
			{0, 1, decision{true, 0}},
			{1500 * time.Millisecond, 1, decision{true, 0}},
			{500 * time.Millisecond, 1, decision{false, 500 * time.Millisecond}}, // as at 1.5 s
		})
	})
	t.Run("n zero, negative, above the limit", func(t *testing.T) {
		w := mustFixedWindow(t, 2, time.Second)
		runSteps(t, func(now time.Time, n int64) (bool, time.Duration) { return w.Allow(now, n) }, []step{
			{0, 0, decision{true, 0}},
			{0, math.MinInt64, decision{true, 0}},
			{0, 3, decision{false, Never}},
			{0, math.MaxInt64, decision{false, Never}},
			{0, 2, decision{true, 0}},
		})
	})
	t.Run("largest limit", func(t *testing.T) {
		w := mustFixedWindow(t, math.MaxInt64, time.Second)
		runSteps(t, func(now time.Time, n int64) (bool, time.Duration) { return w.Allow(now, n) }, []step{
			{0, math.MaxInt64 - 1, decision{true, 0}},
			{0, 2, decision{false, time.Second}},
			{0, 1, decision{true, 0}},
		})
	})
	t.Run("zero value refuses", func(t *testing.T) {
		var w FixedWindow
		if ok, retry := w.Allow(t0, 1); ok || retry != Never {
			t.Errorf("zero value Allow = (%v, %v), want (false, Never)", ok, retry)
		}
	})
	t.Run("from the zero time to the end of time", func(t *testing.T) {
		w := mustFixedWindow(t, 1, time.Nanosecond)
		if ok, _ := w.Allow(time.Time{}, 1); !ok {
			t.Fatal("first call refused")
		}
		if ok, _ := w.Allow(endOfTime, 1); !ok {
			t.Error("call billions of years later refused")
		}
	})
}

func TestFixedWindow_errors(t *testing.T) {
	tests := []struct {
		name    string
		limit   int64
		window  time.Duration
		wantErr error
	}{
		{"limit zero", 0, time.Second, ErrLimit},
		{"limit negative", math.MinInt64, time.Second, ErrLimit},
		{"window zero", 1, 0, ErrWindow},
		{"window negative", 1, -time.Second, ErrWindow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := NewFixedWindow(tt.limit, tt.window)
			if w != nil || !errors.Is(err, tt.wantErr) {
				t.Errorf("NewFixedWindow(%d, %v) = (%v, %v), want (nil, %v)", tt.limit, tt.window, w, err, tt.wantErr)
			}
		})
	}
}

// Oracle: window index ⌊(t − first)/W⌋ of the clamped time t, one counter per index.
func TestFixedWindow_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		limit := 1 + rng.Int64N(6)
		window := time.Duration(1 + rng.Int64N(1000))
		w := mustFixedWindow(t, limit, window)
		started := false
		var now, first, latest time.Duration
		var curIdx, count int64
		for stepNo := range 300 {
			now += time.Duration(rng.Int64N(int64(window)/2+1)) - window/10 // sometimes backwards
			n := rng.Int64N(limit+3) - 1                                    // -1 … limit+1
			ok, retry := w.Allow(at(now), n)

			want := decision{true, 0}
			switch {
			case n <= 0:
			case n > limit:
				want = decision{false, Never}
			default:
				if !started {
					started, first, latest = true, now, now
				}
				latest = max(latest, now)
				idx := int64((latest - first) / window)
				if idx != curIdx {
					curIdx, count = idx, 0
				}
				if count+n <= limit {
					count += n
				} else {
					want = decision{false, time.Duration(idx+1)*window - (latest - first)}
				}
			}
			if (decision{ok, retry}) != want {
				t.Fatalf("trial %d step %d (limit %d, W=%v): Allow(t=%v, n=%d) = (%v, %v), want %v", trial, stepNo, limit, window, now, n, ok, retry, want)
			}
		}
	}
}

func TestFixedWindow_noAllocation(t *testing.T) {
	w := mustFixedWindow(t, 10, time.Second)
	i := 0
	allocs := testing.AllocsPerRun(100, func() {
		i++
		w.Allow(at(time.Duration(i)*30*time.Millisecond), 1)
	})
	if allocs != 0 {
		t.Errorf("Allow allocates %v times per call, want 0", allocs)
	}
}

// One call is O(1); the sizes are calls per iteration, so ns/op grows linearly with the size.
func BenchmarkFixedWindow(b *testing.B) {
	for _, size := range benchSizes {
		w := mustFixedWindow(b, 100, time.Second)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			now := t0
			for b.Loop() {
				for range size.n {
					now = now.Add(7 * time.Millisecond)
					w.Allow(now, 1)
				}
			}
		})
	}
}

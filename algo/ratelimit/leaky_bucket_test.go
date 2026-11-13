package ratelimit

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
	"time"
)

func mustLeakyBucket(t testing.TB, emission, tolerance time.Duration) *LeakyBucket {
	t.Helper()
	b, err := NewLeakyBucket(emission, tolerance)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Hand-derived from ITU-T I.371 Annex A.1 (virtual scheduling): T = 100 ms, τ = 200 ms. At t = 0, TAT
// starts at 0 and three cells conform (TAT − t = 0, 100, 200 ms ≤ τ), moving TAT to 100, 200, 300 ms; the
// fourth is early by 300 − 200 = 100 ms; at t = 100 ms a cell conforms and TAT becomes 400 ms.
func TestLeakyBucket_i371Example(t *testing.T) {
	b := mustLeakyBucket(t, 100*time.Millisecond, 200*time.Millisecond)
	steps := []struct {
		now  time.Duration
		want decision
	}{
		{0, decision{true, 0}},
		{0, decision{true, 0}},
		{0, decision{true, 0}},
		{0, decision{false, 100 * time.Millisecond}},
		{100 * time.Millisecond, decision{true, 0}},
		{100 * time.Millisecond, decision{false, 100 * time.Millisecond}},
	}
	for i, s := range steps {
		ok, retry := b.Allow(at(s.now))
		if (decision{ok, retry}) != s.want {
			t.Fatalf("step %d: Allow(t=%v) = (%v, %v), want %v", i, s.now, ok, retry, s.want)
		}
	}
}

func TestLeakyBucket_cases(t *testing.T) {
	t.Run("zero tolerance spaces cells exactly T apart", func(t *testing.T) {
		b := mustLeakyBucket(t, time.Second, 0)
		if ok, _ := b.Allow(t0); !ok {
			t.Fatal("first cell refused")
		}
		if ok, retry := b.Allow(at(time.Second - time.Nanosecond)); ok || retry != time.Nanosecond {
			t.Errorf("cell 1 ns early = (%v, %v), want (false, 1ns)", ok, retry)
		}
		if ok, _ := b.Allow(at(time.Second)); !ok {
			t.Error("cell exactly T later refused")
		}
	})
	t.Run("a late cell restarts the schedule from its arrival", func(t *testing.T) {
		b := mustLeakyBucket(t, time.Second, 0)
		b.Allow(t0)
		b.Allow(at(10 * time.Second)) // TAT = 11 s, not 2 s
		if ok, retry := b.Allow(at(10*time.Second + 500*time.Millisecond)); ok || retry != 500*time.Millisecond {
			t.Errorf("Allow = (%v, %v), want (false, 500ms)", ok, retry)
		}
	})
	t.Run("clock going back counts as the latest time", func(t *testing.T) {
		b := mustLeakyBucket(t, time.Second, 0)
		b.Allow(at(5 * time.Second))
		if ok, retry := b.Allow(at(time.Second)); ok || retry != time.Second {
			t.Errorf("Allow at an earlier time = (%v, %v), want (false, 1s)", ok, retry)
		}
	})
	t.Run("zero value refuses", func(t *testing.T) {
		var b LeakyBucket
		if ok, retry := b.Allow(t0); ok || retry != Never {
			t.Errorf("zero value Allow = (%v, %v), want (false, Never)", ok, retry)
		}
	})
	t.Run("end of time", func(t *testing.T) {
		b := mustLeakyBucket(t, time.Second, 0)
		if ok, retry := b.Allow(endOfTime); ok || retry != Never {
			t.Errorf("Allow at the end of time = (%v, %v), want (false, Never)", ok, retry)
		}
	})
	t.Run("largest tolerance", func(t *testing.T) {
		b := mustLeakyBucket(t, time.Nanosecond, math.MaxInt64)
		for i := range 1000 {
			if ok, _ := b.Allow(t0); !ok {
				t.Fatalf("cell %d refused under the largest tolerance", i)
			}
		}
	})
}

func TestLeakyBucket_errors(t *testing.T) {
	tests := []struct {
		name                string
		emission, tolerance time.Duration
		wantErr             error
	}{
		{"emission zero", 0, 0, ErrInterval},
		{"emission negative", -time.Second, 0, ErrInterval},
		{"emission at the minimum", math.MinInt64, 0, ErrInterval},
		{"tolerance negative", time.Second, -1, ErrTolerance},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewLeakyBucket(tt.emission, tt.tolerance)
			if b != nil || !errors.Is(err, tt.wantErr) {
				t.Errorf("NewLeakyBucket(%v, %v) = (%v, %v), want (nil, %v)", tt.emission, tt.tolerance, b, err, tt.wantErr)
			}
		})
	}
}

// Oracle: the continuous-state leaky bucket of I.371 Annex A.2, which the NOTE there says decides the same
// cells as A.1. X′ = X − (ta − LCT); conforming iff X′ ≤ τ, then X = max(X′, 0) + T and LCT = ta. The wait
// of a refused cell is X′ − τ (TAT = X + LCT, so TAT − ta = X′).
func TestLeakyBucket_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		emission := time.Duration(1 + rng.Int64N(1000))
		tolerance := time.Duration(rng.Int64N(int64(emission) * 5))
		b := mustLeakyBucket(t, emission, tolerance)
		var now, latest, x, lct time.Duration
		for step := range 300 {
			now += time.Duration(rng.Int64N(int64(emission)*2+1)) - emission/3 // sometimes backwards
			ok, retry := b.Allow(at(now))

			ta := now
			if step == 0 {
				x, lct, latest = 0, ta, ta
			} else if ta < latest {
				ta = latest
			}
			latest = ta
			want := decision{true, 0}
			if xp := x - (ta - lct); xp <= tolerance {
				x, lct = max(xp, 0)+emission, ta
			} else {
				want = decision{false, xp - tolerance}
			}
			if (decision{ok, retry}) != want {
				t.Fatalf("trial %d step %d (T=%v, τ=%v): Allow(t=%v) = (%v, %v), want %v", trial, step, emission, tolerance, now, ok, retry, want)
			}
		}
	}
}

func TestLeakyBucket_noAllocation(t *testing.T) {
	b := mustLeakyBucket(t, time.Millisecond, 5*time.Millisecond)
	i := 0
	allocs := testing.AllocsPerRun(100, func() {
		i++
		b.Allow(at(time.Duration(i) * 700 * time.Microsecond))
	})
	if allocs != 0 {
		t.Errorf("Allow allocates %v times per call, want 0", allocs)
	}
}

// One call is O(1); the sizes are calls per iteration, so ns/op grows linearly with the size.
func BenchmarkLeakyBucket(b *testing.B) {
	for _, size := range benchSizes {
		lb := mustLeakyBucket(b, time.Millisecond, 50*time.Millisecond)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			now := t0
			for b.Loop() {
				for range size.n {
					now = now.Add(700 * time.Microsecond)
					lb.Allow(now)
				}
			}
		})
	}
}

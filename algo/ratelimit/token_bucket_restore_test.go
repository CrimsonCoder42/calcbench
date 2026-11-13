package ratelimit

import (
	"math"
	"testing"
	"time"
)

// The worked example of the doc comment: burst 5, take 3 → 2 left; give 3 back → 5; give 10 back → 5.
func TestTokenBucketRestore_example(t *testing.T) {
	b := mustTokenBucket(t, 1, time.Second, 5)
	if ok, _ := b.Allow(t0, 3); !ok {
		t.Fatal("Allow(t0, 3) refused on a full bucket")
	}
	if got := b.Available(t0); got != 2 {
		t.Fatalf("after taking 3: Available = %d, want 2", got)
	}
	b.Restore(3, t0)
	if got := b.Available(t0); got != 5 {
		t.Fatalf("after Restore(3): Available = %d, want 5", got)
	}
	b.Restore(10, t0)
	if got := b.Available(t0); got != 5 {
		t.Fatalf("after Restore(10): Available = %d, want 5 (capped)", got)
	}
}

func TestTokenBucketRestore_cases(t *testing.T) {
	t.Run("capped", func(t *testing.T) {
		b := mustTokenBucket(t, 1, time.Second, 5)
		b.Allow(t0, 1)
		b.Restore(4, t0)
		if got := b.Available(t0); got != 5 {
			t.Fatalf("Available = %d, want 5", got)
		}
	})
	t.Run("after a later take", func(t *testing.T) {
		b := mustTokenBucket(t, 1, time.Second, 5)
		b.Allow(t0, 3)
		b.Allow(t0, 2)
		if got := b.Available(t0); got != 0 {
			t.Fatalf("after taking 5: Available = %d, want 0", got)
		}
		b.Restore(3, t0)
		if got := b.Available(t0); got != 3 {
			t.Fatalf("after Restore(3): Available = %d, want 3", got)
		}
	})
	t.Run("n zero and negative are no-ops", func(t *testing.T) {
		for _, n := range []int64{0, -1, math.MinInt64} {
			b := mustTokenBucket(t, 1, time.Second, 5)
			b.Allow(t0, 3)
			before := *b
			b.Restore(n, at(time.Hour)) // not even the refill happens
			if *b != before {
				t.Fatalf("Restore(%d) changed the bucket: %+v, want %+v", n, *b, before)
			}
		}
	})
	t.Run("zero value", func(t *testing.T) {
		var b TokenBucket
		b.Restore(3, t0)
		if b != (TokenBucket{}) {
			t.Fatalf("zero value changed: %+v", b)
		}
		if ok, retry := b.Allow(t0, 1); ok || retry != Never {
			t.Fatalf("zero value Allow = (%v, %v), want (false, Never)", ok, retry)
		}
	})
	t.Run("unused bucket stays full", func(t *testing.T) {
		b := mustTokenBucket(t, 1, time.Second, 5)
		b.Restore(2, t0)
		if got := b.Available(t0); got != 5 {
			t.Fatalf("Available = %d, want 5", got)
		}
		if ok, _ := b.Allow(t0, 5); !ok {
			t.Fatal("Allow(t0, 5) refused after Restore on an unused bucket")
		}
	})
	t.Run("huge n does not overflow", func(t *testing.T) {
		b := mustTokenBucket(t, 1, time.Hour, 1000)
		b.Allow(t0, 1000)
		b.Restore(math.MaxInt64, t0)
		if got := b.Available(t0); got != 1000 {
			t.Fatalf("Available = %d, want 1000", got)
		}
		if b.level != b.capacity {
			t.Fatalf("level = %d, want capacity %d", b.level, b.capacity)
		}
	})
	t.Run("restore then refill stays capped", func(t *testing.T) {
		b := mustTokenBucket(t, 1, time.Second, 5)
		b.Allow(t0, 2)
		b.Restore(2, t0)
		if got := b.Available(at(time.Hour)); got != 5 {
			t.Fatalf("Available an hour later = %d, want 5", got)
		}
	})
	t.Run("refill happens before the restore", func(t *testing.T) {
		b := mustTokenBucket(t, 1, time.Second, 5)
		b.Allow(t0, 5)
		b.Restore(1, at(2*time.Second)) // 2 refilled + 1 restored
		if got := b.Available(at(2 * time.Second)); got != 3 {
			t.Fatalf("Available = %d, want 3", got)
		}
	})
	t.Run("now earlier than latest counts as latest", func(t *testing.T) {
		b := mustTokenBucket(t, 1, time.Second, 5)
		b.Allow(at(10*time.Second), 5)
		b.Restore(2, t0) // earlier now: no time passes, no refill
		if got := b.Available(at(10 * time.Second)); got != 2 {
			t.Fatalf("Available = %d, want 2", got)
		}
		if !b.last.Equal(at(10 * time.Second)) {
			t.Fatalf("last = %v, want %v", b.last, at(10*time.Second))
		}
	})
	t.Run("Allow semantics unchanged after Restore", func(t *testing.T) {
		// 2 tokens per second, burst 3: after a restore the bucket meters exactly as one that never lent.
		got := mustTokenBucket(t, 2, time.Second, 3)
		want := mustTokenBucket(t, 2, time.Second, 3)
		got.Allow(t0, 2)
		got.Restore(2, t0)
		want.Allow(t0, 0)
		steps := []struct {
			now time.Duration
			n   int64
		}{{0, 3}, {0, 1}, {250 * time.Millisecond, 1}, {500 * time.Millisecond, 1}, {2 * time.Second, 3}, {2 * time.Second, 4}}
		for i, s := range steps {
			okG, rG := got.Allow(at(s.now), s.n)
			okW, rW := want.Allow(at(s.now), s.n)
			if okG != okW || rG != rW {
				t.Fatalf("step %d: Allow(t=%v, n=%d) = (%v, %v), want (%v, %v)", i, s.now, s.n, okG, rG, okW, rW)
			}
		}
	})
	t.Run("fractional level kept", func(t *testing.T) {
		// Half a token present (500 ms at 1/s) plus 1 restored is 1.5: one whole token, and the half stays.
		b := mustTokenBucket(t, 1, time.Second, 5)
		b.Allow(t0, 5)
		b.Restore(1, at(500*time.Millisecond))
		if got := b.Available(at(500 * time.Millisecond)); got != 1 {
			t.Fatalf("Available = %d, want 1", got)
		}
		if got := b.Available(at(time.Second)); got != 2 {
			t.Fatalf("Available 500 ms later = %d, want 2", got)
		}
	})
}

func TestTokenBucketRestore_noAllocation(t *testing.T) {
	b := mustTokenBucket(t, 1000, time.Second, 10)
	i := 0
	allocs := testing.AllocsPerRun(100, func() {
		i++
		now := at(time.Duration(i) * time.Millisecond)
		b.Allow(now, 2)
		b.Restore(1, now)
	})
	if allocs != 0 {
		t.Errorf("Restore allocates %v times per call, want 0", allocs)
	}
}

// One call is O(1); the sizes are calls per iteration, so ns/op grows linearly with the size.
func BenchmarkTokenBucketRestore(b *testing.B) {
	for _, size := range benchSizes {
		tb := mustTokenBucket(b, 1000, time.Second, 100)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			now := t0
			for b.Loop() {
				for range size.n {
					now = now.Add(700 * time.Microsecond)
					tb.Allow(now, 2)   // take 2 so that giving 1 back takes the add path, not the cap
					tb.Restore(1, now) // level + 1 token < capacity
				}
			}
		})
	}
}

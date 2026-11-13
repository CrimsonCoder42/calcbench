package ratelimit

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
	"time"
)

// t0 is the reference instant of every test: wall clock only (no monotonic reading), so arithmetic is exact.
var t0 = time.Date(2026, 11, 12, 0, 0, 0, 0, time.UTC)

// at is t0 + d.
func at(d time.Duration) time.Time { return t0.Add(d) }

// endOfTime is the latest representable time.Time: Add of any positive duration saturates there.
var endOfTime = time.Unix(math.MaxInt64-62135596800, 999_999_999)

// benchSizes are the collection sizes of .claude/rules/go-complexity.md. Shared by the benchmarks.
var benchSizes = []struct {
	name string
	n    int
}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}}

type decision struct {
	ok    bool
	retry time.Duration
}

func mustTokenBucket(t testing.TB, tokensPerInterval int64, interval time.Duration, burst int64) *TokenBucket {
	t.Helper()
	b, err := NewTokenBucket(tokensPerInterval, interval, burst)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Hand-derived from RFC 2697 §3 (the committed bucket alone: starts full at CBS, refills at CIR up to CBS,
// a request of B conforms when Tc − B ≥ 0): 2 tokens per second, burst 3.
func TestTokenBucket_rfc2697Example(t *testing.T) {
	b := mustTokenBucket(t, 2, time.Second, 3)
	steps := []struct {
		now  time.Duration
		n    int64
		want decision
	}{
		{0, 3, decision{true, 0}},                       // starts full: all three
		{0, 1, decision{false, 500 * time.Millisecond}}, // empty: one token every 500 ms
		{500 * time.Millisecond, 1, decision{true, 0}},
		{500 * time.Millisecond, 2, decision{false, time.Second}},
	}
	for i, s := range steps {
		ok, retry := b.Allow(at(s.now), s.n)
		if (decision{ok, retry}) != s.want {
			t.Fatalf("step %d: Allow(t=%v, n=%d) = (%v, %v), want %v", i, s.now, s.n, ok, retry, s.want)
		}
	}
	if got := b.Available(at(10 * time.Second)); got != 3 {
		t.Errorf("Available after 10 s = %d, want 3 (capped at burst)", got)
	}
}

func TestTokenBucket_cases(t *testing.T) {
	t.Run("n zero or negative conforms and takes nothing", func(t *testing.T) {
		b := mustTokenBucket(t, 1, time.Second, 1)
		for _, n := range []int64{0, -1, math.MinInt64} {
			if ok, retry := b.Allow(t0, n); !ok || retry != 0 {
				t.Errorf("Allow(n=%d) = (%v, %v), want (true, 0)", n, ok, retry)
			}
		}
		if got := b.Available(t0); got != 1 {
			t.Errorf("Available = %d, want 1", got)
		}
	})
	t.Run("n above burst never conforms", func(t *testing.T) {
		b := mustTokenBucket(t, 1, time.Second, 2)
		for _, n := range []int64{3, math.MaxInt64} {
			if ok, retry := b.Allow(t0, n); ok || retry != Never {
				t.Errorf("Allow(n=%d) = (%v, %v), want (false, Never)", n, ok, retry)
			}
		}
		if got := b.Available(t0); got != 2 {
			t.Errorf("Available = %d, want 2 (nothing taken)", got)
		}
	})
	t.Run("zero value refuses", func(t *testing.T) {
		var b TokenBucket
		if ok, retry := b.Allow(t0, 1); ok || retry != Never {
			t.Errorf("zero value Allow = (%v, %v), want (false, Never)", ok, retry)
		}
		if got := b.Available(t0); got != 0 {
			t.Errorf("zero value Available = %d, want 0", got)
		}
	})
	t.Run("retry rounds up to the nanosecond", func(t *testing.T) {
		b := mustTokenBucket(t, 3, time.Second, 1) // one token every 333 333 333.3… ns
		b.Allow(t0, 1)
		ok, retry := b.Allow(t0, 1)
		if ok || retry != 333_333_334 {
			t.Fatalf("Allow = (%v, %v), want (false, 333.333334ms)", ok, retry)
		}
		if ok, _ := b.Allow(at(retry-1), 1); ok {
			t.Errorf("Allow 1 ns before retryAfter conformed")
		}
		if ok, _ := b.Allow(at(retry), 1); !ok {
			t.Errorf("Allow at retryAfter refused")
		}
	})
	t.Run("no drift over a long run", func(t *testing.T) {
		// 7 tokens per 3 ns, burst 7: taking 7 every 3 ns conforms forever; a float level would drift.
		b := mustTokenBucket(t, 7, 3*time.Nanosecond, 7)
		for i := range 100_000 {
			if ok, retry := b.Allow(at(time.Duration(3*i)), 7); !ok {
				t.Fatalf("call %d refused, retry %v", i, retry)
			}
		}
	})
	t.Run("clock going back adds no tokens", func(t *testing.T) {
		b := mustTokenBucket(t, 1, time.Second, 1)
		b.Allow(at(10*time.Second), 1)
		if ok, retry := b.Allow(at(5*time.Second), 1); ok || retry != time.Second {
			t.Errorf("Allow at an earlier time = (%v, %v), want (false, 1s)", ok, retry)
		}
		// Time is measured from the latest now (10 s), not from the earlier one.
		if ok, _ := b.Allow(at(10*time.Second+999*time.Millisecond), 1); ok {
			t.Errorf("Allow 999 ms after the latest now conformed")
		}
		if ok, _ := b.Allow(at(11*time.Second), 1); !ok {
			t.Errorf("Allow 1 s after the latest now refused")
		}
	})
	t.Run("centuries elapsed saturate to full", func(t *testing.T) {
		b := mustTokenBucket(t, math.MaxInt64, time.Nanosecond, 1)
		b.Allow(time.Time{}, 1)
		if got := b.Available(endOfTime); got != 1 {
			t.Errorf("Available at the end of time = %d, want 1", got)
		}
	})
	t.Run("largest capacity", func(t *testing.T) {
		b := mustTokenBucket(t, math.MaxInt64, time.Nanosecond, math.MaxInt64)
		if ok, _ := b.Allow(t0, math.MaxInt64); !ok {
			t.Fatal("full bucket refused its burst")
		}
		if ok, retry := b.Allow(t0, math.MaxInt64); ok || retry != 1 {
			t.Errorf("empty bucket Allow(MaxInt64) = (%v, %v), want (false, 1ns)", ok, retry)
		}
		if got := b.Available(at(time.Nanosecond)); got != math.MaxInt64 {
			t.Errorf("Available after 1 ns = %d, want MaxInt64", got)
		}
	})
}

func TestTokenBucket_errors(t *testing.T) {
	tests := []struct {
		name     string
		rate     int64
		interval time.Duration
		burst    int64
		wantErr  error
	}{
		{"rate zero", 0, time.Second, 1, ErrRate},
		{"rate negative", -1, time.Second, 1, ErrRate},
		{"interval zero", 1, 0, 1, ErrInterval},
		{"interval negative", 1, -time.Second, 1, ErrInterval},
		{"burst zero", 1, time.Second, 0, ErrBurst},
		{"burst at the minimum of int64", 1, time.Second, math.MinInt64, ErrBurst},
		{"burst × interval overflows", 1, 2 * time.Nanosecond, math.MaxInt64, ErrTooLarge},
		{"a year of a billion tokens", 1, 365 * 24 * time.Hour, 1_000_000_000, ErrTooLarge},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewTokenBucket(tt.rate, tt.interval, tt.burst)
			if b != nil || !errors.Is(err, tt.wantErr) {
				t.Errorf("NewTokenBucket(%d, %v, %d) = (%v, %v), want (nil, %v)", tt.rate, tt.interval, tt.burst, b, err, tt.wantErr)
			}
		})
	}
}

// Oracle: an exact rational token bucket. tokens = min(burst, tokens + elapsed·rate/interval); a request
// conforms when tokens ≥ n; the wait is ⌈(n − tokens)·interval/rate⌉ ns.
func TestTokenBucket_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 200 {
		rate := 1 + rng.Int64N(10)
		interval := time.Duration(1 + rng.Int64N(1000))
		burst := 1 + rng.Int64N(8)
		b := mustTokenBucket(t, rate, interval, burst)
		tokens := new(big.Rat).SetInt64(burst)
		capRat := new(big.Rat).SetInt64(burst)
		var now, latest time.Duration
		started := false
		for step := range 200 {
			now += time.Duration(rng.Int64N(int64(interval)*2+1)) - interval/4 // sometimes backwards
			n := rng.Int64N(burst+2) - 1                                       // -1 … burst
			ok, retry := b.Allow(at(now), n)

			want := decision{true, 0}
			if n > 0 { // n <= 0 touches no state
				switch {
				case !started:
					started, latest = true, now
				case now > latest:
					tokens.Add(tokens, new(big.Rat).SetFrac64(int64(now-latest)*rate, int64(interval)))
					if tokens.Cmp(capRat) > 0 {
						tokens.Set(capRat)
					}
					latest = now
				}
				need := new(big.Rat).SetInt64(n)
				if tokens.Cmp(need) >= 0 {
					tokens.Sub(tokens, need)
				} else {
					d := new(big.Rat).Sub(need, tokens)
					d.Mul(d, new(big.Rat).SetFrac64(int64(interval), rate))
					q, r := new(big.Int).QuoRem(d.Num(), d.Denom(), new(big.Int))
					if r.Sign() != 0 {
						q.Add(q, big.NewInt(1))
					}
					want = decision{false, time.Duration(q.Int64())}
				}
			}
			if (decision{ok, retry}) != want {
				t.Fatalf("trial %d step %d (rate %d per %v, burst %d): Allow(t=%v, n=%d) = (%v, %v), want %v",
					trial, step, rate, interval, burst, now, n, ok, retry, want)
			}
		}
	}
}

// TokenBucket(1 per T, burst b) and GCRA(T, (b−1)·T) decide the same requests with the same waits.
// Proof (I.371 A.2): with bucket content X′ = X − (ta − LCT) at an arrival, the token count is
// b − max(X′, 0)/T. Refill: max(X − elapsed, 0) gives min(b, tokens + elapsed/T). A cell conforms iff
// X′ ≤ (b−1)T iff tokens ≥ 1, and afterwards X = max(X′, 0) + T: one token fewer. The GCRA wait
// X′ − (b−1)T equals the token wait (1 − tokens)·T.
func TestTokenBucket_matchesGCRA(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 200 {
		emission := time.Duration(1 + rng.Int64N(1000))
		burst := 1 + rng.Int64N(6)
		tb := mustTokenBucket(t, 1, emission, burst)
		lb := mustLeakyBucket(t, emission, time.Duration(burst-1)*emission)
		var now time.Duration
		for step := range 300 {
			now += time.Duration(rng.Int64N(int64(emission)*2+1)) - emission/3
			ok1, r1 := tb.Allow(at(now), 1)
			ok2, r2 := lb.Allow(at(now))
			if ok1 != ok2 || r1 != r2 {
				t.Fatalf("trial %d step %d (T=%v, b=%d), t=%v: TokenBucket (%v, %v), LeakyBucket (%v, %v)",
					trial, step, emission, burst, now, ok1, r1, ok2, r2)
			}
		}
	}
}

func TestTokenBucket_noAllocation(t *testing.T) {
	b := mustTokenBucket(t, 1000, time.Second, 10)
	i := 0
	allocs := testing.AllocsPerRun(100, func() {
		i++
		b.Allow(at(time.Duration(i)*time.Millisecond), 1)
		b.Available(at(time.Duration(i) * time.Millisecond))
	})
	if allocs != 0 {
		t.Errorf("Allow allocates %v times per call, want 0", allocs)
	}
}

// One call is O(1); the sizes are calls per iteration, so ns/op grows linearly with the size.
func BenchmarkTokenBucket(b *testing.B) {
	for _, size := range benchSizes {
		tb := mustTokenBucket(b, 1000, time.Second, 100)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			now := t0
			for b.Loop() {
				for range size.n {
					now = now.Add(700 * time.Microsecond)
					tb.Allow(now, 1)
				}
			}
		})
	}
}

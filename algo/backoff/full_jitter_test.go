package backoff

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
	"time"
)

// Worked example: numbers from toolbench's request R-1 (doc/toolbench_requests_2026-11-13/BRIEF.md); formula
// Brooker 2015, figure 6 (Full Jitter).
func TestFullJitter_Worked(t *testing.T) {
	got, err := FullJitter(3, 100*time.Millisecond, 10*time.Second, 0.5)
	if err != nil || got != 400*time.Millisecond {
		t.Errorf("FullJitter(3, 100ms, 10s, 0.5) = %v, %v; want 400ms", got, err)
	}
}

func TestFullJitter_Cases(t *testing.T) {
	below1 := math.Nextafter(1, 0)
	cases := []struct {
		name    string
		attempt int
		base    time.Duration
		max     time.Duration
		u       float64
		want    time.Duration
	}{
		{"uZero", 3, 100 * time.Millisecond, 10 * time.Second, 0, 0},
		{"uNegativeZero", 3, 100 * time.Millisecond, 10 * time.Second, math.Copysign(0, -1), 0},
		{"uLargestBelow1", 3, 100 * time.Millisecond, 10 * time.Second, below1, 800*time.Millisecond - 1},
		{"capped", 10, 100 * time.Millisecond, 10 * time.Second, 0.25, 2500 * time.Millisecond},
		{"attempt0", 0, 100 * time.Millisecond, 10 * time.Second, 0.5, 50 * time.Millisecond},
		{"base1ns", 0, 1, 1, below1, 0},
		{"base1ns_62", 62, 1, maxDur, 0.5, 1 << 61},
		{"base1ns_63", 63, 1, maxDur, 0.5, 1<<62 - 1},
		{"base1ns_maxInt", math.MaxInt, 1, maxDur, 0.5, 1<<62 - 1},
		{"baseMaxInt64", 0, maxDur, maxDur, 0.5, 1<<62 - 1},
		{"tinyU", 3, 100 * time.Millisecond, 10 * time.Second, 0x1p-60, 0},
	}
	for _, c := range cases {
		got, err := FullJitter(c.attempt, c.base, c.max, c.u)
		if err != nil || got != c.want {
			t.Errorf("%s: FullJitter(%d, %d, %d, %v) = %d, %v; want %d", c.name, c.attempt, c.base, c.max, c.u, got, err, c.want)
		}
	}
}

// For u = the largest float64 below 1 and sleep ≤ 2⁵³ ns, the result is exactly sleep − 1; for larger sleeps it
// is still below sleep.
func TestFullJitter_TopOfRange(t *testing.T) {
	below1 := math.Nextafter(1, 0)
	for attempt := 0; attempt <= 70; attempt++ {
		sleep, _ := Exponential(attempt, 1, maxDur)
		got, err := FullJitter(attempt, 1, maxDur, below1)
		if err != nil || got >= sleep {
			t.Fatalf("attempt %d: FullJitter = %d, %v; want < %d", attempt, got, err, sleep)
		}
		if sleep <= 1<<53 && got != sleep-1 {
			t.Fatalf("attempt %d: FullJitter = %d; want sleep − 1 = %d", attempt, got, sleep-1)
		}
	}
}

func TestFullJitter_Errors(t *testing.T) {
	cases := []struct {
		name    string
		attempt int
		base    time.Duration
		max     time.Duration
		u       float64
		want    error
	}{
		{"negativeAttempt", -1, time.Second, time.Minute, 0.5, ErrAttempt},
		{"attemptFirst", -1, 0, -1, math.NaN(), ErrAttempt},
		{"zeroBase", 0, 0, time.Minute, 0.5, ErrBase},
		{"baseBeforeDraw", 0, 0, time.Minute, math.NaN(), ErrBase},
		{"maxBelowBase", 0, time.Second, time.Millisecond, 0.5, ErrMax},
		{"maxBeforeDraw", 0, time.Second, time.Millisecond, 2, ErrMax},
		{"uNaN", 0, time.Second, time.Minute, math.NaN(), ErrDraw},
		{"uNegativeTiny", 0, time.Second, time.Minute, -math.SmallestNonzeroFloat64, ErrDraw},
		{"uOne", 0, time.Second, time.Minute, 1, ErrDraw},
		{"uPlusInf", 0, time.Second, time.Minute, math.Inf(1), ErrDraw},
		{"uMinusInf", 0, time.Second, time.Minute, math.Inf(-1), ErrDraw},
	}
	for _, c := range cases {
		got, err := FullJitter(c.attempt, c.base, c.max, c.u)
		if !errors.Is(err, c.want) || got != 0 {
			t.Errorf("%s: FullJitter = %v, %v; want 0, %v", c.name, got, err, c.want)
		}
	}
}

// The result is monotone non-decreasing in u for a fixed schedule (fixed seed, 500 pairs).
func TestFullJitter_MonotoneInU(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for i := 0; i < 500; i++ {
		u1, u2 := r.Float64(), r.Float64()
		if u1 > u2 {
			u1, u2 = u2, u1
		}
		base := 1 + time.Duration(r.Int64N(1<<40))
		d1, err1 := FullJitter(i%70, base, maxDur, u1)
		d2, err2 := FullJitter(i%70, base, maxDur, u2)
		if err1 != nil || err2 != nil || d1 > d2 {
			t.Fatalf("attempt %d, base %d: u %v → %d, u %v → %d (%v, %v)", i%70, base, u1, d1, u2, d2, err1, err2)
		}
	}
}

func TestFullJitter_Allocs(t *testing.T) {
	var sink time.Duration
	n := testing.AllocsPerRun(100, func() {
		d, _ := FullJitter(7, 100*time.Millisecond, 10*time.Second, 0.37)
		sink += d
	})
	if n != 0 {
		t.Errorf("FullJitter allocates %v times per call; want 0", n)
	}
	benchSink = sink
}

// FuzzFullJitter checks 0 ≤ result < sleep (sleep = Exponential ≥ base > 0) and no error on valid input, and an
// error on invalid input. Run with -fuzztime.
func FuzzFullJitter(f *testing.F) {
	f.Add(3, int64(100*time.Millisecond), int64(10*time.Second), 0.5)
	f.Add(63, int64(1), int64(math.MaxInt64), math.Nextafter(1, 0))
	f.Add(0, int64(math.MaxInt64), int64(math.MaxInt64), 0.0)
	f.Add(0, int64(1), int64(1), math.NaN())
	f.Fuzz(func(t *testing.T, attempt int, base, max int64, u float64) {
		got, err := FullJitter(attempt, time.Duration(base), time.Duration(max), u)
		valid := attempt >= 0 && base > 0 && max >= base && u >= 0 && u < 1
		if !valid {
			if err == nil {
				t.Fatalf("FullJitter(%d, %d, %d, %v) = %d, nil; want an error", attempt, base, max, u, got)
			}
			return
		}
		if err != nil {
			t.Fatalf("FullJitter(%d, %d, %d, %v): unexpected error %v", attempt, base, max, u, err)
		}
		sleep, _ := Exponential(attempt, time.Duration(base), time.Duration(max))
		if got < 0 || got >= sleep {
			t.Fatalf("FullJitter(%d, %d, %d, %v) = %d outside [0, %d)", attempt, base, max, u, got, sleep)
		}
	})
}

func BenchmarkFullJitter(b *testing.B) {
	b.ReportAllocs()
	r := rand.New(rand.NewPCG(1, 2))
	var us [256]float64
	for i := range us {
		us[i] = r.Float64()
	}
	var s time.Duration
	for i := 0; b.Loop(); i++ {
		d, _ := FullJitter(i&15, 100*time.Millisecond, 10*time.Second, us[i&255])
		s += d
	}
	benchSink = s
}

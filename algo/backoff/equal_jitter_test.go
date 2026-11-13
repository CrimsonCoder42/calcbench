package backoff

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
	"time"
)

// Worked example: numbers from toolbench's request R-1 (doc/toolbench_requests_2026-11-13/BRIEF.md); formula
// Brooker 2015, figure 10 (Equal Jitter).
func TestEqualJitter_Worked(t *testing.T) {
	got, err := EqualJitter(3, 100*time.Millisecond, 10*time.Second, 0.5)
	if err != nil || got != 600*time.Millisecond {
		t.Errorf("EqualJitter(3, 100ms, 10s, 0.5) = %v, %v; want 600ms", got, err)
	}
}

func TestEqualJitter_Cases(t *testing.T) {
	below1 := math.Nextafter(1, 0)
	cases := []struct {
		name    string
		attempt int
		base    time.Duration
		max     time.Duration
		u       float64
		want    time.Duration
	}{
		{"uZero", 3, 100 * time.Millisecond, 10 * time.Second, 0, 400 * time.Millisecond},
		{"uNegativeZero", 3, 100 * time.Millisecond, 10 * time.Second, math.Copysign(0, -1), 400 * time.Millisecond},
		{"uLargestBelow1", 3, 100 * time.Millisecond, 10 * time.Second, below1, 800*time.Millisecond - 1},
		{"capped", 10, 100 * time.Millisecond, 10 * time.Second, 0.5, 7500 * time.Millisecond},
		{"attempt0", 0, 100 * time.Millisecond, 10 * time.Second, 0, 50 * time.Millisecond},
		{"oddSleepKeepsCeil", 0, 7, 7, 0, 4},
		{"oddSleepTop", 0, 7, 7, below1, 6},
		{"sleep1ns", 0, 1, 1, below1, 1},
		{"base1ns_62", 62, 1, maxDur, 0, 1 << 61},
		{"base1ns_63", 63, 1, maxDur, 0, 1 << 62},
		{"base1ns_64", 64, 1, maxDur, 0.5, 1<<62 + 1<<61 - 1},
		{"base1ns_maxInt", math.MaxInt, 1, maxDur, 0, 1 << 62},
		{"baseMaxInt64Top", 0, maxDur, maxDur, below1, maxDur - 512},
	}
	for _, c := range cases {
		got, err := EqualJitter(c.attempt, c.base, c.max, c.u)
		if err != nil || got != c.want {
			t.Errorf("%s: EqualJitter(%d, %d, %d, %v) = %d, %v; want %d", c.name, c.attempt, c.base, c.max, c.u, got, err, c.want)
		}
	}
}

func TestEqualJitter_Errors(t *testing.T) {
	cases := []struct {
		name    string
		attempt int
		base    time.Duration
		max     time.Duration
		u       float64
		want    error
	}{
		{"negativeAttempt", -1, time.Second, time.Minute, 0.5, ErrAttempt},
		{"attemptFirst", math.MinInt, 0, -1, math.NaN(), ErrAttempt},
		{"zeroBase", 0, 0, time.Minute, 0.5, ErrBase},
		{"minBase", 0, time.Duration(math.MinInt64), time.Minute, 0.5, ErrBase},
		{"maxBelowBase", 0, time.Second, time.Millisecond, 0.5, ErrMax},
		{"maxBeforeDraw", 0, time.Second, 0, -1, ErrMax},
		{"uNaN", 0, time.Second, time.Minute, math.NaN(), ErrDraw},
		{"uNegativeTiny", 0, time.Second, time.Minute, -math.SmallestNonzeroFloat64, ErrDraw},
		{"uOne", 0, time.Second, time.Minute, 1, ErrDraw},
		{"uPlusInf", 0, time.Second, time.Minute, math.Inf(1), ErrDraw},
	}
	for _, c := range cases {
		got, err := EqualJitter(c.attempt, c.base, c.max, c.u)
		if !errors.Is(err, c.want) || got != 0 {
			t.Errorf("%s: EqualJitter = %v, %v; want 0, %v", c.name, got, err, c.want)
		}
	}
}

// EqualJitter is the kept half plus FullJitter's draw on the other half: checked on 500 seeded schedules.
func TestEqualJitter_HalfPlusDraw(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for i := 0; i < 500; i++ {
		attempt := r.IntN(70)
		base := 1 + time.Duration(r.Int64N(1<<40))
		u := r.Float64()
		sleep, _ := Exponential(attempt, base, maxDur)
		got, err := EqualJitter(attempt, base, maxDur, u)
		want := (sleep - sleep/2) + scale(sleep/2, u)
		if err != nil || got != want || got < sleep-sleep/2 || got > sleep {
			t.Fatalf("EqualJitter(%d, %d, max, %v) = %d, %v; want %d in [⌈%d/2⌉, %d]", attempt, base, u, got, err, want, sleep, sleep)
		}
	}
}

func TestEqualJitter_Allocs(t *testing.T) {
	var sink time.Duration
	n := testing.AllocsPerRun(100, func() {
		d, _ := EqualJitter(7, 100*time.Millisecond, 10*time.Second, 0.37)
		sink += d
	})
	if n != 0 {
		t.Errorf("EqualJitter allocates %v times per call; want 0", n)
	}
	benchSink = sink
}

// FuzzEqualJitter checks ⌈sleep/2⌉ ≤ result ≤ sleep and no error on valid input, and an error on invalid input.
// Run with -fuzztime.
func FuzzEqualJitter(f *testing.F) {
	f.Add(3, int64(100*time.Millisecond), int64(10*time.Second), 0.5)
	f.Add(63, int64(1), int64(math.MaxInt64), math.Nextafter(1, 0))
	f.Add(0, int64(1), int64(1), 0.0)
	f.Add(0, int64(1), int64(1), math.Inf(1))
	f.Fuzz(func(t *testing.T, attempt int, base, max int64, u float64) {
		got, err := EqualJitter(attempt, time.Duration(base), time.Duration(max), u)
		valid := attempt >= 0 && base > 0 && max >= base && u >= 0 && u < 1
		if !valid {
			if err == nil {
				t.Fatalf("EqualJitter(%d, %d, %d, %v) = %d, nil; want an error", attempt, base, max, u, got)
			}
			return
		}
		if err != nil {
			t.Fatalf("EqualJitter(%d, %d, %d, %v): unexpected error %v", attempt, base, max, u, err)
		}
		sleep, _ := Exponential(attempt, time.Duration(base), time.Duration(max))
		if got < sleep-sleep/2 || got > sleep {
			t.Fatalf("EqualJitter(%d, %d, %d, %v) = %d outside [⌈%d/2⌉, %d]", attempt, base, max, u, got, sleep, sleep)
		}
	})
}

func BenchmarkEqualJitter(b *testing.B) {
	b.ReportAllocs()
	r := rand.New(rand.NewPCG(1, 2))
	var us [256]float64
	for i := range us {
		us[i] = r.Float64()
	}
	var s time.Duration
	for i := 0; b.Loop(); i++ {
		d, _ := EqualJitter(i&15, 100*time.Millisecond, 10*time.Second, us[i&255])
		s += d
	}
	benchSink = s
}

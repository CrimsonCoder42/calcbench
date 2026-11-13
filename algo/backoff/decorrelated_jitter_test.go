package backoff

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
	"time"
)

// Worked example: numbers from toolbench's request R-1 (doc/toolbench_requests_2026-11-13/BRIEF.md); formula
// Brooker 2015, figure 11 (Decorrelated Jitter).
func TestDecorrelatedJitter_Worked(t *testing.T) {
	got, err := DecorrelatedJitter(300*time.Millisecond, 100*time.Millisecond, 10*time.Second, 0.5)
	if err != nil || got != 500*time.Millisecond {
		t.Errorf("DecorrelatedJitter(300ms, 100ms, 10s, 0.5) = %v, %v; want 500ms", got, err)
	}
}

func TestDecorrelatedJitter_Cases(t *testing.T) {
	below1 := math.Nextafter(1, 0)
	cases := []struct {
		name string
		prev time.Duration
		base time.Duration
		max  time.Duration
		u    float64
		want time.Duration
	}{
		{"uZeroIsBase", 300 * time.Millisecond, 100 * time.Millisecond, 10 * time.Second, 0, 100 * time.Millisecond},
		{"uNegativeZero", 300 * time.Millisecond, 100 * time.Millisecond, 10 * time.Second, math.Copysign(0, -1), 100 * time.Millisecond},
		{"firstRetryPrevIsBase", 100 * time.Millisecond, 100 * time.Millisecond, 10 * time.Second, 0.5, 200 * time.Millisecond},
		{"uLargestBelow1", 300 * time.Millisecond, 100 * time.Millisecond, 10 * time.Second, below1, 900*time.Millisecond - 1},
		{"capped", 5 * time.Second, 100 * time.Millisecond, 10 * time.Second, 0.75, 10 * time.Second},
		{"maxEqualsBase", time.Second, time.Second, time.Second, 0.5, time.Second},
		{"base1ns", 1, 1, maxDur, below1, 2},
		{"prevMaxInt64Saturates", maxDur, 1, maxDur, 0.5, 1 << 62},
		{"prevMaxInt64Top", maxDur, 1, maxDur, below1, maxDur - 1024},
		{"prevLastUnsaturated", math.MaxInt64 / 3, 1, maxDur, 0.5, (maxDur - 1) / 2},
		{"baseMaxInt64", maxDur, maxDur, maxDur, 0.5, maxDur},
	}
	for _, c := range cases {
		got, err := DecorrelatedJitter(c.prev, c.base, c.max, c.u)
		if err != nil || got != c.want {
			t.Errorf("%s: DecorrelatedJitter(%d, %d, %d, %v) = %d, %v; want %d", c.name, c.prev, c.base, c.max, c.u, got, err, c.want)
		}
	}
}

func TestDecorrelatedJitter_Errors(t *testing.T) {
	cases := []struct {
		name string
		prev time.Duration
		base time.Duration
		max  time.Duration
		u    float64
		want error
	}{
		{"prevBelowBase", 99 * time.Millisecond, 100 * time.Millisecond, 10 * time.Second, 0.5, ErrPrev},
		{"prevFirst", -1, 0, -1, math.NaN(), ErrPrev},
		{"zeroBase", 5, 0, time.Minute, 0.5, ErrBase},
		{"negativeBase", -3, -5, time.Minute, 0.5, ErrBase},
		{"maxBelowBase", time.Second, time.Second, time.Millisecond, 0.5, ErrMax},
		{"maxBeforeDraw", time.Second, time.Second, 0, 2, ErrMax},
		{"uNaN", time.Second, time.Second, time.Minute, math.NaN(), ErrDraw},
		{"uNegativeTiny", time.Second, time.Second, time.Minute, -math.SmallestNonzeroFloat64, ErrDraw},
		{"uOne", time.Second, time.Second, time.Minute, 1, ErrDraw},
		{"uPlusInf", time.Second, time.Second, time.Minute, math.Inf(1), ErrDraw},
	}
	for _, c := range cases {
		got, err := DecorrelatedJitter(c.prev, c.base, c.max, c.u)
		if !errors.Is(err, c.want) || got != 0 {
			t.Errorf("%s: DecorrelatedJitter = %v, %v; want 0, %v", c.name, got, err, c.want)
		}
	}
}

// A chain of 500 seeded retries, each feeding its result back as prev, stays in [base, max] and below 3·prev.
func TestDecorrelatedJitter_Chain(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	base, max := 100*time.Millisecond, 10*time.Second
	prev := base
	for i := 0; i < 500; i++ {
		got, err := DecorrelatedJitter(prev, base, max, r.Float64())
		if err != nil || got < base || got > max || got >= 3*prev {
			t.Fatalf("step %d: DecorrelatedJitter(%v, ...) = %v, %v", i, prev, got, err)
		}
		prev = got
	}
}

func TestDecorrelatedJitter_Allocs(t *testing.T) {
	var sink time.Duration
	n := testing.AllocsPerRun(100, func() {
		d, _ := DecorrelatedJitter(300*time.Millisecond, 100*time.Millisecond, 10*time.Second, 0.37)
		sink += d
	})
	if n != 0 {
		t.Errorf("DecorrelatedJitter allocates %v times per call; want 0", n)
	}
	benchSink = sink
}

// FuzzDecorrelatedJitter checks base ≤ result ≤ max, result < hi (hi = 3·prev saturated) unless base = MaxInt64,
// and no error on valid input; an error on invalid input. Run with -fuzztime.
func FuzzDecorrelatedJitter(f *testing.F) {
	f.Add(int64(300*time.Millisecond), int64(100*time.Millisecond), int64(10*time.Second), 0.5)
	f.Add(int64(math.MaxInt64), int64(1), int64(math.MaxInt64), math.Nextafter(1, 0))
	f.Add(int64(math.MaxInt64), int64(math.MaxInt64), int64(math.MaxInt64), 0.0)
	f.Add(int64(-1), int64(0), int64(-1), math.NaN())
	f.Fuzz(func(t *testing.T, prev, base, max int64, u float64) {
		got, err := DecorrelatedJitter(time.Duration(prev), time.Duration(base), time.Duration(max), u)
		valid := prev >= base && base > 0 && max >= base && u >= 0 && u < 1
		if !valid {
			if err == nil {
				t.Fatalf("DecorrelatedJitter(%d, %d, %d, %v) = %d, nil; want an error", prev, base, max, u, got)
			}
			return
		}
		if err != nil {
			t.Fatalf("DecorrelatedJitter(%d, %d, %d, %v): unexpected error %v", prev, base, max, u, err)
		}
		hi := int64(math.MaxInt64)
		if prev <= math.MaxInt64/3 {
			hi = 3 * prev
		}
		if int64(got) < base || int64(got) > max || (base != math.MaxInt64 && int64(got) >= hi) {
			t.Fatalf("DecorrelatedJitter(%d, %d, %d, %v) = %d outside [base, max] ∩ [base, %d)", prev, base, max, u, got, hi)
		}
	})
}

func BenchmarkDecorrelatedJitter(b *testing.B) {
	b.ReportAllocs()
	r := rand.New(rand.NewPCG(1, 2))
	var us [256]float64
	for i := range us {
		us[i] = r.Float64()
	}
	base, max := 100*time.Millisecond, 10*time.Second
	prev := base
	for i := 0; b.Loop(); i++ {
		prev, _ = DecorrelatedJitter(prev, base, max, us[i&255])
	}
	benchSink = prev
}

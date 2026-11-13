package backoff

import (
	"errors"
	"math"
	"math/big"
	"testing"
	"time"
)

// maxDur is the largest time.Duration.
const maxDur = time.Duration(math.MaxInt64)

// benchSink keeps benchmark results alive.
var benchSink time.Duration

// Worked example: numbers from toolbench's request R-1 (doc/toolbench_requests_2026-11-13/BRIEF.md); formula
// Brooker 2015, figure 3 (capped exponential).
func TestExponential_Worked(t *testing.T) {
	cases := []struct {
		attempt int
		want    time.Duration
	}{
		{3, 800 * time.Millisecond},
		{10, 10 * time.Second},
	}
	for _, c := range cases {
		got, err := Exponential(c.attempt, 100*time.Millisecond, 10*time.Second)
		if err != nil || got != c.want {
			t.Errorf("Exponential(%d) = %v, %v; want %v", c.attempt, got, err, c.want)
		}
	}
}

func TestExponential_Cases(t *testing.T) {
	cases := []struct {
		name    string
		attempt int
		base    time.Duration
		max     time.Duration
		want    time.Duration
	}{
		{"attempt0IsBase", 0, 100 * time.Millisecond, 10 * time.Second, 100 * time.Millisecond},
		{"justBelowCap", 6, 100 * time.Millisecond, 10 * time.Second, 6400 * time.Millisecond},
		{"firstAtCap", 7, 100 * time.Millisecond, 10 * time.Second, 10 * time.Second},
		{"exactlyCap", 2, 1, 4, 4},
		{"base1ns_62", 62, 1, maxDur, 1 << 62},
		{"base1ns_63", 63, 1, maxDur, maxDur},
		{"base1ns_64", 64, 1, maxDur, maxDur},
		{"base1ns_maxInt", math.MaxInt, 1, maxDur, maxDur},
		{"base1ns_62_smallCap", 62, 1, 1000, 1000},
		{"baseMaxInt64_0", 0, maxDur, maxDur, maxDur},
		{"baseMaxInt64_1", 1, maxDur, maxDur, maxDur},
		{"maxEqualsBase", 5, time.Second, time.Second, time.Second},
		{"largestFitting", 1, 1 << 61, maxDur, 1 << 62},
		{"shiftFitsAboveCap", 1, 1 << 61, 1<<62 - 1, 1<<62 - 1},
	}
	for _, c := range cases {
		got, err := Exponential(c.attempt, c.base, c.max)
		if err != nil || got != c.want {
			t.Errorf("%s: Exponential(%d, %d, %d) = %d, %v; want %d", c.name, c.attempt, c.base, c.max, got, err, c.want)
		}
	}
}

func TestExponential_Errors(t *testing.T) {
	cases := []struct {
		name    string
		attempt int
		base    time.Duration
		max     time.Duration
		want    error
	}{
		{"negativeAttempt", -1, time.Second, time.Minute, ErrAttempt},
		{"minIntAttempt", math.MinInt, time.Second, time.Minute, ErrAttempt},
		{"attemptBeforeBase", -1, 0, -5, ErrAttempt},
		{"zeroBase", 0, 0, time.Minute, ErrBase},
		{"negativeBase", 0, -time.Second, time.Minute, ErrBase},
		{"baseBeforeMax", 0, -time.Second, -2 * time.Second, ErrBase},
		{"maxBelowBase", 0, time.Second, time.Second - 1, ErrMax},
		{"maxNegative", 0, time.Second, -time.Second, ErrMax},
	}
	for _, c := range cases {
		got, err := Exponential(c.attempt, c.base, c.max)
		if !errors.Is(err, c.want) || got != 0 {
			t.Errorf("%s: Exponential = %v, %v; want 0, %v", c.name, got, err, c.want)
		}
	}
}

func TestExponential_Allocs(t *testing.T) {
	var sink time.Duration
	n := testing.AllocsPerRun(100, func() {
		d, _ := Exponential(7, 100*time.Millisecond, 10*time.Second)
		sink += d
	})
	if n != 0 {
		t.Errorf("Exponential allocates %v times per call; want 0", n)
	}
	benchSink = sink
}

// exponentialOracle is min(max, base·2^attempt) in math/big.
func exponentialOracle(attempt int, base, max time.Duration) time.Duration {
	if attempt > 200 {
		return max
	}
	v := new(big.Int).Lsh(big.NewInt(int64(base)), uint(attempt))
	if v.Cmp(big.NewInt(int64(max))) >= 0 {
		return max
	}
	return time.Duration(v.Int64())
}

// FuzzExponential checks the result against a math/big oracle and the bounds base ≤ result ≤ max on valid
// input, and that invalid input errors. Run with -fuzztime.
func FuzzExponential(f *testing.F) {
	f.Add(3, int64(100*time.Millisecond), int64(10*time.Second))
	f.Add(62, int64(1), int64(math.MaxInt64))
	f.Add(0, int64(math.MaxInt64), int64(math.MaxInt64))
	f.Add(-1, int64(0), int64(-1))
	f.Fuzz(func(t *testing.T, attempt int, base, max int64) {
		got, err := Exponential(attempt, time.Duration(base), time.Duration(max))
		valid := attempt >= 0 && base > 0 && max >= base
		if !valid {
			if err == nil {
				t.Fatalf("Exponential(%d, %d, %d) = %d, nil; want an error", attempt, base, max, got)
			}
			return
		}
		if err != nil {
			t.Fatalf("Exponential(%d, %d, %d): unexpected error %v", attempt, base, max, err)
		}
		if want := exponentialOracle(attempt, time.Duration(base), time.Duration(max)); got != want {
			t.Fatalf("Exponential(%d, %d, %d) = %d; want %d", attempt, base, max, got, want)
		}
		if got < time.Duration(base) || got > time.Duration(max) {
			t.Fatalf("Exponential(%d, %d, %d) = %d outside [base, max]", attempt, base, max, got)
		}
	})
}

func BenchmarkExponential(b *testing.B) {
	b.ReportAllocs()
	var s time.Duration
	for i := 0; b.Loop(); i++ {
		d, _ := Exponential(i&63, 100*time.Millisecond, 10*time.Second)
		s += d
	}
	benchSink = s
}

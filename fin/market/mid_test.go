package market

import (
	"errors"
	"math"
	"testing"
)

var (
	sinkFloat float64
	sinkErr   error
)

// near reports whether got is within tol of want (absolute).
func near(got, want, tol float64) bool { return math.Abs(got-want) <= tol }

// wantErr fails the test unless err wraps target.
func wantErr(t *testing.T, name string, err, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Errorf("%s: err = %v, want %v", name, err, target)
	}
}

// Hand-checked: bid 2419, ask 2421, tick 0.01 → 2420·0.01 = 24.20 (tolerance: one rounding of the product).
func TestMid_handChecked(t *testing.T) {
	got, err := Mid(2419, 2421, 0.01)
	if err != nil || !near(got, 24.20, 1e-12) {
		t.Errorf("Mid = %v, %v; want 24.20", got, err)
	}
}

func TestMid_cases(t *testing.T) {
	for _, c := range []struct {
		name     string
		bid, ask int64
		tick     float64
		want     float64
	}{
		{"locked", 100, 100, 1, 100},
		{"half tick", 100, 101, 1, 100.5},
		{"zero", 0, 0, 0.01, 0},
		{"large exact", 1 << 52, 1<<52 + 2, 1, 1<<52 + 1},
	} {
		got, err := Mid(c.bid, c.ask, c.tick)
		if err != nil || got != c.want {
			t.Errorf("%s: %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestMid_errors(t *testing.T) {
	_, err := Mid(101, 100, 1)
	wantErr(t, "crossed", err, ErrDomain)
	_, err = Mid(-1, 100, 1)
	wantErr(t, "negative bid", err, ErrNegative)
	_, err = Mid(1, maxTicks+1, 1)
	wantErr(t, "ask > 2^53", err, ErrOutOfRange)
	_, err = Mid(1, 2, 0)
	wantErr(t, "tick 0", err, ErrOutOfRange)
	_, err = Mid(1, 2, math.NaN())
	wantErr(t, "tick NaN", err, ErrNotFinite)
	_, err = Mid(1, maxTicks, math.MaxFloat64)
	wantErr(t, "overflow", err, ErrOverflow)
}

func TestMid_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = Mid(2419, 2421, 0.01) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkMid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkErr = Mid(2419, 2421, 0.01)
	}
}

package invest

import (
	"errors"
	"math"
	"testing"
)

// Shared test helpers and benchmark sinks for package invest.

var (
	sinkFloat float64
	sinkCents int64
	sinkErr   error
)

const maxCents = int64(1) << 53

// near reports |got − want| ≤ tol.
func near(got, want, tol float64) bool { return math.Abs(got-want) <= tol }

// wantErr fails unless err wraps sentinel.
func wantErr(t *testing.T, name string, err, sentinel error) {
	t.Helper()
	if !errors.Is(err, sentinel) {
		t.Errorf("%s: err = %v, want %v", name, err, sentinel)
	}
}

// checkRounded fails unless got is exact rounded half-even, allowing 1 cent of slack only when exact
// lies within 1e-6 cent of a half cent (the F-5 accuracy contract; exact must be ≤ 2⁴³ cents).
func checkRounded(t *testing.T, name string, got int64, exact float64) {
	t.Helper()
	want := int64(math.RoundToEven(exact))
	if got == want {
		return
	}
	nearHalf := math.Abs(math.Abs(exact-math.Trunc(exact))-0.5) < 1e-6
	if !nearHalf || got-want > 1 || want-got > 1 {
		t.Errorf("%s = %d, exact %v (want %d)", name, got, exact, want)
	}
}

// noErr fails the test now on an unexpected error.
func noErr(t *testing.T, name string, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: unexpected error %v", name, err)
	}
}

// Hand-checked (no worked example in MSRB "Municipal Bond Basics"): $50.00 / $952.06.
func TestCurrentYield_handChecked(t *testing.T) {
	got, err := CurrentYield(5000, 95206)
	noErr(t, "CurrentYield", err)
	if got != 5000.0/95206.0 {
		t.Errorf("CurrentYield(5000, 95206) = %v, want %v", got, 5000.0/95206.0)
	}
	if !near(got, 0.0525176984643825, 1e-16) {
		t.Errorf("CurrentYield = %v", got)
	}
}

func TestCurrentYield_cases(t *testing.T) {
	cases := []struct {
		name          string
		coupon, price int64
		want          float64
	}{
		{"zero coupon", 0, 95000, 0},
		{"at par", 5000, 100000, 0.05},
		{"max", maxCents, maxCents, 1},
		{"one cent", 1, 1, 1},
	}
	for _, c := range cases {
		got, err := CurrentYield(c.coupon, c.price)
		if err != nil || got != c.want {
			t.Errorf("%s: CurrentYield(%d, %d) = %v, %v; want %v", c.name, c.coupon, c.price, got, err, c.want)
		}
	}
}

func TestCurrentYield_errors(t *testing.T) {
	_, err := CurrentYield(5000, 0)
	wantErr(t, "price 0", err, ErrZeroDenominator)
	_, err = CurrentYield(-1, 100)
	wantErr(t, "coupon < 0", err, ErrNegative)
	_, err = CurrentYield(1, -100)
	wantErr(t, "price < 0", err, ErrNegative)
	_, err = CurrentYield(maxCents+1, 100)
	wantErr(t, "coupon > 2^53", err, ErrOutOfRange)
	_, err = CurrentYield(1, maxCents+1)
	wantErr(t, "price > 2^53", err, ErrOutOfRange)
}

func TestCurrentYield_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = CurrentYield(5000, 95206) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkCurrentYield(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkErr = CurrentYield(5000, 95206)
	}
}

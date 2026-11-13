package invest

import (
	"math"
	"testing"
)

// Hand-checked: the 6 % semiannual 10-year bond at par, Σ k(k+1)·cₖvᵏ⁺²/(P·4) = 68.774822 years²
// (computed independently in double precision; ±1e-6, the hand value's last digit).
func TestConvexity_handChecked(t *testing.T) {
	got, err := Convexity(100000, 0.06, 0.06, 20, 2)
	if err != nil || !near(got, 68.774822, 1e-6) {
		t.Errorf("Convexity = %v, %v; want 68.774822", got, err)
	}
}

func TestConvexity_cases(t *testing.T) {
	// Zero-coupon, annual: n(n+1)/(1+y)².
	got, err := Convexity(100000, 0, 0.05, 10, 1)
	if want := 110 / (1.05 * 1.05); err != nil || !near(got, want, 1e-12) {
		t.Errorf("zero coupon: %v, %v; want %v", got, err, want)
	}
}

// Oracle: (1/P)·d²P/dy² by a second central difference, h = 1e-4 (truncation ≈ 1e-8 relative,
// rounding ≈ 1e-16/h² = 1e-8), compared to 1e-5 relative.
func TestConvexity_secondDifference(t *testing.T) {
	for _, c := range []struct {
		cr, y float64
		n, m  int
	}{{0.04, 0.06, 20, 2}, {0.08, 0.03, 30, 1}, {0, 0.05, 40, 4}, {0.05, 0, 10, 1}} {
		got, err := Convexity(100000, c.cr, c.y, c.n, c.m)
		noErr(t, "Convexity", err)
		const h = 1e-4
		p := priceAt(100000, c.cr, c.y, c.n, c.m)
		d2 := (priceAt(100000, c.cr, c.y+h, c.n, c.m) - 2*p + priceAt(100000, c.cr, c.y-h, c.n, c.m)) / (h * h) / p
		if math.Abs(got-d2) > 1e-5*d2 {
			t.Errorf("%+v: %v, finite difference %v", c, got, d2)
		}
	}
}

func TestConvexity_errors(t *testing.T) {
	_, err := Convexity(0, 0, 0.05, 10, 1)
	wantErr(t, "no flows", err, ErrDomain)
	_, err = Convexity(100000, -0.05, 0.05, 10, 1)
	wantErr(t, "coupon < 0", err, ErrNegative)
	_, err = Convexity(maxCents+1, 0.05, 0.05, 10, 1)
	wantErr(t, "face > 2^53", err, ErrOutOfRange)
	_, err = Convexity(100000, 0.05, math.Inf(-1), 10, 1)
	wantErr(t, "yield −Inf", err, ErrRate)
	_, err = Convexity(100000, 0.05, 0.05, 10, 0)
	wantErr(t, "periodsPerYear 0", err, ErrPeriods)
	_, err = Convexity(maxCents, 0.05, -0.999999, 2000, 1)
	wantErr(t, "overflow", err, ErrOverflow)
}

// periods is capped at 10 000 000 (maxPeriods) so one call is bounded.
func TestConvexity_periodsCap(t *testing.T) {
	_, err := Convexity(100000, 0.05, 0.05, maxPeriods+1, 12)
	wantErr(t, "periods > cap", err, ErrPeriods)
	_, err = Convexity(100000, 0.05, 0.05, math.MaxInt, 12)
	wantErr(t, "periods MaxInt", err, ErrPeriods)
}

func TestConvexity_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = Convexity(100000, 0.06, 0.06, 20, 2) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkConvexity(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run(itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, sinkErr = Convexity(100000, 0.06, 0.05, n, 12)
			}
		})
	}
}

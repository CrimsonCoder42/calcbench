package invest

import (
	"math"
	"testing"
)

// priceAt returns the unrounded bond price at nominal yield y (the oracle for the derivatives).
func priceAt(face int64, cr, y float64, n, m int) float64 {
	i := y / float64(m)
	p, _, _ := bondMoments(face, i, float64(face)*cr/float64(m), n)
	return p
}

// Hand-checked from Macaulay 1938 p. 51: the 6 % semiannual 10-year bond at par has D = 7.661900 years,
// so the modified duration is 7.661900/1.03 = 7.438737 (±1e-6, the hand value's last digit).
func TestModifiedDuration_handChecked(t *testing.T) {
	got, err := ModifiedDuration(100000, 0.06, 0.06, 20, 2)
	if err != nil || !near(got, 7.438737, 1e-6) {
		t.Errorf("ModifiedDuration = %v, %v; want 7.438737", got, err)
	}
	mac, _ := MacaulayDuration(100000, 0.06, 0.06, 20, 2)
	if !near(got, mac/1.03, 1e-14) {
		t.Errorf("ModifiedDuration %v ≠ Macaulay/1.03 %v", got, mac/1.03)
	}
}

// Oracle: −(1/P)·dP/dy by a central difference, h = 1e-6 (truncation O(h²) ≈ 1e-12, rounding
// ≈ 1e-16/h = 1e-10), compared to 1e-7 relative.
func TestModifiedDuration_centralDifference(t *testing.T) {
	for _, c := range []struct {
		cr, y float64
		n, m  int
	}{{0.04, 0.06, 20, 2}, {0.08, 0.03, 30, 1}, {0, 0.05, 40, 4}, {0.05, 0, 10, 1}, {0.12, 0.2, 360, 12}} {
		got, err := ModifiedDuration(100000, c.cr, c.y, c.n, c.m)
		noErr(t, "ModifiedDuration", err)
		const h = 1e-6
		p := priceAt(100000, c.cr, c.y, c.n, c.m)
		d := -(priceAt(100000, c.cr, c.y+h, c.n, c.m) - priceAt(100000, c.cr, c.y-h, c.n, c.m)) / (2 * h) / p
		if math.Abs(got-d) > 1e-7*d {
			t.Errorf("%+v: %v, finite difference %v", c, got, d)
		}
	}
}

func TestModifiedDuration_errors(t *testing.T) {
	_, err := ModifiedDuration(0, 0, 0.05, 10, 1)
	wantErr(t, "no flows", err, ErrDomain)
	_, err = ModifiedDuration(-1, 0.05, 0.05, 10, 1)
	wantErr(t, "face < 0", err, ErrNegative)
	_, err = ModifiedDuration(maxCents+1, 0.05, 0.05, 10, 1)
	wantErr(t, "face > 2^53", err, ErrOutOfRange)
	_, err = ModifiedDuration(100000, 0.05, math.Inf(1), 10, 1)
	wantErr(t, "yield Inf", err, ErrRate)
	_, err = ModifiedDuration(100000, 0.05, 0.05, 0, 1)
	wantErr(t, "periods 0", err, ErrPeriods)
	_, err = ModifiedDuration(maxCents, 0.05, -0.999999, 2000, 1)
	wantErr(t, "overflow", err, ErrOverflow)
}

// periods is capped at 10 000 000 (maxPeriods) so one call is bounded.
func TestModifiedDuration_periodsCap(t *testing.T) {
	_, err := ModifiedDuration(100000, 0.05, 0.05, maxPeriods+1, 12)
	wantErr(t, "periods > cap", err, ErrPeriods)
	_, err = ModifiedDuration(100000, 0.05, 0.05, math.MaxInt, 12)
	wantErr(t, "periods MaxInt", err, ErrPeriods)
}

func TestModifiedDuration_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = ModifiedDuration(100000, 0.06, 0.06, 20, 2) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkModifiedDuration(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run(itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, sinkErr = ModifiedDuration(100000, 0.06, 0.05, n, 12)
			}
		})
	}
}

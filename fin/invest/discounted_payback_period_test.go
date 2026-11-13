package invest

import (
	"math"
	"testing"
)

// OpenStax Principles of Finance §16.4, embroidery machine at 9 %: cumulative discounted flow
// −$145.72 after year 5, year-6 discounted flow $2,981.34 → 5 + 145.72/2,981.34 = 5.05 years
// (two decimals: ±5e-3); hand: 5 + 14571.73/298133.66 = 5.0488765 (±1e-6).
func TestDiscountedPaybackPeriod_openStax(t *testing.T) {
	got, err := DiscountedPaybackPeriod(0.09, embroidery)
	if err != nil || !near(got, 5.05, 5e-3) || !near(got, 5.048876484731, 1e-6) {
		t.Errorf("DiscountedPaybackPeriod = %v, %v; want 5.0488765", got, err)
	}
}

func TestDiscountedPaybackPeriod_cases(t *testing.T) {
	// Zero rate equals the plain payback period.
	flows := []int64{-1600000, 200000, 400000, 500000, 500000, 500000}
	got, err := DiscountedPaybackPeriod(0, flows)
	plain, _ := PaybackPeriod(flows)
	if err != nil || got != plain {
		t.Errorf("zero rate: %v, %v; want %v", got, err, plain)
	}
	// Rate 0.25: 125·fl(1/1.25) rounds to exactly 100, so recovery is at exactly t = 1.
	got, err = DiscountedPaybackPeriod(0.25, []int64{-100, 125})
	if err != nil || got != 1 {
		t.Errorf("one period: %v, %v; want 1", got, err)
	}
	got, err = DiscountedPaybackPeriod(0.1, []int64{5, -1})
	if err != nil || got != 0 {
		t.Errorf("first flow positive: %v, %v; want 0", got, err)
	}
	// A discount rate can push recovery later than the plain payback.
	got, err = DiscountedPaybackPeriod(0.5, []int64{-100, 100, 100})
	if err != nil || !near(got, 1+(100-100/1.5)/(100/2.25), 1e-14) {
		t.Errorf("later: %v, %v", got, err)
	}
}

func TestDiscountedPaybackPeriod_errors(t *testing.T) {
	_, err := DiscountedPaybackPeriod(0.1, nil)
	wantErr(t, "empty", err, ErrDomain)
	_, err = DiscountedPaybackPeriod(0.1, []int64{-1000, 500, 500})
	wantErr(t, "never recovered at 10 %", err, ErrNoSolution)
	_, err = DiscountedPaybackPeriod(0.1, []int64{maxCents + 1})
	wantErr(t, "flow > 2^53", err, ErrOutOfRange)
	for _, bad := range []float64{math.NaN(), math.Inf(1), -1, -3} {
		_, err = DiscountedPaybackPeriod(bad, embroidery)
		wantErr(t, "bad rate", err, ErrRate)
	}
}

func TestDiscountedPaybackPeriod_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = DiscountedPaybackPeriod(0.09, embroidery) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkDiscountedPaybackPeriod(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		flows := benchFlows(n)
		flows[0] = -100_000 * int64(n-1) // at a zero rate, recovered only at the last period
		b.Run(itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, sinkErr = DiscountedPaybackPeriod(0, flows)
			}
		})
	}
}

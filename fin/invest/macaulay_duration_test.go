package invest

import (
	"math"
	"testing"
)

// Macaulay 1938, ch. 2, p. 51, table "Duration in years": semiannual bonds selling at par (yield =
// coupon). The table shows three decimals, so the comparison is to ±5e-4 years. Two printed entries
// disagree with the formula on the same page and are left out: 5 % 1 year (.987; the formula gives
// 0.98780) and 6 % 25 years (13.254; the formula gives 13.25083), recorded in research/invest.md.
func TestMacaulayDuration_macaulayTable(t *testing.T) {
	cases := []struct {
		coupon float64
		years  int
		want   float64
	}{
		{0.04, 1, 0.990}, {0.06, 1, 0.985},
		{0.04, 3, 2.857}, {0.05, 3, 2.823}, {0.06, 3, 2.790},
		{0.04, 10, 8.339}, {0.05, 10, 7.989}, {0.06, 10, 7.662},
		{0.04, 25, 16.026}, {0.05, 25, 14.536},
		{0.04, 100, 25.014}, {0.05, 100, 20.353}, {0.06, 100, 17.120},
	}
	for _, c := range cases {
		got, err := MacaulayDuration(100000, c.coupon, c.coupon, 2*c.years, 2)
		if err != nil || !near(got, c.want, 5e-4) {
			t.Errorf("%v %% %d years: D = %v, %v; want %v", c.coupon*100, c.years, got, err, c.want)
		}
	}
}

func TestMacaulayDuration_cases(t *testing.T) {
	// A zero-coupon bond's duration is its maturity (Macaulay p. 49: "D = n").
	got, err := MacaulayDuration(100000, 0, 0.07, 20, 2)
	if err != nil || !near(got, 10, 1e-12) {
		t.Errorf("zero coupon: %v, %v; want 10", got, err)
	}
	// One period: D = 1 period (Macaulay p. 50).
	got, err = MacaulayDuration(100000, 0.05, 0.08, 1, 4)
	if err != nil || !near(got, 0.25, 1e-15) {
		t.Errorf("one period: %v, %v; want 0.25", got, err)
	}
	// Zero yield: equal-weighted mean time; 2 periods, coupon 10, face 100 → (10·1 + 110·2)/120 years.
	got, err = MacaulayDuration(10000, 0.10, 0, 2, 1)
	if err != nil || !near(got, 230.0/120, 1e-14) {
		t.Errorf("zero yield: %v, %v; want %v", got, err, 230.0/120)
	}
}

// Oracle: Macaulay's closed form (p. 49), D = R/(R−1) − (QR + n(1+Q−QR))/(Rⁿ − 1 − Q + QR) half years,
// Q = face/coupon, R = 1 + i, away from i = 0 where it cancels.
func TestMacaulayDuration_closedForm(t *testing.T) {
	for _, c := range []struct {
		cr, y float64
		n     int
	}{{0.04, 0.06, 20}, {0.06, 0.04, 40}, {0.08, 0.10, 60}, {0.03, 0.12, 200}} {
		got, err := MacaulayDuration(100000, c.cr, c.y, c.n, 2)
		noErr(t, "MacaulayDuration", err)
		q := 1 / (c.cr / 2)
		r := 1 + c.y/2
		n := float64(c.n)
		d := r/(r-1) - (q*r+n*(1+q-q*r))/(math.Pow(r, n)-1-q+q*r)
		if !near(got, d/2, 1e-10) {
			t.Errorf("%+v: %v, closed form %v", c, got, d/2)
		}
	}
}

func TestMacaulayDuration_errors(t *testing.T) {
	_, err := MacaulayDuration(0, 0, 0.05, 10, 1)
	wantErr(t, "no flows", err, ErrDomain)
	_, err = MacaulayDuration(-1, 0.05, 0.05, 10, 1)
	wantErr(t, "face < 0", err, ErrNegative)
	_, err = MacaulayDuration(maxCents+1, 0.05, 0.05, 10, 1)
	wantErr(t, "face > 2^53", err, ErrOutOfRange)
	_, err = MacaulayDuration(100000, 0.05, math.NaN(), 10, 1)
	wantErr(t, "yield NaN", err, ErrRate)
	_, err = MacaulayDuration(100000, 0.05, -1, 10, 1)
	wantErr(t, "yield −1", err, ErrRate)
	_, err = MacaulayDuration(100000, 0.05, 0.05, 0, 1)
	wantErr(t, "periods 0", err, ErrPeriods)
	_, err = MacaulayDuration(100000, 0.05, 0.05, 10, 0)
	wantErr(t, "periodsPerYear 0", err, ErrPeriods)
	_, err = MacaulayDuration(maxCents, 0.05, -0.999999, 2000, 1)
	wantErr(t, "overflow", err, ErrOverflow)
}

// periods is capped at 10 000 000 (maxPeriods) so one call is bounded.
func TestMacaulayDuration_periodsCap(t *testing.T) {
	_, err := MacaulayDuration(100000, 0.05, 0.05, maxPeriods+1, 12)
	wantErr(t, "periods > cap", err, ErrPeriods)
	_, err = MacaulayDuration(100000, 0.05, 0.05, math.MaxInt, 12)
	wantErr(t, "periods MaxInt", err, ErrPeriods)
}

func TestMacaulayDuration_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = MacaulayDuration(100000, 0.06, 0.06, 20, 2) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkMacaulayDuration(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run(itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, sinkErr = MacaulayDuration(100000, 0.06, 0.05, n, 12)
			}
		})
	}
}

// itoa formats a benchmark size without strconv in every file.
func itoa(n int) string {
	switch n {
	case 1_000:
		return "1e3"
	case 10_000:
		return "1e4"
	case 100_000:
		return "1e5"
	}
	return "n"
}

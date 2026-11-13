package invest

import (
	"math"
	"math/rand/v2"
	"testing"
)

// OpenStax Principles of Finance §10.2 "Bond Valuation", Coca-Cola bond: price $952.06, par $1,000,
// $5.00 coupon per half-year, 14 periods → 0.8651 % per half-year, 1.73 % per year. The source shows
// four decimals of a percent, so the periodic yield is compared to ±5e-7.
func TestYieldToMaturity_openStax(t *testing.T) {
	got, err := YieldToMaturity(95206, 100000, 0.01, 14, 2)
	noErr(t, "Coca-Cola", err)
	if !near(got/2, 0.008651, 5e-7) {
		t.Errorf("periodic YTM = %v, want 0.008651 ± 5e-7", got/2)
	}
	if !near(got, 0.0173, 5e-5) {
		t.Errorf("annual YTM = %v, want 0.0173 ± 5e-5", got)
	}
}

func TestYieldToMaturity_cases(t *testing.T) {
	// At par the yield equals the coupon rate (tolerance: the root finder's 1e-12 per period).
	got, err := YieldToMaturity(100000, 100000, 0.06, 20, 2)
	if err != nil || !near(got, 0.06, 1e-11) {
		t.Errorf("at par: %v, %v; want 0.06", got, err)
	}
	// One-period zero-coupon: 90000 → 100000 is 1/9.
	got, err = YieldToMaturity(90000, 100000, 0, 1, 1)
	if err != nil || !near(got, 1.0/9, 1e-12) {
		t.Errorf("zero coupon: %v, %v; want 1/9", got, err)
	}
	// Price above all flows: negative yield.
	got, err = YieldToMaturity(110000, 100000, 0, 1, 1)
	if err != nil || !near(got, -1.0/11, 1e-12) {
		t.Errorf("negative yield: %v, %v; want −1/11", got, err)
	}
	// Price equal to the sum of the flows: zero yield.
	got, err = YieldToMaturity(150000, 100000, 0.05, 10, 1)
	if err != nil || !near(got, 0, 1e-12) {
		t.Errorf("zero yield: %v, %v", got, err)
	}
}

func TestYieldToMaturity_errors(t *testing.T) {
	_, err := YieldToMaturity(0, 100000, 0.05, 10, 1)
	wantErr(t, "price 0", err, ErrDomain)
	_, err = YieldToMaturity(-1, 100000, 0.05, 10, 1)
	wantErr(t, "price < 0", err, ErrNegative)
	_, err = YieldToMaturity(maxCents+1, 100000, 0.05, 10, 1)
	wantErr(t, "price > 2^53", err, ErrOutOfRange)
	_, err = YieldToMaturity(100000, maxCents+1, 0.05, 10, 1)
	wantErr(t, "face > 2^53", err, ErrOutOfRange)
	_, err = YieldToMaturity(100000, maxCents, 2, 10, 1)
	wantErr(t, "coupon > 2^53", err, ErrOutOfRange)
	_, err = YieldToMaturity(100000, -1, 0.05, 10, 1)
	wantErr(t, "face < 0", err, ErrNegative)
	_, err = YieldToMaturity(100000, 100000, math.NaN(), 10, 1)
	wantErr(t, "coupon NaN", err, ErrRate)
	_, err = YieldToMaturity(100000, 100000, 0.05, 0, 1)
	wantErr(t, "periods 0", err, ErrPeriods)
	_, err = YieldToMaturity(100000, 100000, 0.05, 10, 0)
	wantErr(t, "periodsPerYear 0", err, ErrPeriods)
	_, err = YieldToMaturity(100000, 0, 0, 10, 1)
	wantErr(t, "pays nothing", err, ErrDomain)
}

// periods = math.MaxInt must not panic or allocate: the closed-form solver returns a value or an error.
func TestYieldToMaturity_maxPeriods(t *testing.T) {
	got, err := YieldToMaturity(95206, 100000, 0.01, math.MaxInt, 2)
	if err == nil && (math.IsNaN(got) || math.IsInf(got, 0)) {
		t.Errorf("MaxInt periods: %v with nil error", got)
	}
	// A perpetual 1 % coupon priced at $952.06 yields about 500/95206 per half-year.
	if err == nil && !near(got/2, 500.0/95206, 1e-9) {
		t.Errorf("MaxInt periods: periodic %v, want ≈ %v", got/2, 500.0/95206)
	}
}

func TestYieldToMaturity_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = YieldToMaturity(95206, 100000, 0.01, 14, 2) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

// Oracle: pricing the bond at the returned yield gives back the price to the cent (NPV(YTM) ≈ 0).
func TestYieldToMaturity_roundTrip(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	ms := []int{1, 2, 4}
	for range 200 {
		m := ms[r.IntN(len(ms))]
		n := r.IntN(120) + 1
		face := int64(100000) * (r.Int64N(100) + 1)
		cr := float64(r.IntN(1200)) / 10000 // whole basis points per year
		// Use coupons that are whole cents so the yield is the exact root (0.04 always is for m ≤ 4).
		if c := float64(face) * cr / float64(m); c != math.Trunc(c) {
			cr = 0.04
		}
		y := float64(r.IntN(1500)+1) / 10000
		price, err := BondPrice(face, cr, y, n, m)
		noErr(t, "BondPrice", err)
		ytm, err := YieldToMaturity(price, face, cr, n, m)
		noErr(t, "YieldToMaturity", err)
		back, err := BondPrice(face, cr, ytm, n, m)
		noErr(t, "BondPrice(ytm)", err)
		if d := back - price; d < -1 || d > 1 {
			t.Errorf("face %d cr %v y %v n %d m %d: price %d, YTM %v reprices to %d", face, cr, y, n, m, price, ytm, back)
		}
	}
}

func BenchmarkYieldToMaturity(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkErr = YieldToMaturity(95206, 100000, 0.01, 14, 2)
	}
}

package interest

import (
	"math"
	"testing"
)

func TestAnnuityRate_SourceExample(t *testing.T) {
	// Dahlquist & Knight, Principles of Finance (OpenStax 2022) §7.3: PV $1,000, FV $1,125.51, N 4 →
	// I/Y 3 % (FV is rounded to the cent, so the rate is 3 % to about 3e-7).
	got, err := AnnuityRate(100000, 0, 112551, 4)
	wantFloat(t, "§7.3 solve for I/Y", got, err, 0.03, 1e-6)
	// DeFusco et al., Quantitative Investment Analysis 3rd ed. ch. 1 Example 20 reversed: $100,000
	// repaid by 360 payments of $733.76 → 0.08/12 a month.
	got, err = AnnuityRate(10000000, 73376, 0, 360)
	wantFloat(t, "QIA Example 20 reversed", got, err, 0.08/12, 1e-6)
}

func TestAnnuityRate_Table(t *testing.T) {
	cases := []struct {
		name        string
		pv, pmt, fv int64
		periods     int
		want, tol   float64
	}{
		{"bond at par yields its coupon", 100000, 2500, 100000, 20, 0.025, 1e-12},
		{"one payment", 100000, 110000, 0, 1, 0.1, 1e-12},
		{"payments sum to PV: 0 exactly", 120000, 10000, 0, 12, 0, 0},
		{"lump-sum decline: negative rate", 100000, 0, 90000, 1, -0.1, 1e-12},
		{"premium bond yields below coupon", 110000, 5000, 100000, 10, 0.0378052413506336, 1e-11},
		{"very long annuity tends to PMT/PV", 100000, 5000, 0, 1 << 50, 0.05, 1e-12},
	}
	for _, c := range cases {
		got, err := AnnuityRate(c.pv, c.pmt, c.fv, c.periods)
		wantFloat(t, c.name, got, err, c.want, c.tol)
	}
}

func TestAnnuityRate_AnyPeriods(t *testing.T) {
	// periods = math.MaxInt: the closed forms cost O(1); the root ln 2/n is below the 1e-12 tolerance.
	got, err := AnnuityRate(100000, 0, 200000, math.MaxInt)
	wantFloat(t, "MaxInt periods", got, err, math.Ln2/float64(math.MaxInt), 1e-12)
	got, err = AnnuityRate(100000, 1, 0, math.MaxInt)
	wantFloat(t, "MaxInt periods with a payment", got, err, 1e-5, 1e-12)
}

func TestAnnuityRate_Errors(t *testing.T) {
	_, err := AnnuityRate(maxExactCents+1, 1, 1, 1)
	wantErr(t, "pv 2^53+1", err, ErrOutOfRange)
	_, err = AnnuityRate(1, maxExactCents+1, 1, 1)
	wantErr(t, "pmt 2^53+1", err, ErrOutOfRange)
	_, err = AnnuityRate(1, 1, maxExactCents+1, 1)
	wantErr(t, "fv 2^53+1", err, ErrOutOfRange)
	_, err = AnnuityRate(-1, 1, 1, 1)
	wantErr(t, "negative pv", err, ErrNegative)
	_, err = AnnuityRate(1, -1, 1, 1)
	wantErr(t, "negative pmt", err, ErrNegative)
	_, err = AnnuityRate(1, 1, -1, 1)
	wantErr(t, "negative fv", err, ErrNegative)
	_, err = AnnuityRate(0, 1, 1, 1)
	wantErr(t, "zero pv", err, ErrDomain)
	_, err = AnnuityRate(1, 0, 0, 1)
	wantErr(t, "nothing received", err, ErrDomain)
	_, err = AnnuityRate(1, 1, 1, 0)
	wantErr(t, "zero periods", err, ErrPeriods)
	// 1 cent now for 2^53 cents in one period needs a rate of about 9e15, above the 1e6 scan limit.
	_, err = AnnuityRate(1, maxExactCents, 0, 1)
	wantErr(t, "rate above 1e6", err, ErrNoSolution)
}

func TestAnnuityRate_Oracle(t *testing.T) {
	rng := testRNG()
	for k := 0; k < 300; k++ {
		pmt := 10_000_000 + rng.Int64N(10_000_000)
		fv := rng.Int64N(1_000_000_000)
		r := rng.Float64() * 0.2
		n := 1 + rng.IntN(400)
		// Round trip through PresentValueAnnuity (016) and PresentValue (002): PV is rounded to the
		// cent twice, a relative change ≤ 2e-7, which moves the rate by less than 1e-6.
		a, err1 := PresentValueAnnuity(pmt, r, n)
		b, err2 := PresentValue(fv, r, n)
		if err1 != nil || err2 != nil {
			t.Fatal(err1, err2)
		}
		got, err := AnnuityRate(a+b, pmt, fv, n)
		wantFloat(t, "round trip", got, err, r, 1e-6)
		// Agreement with InternalRateOfReturn on the same flows (−PV, PMT, …, PMT + FV).
		if n <= 60 {
			flows := make([]int64, n+1)
			flows[0] = -(a + b)
			for i := 1; i <= n; i++ {
				flows[i] = pmt
			}
			flows[n] += fv
			irr, err := InternalRateOfReturn(flows)
			wantFloat(t, "agrees with IRR", got, err, irr, 1e-10)
		}
	}
}

func TestAnnuityRate_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "AnnuityRate", func() { sinkFloat, _ = AnnuityRate(95000, 2500, 100000, 20) })
}

func BenchmarkAnnuityRate(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = AnnuityRate(95000, 2500, 100000, 20)
	}
}

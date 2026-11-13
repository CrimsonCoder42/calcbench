package interest

import (
	"math"
	"testing"
)

func TestAnnualPercentageRate_SourceExample(t *testing.T) {
	// 12 CFR 1026 Appendix J (c)(1), Example (i): amount advanced $5,000, 24 monthly payments of $230,
	// regular first period (t = 1, f = 0), w = 12 → APR 9.69 %.
	got, err := AnnualPercentageRate(500000, 23000, 24, 12)
	if err != nil {
		t.Fatal(err)
	}
	if pct := math.Round(got*10000) / 100; pct != 9.69 {
		t.Errorf("App. J (c)(1)(i): APR %v (%v %%), want 9.69 %%", got, pct)
	}
}

func TestAnnualPercentageRate_Table(t *testing.T) {
	got, err := AnnualPercentageRate(120000, 10000, 12, 12)
	wantFloat(t, "payments sum to the amount: 0 %", got, err, 0, 1e-12)
	// One payment of $1,050 a year for $1,000: the unit-period rate is 5 %.
	got, err = AnnualPercentageRate(100000, 105000, 1, 1)
	wantFloat(t, "one annual payment", got, err, 0.05, 1e-12)
	// Payments below the amount financed give a negative rate.
	got, err = AnnualPercentageRate(120000, 9000, 12, 12)
	if err != nil || got >= 0 {
		t.Errorf("payments short of the amount: got %v, %v; want a negative rate", got, err)
	}
	// Weekly payments: the unit-period is the week, w = 52.
	got, err = AnnualPercentageRate(100000, 2000, 52, 52)
	if err != nil || got <= 0 {
		t.Errorf("weekly: got %v, %v", got, err)
	}
}

func TestAnnualPercentageRate_Errors(t *testing.T) {
	_, err := AnnualPercentageRate(maxExactCents+1, 100, 12, 12)
	wantErr(t, "amount 2^53+1", err, ErrOutOfRange)
	_, err = AnnualPercentageRate(100, maxExactCents+1, 12, 12)
	wantErr(t, "payment 2^53+1", err, ErrOutOfRange)
	_, err = AnnualPercentageRate(0, 100, 12, 12)
	wantErr(t, "zero amount", err, ErrDomain)
	_, err = AnnualPercentageRate(100, -1, 12, 12)
	wantErr(t, "negative payment", err, ErrDomain)
	_, err = AnnualPercentageRate(100, 10, 0, 12)
	wantErr(t, "zero payments", err, ErrPeriods)
	_, err = AnnualPercentageRate(100, 10, 12, 0)
	wantErr(t, "zero periods per year", err, ErrPeriods)
	// One payment of 1 cent for 2^53 cents needs a rate below −1 + 1e-12: no bracket.
	_, err = AnnualPercentageRate(maxExactCents, 1, 1, 12)
	wantErr(t, "no rate in range", err, ErrNoSolution)
}

func TestAnnualPercentageRate_Oracle(t *testing.T) {
	// Round trip: the payment computed from a rate gives that rate back (within the cent's rounding).
	rng := testRNG()
	for k := 0; k < 200; k++ {
		amount := 1_000_000 + rng.Int64N(100_000_000)
		apr := 0.01 + rng.Float64()*0.3
		n := 6 + rng.IntN(355)
		pmt, err := AnnuityPaymentFromPresentValue(amount, apr/12, n)
		if err != nil {
			t.Fatal(err)
		}
		got, err := AnnualPercentageRate(amount, pmt, n, 12)
		// The payment is rounded to ½ cent, a relative change of at most 1.6e-4 on a $10,000 loan; with
		// |d ln(P/A)/di| ≥ about (n+1)/2·0.8 that moves the APR by under 1e-4 (0.01 percentage point).
		wantFloat(t, "round trip", got, err, apr, 1e-4)
	}
}

func TestAnnualPercentageRate_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "AnnualPercentageRate", func() { sinkFloat, _ = AnnualPercentageRate(500000, 23000, 24, 12) })
}

func BenchmarkAnnualPercentageRate(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = AnnualPercentageRate(500000, 23000, 24, 12)
	}
}

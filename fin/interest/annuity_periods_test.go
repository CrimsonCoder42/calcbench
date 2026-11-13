package interest

import (
	"math"
	"testing"
)

func TestAnnuityPeriods_SourceExample(t *testing.T) {
	// Hand computation from (P/A,i,N) = PV/PMT (NCEES FE Reference Handbook, p. 230): PV $1,000,
	// i = 10 %, PMT $200 → 1 − 0.1·1000/200 = 0.5, N = ln 2/ln 1.1 = 7.272540897341713.
	got, err := AnnuityPeriods(100000, 20000, 0.1)
	wantFloat(t, "hand example", got, err, math.Ln2/math.Log(1.1), 1e-12)
}

func TestAnnuityPeriods_Table(t *testing.T) {
	// $200,000 at 0.5 %/month with the exact 30-year payment 1,199.101050...: N = 360.
	pmt := 20000000 * 0.005 / (1 - math.Pow(1.005, -360)) // 119910.105... cents
	n, err := AnnuityPeriods(20000000, int64(math.Round(pmt)), 0.005)
	wantFloat(t, "30-year mortgage, payment rounded to the cent", n, err, 360, 0.01)
	got, err := AnnuityPeriods(120000, 10000, 0)
	wantFloat(t, "zero rate is PV/PMT", got, err, 12, 0)
	got, err = AnnuityPeriods(0, 10000, 0.01)
	wantFloat(t, "nothing owed", got, err, 0, 0)
	got, err = AnnuityPeriods(100000, 100000, 0.05)
	wantFloat(t, "payment equal to PV takes just over one period", got, err, -math.Log1p(-0.05)/math.Log1p(0.05), 1e-15)
	got, err = AnnuityPeriods(100000, 50000, -0.1)
	wantFloat(t, "negative rate", got, err, -math.Log1p(0.2)/math.Log1p(-0.1), 1e-12)
	got, err = AnnuityPeriods(100000, 10000, 1e-15)
	wantFloat(t, "tiny rate tends to PV/PMT", got, err, 10, 1e-9)
}

func TestAnnuityPeriods_Errors(t *testing.T) {
	for _, r := range badRates {
		_, err := AnnuityPeriods(100000, 1000, r)
		wantErr(t, "bad rate", err, ErrRate)
	}
	_, err := AnnuityPeriods(-1, 1000, 0.01)
	wantErr(t, "negative PV", err, ErrNegative)
	_, err = AnnuityPeriods(100000, 0, 0.01)
	wantErr(t, "zero payment", err, ErrDomain)
	_, err = AnnuityPeriods(100000, -5, 0.01)
	wantErr(t, "negative payment", err, ErrDomain)
	_, err = AnnuityPeriods(100000, 1000, 0.01)
	wantErr(t, "payment equals interest", err, ErrNoSolution)
	_, err = AnnuityPeriods(100000, 999, 0.01)
	wantErr(t, "payment below interest", err, ErrNoSolution)
	_, err = AnnuityPeriods(maxExactCents+1, 1000, 0.01)
	wantErr(t, "PV beyond 2^53", err, ErrOutOfRange)
	_, err = AnnuityPeriods(1000, maxExactCents+1, 0.01)
	wantErr(t, "payment beyond 2^53", err, ErrOutOfRange)
}

func TestAnnuityPeriods_Oracle(t *testing.T) {
	// Definition: PMT·(1 − (1+i)^−N)/i = PV at the returned (fractional) N, to relative 1e-9.
	// Inputs are capped (i ≤ 2 %, payment ≥ 1.0001 × the first period's interest) so the inverse is well conditioned.
	rng := testRNG()
	for k := 0; k < 300; k++ {
		i := 1e-6 + rng.Float64()*0.02
		pv := int64(100_000 + rng.IntN(1_000_000_000))
		pmt := int64(math.Ceil(i*float64(pv)*(1.0001+rng.Float64()*5))) + 1
		n, err := AnnuityPeriods(pv, pmt, i)
		if err != nil {
			t.Fatalf("i=%v pv=%d pmt=%d: %v", i, pv, pmt, err)
		}
		back := float64(pmt) * -math.Expm1(-n*math.Log1p(i)) / i
		if math.Abs(back-float64(pv)) > 1e-9*float64(pv) {
			t.Errorf("i=%v pv=%d pmt=%d: N=%v gives PV %v", i, pv, pmt, n, back)
		}
	}
}

func TestAnnuityPeriods_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "AnnuityPeriods", func() { _, _ = AnnuityPeriods(20000000, 119910, 0.005) })
}

func BenchmarkAnnuityPeriods(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = AnnuityPeriods(20000000, 119910, 0.005)
	}
}

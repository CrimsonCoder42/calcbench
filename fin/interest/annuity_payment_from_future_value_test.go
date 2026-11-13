package interest

import (
	"math/big"
	"testing"
)

func TestAnnuityPaymentFromFutureValue_SourceExample(t *testing.T) {
	// NCEES FE Reference Handbook 10.0.1 p. 230: A = F·(A/F, i%, n). Hand-checked against QIA ch. 1
	// Figure 3: $5,525.63 at 5 % over 5 periods needs $1,000.00 a period.
	got, err := AnnuityPaymentFromFutureValue(552563, 0.05, 5)
	wantCents(t, "inverse of QIA Figure 3", got, err, 100000)
}

func TestAnnuityPaymentFromFutureValue_Table(t *testing.T) {
	cases := []struct {
		name    string
		amount  int64
		rate    float64
		periods int
		want    int64
	}{
		{"zero rate is FV/n", 120000, 0, 12, 10000},
		{"one period", 10000, 0.05, 1, 10000},
		{"negative amount", -552563, 0.05, 5, -100000},
	}
	for _, c := range cases {
		got, err := AnnuityPaymentFromFutureValue(c.amount, c.rate, c.periods)
		wantCents(t, c.name, got, err, c.want)
	}
}

func TestAnnuityPaymentFromFutureValue_Errors(t *testing.T) {
	_, err := AnnuityPaymentFromFutureValue(maxExactCents+1, 0.05, 5)
	wantErr(t, "2^53+1", err, ErrOutOfRange)
	_, err = AnnuityPaymentFromFutureValue(-maxExactCents-1, 0.05, 5)
	wantErr(t, "-2^53-1", err, ErrOutOfRange)
	for _, r := range badRates {
		_, err = AnnuityPaymentFromFutureValue(100, r, 5)
		wantErr(t, "bad rate", err, ErrRate)
	}
	_, err = AnnuityPaymentFromFutureValue(100, 0.05, 0)
	wantErr(t, "periods below minimum", err, ErrPeriods)
}

func TestAnnuityPaymentFromFutureValue_Oracle(t *testing.T) {
	rng := testRNG()
	for k := 0; k < 300; k++ {
		amount := rng.Int64N(2_000_000_000) - 1_000_000_000
		rate := rng.Float64()*0.15 - 0.03
		n := 1 + rng.IntN(120)
		got, err := AnnuityPaymentFromFutureValue(amount, rate, n)
		if err != nil {
			t.Fatalf("%d %v %d: %v", amount, rate, n, err)
		}
		f := bigFactors(rate, n)
		checkCentsOracle(t, "oracle", got, new(big.Float).SetPrec(bigPrec).Mul(bigI(amount), bigF(f.af)))
	}
}

func TestAnnuityPaymentFromFutureValue_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "AnnuityPaymentFromFutureValue", func() { sinkCents, _ = AnnuityPaymentFromFutureValue(552563, 0.05, 5) })
}

func BenchmarkAnnuityPaymentFromFutureValue(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkCents, _ = AnnuityPaymentFromFutureValue(552563, 0.05, 5)
	}
}

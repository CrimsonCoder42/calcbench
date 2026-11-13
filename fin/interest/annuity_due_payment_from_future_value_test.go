package interest

import (
	"math/big"
	"testing"
)

func TestAnnuityDuePaymentFromFutureValue_SourceExample(t *testing.T) {
	// Principles of Finance §8.2 Eq. 8.25 solved for PYMT: $5,801.91 at 5 % over 5 periods in advance
	// → $1,000.00 (the inverse of FutureValueAnnuityDue's example).
	got, err := AnnuityDuePaymentFromFutureValue(580191, 0.05, 5)
	wantCents(t, "inverse of Eq. 8.25", got, err, 100000)
}

func TestAnnuityDuePaymentFromFutureValue_Table(t *testing.T) {
	cases := []struct {
		name    string
		amount  int64
		rate    float64
		periods int
		want    int64
	}{
		{"zero rate is FV/n", 120000, 0, 12, 10000},
		{"one period", 10500, 0.05, 1, 10000},
		{"negative amount", -580191, 0.05, 5, -100000},
	}
	for _, c := range cases {
		got, err := AnnuityDuePaymentFromFutureValue(c.amount, c.rate, c.periods)
		wantCents(t, c.name, got, err, c.want)
	}
}

func TestAnnuityDuePaymentFromFutureValue_Errors(t *testing.T) {
	_, err := AnnuityDuePaymentFromFutureValue(maxExactCents+1, 0.05, 5)
	wantErr(t, "2^53+1", err, ErrOutOfRange)
	_, err = AnnuityDuePaymentFromFutureValue(-maxExactCents-1, 0.05, 5)
	wantErr(t, "-2^53-1", err, ErrOutOfRange)
	for _, r := range badRates {
		_, err = AnnuityDuePaymentFromFutureValue(100, r, 5)
		wantErr(t, "bad rate", err, ErrRate)
	}
	_, err = AnnuityDuePaymentFromFutureValue(100, 0.05, 0)
	wantErr(t, "periods below minimum", err, ErrPeriods)
	_, err = AnnuityDuePaymentFromFutureValue(maxExactCents, -0.99, 1)
	wantErr(t, "result above 2^53", err, ErrOverflow)
}

func TestAnnuityDuePaymentFromFutureValue_Oracle(t *testing.T) {
	rng := testRNG()
	for k := 0; k < 300; k++ {
		amount := rng.Int64N(2_000_000_000) - 1_000_000_000
		rate := rng.Float64()*0.15 - 0.03
		n := 1 + rng.IntN(120)
		got, err := AnnuityDuePaymentFromFutureValue(amount, rate, n)
		if err != nil {
			t.Fatalf("%d %v %d: %v", amount, rate, n, err)
		}
		f := bigFactors(rate, n)
		checkCentsOracle(t, "oracle", got, new(big.Float).SetPrec(bigPrec).Mul(bigI(amount), bigF(f.af/(1+rate))))
	}
}

func TestAnnuityDuePaymentFromFutureValue_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "AnnuityDuePaymentFromFutureValue", func() { sinkCents, _ = AnnuityDuePaymentFromFutureValue(580191, 0.05, 5) })
}

func BenchmarkAnnuityDuePaymentFromFutureValue(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkCents, _ = AnnuityDuePaymentFromFutureValue(580191, 0.05, 5)
	}
}

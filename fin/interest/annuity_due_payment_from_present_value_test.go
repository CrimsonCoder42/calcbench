package interest

import (
	"math/big"
	"testing"
)

func TestAnnuityDuePaymentFromPresentValue_SourceExample(t *testing.T) {
	// QIA ch. 1 §6.1 footnote 7 solved for A: $340.18 at 12 % over 4 periods in advance → $100.00.
	got, err := AnnuityDuePaymentFromPresentValue(34018, 0.12, 4)
	wantCents(t, "inverse of QIA §6.1", got, err, 10000)
}

func TestAnnuityDuePaymentFromPresentValue_Table(t *testing.T) {
	cases := []struct {
		name    string
		amount  int64
		rate    float64
		periods int
		want    int64
	}{
		{"zero rate is PV/n", 120000, 0, 12, 10000},
		{"one period is paid now", 10000, 0.05, 1, 10000},
		{"negative amount", -34018, 0.12, 4, -10000},
	}
	for _, c := range cases {
		got, err := AnnuityDuePaymentFromPresentValue(c.amount, c.rate, c.periods)
		wantCents(t, c.name, got, err, c.want)
	}
}

func TestAnnuityDuePaymentFromPresentValue_Errors(t *testing.T) {
	_, err := AnnuityDuePaymentFromPresentValue(maxExactCents+1, 0.05, 5)
	wantErr(t, "2^53+1", err, ErrOutOfRange)
	_, err = AnnuityDuePaymentFromPresentValue(-maxExactCents-1, 0.05, 5)
	wantErr(t, "-2^53-1", err, ErrOutOfRange)
	for _, r := range badRates {
		_, err = AnnuityDuePaymentFromPresentValue(100, r, 5)
		wantErr(t, "bad rate", err, ErrRate)
	}
	_, err = AnnuityDuePaymentFromPresentValue(100, 0.05, 0)
	wantErr(t, "periods below minimum", err, ErrPeriods)
}

func TestAnnuityDuePaymentFromPresentValue_Oracle(t *testing.T) {
	rng := testRNG()
	for k := 0; k < 300; k++ {
		amount := rng.Int64N(2_000_000_000) - 1_000_000_000
		rate := rng.Float64()*0.15 - 0.03
		n := 1 + rng.IntN(120)
		got, err := AnnuityDuePaymentFromPresentValue(amount, rate, n)
		if err != nil {
			t.Fatalf("%d %v %d: %v", amount, rate, n, err)
		}
		f := bigFactors(rate, n)
		checkCentsOracle(t, "oracle", got, new(big.Float).SetPrec(bigPrec).Mul(bigI(amount), bigF(f.ap/(1+rate))))
	}
}

func TestAnnuityDuePaymentFromPresentValue_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "AnnuityDuePaymentFromPresentValue", func() { sinkCents, _ = AnnuityDuePaymentFromPresentValue(34018, 0.12, 4) })
}

func BenchmarkAnnuityDuePaymentFromPresentValue(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkCents, _ = AnnuityDuePaymentFromPresentValue(34018, 0.12, 4)
	}
}

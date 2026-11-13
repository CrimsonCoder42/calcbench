package interest

import (
	"math/big"
	"testing"
)

func TestGrowingAnnuityPaymentFromPresentValue_SourceExample(t *testing.T) {
	// Inverse of Damodaran App. 3's gold mine: $16,145,979.98 at 10 %, growing 3 %, 20 years
	// → first payment $1,545,000.00 (hand-checked).
	got, err := GrowingAnnuityPaymentFromPresentValue(1614597998, 0.10, 0.03, 20)
	wantCents(t, "inverse gold mine", got, err, 154500000)
}

func TestGrowingAnnuityPaymentFromPresentValue_Table(t *testing.T) {
	cases := []struct {
		name    string
		pv      int64
		r, g    float64
		periods int
		want    int64
	}{
		{"one period", 10000, 0.05, 0.02, 1, 10500},
		{"r = g gives PV·(1+r)/n", 95238, 0.05, 0.05, 10, 10000},
		{"zero growth is the level payment", 10000000, 0.08 / 12, 0, 360, 73376},
		{"negative amount", -1614597998, 0.10, 0.03, 20, -154500000},
		{"zero amount", 0, 0.10, 0.03, 20, 0},
	}
	for _, c := range cases {
		got, err := GrowingAnnuityPaymentFromPresentValue(c.pv, c.r, c.g, c.periods)
		wantCents(t, c.name, got, err, c.want)
	}
}

func TestGrowingAnnuityPaymentFromPresentValue_Errors(t *testing.T) {
	_, err := GrowingAnnuityPaymentFromPresentValue(maxExactCents+1, 0.05, 0.02, 5)
	wantErr(t, "2^53+1", err, ErrOutOfRange)
	for _, r := range badRates {
		_, err = GrowingAnnuityPaymentFromPresentValue(100, r, 0.02, 5)
		wantErr(t, "bad rate", err, ErrRate)
		_, err = GrowingAnnuityPaymentFromPresentValue(100, 0.05, r, 5)
		wantErr(t, "bad growth", err, ErrRate)
	}
	_, err = GrowingAnnuityPaymentFromPresentValue(100, 0.05, 0.02, 0)
	wantErr(t, "zero periods", err, ErrPeriods)
	_, err = GrowingAnnuityPaymentFromPresentValue(maxExactCents, 1, 0, 1)
	wantErr(t, "overflow", err, ErrOverflow)
}

func TestGrowingAnnuityPaymentFromPresentValue_Oracle(t *testing.T) {
	rng := testRNG()
	for k := 0; k < 300; k++ {
		pv := rng.Int64N(1_000_000_000)
		r := rng.Float64() * 0.15
		g := rng.Float64()*0.12 - 0.02
		n := 1 + rng.IntN(60)
		got, err := GrowingAnnuityPaymentFromPresentValue(pv, r, g, n)
		if err != nil {
			t.Fatal(err)
		}
		perUnit := bigGrowingAnnuityPV(1, r, g, n)
		checkCentsOracle(t, "oracle", got, new(big.Float).SetPrec(bigPrec).Quo(bigI(pv), perUnit))
	}
}

func TestGrowingAnnuityPaymentFromPresentValue_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "GrowingAnnuityPaymentFromPresentValue", func() {
		sinkCents, _ = GrowingAnnuityPaymentFromPresentValue(1614597998, 0.10, 0.03, 20)
	})
}

func BenchmarkGrowingAnnuityPaymentFromPresentValue(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkCents, _ = GrowingAnnuityPaymentFromPresentValue(1614597998, 0.10, 0.03, 20)
	}
}

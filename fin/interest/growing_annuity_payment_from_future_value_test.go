package interest

import (
	"math/big"
	"testing"
)

func TestGrowingAnnuityPaymentFromFutureValue_SourceExample(t *testing.T) {
	// Inverse of FutureValueGrowingAnnuity on Damodaran App. 3's gold mine: $108,622,079.49 at 10 %,
	// growing 3 %, 20 years → first payment $1,545,000.00 (hand-checked).
	got, err := GrowingAnnuityPaymentFromFutureValue(10862207949, 0.10, 0.03, 20)
	wantCents(t, "inverse gold mine", got, err, 154500000)
}

func TestGrowingAnnuityPaymentFromFutureValue_Table(t *testing.T) {
	cases := []struct {
		name    string
		fv      int64
		r, g    float64
		periods int
		want    int64
	}{
		{"one period", 10000, 0.05, 0.02, 1, 10000},
		{"two periods", 21300, 0.10, 0.03, 2, 10000},
		{"r = g", 22000, 0.10, 0.10, 2, 10000},
		{"zero growth is the sinking fund", 552563, 0.05, 0, 5, 100000},
		{"negative amount", -21300, 0.10, 0.03, 2, -10000},
	}
	for _, c := range cases {
		got, err := GrowingAnnuityPaymentFromFutureValue(c.fv, c.r, c.g, c.periods)
		wantCents(t, c.name, got, err, c.want)
	}
}

func TestGrowingAnnuityPaymentFromFutureValue_Errors(t *testing.T) {
	_, err := GrowingAnnuityPaymentFromFutureValue(maxExactCents+1, 0.05, 0.02, 5)
	wantErr(t, "2^53+1", err, ErrOutOfRange)
	for _, r := range badRates {
		_, err = GrowingAnnuityPaymentFromFutureValue(100, r, 0.02, 5)
		wantErr(t, "bad rate", err, ErrRate)
		_, err = GrowingAnnuityPaymentFromFutureValue(100, 0.05, r, 5)
		wantErr(t, "bad growth", err, ErrRate)
	}
	_, err = GrowingAnnuityPaymentFromFutureValue(100, 0.05, 0.02, 0)
	wantErr(t, "zero periods", err, ErrPeriods)
	// Payments C and 0.01·C at a rate of −99 %: FV = 0.02·C, so C = 50·FV.
	_, err = GrowingAnnuityPaymentFromFutureValue(maxExactCents, -0.99, -0.99, 2)
	wantErr(t, "overflow", err, ErrOverflow)
}

func TestGrowingAnnuityPaymentFromFutureValue_Oracle(t *testing.T) {
	rng := testRNG()
	for k := 0; k < 300; k++ {
		fv := rng.Int64N(1_000_000_000)
		r := rng.Float64() * 0.15
		g := rng.Float64()*0.12 - 0.02
		n := 1 + rng.IntN(60)
		got, err := GrowingAnnuityPaymentFromFutureValue(fv, r, g, n)
		if err != nil {
			t.Fatal(err)
		}
		perUnit := new(big.Float).SetPrec(bigPrec).Mul(bigGrowingAnnuityPV(1, r, g, n), bigCompound(r, n))
		checkCentsOracle(t, "oracle", got, new(big.Float).SetPrec(bigPrec).Quo(bigI(fv), perUnit))
	}
}

func TestGrowingAnnuityPaymentFromFutureValue_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "GrowingAnnuityPaymentFromFutureValue", func() {
		sinkCents, _ = GrowingAnnuityPaymentFromFutureValue(10862207949, 0.10, 0.03, 20)
	})
}

func BenchmarkGrowingAnnuityPaymentFromFutureValue(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkCents, _ = GrowingAnnuityPaymentFromFutureValue(10862207949, 0.10, 0.03, 20)
	}
}

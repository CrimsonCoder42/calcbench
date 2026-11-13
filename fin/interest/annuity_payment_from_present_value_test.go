package interest

import (
	"math/big"
	"testing"
)

func TestAnnuityPaymentFromPresentValue_SourceExample(t *testing.T) {
	// DeFusco et al., Quantitative Investment Analysis 3rd ed. (Wiley/CFA Institute 2015) ch. 1 Example 20: $100,000 at 8 %/12 for 360 months → $733.76.
	got, err := AnnuityPaymentFromPresentValue(10000000, 0.08/12, 360)
	wantCents(t, "QIA Example 20", got, err, 73376)
	// Principles of Finance (OpenStax 2022) §8.3: $32,000 at 6 %/12 for 36 months → $973.50;
	// $140,000 at 3.6 %/12 for 240 months → $819.16.
	got, err = AnnuityPaymentFromPresentValue(3200000, 0.06/12, 36)
	wantCents(t, "§8.3 auto loan", got, err, 97350)
	got, err = AnnuityPaymentFromPresentValue(14000000, 0.036/12, 240)
	wantCents(t, "§8.3 mortgage", got, err, 81916)
	// §8.2 Eq. 8.30 divides $75,000 by the factor rounded to 4.579705 and shows $16,376.60; the exact
	// payment is $16,376.5929, so the half-even cent is $16,376.59.
	got, err = AnnuityPaymentFromPresentValue(7500000, 0.03, 5)
	wantCents(t, "§8.2 exact", got, err, 1637659)
}

func TestAnnuityPaymentFromPresentValue_Table(t *testing.T) {
	cases := []struct {
		name    string
		amount  int64
		rate    float64
		periods int
		want    int64
	}{
		{"zero rate is PV/n", 120000, 0, 12, 10000},
		{"one period", 10000, 0.05, 1, 10500},
		{"negative amount", -10000000, 0.08 / 12, 360, -73376},
		{"zero amount", 0, 0.05, 10, 0},
	}
	for _, c := range cases {
		got, err := AnnuityPaymentFromPresentValue(c.amount, c.rate, c.periods)
		wantCents(t, c.name, got, err, c.want)
	}
}

func TestAnnuityPaymentFromPresentValue_Errors(t *testing.T) {
	_, err := AnnuityPaymentFromPresentValue(maxExactCents+1, 0.05, 5)
	wantErr(t, "2^53+1", err, ErrOutOfRange)
	_, err = AnnuityPaymentFromPresentValue(-maxExactCents-1, 0.05, 5)
	wantErr(t, "-2^53-1", err, ErrOutOfRange)
	for _, r := range badRates {
		_, err = AnnuityPaymentFromPresentValue(100, r, 5)
		wantErr(t, "bad rate", err, ErrRate)
	}
	_, err = AnnuityPaymentFromPresentValue(100, 0.05, 0)
	wantErr(t, "periods below minimum", err, ErrPeriods)
	_, err = AnnuityPaymentFromPresentValue(maxExactCents, 0.5, 1)
	wantErr(t, "result above 2^53", err, ErrOverflow)
}

func TestAnnuityPaymentFromPresentValue_Oracle(t *testing.T) {
	rng := testRNG()
	for k := 0; k < 300; k++ {
		amount := rng.Int64N(2_000_000_000) - 1_000_000_000
		rate := rng.Float64()*0.15 - 0.03
		n := 1 + rng.IntN(120)
		got, err := AnnuityPaymentFromPresentValue(amount, rate, n)
		if err != nil {
			t.Fatalf("%d %v %d: %v", amount, rate, n, err)
		}
		f := bigFactors(rate, n)
		checkCentsOracle(t, "oracle", got, new(big.Float).SetPrec(bigPrec).Mul(bigI(amount), bigF(f.ap)))
	}
}

func TestAnnuityPaymentFromPresentValue_RepaysToTheCent(t *testing.T) {
	// Amortize with the rounded payment and half-even interest each period. The payment's rounding
	// (≤ ½ cent) and each period's interest rounding (≤ ½ cent) compound, so the balance left after
	// the last payment is at most (F/A,i,n) cents from 0, the residual a schedule's last payment absorbs.
	for _, c := range []struct {
		pv   int64
		rate float64
		n    int
	}{{10000000, 0.08 / 12, 360}, {3200000, 0.005, 36}, {14000000, 0.003, 240}, {99999, 0, 7}} {
		pmt, err := AnnuityPaymentFromPresentValue(c.pv, c.rate, c.n)
		if err != nil {
			t.Fatal(err)
		}
		bal := c.pv
		for k := 0; k < c.n; k++ {
			interest, err := roundCents(float64(bal) * c.rate)
			if err != nil {
				t.Fatal(err)
			}
			bal += interest - pmt
		}
		fa, err := UniformSeriesCompoundAmount(c.rate, c.n)
		if err != nil {
			t.Fatal(err)
		}
		if limit := int64(fa) + 1; bal < -limit || bal > limit {
			t.Errorf("pv %d rate %v n %d: residual %d cents after %d payments of %d", c.pv, c.rate, c.n, bal, c.n, pmt)
		}
	}
}

func TestAnnuityPaymentFromPresentValue_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "AnnuityPaymentFromPresentValue", func() { sinkCents, _ = AnnuityPaymentFromPresentValue(10000000, 0.08/12, 360) })
}

func BenchmarkAnnuityPaymentFromPresentValue(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkCents, _ = AnnuityPaymentFromPresentValue(10000000, 0.08/12, 360)
	}
}

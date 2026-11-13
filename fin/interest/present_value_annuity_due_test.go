package interest

import (
	"math/big"
	"testing"
)

func TestPresentValueAnnuityDue_SourceExample(t *testing.T) {
	// DeFusco et al., Quantitative Investment Analysis 3rd ed. (Wiley/CFA Institute 2015) ch. 1 §6.1: four payments of $100, the first today, at 12 % → $340.18.
	got, err := PresentValueAnnuityDue(10000, 0.12, 4)
	wantCents(t, "QIA §6.1", got, err, 34018)
}

func TestPresentValueAnnuityDue_Table(t *testing.T) {
	cases := []struct {
		name    string
		amount  int64
		rate    float64
		periods int
		want    int64
	}{
		{"zero rate is PMT·n", 1000, 0, 12, 12000},
		{"zero periods", 1000, 0.05, 0, 0},
		{"one period is not discounted", 1000, 0.05, 1, 1000},
		{"negative payment", -10000, 0.12, 4, -34018},
	}
	for _, c := range cases {
		got, err := PresentValueAnnuityDue(c.amount, c.rate, c.periods)
		wantCents(t, c.name, got, err, c.want)
	}
}

func TestPresentValueAnnuityDue_Errors(t *testing.T) {
	_, err := PresentValueAnnuityDue(maxExactCents+1, 0.05, 5)
	wantErr(t, "2^53+1", err, ErrOutOfRange)
	_, err = PresentValueAnnuityDue(-maxExactCents-1, 0.05, 5)
	wantErr(t, "-2^53-1", err, ErrOutOfRange)
	for _, r := range badRates {
		_, err = PresentValueAnnuityDue(100, r, 5)
		wantErr(t, "bad rate", err, ErrRate)
	}
	_, err = PresentValueAnnuityDue(100, 0.05, -1)
	wantErr(t, "periods below minimum", err, ErrPeriods)
	_, err = PresentValueAnnuityDue(maxExactCents, -0.5, 2)
	wantErr(t, "result above 2^53", err, ErrOverflow)
}

func TestPresentValueAnnuityDue_Oracle(t *testing.T) {
	rng := testRNG()
	for k := 0; k < 300; k++ {
		amount := rng.Int64N(2_000_000_000) - 1_000_000_000
		rate := rng.Float64()*0.15 - 0.03
		n := 1 + rng.IntN(120)
		got, err := PresentValueAnnuityDue(amount, rate, n)
		if err != nil {
			t.Fatalf("%d %v %d: %v", amount, rate, n, err)
		}
		f := bigFactors(rate, n)
		checkCentsOracle(t, "oracle", got, new(big.Float).SetPrec(bigPrec).Mul(bigI(amount), bigF(f.pa*(1+rate))))
	}
}

func TestPresentValueAnnuityDue_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "PresentValueAnnuityDue", func() { sinkCents, _ = PresentValueAnnuityDue(10000, 0.12, 4) })
}

func BenchmarkPresentValueAnnuityDue(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkCents, _ = PresentValueAnnuityDue(10000, 0.12, 4)
	}
}

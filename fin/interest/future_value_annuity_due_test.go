package interest

import (
	"math/big"
	"testing"
)

func TestFutureValueAnnuityDue_SourceExample(t *testing.T) {
	// Principles of Finance (OpenStax 2022) §8.2, Eq. 8.25: annuity due = ordinary annuity × (1 + i).
	// QIA ch. 1 Figure 3's $1,000 for 5 years at 5 % ($5,525.63) becomes $5,525.63125 × 1.05 = $5,801.91.
	got, err := FutureValueAnnuityDue(100000, 0.05, 5)
	wantCents(t, "Eq. 8.25 on QIA Figure 3", got, err, 580191)
}

func TestFutureValueAnnuityDue_Table(t *testing.T) {
	cases := []struct {
		name    string
		amount  int64
		rate    float64
		periods int
		want    int64
	}{
		{"zero rate is PMT·n", 1000, 0, 12, 12000},
		{"zero periods", 1000, 0.05, 0, 0},
		{"one period earns one period", 1000, 0.05, 1, 1050},
		{"negative payment", -100000, 0.05, 5, -580191},
	}
	for _, c := range cases {
		got, err := FutureValueAnnuityDue(c.amount, c.rate, c.periods)
		wantCents(t, c.name, got, err, c.want)
	}
}

func TestFutureValueAnnuityDue_Errors(t *testing.T) {
	_, err := FutureValueAnnuityDue(maxExactCents+1, 0.05, 5)
	wantErr(t, "2^53+1", err, ErrOutOfRange)
	_, err = FutureValueAnnuityDue(-maxExactCents-1, 0.05, 5)
	wantErr(t, "-2^53-1", err, ErrOutOfRange)
	for _, r := range badRates {
		_, err = FutureValueAnnuityDue(100, r, 5)
		wantErr(t, "bad rate", err, ErrRate)
	}
	_, err = FutureValueAnnuityDue(100, 0.05, -1)
	wantErr(t, "periods below minimum", err, ErrPeriods)
	_, err = FutureValueAnnuityDue(maxExactCents, 0.05, 2)
	wantErr(t, "result above 2^53", err, ErrOverflow)
}

func TestFutureValueAnnuityDue_Oracle(t *testing.T) {
	rng := testRNG()
	for k := 0; k < 300; k++ {
		amount := rng.Int64N(2_000_000) - 1_000_000 // |FV| ≤ 1e6·(F/A at 12 %, 120)·1.12 < 2⁴³
		rate := rng.Float64()*0.15 - 0.03
		n := 1 + rng.IntN(120)
		got, err := FutureValueAnnuityDue(amount, rate, n)
		if err != nil {
			t.Fatalf("%d %v %d: %v", amount, rate, n, err)
		}
		f := bigFactors(rate, n)
		checkCentsOracle(t, "oracle", got, new(big.Float).SetPrec(bigPrec).Mul(bigI(amount), bigF(f.fa*(1+rate))))
	}
}

func TestFutureValueAnnuityDue_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "FutureValueAnnuityDue", func() { sinkCents, _ = FutureValueAnnuityDue(100000, 0.05, 5) })
}

func BenchmarkFutureValueAnnuityDue(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkCents, _ = FutureValueAnnuityDue(100000, 0.05, 5)
	}
}

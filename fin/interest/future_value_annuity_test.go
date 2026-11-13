package interest

import (
	"math/big"
	"testing"
)

func TestFutureValueAnnuity_SourceExample(t *testing.T) {
	// DeFusco et al., Quantitative Investment Analysis 3rd ed. (Wiley/CFA Institute 2015) ch. 1 Example 7: €20,000 a year for 30 years at 9 % → €2,726,150.77.
	got, err := FutureValueAnnuity(2000000, 0.09, 30)
	wantCents(t, "QIA Example 7", got, err, 272615077)
	// QIA ch. 1 Figure 3: $1,000 a year for 5 years at 5 % → $5,525.63.
	got, err = FutureValueAnnuity(100000, 0.05, 5)
	wantCents(t, "QIA Figure 3", got, err, 552563)
	// Principles of Finance (OpenStax 2022) §8.2 shows $3,000 for 5 years at 4 % as $16,248.98; the exact
	// value is $16,248.9677, so the half-even cent is $16,248.97 (the book rounds its factor first).
	got, err = FutureValueAnnuity(300000, 0.04, 5)
	wantCents(t, "§8.2 exact", got, err, 1624897)
}

func TestFutureValueAnnuity_Table(t *testing.T) {
	cases := []struct {
		name    string
		amount  int64
		rate    float64
		periods int
		want    int64
	}{
		{"zero rate is PMT·n", 1000, 0, 12, 12000},
		{"zero periods", 1000, 0.05, 0, 0},
		{"one period", 1000, 0.05, 1, 1000},
		{"negative payment", -100000, 0.05, 5, -552563},
		{"zero payment", 0, 0.05, 5, 0},
	}
	for _, c := range cases {
		got, err := FutureValueAnnuity(c.amount, c.rate, c.periods)
		wantCents(t, c.name, got, err, c.want)
	}
}

func TestFutureValueAnnuity_Errors(t *testing.T) {
	_, err := FutureValueAnnuity(maxExactCents+1, 0.05, 5)
	wantErr(t, "2^53+1", err, ErrOutOfRange)
	_, err = FutureValueAnnuity(-maxExactCents-1, 0.05, 5)
	wantErr(t, "-2^53-1", err, ErrOutOfRange)
	for _, r := range badRates {
		_, err = FutureValueAnnuity(100, r, 5)
		wantErr(t, "bad rate", err, ErrRate)
	}
	_, err = FutureValueAnnuity(100, 0.05, -1)
	wantErr(t, "periods below minimum", err, ErrPeriods)
	_, err = FutureValueAnnuity(maxExactCents, 0.05, 2)
	wantErr(t, "result above 2^53", err, ErrOverflow)
}

func TestFutureValueAnnuity_Oracle(t *testing.T) {
	rng := testRNG()
	for k := 0; k < 300; k++ {
		amount := rng.Int64N(2_000_000) - 1_000_000 // |FV| ≤ 1e6·(F/A at 12 %, 120) < 2⁴³
		rate := rng.Float64()*0.15 - 0.03
		n := 1 + rng.IntN(120)
		got, err := FutureValueAnnuity(amount, rate, n)
		if err != nil {
			t.Fatalf("%d %v %d: %v", amount, rate, n, err)
		}
		f := bigFactors(rate, n)
		checkCentsOracle(t, "oracle", got, new(big.Float).SetPrec(bigPrec).Mul(bigI(amount), bigF(f.fa)))
	}
}

func TestFutureValueAnnuity_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "FutureValueAnnuity", func() { sinkCents, _ = FutureValueAnnuity(2000000, 0.09, 30) })
}

func BenchmarkFutureValueAnnuity(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkCents, _ = FutureValueAnnuity(2000000, 0.09, 30)
	}
}

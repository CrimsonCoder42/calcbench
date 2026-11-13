package interest

import (
	"math/big"
	"testing"
)

func TestPresentValueAnnuity_SourceExample(t *testing.T) {
	// DeFusco et al., Quantitative Investment Analysis 3rd ed. (Wiley/CFA Institute 2015) ch. 1 Example 11: €1,000 a year for 5 years at 12 % → €3,604.78.
	got, err := PresentValueAnnuity(100000, 0.12, 5)
	wantCents(t, "QIA Example 11", got, err, 360478)
	// Principles of Finance (OpenStax 2022) §8.3: the auto loan's 36 payments of $973.50 at 0.5 % a month
	// are worth $32,000 (to the cent's rounding of the payment: $31,999.93).
	got, err = PresentValueAnnuity(97350, 0.005, 36)
	wantCents(t, "§8.3 auto loan", got, err, 3199993)
}

func TestPresentValueAnnuity_Table(t *testing.T) {
	cases := []struct {
		name    string
		amount  int64
		rate    float64
		periods int
		want    int64
	}{
		{"zero rate is PMT·n", 1000, 0, 12, 12000},
		{"zero periods", 1000, 0.05, 0, 0},
		{"one period", 10500, 0.05, 1, 10000},
		{"negative payment", -100000, 0.12, 5, -360478},
	}
	for _, c := range cases {
		got, err := PresentValueAnnuity(c.amount, c.rate, c.periods)
		wantCents(t, c.name, got, err, c.want)
	}
}

func TestPresentValueAnnuity_Errors(t *testing.T) {
	_, err := PresentValueAnnuity(maxExactCents+1, 0.05, 5)
	wantErr(t, "2^53+1", err, ErrOutOfRange)
	_, err = PresentValueAnnuity(-maxExactCents-1, 0.05, 5)
	wantErr(t, "-2^53-1", err, ErrOutOfRange)
	for _, r := range badRates {
		_, err = PresentValueAnnuity(100, r, 5)
		wantErr(t, "bad rate", err, ErrRate)
	}
	_, err = PresentValueAnnuity(100, 0.05, -1)
	wantErr(t, "periods below minimum", err, ErrPeriods)
	_, err = PresentValueAnnuity(maxExactCents, -0.5, 2)
	wantErr(t, "result above 2^53", err, ErrOverflow)
}

func TestPresentValueAnnuity_Oracle(t *testing.T) {
	rng := testRNG()
	for k := 0; k < 300; k++ {
		amount := rng.Int64N(2_000_000_000) - 1_000_000_000
		rate := rng.Float64()*0.15 - 0.03
		n := 1 + rng.IntN(120)
		got, err := PresentValueAnnuity(amount, rate, n)
		if err != nil {
			t.Fatalf("%d %v %d: %v", amount, rate, n, err)
		}
		f := bigFactors(rate, n)
		checkCentsOracle(t, "oracle", got, new(big.Float).SetPrec(bigPrec).Mul(bigI(amount), bigF(f.pa)))
	}
}

func TestPresentValueAnnuity_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "PresentValueAnnuity", func() { sinkCents, _ = PresentValueAnnuity(100000, 0.12, 5) })
}

func BenchmarkPresentValueAnnuity(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkCents, _ = PresentValueAnnuity(100000, 0.12, 5)
	}
}

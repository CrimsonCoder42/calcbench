package interest

import (
	"math/big"
	"testing"
)

// bigGrowingAnnuityPV returns C·(1 − ((1+g)/(1+r))ⁿ)/(r − g) at 256 bits (r ≠ g), from Damodaran,
// Applied Corporate Finance 3rd ed., App. 3 "Growing Annuities", written per first payment C = A(1+g).
func bigGrowingAnnuityPV(c int64, r, g float64, n int) *big.Float {
	q := func() *big.Float { return new(big.Float).SetPrec(bigPrec) }
	ratio := q().Quo(bigOnePlus(g), bigOnePlus(r))
	num := q().Sub(bigI(1), bigPow(ratio, n))
	den := q().Sub(bigF(r), bigF(g))
	return q().Mul(bigI(c), q().Quo(num, den))
}

func TestPresentValueGrowingAnnuity_SourceExample(t *testing.T) {
	// Damodaran, Applied Corporate Finance 3rd ed. (Wiley 2010) App. 3, gold mine: 5,000 oz a year for
	// 20 years at $300 now, growing 3 %, discounted at 10 % → $16,145,980 (to the dollar). The first
	// payment is 5,000 × $300 × 1.03 = $1,545,000; to the cent the value is $16,145,979.98.
	got, err := PresentValueGrowingAnnuity(154500000, 0.10, 0.03, 20)
	wantCents(t, "gold mine", got, err, 1614597998)
}

func TestPresentValueGrowingAnnuity_Table(t *testing.T) {
	cases := []struct {
		name    string
		c       int64
		r, g    float64
		periods int
		want    int64
	}{
		{"zero periods", 10000, 0.05, 0.02, 0, 0},
		{"one period", 10500, 0.05, 0.02, 1, 10000},
		{"r = g gives C·n/(1+r)", 10000, 0.05, 0.05, 10, 95238},
		{"g > r", 10000, 0.02, 0.05, 3, 30285},
		{"zero rate and growth", 1000, 0, 0, 12, 12000},
		{"zero growth is the level annuity", 100000, 0.12, 0, 5, 360478},
		{"negative payment", -154500000, 0.10, 0.03, 20, -1614597998},
		{"zero payment", 0, 0.1, 0.03, 20, 0},
	}
	for _, c := range cases {
		got, err := PresentValueGrowingAnnuity(c.c, c.r, c.g, c.periods)
		wantCents(t, c.name, got, err, c.want)
	}
}

func TestPresentValueGrowingAnnuity_Errors(t *testing.T) {
	_, err := PresentValueGrowingAnnuity(maxExactCents+1, 0.05, 0.02, 5)
	wantErr(t, "2^53+1", err, ErrOutOfRange)
	for _, r := range badRates {
		_, err = PresentValueGrowingAnnuity(100, r, 0.02, 5)
		wantErr(t, "bad rate", err, ErrRate)
		_, err = PresentValueGrowingAnnuity(100, 0.05, r, 5)
		wantErr(t, "bad growth", err, ErrRate)
	}
	_, err = PresentValueGrowingAnnuity(100, 0.05, 0.02, -1)
	wantErr(t, "negative periods", err, ErrPeriods)
	_, err = PresentValueGrowingAnnuity(maxExactCents, 0.05, 0.05, 2)
	wantErr(t, "overflow", err, ErrOverflow)
}

func TestPresentValueGrowingAnnuity_Oracle(t *testing.T) {
	rng := testRNG()
	for k := 0; k < 300; k++ {
		c := rng.Int64N(100_000_000)
		r := rng.Float64() * 0.15
		g := rng.Float64()*0.12 - 0.02
		if k%10 == 0 {
			g = r + 1e-9 // near r = g, where the textbook form cancels
		}
		n := 1 + rng.IntN(60)
		got, err := PresentValueGrowingAnnuity(c, r, g, n)
		if err != nil {
			t.Fatal(err)
		}
		checkCentsOracle(t, "oracle", got, bigGrowingAnnuityPV(c, r, g, n))
	}
}

func TestPresentValueGrowingAnnuity_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "PresentValueGrowingAnnuity", func() { sinkCents, _ = PresentValueGrowingAnnuity(154500000, 0.10, 0.03, 20) })
}

func BenchmarkPresentValueGrowingAnnuity(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkCents, _ = PresentValueGrowingAnnuity(154500000, 0.10, 0.03, 20)
	}
}

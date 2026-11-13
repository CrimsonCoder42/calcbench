package interest

import (
	"math/big"
	"testing"
)

func TestFutureValueGrowingAnnuity_SourceExample(t *testing.T) {
	// Damodaran App. 3 gold mine ($1,545,000 first payment, 10 %, 3 %, 20 years) carried to year 20:
	// $16,145,979.98 × 1.1²⁰ = $108,622,079.49 (hand-checked in exact rational arithmetic).
	got, err := FutureValueGrowingAnnuity(154500000, 0.10, 0.03, 20)
	wantCents(t, "gold mine at year 20", got, err, 10862207949)
}

func TestFutureValueGrowingAnnuity_Table(t *testing.T) {
	cases := []struct {
		name    string
		c       int64
		r, g    float64
		periods int
		want    int64
	}{
		{"zero periods", 10000, 0.05, 0.02, 0, 0},
		{"one period", 10000, 0.05, 0.02, 1, 10000},
		{"two periods: 11,000 + 10,300", 10000, 0.10, 0.03, 2, 21300},
		{"g > r: 10,404 + 10,710 + 11,025", 10000, 0.02, 0.05, 3, 32139},
		{"r = g gives C·n·(1+r)^(n−1)", 10000, 0.10, 0.10, 2, 22000},
		{"zero growth is the level annuity", 100000, 0.05, 0, 5, 552563},
		{"negative payment", -10000, 0.10, 0.03, 2, -21300},
	}
	for _, c := range cases {
		got, err := FutureValueGrowingAnnuity(c.c, c.r, c.g, c.periods)
		wantCents(t, c.name, got, err, c.want)
	}
}

func TestFutureValueGrowingAnnuity_Errors(t *testing.T) {
	_, err := FutureValueGrowingAnnuity(-maxExactCents-1, 0.05, 0.02, 5)
	wantErr(t, "-2^53-1", err, ErrOutOfRange)
	for _, r := range badRates {
		_, err = FutureValueGrowingAnnuity(100, r, 0.02, 5)
		wantErr(t, "bad rate", err, ErrRate)
		_, err = FutureValueGrowingAnnuity(100, 0.05, r, 5)
		wantErr(t, "bad growth", err, ErrRate)
	}
	_, err = FutureValueGrowingAnnuity(100, 0.05, 0.02, -1)
	wantErr(t, "negative periods", err, ErrPeriods)
	_, err = FutureValueGrowingAnnuity(maxExactCents, 0.05, 0.05, 2)
	wantErr(t, "overflow", err, ErrOverflow)
}

func TestFutureValueGrowingAnnuity_Oracle(t *testing.T) {
	rng := testRNG()
	for k := 0; k < 300; k++ {
		c := rng.Int64N(10_000_000) // |FV| ≤ 1e7·60·1.15⁶⁰ < 2⁴³
		r := rng.Float64() * 0.15
		g := rng.Float64()*0.12 - 0.02
		n := 1 + rng.IntN(60)
		got, err := FutureValueGrowingAnnuity(c, r, g, n)
		if err != nil {
			t.Fatal(err)
		}
		exact := new(big.Float).SetPrec(bigPrec).Mul(bigGrowingAnnuityPV(c, r, g, n), bigCompound(r, n))
		checkCentsOracle(t, "oracle", got, exact)
	}
}

func TestFutureValueGrowingAnnuity_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "FutureValueGrowingAnnuity", func() { sinkCents, _ = FutureValueGrowingAnnuity(154500000, 0.10, 0.03, 20) })
}

func BenchmarkFutureValueGrowingAnnuity(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkCents, _ = FutureValueGrowingAnnuity(154500000, 0.10, 0.03, 20)
	}
}

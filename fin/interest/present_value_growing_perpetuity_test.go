package interest

import (
	"math/big"
	"testing"
)

func TestPresentValueGrowingPerpetuity_SourceExample(t *testing.T) {
	// Principles of Finance (OpenStax 2022) §8.1, Eq. 8.3: Rooney Corp. $1.75/(0.058 − 0.012) = $38.04;
	// Savo Corp. $4.00/(0.03 − 0.002) = $142.86.
	got, err := PresentValueGrowingPerpetuity(175, 0.058, 0.012)
	wantCents(t, "§8.1 Rooney", got, err, 3804)
	got, err = PresentValueGrowingPerpetuity(400, 0.03, 0.002)
	wantCents(t, "§8.1 Savo", got, err, 14286)
}

func TestPresentValueGrowingPerpetuity_Table(t *testing.T) {
	got, err := PresentValueGrowingPerpetuity(175, 0.058, 0)
	wantCents(t, "zero growth is the level perpetuity (§8.1: $30.17)", got, err, 3017)
	got, err = PresentValueGrowingPerpetuity(1000, 0.05, -0.05)
	wantCents(t, "negative growth", got, err, 10000)
	got, err = PresentValueGrowingPerpetuity(-175, 0.058, 0.012)
	wantCents(t, "negative payment", got, err, -3804)
}

func TestPresentValueGrowingPerpetuity_Errors(t *testing.T) {
	_, err := PresentValueGrowingPerpetuity(maxExactCents+1, 0.05, 0.01)
	wantErr(t, "2^53+1", err, ErrOutOfRange)
	for _, r := range badRates {
		_, err = PresentValueGrowingPerpetuity(100, r, 0.01)
		wantErr(t, "bad rate", err, ErrRate)
		_, err = PresentValueGrowingPerpetuity(100, 0.05, r)
		wantErr(t, "bad growth", err, ErrRate)
	}
	_, err = PresentValueGrowingPerpetuity(100, 0.05, 0.05)
	wantErr(t, "r = g", err, ErrDomain)
	_, err = PresentValueGrowingPerpetuity(100, 0.05, 0.06)
	wantErr(t, "r < g", err, ErrDomain)
	_, err = PresentValueGrowingPerpetuity(maxExactCents, 0.1, 0.05)
	wantErr(t, "overflow", err, ErrOverflow)
}

func TestPresentValueGrowingPerpetuity_Oracle(t *testing.T) {
	rng := testRNG()
	for k := 0; k < 300; k++ {
		c := rng.Int64N(100_000_000)
		g := rng.Float64()*0.1 - 0.02
		r := g + 0.01 + rng.Float64()*0.1
		got, err := PresentValueGrowingPerpetuity(c, r, g)
		if err != nil {
			t.Fatal(err)
		}
		den := new(big.Float).SetPrec(bigPrec).Sub(bigF(r), bigF(g))
		checkCentsOracle(t, "oracle", got, new(big.Float).SetPrec(bigPrec).Quo(bigI(c), den))
	}
}

func TestPresentValueGrowingPerpetuity_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "PresentValueGrowingPerpetuity", func() { sinkCents, _ = PresentValueGrowingPerpetuity(175, 0.058, 0.012) })
}

func BenchmarkPresentValueGrowingPerpetuity(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkCents, _ = PresentValueGrowingPerpetuity(175, 0.058, 0.012)
	}
}

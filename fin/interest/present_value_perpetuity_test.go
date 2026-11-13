package interest

import (
	"math/big"
	"testing"
)

func TestPresentValuePerpetuity_SourceExample(t *testing.T) {
	// Principles of Finance (OpenStax 2022) §8.1 Eq. 8.2: Shaw Inc. preferred, $2.00/0.07 = $28.57.
	got, err := PresentValuePerpetuity(200, 0.07)
	wantCents(t, "§8.1 Shaw", got, err, 2857)
	// §8.1 endowment: $75,000 a year at 3.5 % needs $2,142,857 (to the dollar; $2,142,857.14).
	got, err = PresentValuePerpetuity(7500000, 0.035)
	wantCents(t, "§8.1 endowment", got, err, 214285714)
	// DeFusco et al., Quantitative Investment Analysis 3rd ed. ch. 1 Example 14: £100 at 5 % → £2,000.
	got, err = PresentValuePerpetuity(10000, 0.05)
	wantCents(t, "QIA Example 14", got, err, 200000)
}

func TestPresentValuePerpetuity_Table(t *testing.T) {
	got, err := PresentValuePerpetuity(-1000, 0.2)
	wantCents(t, "negative payment", got, err, -5000)
	got, err = PresentValuePerpetuity(0, 0.2)
	wantCents(t, "zero payment", got, err, 0)
	got, err = PresentValuePerpetuity(1000, 0.2)
	wantCents(t, "QIA §6.2: $10 at 20 % → $50", got, err, 5000)
}

func TestPresentValuePerpetuity_Errors(t *testing.T) {
	_, err := PresentValuePerpetuity(maxExactCents+1, 0.05)
	wantErr(t, "2^53+1", err, ErrOutOfRange)
	for _, r := range badRates {
		_, err = PresentValuePerpetuity(100, r)
		wantErr(t, "bad rate", err, ErrRate)
	}
	_, err = PresentValuePerpetuity(100, 0)
	wantErr(t, "zero rate", err, ErrDomain)
	_, err = PresentValuePerpetuity(100, -0.01)
	wantErr(t, "negative rate", err, ErrDomain)
	_, err = PresentValuePerpetuity(maxExactCents, 0.5)
	wantErr(t, "overflow", err, ErrOverflow)
}

func TestPresentValuePerpetuity_Oracle(t *testing.T) {
	// The perpetuity is the limit of the ordinary annuity as n → ∞ (QIA ch. 1 §6.2).
	rng := testRNG()
	for k := 0; k < 300; k++ {
		c := rng.Int64N(1_000_000_000)
		r := 0.001 + rng.Float64()*0.2
		got, err := PresentValuePerpetuity(c, r)
		if err != nil {
			t.Fatal(err)
		}
		checkCentsOracle(t, "oracle", got, new(big.Float).SetPrec(bigPrec).Quo(bigI(c), bigF(r)))
	}
}

func TestPresentValuePerpetuity_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "PresentValuePerpetuity", func() { sinkCents, _ = PresentValuePerpetuity(200, 0.07) })
}

func BenchmarkPresentValuePerpetuity(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkCents, _ = PresentValuePerpetuity(200, 0.07)
	}
}

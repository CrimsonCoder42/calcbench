package interest

import (
	"math"
	"math/big"
	"testing"
)

func TestEffectiveAnnualRate_SourceExample(t *testing.T) {
	// Principles of Finance (OpenStax 2022) §8.4, Eq. 8.38: 1.5 % a month (18 % stated) →
	// (1.015)^12 − 1 = 0.19562.
	got, err := EffectiveAnnualRate(0.18, 12)
	wantFloat(t, "§8.4", got, err, 0.19562, 5e-6)
	// DeFusco et al., Quantitative Investment Analysis 3rd ed. ch. 1 §3.3, Eq. 5: 8 % semiannual → 8.16 %.
	got, err = EffectiveAnnualRate(0.08, 2)
	wantFloat(t, "QIA §3.3", got, err, 0.0816, 1e-15)
	// NCEES FE Reference Handbook 10.0.1 p. 231: i_e = (1 + r/m)^m − 1; QIA Table 1 daily: 8 % → 1.083278.
	got, err = EffectiveAnnualRate(0.08, 365)
	wantFloat(t, "QIA Table 1 daily", got, err, 0.083278, 5e-7)
}

func TestEffectiveAnnualRate_Table(t *testing.T) {
	got, err := EffectiveAnnualRate(0, 12)
	wantFloat(t, "zero rate", got, err, 0, 0)
	got, err = EffectiveAnnualRate(0.07, 1)
	wantRel(t, "annual compounding is the rate", got, err, 0.07, 1e-15)
	got, err = EffectiveAnnualRate(-0.12, 12)
	wantRel(t, "negative rate", got, err, math.Pow(0.99, 12)-1, 1e-14)
	got, err = EffectiveAnnualRate(1e-10, 365)
	wantRel(t, "tiny rate keeps precision", got, err, 1e-10+1e-20*364/(2*365), 1e-15)
}

func TestEffectiveAnnualRate_Errors(t *testing.T) {
	_, err := EffectiveAnnualRate(math.NaN(), 12)
	wantErr(t, "NaN", err, ErrNotFinite)
	_, err = EffectiveAnnualRate(math.Inf(1), 12)
	wantErr(t, "+Inf", err, ErrNotFinite)
	_, err = EffectiveAnnualRate(0.05, 0)
	wantErr(t, "zero periods per year", err, ErrPeriods)
	_, err = EffectiveAnnualRate(-12, 12)
	wantErr(t, "periodic rate −100 %", err, ErrRate)
	_, err = EffectiveAnnualRate(1e6, 1000)
	wantErr(t, "beyond float64", err, ErrOverflow)
}

func TestEffectiveAnnualRate_Oracle(t *testing.T) {
	rng := testRNG()
	for k := 0; k < 300; k++ {
		r := rng.Float64()*0.4 - 0.05
		m := 1 + rng.IntN(365)
		got, err := EffectiveAnnualRate(r, m)
		exact := new(big.Float).SetPrec(bigPrec).Quo(bigF(r), bigI(int64(m)))
		exact.Add(exact, bigI(1))
		v, _ := new(big.Float).Sub(bigPow(exact, m), bigI(1)).Float64()
		wantRel(t, "oracle", got, err, v, 1e-12)
	}
}

func TestEffectiveAnnualRate_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "EffectiveAnnualRate", func() { sinkFloat, _ = EffectiveAnnualRate(0.18, 12) })
}

func BenchmarkEffectiveAnnualRate(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = EffectiveAnnualRate(0.18, 12)
	}
}

package interest

import (
	"math"
	"testing"
)

func TestPeriodicRateFromEffectiveAnnual_SourceExample(t *testing.T) {
	// DeFusco et al., Quantitative Investment Analysis 3rd ed. ch. 1 §3.3, Eq. 5: 4 % a half-year gives
	// EAR 8.16 %; read backwards, 8.16 % effective with m = 2 is 4 % per half-year.
	got, err := PeriodicRateFromEffectiveAnnual(0.0816, 2)
	wantFloat(t, "QIA §3.3 backwards", got, err, 0.04, 1e-15)
	// Principles of Finance (OpenStax 2022) §8.4, Eq. 8.38: 1.5 % a month → (1.015)^12 − 1 = 0.195618...
	got, err = PeriodicRateFromEffectiveAnnual(math.Pow(1.015, 12)-1, 12)
	wantFloat(t, "§8.4 backwards", got, err, 0.015, 1e-15)
}

func TestPeriodicRateFromEffectiveAnnual_Table(t *testing.T) {
	got, err := PeriodicRateFromEffectiveAnnual(0, 12)
	wantFloat(t, "zero rate", got, err, 0, 0)
	got, err = PeriodicRateFromEffectiveAnnual(0.07, 1)
	wantRel(t, "one period a year is the rate", got, err, 0.07, 1e-15)
	got, err = PeriodicRateFromEffectiveAnnual(0.06, 12)
	wantRel(t, "6 % effective, monthly", got, err, 0.004867550565343048, 1e-14)
	got, err = PeriodicRateFromEffectiveAnnual(-0.5, 2)
	wantRel(t, "negative rate", got, err, math.Sqrt(0.5)-1, 1e-15)
	got, err = PeriodicRateFromEffectiveAnnual(1e-12, 365)
	wantRel(t, "tiny rate keeps precision", got, err, 1e-12/365, 1e-9)
}

func TestPeriodicRateFromEffectiveAnnual_Errors(t *testing.T) {
	for _, r := range badRates {
		_, err := PeriodicRateFromEffectiveAnnual(r, 12)
		wantErr(t, "bad rate", err, ErrRate)
	}
	_, err := PeriodicRateFromEffectiveAnnual(0.05, 0)
	wantErr(t, "zero periods", err, ErrPeriods)
	_, err = PeriodicRateFromEffectiveAnnual(0.05, -12)
	wantErr(t, "negative periods", err, ErrPeriods)
}

func TestPeriodicRateFromEffectiveAnnual_Oracle(t *testing.T) {
	// Round trip: EffectiveAnnualRate(m × periodic, m) returns the input.
	rng := testRNG()
	for k := 0; k < 300; k++ {
		ie := rng.Float64()*0.5 - 0.2
		m := 1 + rng.IntN(365)
		p, err := PeriodicRateFromEffectiveAnnual(ie, m)
		if err != nil {
			t.Fatalf("ie=%v m=%d: %v", ie, m, err)
		}
		back, err := EffectiveAnnualRate(p*float64(m), m)
		wantFloat(t, "round trip", back, err, ie, 1e-12)
	}
}

func TestPeriodicRateFromEffectiveAnnual_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "PeriodicRateFromEffectiveAnnual", func() { _, _ = PeriodicRateFromEffectiveAnnual(0.06, 12) })
}

func BenchmarkPeriodicRateFromEffectiveAnnual(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = PeriodicRateFromEffectiveAnnual(0.06, 12)
	}
}

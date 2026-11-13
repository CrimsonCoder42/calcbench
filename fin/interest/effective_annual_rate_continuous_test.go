package interest

import (
	"math"
	"testing"
)

func TestEffectiveAnnualRateContinuous_SourceExample(t *testing.T) {
	// DeFusco et al., Quantitative Investment Analysis 3rd ed. ch. 1 §3.3, Eq. 6: 8 % continuous → 8.33 %;
	// Table 1: $1.00e^0.08 = $1.083287.
	got, err := EffectiveAnnualRateContinuous(0.08)
	wantFloat(t, "QIA Eq. 6", got, err, 0.083287, 5e-7)
}

func TestEffectiveAnnualRateContinuous_Table(t *testing.T) {
	got, err := EffectiveAnnualRateContinuous(0)
	wantFloat(t, "zero", got, err, 0, 0)
	got, err = EffectiveAnnualRateContinuous(1e-12)
	wantRel(t, "tiny rate keeps precision", got, err, 1e-12+5e-25, 1e-15)
	got, err = EffectiveAnnualRateContinuous(-0.5)
	wantRel(t, "negative", got, err, math.Exp(-0.5)-1, 1e-15)
	// Above every finite compounding frequency: e^r − 1 > (1 + r/m)^m − 1.
	disc, err2 := EffectiveAnnualRate(0.08, 365)
	if err2 != nil {
		t.Fatal(err2)
	}
	if got, _ = EffectiveAnnualRateContinuous(0.08); got <= disc {
		t.Errorf("continuous %v not above daily %v", got, disc)
	}
}

func TestEffectiveAnnualRateContinuous_Errors(t *testing.T) {
	for _, r := range badRates {
		_, err := EffectiveAnnualRateContinuous(r)
		wantErr(t, "bad rate", err, ErrRate)
	}
	_, err := EffectiveAnnualRateContinuous(800)
	wantErr(t, "beyond float64", err, ErrOverflow)
}

func TestEffectiveAnnualRateContinuous_Oracle(t *testing.T) {
	// Limit of EffectiveAnnualRate as m grows: at m = 1e6 they agree to about r²/(2m).
	for _, r := range []float64{0.01, 0.05, 0.08, 0.2} {
		c, err := EffectiveAnnualRateContinuous(r)
		d, err2 := EffectiveAnnualRate(r, 1_000_000)
		if err2 != nil {
			t.Fatal(err2)
		}
		wantFloat(t, "limit", c, err, d, r*r/1e6)
	}
}

func TestEffectiveAnnualRateContinuous_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "EffectiveAnnualRateContinuous", func() { sinkFloat, _ = EffectiveAnnualRateContinuous(0.08) })
}

func BenchmarkEffectiveAnnualRateContinuous(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = EffectiveAnnualRateContinuous(0.08)
	}
}

package interest

import (
	"math"
	"testing"
)

func TestUniformSeriesSinkingFund_SourceExample(t *testing.T) {
	// NCEES FE Reference Handbook 10.0.1, factor table p. 233, i = 1.00 %, n = 10: A/F = 0.0956.
	got, err := UniformSeriesSinkingFund(0.01, 10)
	wantFloat(t, "FE table", got, err, 0.0956, 5e-5)
}

func TestUniformSeriesSinkingFund_Table(t *testing.T) {
	got, err := UniformSeriesSinkingFund(0, 8)
	wantFloat(t, "zero rate is 1/n", got, err, 0.125, 0)
	got, err = UniformSeriesSinkingFund(0.07, 1)
	wantRel(t, "one period", got, err, 1, 1e-15)
	got, err = UniformSeriesSinkingFund(1, 1100)
	wantFloat(t, "underflow to 0 is not an error", got, err, 0, 0)
}

func TestUniformSeriesSinkingFund_Errors(t *testing.T) {
	for _, r := range badRates {
		_, err := UniformSeriesSinkingFund(r, 5)
		wantErr(t, "bad rate", err, ErrRate)
	}
	_, err := UniformSeriesSinkingFund(0.05, 0)
	wantErr(t, "periods below minimum", err, ErrPeriods)

}

func TestUniformSeriesSinkingFund_Oracle(t *testing.T) {
	for _, c := range factorCases {
		got, err := UniformSeriesSinkingFund(c.i, c.n)
		wantRel(t, "A/F", got, err, bigFactors(c.i, c.n).af, factorRelTol)
	}
}

func TestUniformSeriesSinkingFund_Identity(t *testing.T) {
	// A/P = A/F + i.
	for _, c := range factorCases {
		af, err1 := UniformSeriesSinkingFund(c.i, c.n)
		ap, err2 := CapitalRecovery(c.i, c.n)
		if err1 != nil || err2 != nil {
			t.Fatal(err1, err2)
		}
		wantRel(t, "A/P = A/F + i", ap, nil, af+c.i, 1e-13)
	}
}

func TestUniformSeriesSinkingFund_NeverOverflows(t *testing.T) {
	// (F/A) ≥ 1, so (A/F) ≤ 1 at the extremes of the rate range.
	for _, r := range []float64{-1 + 1e-15, 1e300, math.MaxFloat64} {
		got, err := UniformSeriesSinkingFund(r, 7)
		if err != nil || got < 0 || got > 1+1e-15 {
			t.Errorf("rate %v: got %v, %v; want a value in [0, 1]", r, got, err)
		}
	}
}

func TestUniformSeriesSinkingFund_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "A/F", func() { sinkFloat, _ = UniformSeriesSinkingFund(0.01, 10) })
}

func BenchmarkUniformSeriesSinkingFund(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = UniformSeriesSinkingFund(0.01, 10)
	}
}

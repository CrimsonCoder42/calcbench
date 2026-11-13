package interest

import "testing"

func TestUniformGradientPresentWorth_SourceExample(t *testing.T) {
	// NCEES FE Reference Handbook 10.0.1, factor table p. 233, i = 1.00 %, n = 10: P/G = 41.8435.
	got, err := UniformGradientPresentWorth(0.01, 10)
	wantFloat(t, "FE table", got, err, 41.8435, 5e-5)
}

func TestUniformGradientPresentWorth_Table(t *testing.T) {
	got, err := UniformGradientPresentWorth(0, 10)
	wantFloat(t, "zero rate is n(n−1)/2", got, err, 45, 0)
	got, err = UniformGradientPresentWorth(0.07, 0)
	wantFloat(t, "zero periods", got, err, 0, 0)
	got, err = UniformGradientPresentWorth(0.07, 1)
	wantFloat(t, "one period", got, err, 0, 0)
	got, err = UniformGradientPresentWorth(0.1, 2)
	wantRel(t, "two periods: 1/1.1²", got, err, 1/1.21, 1e-14)
	got, err = UniformGradientPresentWorth(0.1, 100000)
	wantRel(t, "large n tends to 1/i²", got, err, 100, 1e-12)
}

func TestUniformGradientPresentWorth_Errors(t *testing.T) {
	for _, r := range badRates {
		_, err := UniformGradientPresentWorth(r, 5)
		wantErr(t, "bad rate", err, ErrRate)
	}
	_, err := UniformGradientPresentWorth(0.05, -1)
	wantErr(t, "periods below minimum", err, ErrPeriods)
	_, err = UniformGradientPresentWorth(-0.5, 1100)
	wantErr(t, "beyond float64", err, ErrOverflow)
}

func TestUniformGradientPresentWorth_Oracle(t *testing.T) {
	for _, c := range factorCases {
		got, err := UniformGradientPresentWorth(c.i, c.n)
		wantRel(t, "P/G", got, err, bigFactors(c.i, c.n).pg, factorRelTol)
	}
}

func TestUniformGradientPresentWorth_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "P/G", func() { sinkFloat, _ = UniformGradientPresentWorth(0.01, 10) })
}

func BenchmarkUniformGradientPresentWorth(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = UniformGradientPresentWorth(0.01, 10)
	}
}

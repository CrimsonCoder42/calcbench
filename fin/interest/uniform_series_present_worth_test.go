package interest

import "testing"

func TestUniformSeriesPresentWorth_SourceExample(t *testing.T) {
	// NCEES FE Reference Handbook 10.0.1, factor table p. 233, i = 1.00 %, n = 10: P/A = 9.4713.
	got, err := UniformSeriesPresentWorth(0.01, 10)
	wantFloat(t, "FE table", got, err, 9.4713, 5e-5)
}

func TestUniformSeriesPresentWorth_Table(t *testing.T) {
	got, err := UniformSeriesPresentWorth(0, 37)
	wantFloat(t, "zero rate is n exactly", got, err, 37, 0)
	got, err = UniformSeriesPresentWorth(0.07, 0)
	wantFloat(t, "zero periods", got, err, 0, 0)
	got, err = UniformSeriesPresentWorth(0.25, 1)
	wantRel(t, "one period", got, err, 0.8, 1e-15)
}

func TestUniformSeriesPresentWorth_Errors(t *testing.T) {
	for _, r := range badRates {
		_, err := UniformSeriesPresentWorth(r, 5)
		wantErr(t, "bad rate", err, ErrRate)
	}
	_, err := UniformSeriesPresentWorth(0.05, -1)
	wantErr(t, "periods below minimum", err, ErrPeriods)
	_, err = UniformSeriesPresentWorth(-0.5, 1100)
	wantErr(t, "beyond float64", err, ErrOverflow)
}

func TestUniformSeriesPresentWorth_Oracle(t *testing.T) {
	for _, c := range factorCases {
		got, err := UniformSeriesPresentWorth(c.i, c.n)
		wantRel(t, "P/A", got, err, bigFactors(c.i, c.n).pa, factorRelTol)
	}
}

func TestUniformSeriesPresentWorth_Identity(t *testing.T) {
	// P/A = (F/A)·(P/F).
	for _, c := range factorCases {
		pa, err1 := UniformSeriesPresentWorth(c.i, c.n)
		fa, err2 := UniformSeriesCompoundAmount(c.i, c.n)
		pf, err3 := SinglePaymentPresentWorth(c.i, c.n)
		if err1 != nil || err2 != nil || err3 != nil {
			t.Fatal(err1, err2, err3)
		}
		wantRel(t, "P/A = F/A·P/F", pa, nil, fa*pf, 1e-13)
	}
}

func TestUniformSeriesPresentWorth_SourceExampleQIA(t *testing.T) {
	// Quantitative Investment Analysis 3rd ed. ch. 1 Example 11: (1 − 1/1.12^5)/0.12 = 3.604776.
	got, err := UniformSeriesPresentWorth(0.12, 5)
	wantFloat(t, "QIA Example 11", got, err, 3.604776, 5e-7)
}

func TestUniformSeriesPresentWorth_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "P/A", func() { sinkFloat, _ = UniformSeriesPresentWorth(0.01, 10) })
}

func BenchmarkUniformSeriesPresentWorth(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = UniformSeriesPresentWorth(0.01, 10)
	}
}

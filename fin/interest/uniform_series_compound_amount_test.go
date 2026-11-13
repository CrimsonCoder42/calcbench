package interest

import "testing"

func TestUniformSeriesCompoundAmount_SourceExample(t *testing.T) {
	// NCEES FE Reference Handbook 10.0.1, factor table p. 233, i = 1.00 %, n = 10: F/A = 10.4622.
	got, err := UniformSeriesCompoundAmount(0.01, 10)
	wantFloat(t, "FE table", got, err, 10.4622, 5e-5)
}

func TestUniformSeriesCompoundAmount_Table(t *testing.T) {
	got, err := UniformSeriesCompoundAmount(0, 37)
	wantFloat(t, "zero rate is n exactly", got, err, 37, 0)
	got, err = UniformSeriesCompoundAmount(0.07, 0)
	wantFloat(t, "zero periods", got, err, 0, 0)
	got, err = UniformSeriesCompoundAmount(0.07, 1)
	wantRel(t, "one period", got, err, 1, 1e-15)
	got, err = UniformSeriesCompoundAmount(1e-15, 1000)
	wantRel(t, "tiny rate, no cancellation", got, err, 1000+1000*999/2*1e-15, 1e-15)
}

func TestUniformSeriesCompoundAmount_Errors(t *testing.T) {
	for _, r := range badRates {
		_, err := UniformSeriesCompoundAmount(r, 5)
		wantErr(t, "bad rate", err, ErrRate)
	}
	_, err := UniformSeriesCompoundAmount(0.05, -1)
	wantErr(t, "periods below minimum", err, ErrPeriods)
	_, err = UniformSeriesCompoundAmount(1, 1100)
	wantErr(t, "beyond float64", err, ErrOverflow)
}

func TestUniformSeriesCompoundAmount_Oracle(t *testing.T) {
	for _, c := range factorCases {
		got, err := UniformSeriesCompoundAmount(c.i, c.n)
		wantRel(t, "F/A", got, err, bigFactors(c.i, c.n).fa, factorRelTol)
	}
}

func TestUniformSeriesCompoundAmount_Identity(t *testing.T) {
	// F/A = (F/P − 1)/i (NCEES footnote form), checked where i is not small.
	for _, i := range []float64{0.02, 0.1, 0.3} {
		fa, err1 := UniformSeriesCompoundAmount(i, 25)
		fp, err2 := SinglePaymentCompoundAmount(i, 25)
		if err1 != nil || err2 != nil {
			t.Fatal(err1, err2)
		}
		wantRel(t, "F/A identity", fa, nil, (fp-1)/i, 1e-13)
	}
}

func TestUniformSeriesCompoundAmount_SourceExampleQIA(t *testing.T) {
	// DeFusco et al., Quantitative Investment Analysis 3rd ed. ch. 1 Example 7: (1.09^30 − 1)/0.09 = 136.307539.
	got, err := UniformSeriesCompoundAmount(0.09, 30)
	wantFloat(t, "QIA Example 7", got, err, 136.307539, 5e-7)
}

func TestUniformSeriesCompoundAmount_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "F/A", func() { sinkFloat, _ = UniformSeriesCompoundAmount(0.01, 10) })
}

func BenchmarkUniformSeriesCompoundAmount(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = UniformSeriesCompoundAmount(0.01, 10)
	}
}

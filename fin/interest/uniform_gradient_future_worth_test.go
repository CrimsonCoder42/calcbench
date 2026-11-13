package interest

import "testing"

func TestUniformGradientFutureWorth_SourceExample(t *testing.T) {
	// NCEES FE Reference Handbook 10.0.1, factor table p. 233, i = 1.00 %, n = 10: F/G = 46.2213.
	// The table has no F/G column; 46.2213 = (F/A − n)/i = (10.4622 − 10)/0.01 by the handbook's footnote.
	got, err := UniformGradientFutureWorth(0.01, 10)
	wantFloat(t, "FE table", got, err, 46.2213, 5e-5)
}

func TestUniformGradientFutureWorth_Table(t *testing.T) {
	got, err := UniformGradientFutureWorth(0, 10)
	wantFloat(t, "zero rate is n(n−1)/2", got, err, 45, 0)
	got, err = UniformGradientFutureWorth(0.07, 1)
	wantFloat(t, "one period", got, err, 0, 0)
	got, err = UniformGradientFutureWorth(0.1, 2)
	wantRel(t, "two periods", got, err, 1, 1e-14)
	got, err = UniformGradientFutureWorth(0.1, 3)
	wantRel(t, "three periods: 1.1 + 2", got, err, 3.1, 1e-14)
}

func TestUniformGradientFutureWorth_Errors(t *testing.T) {
	for _, r := range badRates {
		_, err := UniformGradientFutureWorth(r, 5)
		wantErr(t, "bad rate", err, ErrRate)
	}
	_, err := UniformGradientFutureWorth(0.05, -1)
	wantErr(t, "periods below minimum", err, ErrPeriods)
	_, err = UniformGradientFutureWorth(1, 1100)
	wantErr(t, "beyond float64", err, ErrOverflow)
}

func TestUniformGradientFutureWorth_Oracle(t *testing.T) {
	for _, c := range factorCases {
		got, err := UniformGradientFutureWorth(c.i, c.n)
		wantRel(t, "F/G", got, err, bigFactors(c.i, c.n).fg, factorRelTol)
	}
}

func TestUniformGradientFutureWorth_Identity(t *testing.T) {
	// NCEES footnote: F/G = (F/A)·(A/G).
	for _, c := range factorCases {
		fg, err1 := UniformGradientFutureWorth(c.i, c.n)
		fa, err2 := UniformSeriesCompoundAmount(c.i, c.n)
		ag, err3 := UniformGradientUniformSeries(c.i, c.n)
		if err1 != nil || err2 != nil || err3 != nil {
			t.Fatal(err1, err2, err3)
		}
		wantRel(t, "F/G = F/A·A/G", fg, nil, fa*ag, 1e-12)
	}
}

func TestUniformGradientFutureWorth_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "F/G", func() { sinkFloat, _ = UniformGradientFutureWorth(0.01, 10) })
}

func BenchmarkUniformGradientFutureWorth(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = UniformGradientFutureWorth(0.01, 10)
	}
}

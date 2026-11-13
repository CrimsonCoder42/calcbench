package interest

import "testing"

func TestSinglePaymentPresentWorth_SourceExample(t *testing.T) {
	// NCEES FE Reference Handbook 10.0.1, factor table p. 233, i = 1.00 %, n = 10: P/F = 0.9053.
	got, err := SinglePaymentPresentWorth(0.01, 10)
	wantFloat(t, "FE table", got, err, 0.9053, 5e-5)
}

func TestSinglePaymentPresentWorth_Table(t *testing.T) {
	got, err := SinglePaymentPresentWorth(0, 40)
	wantFloat(t, "zero rate", got, err, 1, 0)
	got, err = SinglePaymentPresentWorth(0.07, 0)
	wantFloat(t, "zero periods", got, err, 1, 0)
	got, err = SinglePaymentPresentWorth(0.25, 1)
	wantRel(t, "one period", got, err, 0.8, 1e-15)
	got, err = SinglePaymentPresentWorth(1, 2000)
	wantFloat(t, "underflow to 0 is not an error", got, err, 0, 0)
}

func TestSinglePaymentPresentWorth_Errors(t *testing.T) {
	for _, r := range badRates {
		_, err := SinglePaymentPresentWorth(r, 5)
		wantErr(t, "bad rate", err, ErrRate)
	}
	_, err := SinglePaymentPresentWorth(0.05, -1)
	wantErr(t, "periods below minimum", err, ErrPeriods)
	_, err = SinglePaymentPresentWorth(-0.5, 1100)
	wantErr(t, "beyond float64", err, ErrOverflow)
}

func TestSinglePaymentPresentWorth_Oracle(t *testing.T) {
	for _, c := range factorCases {
		got, err := SinglePaymentPresentWorth(c.i, c.n)
		wantRel(t, "P/F", got, err, bigFactors(c.i, c.n).pf, factorRelTol)
	}
}

func TestSinglePaymentPresentWorth_Identity(t *testing.T) {
	// (P/F) · (F/P) = 1.
	for _, c := range factorCases {
		pf, err1 := SinglePaymentPresentWorth(c.i, c.n)
		fp, err2 := SinglePaymentCompoundAmount(c.i, c.n)
		if err1 != nil || err2 != nil {
			t.Fatal(err1, err2)
		}
		wantRel(t, "P/F·F/P", pf*fp, nil, 1, 1e-14)
	}
}

func TestSinglePaymentPresentWorth_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "P/F", func() { sinkFloat, _ = SinglePaymentPresentWorth(0.01, 10) })
}

func BenchmarkSinglePaymentPresentWorth(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = SinglePaymentPresentWorth(0.01, 10)
	}
}

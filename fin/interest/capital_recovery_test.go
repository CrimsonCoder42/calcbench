package interest

import (
	"math"
	"testing"
)

func TestCapitalRecovery_SourceExample(t *testing.T) {
	// NCEES FE Reference Handbook 10.0.1, factor table p. 233, i = 1.00 %, n = 10: A/P = 0.1056.
	got, err := CapitalRecovery(0.01, 10)
	wantFloat(t, "FE table", got, err, 0.1056, 5e-5)
}

func TestCapitalRecovery_Table(t *testing.T) {
	got, err := CapitalRecovery(0, 4)
	wantFloat(t, "zero rate is 1/n", got, err, 0.25, 0)
	got, err = CapitalRecovery(0.25, 1)
	wantRel(t, "one period", got, err, 1.25, 1e-15)
	got, err = CapitalRecovery(-0.5, 1100)
	wantFloat(t, "underflow to 0 is not an error", got, err, 0, 0)
}

func TestCapitalRecovery_Errors(t *testing.T) {
	for _, r := range badRates {
		_, err := CapitalRecovery(r, 5)
		wantErr(t, "bad rate", err, ErrRate)
	}
	_, err := CapitalRecovery(0.05, 0)
	wantErr(t, "periods below minimum", err, ErrPeriods)

}

func TestCapitalRecovery_Oracle(t *testing.T) {
	for _, c := range factorCases {
		got, err := CapitalRecovery(c.i, c.n)
		wantRel(t, "A/P", got, err, bigFactors(c.i, c.n).ap, factorRelTol)
	}
}

func TestCapitalRecovery_Identity(t *testing.T) {
	// (A/P)·(P/A) = 1.
	for _, c := range factorCases {
		ap, err1 := CapitalRecovery(c.i, c.n)
		pa, err2 := UniformSeriesPresentWorth(c.i, c.n)
		if err1 != nil || err2 != nil {
			t.Fatal(err1, err2)
		}
		wantRel(t, "A/P·P/A", ap*pa, nil, 1, 1e-14)
	}
}

func TestCapitalRecovery_NeverOverflows(t *testing.T) {
	// (A/P) ≤ 1 + i: computed as i + (A/F) for i ≥ 0 it stays finite even at the largest rate.
	got, err := CapitalRecovery(math.MaxFloat64, 1)
	wantRel(t, "MaxFloat64, one period", got, err, math.MaxFloat64, 1e-15)
	got, err = CapitalRecovery(1e300, 3)
	wantRel(t, "huge rate tends to i", got, err, 1e300, 1e-15)
}

func TestCapitalRecovery_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "A/P", func() { sinkFloat, _ = CapitalRecovery(0.01, 10) })
}

func BenchmarkCapitalRecovery(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = CapitalRecovery(0.01, 10)
	}
}

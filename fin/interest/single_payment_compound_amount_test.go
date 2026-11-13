package interest

import (
	"math/big"
	"testing"
)

// bigFactors returns the nine engineering-economy factors for i ≠ 0, n ≥ 1, from their handbook
// definitions evaluated at 256 bits (NCEES FE Reference Handbook 10.0.1, p. 230).
type factorSet struct{ fp, pf, fa, af, pa, ap, pg, fg, ag float64 }

func bigFactors(i float64, n int) factorSet {
	one := bigI(1)
	bi := bigF(i)
	bn := bigI(int64(n))
	q := func() *big.Float { return new(big.Float).SetPrec(bigPrec) }
	fp := bigCompound(i, n)
	pf := q().Quo(one, fp)
	fa := q().Quo(q().Sub(fp, one), bi)
	pa := q().Quo(q().Sub(one, pf), bi)
	fg := q().Quo(q().Sub(fa, bn), bi)
	pg := q().Mul(fg, pf)
	ag := q().Quo(fg, fa)
	f := func(x *big.Float) float64 { v, _ := x.Float64(); return v }
	return factorSet{
		fp: f(fp), pf: f(pf), fa: f(fa), af: 1 / f(fa), pa: f(pa), ap: 1 / f(pa),
		pg: f(pg), fg: f(fg), ag: f(ag),
	}
}

// factorCases are the (i, n) pairs every factor's oracle test runs over, including rates small
// enough that pow(1+i, n) − 1 would lose most of its digits.
var factorCases = []struct {
	i float64
	n int
}{
	{0.01, 10}, {0.05, 1}, {0.05, 2}, {0.08, 30}, {0.005, 360}, {0.25, 50}, {-0.03, 20},
	{1e-6, 12}, {1e-9, 360}, {1e-12, 100}, {0.12, 400},
}

// factorRelTol is the relative tolerance against the 256-bit oracle: the closed forms use a handful
// of correctly rounded operations (exp, log1p, expm1, ×, /) whose errors grow with |n·ln(1+i)|.
const factorRelTol = 1e-12

func TestSinglePaymentCompoundAmount_SourceExample(t *testing.T) {
	// NCEES FE Reference Handbook 10.0.1, factor table p. 233, i = 1.00 %, n = 10: F/P = 1.1046.
	got, err := SinglePaymentCompoundAmount(0.01, 10)
	wantFloat(t, "FE table", got, err, 1.1046, 5e-5)
}

func TestSinglePaymentCompoundAmount_Table(t *testing.T) {
	got, err := SinglePaymentCompoundAmount(0, 50)
	wantFloat(t, "zero rate", got, err, 1, 0)
	got, err = SinglePaymentCompoundAmount(0.07, 0)
	wantFloat(t, "zero periods", got, err, 1, 0)
	got, err = SinglePaymentCompoundAmount(0.07, 1)
	wantRel(t, "one period", got, err, 1.07, 1e-15)
	got, err = SinglePaymentCompoundAmount(-0.5, 3)
	wantRel(t, "negative rate", got, err, 0.125, 1e-15)
}

func TestSinglePaymentCompoundAmount_Errors(t *testing.T) {
	for _, r := range badRates {
		_, err := SinglePaymentCompoundAmount(r, 1)
		wantErr(t, "bad rate", err, ErrRate)
	}
	_, err := SinglePaymentCompoundAmount(0.05, -1)
	wantErr(t, "negative periods", err, ErrPeriods)
	_, err = SinglePaymentCompoundAmount(1, 1100)
	wantErr(t, "beyond float64", err, ErrOverflow)
}

func TestSinglePaymentCompoundAmount_Oracle(t *testing.T) {
	for _, c := range factorCases {
		got, err := SinglePaymentCompoundAmount(c.i, c.n)
		wantRel(t, "F/P", got, err, bigFactors(c.i, c.n).fp, factorRelTol)
	}
}

func TestSinglePaymentCompoundAmount_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "F/P", func() { sinkFloat, _ = SinglePaymentCompoundAmount(0.01, 10) })
}

func BenchmarkSinglePaymentCompoundAmount(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = SinglePaymentCompoundAmount(0.01, 10)
	}
}

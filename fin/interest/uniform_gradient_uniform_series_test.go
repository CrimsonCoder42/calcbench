package interest

import (
	"math"
	"testing"
)

func TestUniformGradientUniformSeries_SourceExample(t *testing.T) {
	// NCEES FE Reference Handbook 10.0.1, factor table p. 233, i = 1.00 %, n = 10: A/G = 4.4179.
	got, err := UniformGradientUniformSeries(0.01, 10)
	wantFloat(t, "FE table", got, err, 4.4179, 5e-5)
}

func TestUniformGradientUniformSeries_Table(t *testing.T) {
	got, err := UniformGradientUniformSeries(0, 10)
	wantFloat(t, "zero rate is (n−1)/2", got, err, 4.5, 0)
	got, err = UniformGradientUniformSeries(0.07, 1)
	wantFloat(t, "one period", got, err, 0, 0)
	got, err = UniformGradientUniformSeries(0.1, 2)
	wantRel(t, "two periods: 1/2.1", got, err, 1/2.1, 1e-14)
	got, err = UniformGradientUniformSeries(0.1, 100000)
	wantRel(t, "large n tends to 1/i", got, err, 10, 1e-12)
}

func TestUniformGradientUniformSeries_Errors(t *testing.T) {
	for _, r := range badRates {
		_, err := UniformGradientUniformSeries(r, 5)
		wantErr(t, "bad rate", err, ErrRate)
	}
	_, err := UniformGradientUniformSeries(0.05, 0)
	wantErr(t, "periods below minimum", err, ErrPeriods)

}

func TestUniformGradientUniformSeries_Oracle(t *testing.T) {
	for _, c := range factorCases {
		got, err := UniformGradientUniformSeries(c.i, c.n)
		wantRel(t, "A/G", got, err, bigFactors(c.i, c.n).ag, factorRelTol)
	}
}

func TestUniformGradientUniformSeries_Identity(t *testing.T) {
	// A/G = (P/G)·(A/P).
	for _, c := range factorCases {
		ag, err1 := UniformGradientUniformSeries(c.i, c.n)
		pg, err2 := UniformGradientPresentWorth(c.i, c.n)
		ap, err3 := CapitalRecovery(c.i, c.n)
		if err1 != nil || err2 != nil || err3 != nil {
			t.Fatal(err1, err2, err3)
		}
		wantRel(t, "A/G = P/G·A/P", ag, nil, pg*ap, 1e-12)
	}
}

func TestUniformGradientUniformSeries_NeverOverflows(t *testing.T) {
	// (A/G) is a weighted mean of 0 … n−1, so it stays in [0, n−1] at the extremes.
	for _, r := range []float64{-1 + 1e-15, -0.5, 1e300} {
		for _, n := range []int{2, 1100, 1 << 40} {
			got, err := UniformGradientUniformSeries(r, n)
			if err != nil || math.IsNaN(got) || got < -1e-9*float64(n) || got > float64(n-1)*(1+1e-9) {
				t.Errorf("rate %v n %d: got %v, %v; want a value in [0, n−1]", r, n, got, err)
			}
		}
	}
}

func TestUniformGradientUniformSeries_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "A/G", func() { sinkFloat, _ = UniformGradientUniformSeries(0.01, 10) })
}

func BenchmarkUniformGradientUniformSeries(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = UniformGradientUniformSeries(0.01, 10)
	}
}

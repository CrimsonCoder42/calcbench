package interest

import (
	"math"
	"testing"
)

func TestDoublingTime_SourceExample(t *testing.T) {
	// DeFusco et al., Quantitative Investment Analysis 3rd ed. ch. 1 Example 19: €10,000,000 at 7 %
	// doubles in ln(2)/ln(1.07) = 10.24 years.
	got, err := DoublingTime(0.07)
	wantFloat(t, "QIA Example 19", got, err, 10.24, 0.005)
}

func TestDoublingTime_Table(t *testing.T) {
	got, err := DoublingTime(1)
	wantFloat(t, "100 % doubles in one period", got, err, 1, 1e-15)
	got, err = DoublingTime(math.Sqrt2 - 1)
	wantFloat(t, "√2 − 1 doubles in two", got, err, 2, 1e-14)
	got, err = DoublingTime(1e-12)
	wantRel(t, "tiny rate", got, err, math.Ln2/1e-12, 1e-11)
}

func TestDoublingTime_Errors(t *testing.T) {
	for _, r := range badRates {
		_, err := DoublingTime(r)
		wantErr(t, "bad rate", err, ErrRate)
	}
	_, err := DoublingTime(0)
	wantErr(t, "zero", err, ErrDomain)
	_, err = DoublingTime(-0.5)
	wantErr(t, "negative", err, ErrDomain)
	_, err = DoublingTime(1e-320)
	wantErr(t, "beyond float64", err, ErrOverflow)
}

func TestDoublingTime_Oracle(t *testing.T) {
	// Definition: (1+i)^N = 2.
	rng := testRNG()
	for k := 0; k < 300; k++ {
		r := 1e-6 + rng.Float64()*0.5
		n, err := DoublingTime(r)
		if err != nil {
			t.Fatal(err)
		}
		wantRel(t, "(1+i)^N", math.Exp(n*math.Log1p(r)), nil, 2, 1e-13)
	}
}

func TestDoublingTime_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "DoublingTime", func() { sinkFloat, _ = DoublingTime(0.07) })
}

func BenchmarkDoublingTime(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = DoublingTime(0.07)
	}
}

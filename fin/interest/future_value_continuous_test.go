package interest

import (
	"math"
	"testing"
)

func TestFutureValueContinuous_SourceExample(t *testing.T) {
	// DeFusco et al., Quantitative Investment Analysis 3rd ed. (2015) ch. 1 Example 6:
	// $10,000 at 8 % compounded continuously for 2 years → $11,735.11.
	got, err := FutureValueContinuous(1000000, 0.08, 2)
	wantCents(t, "QIA Example 6", got, err, 1173511)
	// QIA ch. 1 Table 1: $1 at 8 % continuous for 1 year → $1.083287; $1,000,000.00 → $1,083,287.07.
	got, err = FutureValueContinuous(100000000, 0.08, 1)
	wantCents(t, "QIA Table 1 scaled", got, err, 108328707)
}

func TestFutureValueContinuous_Table(t *testing.T) {
	cases := []struct {
		name  string
		pv    int64
		rate  float64
		years float64
		want  int64
	}{
		{"zero rate", 5000, 0, 10, 5000},
		{"zero years", 5000, 0.1, 0, 5000},
		{"fractional year", 100000, 0.05, 0.5, 102532}, // 1000·e^0.025 = 1025.3151
		{"negative amount", -1000000, 0.08, 2, -1173511},
		{"zero amount", 0, 0.2, 5, 0},
		{"max exact, zero years", maxExactCents, 0.2, 0, maxExactCents},
	}
	for _, c := range cases {
		got, err := FutureValueContinuous(c.pv, c.rate, c.years)
		wantCents(t, c.name, got, err, c.want)
	}
}

func TestFutureValueContinuous_Errors(t *testing.T) {
	_, err := FutureValueContinuous(maxExactCents+1, 0.05, 1)
	wantErr(t, "2^53+1", err, ErrOutOfRange)
	for _, r := range badRates {
		_, err = FutureValueContinuous(100, r, 1)
		wantErr(t, "bad rate", err, ErrRate)
	}
	_, err = FutureValueContinuous(100, 0.05, math.NaN())
	wantErr(t, "NaN years", err, ErrNotFinite)
	_, err = FutureValueContinuous(100, 0.05, math.Inf(1))
	wantErr(t, "Inf years", err, ErrNotFinite)
	_, err = FutureValueContinuous(100, 0.05, -1)
	wantErr(t, "negative years", err, ErrNegative)
	_, err = FutureValueContinuous(maxExactCents, 0.05, 1)
	wantErr(t, "overflow", err, ErrOverflow)
}

func TestFutureValueContinuous_Oracle(t *testing.T) {
	// e^(rt) is the limit of (1 + r/m)^(mt); against math.Exp on random inputs the result must agree
	// with PV·e^(rt) computed independently in two steps to within 1 cent.
	rng := testRNG()
	for k := 0; k < 300; k++ {
		pv := rng.Int64N(1_000_000_000)
		r := rng.Float64() * 0.2
		y := rng.Float64() * 30
		got, err := FutureValueContinuous(pv, r, y)
		if err != nil {
			t.Fatal(err)
		}
		want := float64(pv) * math.Pow(math.E, r*y)
		if math.Abs(float64(got)-want) > 1+1e-12*want {
			t.Errorf("PV %d r %v t %v: got %d, want %.4f", pv, r, y, got, want)
		}
	}
}

func TestFutureValueContinuous_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "FutureValueContinuous", func() { sinkCents, _ = FutureValueContinuous(1000000, 0.08, 2) })
}

func BenchmarkFutureValueContinuous(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkCents, _ = FutureValueContinuous(1000000, 0.08, 2)
	}
}

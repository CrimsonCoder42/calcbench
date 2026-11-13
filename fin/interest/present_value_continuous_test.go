package interest

import (
	"math"
	"testing"
)

func TestPresentValueContinuous_SourceExample(t *testing.T) {
	// Strang & Herman, Calculus Volume 1 (OpenStax 2016) §6.8 Example 6.43: $1,000,000 in 40 years,
	// compounded continuously: at 5 % invest $135,335.28; at 6 % invest $90,717.95.
	got, err := PresentValueContinuous(100000000, 0.05, 40)
	wantCents(t, "Example 6.43 at 5%", got, err, 13533528)
	got, err = PresentValueContinuous(100000000, 0.06, 40)
	wantCents(t, "Example 6.43 at 6%", got, err, 9071795)
}

func TestPresentValueContinuous_Table(t *testing.T) {
	cases := []struct {
		name  string
		fv    int64
		rate  float64
		years float64
		want  int64
	}{
		{"zero rate", 5000, 0, 10, 5000},
		{"zero years", 5000, 0.1, 0, 5000},
		{"inverse of QIA Example 6", 1173511, 0.08, 2, 1000000},
		{"negative amount", -100000000, 0.05, 40, -13533528},
		{"long horizon rounds to 0", 100, 0.5, 100, 0},
	}
	for _, c := range cases {
		got, err := PresentValueContinuous(c.fv, c.rate, c.years)
		wantCents(t, c.name, got, err, c.want)
	}
}

func TestPresentValueContinuous_Errors(t *testing.T) {
	_, err := PresentValueContinuous(-maxExactCents-1, 0.05, 1)
	wantErr(t, "-2^53-1", err, ErrOutOfRange)
	for _, r := range badRates {
		_, err = PresentValueContinuous(100, r, 1)
		wantErr(t, "bad rate", err, ErrRate)
	}
	_, err = PresentValueContinuous(100, 0.05, math.NaN())
	wantErr(t, "NaN years", err, ErrNotFinite)
	_, err = PresentValueContinuous(100, 0.05, -0.5)
	wantErr(t, "negative years", err, ErrNegative)
	_, err = PresentValueContinuous(maxExactCents, -0.5, 2)
	wantErr(t, "negative rate overflows", err, ErrOverflow)
}

func TestPresentValueContinuous_Oracle(t *testing.T) {
	// Identity: PV(FV(x)) returns x within 1 cent.
	rng := testRNG()
	for k := 0; k < 300; k++ {
		x := rng.Int64N(1_000_000_000)
		r := rng.Float64() * 0.15
		y := rng.Float64() * 20
		fv, err := FutureValueContinuous(x, r, y)
		if err != nil {
			t.Fatal(err)
		}
		pv, err := PresentValueContinuous(fv, r, y)
		if err != nil {
			t.Fatal(err)
		}
		if d := pv - x; d < -1 || d > 1 {
			t.Errorf("x %d r %v t %v: PV(FV(x)) = %d", x, r, y, pv)
		}
	}
}

func TestPresentValueContinuous_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "PresentValueContinuous", func() { sinkCents, _ = PresentValueContinuous(100000000, 0.05, 40) })
}

func BenchmarkPresentValueContinuous(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkCents, _ = PresentValueContinuous(100000000, 0.05, 40)
	}
}

package interest

import (
	"math"
	"math/big"
	"testing"
)

func TestPresentValue_SourceExample(t *testing.T) {
	// Principles of Finance (OpenStax 2022) §7.3 Example 2: FV $50,000, I/Y 5, N 5 → $39,176.31.
	got, err := PresentValue(5000000, 0.05, 5)
	wantCents(t, "§7.3 Example 2", got, err, 3917631)
	// §7.4: $1,000 × 1/(1.05)^30 = $231.38.
	got, err = PresentValue(100000, 0.05, 30)
	wantCents(t, "§7.4", got, err, 23138)
}

func TestPresentValue_Table(t *testing.T) {
	cases := []struct {
		name    string
		fv      int64
		rate    float64
		periods int
		want    int64
	}{
		{"zero rate", 98765, 0, 12, 98765},
		{"zero periods", 98765, 0.09, 0, 98765},
		{"one period", 10500, 0.05, 1, 10000},
		{"negative amount", -5000000, 0.05, 5, -3917631},
		{"long horizon rounds to 0", 100, 1, 200, 0},
		{"max exact amount", maxExactCents, 0.01, 0, maxExactCents},
	}
	for _, c := range cases {
		got, err := PresentValue(c.fv, c.rate, c.periods)
		wantCents(t, c.name, got, err, c.want)
	}
}

func TestPresentValue_Errors(t *testing.T) {
	_, err := PresentValue(maxExactCents+1, 0.05, 1)
	wantErr(t, "2^53+1", err, ErrOutOfRange)
	for _, r := range badRates {
		_, err = PresentValue(100, r, 1)
		wantErr(t, "bad rate", err, ErrRate)
	}
	_, err = PresentValue(100, 0.05, -1)
	wantErr(t, "negative periods", err, ErrPeriods)
	_, err = PresentValue(maxExactCents, -0.5, 1)
	wantErr(t, "negative rate overflows", err, ErrOverflow)
}

func TestPresentValue_Oracle(t *testing.T) {
	rng := testRNG()
	for k := 0; k < 300; k++ {
		fv := rng.Int64N(2_000_000_000) - 1_000_000_000
		rate := rng.Float64()*0.15 - 0.05
		n := rng.IntN(120)
		got, err := PresentValue(fv, rate, n)
		if err != nil {
			t.Fatalf("PresentValue(%d, %v, %d): %v", fv, rate, n, err)
		}
		exact := new(big.Float).SetPrec(bigPrec).Mul(bigI(fv), bigCompound(rate, -n))
		checkCentsOracle(t, "oracle", got, exact)
	}
}

func TestPresentValue_InverseOfFutureValue(t *testing.T) {
	for _, n := range []int{1, 5, 30, 360} {
		fv, err := FutureValue(1_000_000_00, 0.004, n)
		if err != nil {
			t.Fatal(err)
		}
		pv, err := PresentValue(fv, 0.004, n)
		if err != nil {
			t.Fatal(err)
		}
		if math.Abs(float64(pv-1_000_000_00)) > 1 {
			t.Errorf("n=%d: PV(FV(x)) = %d, want 100000000 ± 1", n, pv)
		}
	}
}

func TestPresentValue_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "PresentValue", func() { sinkCents, _ = PresentValue(5000000, 0.05, 5) })
}

func BenchmarkPresentValue(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkCents, _ = PresentValue(5000000, 0.05, 360)
	}
}

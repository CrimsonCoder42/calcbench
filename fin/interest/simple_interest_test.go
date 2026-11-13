package interest

import (
	"math"
	"math/big"
	"testing"
)

func TestSimpleInterest_SourceExample(t *testing.T) {
	// Kirk, Contemporary Mathematics (OpenStax 2023) §6.3 Example 6.27: $4,000 at 5.5 % for 4 years → $880.
	got, err := SimpleInterest(400000, 0.055, 4)
	wantCents(t, "Example 6.27", got, err, 88000)
	// Example 6.28: $11,995 at 14.9 % for 2 years → $3,574.51.
	got, err = SimpleInterest(1199500, 0.149, 2)
	wantCents(t, "Example 6.28", got, err, 357451)
	// Example 6.30: $5,600 at 11.75 % for 45 days on Actual/365: $81.1233. Half-even gives $81.12; the
	// book rounds interest up to the next cent and shows $81.13.
	got, err = SimpleInterest(560000, 0.1175, 45.0/365)
	wantCents(t, "Example 6.30 (half-even)", got, err, 8112)
}

func TestSimpleInterest_Table(t *testing.T) {
	cases := []struct {
		name  string
		p     int64
		rate  float64
		years float64
		want  int64
	}{
		{"zero rate", 100000, 0, 3, 0},
		{"zero time", 100000, 0.05, 0, 0},
		{"one year", 100000, 0.05, 1, 5000},
		{"half cent rounds to even (down)", 50, 0.01, 1, 0},
		{"half cent rounds to even (up)", 150, 0.01, 1, 2},
		{"negative principal", -400000, 0.055, 4, -88000},
		{"negative rate", 100000, -0.01, 1, -1000},
	}
	for _, c := range cases {
		got, err := SimpleInterest(c.p, c.rate, c.years)
		wantCents(t, c.name, got, err, c.want)
	}
}

func TestSimpleInterest_Errors(t *testing.T) {
	_, err := SimpleInterest(maxExactCents+1, 0.05, 1)
	wantErr(t, "2^53+1", err, ErrOutOfRange)
	for _, r := range badRates {
		_, err = SimpleInterest(100, r, 1)
		wantErr(t, "bad rate", err, ErrRate)
	}
	_, err = SimpleInterest(100, 0.05, math.NaN())
	wantErr(t, "NaN time", err, ErrNotFinite)
	_, err = SimpleInterest(100, 0.05, -1)
	wantErr(t, "negative time", err, ErrNegative)
	_, err = SimpleInterest(maxExactCents, 0.5, 3)
	wantErr(t, "overflow", err, ErrOverflow)
}

func TestSimpleInterest_Oracle(t *testing.T) {
	rng := testRNG()
	for k := 0; k < 300; k++ {
		p := rng.Int64N(2_000_000_000) - 1_000_000_000
		r := rng.Float64() * 0.3
		y := rng.Float64() * 30
		got, err := SimpleInterest(p, r, y)
		if err != nil {
			t.Fatal(err)
		}
		exact := new(big.Float).SetPrec(bigPrec).Mul(bigI(p), bigF(r))
		exact.Mul(exact, bigF(y))
		checkCentsOracle(t, "oracle", got, exact)
	}
}

func TestSimpleInterest_NoAllocs(t *testing.T) {
	wantNoAllocs(t, "SimpleInterest", func() { sinkCents, _ = SimpleInterest(400000, 0.055, 4) })
}

func BenchmarkSimpleInterest(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkCents, _ = SimpleInterest(400000, 0.055, 4)
	}
}

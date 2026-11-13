package business

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// Hand-checked (the CFA ratio list gives no numbers): dividends $250,000 on net income $1,000,000 →
// (1,000,000 − 250,000) ÷ 1,000,000 = 0.75 = 1 − 0.25. Exact in float64.
func TestRetentionRatio_sourceExample(t *testing.T) {
	got, err := RetentionRatio(25000000, 100000000)
	if err != nil || got != 0.75 {
		t.Fatalf("got %v, %v; want 0.75", got, err)
	}
}

func TestRetentionRatio_table(t *testing.T) {
	cases := []struct {
		name    string
		div, ni int64
		want    float64
	}{
		{"no dividend", 0, 5, 1},
		{"all paid out", 7, 7, 0},
		{"more than earnings", 3, 2, -0.5},
		{"max dividend", maxExactCents, 1, 1 - maxExactCents},
		{"max income", 0, maxExactCents, 1},
		{"max both", maxExactCents, maxExactCents, 0},
	}
	for _, c := range cases {
		got, err := RetentionRatio(c.div, c.ni)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestRetentionRatio_errors(t *testing.T) {
	cases := []struct {
		name    string
		div, ni int64
		want    error
	}{
		{"zero income", 1, 0, ErrZeroDenominator},
		{"loss", 1, -1, ErrNegative},
		{"negative dividends", -1, 1, ErrNegative},
		{"dividends above 2^53", maxExactCents + 1, 1, ErrOutOfRange},
		{"income above 2^53", 1, maxExactCents + 1, ErrOutOfRange},
		{"income min int64", 1, math.MinInt64, ErrOutOfRange},
	}
	for _, c := range cases {
		if _, err := RetentionRatio(c.div, c.ni); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

// Against math/big on 500 fixed-seed inputs: the difference is exact, so the result is the correctly
// rounded quotient and the comparison is exact. Small cases also check the identity with the payout ratio
// to 1 ulp of 1 (1 − payout is rounded twice).
func TestRetentionRatio_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for i := 0; i < 500; i++ {
		div := int64(rng.Uint64N(maxExactCents + 1))
		ni := 1 + int64(rng.Uint64N(maxExactCents))
		got, err := RetentionRatio(div, ni)
		if err != nil {
			t.Fatalf("(%d, %d): %v", div, ni, err)
		}
		if want := ratOracle(ni-div, ni); got != want {
			t.Fatalf("(%d, %d) = %v, want %v", div, ni, got, want)
		}
		if div <= ni {
			payout, _ := DividendPayoutRatio(div, ni)
			if math.Abs(got-(1-payout)) > 0x1p-52 {
				t.Fatalf("(%d, %d): retention %v, 1 − payout %v", div, ni, got, 1-payout)
			}
		}
	}
}

func TestRetentionRatio_allocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = RetentionRatio(25000000, 100000000) }); n != 0 {
		t.Fatalf("allocs = %v, want 0", n)
	}
}

func BenchmarkRetentionRatio(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = RetentionRatio(25000000, 100000000)
	}
}

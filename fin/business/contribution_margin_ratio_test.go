package business

import (
	"errors"
	"math/rand/v2"
	"testing"
)

// Principles of Accounting Vol. 2 §3.1, Hicks Manufacturing: Blue Jay ($100 − $20) ÷ $100 = 0.80; April
// totals ($50,000 − $10,000) ÷ $50,000 = 0.80. Exact: 0.8 is the correctly rounded 8000/10000.
func TestContributionMarginRatio_sourceExample(t *testing.T) {
	for _, c := range [][2]int64{{10000, 2000}, {5000000, 1000000}} {
		got, err := ContributionMarginRatio(c[0], c[1])
		if err != nil || got != 0.8 {
			t.Errorf("%v: got %v, %v; want 0.8", c, got, err)
		}
	}
}

func TestContributionMarginRatio_table(t *testing.T) {
	cases := []struct {
		name           string
		sales, varCost int64
		want           float64
	}{
		{"no variable cost", 100, 0, 1},
		{"all variable", 100, 100, 0},
		{"negative margin", 100, 200, -1},
		{"max", maxExactCents, 0, 1},
	}
	for _, c := range cases {
		got, err := ContributionMarginRatio(c.sales, c.varCost)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestContributionMarginRatio_errors(t *testing.T) {
	cases := []struct {
		name           string
		sales, varCost int64
		want           error
	}{
		{"zero sales", 0, 1, ErrZeroDenominator},
		{"sales above 2^53", maxExactCents + 1, 1, ErrOutOfRange},
		{"variable above 2^53", 1, maxExactCents + 1, ErrOutOfRange},
		{"negative sales", -1, 1, ErrNegative},
		{"negative variable", 1, -1, ErrNegative},
	}
	for _, c := range cases {
		if _, err := ContributionMarginRatio(c.sales, c.varCost); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

// Oracle: (sales − variable)/sales against math/big; exact (the difference is exact in float64).
func TestContributionMarginRatio_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for i := 0; i < 500; i++ {
		sales := 1 + int64(rng.Uint64N(maxExactCents))
		v := int64(rng.Uint64N(maxExactCents + 1))
		got, err := ContributionMarginRatio(sales, v)
		if err != nil {
			t.Fatal(err)
		}
		if want := ratOracle(sales-v, sales); got != want {
			t.Fatalf("(%d, %d) = %v, want %v", sales, v, got, want)
		}
	}
}

func TestContributionMarginRatio_allocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = ContributionMarginRatio(10000, 2000) }); n != 0 {
		t.Fatalf("allocs = %v, want 0", n)
	}
}

func BenchmarkContributionMarginRatio(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = ContributionMarginRatio(10000, 2000)
	}
}

package business

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// Sinks shared by the package's benchmarks and allocation tests.
var (
	sinkFloat float64
	sinkInt   int64
	sinkSlice []int64
)

// checkRatioOracle compares a two-input cents ratio against math/big on 500 inputs drawn with a fixed seed:
// num in [minNum, 2⁵³], den in [minDen, 2⁵³] without 0. Both inputs are exact in float64, so the IEEE
// quotient is the correctly rounded value of the exact rational, which big.Rat.Float64 also returns: the
// comparison is exact.
func checkRatioOracle(t *testing.T, f func(int64, int64) (float64, error), minNum, minDen int64) {
	t.Helper()
	rng := rand.New(rand.NewPCG(1, 2))
	draw := func(lo int64) int64 {
		span := uint64(maxExactCents - lo + 1)
		return lo + int64(rng.Uint64N(span))
	}
	for i := 0; i < 500; i++ {
		num := draw(minNum)
		den := draw(minDen)
		if den == 0 {
			continue
		}
		got, err := f(num, den)
		if err != nil {
			t.Fatalf("(%d, %d): %v", num, den, err)
		}
		want, _ := new(big.Rat).SetFrac(big.NewInt(num), big.NewInt(den)).Float64()
		if got != want {
			t.Fatalf("(%d, %d) = %v, want %v", num, den, got, want)
		}
	}
}

// Principles of Finance §6.3, Clear Lake Sporting Goods: $200,000 ÷ $100,000 = 2. Exact comparison: both
// values are exact in float64 and the quotient is an integer.
func TestCurrentRatio_sourceExample(t *testing.T) {
	got, err := CurrentRatio(20000000, 10000000)
	if err != nil || got != 2 {
		t.Fatalf("got %v, %v; want 2", got, err)
	}
}

func TestCurrentRatio_table(t *testing.T) {
	cases := []struct {
		name     string
		num, den int64
		want     float64
	}{
		{"zero assets", 0, 5, 0},
		{"one", 7, 7, 1},
		{"max assets", maxExactCents, 1, maxExactCents},
		{"max liabilities", 1, maxExactCents, 1.0 / maxExactCents},
		{"below one", 1, 4, 0.25},
	}
	for _, c := range cases {
		got, err := CurrentRatio(c.num, c.den)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestCurrentRatio_errors(t *testing.T) {
	cases := []struct {
		name     string
		num, den int64
		want     error
	}{
		{"zero liabilities", 1, 0, ErrZeroDenominator},
		{"assets above 2^53", maxExactCents + 1, 1, ErrOutOfRange},
		{"liabilities above 2^53", 1, maxExactCents + 1, ErrOutOfRange},
		{"liabilities max int64", 1, math.MaxInt64, ErrOutOfRange},
		{"negative assets", -1, 1, ErrNegative},
		{"negative liabilities", 1, -1, ErrNegative},
	}
	for _, c := range cases {
		if _, err := CurrentRatio(c.num, c.den); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestCurrentRatio_oracle(t *testing.T) {
	checkRatioOracle(t, CurrentRatio, 0, 1)
}

func TestCurrentRatio_allocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = CurrentRatio(20000000, 10000000) }); n != 0 {
		t.Fatalf("allocs = %v, want 0", n)
	}
}

func BenchmarkCurrentRatio(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = CurrentRatio(20000000, 10000000)
	}
}

// ratOracle returns num/den correctly rounded to float64 via math/big (den ≠ 0).
func ratOracle(num, den int64) float64 {
	f, _ := new(big.Rat).SetFrac(big.NewInt(num), big.NewInt(den)).Float64()
	return f
}

package business

import (
	"errors"
	"math/big"
	"math/rand/v2"
	"slices"
	"strconv"
	"testing"
)

// 26 CFR 1.167(b)-3(a), Example (1): cost $1,750, salvage $250, 5 years → 5/15, 4/15, … of $1,500.
func TestDepreciationSumOfYearsDigits_sourceExample(t *testing.T) {
	got, err := DepreciationSumOfYearsDigits(175000, 25000, 5)
	if err != nil {
		t.Fatal(err)
	}
	if want := []int64{50000, 40000, 30000, 20000, 10000}; !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestDepreciationSumOfYearsDigits_table(t *testing.T) {
	cases := []struct {
		name          string
		cost, salvage int64
		life          int
		want          []int64
	}{
		{"one period", 1000, 100, 1, []int64{900}},
		{"nothing to depreciate", 500, 500, 3, []int64{0, 0, 0}},
		// 1000 × 3/6 = 500, × 2/6 = 333.33 → 333, last 167
		{"remainder in last period", 1000, 0, 3, []int64{500, 333, 167}},
		// 10 × 2/3 = 6.67 → 7, last 3
		{"rounding up", 10, 0, 2, []int64{7, 3}},
		// 3 × 2/3 = 2 exact; 9 × 1/6... half-even: 3 × 3/6 = 1.5 → 2
		{"half rounds to even", 3, 0, 3, []int64{2, 1, 0}},
		// 11,000 − 1,000 over 5 years: 3,333.33, 2,666.67, 2,000, 1,333.33, 666.67
		{"five-year asset", 1100000, 100000, 5, []int64{333333, 266667, 200000, 133333, 66667}},
	}
	for _, c := range cases {
		got, err := DepreciationSumOfYearsDigits(c.cost, c.salvage, c.life)
		if err != nil || !slices.Equal(got, c.want) {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestDepreciationSumOfYearsDigits_largest(t *testing.T) {
	// 2⁵³ cents over 1 000 000 periods: the 128-bit products must not overflow and the sum must be exact.
	got, err := DepreciationSumOfYearsDigits(maxExactCents, 0, 1_000_000)
	if err != nil {
		t.Fatal(err)
	}
	var s int64
	for _, v := range got {
		if v < 0 {
			t.Fatalf("negative amount %d", v)
		}
		s += v
	}
	if s != maxExactCents {
		t.Fatalf("sum %d, want %d", s, int64(maxExactCents))
	}
}

func TestDepreciationSumOfYearsDigits_oracle(t *testing.T) {
	// Each non-final amount equals base × (n−k+1)/S rounded half-even, checked with math/big.
	rng := rand.New(rand.NewPCG(1, 2))
	for it := 0; it < 300; it++ {
		cost := rng.Int64N(1 << 40)
		salvage := rng.Int64N(cost + 1)
		life := 1 + rng.IntN(60)
		got, err := DepreciationSumOfYearsDigits(cost, salvage, life)
		if err != nil {
			t.Fatal(err)
		}
		sum := int64(life) * int64(life+1) / 2
		var total int64
		for k, v := range got {
			total += v
			if v < 0 {
				t.Fatalf("negative amount %d", v)
			}
			if k == life-1 {
				continue
			}
			num := new(big.Rat).SetFrac(big.NewInt((cost-salvage)*int64(life-k)), big.NewInt(sum))
			fl, _ := num.Float64()
			if d := float64(v) - fl; d > 0.5 || d < -0.5 {
				t.Fatalf("cost %d salvage %d life %d period %d: %d vs %v", cost, salvage, life, k+1, v, fl)
			}
		}
		if total != cost-salvage {
			t.Fatalf("total %d, want %d", total, cost-salvage)
		}
	}
}

func TestDepreciationSumOfYearsDigits_errors(t *testing.T) {
	cases := []struct {
		name          string
		cost, salvage int64
		life          int
		want          error
	}{
		{"negative cost", -1, 0, 5, ErrNegative},
		{"negative salvage", 100, -1, 5, ErrNegative},
		{"salvage above cost", 100, 101, 5, ErrDomain},
		{"cost beyond 2^53", maxExactCents + 1, 0, 5, ErrOutOfRange},
		{"zero life", 100, 0, 0, ErrPeriods},
		{"life too long", 100, 0, 1_000_001, ErrPeriods},
	}
	for _, c := range cases {
		if _, err := DepreciationSumOfYearsDigits(c.cost, c.salvage, c.life); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestDepreciationSumOfYearsDigits_oneAllocation(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { _, _ = DepreciationSumOfYearsDigits(1100000, 100000, 5) }); a != 1 {
		t.Errorf("allocations = %v, want 1", a)
	}
}

func BenchmarkDepreciationSumOfYearsDigits(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = DepreciationSumOfYearsDigits(1<<50, 0, n)
			}
		})
	}
}

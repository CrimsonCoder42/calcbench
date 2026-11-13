package business

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

var usc1a = struct {
	starts []int64
	rates  []float64
}{[]int64{0, 3690000, 8915000}, []float64{0.15, 0.28, 0.31}}

// 26 U.S.C. § 1(a) table rows: the base amounts $5,535 and $20,165 are the lower tiers in full; $100,000 →
// $20,165 + 31 % × $10,850 = $23,528.50 (hand computation from the table).
func TestTieredMarginalAmount_sourceExample(t *testing.T) {
	cases := []struct {
		amount, want int64
	}{
		{3690000, 553500},   // $5,535 at the first boundary
		{8915000, 2016500},  // $20,165 at the second boundary
		{10000000, 2352850}, // $23,528.50
	}
	for _, c := range cases {
		got, err := TieredMarginalAmount(c.amount, usc1a.starts, usc1a.rates)
		if err != nil || got != c.want {
			t.Errorf("amount %d: got %d, %v; want %d", c.amount, got, err, c.want)
		}
	}
}

func TestTieredMarginalAmount_table(t *testing.T) {
	cases := []struct {
		name   string
		amount int64
		starts []int64
		rates  []float64
		want   int64
	}{
		{"zero amount", 0, []int64{0, 100}, []float64{0.1, 0.2}, 0},
		{"one flat tier", 20000000, []int64{0}, []float64{0.03}, 600000},
		{"graduated commission 3 % to 20k, 5 % above", 20000000, []int64{0, 2000000}, []float64{0.03, 0.05}, 960000},
		{"below the second tier", 150000, []int64{0, 2000000}, []float64{0.03, 0.05}, 4500},
		{"half cent rounds to even (down)", 5, []int64{0}, []float64{0.5}, 2}, // 2.5 → 2
		{"half cent rounds to even (up)", 7, []int64{0}, []float64{0.5}, 4},   // 3.5 → 4
		{"zero-rate first tier", 1000, []int64{0, 500}, []float64{0, 1}, 500},
		{"largest amount", maxExactCents, []int64{0}, []float64{1}, maxExactCents},
	}
	for _, c := range cases {
		got, err := TieredMarginalAmount(c.amount, c.starts, c.rates)
		if err != nil || got != c.want {
			t.Errorf("%s: got %d, %v; want %d", c.name, got, err, c.want)
		}
	}
}

func TestTieredMarginalAmount_oracle(t *testing.T) {
	// Against exact rational arithmetic on the same float64 rates, for random schedules of up to 8 tiers.
	rng := rand.New(rand.NewPCG(1, 2))
	for it := 0; it < 300; it++ {
		n := 1 + rng.IntN(8)
		starts := make([]int64, n)
		rates := make([]float64, n)
		for k := range starts {
			if k > 0 {
				starts[k] = starts[k-1] + 1 + rng.Int64N(1<<30)
			}
			rates[k] = rng.Float64()
		}
		amount := rng.Int64N(starts[n-1] + 1<<31)
		got, err := TieredMarginalAmount(amount, starts, rates)
		if err != nil {
			t.Fatal(err)
		}
		exact := new(big.Rat)
		for k := range starts {
			if amount <= starts[k] {
				break
			}
			top := amount
			if k+1 < n && starts[k+1] < top {
				top = starts[k+1]
			}
			r := new(big.Rat).SetFloat64(rates[k])
			exact.Add(exact, r.Mul(r, new(big.Rat).SetInt64(top-starts[k])))
		}
		f, _ := exact.Float64()
		if math.Abs(float64(got)-f) > 1 {
			t.Fatalf("amount %d: got %d, exact %v", amount, got, f)
		}
	}
}

func TestTieredMarginalAmount_errors(t *testing.T) {
	cases := []struct {
		name   string
		amount int64
		starts []int64
		rates  []float64
		want   error
	}{
		{"negative amount", -1, []int64{0}, []float64{0.1}, ErrNegative},
		{"amount beyond 2^53", maxExactCents + 1, []int64{0}, []float64{0.1}, ErrOutOfRange},
		{"empty schedule", 100, nil, nil, ErrOutOfRange},
		{"length mismatch", 100, []int64{0, 10}, []float64{0.1}, ErrOutOfRange},
		{"first start not 0", 100, []int64{5}, []float64{0.1}, ErrDomain},
		{"starts not increasing", 100, []int64{0, 10, 10}, []float64{0.1, 0.2, 0.3}, ErrDomain},
		{"start beyond 2^53", 100, []int64{0, maxExactCents + 1}, []float64{0.1, 0.2}, ErrOutOfRange},
		{"NaN rate", 100, []int64{0}, []float64{math.NaN()}, ErrNotFinite},
		{"rate above 1", 100, []int64{0}, []float64{1.5}, ErrRate},
		{"negative rate", 100, []int64{0}, []float64{-0.1}, ErrRate},
	}
	for _, c := range cases {
		if _, err := TieredMarginalAmount(c.amount, c.starts, c.rates); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestTieredMarginalAmount_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { _, _ = TieredMarginalAmount(10000000, usc1a.starts, usc1a.rates) }); a != 0 {
		t.Errorf("allocations = %v, want 0", a)
	}
}

func BenchmarkTieredMarginalAmount(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = TieredMarginalAmount(10000000, usc1a.starts, usc1a.rates)
	}
}

package business

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"slices"
	"strconv"
	"testing"
)

// Principles of Accounting Vol. 1 §11.3, Kenzie Company: cost $58,000, salvage $10,000, 5 years, 40 %
// (2 × 20 %): $23,200, $13,920, $8,352, $2,528 (capped at book value − salvage), $0; total $48,000.
// Exact cents.
func TestDepreciationDecliningBalance_sourceExample(t *testing.T) {
	got, err := DepreciationDecliningBalance(5800000, 1000000, 5, 2)
	if err != nil {
		t.Fatal(err)
	}
	if want := []int64{2320000, 1392000, 835200, 252800, 0}; !slices.Equal(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestDepreciationDecliningBalance_table(t *testing.T) {
	cases := []struct {
		name          string
		cost, salvage int64
		life          int
		multiple      float64
		want          []int64
	}{
		{"one period capped at salvage", 1000, 100, 1, 2, []int64{900}},
		{"nothing to depreciate", 500, 500, 3, 2, []int64{0, 0, 0}},
		{"rate 100 % in period one", 1000, 0, 2, 2, []int64{1000, 0}},
		// Salvage 0 is never reached by a declining balance: no final adjustment (source gives none).
		{"salvage never reached", 1000, 0, 4, 2, []int64{500, 250, 125, 62}},
		{"150 % declining", 10000, 0, 3, 1.5, []int64{5000, 2500, 1250}},
		{"tie rounds to even", 5, 0, 2, 1, []int64{2, 2}},
		{"huge multiple", 1000, 10, 3, 1e300, []int64{990, 0, 0}},
		// 9 × 2 ÷ 12 = 1.5 exactly → half-even 2 (dividing the rate first, 9 × (2/12), gives 1.4999… → 1).
		{"exact tie in period one", 9, 0, 12, 2, []int64{2, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0}},
	}
	for _, c := range cases {
		got, err := DepreciationDecliningBalance(c.cost, c.salvage, c.life, c.multiple)
		if err != nil || !slices.Equal(got, c.want) {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

// 5,999,999,999,997 × 2 ÷ 12 = 999,999,999,999.5 exactly → half-even 1,000,000,000,000 (the rate-first
// computation lands below the tie and gives 999,999,999,999).
func TestDepreciationDecliningBalance_largeTie(t *testing.T) {
	got, err := DepreciationDecliningBalance(5999999999997, 0, 12, 2)
	if err != nil || got[0] != 1000000000000 {
		t.Fatalf("period 1 = %v, %v; want 1000000000000", got[0], err)
	}
}

func TestDepreciationDecliningBalance_errors(t *testing.T) {
	cases := []struct {
		name          string
		cost, salvage int64
		life          int
		multiple      float64
		want          error
	}{
		{"cost above 2^53", maxExactCents + 1, 0, 5, 2, ErrOutOfRange},
		{"negative cost", -1, 0, 5, 2, ErrNegative},
		{"negative salvage", 10, -1, 5, 2, ErrNegative},
		{"salvage above cost", 10, 11, 5, 2, ErrDomain},
		{"zero life", 10, 0, 0, 2, ErrPeriods},
		{"life above limit", 10, 0, maxLifePeriods + 1, 2, ErrPeriods},
		{"zero multiple", 10, 0, 5, 0, ErrDomain},
		{"negative multiple", 10, 0, 5, -2, ErrDomain},
		{"NaN multiple", 10, 0, 5, math.NaN(), ErrNotFinite},
		{"Inf multiple", 10, 0, 5, math.Inf(1), ErrNotFinite},
	}
	for _, c := range cases {
		got, err := DepreciationDecliningBalance(c.cost, c.salvage, c.life, c.multiple)
		if !errors.Is(err, c.want) || got != nil {
			t.Errorf("%s: got %v, err = %v, want nil, %v", c.name, got, err, c.want)
		}
	}
}

// Invariants on 300 fixed-seed inputs: no negative amount, book value never below salvage, the sum never
// above cost − salvage, every amount equals round-half-even(book × rate) or the cap book − salvage, and
// once salvage is reached every later amount is 0. When the book value reaches salvage the schedule
// reconciles to cost − salvage to the cent.
func TestDepreciationDecliningBalance_invariants(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for i := 0; i < 300; i++ {
		cost := int64(rng.Uint64N(maxExactCents + 1))
		salvage := int64(rng.Uint64N(uint64(cost) + 1))
		if i%2 == 0 {
			salvage = cost / int64(2+rng.IntN(20))
		}
		life := 1 + rng.IntN(60)
		multiple := 0.5 + 2.5*rng.Float64()
		got, err := DepreciationDecliningBalance(cost, salvage, life, multiple)
		if err != nil || len(got) != life {
			t.Fatalf("i=%d: len %d, err %v", i, len(got), err)
		}
		book := cost
		for j, e := range got {
			remaining := book - salvage
			want := remaining
			if x := float64(book) * multiple / float64(life); x < float64(remaining) {
				want = min(int64(math.RoundToEven(x)), remaining)
			}
			if e != want || e < 0 {
				t.Fatalf("i=%d period %d: got %d, want %d", i, j, e, want)
			}
			book -= e
		}
		if book < salvage {
			t.Fatalf("i=%d: book %d below salvage %d", i, book, salvage)
		}
		if book == salvage {
			var sum int64
			for _, e := range got {
				sum += e
			}
			if sum != cost-salvage {
				t.Fatalf("i=%d: sum %d, want %d", i, sum, cost-salvage)
			}
		}
	}
}

// Accuracy contract against math/big: each amount equals min(halfEven(book × multiple ÷ life), book −
// salvage) computed exactly, with 1 cent of slack only when the exact value lies within 1e-6 cent of a
// half cent. Costs ≤ 2⁴³ cents.
func TestDepreciationDecliningBalance_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	half := big.NewRat(1, 2)
	slack := big.NewRat(1, 1_000_000)
	for i := 0; i < 300; i++ {
		cost := int64(rng.Uint64N(1<<43 + 1))
		salvage := int64(rng.Uint64N(uint64(cost)/4 + 1))
		life := 1 + rng.IntN(40)
		multiple := []float64{2, 1.5, 1.25, 0.5 + 2.5*rng.Float64()}[i%4]
		got, err := DepreciationDecliningBalance(cost, salvage, life, multiple)
		if err != nil {
			t.Fatal(err)
		}
		book := cost
		for j, e := range got {
			remaining := book - salvage
			exact := new(big.Rat).SetFloat64(multiple)
			exact.Mul(exact, new(big.Rat).SetInt64(book))
			exact.Quo(exact, new(big.Rat).SetInt64(int64(life)))
			want := min(halfEvenRat(exact), remaining)
			if e != want {
				// distance of the exact value from the nearest half cent
				fl := new(big.Rat).SetInt64(new(big.Int).Quo(exact.Num(), exact.Denom()).Int64())
				d := new(big.Rat).Sub(exact, fl)
				d.Sub(d, half)
				d.Abs(d)
				if e-want > 1 || want-e > 1 || d.Cmp(slack) > 0 {
					t.Fatalf("i=%d period %d: got %d, want %d", i, j, e, want)
				}
			}
			book -= e
		}
	}
}

// One allocation: the result slice.
func TestDepreciationDecliningBalance_allocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkSlice, _ = DepreciationDecliningBalance(5800000, 1000000, 5, 2) }); n != 1 {
		t.Fatalf("allocs = %v, want 1", n)
	}
}

// Salvage 0 is never reached, so every period runs the full multiply-and-round step at every size.
func BenchmarkDepreciationDecliningBalance(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkSlice, _ = DepreciationDecliningBalance(maxExactCents, 0, n, 2)
			}
		})
	}
}

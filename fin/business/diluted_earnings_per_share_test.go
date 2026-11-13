package business

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// CFA Program Level I sample question Q5, Sublyme Corporation: net income USD200 million, 50,000,000
// weighted average shares, 2,000,000 convertible preferred shares each convertible into two common shares
// (USD5 dividend, assumed converted, so no preferred dividend is deducted): 200,000,000 ÷ 54,000,000 =
// 3.7037… → USD3.70. Exact cents.
// Hand-checked: net income $1,000,000, preferred dividends $100,000, after-tax convertible interest
// $60,000, 1,000,000 shares + 200,000 from conversion → (100,000,000 − 10,000,000 + 6,000,000)¢ ÷
// 1,200,000 = 96,000,000 ÷ 1,200,000 = 80¢.
func TestDilutedEarningsPerShare_sourceExample(t *testing.T) {
	cases := []struct {
		ni, pref, interest, shares, dilutive, want int64
	}{
		{20000000000, 0, 0, 50000000, 4000000, 370},
		{100000000, 10000000, 6000000, 1000000, 200000, 80},
	}
	for _, c := range cases {
		got, err := DilutedEarningsPerShare(c.ni, c.pref, c.interest, c.shares, c.dilutive)
		if err != nil || got != c.want {
			t.Errorf("%v: got %v, %v; want %v", c, got, err, c.want)
		}
	}
}

// With no dilutive instruments the result equals basic EPS; ties go to the even cent in both signs.
func TestDilutedEarningsPerShare_table(t *testing.T) {
	cases := []struct {
		name                                 string
		ni, pref, interest, shares, dilutive int64
		want                                 int64
	}{
		{"equals basic EPS", 3500000, 500000, 0, 8000, 0, 375},
		{"zero income", 0, 0, 0, 7, 3, 0},
		{"tie to even down", 5, 0, 0, 1, 1, 2},
		{"tie to even up", 7, 0, 0, 1, 1, 4},
		{"negative tie", -5, 0, 0, 1, 1, -2},
		{"interest added back", 0, 0, 10, 4, 1, 2},
		{"loss after preferred", 0, 10, 0, 4, 0, -2},
		{"max numerator", maxExactCents, 0, 0, 1, 0, maxExactCents},
		{"max shares", maxExactCents, 0, 0, maxExactCents, maxExactCents, 0},
		{"interest offsets preferred", 1000, maxExactCents, maxExactCents, 1, 0, 1000},
	}
	for _, c := range cases {
		got, err := DilutedEarningsPerShare(c.ni, c.pref, c.interest, c.shares, c.dilutive)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestDilutedEarningsPerShare_errors(t *testing.T) {
	cases := []struct {
		name                                 string
		ni, pref, interest, shares, dilutive int64
		want                                 error
	}{
		{"zero shares", 1, 0, 0, 0, 5, ErrZeroDenominator},
		{"negative shares", 1, 0, 0, -1, 0, ErrNegative},
		{"negative dilutive", 1, 0, 0, 1, -1, ErrNegative},
		{"shares above 2^53", 1, 0, 0, maxExactCents + 1, 0, ErrOutOfRange},
		{"dilutive above 2^53", 1, 0, 0, 1, maxExactCents + 1, ErrOutOfRange},
		{"dilutive max int64", 1, 0, 0, 1, math.MaxInt64, ErrOutOfRange},
		{"income above 2^53", maxExactCents + 1, 0, 0, 1, 0, ErrOutOfRange},
		{"preferred above 2^53", 0, maxExactCents + 1, 0, 1, 0, ErrOutOfRange},
		{"interest above 2^53", 0, 0, maxExactCents + 1, 1, 0, ErrOutOfRange},
		{"negative preferred", 0, -1, 0, 1, 0, ErrNegative},
		{"negative interest", 0, 0, -1, 1, 0, ErrNegative},
		{"result above 2^53", maxExactCents, 0, maxExactCents, 1, 0, ErrOverflow},
		{"result below -2^53", -maxExactCents, maxExactCents, 0, 1, 0, ErrOverflow},
	}
	for _, c := range cases {
		if _, err := DilutedEarningsPerShare(c.ni, c.pref, c.interest, c.shares, c.dilutive); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

// Against math/big (exact rational, rounded half-even by halfEvenRat) on 500 fixed-seed inputs.
func TestDilutedEarningsPerShare_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for i := 0; i < 500; i++ {
		ni := int64(rng.Uint64N(2*maxExactCents+1)) - maxExactCents
		pref := int64(rng.Uint64N(maxExactCents + 1))
		interest := int64(rng.Uint64N(maxExactCents + 1))
		shares := 1 + int64(rng.Uint64N(1<<uint(1+rng.IntN(53))))
		dilutive := int64(rng.Uint64N(1 << uint(1+rng.IntN(53))))
		if i%3 == 0 { // small values exercise ties
			ni, pref, interest = int64(rng.IntN(2001))-1000, int64(rng.IntN(100)), int64(rng.IntN(100))
			shares, dilutive = 1+int64(rng.IntN(8)), int64(rng.IntN(8))
		}
		num := big.NewInt(ni - pref + interest)
		want := halfEvenRat(new(big.Rat).SetFrac(num, big.NewInt(shares+dilutive)))
		got, err := DilutedEarningsPerShare(ni, pref, interest, shares, dilutive)
		if want > maxExactCents || want < -maxExactCents {
			if !errors.Is(err, ErrOverflow) {
				t.Fatalf("(%d,%d,%d,%d,%d): err = %v, want ErrOverflow", ni, pref, interest, shares, dilutive, err)
			}
			continue
		}
		if err != nil || got != want {
			t.Fatalf("(%d,%d,%d,%d,%d) = %v, %v; want %v", ni, pref, interest, shares, dilutive, got, err, want)
		}
	}
}

func TestDilutedEarningsPerShare_allocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() {
		sinkInt, _ = DilutedEarningsPerShare(100000000, 10000000, 6000000, 1000000, 200000)
	}); n != 0 {
		t.Fatalf("allocs = %v, want 0", n)
	}
}

func BenchmarkDilutedEarningsPerShare(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = DilutedEarningsPerShare(100000000, 10000000, 6000000, 1000000, 200000)
	}
}

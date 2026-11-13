package business

import (
	"errors"
	"math/big"
	"math/rand/v2"
	"testing"
)

// Principles of Finance §6.5, Clear Lake Sporting Goods: ($35,000 − $5,000) ÷ 8,000 = $3.75.
// Principles of Accounting Vol. 1 §14.5, Sanaron: ($50,000 − $10,000) ÷ 5,000 = $8.00; Ranadune:
// $15,000 ÷ 17,500 = $0.857…, printed $0.86. Exact cents.
func TestEarningsPerShare_sourceExample(t *testing.T) {
	cases := []struct {
		ni, pref, shares, want int64
	}{
		{3500000, 500000, 8000, 375},
		{5000000, 1000000, 5000, 800},
		{1500000, 0, 20000, 75},
		{1500000, 0, 17500, 86},
	}
	for _, c := range cases {
		got, err := EarningsPerShare(c.ni, c.pref, c.shares)
		if err != nil || got != c.want {
			t.Errorf("%v: got %v, %v; want %v", c, got, err, c.want)
		}
	}
}

// Ties go to the even cent, in both signs.
func TestEarningsPerShare_table(t *testing.T) {
	cases := []struct {
		name             string
		ni, pref, shares int64
		want             int64
	}{
		{"zero income", 0, 0, 7, 0},
		{"tie to even down", 5, 0, 2, 2},
		{"tie to even up", 7, 0, 2, 4},
		{"negative tie to even", -5, 0, 2, -2},
		{"negative tie up", -7, 0, 2, -4},
		{"below half", 4, 0, 3, 1},
		{"above half", 5, 0, 3, 2},
		{"negative above half", -5, 0, 3, -2},
		{"loss after preferred", 0, 10, 4, -2},
		{"max", maxExactCents, 0, 1, maxExactCents},
		{"one share max", -maxExactCents, 0, 1, -maxExactCents},
		{"max shares", maxExactCents, 0, maxExactCents, 1},
	}
	for _, c := range cases {
		got, err := EarningsPerShare(c.ni, c.pref, c.shares)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestEarningsPerShare_errors(t *testing.T) {
	cases := []struct {
		name             string
		ni, pref, shares int64
		want             error
	}{
		{"zero shares", 1, 0, 0, ErrZeroDenominator},
		{"negative shares", 1, 0, -1, ErrNegative},
		{"shares above 2^53", 1, 0, maxExactCents + 1, ErrOutOfRange},
		{"income above 2^53", maxExactCents + 1, 0, 1, ErrOutOfRange},
		{"preferred above 2^53", 0, maxExactCents + 1, 1, ErrOutOfRange},
		{"negative preferred", 0, -1, 1, ErrNegative},
		{"result below -2^53", -maxExactCents, maxExactCents, 1, ErrOverflow},
	}
	for _, c := range cases {
		if _, err := EarningsPerShare(c.ni, c.pref, c.shares); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

// halfEvenRat rounds an exact rational half-even, the oracle for every per-share result.
func halfEvenRat(r *big.Rat) int64 {
	num, den := new(big.Int).Set(r.Num()), r.Denom()    // den > 0
	q, m := new(big.Int).DivMod(num, den, new(big.Int)) // floor division, 0 ≤ m < den
	twice := new(big.Int).Lsh(m, 1)
	switch twice.Cmp(den) {
	case 1:
		q.Add(q, big.NewInt(1))
	case 0:
		if q.Bit(0) == 1 {
			q.Add(q, big.NewInt(1))
		}
	}
	return q.Int64()
}

func TestEarningsPerShare_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for i := 0; i < 500; i++ {
		ni := int64(rng.Uint64N(2*maxExactCents+1)) - maxExactCents
		pref := int64(rng.Uint64N(maxExactCents + 1))
		shares := 1 + int64(rng.Uint64N(1<<uint(1+rng.IntN(53))))
		if i%3 == 0 { // small values exercise ties
			ni, pref, shares = int64(rng.IntN(2001))-1000, int64(rng.IntN(100)), 1+int64(rng.IntN(8))
		}
		want := halfEvenRat(new(big.Rat).SetFrac(big.NewInt(ni-pref), big.NewInt(shares)))
		got, err := EarningsPerShare(ni, pref, shares)
		if want > maxExactCents || want < -maxExactCents {
			if !errors.Is(err, ErrOverflow) {
				t.Fatalf("(%d,%d,%d): err = %v, want ErrOverflow", ni, pref, shares, err)
			}
			continue
		}
		if err != nil || got != want {
			t.Fatalf("(%d,%d,%d) = %v, %v; want %v", ni, pref, shares, got, err, want)
		}
	}
}

func TestEarningsPerShare_allocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = EarningsPerShare(3500000, 500000, 8000) }); n != 0 {
		t.Fatalf("allocs = %v, want 0", n)
	}
}

func BenchmarkEarningsPerShare(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = EarningsPerShare(3500000, 500000, 8000)
	}
}

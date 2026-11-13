package business

import (
	"errors"
	"math/big"
	"math/rand/v2"
	"testing"
)

// Principles of Finance §6.5, Clear Lake Sporting Goods: ($100,000 − $20,000) ÷ 8,000 = $10. Exact cents.
func TestBookValuePerShare_sourceExample(t *testing.T) {
	got, err := BookValuePerShare(10000000, 2000000, 8000)
	if err != nil || got != 1000 {
		t.Fatalf("got %v, %v; want 1000", got, err)
	}
}

func TestBookValuePerShare_table(t *testing.T) {
	cases := []struct {
		name             string
		eq, pref, shares int64
		want             int64
	}{
		{"zero equity", 0, 0, 3, 0},
		{"negative equity", -10, 0, 4, -2},
		{"tie to even", 25, 0, 10, 2},
		{"tie to even up", 35, 0, 10, 4},
		{"max", maxExactCents, 0, 1, maxExactCents},
	}
	for _, c := range cases {
		got, err := BookValuePerShare(c.eq, c.pref, c.shares)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestBookValuePerShare_errors(t *testing.T) {
	cases := []struct {
		name             string
		eq, pref, shares int64
		want             error
	}{
		{"zero shares", 1, 0, 0, ErrZeroDenominator},
		{"negative shares", 1, 0, -1, ErrNegative},
		{"shares above 2^53", 1, 0, maxExactCents + 1, ErrOutOfRange},
		{"equity below -2^53", -maxExactCents - 1, 0, 1, ErrOutOfRange},
		{"preferred above 2^53", 0, maxExactCents + 1, 1, ErrOutOfRange},
		{"negative preferred", 0, -1, 1, ErrNegative},
		{"result below -2^53", -1, maxExactCents, 1, ErrOverflow},
	}
	for _, c := range cases {
		if _, err := BookValuePerShare(c.eq, c.pref, c.shares); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestBookValuePerShare_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for i := 0; i < 500; i++ {
		eq := int64(rng.Uint64N(2*maxExactCents+1)) - maxExactCents
		pref := int64(rng.Uint64N(maxExactCents + 1))
		shares := 1 + int64(rng.Uint64N(1<<uint(1+rng.IntN(53))))
		want := halfEvenRat(new(big.Rat).SetFrac(big.NewInt(eq-pref), big.NewInt(shares)))
		got, err := BookValuePerShare(eq, pref, shares)
		if want > maxExactCents || want < -maxExactCents {
			if !errors.Is(err, ErrOverflow) {
				t.Fatalf("(%d,%d,%d): err = %v, want ErrOverflow", eq, pref, shares, err)
			}
			continue
		}
		if err != nil || got != want {
			t.Fatalf("(%d,%d,%d) = %v, %v; want %v", eq, pref, shares, got, err, want)
		}
	}
}

func TestBookValuePerShare_allocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = BookValuePerShare(10000000, 2000000, 8000) }); n != 0 {
		t.Fatalf("allocs = %v, want 0", n)
	}
}

func BenchmarkBookValuePerShare(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = BookValuePerShare(10000000, 2000000, 8000)
	}
}

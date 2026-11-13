package business

import (
	"errors"
	"math/big"
	"math/rand/v2"
	"testing"
)

// Principles of Accounting Vol. 1 §14.3, La Cantina: 9,200 outstanding shares × $0.50 = $4,600, so
// $4,600 ÷ 9,200 = $0.50. Exact cents.
func TestDividendsPerShare_sourceExample(t *testing.T) {
	got, err := DividendsPerShare(460000, 9200)
	if err != nil || got != 50 {
		t.Fatalf("got %v, %v; want 50", got, err)
	}
}

func TestDividendsPerShare_table(t *testing.T) {
	cases := []struct {
		name          string
		total, shares int64
		want          int64
	}{
		{"no dividend", 0, 10, 0},
		{"one share", 123, 1, 123},
		{"tie to even", 1, 2, 0},
		{"tie to even up", 3, 2, 2},
		{"max", maxExactCents, 1, maxExactCents},
		{"max shares", maxExactCents, maxExactCents, 1},
	}
	for _, c := range cases {
		got, err := DividendsPerShare(c.total, c.shares)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestDividendsPerShare_errors(t *testing.T) {
	cases := []struct {
		name          string
		total, shares int64
		want          error
	}{
		{"zero shares", 1, 0, ErrZeroDenominator},
		{"negative shares", 1, -1, ErrNegative},
		{"shares above 2^53", 1, maxExactCents + 1, ErrOutOfRange},
		{"total above 2^53", maxExactCents + 1, 1, ErrOutOfRange},
		{"negative total", -1, 1, ErrNegative},
	}
	for _, c := range cases {
		if _, err := DividendsPerShare(c.total, c.shares); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestDividendsPerShare_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for i := 0; i < 500; i++ {
		total := int64(rng.Uint64N(maxExactCents + 1))
		shares := 1 + int64(rng.Uint64N(1<<uint(1+rng.IntN(53))))
		want := halfEvenRat(new(big.Rat).SetFrac(big.NewInt(total), big.NewInt(shares)))
		got, err := DividendsPerShare(total, shares)
		if err != nil || got != want {
			t.Fatalf("(%d,%d) = %v, %v; want %v", total, shares, got, err, want)
		}
	}
}

func TestDividendsPerShare_allocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = DividendsPerShare(460000, 9200) }); n != 0 {
		t.Fatalf("allocs = %v, want 0", n)
	}
}

func BenchmarkDividendsPerShare(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = DividendsPerShare(460000, 9200)
	}
}

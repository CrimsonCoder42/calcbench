package business

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// Principles of Finance §6.3, Clear Lake Sporting Goods: ($110,000 + $20,000 + $30,000) ÷ $100,000 = 1.6.
// Exact comparison with the correctly rounded quotient 16000000/10000000.
func TestQuickRatio_sourceExample(t *testing.T) {
	got, err := QuickRatio(11000000, 2000000, 3000000, 10000000)
	if err != nil || got != 1.6 {
		t.Fatalf("got %v, %v; want 1.6", got, err)
	}
}

func TestQuickRatio_table(t *testing.T) {
	cases := []struct {
		name               string
		cash, ms, ar, liab int64
		want               float64
	}{
		{"all zero assets", 0, 0, 0, 5, 0},
		{"cash only", 4, 0, 0, 4, 1},
		{"max cash", maxExactCents, 0, 0, 1, maxExactCents},
		{"sum above 2^53", maxExactCents, maxExactCents, maxExactCents, 3, maxExactCents},
	}
	for _, c := range cases {
		got, err := QuickRatio(c.cash, c.ms, c.ar, c.liab)
		if err != nil || got != c.want {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestQuickRatio_errors(t *testing.T) {
	big := int64(maxExactCents + 1)
	cases := []struct {
		name               string
		cash, ms, ar, liab int64
		want               error
	}{
		{"zero liabilities", 1, 1, 1, 0, ErrZeroDenominator},
		{"cash above 2^53", big, 0, 0, 1, ErrOutOfRange},
		{"securities above 2^53", 0, big, 0, 1, ErrOutOfRange},
		{"receivables above 2^53", 0, 0, big, 1, ErrOutOfRange},
		{"liabilities above 2^53", 0, 0, 0, big, ErrOutOfRange},
		{"negative cash", -1, 0, 0, 1, ErrNegative},
		{"negative securities", 0, -1, 0, 1, ErrNegative},
		{"negative receivables", 0, 0, -1, 1, ErrNegative},
		{"negative liabilities", 0, 0, 0, -1, ErrNegative},
	}
	for _, c := range cases {
		if _, err := QuickRatio(c.cash, c.ms, c.ar, c.liab); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

// When the three-term sum is ≤ 2⁵³ the result is the correctly rounded quotient (exact comparison against
// big.Rat); above it, two roundings: relative error ≤ 2⁻⁵².
func TestQuickRatio_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for i := 0; i < 500; i++ {
		cash := int64(rng.Uint64N(maxExactCents + 1))
		ms := int64(rng.Uint64N(maxExactCents + 1))
		ar := int64(rng.Uint64N(maxExactCents + 1))
		if i%2 == 0 {
			cash, ms, ar = cash>>2, ms>>2, ar>>2
		}
		liab := 1 + int64(rng.Uint64N(maxExactCents))
		got, err := QuickRatio(cash, ms, ar, liab)
		if err != nil {
			t.Fatal(err)
		}
		sum := cash + ms + ar
		want, _ := new(big.Rat).SetFrac(big.NewInt(sum), big.NewInt(liab)).Float64()
		if sum <= maxExactCents {
			if got != want {
				t.Fatalf("(%d,%d,%d,%d) = %v, want %v", cash, ms, ar, liab, got, want)
			}
		} else if math.Abs(got-want) > want*0x1p-52 {
			t.Fatalf("(%d,%d,%d,%d) = %v, want %v within 2^-52 relative", cash, ms, ar, liab, got, want)
		}
	}
}

func TestQuickRatio_allocs(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = QuickRatio(11000000, 2000000, 3000000, 10000000) }); n != 0 {
		t.Fatalf("allocs = %v, want 0", n)
	}
}

func BenchmarkQuickRatio(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = QuickRatio(11000000, 2000000, 3000000, 10000000)
	}
}

package combinatorics

import (
	"errors"
	"math"
	"math/big"
	"testing"
)

// DLMF Table 26.5.1: C(10) = 16 796. OEIS A000108: C(36) = 11 959 798 385 860 453 492 (last below 2⁶⁴),
// C(37) = 45 950 804 324 621 742 364.
func TestCatalan_sourceExample(t *testing.T) {
	if got, err := Catalan(10); err != nil || got != 16796 {
		t.Errorf("Catalan(10) = (%d, %v), want (16796, nil)", got, err)
	}
	if got, err := Catalan(36); err != nil || got != 11959798385860453492 {
		t.Errorf("Catalan(36) = (%d, %v), want (11959798385860453492, nil)", got, err)
	}
	want, _ := new(big.Int).SetString("45950804324621742364", 10)
	if got, err := CatalanBig(37); err != nil || got.Cmp(want) != 0 {
		t.Errorf("CatalanBig(37) = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestCatalan_cases(t *testing.T) {
	for _, c := range []struct {
		name string
		n    int
		want uint64
	}{{"zero", 0, 1}, {"one", 1, 1}, {"two", 2, 2}, {"three", 3, 5}} {
		t.Run(c.name, func(t *testing.T) {
			if got, err := Catalan(c.n); err != nil || got != c.want {
				t.Errorf("Catalan(%d) = (%d, %v), want (%d, nil)", c.n, got, err, c.want)
			}
		})
	}
}

func TestCatalan_errors(t *testing.T) {
	for _, n := range []int{-1, math.MinInt} {
		if got, err := Catalan(n); !errors.Is(err, ErrNegative) || got != 0 {
			t.Errorf("Catalan(%d) = (%d, %v), want (0, ErrNegative)", n, got, err)
		}
		if got, err := CatalanBig(n); !errors.Is(err, ErrNegative) || got != nil {
			t.Errorf("CatalanBig(%d) = (%v, %v), want (nil, ErrNegative)", n, got, err)
		}
	}
	for _, n := range []int{37, math.MaxInt} {
		if got, err := Catalan(n); !errors.Is(err, ErrOverflow) || got != 0 {
			t.Errorf("Catalan(%d) = (%d, %v), want (0, ErrOverflow)", n, got, err)
		}
	}
	if got, err := CatalanBig(math.MaxInt); !errors.Is(err, ErrOverflow) || got != nil {
		t.Errorf("CatalanBig(MaxInt) = (%v, %v), want (nil, ErrOverflow)", got, err)
	}
}

// The table and the big twin against the recurrence C(n+1) = 2(2n+1)/(n+2)·C(n), DLMF 26.5.4, in math/big.
func TestCatalan_oracle(t *testing.T) {
	c := big.NewInt(1)
	for n := 0; n <= 120; n++ {
		got, err := Catalan(n)
		wantU64(t, "Catalan", got, err, c)
		gotBig, err := CatalanBig(n)
		if err != nil || gotBig.Cmp(c) != 0 {
			t.Fatalf("CatalanBig(%d) = (%v, %v), want %v", n, gotBig, err, c)
		}
		c = new(big.Int).Mul(c, big.NewInt(int64(2*(2*n+1))))
		c.Quo(c, big.NewInt(int64(n+2)))
	}
}

func TestCatalan_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkU, _ = Catalan(36) }); a != 0 {
		t.Errorf("Catalan allocates %v times, want 0", a)
	}
}

func BenchmarkCatalan(b *testing.B) {
	b.Run("n=36", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sinkU, _ = Catalan(36)
		}
	})
}

func BenchmarkCatalanBig(b *testing.B) {
	for _, c := range []struct {
		name string
		n    int
	}{{"n=100", 100}, {"n=1000", 1000}, {"n=10000", 10000}} {
		b.Run(c.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkBig, _ = CatalanBig(c.n)
			}
		})
	}
}

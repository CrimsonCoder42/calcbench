package combinatorics

import (
	"errors"
	"math"
	"math/big"
	"testing"
)

// DLMF Table 26.7.1: B(5) = 52, B(10) = 115 975. OEIS A000110: B(25) = 4 638 590 332 229 999 353 (last
// below 2⁶⁴), B(26) = 49 631 246 523 618 756 274.
func TestBell_sourceExample(t *testing.T) {
	for _, c := range []struct {
		n    int
		want uint64
	}{{5, 52}, {10, 115975}, {25, 4638590332229999353}} {
		if got, err := Bell(c.n); err != nil || got != c.want {
			t.Errorf("Bell(%d) = (%d, %v), want (%d, nil)", c.n, got, err, c.want)
		}
	}
	want, _ := new(big.Int).SetString("49631246523618756274", 10)
	if got, err := BellBig(26); err != nil || got.Cmp(want) != 0 {
		t.Errorf("BellBig(26) = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestBell_cases(t *testing.T) {
	for _, c := range []struct {
		name string
		n    int
		want uint64
	}{{"zero", 0, 1}, {"one", 1, 1}, {"two", 2, 2}, {"three", 3, 5}} {
		t.Run(c.name, func(t *testing.T) {
			if got, err := Bell(c.n); err != nil || got != c.want {
				t.Errorf("Bell(%d) = (%d, %v), want (%d, nil)", c.n, got, err, c.want)
			}
		})
	}
}

func TestBell_errors(t *testing.T) {
	for _, n := range []int{-1, math.MinInt} {
		if got, err := Bell(n); !errors.Is(err, ErrNegative) || got != 0 {
			t.Errorf("Bell(%d) = (%d, %v), want (0, ErrNegative)", n, got, err)
		}
		if got, err := BellBig(n); !errors.Is(err, ErrNegative) || got != nil {
			t.Errorf("BellBig(%d) = (%v, %v), want (nil, ErrNegative)", n, got, err)
		}
	}
	for _, n := range []int{26, math.MaxInt} {
		if got, err := Bell(n); !errors.Is(err, ErrOverflow) || got != 0 {
			t.Errorf("Bell(%d) = (%d, %v), want (0, ErrOverflow)", n, got, err)
		}
	}
}

// The table and the big twin against B(n+1) = Σₖ C(n,k)·B(k), DLMF 26.7.6, in math/big (independent of
// the Stirling row BellBig sums).
func TestBell_oracle(t *testing.T) {
	const top = 100
	b := []*big.Int{big.NewInt(1)}
	for n := 0; n < top; n++ {
		next := new(big.Int)
		for k := 0; k <= n; k++ {
			next.Add(next, new(big.Int).Mul(new(big.Int).Binomial(int64(n), int64(k)), b[k]))
		}
		b = append(b, next)
	}
	for n := 0; n <= top; n++ {
		got, err := Bell(n)
		wantU64(t, "Bell", got, err, b[n])
		gotBig, err := BellBig(n)
		if err != nil || gotBig.Cmp(b[n]) != 0 {
			t.Fatalf("BellBig(%d) = (%v, %v), want %v", n, gotBig, err, b[n])
		}
	}
}

func TestBell_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkU, _ = Bell(25) }); a != 0 {
		t.Errorf("Bell allocates %v times, want 0", a)
	}
}

func BenchmarkBell(b *testing.B) {
	b.Run("n=25", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sinkU, _ = Bell(25)
		}
	})
}

func BenchmarkBellBig(b *testing.B) {
	for _, c := range []struct {
		name string
		n    int
	}{{"n=100", 100}, {"n=300", 300}, {"n=1000", 1000}} {
		b.Run(c.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkBig, _ = BellBig(c.n)
			}
		})
	}
}

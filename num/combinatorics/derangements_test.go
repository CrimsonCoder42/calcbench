package combinatorics

import (
	"errors"
	"math"
	"math/big"
	"testing"
)

// bigDerangements returns d(n) = n!·Σⱼ₌₀..ₙ (−1)ʲ/j! = Σⱼ (−1)ʲ·n!/j!, DLMF eq. 26.13.4, exactly.
func bigDerangements(n int) *big.Int {
	sum := new(big.Int)
	for j := 0; j <= n; j++ {
		term := big.NewInt(1)
		if j < n {
			term.MulRange(int64(j+1), int64(n))
		}
		if j%2 == 0 {
			sum.Add(sum, term)
		} else {
			sum.Sub(sum, term)
		}
	}
	return sum
}

// OEIS A000166: d(4) = 9, d(6) = 265, d(10) = 1 334 961, d(20) = 895 014 631 192 902 121.
func TestDerangements_sourceExample(t *testing.T) {
	for _, c := range []struct {
		n    int
		want uint64
	}{{4, 9}, {6, 265}, {10, 1334961}, {20, 895014631192902121}} {
		if got, err := Derangements(c.n); err != nil || got != c.want {
			t.Errorf("Derangements(%d) = (%d, %v), want (%d, nil)", c.n, got, err, c.want)
		}
	}
	// OEIS A000166: d(22) = 413 496 759 611 120 779 881.
	want, _ := new(big.Int).SetString("413496759611120779881", 10)
	if got, err := DerangementsBig(22); err != nil || got.Cmp(want) != 0 {
		t.Errorf("DerangementsBig(22) = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestDerangements_cases(t *testing.T) {
	for _, c := range []struct {
		name string
		n    int
		want uint64
	}{{"zero", 0, 1}, {"one", 1, 0}, {"two", 2, 1}, {"three", 3, 2}} {
		t.Run(c.name, func(t *testing.T) {
			if got, err := Derangements(c.n); err != nil || got != c.want {
				t.Errorf("Derangements(%d) = (%d, %v), want (%d, nil)", c.n, got, err, c.want)
			}
		})
	}
}

func TestDerangements_errors(t *testing.T) {
	for _, n := range []int{-1, math.MinInt} {
		if got, err := Derangements(n); !errors.Is(err, ErrNegative) || got != 0 {
			t.Errorf("Derangements(%d) = (%d, %v), want (0, ErrNegative)", n, got, err)
		}
		if got, err := DerangementsBig(n); !errors.Is(err, ErrNegative) || got != nil {
			t.Errorf("DerangementsBig(%d) = (%v, %v), want (nil, ErrNegative)", n, got, err)
		}
	}
	// d(21) = 18 795 307 255 050 944 540 > 2⁶⁴ − 1.
	for _, n := range []int{21, math.MaxInt} {
		if got, err := Derangements(n); !errors.Is(err, ErrOverflow) || got != 0 {
			t.Errorf("Derangements(%d) = (%d, %v), want (0, ErrOverflow)", n, got, err)
		}
	}
}

// Every table entry, and the big twin up to 80, against DLMF 26.13.4 summed in math/big.
func TestDerangements_oracle(t *testing.T) {
	for n := range 81 {
		exact := bigDerangements(n)
		got, err := Derangements(n)
		wantU64(t, "Derangements", got, err, exact)
		gotBig, err := DerangementsBig(n)
		if err != nil || gotBig.Cmp(exact) != 0 {
			t.Fatalf("DerangementsBig(%d) = (%v, %v), want %v", n, gotBig, err, exact)
		}
	}
}

func TestDerangements_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkU, _ = Derangements(20) }); a != 0 {
		t.Errorf("Derangements allocates %v times, want 0", a)
	}
}

func BenchmarkDerangements(b *testing.B) {
	b.Run("n=20", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sinkU, _ = Derangements(20)
		}
	})
}

func BenchmarkDerangementsBig(b *testing.B) {
	for _, c := range []struct {
		name string
		n    int
	}{{"n=100", 100}, {"n=1000", 1000}, {"n=10000", 10000}} {
		b.Run(c.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkBig, _ = DerangementsBig(c.n)
			}
		})
	}
}

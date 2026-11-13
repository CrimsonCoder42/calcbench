package combinatorics

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"slices"
	"testing"
)

// bigMultinomial returns (Σc)!/Πc! by the factorial definition (DLMF 26.4.2, left side).
func bigMultinomial(counts []int) *big.Int {
	sum := 0
	den := big.NewInt(1)
	for _, c := range counts {
		sum += c
		den.Mul(den, bigFactorial(c))
	}
	return new(big.Int).Quo(bigFactorial(sum), den)
}

// DLMF eq. 26.4.2. Worked examples (hand-checked with exact integers): (5,3,2) → 10!/(5!3!2!) = 2520;
// the letters of MISSISSIPPI, M1 I4 S4 P2 → 11!/(1!4!4!2!) = 34 650.
func TestMultinomial_sourceExample(t *testing.T) {
	for _, c := range []struct {
		counts []int
		want   uint64
	}{{[]int{5, 3, 2}, 2520}, {[]int{1, 4, 4, 2}, 34650}} {
		if got, err := Multinomial(c.counts); err != nil || got != c.want {
			t.Errorf("Multinomial(%v) = (%d, %v), want (%d, nil)", c.counts, got, err, c.want)
		}
		if got, err := MultinomialBig(c.counts); err != nil || got.Uint64() != c.want {
			t.Errorf("MultinomialBig(%v) = (%v, %v), want (%d, nil)", c.counts, got, err, c.want)
		}
	}
}

func TestMultinomial_cases(t *testing.T) {
	ones20 := slices.Repeat([]int{1}, 20)
	tests := []struct {
		name   string
		counts []int
		want   uint64
	}{
		{"one part", []int{7}, 1},
		{"one empty part", []int{0}, 1},
		{"zeros", []int{0, 0, 0}, 1},
		{"two parts is a binomial", []int{3, 4}, 35},
		{"twenty ones is 20!", ones20, 2432902008176640000},
		{"max int then one", []int{math.MaxInt, 1}, 1 << 63},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Multinomial(tt.counts); err != nil || got != tt.want {
				t.Errorf("Multinomial(%v) = (%d, %v), want (%d, nil)", tt.counts, got, err, tt.want)
			}
		})
	}
	for _, counts := range [][]int{slices.Repeat([]int{1}, 21), {30, 30, 30}, {math.MaxInt, 1, math.MaxInt}, {math.MaxInt, math.MaxInt}} {
		if got, err := Multinomial(counts); !errors.Is(err, ErrOverflow) || got != 0 {
			t.Errorf("Multinomial(%v) = (%d, %v), want (0, ErrOverflow)", counts, got, err)
		}
	}
}

func TestMultinomial_errors(t *testing.T) {
	for _, counts := range [][]int{nil, {}} {
		if got, err := Multinomial(counts); !errors.Is(err, ErrEmpty) || got != 0 {
			t.Errorf("Multinomial(%v) = (%d, %v), want (0, ErrEmpty)", counts, got, err)
		}
		if got, err := MultinomialBig(counts); !errors.Is(err, ErrEmpty) || got != nil {
			t.Errorf("MultinomialBig(%v) = (%v, %v), want (nil, ErrEmpty)", counts, got, err)
		}
	}
	// A negative count is reported even when an earlier prefix would already overflow.
	neg := []int{30, 30, 30, -1}
	if got, err := Multinomial(neg); !errors.Is(err, ErrNegative) || got != 0 {
		t.Errorf("Multinomial(%v) = (%d, %v), want (0, ErrNegative)", neg, got, err)
	}
	if got, err := MultinomialBig(neg); !errors.Is(err, ErrNegative) || got != nil {
		t.Errorf("MultinomialBig(%v) = (%v, %v), want (nil, ErrNegative)", neg, got, err)
	}
	if got, err := MultinomialBig([]int{math.MaxInt, 1}); !errors.Is(err, ErrOverflow) || got != nil {
		t.Errorf("MultinomialBig(MaxInt, 1) = (%v, %v), want (nil, ErrOverflow)", got, err)
	}
}

func TestMultinomial_doesNotModifyInput(t *testing.T) {
	counts := []int{4, 0, 3, 2}
	orig := slices.Clone(counts)
	_, _ = Multinomial(counts)
	_, _ = MultinomialBig(counts)
	if !slices.Equal(counts, orig) {
		t.Errorf("input modified: %v, want %v", counts, orig)
	}
}

func TestMultinomial_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		counts := make([]int, 1+rng.IntN(6))
		for i := range counts {
			counts[i] = rng.IntN(25)
		}
		exact := bigMultinomial(counts)
		got, err := Multinomial(counts)
		wantU64(t, "Multinomial", got, err, exact)
		gotBig, err := MultinomialBig(counts)
		if err != nil || gotBig.Cmp(exact) != 0 {
			t.Fatalf("trial %d: MultinomialBig(%v) = (%v, %v), want %v", trial, counts, gotBig, err, exact)
		}
	}
}

func TestMultinomial_noAlloc(t *testing.T) {
	counts := []int{1, 4, 4, 2}
	if a := testing.AllocsPerRun(100, func() { sinkU, _ = Multinomial(counts) }); a != 0 {
		t.Errorf("Multinomial allocates %v times, want 0", a)
	}
}

func BenchmarkMultinomial(b *testing.B) {
	counts := []int{5, 4, 4, 3, 2, 2}
	b.Run("m=6", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sinkU, _ = Multinomial(counts)
		}
	})
}

func BenchmarkMultinomialBig(b *testing.B) {
	for _, c := range []struct {
		name   string
		counts []int
	}{{"sum=100", []int{40, 30, 20, 10}}, {"sum=1000", []int{400, 300, 200, 100}}, {"sum=10000", []int{4000, 3000, 2000, 1000}}} {
		b.Run(c.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkBig, _ = MultinomialBig(c.counts)
			}
		})
	}
}

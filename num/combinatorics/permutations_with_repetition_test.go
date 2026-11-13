package combinatorics

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// DLMF Table 26.17.1, unrestricted maps of k labelled balls into n labelled boxes: nᵏ. Three-letter
// codes over 26 letters: 26³ = 17 576 (hand-checked).
func TestPermutationsWithRepetition_sourceExample(t *testing.T) {
	if got, err := PermutationsWithRepetition(26, 3); err != nil || got != 17576 {
		t.Errorf("PermutationsWithRepetition(26, 3) = (%d, %v), want (17576, nil)", got, err)
	}
}

func TestPermutationsWithRepetition_cases(t *testing.T) {
	tests := []struct {
		name string
		n, k int
		want uint64
	}{
		{"zero to the zero", 0, 0, 1},
		{"zero symbols", 0, 5, 0},
		{"one symbol, huge k", 1, math.MaxInt, 1},
		{"k zero", math.MaxInt, 0, 1},
		{"k one", math.MaxInt, 1, math.MaxInt},
		{"two to the 63", 2, 63, 1 << 63},
		{"three to the 40", 3, 40, 12157665459056928801},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := PermutationsWithRepetition(tt.n, tt.k); err != nil || got != tt.want {
				t.Errorf("PermutationsWithRepetition(%d, %d) = (%d, %v), want (%d, nil)", tt.n, tt.k, got, err, tt.want)
			}
		})
	}
	for _, c := range [][2]int{{2, 64}, {3, 41}, {math.MaxInt, 2}, {2, math.MaxInt}} {
		if got, err := PermutationsWithRepetition(c[0], c[1]); !errors.Is(err, ErrOverflow) || got != 0 {
			t.Errorf("PermutationsWithRepetition(%d, %d) = (%d, %v), want (0, ErrOverflow)", c[0], c[1], got, err)
		}
	}
}

func TestPermutationsWithRepetition_errors(t *testing.T) {
	for _, c := range [][2]int{{-1, 0}, {0, -1}, {math.MinInt, 2}} {
		if got, err := PermutationsWithRepetition(c[0], c[1]); !errors.Is(err, ErrNegative) || got != 0 {
			t.Errorf("PermutationsWithRepetition(%d, %d) = (%d, %v), want (0, ErrNegative)", c[0], c[1], got, err)
		}
	}
}

func TestPermutationsWithRepetition_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 600 {
		n := rng.IntN(300)
		k := rng.IntN(70)
		exact := new(big.Int).Exp(big.NewInt(int64(n)), big.NewInt(int64(k)), nil)
		got, err := PermutationsWithRepetition(n, k)
		wantU64(t, "PermutationsWithRepetition", got, err, exact)
	}
}

func TestPermutationsWithRepetition_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkU, _ = PermutationsWithRepetition(3, 40) }); a != 0 {
		t.Errorf("PermutationsWithRepetition allocates %v times, want 0", a)
	}
}

func BenchmarkPermutationsWithRepetition(b *testing.B) {
	b.Run("n=3,k=40", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sinkU, _ = PermutationsWithRepetition(3, 40)
		}
	})
}

package combinatorics

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// DLMF Table 26.17.1, k unlabelled balls into n labelled boxes: C(n+k−1, k). Five kinds, size three:
// C(7,3) = 35 (hand-checked).
func TestCombinationsWithRepetition_sourceExample(t *testing.T) {
	if got, err := CombinationsWithRepetition(5, 3); err != nil || got != 35 {
		t.Errorf("CombinationsWithRepetition(5, 3) = (%d, %v), want (35, nil)", got, err)
	}
}

func TestCombinationsWithRepetition_cases(t *testing.T) {
	tests := []struct {
		name string
		n, k int
		want uint64
	}{
		{"no kinds, empty multiset", 0, 0, 1},
		{"no kinds, k > 0", 0, 3, 0},
		{"k zero", 3, 0, 1},
		{"one kind", 1, math.MaxInt, 1},
		{"max int kinds, size one", math.MaxInt, 1, math.MaxInt},
		{"two kinds, n+k−1 = 2⁶³ beyond int", 2, math.MaxInt, 1 << 63},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := CombinationsWithRepetition(tt.n, tt.k); err != nil || got != tt.want {
				t.Errorf("CombinationsWithRepetition(%d, %d) = (%d, %v), want (%d, nil)", tt.n, tt.k, got, err, tt.want)
			}
		})
	}
	if got, err := CombinationsWithRepetition(math.MaxInt, math.MaxInt); !errors.Is(err, ErrOverflow) || got != 0 {
		t.Errorf("CombinationsWithRepetition(MaxInt, MaxInt) = (%d, %v), want (0, ErrOverflow)", got, err)
	}
}

func TestCombinationsWithRepetition_errors(t *testing.T) {
	for _, c := range [][2]int{{-1, 0}, {0, -1}, {math.MinInt, 1}} {
		if got, err := CombinationsWithRepetition(c[0], c[1]); !errors.Is(err, ErrNegative) || got != 0 {
			t.Errorf("CombinationsWithRepetition(%d, %d) = (%d, %v), want (0, ErrNegative)", c[0], c[1], got, err)
		}
	}
}

func TestCombinationsWithRepetition_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 500 {
		n := 1 + rng.IntN(120)
		k := rng.IntN(120)
		got, err := CombinationsWithRepetition(n, k)
		wantU64(t, "CombinationsWithRepetition", got, err, new(big.Int).Binomial(int64(n+k-1), int64(k)))
	}
}

func TestCombinationsWithRepetition_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkU, _ = CombinationsWithRepetition(30, 30) }); a != 0 {
		t.Errorf("CombinationsWithRepetition allocates %v times, want 0", a)
	}
}

func BenchmarkCombinationsWithRepetition(b *testing.B) {
	b.Run("n=30,k=30", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sinkU, _ = CombinationsWithRepetition(30, 30)
		}
	})
}

package combinatorics

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// P(10,3) = 10·9·8 = 720 (DLMF §26.1, falling factorial; hand-checked), and P(5,5) = 5! = 120.
func TestPermutations_sourceExample(t *testing.T) {
	if got, err := Permutations(10, 3); err != nil || got != 720 {
		t.Errorf("Permutations(10, 3) = (%d, %v), want (720, nil)", got, err)
	}
	if got, err := Permutations(5, 5); err != nil || got != 120 {
		t.Errorf("Permutations(5, 5) = (%d, %v), want (120, nil)", got, err)
	}
}

func TestPermutations_cases(t *testing.T) {
	tests := []struct {
		name string
		n, k int
		want uint64
	}{
		{"zero of zero", 0, 0, 1},
		{"k zero", 7, 0, 1},
		{"k above n", 5, 6, 0},
		{"k equals n is n!", 20, 20, 2432902008176640000},
		{"largest k for n = 100", 100, 9, 690281878632192000},
		{"max int, k = 1", math.MaxInt, 1, math.MaxInt},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Permutations(tt.n, tt.k); err != nil || got != tt.want {
				t.Errorf("Permutations(%d, %d) = (%d, %v), want (%d, nil)", tt.n, tt.k, got, err, tt.want)
			}
		})
	}
	for _, c := range [][2]int{{100, 10}, {21, 21}, {math.MaxInt, 2}, {math.MaxInt, math.MaxInt}} {
		if got, err := Permutations(c[0], c[1]); !errors.Is(err, ErrOverflow) || got != 0 {
			t.Errorf("Permutations(%d, %d) = (%d, %v), want (0, ErrOverflow)", c[0], c[1], got, err)
		}
	}
}

func TestPermutations_errors(t *testing.T) {
	for _, c := range [][2]int{{-1, 0}, {0, -1}, {math.MinInt, math.MinInt}} {
		if got, err := Permutations(c[0], c[1]); !errors.Is(err, ErrNegative) || got != 0 {
			t.Errorf("Permutations(%d, %d) = (%d, %v), want (0, ErrNegative)", c[0], c[1], got, err)
		}
	}
}

func TestPermutations_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 500 {
		n := rng.IntN(150)
		k := rng.IntN(n + 3)
		exact := new(big.Int)
		switch {
		case k > n:
		case k == 0:
			exact.SetInt64(1)
		default:
			exact.MulRange(int64(n-k+1), int64(n))
		}
		got, err := Permutations(n, k)
		wantU64(t, "Permutations", got, err, exact)
	}
}

func TestPermutations_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { sinkU, _ = Permutations(100, 9) }); a != 0 {
		t.Errorf("Permutations allocates %v times, want 0", a)
	}
}

func BenchmarkPermutations(b *testing.B) {
	b.Run("n=100,k=9", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			sinkU, _ = Permutations(100, 9)
		}
	})
}

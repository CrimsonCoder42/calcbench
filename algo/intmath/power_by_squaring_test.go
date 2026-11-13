package intmath

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// benchSink keeps benchmark results alive. Shared by every benchmark in this package.
var benchSink uint64

// benchSizes are the collection sizes of .claude/rules/go-complexity.md. Shared by the benchmarks.
var benchSizes = []struct {
	name string
	n    int
}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}}

// Knuth Vol. 2 §4.6.3, Algorithm A: x^23 by the right-to-left binary method (23 = 10111₂); here 2^23.
func TestPowerBySquaring_sourceExample(t *testing.T) {
	if got, err := PowerBySquaring(2, 23); err != nil || got != 8_388_608 {
		t.Errorf("PowerBySquaring(2, 23) = (%d, %v), want 8388608", got, err)
	}
}

func TestPowerBySquaring_cases(t *testing.T) {
	tests := []struct{ base, exp, want int64 }{
		{0, 0, 1}, {0, 5, 0}, {5, 0, 1}, {1, math.MaxInt64, 1}, {-1, math.MaxInt64, -1}, {-1, 2, 1},
		{3, 4, 81}, {-3, 3, -27}, {10, 18, 1_000_000_000_000_000_000}, {2, 62, 1 << 62},
		{-2, 63, math.MinInt64}, {-2, 62, 1 << 62}, {math.MaxInt64, 1, math.MaxInt64},
		{math.MinInt64, 1, math.MinInt64}, {3, 39, 4_052_555_153_018_976_267},
	}
	for _, tt := range tests {
		if got, err := PowerBySquaring(tt.base, tt.exp); err != nil || got != tt.want {
			t.Errorf("PowerBySquaring(%d, %d) = (%d, %v), want %d", tt.base, tt.exp, got, err, tt.want)
		}
	}
}

func TestPowerBySquaring_errors(t *testing.T) {
	tests := []struct {
		base, exp int64
		wantErr   error
	}{
		{2, -1, ErrNegative}, {0, math.MinInt64, ErrNegative},
		{2, 63, ErrOverflow}, {-2, 64, ErrOverflow}, {10, 19, ErrOverflow}, {3, 40, ErrOverflow},
		{math.MinInt64, 2, ErrOverflow}, {math.MaxInt64, 2, ErrOverflow}, {2, math.MaxInt64, ErrOverflow},
	}
	for _, tt := range tests {
		if _, err := PowerBySquaring(tt.base, tt.exp); !errors.Is(err, tt.wantErr) {
			t.Errorf("PowerBySquaring(%d, %d) error = %v, want %v", tt.base, tt.exp, err, tt.wantErr)
		}
	}
}

// Oracle: math/big exponentiation; overflow exactly when the exact power is outside int64.
func TestPowerBySquaring_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	lo, hi := big.NewInt(math.MinInt64), big.NewInt(math.MaxInt64)
	for range 600 {
		base := rng.Int64N(41) - 20
		if rng.IntN(4) == 0 {
			base = rng.Int64() >> rng.IntN(64)
		}
		exp := rng.Int64N(70)
		want := new(big.Int).Exp(big.NewInt(base), big.NewInt(exp), nil)
		got, err := PowerBySquaring(base, exp)
		if want.Cmp(lo) < 0 || want.Cmp(hi) > 0 {
			if !errors.Is(err, ErrOverflow) {
				t.Fatalf("PowerBySquaring(%d, %d) = (%d, %v), want ErrOverflow", base, exp, got, err)
			}
			continue
		}
		if err != nil || got != want.Int64() {
			t.Fatalf("PowerBySquaring(%d, %d) = (%d, %v), want %s", base, exp, got, err, want)
		}
	}
}

func BenchmarkPowerBySquaring(b *testing.B) {
	b.ReportAllocs()
	var s int64
	for i := int64(0); b.Loop(); i++ {
		v, _ := PowerBySquaring(3, 39-i&7)
		s += v
	}
	benchSink = uint64(s)
}

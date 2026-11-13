package intmath

import (
	"math"
	"math/big"
	"testing"
)

// Concrete Mathematics §4.4: ε_5(100!) = 20 + 4 = 24, so 100! ends in 24 zeros.
func TestFactorialTrailingZeros_sourceExample(t *testing.T) {
	if got := FactorialTrailingZeros(100); got != 24 {
		t.Errorf("FactorialTrailingZeros(100) = %d, want 24", got)
	}
}

func TestFactorialTrailingZeros_cases(t *testing.T) {
	tests := []struct{ n, want uint64 }{
		{0, 0}, {1, 0}, {4, 0}, {5, 1}, {10, 2}, {24, 4}, {25, 6}, {125, 31}, {1000, 249},
		{math.MaxUint64, 4_611_686_018_427_387_890},
	}
	for _, tt := range tests {
		if got := FactorialTrailingZeros(tt.n); got != tt.want {
			t.Errorf("FactorialTrailingZeros(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}

// Oracle: count the zeros at the end of n! computed with math/big.
func TestFactorialTrailingZeros_oracle(t *testing.T) {
	f := big.NewInt(1)
	for n := uint64(1); n <= 600; n++ {
		f.Mul(f, new(big.Int).SetUint64(n))
		s := f.String()
		zeros := uint64(0)
		for i := len(s) - 1; i >= 0 && s[i] == '0'; i-- {
			zeros++
		}
		if got := FactorialTrailingZeros(n); got != zeros {
			t.Fatalf("FactorialTrailingZeros(%d) = %d, want %d", n, got, zeros)
		}
	}
}

func BenchmarkFactorialTrailingZeros(b *testing.B) {
	b.ReportAllocs()
	var s uint64
	for i := uint64(0); b.Loop(); i++ {
		s += FactorialTrailingZeros(math.MaxUint64 - i)
	}
	benchSink = s
}

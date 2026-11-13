package intmath

import (
	"math"
	"math/rand/v2"
	"testing"
)

// Gauss, Disquisitiones art. 12: a number and its digit sum agree modulo 9. 38 → 3+8 = 11 → 1+1 = 2.
func TestDigitalRoot_cases(t *testing.T) {
	tests := []struct{ n, want uint64 }{
		{38, 2}, {0, 0}, {1, 1}, {9, 9}, {10, 1}, {18, 9}, {99, 9}, {12345, 6},
		{math.MaxUint64, 6}, // 18446744073709551615: digit sum 87 → 15 → 6
	}
	for _, tt := range tests {
		if got := DigitalRoot(tt.n); got != tt.want {
			t.Errorf("DigitalRoot(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}

// Oracle: sum the digits repeatedly.
func TestDigitalRoot_oracle(t *testing.T) {
	slow := func(n uint64) uint64 {
		for n >= 10 {
			var s uint64
			for ; n > 0; n /= 10 {
				s += n % 10
			}
			n = s
		}
		return n
	}
	rng := rand.New(rand.NewPCG(1, 2))
	for range 600 {
		n := rng.Uint64() >> rng.IntN(64)
		if got := DigitalRoot(n); got != slow(n) {
			t.Fatalf("DigitalRoot(%d) = %d, want %d", n, got, slow(n))
		}
	}
}

func BenchmarkDigitalRoot(b *testing.B) {
	b.ReportAllocs()
	var s uint64
	for i := uint64(0); b.Loop(); i++ {
		s += DigitalRoot(i)
	}
	benchSink = s
}

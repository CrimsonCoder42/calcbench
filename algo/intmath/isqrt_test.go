package intmath

import (
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// Cohen §1.7.1: the output x satisfies x² ≤ n < (x+1)².
func TestISqrt_cases(t *testing.T) {
	tests := []struct{ n, want uint64 }{
		{0, 0}, {1, 1}, {2, 1}, {3, 1}, {4, 2}, {8, 2}, {9, 3}, {99, 9}, {100, 10},
		{1<<32 - 1, 65535}, {1 << 32, 65536}, {1<<62 - 1, 1<<31 - 1}, {1 << 62, 1 << 31},
		{(1<<32 - 1) * (1<<32 - 1), 1<<32 - 1}, {(1<<32-1)*(1<<32-1) - 1, 1<<32 - 2},
		{math.MaxUint64, 1<<32 - 1},
	}
	for _, tt := range tests {
		if got := ISqrt(tt.n); got != tt.want {
			t.Errorf("ISqrt(%d) = %d, want %d", tt.n, got, tt.want)
		}
	}
}

// Oracle: math/big's exact integer square root, and squares and their neighbours of every size.
func TestISqrt_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 600 {
		n := rng.Uint64() >> rng.IntN(64)
		want := new(big.Int).Sqrt(new(big.Int).SetUint64(n)).Uint64()
		if got := ISqrt(n); got != want {
			t.Fatalf("ISqrt(%d) = %d, want %d", n, got, want)
		}
	}
	for range 300 {
		r := rng.Uint64() >> 32
		sq := r * r
		if ISqrt(sq) != r || (sq > 0 && ISqrt(sq-1) != r-1) || (r < 1<<32-1 && ISqrt(sq+2*r) != r) {
			t.Fatalf("ISqrt wrong near %d²", r)
		}
	}
}

func TestISqrt_noAllocation(t *testing.T) {
	if a := testing.AllocsPerRun(20, func() { benchSink = ISqrt(math.MaxUint64) }); a != 0 {
		t.Errorf("ISqrt allocates %v times, want 0", a)
	}
}

func BenchmarkISqrt(b *testing.B) {
	b.ReportAllocs()
	var s uint64
	for i := uint64(1); b.Loop(); i++ {
		s += ISqrt(i * 0x9E3779B97F4A7C15)
	}
	benchSink = s
}

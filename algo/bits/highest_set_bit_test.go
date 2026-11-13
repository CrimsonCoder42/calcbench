package bits

import (
	"math"
	"math/rand/v2"
	"testing"
)

// Hacker's Delight §3-2: flp2 is the greatest power of 2 ≤ x, and 0 for x = 0.
func TestHighestSetBit_cases(t *testing.T) {
	tests := []struct{ x, want uint64 }{
		{0b10110, 0b10000}, {0, 0}, {1, 1}, {2, 2}, {3, 2}, {1000, 512},
		{math.MaxUint64, 1 << 63}, {1 << 63, 1 << 63},
	}
	for _, tt := range tests {
		if got := HighestSetBit(tt.x); got != tt.want {
			t.Errorf("HighestSetBit(%d) = %d, want %d", tt.x, got, tt.want)
		}
	}
}

func TestHighestSetBit_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 500 {
		x := rng.Uint64() >> rng.IntN(64)
		var want uint64
		for p := uint64(1); p != 0 && p <= x; p <<= 1 {
			want = p
		}
		if got := HighestSetBit(x); got != want {
			t.Fatalf("HighestSetBit(%d) = %d, want %d", x, got, want)
		}
	}
}

func BenchmarkHighestSetBit(b *testing.B) {
	b.ReportAllocs()
	var s uint64
	for i := uint64(0); b.Loop(); i++ {
		s += HighestSetBit(i)
	}
	benchSink = s
}

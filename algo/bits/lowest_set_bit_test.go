package bits

import (
	"math"
	mathbits "math/bits"
	"math/rand/v2"
	"testing"
)

// Knuth Vol. 4A §7.1.3 and Hacker's Delight §2-1: x & −x keeps only the rightmost 1, e.g. 0101 1000 → 0000 1000.
func TestLowestSetBit_cases(t *testing.T) {
	tests := []struct{ x, want uint64 }{
		{0b0101_1000, 0b0000_1000}, {0b10110, 0b00010}, {0, 0}, {1, 1}, {1 << 63, 1 << 63},
		{math.MaxUint64, 1}, {6, 2},
	}
	for _, tt := range tests {
		if got := LowestSetBit(tt.x); got != tt.want {
			t.Errorf("LowestSetBit(%#b) = %#b, want %#b", tt.x, got, tt.want)
		}
	}
}

func TestLowestSetBit_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 500 {
		x := rng.Uint64() << rng.IntN(64)
		var want uint64
		if x != 0 {
			want = 1 << mathbits.TrailingZeros64(x)
		}
		if got := LowestSetBit(x); got != want {
			t.Fatalf("LowestSetBit(%#x) = %#x, want %#x", x, got, want)
		}
	}
}

func BenchmarkLowestSetBit(b *testing.B) {
	b.ReportAllocs()
	var s uint64
	for i := uint64(0); b.Loop(); i++ {
		s += LowestSetBit(i)
	}
	benchSink = s
}

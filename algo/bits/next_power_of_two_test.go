package bits

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// Hacker's Delight §3-2: clp2 rounds up to the next power of two; clp2(0) is taken as 1 here.
func TestNextPowerOfTwo_cases(t *testing.T) {
	tests := []struct{ x, want uint64 }{
		{0, 1}, {1, 1}, {2, 2}, {3, 4}, {5, 8}, {8, 8}, {9, 16}, {1000, 1024},
		{1<<32 - 1, 1 << 32}, {1<<62 + 1, 1 << 63}, {1 << 63, 1 << 63},
	}
	for _, tt := range tests {
		got, err := NextPowerOfTwo(tt.x)
		if err != nil || got != tt.want {
			t.Errorf("NextPowerOfTwo(%d) = (%d, %v), want (%d, nil)", tt.x, got, err, tt.want)
		}
	}
}

func TestNextPowerOfTwo_errors(t *testing.T) {
	for _, x := range []uint64{1<<63 + 1, math.MaxUint64} {
		if _, err := NextPowerOfTwo(x); !errors.Is(err, ErrOverflow) {
			t.Errorf("NextPowerOfTwo(%d) error = %v, want ErrOverflow", x, err)
		}
	}
}

func TestNextPowerOfTwo_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 500 {
		x := rng.Uint64() >> (1 + rng.IntN(63))
		want := uint64(1)
		for want < x {
			want <<= 1
		}
		if got, err := NextPowerOfTwo(x); err != nil || got != want {
			t.Fatalf("NextPowerOfTwo(%d) = (%d, %v), want %d", x, got, err, want)
		}
	}
}

func BenchmarkNextPowerOfTwo(b *testing.B) {
	b.ReportAllocs()
	var s uint64
	for i := uint64(0); b.Loop(); i++ {
		v, _ := NextPowerOfTwo(i)
		s += v
	}
	benchSink = s
}

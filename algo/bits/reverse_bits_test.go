package bits

import (
	"errors"
	"math"
	mathbits "math/bits"
	"math/rand/v2"
	"testing"
)

// slowReverse reverses the low width bits one at a time.
func slowReverse(x uint64, width int) uint64 {
	var r uint64
	for i := range width {
		r |= (x >> i & 1) << (width - 1 - i)
	}
	return r
}

// Hacker's Delight §7-1: reversing a field of width w is the full reversal shifted right by 64 − w.
func TestReverseBits_cases(t *testing.T) {
	tests := []struct {
		x     uint64
		width int
		want  uint64
	}{
		{0b0011, 4, 0b1100}, {0b1, 1, 0b1}, {0, 8, 0}, {0b0000_0001, 8, 0b1000_0000},
		{0b1011, 4, 0b1101}, {0b110, 3, 0b011}, {1, 64, 1 << 63}, {math.MaxUint64, 64, math.MaxUint64},
		{0x0000_0000_0000_00F0, 8, 0x0F},
	}
	for _, tt := range tests {
		got, err := ReverseBits(tt.x, tt.width)
		if err != nil || got != tt.want {
			t.Errorf("ReverseBits(%#b, %d) = (%#b, %v), want %#b", tt.x, tt.width, got, err, tt.want)
		}
	}
}

func TestReverseBits_errors(t *testing.T) {
	tests := []struct {
		x       uint64
		width   int
		wantErr error
	}{
		{1, 0, ErrWidth}, {1, -1, ErrWidth}, {1, 65, ErrWidth}, {1, math.MinInt, ErrWidth},
		{0b10000, 4, ErrOutOfRange}, {1 << 63, 63, ErrOutOfRange},
	}
	for _, tt := range tests {
		if _, err := ReverseBits(tt.x, tt.width); !errors.Is(err, tt.wantErr) {
			t.Errorf("ReverseBits(%#b, %d) error = %v, want %v", tt.x, tt.width, err, tt.wantErr)
		}
	}
}

func TestReverseBits_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 500 {
		w := 1 + rng.IntN(64)
		x := rng.Uint64()
		if w < 64 {
			x &= 1<<w - 1
		}
		got, err := ReverseBits(x, w)
		if want := slowReverse(x, w); err != nil || got != want {
			t.Fatalf("ReverseBits(%#x, %d) = (%#x, %v), want %#x", x, w, got, err, want)
		}
		if back, _ := ReverseBits(got, w); back != x {
			t.Fatalf("reversing twice gave %#x, want %#x", back, x)
		}
	}
	if got, _ := ReverseBits(0xDEADBEEF, 64); got != mathbits.Reverse64(0xDEADBEEF) {
		t.Error("width 64 differs from math/bits.Reverse64")
	}
}

func BenchmarkReverseBits(b *testing.B) {
	b.ReportAllocs()
	var s uint64
	for i := uint64(0); b.Loop(); i++ {
		v, _ := ReverseBits(i&0xFFFF, 16)
		s += v
	}
	benchSink = s
}

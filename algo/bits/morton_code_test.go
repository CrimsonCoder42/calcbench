package bits

import (
	"math"
	"math/rand/v2"
	"testing"
)

// slowMorton interleaves one bit at a time: x to even positions, y to odd.
func slowMorton(x, y uint32) uint64 {
	var code uint64
	for i := range 32 {
		code |= uint64(x>>i&1) << (2 * i)
		code |= uint64(y>>i&1) << (2*i + 1)
	}
	return code
}

// Morton 1966: the Z-order of the 2×2 cells (0,0), (1,0), (0,1), (1,1) is 0, 1, 2, 3; of (2,3) is 14.
func TestMortonEncode2D_cases(t *testing.T) {
	tests := []struct {
		x, y uint32
		want uint64
	}{
		{0, 0, 0}, {1, 0, 1}, {0, 1, 2}, {1, 1, 3}, {2, 3, 14}, {3, 5, 0b100111},
		{math.MaxUint32, 0, 0x5555_5555_5555_5555}, {0, math.MaxUint32, 0xAAAA_AAAA_AAAA_AAAA},
		{math.MaxUint32, math.MaxUint32, math.MaxUint64},
	}
	for _, tt := range tests {
		if got := MortonEncode2D(tt.x, tt.y); got != tt.want {
			t.Errorf("MortonEncode2D(%d, %d) = %#b, want %#b", tt.x, tt.y, got, tt.want)
		}
		if x, y := MortonDecode2D(tt.want); x != tt.x || y != tt.y {
			t.Errorf("MortonDecode2D(%#b) = (%d, %d), want (%d, %d)", tt.want, x, y, tt.x, tt.y)
		}
	}
}

func TestMortonEncode2D_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 500 {
		x, y := rng.Uint32(), rng.Uint32()
		code := MortonEncode2D(x, y)
		if want := slowMorton(x, y); code != want {
			t.Fatalf("MortonEncode2D(%d, %d) = %#x, want %#x", x, y, code, want)
		}
		if gx, gy := MortonDecode2D(code); gx != x || gy != y {
			t.Fatalf("round trip of (%d, %d) gave (%d, %d)", x, y, gx, gy)
		}
	}
}

func TestMortonDecode2D_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 500 {
		code := rng.Uint64()
		x, y := MortonDecode2D(code)
		if MortonEncode2D(x, y) != code {
			t.Fatalf("MortonEncode2D(MortonDecode2D(%#x)) differs", code)
		}
	}
}

func BenchmarkMortonEncode2D(b *testing.B) {
	b.ReportAllocs()
	var s uint64
	for i := uint32(0); b.Loop(); i++ {
		s += MortonEncode2D(i, ^i)
	}
	benchSink = s
}

func BenchmarkMortonDecode2D(b *testing.B) {
	b.ReportAllocs()
	var s uint64
	for i := uint64(0); b.Loop(); i++ {
		x, y := MortonDecode2D(i * 0x9E3779B97F4A7C15)
		s += uint64(x) + uint64(y)
	}
	benchSink = s
}

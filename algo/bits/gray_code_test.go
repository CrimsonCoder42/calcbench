package bits

import (
	"math"
	mathbits "math/bits"
	"math/rand/v2"
	"testing"
)

// Knuth Vol. 4A §7.2.1.1: the 3-bit Gray binary sequence is 000, 001, 011, 010, 110, 111, 101, 100.
func TestGrayEncode_sourceExample(t *testing.T) {
	want := []uint64{0b000, 0b001, 0b011, 0b010, 0b110, 0b111, 0b101, 0b100}
	for k, w := range want {
		if got := GrayEncode(uint64(k)); got != w {
			t.Errorf("GrayEncode(%d) = %03b, want %03b", k, got, w)
		}
	}
}

func TestGrayEncode_neighboursDifferInOneBit(t *testing.T) {
	for k := uint64(0); k < 4096; k++ {
		if d := mathbits.OnesCount64(GrayEncode(k) ^ GrayEncode(k+1)); d != 1 {
			t.Fatalf("codes of %d and %d differ in %d bits, want 1", k, k+1, d)
		}
	}
	// Reflected: the code of the largest value differs from the code of 0 in one bit too.
	if d := mathbits.OnesCount64(GrayEncode(math.MaxUint64) ^ GrayEncode(0)); d != 1 {
		t.Errorf("cyclic wrap differs in %d bits, want 1", d)
	}
}

func TestGrayDecode_sourceExample(t *testing.T) {
	codes := []uint64{0b000, 0b001, 0b011, 0b010, 0b110, 0b111, 0b101, 0b100}
	for k, g := range codes {
		if got := GrayDecode(g); got != uint64(k) {
			t.Errorf("GrayDecode(%03b) = %d, want %d", g, got, k)
		}
	}
}

// Oracle: decoding bit by bit from the top, b_i = g_i ⊕ b_{i+1}.
func TestGrayDecode_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for _, x := range []uint64{0, 1, math.MaxUint64, 1 << 63} {
		if got := GrayDecode(GrayEncode(x)); got != x {
			t.Errorf("GrayDecode(GrayEncode(%#x)) = %#x", x, got)
		}
	}
	for range 500 {
		g := rng.Uint64()
		var want, prev uint64
		for i := 63; i >= 0; i-- {
			prev ^= g >> i & 1
			want |= prev << i
		}
		if got := GrayDecode(g); got != want {
			t.Fatalf("GrayDecode(%#x) = %#x, want %#x", g, got, want)
		}
		if GrayEncode(GrayDecode(g)) != g {
			t.Fatalf("GrayEncode(GrayDecode(%#x)) != input", g)
		}
	}
}

func BenchmarkGrayEncode(b *testing.B) {
	b.ReportAllocs()
	var s uint64
	for i := uint64(0); b.Loop(); i++ {
		s += GrayEncode(i)
	}
	benchSink = s
}

func BenchmarkGrayDecode(b *testing.B) {
	b.ReportAllocs()
	var s uint64
	for i := uint64(0); b.Loop(); i++ {
		s += GrayDecode(i)
	}
	benchSink = s
}

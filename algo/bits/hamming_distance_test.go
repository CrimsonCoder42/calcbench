package bits

import (
	"math"
	"math/rand/v2"
	"testing"
)

func TestHammingDistance_cases(t *testing.T) {
	tests := []struct {
		x, y uint64
		want int
	}{
		{0b1011101, 0b1001001, 2}, // the two words differ in bits 2 and 4
		{0, 0, 0}, {7, 7, 0}, {1, 4, 2}, {0, math.MaxUint64, 64}, {1 << 63, 0, 1},
	}
	for _, tt := range tests {
		if got := HammingDistance(tt.x, tt.y); got != tt.want {
			t.Errorf("HammingDistance(%#b, %#b) = %d, want %d", tt.x, tt.y, got, tt.want)
		}
	}
}

func TestHammingDistance_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 500 {
		x, y := rng.Uint64(), rng.Uint64()
		want := 0
		for i := range 64 {
			if (x>>i)&1 != (y>>i)&1 {
				want++
			}
		}
		if got := HammingDistance(x, y); got != want || HammingDistance(y, x) != want {
			t.Fatalf("HammingDistance(%#x, %#x) = %d, want %d (symmetric)", x, y, got, want)
		}
	}
}

func BenchmarkHammingDistance(b *testing.B) {
	b.ReportAllocs()
	s := 0
	for i := uint64(0); b.Loop(); i++ {
		s += HammingDistance(i, i*0x9E3779B97F4A7C15)
	}
	benchSink = uint64(s)
}

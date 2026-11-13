package bits

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

func TestBitwiseAndRange_cases(t *testing.T) {
	tests := []struct{ lo, hi, want uint64 }{
		{5, 7, 4},    // 101 & 110 & 111
		{12, 15, 12}, // 1100 … 1111
		{0, 0, 0},    // one value
		{9, 9, 9},    // one value
		{0, 1, 0},    // includes zero
		{1, math.MaxUint64, 0},
		{math.MaxUint64, math.MaxUint64, math.MaxUint64},
		{1 << 63, math.MaxUint64, 1 << 63},
		{8, 16, 0}, // crosses a power of two
	}
	for _, tt := range tests {
		got, err := BitwiseAndRange(tt.lo, tt.hi)
		if err != nil || got != tt.want {
			t.Errorf("BitwiseAndRange(%d, %d) = (%d, %v), want %d", tt.lo, tt.hi, got, err, tt.want)
		}
	}
}

func TestBitwiseAndRange_errors(t *testing.T) {
	if _, err := BitwiseAndRange(3, 2); !errors.Is(err, ErrRange) {
		t.Errorf("BitwiseAndRange(3, 2) error = %v, want ErrRange", err)
	}
}

// Oracle: AND every value in short ranges.
func TestBitwiseAndRange_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 500 {
		lo := rng.Uint64() >> rng.IntN(64)
		span := uint64(rng.IntN(300))
		if lo > math.MaxUint64-span {
			continue
		}
		hi := lo + span
		want := lo
		for v := lo; v <= hi; v++ {
			want &= v
			if v == math.MaxUint64 {
				break
			}
		}
		if got, err := BitwiseAndRange(lo, hi); err != nil || got != want {
			t.Fatalf("BitwiseAndRange(%d, %d) = (%d, %v), want %d", lo, hi, got, err, want)
		}
	}
}

func BenchmarkBitwiseAndRange(b *testing.B) {
	b.ReportAllocs()
	var s uint64
	for i := uint64(0); b.Loop(); i++ {
		v, _ := BitwiseAndRange(i, i+i&0xFF)
		s += v
	}
	benchSink = s
}

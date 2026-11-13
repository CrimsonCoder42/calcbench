package intmath

import (
	"math"
	"testing"
)

// Guy, Unsolved Problems in Number Theory, E34: the happy numbers up to 50 are 1, 7, 10, 13, 19, 23, 28,
// 31, 32, 44, 49; every other number reaches the cycle 4 → 16 → 37 → 58 → 89 → 145 → 42 → 20 → 4.
func TestIsHappyNumber_sourceExample(t *testing.T) {
	happy := map[uint64]bool{1: true, 7: true, 10: true, 13: true, 19: true, 23: true, 28: true, 31: true, 32: true, 44: true, 49: true}
	for n := uint64(0); n <= 50; n++ {
		if got := IsHappyNumber(n); got != happy[n] {
			t.Errorf("IsHappyNumber(%d) = %v, want %v", n, got, happy[n])
		}
	}
	for _, n := range []uint64{4, 16, 37, 58, 89, 145, 42, 20} {
		if IsHappyNumber(n) {
			t.Errorf("IsHappyNumber(%d) = true; it is on the unhappy cycle", n)
		}
	}
}

// Oracle: iterate with a visited set until 1 or a repeat.
func TestIsHappyNumber_oracle(t *testing.T) {
	slow := func(n uint64) bool {
		seen := map[uint64]bool{}
		for n != 1 && !seen[n] {
			seen[n] = true
			n = digitSquareSum(n)
		}
		return n == 1
	}
	for n := uint64(0); n < 20_000; n++ {
		if got := IsHappyNumber(n); got != slow(n) {
			t.Fatalf("IsHappyNumber(%d) = %v, want %v", n, got, slow(n))
		}
	}
	for _, n := range []uint64{math.MaxUint64, 1 << 63, 10_000_000_000_000_000_000} {
		if got := IsHappyNumber(n); got != slow(n) {
			t.Errorf("IsHappyNumber(%d) = %v, want %v", n, got, slow(n))
		}
	}
}

func BenchmarkIsHappyNumber(b *testing.B) {
	b.ReportAllocs()
	n := 0
	for i := uint64(0); b.Loop(); i++ {
		if IsHappyNumber(i) {
			n++
		}
	}
	benchSink = uint64(n)
}

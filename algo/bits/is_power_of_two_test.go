package bits

import (
	"math"
	"math/rand/v2"
	"testing"
)

// benchSink keeps benchmark results alive. Shared by every benchmark in this package.
var benchSink uint64

// benchSizes are the collection sizes of .claude/rules/go-complexity.md. Shared by the benchmarks.
var benchSizes = []struct {
	name string
	n    int
}{{"1e3", 1_000}, {"1e4", 10_000}, {"1e5", 100_000}}

// Hacker's Delight §2-1: x & (x−1) is zero exactly for 0 and the powers of two.
func TestIsPowerOfTwo_sourceExample(t *testing.T) {
	for k := range 64 {
		if !IsPowerOfTwo(uint64(1) << k) {
			t.Errorf("IsPowerOfTwo(2^%d) = false, want true", k)
		}
	}
}

func TestIsPowerOfTwo_cases(t *testing.T) {
	if IsPowerOfTwo(0) || IsPowerOfTwo(-4) || IsPowerOfTwo(6) || IsPowerOfTwo(uint64(math.MaxUint64)) {
		t.Error("IsPowerOfTwo accepted 0, −4, 6 or MaxUint64")
	}
	if IsPowerOfTwo(int64(math.MinInt64)) || IsPowerOfTwo(int8(math.MinInt8)) {
		t.Error("IsPowerOfTwo accepted the minimum of a signed type")
	}
	if !IsPowerOfTwo(int8(64)) || !IsPowerOfTwo(1) || !IsPowerOfTwo(int64(1)<<62) || !IsPowerOfTwo(uint8(128)) {
		t.Error("IsPowerOfTwo rejected a power of two")
	}
}

func TestIsPowerOfTwo_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 500 {
		x := rng.Uint64() >> rng.IntN(64)
		want := false
		for k := range 64 {
			if x == uint64(1)<<k {
				want = true
			}
		}
		if got := IsPowerOfTwo(x); got != want {
			t.Fatalf("IsPowerOfTwo(%d) = %v, want %v", x, got, want)
		}
	}
}

func TestIsPowerOfFour_cases(t *testing.T) {
	for _, x := range []uint64{1, 4, 16, 64, 1 << 62} {
		if !IsPowerOfFour(x) {
			t.Errorf("IsPowerOfFour(%d) = false, want true", x)
		}
	}
	for _, x := range []uint64{0, 2, 8, 12, 32, 1 << 63, math.MaxUint64} {
		if IsPowerOfFour(x) {
			t.Errorf("IsPowerOfFour(%d) = true, want false", x)
		}
	}
	if IsPowerOfFour(-4) || !IsPowerOfFour(int32(1<<30)) || IsPowerOfFour(int32(math.MinInt32)) {
		t.Error("IsPowerOfFour wrong on signed input")
	}
}

func TestIsPowerOfFour_oracle(t *testing.T) {
	for x := uint64(0); x < 5000; x++ {
		want := false
		for p := uint64(1); p <= x; p *= 4 {
			if p == x {
				want = true
			}
		}
		if got := IsPowerOfFour(x); got != want {
			t.Fatalf("IsPowerOfFour(%d) = %v, want %v", x, got, want)
		}
	}
}

func BenchmarkIsPowerOfTwo(b *testing.B) {
	b.ReportAllocs()
	n := 0
	for i := uint64(0); b.Loop(); i++ {
		if IsPowerOfTwo(i) {
			n++
		}
	}
	benchSink = uint64(n)
}

func BenchmarkIsPowerOfFour(b *testing.B) {
	b.ReportAllocs()
	n := 0
	for i := uint64(0); b.Loop(); i++ {
		if IsPowerOfFour(i) {
			n++
		}
	}
	benchSink = uint64(n)
}

package intmath

import (
	"math"
	"math/rand/v2"
	"testing"
)

func TestIsPowerOfThree_cases(t *testing.T) {
	for _, x := range []uint64{1, 3, 9, 27, 81, 3_486_784_401, maxPowerOfThreeUint64} {
		if !IsPowerOfThree(x) {
			t.Errorf("IsPowerOfThree(%d) = false, want true", x)
		}
	}
	for _, x := range []uint64{0, 2, 6, 45, 3_486_784_400, maxPowerOfThreeUint64 - 1, math.MaxUint64} {
		if IsPowerOfThree(x) {
			t.Errorf("IsPowerOfThree(%d) = true, want false", x)
		}
	}
	if IsPowerOfThree(-3) || IsPowerOfThree(int64(math.MinInt64)) || !IsPowerOfThree(int8(81)) || !IsPowerOfThree(int64(4_052_555_153_018_976_267)) {
		t.Error("IsPowerOfThree wrong on signed input")
	}
}

// Oracle: divide by 3 while possible.
func TestIsPowerOfThree_oracle(t *testing.T) {
	slow := func(x uint64) bool {
		for x > 1 && x%3 == 0 {
			x /= 3
		}
		return x == 1
	}
	for x := uint64(0); x < 10_000; x++ {
		if IsPowerOfThree(x) != slow(x) {
			t.Fatalf("IsPowerOfThree(%d) wrong", x)
		}
	}
	rng := rand.New(rand.NewPCG(1, 2))
	for range 500 {
		x := rng.Uint64() >> rng.IntN(64)
		if IsPowerOfThree(x) != slow(x) {
			t.Fatalf("IsPowerOfThree(%d) wrong", x)
		}
	}
}

func BenchmarkIsPowerOfThree(b *testing.B) {
	b.ReportAllocs()
	n := 0
	for i := uint64(1); b.Loop(); i++ {
		if IsPowerOfThree(i) {
			n++
		}
	}
	benchSink = uint64(n)
}

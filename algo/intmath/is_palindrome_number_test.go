package intmath

import (
	"math"
	"math/rand/v2"
	"strconv"
	"testing"
)

func TestIsPalindromeNumber_cases(t *testing.T) {
	yes := []int64{0, 7, 11, 121, 1221, 12321, 1_000_000_001, 900_000_000_000_000_009}
	no := []int64{-121, -1, 10, 12, 100, 1231, math.MaxInt64, math.MinInt64}
	for _, x := range yes {
		if !IsPalindromeNumber(x) {
			t.Errorf("IsPalindromeNumber(%d) = false, want true", x)
		}
	}
	for _, x := range no {
		if IsPalindromeNumber(x) {
			t.Errorf("IsPalindromeNumber(%d) = true, want false", x)
		}
	}
	if !IsPalindromeNumber(uint64(18_446_744_066_044_764_481)) || IsPalindromeNumber(uint64(math.MaxUint64)) {
		t.Error("IsPalindromeNumber wrong on 20-digit uint64")
	}
}

// Oracle: compare the decimal string with its reversal.
func TestIsPalindromeNumber_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	check := func(x uint64) {
		s := strconv.FormatUint(x, 10)
		want := true
		for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
			if s[i] != s[j] {
				want = false
			}
		}
		if got := IsPalindromeNumber(x); got != want {
			t.Fatalf("IsPalindromeNumber(%d) = %v, want %v", x, got, want)
		}
	}
	for x := uint64(0); x < 20_000; x++ {
		check(x)
	}
	for range 500 {
		// Build palindromes and near-palindromes of random length.
		half := strconv.FormatUint(rng.Uint64()>>rng.IntN(64), 10)
		rs := []byte(half)
		for i, j := 0, len(rs)-1; i < j; i, j = i+1, j-1 {
			rs[i], rs[j] = rs[j], rs[i]
		}
		if v, err := strconv.ParseUint(half+string(rs), 10, 64); err == nil {
			check(v)
			check(v + 1)
		}
		check(rng.Uint64())
	}
}

func BenchmarkIsPalindromeNumber(b *testing.B) {
	b.ReportAllocs()
	n := 0
	for i := uint64(0); b.Loop(); i++ {
		if IsPalindromeNumber(123_456_787_654_321 + i&1) {
			n++
		}
	}
	benchSink = uint64(n)
}

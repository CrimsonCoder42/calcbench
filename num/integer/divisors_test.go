package integer

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"strconv"
	"testing"
)

// divisorsTrial is the oracle: trial division up to √n.
func divisorsTrial(n uint64) []uint64 {
	var lo, hi []uint64
	for d := uint64(1); d*d <= n; d++ {
		if n%d == 0 {
			lo = append(lo, d)
			if d != n/d {
				hi = append(hi, n/d)
			}
		}
	}
	slices.Reverse(hi)
	return append(lo, hi...)
}

// Rosen §7.2: n = ∏ pᵢ^aᵢ has τ(n) = ∏(aᵢ + 1) divisors. Hand: 120 = 2³·3·5 → 16 divisors listed below.
func TestDivisors_cases(t *testing.T) {
	tests := []struct {
		n    uint64
		want []uint64
	}{
		{1, []uint64{1}},
		{2, []uint64{1, 2}},
		{12, []uint64{1, 2, 3, 4, 6, 12}},
		{36, []uint64{1, 2, 3, 4, 6, 9, 12, 18, 36}},
		{97, []uint64{1, 97}},
		{120, []uint64{1, 2, 3, 4, 5, 6, 8, 10, 12, 15, 20, 24, 30, 40, 60, 120}},
		{1 << 10, []uint64{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024}},
	}
	for _, tc := range tests {
		if got, err := Divisors(tc.n); err != nil || !slices.Equal(got, tc.want) {
			t.Errorf("Divisors(%d) = %v, %v; want %v", tc.n, got, err, tc.want)
		}
	}
}

// 2⁶⁴ − 1 = 3·5·17·257·641·65537·6700417 (seven distinct primes) has 2⁷ = 128 divisors, the largest itself.
func TestDivisors_maxUint64(t *testing.T) {
	got, err := Divisors(math.MaxUint64)
	if err != nil || len(got) != 128 || got[0] != 1 || got[127] != math.MaxUint64 {
		t.Fatalf("Divisors(2^64-1): len %d, err %v", len(got), err)
	}
	for _, d := range got {
		if math.MaxUint64%d != 0 {
			t.Fatalf("%d does not divide 2^64-1", d)
		}
	}
}

func TestDivisors_errors(t *testing.T) {
	if _, err := Divisors(0); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("Divisors(0): %v", err)
	}
}

func TestDivisors_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		n := r.Uint64N(1<<20) + 1
		got, err := Divisors(n)
		if want := divisorsTrial(n); err != nil || !slices.Equal(got, want) {
			t.Fatalf("Divisors(%d) = %v, %v; want %v", n, got, err, want)
		}
	}
}

func BenchmarkDivisors(b *testing.B) {
	for _, n := range []uint64{720720, 963761198400, math.MaxUint64} {
		b.Run("n="+strconv.FormatUint(n, 10), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = Divisors(n)
			}
		})
	}
}

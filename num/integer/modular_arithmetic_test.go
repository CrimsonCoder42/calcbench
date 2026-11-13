package integer

import (
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// bigOf returns n as a *big.Int, for oracles that must not overflow.
func bigOf(n uint64) *big.Int { return new(big.Int).SetUint64(n) }

// moduli are the n used by the modular oracles: small, word-sized and the extremes of uint64.
func moduli(rng *rand.Rand) []uint64 {
	ns := []uint64{2, 3, 1031, 4294967291, math.MaxUint64, math.MaxUint64 - 58, 1 << 63}
	for range 100 {
		ns = append(ns, 2+rng.Uint64N(math.MaxUint64-1))
	}
	return ns
}

func TestModularArithmetic_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for _, n := range moduli(rng) {
		bn := bigOf(n)
		for range 4 {
			a, b, e := rng.Uint64N(n), rng.Uint64N(n), rng.Uint64()
			want := new(big.Int).Mul(bigOf(a), bigOf(b))
			if got := mulMod(a, b, n); got != want.Mod(want, bn).Uint64() {
				t.Fatalf("mulMod(%d, %d, %d) = %d, want %d", a, b, n, got, want)
			}
			want = new(big.Int).Add(bigOf(a), bigOf(b))
			if got := addMod(a, b, n); got != want.Mod(want, bn).Uint64() {
				t.Fatalf("addMod(%d, %d, %d) = %d, want %d", a, b, n, got, want)
			}
			want = new(big.Int).Exp(bigOf(a), bigOf(e), bn)
			if got := powMod(a, e, n); got != want.Uint64() {
				t.Fatalf("powMod(%d, %d, %d) = %d, want %d", a, e, n, got, want)
			}
			want = new(big.Int).GCD(nil, nil, bigOf(a), bigOf(n))
			if got := GCD(a, n); got != want.Uint64() {
				t.Fatalf("gcd(%d, %d) = %d, want %d", a, n, got, want)
			}
		}
	}
}

func TestGcd_cases(t *testing.T) {
	tests := []struct{ a, b, want uint64 }{
		{0, 0, 0},
		{0, 7, 7},
		{7, 0, 7},
		{12, 18, 6},
		{1 << 40, 1 << 20, 1 << 20},
		{math.MaxUint64, math.MaxUint64, math.MaxUint64},
		{math.MaxUint64, 641 * 6700417, 641 * 6700417}, // 2³²+1 divides 2⁶⁴−1
	}
	for _, tt := range tests {
		if got := GCD(tt.a, tt.b); got != tt.want {
			t.Errorf("gcd(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestPowMod_cases(t *testing.T) {
	// Fermat: 2^(p−1) ≡ 1 (mod p) for the prime p = 2⁶⁴ − 59; anything^0 = 1.
	if got := powMod(2, math.MaxUint64-59, math.MaxUint64-58); got != 1 {
		t.Errorf("powMod(2, p-1, p) = %d, want 1", got)
	}
	if got := powMod(12345, 0, 1031); got != 1 {
		t.Errorf("powMod(12345, 0, 1031) = %d, want 1", got)
	}
}

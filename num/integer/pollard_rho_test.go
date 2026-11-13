package integer

import (
	"math/rand/v2"
	"testing"
)

// composites with no prime factor below 1031, the inputs findFactor accepts.
var rhoComposites = []uint64{
	1031 * 1031,
	1031 * 1033,
	999983 * 1000003,
	2097143 * 2097143 * 2097143,
	4294967279 * 4294967291,
	4294967291 * 4294967291,
}

func TestFindFactor_nontrivial(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	inputs := append([]uint64(nil), rhoComposites...)
	for range 50 {
		inputs = append(inputs, randomPrime32(rng)*randomPrime32(rng))
	}
	for _, m := range inputs {
		d := findFactor(m)
		if d <= 1 || d >= m || m%d != 0 {
			t.Errorf("findFactor(%d) = %d, want a divisor strictly between 1 and m", m, d)
		}
	}
}

func TestRhoBrent_returnsDivisor(t *testing.T) {
	m := uint64(4294967279 * 4294967291)
	for c := uint64(1); c <= 8; c++ {
		if d := rhoBrent(m, c); d <= 1 || m%d != 0 {
			t.Errorf("rhoBrent(%d, %d) = %d, want a divisor > 1", m, c, d)
		}
	}
}

// trialFactor is the fallback that guarantees findFactor ends; it returns the least divisor ≥ 1031.
func TestTrialFactor_cases(t *testing.T) {
	tests := []struct{ m, want uint64 }{
		{1031 * 1031, 1031},
		{1033 * 1031, 1031},
		{999983 * 1000003, 999983},
		{1000003, 1000003}, // prime: returns m
	}
	for _, tt := range tests {
		if got := trialFactor(tt.m); got != tt.want {
			t.Errorf("trialFactor(%d) = %d, want %d", tt.m, got, tt.want)
		}
	}
}

func TestAbsDiff(t *testing.T) {
	if absDiff(3, 10) != 7 || absDiff(10, 3) != 7 || absDiff(5, 5) != 0 {
		t.Error("absDiff is wrong")
	}
}

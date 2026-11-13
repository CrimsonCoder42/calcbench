package integer

import (
	"math"
	"math/rand/v2"
	"testing"
)

// gcdEuclid is the oracle: Euclid's algorithm by remainders (Knuth §4.5.2, Algorithm A).
func gcdEuclid(a, b uint64) uint64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// Knuth §4.5.2: gcd(40902, 24140) = 34 (the worked example of Algorithms A and B).
func TestGCD_knuthExample(t *testing.T) {
	if got := GCD(40902, 24140); got != 34 {
		t.Errorf("GCD(40902, 24140) = %d, want 34", got)
	}
}

func TestGCD_cases(t *testing.T) {
	tests := []struct {
		name       string
		a, b, want uint64
	}{
		{"both zero", 0, 0, 0},
		{"zero left", 0, 12, 12},
		{"zero right", 12, 0, 12},
		{"one", 1, 999, 1},
		{"equal", 18, 18, 18},
		{"coprime", 35, 64, 1},
		{"powers of two", 1 << 40, 1 << 20, 1 << 20},
		{"max", math.MaxUint64, math.MaxUint64, math.MaxUint64},
		{"max and 3", math.MaxUint64, 3, 3},
		{"Fibonacci neighbours", 7540113804746346429, 4660046610375530309, 1},
	}
	for _, tc := range tests {
		if got := GCD(tc.a, tc.b); got != tc.want {
			t.Errorf("%s: GCD(%d, %d) = %d, want %d", tc.name, tc.a, tc.b, got, tc.want)
		}
	}
}

func TestGCD_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for range 500 {
		k := r.Uint64N(1<<20) + 1
		a, b := r.Uint64N(1<<40)*k, r.Uint64N(1<<40)*k
		if got, want := GCD(a, b), gcdEuclid(a, b); got != want {
			t.Fatalf("GCD(%d, %d) = %d, want %d", a, b, got, want)
		}
	}
}

func BenchmarkGCD(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = GCD(7540113804746346429, 4660046610375530309)
	}
}

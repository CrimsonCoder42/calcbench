package integer

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// Worked examples, hand-checked (catalog, doc/catalog/algebra.md): 2⁶⁴ − 1 → 3; the product of the two
// largest 32-bit primes, 4294967291 · 4294967279 = 18446743979220271189 → 4294967279.
func TestSmallestPrimeFactor_example(t *testing.T) {
	if 4294967291*4294967279 != 18446743979220271189 {
		t.Fatal("product in the worked example is wrong")
	}
	for _, tt := range []struct{ n, want uint64 }{{math.MaxUint64, 3}, {18446743979220271189, 4294967279}} {
		if got, err := SmallestPrimeFactor(tt.n); err != nil || got != tt.want {
			t.Errorf("SmallestPrimeFactor(%d) = (%d, %v), want (%d, nil)", tt.n, got, err, tt.want)
		}
	}
}

func TestSmallestPrimeFactor_cases(t *testing.T) {
	tests := []struct {
		name    string
		n, want uint64
	}{
		{"two", 2, 2},
		{"three", 3, 3},
		{"four", 4, 2},
		{"2⁶³", 1 << 63, 2},
		{"largest table prime", 1021, 1021},
		{"1031, first prime above the table", 1031, 1031},
		{"1031²", 1031 * 1031, 1031},
		{"1021·1031", 1021 * 1031, 1021},
		{"two six-digit primes", 1000003 * 999983, 999983},
		{"square of the largest 32-bit prime", 4294967291 * 4294967291, 4294967291},
		{"cube of the prime 2097143", 2097143 * 2097143 * 2097143, 2097143},
		{"largest prime below 2⁶⁴", math.MaxUint64 - 58, math.MaxUint64 - 58},
		{"Mersenne prime 2⁶¹−1", 1<<61 - 1, 1<<61 - 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := SmallestPrimeFactor(tt.n); err != nil || got != tt.want {
				t.Errorf("SmallestPrimeFactor(%d) = (%d, %v), want (%d, nil)", tt.n, got, err, tt.want)
			}
		})
	}
}

func TestSmallestPrimeFactor_errors(t *testing.T) {
	for _, n := range []uint64{0, 1} {
		if got, err := SmallestPrimeFactor(n); !errors.Is(err, ErrOutOfRange) || got != 0 {
			t.Errorf("SmallestPrimeFactor(%d) = (%d, %v), want (0, ErrOutOfRange)", n, got, err)
		}
	}
}

// Oracle: trial division below 20000, and the first element of PrimeFactors (itself checked by math/big)
// on random inputs of every bit length and on semiprimes of two random 32-bit primes.
func TestSmallestPrimeFactor_oracle(t *testing.T) {
	for n := uint64(2); n < 20_000; n++ {
		want := n
		for d := uint64(2); d*d <= n; d++ {
			if n%d == 0 {
				want = d
				break
			}
		}
		if got, err := SmallestPrimeFactor(n); err != nil || got != want {
			t.Fatalf("SmallestPrimeFactor(%d) = (%d, %v), want (%d, nil)", n, got, err, want)
		}
	}
	rng := rand.New(rand.NewPCG(1, 2))
	inputs := make([]uint64, 0, 500)
	for range 400 {
		inputs = append(inputs, max(2, rng.Uint64()>>rng.UintN(64)))
	}
	for range 100 {
		inputs = append(inputs, randomPrime32(rng)*randomPrime32(rng))
	}
	for _, n := range inputs {
		factors, err := PrimeFactors(n)
		if err != nil {
			t.Fatalf("PrimeFactors(%d): %v", n, err)
		}
		checkFactorization(t, n, factors)
		if got, err := SmallestPrimeFactor(n); err != nil || got != factors[0] {
			t.Fatalf("SmallestPrimeFactor(%d) = (%d, %v), want (%d, nil)", n, got, err, factors[0])
		}
	}
}

func TestSmallestPrimeFactor_noAllocation(t *testing.T) {
	if allocs := testing.AllocsPerRun(20, func() { benchSink, _ = SmallestPrimeFactor(4294967279 * 4294967291) }); allocs != 0 {
		t.Errorf("SmallestPrimeFactor allocates %v times per call, want 0", allocs)
	}
}

func BenchmarkSmallestPrimeFactor(b *testing.B) {
	for _, size := range factorBenchSizes {
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchSink, _ = SmallestPrimeFactor(size.n)
			}
		})
	}
}

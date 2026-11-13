package integer

import (
	"math"
	"math/rand/v2"
	"testing"
)

// benchSinkBool and benchSink keep benchmark results observable.
var (
	benchSinkBool bool
	benchSink     uint64
)

// Worked example: Pomerance, Selfridge, Wagstaff, "The pseudoprimes to 25·10⁹", Mathematics of Computation,
// 1980, 35(151):1003–1026: 3215031751 = 151·751·28351 is a strong pseudoprime to bases 2, 3, 5 and 7, so a
// test with only those bases calls it prime. Ours must not.
func TestIsPrime_pomeranceSelfridgeWagstaffExample(t *testing.T) {
	if 151*751*28351 != 3215031751 {
		t.Fatal("factorisation of the worked example is wrong")
	}
	for _, a := range []uint64{2, 3, 5, 7} {
		n := uint64(3215031751)
		s := 0
		d := n - 1
		for d%2 == 0 {
			d /= 2
			s++
		}
		if !strongProbablePrime(n, a, d, s) {
			t.Errorf("3215031751 fails the strong test to base %d; the source says it passes", a)
		}
	}
	if IsPrime(3215031751) {
		t.Error("IsPrime(3215031751) = true, want false")
	}
}

func TestIsPrime_cases(t *testing.T) {
	tests := []struct {
		name string
		n    uint64
		want bool
	}{
		{"zero", 0, false},
		{"one", 1, false},
		{"two, the only even prime", 2, true},
		{"three", 3, true},
		{"four", 4, false},
		{"largest base, 37", 37, true},
		{"square of a base, 37²", 37 * 37, false},
		{"41, first prime above the bases", 41, true},
		{"41², least composite with no factor ≤ 37", 41 * 41, false},
		{"1031, first prime above the trial table", 1031, true},
		{"Carmichael 561", 561, false},
		{"Carmichael 1105", 1105, false},
		{"largest 32-bit prime 2³²−5", 4294967291, true},
		{"2³²+1 = 641·6700417", 4294967297, false},
		{"Mersenne prime 2⁶¹−1", 1<<61 - 1, true},
		{"2⁶³", 1 << 63, false},
		{"largest prime below 2⁶⁴, 2⁶⁴−59 (hand-checked)", math.MaxUint64 - 58, true},
		{"2⁶⁴−1, the maximum of the type", math.MaxUint64, false},
		{"square of the largest 32-bit prime", 4294967291 * 4294967291, false},
		{"semiprime of two 32-bit primes", 4294967291 * 4294967279, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPrime(tt.n); got != tt.want {
				t.Errorf("IsPrime(%d) = %v, want %v", tt.n, got, tt.want)
			}
		})
	}
}

// Strong pseudoprimes to the first k prime bases, each composite (the oracle confirms it). They defeat any
// shorter base list, and so test that all 12 bases are used. Numbers re-derived here by script, 2026-11-12.
func TestIsPrime_strongPseudoprimes(t *testing.T) {
	for _, n := range []uint64{
		2047,                // base 2
		1373653,             // bases 2, 3
		25326001,            // bases 2, 3, 5
		3215031751,          // bases 2, 3, 5, 7
		2152302898747,       // bases 2 … 11
		3474749660383,       // bases 2 … 13
		341550071728321,     // bases 2 … 17
		3825123056546413051, // bases 2 … 23
	} {
		if bigOf(n).ProbablyPrime(0) {
			t.Fatalf("%d is prime per math/big; the test table is wrong", n)
		}
		if IsPrime(n) {
			t.Errorf("IsPrime(%d) = true, want false (strong pseudoprime)", n)
		}
	}
}

// Oracle: trial division for every n below 2·10⁵.
func TestIsPrime_oracleSmall(t *testing.T) {
	for n := range uint64(200_000) {
		want := n >= 2
		for d := uint64(2); d*d <= n; d++ {
			if n%d == 0 {
				want = false
				break
			}
		}
		if got := IsPrime(n); got != want {
			t.Fatalf("IsPrime(%d) = %v, want %v", n, got, want)
		}
	}
}

// Oracle: math/big.Int.ProbablyPrime(0), which Go documents as 100% accurate below 2⁶⁴.
func TestIsPrime_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 3000 {
		n := rng.Uint64() >> rng.UintN(64) // every bit length
		if trial%2 == 0 {
			n |= 1 // odd numbers are prime far more often
		}
		if got, want := IsPrime(n), bigOf(n).ProbablyPrime(0); got != want {
			t.Fatalf("trial %d: IsPrime(%d) = %v, want %v", trial, n, got, want)
		}
	}
}

func TestIsPrime_noAllocation(t *testing.T) {
	if allocs := testing.AllocsPerRun(100, func() { benchSinkBool = IsPrime(math.MaxUint64 - 58) }); allocs != 0 {
		t.Errorf("IsPrime allocates %v times per call, want 0", allocs)
	}
}

// Benchmarks use primes, the worst case: all 12 bases run.
func BenchmarkIsPrime(b *testing.B) {
	for _, size := range []struct {
		name string
		n    uint64
	}{{"32bit", 4294967291}, {"64bit", math.MaxUint64 - 58}} {
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchSinkBool = IsPrime(size.n)
			}
		})
	}
}

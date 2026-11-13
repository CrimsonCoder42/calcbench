package integer

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"slices"
	"testing"
)

// checkFactorization fails t unless factors is ascending, every element is prime (math/big oracle) and
// their product, formed in math/big, is n.
func checkFactorization(t *testing.T, n uint64, factors []uint64) {
	t.Helper()
	if !slices.IsSorted(factors) {
		t.Fatalf("PrimeFactors(%d) = %v, not ascending", n, factors)
	}
	product := big.NewInt(1)
	for _, f := range factors {
		if !bigOf(f).ProbablyPrime(0) {
			t.Fatalf("PrimeFactors(%d) = %v: %d is not prime", n, factors, f)
		}
		product.Mul(product, bigOf(f))
	}
	if product.Cmp(bigOf(n)) != 0 {
		t.Fatalf("PrimeFactors(%d) = %v: product is %v", n, factors, product)
	}
}

// randomPrime32 returns a prime in [2³¹, 2³²), drawn with rng and checked by math/big.
func randomPrime32(rng *rand.Rand) uint64 {
	for {
		p := uint64(rng.Uint32()) | 1<<31 | 1
		if bigOf(p).ProbablyPrime(0) {
			return p
		}
	}
}

// Worked example, hand-checked: 2⁶⁴ − 1 = (2³² − 1)(2³² + 1) = (3·5·17·257·65537)(641·6700417).
func TestPrimeFactors_example(t *testing.T) {
	want := []uint64{3, 5, 17, 257, 641, 65537, 6700417}
	got, err := PrimeFactors(math.MaxUint64)
	if err != nil || !slices.Equal(got, want) {
		t.Errorf("PrimeFactors(2⁶⁴−1) = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestPrimeFactors_cases(t *testing.T) {
	sixtyThreeTwos := slices.Repeat([]uint64{2}, 63)
	tests := []struct {
		name string
		n    uint64
		want []uint64
	}{
		{"two", 2, []uint64{2}},
		{"three", 3, []uint64{3}},
		{"four", 4, []uint64{2, 2}},
		{"twelve", 12, []uint64{2, 2, 3}},
		{"2¹⁰", 1024, slices.Repeat([]uint64{2}, 10)},
		{"2⁶³, the most factors a uint64 has", 1 << 63, sixtyThreeTwos},
		{"largest table prime times the next prime", 1021 * 1031, []uint64{1021, 1031}},
		{"1031², least composite left to rho", 1031 * 1031, []uint64{1031, 1031}},
		{"two six-digit primes", 999983 * 1000003, []uint64{999983, 1000003}},
		{"two 16-bit primes", 65519 * 65521, []uint64{65519, 65521}},
		{"two 32-bit primes", 4294967279 * 4294967291, []uint64{4294967279, 4294967291}},
		{"square of the largest 32-bit prime", 4294967291 * 4294967291, []uint64{4294967291, 4294967291}},
		{"cube of the prime 2097143", 2097143 * 2097143 * 2097143, []uint64{2097143, 2097143, 2097143}},
		{"small times a large prime", 3 * 4294967291, []uint64{3, 4294967291}},
		{"largest prime below 2⁶⁴", math.MaxUint64 - 58, []uint64{math.MaxUint64 - 58}},
		{"Mersenne prime 2⁶¹−1", 1<<61 - 1, []uint64{1<<61 - 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PrimeFactors(tt.n)
			if err != nil || !slices.Equal(got, tt.want) {
				t.Errorf("PrimeFactors(%d) = (%v, %v), want (%v, nil)", tt.n, got, err, tt.want)
			}
		})
	}
}

func TestPrimeFactors_errors(t *testing.T) {
	for _, n := range []uint64{0, 1} {
		if got, err := PrimeFactors(n); !errors.Is(err, ErrOutOfRange) || got != nil {
			t.Errorf("PrimeFactors(%d) = (%v, %v), want (nil, ErrOutOfRange)", n, got, err)
		}
	}
}

// Oracle: trial division for every n in [2, 20000).
func TestPrimeFactors_oracleSmall(t *testing.T) {
	for n := uint64(2); n < 20_000; n++ {
		var want []uint64
		m := n
		for d := uint64(2); d*d <= m; d++ {
			for m%d == 0 {
				want = append(want, d)
				m /= d
			}
		}
		if m > 1 {
			want = append(want, m)
		}
		if got, err := PrimeFactors(n); err != nil || !slices.Equal(got, want) {
			t.Fatalf("PrimeFactors(%d) = (%v, %v), want (%v, nil)", n, got, err, want)
		}
	}
}

// Oracle: random inputs of every bit length, and semiprimes of two random 32-bit primes (the hardest case
// for rho), checked by math/big.
func TestPrimeFactors_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 400 {
		n := max(2, rng.Uint64()>>rng.UintN(64))
		got, err := PrimeFactors(n)
		if err != nil {
			t.Fatalf("PrimeFactors(%d): %v", n, err)
		}
		checkFactorization(t, n, got)
	}
	for range 100 {
		p, q := randomPrime32(rng), randomPrime32(rng)
		got, err := PrimeFactors(p * q)
		if want := []uint64{min(p, q), max(p, q)}; err != nil || !slices.Equal(got, want) {
			t.Fatalf("PrimeFactors(%d·%d) = (%v, %v), want (%v, nil)", p, q, got, err, want)
		}
	}
}

func TestPrimeFactors_oneAllocation(t *testing.T) {
	var sink []uint64
	if allocs := testing.AllocsPerRun(20, func() { sink, _ = PrimeFactors(4294967279 * 4294967291) }); allocs != 1 {
		t.Errorf("PrimeFactors allocates %v times per call, want 1 (the result)", allocs)
	}
	_ = sink
}

// The trial table holds exactly the primes below 1024, in order.
func TestSmallPrimes_table(t *testing.T) {
	var want []uint64
	for n := range uint64(1024) {
		if bigOf(n).ProbablyPrime(0) {
			want = append(want, n)
		}
	}
	if !slices.Equal(smallPrimes[:], want) {
		t.Errorf("smallPrimes differs from the primes below 1024 (%d of them)", len(want))
	}
	if trialBound != 1031*1031 || !bigOf(1031).ProbablyPrime(0) {
		t.Error("trialBound must be the square of 1031, the first prime above the table")
	}
	for n := uint64(1022); n < 1031; n++ {
		if bigOf(n).ProbablyPrime(0) {
			t.Errorf("%d is prime: 1031 is not the next prime after the table", n)
		}
	}
}

// Benchmarks use semiprimes of two equal-sized primes, the worst case for rho.
var factorBenchSizes = []struct {
	name string
	n    uint64
}{{"32bit", 65519 * 65521}, {"64bit", 4294967279 * 4294967291}}

func BenchmarkPrimeFactors(b *testing.B) {
	for _, size := range factorBenchSizes {
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				f, _ := PrimeFactors(size.n)
				benchSink = f[0]
			}
		})
	}
}

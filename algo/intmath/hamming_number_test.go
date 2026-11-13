package intmath

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"slices"
	"testing"
)

// Dijkstra, A Discipline of Programming, ch. 17: the sequence begins 1, 2, 3, 4, 5, 6, 8, 9, 10, 12, 15,
// 16, 18, 20, 24, 25, 27, 30.
var hammingStart = []uint64{1, 2, 3, 4, 5, 6, 8, 9, 10, 12, 15, 16, 18, 20, 24, 25, 27, 30}

func TestNthHammingNumber_sourceExample(t *testing.T) {
	for i, want := range hammingStart {
		if got, err := NthHammingNumber(i + 1); err != nil || got != want {
			t.Errorf("NthHammingNumber(%d) = (%d, %v), want %d", i+1, got, err, want)
		}
	}
}

func TestIsHammingNumber_sourceExample(t *testing.T) {
	for n := uint64(0); n <= 30; n++ {
		if got, want := IsHammingNumber(n), slices.Contains(hammingStart, n); got != want {
			t.Errorf("IsHammingNumber(%d) = %v, want %v", n, got, want)
		}
	}
}

func TestIsHammingNumber_cases(t *testing.T) {
	for _, n := range []uint64{1 << 63, 18_432_000_000_000_000_000, 12_157_665_459_056_928_801, 7_450_580_596_923_828_125} {
		if !IsHammingNumber(n) {
			t.Errorf("IsHammingNumber(%d) = false, want true", n)
		}
	}
	for _, n := range []uint64{0, 7, 14, 49, math.MaxUint64, 1<<63 + 1} {
		if IsHammingNumber(n) {
			t.Errorf("IsHammingNumber(%d) = true, want false", n)
		}
	}
}

// The largest Hamming number below 2^64 is the 13282nd: recount them with math/big.
func TestNthHammingNumber_extreme(t *testing.T) {
	limit := new(big.Int).SetUint64(math.MaxUint64)
	count := 0
	for p2 := big.NewInt(1); p2.Cmp(limit) <= 0; p2.Lsh(p2, 1) {
		for p3 := new(big.Int).Set(p2); p3.Cmp(limit) <= 0; p3.Mul(p3, big.NewInt(3)) {
			for p5 := new(big.Int).Set(p3); p5.Cmp(limit) <= 0; p5.Mul(p5, big.NewInt(5)) {
				count++
			}
		}
	}
	if count != hammingCountUint64 {
		t.Fatalf("Hamming numbers ≤ 2^64−1: %d, constant says %d", count, hammingCountUint64)
	}
	if got, err := NthHammingNumber(hammingCountUint64); err != nil || got != 18_432_000_000_000_000_000 {
		t.Errorf("NthHammingNumber(13282) = (%d, %v), want 18432000000000000000", got, err)
	}
}

func TestNthHammingNumber_errors(t *testing.T) {
	for _, tt := range []struct {
		n       int
		wantErr error
	}{{0, ErrOutOfRange}, {-1, ErrOutOfRange}, {math.MinInt, ErrOutOfRange}, {hammingCountUint64 + 1, ErrOverflow}, {math.MaxInt, ErrOverflow}} {
		if _, err := NthHammingNumber(tt.n); !errors.Is(err, tt.wantErr) {
			t.Errorf("NthHammingNumber(%d) error = %v, want %v", tt.n, err, tt.wantErr)
		}
	}
}

// Oracle: scan the integers with IsHammingNumber's slow twin (trial division) for the first 2000.
func TestNthHammingNumber_oracle(t *testing.T) {
	slow := func(n uint64) bool {
		for _, p := range []uint64{2, 3, 5} {
			for n > 0 && n%p == 0 {
				n /= p
			}
		}
		return n == 1
	}
	// The first 300 Hamming numbers end near 4·10⁵, so this exhaustive scan stays well inside the
	// 60-second test cap (2000 reached 8·10⁹ and took 142 s).
	k := 0
	for n := uint64(1); k < 300; n++ {
		if slow(n) {
			k++
			if got, err := NthHammingNumber(k); err != nil || got != n {
				t.Fatalf("NthHammingNumber(%d) = (%d, %v), want %d", k, got, err, n)
			}
		}
		if IsHammingNumber(n) != slow(n) {
			t.Fatalf("IsHammingNumber(%d) disagrees with trial division", n)
		}
	}
	rng := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		n := rng.Uint64()
		if IsHammingNumber(n) != slow(n) {
			t.Fatalf("IsHammingNumber(%d) disagrees with trial division", n)
		}
	}
}

func BenchmarkIsHammingNumber(b *testing.B) {
	b.ReportAllocs()
	n := 0
	for i := uint64(0); b.Loop(); i++ {
		if IsHammingNumber(18_432_000_000_000_000_000 - i&1) {
			n++
		}
	}
	benchSink = uint64(n)
}

func BenchmarkNthHammingNumber(b *testing.B) {
	for _, size := range []struct {
		name string
		n    int
	}{{"1e2", 100}, {"1e3", 1_000}, {"1e4", 10_000}} {
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			var s uint64
			for b.Loop() {
				v, _ := NthHammingNumber(size.n)
				s += v
			}
			benchSink = s
		})
	}
}

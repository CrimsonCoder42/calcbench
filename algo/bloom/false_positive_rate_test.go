package bloom

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// Broder–Mitzenmacher §2.1: at m/n = 10 bits per item and k = 7 (the integer nearest ln 2 · 10 = 6.93),
// f = (1 − e^(−0.7))^7 ≈ 0.008194 (hand-computed). Relative tolerance 1e-12: two library calls, each
// within a few ulps.
func TestFalsePositiveRate_surveyExample(t *testing.T) {
	got, err := FalsePositiveRate(10_000, 1_000, 7)
	want := 0.008193722065862417
	if err != nil || math.Abs(got-want) > 1e-12*want {
		t.Errorf("FalsePositiveRate(10000, 1000, 7) = (%v, %v), want %v", got, err, want)
	}
	// At the optimum k = ln 2 · m/n the rate is (1/2)^k (§2.1): m = 1 000 000, n = 69 315 gives k·n/m =
	// 0.693150 ≈ ln 2 (off by 3e-6) for k = 10, so the rate is 2^-10 within 1e-4 relative (k·3e-6 = 3e-5).
	got, _ = FalsePositiveRate(1_000_000, 69_315, 10)
	if want := math.Pow(0.5, 10); math.Abs(got-want) > 1e-4*want {
		t.Errorf("FalsePositiveRate at the optimum = %v, want %v", got, want)
	}
}

func TestFalsePositiveRate_cases(t *testing.T) {
	tests := []struct {
		name            string
		bitCount, items uint64
		hashCount       int
		want, relTol    float64
	}{
		{"no items", 100, 0, 3, 0, 0},
		{"one bit, one item, one hash", 1, 1, 1, 1 - math.Exp(-1), 1e-15},
		// x = 2⁻³⁰: 1 − e^(−x) = x − x²/2 + …; a plain 1 − Exp(−x) would lose about 8 digits here.
		{"tiny load", MaxBits, 1, 1, 0x1p-30 - 0x1p-61, 1e-15},
		{"saturated", 1, math.MaxUint64, 1, 1, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FalsePositiveRate(tt.bitCount, tt.items, tt.hashCount)
			if err != nil || math.Abs(got-tt.want) > tt.relTol*tt.want {
				t.Errorf("FalsePositiveRate(%d, %d, %d) = (%v, %v), want %v", tt.bitCount, tt.items, tt.hashCount, got, err, tt.want)
			}
		})
	}
}

func TestFalsePositiveRate_errors(t *testing.T) {
	if _, err := FalsePositiveRate(0, 10, 3); !errors.Is(err, ErrBitCount) {
		t.Errorf("bitCount 0: err = %v, want ErrBitCount", err)
	}
	if _, err := FalsePositiveRate(100, 10, 0); !errors.Is(err, ErrHashCount) {
		t.Errorf("hashCount 0: err = %v, want ErrHashCount", err)
	}
	if _, err := FalsePositiveRate(100, 10, -1); !errors.Is(err, ErrHashCount) {
		t.Errorf("hashCount -1: err = %v, want ErrHashCount", err)
	}
}

// Oracle: the exact form (1 − (1 − 1/m)^(kn))^k of §2.1, computed in log space. The exponents differ by
// at most kn/(2m²) + …, which moves the rate by at most about k/(2m) relative; the tolerance is k/m.
// Also: the rate never decreases as items grow.
func TestFalsePositiveRate_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		m := 1000 + rng.Uint64N(1_000_000)
		k := 1 + rng.IntN(20)
		n := rng.Uint64N(2 * m / uint64(k))
		got, err := FalsePositiveRate(m, n, k)
		zeroBit := math.Exp(float64(k) * float64(n) * math.Log1p(-1/float64(m)))
		exact := math.Pow(1-zeroBit, float64(k))
		if err != nil || math.Abs(got-exact) > float64(k)/float64(m)*exact {
			t.Fatalf("trial %d: FalsePositiveRate(%d, %d, %d) = (%v, %v), exact form %v", trial, m, n, k, got, err, exact)
		}
		more, _ := FalsePositiveRate(m, n+1, k)
		if more < got {
			t.Fatalf("trial %d: rate fell from %v to %v when an item was added", trial, got, more)
		}
	}
}

func BenchmarkFalsePositiveRate(b *testing.B) {
	b.Run("scalar", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			FalsePositiveRate(9586, 1000, 7)
		}
	})
}

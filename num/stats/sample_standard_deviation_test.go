package stats

import (
	"errors"
	"math"
	"math/big"
	"math/rand/v2"
	"testing"
)

// Worked example: the certified sample SDs of StRD Univariate (Michelso → 0.0790105478190518).
func TestSampleStandardDeviation_strd(t *testing.T) {
	for _, tt := range strdSD {
		t.Run(tt.name, func(t *testing.T) {
			// Tolerance per dataset, explained at strdSD.
			if got, err := SampleStandardDeviation(tt.values); err != nil || !near(got, tt.sd, tt.tol) {
				t.Errorf("SampleStandardDeviation(%s) = %v, %v; want %v", tt.name, got, err, tt.sd)
			}
		})
	}
}

func TestSampleStandardDeviation_cases(t *testing.T) {
	m := math.MaxFloat64
	tests := []struct {
		name   string
		values []float64
		want   float64
	}{
		{"two elements", []float64{1, 3}, math.Sqrt2},
		{"all equal", []float64{9, 9}, 0},
		{"textbook", []float64{2, 4, 4, 4, 5, 5, 7, 9}, math.Sqrt(32.0 / 7)},
		// The variance (m²/2) overflows, the SD (m/√2) does not.
		{"variance overflows, SD finite", []float64{-m / 2, m / 2}, m / math.Sqrt2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Tolerance 1e-15: a few roundings.
			if got, err := SampleStandardDeviation(tt.values); err != nil || !near(got, tt.want, 1e-15) {
				t.Errorf("SampleStandardDeviation(%v) = %v, %v; want %v", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestSampleStandardDeviation_errors(t *testing.T) {
	m := math.MaxFloat64
	tests := []struct {
		name   string
		values []float64
		want   error
	}{
		{"empty", []float64{}, ErrEmpty},
		{"one element", []float64{3}, ErrTooFew},
		{"NaN", []float64{math.NaN(), 1}, ErrNotFinite},
		{"overflow", []float64{-m, m}, ErrOverflow}, // s = √2·m
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := SampleStandardDeviation(tt.values); !errors.Is(err, tt.want) || got != 0 {
				t.Errorf("SampleStandardDeviation(%v) = %v, %v; want 0, %v", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestSampleStandardDeviation_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		n := 2 + rng.IntN(400)
		values := randomFloats(rng, n)
		want := math.Sqrt(bigSumOfSquares(values) / float64(n-1))
		// Tolerance 1e-13: O(n·u), n ≤ 401, halved by the square root.
		if got, err := SampleStandardDeviation(values); err != nil || !near(got, want, 1e-13) {
			t.Fatalf("trial %d: SampleStandardDeviation(n=%d) = %v, %v; oracle %v", trial, n, got, err, want)
		}
	}
}

func TestSampleStandardDeviation_noAllocation(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { benchSink, _ = SampleStandardDeviation(michelso) }); a != 0 {
		t.Errorf("SampleStandardDeviation allocates %v times per call, want 0", a)
	}
}

func BenchmarkSampleStandardDeviation(b *testing.B) {
	for _, size := range benchSizes {
		inputs := benchInputs(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; b.Loop(); i++ {
				benchSink, _ = SampleStandardDeviation(inputs[i&7])
			}
		})
	}
}

// exactCentred returns the exact rational mean of x and the exact centred sums Σ(x − x̄)², Σ(y − ȳ)² and
// Σ(x − x̄)(y − ȳ) of the float64 inputs (y may be nil). It is the oracle for the tiny-scale tests: no
// rounding, no underflow, independent of the package's scaling.
func exactCentred(x, y []float64) (xMean, yMean, sxx, syy, sxy *big.Rat) {
	mean := func(v []float64) *big.Rat {
		s := new(big.Rat)
		for _, f := range v {
			s.Add(s, new(big.Rat).SetFloat64(f))
		}
		return s.Quo(s, big.NewRat(int64(len(v)), 1))
	}
	xMean, sxx, syy, sxy = mean(x), new(big.Rat), new(big.Rat), new(big.Rat)
	yMean = new(big.Rat)
	if y != nil {
		yMean = mean(y)
	}
	for i, f := range x {
		dx := new(big.Rat).Sub(new(big.Rat).SetFloat64(f), xMean)
		sxx.Add(sxx, new(big.Rat).Mul(dx, dx))
		if y != nil {
			dy := new(big.Rat).Sub(new(big.Rat).SetFloat64(y[i]), yMean)
			syy.Add(syy, new(big.Rat).Mul(dy, dy))
			sxy.Add(sxy, new(big.Rat).Mul(dx, dy))
		}
	}
	return xMean, yMean, sxx, syy, sxy
}

// ratSqrt returns √r rounded to float64 via 512-bit big.Float (correctly rounded to well within 1 ulp).
func ratSqrt(r *big.Rat) float64 {
	f := new(big.Float).SetPrec(512).SetRat(r)
	v, _ := f.Sqrt(f).Float64()
	return v
}

// ratFloat returns r rounded to float64 (subnormals included).
func ratFloat(r *big.Rat) float64 {
	v, _ := r.Float64()
	return v
}

// Tiny spreads: before the underflow fix, squared deviations below ~1e-308 vanished and these returned 0.
// Expected values: the 2^k-scaled cases are exact (s(c·x) = c·s(x) for c a power of two, and s(1, 2, 3)
// = 1 by hand: deviations −1, 0, 1, Σd² = 2, 2/(3 − 1) = 1); the decimal case is checked against exact
// rational arithmetic on the float64 inputs.
func TestSampleStandardDeviation_tinySpread(t *testing.T) {
	for _, tt := range []struct {
		name   string
		values []float64
		want   float64
	}{
		{"1,2,3 × 2^-700", []float64{0x1p-700, 0x2p-700, 0x3p-700}, 0x1p-700},
		{"1,2,3 × 2^-1000 (about 1e-301)", []float64{0x1p-1000, 0x2p-1000, 0x3p-1000}, 0x1p-1000},
		{"1,2,3 × 2^-1070 (subnormal)", []float64{0x1p-1070, 0x2p-1070, 0x3p-1070}, 0x1p-1070},
		// √(1/2)·2^-1074 rounds to the smallest subnormal, 2^-1074.
		{"two smallest subnormals", []float64{0x1p-1074, 0x1p-1073}, 0x1p-1074},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := SampleStandardDeviation(tt.values); err != nil || got != tt.want {
				t.Errorf("SampleStandardDeviation(%v) = %v, %v; want %v", tt.values, got, err, tt.want)
			}
		})
	}
	values := []float64{1e-200, 2e-200, 3e-200}
	_, _, sxx, _, _ := exactCentred(values, nil)
	want := ratSqrt(sxx.Quo(sxx, big.NewRat(2, 1))) // ≈ 1e-200
	// Tolerance 4e-16 relative: a few roundings in the scaled two-pass sum, the square root and the divide.
	if got, err := SampleStandardDeviation(values); err != nil || !near(got, want, 4e-16) || !near(got, 1e-200, 1e-15) {
		t.Errorf("SampleStandardDeviation(1e-200, 2e-200, 3e-200) = %v, %v; want %v", got, err, want)
	}
}

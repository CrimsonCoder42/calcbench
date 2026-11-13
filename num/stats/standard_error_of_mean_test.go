package stats

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// Worked example: StRD Michelso, certified s = 0.0790105478190518, n = 100 → s/√n = 0.00790105478190518.
// The e-Handbook §1.3.5.2 ZARR13 example (N = 195, s = 0.022789 → 0.0016320) is checked on its printed
// summary values in TestStandardErrorOfMean_zarr13Arithmetic.
func TestStandardErrorOfMean_strd(t *testing.T) {
	// Tolerance 1e-13: as the SD (strdSD).
	if got, err := StandardErrorOfMean(michelso); err != nil || !near(got, 0.00790105478190518, 1e-13) {
		t.Errorf("StandardErrorOfMean(Michelso) = %v, %v; want 0.00790105478190518", got, err)
	}
}

// e-Handbook §1.3.5.2: the printed 95% interval (9.258242, 9.264679) has half-width / t(0.975, 194) =
// 0.0032185/1.9723 = 0.0016319, and s/√N = 0.022789/√195 = 0.0016320: s/√N is the formula behind it.
func TestStandardErrorOfMean_zarr13Arithmetic(t *testing.T) {
	sem := 0.022789 / math.Sqrt(195)
	halfWidth := (9.264679 - 9.258242) / 2
	// Tolerance 1e-4 relative: the printed values have 5 to 7 significant digits.
	if !near(sem, 0.0016320, 1e-4) || !near(halfWidth/1.9723, sem, 1e-4) {
		t.Errorf("ZARR13: s/√N = %v, half-width/t = %v; want both 0.0016320", sem, halfWidth/1.9723)
	}
}

func TestStandardErrorOfMean_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   float64
	}{
		{"two elements", []float64{1, 3}, 1},
		{"all equal", []float64{2, 2, 2}, 0},
		{"textbook", []float64{2, 4, 4, 4, 5, 5, 7, 9}, math.Sqrt(32.0/7) / math.Sqrt(8)},
		{"s overflows, SEM finite", []float64{-math.MaxFloat64, math.MaxFloat64}, math.MaxFloat64},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Tolerance 1e-15: a few roundings.
			if got, err := StandardErrorOfMean(tt.values); err != nil || !near(got, tt.want, 1e-15) {
				t.Errorf("StandardErrorOfMean(%v) = %v, %v; want %v", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestStandardErrorOfMean_errors(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   error
	}{
		{"empty", nil, ErrEmpty},
		{"one element", []float64{1}, ErrTooFew},
		{"Inf", []float64{1, math.Inf(1)}, ErrNotFinite},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := StandardErrorOfMean(tt.values); !errors.Is(err, tt.want) || got != 0 {
				t.Errorf("StandardErrorOfMean(%v) = %v, %v; want 0, %v", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestStandardErrorOfMean_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		n := 2 + rng.IntN(400)
		values := randomFloats(rng, n)
		want := math.Sqrt(bigSumOfSquares(values)/float64(n-1)) / math.Sqrt(float64(n))
		// Tolerance 1e-13: O(n·u).
		if got, err := StandardErrorOfMean(values); err != nil || !near(got, want, 1e-13) {
			t.Fatalf("trial %d: StandardErrorOfMean(n=%d) = %v, %v; oracle %v", trial, n, got, err, want)
		}
	}
}

func TestStandardErrorOfMean_noAllocation(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { benchSink, _ = StandardErrorOfMean(michelso) }); a != 0 {
		t.Errorf("StandardErrorOfMean allocates %v times per call, want 0", a)
	}
}

func BenchmarkStandardErrorOfMean(b *testing.B) {
	for _, size := range benchSizes {
		inputs := benchInputs(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; b.Loop(); i++ {
				benchSink, _ = StandardErrorOfMean(inputs[i&7])
			}
		})
	}
}

// Tiny spread: s(1, 2, 3) = 1 (hand), so the SEM of 2^-700·(1, 2, 3) is 2^-700/√3, exactly the float64
// quotient (both sides round the same real number once). At 1e-200 it is 1e-200/√3 = 5.773502691896258e-201
// to within the inputs' rounding. Before the underflow fix both returned 0.
func TestStandardErrorOfMean_tinySpread(t *testing.T) {
	want := 0x1p-700 / math.Sqrt(3)
	if got, err := StandardErrorOfMean([]float64{0x1p-700, 0x2p-700, 0x3p-700}); err != nil || got != want {
		t.Errorf("StandardErrorOfMean(2^-700·(1 2 3)) = %v, %v; want %v", got, err, want)
	}
	// Tolerance 1e-15 relative: the decimal inputs are each within 2^-53 of k·1e-200.
	if got, err := StandardErrorOfMean([]float64{1e-200, 2e-200, 3e-200}); err != nil || !near(got, 1e-200/math.Sqrt(3), 1e-15) {
		t.Errorf("StandardErrorOfMean(1e-200 2e-200 3e-200) = %v, %v; want %v", got, err, 1e-200/math.Sqrt(3))
	}
}

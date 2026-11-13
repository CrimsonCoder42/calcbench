package stats

import (
	"errors"
	"math"
	"testing"
)

// Worked example: StRD Michelso, certified s = 0.0790105478190518 and mean 299.8524 →
// cv = 2.6349813381200816e-04 (Dataplot "COEFFICIENT OF VARIATION": cv = s/x̄).
func TestCoefficientOfVariation_strd(t *testing.T) {
	// Tolerance 1e-13: as the SD (strdSD).
	if got, err := CoefficientOfVariation(michelso); err != nil || !near(got, 0.0790105478190518/299.8524, 1e-13) {
		t.Errorf("CoefficientOfVariation(Michelso) = %v, %v", got, err)
	}
}

func TestCoefficientOfVariation_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   float64
	}{
		{"two elements", []float64{1, 3}, math.Sqrt2 / 2},
		{"all equal", []float64{4, 4, 4}, 0},
		{"negative mean keeps the sign", []float64{-1, -3}, -math.Sqrt2 / 2},
		{"NumAcc1", numAcc1, 1.0 / 10000002},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Tolerance 1e-15: a few roundings.
			if got, err := CoefficientOfVariation(tt.values); err != nil || !near(got, tt.want, 1e-15) {
				t.Errorf("CoefficientOfVariation(%v) = %v, %v; want %v", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestCoefficientOfVariation_errors(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   error
	}{
		{"empty", nil, ErrEmpty},
		{"one element", []float64{2}, ErrTooFew},
		{"NaN", []float64{1, math.NaN()}, ErrNotFinite},
		{"zero mean", []float64{-1, 1}, ErrZeroMean},
		{"overflow", []float64{-1e300, 1e300, 1e-10}, ErrOverflow}, // s ≈ 1e300, mean ≈ 3.3e-11
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := CoefficientOfVariation(tt.values); !errors.Is(err, tt.want) || got != 0 {
				t.Errorf("CoefficientOfVariation(%v) = %v, %v; want 0, %v", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestCoefficientOfVariation_noAllocation(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { benchSink, _ = CoefficientOfVariation(michelso) }); a != 0 {
		t.Errorf("CoefficientOfVariation allocates %v times per call, want 0", a)
	}
}

func BenchmarkCoefficientOfVariation(b *testing.B) {
	for _, size := range benchSizes {
		inputs := benchInputs(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; b.Loop(); i++ {
				benchSink, _ = CoefficientOfVariation(inputs[i&7])
			}
		})
	}
}

// Tiny spread: 1, 2, 3 has s = 1 and mean 2 (hand), so CV = 1/2 at any power-of-two scale; at 1e-200 the
// decimal inputs give 1/2 to within their rounding. Before the underflow fix both returned 0.
func TestCoefficientOfVariation_tinySpread(t *testing.T) {
	if got, err := CoefficientOfVariation([]float64{0x1p-700, 0x2p-700, 0x3p-700}); err != nil || got != 0.5 {
		t.Errorf("CoefficientOfVariation(2^-700·(1 2 3)) = %v, %v; want 0.5", got, err)
	}
	if got, err := CoefficientOfVariation([]float64{0x1p-1070, 0x2p-1070, 0x3p-1070}); err != nil || got != 0.5 {
		t.Errorf("CoefficientOfVariation(2^-1070·(1 2 3)) = %v, %v; want 0.5", got, err)
	}
	// Tolerance 1e-15 relative: the inputs 1e-200, 2e-200, 3e-200 are each within 2^-53 of k·1e-200.
	if got, err := CoefficientOfVariation([]float64{1e-200, 2e-200, 3e-200}); err != nil || !near(got, 0.5, 1e-15) {
		t.Errorf("CoefficientOfVariation(1e-200 2e-200 3e-200) = %v, %v; want 0.5", got, err)
	}
	// Mean far below the spread: −2^-1000, 2^-1000, 3·2^-1074 has mean exactly 2^-1074 and
	// s² = (Σx² − 3·mean²)/2 = 2^-2000 + 3·2^-2148 (hand), so s/mean = 2^74·(1 + 1.5·2^-148) → 2^74. The
	// unscaled quotient s/mean would overflow on the way (0.5/2^-1074); the result itself is finite.
	if got, err := CoefficientOfVariation([]float64{-0x1p-1000, 0x1p-1000, 0x3p-1074}); err != nil || !near(got, 0x1p74, 1e-15) {
		t.Errorf("CoefficientOfVariation(−2^-1000, 2^-1000, 3·2^-1074) = %v, %v; want 2^74", got, err)
	}
}

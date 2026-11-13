package stats

import (
	"errors"
	"math"
	"testing"
)

// Worked example: StRD Michelso → 100·0.0790105478190518/299.8524 = 0.026349813381200816 (percent), from
// the certified s and mean (Dataplot "RELATIVE STANDARD DEVIATION": relsd = 100·(s/|x̄|)).
func TestPercentRelativeStandardDeviation_strd(t *testing.T) {
	// Tolerance 1e-13: as the SD (strdSD).
	if got, err := PercentRelativeStandardDeviation(michelso); err != nil || !near(got, 100*0.0790105478190518/299.8524, 1e-13) {
		t.Errorf("PercentRelativeStandardDeviation(Michelso) = %v, %v", got, err)
	}
}

func TestPercentRelativeStandardDeviation_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   float64
	}{
		{"two elements", []float64{1, 3}, 50 * math.Sqrt2},
		{"negative mean gives a positive result", []float64{-1, -3}, 50 * math.Sqrt2},
		{"all equal", []float64{-7, -7}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Tolerance 1e-15: a few roundings.
			if got, err := PercentRelativeStandardDeviation(tt.values); err != nil || !near(got, tt.want, 1e-15) {
				t.Errorf("PercentRelativeStandardDeviation(%v) = %v, %v; want %v", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestPercentRelativeStandardDeviation_errors(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   error
	}{
		{"empty", nil, ErrEmpty},
		{"one element", []float64{2}, ErrTooFew},
		{"Inf", []float64{math.Inf(1), 1}, ErrNotFinite},
		{"zero mean", []float64{-2, 0, 2}, ErrZeroMean},
		{"overflow", []float64{-1e300, 1e300, 1e-10}, ErrOverflow},                // s ≈ 1e300, mean ≈ 3.3e-11
		{"overflow only after ×100", []float64{-1e300, 1e300, 3e-7}, ErrOverflow}, // s/|x̄| ≈ 1e307
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := PercentRelativeStandardDeviation(tt.values); !errors.Is(err, tt.want) || got != 0 {
				t.Errorf("PercentRelativeStandardDeviation(%v) = %v, %v; want 0, %v", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestPercentRelativeStandardDeviation_noAllocation(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { benchSink, _ = PercentRelativeStandardDeviation(michelso) }); a != 0 {
		t.Errorf("PercentRelativeStandardDeviation allocates %v times per call, want 0", a)
	}
}

func BenchmarkPercentRelativeStandardDeviation(b *testing.B) {
	for _, size := range benchSizes {
		inputs := benchInputs(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; b.Loop(); i++ {
				benchSink, _ = PercentRelativeStandardDeviation(inputs[i&7])
			}
		})
	}
}

// Tiny spread: 1, 2, 3 has s/|mean| = 1/2 (hand), so %RSD = 50 at any power-of-two scale and, to within
// the inputs' rounding, at 1e-200. Before the underflow fix both returned 0.
func TestPercentRelativeStandardDeviation_tinySpread(t *testing.T) {
	if got, err := PercentRelativeStandardDeviation([]float64{0x1p-1000, 0x2p-1000, 0x3p-1000}); err != nil || got != 50 {
		t.Errorf("PercentRelativeStandardDeviation(2^-1000·(1 2 3)) = %v, %v; want 50", got, err)
	}
	// Tolerance 1e-15 relative: the decimal inputs are each within 2^-53 of k·1e-200.
	if got, err := PercentRelativeStandardDeviation([]float64{1e-200, 2e-200, 3e-200}); err != nil || !near(got, 50, 1e-15) {
		t.Errorf("PercentRelativeStandardDeviation(1e-200 2e-200 3e-200) = %v, %v; want 50", got, err)
	}
}

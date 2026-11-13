package stats

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Worked example: e-Handbook §1.3.5.6, MAD = median(|Y − Ỹ|). For 1, 1, 2, 2, 4, 6, 9: Ỹ = 2, absolute
// deviations 1, 1, 0, 0, 2, 4, 7 → 1 (hand-checked). The 12 resistivity values of §7.2.6.2: Ỹ = 95.1579,
// sorted deviations put 0.0358 and 0.0380 in the middle → 0.0369 (hand-checked).
func TestMedianAbsoluteDeviation_worked(t *testing.T) {
	if got, err := MedianAbsoluteDeviation([]float64{1, 1, 2, 2, 4, 6, 9}); err != nil || got != 1 {
		t.Errorf("MedianAbsoluteDeviation(1 1 2 2 4 6 9) = %v, %v; want 1", got, err)
	}
	// Tolerance 1e-12: deviations of ~0.04 from values near 95 carry the inputs' decimal rounding.
	if got, err := MedianAbsoluteDeviation(resistivity); err != nil || !near(got, 0.0369, 1e-12) {
		t.Errorf("MedianAbsoluteDeviation(resistivity) = %v, %v; want 0.0369", got, err)
	}
}

func TestMedianAbsoluteDeviation_cases(t *testing.T) {
	m := math.MaxFloat64
	for _, tt := range []struct {
		name   string
		values []float64
		want   float64
	}{
		{"one element", []float64{3}, 0},
		{"two elements", []float64{1, 5}, 2},
		{"all equal", []float64{7, 7, 7}, 0},
		{"one outlier", []float64{1, 2, 3, 4, 1000}, 1},
		{"deviation overflows", []float64{-m, m}, m},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := MedianAbsoluteDeviation(tt.values); err != nil || got != tt.want {
				t.Errorf("MedianAbsoluteDeviation(%v) = %v, %v; want %v", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestMedianAbsoluteDeviation_errors(t *testing.T) {
	for _, tt := range []struct {
		name   string
		values []float64
		want   error
	}{
		{"empty", nil, ErrEmpty},
		{"NaN", []float64{math.NaN()}, ErrNotFinite},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := MedianAbsoluteDeviation(tt.values); !errors.Is(err, tt.want) || got != 0 {
				t.Errorf("MedianAbsoluteDeviation(%v) = %v, %v; want 0, %v", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestMedianAbsoluteDeviation_oracle(t *testing.T) {
	sortedMedian := func(v []float64) float64 {
		s := slices.Sorted(slices.Values(v))
		return (s[(len(s)-1)/2] + s[len(s)/2]) / 2
	}
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		values := randomFloats(rng, 1+rng.IntN(300))
		med := sortedMedian(values)
		dev := make([]float64, len(values))
		for i, v := range values {
			dev[i] = math.Abs(v - med)
		}
		want := sortedMedian(dev)
		// Tolerance 1e-13 relative (absolute when 0): the two medians round their averages differently.
		if got, err := MedianAbsoluteDeviation(values); err != nil || !near(got, want, 1e-13) {
			t.Fatalf("trial %d: MedianAbsoluteDeviation = %v, %v; want %v", trial, got, err, want)
		}
	}
}

func TestMedianAbsoluteDeviation_doesNotModifyInput(t *testing.T) {
	values := slices.Clone(resistivity)
	if _, err := MedianAbsoluteDeviation(values); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(values, resistivity) {
		t.Error("MedianAbsoluteDeviation modified its input")
	}
}

func BenchmarkMedianAbsoluteDeviation(b *testing.B) {
	for _, size := range benchSizes {
		inputs := benchInputs(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; b.Loop(); i++ {
				benchSink, _ = MedianAbsoluteDeviation(inputs[i&7])
			}
		})
	}
}

// Space line: one allocation (the buffer that holds the copy, then the deviations).
func TestMedianAbsoluteDeviation_allocs(t *testing.T) {
	values := randomFloats(rand.New(rand.NewPCG(1, 2)), 1000)
	if got := testing.AllocsPerRun(100, func() { _, _ = MedianAbsoluteDeviation(values) }); got != 1 {
		t.Errorf("MedianAbsoluteDeviation allocations = %v; want 1", got)
	}
}

// The input is never modified, though one buffer is reordered twice.
func TestMedianAbsoluteDeviation_inputUnchanged(t *testing.T) {
	values := randomFloats(rand.New(rand.NewPCG(3, 4)), 101)
	before := slices.Clone(values)
	if _, err := MedianAbsoluteDeviation(values); err != nil || !slices.Equal(values, before) {
		t.Errorf("MedianAbsoluteDeviation modified its input (err %v)", err)
	}
}

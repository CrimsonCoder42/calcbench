package stats

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// Worked example: Dataplot "WEIGHTED MEAN", x̄w = Σw·x/Σw. Values 80, 90 with weights 1, 3 →
// (80 + 270)/4 = 87.5 (hand-checked).
func TestWeightedMean_worked(t *testing.T) {
	if got, err := WeightedMean([]float64{80, 90}, []float64{1, 3}); err != nil || got != 87.5 {
		t.Errorf("WeightedMean = %v, %v; want 87.5", got, err)
	}
}

func TestWeightedMean_cases(t *testing.T) {
	m := math.MaxFloat64
	for _, tt := range []struct {
		name            string
		values, weights []float64
		want            float64
	}{
		{"one element", []float64{-3}, []float64{0.5}, -3},
		{"equal weights is the mean", []float64{1, 2, 3, 4}, []float64{2, 2, 2, 2}, 2.5},
		{"zero weight ignored", []float64{1, 100}, []float64{1, 0}, 1},
		{"fractional weights", []float64{10, 20}, []float64{0.25, 0.75}, 17.5},
		{"products overflow", []float64{m, m / 2}, []float64{4, 4}, 0.75 * m},
		{"weights overflow in sum", []float64{1, 3}, []float64{m, m}, 2},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// Tolerance 1e-15: a few roundings.
			if got, err := WeightedMean(tt.values, tt.weights); err != nil || !near(got, tt.want, 1e-15) {
				t.Errorf("WeightedMean(%v, %v) = %v, %v; want %v", tt.values, tt.weights, got, err, tt.want)
			}
		})
	}
}

func TestWeightedMean_errors(t *testing.T) {
	for _, tt := range []struct {
		name            string
		values, weights []float64
		want            error
	}{
		{"empty", nil, nil, ErrEmpty},
		{"length mismatch", []float64{1, 2}, []float64{1}, ErrLengthMismatch},
		{"NaN value", []float64{math.NaN()}, []float64{1}, ErrNotFinite},
		{"Inf weight", []float64{1}, []float64{math.Inf(1)}, ErrNotFinite},
		{"negative weight", []float64{1, 2}, []float64{1, -1}, ErrNegative},
		{"all weights zero", []float64{1, 2}, []float64{0, 0}, ErrZeroTotal},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := WeightedMean(tt.values, tt.weights); !errors.Is(err, tt.want) || got != 0 {
				t.Errorf("WeightedMean = %v, %v; want 0, %v", got, err, tt.want)
			}
		})
	}
}

// With integer weights the weighted mean equals the arithmetic mean of the values repeated.
func TestWeightedMean_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		n := 1 + rng.IntN(100)
		values := randomFloats(rng, n)
		weights := make([]float64, n)
		var expanded []float64
		for i := range weights {
			k := rng.IntN(5)
			weights[i] = float64(k)
			for range k {
				expanded = append(expanded, values[i])
			}
		}
		if len(expanded) == 0 {
			continue
		}
		want, _ := ArithmeticMean(expanded)
		// Tolerance 1e-14: both compensated; the products w·x round once each.
		if got, err := WeightedMean(values, weights); err != nil || !near(got, want, 1e-14) {
			t.Fatalf("trial %d: WeightedMean = %v, %v; want %v", trial, got, err, want)
		}
	}
}

func TestWeightedMean_noAllocation(t *testing.T) {
	w := make([]float64, len(michelso))
	for i := range w {
		w[i] = float64(i % 3)
	}
	if a := testing.AllocsPerRun(100, func() { benchSink, _ = WeightedMean(michelso, w) }); a != 0 {
		t.Errorf("WeightedMean allocates %v times per call, want 0", a)
	}
}

func BenchmarkWeightedMean(b *testing.B) {
	for _, size := range benchSizes {
		inputs := benchInputs(size.n)
		weights := make([]float64, size.n)
		for i := range weights {
			weights[i] = float64(1 + i%5)
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; b.Loop(); i++ {
				benchSink, _ = WeightedMean(inputs[i&7], weights)
			}
		})
	}
}

// Tiny products. Values 2^-1000·(1, 3) with weights 2^-1000·(1, 1): products 2^-2000 vanish, but the
// mean is (1 + 3)/2·2^-1000 = 2^-999 exactly (hand). Values 2^-600·(1, 2) with weights 2^-1000·(1, 3):
// (1 + 6)/4·2^-600 = 1.75·2^-600 (hand). Before the fix both returned 0.
func TestWeightedMean_tinyProducts(t *testing.T) {
	if got, err := WeightedMean([]float64{0x1p-1000, 0x3p-1000}, []float64{0x1p-1000, 0x1p-1000}); err != nil || got != 0x1p-999 {
		t.Errorf("WeightedMean(2^-1000·(1 3), 2^-1000·(1 1)) = %v, %v; want 2^-999", got, err)
	}
	// Tolerance 4e-16 relative: the weight 1/3 after normalising is rounded once.
	if got, err := WeightedMean([]float64{0x1p-600, 0x2p-600}, []float64{0x1p-1000, 0x3p-1000}); err != nil || !near(got, 1.75*0x1p-600, 4e-16) {
		t.Errorf("WeightedMean(2^-600·(1 2), 2^-1000·(1 3)) = %v, %v; want 1.75·2^-600", got, err)
	}
}

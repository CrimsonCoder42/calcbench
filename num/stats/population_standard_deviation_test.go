package stats

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// Worked example: 2, 4, 4, 4, 5, 5, 7, 9: mean 5, squared deviations sum to 32, σ² = 32/8 = 4, σ = 2
// (hand-checked from the definition, Hogg, McKean, Craig §1.9).
func TestPopulationStandardDeviation_worked(t *testing.T) {
	if got, err := PopulationStandardDeviation([]float64{2, 4, 4, 4, 5, 5, 7, 9}); err != nil || got != 2 {
		t.Errorf("PopulationStandardDeviation(textbook) = %v, %v; want 2", got, err)
	}
}

func TestPopulationStandardDeviation_cases(t *testing.T) {
	m := math.MaxFloat64
	tests := []struct {
		name   string
		values []float64
		want   float64
	}{
		{"one element", []float64{-8}, 0},
		{"two elements", []float64{1, 3}, 1},
		{"NumAcc1", numAcc1, math.Sqrt(2.0 / 3)},
		// σ² = m² overflows, σ = m does not.
		{"variance overflows, SD finite", []float64{-m, m}, m},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Tolerance 1e-15: a few roundings.
			if got, err := PopulationStandardDeviation(tt.values); err != nil || !near(got, tt.want, 1e-15) {
				t.Errorf("PopulationStandardDeviation(%v) = %v, %v; want %v", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestPopulationStandardDeviation_errors(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   error
	}{
		{"empty", nil, ErrEmpty},
		{"NaN", []float64{1, 2, math.NaN()}, ErrNotFinite},
		{"-Inf", []float64{math.Inf(-1)}, ErrNotFinite},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := PopulationStandardDeviation(tt.values); !errors.Is(err, tt.want) || got != 0 {
				t.Errorf("PopulationStandardDeviation(%v) = %v, %v; want 0, %v", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestPopulationStandardDeviation_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		n := 1 + rng.IntN(400)
		values := randomFloats(rng, n)
		want := math.Sqrt(bigSumOfSquares(values) / float64(n))
		// Tolerance 1e-13 (absolute when n = 1, where σ = 0).
		if got, err := PopulationStandardDeviation(values); err != nil || !near(got, want, 1e-13) {
			t.Fatalf("trial %d: PopulationStandardDeviation(n=%d) = %v, %v; oracle %v", trial, n, got, err, want)
		}
	}
}

func TestPopulationStandardDeviation_noAllocation(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { benchSink, _ = PopulationStandardDeviation(michelso) }); a != 0 {
		t.Errorf("PopulationStandardDeviation allocates %v times per call, want 0", a)
	}
}

func BenchmarkPopulationStandardDeviation(b *testing.B) {
	for _, size := range benchSizes {
		inputs := benchInputs(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; b.Loop(); i++ {
				benchSink, _ = PopulationStandardDeviation(inputs[i&7])
			}
		})
	}
}

// Tiny spread: the hand-checked 2, 4, 4, 4, 5, 5, 7, 9 → 2, scaled by 2^-1000 (exact), gives 2·2^-1000.
// Before the underflow fix the squared deviations (~1e-602) vanished and this returned 0.
func TestPopulationStandardDeviation_tinySpread(t *testing.T) {
	values := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	for i := range values {
		values[i] = math.Ldexp(values[i], -1000)
	}
	if got, err := PopulationStandardDeviation(values); err != nil || got != 0x2p-1000 {
		t.Errorf("PopulationStandardDeviation(2^-1000·(2 4 4 4 5 5 7 9)) = %v, %v; want %v", got, err, 0x2p-1000)
	}
}

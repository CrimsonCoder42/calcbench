package stats

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// Worked example: NumAcc1 (10000001, 10000003, 10000002): deviations −1, 1, 0, so σ² = 2/3 (hand-checked
// from the definition σ² = Σ(x − μ)²/N, Hogg, McKean, Craig §1.9 applied to N equally likely values).
func TestPopulationVariance_worked(t *testing.T) {
	// Tolerance 1e-15: 2/3 is not a float64; one rounding.
	if got, err := PopulationVariance(numAcc1); err != nil || !near(got, 2.0/3, 1e-15) {
		t.Errorf("PopulationVariance(NumAcc1) = %v, %v; want 2/3", got, err)
	}
}

func TestPopulationVariance_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   float64
	}{
		{"one element", []float64{5}, 0},
		{"two elements", []float64{1, 3}, 1},
		{"textbook", []float64{2, 4, 4, 4, 5, 5, 7, 9}, 4},
		{"negative", []float64{-3, -1}, 1},
		{"large but finite", []float64{-1e153, 1e153}, 1e306},
		{"Michelso, from the certified SD", michelso, 0.0790105478190518 * 0.0790105478190518 * 99 / 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Tolerance 1e-13: a few roundings; the certified SD has 15 digits.
			if got, err := PopulationVariance(tt.values); err != nil || !near(got, tt.want, 1e-13) {
				t.Errorf("PopulationVariance(%v) = %v, %v; want %v", tt.name, got, err, tt.want)
			}
		})
	}
}

func TestPopulationVariance_errors(t *testing.T) {
	m := math.MaxFloat64
	tests := []struct {
		name   string
		values []float64
		want   error
	}{
		{"empty", nil, ErrEmpty},
		{"NaN", []float64{math.NaN()}, ErrNotFinite},
		{"Inf", []float64{0, math.Inf(1)}, ErrNotFinite},
		{"overflow", []float64{-m, m}, ErrOverflow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := PopulationVariance(tt.values); !errors.Is(err, tt.want) || got != 0 {
				t.Errorf("PopulationVariance(%v) = %v, %v; want 0, %v", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestPopulationVariance_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		n := 1 + rng.IntN(400)
		values := randomFloats(rng, n)
		want := bigSumOfSquares(values) / float64(n)
		// Tolerance 1e-13 (absolute 1e-13 when n = 1, where the variance is 0).
		if got, err := PopulationVariance(values); err != nil || !near(got, want, 1e-13) {
			t.Fatalf("trial %d: PopulationVariance(n=%d) = %v, %v; oracle %v", trial, n, got, err, want)
		}
	}
}

func TestPopulationVariance_noAllocation(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { benchSink, _ = PopulationVariance(michelso) }); a != 0 {
		t.Errorf("PopulationVariance allocates %v times per call, want 0", a)
	}
}

func BenchmarkPopulationVariance(b *testing.B) {
	for _, size := range benchSizes {
		inputs := benchInputs(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; b.Loop(); i++ {
				benchSink, _ = PopulationVariance(inputs[i&7])
			}
		})
	}
}

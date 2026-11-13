package stats

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// Worked example: e-Handbook §1.3.5.6, AAD = Σ|Y − Ȳ|/N. NumAcc1: deviations 1, 1, 0 → 2/3. The 12
// resistivity values of §7.2.6.2 → 0.0358763888… (hand-checked in exact decimals).
func TestMeanAbsoluteDeviation_worked(t *testing.T) {
	// Tolerance 1e-15 for NumAcc1 (one rounding of 2/3); 1e-12 for resistivity, whose deviations of
	// ~0.04 from values near 95 carry the inputs' decimal rounding (ulp 1.4e-14).
	if got, err := MeanAbsoluteDeviation(numAcc1); err != nil || !near(got, 2.0/3, 1e-15) {
		t.Errorf("MeanAbsoluteDeviation(NumAcc1) = %v, %v; want 2/3", got, err)
	}
	if got, err := MeanAbsoluteDeviation(resistivity); err != nil || !near(got, 0.035876388888888888, 1e-12) {
		t.Errorf("MeanAbsoluteDeviation(resistivity) = %v, %v; want 0.0358763889", got, err)
	}
}

func TestMeanAbsoluteDeviation_cases(t *testing.T) {
	m := math.MaxFloat64
	tests := []struct {
		name   string
		values []float64
		want   float64
	}{
		{"one element", []float64{3}, 0},
		{"two elements", []float64{1, 3}, 1},
		{"negative and zero", []float64{-2, 0, 2}, 4.0 / 3},
		{"deviations overflow, result finite", []float64{-m, m}, m},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Tolerance 1e-15: a few roundings.
			if got, err := MeanAbsoluteDeviation(tt.values); err != nil || !near(got, tt.want, 1e-15) {
				t.Errorf("MeanAbsoluteDeviation(%v) = %v, %v; want %v", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestMeanAbsoluteDeviation_errors(t *testing.T) {
	for _, tt := range []struct {
		name   string
		values []float64
		want   error
	}{
		{"empty", nil, ErrEmpty},
		{"NaN", []float64{math.NaN()}, ErrNotFinite},
		{"Inf", []float64{1, math.Inf(-1)}, ErrNotFinite},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := MeanAbsoluteDeviation(tt.values); !errors.Is(err, tt.want) || got != 0 {
				t.Errorf("MeanAbsoluteDeviation(%v) = %v, %v; want 0, %v", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestMeanAbsoluteDeviation_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		n := 1 + rng.IntN(400)
		values := randomFloats(rng, n)
		mean, _ := bigMean(values).Float64()
		var want float64 // naive two-pass in float64 with the oracle mean
		for _, v := range values {
			want += math.Abs(v - mean)
		}
		want /= float64(n)
		// Tolerance 1e-12 relative (absolute when n = 1): the naive oracle loses up to n·u.
		if got, err := MeanAbsoluteDeviation(values); err != nil || !near(got, want, 1e-12) {
			t.Fatalf("trial %d: MeanAbsoluteDeviation(n=%d) = %v, %v; oracle %v", trial, n, got, err, want)
		}
	}
}

func TestMeanAbsoluteDeviation_noAllocation(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { benchSink, _ = MeanAbsoluteDeviation(michelso) }); a != 0 {
		t.Errorf("MeanAbsoluteDeviation allocates %v times per call, want 0", a)
	}
}

func BenchmarkMeanAbsoluteDeviation(b *testing.B) {
	for _, size := range benchSizes {
		inputs := benchInputs(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; b.Loop(); i++ {
				benchSink, _ = MeanAbsoluteDeviation(inputs[i&7])
			}
		})
	}
}

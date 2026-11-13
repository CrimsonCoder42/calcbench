package stats

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Worked example: e-Handbook §1.3.5.1, mid-range = (smallest + largest)/2; on the 12 resistivity values
// of §7.2.6.2, (95.0610 + 95.1990)/2 = 95.1300 (hand-checked).
func TestMidrange_worked(t *testing.T) {
	// Tolerance 1e-15: one rounding of the stored decimals.
	if got, err := Midrange(resistivity); err != nil || !near(got, 95.13, 1e-15) {
		t.Errorf("Midrange(resistivity) = %v, %v; want 95.13", got, err)
	}
}

func TestMidrange_cases(t *testing.T) {
	m := math.MaxFloat64
	for _, tt := range []struct {
		name   string
		values []float64
		want   float64
	}{
		{"one element", []float64{5}, 5},
		{"two elements", []float64{3, 9}, 6},
		{"negative and zero", []float64{-4, 0, -1}, -2},
		{"max of the type", []float64{m, m}, m},
		{"opposite extremes", []float64{-m, m}, 0},
		{"lowest of the type", []float64{-m, -m / 2}, -0.75 * m},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Midrange(tt.values); err != nil || got != tt.want {
				t.Errorf("Midrange(%v) = %v, %v; want %v", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestMidrange_errors(t *testing.T) {
	for _, tt := range []struct {
		name   string
		values []float64
		want   error
	}{
		{"empty", []float64{}, ErrEmpty},
		{"Inf", []float64{math.Inf(1)}, ErrNotFinite},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Midrange(tt.values); !errors.Is(err, tt.want) || got != 0 {
				t.Errorf("Midrange(%v) = %v, %v; want 0, %v", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestMidrange_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		values := randomFloats(rng, 1+rng.IntN(300))
		want := (slices.Max(values) + slices.Min(values)) / 2
		// Tolerance 1e-15: the oracle rounds the sum, Midrange the difference.
		if got, err := Midrange(values); err != nil || !near(got, want, 1e-15) {
			t.Fatalf("trial %d: Midrange = %v, %v; want %v", trial, got, err, want)
		}
	}
}

func TestMidrange_noAllocation(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { benchSink, _ = Midrange(michelso) }); a != 0 {
		t.Errorf("Midrange allocates %v times per call, want 0", a)
	}
}

func BenchmarkMidrange(b *testing.B) {
	for _, size := range benchSizes {
		inputs := benchInputs(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; b.Loop(); i++ {
				benchSink, _ = Midrange(inputs[i&7])
			}
		})
	}
}

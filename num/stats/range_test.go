package stats

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Worked example: e-Handbook §1.3.5.6 defines the range as largest minus smallest; on the 12 resistivity
// values of §7.2.6.2 that is 95.1990 − 95.0610 = 0.1380 (hand-checked).
func TestRange_worked(t *testing.T) {
	// Exact: the result is the float64 difference of the two stored extremes. The extremes are float64
	// variables so the subtraction is done in float64, not as an exact untyped constant (which is 0.138).
	hi, lo := 95.1990, 95.0610
	if got, err := Range(resistivity); err != nil || got != hi-lo {
		t.Errorf("Range(resistivity) = %v, %v; want %v", got, err, hi-lo)
	}
}

func TestRange_cases(t *testing.T) {
	for _, tt := range []struct {
		name   string
		values []float64
		want   float64
	}{
		{"one element", []float64{5}, 0},
		{"two elements", []float64{9, 3}, 6},
		{"duplicates", []float64{2, 2, 2}, 0},
		{"negative", []float64{-5, -1, -3}, 4},
		{"reversed", []float64{5, 4, 3, 2, 1}, 4},
		{"max of the type", []float64{math.MaxFloat64, 0}, math.MaxFloat64},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Range(tt.values); err != nil || got != tt.want {
				t.Errorf("Range(%v) = %v, %v; want %v", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestRange_errors(t *testing.T) {
	for _, tt := range []struct {
		name   string
		values []float64
		want   error
	}{
		{"empty", nil, ErrEmpty},
		{"NaN", []float64{1, math.NaN()}, ErrNotFinite},
		{"overflow", []float64{-math.MaxFloat64, math.MaxFloat64}, ErrOverflow},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Range(tt.values); !errors.Is(err, tt.want) || got != 0 {
				t.Errorf("Range(%v) = %v, %v; want 0, %v", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestRange_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		values := randomFloats(rng, 1+rng.IntN(300))
		want := slices.Max(values) - slices.Min(values)
		if got, err := Range(values); err != nil || got != want {
			t.Fatalf("trial %d: Range = %v, %v; want %v", trial, got, err, want)
		}
	}
}

func TestRange_noAllocation(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { benchSink, _ = Range(michelso) }); a != 0 {
		t.Errorf("Range allocates %v times per call, want 0", a)
	}
}

func BenchmarkRange(b *testing.B) {
	for _, size := range benchSizes {
		inputs := benchInputs(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; b.Loop(); i++ {
				benchSink, _ = Range(inputs[i&7])
			}
		})
	}
}

package stats

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Worked example: e-Handbook §1.3.5.1 median, even N: (Y(N/2) + Y(N/2+1))/2. The 12 resistivity values
// of §7.2.6.2 sorted have 95.1567 and 95.1591 in the middle → 95.1579 (hand-checked).
func TestMedian_worked(t *testing.T) {
	// Tolerance 1e-15: one rounding of the average of two stored decimals.
	if got, err := Median(resistivity); err != nil || !near(got, 95.1579, 1e-15) {
		t.Errorf("Median(resistivity) = %v, %v; want 95.1579", got, err)
	}
	if got, err := Median(numAcc1); err != nil || got != 10000002 {
		t.Errorf("Median(NumAcc1) = %v, %v; want 10000002", got, err)
	}
}

func TestMedian_cases(t *testing.T) {
	m := math.MaxFloat64
	for _, tt := range []struct {
		name   string
		values []float64
		want   float64
	}{
		{"one element", []float64{4}, 4},
		{"two elements", []float64{4, 1}, 2.5},
		{"odd", []float64{9, 1, 5}, 5},
		{"duplicates", []float64{2, 2, 1, 2}, 2},
		{"sorted", []float64{1, 2, 3, 4}, 2.5},
		{"reversed", []float64{4, 3, 2, 1}, 2.5},
		{"negative and zero", []float64{-3, 0, -1, -2}, -1.5},
		{"extremes", []float64{-m, m}, 0},
		{"max of the type", []float64{m, m}, m},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Median(tt.values); err != nil || got != tt.want {
				t.Errorf("Median(%v) = %v, %v; want %v", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestMedian_errors(t *testing.T) {
	for _, tt := range []struct {
		name   string
		values []float64
		want   error
	}{
		{"empty", nil, ErrEmpty},
		{"NaN", []float64{1, math.NaN(), 2}, ErrNotFinite},
		{"Inf", []float64{math.Inf(-1)}, ErrNotFinite},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Median(tt.values); !errors.Is(err, tt.want) || got != 0 {
				t.Errorf("Median(%v) = %v, %v; want 0, %v", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestMedian_doesNotModifyInput(t *testing.T) {
	values := slices.Clone(resistivity)
	if _, err := Median(values); err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(values, resistivity) {
		t.Error("Median modified its input")
	}
}

func TestMedian_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		n := 1 + rng.IntN(300)
		values := randomFloats(rng, n)
		s := slices.Sorted(slices.Values(values))
		want := (s[(n-1)/2] + s[n/2]) / 2
		// Tolerance 1e-15: the oracle rounds the sum, Median the difference.
		if got, err := Median(values); err != nil || !near(got, want, 1e-15) {
			t.Fatalf("trial %d: Median(n=%d) = %v, %v; want %v", trial, n, got, err, want)
		}
	}
}

func TestMedian_allocations(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { benchSink, _ = Median(michelso) }); a != 1 {
		t.Errorf("Median allocates %v times per call, want 1 (the copy)", a)
	}
}

func BenchmarkMedian(b *testing.B) {
	for _, size := range benchSizes {
		inputs := benchInputs(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; b.Loop(); i++ {
				benchSink, _ = Median(inputs[i&7])
			}
		})
	}
}

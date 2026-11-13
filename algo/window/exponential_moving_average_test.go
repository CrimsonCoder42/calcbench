package window

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// NIST/SEMATECH e-Handbook §6.4.3.1 "Single Exponential Smoothing": the 12 observations and, for
// α = 0.1, the printed smoothed values S₂ … S₁₂ (S₂ = y₁ = 71). Our result i is S_{i+2}. The handbook
// prints two decimals (70.9 with one), so the tolerance is half a unit in the last printed place.
func TestExponentialMovingAverage_nistExample(t *testing.T) {
	y := []float64{71, 70, 69, 68, 64, 65, 72, 78, 75, 75, 75, 70}
	printed := []float64{71, 70.9, 70.71, 70.44, 69.80, 69.32, 69.58, 70.43, 70.88, 71.29, 71.67}
	got, err := ExponentialMovingAverage(nil, y, 0.1)
	if err != nil || len(got) != len(y) {
		t.Fatalf("ExponentialMovingAverage(nist, 0.1) = (%v, %v)", got, err)
	}
	for i, p := range printed {
		if math.Abs(got[i]-p) > 0.005 {
			t.Errorf("result %d (S%d) = %.4f, NIST prints %.2f", i, i+2, got[i], p)
		}
	}
	// Two exact steps, hand-derived: 0.1·70 + 0.9·71 = 70.9; 0.1·69 + 0.9·70.9 = 70.71.
	if math.Abs(got[1]-70.9) > 1e-12 || math.Abs(got[2]-70.71) > 1e-12 {
		t.Errorf("first steps = %v, %v, want 70.9, 70.71", got[1], got[2])
	}
}

func TestExponentialMovingAverage_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		alpha  float64
		want   []float64
	}{
		{"empty", nil, 0.5, nil},
		{"one", []float64{4}, 0.3, []float64{4}},
		{"alpha 1 is the identity", []float64{3, -1, 2}, 1, []float64{3, -1, 2}},
		{"alpha 0.5", []float64{0, 4, 8}, 0.5, []float64{0, 2, 5}},
		{"constant stays constant", []float64{7, 7, 7}, 0.25, []float64{7, 7, 7}},
		{"negative and zero", []float64{-4, 0}, 0.25, []float64{-4, -3}},
		{"largest finite values", []float64{math.MaxFloat64, math.MaxFloat64}, 0.5, []float64{math.MaxFloat64, math.MaxFloat64}},
		{"opposite extremes", []float64{-math.MaxFloat64, math.MaxFloat64}, 0.5, []float64{-math.MaxFloat64, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExponentialMovingAverage(nil, tt.values, tt.alpha)
			if err != nil || !slices.Equal(got, tt.want) {
				t.Errorf("ExponentialMovingAverage(%v, %v) = (%v, %v), want (%v, nil)", tt.values, tt.alpha, got, err, tt.want)
			}
		})
	}
}

func TestExponentialMovingAverage_errors(t *testing.T) {
	tests := []struct {
		name    string
		values  []float64
		alpha   float64
		wantErr error
	}{
		{"alpha zero", []float64{1}, 0, ErrAlpha},
		{"alpha negative", []float64{1}, -0.1, ErrAlpha},
		{"alpha above one", []float64{1}, 1.0000001, ErrAlpha},
		{"alpha NaN", []float64{1}, math.NaN(), ErrAlpha},
		{"alpha +Inf", []float64{1}, math.Inf(1), ErrAlpha},
		{"NaN value", []float64{1, math.NaN()}, 0.5, ErrNotFinite},
		{"-Inf value", []float64{math.Inf(-1)}, 0.5, ErrNotFinite},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := []float64{9}
			got, err := ExponentialMovingAverage(dst, tt.values, tt.alpha)
			if !errors.Is(err, tt.wantErr) || !slices.Equal(got, []float64{9}) {
				t.Errorf("ExponentialMovingAverage([9], %v, %v) = (%v, %v), want ([9], %v)", tt.values, tt.alpha, got, err, tt.wantErr)
			}
		})
	}
}

func TestExponentialMovingAverage_appendsToDst(t *testing.T) {
	got, err := ExponentialMovingAverage([]float64{-1}, []float64{2, 4}, 0.5)
	if want := []float64{-1, 2, 3}; err != nil || !slices.Equal(got, want) {
		t.Errorf("ExponentialMovingAverage([-1], [2 4], 0.5) = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestExponentialMovingAverage_noAllocationWithBuffer(t *testing.T) {
	values := benchFloats(1000)
	buf := make([]float64, 0, len(values))
	allocs := testing.AllocsPerRun(20, func() { buf, _ = ExponentialMovingAverage(buf[:0], values, 0.1) })
	if allocs != 0 {
		t.Errorf("ExponentialMovingAverage with a reused buffer allocates %v times per call, want 0", allocs)
	}
}

// Oracle: the closed form Sₜ = α Σ (1−α)^j y_{t−j} + (1−α)^t y₀, obtained by expanding the recursion
// of NIST §6.4.3.1, summed directly. Both are sums of at most 31 terms of size ≤ 1000, each rounding at most
// a few ulps, so 1e-9 is a safe bound.
func TestExponentialMovingAverage_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		values := make([]float64, 1+rng.IntN(30))
		for i := range values {
			values[i] = 2000*rng.Float64() - 1000
		}
		before := slices.Clone(values)
		alpha := 0.01 + 0.99*rng.Float64()
		got, err := ExponentialMovingAverage(nil, values, alpha)
		if err != nil || len(got) != len(values) {
			t.Fatalf("trial %d: error %v, %d results", trial, err, len(got))
		}
		for t2 := range values {
			want := math.Pow(1-alpha, float64(t2)) * values[0]
			for j := 0; j < t2; j++ {
				want += alpha * math.Pow(1-alpha, float64(j)) * values[t2-j]
			}
			if math.Abs(got[t2]-want) > 1e-9 {
				t.Fatalf("trial %d: result %d = %v, closed form %v", trial, t2, got[t2], want)
			}
		}
		if !slices.Equal(values, before) {
			t.Fatalf("trial %d: input modified", trial)
		}
	}
}

func BenchmarkExponentialMovingAverage(b *testing.B) {
	for _, size := range benchSizes {
		values := benchFloats(size.n)
		buf := make([]float64, 0, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				buf, _ = ExponentialMovingAverage(buf[:0], values, 0.1)
			}
		})
	}
}

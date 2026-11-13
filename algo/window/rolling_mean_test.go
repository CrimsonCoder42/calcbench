package window

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// NIST/SEMATECH e-Handbook of Statistical Methods, §6.4.2.1 "Single Moving Average": N = 3 over the
// supplier data. The handbook prints 8.667, 9.667, 10.000, 11.000, 10.667, 10.000, 10.333, 9.667,
// 11.000, 10.000 (three decimals). The exact values are the hand-derived sums over 3.
// Tolerances: 1e-15 relative to the exact fraction (one rounding of a sum of small integers, then one
// division); 0.0005 to the printed value (half a unit in its last printed digit).
func TestRollingMean_nistExample(t *testing.T) {
	printed := []float64{8.667, 9.667, 10.000, 11.000, 10.667, 10.000, 10.333, 9.667, 11.000, 10.000}
	exact := []float64{26.0 / 3, 29.0 / 3, 10, 11, 32.0 / 3, 10, 31.0 / 3, 29.0 / 3, 11, 10}
	got, err := RollingMean(nil, nistSupplier, 3)
	if err != nil || len(got) != len(exact) {
		t.Fatalf("RollingMean(nistSupplier, 3) = (%v, %v), want %d means", got, err, len(exact))
	}
	for i := range got {
		if math.Abs(got[i]-exact[i]) > 1e-15*exact[i] {
			t.Errorf("mean %d = %v, want %v (exact)", i, got[i], exact[i])
		}
		if math.Abs(got[i]-printed[i]) > 0.0005 {
			t.Errorf("mean %d = %v, NIST prints %v", i, got[i], printed[i])
		}
	}
}

func TestRollingMean_cases(t *testing.T) {
	tests := []struct {
		name      string
		values    []float64
		windowLen int
		want      []float64
	}{
		{"empty", nil, 1, nil},
		{"window longer than the data", []float64{1, 2}, 3, nil},
		{"one element", []float64{4}, 1, []float64{4}},
		{"window equals the length", []float64{3, -1, 4}, 3, []float64{2}},
		{"two windows", []float64{1, 3, 5}, 2, []float64{2, 4}},
		{"negative and zero", []float64{-2, 0, -4}, 2, []float64{-1, -2}},
		{"large value leaves the window", []float64{1e16, 1, 1, 1}, 2, []float64{5e15, 1, 1}},
		{"largest finite value, window 1", []float64{math.MaxFloat64}, 1, []float64{math.MaxFloat64}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RollingMean(nil, tt.values, tt.windowLen)
			if err != nil || len(got) != len(tt.want) || !slices.Equal(got, tt.want) {
				t.Errorf("RollingMean(%v, %d) = (%v, %v), want (%v, nil)", tt.values, tt.windowLen, got, err, tt.want)
			}
		})
	}
}

func TestRollingMean_errors(t *testing.T) {
	tests := []struct {
		name      string
		values    []float64
		windowLen int
		wantErr   error
	}{
		{"window zero", []float64{1}, 0, ErrWindowLen},
		{"window negative", []float64{1}, -3, ErrWindowLen},
		{"NaN", []float64{math.NaN()}, 1, ErrNotFinite},
		{"-Inf", []float64{2, math.Inf(-1)}, 1, ErrNotFinite},
		// A-007: the sum overflows although the mean, MaxFloat64, would fit.
		{"sum overflows, mean would fit", []float64{math.MaxFloat64, math.MaxFloat64}, 2, ErrOverflow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RollingMean(nil, tt.values, tt.windowLen)
			if !errors.Is(err, tt.wantErr) || len(got) != 0 {
				t.Errorf("RollingMean(%v, %d) = (%v, %v), want ([], %v)", tt.values, tt.windowLen, got, err, tt.wantErr)
			}
		})
	}
}

func TestRollingMean_noAllocationWithBuffer(t *testing.T) {
	values := benchFloats(1000)
	buf := make([]float64, 0, len(values))
	allocs := testing.AllocsPerRun(20, func() { buf, _ = RollingMean(buf[:0], values, benchWindow) })
	if allocs != 0 {
		t.Errorf("RollingMean with a reused buffer allocates %v times per call, want 0", allocs)
	}
}

// Oracle: the exact window sum divided by w. Tolerance: sumTolerance over w, plus the roundings of
// the division in the code and in the oracle (2ε·|mean|).
func TestRollingMean_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		n := 1 + rng.IntN(200)
		values := mixedMagnitudes(rng, n)
		w := 1 + rng.IntN(n)
		got, err := RollingMean(nil, values, w)
		if err != nil || len(got) != n-w+1 {
			t.Fatalf("trial %d: RollingMean(n=%d, w=%d) = %d means, %v", trial, n, w, len(got), err)
		}
		for i, g := range got {
			exactSum := exactWindowSum(values[i : i+w])
			want := exactSum / float64(w)
			if tol := sumTolerance(values, i, w, exactSum)/float64(w) + 2*eps*math.Abs(want); math.Abs(g-want) > tol {
				t.Fatalf("trial %d: window %d (w=%d) mean = %v, want %v, error %g > tolerance %g", trial, i, w, g, want, math.Abs(g-want), tol)
			}
		}
	}
}

func BenchmarkRollingMean(b *testing.B) {
	for _, size := range benchSizes {
		values := benchFloats(size.n)
		buf := make([]float64, 0, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				buf, _ = RollingMean(buf[:0], values, benchWindow)
			}
		})
	}
}

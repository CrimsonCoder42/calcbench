package stats

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// strdSD lists the StRD Univariate datasets with their certified sample SD (divisor n − 1) and the
// tolerance the float64 inputs allow. NumAcc3 and NumAcc4 store decimals near 1e6 and 1e7, whose ulps
// (1.2e-10, 1.9e-9) shift the ±0.1 deviations: the exact SD of the stored inputs differs from the
// certified 0.1 by 3.5e-10 and 5.6e-9 relative (computed in exact rationals), so 1e-8 is allowed there.
var strdSD = []struct {
	name   string
	values []float64
	sd     float64
	tol    float64
}{
	{"NumAcc1", numAcc1, 1, 1e-15},
	{"NumAcc2", numAcc2, 0.1, 1e-14},
	{"NumAcc3", numAcc3, 0.1, 1e-8},
	{"NumAcc4", numAcc4, 0.1, 1e-8},
	{"Michelso", michelso, 0.0790105478190518, 1e-13},
	{"Mavro", mavro, 0.000429123454003053, 1e-12},
	{"Lew", lew, 277.332168044316, 1e-14},
}

// Worked example: StRD certified sample SDs, squared (NumAcc4 → 0.01).
func TestSampleVariance_strd(t *testing.T) {
	for _, tt := range strdSD {
		t.Run(tt.name, func(t *testing.T) {
			// Tolerance: twice the SD's (a square doubles the relative error).
			if got, err := SampleVariance(tt.values); err != nil || !near(got, tt.sd*tt.sd, 2*tt.tol+1e-15) {
				t.Errorf("SampleVariance(%s) = %v, %v; want %v", tt.name, got, err, tt.sd*tt.sd)
			}
		})
	}
}

func TestSampleVariance_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		want   float64
	}{
		{"two elements", []float64{1, 3}, 2},
		{"all equal", []float64{-4, -4, -4}, 0},
		{"negative and zero", []float64{-1, 0, 1}, 1},
		{"textbook", []float64{2, 4, 4, 4, 5, 5, 7, 9}, 32.0 / 7},
		{"large but finite", []float64{-1e153, 1e153}, 2e306},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Tolerance 1e-15: a few roundings.
			if got, err := SampleVariance(tt.values); err != nil || !near(got, tt.want, 1e-15) {
				t.Errorf("SampleVariance(%v) = %v, %v; want %v", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestSampleVariance_errors(t *testing.T) {
	m := math.MaxFloat64
	tests := []struct {
		name   string
		values []float64
		want   error
	}{
		{"empty", nil, ErrEmpty},
		{"one element", []float64{1}, ErrTooFew},
		{"NaN", []float64{1, math.NaN()}, ErrNotFinite},
		{"Inf", []float64{math.Inf(-1), 1}, ErrNotFinite},
		{"overflow", []float64{-m, m}, ErrOverflow},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := SampleVariance(tt.values); !errors.Is(err, tt.want) || got != 0 {
				t.Errorf("SampleVariance(%v) = %v, %v; want 0, %v", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestSampleVariance_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 300 {
		n := 2 + rng.IntN(400)
		values := randomFloats(rng, n)
		want := bigSumOfSquares(values) / float64(n-1)
		// Tolerance 1e-13: O(n·u) for the corrected two-pass sum, n ≤ 401.
		if got, err := SampleVariance(values); err != nil || !near(got, want, 1e-13) {
			t.Fatalf("trial %d: SampleVariance(n=%d) = %v, %v; oracle %v", trial, n, got, err, want)
		}
	}
}

func TestSampleVariance_noAllocation(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { benchSink, _ = SampleVariance(michelso) }); a != 0 {
		t.Errorf("SampleVariance allocates %v times per call, want 0", a)
	}
}

func BenchmarkSampleVariance(b *testing.B) {
	for _, size := range benchSizes {
		inputs := benchInputs(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; b.Loop(); i++ {
				benchSink, _ = SampleVariance(inputs[i&7])
			}
		})
	}
}

package stats

import (
	"errors"
	"math"
	"testing"
)

// Worked example: Dataplot "ISO 13528 ZSCORE", Z = (X − Xref)/sigma: x = 12, mean 10, sd 2 → 1
// (hand-checked).
func TestZScore_worked(t *testing.T) {
	if got, err := ZScore(12, 10, 2); err != nil || got != 1 {
		t.Errorf("ZScore(12, 10, 2) = %v, %v; want 1", got, err)
	}
}

func TestZScore_cases(t *testing.T) {
	m := math.MaxFloat64
	for _, tt := range []struct {
		name        string
		x, mean, sd float64
		want        float64
	}{
		{"at the mean", 10, 10, 3, 0},
		{"below", 4, 10, 3, -2},
		{"tiny sd", 1e-300, 0, 1e-300, 1},
		{"difference overflows, z finite", m, -m, m, 2},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := ZScore(tt.x, tt.mean, tt.sd); err != nil || got != tt.want {
				t.Errorf("ZScore(%v, %v, %v) = %v, %v; want %v", tt.x, tt.mean, tt.sd, got, err, tt.want)
			}
		})
	}
}

func TestZScore_errors(t *testing.T) {
	for _, tt := range []struct {
		name        string
		x, mean, sd float64
		want        error
	}{
		{"NaN x", math.NaN(), 0, 1, ErrNotFinite},
		{"Inf mean", 0, math.Inf(1), 1, ErrNotFinite},
		{"Inf sd", 0, 0, math.Inf(1), ErrNotFinite},
		{"negative sd", 1, 0, -2, ErrNegative},
		{"zero sd", 1, 0, 0, ErrZeroVariance},
		{"overflow", 1e300, 0, 1e-10, ErrOverflow},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := ZScore(tt.x, tt.mean, tt.sd); !errors.Is(err, tt.want) || got != 0 {
				t.Errorf("ZScore = %v, %v; want 0, %v", got, err, tt.want)
			}
		})
	}
}

// RawScoreFromZScore inverts ZScore.
func TestZScore_roundTrip(t *testing.T) {
	for _, x := range []float64{-7.5, 0, 3, 1e6} {
		z, err := ZScore(x, 2.5, 1.25)
		back, err2 := RawScoreFromZScore(z, 2.5, 1.25)
		// Tolerance 1e-15 relative plus 1e-15 absolute: two roundings.
		if err != nil || err2 != nil || math.Abs(back-x) > 1e-15*(1+math.Abs(x)) {
			t.Errorf("RawScoreFromZScore(ZScore(%v)) = %v, %v, %v", x, back, err, err2)
		}
	}
}

func TestZScore_noAllocation(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { benchSink, _ = ZScore(12, 10, 2) }); a != 0 {
		t.Errorf("ZScore allocates %v times per call, want 0", a)
	}
}

func BenchmarkZScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchSink, _ = ZScore(12, 10, 2)
	}
}

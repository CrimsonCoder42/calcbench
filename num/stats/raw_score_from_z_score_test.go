package stats

import (
	"errors"
	"math"
	"testing"
)

// Worked example: the inverse of the ZScore example, Z = (X − Xref)/sigma solved for X: z = 1, mean
// 10, sd 2 → 12 (hand-checked).
func TestRawScoreFromZScore_worked(t *testing.T) {
	if got, err := RawScoreFromZScore(1, 10, 2); err != nil || got != 12 {
		t.Errorf("RawScoreFromZScore(1, 10, 2) = %v, %v; want 12", got, err)
	}
}

func TestRawScoreFromZScore_cases(t *testing.T) {
	for _, tt := range []struct {
		name        string
		z, mean, sd float64
		want        float64
	}{
		{"z zero", 0, 7, 3, 7},
		{"negative z", -2, 10, 3, 4},
		{"sd zero", 5, 7, 0, 7},
		{"one rounding", 0.1, 0, 3, 0.30000000000000004},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := RawScoreFromZScore(tt.z, tt.mean, tt.sd); err != nil || got != tt.want {
				t.Errorf("RawScoreFromZScore(%v, %v, %v) = %v, %v; want %v", tt.z, tt.mean, tt.sd, got, err, tt.want)
			}
		})
	}
}

func TestRawScoreFromZScore_errors(t *testing.T) {
	for _, tt := range []struct {
		name        string
		z, mean, sd float64
		want        error
	}{
		{"NaN z", math.NaN(), 0, 1, ErrNotFinite},
		{"Inf mean", 0, math.Inf(-1), 1, ErrNotFinite},
		{"NaN sd", 0, 0, math.NaN(), ErrNotFinite},
		{"negative sd", 1, 0, -1, ErrNegative},
		{"overflow", 1e300, 0, 1e300, ErrOverflow},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := RawScoreFromZScore(tt.z, tt.mean, tt.sd); !errors.Is(err, tt.want) || got != 0 {
				t.Errorf("RawScoreFromZScore = %v, %v; want 0, %v", got, err, tt.want)
			}
		})
	}
}

func TestRawScoreFromZScore_noAllocation(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { benchSink, _ = RawScoreFromZScore(1, 10, 2) }); a != 0 {
		t.Errorf("RawScoreFromZScore allocates %v times per call, want 0", a)
	}
}

func BenchmarkRawScoreFromZScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchSink, _ = RawScoreFromZScore(1, 10, 2)
	}
}

package prob

import (
	"math"
	"testing"
)

// e-Handbook §1.3.6.7.4: the upper-tail critical value at 0.95 leaves 0.05 above it (table to 3 decimals, so
// tolerance 0.0005 as in TestChiSquareCDF_sourceTable).
func TestChiSquareSurvival_sourceTable(t *testing.T) {
	for _, tc := range chiSquareTable {
		got, err := ChiSquareSurvival(tc.x, tc.df)
		if err != nil {
			t.Fatal(err)
		}
		if math.Abs(got-(1-tc.prob)) > 0.0005+1e-9 {
			t.Errorf("ChiSquareSurvival(%v, %v) = %v, want %v ± 0.0005", tc.x, tc.df, got, 1-tc.prob)
		}
	}
}

// Closed forms in the far upper tail, where 1 − F would be 0: Q(1, z) = e^(−z) (ν = 2), Q(½, z) = erfc √z
// (ν = 1). Tolerance 1e-12 relative (exp of a large argument amplifies its rounding).
func TestChiSquareSurvival_upperTail(t *testing.T) {
	for _, x := range []float64{2, 10, 100, 700, 1400} {
		got, err := ChiSquareSurvival(x, 2)
		if want := math.Exp(-x / 2); err != nil || !within(got, want, 1e-12) {
			t.Errorf("ChiSquareSurvival(%v, 2) = (%v, %v), want %v", x, got, err, want)
		}
		got, err = ChiSquareSurvival(x, 1)
		if want := math.Erfc(math.Sqrt(x / 2)); err != nil || !within(got, want, 1e-12) {
			t.Errorf("ChiSquareSurvival(%v, 1) = (%v, %v), want %v", x, got, err, want)
		}
	}
}

func TestChiSquareSurvival_cases(t *testing.T) {
	tests := []struct {
		name  string
		x, df float64
		want  float64
	}{
		{"zero", 0, 4, 1},
		{"negative", -1, 4, 1},
		{"beyond the float64 range of the tail", 1e5, 2, 0},
	}
	for _, tc := range tests {
		got, err := ChiSquareSurvival(tc.x, tc.df)
		if err != nil || got != tc.want {
			t.Errorf("%s: ChiSquareSurvival(%v, %v) = (%v, %v), want %v", tc.name, tc.x, tc.df, got, err, tc.want)
		}
	}
}

func TestChiSquareSurvival_errors(t *testing.T) {
	_, err := ChiSquareSurvival(nan, 2)
	requireErr(t, "ChiSquareSurvival(NaN, 2)", err, ErrNotFinite)
	_, err = ChiSquareSurvival(-inf, 2)
	requireErr(t, "ChiSquareSurvival(-Inf, 2)", err, ErrNotFinite)
	_, err = ChiSquareSurvival(1, 0)
	requireErr(t, "ChiSquareSurvival(1, 0)", err, ErrOutOfRange)
	_, err = ChiSquareSurvival(1, 2e6)
	requireErr(t, "ChiSquareSurvival(1, 2e6)", err, ErrOutOfRange)
}

func TestChiSquareSurvival_noAlloc(t *testing.T) {
	requireNoAlloc(t, "ChiSquareSurvival", func() { sinkFloat, _ = ChiSquareSurvival(12.592, 6) })
}

func BenchmarkChiSquareSurvival(b *testing.B) {
	for _, df := range []float64{3, 100, 1e4, 1e6} {
		b.Run(fmtDF(df), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, _ = ChiSquareSurvival(1.05*df, df)
			}
		})
	}
}

package prob

import (
	"math"
	"testing"
)

// ν = 1: ½ + atan(t)/π; ν = 2: ½ + t/(2√(2 + t²)) (JKB ch. 28 §2, closed forms; A&S 26.7.3).
func TestStudentTCDF_sourceExample(t *testing.T) {
	for _, x := range []float64{-30, -1, 0, 0.5, 2, 1e5} {
		c1, err := StudentTCDF(x, 1)
		if err != nil || !within(c1, 0.5+math.Atan(x)/math.Pi, 1e-12) {
			t.Errorf("ν = 1, t = %v: (%v, %v)", x, c1, err)
		}
		c2, _ := StudentTCDF(x, 2)
		if !within(c2, 0.5+x/(2*math.Sqrt(2+x*x)), 1e-12) {
			t.Errorf("ν = 2, t = %v: %v", x, c2)
		}
	}
}

// The lower tail keeps relative accuracy: ν = 1, t = −1e10 → atan2(1, 1e10)/π ≈ 1/(π·1e10).
func TestStudentTCDF_cases(t *testing.T) {
	if got, _ := StudentTCDF(-1e10, 1); !within(got, 1/(math.Pi*1e10), 1e-9) {
		t.Errorf("far lower tail: %v", got)
	}
	if got, _ := StudentTCDF(0, 7); got != 0.5 {
		t.Errorf("t = 0: %v", got)
	}
	if got, _ := StudentTCDF(1e200, 3); got != 1 {
		t.Errorf("t huge: %v", got)
	}
}

func TestStudentTCDF_errors(t *testing.T) {
	_, err := StudentTCDF(1, -2)
	requireErr(t, "df < 0", err, ErrOutOfRange)
	_, err = StudentTCDF(inf, 2)
	requireErr(t, "t Inf", err, ErrNotFinite)
}

func BenchmarkStudentTCDF(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = StudentTCDF(2.1, 10)
	}
}

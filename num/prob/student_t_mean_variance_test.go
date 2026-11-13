package prob

import "testing"

// JKB vol. 2 ch. 28 §2: variance ν/(ν − 2); ν = 4 → (0, 2).
func TestStudentTMeanVariance_sourceExample(t *testing.T) {
	if m, v, err := StudentTMeanVariance(4); err != nil || m != 0 || v != 2 {
		t.Errorf("StudentTMeanVariance(4) = (%v, %v, %v)", m, v, err)
	}
}

func TestStudentTMeanVariance_errors(t *testing.T) {
	for _, df := range []float64{1, 2, -1} {
		_, _, err := StudentTMeanVariance(df)
		requireErr(t, "df ≤ 2", err, ErrOutOfRange)
	}
	_, _, err := StudentTMeanVariance(nan)
	requireErr(t, "df NaN", err, ErrNotFinite)
}

func BenchmarkStudentTMeanVariance(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _, _ = StudentTMeanVariance(4)
	}
}

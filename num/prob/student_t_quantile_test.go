package prob

import (
	"math"
	"testing"
)

// ν = 1: tan(π(u − ½)); NIST e-Handbook §1.3.6.7.2 t table: ν = 10, 0.975 → 2.228; ν = 30, 0.995 → 2.750.
func TestStudentTQuantile_sourceExample(t *testing.T) {
	if got, err := StudentTQuantile(0.975, 1); err != nil || !within(got, math.Tan(0.475*math.Pi), 1e-12) {
		t.Errorf("ν = 1: (%v, %v)", got, err)
	}
	cases := []struct{ u, df, want float64 }{{0.975, 10, 2.228}, {0.995, 30, 2.750}, {0.95, 5, 2.015}}
	for _, c := range cases {
		if got, _ := StudentTQuantile(c.u, c.df); math.Abs(got-c.want) > 5e-4 {
			t.Errorf("u %v ν %v: %v, want %v", c.u, c.df, got, c.want)
		}
	}
}

func TestStudentTQuantile_roundTrip(t *testing.T) {
	for _, df := range []float64{0.5, 1, 3, 100} {
		for _, u := range []float64{1e-12, 0.01, 0.3, 0.7, 0.99, 1 - 1e-10} {
			x, err := StudentTQuantile(u, df)
			tail, _, _ := studentTTail(x, df)
			p := min(u, 1-u)
			if err != nil || !within(tail, p, 1e-10) {
				t.Errorf("ν %v u %v: x %v tail %v err %v", df, u, x, tail, err)
			}
		}
		// 0.25 and 0.75 are exact in float64 and 1 − 0.75 == 0.25, so exact symmetry is a fair demand
		// (1 − 0.8 != 0.2 in float64, so 0.2/0.8 are not mirror images).
		lo, _ := StudentTQuantile(0.25, df)
		hi, _ := StudentTQuantile(0.75, df)
		if lo != -hi {
			t.Errorf("ν %v: not symmetric %v %v", df, lo, hi)
		}
	}
}

func TestStudentTQuantile_errors(t *testing.T) {
	if got, err := StudentTQuantile(0.5, 3); err != nil || got != 0 {
		t.Errorf("u = ½: (%v, %v)", got, err)
	}
	for _, u := range []float64{0, 1} {
		_, err := StudentTQuantile(u, 3)
		requireErr(t, "infinite quantile", err, ErrOutOfRange)
	}
	_, err := StudentTQuantile(2, 3)
	requireErr(t, "u > 1", err, ErrProbability)
	_, err = StudentTQuantile(0.9, 0)
	requireErr(t, "df 0", err, ErrOutOfRange)
	_, err = StudentTQuantile(nan, 3)
	requireErr(t, "u NaN", err, ErrNotFinite)
}

func BenchmarkStudentTQuantile(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = StudentTQuantile(0.975, 10)
	}
}

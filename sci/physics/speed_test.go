package physics

import (
	"errors"
	"math"
	"testing"
)

// Hand computations from average speed = distance / time (OpenStax University Physics Vol. 1, §3.1):
// 100 000 m / 3600 s = 27.7777… m/s (100 km/h). Tolerance 1e-12 relative: one operation.
func TestAverageSpeed_cases(t *testing.T) {
	got, err := AverageSpeed(100000, 3600)
	if err != nil || math.Abs(got-250.0/9) > 1e-12*got {
		t.Errorf("got %v %v", got, err)
	}
	for _, c := range []struct {
		d, s float64
		want error
	}{
		{1, 0, ErrZero},
		{-1, 1, ErrNegative},
		{math.NaN(), 1, ErrNotFinite},
		{math.MaxFloat64, 0.5, ErrOverflow},
	} {
		if _, err := AverageSpeed(c.d, c.s); !errors.Is(err, c.want) {
			t.Errorf("AverageSpeed(%v, %v): %v, want %v", c.d, c.s, err, c.want)
		}
	}
}

func TestDistanceFromSpeed_cases(t *testing.T) {
	if got, err := DistanceFromSpeed(25, 7200); err != nil || got != 180000 {
		t.Errorf("got %v %v", got, err)
	}
	if _, err := DistanceFromSpeed(-25, 7200); !errors.Is(err, ErrNegative) {
		t.Errorf("negative: %v", err)
	}
}

func TestDurationFromSpeed_cases(t *testing.T) {
	if got, err := DurationFromSpeed(180000, 25); err != nil || got != 7200 {
		t.Errorf("got %v %v", got, err)
	}
	if _, err := DurationFromSpeed(1, 0); !errors.Is(err, ErrZero) {
		t.Errorf("zero: %v", err)
	}
}

func BenchmarkAverageSpeed(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = AverageSpeed(100000, 3600)
	}
}

func BenchmarkDistanceFromSpeed(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = DistanceFromSpeed(25, 7200)
	}
}

func BenchmarkDurationFromSpeed(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = DurationFromSpeed(180000, 25)
	}
}

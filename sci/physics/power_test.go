package physics

import (
	"errors"
	"math"
	"testing"
)

// Hand computation from W = F·d and P = W/t (OpenStax University Physics Vol. 1, §7.1, §7.4):
// 10 N × 3 m / 2 s = 15 W; 4448.2216152605 N (1000 lbf) × 1609.344 m / 60 s = 119 311.979... W (= 160 hp of 550 ft·lbf/s). Exact in the first
// case; 1e-12 relative in the second.
func TestPowerFromWork_cases(t *testing.T) {
	if got, err := PowerFromWork(10, 3, 2); err != nil || got != 15 {
		t.Errorf("got %v %v", got, err)
	}
	want := 4448.2216152605 * 1609.344 / 60
	if got, err := PowerFromWork(4448.2216152605, 1609.344, 60); err != nil || math.Abs(got-want) > 1e-12*want {
		t.Errorf("got %v %v, want %v", got, err, want)
	}
	for _, c := range []struct {
		f, d, s float64
		want    error
	}{
		{1, 1, 0, ErrZero},
		{-1, 1, 1, ErrNegative},
		{1, math.NaN(), 1, ErrNotFinite},
		{math.MaxFloat64, 2, 1, ErrOverflow},
		{math.MaxFloat64, 1, 0.5, ErrOverflow},
	} {
		if _, err := PowerFromWork(c.f, c.d, c.s); !errors.Is(err, c.want) {
			t.Errorf("(%v, %v, %v): %v, want %v", c.f, c.d, c.s, err, c.want)
		}
	}
}

func BenchmarkPowerFromWork(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = PowerFromWork(10, 3, 2)
	}
}

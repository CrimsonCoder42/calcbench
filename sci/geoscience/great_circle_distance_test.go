package geoscience

import (
	"errors"
	"math"
	"testing"
)

// Hand computations: a quarter meridian on the mean-radius sphere is R·π/2 = 10 007 557.176... m; one degree
// of the equator is R·π/180. Tolerance 1e-9 relative: one central angle and one product.
func TestGreatCircleDistance_cases(t *testing.T) {
	R := MeanEarthRadiusMeters
	cases := []struct {
		name                     string
		la1, lo1, la2, lo2, want float64
	}{
		{"pole to equator", 90, 0, 0, 0, R * math.Pi / 2},
		{"equator 1 degree", 0, 0, 0, 1, R * math.Pi / 180},
		{"same point", 12.5, -45, 12.5, -45, 0},
		{"antipodes", 0, 0, 0, 180, R * math.Pi},
	}
	for _, c := range cases {
		got, err := GreatCircleDistance(c.la1, c.lo1, c.la2, c.lo2, R)
		if err != nil || math.Abs(got-c.want) > 1e-9*math.Max(1, c.want) {
			t.Errorf("%s: got %v %v, want %v", c.name, got, err, c.want)
		}
	}
}

func TestGreatCircleDistance_errors(t *testing.T) {
	cases := []struct {
		la1, lo1, la2, lo2, r float64
		want                  error
	}{
		{91, 0, 0, 0, 1, ErrLatitude},
		{0, 0, -90.5, 0, 1, ErrLatitude},
		{math.NaN(), 0, 0, 0, 1, ErrNotFinite},
		{0, math.Inf(1), 0, 0, 1, ErrNotFinite},
		{0, 0, 1, 1, 0, ErrEllipsoid},
		{0, 0, 1, 1, -1, ErrEllipsoid},
		{0, 0, 1, 1, math.Inf(1), ErrEllipsoid},
	}
	for _, c := range cases {
		if _, err := GreatCircleDistance(c.la1, c.lo1, c.la2, c.lo2, c.r); !errors.Is(err, c.want) {
			t.Errorf("%+v: %v, want %v", c, err, c.want)
		}
	}
}

func BenchmarkGreatCircleDistance(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = GreatCircleDistance(38.8976, -77.0366, 39.9496, -75.1503, MeanEarthRadiusMeters)
	}
}

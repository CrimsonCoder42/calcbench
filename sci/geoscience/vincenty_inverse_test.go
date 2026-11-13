package geoscience

import (
	"errors"
	"math"
	"testing"
)

// Vincenty, Survey Review 1975;23(176):88-93, test line (a): Bessel ellipsoid a = 6 377 397.155 m,
// f = 1/299.1528128, φ1 = 55°45′, φ2 = −33°26′, L = 108°13′ → s = 14 110 526.170 m. Tolerance 1 mm (the
// paper's printed precision).
func TestVincentyInverseDistance_vincentyLineA(t *testing.T) {
	got, err := VincentyInverseDistance(55.75, 0, -(33 + 26.0/60), 108+13.0/60, 6377397.155, 1/299.1528128)
	if err != nil || math.Abs(got-14110526.170) > 0.001 {
		t.Errorf("line (a): got %.4f %v, want 14110526.170", got, err)
	}
}

// Along the equator the geodesic is the equatorial arc: one degree is a·π/180 = 111 319.49079327357 m on WGS 84
// (NIMA TR8350.2 a = 6 378 137 m). Coincident points give 0. Symmetry: swapping the points gives the same length.
func TestWGS84Distance_cases(t *testing.T) {
	got, err := WGS84Distance(0, 0, 0, 1)
	if err != nil || math.Abs(got-111319.49079327357) > 1e-6 {
		t.Errorf("equator 1°: %v %v", got, err)
	}
	if got, err := WGS84Distance(10, 20, 10, 20); err != nil || got != 0 {
		t.Errorf("same point: %v %v", got, err)
	}
	ab, _ := WGS84Distance(51.5074, -0.1278, 40.7128, -74.0060)
	ba, _ := WGS84Distance(40.7128, -74.0060, 51.5074, -0.1278)
	if math.Abs(ab-ba) > 1e-6 {
		t.Errorf("asymmetric: %v vs %v", ab, ba)
	}
	// The ellipsoidal length lies within 0.6 % of the mean-radius sphere's for these routes (flattening ≈ 0.34 %).
	for _, p := range [][4]float64{{38.8976, -77.0366, 39.9496, -75.1503}, {51.5074, -0.1278, 40.7128, -74.0060}, {-33.8688, 151.2093, 35.6762, 139.6503}} {
		w, err := WGS84Distance(p[0], p[1], p[2], p[3])
		g, _ := GreatCircleDistance(p[0], p[1], p[2], p[3], MeanEarthRadiusMeters)
		if err != nil || math.Abs(w-g) > 0.006*g {
			t.Errorf("%v: wgs84 %v sphere %v %v", p, w, g, err)
		}
	}
}

func TestVincentyInverseDistance_errors(t *testing.T) {
	if _, err := WGS84Distance(0, 0, 0.5, 179.7); !errors.Is(err, ErrNoConvergence) {
		t.Errorf("nearly antipodal: %v", err)
	}
	if _, err := WGS84Distance(95, 0, 0, 0); !errors.Is(err, ErrLatitude) {
		t.Errorf("latitude: %v", err)
	}
	if _, err := WGS84Distance(0, math.NaN(), 0, 0); !errors.Is(err, ErrNotFinite) {
		t.Errorf("NaN: %v", err)
	}
	for _, e := range [][2]float64{{0, 0.003}, {6.4e6, 1}, {6.4e6, -0.1}, {math.NaN(), 0}} {
		if _, err := VincentyInverseDistance(0, 0, 1, 1, e[0], e[1]); !errors.Is(err, ErrEllipsoid) {
			t.Errorf("ellipsoid %v: %v", e, err)
		}
	}
}

func BenchmarkVincentyInverseDistance(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = VincentyInverseDistance(55.75, 0, -33.43333333333333, 108.21666666666667, 6377397.155, 1/299.1528128)
	}
}

func BenchmarkWGS84Distance(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_, _ = WGS84Distance(38.8976, -77.0366, 39.9496, -75.1503)
	}
}

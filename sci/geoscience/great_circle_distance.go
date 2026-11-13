package geoscience

import (
	"errors"
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/num/trig"
)

var (
	// ErrNotFinite reports a NaN or infinite input.
	ErrNotFinite = errors.New("geoscience: value is NaN or infinite")
	// ErrLatitude reports a latitude outside −90..90 degrees.
	ErrLatitude = errors.New("geoscience: latitude outside -90..90 degrees")
	// ErrEllipsoid reports a non-positive radius or semi-major axis, or a flattening outside [0, 1).
	ErrEllipsoid = errors.New("geoscience: invalid radius or ellipsoid")
	// ErrNoConvergence reports that Vincenty's inverse iteration did not converge (nearly antipodal points).
	ErrNoConvergence = errors.New("geoscience: Vincenty iteration did not converge")
)

// MeanEarthRadiusMeters is the mean radius R1 = (2a + b)/3 of the GRS 80 ellipsoid, 6 371 008.7714 m
// (Moritz H. Geodetic Reference System 1980. Bulletin Géodésique 1980;54(3):395-405; J Geodesy 2000;74:128-162).
const MeanEarthRadiusMeters = 6371008.7714

func checkPoint(latDeg, lonDeg float64, which string) error {
	if math.IsNaN(latDeg) || math.IsInf(latDeg, 0) || math.IsNaN(lonDeg) || math.IsInf(lonDeg, 0) {
		return fmt.Errorf("%w: point %s (%v, %v)", ErrNotFinite, which, latDeg, lonDeg)
	}
	if latDeg < -90 || latDeg > 90 {
		return fmt.Errorf("%w: point %s latitude %v", ErrLatitude, which, latDeg)
	}
	return nil
}

// GreatCircleDistance returns the great-circle distance between two points on a sphere of the given radius,
// s = R·σ, with σ the central angle by the haversine form (NUM-TRIG-031). With MeanEarthRadiusMeters it is
// the spherical approximation of a distance on the Earth (error up to about 0.5 % against the ellipsoid).
// Errors: ErrNotFinite; ErrLatitude; ErrEllipsoid for radiusMeters ≤ 0, NaN or ±Inf.
//
// ID:      SCI-GEO-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, eq. 4.42.8 (spherical law, via NUM-TRIG-031 CentralAngle); arc length s = Rσ; mean radius R1 of GRS 80 from Moritz H, Geodetic Reference System 1980, J Geodesy 2000;74:128-162
// Tags:    great-circle, haversine, distance, sphere, earth, geodesy, geoscience
// Units:   lat1Deg, lon1Deg, lat2Deg, lon2Deg decimal degrees; radiusMeters m; result m
// Range:   latitudes −90..90; radiusMeters > 0
// Example: (38.8976, −77.0366) to (39.9496, −75.1503) on R = 6 371 008.7714 m → 199 830 m (central angle 0.0313656 rad)
func GreatCircleDistance(lat1Deg, lon1Deg, lat2Deg, lon2Deg, radiusMeters float64) (float64, error) {
	if err := checkPoint(lat1Deg, lon1Deg, "1"); err != nil {
		return 0, err
	}
	if err := checkPoint(lat2Deg, lon2Deg, "2"); err != nil {
		return 0, err
	}
	if math.IsNaN(radiusMeters) || math.IsInf(radiusMeters, 0) || radiusMeters <= 0 {
		return 0, fmt.Errorf("%w: radiusMeters = %v", ErrEllipsoid, radiusMeters)
	}
	sigma, err := trig.CentralAngle(lat1Deg*math.Pi/180, lon1Deg*math.Pi/180, lat2Deg*math.Pi/180, lon2Deg*math.Pi/180)
	if err != nil {
		return 0, fmt.Errorf("geoscience: %w", err)
	}
	return radiusMeters * sigma, nil
}

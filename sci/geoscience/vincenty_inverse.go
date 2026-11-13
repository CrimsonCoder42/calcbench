package geoscience

import (
	"fmt"
	"math"
)

// WGS 84 defining parameters (NIMA TR8350.2, 3rd ed., 2000, Table 3.1).
const (
	wgs84SemiMajorMeters = 6378137.0
	wgs84Flattening      = 1 / 298.257223563
)

const (
	vincentyTolerance = 1e-12 // radians on λ, about 0.006 mm on the Earth
	vincentyMaxIter   = 200
)

// VincentyInverseDistance returns the length in meters of the geodesic between two points on an ellipsoid of
// revolution with semi-major axis a and flattening f, by Vincenty's inverse method: iterate on the longitude λ
// on the auxiliary sphere (eqs. 13-18 with reduced latitudes U = atan((1 − f) tan φ)) until |Δλ| < 1e-12 rad,
// then s = b·A·(σ − Δσ) (eqs. 3-6, u² = cos²α (a² − b²)/b²). Coincident points give 0.
// Errors: ErrNotFinite; ErrLatitude; ErrEllipsoid for a ≤ 0 or f outside [0, 1); ErrNoConvergence after 200
// iterations (nearly antipodal points, where the method is known to fail; a Karney 2013 solver is the remedy).
//
// ID:      SCI-GEO-002
// Time:    O(k)   k = iterations, at most 200 (typically under 10)
// Space:   O(1), no allocation
// Source:  Vincenty T. Direct and inverse solutions of geodesics on the ellipsoid with application of nested equations. Survey Review 1975;23(176):88-93, doi:10.1179/sre.1975.23.176.88 (inverse formulae, eqs. 3-6 and 13-18; test line (a))
// Tags:    geodesic, vincenty, ellipsoid, distance, earth, geodesy, geoscience
// Units:   lat1Deg, lon1Deg, lat2Deg, lon2Deg decimal degrees; semiMajorMeters m; flattening dimensionless; result m
// Range:   latitudes −90..90; semiMajorMeters > 0; 0 ≤ flattening < 1; not nearly antipodal
// Example: Bessel ellipsoid (a = 6 377 397.155 m, 1/f = 299.1528128), φ1 = 55°45′, φ2 = −33°26′, L = 108°13′ → 14 110 526.170 m (Vincenty's line (a))
func VincentyInverseDistance(lat1Deg, lon1Deg, lat2Deg, lon2Deg, semiMajorMeters, flattening float64) (float64, error) {
	if err := checkPoint(lat1Deg, lon1Deg, "1"); err != nil {
		return 0, err
	}
	if err := checkPoint(lat2Deg, lon2Deg, "2"); err != nil {
		return 0, err
	}
	if math.IsNaN(semiMajorMeters) || math.IsInf(semiMajorMeters, 0) || semiMajorMeters <= 0 ||
		math.IsNaN(flattening) || flattening < 0 || flattening >= 1 {
		return 0, fmt.Errorf("%w: a = %v, f = %v", ErrEllipsoid, semiMajorMeters, flattening)
	}
	a, f := semiMajorMeters, flattening
	b := a * (1 - f)
	rad := math.Pi / 180
	L := math.Remainder((lon2Deg-lon1Deg)*rad, 2*math.Pi)
	u1 := math.Atan((1 - f) * math.Tan(lat1Deg*rad))
	u2 := math.Atan((1 - f) * math.Tan(lat2Deg*rad))
	sinU1, cosU1 := math.Sincos(u1)
	sinU2, cosU2 := math.Sincos(u2)

	lambda := L
	var sinSigma, cosSigma, sigma, cos2Alpha, cos2SigmaM float64
	converged := false
	for i := 0; i < vincentyMaxIter; i++ {
		sinLambda, cosLambda := math.Sincos(lambda)
		t1 := cosU2 * sinLambda
		t2 := cosU1*sinU2 - sinU1*cosU2*cosLambda
		sinSigma = math.Sqrt(t1*t1 + t2*t2) // eq. 14
		if sinSigma == 0 {
			return 0, nil // coincident points
		}
		cosSigma = sinU1*sinU2 + cosU1*cosU2*cosLambda   // eq. 15
		sigma = math.Atan2(sinSigma, cosSigma)           // eq. 16
		sinAlpha := cosU1 * cosU2 * sinLambda / sinSigma // eq. 17
		cos2Alpha = 1 - sinAlpha*sinAlpha
		if cos2Alpha != 0 {
			cos2SigmaM = cosSigma - 2*sinU1*sinU2/cos2Alpha // eq. 18
		} else {
			cos2SigmaM = 0 // both points on the equator
		}
		c := f / 16 * cos2Alpha * (4 + f*(4-3*cos2Alpha)) // eq. 10
		prev := lambda
		lambda = L + (1-c)*f*sinAlpha*(sigma+c*sinSigma*(cos2SigmaM+c*cosSigma*(-1+2*cos2SigmaM*cos2SigmaM))) // eq. 11
		if math.Abs(lambda-prev) < vincentyTolerance {
			converged = true
			break
		}
	}
	if !converged {
		return 0, fmt.Errorf("%w: (%v, %v) to (%v, %v)", ErrNoConvergence, lat1Deg, lon1Deg, lat2Deg, lon2Deg)
	}
	uSq := cos2Alpha * (a*a - b*b) / (b * b)
	bigA := 1 + uSq/16384*(4096+uSq*(-768+uSq*(320-175*uSq))) // eq. 3
	bigB := uSq / 1024 * (256 + uSq*(-128+uSq*(74-47*uSq)))   // eq. 4
	deltaSigma := bigB * sinSigma * (cos2SigmaM + bigB/4*(cosSigma*(-1+2*cos2SigmaM*cos2SigmaM)-
		bigB/6*cos2SigmaM*(-3+4*sinSigma*sinSigma)*(-3+4*cos2SigmaM*cos2SigmaM))) // eq. 6
	return b * bigA * (sigma - deltaSigma), nil // eq. 19
}

// WGS84Distance returns the geodesic length in meters between two points on the WGS 84 ellipsoid
// (a = 6 378 137 m, 1/f = 298.257223563): VincentyInverseDistance with those parameters.
// Errors: as VincentyInverseDistance (ErrNotFinite, ErrLatitude, ErrNoConvergence).
//
// ID:      SCI-GEO-003
// Time:    O(k)   k = iterations, at most 200 (typically under 10)
// Space:   O(1), no allocation
// Source:  National Imagery and Mapping Agency. Department of Defense World Geodetic System 1984. NIMA TR8350.2, 3rd ed., amendment 1, 2000, Table 3.1 (a, 1/f); Vincenty T, Survey Review 1975;23(176):88-93 (inverse method)
// Tags:    geodesic, wgs84, vincenty, ellipsoid, distance, earth, gps, geodesy, geoscience
// Units:   lat1Deg, lon1Deg, lat2Deg, lon2Deg decimal degrees; result m
// Range:   latitudes −90..90; not nearly antipodal
// Example: (38.8976, −77.0366) to (39.9496, −75.1503) → 200 050.08 m; along the equator 1° → a·π/180 = 111 319.4908 m
func WGS84Distance(lat1Deg, lon1Deg, lat2Deg, lon2Deg float64) (float64, error) {
	return VincentyInverseDistance(lat1Deg, lon1Deg, lat2Deg, lon2Deg, wgs84SemiMajorMeters, wgs84Flattening)
}

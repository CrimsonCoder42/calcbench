// Package geoscience holds the geometry of the Earth: the distance between two points on its surface, on a
// sphere of stated radius and on an ellipsoid (WGS 84) by Vincenty's inverse method.
//
// ID short form: GEO (IDs are SCI-GEO-NNN).
//
// Units: latitudes and longitudes in decimal degrees (north and east positive), distances and radii in meters.
// NaN or ±Inf returns ErrNotFinite; a latitude outside −90..90 returns ErrLatitude; a non-positive radius or an
// invalid ellipsoid returns ErrEllipsoid; Vincenty's iteration failing to converge (nearly antipodal points)
// returns ErrNoConvergence. Results are not rounded.
//
// Functions:
//
//	SCI-GEO-001  GreatCircleDistance       s = R·σ, σ the central angle (haversine form, NUM-TRIG-031)
//	SCI-GEO-002  VincentyInverseDistance   geodesic length on an ellipsoid (a, f), Vincenty 1975 inverse
//	SCI-GEO-003  WGS84Distance             VincentyInverseDistance on WGS 84 (a = 6 378 137 m, 1/f = 298.257223563)
package geoscience

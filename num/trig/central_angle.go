package trig

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// CentralAngle returns the central angle σ, in [0, π], between two points on a sphere given by latitude
// φ and longitude λ: the great-circle distance on the unit sphere (multiply by a radius for a length).
// It uses the haversine form hav σ = hav(φ2 − φ1) + cos φ1 cos φ2 hav(λ2 − λ1), hav θ = sin²(θ/2),
// σ = 2 arcsin(min(1, √hav σ)), which is accurate for small σ where the spherical law of cosines
// (cos σ = sin φ1 sin φ2 + cos φ1 cos φ2 cos Δλ) loses digits; it loses some accuracy only for nearly
// antipodal points. The half-differences are formed as φ2/2 − φ1/2, so they cannot overflow, and
// longitudes need no normalising. (The haversine form is old navigation practice; R. W. Sinnott's 1984
// Sky and Telescope note "Virtues of the Haversine" is often credited for it, but that note was not seen
// and is not relied on here.)
// Errors: ErrNotFinite when an input is NaN or ±Inf; ErrDomain when a latitude is outside [−π/2, π/2].
//
// ID:      NUM-TRIG-031
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  NIST Digital Library of Mathematical Functions, Release 1.2.8 (2026-09-15), eq. 4.42.8 (with b, c the colatitudes and A = Δλ) and 4.21.28 (cos 2z = 1 − 2 sin² z), which give the haversine form, https://dlmf.nist.gov/4.42 (read 2026-11-12)
// Source:  Abramowitz and Stegun (eds.), Handbook of Mathematical Functions, NBS Applied Mathematics Series 55, 10th printing, 1972, eq. 4.3.147 ("haversine A = hav A = ½ vers A"), p. 78
// Tags:    trig, spherical, great-circle, haversine
// Units:   lat1Rad, lon1Rad, lat2Rad, lon2Rad in radians; result in radians, in [0, π]
// Example: (0, 0) to (0, π/2) → π/2; a pole to the equator → π/2; two equator points 1° apart → π/180 (hand-derived)
func CentralAngle(lat1Rad, lon1Rad, lat2Rad, lon2Rad float64) (float64, error) {
	for _, p := range [4]struct {
		v    float64
		name string
	}{{lat1Rad, "lat1Rad"}, {lon1Rad, "lon1Rad"}, {lat2Rad, "lat2Rad"}, {lon2Rad, "lon2Rad"}} {
		if !reduce.AllFinite(p.v) {
			return 0, fmt.Errorf("%w: %s", ErrNotFinite, p.name)
		}
	}
	if math.Abs(lat1Rad) > math.Pi/2 {
		return 0, fmt.Errorf("%w: lat1Rad outside [−π/2, π/2]", ErrDomain)
	}
	if math.Abs(lat2Rad) > math.Pi/2 {
		return 0, fmt.Errorf("%w: lat2Rad outside [−π/2, π/2]", ErrDomain)
	}
	sLat := math.Sin(lat2Rad/2 - lat1Rad/2)
	sLon := math.Sin(lon2Rad/2 - lon1Rad/2)
	h := sLat*sLat + math.Cos(lat1Rad)*math.Cos(lat2Rad)*sLon*sLon
	return 2 * math.Asin(math.Min(1, math.Sqrt(h))), nil
}

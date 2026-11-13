package meteorology

import (
	"fmt"
	"math"
)

// HeatIndexF returns the NWS heat index, in °F, of an air temperature (°F) and relative humidity (%).
// First Steadman's simple form, HI = 0.5·(T + 61 + 1.2·(T − 68) + 0.094·RH); if its mean with T is 80 °F or more,
// the Rothfusz regression
// HI = −42.379 + 2.04901523T + 10.14333127RH − 0.22475541T·RH − 0.00683783T² − 0.05481717RH²
// + 0.00122874T²·RH + 0.00085282T·RH² − 0.00000199T²·RH², minus ((13 − RH)/4)·√((17 − |T − 95|)/17) when
// RH < 13 % and 80 ≤ T ≤ 112 °F, plus ((RH − 85)/10)·((87 − T)/5) when RH > 85 % and 80 ≤ T ≤ 87 °F.
// Errors: ErrNotFinite; ErrOutOfRange for relativeHumidityPercent outside 0..100 or airTempF below absolute
// zero or above 200 °F.
//
// ID:      SCI-METEO-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Rothfusz LP. The Heat Index "Equation" (or, More Than You Ever Wanted to Know About Heat Index). NWS Southern Region Technical Attachment SR 90-23, 1990 (the regression and its coefficients)
// Source:  NOAA NWS Weather Prediction Center, "The Heat Index Equation", https://www.wpc.ncep.noaa.gov/html/heatindex_equation.shtml (read 2026-11-13): the simple form, the 80 °F switch, and the two adjustments
// Tags:    heat-index, apparent-temperature, humidity, weather, nws, meteorology
// Units:   airTempF °F; relativeHumidityPercent %; result °F
// Range:   −459.67 ≤ airTempF ≤ 200; 0 ≤ relativeHumidityPercent ≤ 100
// Example: 85 °F, 70 % → 92.70 °F; 70 °F, 50 % → 69.05 °F (simple form)
func HeatIndexF(airTempF, relativeHumidityPercent float64) (float64, error) {
	if err := checkFinite("airTempF", airTempF); err != nil {
		return 0, err
	}
	if err := checkFinite("relativeHumidityPercent", relativeHumidityPercent); err != nil {
		return 0, err
	}
	if airTempF < -459.67 || airTempF > 200 {
		return 0, fmt.Errorf("%w: airTempF = %v, want −459.67..200", ErrOutOfRange, airTempF)
	}
	if relativeHumidityPercent < 0 || relativeHumidityPercent > 100 {
		return 0, fmt.Errorf("%w: relativeHumidityPercent = %v, want 0..100", ErrOutOfRange, relativeHumidityPercent)
	}
	t, rh := airTempF, relativeHumidityPercent
	simple := 0.5 * (t + 61 + (t-68)*1.2 + rh*0.094)
	if (simple+t)/2 < 80 {
		return simple, nil
	}
	hi := -42.379 + 2.04901523*t + 10.14333127*rh - 0.22475541*t*rh - 0.00683783*t*t -
		0.05481717*rh*rh + 0.00122874*t*t*rh + 0.00085282*t*rh*rh - 0.00000199*t*t*rh*rh
	if rh < 13 && t >= 80 && t <= 112 {
		hi -= ((13 - rh) / 4) * math.Sqrt((17-math.Abs(t-95))/17)
	} else if rh > 85 && t >= 80 && t <= 87 {
		hi += ((rh - 85) / 10) * ((87 - t) / 5)
	}
	return hi, nil
}

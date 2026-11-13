package meteorology

import (
	"errors"
	"fmt"
	"math"
)

var (
	// ErrNotFinite reports a NaN or infinite input.
	ErrNotFinite = errors.New("meteorology: value is NaN or infinite")
	// ErrOutOfRange reports an input outside the formula's stated range.
	ErrOutOfRange = errors.New("meteorology: value outside the formula's range")
)

func checkFinite(name string, v float64) error {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return fmt.Errorf("%w: %s = %v", ErrNotFinite, name, v)
	}
	return nil
}

// WindChillF returns the NWS wind chill temperature, in °F:
// WCT = 35.74 + 0.6215·T − 35.75·V^0.16 + 0.4275·T·V^0.16, T the air temperature in °F and V the wind speed in mph
// at 10 m (the 2001 JAG/TI index adopted by the NWS and the Meteorological Service of Canada).
// Errors: ErrNotFinite; ErrOutOfRange for airTempF > 50 °F or windSpeedMph < 3 mph (outside the index's defined
// range), or below absolute zero, or windSpeedMph above 200 mph.
//
// ID:      SCI-METEO-001
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Office of the Federal Coordinator for Meteorology. Report on Wind Chill Temperature and Extreme Heat Indices: Evaluation and Improvement Projects, FCM-R19-2003, 2003, ch. 2 and Appendix C (the WCT equation and its range: T ≤ 50 °F, V ≥ 3 mph)
// Source:  Osczevski R, Bluestein M. The new wind chill equivalent temperature chart. Bull Am Meteorol Soc. 2005;86(10):1453-1458. doi:10.1175/BAMS-86-10-1453
// Tags:    wind-chill, apparent-temperature, weather, nws, meteorology
// Units:   airTempF °F; windSpeedMph mph; result °F
// Range:   −459.67 ≤ airTempF ≤ 50; 3 ≤ windSpeedMph ≤ 200
// Example: 10 °F, 15 mph → −6.59 °F (the NWS chart prints −7)
func WindChillF(airTempF, windSpeedMph float64) (float64, error) {
	if err := checkFinite("airTempF", airTempF); err != nil {
		return 0, err
	}
	if err := checkFinite("windSpeedMph", windSpeedMph); err != nil {
		return 0, err
	}
	if airTempF > 50 || airTempF < -459.67 {
		return 0, fmt.Errorf("%w: airTempF = %v, want −459.67..50", ErrOutOfRange, airTempF)
	}
	if windSpeedMph < 3 || windSpeedMph > 200 {
		return 0, fmt.Errorf("%w: windSpeedMph = %v, want 3..200", ErrOutOfRange, windSpeedMph)
	}
	v16 := math.Pow(windSpeedMph, 0.16)
	return 35.74 + 0.6215*airTempF - 35.75*v16 + 0.4275*airTempF*v16, nil
}

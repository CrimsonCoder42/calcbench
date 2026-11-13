package meteorology

import (
	"fmt"
	"math"
)

// Magnus coefficients over water recommended by Alduchov and Eskridge (1996): e_s = 6.1094·exp(17.625·t/(t + 243.04)).
const (
	magnusA = 17.625
	magnusB = 243.04 // °C
)

// DewPointC returns the dew point, in °C, of air at a temperature and relative humidity, by inverting the
// Magnus form: γ = ln(RH/100) + a·t/(b + t), Td = b·γ/(a − γ), a = 17.625, b = 243.04 °C. At RH = 100 % the
// result is t exactly, and it is never above t (floating-point rounding near saturation is clamped).
// Errors: ErrNotFinite; ErrOutOfRange for airTempC outside −40..50 °C (the range Alduchov and Eskridge fitted)
// or relativeHumidityPercent outside (0, 100].
//
// ID:      SCI-METEO-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Alduchov OA, Eskridge RE. Improved Magnus form approximation of saturation vapor pressure. J Appl Meteorol. 1996;35(4):601-609, eq. 21 (AERK coefficients 17.625, 243.04 °C). doi:10.1175/1520-0450(1996)035<0601:IMFAOS>2.0.CO;2
// Source:  Lawrence MG. The relationship between relative humidity and the dewpoint temperature in moist air. Bull Am Meteorol Soc. 2005;86(2):225-233, eq. 8 (the inversion for Td). doi:10.1175/BAMS-86-2-225
// Tags:    dew-point, humidity, relative-humidity, magnus, weather, meteorology
// Units:   airTempC °C; relativeHumidityPercent %; result °C
// Range:   −40 ≤ airTempC ≤ 50; 0 < relativeHumidityPercent ≤ 100
// Example: 21.111 °C (70 °F), 65 % → 14.28 °C (57.7 °F); RH 100 % → Td = t
func DewPointC(airTempC, relativeHumidityPercent float64) (float64, error) {
	if err := checkFinite("airTempC", airTempC); err != nil {
		return 0, err
	}
	if err := checkFinite("relativeHumidityPercent", relativeHumidityPercent); err != nil {
		return 0, err
	}
	if airTempC < -40 || airTempC > 50 {
		return 0, fmt.Errorf("%w: airTempC = %v, want −40..50", ErrOutOfRange, airTempC)
	}
	if relativeHumidityPercent <= 0 || relativeHumidityPercent > 100 {
		return 0, fmt.Errorf("%w: relativeHumidityPercent = %v, want (0, 100]", ErrOutOfRange, relativeHumidityPercent)
	}
	if relativeHumidityPercent == 100 {
		return airTempC, nil // saturated air: Td = t exactly
	}
	g := math.Log(relativeHumidityPercent/100) + magnusA*airTempC/(magnusB+airTempC)
	td := magnusB * g / (magnusA - g)
	if td > airTempC {
		td = airTempC // rounding near saturation can not lift Td above t
	}
	return td, nil
}

// RelativeHumidityFromDewPoint returns the relative humidity, in %, of air at a temperature and dew point, by
// the Magnus form: RH = 100·exp(a·Td/(b + Td) − a·t/(b + t)), a = 17.625, b = 243.04 °C.
// Errors: ErrNotFinite; ErrOutOfRange for a temperature outside −40..50 °C or a dew point above the temperature.
//
// ID:      SCI-METEO-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Alduchov OA, Eskridge RE. Improved Magnus form approximation of saturation vapor pressure. J Appl Meteorol. 1996;35(4):601-609, eq. 21. doi:10.1175/1520-0450(1996)035<0601:IMFAOS>2.0.CO;2
// Source:  Lawrence MG. Bull Am Meteorol Soc. 2005;86(2):225-233, eq. 6 (RH as the ratio of saturation pressures at Td and t). doi:10.1175/BAMS-86-2-225
// Tags:    relative-humidity, dew-point, humidity, magnus, weather, meteorology
// Units:   airTempC °C; dewPointC °C; result %
// Range:   −40 ≤ dewPointC ≤ airTempC ≤ 50
// Example: 29.444 °C (85 °F), dew point 21.111 °C (70 °F) → 60.9 %
func RelativeHumidityFromDewPoint(airTempC, dewPointC float64) (float64, error) {
	if err := checkFinite("airTempC", airTempC); err != nil {
		return 0, err
	}
	if err := checkFinite("dewPointC", dewPointC); err != nil {
		return 0, err
	}
	if airTempC < -40 || airTempC > 50 || dewPointC < -40 {
		return 0, fmt.Errorf("%w: airTempC = %v, dewPointC = %v, want −40..50", ErrOutOfRange, airTempC, dewPointC)
	}
	if dewPointC > airTempC {
		return 0, fmt.Errorf("%w: dewPointC = %v above airTempC = %v", ErrOutOfRange, dewPointC, airTempC)
	}
	return 100 * math.Exp(magnusA*dewPointC/(magnusB+dewPointC)-magnusA*airTempC/(magnusB+airTempC)), nil
}

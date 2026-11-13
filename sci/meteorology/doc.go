// Package meteorology holds the apparent-temperature and humidity formulas of operational meteorology: the
// NWS wind chill temperature, the NWS heat index, and dew point / relative humidity by the Magnus form.
//
// ID short form: METEO (IDs are SCI-METEO-NNN).
//
// Units: each formula takes the units its source states (°F and mph for the NWS indices, °C for the Magnus form),
// named in every parameter; conversions are unit/temperature's. NaN or ±Inf returns ErrNotFinite; inputs outside
// the source's stated range return ErrOutOfRange. Results are not rounded.
//
// Functions:
//
//	SCI-METEO-001  WindChillF                  NWS 2001 wind chill temperature, °F
//	SCI-METEO-002  HeatIndexF                  NWS heat index (Steadman simple form, Rothfusz regression, adjustments), °F
//	SCI-METEO-003  DewPointC                   dew point from temperature and RH, Magnus form, °C
//	SCI-METEO-004  RelativeHumidityFromDewPoint RH from temperature and dew point, Magnus form, %
package meteorology

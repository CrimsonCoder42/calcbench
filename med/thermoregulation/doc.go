// Package thermoregulation computes body-temperature categories from their primary sources: the staging of
// accidental hypothermia by core temperature (Paal 2016, Table 1).
//
// ID short form: THERMO (IDs are MED-THERMO-NNN).
//
// Units: every input names its unit (coreTemperatureC, coreTemperatureF). A function that takes degrees
// Fahrenheit compares against the exact Fahrenheit images of the source's Celsius edges (35 °C = 95 °F,
// 32 °C = 89.6 °F, 28 °C = 82.4 °F) rather than converting the input, so the edges stay exact. NaN and ±Inf
// return ErrNotFinite. Category functions return the source's own label, verbatim, as a named string type.
//
// Results are numbers and the source's own category labels; nothing here diagnoses, interprets or advises.
//
// Functions:
//
//	MED-THERMO-001  HypothermiaStageSwiss   accidental-hypothermia stage I–III for a core temperature in °C
//	MED-THERMO-002  HypothermiaStageSwissF  same stage for a core temperature in °F
package thermoregulation

// Package nutrition holds energy and nutrient-intake calculators: the Mifflin–St Jeor resting energy
// expenditure (both forms), macronutrient grams from the Acceptable Macronutrient Distribution Ranges, the DRI
// tables for total water and for protein and carbohydrate, the percent Daily Value of a nutrition label, and
// the ISSN 2017 protein range for exercising individuals.
//
// ID short form: NUTR (IDs are LIF-NUTR-NNN).
//
// Units: every input and result names its unit: body weight kg, height cm, energy kcal/day, nutrients g/day,
// water L/day (the unit its table prints). Ages are the caller's age in decimal years, or in decimal months
// where the source's table starts in months (the DRI tables); groups follow the table's own life stage rows.
// Sex is female bool where a coefficient or row depends on it. NaN or ±Inf in any input returns ErrNotFinite
// before any range check; an input outside the range the source states returns ErrOutOfRange; inputs valid
// one by one that have no row together, or a result beyond float64, return ErrDomain. Nothing is rounded,
// except the percent Daily Value, which the regulation rounds to a whole percent. Embedded government tables
// are copied from the official text, cited to the table, with the date read.
//
// Results are numbers and the source's own labels; nothing here diagnoses, interprets or advises.
//
// Functions:
//
//	LIF-NUTR-001  RestingEnergyExpenditureMifflinStJeor            REE (kcal/day), Mifflin 1990 equation with 166 × sex
//	LIF-NUTR-002  RestingEnergyExpenditureMifflinStJeorSimplified  REE (kcal/day), Mifflin 1990 sex-specific simplified forms
//	LIF-NUTR-003  AMDRGrams                                        protein, carbohydrate and fat g/day at the AMDR low and high ends
//	LIF-NUTR-004  TotalWaterAI                                     total-water and beverage AI (L/day), IOM 2005 Table S-2
//	LIF-NUTR-005  ProteinCarbohydrateDRIFor                        protein and carbohydrate RDA or AI (g/day), NASEM DRI summary table
//	LIF-NUTR-006  PercentDailyValue                                percent Daily Value of a DRV nutrient, 21 CFR 101.9(d)(7)(ii)
//	LIF-NUTR-007  ProteinRangeISSN                                 ISSN 2017 protein range 1.4–2.0 g/kg/day (g/day)
package nutrition

package nutrition

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// MacronutrientGrams is a day's macronutrient range in grams: the low and high ends for protein, carbohydrate
// and fat.
type MacronutrientGrams struct {
	ProteinLowG, ProteinHighG           float64
	CarbohydrateLowG, CarbohydrateHighG float64
	FatLowG, FatHighG                   float64
}

// amdrPercent is one age group's AMDR, in percent of energy (low, high).
type amdrPercent struct{ fatLo, fatHi, carbLo, carbHi, proteinLo, proteinHi float64 }

// NASEM 2024 Table 3-1 (SOURCE: IOM, 2002/2005), percent of energy.
var (
	amdrChildren1to3  = amdrPercent{30, 40, 45, 65, 5, 20}  // "Children 1–3 years"
	amdrChildren4to18 = amdrPercent{25, 35, 45, 65, 10, 30} // "Children 4–18 years"
	amdrAdults        = amdrPercent{20, 35, 45, 65, 10, 35} // "Adults" (the DRI adult life stages, 19 y and older)
)

// 21 CFR 101.9(c)(1)(i)(B): "the general factors of 4, 4, and 9 calories per gram for protein, total
// carbohydrate, and total fat, respectively". Each is stored ×100 so that grams = kcal × percent / factor100.
const (
	kcalPerGProteinX100      = 400.0
	kcalPerGCarbohydrateX100 = 400.0
	kcalPerGFatX100          = 900.0
)

// AMDRGrams returns the grams per day of protein, carbohydrate and fat at the low and high ends of the
// Acceptable Macronutrient Distribution Ranges (AMDR) for an energy intake and an age:
// grams = energy × AMDR percent / 100 / (4 kcal/g protein, 4 kcal/g carbohydrate, 9 kcal/g fat).
// AMDR (percent of energy): ages 1–3 fat 30–40, carbohydrate 45–65, protein 5–20; ages 4–18 fat 25–35,
// carbohydrate 45–65, protein 10–30; adults (19 and older) fat 20–35, carbohydrate 45–65, protein 10–35.
// Not rounded.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when a result is beyond float64.
//
// ID:      LIF-NUTR-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  National Academies of Sciences, Engineering, and Medicine. Rethinking the Acceptable Macronutrient Distribution Range for the 21st Century: A Letter Report. Washington, DC: The National Academies Press; 2024. doi:10.17226/27957, Table 3-1 "Acceptable Macronutrient Distribution Ranges (AMDRs) by Age Group" (SOURCE: IOM, 2002/2005), read 2026-11-12 at nationalacademies.org/read/27957/chapter/5
// Source:  21 CFR 101.9(c)(1)(i)(B), eCFR current text read 2026-11-12: "the general factors of 4, 4, and 9 calories per gram for protein, total carbohydrate, and total fat"
// Tags:    nutrition, dri, amdr, macronutrient, protein, carbohydrate, fat, energy
// Units:   energyKcalPerDay kcal/day; ageYears years (decimal; groups by completed years); results g/day
// Range:   energyKcalPerDay > 0; ageYears ≥ 1 (Table 3-1 starts at 1–3 years; infants have no AMDR)
// Example: adult, 2000 kcal/day → protein 50–175 g, carbohydrate 225–325 g, fat 44.4…–77.7… g (HAND-DERIVED: 2000 × 10 / 400 = 50; 2000 × 20 / 900 = 44.44)
func AMDRGrams(energyKcalPerDay, ageYears float64) (MacronutrientGrams, error) {
	if err := checkFinitePositive("energyKcalPerDay", energyKcalPerDay); err != nil {
		return MacronutrientGrams{}, err
	}
	if err := checkFinite("ageYears", ageYears); err != nil {
		return MacronutrientGrams{}, err
	}
	var p amdrPercent
	switch {
	case !(ageYears >= 1):
		return MacronutrientGrams{}, fmt.Errorf("%w: ageYears = %v, want ≥ 1", ErrOutOfRange, ageYears)
	case ageYears < 4:
		p = amdrChildren1to3
	case ageYears < 19:
		p = amdrChildren4to18
	default:
		p = amdrAdults
	}
	g := MacronutrientGrams{
		ProteinLowG:       energyKcalPerDay * p.proteinLo / kcalPerGProteinX100,
		ProteinHighG:      energyKcalPerDay * p.proteinHi / kcalPerGProteinX100,
		CarbohydrateLowG:  energyKcalPerDay * p.carbLo / kcalPerGCarbohydrateX100,
		CarbohydrateHighG: energyKcalPerDay * p.carbHi / kcalPerGCarbohydrateX100,
		FatLowG:           energyKcalPerDay * p.fatLo / kcalPerGFatX100,
		FatHighG:          energyKcalPerDay * p.fatHi / kcalPerGFatX100,
	}
	// Each high end is at least its low end, so finite high ends mean finite low ends.
	if !reduce.AllFinite(g.CarbohydrateHighG, g.ProteinHighG, g.FatHighG) {
		return MacronutrientGrams{}, fmt.Errorf("%w: result is beyond float64", ErrDomain)
	}
	return g, nil
}

package nutrition

import "fmt"

// LifeStage selects the pregnancy or lactation rows of a Dietary Reference Intake table.
type LifeStage int

const (
	// LifeStageNone selects the age-and-sex rows.
	LifeStageNone LifeStage = iota
	// LifeStagePregnancy selects the "Pregnancy" rows (females 14 through 50 y).
	LifeStagePregnancy
	// LifeStageLactation selects the "Lactation" rows (females 14 through 50 y).
	LifeStageLactation
)

// driGroup is a DRI life stage group, in table order.
type driGroup int

const (
	driInfant0to6mo driGroup = iota
	driInfant7to12mo
	driChild1to3y
	driChild4to8y
	dri9to13y
	dri14to18y
	dri19yPlus
	driPregnancy
	driLactation
)

// driLifeStageGroup returns the DRI life stage group for an age in months, sex and life stage.
// Groups: 0 through 6 mo (ageMonths < 7), 7 through 12 mo (< 12), then by years: 1 through 3, 4 through 8,
// 9 through 13, 14 through 18, 19 and older; pregnancy and lactation for females 14 ≤ years < 51.
func driLifeStageGroup(ageMonths float64, female bool, stage LifeStage) (driGroup, error) {
	if err := checkFinite("ageMonths", ageMonths); err != nil {
		return 0, err
	}
	if !(ageMonths >= 0) {
		return 0, fmt.Errorf("%w: ageMonths = %v, want ≥ 0", ErrOutOfRange, ageMonths)
	}
	years := ageMonths / 12
	switch stage {
	case LifeStageNone:
	case LifeStagePregnancy, LifeStageLactation:
		if !female {
			return 0, fmt.Errorf("%w: stage = %d (pregnancy or lactation) with female = false", ErrDomain, stage)
		}
		if !(years >= 14 && years < 51) {
			return 0, fmt.Errorf("%w: ageMonths = %v, want 14 ≤ years < 51 for pregnancy or lactation", ErrOutOfRange, ageMonths)
		}
		if stage == LifeStagePregnancy {
			return driPregnancy, nil
		}
		return driLactation, nil
	default:
		return 0, fmt.Errorf("%w: stage = %d, want LifeStageNone, LifeStagePregnancy or LifeStageLactation", ErrOutOfRange, stage)
	}
	switch {
	case ageMonths < 7:
		return driInfant0to6mo, nil
	case ageMonths < 12:
		return driInfant7to12mo, nil
	case years < 4:
		return driChild1to3y, nil
	case years < 9:
		return driChild4to8y, nil
	case years < 14:
		return dri9to13y, nil
	case years < 19:
		return dri14to18y, nil
	default:
		return dri19yPlus, nil
	}
}

// WaterAI is a Table S-2 Adequate Intake for water, in liters per day: total water (from foods and beverages)
// and the part from beverages.
type WaterAI struct {
	TotalLPerDay, BeveragesLPerDay float64
}

// Table S-2, AI (L/d) Total Water and From Beverages, indexed by driGroup; [0] male, [1] female.
var waterAITable = [...][2]WaterAI{
	driInfant0to6mo:  {{0.7, 0.7}, {0.7, 0.7}},
	driInfant7to12mo: {{0.8, 0.6}, {0.8, 0.6}},
	driChild1to3y:    {{1.3, 0.9}, {1.3, 0.9}},
	driChild4to8y:    {{1.7, 1.2}, {1.7, 1.2}},
	dri9to13y:        {{2.4, 1.8}, {2.1, 1.6}},
	dri14to18y:       {{3.3, 2.6}, {2.3, 1.8}},
	dri19yPlus:       {{3.7, 3.0}, {2.7, 2.2}},
	driPregnancy:     {{3.0, 2.3}, {3.0, 2.3}}, // females only; the male cell is never reached
	driLactation:     {{3.8, 3.1}, {3.8, 3.1}}, // females only; the male cell is never reached
}

// TotalWaterAI returns the Adequate Intake (AI) for total water and for water from beverages, in L/day, for a
// life stage group, from IOM 2005 Table S-2. For 0 through 12 months the beverage AI is water from human milk
// (the table's criterion). The source states an AI "should not be interpreted as a specific requirement".
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range or for an unknown stage; ErrDomain for
// pregnancy or lactation with female false.
//
// ID:      LIF-NUTR-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Institute of Medicine. Dietary Reference Intakes for Water, Potassium, Sodium, Chloride, and Sulfate. Washington, DC: The National Academies Press; 2005. doi:10.17226/10925, Summary, Table S-2 "Criteria and Dietary Reference Intake Values for Total Water", pp. 6–7, read 2026-11-12 at nationalacademies.org/read/10925/chapter/2
// Source:  NASEM. Dietary Reference Intakes for Sodium and Potassium. 2019. doi:10.17226/25353, Appendix J, p. 572, Total Water (L/d) column: the same totals
// Tags:    nutrition, dri, water, hydration, adequate-intake, pregnancy, lactation
// Units:   ageMonths months (decimal; the caller's age); female true for the table's Female columns; result L/day
// Range:   ageMonths ≥ 0; stage LifeStagePregnancy or LifeStageLactation only for female and 14 ≤ ageMonths/12 < 51 ("14 through 50 y")
// Example: man, 30 y (360 months) → total 3.7 L/day, beverages 3.0 L/day; pregnancy, 28 y → 3.0 and 2.3 (Table S-2 rows "> 19 y" and "Pregnancy 14 through 50 y")
func TotalWaterAI(ageMonths float64, female bool, stage LifeStage) (WaterAI, error) {
	g, err := driLifeStageGroup(ageMonths, female, stage)
	if err != nil {
		return WaterAI{}, err
	}
	sex := 0
	if female {
		sex = 1
	}
	return waterAITable[g][sex], nil
}

package nutrition

// ProteinCarbohydrateDRI is a life stage group's Dietary Reference Intake for protein and carbohydrate, in
// grams per day. ProteinIsAI and CarbohydrateIsAI report that the value is an Adequate Intake (the table's
// asterisk) rather than a Recommended Dietary Allowance (the table's bold type).
type ProteinCarbohydrateDRI struct {
	ProteinGPerDay, CarbohydrateGPerDay float64
	ProteinIsAI, CarbohydrateIsAI       bool
}

// Appendix J p. 572, Carbohydrate (g/d) and Protein (g/d), indexed by driGroup; [0] male, [1] female.
var proteinCarbohydrateTable = [...][2]ProteinCarbohydrateDRI{
	driInfant0to6mo:  {{9.1, 60, true, true}, {9.1, 60, true, true}}, // 9.1*, 60*
	driInfant7to12mo: {{11.0, 95, false, true}, {11.0, 95, false, true}},
	driChild1to3y:    {{13, 130, false, false}, {13, 130, false, false}},
	driChild4to8y:    {{19, 130, false, false}, {19, 130, false, false}},
	dri9to13y:        {{34, 130, false, false}, {34, 130, false, false}},
	dri14to18y:       {{52, 130, false, false}, {46, 130, false, false}},
	dri19yPlus:       {{56, 130, false, false}, {46, 130, false, false}}, // 19–30, 31–50, 51–70, > 70 y
	driPregnancy:     {{71, 175, false, false}, {71, 175, false, false}}, // females only; male cell never reached
	driLactation:     {{71, 210, false, false}, {71, 210, false, false}}, // females only; male cell never reached
}

// ProteinCarbohydrateDRIFor returns the Dietary Reference Intake for protein and carbohydrate, in g/day, for a
// life stage group: the RDA, or the AI where the table marks one (protein 0–6 months; carbohydrate 0–12
// months). Rows: 0–6 mo 9.1* / 60*; 7–12 mo 11.0 / 95*; 1–3 y 13 / 130; 4–8 y 19 / 130; 9–13 y 34 / 130;
// males 14–18 y 52 / 130, 19 y and older 56 / 130; females 14–18 y and 19 y and older 46 / 130; pregnancy
// 71 / 175; lactation 71 / 210 (protein / carbohydrate g/d; * = AI).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range or for an unknown stage; ErrDomain for
// pregnancy or lactation with female false.
//
// ID:      LIF-NUTR-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  National Academies of Sciences, Engineering, and Medicine. Dietary Reference Intakes for Sodium and Potassium. Washington, DC: The National Academies Press; 2019. doi:10.17226/25353, Appendix J "Dietary Reference Intakes Summary Tables", p. 572, "Recommended Dietary Allowances and Adequate Intakes, Total Water and Macronutrients", Carbohydrate (g/d) and Protein (g/d) columns, read 2026-11-12 at nationalacademies.org/read/25353/chapter/28
// Tags:    nutrition, dri, rda, adequate-intake, protein, carbohydrate, pregnancy, lactation
// Units:   ageMonths months (decimal; the caller's age); female true for the table's Females rows; results g/day
// Range:   ageMonths ≥ 0; stage LifeStagePregnancy or LifeStageLactation only for female and 14 ≤ ageMonths/12 < 51 (the 14–18, 19–30 and 31–50 y rows)
// Example: male, 16 y → protein 52 g/day, carbohydrate 130 g/day (RDAs); 9 months → 11.0 g (RDA) and 95 g (AI) (Appendix J p. 572)
func ProteinCarbohydrateDRIFor(ageMonths float64, female bool, stage LifeStage) (ProteinCarbohydrateDRI, error) {
	g, err := driLifeStageGroup(ageMonths, female, stage)
	if err != nil {
		return ProteinCarbohydrateDRI{}, err
	}
	sex := 0
	if female {
		sex = 1
	}
	return proteinCarbohydrateTable[g][sex], nil
}

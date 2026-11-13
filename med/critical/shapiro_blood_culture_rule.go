package critical

// ShapiroFindings are the Shapiro rule's criteria besides age, each true when present as the abstract words it.
type ShapiroFindings struct {
	// Major criteria.
	TemperatureAbove39_5C bool // "temperature > 39.5 degrees C (103.0 degrees F)"
	IndwellingCatheter    bool // "indwelling vascular catheter"
	SuspectedEndocarditis bool // "clinical suspicion of endocarditis"
	// Minor criteria (age > 65 years is taken from ageYears).
	Temperature38_3To39_4C bool // "temperature 38.3-39.4 degrees C (101-102.9 degrees F)"
	Chills                 bool // "chills"
	Vomiting               bool // "vomiting"
	Hypotension            bool // "hypotension (systolic blood pressure < 90 mm Hg)"
	NeutrophilsAbove80Pct  bool // "neutrophil% > 80"
	WBCAbove18k            bool // "white blood cell count > 18 k"
	BandsAbove5Pct         bool // "bands > 5%"
	PlateletsBelow150k     bool // "platelets < 150 k"
	CreatinineAbove2       bool // "creatinine > 2.0" (the abstract gives no unit)
}

// ShapiroResult holds the counts of major and minor criteria present and the rule's result.
type ShapiroResult struct {
	Major     int  // 0–3
	Minor     int  // 0–10
	Indicated bool // "A blood culture is indicated by the rule if at least one major criterion or two minor criteria are present"
}

// ShapiroBloodCultureRule applies the Shapiro decision rule for blood cultures in adults in the emergency department
// with suspected infection: three major criteria and ten minor criteria (age > 65 years among them); the rule's
// result is true when at least one major or at least two minor criteria are present ("Otherwise, patients are
// classified as 'low risk'").
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when ageYears is outside 0–130.
//
// ID:      MED-CRIT-010
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Shapiro NI, Wolfe RE, Wright SB, Moore R, Bates DW. Who needs a blood culture? A prospectively derived and validated prediction rule. J Emerg Med. 2008;35(3):255-264. doi:10.1016/j.jemermed.2008.04.001. PMID 18486413. Abstract: "major criteria" defined as: temperature > 39.5 degrees C (103.0 degrees F), indwelling vascular catheter, or clinical suspicion of endocarditis. "Minor criteria" were: temperature 38.3-39.4 degrees C (101-102.9 degrees F), age > 65 years, chills, vomiting, hypotension (systolic blood pressure < 90 mm Hg), neutrophil% > 80, white blood cell count > 18 k, bands > 5%, platelets < 150 k, and creatinine > 2.0. A blood culture is indicated by the rule if at least one major criterion or two minor criteria are present.
// Tags:    critical, sepsis, bacteremia, blood culture, shapiro, decision rule
// Units:   ageYears years; result counts and boolean
// Range:   0 ≤ ageYears ≤ 130 (physical limits; adults in the ED with suspected infection, the study's population)
// Example: age 70 and chills → Minor 2, Indicated; indwelling catheter alone → Major 1, Indicated; age 65 and chills → Minor 1, not indicated (from the abstract)
func ShapiroBloodCultureRule(ageYears float64, f ShapiroFindings) (ShapiroResult, error) {
	if err := finite("ageYears", ageYears); err != nil {
		return ShapiroResult{}, err
	}
	if err := within("ageYears", ageYears, 0, 130); err != nil {
		return ShapiroResult{}, err
	}
	var r ShapiroResult
	for _, m := range [...]bool{f.TemperatureAbove39_5C, f.IndwellingCatheter, f.SuspectedEndocarditis} {
		if m {
			r.Major++
		}
	}
	for _, m := range [...]bool{f.Temperature38_3To39_4C, ageYears > 65, f.Chills, f.Vomiting, f.Hypotension,
		f.NeutrophilsAbove80Pct, f.WBCAbove18k, f.BandsAbove5Pct, f.PlateletsBelow150k, f.CreatinineAbove2} {
		if m {
			r.Minor++
		}
	}
	r.Indicated = r.Major >= 1 || r.Minor >= 2
	return r, nil
}

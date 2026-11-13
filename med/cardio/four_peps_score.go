package cardio

// FourPEPSFindings holds the 4PEPS items (Roy 2021, Table 3).
type FourPEPSFindings struct {
	AgeYears                  float64 // completed age, years: < 50: −2, 50–64: −1, ≥ 65: 0
	Male                      bool    // male sex: +2
	ChronicRespiratoryDisease bool    // chronic respiratory disease: −1
	HeartRateBpm              float64 // heart rate, beats/min: < 80: −1
	ChestPainAndAcuteDyspnea  bool    // chest pain and acute dyspnea: +1
	EstrogenicTreatment       bool    // hormonal estrogenic treatment: +2
	PersonalHistoryVTE        bool    // personal history of VTE: +2
	Syncope                   bool    // syncope: +2
	ImmobilityWithin4Weeks    bool    // immobility within last 4 weeks: +2
	SpO2Percent               float64 // pulse oxygen saturation, %: < 95: +3
	CalfPainOrUnilateralEdema bool    // calf pain and/or unilateral lower limb edema: +3
	PEMostLikelyDiagnosis     bool    // PE is the most likely diagnosis: +5
}

// FourPEPSScore returns the 4-Level Pulmonary Embolism Clinical Probability Score (4PEPS), −4 to 22.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for age < 18 years, heart rate < 0 or SpO2 outside
// 0–100 %.
//
// ID:      MED-CARDIO-064
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Roy PM, Friou E, Germeau B, et al. Derivation and Validation of a 4-Level Clinical Pretest Probability Score for Suspected Pulmonary Embolism to Safely Decrease Imaging Testing. JAMA Cardiol. 2021;6(6):669-677. doi:10.1001/jamacardio.2021.0064, PMID 33656522, PMC7931139, Table 3
// Tags:    cardio, pulmonary embolism, vte, 4peps, score
// Units:   AgeYears years; HeartRateBpm beats/min; SpO2Percent %; result points −4 to 22
// Range:   AgeYears ≥ 18 (adults suspected of PE, the source population)
// Example: man of 70, HR 95, SpO2 93 %, PE most likely → 0 + 2 + 3 + 5 = 10 (hand-derived from Table 3)
func FourPEPSScore(f FourPEPSFindings) (int, error) {
	for _, in := range [...]struct {
		name string
		v    float64
	}{{"AgeYears", f.AgeYears}, {"HeartRateBpm", f.HeartRateBpm}, {"SpO2Percent", f.SpO2Percent}} {
		if err := checkFinite(in.name, in.v); err != nil {
			return 0, err
		}
	}
	if err := checkRange("AgeYears", f.AgeYears, 18, 150); err != nil {
		return 0, err
	}
	if err := checkRange("HeartRateBpm", f.HeartRateBpm, 0, 400); err != nil {
		return 0, err
	}
	if err := checkRange("SpO2Percent", f.SpO2Percent, 0, 100); err != nil {
		return 0, err
	}
	s := 0
	switch {
	case f.AgeYears < 50:
		s -= 2
	case f.AgeYears < 65:
		s--
	}
	if f.ChronicRespiratoryDisease {
		s--
	}
	if f.HeartRateBpm < 80 {
		s--
	}
	if f.ChestPainAndAcuteDyspnea {
		s++
	}
	for _, two := range [...]bool{f.Male, f.EstrogenicTreatment, f.PersonalHistoryVTE, f.Syncope, f.ImmobilityWithin4Weeks} {
		if two {
			s += 2
		}
	}
	if f.SpO2Percent < 95 {
		s += 3
	}
	if f.CalfPainOrUnilateralEdema {
		s += 3
	}
	if f.PEMostLikelyDiagnosis {
		s += 5
	}
	return s, nil
}

// FourPEPSProbability is a 4PEPS clinical pretest probability (CPP) level; values are the source's labels.
type FourPEPSProbability string

const (
	FourPEPSVeryLow  FourPEPSProbability = "Very low CPP" // < 0
	FourPEPSLow      FourPEPSProbability = "Low CPP"      // 0–5
	FourPEPSModerate FourPEPSProbability = "Moderate CPP" // 6–12
	FourPEPSHigh     FourPEPSProbability = "High CPP"     // ≥ 13
)

// FourPEPSLevel returns the 4PEPS clinical pretest probability level for a score: "Very low CPP (<0)",
// "Low CPP (0-5)", "Moderate CPP (6-12)", "High CPP (≥13)". The table's testing strategy per level is not
// returned.
// Errors: ErrOutOfRange for a score outside −4 to 22.
//
// ID:      MED-CARDIO-065
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Roy PM, Friou E, Germeau B, et al. JAMA Cardiol. 2021;6(6):669-677. doi:10.1001/jamacardio.2021.0064, PMID 33656522, PMC7931139, Table 3 (CPP levels)
// Tags:    cardio, pulmonary embolism, vte, 4peps, category
// Units:   score 4PEPS points; result the source's label
// Range:   score −4 to 22
// Example: −1 → Very low CPP; 0 → Low CPP; 6 → Moderate CPP; 13 → High CPP (Table 3)
func FourPEPSLevel(score int) (FourPEPSProbability, error) {
	switch {
	case score < -4 || score > 22:
		return "", errScore(score, -4, 22)
	case score < 0:
		return FourPEPSVeryLow, nil
	case score <= 5:
		return FourPEPSLow, nil
	case score <= 12:
		return FourPEPSModerate, nil
	}
	return FourPEPSHigh, nil
}

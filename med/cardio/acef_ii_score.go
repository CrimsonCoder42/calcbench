package cardio

// ACEFII returns the ACEF II score for cardiac surgery: age (years) / ejection fraction (%), plus 2 when serum
// creatinine is above 2 mg/dL, plus 3 for emergency surgery, plus 0.2 per haematocrit point below 36%.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-CARDIO-034
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ranucci M, Pistuddi V, Scolletta S, de Vincentiis C, Menicanti L. The ACEF II Risk Score for cardiac surgery: updated but still parsimonious. Eur Heart J. 2018;39(23):2183-2189. doi:10.1093/eurheartj/ehx228. PMID 28498904. Abstract: "age(years)/ejection fraction (%). Additional points are attributed to a serum creatinine level > 2 mg/dL (2 points), emergency surgery (3 points) and anaemia [haematocrit (HCT) < 36%, 0.2 points per each HCT point below 36%]"
// Tags:    cardio, cardiac-surgery, mortality, acef, ranucci, risk-score
// Units:   ageYears years; ejectionFractionPct %; creatinineMgDL mg/dL; hematocritPct %; result points
// Range:   ageYears 18–120; ejectionFractionPct 1–100; creatinineMgDL > 0; hematocritPct 1–100 (adult cardiac surgery; no limits stated in the abstract)
// Example: age 70, EF 50%, creatinine 1.0, elective, HCT 33% → 1.4 + 0.6 = 2.0 (HAND-DERIVED)
func ACEFII(ageYears, ejectionFractionPct, creatinineMgDL float64, emergency bool, hematocritPct float64) (float64, error) {
	for _, p := range [...]struct {
		n      string
		v      float64
		lo, hi float64
	}{{"ageYears", ageYears, 18, 120}, {"ejectionFractionPct", ejectionFractionPct, 1, 100},
		{"creatinineMgDL", creatinineMgDL, 0, maxFinite}, {"hematocritPct", hematocritPct, 1, 100}} {
		if err := checkFinite(p.n, p.v); err != nil {
			return 0, err
		}
		if err := checkRange(p.n, p.v, p.lo, p.hi); err != nil {
			return 0, err
		}
	}
	if err := checkPositive("creatinineMgDL", creatinineMgDL); err != nil {
		return 0, err
	}
	s := ageYears / ejectionFractionPct
	if creatinineMgDL > 2 {
		s += 2
	}
	if emergency {
		s += 3
	}
	if hematocritPct < 36 {
		s += 0.2 * (36 - hematocritPct)
	}
	return s, nil
}

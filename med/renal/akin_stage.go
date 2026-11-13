package renal

import (
	"fmt"
	"math"
)

// AKINStage returns the Acute Kidney Injury Network stage, 0–3 (0: no stage criterion met), as the higher of
// the serum-creatinine stage and the urine-output stage of Mehta 2007 Table 2, and 3 for anyone receiving renal
// replacement therapy. Creatinine: stage 1 an increase ≥ 0.3 mg/dL or to ≥ 150% of baseline; stage 2 to > 200%;
// stage 3 to > 300%, or a value ≥ 4.0 mg/dL with an acute increase of at least 0.5 mg/dL. Urine output (mL/kg
// per hour, sustained for urineOutputHours): stage 1 < 0.5 for more than 6 h; stage 2 < 0.5 for more than 12 h;
// stage 3 < 0.3 for 24 h, or anuria (0) for 12 h. The change is the caller's within the source's 48-hour window.
// Creatinine values are compared at 0.001 mg/dL resolution (each is rounded to the nearest 0.001 mg/dL first) so
// that decimal inputs such as 0.9 and 1.2 compare as written; the source states its thresholds to 0.1 mg/dL.
// Errors: ErrNotFinite for NaN or ±Inf; ErrDomain when a creatinine × 1000 overflows float64; ErrOutOfRange for
// a creatinine not above zero or a negative urine
// output or duration.
//
// ID:      MED-RENAL-012
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Mehta RL, Kellum JA, Shah SV, et al. Acute Kidney Injury Network: report of an initiative to improve outcomes in acute kidney injury. Crit Care. 2007;11(2):R31. doi:10.1186/cc5713. PMID 17331245, PMC2206446, Table 2 "Classification/staging system for acute kidney injury" and its footnotes ("Only one criterion (creatinine or urine output) has to be fulfilled to qualify for a stage"; "individuals who receive RRT are considered to have met the criteria for stage 3 irrespective of the stage they are in at the time of RRT")
// Tags:    renal, acute kidney injury, akin, staging, creatinine, urine output
// Units:   baselineCreatinineMgDL, currentCreatinineMgDL mg/dL; urineOutputMLPerKgPerH mL/kg/h; urineOutputHours h; result stage 0–3
// Range:   creatinine > 0; urineOutputMLPerKgPerH ≥ 0; urineOutputHours ≥ 0
// Example: baseline 1.0, current 1.3, urine 1.0 mL/kg/h for 24 h → 1; baseline 1.0, current 2.1 → 2; baseline 3.6, current 4.1 → 3 (≥ 4.0 with rise ≥ 0.5); urine 0.4 for 13 h → 2; urine 0 for 12 h → 3 (HAND-DERIVED from Table 2)
func AKINStage(baselineCreatinineMgDL, currentCreatinineMgDL, urineOutputMLPerKgPerH, urineOutputHours float64, renalReplacementTherapy bool) (int, error) {
	for _, e := range [...]error{
		checkFinite("baselineCreatinineMgDL", baselineCreatinineMgDL), checkFinite("currentCreatinineMgDL", currentCreatinineMgDL),
		checkFinite("urineOutputMLPerKgPerH", urineOutputMLPerKgPerH), checkFinite("urineOutputHours", urineOutputHours),
		checkPositive("baselineCreatinineMgDL", baselineCreatinineMgDL), checkPositive("currentCreatinineMgDL", currentCreatinineMgDL),
	} {
		if e != nil {
			return 0, e
		}
	}
	if urineOutputMLPerKgPerH < 0 {
		return 0, fmt.Errorf("%w: urineOutputMLPerKgPerH = %v, want ≥ 0", ErrOutOfRange, urineOutputMLPerKgPerH)
	}
	if urineOutputHours < 0 {
		return 0, fmt.Errorf("%w: urineOutputHours = %v, want ≥ 0", ErrOutOfRange, urineOutputHours)
	}
	if renalReplacementTherapy {
		return 3, nil
	}
	// Thousandths of a mg/dL: integers held exactly in float64 for any physiological value.
	b := math.Round(baselineCreatinineMgDL * 1000)
	c := math.Round(currentCreatinineMgDL * 1000)
	if math.IsInf(b, 0) || math.IsInf(c, 0) {
		return 0, fmt.Errorf("%w: a creatinine × 1000 overflows float64", ErrDomain)
	}
	cr := 0
	switch {
	case c > 3*b || (c >= 4000 && c-b >= 500):
		cr = 3
	case c > 2*b:
		cr = 2
	case c-b >= 300 || 2*c >= 3*b:
		cr = 1
	}
	uo := 0
	switch r, h := urineOutputMLPerKgPerH, urineOutputHours; {
	case (r < 0.3 && h >= 24) || (r == 0 && h >= 12):
		uo = 3
	case r < 0.5 && h > 12:
		uo = 2
	case r < 0.5 && h > 6:
		uo = 1
	}
	return max(cr, uo), nil
}

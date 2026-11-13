package renal

import (
	"fmt"
	"math"
)

// RIFLEClass is a class of the ADQI RIFLE classification, labelled as Bellomo 2004 labels it.
type RIFLEClass string

// RIFLE classes (Bellomo 2004 Figure 1 and text), worst last.
const (
	RIFLENone    RIFLEClass = ""        // no RIFLE criterion met
	RIFLERisk    RIFLEClass = "Risk"    // creatinine ≥ 1.5 × baseline or urine < 0.5 mL/kg/h × 6 h
	RIFLEInjury  RIFLEClass = "Injury"  // creatinine ≥ 2 × baseline or urine < 0.5 mL/kg/h × 12 h
	RIFLEFailure RIFLEClass = "Failure" // creatinine ≥ 3 × baseline, or > 4.0 mg/dL with an acute rise ≥ 0.5, or urine < 0.3 × 24 h or anuria × 12 h
	RIFLELoss    RIFLEClass = "Loss"    // persistent acute renal failure: renal replacement therapy > 4 weeks
	RIFLEESRD    RIFLEClass = "ESRD"    // end-stage renal disease: dialysis > 3 months
)

// RIFLEUrineOutput is the urine-output criterion of RIFLE met by the caller's observation.
type RIFLEUrineOutput int

// RIFLE urine-output criteria.
const (
	RIFLEUrineNone                     RIFLEUrineOutput = iota // no urine-output criterion met
	RIFLEUrineBelow05For6h                                     // < 0.5 mL/kg/h for 6 h (Risk)
	RIFLEUrineBelow05For12h                                    // < 0.5 mL/kg/h for 12 h (Injury)
	RIFLEUrineBelow03For24hOrAnuria12h                         // < 0.3 mL/kg/h for 24 h or anuria for 12 h (Failure)
)

// RIFLEClassify returns the RIFLE class from serum creatinine (current against baseline, mg/dL), the urine-output
// criterion met, and the two outcome classes: dialysisOver3Months → ESRD; rrtOver4Weeks → Loss; otherwise the
// worst of the creatinine and urine-output classes ("The criteria that lead to the worst possible classification
// should be used"). Creatinine: Risk ≥ 1.5 ×, Injury ≥ 2 ×, Failure ≥ 3 × baseline, or Failure when the value is
// above 4.0 mg/dL "in the setting of an acute increase of at least 0.5 mg/dl". The GFR-decrease arm (> 25%, > 50%,
// > 75%) is not built: it was seen in one secondary only. Creatinine values are compared at 0.001 mg/dL
// resolution, as MED-RENAL-012 does.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a creatinine not above zero or an unknown urine
// criterion; ErrDomain when a creatinine × 1000 overflows float64.
//
// ID:      MED-RENAL-025
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Bellomo R, Ronco C, Kellum JA, Mehta RL, Palevsky P; ADQI workgroup. Acute renal failure - definition, outcome measures, animal models, fluid therapy and information technology needs: the Second International Consensus Conference of the Acute Dialysis Quality Initiative (ADQI) Group. Crit Care. 2004;8(4):R204-R212. doi:10.1186/cc2872. PMID 15312219. PMC522841, text and Figure 1 legend (Loss, ESRD, the > 4.0 mg/dL rule, worst criterion); thresholds of Figure 1 (an image) read in the two sources below
// Source:  Hoste EA, Clermont G, Kersten A, et al. RIFLE criteria for acute kidney injury are associated with hospital mortality in critically ill patients: a cohort analysis. Crit Care. 2006;10(3):R73. doi:10.1186/cc4915. PMID 16696865. PMC1550961, Table 1 (RIFLE, after Bellomo 2004)
// Source:  Wlodzimirow KA, Abu-Hanna A, Slabbekoorn M, et al. A comparison of RIFLE with and without urine output criteria for acute kidney injury in critically ill patients. Crit Care. 2012;16(5):R200. doi:10.1186/cc11808. PMID 23078781. PMC3682302 (RIFLE criteria, after Bellomo 2004)
// Verified: secondary (Hoste 2006 Table 1; Wlodzimirow 2012)
// Tags:    renal, acute kidney injury, rifle, adqi, creatinine, urine output
// Units:   baselineCreatinineMgDL, currentCreatinineMgDL mg/dL; result the source's class label
// Range:   creatinine > 0
// Example: baseline 1.0, current 1.6 → Risk; 2.0 → Injury; 3.0 → Failure; baseline 3.6, current 4.1 → Failure (> 4.0 with rise ≥ 0.5) (HAND-DERIVED)
func RIFLEClassify(baselineCreatinineMgDL, currentCreatinineMgDL float64, urine RIFLEUrineOutput, rrtOver4Weeks, dialysisOver3Months bool) (RIFLEClass, error) {
	for _, e := range [...]error{
		checkFinite("baselineCreatinineMgDL", baselineCreatinineMgDL), checkFinite("currentCreatinineMgDL", currentCreatinineMgDL),
		checkPositive("baselineCreatinineMgDL", baselineCreatinineMgDL), checkPositive("currentCreatinineMgDL", currentCreatinineMgDL),
	} {
		if e != nil {
			return RIFLENone, e
		}
	}
	return rifleClass(baselineCreatinineMgDL, currentCreatinineMgDL, urine, rrtOver4Weeks, dialysisOver3Months)
}

// RIFLEClassifySI returns the RIFLE class (MED-RENAL-025) from creatinine in µmol/L: each value ÷ 88.4 gives mg/dL,
// then the same rules.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a creatinine not above zero or an unknown urine
// criterion; ErrDomain when a creatinine × 1000 overflows float64.
//
// ID:      MED-RENAL-026
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Bellomo R, Ronco C, Kellum JA, Mehta RL, Palevsky P; ADQI workgroup. Crit Care. 2004;8(4):R204-R212. doi:10.1186/cc2872. PMID 15312219. PMC522841 (as MED-RENAL-025)
// Source:  Kidney Disease: Improving Global Outcomes (KDIGO) CKD Work Group. KDIGO 2024 Clinical Practice Guideline for the Evaluation and Management of Chronic Kidney Disease. Kidney Int. 2024;105(4S):S117–S314. doi:10.1016/j.kint.2023.10.018. p. S126, "Conversion factors of conventional units to SI units" (creatinine 88.4)
// Source:  Hoste EA, Clermont G, Kersten A, et al. Crit Care. 2006;10(3):R73. doi:10.1186/cc4915 (as MED-RENAL-025)
// Source:  Wlodzimirow KA, Abu-Hanna A, Slabbekoorn M, et al. Crit Care. 2012;16(5):R200. doi:10.1186/cc11808 (as MED-RENAL-025)
// Verified: secondary (Hoste 2006 Table 1; Wlodzimirow 2012)
// Tags:    renal, acute kidney injury, rifle, adqi, creatinine, urine output, si
// Units:   baselineCreatinineUmolL, currentCreatinineUmolL µmol/L; result the source's class label
// Range:   creatinine > 0
// Example: baseline 88.4, current 176.8 µmol/L (1.0 → 2.0 mg/dL) → Injury (HAND-DERIVED)
func RIFLEClassifySI(baselineCreatinineUmolL, currentCreatinineUmolL float64, urine RIFLEUrineOutput, rrtOver4Weeks, dialysisOver3Months bool) (RIFLEClass, error) {
	for _, e := range [...]error{
		checkFinite("baselineCreatinineUmolL", baselineCreatinineUmolL), checkFinite("currentCreatinineUmolL", currentCreatinineUmolL),
		checkPositive("baselineCreatinineUmolL", baselineCreatinineUmolL), checkPositive("currentCreatinineUmolL", currentCreatinineUmolL),
	} {
		if e != nil {
			return RIFLENone, e
		}
	}
	return rifleClass(baselineCreatinineUmolL/creatinineUmolLPerMgDL, currentCreatinineUmolL/creatinineUmolLPerMgDL, urine, rrtOver4Weeks, dialysisOver3Months)
}

// rifleClass applies the rules to creatinine in mg/dL.
func rifleClass(baselineMgDL, currentMgDL float64, urine RIFLEUrineOutput, rrtOver4Weeks, dialysisOver3Months bool) (RIFLEClass, error) {
	if urine < RIFLEUrineNone || urine > RIFLEUrineBelow03For24hOrAnuria12h {
		return RIFLENone, fmt.Errorf("%w: urine = %d, want 0–3", ErrOutOfRange, urine)
	}
	if dialysisOver3Months {
		return RIFLEESRD, nil
	}
	if rrtOver4Weeks {
		return RIFLELoss, nil
	}
	// Thousandths of a mg/dL, as MED-RENAL-012.
	b := math.Round(baselineMgDL * 1000)
	c := math.Round(currentMgDL * 1000)
	if math.IsInf(b, 0) || math.IsInf(c, 0) {
		return RIFLENone, fmt.Errorf("%w: a creatinine × 1000 overflows float64", ErrDomain)
	}
	cr := 0
	switch {
	case c >= 3*b || (c > 4000 && c-b >= 500):
		cr = 3
	case c >= 2*b:
		cr = 2
	case 2*c >= 3*b:
		cr = 1
	}
	if u := int(urine); u > cr {
		cr = u
	}
	return [...]RIFLEClass{RIFLENone, RIFLERisk, RIFLEInjury, RIFLEFailure}[cr], nil
}

package trauma

import "fmt"

// PECARNCSpineInput holds the findings of the PECARN cervical spine injury prediction rule (Leonard 2024):
// the four high-risk factors and the five CART-identified factors. Each is true when present.
type PECARNCSpineInput struct {
	AgeYears float64 // the child's age; the rule enrolled "Children <18 years"

	// High-risk factors.
	GCS3To8                            bool // Glasgow Coma Scale score of 3–8
	AVPUUnresponsive                   bool // AVPU scale of U
	AbnormalAirwayBreathingCirculation bool // abnormal airway, breathing or circulation findings
	FocalNeurologicDeficit             bool // focal neurologic deficits including paresthesia, numbness or weakness

	// CART-identified factors.
	NeckPain               bool // neck pain
	AlteredMentalStatus    bool // GCS score of 9–14; AVPU scale of V or P; or other signs of altered mental status
	SubstantialHeadInjury  bool // an injury warranting inpatient observation or surgical intervention
	SubstantialTorsoInjury bool // an injury warranting inpatient observation or surgical intervention
	MidlineNeckTenderness  bool // (posterior) midline neck tenderness
}

// PECARNCSpineGroup is the group the rule's factors place a child in, named with the source's own terms.
type PECARNCSpineGroup string

const (
	// PECARNCSpineHighRisk: one or more of the four high-risk factors.
	PECARNCSpineHighRisk PECARNCSpineGroup = "high-risk factor present"
	// PECARNCSpineCART: no high-risk factor and one or more of the five CART-identified factors.
	PECARNCSpineCART PECARNCSpineGroup = "CART-identified factor present, no high-risk factor"
	// PECARNCSpineNone: none of the nine factors.
	PECARNCSpineNone PECARNCSpineGroup = "no high-risk or CART-identified factor"
)

// PECARNCSpineRule returns the group of the PECARN cervical spine injury prediction rule for children after
// blunt trauma: PECARNCSpineHighRisk when any of the four high-risk factors is present (GCS 3–8, AVPU U,
// abnormal airway/breathing/circulation, focal neurologic deficit); otherwise PECARNCSpineCART when any of the
// five CART-identified factors is present (neck pain, altered mental status, substantial head injury,
// substantial torso injury, midline neck tenderness); otherwise PECARNCSpineNone. The rule is positive in
// the first two groups ("combining the four high-risk factors and the five CART-identified factors").
// Errors: ErrNotFinite for a NaN or ±Inf AgeYears; ErrOutOfRange unless 0 ≤ AgeYears < 18.
//
// ID:      MED-TRAUMA-023
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Leonard JC, Harding M, Cook LJ, et al. PECARN prediction rule for cervical spine imaging of children presenting to the emergency department with blunt trauma: a multicentre prospective observational study. Lancet Child Adolesc Health. 2024;8(7):482-490. doi:10.1016/S2352-4642(24)00104-4. PMID 38843852, PMC11261431 (author manuscript): Summary Findings, Results ("combining the four high-risk factors and the five CART-identified factors"), Figure 2, Methods ("Children <18 years experiencing blunt trauma")
// Tags:    trauma, pecarn, cervical spine, c-spine, pediatric, blunt trauma, decision rule
// Units:   AgeYears years; findings true/false; result the group label
// Range:   0 ≤ AgeYears < 18 (the enrolled population)
// Example: neck pain only → PECARNCSpineCART; neck pain and GCS 7 → PECARNCSpineHighRisk; nothing → PECARNCSpineNone (from the rule's definition)
func PECARNCSpineRule(in PECARNCSpineInput) (PECARNCSpineGroup, error) {
	if err := finite("AgeYears", in.AgeYears); err != nil {
		return "", err
	}
	if in.AgeYears < 0 || in.AgeYears >= 18 {
		return "", fmt.Errorf("%w: AgeYears = %v, want 0 ≤ age < 18", ErrOutOfRange, in.AgeYears)
	}
	if in.GCS3To8 || in.AVPUUnresponsive || in.AbnormalAirwayBreathingCirculation || in.FocalNeurologicDeficit {
		return PECARNCSpineHighRisk, nil
	}
	if in.NeckPain || in.AlteredMentalStatus || in.SubstantialHeadInjury || in.SubstantialTorsoInjury || in.MidlineNeckTenderness {
		return PECARNCSpineCART, nil
	}
	return PECARNCSpineNone, nil
}

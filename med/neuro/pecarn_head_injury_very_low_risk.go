package neuro

// PECARNHeadFindings are the predictors of the PECARN head-trauma rules (Kuppermann 2009 abstract). Each field
// is true when the finding named is present; the rule for the child's age reads only its own six.
type PECARNHeadFindings struct {
	AbnormalMentalStatus     bool // both ages: the opposite of "normal mental status"
	SevereInjuryMechanism    bool // both ages: the opposite of "non-severe injury mechanism"
	NonFrontalScalpHematoma  bool // younger than 2: "no scalp haematoma except frontal"
	LOCFiveSecondsOrMore     bool // younger than 2: "no loss of consciousness or loss of consciousness for less than 5 s"
	PalpableSkullFracture    bool // younger than 2: "no palpable skull fracture"
	NotActingNormally        bool // younger than 2: "acting normally"
	AnyLossOfConsciousness   bool // 2 and older: "no loss of consciousness"
	Vomiting                 bool // 2 and older: "no vomiting"
	BasilarSkullFractureSign bool // 2 and older: "no signs of basilar skull fracture"
	SevereHeadache           bool // 2 and older: "no severe headache"
}

// PECARNHeadInjuryVeryLowRisk reports whether a child meets the PECARN prediction rule for very low risk of
// clinically important traumatic brain injury: under 2 years, normal mental status, no scalp hematoma except
// frontal, no loss of consciousness or one under 5 s, non-severe mechanism, no palpable skull fracture and
// acting normally; 2 years and older, normal mental status, no loss of consciousness, no vomiting, non-severe
// mechanism, no signs of basilar skull fracture and no severe headache. The rules were derived for children
// presenting within 24 h of head trauma with a Glasgow Coma Scale score of 14–15; that is the caller's.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when ageYears is outside 0 ≤ age < 18.
//
// ID:      MED-NEURO-022
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Kuppermann N, Holmes JF, Dayan PS, et al. Identification of children at very low risk of clinically-important brain injuries after head trauma: a prospective cohort study. Lancet. 2009;374(9696):1160-1170. doi:10.1016/S0140-6736(09)61558-0, PMID 19758692, abstract (the two six-item rules, quoted on PECARNHeadFindings; population "younger than 18 years ... within 24 h of head trauma with Glasgow Coma Scale scores of 14-15")
// Tags:    neuro, head injury, pediatric, pecarn, decision rule
// Units:   ageYears years (decimal); findings yes/no; result true when every criterion of the age's rule is met
// Range:   0 ≤ ageYears < 18 (the source's population; GCS 14–15 within 24 h is the caller's)
// Example: age 1, nothing present → true; age 1, frontal-only hematoma → true; age 1, non-frontal hematoma → false; age 5, vomiting → false (HAND-DERIVED from the abstract)
func PECARNHeadInjuryVeryLowRisk(ageYears float64, f PECARNHeadFindings) (bool, error) {
	if err := checkChildAge(ageYears); err != nil {
		return false, err
	}
	if f.AbnormalMentalStatus || f.SevereInjuryMechanism {
		return false, nil
	}
	if ageYears < 2 {
		return !(f.NonFrontalScalpHematoma || f.LOCFiveSecondsOrMore || f.PalpableSkullFracture || f.NotActingNormally), nil
	}
	return !(f.AnyLossOfConsciousness || f.Vomiting || f.BasilarSkullFractureSign || f.SevereHeadache), nil
}

package obstetric

import "fmt"

// IOTASimpleRulesResult is the classification of the IOTA simple rules, with the source's words.
type IOTASimpleRulesResult string

// The three outcomes of the IOTA simple rules (Timmerman 2010, Methods), labels verbatim.
const (
	IOTAMalignant    IOTASimpleRulesResult = "malignant"    // rule 1: ≥ 1 M feature and no B feature
	IOTABenign       IOTASimpleRulesResult = "benign"       // rule 2: ≥ 1 B feature and no M feature
	IOTAInconclusive IOTASimpleRulesResult = "inconclusive" // rule 3: both kinds present, or none
)

// iotaFeaturesPerKind is the number of M features (M1–M5) and of B features (B1–B5) the rules define.
const iotaFeaturesPerKind = 5

// IOTASimpleRules returns the IOTA simple-rules classification of an adnexal mass from the number of M
// features (M1–M5) and B features (B1–B5) present: "malignant" when one or more M features are present and
// no B feature, "benign" when one or more B features are present and no M feature, and "inconclusive" when
// both kinds are present or none is. Which features are present is judged by the caller from the ultrasound
// examination; the features are M1 irregular solid tumour, M2 ascites, M3 at least four papillary
// structures, M4 irregular multilocular solid tumour with a largest diameter of at least 100 mm, M5 very high
// colour content on colour Doppler examination; B1 unilocular cyst, B2 solid components whose largest is
// < 7 mm in largest diameter, B3 acoustic shadows, B4 smooth multilocular tumour, B5 no detectable blood flow
// on Doppler examination. The 2016 risk-of-malignancy model built on these features is a different
// calculation and is not this function.
// Errors: ErrOutOfRange when a count is outside 0–5.
//
// ID:      MED-OB-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Timmerman D, Ameye L, Fischerova D, et al. Simple ultrasound rules to distinguish between benign and malignant adnexal masses before surgery: prospective validation by IOTA group. BMJ. 2010;341:c6839. doi:10.1136/bmj.c6839. PMID 21156740, PMC3001703, Methods: "If one or more M features were present in the absence of a B feature, we classified the mass as malignant (rule 1). If one or more B features were present in the absence of an M feature, we classified the mass as benign (rule 2). If both M features and B features were present, or if none of the features was present, the simple rules were inconclusive (rule 3)."
// Source:  Timmerman D, Testa AC, Bourne T, et al. Simple ultrasound-based rules for the diagnosis of ovarian cancer. Ultrasound Obstet Gynecol. 2008;31(6):681-690. doi:10.1002/uog.5365. PMID 18504770 (derivation; abstract lists the five M and five B features)
// Tags:    obstetric, gynecology, adnexal-mass, ovarian, ultrasound, iota, simple-rules
// Units:   mFeatures count of M1–M5 present; bFeatures count of B1–B5 present; result category
// Range:   mFeatures 0–5; bFeatures 0–5 (adnexal masses examined before surgery, the population of the source)
// Example: M 1, B 0 → "malignant"; M 0, B 2 → "benign"; M 1, B 1 → "inconclusive"; M 0, B 0 → "inconclusive" (rules 1–3, BMJ 2010)
func IOTASimpleRules(mFeatures, bFeatures int) (IOTASimpleRulesResult, error) {
	if mFeatures < 0 || mFeatures > iotaFeaturesPerKind {
		return "", fmt.Errorf("%w: mFeatures = %d, want 0–5", ErrOutOfRange, mFeatures)
	}
	if bFeatures < 0 || bFeatures > iotaFeaturesPerKind {
		return "", fmt.Errorf("%w: bFeatures = %d, want 0–5", ErrOutOfRange, bFeatures)
	}
	switch {
	case mFeatures > 0 && bFeatures == 0:
		return IOTAMalignant, nil
	case bFeatures > 0 && mFeatures == 0:
		return IOTABenign, nil
	default:
		return IOTAInconclusive, nil
	}
}

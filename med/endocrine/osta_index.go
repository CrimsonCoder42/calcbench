package endocrine

import (
	"fmt"
	"math"
)

// ostaKgPerPound converts pounds to kilograms; exact by the 1959 international yard and pound agreement
// (NIST SP 811 (2008) §B.8).
const ostaKgPerPound = 0.45359237

// OSTAIndex returns the Osteoporosis Self-assessment Tool for Asians index: "(body weight in kg − age in
// years) × 0.2", "truncated to an integer" (toward zero).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for weight ≤ 0 or age < 0.
//
// ID:      MED-ENDO-029
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Koh LK, Sedrine WB, Torralba TP, et al. A simple tool to identify Asian women at increased risk of osteoporosis. Osteoporos Int. 2001;12(8):699-705. doi:10.1007/s001980170070, PMID 11580084 (primary; not opened)
// Source:  Subramaniam S, Ima-Nirwana S, Chin KY. Performance of Osteoporosis Self-Assessment Tool (OST) in Predicting Osteoporosis: A Review. Int J Environ Res Public Health. 2018;15(7):1445. doi:10.3390/ijerph15071445, PMID 29987247, PMC6068473, text citing ref 11
// Source:  Lu C, Chen D, Cai Y, Wei S. Concordance of OSTA and lumbar spine BMD by DXA in identifying risk of osteoporosis. J Orthop Surg Res. 2006;1:14. doi:10.1186/1749-799X-1-14, PMID 17150121, PMC1693545, Methods
// Verified: secondary (Subramaniam 2018 PMC6068473; Lu 2006 PMC1693545)
// Tags:    endocrine, bone, osteoporosis, osta, ost, index
// Units:   weightKg kg; ageYears years; result index points (integer)
// Range:   postmenopausal Asian women (the derivation population); weightKg > 0, ageYears ≥ 0
// Example: 50 kg, 65 years → trunc(0.2 × (50 − 65)) = trunc(−3.0) = −3 (hand-derived)
func OSTAIndex(weightKg, ageYears float64) (int, error) {
	if err := checkFinite("weightKg", weightKg); err != nil {
		return 0, err
	}
	if err := checkFinite("ageYears", ageYears); err != nil {
		return 0, err
	}
	if weightKg <= 0 {
		return 0, fmt.Errorf("%w: weightKg = %v, want > 0", ErrOutOfRange, weightKg)
	}
	if ageYears < 0 {
		return 0, fmt.Errorf("%w: ageYears = %v, want ≥ 0", ErrOutOfRange, ageYears)
	}
	return int(math.Trunc(0.2 * (weightKg - ageYears))), nil
}

// OSTAIndexImperial returns the OSTA index for a weight in pounds: the weight is converted to kilograms once
// (× 0.45359237, exact) and OSTAIndex applies. The sources state no pound form.
// Errors: as OSTAIndex, with the parameter named weightLb.
//
// ID:      MED-ENDO-030
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Koh LK, Sedrine WB, Torralba TP, et al. Osteoporos Int. 2001;12(8):699-705. doi:10.1007/s001980170070, PMID 11580084 (formula in kg; see OSTAIndex)
// Source:  NIST Special Publication 811 (2008), §B.8 (1 lb = 0.45359237 kg exactly)
// Source:  Subramaniam S, Ima-Nirwana S, Chin KY. Int J Environ Res Public Health. 2018;15(7):1445. doi:10.3390/ijerph15071445, PMID 29987247, PMC6068473 (as MED-ENDO-029)
// Source:  Lu C, Chen D, Cai Y, Wei S. J Orthop Surg Res. 2006;1:14. doi:10.1186/1749-799X-1-14, PMID 17150121, PMC1693545 (as MED-ENDO-029)
// Verified: secondary (Subramaniam 2018 PMC6068473; Lu 2006 PMC1693545)
// Tags:    endocrine, bone, osteoporosis, osta, ost, index, imperial
// Units:   weightLb lb; ageYears years; result index points (integer)
// Range:   as OSTAIndex; weightLb > 0
// Example: 130 lb, 65 years → 58.967 kg → trunc(−1.207) = −1 (hand-derived)
func OSTAIndexImperial(weightLb, ageYears float64) (int, error) {
	if err := checkFinite("weightLb", weightLb); err != nil {
		return 0, err
	}
	if weightLb <= 0 {
		return 0, fmt.Errorf("%w: weightLb = %v, want > 0", ErrOutOfRange, weightLb)
	}
	return OSTAIndex(weightLb*ostaKgPerPound, ageYears)
}

// OSTARisk is an OSTA risk category; values are the review's words for Koh 2001's three groups.
type OSTARisk string

const (
	OSTALowRisk      OSTARisk = "low-risk"      // index > −1
	OSTAModerateRisk OSTARisk = "moderate-risk" // −1 to −4
	OSTAHighRisk     OSTARisk = "high-risk"     // < −4
)

// OSTARiskCategory returns the OSTA risk category of an index: "low-risk (>−1), moderate-risk (−1 to −4),
// and high-risk (<−4)". The index is an integer, so −1 and −4 are moderate-risk. These are the women's bands;
// the male bands of Adler 2003 (≥ 4, −1 to 3, ≤ −2) are different and are not built here.
//
// ID:      MED-ENDO-031
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Koh LK, Sedrine WB, Torralba TP, et al. Osteoporos Int. 2001;12(8):699-705. doi:10.1007/s001980170070, PMID 11580084, abstract (high-, intermediate- and low-risk groups)
// Source:  Subramaniam S, Ima-Nirwana S, Chin KY. Int J Environ Res Public Health. 2018;15(7):1445. doi:10.3390/ijerph15071445, PMID 29987247, PMC6068473 (cut-offs and labels)
// Source:  Lu C, Chen D, Cai Y, Wei S. J Orthop Surg Res. 2006;1:14. doi:10.1186/1749-799X-1-14, PMID 17150121, PMC1693545 (as MED-ENDO-029)
// Verified: secondary (Subramaniam 2018 PMC6068473; Lu 2006 PMC1693545 "Three risk categories … according to its developer's recommendations")
// Tags:    endocrine, bone, osteoporosis, osta, category
// Units:   index OSTA points; result the source's label
// Range:   any index
// Example: 0 → low-risk; −1 → moderate-risk; −4 → moderate-risk; −5 → high-risk
func OSTARiskCategory(index int) OSTARisk {
	switch {
	case index > -1:
		return OSTALowRisk
	case index >= -4:
		return OSTAModerateRisk
	}
	return OSTAHighRisk
}

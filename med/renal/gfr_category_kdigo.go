package renal

import "fmt"

// GFRCategory is a KDIGO GFR category code, "G1" to "G5" (KDIGO 2024, p. S125).
type GFRCategory string

// The KDIGO 2024 GFR categories (p. S125, "Current chronic kidney disease (CKD) nomenclature used by KDIGO").
const (
	G1  GFRCategory = "G1"
	G2  GFRCategory = "G2"
	G3a GFRCategory = "G3a"
	G3b GFRCategory = "G3b"
	G4  GFRCategory = "G4"
	G5  GFRCategory = "G5"
)

// Term returns the source's own term for the category, verbatim from KDIGO 2024 p. S125 ("Normal or high",
// "Mildly decreased", "Mildly to moderately decreased", "Moderately to severely decreased", "Severely
// decreased", "Kidney failure"), or "" for a value that is not one of the six constants.
func (c GFRCategory) Term() string {
	switch c {
	case G1:
		return "Normal or high"
	case G2:
		return "Mildly decreased"
	case G3a:
		return "Mildly to moderately decreased"
	case G3b:
		return "Moderately to severely decreased"
	case G4:
		return "Severely decreased"
	case G5:
		return "Kidney failure"
	}
	return ""
}

// GFRCategoryKDIGO returns the KDIGO GFR category for a GFR in mL/min/1.73 m²: ≥ 90 G1, ≥ 60 G2, ≥ 45 G3a,
// ≥ 30 G3b, ≥ 15 G4, else G5. The table is printed in integers (90+, 60–89, 45–59, 30–44, 15–29, < 15);
// a value between two printed integers (89.5) takes the lower category (ASSUMPTION A-M3).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for egfrMLMin173 < 0.
//
// ID:      MED-RENAL-011
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Kidney Disease: Improving Global Outcomes (KDIGO) CKD Work Group. KDIGO 2024 Clinical Practice Guideline for the Evaluation and Management of Chronic Kidney Disease. Kidney Int. 2024;105(4S):S117–S314. doi:10.1016/j.kint.2023.10.018. p. S125, "Current chronic kidney disease (CKD) nomenclature used by KDIGO", GFR categories
// Tags:    renal, egfr, kdigo, ckd, category
// Units:   egfrMLMin173 mL/min/1.73 m²; result a category code (Term gives the source's label)
// Range:   egfrMLMin173 ≥ 0
// Example: 89.999 → G2 "Mildly decreased"; 90 → G1 "Normal or high" (KDIGO 2024 p. S125)
func GFRCategoryKDIGO(egfrMLMin173 float64) (GFRCategory, error) {
	if err := checkFinite("egfrMLMin173", egfrMLMin173); err != nil {
		return "", err
	}
	if egfrMLMin173 < 0 {
		return "", fmt.Errorf("%w: egfrMLMin173 = %v, want ≥ 0", ErrOutOfRange, egfrMLMin173)
	}
	switch {
	case egfrMLMin173 >= 90:
		return G1, nil
	case egfrMLMin173 >= 60:
		return G2, nil
	case egfrMLMin173 >= 45:
		return G3a, nil
	case egfrMLMin173 >= 30:
		return G3b, nil
	case egfrMLMin173 >= 15:
		return G4, nil
	}
	return G5, nil
}

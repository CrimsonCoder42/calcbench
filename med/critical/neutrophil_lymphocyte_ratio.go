package critical

import "fmt"

// NeutrophilLymphocyteRatio returns the ratio of the neutrophil count to the lymphocyte count (NLR), both in the
// same unit (for example 10⁹/L), dimensionless.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a negative neutrophil count or a lymphocyte count not
// above zero, or when the ratio is beyond float64.
//
// ID:      MED-CRIT-006
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Zahorec R. Ratio of neutrophil to lymphocyte counts--rapid and simple parameter of systemic inflammation and stress in critically ill. Bratisl Lek Listy. 2001;102(1):5-14. PMID 11723675. Title: "Ratio of neutrophil to lymphocyte counts"
// Tags:    critical, hematology, neutrophil, lymphocyte, nlr, ratio
// Units:   neutrophils and lymphocytes in one count unit (e.g. 10⁹/L); result dimensionless
// Range:   neutrophils ≥ 0; lymphocytes > 0
// Example: neutrophils 8.0, lymphocytes 1.0 → 8.0 (from the definition)
func NeutrophilLymphocyteRatio(neutrophils, lymphocytes float64) (float64, error) {
	for _, e := range [...]error{
		finite("neutrophils", neutrophils), finite("lymphocytes", lymphocytes), atLeastZero("neutrophils", neutrophils),
	} {
		if e != nil {
			return 0, e
		}
	}
	if !(lymphocytes > 0) {
		return 0, fmt.Errorf("%w: lymphocytes = %v, want > 0", ErrOutOfRange, lymphocytes)
	}
	r := neutrophils / lymphocytes
	if err := finite("result", r); err != nil { // e.g. 1e308 / 1e-308
		return 0, fmt.Errorf("%w: neutrophils / lymphocytes is beyond float64", ErrOutOfRange)
	}
	return r, nil
}

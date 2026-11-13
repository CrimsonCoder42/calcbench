package renal

import "fmt"

// UrineAnionGap returns the urine anion gap, Na⁺ + K⁺ − Cl⁻, in mmol/L, from the urine sodium, potassium and
// chloride concentrations. It can be negative.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a negative concentration; ErrDomain when the result is
// beyond float64.
//
// ID:      MED-RENAL-016
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Goldstein MB, Bear R, Richardson RM, Marsden PA, Halperin ML. The urine anion gap: a clinically useful index of ammonium excretion. Am J Med Sci. 1986;292(4):198-202. doi:10.1097/00000441-198610000-00003. PMID 3752165. Abstract: "the urine anion gap (Na+ + K+ - Cl-)"
// Tags:    renal, acid-base, urine anion gap, electrolytes, urine
// Units:   urineSodiumMmolL, urinePotassiumMmolL, urineChlorideMmolL mmol/L (= mEq/L for monovalent ions); result mmol/L
// Range:   every concentration ≥ 0
// Example: Na 50, K 30, Cl 100 → −20 mmol/L (from the definition)
func UrineAnionGap(urineSodiumMmolL, urinePotassiumMmolL, urineChlorideMmolL float64) (float64, error) {
	for _, e := range [...]error{
		checkFinite("urineSodiumMmolL", urineSodiumMmolL), checkFinite("urinePotassiumMmolL", urinePotassiumMmolL),
		checkFinite("urineChlorideMmolL", urineChlorideMmolL),
	} {
		if e != nil {
			return 0, e
		}
	}
	for _, v := range [...]struct {
		name string
		v    float64
	}{{"urineSodiumMmolL", urineSodiumMmolL}, {"urinePotassiumMmolL", urinePotassiumMmolL}, {"urineChlorideMmolL", urineChlorideMmolL}} {
		if v.v < 0 {
			return 0, fmt.Errorf("%w: %s = %v, want ≥ 0", ErrOutOfRange, v.name, v.v)
		}
	}
	g := urineSodiumMmolL + urinePotassiumMmolL - urineChlorideMmolL
	if err := checkResult(g); err != nil {
		return 0, err
	}
	return g, nil
}

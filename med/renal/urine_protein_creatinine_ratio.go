package renal

import "fmt"

// UrineProteinCreatinineRatio returns the protein/creatinine ratio of a single voided urine sample, in mg/mg,
// from the urine protein and urine creatinine concentrations in the same unit (mg/dL here). Ginsberg 1983
// reports the ratio against 24-hour excretion: "a protein/creatinine ratio of more than 3.5 (mg/mg) can be taken
// to represent 'nephrotic-range' proteinuria, and a ratio of less than 0.2 is within normal limits."
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a negative protein or a creatinine not above zero;
// ErrDomain when the ratio is beyond float64.
//
// ID:      MED-RENAL-015
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ginsberg JM, Chang BS, Matarese RA, Garella S. Use of single voided urine samples to estimate quantitative proteinuria. N Engl J Med. 1983;309(25):1543-1546. doi:10.1056/NEJM198312223092503. PMID 6656849. Abstract: "a protein/creatinine ratio of more than 3.5 (mg/mg) can be taken to represent 'nephrotic-range' proteinuria, and a ratio of less than 0.2 is within normal limits"
// Tags:    renal, proteinuria, protein creatinine ratio, upcr, urine
// Units:   urineProteinMgDL, urineCreatinineMgDL mg/dL; result mg/mg
// Range:   urineProteinMgDL ≥ 0; urineCreatinineMgDL > 0
// Example: protein 350 mg/dL, creatinine 100 mg/dL → 3.5 mg/mg (from the definition)
func UrineProteinCreatinineRatio(urineProteinMgDL, urineCreatinineMgDL float64) (float64, error) {
	for _, e := range [...]error{
		checkFinite("urineProteinMgDL", urineProteinMgDL), checkFinite("urineCreatinineMgDL", urineCreatinineMgDL),
		checkPositive("urineCreatinineMgDL", urineCreatinineMgDL),
	} {
		if e != nil {
			return 0, e
		}
	}
	if urineProteinMgDL < 0 {
		return 0, fmt.Errorf("%w: urineProteinMgDL = %v, want ≥ 0", ErrOutOfRange, urineProteinMgDL)
	}
	r := urineProteinMgDL / urineCreatinineMgDL
	if err := checkResult(r); err != nil {
		return 0, err
	}
	return r, nil
}

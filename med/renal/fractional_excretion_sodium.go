package renal

import "github.com/CrimsonCoder42/calcbench/algo/reduce"

// FractionalExcretionSodium returns the fractional excretion of sodium (FENa), in percent:
// (UNa × PCr) / (PNa × UCr) × 100. The two creatinine values may be in any one unit (mg/dL or µmol/L),
// the same for both: the ratio cancels it. The source reports group findings in oliguric acute renal
// failure, not a classification table, so no category is returned.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      MED-RENAL-009
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Espinel CH. The FENa test. Use in the differential diagnosis of acute renal failure. JAMA. 1976;236(6):579-581. doi:10.1001/jama.236.6.579. PMID 947239
// Tags:    renal, fena, sodium, creatinine, urine
// Units:   urineSodiumMmolL, plasmaSodiumMmolL mmol/L; urineCreatinine, plasmaCreatinine any one unit, the same for both; result percent
// Range:   every input > 0 (population: oliguric acute renal failure; no input limits stated by the source)
// Example: UNa 20, PNa 140, UCr 100, PCr 2 → (20 × 2) / (140 × 100) × 100 = 0.2857142857142857 % (hand-derived)
func FractionalExcretionSodium(urineSodiumMmolL, plasmaSodiumMmolL, urineCreatinine, plasmaCreatinine float64) (float64, error) {
	if err := checkFractionalExcretion("urineSodiumMmolL", urineSodiumMmolL, "plasmaSodiumMmolL", plasmaSodiumMmolL,
		urineCreatinine, plasmaCreatinine); err != nil {
		return 0, err
	}
	return fractionalExcretion(urineSodiumMmolL, plasmaSodiumMmolL, urineCreatinine, plasmaCreatinine)
}

// checkFractionalExcretion validates the inputs of MED-RENAL-009 and -010: every finiteness check first,
// then ranges.
func checkFractionalExcretion(urineName string, urine float64, plasmaName string, plasma, urineCreatinine, plasmaCreatinine float64) error {
	if err := checkFinite(urineName, urine); err != nil {
		return err
	}
	if err := checkFinite(plasmaName, plasma); err != nil {
		return err
	}
	if err := checkFinite("urineCreatinine", urineCreatinine); err != nil {
		return err
	}
	if err := checkFinite("plasmaCreatinine", plasmaCreatinine); err != nil {
		return err
	}
	if err := checkPositive(urineName, urine); err != nil {
		return err
	}
	if err := checkPositive(plasmaName, plasma); err != nil {
		return err
	}
	if err := checkPositive("urineCreatinine", urineCreatinine); err != nil {
		return err
	}
	return checkPositive("plasmaCreatinine", plasmaCreatinine)
}

// fractionalExcretion returns (U × PCr) / (P × UCr) × 100 for validated positive inputs. A product that
// overflows or underflows on its own (1e200 × 1e200) is avoided by forming the two ratios first; a result
// beyond float64 returns ErrDomain.
func fractionalExcretion(urine, plasma, urineCreatinine, plasmaCreatinine float64) (float64, error) {
	r := (urine * plasmaCreatinine) / (plasma * urineCreatinine) * 100
	if !reduce.AllFinite(r) || r == 0 {
		r = (urine / plasma) * (plasmaCreatinine / urineCreatinine) * 100
	}
	if err := checkResult(r); err != nil {
		return 0, err
	}
	return r, nil
}

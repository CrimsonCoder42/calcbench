package respiratory

import (
	"fmt"
	"math"
)

// ROXIndex returns the ROX index: (SpO₂ / FiO₂) / respiratory rate, "the ratio of pulse oximetry/fraction of
// inspired oxygen to respiratory rate". With SpO₂ in percent and FiO₂ as a fraction, the source's cut-off is
// "a ROX index greater than or equal to 4.88 measured after 12 hours of HFNC".
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the quotient is beyond
// float64 (a respiratory rate near the smallest float64).
//
// ID:      MED-RESP-012
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Roca O, Messika J, Caralt B, et al. Predicting success of high-flow nasal cannula in pneumonia patients with hypoxemic respiratory failure: The utility of the ROX index. J Crit Care. 2016;35:200-205. doi:10.1016/j.jcrc.2016.05.022. PMID 27481760, abstract
// Tags:    respiratory, rox, hfnc, high-flow, oxygenation, index
// Units:   spo2Percent %; fio2Fraction fraction 0.21–1.0; respiratoryRatePerMin breaths/min; result (%/fraction) per breath/min, unitless as reported
// Range:   0 < spo2Percent ≤ 100; 0.21 ≤ fio2Fraction ≤ 1.0; respiratoryRatePerMin > 0
// Example: SpO₂ 92 %, FiO₂ 0.5, RR 25 → 184 / 25 = 7.36 (hand-derived from the definition)
func ROXIndex(spo2Percent, fio2Fraction, respiratoryRatePerMin float64) (float64, error) {
	if err := checkPositiveAtMost("spo2Percent", spo2Percent, 100); err != nil {
		return 0, err
	}
	if err := checkClosedRange("fio2Fraction", fio2Fraction, 0.21, 1.0); err != nil {
		return 0, err
	}
	if err := checkPositiveFinite("respiratoryRatePerMin", respiratoryRatePerMin); err != nil {
		return 0, err
	}
	rox := spo2Percent / fio2Fraction / respiratoryRatePerMin
	if math.IsInf(rox, 0) {
		return 0, fmt.Errorf("%w: ROX index overflows for respiratoryRatePerMin = %v", ErrDomain, respiratoryRatePerMin)
	}
	return rox, nil
}

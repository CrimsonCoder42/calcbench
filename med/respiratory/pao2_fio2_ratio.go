package respiratory

// checkFiO2Fraction returns ErrNotFinite for NaN or ±Inf and ErrOutOfRange unless 0.21 ≤ fio2Fraction ≤ 1.0
// (room air to pure oxygen).
func checkFiO2Fraction(fio2Fraction float64) error {
	return checkClosedRange("fio2Fraction", fio2Fraction, 0.21, 1.0)
}

// PaO2FiO2Ratio returns the PaO₂/FiO₂ ratio, in mmHg: arterial oxygen tension divided by the fraction of
// inspired oxygen, the oxygenation measure of the Berlin Definition of ARDS.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-RESP-013
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ranieri VM, Rubenfeld GD, Thompson BT, et al. Acute respiratory distress syndrome: the Berlin Definition. JAMA. 2012;307(23):2526-2533. doi:10.1001/jama.2012.5669. PMID 22797452, abstract ("PaO2/FIO2", mm Hg)
// Tags:    respiratory, pf-ratio, pao2-fio2, horowitz, oxygenation, ards
// Units:   paO2MmHg mmHg; fio2Fraction fraction 0.21–1.0; result mmHg
// Range:   paO2MmHg > 0; 0.21 ≤ fio2Fraction ≤ 1.0
// Example: PaO₂ 80 mmHg on FiO₂ 0.4 → 200 mmHg (hand-derived)
func PaO2FiO2Ratio(paO2MmHg, fio2Fraction float64) (float64, error) {
	if err := checkPositiveFinite("paO2MmHg", paO2MmHg); err != nil {
		return 0, err
	}
	if err := checkFiO2Fraction(fio2Fraction); err != nil {
		return 0, err
	}
	return paO2MmHg / fio2Fraction, nil
}

// BerlinARDSOxygenation is the Berlin Definition's oxygenation category for a PaO₂/FiO₂ ratio.
type BerlinARDSOxygenation string

// The Berlin categories, verbatim from the abstract: "mild (200 mm Hg < PaO2/FIO2 ≤ 300 mm Hg), moderate
// (100 mm Hg < PaO2/FIO2 ≤ 200 mm Hg), and severe (PaO2/FIO2 ≤ 100 mm Hg)". BerlinAbove300 is a ratio above
// every category's bound; its value states that bound, not a category of the source.
const (
	BerlinSevere   BerlinARDSOxygenation = "severe"
	BerlinModerate BerlinARDSOxygenation = "moderate"
	BerlinMild     BerlinARDSOxygenation = "mild"
	BerlinAbove300 BerlinARDSOxygenation = "PaO2/FIO2 > 300 mm Hg"
)

// BerlinARDSOxygenationCategory returns the Berlin Definition's oxygenation category for a PaO₂/FiO₂ ratio:
// ≤ 100 severe, (100, 200] moderate, (200, 300] mild, > 300 BerlinAbove300. Only the oxygenation criterion is
// computed; the definition's timing, chest imaging, origin of oedema and PEEP/CPAP criteria are the caller's.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a ratio ≤ 0.
//
// ID:      MED-RESP-014
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ranieri VM, Rubenfeld GD, Thompson BT, et al. Acute respiratory distress syndrome: the Berlin Definition. JAMA. 2012;307(23):2526-2533. doi:10.1001/jama.2012.5669. PMID 22797452, abstract (the three categories and their bounds)
// Tags:    respiratory, ards, berlin, pf-ratio, oxygenation, category
// Units:   pfRatioMmHg mmHg (PaO2FiO2Ratio's result); result the source's category
// Range:   pfRatioMmHg > 0
// Example: 200 → "moderate"; 200.5 → "mild"; 100 → "severe" (the abstract's bounds)
func BerlinARDSOxygenationCategory(pfRatioMmHg float64) (BerlinARDSOxygenation, error) {
	if err := checkPositiveFinite("pfRatioMmHg", pfRatioMmHg); err != nil {
		return "", err
	}
	switch {
	case pfRatioMmHg <= 100:
		return BerlinSevere, nil
	case pfRatioMmHg <= 200:
		return BerlinModerate, nil
	case pfRatioMmHg <= 300:
		return BerlinMild, nil
	default:
		return BerlinAbove300, nil
	}
}

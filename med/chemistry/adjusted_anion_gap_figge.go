package chemistry

// AdjustedAnionGapFigge returns the albumin-adjusted anion gap, in mmol/L, by Figge's correction:
// observed anion gap + 0.25 × (normal albumin − observed albumin), albumin in g/L. "Normal albumin" is the
// source's own term for the reference albumin; its abstract gives no number for it, so the caller passes it.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      MED-CHEM-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Figge J, Jabor A, Kazda A, Fencl V. Anion gap and hypoalbuminemia. Crit Care Med. 1998;26(11):1807-1810. doi:10.1097/00003246-199811000-00019. PMID 9824071. Abstract
// Tags:    chemistry, anion-gap, albumin, acid-base, figge
// Units:   observedAnionGapMmolL mmol/L; albuminGL g/L; normalAlbuminGL g/L; result mmol/L
// Range:   albuminGL > 0; normalAlbuminGL > 0; observedAnionGapMmolL any finite value (limits not stated by the source)
// Example: anion gap 12, albumin 20 g/L, normal albumin 44 g/L → 12 + 0.25 × 24 = 18 mmol/L (hand-derived)
func AdjustedAnionGapFigge(observedAnionGapMmolL, albuminGL, normalAlbuminGL float64) (float64, error) {
	if err := checkFigge(observedAnionGapMmolL, "albuminGL", albuminGL, "normalAlbuminGL", normalAlbuminGL); err != nil {
		return 0, err
	}
	r := observedAnionGapMmolL + 0.25*(normalAlbuminGL-albuminGL)
	if err := checkResult(r); err != nil {
		return 0, err
	}
	return r, nil
}

// AdjustedAnionGapFiggeGDL returns the albumin-adjusted anion gap, in mmol/L, from albumin in g/dL:
// observed anion gap + 2.5 × (normal albumin − observed albumin), the source's factor for g/dL ("if given
// in g/dL, the factor is 2.5"). "Normal albumin" is the source's own term; the caller passes it.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      MED-CHEM-008
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Figge J, Jabor A, Kazda A, Fencl V. Anion gap and hypoalbuminemia. Crit Care Med. 1998;26(11):1807-1810. doi:10.1097/00003246-199811000-00019. PMID 9824071. Abstract
// Tags:    chemistry, anion-gap, albumin, acid-base, figge
// Units:   observedAnionGapMmolL mmol/L; albuminGDL g/dL; normalAlbuminGDL g/dL; result mmol/L
// Range:   albuminGDL > 0; normalAlbuminGDL > 0; observedAnionGapMmolL any finite value (limits not stated by the source)
// Example: anion gap 12, albumin 2.0 g/dL, normal albumin 4.4 g/dL → 12 + 2.5 × 2.4 = 18 mmol/L (hand-derived)
func AdjustedAnionGapFiggeGDL(observedAnionGapMmolL, albuminGDL, normalAlbuminGDL float64) (float64, error) {
	if err := checkFigge(observedAnionGapMmolL, "albuminGDL", albuminGDL, "normalAlbuminGDL", normalAlbuminGDL); err != nil {
		return 0, err
	}
	r := observedAnionGapMmolL + 2.5*(normalAlbuminGDL-albuminGDL)
	if err := checkResult(r); err != nil {
		return 0, err
	}
	return r, nil
}

// checkFigge validates the inputs of MED-CHEM-007 and -008: every finiteness check first, then ranges.
func checkFigge(observedAnionGapMmolL float64, albuminName string, albumin float64, normalAlbuminName string, normalAlbumin float64) error {
	if err := checkFinite("observedAnionGapMmolL", observedAnionGapMmolL); err != nil {
		return err
	}
	if err := checkFinite(albuminName, albumin); err != nil {
		return err
	}
	if err := checkFinite(normalAlbuminName, normalAlbumin); err != nil {
		return err
	}
	if err := checkPositive(albuminName, albumin); err != nil {
		return err
	}
	return checkPositive(normalAlbuminName, normalAlbumin)
}

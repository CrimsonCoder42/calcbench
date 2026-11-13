package cardio

// SubtleAnteriorSTEMI4Variable returns the value of Driver's 4-variable formula,
// 0.052 × QTc-B − 0.151 × QRSV2 − 0.268 × RV4 + 1.062 × STE60V3. The paper's cut-point: "≥18.2".
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when an amplitude is negative or qtcBazettMs ≤ 0.
//
// ID:      MED-CARDIO-013
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Driver BE, Khalil A, Henry T, Kazmi F, Adil A, Smith SW. A new 4-variable formula to differentiate normal variant ST segment elevation in V2-V4 (early repolarization) from subtle left anterior descending coronary occlusion - adding QRS amplitude of V2 improves the model. J Electrocardiol. 2017;50(5):561-569. doi:10.1016/j.jelectrocard.2017.04.005. PMID 28460689. Abstract: "0.052*QTc-B - 0.151*QRSV2 - 0.268*RV4 + 1.062*STE60V3"; cutpoint "≥18.2"; QTc-B the "computerized Bazett-corrected QT interval"
// Tags:    cardio, ecg, stemi, lad-occlusion, early-repolarization, smith
// Units:   qtcBazettMs ms (the ECG machine's Bazett QTc); qrsAmplitudeV2Mm, rWaveV4Mm, stElevation60msV3Mm mm (the units of the 3-variable paper; the 2017 abstract does not restate them); result unitless
// Range:   qtcBazettMs > 0; amplitudes ≥ 0 (ECGs with 1 mm or more ST elevation in V2–V4; the abstract states no input limits)
// Example: QTc 400 ms, QRS V2 15 mm, R V4 10 mm, STE60 V3 2 mm → 20.8 − 2.265 − 2.68 + 2.124 = 17.979 (HAND-DERIVED)
func SubtleAnteriorSTEMI4Variable(qtcBazettMs, qrsAmplitudeV2Mm, rWaveV4Mm, stElevation60msV3Mm float64) (float64, error) {
	if err := checkSTEMIInputs(qtcBazettMs, rWaveV4Mm, stElevation60msV3Mm); err != nil {
		return 0, err
	}
	if err := checkFinite("qrsAmplitudeV2Mm", qrsAmplitudeV2Mm); err != nil {
		return 0, err
	}
	if err := checkRange("qrsAmplitudeV2Mm", qrsAmplitudeV2Mm, 0, maxFinite); err != nil {
		return 0, err
	}
	return 0.052*qtcBazettMs - 0.151*qrsAmplitudeV2Mm - 0.268*rWaveV4Mm + 1.062*stElevation60msV3Mm, nil
}

// SubtleAnteriorSTEMI3Variable returns the value of Smith's 3-variable formula,
// 1.196 × STE60V3 + 0.059 × QTc − 0.326 × RV4. The paper's rule: a value "greater than 23.4 predicted STEMI".
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when an amplitude is negative or qtcBazettMs ≤ 0.
//
// ID:      MED-CARDIO-014
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Smith SW, Khalil A, Henry TD, et al. Electrocardiographic differentiation of early repolarization from subtle anterior ST-segment elevation myocardial infarction. Ann Emerg Med. 2012;60(1):45-56.e2. doi:10.1016/j.annemergmed.2012.02.015. PMID 22520989. Abstract: "([1.196 x ST-segment elevation 60 ms after the J point in lead V3 in mm]+[0.059 x QTc in ms]-[0.326 x R-wave amplitude in lead V4 in mm]) is greater than 23.4"
// Tags:    cardio, ecg, stemi, lad-occlusion, early-repolarization, smith
// Units:   qtcBazettMs ms; rWaveV4Mm, stElevation60msV3Mm mm; result unitless
// Range:   qtcBazettMs > 0; amplitudes ≥ 0 (143 subtle LAD occlusions and 171 early repolarization ECGs; no input limits stated)
// Example: STE60 V3 2 mm, QTc 400 ms, R V4 10 mm → 2.392 + 23.6 − 3.26 = 22.732 (HAND-DERIVED)
func SubtleAnteriorSTEMI3Variable(qtcBazettMs, rWaveV4Mm, stElevation60msV3Mm float64) (float64, error) {
	if err := checkSTEMIInputs(qtcBazettMs, rWaveV4Mm, stElevation60msV3Mm); err != nil {
		return 0, err
	}
	return 1.196*stElevation60msV3Mm + 0.059*qtcBazettMs - 0.326*rWaveV4Mm, nil
}

// maxFinite is the largest float64, the open upper end of a "≥ 0" range.
const maxFinite = 1.7976931348623157e308

// checkSTEMIInputs validates the three inputs the two subtle-anterior-STEMI formulas share.
func checkSTEMIInputs(qtcBazettMs, rWaveV4Mm, stElevation60msV3Mm float64) error {
	for _, p := range [...]struct {
		n string
		v float64
	}{{"qtcBazettMs", qtcBazettMs}, {"rWaveV4Mm", rWaveV4Mm}, {"stElevation60msV3Mm", stElevation60msV3Mm}} {
		if err := checkFinite(p.n, p.v); err != nil {
			return err
		}
	}
	if err := checkPositive("qtcBazettMs", qtcBazettMs); err != nil {
		return err
	}
	if err := checkRange("rWaveV4Mm", rWaveV4Mm, 0, maxFinite); err != nil {
		return err
	}
	return checkRange("stElevation60msV3Mm", stElevation60msV3Mm, 0, maxFinite)
}

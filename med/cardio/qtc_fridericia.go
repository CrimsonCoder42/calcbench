package cardio

import "math"

// QTcFridericia returns the heart-rate-corrected QT interval, in milliseconds, by Fridericia's cube-root
// relation QTc = QT / RR^(1/3), RR in seconds: the QT interval scaled to a pulse period of 1 s.
// Fridericia found the duration of electrical systole proportional to the cube root of the pulse period;
// dividing by RR^(1/3) with RR in seconds is that proportionality normalised to RR = 1 s.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-CARDIO-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Fridericia LS. The duration of systole in an electrocardiogram in normal humans and in patients with heart disease (English translation of Acta Med Scand. 1920;53:469-486). Ann Noninvasive Electrocardiol. 2003;8(4):343-351. doi:10.1046/j.1542-474X.2003.08413.x. PMID 14516292. PMC7328879, "Relationship of the Pulse Frequency and the Duration of the Ventricular Electrocardiogram": "the duration of electrical systole in an electrocardiogram is proportional to the cube root of the duration of the pulse period"
// Source:  Fridericia LS. Die Systolendauer im Elektrokardiogramm bei normalen Menschen und bei Herzkranken. Acta Med Scand. 1920;53(1):469-486. doi:10.1111/j.0954-6820.1920.tb18266.x
// Tags:    cardio, ecg, qt, qtc, fridericia
// Units:   qtMs milliseconds; rrSeconds seconds; ageYears years; result milliseconds
// Range:   rrSeconds 60/135 to 60/51 s (the paper's pulse frequencies, 51 to 135 per minute); ageYears 2–87 (the paper's age range; 50 persons); qtMs > 0. ageYears only enforces the Range
// Example: QT 400 ms at RR 0.8 s → 400 / 0.8^(1/3) = 430.887… ms (hand-derived)
func QTcFridericia(qtMs, rrSeconds, ageYears float64) (float64, error) {
	if err := checkFinite("qtMs", qtMs); err != nil {
		return 0, err
	}
	if err := checkFinite("rrSeconds", rrSeconds); err != nil {
		return 0, err
	}
	if err := checkFinite("ageYears", ageYears); err != nil {
		return 0, err
	}
	if err := checkPositive("qtMs", qtMs); err != nil {
		return 0, err
	}
	if err := checkRange("rrSeconds", rrSeconds, 60.0/135, 60.0/51); err != nil {
		return 0, err
	}
	if err := checkRange("ageYears", ageYears, 2, 87); err != nil {
		return 0, err
	}
	return qtMs / math.Cbrt(rrSeconds), nil
}

package respiratory

// QCSIScore returns the quick COVID-19 Severity Index (qCSI) (0–12) of Haimovich 2020, Table 2: respiratory
// rate ≤ 22 → 0, 23–28 → 1, > 28 → 2 breaths/min; pulse oximetry > 92 → 0, 89–92 → 2, ≤ 88 → 5 %; oxygen flow
// rate ≤ 2 → 0, 3–4 → 4, 5–6 → 5 L/min. The table prints whole numbers; a rate between 22 and 23 or a flow
// between 2 and 3 is compared with the upper edge of the lower band (rate > 22 scores 1, flow > 2 scores 4).
// The paper's pulse oximetry is "the lowest value recorded during the first 4 hours". The paper's risk groups
// are in a figure not transcribed here, so no category is returned.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-RESP-036
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Haimovich AD, Ravindra NG, Stoytchev S, et al. Development and validation of the quick COVID-19 severity index: a prognostic tool for early clinical decompensation. Ann Emerg Med. 2020;76(4):442-453. doi:10.1016/j.annemergmed.2020.07.022. PMID 33012378, PMC7373004, Table 2
// Tags:    respiratory, covid-19, qcsi, oxygen, score
// Units:   respiratoryRatePerMin breaths/min; spo2Percent % (integer); oxygenFlowLPerMin L/min; result points 0–12
// Range:   respiratoryRatePerMin ≥ 0; spo2Percent 0–100; oxygenFlowLPerMin 0–6 (the table's top band); admitted adults ≥ 18 years
// Example: rate 26, SpO₂ 90, flow 4 → 1 + 2 + 4 = 7 (hand-derived from Table 2)
func QCSIScore(respiratoryRatePerMin float64, spo2Percent int, oxygenFlowLPerMin float64) (int, error) {
	if err := checkFiniteNonNegative("respiratoryRatePerMin", respiratoryRatePerMin); err != nil {
		return 0, err
	}
	if err := checkSpO2Percent(spo2Percent); err != nil {
		return 0, err
	}
	if err := checkClosedRange("oxygenFlowLPerMin", oxygenFlowLPerMin, 0, 6); err != nil {
		return 0, err
	}
	s := 0
	switch {
	case respiratoryRatePerMin > 28:
		s += 2
	case respiratoryRatePerMin > 22:
		s++
	}
	switch {
	case spo2Percent <= 88:
		s += 5
	case spo2Percent <= 92:
		s += 2
	}
	switch {
	case oxygenFlowLPerMin > 4:
		s += 5
	case oxygenFlowLPerMin > 2:
		s += 4
	}
	return s, nil
}

package cardio

// DiastolicShockIndex returns the diastolic shock index, heart rate / diastolic arterial pressure, in
// beats/min per mmHg.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when an input is ≤ 0.
//
// ID:      MED-CARDIO-032
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ospina-Tascón GA, Teboul JL, Hernandez G, et al. Diastolic shock index and clinical outcomes in patients with septic shock. Ann Intensive Care. 2020;10(1):41. doi:10.1186/s13613-020-00658-8. PMID 32296976. PMC7160223. Abstract: "Diastolic shock index (DSI) was defined as the ratio between HR and DAP"
// Tags:    cardio, hemodynamics, septic-shock, shock-index
// Units:   heartRateBpm beats/min; diastolicMmHg mmHg; result beats/min per mmHg
// Range:   inputs > 0 (337 + 424 persons in septic shock; no limits stated)
// Example: HR 110, DAP 50 → 2.2 (HAND-DERIVED)
func DiastolicShockIndex(heartRateBpm, diastolicMmHg float64) (float64, error) {
	if err := checkPositivePair("heartRateBpm", heartRateBpm, "diastolicMmHg", diastolicMmHg); err != nil {
		return 0, err
	}
	return heartRateBpm / diastolicMmHg, nil
}

package cardio

// CardiacPowerOutputW returns cardiac power output, in watts, as mean arterial pressure × cardiac output / 451.
// The mean arterial pressure is the measured one, not the one-third estimate.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when an input is ≤ 0.
//
// ID:      MED-CARDIO-031
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Fincke R, Hochman JS, Lowe AM, et al. Cardiac power is the strongest hemodynamic correlate of mortality in cardiogenic shock. J Am Coll Cardiol. 2004;44(2):340-348. doi:10.1016/j.jacc.2004.03.060. PMID 15261929. Abstract: "Cardiac power output (CPO) (W) was calculated as mean arterial pressure x cardiac output/451"
// Tags:    cardio, hemodynamics, cardiogenic-shock, cardiac-power
// Units:   meanArterialMmHg mmHg; cardiacOutputLMin L/min; result W
// Range:   inputs > 0 (541 persons in the SHOCK registry; no limits stated)
// Example: MAP 70 mmHg, CO 4 L/min → 280/451 = 0.6208 W (HAND-DERIVED)
func CardiacPowerOutputW(meanArterialMmHg, cardiacOutputLMin float64) (float64, error) {
	if err := checkPositivePair("meanArterialMmHg", meanArterialMmHg, "cardiacOutputLMin", cardiacOutputLMin); err != nil {
		return 0, err
	}
	return meanArterialMmHg * cardiacOutputLMin / 451, nil
}

// checkPositivePair checks two inputs for finiteness, then positivity.
func checkPositivePair(n1 string, v1 float64, n2 string, v2 float64) error {
	if err := checkFinite(n1, v1); err != nil {
		return err
	}
	if err := checkFinite(n2, v2); err != nil {
		return err
	}
	if err := checkPositive(n1, v1); err != nil {
		return err
	}
	return checkPositive(n2, v2)
}

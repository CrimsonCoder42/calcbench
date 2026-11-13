package respiratory

import "math"

// PeakFlowZone is a peak-flow zone of the NHLBI Asthma Action Plan; its value is the plan's own zone label,
// verbatim ("GREEN ZONE", "YELLOW ZONE", "RED ZONE").
type PeakFlowZone string

// The three zones of the NHLBI Asthma Action Plan (NIH Publication No. 20-HL-5251, February 2021), labels
// verbatim. The plan prints the peak-flow line of each zone as "(80 percent or more of my best peak flow)",
// "(50 to 79 percent of my best peak flow)" and "less than … (50 percent of my best peak flow)".
const (
	PeakFlowGreenZone  PeakFlowZone = "GREEN ZONE"  // 80 percent or more of personal best
	PeakFlowYellowZone PeakFlowZone = "YELLOW ZONE" // 50 to 79 percent of personal best
	PeakFlowRedZone    PeakFlowZone = "RED ZONE"    // less than 50 percent of personal best
)

// Heading returns the plan's heading printed beside the zone, verbatim ("DOING WELL", "ASTHMA IS GETTING
// WORSE", "MEDICAL ALERT!"), or "" for a value that is not one of the three constants.
func (z PeakFlowZone) Heading() string {
	switch z {
	case PeakFlowGreenZone:
		return "DOING WELL"
	case PeakFlowYellowZone:
		return "ASTHMA IS GETTING WORSE"
	case PeakFlowRedZone:
		return "MEDICAL ALERT!"
	}
	return ""
}

// PeakFlowZoneNHLBI returns the NHLBI Asthma Action Plan zone for a peak expiratory flow measured against the
// person's own best peak flow: ≥ 80 percent "GREEN ZONE"; 50 to < 80 percent "YELLOW ZONE"; < 50 percent
// "RED ZONE". The plan prints whole percents (80 or more; 50 to 79; less than 50); a value between 79 and 80
// takes the lower zone, as KDIGO's integer table does in this repository (ASSUMPTION A-M3). The edges are
// compared on the exact quotient pef/best, so 400 of 500 is 80 percent with no rounding error.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for pefLPerMin < 0 or personalBestLPerMin ≤ 0.
//
// ID:      MED-RESP-026
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  National Heart, Lung, and Blood Institute. Asthma Action Plan. NIH Publication No. 20-HL-5251. February 2021. https://www.nhlbi.nih.gov/sites/default/files/publications/Asthma-Action-Plan-2020_rev_508.pdf (read 2026-11-12), zone peak-flow lines
// Tags:    respiratory, asthma, peak-flow, pef, zone, nhlbi, action-plan
// Units:   pefLPerMin, personalBestLPerMin L/min (any one flow unit for both: only their ratio is used); result the plan's zone label
// Range:   pefLPerMin ≥ 0; personalBestLPerMin > 0; a reading above the personal best is ≥ 80 percent
// Example: 400 of 500 → 80 percent → "GREEN ZONE"; 399 of 500 → "YELLOW ZONE"; 250 of 500 → "YELLOW ZONE"; 249 of 500 → "RED ZONE"
func PeakFlowZoneNHLBI(pefLPerMin, personalBestLPerMin float64) (PeakFlowZone, error) {
	if err := checkFiniteNonNegative("pefLPerMin", pefLPerMin); err != nil {
		return "", err
	}
	if err := checkPositiveFinite("personalBestLPerMin", personalBestLPerMin); err != nil {
		return "", err
	}
	// Division is correctly rounded, so an exact 4/5 or 1/2 quotient equals the literals 0.8 and 0.5.
	fraction := pefLPerMin / personalBestLPerMin
	switch {
	case fraction >= 0.8:
		return PeakFlowGreenZone, nil
	case fraction >= 0.5:
		return PeakFlowYellowZone, nil
	}
	return PeakFlowRedZone, nil
}

// checkFiniteNonNegative returns ErrNotFinite for NaN or ±Inf and ErrOutOfRange for v < 0.
func checkFiniteNonNegative(name string, v float64) error {
	return checkClosedRange(name, v, 0, math.MaxFloat64)
}

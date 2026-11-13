package respiratory

import "fmt"

// checkSpO2Percent returns ErrOutOfRange unless 0 ≤ spo2Percent ≤ 100.
func checkSpO2Percent(spo2Percent int) error {
	if spo2Percent < 0 || spo2Percent > 100 {
		return fmt.Errorf("%w: spo2Percent = %d, want 0 ≤ spo2Percent ≤ 100", ErrOutOfRange, spo2Percent)
	}
	return nil
}

// NEWS2SpO2Scale1Points returns the NEWS2 points (0–3) for an oxygen saturation on SpO₂ Scale 1:
// ≤91 → 3, 92–93 → 2, 94–95 → 1, ≥96 → 0. Scale 1 does not depend on air or oxygen; the separate
// "Air or oxygen?" row of the chart (Oxygen 2 points) is not included here.
// Errors: ErrOutOfRange outside 0–100.
//
// ID:      MED-RESP-009
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Royal College of Physicians. National Early Warning Score (NEWS) 2. London: RCP, 2017 (chart footer "© Royal College of Physicians 2017"). Chart 1: The NEWS scoring system, row "SpO2 Scale 1 (%)". https://www.rcp.ac.uk/media/alxev00t/news2-chart-1_the-news-scoring-system_0_0.pdf (read 2026-11-12)
// Tags:    respiratory, news2, spo2, oxygen-saturation, early-warning, points
// Units:   spo2Percent whole percent, as the chart's bands are; result points 0–3
// Range:   spo2Percent 0–100
// Example: 93 → 2; 96 → 0 (Chart 1)
func NEWS2SpO2Scale1Points(spo2Percent int) (int, error) {
	if err := checkSpO2Percent(spo2Percent); err != nil {
		return 0, err
	}
	switch {
	case spo2Percent <= 91:
		return 3, nil
	case spo2Percent <= 93:
		return 2, nil
	case spo2Percent <= 95:
		return 1, nil
	default:
		return 0, nil
	}
}

// NEWS2SpO2Scale2Points returns the NEWS2 points (0–3) for an oxygen saturation on SpO₂ Scale 2:
// ≤83 → 3, 84–85 → 2, 86–87 → 1, 88–92 → 0, ≥93 on air → 0, 93–94 on oxygen → 1, 95–96 on oxygen → 2, ≥97 on oxygen → 3.
// Errors: ErrOutOfRange outside 0–100.
//
// ID:      MED-RESP-010
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Royal College of Physicians. National Early Warning Score (NEWS) 2. London: RCP, 2017 (chart footer "© Royal College of Physicians 2017"). Chart 1: The NEWS scoring system, row "SpO2 Scale 2 (%)" with its "on air" / "on oxygen" cells. https://www.rcp.ac.uk/media/alxev00t/news2-chart-1_the-news-scoring-system_0_0.pdf (read 2026-11-12)
// Tags:    respiratory, news2, spo2, oxygen-saturation, early-warning, points
// Units:   spo2Percent whole percent; onOxygen true when breathing supplemental oxygen; result points 0–3
// Range:   spo2Percent 0–100
// Example: 95 on oxygen → 2; 95 on air → 0 (Chart 1)
func NEWS2SpO2Scale2Points(spo2Percent int, onOxygen bool) (int, error) {
	if err := checkSpO2Percent(spo2Percent); err != nil {
		return 0, err
	}
	switch {
	case spo2Percent <= 83:
		return 3, nil
	case spo2Percent <= 85:
		return 2, nil
	case spo2Percent <= 87:
		return 1, nil
	case spo2Percent <= 92 || !onOxygen:
		return 0, nil
	case spo2Percent <= 94:
		return 1, nil
	case spo2Percent <= 96:
		return 2, nil
	default:
		return 3, nil
	}
}

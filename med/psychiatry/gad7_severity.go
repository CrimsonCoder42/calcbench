package psychiatry

import "fmt"

// GAD7Severity is a GAD-7 severity band; its value is the manual's words. The manual names no band below 5.
type GAD7Severity string

// PHQ and GAD-7 instructions manual, p. 6: "Scores of 5, 10, and 15 represent cutpoints for mild, moderate, and
// severe anxiety, respectively."
const (
	GAD7BelowMild GAD7Severity = ""                 // 0–4: no label in the source
	GAD7Mild      GAD7Severity = "mild anxiety"     // 5–9
	GAD7Moderate  GAD7Severity = "moderate anxiety" // 10–14
	GAD7Severe    GAD7Severity = "severe anxiety"   // 15–21
)

// GAD7SeverityOf returns the GAD-7 severity band for a total (0–21): 5–9 "mild anxiety", 10–14 "moderate
// anxiety", 15–21 "severe anxiety"; below 5 it returns GAD7BelowMild (the empty string), since the source names
// no band there.
// Errors: ErrOutOfRange when total is outside 0–21.
//
// ID:      MED-PSYCH-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Instructions for Patient Health Questionnaire (PHQ) and GAD-7 Measures (developers' manual, phqscreeners.com, instructions.pdf, read 2026-11-12), p. 6 ("Scores of 5, 10, and 15 represent cutpoints for mild, moderate, and severe anxiety, respectively")
// Tags:    psychiatry, anxiety, gad-7, category
// Units:   total points 0–21; result the source's words
// Range:   total 0–21
// Example: 4 → ""; 5 → "mild anxiety"; 10 → "moderate anxiety"; 15 → "severe anxiety" (manual p. 6)
func GAD7SeverityOf(total int) (GAD7Severity, error) {
	if total < 0 || total > 21 {
		return "", fmt.Errorf("%w: total = %d, want 0–21", ErrOutOfRange, total)
	}
	switch {
	case total < 5:
		return GAD7BelowMild, nil
	case total < 10:
		return GAD7Mild, nil
	case total < 15:
		return GAD7Moderate, nil
	}
	return GAD7Severe, nil
}

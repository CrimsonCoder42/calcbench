package endocrine

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// CGMRangePercents is the percent of glucose readings in each of the five ranges of the 2019
// International Consensus on Time in Range (Battelino 2019, Table 2), labelled as the source labels them:
// TBR Level 2, TBR Level 1, TIR, TAR Level 1, TAR Level 2. The five add to 100 up to rounding.
type CGMRangePercents struct {
	BelowRangeLevel2Percent float64 // "TBR: Level 2"
	BelowRangeLevel1Percent float64 // "TBR: Level 1"
	InRangePercent          float64 // "TIR"
	AboveRangeLevel1Percent float64 // "TAR: Level 1"
	AboveRangeLevel2Percent float64 // "TAR: Level 2"
	Readings                int
}

// TimeInRangesMgDL returns the percent of glucose readings (mg/dL) in each consensus range: TBR Level 2
// < 54, TBR Level 1 54–69, TIR 70–180, TAR Level 1 181–250, TAR Level 2 > 250 mg/dL (Battelino 2019,
// Table 2, "% of readings and time"). The table prints whole-number ranges; a fractional reading is
// placed by the printed lower edges 54 and 70 and upper edges 180 and 250 (69.5 is TBR Level 1, 180.5 is
// TAR Level 1). Readings are counted, not weighted by time: equal to the time percent when the readings
// are evenly spaced, as CGM readings are. glucoseMgDL is read, never modified.
// Errors: ErrNotFinite for a NaN or ±Inf reading; ErrOutOfRange when there is no reading or a reading
// is ≤ 0.
//
// ID:      MED-ENDO-011
// Time:    O(n)   n = len(glucoseMgDL)
// Space:   O(1), no allocation; the result is a struct
// Source:  Battelino T, Danne T, Bergenstal RM, et al. Clinical Targets for Continuous Glucose Monitoring Data Interpretation: Recommendations From the International Consensus on Time in Range. Diabetes Care. 2019;42(8):1593-1603. doi:10.2337/dci19-0028. PMID 31177185. PMC6973648, Table 2 "Standardized CGM metrics for clinical care: 2019"
// Tags:    endocrine, diabetes, cgm, glucose, time-in-range, tir, tbr, tar
// Units:   glucoseMgDL mg/dL; results percent of readings (0–100); Readings a count
// Range:   at least one reading; every reading > 0 mg/dL
// Example: 110, 180, 95, 150, 152, 90, 200, 120, 130, 85 mg/dL → TIR 90 %, TAR Level 1 10 %, the others 0 % (HAND-DERIVED)
func TimeInRangesMgDL(glucoseMgDL []float64) (CGMRangePercents, error) {
	// Battelino 2019 Table 2 edges, mg/dL.
	return timeInRanges("glucoseMgDL", glucoseMgDL, 54, 70, 180, 250)
}

// TimeInRangesMmolL returns the percent of glucose readings (mmol/L) in each consensus range: TBR Level 2
// < 3.0, TBR Level 1 3.0–3.8, TIR 3.9–10.0, TAR Level 1 10.1–13.9, TAR Level 2 > 13.9 mmol/L (Battelino
// 2019, Table 2, the source's own mmol/L numbers; no conversion is applied). A reading between the
// printed ranges is placed by the lower edges 3.0 and 3.9 and the upper edges 10.0 and 13.9 (3.85 is TBR
// Level 1, 10.05 is TAR Level 1). Readings are counted, not weighted by time. glucoseMmolL is read, never
// modified.
// Errors: ErrNotFinite for a NaN or ±Inf reading; ErrOutOfRange when there is no reading or a reading
// is ≤ 0.
//
// ID:      MED-ENDO-012
// Time:    O(n)   n = len(glucoseMmolL)
// Space:   O(1), no allocation; the result is a struct
// Source:  Battelino T, Danne T, Bergenstal RM, et al. Clinical Targets for Continuous Glucose Monitoring Data Interpretation: Recommendations From the International Consensus on Time in Range. Diabetes Care. 2019;42(8):1593-1603. doi:10.2337/dci19-0028. PMID 31177185. PMC6973648, Table 2 "Standardized CGM metrics for clinical care: 2019"
// Tags:    endocrine, diabetes, cgm, glucose, time-in-range, tir, tbr, tar
// Units:   glucoseMmolL mmol/L; results percent of readings (0–100); Readings a count
// Range:   at least one reading; every reading > 0 mmol/L
// Example: 2.9, 3.0, 3.9, 10.0, 10.1, 14.0 mmol/L → one reading in each range except two in TIR: 16.67 / 16.67 / 33.33 / 16.67 / 16.67 % (HAND-DERIVED)
func TimeInRangesMmolL(glucoseMmolL []float64) (CGMRangePercents, error) {
	// Battelino 2019 Table 2 edges, mmol/L.
	return timeInRanges("glucoseMmolL", glucoseMmolL, 3.0, 3.9, 10.0, 13.9)
}

// timeInRanges counts readings into g < low2, low2 ≤ g < low1, low1 ≤ g ≤ high1, high1 < g ≤ high2,
// g > high2 and returns the percents.
func timeInRanges(name string, readings []float64, low2, low1, high1, high2 float64) (CGMRangePercents, error) {
	if len(readings) == 0 {
		return CGMRangePercents{}, fmt.Errorf("%w: %s is empty, want ≥ 1 reading", ErrOutOfRange, name)
	}
	var c [5]int
	for i, g := range readings {
		if !reduce.AllFinite(g) {
			return CGMRangePercents{}, fmt.Errorf("%w: %s[%d] = %v", ErrNotFinite, name, i, g)
		}
		switch {
		case !(g > 0):
			return CGMRangePercents{}, fmt.Errorf("%w: %s[%d] = %v, want > 0", ErrOutOfRange, name, i, g)
		case g < low2:
			c[0]++
		case g < low1:
			c[1]++
		case g <= high1:
			c[2]++
		case g <= high2:
			c[3]++
		default:
			c[4]++
		}
	}
	n := float64(len(readings))
	pct := func(k int) float64 { return 100 * float64(c[k]) / n }
	return CGMRangePercents{
		BelowRangeLevel2Percent: pct(0), BelowRangeLevel1Percent: pct(1), InRangePercent: pct(2),
		AboveRangeLevel1Percent: pct(3), AboveRangeLevel2Percent: pct(4), Readings: len(readings),
	}, nil
}

package endocrine

import "fmt"

// CGMMetricPercents holds the continuous glucose monitoring summary a caller has: the percent of readings
// (or time) in each range and the glucose coefficient of variation. The ranges nest: below 54 mg/dL
// (3.0 mmol/L) is part of below 70 mg/dL (3.9 mmol/L), and above 250 mg/dL (13.9 mmol/L) is part of above
// 180 mg/dL (10.0 mmol/L). From MED-ENDO-011/012: BelowRange70Percent = BelowRangeLevel1Percent +
// BelowRangeLevel2Percent, BelowRange54Percent = BelowRangeLevel2Percent, AboveRange180Percent =
// AboveRangeLevel1Percent + AboveRangeLevel2Percent, AboveRange250Percent = AboveRangeLevel2Percent.
type CGMMetricPercents struct {
	InRangePercent       float64 // TIR, 70–180 mg/dL (3.9–10.0 mmol/L), percent
	BelowRange70Percent  float64 // TBR < 70 mg/dL (< 3.9 mmol/L), percent
	BelowRange54Percent  float64 // TBR < 54 mg/dL (< 3.0 mmol/L), percent
	AboveRange180Percent float64 // TAR > 180 mg/dL (> 10.0 mmol/L), percent
	AboveRange250Percent float64 // TAR > 250 mg/dL (> 13.9 mmol/L), percent
	CVPercent            float64 // glucose coefficient of variation (%CV), percent
}

// CGMTargetsMet reports, for each metric, whether it meets its consensus target.
type CGMTargetsMet struct {
	InRange       bool // TIR > 70 %
	BelowRange70  bool // TBR < 70 mg/dL: < 4 %
	BelowRange54  bool // TBR < 54 mg/dL: < 1 %
	AboveRange180 bool // TAR > 180 mg/dL: < 25 %
	AboveRange250 bool // TAR > 250 mg/dL: < 5 %
	CV            bool // %CV ≤ 36 %
}

// CGMTargetsMetBattelino reports which of the 2019 International Consensus on Time in Range targets for
// people with type 1 or type 2 diabetes the metrics meet: TIR ">70%", TBR < 70 mg/dL "<4%", TBR < 54 mg/dL
// "<1%", TAR > 180 mg/dL "<25%", TAR > 250 mg/dL "<5%" (Table 3), and "Glycemic variability (%CV) target
// ≤36%" (Table 2). The inequalities are the source's: exactly 70 % TIR does not meet ">70%", exactly 36 %
// CV meets "≤36%". The targets for older/high-risk individuals and for pregnancy (also in Table 3) are not
// applied. The targets are percents, so the function is the same for mg/dL and mmol/L data.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when a nested range
// exceeds the range it is part of (BelowRange54Percent > BelowRange70Percent, or AboveRange250Percent >
// AboveRange180Percent).
//
// ID:      MED-ENDO-023
// Time:    O(1)
// Space:   O(1), no allocation; the result is a struct
// Source:  Battelino T, Danne T, Bergenstal RM, et al. Clinical Targets for Continuous Glucose Monitoring Data Interpretation: Recommendations From the International Consensus on Time in Range. Diabetes Care. 2019;42(8):1593-1603. doi:10.2337/dci19-0028. PMID 31177185. PMC6973648, Table 3 (type 1 and type 2 targets) and Table 2 ("Glycemic variability (%CV) target ≤36%")
// Tags:    endocrine, diabetes, cgm, glucose, time-in-range, tir, tbr, tar, cv, target
// Units:   every field of m percent; result one bool per target
// Range:   the five range percents 0–100; CVPercent ≥ 0; people with type 1 or type 2 diabetes, not pregnant, not older/high-risk (Table 3's first group)
// Example: TIR 71, TBR70 3.9, TBR54 0.9, TAR180 24.9, TAR250 4.9, CV 36 → all met; TIR 70, TBR70 4, TBR54 1, TAR180 25, TAR250 5, CV 36.1 → none met (Tables 2, 3)
func CGMTargetsMetBattelino(m CGMMetricPercents) (CGMTargetsMet, error) {
	vals := [...]struct {
		name string
		v    float64
	}{
		{"InRangePercent", m.InRangePercent}, {"BelowRange70Percent", m.BelowRange70Percent},
		{"BelowRange54Percent", m.BelowRange54Percent}, {"AboveRange180Percent", m.AboveRange180Percent},
		{"AboveRange250Percent", m.AboveRange250Percent}, {"CVPercent", m.CVPercent},
	}
	for _, x := range vals {
		if err := checkFinite(x.name, x.v); err != nil {
			return CGMTargetsMet{}, err
		}
	}
	for _, x := range vals[:5] {
		if x.v < 0 || x.v > 100 {
			return CGMTargetsMet{}, fmt.Errorf("%w: %s = %v, want 0 ≤ percent ≤ 100", ErrOutOfRange, x.name, x.v)
		}
	}
	if m.CVPercent < 0 {
		return CGMTargetsMet{}, fmt.Errorf("%w: CVPercent = %v, want ≥ 0", ErrOutOfRange, m.CVPercent)
	}
	if m.BelowRange54Percent > m.BelowRange70Percent {
		return CGMTargetsMet{}, fmt.Errorf("%w: BelowRange54Percent = %v > BelowRange70Percent = %v", ErrDomain,
			m.BelowRange54Percent, m.BelowRange70Percent)
	}
	if m.AboveRange250Percent > m.AboveRange180Percent {
		return CGMTargetsMet{}, fmt.Errorf("%w: AboveRange250Percent = %v > AboveRange180Percent = %v", ErrDomain,
			m.AboveRange250Percent, m.AboveRange180Percent)
	}
	// Battelino 2019 Table 3 (type 1 and type 2) and Table 2 (%CV).
	return CGMTargetsMet{
		InRange:       m.InRangePercent > 70,
		BelowRange70:  m.BelowRange70Percent < 4,
		BelowRange54:  m.BelowRange54Percent < 1,
		AboveRange180: m.AboveRange180Percent < 25,
		AboveRange250: m.AboveRange250Percent < 5,
		CV:            m.CVPercent <= 36,
	}, nil
}

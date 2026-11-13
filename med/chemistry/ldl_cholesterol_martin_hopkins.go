package chemistry

import (
	"fmt"
	"math"
)

// martinHopkinsTGUpper holds the printed upper bound, in mg/dL, of each triglyceride stratum (row) of
// Martin et al. 2013, Figure 2: 7–49, 50–56, …, 400–13975.
var martinHopkinsTGUpper = [30]float64{49, 56, 61, 66, 71, 75, 79, 83, 87, 92, 96, 100, 105, 110, 115, 120,
	126, 132, 138, 146, 154, 163, 173, 185, 201, 220, 247, 292, 399, 13975}

// martinHopkinsNonHDLUpper holds the printed upper bound, in mg/dL, of each non-HDL-C stratum (column) of
// Figure 2: < 100, 100–129, 130–159, 160–189, 190–219, ≥ 220.
var martinHopkinsNonHDLUpper = [6]float64{99, 129, 159, 189, 219, math.Inf(1)}

// martinHopkinsFactor is the 180-cell factor table of Martin et al. 2013, Figure 2 (PMC4226221, transcribed
// twice independently 2026-11-12, all cells agree): rows are the TG strata of martinHopkinsTGUpper, columns
// the non-HDL-C strata of martinHopkinsNonHDLUpper.
var martinHopkinsFactor = [30][6]float64{
	{3.5, 3.4, 3.3, 3.3, 3.2, 3.1},   // TG 7–49
	{4.0, 3.9, 3.7, 3.6, 3.6, 3.4},   // TG 50–56
	{4.3, 4.1, 4.0, 3.9, 3.8, 3.6},   // TG 57–61
	{4.5, 4.3, 4.1, 4.0, 3.9, 3.9},   // TG 62–66
	{4.7, 4.4, 4.3, 4.2, 4.1, 3.9},   // TG 67–71
	{4.8, 4.6, 4.4, 4.2, 4.2, 4.1},   // TG 72–75
	{4.9, 4.6, 4.5, 4.3, 4.3, 4.2},   // TG 76–79
	{5.0, 4.8, 4.6, 4.4, 4.3, 4.2},   // TG 80–83
	{5.1, 4.8, 4.6, 4.5, 4.4, 4.3},   // TG 84–87
	{5.2, 4.9, 4.7, 4.6, 4.4, 4.3},   // TG 88–92
	{5.3, 5.0, 4.8, 4.7, 4.5, 4.4},   // TG 93–96
	{5.4, 5.1, 4.8, 4.7, 4.5, 4.3},   // TG 97–100
	{5.5, 5.2, 5.0, 4.7, 4.6, 4.5},   // TG 101–105
	{5.6, 5.3, 5.0, 4.8, 4.6, 4.5},   // TG 106–110
	{5.7, 5.4, 5.1, 4.9, 4.7, 4.5},   // TG 111–115
	{5.8, 5.5, 5.2, 5.0, 4.8, 4.6},   // TG 116–120
	{6.0, 5.5, 5.3, 5.0, 4.8, 4.6},   // TG 121–126
	{6.1, 5.7, 5.3, 5.1, 4.9, 4.7},   // TG 127–132
	{6.2, 5.8, 5.4, 5.2, 5.0, 4.7},   // TG 133–138
	{6.3, 5.9, 5.6, 5.3, 5.0, 4.8},   // TG 139–146
	{6.5, 6.0, 5.7, 5.4, 5.1, 4.8},   // TG 147–154
	{6.7, 6.2, 5.8, 5.4, 5.2, 4.9},   // TG 155–163
	{6.8, 6.3, 5.9, 5.5, 5.3, 5.0},   // TG 164–173
	{7.0, 6.5, 6.0, 5.7, 5.4, 5.1},   // TG 174–185
	{7.3, 6.7, 6.2, 5.8, 5.5, 5.2},   // TG 186–201
	{7.6, 6.9, 6.4, 6.0, 5.6, 5.3},   // TG 202–220
	{8.0, 7.2, 6.6, 6.2, 5.9, 5.4},   // TG 221–247
	{8.5, 7.6, 7.0, 6.5, 6.1, 5.6},   // TG 248–292
	{9.5, 8.3, 7.5, 7.0, 6.5, 5.9},   // TG 293–399
	{11.9, 10.0, 8.8, 8.1, 7.5, 6.7}, // TG 400–13975
}

// Martin et al. 2013, Figure 2 caption: "To convert HDL-C, LDL-C, and total cholesterol, multiply by 0.0259;
// triglycerides, multiply by 0.0113" (mg/dL → mmol/L).
const (
	cholesterolMmolLPerMgDL  = 0.0259
	triglycerideMmolLPerMgDL = 0.0113
)

// LDLCholesterolMartinHopkins returns the estimated LDL cholesterol, in mg/dL, by the Martin–Hopkins method:
// non-HDL-C − TG / factor, with non-HDL-C = total cholesterol − HDL-C and the factor read from the 180-cell
// table of Figure 2 by triglyceride stratum (row) and non-HDL-C stratum (column). A value belongs to the
// first stratum whose printed upper bound it does not exceed: the integer strata are read as half-open at
// the top, so TG 49.5 is in "50–56" and non-HDL-C 99.5 is in "100–129" (ASSUMPTION A-M5). High TG with low
// non-HDL-C can give a result ≤ 0, which is returned as computed.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-CHEM-009
// Time:    O(1) (a scan of at most 30 + 6 fixed strata)
// Space:   O(1), no allocation
// Source:  Martin SS, Blaha MJ, Elshazly MB, et al. Comparison of a novel method vs the Friedewald equation for estimating low-density lipoprotein cholesterol levels from the standard lipid profile. JAMA. 2013;310(19):2061-2068. doi:10.1001/jama.2013.280532. PMID 24240933. Figure 2 (180-cell table; PMC4226221, transcribed twice independently 2026-11-12, all cells agree)
// Tags:    chemistry, lipids, ldl, cholesterol, triglycerides, martin-hopkins
// Units:   totalCholesterolMgDL, hdlCholesterolMgDL, triglyceridesMgDL mg/dL; result mg/dL
// Range:   triglyceridesMgDL 7–13975 (the table's span); totalCholesterolMgDL > 0; hdlCholesterolMgDL > 0; non-HDL-C = total − HDL > 0
// Example: TC 200, HDL 50, TG 150 → non-HDL 150, factor 5.7 (TG 147–154, non-HDL 130–159) → 150 − 150/5.7 = 123.68421052631578 mg/dL (hand-derived)
func LDLCholesterolMartinHopkins(totalCholesterolMgDL, hdlCholesterolMgDL, triglyceridesMgDL float64) (float64, error) {
	if err := checkFinite("totalCholesterolMgDL", totalCholesterolMgDL); err != nil {
		return 0, err
	}
	if err := checkFinite("hdlCholesterolMgDL", hdlCholesterolMgDL); err != nil {
		return 0, err
	}
	if err := checkFinite("triglyceridesMgDL", triglyceridesMgDL); err != nil {
		return 0, err
	}
	if err := checkPositive("totalCholesterolMgDL", totalCholesterolMgDL); err != nil {
		return 0, err
	}
	if err := checkPositive("hdlCholesterolMgDL", hdlCholesterolMgDL); err != nil {
		return 0, err
	}
	if err := checkMartinHopkinsTG("triglyceridesMgDL", triglyceridesMgDL, triglyceridesMgDL); err != nil {
		return 0, err
	}
	nonHDL := totalCholesterolMgDL - hdlCholesterolMgDL
	if !(nonHDL > 0) {
		return 0, fmt.Errorf("%w: totalCholesterolMgDL − hdlCholesterolMgDL = %v, want > 0", ErrOutOfRange, nonHDL)
	}
	return ldlMartinHopkins(nonHDL, triglyceridesMgDL), nil
}

// LDLCholesterolMartinHopkinsSI returns the estimated LDL cholesterol, in mmol/L, by the Martin–Hopkins
// method (MED-CHEM-009) with the paper's own Figure 2 caption factors: cholesterol mg/dL = mmol/L ÷ 0.0259,
// triglycerides mg/dL = mmol/L ÷ 0.0113; the mg/dL result is multiplied by 0.0259 back to mmol/L. The
// strata are applied to the converted mg/dL values.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      MED-CHEM-010
// Time:    O(1) (a scan of at most 30 + 6 fixed strata)
// Space:   O(1), no allocation
// Source:  Martin SS, Blaha MJ, Elshazly MB, et al. Comparison of a novel method vs the Friedewald equation for estimating low-density lipoprotein cholesterol levels from the standard lipid profile. JAMA. 2013;310(19):2061-2068. doi:10.1001/jama.2013.280532. PMID 24240933. Figure 2 (180-cell table; PMC4226221, transcribed twice independently 2026-11-12, all cells agree); Figure 2 caption conversion factors 0.0259 and 0.0113
// Tags:    chemistry, lipids, ldl, cholesterol, triglycerides, martin-hopkins, si
// Units:   totalCholesterolMmolL, hdlCholesterolMmolL, triglyceridesMmolL mmol/L; result mmol/L
// Range:   triglyceridesMmolL ÷ 0.0113 within 7–13975 mg/dL; totalCholesterolMmolL > 0; hdlCholesterolMmolL > 0; total − HDL > 0
// Example: TC 5.18, HDL 1.295, TG 1.695 mmol/L (= 200, 50, 150 mg/dL) → (150 − 150/5.7) × 0.0259 = 3.2034210526315787 mmol/L (hand-derived)
func LDLCholesterolMartinHopkinsSI(totalCholesterolMmolL, hdlCholesterolMmolL, triglyceridesMmolL float64) (float64, error) {
	if err := checkFinite("totalCholesterolMmolL", totalCholesterolMmolL); err != nil {
		return 0, err
	}
	if err := checkFinite("hdlCholesterolMmolL", hdlCholesterolMmolL); err != nil {
		return 0, err
	}
	if err := checkFinite("triglyceridesMmolL", triglyceridesMmolL); err != nil {
		return 0, err
	}
	if err := checkPositive("totalCholesterolMmolL", totalCholesterolMmolL); err != nil {
		return 0, err
	}
	if err := checkPositive("hdlCholesterolMmolL", hdlCholesterolMmolL); err != nil {
		return 0, err
	}
	tgMgDL := triglyceridesMmolL / triglycerideMmolLPerMgDL
	if err := checkMartinHopkinsTG("triglyceridesMmolL", triglyceridesMmolL, tgMgDL); err != nil {
		return 0, err
	}
	if !(totalCholesterolMmolL > hdlCholesterolMmolL) {
		return 0, fmt.Errorf("%w: totalCholesterolMmolL − hdlCholesterolMmolL = %v, want > 0",
			ErrOutOfRange, totalCholesterolMmolL-hdlCholesterolMmolL)
	}
	nonHDL := totalCholesterolMmolL/cholesterolMmolLPerMgDL - hdlCholesterolMmolL/cholesterolMmolLPerMgDL
	r := ldlMartinHopkins(nonHDL, tgMgDL) * cholesterolMmolLPerMgDL
	if err := checkResult(r); err != nil {
		return 0, err
	}
	return r, nil
}

// checkMartinHopkinsTG returns ErrOutOfRange, naming the caller's parameter, unless tgMgDL is in 7–13975.
func checkMartinHopkinsTG(name string, given, tgMgDL float64) error {
	if tgMgDL < 7 || tgMgDL > 13975 {
		return fmt.Errorf("%w: %s = %v (%v mg/dL), want 7 ≤ TG ≤ 13975 mg/dL", ErrOutOfRange, name, given, tgMgDL)
	}
	return nil
}

// ldlMartinHopkins returns non-HDL-C − TG / factor for validated mg/dL inputs (TG in 7–13975, non-HDL-C > 0).
func ldlMartinHopkins(nonHDLMgDL, triglyceridesMgDL float64) float64 {
	row := 0
	for row < len(martinHopkinsTGUpper)-1 && triglyceridesMgDL > martinHopkinsTGUpper[row] {
		row++
	}
	col := 0
	for col < len(martinHopkinsNonHDLUpper)-1 && nonHDLMgDL > martinHopkinsNonHDLUpper[col] {
		col++
	}
	return nonHDLMgDL - triglyceridesMgDL/martinHopkinsFactor[row][col]
}

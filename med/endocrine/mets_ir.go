package endocrine

import (
	"fmt"
	"math"
)

// METSIR returns the metabolic score for insulin resistance (METS-IR), dimensionless:
// ln(2 × fasting glucose + fasting triglycerides) × BMI / ln(HDL cholesterol), every concentration in mg/dL
// and natural logarithms. Not rounded; no band is returned.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      MED-ENDO-027
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Bello-Chavolla OY, Almeda-Valdes P, Gomez-Velasco D, et al. METS-IR, a novel score to evaluate insulin sensitivity, is predictive of visceral adiposity and incident type 2 diabetes. Eur J Endocrinol. 2018;178(5):533-544. doi:10.1530/EJE-17-0883. PMID 29535168. Abstract: "METS-IR as Ln((2*G0)+TG0)*BMI)/(Ln(HDL-c))" (units not in the abstract; full text not opened)
// Source:  Zhang Y, Wang R, Fu X, Song H. Non-insulin-based insulin resistance indexes in predicting severity for coronary artery disease. Diabetol Metab Syndr. 2022. PMID 36528713. PMC9759860, Methods: "METS-IR = Ln [(2 × FPG (mg/dl)) + TG (mg/dl)] × BMI (kg/m²) ÷ Ln [HDL-C (mg/dl)]", cites Bello-Chavolla 2018
// Source:  Wu Z, Cui H, Li W, et al. Comparison of three non-insulin-based insulin resistance indexes in predicting the presence and severity of coronary artery disease. Front Cardiovasc Med. 2022. PMID 35966520. PMC9374164, Methods: the same formula, FPG, TG and HDL-C in mg/dL
// Verified: secondary (Zhang 2022 PMC9759860; Wu 2022 PMC9374164)
// Tags:    endocrine, diabetes, insulin-resistance, mets-ir, glucose, triglycerides, hdl, bmi
// Units:   fastingGlucoseMgDL, triglyceridesMgDL, hdlCholesterolMgDL mg/dL; bmiKgM2 kg/m²; result dimensionless
// Range:   every input > 0 and hdlCholesterolMgDL > 1 mg/dL (ln HDL must be positive; the sources state no further limits)
// Example: glucose 100, TG 150, BMI 28, HDL 50 → ln(350) × 28 / ln(50) = 41.927… (hand-derived)
func METSIR(fastingGlucoseMgDL, triglyceridesMgDL, bmiKgM2, hdlCholesterolMgDL float64) (float64, error) {
	vals := [...]struct {
		name string
		v    float64
	}{
		{"fastingGlucoseMgDL", fastingGlucoseMgDL}, {"triglyceridesMgDL", triglyceridesMgDL},
		{"bmiKgM2", bmiKgM2}, {"hdlCholesterolMgDL", hdlCholesterolMgDL},
	}
	for _, x := range vals {
		if err := checkFinite(x.name, x.v); err != nil {
			return 0, err
		}
	}
	for _, x := range vals {
		if err := checkPositive(x.name, x.v); err != nil {
			return 0, err
		}
	}
	if !(hdlCholesterolMgDL > 1) {
		return 0, fmt.Errorf("%w: hdlCholesterolMgDL = %v, want > 1 (ln HDL in the denominator)", ErrOutOfRange, hdlCholesterolMgDL)
	}
	r := math.Log(2*fastingGlucoseMgDL+triglyceridesMgDL) * bmiKgM2 / math.Log(hdlCholesterolMgDL)
	if err := checkResult(r); err != nil {
		return 0, err
	}
	return r, nil
}

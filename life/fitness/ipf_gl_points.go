package fitness

import (
	"fmt"
	"math"
)

// IPFEvent is an IPF competition type, in the row order of Table 1 of the IPF GL document.
type IPFEvent int

// The four IPF competition types of Table 1 ("The parameters of the equation").
const (
	IPFEquippedPowerlifting IPFEvent = iota // "Equipped Powerlifting"
	IPFClassicPowerlifting                  // "Classic Powerlifting"
	IPFEquippedBenchPress                   // "Equipped Bench Press"
	IPFClassicBenchPress                    // "Classic Bench Press"
)

// ipfGLParams is Table 1 of the IPF GL document, p. 1: A, B, C by [female][event]; female 0 = "Men's",
// 1 = "Women's"; event in the table's row order. Copied digit for digit from the official PDF.
var ipfGLParams = [2][4][3]float64{
	{ // Men's
		{1236.25115, 1449.21864, 0.01644}, // Equipped Powerlifting
		{1199.72839, 1025.18162, 0.00921}, // Classic Powerlifting
		{381.22073, 733.79378, 0.02398},   // Equipped Bench Press
		{320.98041, 281.40258, 0.01008},   // Classic Bench Press
	},
	{ // Women's
		{758.63878, 949.31382, 0.02435},  // Equipped Powerlifting
		{610.32796, 1045.59282, 0.03048}, // Classic Powerlifting
		{221.82209, 357.00377, 0.02937},  // Equipped Bench Press
		{142.40398, 442.52671, 0.04724},  // Classic Bench Press
	},
}

// The document's body-weight condition under equation (2), p. 1: "Bwt ≥ 40kg for men and Bwt ≥ 35kg for women".
const (
	ipfMinBodyWeightMenKg   = 40
	ipfMinBodyWeightWomenKg = 35
)

// ipfGLCoefficientMicro returns the IPF GL coefficient in millionths, rounded half away from zero: the
// document's "Rounding is performed with an accuracy of 6 decimal places" (p. 4), kept as an integer so the
// points can be rounded from an exact product.
func ipfGLCoefficientMicro(event IPFEvent, female bool, bodyWeightKg float64) (float64, error) {
	if err := checkFinite("bodyWeightKg", bodyWeightKg); err != nil {
		return 0, err
	}
	if event < IPFEquippedPowerlifting || event > IPFClassicBenchPress {
		return 0, fmt.Errorf("%w: event = %d, want 0–3 (Table 1 rows)", ErrOutOfRange, int(event))
	}
	sex, minKg := 0, float64(ipfMinBodyWeightMenKg)
	if female {
		sex, minKg = 1, ipfMinBodyWeightWomenKg
	}
	if bodyWeightKg < minKg {
		return 0, fmt.Errorf("%w: bodyWeightKg = %v, want ≥ %v", ErrOutOfRange, bodyWeightKg, minKg)
	}
	p := &ipfGLParams[sex][event]
	// A − B·e^(−C·Bwt) is increasing in Bwt and positive at every row's minimum body weight, so it is
	// positive here; 100 / it × 1e6 = 1e8 / it.
	denom := p[0] - p[1]*math.Exp(-p[2]*bodyWeightKg)
	return math.Round(1e8 / denom), nil
}

// IPFGLCoefficient returns the IPF GL coefficient for a body weight, sex and competition type: equation (2)
// of the IPF GL document, 100 / (A − B·e^(−C·Bwt)), with A, B, C from its Table 1, rounded to 6 decimal
// places (p. 4).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for an event outside the four of Table 1 or a body
// weight below the document's condition (40 kg men, 35 kg women).
//
// ID:      LIF-FIT-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  International Powerlifting Federation. The IPF GL Coefficients for Relative Scoring, actual from May 01, 2020 to Dec 31, 2023 (v.0.6.5, May 2020). https://www.powerlifting.sport/fileadmin/ipf/data/ipf-formula/IPF_GL_Coefficients-2020.pdf (read 2026-11-12), eq. (2) and Table 1 p. 1, rounding p. 4
// Tags:    fitness, powerlifting, ipf, gl, coefficient, relative-scoring
// Units:   bodyWeightKg kg; female the document's Men's/Women's split; result dimensionless, 6 decimals
// Range:   bodyWeightKg ≥ 40 (men), ≥ 35 (women), the document's condition; no upper limit stated
// Example: men, Equipped Powerlifting, 92.04 kg → 0.109039 (the document's practical example 1, p. 4)
func IPFGLCoefficient(event IPFEvent, female bool, bodyWeightKg float64) (float64, error) {
	micro, err := ipfGLCoefficientMicro(event, female, bodyWeightKg)
	if err != nil {
		return 0, err
	}
	return micro / 1e6, nil
}

// IPFGLPoints returns the IPF GL points for a result: equation (1) of the IPF GL document, the coefficient
// (rounded to 6 decimals, LIF-FIT-002) × the result, rounded again to 6 decimal places half up (p. 4); a
// result of 0 gives 0 points. The coefficient is held in millionths so that, for results in the 0.5 kg steps
// the document names, the product before rounding is exact.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for an event outside Table 1, a body weight below the
// document's condition, or resultKg < 0; ErrDomain when the points are beyond float64.
//
// ID:      LIF-FIT-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  International Powerlifting Federation. The IPF GL Coefficients for Relative Scoring, actual from May 01, 2020 to Dec 31, 2023 (v.0.6.5, May 2020). https://www.powerlifting.sport/fileadmin/ipf/data/ipf-formula/IPF_GL_Coefficients-2020.pdf (read 2026-11-12), eq. (1) and (2), Table 1 p. 1, rounding and practical examples p. 4
// Tags:    fitness, powerlifting, ipf, gl, points, relative-scoring
// Units:   bodyWeightKg kg; resultKg kg (total or bench press); result IPF GL points, 6 decimals
// Range:   bodyWeightKg ≥ 40 (men), ≥ 35 (women), the document's condition; resultKg ≥ 0
// Example: women, Classic Bench Press, 70.50 kg, 122.5 kg → 0.790069 × 122.5 = 96.783453 (the document's practical example 2, p. 4)
func IPFGLPoints(event IPFEvent, female bool, bodyWeightKg, resultKg float64) (float64, error) {
	if err := checkFinite("resultKg", resultKg); err != nil {
		return 0, err
	}
	micro, err := ipfGLCoefficientMicro(event, female, bodyWeightKg)
	if err != nil {
		return 0, err
	}
	if resultKg < 0 {
		return 0, fmt.Errorf("%w: resultKg = %v, want ≥ 0", ErrOutOfRange, resultKg)
	}
	if resultKg == 0 {
		return 0, nil // equation (1): "IPF GL Points = 0, if Result = 0"
	}
	pts := math.Round(micro*resultKg) / 1e6
	if err := checkResult(pts); err != nil {
		return 0, err
	}
	return pts, nil
}

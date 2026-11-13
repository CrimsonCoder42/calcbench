package endocrine

import (
	"fmt"
	"math"
)

// QUICKI returns the quantitative insulin sensitivity check index, dimensionless:
// 1 / (log10 fasting insulin (µU/mL) + log10 fasting glucose (mg/dL)). Not rounded.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when insulin × glucose ≤ 1
// (the denominator is zero or negative).
//
// ID:      MED-ENDO-026
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Katz A, Nambi SS, Mather K, Baron AD, Follmann DA, Sullivan G, Quon MJ. Quantitative insulin sensitivity check index: a simple, accurate method for assessing insulin sensitivity in humans. J Clin Endocrinol Metab. 2000;85(7):2402-2410. doi:10.1210/jcem.85.7.6661. PMID 10902785. Abstract: "QUICKI = 1/[log(I(0)) + log(G(0))]" (units and log base not in the abstract)
// Source:  Singh B, Saxena A. Surrogate markers of insulin resistance: a review. World J Diabetes. 2010;1(2):36-47. doi:10.4239/wjd.v1.i2.36. PMC3083884, QUICKI section: "QUICKI = 1/[log (Insulin μU/mL) + log (Glucose mg/dL)]", cites Katz 2000
// Source:  Gutch M, Kumar S, Razi SM, Gupta KK, Gupta A. Assessment of insulin sensitivity/resistance. Indian J Endocrinol Metab. 2015;19(1):160-164. PMC4287763: "QUICKI = 1/(logI0 + logG0)", fasting glucose mg/dL, insulin μIU/mL
// Verified: secondary (Singh 2010 PMC3083884; Gutch 2015 PMC4287763)
// Tags:    endocrine, diabetes, insulin-sensitivity, quicki, glucose, insulin
// Units:   fastingInsulinMicroUnitsPerML µU/mL; fastingGlucoseMgDL mg/dL; result dimensionless (common logarithms)
// Range:   every input > 0 and insulin × glucose > 1 (the formula's own domain; the sources state no further limits)
// Example: insulin 10 µU/mL, glucose 100 mg/dL → 1 / (1 + 2) = 0.3333… (hand-derived)
func QUICKI(fastingInsulinMicroUnitsPerML, fastingGlucoseMgDL float64) (float64, error) {
	if err := checkFinite("fastingInsulinMicroUnitsPerML", fastingInsulinMicroUnitsPerML); err != nil {
		return 0, err
	}
	if err := checkFinite("fastingGlucoseMgDL", fastingGlucoseMgDL); err != nil {
		return 0, err
	}
	if err := checkPositive("fastingInsulinMicroUnitsPerML", fastingInsulinMicroUnitsPerML); err != nil {
		return 0, err
	}
	if err := checkPositive("fastingGlucoseMgDL", fastingGlucoseMgDL); err != nil {
		return 0, err
	}
	// Two logs, not log of the product, so a product beyond float64 never overflows.
	d := math.Log10(fastingInsulinMicroUnitsPerML) + math.Log10(fastingGlucoseMgDL)
	if !(d > 0) {
		return 0, fmt.Errorf("%w: insulin × glucose = %v × %v ≤ 1, the denominator is not positive", ErrDomain, fastingInsulinMicroUnitsPerML, fastingGlucoseMgDL)
	}
	return 1 / d, nil
}

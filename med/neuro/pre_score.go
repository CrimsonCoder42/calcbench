package neuro

import (
	"fmt"
	"math"
)

// PREScore returns the Pittsburgh Response to Endovascular therapy score: age (years) + 2 × NIHSS − 10 ×
// ASPECTS. The source reports benefit of recanalization for "PRE scores -24 to +49" and not for "PRE scores
// <-24 or ≥ 50"; applying those ranges is the caller's.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-NEURO-011
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Rangaraju S, Aghaebrahim A, Streib C, et al. Pittsburgh Response to Endovascular therapy (PRE) score: optimizing patient selection for endovascular therapy for large vessel occlusion strokes. J Neurointerv Surg. 2015;7(11):783-788. doi:10.1136/neurintsurg-2014-011351, PMID 25320054, abstract ("PRE score=age (years)+2×NIHSS-10 × ASPECTS")
// Source:  Barber PA, Demchuk AM, Zhang J, Buchan AM. Validity and reliability of a quantitative computed tomography score in predicting outcome of hyperacute stroke before thrombolytic therapy. Lancet. 2000;355(9216):1670-1674. doi:10.1016/S0140-6736(00)02237-6, PMID 10905241, abstract (ASPECTS "divides the middle-cerebral-artery territory into ten regions of interest": the 0–10 bound)
// Tags:    neuro, stroke, endovascular, aspects, nihss, pre, score
// Units:   ageYears years; nihss NIHSS points; aspects ASPECTS points 0–10; result points (may be negative)
// Range:   ageYears ≥ 0; nihss ≥ 0; aspects 0–10
// Example: age 70, NIHSS 20, ASPECTS 8 → 70 + 40 − 80 = 30 (hand-derived from the abstract's formula)
func PREScore(ageYears float64, nihss, aspects int) (float64, error) {
	if err := checkRange("ageYears", ageYears, 0, math.MaxFloat64/2); err != nil {
		return 0, err
	}
	if nihss < 0 {
		return 0, fmt.Errorf("%w: nihss = %d, want ≥ 0", ErrOutOfRange, nihss)
	}
	if err := checkInt("aspects", aspects, 0, 10); err != nil {
		return 0, err
	}
	return ageYears + 2*float64(nihss) - 10*float64(aspects), nil
}

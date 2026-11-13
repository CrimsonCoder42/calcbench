package stats

import (
	"fmt"
	"math"
	"math/bits"
)

// MatthewsCorrelation returns the Matthews correlation coefficient of a 2×2 table of predictions against
// observations, MCC = (TP·TN − FP·FN)/√((TP + FP)(TP + FN)(TN + FP)(TN + FN)), in [−1, 1]: +1 for perfect
// agreement, 0 for none beyond chance, −1 for total disagreement. The numerator is computed exactly in
// 128-bit integers, so it never overflows or cancels. A zero row or column total leaves the coefficient
// undefined, and is an error, never a silent 0.
// Errors: ErrNegative when a count is negative; ErrZeroTotal when TP + FP, TP + FN, TN + FP or TN + FN is 0.
//
// ID:      NUM-STATS-010
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Matthews, "Comparison of the predicted and observed secondary structure of T4 phage lysozyme", Biochimica et Biophysica Acta (BBA) – Protein Structure, 1975, 405(2):442–451, doi:10.1016/0005-2795(75)90109-9
// Source:  Baldi, Brunak, Chauvin, Andersen, Nielsen, "Assessing the accuracy of prediction algorithms for classification: an overview", Bioinformatics, 2000, 16(5):412–424, doi:10.1093/bioinformatics/16.5.412 (eq. 11; undefined when a margin is 0, p. 415)
// Source:  Chicco, Jurman, "The advantages of the Matthews correlation coefficient (MCC) over F1 score and accuracy in binary classification evaluation", BMC Genomics, 2020, 21:6, doi:10.1186/s12864-019-6413-7 (eq. 2)
// Tags:    matthews-correlation, mcc, phi-coefficient, confusion-matrix, classification, association
// Units:   counts dimensionless; result dimensionless in [−1, 1]
// Example: TP = 90, TN = 80, FP = 20, FN = 10 → (7200 − 200)/√(110·100·100·90) = 0.70353 (hand-checked)
func MatthewsCorrelation(truePositives, trueNegatives, falsePositives, falseNegatives int) (float64, error) {
	for _, c := range [...]struct {
		name  string
		count int
	}{
		{"truePositives", truePositives}, {"trueNegatives", trueNegatives},
		{"falsePositives", falsePositives}, {"falseNegatives", falseNegatives},
	} {
		if c.count < 0 {
			return 0, fmt.Errorf("%w: %s = %d", ErrNegative, c.name, c.count)
		}
	}
	tp, tn := float64(truePositives), float64(trueNegatives)
	fp, fn := float64(falsePositives), float64(falseNegatives)
	// Margins as float64: no int overflow; exact below 2^53, within 1 ulp above.
	for _, m := range [...]struct {
		name string
		sum  float64
	}{
		{"truePositives + falsePositives", tp + fp}, {"truePositives + falseNegatives", tp + fn},
		{"trueNegatives + falsePositives", tn + fp}, {"trueNegatives + falseNegatives", tn + fn},
	} {
		if m.sum == 0 {
			return 0, fmt.Errorf("%w: %s = 0", ErrZeroTotal, m.name)
		}
	}
	num := exactProductDifference(uint64(truePositives), uint64(trueNegatives), uint64(falsePositives), uint64(falseNegatives))
	// Each square root separately: the product of the four margins could exceed the float64 range.
	den := math.Sqrt(tp+fp) * math.Sqrt(tp+fn) * math.Sqrt(tn+fp) * math.Sqrt(tn+fn)
	return min(max(num/den, -1), 1), nil
}

// exactProductDifference returns a·b − c·d, computed exactly in 128 bits and rounded once or twice to
// float64.
func exactProductDifference(a, b, c, d uint64) float64 {
	hi1, lo1 := bits.Mul64(a, b)
	hi2, lo2 := bits.Mul64(c, d)
	sign := 1.0
	if hi1 < hi2 || (hi1 == hi2 && lo1 < lo2) {
		hi1, lo1, hi2, lo2 = hi2, lo2, hi1, lo1
		sign = -1
	}
	lo, borrow := bits.Sub64(lo1, lo2, 0)
	hi, _ := bits.Sub64(hi1, hi2, borrow)
	return sign * (math.Ldexp(float64(hi), 64) + float64(lo))
}

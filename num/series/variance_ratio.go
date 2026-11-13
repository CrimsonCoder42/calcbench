package series

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// VarianceRatio returns the Lo–MacKinlay variance ratio VR(q) of a log-price series and its z statistic
// under the homoskedastic random-walk null. With T = len(logPrices) − 1 one-period returns:
//
//	μ̂   = (p_T − p_0)/T
//	σ̂²_a = (1/(T − 1))·Σ_{k=1}^{T} (p_k − p_{k−1} − μ̂)²
//	σ̄²_c = (1/m)·Σ_{k=q}^{T} (p_k − p_{k−q} − q·μ̂)²,   m = q(T − q + 1)(1 − q/T)
//	VR   = σ̄²_c/σ̂²_a,   z = √T·(VR − 1)·(2(2q − 1)(q − 1)/(3q))^(−1/2)
//
// the bias-corrected estimators with overlapping q-period returns. The paper writes T = nq; the formulas
// are used as written for any T > q. VR ≈ 1 under a random walk, < 1 for mean reversion, > 1 for momentum;
// z is asymptotically N(0, 1) under the null with i.i.d. Gaussian increments. Sums are compensated.
// logPrices is read, never modified.
// Errors: ErrNotFinite when a value is NaN or ±Inf; ErrOutOfRange when q < 2; ErrTooFew when T ≤ q
// (len(logPrices) ≤ q + 1); ErrZeroVariance when every one-period return equals μ̂ (σ̂²_a = 0); ErrOverflow
// when a difference or a sum exceeds the float64 range.
//
// ID:      NUM-SERIES-016
// Time:    O(T)   T = len(logPrices) − 1; three passes
// Space:   O(1), no allocation
// Source:  Lo & MacKinlay, "Stock Market Prices Do Not Follow Random Walks: Evidence from a Simple Specification Test", The Review of Financial Studies, 1988, 1(1):41–66, doi:10.1093/rfs/1.1.41, §1.1–1.2 (the unbiased estimators σ̂²_a and σ̄²_c, M̄_r(q) = σ̄²_c/σ̂²_a − 1, and the homoskedastic z(q) of Theorem 1)
// Tags:    time-series, variance-ratio, lo-mackinlay, random-walk, mean-reversion, momentum, hypothesis-test, finance
// Units:   logPrices natural logarithms of prices (dimensionless); q periods (an integer ≥ 2); vr and z dimensionless
// Example: logPrices 0, 1, 3, 2, 4, q = 2 → μ̂ = 1, σ̂²_a = 6/3 = 2, m = 3, σ̄²_c = 3/3 = 1, VR = 0.5, z = 2·(−0.5)/1 = −1 (hand-checked)
func VarianceRatio(logPrices []float64, q int) (vr, z float64, err error) {
	for i, p := range logPrices {
		if !reduce.AllFinite(p) {
			return 0, 0, fmt.Errorf("%w: logPrices[%d] = %v", ErrNotFinite, i, p)
		}
	}
	if q < 2 {
		return 0, 0, fmt.Errorf("%w: q = %d, want ≥ 2", ErrOutOfRange, q)
	}
	t := len(logPrices) - 1
	if t <= q {
		return 0, 0, fmt.Errorf("%w: %d returns, want more than q = %d", ErrTooFew, t, q)
	}
	tf, qf := float64(t), float64(q)
	mu := (logPrices[t] - logPrices[0]) / tf
	var sa, ca, sc, cc float64
	for k := 1; k <= t; k++ {
		d := logPrices[k] - logPrices[k-1] - mu
		sa, ca = reduce.NeumaierAdd(sa, ca, d*d)
	}
	for k := q; k <= t; k++ {
		d := logPrices[k] - logPrices[k-q] - qf*mu
		sc, cc = reduce.NeumaierAdd(sc, cc, d*d)
	}
	varA := (sa + ca) / (tf - 1)
	m := qf * (tf - qf + 1) * (1 - qf/tf)
	varC := (sc + cc) / m
	if !reduce.AllFinite(mu, varA, varC) {
		return 0, 0, fmt.Errorf("%w: return variance", ErrOverflow)
	}
	if varA == 0 {
		return 0, 0, fmt.Errorf("%w: every one-period return is equal", ErrZeroVariance)
	}
	vr = varC / varA
	z = math.Sqrt(tf) * (vr - 1) / math.Sqrt(2*(2*qf-1)*(qf-1)/(3*qf))
	if !reduce.AllFinite(vr, z) {
		return 0, 0, fmt.Errorf("%w: variance ratio", ErrOverflow)
	}
	return vr, z, nil
}

package invest

import "math"

// CumulativeReturn returns the compound return of a series of simple returns, Π(1 + Rₜ) − 1, computed
// as expm1(Σ ln(1 + Rₜ)) with a compensated sum. A return of −1 (total loss) gives −1.
// returns is read, never modified.
// Errors: ErrTooFew when returns is empty; ErrNotFinite for NaN or ±Inf; ErrRate for a return below −1;
// ErrOverflow when the compound return exceeds the float64 range.
//
// ID:      FIN-INVEST-034
// Time:    O(n)   n = len(returns)
// Space:   O(1), no allocation
// Source:  Campbell, Lo, MacKinlay, The Econometrics of Financial Markets, Princeton University Press, 1997, §1.4.1 "Definitions and Notation", eq. (1.4.2) (multiperiod gross return as the product of single-period gross returns)
// Tags:    invest, returns, cumulative-return, compound-return, multiperiod, time-series, no-alloc
// Units:   returns decimal per period (0.01 = 1 %); result decimal over the whole series
// Example: hand-checked: 0.10, −0.10 → 1.1 × 0.9 − 1 = −0.01
func CumulativeReturn(returns []float64) (float64, error) {
	if err := checkReturnSeries("returns", returns, 1); err != nil {
		return 0, err
	}
	s, total := sumLog1p(returns)
	if total {
		return -1, nil
	}
	return finiteResult(math.Expm1(s))
}

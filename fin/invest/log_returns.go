package invest

// LogReturns writes the continuously compounded returns rₜ = ln(Pₜ/Pₜ₋₁) of a price series into dst
// and returns it, len(prices) − 1 values. When Pₜ/Pₜ₋₁ is beyond float64 it uses ln Pₜ − ln Pₜ₋₁, so
// every result is finite. dst is reused when cap(dst) ≥ len(prices) − 1, else one new slice is made;
// dst must not overlap prices.
// Errors: ErrTooFew when there are fewer than 2 prices; ErrNotFinite for NaN or ±Inf; ErrDomain for a
// price ≤ 0.
//
// ID:      FIN-INVEST-033
// Time:    O(n)   n = len(prices)
// Space:   O(1) extra; no allocation when cap(dst) ≥ n − 1, else one allocation for the result
// Source:  Campbell, Lo, MacKinlay, The Econometrics of Financial Markets, Princeton University Press, 1997, §1.4.1 "Definitions and Notation", eq. (1.4.4) (continuously compounded return rₜ = log(1 + Rₜ) = pₜ − pₜ₋₁)
// Tags:    invest, returns, log-return, continuously-compounded, price-series, time-series, dst
// Units:   prices in any one currency unit per share (or index points); results decimal per period, continuously compounded
// Example: hand-checked: prices 100, 110 → ln 1.1 = 0.0953101798…
func LogReturns(dst, prices []float64) ([]float64, error) {
	if err := checkPriceSeries("prices", prices, 2); err != nil {
		return dst, err
	}
	out := sizedBuffer(dst, len(prices)-1)
	for i := 1; i < len(prices); i++ {
		out[i-1] = logRatio(prices[i], prices[i-1])
	}
	return out, nil
}

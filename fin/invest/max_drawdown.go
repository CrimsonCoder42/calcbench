package invest

// MaxDrawdown returns the largest peak-to-trough decline of a wealth (equity) series as a fraction of the
// peak, max over t of (max₍ₛ≤ₜ₎ Wₛ − Wₜ)/max₍ₛ≤ₜ₎ Wₛ, with the index of that peak and of that trough. The
// first peak and trough of the largest decline are reported. A series that never declines returns 0 with
// both indices 0. Build wealth from simple returns with algo/reduce CumProd of 1 + Rₜ. wealth is read,
// never modified.
// Errors: ErrTooFew when wealth is empty; ErrNotFinite for NaN or ±Inf; ErrDomain for a value ≤ 0.
//
// ID:      FIN-INVEST-047
// Time:    O(n)   n = len(wealth); one pass
// Space:   O(1), no allocation
// Source:  Magdon-Ismail, Atiya, Pratap, Abu-Mostafa, "On the Maximum Drawdown of a Brownian Motion", Journal of Applied Probability, 2004, 41(1):147–161, doi:10.1239/jap/1077134674, §1 (MDD = the largest drop from a peak to a later bottom; here relative to the peak)
// Tags:    invest, risk, drawdown, max-drawdown, mdd, equity-curve, no-alloc
// Units:   wealth in any one currency unit (or an index); result decimal (0.25 = a 25 % decline); indices into wealth
// Example: hand-checked: 100, 120, 90, 110, 80, 130, 125 → (120 − 80)/120 = 1/3, peak 1, trough 4
func MaxDrawdown(wealth []float64) (drawdown float64, peakIndex, troughIndex int, err error) {
	if err := checkPriceSeries("wealth", wealth, 1); err != nil {
		return 0, 0, 0, err
	}
	runPeak := 0
	for i, w := range wealth {
		if w > wealth[runPeak] {
			runPeak = i
			continue
		}
		if dd := (wealth[runPeak] - w) / wealth[runPeak]; dd > drawdown {
			drawdown, peakIndex, troughIndex = dd, runPeak, i
		}
	}
	return drawdown, peakIndex, troughIndex, nil
}

package market

import (
	"fmt"
	"math"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// checkImpact validates the Almgren–Chriss market parameters shared by the trajectory and its cost and
// returns η̃ = η − ½γτ, which must be > 0.
func checkImpact(tau, sigma, eta, gamma float64) (etaTilde float64, err error) {
	if !reduce.AllFinite(tau, sigma, eta, gamma) {
		return 0, fmt.Errorf("%w: tau, sigma, eta, gamma = %v, %v, %v, %v", ErrNotFinite, tau, sigma, eta, gamma)
	}
	if tau <= 0 {
		return 0, fmt.Errorf("%w: period length = %v, want > 0", ErrOutOfRange, tau)
	}
	if sigma < 0 || eta <= 0 || gamma < 0 {
		return 0, fmt.Errorf("%w: want sigma ≥ 0, eta > 0, gamma ≥ 0; got %v, %v, %v", ErrOutOfRange, sigma, eta, gamma)
	}
	etaTilde = eta - gamma*tau/2
	if etaTilde <= 0 {
		return 0, fmt.Errorf("%w: η̃ = η − γτ/2 = %v, want > 0", ErrDomain, etaTilde)
	}
	return etaTilde, nil
}

// AlmgrenChrissTrajectory returns the Almgren–Chriss optimal liquidation schedule under linear
// temporary and permanent impact: the holdings x_j = sinh(κ(T − t_j))/sinh(κT)·X at t_j = jτ, τ = T/N,
// j = 0 … N (x_0 = X, x_N = 0), where κ solves cosh(κτ) = 1 + ½κ̃²τ², κ̃² = λσ²/η̃ and η̃ = η − ½γτ.
// λ = 0 (risk-neutral) gives the straight line x_j = (1 − j/N)·X. The sinh ratio is evaluated as
// e^{−κt_j}·(1 − e^{−2κ(T−t_j)})/(1 − e^{−2κT}), so a large κT never overflows. Holdings are written into
// dst, reused when cap(dst) ≥ N + 1, else one new slice is allocated.
// Errors: ErrNotFinite if a float input is NaN or ±Inf; ErrOutOfRange if periods < 1, horizon ≤ 0,
// sigma < 0, eta ≤ 0, gamma < 0 or riskAversion < 0; ErrDomain if η̃ ≤ 0.
//
// ID:      FIN-MKT-016
// Time:    O(N)   N = periods
// Space:   O(1) extra; O(N) for the result when dst is too small
// Source:  Almgren, Chriss, "Optimal execution of portfolio transactions", Journal of Risk 3(2):5–39, 2001, doi:10.21314/JOR.2001.041, §2 "The efficient frontier of optimal execution", §2.2 (x_j = sinh(κ(T − t_j))/sinh(κT)·X; κ̃² = λσ²/η̃; (2/τ²)(cosh(κτ) − 1) = κ̃²) and §1 (η̃ = η − ½γτ)
// Tags:    market, execution, almgren-chriss, optimal-execution, trading-trajectory, market-impact, liquidation
// Units:   shares in shares; horizon in time units; sigma in price units per share per √time unit; eta in price units·time per share²; gamma in price units per share²; riskAversion per price unit (currency⁻¹); result in shares
// Example: hand-checked: X = 100, T = 2, N = 2, σ = 1, η = 1, γ = 0, λ = 0.5 → κ̃² = 0.5, cosh κ = 1.25, κ = ln 2, x = 100, 40, 0
func AlmgrenChrissTrajectory(dst []float64, shares, horizon float64, periods int, sigma, eta, gamma, riskAversion float64) ([]float64, error) {
	if periods < 1 {
		return dst[:0], fmt.Errorf("%w: periods = %d, want ≥ 1", ErrOutOfRange, periods)
	}
	if !reduce.AllFinite(shares, horizon, riskAversion) {
		return dst[:0], fmt.Errorf("%w: shares, horizon, riskAversion = %v, %v, %v", ErrNotFinite, shares, horizon, riskAversion)
	}
	if horizon <= 0 {
		return dst[:0], fmt.Errorf("%w: horizon = %v, want > 0", ErrOutOfRange, horizon)
	}
	if riskAversion < 0 {
		return dst[:0], fmt.Errorf("%w: riskAversion = %v, want ≥ 0", ErrOutOfRange, riskAversion)
	}
	tau := horizon / float64(periods)
	etaTilde, err := checkImpact(tau, sigma, eta, gamma)
	if err != nil {
		return dst[:0], err
	}
	if cap(dst) >= periods+1 {
		dst = dst[:periods+1]
	} else {
		dst = make([]float64, periods+1)
	}
	// u = ½κ̃²τ²; κτ = acosh(1 + u) = log1p(u + √(u(u + 2))), accurate for small u.
	u := riskAversion * sigma * sigma / etaTilde * tau * tau / 2
	kt := math.Log1p(u + math.Sqrt(u*(u+2)))
	if !reduce.AllFinite(kt) {
		return dst[:0], fmt.Errorf("%w: κτ", ErrOverflow)
	}
	n := float64(periods)
	for j := range periods + 1 {
		if kt == 0 {
			dst[j] = shares * (1 - float64(j)/n)
			continue
		}
		// e^{κ(T−t_j) − κT} = e^{−κt_j}, taken directly so no a − b cancellation (or FMA) enters.
		a := kt * (n - float64(j)) // κ(T − t_j)
		b := kt * n                // κT
		dst[j] = shares * math.Exp(-kt*float64(j)) * (math.Expm1(-2*a) / math.Expm1(-2*b))
	}
	dst[0] = shares
	dst[periods] = 0
	return dst, nil
}

// AlmgrenChrissExpectedCost returns the expected cost E[x] and its variance V[x] of a liquidation
// trajectory x_0 = X, x_1, …, x_N = 0 under the Almgren–Chriss linear impact model, with trades
// n_k = x_{k−1} − x_k:
//
//	E[x] = ½γX² + ε·Σ|n_k| + (η̃/τ)·Σ n_k²,   η̃ = η − ½γτ
//	V[x] = σ²·τ·Σ_{k=1}^{N} x_k²
//
// Any schedule may be passed (AlmgrenChrissTrajectory's, or a TWAP line). holdings is read, never
// modified.
// Errors: ErrTooFew if len(holdings) < 2; ErrNotFinite if a value is NaN or ±Inf; ErrOutOfRange if
// tau ≤ 0, sigma < 0, eta ≤ 0, gamma < 0 or epsilon < 0; ErrDomain if the last holding is not 0 or
// η̃ ≤ 0; ErrOverflow if a result exceeds float64.
//
// ID:      FIN-MKT-017
// Time:    O(N)   N = len(holdings) − 1 periods
// Space:   O(1), no allocation
// Source:  Almgren, Chriss, "Optimal execution of portfolio transactions", Journal of Risk 3(2):5–39, 2001, doi:10.21314/JOR.2001.041, §1 "The trading model", linear impact (E(x) = ½γX² + εΣ|n_k| + (η̃/τ)Σn_k², V(x) = σ²Σ τx_k²)
// Tags:    market, execution, almgren-chriss, expected-cost, variance, market-impact, transaction-cost
// Units:   holdings in shares; tau in time units; sigma in price units per share per √time unit; eta in price units·time per share²; gamma in price units per share²; epsilon in price units per share; E in currency, V in currency²
// Example: hand-checked: x = 100, 40, 0, τ = 1, σ = 1, η = 1, γ = 0.1, ε = 0.5 → E = 500 + 50 + 0.95·5200 = 5490, V = 1600
func AlmgrenChrissExpectedCost(holdings []float64, tau, sigma, eta, gamma, epsilon float64) (mean, variance float64, err error) {
	if len(holdings) < 2 {
		return 0, 0, fmt.Errorf("%w: len(holdings) = %d, want ≥ 2", ErrTooFew, len(holdings))
	}
	if !reduce.AllFinite(holdings...) {
		return 0, 0, fmt.Errorf("%w: holdings", ErrNotFinite)
	}
	if !reduce.AllFinite(epsilon) {
		return 0, 0, fmt.Errorf("%w: epsilon = %v", ErrNotFinite, epsilon)
	}
	if epsilon < 0 {
		return 0, 0, fmt.Errorf("%w: epsilon = %v, want ≥ 0", ErrOutOfRange, epsilon)
	}
	etaTilde, err := checkImpact(tau, sigma, eta, gamma)
	if err != nil {
		return 0, 0, err
	}
	last := len(holdings) - 1
	if holdings[last] != 0 {
		return 0, 0, fmt.Errorf("%w: holdings[%d] = %v, want 0 (fully liquidated)", ErrDomain, last, holdings[last])
	}
	var abs, absC, sq, sqC, xx, xxC float64
	for k := 1; k <= last; k++ {
		nk := holdings[k-1] - holdings[k]
		abs, absC = reduce.NeumaierAdd(abs, absC, math.Abs(nk))
		sq, sqC = reduce.NeumaierAdd(sq, sqC, nk*nk)
		xx, xxC = reduce.NeumaierAdd(xx, xxC, holdings[k]*holdings[k])
	}
	x0 := holdings[0]
	mean = gamma*x0*x0/2 + epsilon*(abs+absC) + etaTilde/tau*(sq+sqC)
	variance = sigma * sigma * tau * (xx + xxC)
	if !reduce.AllFinite(mean, variance) {
		return 0, 0, fmt.Errorf("%w: E or V", ErrOverflow)
	}
	return mean, variance, nil
}

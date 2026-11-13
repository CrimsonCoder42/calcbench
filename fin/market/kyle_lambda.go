package market

import (
	"errors"
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
	"github.com/CrimsonCoder42/calcbench/num/stats"
)

// KyleLambda returns Kyle's λ, the price impact per unit of signed order flow, estimated as the
// least-squares slope of the price changes on the signed volumes: Δpₜ = c + λ·(signed volume)ₜ + error
// (Kyle's p = μ + λ(x + u), fitted by NUM-STATS-009 LinearRegression). Sign the volume with
// FIN-MKT-006, FIN-MKT-007 or FIN-MKT-008. The slices are read, never modified.
// Errors: ErrEmpty if priceChanges is empty; ErrLengthMismatch if the lengths differ; ErrNotFinite if a
// value is NaN or ±Inf; ErrDomain (also wrapping stats.ErrTooFew or stats.ErrZeroVariance) for one
// interval or constant signed volume; ErrOverflow (also wrapping stats.ErrOverflow) if the slope
// exceeds float64.
//
// ID:      FIN-MKT-013
// Time:    O(n)   n = number of intervals
// Space:   O(1), no allocation
// Source:  Kyle, "Continuous auctions and insider trading", Econometrica 53(6):1315–1335, 1985, doi:10.2307/1913210, §2 Theorem 1 (linear price rule p = μ + λ(x + u))
// Tags:    market, microstructure, kyle-lambda, price-impact, market-depth, regression, liquidity
// Units:   priceChanges in one price unit; signedVolumes in one size unit; result in price units per size unit
// Example: hand-checked: signed volumes −2, 0, 2, Δp −1, 0.5, 2 → λ = 0.75
func KyleLambda(priceChanges, signedVolumes []float64) (float64, error) {
	if len(priceChanges) == 0 {
		return 0, fmt.Errorf("%w: priceChanges", ErrEmpty)
	}
	if len(signedVolumes) != len(priceChanges) {
		return 0, fmt.Errorf("%w: len(priceChanges) = %d, len(signedVolumes) = %d", ErrLengthMismatch, len(priceChanges), len(signedVolumes))
	}
	if !reduce.AllFinite(priceChanges...) {
		return 0, fmt.Errorf("%w: priceChanges", ErrNotFinite)
	}
	if !reduce.AllFinite(signedVolumes...) {
		return 0, fmt.Errorf("%w: signedVolumes", ErrNotFinite)
	}
	_, slope, err := stats.LinearRegression(signedVolumes, priceChanges)
	if err != nil {
		if errors.Is(err, stats.ErrOverflow) {
			return 0, fmt.Errorf("%w: %w", ErrOverflow, err)
		}
		return 0, fmt.Errorf("%w: %w", ErrDomain, err)
	}
	return slope, nil
}

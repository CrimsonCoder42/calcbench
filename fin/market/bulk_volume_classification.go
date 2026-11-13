package market

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
	"github.com/CrimsonCoder42/calcbench/num/prob"
)

// checkBar validates one bar for bulk volume classification: volume ≥ 0 and finite, priceChange
// finite, priceChangeSD finite and > 0.
func checkBar(volume, priceChange, priceChangeSD float64) error {
	if err := checkSize("volume", volume); err != nil {
		return err
	}
	if !reduce.AllFinite(priceChange) {
		return fmt.Errorf("%w: priceChange = %v", ErrNotFinite, priceChange)
	}
	if !reduce.AllFinite(priceChangeSD) {
		return fmt.Errorf("%w: priceChangeSD = %v", ErrNotFinite, priceChangeSD)
	}
	if priceChangeSD <= 0 {
		return fmt.Errorf("%w: priceChangeSD = %v, want > 0", ErrOutOfRange, priceChangeSD)
	}
	return nil
}

// BulkVolumeClassification splits one bar's volume into buy and sell volume by bulk volume
// classification with the standard normal CDF Z: buy = V·Z(ΔP/σ_ΔP), sell = V·Z(−ΔP/σ_ΔP) = V − buy.
// ΔP is the bar's price change and σ_ΔP the standard deviation of bar price changes, estimated by the
// caller over its chosen sample. Summing the bars of a volume bucket gives V^B_τ and V^S_τ for VPIN
// (FIN-MKT-010). For the Student-t variant use BulkVolumeClassificationT (FIN-MKT-009).
// Errors: ErrNotFinite if an input is NaN or ±Inf; ErrNegative if volume < 0; ErrOutOfRange if
// priceChangeSD ≤ 0.
//
// ID:      FIN-MKT-008
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Easley, López de Prado, O'Hara, "Flow toxicity and liquidity in a high-frequency world", Review of Financial Studies 25(5):1457–1493, 2012, doi:10.1093/rfs/hhs053, §2 (bulk volume classification: V^B_τ = Σ V_i·Z((P_i − P_{i−1})/σ_ΔP), V^S_τ = V_τ − V^B_τ)
// Tags:    market, microstructure, bulk-volume-classification, bvc, trade-classification, vpin, order-flow-toxicity
// Units:   volume in one size unit; priceChange and priceChangeSD in one price unit; results in the volume unit
// Example: hand-checked: V = 1000, ΔP = 0 → buy 500, sell 500; ΔP = 1.96, σ = 1 → buy 975.0021…, sell 24.9978…
func BulkVolumeClassification(volume, priceChange, priceChangeSD float64) (buyVolume, sellVolume float64, err error) {
	if err := checkBar(volume, priceChange, priceChangeSD); err != nil {
		return 0, 0, err
	}
	up, err := prob.NormalCDF(priceChange, 0, priceChangeSD)
	if err != nil {
		return 0, 0, fmt.Errorf("%w: %w", ErrDomain, err)
	}
	down, err := prob.NormalCDF(-priceChange, 0, priceChangeSD)
	if err != nil {
		return 0, 0, fmt.Errorf("%w: %w", ErrDomain, err)
	}
	return volume * up, volume * down, nil
}

// BulkVolumeClassificationT splits one bar's volume into buy and sell volume by bulk volume
// classification with the Student-t CDF t_ν in place of the normal: buy = V·t_ν(ΔP/σ_ΔP),
// sell = V·t_ν(−ΔP/σ_ΔP). The degrees of freedom ν is a parameter (the source fits it to the data);
// as ν → ∞ the result tends to BulkVolumeClassification (FIN-MKT-008).
// Errors: ErrNotFinite if an input is NaN or ±Inf; ErrNegative if volume < 0; ErrOutOfRange if
// priceChangeSD ≤ 0 or degreesOfFreedom ≤ 0; ErrDomain (also wrapping the num/prob error) if the t CDF
// does not converge.
//
// ID:      FIN-MKT-009
// Time:    O(√ν + log(1/ε))   ν = degreesOfFreedom, ε = 2⁻⁵² (the t CDF, NUM-PROB-078)
// Space:   O(1), no allocation
// Source:  Easley, López de Prado, O'Hara, "Discerning information from trade data", Journal of Financial Economics 120(2):269–285, 2016, doi:10.1016/j.jfineco.2016.01.018, §2 (bulk volume classification with the Student-t distribution)
// Source:  Easley, López de Prado, O'Hara, Review of Financial Studies 25(5):1457–1493, 2012, §2 (the normal form this generalises)
// Tags:    market, microstructure, bulk-volume-classification, bvc, student-t, trade-classification, vpin
// Units:   volume in one size unit; priceChange and priceChangeSD in one price unit; degreesOfFreedom dimensionless; results in the volume unit
// Example: hand-checked: V = 1000, ΔP = 1, σ = 1, ν = 1 (Cauchy, t₁(1) = 0.75) → buy 750, sell 250
func BulkVolumeClassificationT(volume, priceChange, priceChangeSD, degreesOfFreedom float64) (buyVolume, sellVolume float64, err error) {
	if err := checkBar(volume, priceChange, priceChangeSD); err != nil {
		return 0, 0, err
	}
	if !reduce.AllFinite(degreesOfFreedom) {
		return 0, 0, fmt.Errorf("%w: degreesOfFreedom = %v", ErrNotFinite, degreesOfFreedom)
	}
	if degreesOfFreedom <= 0 {
		return 0, 0, fmt.Errorf("%w: degreesOfFreedom = %v, want > 0", ErrOutOfRange, degreesOfFreedom)
	}
	z := priceChange / priceChangeSD
	if !reduce.AllFinite(z) {
		return 0, 0, fmt.Errorf("%w: priceChange/priceChangeSD", ErrOverflow)
	}
	up, err := prob.StudentTCDF(z, degreesOfFreedom)
	if err != nil {
		return 0, 0, fmt.Errorf("%w: %w", ErrDomain, err)
	}
	down, err := prob.StudentTCDF(-z, degreesOfFreedom)
	if err != nil {
		return 0, 0, fmt.Errorf("%w: %w", ErrDomain, err)
	}
	return volume * up, volume * down, nil
}

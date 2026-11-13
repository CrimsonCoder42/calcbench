package prob

import "slices"

// NormalizedImpliedProbabilities appends to dst πᵢ = (1/dᵢ) / Σⱼ (1/dⱼ) for each outcome of one market
// and returns the extended slice: the implied probabilities with the overround removed proportionally
// ("basic normalisation"), summing to 1 up to rounding. Shin's method is reported to give more accurate
// probabilities (Štrumbelj 2014, as reported by Egidi et al. 2018, §1 of the preprint); it is not provided.
// decimalOdds is read, never modified.
// Errors: as Overround. On error dst is returned with its length unchanged.
//
// ID:      NUM-PROB-017
// Time:    O(n)   n = len(decimalOdds)
// Space:   O(1) extra; O(n) for the result, one allocation only if dst lacks capacity
// Source:  Egidi, Pauli, Torelli, "Combining historical data and bookmakers' odds in modelling football scores", Statistical Modelling, 2018, 18(5–6):436–459, doi:10.1177/1471082X18798414; author preprint arXiv:1802.08848 (read 2026-11-12), §2 "Transforming the betting odds into probabilities", eq. (2.1) (basic normalisation πᵢ = oᵢ/β, oᵢ the inverse odds, β the booksum)
// Source:  Wunderlich, Memmert, "The Betting Odds Rating System: Using soccer forecasts to forecast soccer", PLOS ONE, 2018, 13(6):e0198668, doi:10.1371/journal.pone.0198668 (read 2026-11-12), Method, "Transferring betting odds to probabilities" (basic normalisation: each inverse odd divided by their sum)
// Tags:    probability, odds, implied-probability, overround, normalization, betting
// Units:   decimalOdds: units returned per unit staked, each ≥ 1; result: probabilities in (0, 1]
// Example: [1.5, 2.75] → [11/17, 6/17] ≈ [0.6471, 0.3529]
func NormalizedImpliedProbabilities(dst, decimalOdds []float64) ([]float64, error) {
	if err := checkDecimalOdds(decimalOdds); err != nil {
		return dst, err
	}
	var s compensatedSum
	for _, d := range decimalOdds {
		s.add(1 / d)
	}
	booksum := s.total()
	out := slices.Grow(dst, len(decimalOdds))
	for _, d := range decimalOdds {
		// (1/d)/booksum rather than 1/(d·booksum): the product could overflow for a huge d.
		out = append(out, min(1, (1/d)/booksum))
	}
	return out, nil
}

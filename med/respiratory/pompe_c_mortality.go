package respiratory

import "math"

// POMPECFindings are the predictors of the POMPE-C model (Kline 2012, Table 2), worded as in the table.
type POMPECFindings struct {
	WeightLb               float64 // "Patient weight (lb)"
	HighestRespRatePerMin  float64 // "Highest respiratory rate (breaths/min)" in the previous 6 hours
	SaO2Percent            float64 // "SaO2 (%)", pulse oximetry on room air
	HeartRateAbove99       bool    // "Heart rate >99 beats/min"
	AlteredMentalStatus    bool    // "Altered mental status"
	RespiratoryDistress    bool    // "Respiratory distress"
	DoNotResuscitate       bool    // "Do not resuscitate status"
	UnilateralLimbSwelling bool    // "Unilateral limb swelling"
}

// POMPECMortalityPercent returns the POMPE-C estimate of the probability of death within 30 days, in percent,
// for pulmonary embolism with active cancer, by the paper's printed equation
// P = 100 × (1 − 1/(1 + exp(3.718 + 1.551·DNR + 0.800·RespDist + 0.734·ULS + 1.473·AMS + HR + 0.044·RR −
// 0.063·SaO₂ − 0.012·weight))). The printed equation writes the heart-rate term as "HR" with no coefficient;
// Table 2's coefficient column is shifted by one row against its odds-ratio column, and matching each odds ratio
// to its logarithm (ln 2.8 = 1.03) gives the heart-rate coefficient 1.027 used here (ledger pompe_c).
// It is computed as 100/(1 + exp(−L)), the same value without the cancellation of 1 − 1/(1 + e^L).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-RESP-037
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Kline JA, Roy PM, Than MP, et al. Derivation and validation of a multivariate model to predict mortality from pulmonary embolism with cancer: the POMPE-C tool. Thromb Res. 2012;129(5):e194-e199. doi:10.1016/j.thromres.2012.03.015. PMID 22475313, PMC3350789, Table 2 and the equation for P
// Tags:    respiratory, pulmonary-embolism, cancer, pompe-c, mortality, logistic
// Units:   WeightLb lb; HighestRespRatePerMin breaths/min; SaO2Percent %; the other items bool; result percent 0–100
// Range:   WeightLb > 0; HighestRespRatePerMin ≥ 0; SaO2Percent 0–100; adults with active cancer and pulmonary embolism
// Example: 150 lb, rate 20, SaO₂ 95 %, no other item → L = 3.718 + 0.88 − 5.985 − 1.8 = −3.187 → 3.96 % (hand-derived from the equation)
func POMPECMortalityPercent(f POMPECFindings) (float64, error) {
	if err := checkPositiveFinite("WeightLb", f.WeightLb); err != nil {
		return 0, err
	}
	if err := checkFiniteNonNegative("HighestRespRatePerMin", f.HighestRespRatePerMin); err != nil {
		return 0, err
	}
	if err := checkClosedRange("SaO2Percent", f.SaO2Percent, 0, 100); err != nil {
		return 0, err
	}
	l := 3.718 + 0.044*f.HighestRespRatePerMin - 0.063*f.SaO2Percent - 0.012*f.WeightLb
	if f.DoNotResuscitate {
		l += 1.551
	}
	if f.RespiratoryDistress {
		l += 0.800
	}
	if f.UnilateralLimbSwelling {
		l += 0.734
	}
	if f.AlteredMentalStatus {
		l += 1.473
	}
	if f.HeartRateAbove99 {
		l += 1.027
	}
	return 100 / (1 + math.Exp(-l)), nil
}

// POMPECMortalityPercentKg returns the same POMPE-C estimate as POMPECMortalityPercent with the weight in kg,
// converted once at entry to lb (1 lb = 0.45359237 kg exactly); f.WeightLb is ignored.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range.
//
// ID:      MED-RESP-038
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Kline JA, Roy PM, Than MP, et al. Derivation and validation of a multivariate model to predict mortality from pulmonary embolism with cancer: the POMPE-C tool. Thromb Res. 2012;129(5):e194-e199. doi:10.1016/j.thromres.2012.03.015. PMID 22475313, PMC3350789, Table 2 and the equation for P
// Tags:    respiratory, pulmonary-embolism, cancer, pompe-c, mortality, logistic, metric
// Units:   weightKg kg; the rest as POMPECMortalityPercent; result percent 0–100
// Range:   weightKg > 0; otherwise as POMPECMortalityPercent
// Example: 68.0388555 kg (150 lb), rate 20, SaO₂ 95 % → 3.96 % (as the lb example)
func POMPECMortalityPercentKg(f POMPECFindings, weightKg float64) (float64, error) {
	if err := checkPositiveFinite("weightKg", weightKg); err != nil {
		return 0, err
	}
	f.WeightLb = weightKg / kgPerPound
	return POMPECMortalityPercent(f)
}

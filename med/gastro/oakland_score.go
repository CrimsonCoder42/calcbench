package gastro

import "fmt"

// OaklandFindings holds the Oakland score items for acute lower gastrointestinal bleeding (Oakland 2017,
// Table 2).
type OaklandFindings struct {
	AgeYears          float64 // completed age, years: < 40: 0, 40–69: 1, "> 70": 2
	Male              bool    // gender male: 1
	PreviousLGIBAdmit bool    // previous admission with LGIB: 1
	BloodOnDRE        bool    // digital rectal examination finding "Blood": 1
	HeartRateBpm      float64 // heart rate, beats/min: < 70: 0, 70–89: 1, 90–109: 2, "> 110": 3
	SystolicMmHg      float64 // systolic blood pressure, mmHg: 50–89: 5, 90–119: 4, 120–129: 3, 130–159: 2, "> 160": 0
	HaemoglobinGL     float64 // haemoglobin, g/L: 36–69: 22, 70–89: 17, 90–109: 13, 110–129: 8, 130–159: 4, "> 160": 0
}

// OaklandScore returns the Oakland score for acute lower gastrointestinal bleeding, 0–35 by the table (the
// source reports "The score ranges from 3-33", the range observed in its cohort). The printed bands are integer ranges with gaps at 70 years, 110 beats/min and
// 160 mmHg or g/L ("> 70" after "40-69"); they are read here as consecutive half-open intervals, so a value on
// an edge scores with the band that starts there (age 70 → 2, heart rate 110 → 3, SBP 160 → 0, Hb 160 → 0).
// Table 2 prints no Hb unit while Methods says g/dL; the band numbers (36–160) are g/L and are taken as g/L.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for age < 16 years (the cohort was adults ≥ 16),
// heart rate < 0, SBP < 50 mmHg or Hb < 36 g/L (below the lowest printed band).
//
// ID:      MED-GASTRO-013
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Oakland K, Jairath V, Uberoi R, Guy R, Ayaru L, Mortensen N, Murphy MF, Collins GS. Derivation and validation of a novel risk score for safe discharge after acute lower gastrointestinal bleeding: a modelling study. Lancet Gastroenterol Hepatol. 2017;2(9):635-643. doi:10.1016/S2468-1253(17)30150-4, PMID 28651935, Table 2 (author accepted manuscript, Oxford ORA)
// Tags:    gastro, gi-bleeding, lower-gi, oakland, score
// Units:   AgeYears years; HeartRateBpm beats/min; SystolicMmHg mmHg; HaemoglobinGL g/L; result points 0–35
// Range:   AgeYears ≥ 16; SystolicMmHg ≥ 50; HaemoglobinGL ≥ 36
// Example: man 75, no prior LGIB, no blood on DRE, HR 95, SBP 125, Hb 100 → 2+1+0+0+2+3+13 = 21 (hand-derived from Table 2)
func OaklandScore(f OaklandFindings) (int, error) {
	for _, in := range [...]struct {
		name string
		v    float64
	}{{"AgeYears", f.AgeYears}, {"HeartRateBpm", f.HeartRateBpm}, {"SystolicMmHg", f.SystolicMmHg}, {"HaemoglobinGL", f.HaemoglobinGL}} {
		if err := checkFinite(in.name, in.v); err != nil {
			return 0, err
		}
	}
	switch {
	case f.AgeYears < 16:
		return 0, fmt.Errorf("%w: AgeYears = %v, want ≥ 16", ErrOutOfRange, f.AgeYears)
	case f.HeartRateBpm < 0:
		return 0, fmt.Errorf("%w: HeartRateBpm = %v, want ≥ 0", ErrOutOfRange, f.HeartRateBpm)
	case f.SystolicMmHg < 50:
		return 0, fmt.Errorf("%w: SystolicMmHg = %v, want ≥ 50", ErrOutOfRange, f.SystolicMmHg)
	case f.HaemoglobinGL < 36:
		return 0, fmt.Errorf("%w: HaemoglobinGL = %v, want ≥ 36", ErrOutOfRange, f.HaemoglobinGL)
	}
	s := 0
	switch {
	case f.AgeYears >= 70:
		s += 2
	case f.AgeYears >= 40:
		s++
	}
	if f.Male {
		s++
	}
	if f.PreviousLGIBAdmit {
		s++
	}
	if f.BloodOnDRE {
		s++
	}
	switch hr := f.HeartRateBpm; {
	case hr >= 110:
		s += 3
	case hr >= 90:
		s += 2
	case hr >= 70:
		s++
	}
	switch sbp := f.SystolicMmHg; {
	case sbp >= 160:
	case sbp >= 130:
		s += 2
	case sbp >= 120:
		s += 3
	case sbp >= 90:
		s += 4
	default:
		s += 5
	}
	switch hb := f.HaemoglobinGL; {
	case hb >= 160:
	case hb >= 130:
		s += 4
	case hb >= 110:
		s += 8
	case hb >= 90:
		s += 13
	case hb >= 70:
		s += 17
	default:
		s += 22
	}
	return s, nil
}

package critical

// SIRSResult holds which of the four SIRS criteria are met, their count and whether the count reaches the
// definition's "two or more".
type SIRSResult struct {
	Temperature bool // temperature > 38 °C or < 36 °C
	HeartRate   bool // heart rate > 90/min
	Respiratory bool // respiratory rate > 20/min or PaCO2 < 32 mmHg
	WhiteCells  bool // white blood cell count > 12 000/mm³ or < 4000/mm³ or > 10% immature bands
	Count       int  // criteria met, 0–4
	SIRS        bool // Count ≥ 2 ("Two or more of")
}

// SIRSCriteria returns the systemic inflammatory response syndrome criteria of the 1991 ACCP/SCCM consensus
// conference: temperature > 38 °C or < 36 °C; heart rate > 90/min; respiratory rate > 20/min or PaCO2 < 32 mmHg;
// white blood cell count > 12 000/mm³ or < 4000/mm³ or > 10% immature bands; SIRS is two or more. PaCO2 enters
// only when paco2Measured is true; an unmeasured band count is passed as 0 (it then meets nothing).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a negative value, a temperature outside 10–50 °C or a
// band fraction above 100%.
//
// ID:      MED-CRIT-013
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Bone RC, Balk RA, Cerra FB, et al. Definitions for sepsis and organ failure and guidelines for the use of innovative therapies in sepsis. The ACCP/SCCM Consensus Conference Committee. Chest. 1992;101(6):1644-1655. doi:10.1378/chest.101.6.1644. PMID 1303622 (primary; not opened)
// Source:  Singer M, Deutschman CS, Seymour CW, et al. The Third International Consensus Definitions for Sepsis and Septic Shock (Sepsis-3). JAMA. 2016;315(8):801-810. doi:10.1001/jama.2016.0287. PMID 26903338. PMC4968574, Box 1 "SIRS (Systemic Inflammatory Response Syndrome)", "From Bone et al"
// Source:  Balk RA. Systemic inflammatory response syndrome (SIRS): where did it come from and is it still relevant today? Virulence. 2014;5(1):20-26. doi:10.4161/viru.27135. PMID 24280933. PMC3916374, the 1991 ACCP/SCCM SIRS criteria, citing Bone 1992
// Verified: secondary (Singer 2016 Box 1; Balk 2014)
// Tags:    critical, sepsis, sirs, inflammation, criteria
// Units:   temperatureC °C; heartRatePerMin beats/min; respirationRatePerMin breaths/min; paco2MmHg mmHg; wbcPerMm3 cells/mm³ (= cells/µL); immatureBandsPercent % of white cells
// Range:   10 ≤ temperatureC ≤ 50; heartRatePerMin, respirationRatePerMin, paco2MmHg, wbcPerMm3 ≥ 0; 0 ≤ immatureBandsPercent ≤ 100 (physical limits; the definition has no age limit, adults in the source)
// Example: 38.5 °C, HR 95, RR 18, PaCO2 not measured, WBC 9000, bands 0 → Count 2, SIRS true (from Box 1)
func SIRSCriteria(temperatureC, heartRatePerMin, respirationRatePerMin, paco2MmHg float64, paco2Measured bool, wbcPerMm3, immatureBandsPercent float64) (SIRSResult, error) {
	for _, e := range [...]error{
		finite("temperatureC", temperatureC), finite("heartRatePerMin", heartRatePerMin),
		finite("respirationRatePerMin", respirationRatePerMin), finite("paco2MmHg", paco2MmHg),
		finite("wbcPerMm3", wbcPerMm3), finite("immatureBandsPercent", immatureBandsPercent),
		within("temperatureC", temperatureC, 10, 50), atLeastZero("heartRatePerMin", heartRatePerMin),
		atLeastZero("respirationRatePerMin", respirationRatePerMin), atLeastZero("paco2MmHg", paco2MmHg),
		atLeastZero("wbcPerMm3", wbcPerMm3), within("immatureBandsPercent", immatureBandsPercent, 0, 100),
	} {
		if e != nil {
			return SIRSResult{}, e
		}
	}
	var r SIRSResult
	r.Temperature = temperatureC > 38 || temperatureC < 36
	r.HeartRate = heartRatePerMin > 90
	r.Respiratory = respirationRatePerMin > 20 || (paco2Measured && paco2MmHg < 32)
	r.WhiteCells = wbcPerMm3 > 12000 || wbcPerMm3 < 4000 || immatureBandsPercent > 10
	for _, met := range [...]bool{r.Temperature, r.HeartRate, r.Respiratory, r.WhiteCells} {
		if met {
			r.Count++
		}
	}
	r.SIRS = r.Count >= 2
	return r, nil
}

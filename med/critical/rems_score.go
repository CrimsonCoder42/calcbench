package critical

// REMSScore returns the Rapid Emergency Medicine Score, 0–26, the sum of six item scores. Age (years): < 45 → 0;
// 45–54 → 2; 55–64 → 3; 65–74 → 5; > 74 → 6. Mean arterial pressure (mmHg): 70–109 → 0; 110–129 or 50–69 → 2;
// 130–159 → 3; > 159 or ≤ 49 → 4. Heart rate (/min): 70–109 → 0; 110–139 or 55–69 → 2; 140–179 or 40–54 → 3;
// > 179 or ≤ 39 → 4. Respiratory rate (/min): 12–24 → 0; 25–34 or 10–11 → 1; 6–9 → 2; 35–49 → 3; > 49 or ≤ 5 → 4.
// SpO2 (%): > 89 → 0; 86–89 → 1; 75–85 → 3; < 75 → 4. GCS: 14–15 → 0; 11–13 → 1; 8–10 → 2; 5–7 → 3; 3–4 → 4.
// The table prints integer bands; the open-ended bands apply as printed (≤ 49, > 159, > 74, > 89 …) and any other
// value between two printed bands scores with the band whose lower edge it has reached (e.g. MAP 49.5 → 2, age
// 74.5 → 6, age 44.9 → 0, SpO2 89.5 → 0).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a negative value, an SpO2 above 100 or a GCS outside 3–15.
//
// ID:      MED-CRIT-018
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Olsson T, Terent A, Lind L. Rapid Emergency Medicine score: a new prognostic tool for in-hospital mortality in nonsurgical emergency department patients. J Intern Med. 2004;255(5):579-587. doi:10.1111/j.1365-2796.2004.01321.x. PMID 15078500 (primary; not opened)
// Source:  Imhoff BF, Thompson NJ, Hastings MA, et al. Rapid Emergency Medicine Score (REMS) in the trauma population: a retrospective study. BMJ Open. 2014;4(5):e004738. doi:10.1136/bmjopen-2013-004738. PMID 24793256. PMC4024603, table "REMS scoring system" (after Olsson 2004)
// Source:  Ala A, Vahdati SS, Jalali M, et al. Rapid Emergency Medicine Score as a Predictive Value for 30-day Outcome of Nonsurgical Patients Referred to the Emergency Department. Indian J Crit Care Med. 2020;24(6):418-422. doi:10.5005/jp-journals-10071-23456. PMID 32863634. PMC7435099, REMS table (after Olsson 2004)
// Verified: secondary (Imhoff 2014; Ala 2020)
// Tags:    critical, emergency, rems, early warning, score
// Units:   ageYears years; mapMmHg mmHg; heartRatePerMin beats/min; respirationRatePerMin breaths/min; spo2Percent %; gcs points; result points 0–26
// Range:   every numeric input ≥ 0; spo2Percent ≤ 100; 3 ≤ gcs ≤ 15 (physical limits; adults in non-surgical emergency care, the source's population)
// Example: age 70, MAP 65, HR 115, RR 26, SpO2 88, GCS 13 → 5+2+2+1+1+1 = 12 (HAND-DERIVED from the table)
func REMSScore(ageYears, mapMmHg, heartRatePerMin, respirationRatePerMin, spo2Percent float64, gcs int) (int, error) {
	for _, e := range [...]error{
		finite("ageYears", ageYears), finite("mapMmHg", mapMmHg), finite("heartRatePerMin", heartRatePerMin),
		finite("respirationRatePerMin", respirationRatePerMin), finite("spo2Percent", spo2Percent),
		atLeastZero("ageYears", ageYears), atLeastZero("mapMmHg", mapMmHg), atLeastZero("heartRatePerMin", heartRatePerMin),
		atLeastZero("respirationRatePerMin", respirationRatePerMin), within("spo2Percent", spo2Percent, 0, 100),
		checkGCS(gcs),
	} {
		if e != nil {
			return 0, e
		}
	}
	s := 0
	switch a := ageYears; {
	case a > 74:
		s += 6
	case a >= 65:
		s += 5
	case a >= 55:
		s += 3
	case a >= 45:
		s += 2
	}
	switch m := mapMmHg; {
	case m <= 49 || m > 159:
		s += 4
	case m >= 130:
		s += 3
	case m < 70 || m >= 110:
		s += 2
	}
	switch h := heartRatePerMin; {
	case h <= 39 || h > 179:
		s += 4
	case h < 55 || h >= 140:
		s += 3
	case h < 70 || h >= 110:
		s += 2
	}
	switch r := respirationRatePerMin; {
	case r <= 5 || r > 49:
		s += 4
	case r >= 35:
		s += 3
	case r < 10:
		s += 2
	case r < 12 || r >= 25:
		s++
	}
	switch o := spo2Percent; {
	case o < 75:
		s += 4
	case o < 86:
		s += 3
	case o <= 89:
		s++
	}
	switch {
	case gcs <= 4:
		s += 4
	case gcs <= 7:
		s += 3
	case gcs <= 10:
		s += 2
	case gcs <= 13:
		s++
	}
	return s, nil
}

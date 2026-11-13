package neuro

// CAMICUPositive applies the Confusion Assessment Method for the ICU algorithm: positive when feature 1 (acute
// onset of mental status change or a fluctuating course) and feature 2 (inattention) are present, plus either
// feature 3 (disorganized thinking) or feature 4 (altered level of consciousness). Whether each feature is
// present, and whether the person can be assessed at all, is the caller's.
//
// ID:      MED-NEURO-027
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Ely EW, Inouye SK, Bernard GR, et al. Delirium in mechanically ventilated patients: validity and reliability of the confusion assessment method for the intensive care unit (CAM-ICU). JAMA. 2001;286(21):2703-2710. doi:10.1001/jama.286.21.2703, PMID 11730446 (full text not open)
// Source:  Pipanmekaporn T, Wongpakaran N, Mueankwan S, et al. Validity and reliability of the Thai version of the Confusion Assessment Method for the Intensive Care Unit (CAM-ICU). Clin Interv Aging. 2014;9:879-885. doi:10.2147/CIA.S62660, PMC4043427, Methods ("Delirium is diagnosed when a participant manifests features 1 and 2, plus either feature 3 or 4")
// Source:  Kim JY, Kang YA, Kang PJ. Eur J Med Res. 2026, PMC12964761 ("A patient is considered CAM-ICU positive when features 1 and 2 are present with either feature 3 or 4"; ref 4 Ely 2001 JAMA)
// Verified: secondary (Pipanmekaporn 2014 PMC4043427, ref 30 Ely 2001; Kim 2026 PMC12964761, ref 4 Ely 2001)
// Tags:    neuro, delirium, icu, cam-icu, decision rule
// Units:   features present yes/no; result true when the algorithm is positive
// Range:   all inputs boolean
// Example: 1, 2 and 4 present → true; 1 and 2 only → false; 2, 3 and 4 without 1 → false
func CAMICUPositive(acuteOnsetOrFluctuating, inattention, disorganizedThinking, alteredConsciousness bool) bool {
	return acuteOnsetOrFluctuating && inattention && (disorganizedThinking || alteredConsciousness)
}

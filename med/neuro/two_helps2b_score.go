package neuro

// TwoHELPS2BScore returns the 2HELPS2B score (0–7) from EEG findings and history: brief (ictal) rhythmic
// discharges 2 points; lateralized periodic discharges, lateralized rhythmic delta activity or bilateral
// independent periodic discharges 1; prior seizure 1; sporadic epileptiform discharges 1; frequency > 2.0 Hz
// for any periodic or rhythmic pattern 1; "plus" features 1.
//
// ID:      MED-NEURO-012
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Struck AF, Ustun B, Ruiz AR, et al. Association of an electroencephalography-based risk score with seizure probability in hospitalized patients. JAMA Neurol. 2017;74(12):1419-1424. doi:10.1001/jamaneurol.2017.2459, PMID 29052706, PMC5822188 ("brief (ictal) rhythmic discharges (B[I]RDs) (2 points)", the five 1-point findings, score 0–7)
// Tags:    neuro, seizure, eeg, 2helps2b, score
// Units:   each input present or absent; result points 0–7
// Range:   all inputs boolean
// Example: BIRDs, LPDs, frequency > 2 Hz → 2 + 1 + 1 = 4 (hand-derived from the source)
func TwoHELPS2BScore(briefRhythmicDischarges, lpdLRDAOrBIPD, priorSeizure, sporadicEpileptiformDischarges, frequencyAbove2Hz, plusFeatures bool) int {
	n := 0
	if briefRhythmicDischarges {
		n += 2
	}
	for _, b := range [5]bool{lpdLRDAOrBIPD, priorSeizure, sporadicEpileptiformDischarges, frequencyAbove2Hz, plusFeatures} {
		if b {
			n++
		}
	}
	return n
}

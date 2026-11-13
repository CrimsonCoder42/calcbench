package cardio

import "fmt"

// ShanghaiECG is the ECG item of the Shanghai Brugada score; one of the three must apply.
type ShanghaiECG int

// ECG items (Antzelevitch 2016, Table 2, section I). The zero value is not a valid item.
const (
	ShanghaiECGSpontaneousType1      ShanghaiECG = iota + 1 // "Spontaneous Type 1 Brugada ECG pattern at nominal or high leads": 3.5
	ShanghaiECGFeverInducedType1                            // "Fever-induced type 1": 3
	ShanghaiECGDrugConvertedType2Or3                        // "Type 2 or 3 that converts with provocative drug challenge": 2
)

// ShanghaiHistory are the clinical-history and family-history items (Table 2, sections II and III); within
// each section only the highest item present counts.
type ShanghaiHistory struct {
	CardiacArrestOrVFPolyVT    bool // II "Unexplained cardiac arrest or documented VF/polymorphic VT": 3
	NocturnalAgonalRespiration bool // II "Nocturnal agonal respirations": 2
	SuspectedArrhythmicSyncope bool // II "Suspected arrhythmic syncope": 2
	SyncopeUnclearMechanism    bool // II "Syncope of unclear mechanism/unclear etiology": 1
	AFFlutterUnder30           bool // II AF or flutter under 30 years without alternative etiology: 0.5
	RelativeDefiniteBrS        bool // III first- or second-degree relative with definite BrS: 2
	RelativeSuspiciousSCD      bool // III suspicious SCD (fever, nocturnal, Brugada aggravating drugs) in a 1st/2nd-degree relative: 1
	RelativeUnexplainedSCDU45  bool // III unexplained SCD < 45 y in a 1st/2nd-degree relative with negative autopsy: 0.5
	ProbablePathogenicMutation bool // IV "Probable pathogenic mutation in BrS susceptibility gene": 0.5
}

// ShanghaiBrugadaScore returns the Shanghai score for Brugada syndrome, in points (multiples of 0.5, 2–9):
// the ECG item plus the highest item present in each of sections II and III, plus section IV. The probability
// bands of the table are not applied (the open text renders the top band as "> 3.5", which conflicts with the
// minimum of 3.5 for a spontaneous type 1 pattern; to be read in the PDF).
// Errors: ErrOutOfRange when ecg is not one of the three ECG items ("One item from this category must apply").
//
// ID:      MED-CARDIO-015
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Antzelevitch C, Yan GX, Ackerman MJ, et al. J-Wave syndromes expert consensus conference report: emerging concepts and gaps in knowledge. Heart Rhythm. 2016;13(10):e295-e324. doi:10.1016/j.hrthm.2016.05.024. PMID 27423412. PMC5035208, Table 2 (Proposed Shanghai Score System for diagnosis of Brugada syndrome); footnote "Only award points once for highest score within this category. One item from this category must apply."
// Tags:    cardio, ecg, brugada, shanghai, channelopathy, score
// Units:   result points
// Range:   ecg one of the three ShanghaiECG items
// Example: spontaneous type 1 + suspected arrhythmic syncope + relative with definite BrS → 3.5 + 2 + 2 = 7.5 (HAND-DERIVED from Table 2)
func ShanghaiBrugadaScore(ecg ShanghaiECG, h ShanghaiHistory) (float64, error) {
	var s float64
	switch ecg {
	case ShanghaiECGSpontaneousType1:
		s = 3.5
	case ShanghaiECGFeverInducedType1:
		s = 3
	case ShanghaiECGDrugConvertedType2Or3:
		s = 2
	default:
		return 0, fmt.Errorf("%w: ecg = %d, want one of the three ECG items", ErrOutOfRange, ecg)
	}
	switch {
	case h.CardiacArrestOrVFPolyVT:
		s += 3
	case h.NocturnalAgonalRespiration, h.SuspectedArrhythmicSyncope:
		s += 2
	case h.SyncopeUnclearMechanism:
		s++
	case h.AFFlutterUnder30:
		s += 0.5
	}
	switch {
	case h.RelativeDefiniteBrS:
		s += 2
	case h.RelativeSuspiciousSCD:
		s++
	case h.RelativeUnexplainedSCDU45:
		s += 0.5
	}
	if h.ProbablePathogenicMutation {
		s += 0.5
	}
	return s, nil
}

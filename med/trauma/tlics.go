package trauma

import "fmt"

// TLICSMorphology is the injury-morphology category of the Thoracolumbar Injury Classification and Severity Score.
type TLICSMorphology int

// TLICS morphology categories (Lee 2005 Table 1).
const (
	TLICSCompression             TLICSMorphology = iota + 1 // "Compression" (1 point)
	TLICSBurst                                              // "Burst" (2 points)
	TLICSTranslationalRotational                            // "Translational/rotational" (3 points)
	TLICSDistraction                                        // "Distraction" (4 points)
)

// TLICSNeurologic is the neurologic-involvement category of TLICS.
type TLICSNeurologic int

// TLICS neurologic categories (Lee 2005 Table 1).
const (
	TLICSNeuroIntact         TLICSNeurologic = iota // "Intact" (0)
	TLICSNerveRoot                                  // "Nerve root" (2)
	TLICSCordConusIncomplete                        // "Cord, conus medullaris incomplete" (3)
	TLICSCordConusComplete                          // "Cord, conus medullaris complete" (2)
	TLICSCaudaEquina                                // "Cauda equina" (3)
)

// TLICSPLC is the posterior-ligamentous-complex category of TLICS.
type TLICSPLC int

// TLICS posterior ligamentous complex categories (Lee 2005 Table 1).
const (
	TLICSPLCIntact        TLICSPLC = iota // "Intact" (0)
	TLICSPLCIndeterminate                 // "Injury suspected/indeterminate" (2)
	TLICSPLCInjured                       // "Injured" (3)
)

// TLICS returns the Thoracolumbar Injury Classification and Severity Score, 1–10: morphology (compression 1,
// burst 2, translational/rotational 3, distraction 4) + neurologic involvement (intact 0, nerve root 2, cord or
// conus incomplete 3, cord or conus complete 2, cauda equina 3) + posterior ligamentous complex (intact 0,
// suspected/indeterminate 2, injured 3). The caller passes one category per item.
// Errors: ErrOutOfRange for a category value outside its constants.
//
// ID:      MED-TRAUMA-015
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Lee JY, Vaccaro AR, Lim MR, et al. Thoracolumbar injury classification and severity score: a new paradigm for the treatment of thoracolumbar spine trauma. J Orthop Sci. 2005;10(6):671-675. doi:10.1007/s00776-005-0956-y. PMID 16307197. PMC2779435, Table 1 (points) and Table 2 (total-score thresholds ≤3, 4, ≥5)
// Tags:    trauma, spine, thoracolumbar, tlics, score
// Units:   categories; result points 1–10
// Range:   the categories of Table 1 (thoracolumbar spine trauma, the source's population)
// Example: burst fracture, intact neurology, PLC indeterminate → 2 + 0 + 2 = 4 (Table 1)
func TLICS(morphology TLICSMorphology, neuro TLICSNeurologic, plc TLICSPLC) (int, error) {
	var m, n, p int
	switch morphology {
	case TLICSCompression:
		m = 1
	case TLICSBurst:
		m = 2
	case TLICSTranslationalRotational:
		m = 3
	case TLICSDistraction:
		m = 4
	default:
		return 0, fmt.Errorf("%w: morphology = %d, want 1–4", ErrOutOfRange, morphology)
	}
	switch neuro {
	case TLICSNeuroIntact:
		n = 0
	case TLICSNerveRoot, TLICSCordConusComplete:
		n = 2
	case TLICSCordConusIncomplete, TLICSCaudaEquina:
		n = 3
	default:
		return 0, fmt.Errorf("%w: neuro = %d, want 0–4", ErrOutOfRange, neuro)
	}
	switch plc {
	case TLICSPLCIntact:
		p = 0
	case TLICSPLCIndeterminate:
		p = 2
	case TLICSPLCInjured:
		p = 3
	default:
		return 0, fmt.Errorf("%w: plc = %d, want 0–2", ErrOutOfRange, plc)
	}
	return m + n + p, nil
}

package rheumatology

import (
	"fmt"
	"math"
)

// GoutACREULAR2015Score returns the 2015 ACR/EULAR gout classification criteria domain score, −6 to 23 (Table 2,
// "The maximum possible score in the final criteria is 23"), applied when the entry criterion is met and the
// sufficient criterion (MSU crystals or tophus) is not; both are the caller's to apply. Points (Table 2):
//   - jointPattern: 0 neither; 1 ankle or midfoot (mono/oligoarticular without first MTP); 2 first MTP → 0, 1, 2
//   - episodeCharacteristics: how many of erythema, can't bear touch, great difficulty walking (0–3) → 0–3
//   - timeCourse: 0 no typical episode; 1 one typical episode; 2 recurrent typical episodes → 0, 1, 2
//   - clinicalTophus → 4
//   - serumUrateMgDL: < 4 → −4; 4–< 6 → 0; 6–< 8 → 2; 8–< 10 → 3; ≥ 10 → 4 (the highest value, Table 2 note)
//   - synovialFluidMSUNegative (examined and no MSU crystals) → −2; "If synovial fluid was not assessed, score this item as 0"
//   - imagingUrateDeposition (ultrasound double-contour sign or DECT) → 4
//   - imagingGoutErosion (conventional radiography, at least 1 erosion) → 4
//
// The source classifies at a score ≥ 8; that is the caller's to apply.
// Errors: ErrNotFinite for a NaN or ±Inf urate; ErrOutOfRange outside Range, naming the item.
//
// ID:      MED-RHEUM-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Neogi T, Jansen TL, Dalbeth N, et al. 2015 Gout Classification Criteria: an American College of Rheumatology/European League Against Rheumatism collaborative initiative. Arthritis Rheumatol. 2015;67(10):2557-2568. doi:10.1002/art.39254, PMID 26352873, Table 2
// Tags:    rheumatology, gout, classification criteria, acr, eular, urate, score
// Units:   serumUrateMgDL mg/dL; other items category codes or present/absent; result points −6 to 23
// Range:   0 ≤ jointPattern ≤ 2; 0 ≤ episodeCharacteristics ≤ 3; 0 ≤ timeCourse ≤ 2; 0 ≤ serumUrateMgDL ≤ 30
// Example: first MTP (2), 3 characteristics (3), recurrent (2), urate 9 (3) → 10; urate 3.9 alone → −4; maximum 2+3+2+4+4+4+4 = 23 (Table 2)
func GoutACREULAR2015Score(jointPattern, episodeCharacteristics, timeCourse int, clinicalTophus bool, serumUrateMgDL float64, synovialFluidMSUNegative, imagingUrateDeposition, imagingGoutErosion bool) (int, error) {
	if math.IsNaN(serumUrateMgDL) || math.IsInf(serumUrateMgDL, 0) {
		return 0, fmt.Errorf("%w: serumUrateMgDL = %v", ErrNotFinite, serumUrateMgDL)
	}
	switch {
	case jointPattern < 0 || jointPattern > 2:
		return 0, fmt.Errorf("%w: jointPattern = %d, want 0–2", ErrOutOfRange, jointPattern)
	case episodeCharacteristics < 0 || episodeCharacteristics > 3:
		return 0, fmt.Errorf("%w: episodeCharacteristics = %d, want 0–3", ErrOutOfRange, episodeCharacteristics)
	case timeCourse < 0 || timeCourse > 2:
		return 0, fmt.Errorf("%w: timeCourse = %d, want 0–2", ErrOutOfRange, timeCourse)
	case serumUrateMgDL < 0 || serumUrateMgDL > 30:
		return 0, fmt.Errorf("%w: serumUrateMgDL = %v, want 0–30", ErrOutOfRange, serumUrateMgDL)
	}
	var urate int
	switch u := serumUrateMgDL; {
	case u < 4:
		urate = -4
	case u < 6:
		urate = 0
	case u < 8:
		urate = 2
	case u < 10:
		urate = 3
	default:
		urate = 4
	}
	return jointPattern + episodeCharacteristics + timeCourse + pts(clinicalTophus, 4) + urate +
		pts(synovialFluidMSUNegative, -2) + pts(imagingUrateDeposition, 4) + pts(imagingGoutErosion, 4), nil
}

package infectious

import (
	"errors"
	"fmt"
)

var (
	// ErrOutOfRange reports an input outside the values the source defines.
	ErrOutOfRange = errors.New("infectious: input outside the range the source validated")
)

// RASPSource is Step A of the RASP (source population); the constant's value is the source's score.
type RASPSource int

// Step A of Vertesi 2003, Fig. 1.
const (
	RASPSourceAcuteAIDS       RASPSource = 1    // "Known HIV carrier: Acute AIDS illness" (end-stage AIDS, hospitalized, known high-viral load)
	RASPSourceAsymptomatic    RASPSource = 10   // "Known HIV carrier: Asymptomatic"
	RASPSourceUnknownHighRisk RASPSource = 100  // "Unknown HIV status: High-risk situation"
	RASPSourceUnknownLowRisk  RASPSource = 1000 // "Unknown HIV status: Low-risk situation (other)"
)

// RASPInoculum is Step B of the RASP (inoculum type); the constant's value is the source's score.
type RASPInoculum int

// Step B of Vertesi 2003, Fig. 1.
const (
	RASPInoculumFreshBlood        RASPInoculum = 1    // "Fresh blood"
	RASPInoculumBodyFluidsAtRisk  RASPInoculum = 10   // "Body fluids at risk (e.g., semen)"
	RASPInoculumDriedOldBlood     RASPInoculum = 100  // "Dried old blood"
	RASPInoculumLowRiskSecretions RASPInoculum = 1000 // "Low-risk secretions (tears, saliva, urine)"
)

// RASPMethod is Step C of the RASP (method of transmission); the constant's value is the source's score.
type RASPMethod int

// Step C of Vertesi 2003, Fig. 1.
const (
	RASPMethodIntravenous           RASPMethod = 1    // "Intravenous"
	RASPMethodDeepIntramuscular     RASPMethod = 10   // "Deep intramuscular"
	RASPMethodDeepWithBleeding      RASPMethod = 100  // "Deep transcutaneous with visible bleeding at site"
	RASPMethodSuperficialNoBleeding RASPMethod = 200  // "Superficial transcutaneous with no visible bleeding"
	RASPMethodMucosal               RASPMethod = 500  // "Mucosal contact only"
	RASPMethodIntactSkin            RASPMethod = 1000 // "Intact skin"
)

// RASPVolume is Step D of the RASP (volume of inoculum); the constant's value is the source's modifier.
type RASPVolume int

// Step D of Vertesi 2003, Fig. 1.
const (
	RASPVolumeMassive    RASPVolume = 100 // "Massive (e.g., transfusion)"
	RASPVolumeMeasurable RASPVolume = 10  // "Measurable ( > 1 mL)"
	RASPVolumeModerate   RASPVolume = 5   // "Moderate (large-bore hollow needle > 22 g)"
	RASPVolumeSmall      RASPVolume = 3   // "Small (small-bore hollow needle < 22 g)"
	RASPVolumeTrace      RASPVolume = 1   // "Trace surface only (e.g., suture needle)"
)

// HIVRASPTotalRisk returns the RASP "Total risk" of HIV transmission for one exposure, as a probability: "Total
// score (Z) = A × B × C", "Basic risk = 1 / Z", "Total risk = Basic risk × Modifier (D)", so the result is
// D / (A × B × C). The source's Table 2 attaches treatment suggestions to risk levels; those are advice and are
// not returned.
// Errors: ErrOutOfRange for a step value other than the constants of its type.
//
// ID:      MED-INFECT-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Vertesi L. Risk Assessment Stratification Protocol (RASP) to help patients decide on the use of postexposure prophylaxis for HIV exposure. CJEM. 2003;5(1):46-48. doi:10.1017/s1481803500008113, PMID 17659153, Fig. 1 (Steps A–D and the formulas) and Examples 1 and 2
// Tags:    infectious, hiv, postexposure, needlestick, rasp, risk
// Units:   each step as the source's constant; result probability 0–1 (1 in 1/result)
// Range:   each step one of its constants (values adjusted to "an average Canadian community", the source's setting)
// Example: asymptomatic carrier, fresh blood, deep intramuscular, small-bore needle → 3/1000; unknown low-risk, dried old blood, superficial, large-bore → 1/4,000,000 (Examples 1 and 2)
func HIVRASPTotalRisk(source RASPSource, inoculum RASPInoculum, method RASPMethod, volume RASPVolume) (float64, error) {
	switch source {
	case RASPSourceAcuteAIDS, RASPSourceAsymptomatic, RASPSourceUnknownHighRisk, RASPSourceUnknownLowRisk:
	default:
		return 0, fmt.Errorf("%w: source = %d, want 1, 10, 100 or 1000", ErrOutOfRange, source)
	}
	switch inoculum {
	case RASPInoculumFreshBlood, RASPInoculumBodyFluidsAtRisk, RASPInoculumDriedOldBlood, RASPInoculumLowRiskSecretions:
	default:
		return 0, fmt.Errorf("%w: inoculum = %d, want 1, 10, 100 or 1000", ErrOutOfRange, inoculum)
	}
	switch method {
	case RASPMethodIntravenous, RASPMethodDeepIntramuscular, RASPMethodDeepWithBleeding,
		RASPMethodSuperficialNoBleeding, RASPMethodMucosal, RASPMethodIntactSkin:
	default:
		return 0, fmt.Errorf("%w: method = %d, want 1, 10, 100, 200, 500 or 1000", ErrOutOfRange, method)
	}
	switch volume {
	case RASPVolumeMassive, RASPVolumeMeasurable, RASPVolumeModerate, RASPVolumeSmall, RASPVolumeTrace:
	default:
		return 0, fmt.Errorf("%w: volume = %d, want 100, 10, 5, 3 or 1", ErrOutOfRange, volume)
	}
	z := float64(source) * float64(inoculum) * float64(method)
	return float64(volume) / z, nil
}

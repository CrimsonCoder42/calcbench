package cardio

import "fmt"

// RESCUEIHCATimeOfDay is the time-of-day row of the RESCUE-IHCA score (Tonna 2022 Table 2).
type RESCUEIHCATimeOfDay int

// Time-of-day levels; the reference level (no points in Table 2) is 7:00 am to 2:59 pm.
const (
	RESCUEIHCADay     RESCUEIHCATimeOfDay = iota // 7am–2:59pm: 0 (reference)
	RESCUEIHCAEvening                            // "3pm - 10:59pm": +4
	RESCUEIHCANight                              // "11pm - 6:59am": +13
)

// RESCUEIHCAIllness is the illness-category row of the RESCUE-IHCA score.
type RESCUEIHCAIllness int

// Illness categories; the reference (no points in Table 2) is medical non-cardiac.
const (
	RESCUEIHCAMedicalNonCardiac  RESCUEIHCAIllness = iota // 0 (reference)
	RESCUEIHCAMedicalCardiac                              // −2
	RESCUEIHCASurgicalNonCardiac                          // −6
	RESCUEIHCASurgicalCardiac                             // −11
)

// RESCUEIHCARhythm is the presenting-rhythm row of the RESCUE-IHCA score.
type RESCUEIHCARhythm int

// Presenting rhythms; the reference (no points in Table 2) is asystole.
const (
	RESCUEIHCAAsystole      RESCUEIHCARhythm = iota // 0 (reference)
	RESCUEIHCAPEA                                   // −1
	RESCUEIHCAPulselessVT                           // "pVT": −5
	RESCUEIHCAPalpablePulse                         // "Palpable pulse initially": −5
	RESCUEIHCAVF                                    // −8
)

var (
	rescueTimePoints    = [3]int{0, 4, 13}
	rescueIllnessPoints = [4]int{0, -2, -6, -11}
	rescueRhythmPoints  = [5]int{0, -1, -5, -5, -8}
)

// RESCUEIHCAScore returns the RESCUE-IHCA mortality prediction score for in-hospital cardiac arrest treated
// with ECPR (Tonna 2022 Table 2): age ≤ 20 years 2 points and +2 for "each additional 10" years; pre-existing
// renal insufficiency +8; the time-of-day, illness and rhythm rows; +2 for each 10 minutes of arrest. Table 2
// does not say how a partial 10-year or 10-minute interval counts, so the caller passes whole intervals. The
// paper's mortality model is Bayesian and its probability for a score is not printed as a formula; none is
// returned here.
// Errors: ErrOutOfRange when an interval count or an enum is outside Range.
//
// ID:      MED-CARDIO-052
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Tonna JE, Selzman CH, Girotra S, et al. Resuscitation using ECPR during in-hospital cardiac arrest (RESCUE-IHCA) mortality prediction score and external validation. JACC Cardiovasc Interv. 2022;15(3):237-247. doi:10.1016/j.jcin.2021.09.032. PMID 35033471. PMC8837656, Table 2
// Tags:    cardio, cardiac-arrest, ecpr, in-hospital, rescue-ihca, risk-score
// Units:   ageTensAbove20 whole 10-year intervals above 20 years; arrestTensOfMinutes whole 10-minute intervals; result points
// Range:   ageTensAbove20 0–9; arrestTensOfMinutes 0–30; enums as declared (GWTG-Resuscitation ECPR cohort)
// Example: age 50 (3 intervals), renal insufficiency, night, surgical cardiac, VF, 40 min (4) → 8 + 8 + 13 − 11 − 8 + 8 = 18 (HAND-DERIVED from Table 2)
func RESCUEIHCAScore(ageTensAbove20 int, renalInsufficiency bool, t RESCUEIHCATimeOfDay, ill RESCUEIHCAIllness, r RESCUEIHCARhythm, arrestTensOfMinutes int) (int, error) {
	switch {
	case ageTensAbove20 < 0 || ageTensAbove20 > 9:
		return 0, fmt.Errorf("%w: ageTensAbove20 = %d, want 0–9", ErrOutOfRange, ageTensAbove20)
	case arrestTensOfMinutes < 0 || arrestTensOfMinutes > 30:
		return 0, fmt.Errorf("%w: arrestTensOfMinutes = %d, want 0–30", ErrOutOfRange, arrestTensOfMinutes)
	case t < RESCUEIHCADay || t > RESCUEIHCANight:
		return 0, fmt.Errorf("%w: time of day = %d, want 0–2", ErrOutOfRange, t)
	case ill < RESCUEIHCAMedicalNonCardiac || ill > RESCUEIHCASurgicalCardiac:
		return 0, fmt.Errorf("%w: illness = %d, want 0–3", ErrOutOfRange, ill)
	case r < RESCUEIHCAAsystole || r > RESCUEIHCAVF:
		return 0, fmt.Errorf("%w: rhythm = %d, want 0–4", ErrOutOfRange, r)
	}
	s := 2 + 2*ageTensAbove20 + 2*arrestTensOfMinutes
	if renalInsufficiency {
		s += 8
	}
	return s + rescueTimePoints[t] + rescueIllnessPoints[ill] + rescueRhythmPoints[r], nil
}

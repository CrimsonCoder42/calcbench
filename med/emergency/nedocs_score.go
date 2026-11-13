package emergency

import (
	"errors"
	"fmt"
	"math"
)

var (
	// ErrNotFinite reports a NaN or ±Inf input.
	ErrNotFinite = errors.New("emergency: value is NaN or infinite")
	// ErrOutOfRange reports an input outside the range the source validated.
	ErrOutOfRange = errors.New("emergency: input outside the range the source validated")
)

// NEDOCSInputs holds the NEDOCS variables at one moment (Ilhan 2020; Dzelebdzic 2026).
type NEDOCSInputs struct {
	EDCensus                float64 // C: total persons under care in the emergency department (census)
	EDBeds                  float64 // A: number of ED beds
	AdmitsWaiting           float64 // F: admitted persons waiting in the ED
	HospitalBeds            float64 // B: number of hospital beds
	Ventilators             float64 // D: persons on ventilators in the ED
	LongestAdmitHours       float64 // E: hours waited by the longest-waiting admitted person
	LongestWaitingRoomHours float64 // G: hours from registration of the longest-waiting person in the waiting room
}

// NEDOCSScore returns the National Emergency Department Overcrowding Score, "NEDOCS score = 85.8(C/A) +
// 600(F/B) + 13.4(D) + 0.93(E) + 5.64(G) – 20", not capped and not rounded (one secondary reports capping to
// 0–200; NEDOCSLevel accepts any value).
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for fewer than 1 ED or hospital bed, or a negative
// count or time.
//
// ID:      MED-EMERG-002
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Weiss SJ, Derlet R, Arndahl J, et al. Estimating the degree of emergency department overcrowding in academic medical centers: results of the National ED Overcrowding Study (NEDOCS). Acad Emerg Med. 2004;11(1):38-50. doi:10.1197/j.aem.2003.07.017, PMID 14709427 (primary; not opened)
// Source:  Ilhan B, Kunt MM, Damarsoy FF, Demir MC, Aksu NM. NEDOCS: is it really useful for detecting emergency department overcrowding today? Medicine (Baltimore). 2020;99(28):e20478. doi:10.1097/MD.0000000000020478, PMID 32664059, PMC7360290, Methods (equation)
// Source:  Dzelebdzic S, Djikanovic B, Bukumiric Z, et al. Daily, weekly, and monthly dynamics of crowding in emergency department of university clinical centre Serbia. Int J Emerg Med. 2026;19(1):195. doi:10.1186/s12245-026-01313-y, PMID 42552518, PMC13435447, Methods (equation and variable definitions)
// Verified: secondary (Ilhan 2020 PMC7360290; Dzelebdzic 2026 PMC13435447; also Getachew 2024 PMC10787261)
// Tags:    emergency, crowding, operations, nedocs, score
// Units:   counts as given; hours; result NEDOCS points
// Range:   EDBeds ≥ 1; HospitalBeds ≥ 1; other inputs ≥ 0
// Example: census 45, 30 ED beds, 10 admits, 300 beds, 1 ventilator, 6 h, 2 h → −20 + 128.7 + 20 + 13.4 + 5.58 + 11.28 = 158.96 (hand-derived)
func NEDOCSScore(in NEDOCSInputs) (float64, error) {
	vals := [...]struct {
		name string
		v    float64
		min  float64
	}{
		{"EDCensus", in.EDCensus, 0}, {"EDBeds", in.EDBeds, 1}, {"AdmitsWaiting", in.AdmitsWaiting, 0},
		{"HospitalBeds", in.HospitalBeds, 1}, {"Ventilators", in.Ventilators, 0},
		{"LongestAdmitHours", in.LongestAdmitHours, 0}, {"LongestWaitingRoomHours", in.LongestWaitingRoomHours, 0},
	}
	for _, x := range vals {
		if math.IsNaN(x.v) || math.IsInf(x.v, 0) {
			return 0, fmt.Errorf("%w: %s = %v", ErrNotFinite, x.name, x.v)
		}
		if x.v < x.min {
			return 0, fmt.Errorf("%w: %s = %v, want ≥ %v", ErrOutOfRange, x.name, x.v, x.min)
		}
	}
	return 85.8*(in.EDCensus/in.EDBeds) + 600*(in.AdmitsWaiting/in.HospitalBeds) + 13.4*in.Ventilators +
		0.93*in.LongestAdmitHours + 5.64*in.LongestWaitingRoomHours - 20, nil
}

// NEDOCSCategory is a NEDOCS level; values are the source's labels.
type NEDOCSCategory string

const (
	NEDOCSNotBusy                NEDOCSCategory = "not busy"                           // level 1: 0–20
	NEDOCSBusy                   NEDOCSCategory = "busy"                               // level 2: 21–60
	NEDOCSExtremelyBusy          NEDOCSCategory = "extremely busy but not overcrowded" // level 3: 61–100
	NEDOCSOvercrowded            NEDOCSCategory = "overcrowded"                        // level 4: 101–140
	NEDOCSSeverelyOvercrowded    NEDOCSCategory = "severely overcrowded"               // level 5: 141–180
	NEDOCSDangerouslyOvercrowded NEDOCSCategory = "dangerously overcrowded"            // level 6: 181–200
)

// NEDOCSLevel returns the NEDOCS level for a score: 0–20 not busy, 21–60 busy, 61–100 extremely busy but not
// overcrowded, 101–140 overcrowded, 141–180 severely overcrowded, 181–200 dangerously overcrowded. The
// printed bands are integer ranges; a non-integer score goes to the band whose upper edge it does not
// exceed (20.5 → busy). A score below 0 is "not busy" and above 200 "dangerously overcrowded", as the
// secondary that caps the score to 0–200 implies.
// Errors: ErrNotFinite for NaN or ±Inf.
//
// ID:      MED-EMERG-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Weiss SJ, Derlet R, Arndahl J, et al. Acad Emerg Med. 2004;11(1):38-50. doi:10.1197/j.aem.2003.07.017, PMID 14709427 (primary; not opened)
// Source:  Ilhan B, Kunt MM, Damarsoy FF, Demir MC, Aksu NM. Medicine (Baltimore). 2020;99(28):e20478. doi:10.1097/MD.0000000000020478, PMID 32664059, PMC7360290 (level names)
// Source:  Dzelebdzic S, Djikanovic B, Bukumiric Z, et al. Int J Emerg Med. 2026;19(1):195. doi:10.1186/s12245-026-01313-y, PMID 42552518, PMC13435447 (score bands and the 0–200 cap)
// Verified: secondary (Ilhan 2020 PMC7360290; Dzelebdzic 2026 PMC13435447)
// Tags:    emergency, crowding, operations, nedocs, category
// Units:   score NEDOCS points; result the source's label
// Range:   any finite score
// Example: 20 → not busy; 21 → busy; 100 → extremely busy but not overcrowded; 181 → dangerously overcrowded
func NEDOCSLevel(score float64) (NEDOCSCategory, error) {
	if math.IsNaN(score) || math.IsInf(score, 0) {
		return "", fmt.Errorf("%w: score = %v", ErrNotFinite, score)
	}
	switch {
	case score <= 20:
		return NEDOCSNotBusy, nil
	case score <= 60:
		return NEDOCSBusy, nil
	case score <= 100:
		return NEDOCSExtremelyBusy, nil
	case score <= 140:
		return NEDOCSOvercrowded, nil
	case score <= 180:
		return NEDOCSSeverelyOvercrowded, nil
	}
	return NEDOCSDangerouslyOvercrowded, nil
}

package endocrine

import "fmt"

// MyxedemaCNS is the central nervous system row of the myxedema coma diagnostic score.
type MyxedemaCNS int

const (
	MyxedemaCNSAbsent         MyxedemaCNS = iota // "Absent": 0
	MyxedemaCNSSomnolent                         // "Somnolent/Lethargy": 10
	MyxedemaCNSObtunded                          // "Obtunded": 15
	MyxedemaCNSStupor                            // "Stupor": 20
	MyxedemaCNSComaOrSeizures                    // "Coma/seizures": 30
)

// MyxedemaGI is the gastrointestinal row of the myxedema coma diagnostic score.
type MyxedemaGI int

const (
	MyxedemaGIAbsent                   MyxedemaGI = iota // "Absent": 0
	MyxedemaGIAnorexiaPainConstipation                   // "Anorexia/abdominal pain/constipation": 5
	MyxedemaGIDecreasedMotility                          // "Decreased intestinal motility": 15
	MyxedemaGIParalyticIleus                             // "Paralytic ileus": 20
)

// MyxedemaComaFindings holds the items of the diagnostic scoring system for myxedema coma (Popoveniuc 2014,
// as reproduced by Çetin 2025 and Amin 2023).
type MyxedemaComaFindings struct {
	TemperatureC       float64     // "> 35" °C: 0; 32–35: 10; "< 32": 20
	CNS                MyxedemaCNS // 0–30
	GI                 MyxedemaGI  // 0–20
	PrecipitatingEvent bool        // present: 10
	HeartRateBpm       float64     // bradycardia: absent (≥ 60) 0; 50–59: 10; 40–49: 20; < 40: 30
	OtherEKGChanges    bool        // 10 ("bradycardia, low QRS voltages, and widespread T wave inversions", Amin 2023)
	PericardialPleural bool        // pericardial/pleural effusion: 10
	PulmonaryEdema     bool        // 15
	Cardiomegaly       bool        // 15
	Hypotension        bool        // 20
	Hyponatremia       bool        // 10
	Hypoglycemia       bool        // 10
	Hypoxemia          bool        // 10
	Hypercarbia        bool        // 10
	DecreasedGFR       bool        // "Decrease in GFR": 10
}

// MyxedemaComaScore returns the myxedema coma diagnostic score, 0–230. Temperature 35 °C and 32 °C score 10
// (the row reads "> 95/35" 0, "89.6-95/32-35" 10, "< 89.6/32" 20); heart-rate bands are half-open (50 ≤ HR <
// 60 scores 10). The cut-off bands are not built here: the primary abstract gives "scores between 45 and 59"
// for at-risk while the secondaries give 25–59.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for an unknown level, temperature outside 20–45 °C or
// a negative heart rate.
//
// ID:      MED-ENDO-028
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Popoveniuc G, Chandra T, Sud A, et al. A diagnostic scoring system for myxedema coma. Endocr Pract. 2014;20(8):808-817. doi:10.4158/EP13460.OR, PMID 24518183 (primary; not opened)
// Source:  Çetin K, Donbaloğlu Z, Bahar YF, et al. A rare endocrinological emergency in children with hypothyroidism: clinical evaluation of cases presenting with myxedema coma. Turk Arch Pediatr. 2025;60(5):474-482. doi:10.5152/TurkArchPediatr.2025.25061, PMID 40959820, PMC12432182, table "The Diagnostic Scoring System for Myxedema Coma"
// Source:  Amin MS, Pabani UK, Lohano S, Khan Z. Myxedema coma precipitated by sepsis in a patient with a complex mental health history. Cureus. 2023;15(8):e43574. doi:10.7759/cureus.43574, PMID 37719600, PMC10503452, Table
// Verified: secondary (Çetin 2025 PMC12432182; Amin 2023 PMC10503452)
// Tags:    endocrine, thyroid, myxedema coma, hypothyroidism, score
// Units:   TemperatureC °C; HeartRateBpm beats/min; result points 0–230
// Range:   TemperatureC 20–45 (plausibility bound, not from the source); HeartRateBpm ≥ 0
// Example: 33 °C, obtunded, constipation, infection, HR 45, hyponatremia → 10+15+5+10+20+10 = 70 (hand-derived from the table)
func MyxedemaComaScore(f MyxedemaComaFindings) (int, error) {
	if err := checkFinite("TemperatureC", f.TemperatureC); err != nil {
		return 0, err
	}
	if err := checkFinite("HeartRateBpm", f.HeartRateBpm); err != nil {
		return 0, err
	}
	if f.TemperatureC < 20 || f.TemperatureC > 45 {
		return 0, fmt.Errorf("%w: TemperatureC = %v, want 20–45", ErrOutOfRange, f.TemperatureC)
	}
	if f.HeartRateBpm < 0 {
		return 0, fmt.Errorf("%w: HeartRateBpm = %v, want ≥ 0", ErrOutOfRange, f.HeartRateBpm)
	}
	cnsPoints := [...]int{0, 10, 15, 20, 30}
	giPoints := [...]int{0, 5, 15, 20}
	if f.CNS < 0 || int(f.CNS) >= len(cnsPoints) {
		return 0, fmt.Errorf("%w: CNS = %d, want 0–4", ErrOutOfRange, f.CNS)
	}
	if f.GI < 0 || int(f.GI) >= len(giPoints) {
		return 0, fmt.Errorf("%w: GI = %d, want 0–3", ErrOutOfRange, f.GI)
	}
	s := cnsPoints[f.CNS] + giPoints[f.GI]
	switch t := f.TemperatureC; {
	case t < 32:
		s += 20
	case t <= 35:
		s += 10
	}
	switch hr := f.HeartRateBpm; {
	case hr < 40:
		s += 30
	case hr < 50:
		s += 20
	case hr < 60:
		s += 10
	}
	for _, it := range [...]struct {
		present bool
		points  int
	}{
		{f.PrecipitatingEvent, 10}, {f.OtherEKGChanges, 10}, {f.PericardialPleural, 10}, {f.PulmonaryEdema, 15},
		{f.Cardiomegaly, 15}, {f.Hypotension, 20}, {f.Hyponatremia, 10}, {f.Hypoglycemia, 10}, {f.Hypoxemia, 10},
		{f.Hypercarbia, 10}, {f.DecreasedGFR, 10},
	} {
		if it.present {
			s += it.points
		}
	}
	return s, nil
}

package trauma

// WisconsinFindings are the Wisconsin criteria besides the GCS (Sitzman 2011 abstract).
type WisconsinFindings struct {
	BonyStepoffOrInstability     bool // "bony stepoff or instability"
	PeriorbitalSwellingContusion bool // "periorbital swelling or contusion"
	Malocclusion                 bool // "malocclusion"
	ToothAbsence                 bool // "tooth absence"
}

// WisconsinCriteria returns the number of the five Wisconsin criteria for maxillofacial CT present (0–5): bony
// stepoff or instability, periorbital swelling or contusion, Glasgow Coma Scale score less than 14, malocclusion,
// or tooth absence.
// Errors: ErrOutOfRange when gcs is outside 3–15.
//
// ID:      MED-TRAUMA-014
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Sitzman TJ, Hanson SE, Alsheik NH, Gentry LR, Doyle JF, Gutowski KA. Clinical criteria for obtaining maxillofacial computed tomographic scans in trauma patients. Plast Reconstr Surg. 2011;127(3):1270-1278. doi:10.1097/PRS.0b013e3182043ad8. PMID 21364427. Abstract: "bony stepoff or instability, periorbital swelling or contusion, Glasgow Coma Scale score less than 14, malocclusion, or tooth absence"
// Tags:    trauma, maxillofacial, facial fracture, wisconsin criteria, ct, decision rule
// Units:   gcs points 3–15; result count 0–5
// Range:   3 ≤ gcs ≤ 15 (trauma subjects, the study's population, which the caller establishes)
// Example: GCS 13, nothing else → 1; GCS 14, none → 0 (from "less than 14")
func WisconsinCriteria(gcs int, f WisconsinFindings) (int, error) {
	if err := checkGCS(gcs); err != nil {
		return 0, err
	}
	return count(f.BonyStepoffOrInstability, f.PeriorbitalSwellingContusion, gcs < 14, f.Malocclusion, f.ToothAbsence), nil
}

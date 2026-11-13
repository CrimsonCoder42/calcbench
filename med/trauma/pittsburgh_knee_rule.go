package trauma

// PittsburghKneeRule reports whether the Pittsburgh knee rule's combination is present: a fall or blunt trauma
// together with either inability to ambulate or age younger than 12 or older than 50 years. When fallOrBluntTrauma
// is false the combination is absent.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange when ageYears is outside 0–130.
//
// ID:      MED-TRAUMA-011
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Seaberg DC, Jackson R. Clinical decision rule for knee radiographs. Am J Emerg Med. 1994;12(5):541-543. doi:10.1016/0735-6757(94)90274-7. PMID 8060409. Abstract: "the combination of fall or blunt trauma with either inability to ambulate or age (younger than 12 to older than 50 years of age) was 100% sensitive, with a specificity of 79%"
// Tags:    trauma, knee, pittsburgh knee rule, decision rule, radiography
// Units:   ageYears years; result boolean
// Range:   0 ≤ ageYears ≤ 130 (physical limits; knee injuries, the study's population, which the caller establishes)
// Example: fall, age 51, walks → true; fall, age 50, walks → false; no fall or blunt trauma → false (from the abstract's combination)
func PittsburghKneeRule(ageYears float64, fallOrBluntTrauma, unableToAmbulate bool) (bool, error) {
	if err := checkAge(ageYears); err != nil {
		return false, err
	}
	return fallOrBluntTrauma && (unableToAmbulate || ageYears < 12 || ageYears > 50), nil
}

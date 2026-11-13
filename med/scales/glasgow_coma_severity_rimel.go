package scales

import "fmt"

// GCSSeverity is one of Rimel et al.'s 1982 head-injury groups by admission Glasgow Coma Scale total; its
// value is the source's own word, verbatim and lower-case.
type GCSSeverity string

// Rimel et al. 1982 groups, verbatim from the abstract: "severe, 3-8; moderate, 9-12; and minor, 13-15".
const (
	GCSSevere   GCSSeverity = "severe"   // 3–8
	GCSModerate GCSSeverity = "moderate" // 9–12
	GCSMinor    GCSSeverity = "minor"    // 13–15
)

// GlasgowComaSeverityRimel returns Rimel et al.'s 1982 head-injury group for a Glasgow Coma Scale total, with
// the source's own word: 3–8 "severe"; 9–12 "moderate"; 13–15 "minor". These are Rimel's head-injury groups;
// the official Glasgow Coma Scale site states no bands.
// Errors: ErrOutOfRange when total is outside 3–15.
//
// ID:      MED-SCALE-004
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Rimel RW, Giordani B, Barth JT, Jane JA. Moderate head injury: completing the clinical spectrum of brain trauma. Neurosurgery. 1982;11(3):344-351. doi:10.1227/00006123-198209000-00002. PMID 7133350. Abstract: "severe, 3-8; moderate, 9-12; and minor, 13-15"
// Tags:    scales, gcs, glasgow-coma-scale, head-injury, category
// Units:   total points 3–15; result the source's group word
// Range:   3 ≤ total ≤ 15
// Example: 8 → "severe"; 9 → "moderate"; 13 → "minor" (Rimel 1982 abstract)
func GlasgowComaSeverityRimel(total int) (GCSSeverity, error) {
	if total < 3 || total > 15 {
		return "", fmt.Errorf("%w: total = %d, want 3–15", ErrOutOfRange, total)
	}
	switch {
	case total <= 8:
		return GCSSevere, nil
	case total <= 12:
		return GCSModerate, nil
	}
	return GCSMinor, nil
}

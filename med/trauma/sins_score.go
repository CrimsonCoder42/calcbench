package trauma

import "fmt"

// SINSScore returns the Spinal Instability Neoplastic Score, 0–18, the sum of six item points as the caller
// grades them: location (rigid S2–S5 0, semirigid T3–T10 1, mobile C3–C6 or L2–L4 2, junctional occiput–C2,
// C7–T2, T11–L1 or L5–S1 3); pain (pain-free lesion 0, occasional pain but not mechanical 1, mechanical pain 3);
// bone lesion (blastic 0, mixed 1, lytic 2); radiographic alignment (normal 0, de novo deformity 2,
// subluxation/translation 4); vertebral body collapse (none 0, no collapse with > 50% body involved 1, < 50%
// collapse 2, > 50% collapse 3); posterolateral involvement (none 0, unilateral 1, bilateral 3). The sources band
// the total 0–6, 7–12 and 13–18 (stable, potentially unstable, unstable); the band is not returned.
// Errors: ErrOutOfRange for a point value that its item does not print.
//
// ID:      MED-TRAUMA-017
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Fisher CG, DiPaola CP, Ryken TC, et al. A novel classification system for spinal instability in neoplastic disease: an evidence-based approach and expert consensus from the Spine Oncology Study Group. Spine (Phila Pa 1976). 2010;35(22):E1221-E1229. doi:10.1097/BRS.0b013e3181e16ae2. PMID 20562730 (primary; not opened)
// Source:  Fisher CG, Schouten R, Versteeg AL, et al. Reliability of the Spinal Instability Neoplastic Score (SINS) among radiation oncologists. Radiat Oncol. 2014;9:69. doi:10.1186/1748-717X-9-69. PMID 24594004. PMC3995991, SINS table (cites Fisher 2010)
// Source:  Kim YR, Lee CH, Yang SH, et al. Accuracy and precision of the spinal instability neoplastic score (SINS) for predicting vertebral compression fractures after radiotherapy. Sci Rep. 2021;11(1):5553. doi:10.1038/s41598-021-84975-3. PMID 33692442. PMC7947012, Table 1 (cites Fisher 2010)
// Verified: secondary (Fisher 2014; Kim 2021 Table 1)
// Tags:    trauma, spine, oncology, sins, instability, score
// Units:   each argument the item's printed points; result points 0–18
// Range:   location 0–3; pain 0, 1 or 3; lesion 0–2; alignment 0, 2 or 4; collapse 0–3; posterolateral 0, 1 or 3
// Example: mobile 2 + mechanical pain 3 + lytic 2 + normal 0 + < 50% collapse 2 + unilateral 1 = 10 (HAND-DERIVED)
func SINSScore(location, pain, lesion, alignment, collapse, posterolateral int) (int, error) {
	type item struct {
		name  string
		v     int
		valid uint8 // bit k set when k points is printed
	}
	for _, it := range [...]item{
		{"location", location, 0b1111}, {"pain", pain, 0b1011}, {"lesion", lesion, 0b0111},
		{"alignment", alignment, 0b10101}, {"collapse", collapse, 0b1111}, {"posterolateral", posterolateral, 0b1011},
	} {
		if it.v < 0 || it.v > 4 || it.valid&(1<<it.v) == 0 {
			return 0, fmt.Errorf("%w: %s = %d, not a printed point value", ErrOutOfRange, it.name, it.v)
		}
	}
	return location + pain + lesion + alignment + collapse + posterolateral, nil
}

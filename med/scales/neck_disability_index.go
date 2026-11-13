package scales

import "fmt"

// NeckDisabilityIndex returns the Neck Disability Index total, 0–50: "10 sections, each with six possible
// answers representing six progressive levels of functional disability, 0 being the lowest level and five
// being the highest". Items are not embedded; the caller supplies each section's 0–5 rating.
// Errors: ErrOutOfRange for a section outside 0–5.
//
// ID:      MED-SCALE-011
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Vernon H, Mior S. The Neck Disability Index: a study of reliability and validity. J Manipulative Physiol Ther. 1991;14(7):409-415. PMID 1834753 (primary; not opened)
// Source:  Rodríguez-Sanz J, Malo-Urriés M, Lucha-López MO, et al. Comparison of an exercise program with and without manual therapy for patients with chronic neck pain. PeerJ. 2021;9:e12546. doi:10.7717/peerj.12546, PMID 34900443, PMC8627131, Methods
// Source:  Kasemsan A, Joseph L, Paungmali A, et al. Prevalence of musculoskeletal pain and associated disability among professional bus drivers. Int Arch Occup Environ Health. 2021;94(6):1263-1270. doi:10.1007/s00420-021-01683-1, PMID 33856539, PMC8047600, Methods
// Verified: secondary (Rodríguez-Sanz 2021 PMC8627131; Kasemsan 2021 PMC8047600; also Geoghegan 2023 PMC9977990)
// Tags:    scales, neck pain, disability, ndi, score
// Units:   sections points 0–5 each; result points 0–50 (× 2 gives the percentage)
// Range:   each section 0–5
// Example: sections 2,1,2,1,0,1,2,1,1,0 → 11 (hand-derived)
func NeckDisabilityIndex(sections [10]int) (int, error) {
	s := 0
	for i, v := range sections {
		if v < 0 || v > 5 {
			return 0, fmt.Errorf("%w: sections[%d] = %d, want 0–5", ErrOutOfRange, i, v)
		}
		s += v
	}
	return s, nil
}

// NDICategory is a Neck Disability Index disability band; values are the source's words.
type NDICategory string

const (
	NDINoDisability       NDICategory = "no disability"       // 0–4
	NDIMildDisability     NDICategory = "mild disability"     // 5–14
	NDIModerateDisability NDICategory = "moderate disability" // 15–24
	NDISevereDisability   NDICategory = "severe disability"   // 25–34
	NDICompleteDisability NDICategory = "complete disability" // ≥ 35
)

// NDIDisability returns the Neck Disability Index band of a total: "0–4, no disability; 5–14, mild
// disability; 15–24, moderate disability; 25–34, severe disability"; complete disability from 35 ("Above
// 34", Rodríguez-Sanz; "≥35", Geoghegan; Kasemsan's "greater than 35" would leave 35 unassigned).
// Errors: ErrOutOfRange for a total outside 0–50.
//
// ID:      MED-SCALE-012
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Vernon H, Mior S. J Manipulative Physiol Ther. 1991;14(7):409-415. PMID 1834753 (primary; not opened)
// Source:  Kasemsan A, Joseph L, Paungmali A, et al. Int Arch Occup Environ Health. 2021;94(6):1263-1270. doi:10.1007/s00420-021-01683-1, PMID 33856539, PMC8047600 (labels)
// Source:  Rodríguez-Sanz J, Malo-Urriés M, Lucha-López MO, et al. PeerJ. 2021;9:e12546. doi:10.7717/peerj.12546, PMID 34900443, PMC8627131 (bands, "Above 34")
// Verified: secondary (Kasemsan 2021 PMC8047600; Rodríguez-Sanz 2021 PMC8627131; Geoghegan 2023 PMC9977990 "≥35")
// Tags:    scales, neck pain, disability, ndi, category
// Units:   total NDI points; result the source's label
// Range:   total 0–50
// Example: 4 → no disability; 5 → mild disability; 35 → complete disability
func NDIDisability(total int) (NDICategory, error) {
	switch {
	case total < 0 || total > 50:
		return "", fmt.Errorf("%w: total = %d, want 0–50", ErrOutOfRange, total)
	case total <= 4:
		return NDINoDisability, nil
	case total <= 14:
		return NDIMildDisability, nil
	case total <= 24:
		return NDIModerateDisability, nil
	case total <= 34:
		return NDISevereDisability, nil
	}
	return NDICompleteDisability, nil
}

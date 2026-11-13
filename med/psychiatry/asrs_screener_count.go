package psychiatry

import "fmt"

// ASRSScreenerCount returns how many of the six ASRS v1.1 screener (Part A) items reach their threshold, 0–6.
// Each item is rated "0 (never), 1 (rarely), 2 (sometimes), 3 (often) or 4 (very often)"; items 1–3 count at
// 2 or more and items 4–6 at 3 or more. The cited screen is "four or more"; applying it is the caller's. Items
// not embedded; licence: the ASRS is the World Health Organization's.
// Errors: ErrOutOfRange when an item is outside 0–4, naming its index.
//
// ID:      MED-PSYCH-013
// Time:    O(1)   6 items, fixed
// Space:   O(1), no allocation
// Source:  Kessler RC, Adler L, Ames M, et al. The World Health Organization Adult ADHD Self-Report Scale (ASRS): a short screening scale for use in the general population. Psychol Med. 2005;35(2):245-256. doi:10.1017/s0033291704002892, PMID 15841682 (full text not open)
// Source:  Chamberlain SR, Cortese S, Grant JE. Compr Psychiatry. 2021;106:152224. doi:10.1016/j.comppsych.2021.152224, PMC7116749 ("endorsing sometimes/often/very often for questions 1–3; and endorsing often/very often for questions 4–6"; "a threshold of four or more")
// Source:  Kittel-Schneider S, et al. J Clin Med. 2019;8(10):1643. doi:10.3390/jcm8101643, PMC6832520 ("four or more questions answered with ≥2 (questions 1–3) and/or ≥3 (questions 4–6)")
// Verified: secondary (Chamberlain 2021 PMC7116749; Kittel-Schneider 2019 PMC6832520)
// Tags:    psychiatry, adhd, asrs, screener, questionnaire, count
// Units:   item ratings 0–4; result count 0–6
// Range:   each item 0–4
// Example: ratings 2, 2, 2, 3, 2, 2 → 4 (items 5 and 6 below 3); all 1 → 0
func ASRSScreenerCount(items [6]int) (int, error) {
	n := 0
	for i, v := range items {
		if v < 0 || v > 4 {
			return 0, fmt.Errorf("%w: items[%d] = %d, want 0–4", ErrOutOfRange, i, v)
		}
		threshold := 3
		if i < 3 {
			threshold = 2
		}
		if v >= threshold {
			n++
		}
	}
	return n, nil
}

package scales

// PAINADScore returns the Pain Assessment in Advanced Dementia (PAINAD) score, 0–10: the sum of the five items
// breathing (independent of vocalization), negative vocalization, facial expression, body language and
// consolability, each rated 0, 1 or 2. Items not embedded: the rating descriptors are the instrument's text
// (Warden 2003); only the arithmetic is here. No bands are returned: the only published bands found overlap
// at 7 as printed and are not used.
// Errors: ErrOutOfRange when an item is outside 0–2, naming the item.
//
// ID:      MED-SCALE-015
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Warden V, Hurley AC, Volicer L. Development and psychometric evaluation of the Pain Assessment in Advanced Dementia (PAINAD) scale. J Am Med Dir Assoc. 2003;4(1):9-15. doi:10.1097/01.JAM.0000043422.31640.F7. PMID 12807591, abstract ("a five-item observational tool with a range of 0 to 10"); per-item 0–2 not in the abstract
// Source:  Shaw C, et al. Innov Aging. 2023;7(10):igad076. PMC10714906, Methods ("a five-item scale in which each item is scored from 0 to 2 for breathing, vocalization, facial expression, body language, and consolability"; "summed for a pain severity rating of 0–10"); cites Warden 2003
// Source:  Pinto MCM, et al. Einstein (São Paulo). 2015;13(1):14-19. PMC4977604, Chart 1 (PAINAD-Brazil: the five items each 0, 1, 2; 0–10); cites Warden 2003
// Verified: secondary (Shaw 2023 PMC10714906; Pinto 2015 PMC4977604)
// Tags:    scales, pain, painad, dementia, observational, score
// Units:   each item points 0–2; result points 0–10
// Range:   breathing, negativeVocalization, facialExpression, bodyLanguage, consolability each 0–2
// Example: breathing 1, vocalization 1, facial 2, body language 1, consolability 0 → 5; all 2 → 10 (HAND-DERIVED from the sum)
func PAINADScore(breathing, negativeVocalization, facialExpression, bodyLanguage, consolability int) (int, error) {
	for _, e := range [...]error{
		itemZeroToTwo("breathing", breathing), itemZeroToTwo("negativeVocalization", negativeVocalization),
		itemZeroToTwo("facialExpression", facialExpression), itemZeroToTwo("bodyLanguage", bodyLanguage),
		itemZeroToTwo("consolability", consolability),
	} {
		if e != nil {
			return 0, e
		}
	}
	return breathing + negativeVocalization + facialExpression + bodyLanguage + consolability, nil
}

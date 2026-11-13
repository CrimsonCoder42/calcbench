package neuro

// ASPECTSScore returns the Alberta Stroke Program Early CT Score, 0–10: 10 minus one point for each of the ten
// middle-cerebral-artery regions (caudate, lentiform, internal capsule, insular ribbon, M1–M6) with early
// ischemic change. The order of the ten flags does not matter.
//
// ID:      MED-NEURO-024
// Time:    O(1)   10 regions, fixed
// Space:   O(1), no allocation
// Source:  Barber PA, Demchuk AM, Zhang J, Buchan AM. Validity and reliability of a quantitative computed tomography score in predicting outcome of hyperacute stroke before thrombolytic therapy. ASPECTS Study Group. Lancet. 2000;355(9216):1670-1674. doi:10.1016/s0140-6736(00)02237-6, PMID 10905241 (abstract: "ten regions of interest"; the subtraction rule not in the abstract)
// Source:  Teekaput C, et al. J Multidiscip Healthc. 2024;17:5487-5499. doi:10.2147/JMDH.S495952, PMC11590667, Methods ("One point was subtracted from a total of ten for each observed sign of EIC within the defined regions")
// Source:  Nukovic JJ, et al. Medicina (Kaunas). 2023;59(11):1908. doi:10.3390/medicina59111908, PMC10673396, Methods ("starting from 10 points, one point is subtracted for each marked area affected by ischemic changes")
// Verified: secondary (Teekaput 2024 PMC11590667; Nukovic 2023 PMC10673396)
// Tags:    neuro, stroke, ct, aspects, score
// Units:   region flags yes/no (early ischemic change present); result points 0–10
// Range:   all inputs boolean
// Example: no region affected → 10; caudate and M1 affected → 8 (10 − 2)
func ASPECTSScore(regionAffected [10]bool) int {
	s := 10
	for _, hit := range regionAffected {
		if hit {
			s--
		}
	}
	return s
}

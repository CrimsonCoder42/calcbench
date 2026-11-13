package neuro

// AcuteStrokeDysphagiaScreenItems are the five items of the Acute Stroke Dysphagia Screen (the Barnes-Jewish
// Hospital stroke dysphagia screen; Edmiaston 2010). Each is true when the abnormal finding is present; the
// thresholds (the Glasgow Coma Scale level, the 3-oz water swallow signs) are in the source's figure and are the
// caller's.
type AcuteStrokeDysphagiaScreenItems struct {
	ReducedConsciousness bool // level of consciousness by the Glasgow Coma Scale
	FacialAsymmetry      bool // "Facial ... symmetry" (dysarthria item)
	LingualAsymmetry     bool // "lingual ... symmetry" (dysarthria item)
	PalatalAsymmetry     bool // "palatal symmetry" (dysarthria item)
	WaterSwallowFailed   bool // "a 3-oz water swallow test"
}

// AcuteStrokeDysphagiaScreenFailed reports whether a person fails the Acute Stroke Dysphagia Screen: each item
// is scored present or absent, and any item present fails the screening (Edmiaston 2010, Methods).
//
// ID:      MED-NEURO-032
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Edmiaston J, Connor LT, Loehr L, Nassief A. Validation of a dysphagia screening tool in acute stroke patients. Am J Crit Care. 2010;19(4):357-364. doi:10.4037/ajcc2009961, PMID 19875722, PMC2896456, Methods ("a measure of level of consciousness, several items that together indicated whether a patient had dysarthria, and a 3-oz water swallow test"; "Facial, lingual, and palatal symmetry were all included"; each item "scored as either present or absent", and any item marked present fails the screening)
// Tags:    neuro, stroke, dysphagia, screening, decision rule
// Units:   yes/no items; result true when any item is present
// Range:   all inputs boolean (acute stroke, the source's population; the caller's)
// Example: water swallow failed alone → true; nothing present → false
func AcuteStrokeDysphagiaScreenFailed(it AcuteStrokeDysphagiaScreenItems) bool {
	return it.ReducedConsciousness || it.FacialAsymmetry || it.LingualAsymmetry || it.PalatalAsymmetry ||
		it.WaterSwallowFailed
}

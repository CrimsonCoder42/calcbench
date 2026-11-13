package neuro

// HINTSAnyFinding reports whether any of the three findings of the HINTS bedside oculomotor examination is
// present: a normal horizontal head impulse test, direction-changing nystagmus in eccentric gaze, or skew
// deviation. The source reports the presence of any of the three as "100% sensitive and 96% specific for
// stroke" in acute vestibular syndrome with at least one stroke risk factor; applying that is the caller's.
//
// ID:      MED-NEURO-023
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Kattah JC, Talkad AV, Wang DZ, Hsieh YH, Newman-Toker DE. HINTS to diagnose stroke in the acute vestibular syndrome: three-step bedside oculomotor examination more sensitive than early MRI diffusion-weighted imaging. Stroke. 2009;40(11):3504-3510. doi:10.1161/STROKEAHA.109.551234, PMID 19762709, abstract ("normal horizontal head impulse test, direction-changing nystagmus in eccentric gaze, or skew deviation (vertical ocular misalignment) was 100% sensitive and 96% specific for stroke")
// Tags:    neuro, vertigo, stroke, hints, oculomotor, decision rule
// Units:   yes/no findings; result true when at least one is present
// Range:   acute vestibular syndrome with at least one stroke risk factor (the source's population; the caller's)
// Example: abnormal head impulse, direction-fixed nystagmus, no skew → false; skew alone → true (HAND-DERIVED)
func HINTSAnyFinding(normalHeadImpulse, directionChangingNystagmus, skewDeviation bool) bool {
	return normalHeadImpulse || directionChangingNystagmus || skewDeviation
}

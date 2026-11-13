package neuro

// OttawaSAHFindings are the six findings of the Ottawa SAH Rule (Perry 2013 abstract).
type OttawaSAHFindings struct {
	Age40OrOlder                 bool // "age 40 years or older"
	NeckPainOrStiffness          bool // "neck pain or stiffness"
	WitnessedLossOfConsciousness bool // "witnessed loss of consciousness"
	OnsetDuringExertion          bool // "onset during exertion"
	ThunderclapHeadache          bool // "thunderclap headache"
	LimitedNeckFlexion           bool // "limited neck flexion on examination"
}

// OttawaSAHRuleCount returns how many of the six Ottawa SAH Rule findings are present, 0–6. The source reports
// the rule (any finding present) as "100% sensitivity and 15.3% specificity" for subarachnoid hemorrhage in
// alert adults with a headache peaking within 1 hour and no neurologic deficit; applying it is the caller's.
//
// ID:      MED-NEURO-020
// Time:    O(1)   6 findings, fixed
// Space:   O(1), no allocation
// Source:  Perry JJ, Stiell IG, Sivilotti ML, et al. Clinical decision rules to rule out subarachnoid hemorrhage for acute headache. JAMA. 2013;310(12):1248-1255. doi:10.1001/jama.2013.278018, PMID 24065011, abstract ("age 40 years or older, neck pain or stiffness, witnessed loss of consciousness, or onset during exertion" with "thunderclap headache" and "limited neck flexion on examination"; "100% sensitivity and 15.3% specificity")
// Tags:    neuro, headache, subarachnoid hemorrhage, ottawa sah rule, decision rule, count
// Units:   yes/no findings; result count 0–6
// Range:   adults with headache peaking within 1 hour and no neurologic deficit (the source's population; the caller's)
// Example: age 40 or older alone → 1; none → 0 (the count)
func OttawaSAHRuleCount(f OttawaSAHFindings) int {
	n := 0
	for _, yes := range [...]bool{f.Age40OrOlder, f.NeckPainOrStiffness, f.WitnessedLossOfConsciousness,
		f.OnsetDuringExertion, f.ThunderclapHeadache, f.LimitedNeckFlexion} {
		if yes {
			n++
		}
	}
	return n
}

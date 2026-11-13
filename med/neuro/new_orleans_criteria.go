package neuro

// NewOrleansFindings are the seven findings of the New Orleans criteria (Haydel 2000 abstract).
type NewOrleansFindings struct {
	Headache                  bool // "headache"
	Vomiting                  bool // "vomiting"
	AgeOver60                 bool // "an age over 60 years"
	DrugOrAlcoholIntoxication bool // "drug or alcohol intoxication"
	ShortTermMemoryDeficit    bool // "deficits in short-term memory"
	TraumaAboveClavicles      bool // "physical evidence of trauma above the clavicles"
	Seizure                   bool // "seizure"
}

// NewOrleansCriteriaCount returns how many of the seven New Orleans findings are present, 0–7. The source
// reports the seven findings combined (any one present) as "100 percent" sensitive for a positive CT scan in
// minor head injury with a normal Glasgow Coma Scale score; applying that is the caller's.
//
// ID:      MED-NEURO-018
// Time:    O(1)   7 findings, fixed
// Space:   O(1), no allocation
// Source:  Haydel MJ, Preston CA, Mills TJ, Luber S, Blaudeau E, DeBlieux PM. Indications for computed tomography in patients with minor head injury. N Engl J Med. 2000;343(2):100-105. doi:10.1056/NEJM200007133430204, PMID 10891517, abstract ("seven findings: headache, vomiting, an age over 60 years, drug or alcohol intoxication, deficits in short-term memory, physical evidence of trauma above the clavicles, and seizure"; "sensitivity of the seven findings combined was 100 percent")
// Tags:    neuro, head injury, ct, new orleans criteria, decision rule, count
// Units:   yes/no findings; result count 0–7
// Range:   minor head injury with a normal GCS score (the source's population; the caller's)
// Example: headache and vomiting → 2; none → 0 (the count)
func NewOrleansCriteriaCount(f NewOrleansFindings) int {
	n := 0
	for _, yes := range [...]bool{f.Headache, f.Vomiting, f.AgeOver60, f.DrugOrAlcoholIntoxication,
		f.ShortTermMemoryDeficit, f.TraumaAboveClavicles, f.Seizure} {
		if yes {
			n++
		}
	}
	return n
}

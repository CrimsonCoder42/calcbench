package neuro

import "fmt"

// SAH thickness on CT for ModifiedFisherGrade.
const (
	SAHNone  = 0 // no subarachnoid hemorrhage
	SAHThin  = 1 // focal or diffuse, thin SAH
	SAHThick = 2 // focal or diffuse, thick SAH
)

// ModifiedFisherGrade returns the modified Fisher grade, 0–4, from the subarachnoid blood on CT and the presence
// of intraventricular hemorrhage (IVH): 0 no SAH and no IVH; 1 thin SAH, no IVH; 2 thin SAH with IVH; 3 thick SAH,
// no IVH; 4 thick SAH with IVH. The source grades no case of IVH without SAH.
// Errors: ErrOutOfRange when sahThickness is not SAHNone, SAHThin or SAHThick, or for IVH without SAH.
//
// ID:      MED-NEURO-028
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Frontera JA, Claassen J, Schmidt JM, et al. Prediction of symptomatic vasospasm after subarachnoid hemorrhage: the modified Fisher scale. Neurosurgery. 2006;59(1):21-27. doi:10.1227/01.neu.0000243277.86222.6c, PMID 16823296 (full text not open)
// Source:  Patel S, Parikh A, Okorie ON. Subarachnoid hemorrhage in the emergency department. Int J Emerg Med. 2021;14:31. doi:10.1186/s12245-021-00353-w, PMC8117305, Table 6 (0 "No SAH; no IVH" … 4 "Thick, focal, or diffuse SAH; IVH")
// Source:  Zhou Z, Yu J. Stroke Res Treat. 2022;2022:5774735. doi:10.1155/2022/5774735, PMC9492435, §2.3 ("Grade 0, no SAH or intraventricular hemorrhage (IVH)" … "Grade 4, focal or diffuse, thick SAH, and with IVH")
// Verified: secondary (Patel 2021 PMC8117305; Zhou 2022 PMC9492435)
// Tags:    neuro, subarachnoid hemorrhage, ct, modified fisher, grade
// Units:   sahThickness code (SAHNone, SAHThin, SAHThick); ivh yes/no; result grade 0–4
// Range:   sahThickness 0–2; IVH only with SAH present
// Example: thick SAH with IVH → 4; thin SAH without IVH → 1; neither → 0
func ModifiedFisherGrade(sahThickness int, ivh bool) (int, error) {
	if err := checkInt("sahThickness", sahThickness, SAHNone, SAHThick); err != nil {
		return 0, err
	}
	switch sahThickness {
	case SAHNone:
		if ivh {
			return 0, fmt.Errorf("%w: IVH without SAH is not graded by the source", ErrOutOfRange)
		}
		return 0, nil
	case SAHThin:
		if ivh {
			return 2, nil
		}
		return 1, nil
	}
	if ivh {
		return 4, nil
	}
	return 3, nil
}

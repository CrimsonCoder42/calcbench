package hepatic

import (
	"fmt"

	"github.com/CrimsonCoder42/calcbench/algo/reduce"
)

// MaddreyDiscriminantFunction returns Maddrey's discriminant function for alcoholic hepatitis, dimensionless:
// 4.6 × (prothrombin time − control prothrombin time, seconds) + total bilirubin (mg/dL). A prothrombin time
// shorter than the control gives a smaller, possibly negative, result, returned as computed; no band is
// returned.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange outside Range; ErrDomain when the result is beyond float64.
//
// ID:      MED-HEPAT-007
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Maddrey WC, Boitnott JK, Bedine MS, Weber FL Jr, Mezey E, White RI Jr. Corticosteroid therapy of alcoholic hepatitis. Gastroenterology. 1978;75(2):193-199. PMID 352788 (full text not opened; abstract states no formula)
// Source:  Mitri J, Almeqdadi M, Karagozian R. Prognostic and diagnostic scoring models in acute alcohol-associated hepatitis. World J Hepatol. 2023. PMID 37701919. PMC10494564: "Discriminant Function = 4.6 × (Pt's PT-control PT) + TBili." [ref 6, Maddrey 1978]
// Source:  Dugum M, McCullough A. Diagnosis and management of alcoholic liver disease. J Clin Transl Hepatol. 2015. PMID 26356792. PMC4548353: "4.6 × [patient's prothrombin time (seconds) – control prothrombin time (seconds)] + bilirubin (mg/dL)" [Maddrey 1978]
// Verified: secondary (Mitri 2023 PMC10494564; Dugum 2015 PMC4548353)
// Tags:    hepatic, liver, alcoholic-hepatitis, maddrey, discriminant-function, prothrombin, bilirubin
// Units:   ptSeconds s; controlPTSeconds s; bilirubinMgDL mg/dL; result dimensionless
// Range:   every input > 0 (input limits not stated in the sources)
// Example: PT 20 s, control 12 s, bilirubin 10 mg/dL → 4.6 × 8 + 10 = 46.8 (hand-derived)
func MaddreyDiscriminantFunction(ptSeconds, controlPTSeconds, bilirubinMgDL float64) (float64, error) {
	if err := checkFinite("ptSeconds", ptSeconds); err != nil {
		return 0, err
	}
	if err := checkFinite("controlPTSeconds", controlPTSeconds); err != nil {
		return 0, err
	}
	if err := checkFinite("bilirubinMgDL", bilirubinMgDL); err != nil {
		return 0, err
	}
	if err := checkPositive("ptSeconds", ptSeconds); err != nil {
		return 0, err
	}
	if err := checkPositive("controlPTSeconds", controlPTSeconds); err != nil {
		return 0, err
	}
	if err := checkPositive("bilirubinMgDL", bilirubinMgDL); err != nil {
		return 0, err
	}
	r := 4.6*(ptSeconds-controlPTSeconds) + bilirubinMgDL
	if !reduce.AllFinite(r) {
		return 0, fmt.Errorf("%w: result is beyond float64", ErrDomain)
	}
	return r, nil
}

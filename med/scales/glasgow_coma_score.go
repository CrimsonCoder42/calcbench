package scales

import "fmt"

// GCSNotTestable is the rating passed for a Glasgow Coma Scale component the assessor rated "Non testable"
// (NT) (Teasdale 2015, GCS Assessment Aid). With any component NT, GlasgowComaScore returns ErrNotTestable.
const GCSNotTestable = 0

// Glasgow Coma Scale ratings, with the rating words of Teasdale G, Glasgow Coma Scale: Assessment Aid
// (English), 2015, https://www.glasgowcomascale.org/downloads/GCS-Assessment-Aid-English.pdf (read 2026-11-12).
const (
	GCSEyeSpontaneous = 4 // eye opening: "Spontaneous"
	GCSEyeToSound     = 3 // eye opening: "To sound"
	GCSEyeToPressure  = 2 // eye opening: "To pressure"
	GCSEyeNone        = 1 // eye opening: "None"

	GCSVerbalOrientated = 5 // verbal response: "Orientated"
	GCSVerbalConfused   = 4 // verbal response: "Confused"
	GCSVerbalWords      = 3 // verbal response: "Words"
	GCSVerbalSounds     = 2 // verbal response: "Sounds"
	GCSVerbalNone       = 1 // verbal response: "None"

	GCSMotorObeysCommands   = 6 // best motor response: "Obeys commands"
	GCSMotorLocalising      = 5 // best motor response: "Localising"
	GCSMotorNormalFlexion   = 4 // best motor response: "Normal flexion" (the source's label)
	GCSMotorAbnormalFlexion = 3 // best motor response: "Abnormal flexion" (the source's label)
	GCSMotorExtension       = 2 // best motor response: "Extension"
	GCSMotorNone            = 1 // best motor response: "None"
)

// checkGCSComponent returns ErrOutOfRange unless v is GCSNotTestable or 1 ≤ v ≤ maxRating, and reports
// whether v is GCSNotTestable.
func checkGCSComponent(name string, v, maxRating int) (notTestable bool, err error) {
	if v == GCSNotTestable {
		return true, nil
	}
	if v < 1 || v > maxRating {
		return false, fmt.Errorf("%w: %s = %d, want 1–%d or GCSNotTestable (0)", ErrOutOfRange, name, v, maxRating)
	}
	return false, nil
}

// GlasgowComaScore returns the Glasgow Coma Scale total, 3–15: eye (1–4) + verbal (1–5) + motor (1–6). A
// component rated "Non testable" is passed as GCSNotTestable (0); then no total exists and the function
// returns ErrNotTestable. Teasdale et al. 2014 (abstract): "Individual patients are best described by the
// three components of the coma scale; whereas the derived total coma score should be used to characterise
// groups."
// Errors: ErrOutOfRange when a component is outside its range and not GCSNotTestable (checked for all three
// first); ErrNotTestable when any component is GCSNotTestable.
//
// ID:      MED-SCALE-003
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Teasdale G. Glasgow Coma Scale: Assessment Aid (English). 2015. https://www.glasgowcomascale.org/downloads/GCS-Assessment-Aid-English.pdf (read 2026-11-12): eye 4 Spontaneous, 3 To sound, 2 To pressure, 1 None; verbal 5 Orientated, 4 Confused, 3 Words, 2 Sounds, 1 None; motor 6 Obeys commands, 5 Localising, 4 Normal flexion, 3 Abnormal flexion, 2 Extension, 1 None; NT Non testable.
// Source:  Teasdale G, Jennett B. Assessment of coma and impaired consciousness. A practical scale. Lancet. 1974;2(7872):81-84. doi:10.1016/s0140-6736(74)91639-0. PMID 4136544.
// Source:  Teasdale G, Maas A, Lecky F, Manley G, Stocchetti N, Murray G. The Glasgow Coma Scale at 40 years: standing the test of time. Lancet Neurol. 2014;13(8):844-854. doi:10.1016/S1474-4422(14)70120-6. PMID 25030516.
// Tags:    scales, gcs, glasgow-coma-scale, consciousness, neurology, score
// Units:   eye points 1–4; verbal points 1–5; motor points 1–6 (0 = GCSNotTestable); result points 3–15
// Range:   eye 1–4, verbal 1–5, motor 1–6, each or GCSNotTestable (0)
// Example: eye 3 (To sound), verbal 4 (Confused), motor 5 (Localising) → 12; 4 + 5 + 6 → 15; 1 + 1 + 1 → 3 (HAND-DERIVED from the Assessment Aid)
func GlasgowComaScore(eye, verbal, motor int) (int, error) {
	eyeNT, err := checkGCSComponent("eye", eye, 4)
	if err != nil {
		return 0, err
	}
	verbalNT, err := checkGCSComponent("verbal", verbal, 5)
	if err != nil {
		return 0, err
	}
	motorNT, err := checkGCSComponent("motor", motor, 6)
	if err != nil {
		return 0, err
	}
	if eyeNT || verbalNT || motorNT {
		return 0, fmt.Errorf("%w: eye = %d, verbal = %d, motor = %d", ErrNotTestable, eye, verbal, motor)
	}
	return eye + verbal + motor, nil
}

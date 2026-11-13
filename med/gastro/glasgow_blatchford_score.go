package gastro

import "fmt"

// GlasgowBlatchfordFindings holds the Glasgow-Blatchford admission risk markers (Blatchford 2000, as reproduced
// in Banister 2018 Table 1 and Orpen-Palmer 2023 Table 1). Laboratory values are in the units of the table.
type GlasgowBlatchfordFindings struct {
	UreaMmolL      float64 // blood urea, mmol/L: 6.5–7.9: 2, 8.0–9.9: 3, 10.0–24.9: 4, ≥ 25: 6
	HaemoglobinGL  float64 // haemoglobin, g/L; men 120–129: 1, 100–119: 3, < 100: 6; women 100–119: 1, < 100: 6
	SystolicMmHg   float64 // systolic blood pressure, mmHg: 100–109: 1, 90–99: 2, < 90: 3
	PulseBpm       float64 // pulse, beats/min: ≥ 100: 1
	Male           bool    // selects the haemoglobin row for men
	Melaena        bool    // "Presentation with melaena": 1
	Syncope        bool    // "Presentation with syncope": 2
	HepaticDisease bool    // "Hepatic disease": 2
	CardiacFailure bool    // "Cardiac failure": 2
}

// GlasgowBlatchfordScore returns the Glasgow-Blatchford score, 0–23, for upper gastrointestinal haemorrhage.
// Band edges are read as half-open intervals (a value on an edge scores with the band that starts there:
// urea 25.0 → 6, haemoglobin 130 g/L in men → 0, 120 g/L in women → 0); the secondaries print the edges
// differently (see Verified). No mg/dL form is given here because neither source states one.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a negative measurement.
//
// ID:      MED-GASTRO-010
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Blatchford O, Murray WR, Blatchford M. A risk score to predict need for treatment for upper-gastrointestinal haemorrhage. Lancet. 2000;356(9238):1318-1321. doi:10.1016/S0140-6736(00)02816-6, PMID 11073021 (primary; not opened)
// Source:  Banister T, Spiking J, Ayaru L. Discharge of patients with an acute upper gastrointestinal bleed from the emergency department using an extended Glasgow-Blatchford Score. BMJ Open Gastroenterol. 2018;5(1):e000225. doi:10.1136/bmjgast-2018-000225, PMID 30233807, PMC6135483, Table 1
// Source:  Orpen-Palmer J, Stanley AJ. A Review of Risk Scores within Upper Gastrointestinal Bleeding. J Clin Med. 2023;12(11):3678. PMID 37297873, PMC10253886, Table 1 "The Glasgow Blatchford Score"
// Verified: secondary (Banister 2018 PMC6135483 Table 1; Orpen-Palmer 2023 PMC10253886 Table 1; edges also as Adrych 2025 PMC12746585 Table I: urea "≥25" 6, Hb men "120–129")
// Tags:    gastro, gi-bleeding, upper-gi, glasgow-blatchford, gbs, score
// Units:   UreaMmolL mmol/L; HaemoglobinGL g/L; SystolicMmHg mmHg; PulseBpm beats/min; result points 0–23
// Range:   any non-negative measurement
// Example: man, urea 9.0, Hb 115, SBP 95, pulse 105, melaena → 3 + 3 + 2 + 1 + 1 = 10 (hand-derived from Table 1)
func GlasgowBlatchfordScore(f GlasgowBlatchfordFindings) (int, error) {
	for _, in := range [...]struct {
		name string
		v    float64
	}{{"UreaMmolL", f.UreaMmolL}, {"HaemoglobinGL", f.HaemoglobinGL}, {"SystolicMmHg", f.SystolicMmHg}, {"PulseBpm", f.PulseBpm}} {
		if err := checkFinite(in.name, in.v); err != nil {
			return 0, err
		}
		if in.v < 0 {
			return 0, fmt.Errorf("%w: %s = %v, want ≥ 0", ErrOutOfRange, in.name, in.v)
		}
	}
	s := 0
	switch u := f.UreaMmolL; {
	case u >= 25:
		s += 6
	case u >= 10:
		s += 4
	case u >= 8:
		s += 3
	case u >= 6.5:
		s += 2
	}
	switch hb := f.HaemoglobinGL; {
	case hb < 100:
		s += 6
	case f.Male && hb < 120:
		s += 3
	case f.Male && hb < 130:
		s++
	case !f.Male && hb < 120:
		s++
	}
	switch sbp := f.SystolicMmHg; {
	case sbp < 90:
		s += 3
	case sbp < 100:
		s += 2
	case sbp < 110:
		s++
	}
	if f.PulseBpm >= 100 {
		s++
	}
	if f.Melaena {
		s++
	}
	if f.Syncope {
		s += 2
	}
	if f.HepaticDisease {
		s += 2
	}
	if f.CardiacFailure {
		s += 2
	}
	return s, nil
}

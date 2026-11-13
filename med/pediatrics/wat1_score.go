package pediatrics

// WAT1Findings holds the eleven Withdrawal Assessment Tool-1 items (Franck 2008). The first three are taken
// from the record of the previous 12 h; the next five from a 2-minute pre-stimulus observation; two from a
// 1-minute stimulus observation; the last from post-stimulus recovery.
type WAT1Findings struct {
	LooseWateryStools       bool    // 0/1
	VomitingRetchingGagging bool    // 0/1
	TemperatureAbove378C    bool    // temperature > 37.8 °C: 0/1
	StateAwakeDistressed    bool    // behavioural state item: 0/1
	TremorModerateSevere    bool    // 0/1
	Sweating                bool    // 0/1
	UncoordinatedMovement   bool    // uncoordinated/repetitive movement, moderate/severe: 0/1
	YawningOrSneezing2Plus  bool    // yawning or sneezing ≥ 2: 0/1
	StartleToTouch          bool    // moderate/severe: 0/1
	IncreasedMuscleTone     bool    // 0/1
	MinutesToCalmState      float64 // time to gain calm state: < 2 min 0; 2–5 min 1; > 5 min 2
}

// WAT1Score returns the Withdrawal Assessment Tool-1 total, 0–12 ("11 items", "12-point scale"). The source
// reports that "A WAT-1 score of 3 or higher had the best sensitivity and specificity"; no category label is
// built.
// Errors: ErrNotFinite for NaN or ±Inf; ErrOutOfRange for a negative time.
//
// ID:      MED-PEDS-005
// Time:    O(1)
// Space:   O(1), no allocation
// Source:  Franck LS, Harris SK, Soetenga DJ, Amling JK, Curley MAQ. The Withdrawal Assessment Tool-1 (WAT-1): an assessment instrument for monitoring opioid and benzodiazepine withdrawal symptoms in pediatric patients. Pediatr Crit Care Med. 2008;9(6):573-580. doi:10.1097/PCC.0b013e31818c8328, PMID 18838937, PMC2775493 (author manuscript), WAT-1 form and Results
// Tags:    pediatrics, critical care, withdrawal, opioid, benzodiazepine, wat-1, score
// Units:   every item bool; MinutesToCalmState minutes; result points 0–12
// Range:   critically ill children being weaned from opioids and benzodiazepines (the source population); MinutesToCalmState ≥ 0
// Example: watery stools, tremor, sweating, startle, calm in 3 min → 1+1+1+1+1 = 5 (hand-derived from the form)
func WAT1Score(f WAT1Findings) (int, error) {
	if err := checkFloat("MinutesToCalmState", f.MinutesToCalmState, 0, 1e6); err != nil {
		return 0, err
	}
	s := 0
	for _, one := range [...]bool{
		f.LooseWateryStools, f.VomitingRetchingGagging, f.TemperatureAbove378C, f.StateAwakeDistressed, f.TremorModerateSevere,
		f.Sweating, f.UncoordinatedMovement, f.YawningOrSneezing2Plus, f.StartleToTouch, f.IncreasedMuscleTone,
	} {
		if one {
			s++
		}
	}
	switch m := f.MinutesToCalmState; {
	case m > 5:
		s += 2
	case m >= 2:
		s++
	}
	return s, nil
}

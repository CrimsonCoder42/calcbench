package critical

import "testing"

var sinkPhoenix PhoenixScores

// phoenixBase is a 30-month-old with nothing measured: 0 points.
var phoenixBase = PhoenixInput{AgeMonths: 30}

// Schlapbach 2024 Phoenix Sepsis Score table: both sides of every edge. HAND-DERIVED.
func TestPhoenixSepsisScore_table(t *testing.T) {
	r := func(s PhoenixScores) int { return s.Respiratory }
	c := func(s PhoenixScores) int { return s.Cardiovascular }
	g := func(s PhoenixScores) int { return s.Coagulation }
	n := func(s PhoenixScores) int { return s.Neurological }
	pf := func(v float64, sup, imv bool) func(*PhoenixInput) {
		return func(i *PhoenixInput) {
			i.PaO2FiO2Measured, i.PaO2FiO2, i.AnyRespiratorySupport, i.IMV = true, v, sup, imv
		}
	}
	sf := func(v float64, sup, imv bool) func(*PhoenixInput) {
		return func(i *PhoenixInput) {
			i.SpO2FiO2Measured, i.SpO2FiO2, i.AnyRespiratorySupport, i.IMV = true, v, sup, imv
		}
	}
	mapAt := func(ageMonths, m float64) func(*PhoenixInput) {
		return func(i *PhoenixInput) { i.AgeMonths, i.MAPMeasured, i.MAPMmHg = ageMonths, true, m }
	}
	tests := []struct {
		name string
		mod  func(*PhoenixInput)
		pick func(PhoenixScores) int
		want int
	}{
		{"nothing", func(*PhoenixInput) {}, r, 0},
		{"pf 399 no support", pf(399, false, false), r, 0},
		{"pf 400 support", pf(400, true, false), r, 0},
		{"pf 399 support", pf(399, true, false), r, 1},
		{"pf 150 support no imv", pf(150, true, false), r, 1},
		{"pf 201 imv", pf(201, false, true), r, 1},
		{"pf 200 imv", pf(200, false, true), r, 2},
		{"pf 100 imv", pf(100, false, true), r, 2},
		{"pf 99 imv", pf(99, false, true), r, 3},
		{"sf 292 support", sf(292, true, false), r, 0},
		{"sf 291 support", sf(291, true, false), r, 1},
		{"sf 221 imv", sf(221, false, true), r, 1},
		{"sf 220 imv", sf(220, false, true), r, 2},
		{"sf 148 imv", sf(148, false, true), r, 2},
		{"sf 147 imv", sf(147, false, true), r, 3},
		{"vaso 1", func(i *PhoenixInput) { i.VasoactiveMedications = 1 }, c, 1},
		{"vaso 2", func(i *PhoenixInput) { i.VasoactiveMedications = 2 }, c, 2},
		{"vaso 3", func(i *PhoenixInput) { i.VasoactiveMedications = 3 }, c, 2},
		{"lactate 4.9", func(i *PhoenixInput) { i.LactateMeasured, i.LactateMmolL = true, 4.9 }, c, 0},
		{"lactate 5", func(i *PhoenixInput) { i.LactateMeasured, i.LactateMmolL = true, 5 }, c, 1},
		{"lactate 10.9", func(i *PhoenixInput) { i.LactateMeasured, i.LactateMmolL = true, 10.9 }, c, 1},
		{"lactate 11", func(i *PhoenixInput) { i.LactateMeasured, i.LactateMmolL = true, 11 }, c, 2},
		{"<1 mo map 31", mapAt(0.5, 31), c, 0},
		{"<1 mo map 30", mapAt(0.5, 30), c, 1},
		{"<1 mo map 17", mapAt(0.5, 17), c, 1},
		{"<1 mo map 16", mapAt(0.5, 16), c, 2},
		{"1-11 mo map 39", mapAt(1, 39), c, 0},
		{"1-11 mo map 38", mapAt(1, 38), c, 1},
		{"1-11 mo map 24", mapAt(11.9, 24), c, 2},
		{"1-<2 y map 44", mapAt(12, 44), c, 0},
		{"1-<2 y map 43", mapAt(12, 43), c, 1},
		{"1-<2 y map 30", mapAt(23.9, 30), c, 2},
		{"2-<5 y map 45", mapAt(24, 45), c, 0},
		{"2-<5 y map 44", mapAt(24, 44), c, 1},
		{"2-<5 y map 31", mapAt(59.9, 31), c, 2},
		{"5-<12 y map 49", mapAt(60, 49), c, 0},
		{"5-<12 y map 48", mapAt(60, 48), c, 1},
		{"5-<12 y map 35", mapAt(143.9, 35), c, 2},
		{"12-17 y map 52", mapAt(144, 52), c, 0},
		{"12-17 y map 51", mapAt(144, 51), c, 1},
		{"12-17 y map 37", mapAt(215.9, 37), c, 2},
		{"cv maximum 6", func(i *PhoenixInput) {
			i.VasoactiveMedications, i.LactateMeasured, i.LactateMmolL, i.MAPMeasured, i.MAPMmHg = 2, true, 12, true, 20
		}, c, 6},
		{"plt 100", func(i *PhoenixInput) { i.PlateletsMeasured, i.PlateletsKPerUL = true, 100 }, g, 0},
		{"plt 99", func(i *PhoenixInput) { i.PlateletsMeasured, i.PlateletsKPerUL = true, 99 }, g, 1},
		{"inr 1.3", func(i *PhoenixInput) { i.INRMeasured, i.INR = true, 1.3 }, g, 0},
		{"inr 1.31", func(i *PhoenixInput) { i.INRMeasured, i.INR = true, 1.31 }, g, 1},
		{"ddimer 2", func(i *PhoenixInput) { i.DDimerMeasured, i.DDimerMgLFEU = true, 2 }, g, 0},
		{"ddimer 2.1", func(i *PhoenixInput) { i.DDimerMeasured, i.DDimerMgLFEU = true, 2.1 }, g, 1},
		{"fib 100", func(i *PhoenixInput) { i.FibrinogenMeasured, i.FibrinogenMgDL = true, 100 }, g, 0},
		{"fib 99", func(i *PhoenixInput) { i.FibrinogenMeasured, i.FibrinogenMgDL = true, 99 }, g, 1},
		{"coag cap 2", func(i *PhoenixInput) {
			i.PlateletsMeasured, i.PlateletsKPerUL, i.INRMeasured, i.INR = true, 50, true, 2
			i.DDimerMeasured, i.DDimerMgLFEU, i.FibrinogenMeasured, i.FibrinogenMgDL = true, 5, true, 50
		}, g, 2},
		{"gcs 11", func(i *PhoenixInput) { i.GCSMeasured, i.GCS = true, 11 }, n, 0},
		{"gcs 10", func(i *PhoenixInput) { i.GCSMeasured, i.GCS = true, 10 }, n, 1},
		{"fixed pupils", func(i *PhoenixInput) { i.FixedPupilsBilateral = true }, n, 2},
	}
	for _, tc := range tests {
		in := phoenixBase
		tc.mod(&in)
		s, err := PhoenixSepsisScore(in)
		if err != nil || tc.pick(s) != tc.want {
			t.Errorf("%s: (%+v, %v), want %d", tc.name, s, err, tc.want)
		}
	}
}

// "Sepsis: … Phoenix Sepsis Score ≥2"; "Septic shock: Sepsis with ≥1 cardiovascular point(s)".
func TestPhoenixSepsisScore_criteria(t *testing.T) {
	in := phoenixBase
	in.GCSMeasured, in.GCS = true, 9 // 1 point
	s, _ := PhoenixSepsisScore(in)
	if s.Total != 1 || s.ScoreAtLeast2 || s.SepticShockCV {
		t.Errorf("1 point: %+v", s)
	}
	in.PlateletsMeasured, in.PlateletsKPerUL = true, 50 // 2 points, no cardiovascular
	s, _ = PhoenixSepsisScore(in)
	if s.Total != 2 || !s.ScoreAtLeast2 || s.SepticShockCV {
		t.Errorf("2 points: %+v", s)
	}
	in.VasoactiveMedications = 1
	s, _ = PhoenixSepsisScore(in)
	if !s.SepticShockCV {
		t.Errorf("with cardiovascular point: %+v", s)
	}
}

// Doc comment example: 2 + 3 + 1 + 0 = 6 (HAND-DERIVED).
func TestPhoenixSepsisScore_example(t *testing.T) {
	in := PhoenixInput{AgeMonths: 30, SpO2FiO2Measured: true, SpO2FiO2: 200, IMV: true, VasoactiveMedications: 1,
		LactateMeasured: true, LactateMmolL: 6, MAPMeasured: true, MAPMmHg: 40, PlateletsMeasured: true,
		PlateletsKPerUL: 80, GCSMeasured: true, GCS: 12}
	s, err := PhoenixSepsisScore(in)
	if err != nil || s.Total != 6 || s.Respiratory != 2 || s.Cardiovascular != 3 || s.Coagulation != 1 || !s.SepticShockCV {
		t.Errorf("got (%+v, %v), want 6", s, err)
	}
}

// parity: ledger phoenix_sepsis, source case 1. The app takes respiratory, vasoactive, lactate and neuro points as
// picked values and computes the MAP points by the same age rows (scoreMAPByAge): 6 months, MAP 30 → 1 (25–38),
// with respiratory 1, one vasoactive, lactate < 5, INR > 1.3, GCS ≤ 10 → 1 + 2 + 1 + 1 = 5, septic shock.
func TestPhoenixSepsisScore_parity(t *testing.T) {
	in := PhoenixInput{AgeMonths: 6, PaO2FiO2Measured: true, PaO2FiO2: 300, AnyRespiratorySupport: true,
		VasoactiveMedications: 1, MAPMeasured: true, MAPMmHg: 30, INRMeasured: true, INR: 1.5, GCSMeasured: true, GCS: 9}
	if s, err := PhoenixSepsisScore(in); err != nil || s.Total != 5 || !s.SepticShockCV {
		t.Errorf("got (%+v, %v), want 5", s, err)
	}
}

func TestPhoenixSepsisScore_errors(t *testing.T) {
	fields := []func(*PhoenixInput, float64){
		func(i *PhoenixInput, v float64) { i.AgeMonths = v },
		func(i *PhoenixInput, v float64) { i.PaO2FiO2 = v },
		func(i *PhoenixInput, v float64) { i.SpO2FiO2 = v },
		func(i *PhoenixInput, v float64) { i.LactateMmolL = v },
		func(i *PhoenixInput, v float64) { i.MAPMmHg = v },
		func(i *PhoenixInput, v float64) { i.PlateletsKPerUL = v },
		func(i *PhoenixInput, v float64) { i.INR = v },
		func(i *PhoenixInput, v float64) { i.DDimerMgLFEU = v },
		func(i *PhoenixInput, v float64) { i.FibrinogenMgDL = v },
	}
	for _, set := range fields {
		for _, v := range nonFiniteB6 {
			in := phoenixBase
			set(&in, v)
			_, err := PhoenixSepsisScore(in)
			wantErrB6(t, "nonfinite", err, ErrNotFinite)
		}
		in := phoenixBase
		set(&in, -1)
		_, err := PhoenixSepsisScore(in)
		wantErrB6(t, "-1", err, ErrOutOfRange)
	}
	in := phoenixBase
	in.AgeMonths = 216
	_, err := PhoenixSepsisScore(in)
	wantErrB6(t, "age 216", err, ErrOutOfRange)
	in = phoenixBase
	in.VasoactiveMedications = -1
	_, err = PhoenixSepsisScore(in)
	wantErrB6(t, "vaso -1", err, ErrOutOfRange)
	in = phoenixBase
	in.GCSMeasured, in.GCS = true, 2
	_, err = PhoenixSepsisScore(in)
	wantErrB6(t, "gcs 2", err, ErrOutOfRange)
	in = phoenixBase
	in.AgeMonths = 0
	if _, err := PhoenixSepsisScore(in); err != nil {
		t.Errorf("age 0: %v", err)
	}
}

func TestPhoenixSepsisScore_noAlloc(t *testing.T) {
	in := PhoenixInput{AgeMonths: 30, SpO2FiO2Measured: true, SpO2FiO2: 200, IMV: true, VasoactiveMedications: 1,
		MAPMeasured: true, MAPMmHg: 40, GCSMeasured: true, GCS: 12}
	if n := testing.AllocsPerRun(100, func() { sinkPhoenix, _ = PhoenixSepsisScore(in) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkPhoenixSepsisScore(b *testing.B) {
	in := PhoenixInput{AgeMonths: 30, SpO2FiO2Measured: true, SpO2FiO2: 200, IMV: true, VasoactiveMedications: 1,
		MAPMeasured: true, MAPMmHg: 40, GCSMeasured: true, GCS: 12}
	b.ReportAllocs()
	for b.Loop() {
		sinkPhoenix, _ = PhoenixSepsisScore(in)
	}
}

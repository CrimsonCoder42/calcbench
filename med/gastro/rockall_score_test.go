package gastro

import "testing"

// rockallBase scores 0 on every pre-endoscopy item.
var rockallBase = RockallFindings{AgeYears: 50, SystolicMmHg: 120, PulseBpm: 80}

// HAND-DERIVED from Rockall 1996 Table III; both sides of every edge.
func TestRockallScorePreEndoscopy_edges(t *testing.T) {
	with := func(mut func(*RockallFindings)) RockallFindings { f := rockallBase; mut(&f); return f }
	tests := []struct {
		name string
		f    RockallFindings
		want int
	}{
		{"all zero", rockallBase, 0},
		{"age 59.99", with(func(f *RockallFindings) { f.AgeYears = 59.99 }), 0},
		{"age 60", with(func(f *RockallFindings) { f.AgeYears = 60 }), 1},
		{"age 79.99", with(func(f *RockallFindings) { f.AgeYears = 79.99 }), 1},
		{"age 80", with(func(f *RockallFindings) { f.AgeYears = 80 }), 2},
		{"pulse 99", with(func(f *RockallFindings) { f.PulseBpm = 99 }), 0},
		{"pulse 100", with(func(f *RockallFindings) { f.PulseBpm = 100 }), 1},
		{"systolic 100", with(func(f *RockallFindings) { f.SystolicMmHg = 100 }), 0},
		{"systolic 99", with(func(f *RockallFindings) { f.SystolicMmHg = 99 }), 2},
		{"hypotension and tachycardia", with(func(f *RockallFindings) { f.SystolicMmHg = 80; f.PulseBpm = 130 }), 2},
		{"major comorbidity", with(func(f *RockallFindings) { f.Comorbidity = RockallComorbidityMajor }), 2},
		{"renal/liver/malignancy", with(func(f *RockallFindings) { f.Comorbidity = RockallComorbidityRenalLiverDisseminatedMalignancy }), 3},
		{"diagnosis ignored", with(func(f *RockallFindings) { f.Diagnosis = RockallDiagnosisUpperGIMalignancy; f.MajorStigmata = true }), 0},
		{"worked", RockallFindings{AgeYears: 72, SystolicMmHg: 95, PulseBpm: 110, Comorbidity: RockallComorbidityMajor}, 5},
		{"maximum 7", RockallFindings{AgeYears: 85, SystolicMmHg: 90, PulseBpm: 120, Comorbidity: RockallComorbidityRenalLiverDisseminatedMalignancy}, 7},
	}
	for _, tc := range tests {
		if got, err := RockallScorePreEndoscopy(tc.f); err != nil || got != tc.want {
			t.Errorf("%s: got (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

func TestRockallScoreComplete_example(t *testing.T) {
	with := func(mut func(*RockallFindings)) RockallFindings { f := rockallBase; mut(&f); return f }
	tests := []struct {
		name string
		f    RockallFindings
		want int
	}{
		{"all zero", rockallBase, 0},
		{"other diagnosis", with(func(f *RockallFindings) { f.Diagnosis = RockallDiagnosisOther }), 1},
		{"upper GI malignancy", with(func(f *RockallFindings) { f.Diagnosis = RockallDiagnosisUpperGIMalignancy }), 2},
		{"major SRH", with(func(f *RockallFindings) { f.MajorStigmata = true }), 2},
		{"maximum 11", RockallFindings{85, 90, 120, RockallComorbidityRenalLiverDisseminatedMalignancy, RockallDiagnosisUpperGIMalignancy, true}, 11},
	}
	for _, tc := range tests {
		if got, err := RockallScoreComplete(tc.f); err != nil || got != tc.want {
			t.Errorf("%s: got (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// parity: ledger rockall_pre_endoscopy and rockall_score_gi. The app takes each item pre-scored (age 0–2,
// shock 0–2, comorbidities 0–2 as picker indices, diagnosis 0–2, stigmata 0–2); its sum for picker values
// (2, 2, 2, 2, 2) is 11 when comorbidity index 2 means 3 points (Table III). Representative inputs here.
func TestRockallScore_parity(t *testing.T) {
	f := RockallFindings{AgeYears: 65, SystolicMmHg: 110, PulseBpm: 105, Comorbidity: RockallComorbidityMajor,
		Diagnosis: RockallDiagnosisOther, MajorStigmata: false}
	if got, err := RockallScorePreEndoscopy(f); err != nil || got != 4 {
		t.Errorf("pre: got (%v, %v), want 4", got, err)
	}
	if got, err := RockallScoreComplete(f); err != nil || got != 5 {
		t.Errorf("complete: got (%v, %v), want 5", got, err)
	}
}

func TestRockallScore_errors(t *testing.T) {
	with := func(mut func(*RockallFindings)) RockallFindings { f := rockallBase; mut(&f); return f }
	pre := func(f RockallFindings) func() error {
		return func() error { _, err := RockallScorePreEndoscopy(f); return err }
	}
	full := func(f RockallFindings) func() error {
		return func() error { _, err := RockallScoreComplete(f); return err }
	}
	runErrCases(t, []errCase{
		{"age NaN", pre(with(func(f *RockallFindings) { f.AgeYears = nan })), ErrNotFinite},
		{"systolic +Inf", pre(with(func(f *RockallFindings) { f.SystolicMmHg = pInf })), ErrNotFinite},
		{"pulse -Inf", full(with(func(f *RockallFindings) { f.PulseBpm = nInf })), ErrNotFinite},
		{"age negative", pre(with(func(f *RockallFindings) { f.AgeYears = -1 })), ErrOutOfRange},
		{"systolic negative", pre(with(func(f *RockallFindings) { f.SystolicMmHg = -1 })), ErrOutOfRange},
		{"pulse negative", pre(with(func(f *RockallFindings) { f.PulseBpm = -1 })), ErrOutOfRange},
		{"comorbidity unknown", pre(with(func(f *RockallFindings) { f.Comorbidity = 3 })), ErrOutOfRange},
		{"diagnosis unknown", full(with(func(f *RockallFindings) { f.Diagnosis = 3 })), ErrOutOfRange},
		{"diagnosis negative", full(with(func(f *RockallFindings) { f.Diagnosis = -1 })), ErrOutOfRange},
	})
}

func TestRockallScore_noAllocation(t *testing.T) {
	f := RockallFindings{72, 95, 110, RockallComorbidityMajor, RockallDiagnosisOther, true}
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = RockallScorePreEndoscopy(f) }); n != 0 {
		t.Errorf("RockallScorePreEndoscopy allocates %v times, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = RockallScoreComplete(f) }); n != 0 {
		t.Errorf("RockallScoreComplete allocates %v times, want 0", n)
	}
}

func BenchmarkRockallScorePreEndoscopy(b *testing.B) {
	b.ReportAllocs()
	f := RockallFindings{72, 95, 110, RockallComorbidityMajor, RockallDiagnosisOther, true}
	for b.Loop() {
		sinkInt, _ = RockallScorePreEndoscopy(f)
	}
}

func BenchmarkRockallScoreComplete(b *testing.B) {
	b.ReportAllocs()
	f := RockallFindings{72, 95, 110, RockallComorbidityMajor, RockallDiagnosisOther, true}
	for b.Loop() {
		sinkInt, _ = RockallScoreComplete(f)
	}
}

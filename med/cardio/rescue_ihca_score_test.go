package cardio

import "testing"

// Tonna 2022 Table 2: each row's points from the reference levels (age ≤ 20 → 2 points).
func TestRESCUEIHCAScore_table2(t *testing.T) {
	tests := []struct {
		age   int
		renal bool
		tod   RESCUEIHCATimeOfDay
		ill   RESCUEIHCAIllness
		rh    RESCUEIHCARhythm
		dur   int
		want  int
	}{
		{0, false, RESCUEIHCADay, RESCUEIHCAMedicalNonCardiac, RESCUEIHCAAsystole, 0, 2},
		{1, false, RESCUEIHCADay, RESCUEIHCAMedicalNonCardiac, RESCUEIHCAAsystole, 0, 4},
		{0, true, RESCUEIHCADay, RESCUEIHCAMedicalNonCardiac, RESCUEIHCAAsystole, 0, 10},
		{0, false, RESCUEIHCAEvening, RESCUEIHCAMedicalNonCardiac, RESCUEIHCAAsystole, 0, 6},
		{0, false, RESCUEIHCANight, RESCUEIHCAMedicalNonCardiac, RESCUEIHCAAsystole, 0, 15},
		{0, false, RESCUEIHCADay, RESCUEIHCAMedicalCardiac, RESCUEIHCAAsystole, 0, 0},
		{0, false, RESCUEIHCADay, RESCUEIHCASurgicalNonCardiac, RESCUEIHCAAsystole, 0, -4},
		{0, false, RESCUEIHCADay, RESCUEIHCASurgicalCardiac, RESCUEIHCAAsystole, 0, -9},
		{0, false, RESCUEIHCADay, RESCUEIHCAMedicalNonCardiac, RESCUEIHCAPEA, 0, 1},
		{0, false, RESCUEIHCADay, RESCUEIHCAMedicalNonCardiac, RESCUEIHCAPulselessVT, 0, -3},
		{0, false, RESCUEIHCADay, RESCUEIHCAMedicalNonCardiac, RESCUEIHCAPalpablePulse, 0, -3},
		{0, false, RESCUEIHCADay, RESCUEIHCAMedicalNonCardiac, RESCUEIHCAVF, 0, -6},
		{0, false, RESCUEIHCADay, RESCUEIHCAMedicalNonCardiac, RESCUEIHCAAsystole, 1, 4},
		{3, true, RESCUEIHCANight, RESCUEIHCASurgicalCardiac, RESCUEIHCAVF, 4, 18}, // doc Example
	}
	for i, tc := range tests {
		got, err := RESCUEIHCAScore(tc.age, tc.renal, tc.tod, tc.ill, tc.rh, tc.dur)
		if err != nil || got != tc.want {
			t.Errorf("case %d = (%d, %v), want %d", i, got, err, tc.want)
		}
	}
}

// parity: ledger rescue_ihca — no source test file or JSON example. The app floors (age − 20) / 10 and
// duration / 10; source case 1: age 45 (2 intervals), evening, medical cardiac, PEA, 25 min (2) → 2+4+4−2−1+4 = 11.
// The app's mortality percentage (logit −1.5 + 0.08 × score) has no source and is not reproduced.
func TestRESCUEIHCAScore_parity(t *testing.T) {
	got, err := RESCUEIHCAScore(2, false, RESCUEIHCAEvening, RESCUEIHCAMedicalCardiac, RESCUEIHCAPEA, 2) // parity: ledger rescue_ihca, source case 1
	if err != nil || got != 11 {
		t.Errorf("(%d, %v), want 11", got, err)
	}
}

func TestRESCUEIHCAScore_errors(t *testing.T) {
	bad := []func() (int, error){
		func() (int, error) { return RESCUEIHCAScore(-1, false, 0, 0, 0, 0) },
		func() (int, error) { return RESCUEIHCAScore(10, false, 0, 0, 0, 0) },
		func() (int, error) { return RESCUEIHCAScore(0, false, 0, 0, 0, -1) },
		func() (int, error) { return RESCUEIHCAScore(0, false, 0, 0, 0, 31) },
		func() (int, error) { return RESCUEIHCAScore(0, false, -1, 0, 0, 0) },
		func() (int, error) { return RESCUEIHCAScore(0, false, 3, 0, 0, 0) },
		func() (int, error) { return RESCUEIHCAScore(0, false, 0, -1, 0, 0) },
		func() (int, error) { return RESCUEIHCAScore(0, false, 0, 4, 0, 0) },
		func() (int, error) { return RESCUEIHCAScore(0, false, 0, 0, -1, 0) },
		func() (int, error) { return RESCUEIHCAScore(0, false, 0, 0, 5, 0) },
	}
	for i, f := range bad {
		_, err := f()
		checkErr(t, "range", err, ErrOutOfRange)
		_ = i
	}
	if _, err := RESCUEIHCAScore(9, true, RESCUEIHCANight, RESCUEIHCASurgicalCardiac, RESCUEIHCAVF, 30); err != nil {
		t.Errorf("upper edges: %v", err)
	}
}

func TestRESCUEIHCAScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() {
		sinkInt, _ = RESCUEIHCAScore(3, true, RESCUEIHCANight, RESCUEIHCASurgicalCardiac, RESCUEIHCAVF, 4)
	}); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkRESCUEIHCAScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = RESCUEIHCAScore(3, true, RESCUEIHCANight, RESCUEIHCASurgicalCardiac, RESCUEIHCAVF, 4)
	}
}

package critical

import "testing"

var msofaZero = MSOFAInput{SpO2FiO2: 450, MAPMmHg: 80, GCS: 15, CreatinineMgDL: 0.9}

// Grissom 2010 Table 2: both sides of every edge. HAND-DERIVED. The liver cell is 3 points: four systems at 4
// plus 3 give the stated "maximum score of 19".
func TestMSOFAScore_table2(t *testing.T) {
	tests := []struct {
		name string
		mod  func(*MSOFAInput)
		want int
	}{
		{"zero", func(*MSOFAInput) {}, 0},
		{"sf 400.1", func(i *MSOFAInput) { i.SpO2FiO2 = 400.1 }, 0},
		{"sf 400", func(i *MSOFAInput) { i.SpO2FiO2 = 400 }, 1},
		{"sf 315.1", func(i *MSOFAInput) { i.SpO2FiO2 = 315.1 }, 1},
		{"sf 315", func(i *MSOFAInput) { i.SpO2FiO2 = 315 }, 2},
		{"sf 235.1", func(i *MSOFAInput) { i.SpO2FiO2 = 235.1 }, 2},
		{"sf 235", func(i *MSOFAInput) { i.SpO2FiO2 = 235 }, 3},
		{"sf 150.1", func(i *MSOFAInput) { i.SpO2FiO2 = 150.1 }, 3},
		{"sf 150", func(i *MSOFAInput) { i.SpO2FiO2 = 150 }, 4},
		{"jaundice", func(i *MSOFAInput) { i.Jaundice = true }, 3},
		{"map 69", func(i *MSOFAInput) { i.MAPMmHg = 69 }, 1},
		{"dopa 5", func(i *MSOFAInput) { i.DopamineUgKgMin = 5 }, 2},
		{"dobutamine", func(i *MSOFAInput) { i.Dobutamine = true }, 2},
		{"dopa 5.1", func(i *MSOFAInput) { i.DopamineUgKgMin = 5.1 }, 3},
		{"dopa 15.1", func(i *MSOFAInput) { i.DopamineUgKgMin = 15.1 }, 4},
		{"epi 0.1", func(i *MSOFAInput) { i.EpinephrineUgKgMin = 0.1 }, 3},
		{"norepi 0.11", func(i *MSOFAInput) { i.NorepinephrineUgKgMin = 0.11 }, 4},
		{"gcs 13", func(i *MSOFAInput) { i.GCS = 13 }, 1},
		{"gcs 12", func(i *MSOFAInput) { i.GCS = 12 }, 2},
		{"gcs 9", func(i *MSOFAInput) { i.GCS = 9 }, 3},
		{"gcs 5", func(i *MSOFAInput) { i.GCS = 5 }, 4},
		{"cr 1.19", func(i *MSOFAInput) { i.CreatinineMgDL = 1.19 }, 0},
		{"cr 1.2", func(i *MSOFAInput) { i.CreatinineMgDL = 1.2 }, 1},
		{"cr 2.0", func(i *MSOFAInput) { i.CreatinineMgDL = 2.0 }, 2},
		{"cr 3.5", func(i *MSOFAInput) { i.CreatinineMgDL = 3.5 }, 3},
		{"cr 4.9", func(i *MSOFAInput) { i.CreatinineMgDL = 4.9 }, 3},
		{"cr 5.0", func(i *MSOFAInput) { i.CreatinineMgDL = 5.0 }, 4},
		{"maximum 19", func(i *MSOFAInput) {
			*i = MSOFAInput{SpO2FiO2: 100, Jaundice: true, MAPMmHg: 50, NorepinephrineUgKgMin: 0.3, GCS: 3, CreatinineMgDL: 6}
		}, 19},
	}
	for _, tc := range tests {
		in := msofaZero
		tc.mod(&in)
		if s, err := MSOFAScore(in); err != nil || s.Total != tc.want {
			t.Errorf("%s: (%+v, %v), want %d", tc.name, s, err, tc.want)
		}
	}
}

// Doc comment example: 3+3+1+2+1 = 10 (HAND-DERIVED).
func TestMSOFAScore_example(t *testing.T) {
	in := MSOFAInput{SpO2FiO2: 94 / 0.40, Jaundice: true, MAPMmHg: 65, GCS: 12, CreatinineMgDL: 1.5}
	if s, err := MSOFAScore(in); err != nil || s.Total != 10 {
		t.Errorf("got (%+v, %v), want 10", s, err)
	}
}

// parity: ledger msofa, source case 1. The app sums picked component scores with the liver item worth 3
// (LIVER_PTS = [0, 3]), as here: respiratory 2, jaundice, cardiovascular 3, CNS 1, renal 2 → 11.
func TestMSOFAScore_parity(t *testing.T) {
	in := MSOFAInput{SpO2FiO2: 300, Jaundice: true, MAPMmHg: 60, NorepinephrineUgKgMin: 0.05, GCS: 14, CreatinineMgDL: 2.5}
	if s, err := MSOFAScore(in); err != nil || s.Total != 11 {
		t.Errorf("got (%+v, %v), want 11", s, err)
	}
}

func TestMSOFAScore_errors(t *testing.T) {
	fields := []func(*MSOFAInput, float64){
		func(i *MSOFAInput, v float64) { i.SpO2FiO2 = v },
		func(i *MSOFAInput, v float64) { i.MAPMmHg = v },
		func(i *MSOFAInput, v float64) { i.DopamineUgKgMin = v },
		func(i *MSOFAInput, v float64) { i.EpinephrineUgKgMin = v },
		func(i *MSOFAInput, v float64) { i.NorepinephrineUgKgMin = v },
		func(i *MSOFAInput, v float64) { i.CreatinineMgDL = v },
	}
	for _, set := range fields {
		for _, v := range nonFiniteB6 {
			in := msofaZero
			set(&in, v)
			_, err := MSOFAScore(in)
			wantErrB6(t, "nonfinite", err, ErrNotFinite)
		}
		in := msofaZero
		set(&in, -1)
		_, err := MSOFAScore(in)
		wantErrB6(t, "-1", err, ErrOutOfRange)
	}
	in := msofaZero
	in.GCS = 16
	_, err := MSOFAScore(in)
	wantErrB6(t, "gcs 16", err, ErrOutOfRange)
}

func TestMSOFAScore_noAlloc(t *testing.T) {
	in := MSOFAInput{SpO2FiO2: 235, Jaundice: true, MAPMmHg: 65, GCS: 12, CreatinineMgDL: 1.5}
	if n := testing.AllocsPerRun(100, func() { sinkMSOFA, _ = MSOFAScore(in) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkMSOFAScore(b *testing.B) {
	in := MSOFAInput{SpO2FiO2: 235, Jaundice: true, MAPMmHg: 65, GCS: 12, CreatinineMgDL: 1.5}
	b.ReportAllocs()
	for b.Loop() {
		sinkMSOFA, _ = MSOFAScore(in)
	}
}

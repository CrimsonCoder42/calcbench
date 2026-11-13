package critical

import "testing"

// Imhoff 2014 / Ala 2020 REMS table: both sides of every edge, one item at a time from an all-zero base.
// HAND-DERIVED.
func TestREMSScore_table(t *testing.T) {
	type in struct {
		age, mapv, hr, rr, spo2 float64
		gcs                     int
	}
	base := in{30, 90, 80, 16, 98, 15}
	tests := []struct {
		name string
		mod  func(*in)
		want int
	}{
		{"zero", func(*in) {}, 0},
		{"age 44.9", func(v *in) { v.age = 44.9 }, 0},
		{"age 45", func(v *in) { v.age = 45 }, 2},
		{"age 55", func(v *in) { v.age = 55 }, 3},
		{"age 65", func(v *in) { v.age = 65 }, 5},
		{"age 74", func(v *in) { v.age = 74 }, 5},
		{"age 74.5", func(v *in) { v.age = 74.5 }, 6},
		{"map 49", func(v *in) { v.mapv = 49 }, 4},
		{"map 49.5", func(v *in) { v.mapv = 49.5 }, 2},
		{"map 69", func(v *in) { v.mapv = 69 }, 2},
		{"map 70", func(v *in) { v.mapv = 70 }, 0},
		{"map 109", func(v *in) { v.mapv = 109 }, 0},
		{"map 110", func(v *in) { v.mapv = 110 }, 2},
		{"map 130", func(v *in) { v.mapv = 130 }, 3},
		{"map 159", func(v *in) { v.mapv = 159 }, 3},
		{"map 160", func(v *in) { v.mapv = 160 }, 4},
		{"hr 39", func(v *in) { v.hr = 39 }, 4},
		{"hr 40", func(v *in) { v.hr = 40 }, 3},
		{"hr 55", func(v *in) { v.hr = 55 }, 2},
		{"hr 70", func(v *in) { v.hr = 70 }, 0},
		{"hr 110", func(v *in) { v.hr = 110 }, 2},
		{"hr 140", func(v *in) { v.hr = 140 }, 3},
		{"hr 179", func(v *in) { v.hr = 179 }, 3},
		{"hr 180", func(v *in) { v.hr = 180 }, 4},
		{"rr 5", func(v *in) { v.rr = 5 }, 4},
		{"rr 6", func(v *in) { v.rr = 6 }, 2},
		{"rr 10", func(v *in) { v.rr = 10 }, 1},
		{"rr 12", func(v *in) { v.rr = 12 }, 0},
		{"rr 24", func(v *in) { v.rr = 24 }, 0},
		{"rr 25", func(v *in) { v.rr = 25 }, 1},
		{"rr 35", func(v *in) { v.rr = 35 }, 3},
		{"rr 49", func(v *in) { v.rr = 49 }, 3},
		{"rr 50", func(v *in) { v.rr = 50 }, 4},
		{"spo2 90", func(v *in) { v.spo2 = 90 }, 0},
		{"spo2 89", func(v *in) { v.spo2 = 89 }, 1},
		{"spo2 86", func(v *in) { v.spo2 = 86 }, 1},
		{"spo2 85", func(v *in) { v.spo2 = 85 }, 3},
		{"spo2 75", func(v *in) { v.spo2 = 75 }, 3},
		{"spo2 74", func(v *in) { v.spo2 = 74 }, 4},
		{"gcs 14", func(v *in) { v.gcs = 14 }, 0},
		{"gcs 13", func(v *in) { v.gcs = 13 }, 1},
		{"gcs 10", func(v *in) { v.gcs = 10 }, 2},
		{"gcs 7", func(v *in) { v.gcs = 7 }, 3},
		{"gcs 4", func(v *in) { v.gcs = 4 }, 4},
		{"maximum 26", func(v *in) { *v = in{80, 40, 30, 4, 70, 3} }, 26},
	}
	for _, tc := range tests {
		v := base
		tc.mod(&v)
		if got, err := REMSScore(v.age, v.mapv, v.hr, v.rr, v.spo2, v.gcs); err != nil || got != tc.want {
			t.Errorf("%s: (%d, %v), want %d", tc.name, got, err, tc.want)
		}
	}
}

// Doc comment example: 5+2+2+1+1+1 = 12 (HAND-DERIVED).
func TestREMSScore_example(t *testing.T) {
	if got, err := REMSScore(70, 65, 115, 26, 88, 13); err != nil || got != 12 {
		t.Errorf("got (%d, %v), want 12", got, err)
	}
}

// parity: ledger rems, source case 1 (the app's item bands equal the table; the app derives MAP from SBP/DBP,
// here passed directly): age 60, MAP 100, HR 120, RR 30, SpO2 92, GCS 15 → 3+0+2+1+0+0 = 6. The app's risk
// bands (3/6/9/13) are not in either secondary and are not ported.
func TestREMSScore_parity(t *testing.T) {
	if got, err := REMSScore(60, 100, 120, 30, 92, 15); err != nil || got != 6 {
		t.Errorf("got (%d, %v), want 6", got, err)
	}
}

func TestREMSScore_errors(t *testing.T) {
	for _, v := range nonFiniteB6 {
		for i := range 5 {
			a := [5]float64{30, 90, 80, 16, 98}
			a[i] = v
			_, err := REMSScore(a[0], a[1], a[2], a[3], a[4], 15)
			wantErrB6(t, "nonfinite", err, ErrNotFinite)
		}
	}
	for i := range 5 {
		a := [5]float64{30, 90, 80, 16, 98}
		a[i] = -1
		_, err := REMSScore(a[0], a[1], a[2], a[3], a[4], 15)
		wantErrB6(t, "-1", err, ErrOutOfRange)
	}
	_, err := REMSScore(30, 90, 80, 16, 100.1, 15)
	wantErrB6(t, "spo2 100.1", err, ErrOutOfRange)
	_, err = REMSScore(30, 90, 80, 16, 98, 2)
	wantErrB6(t, "gcs 2", err, ErrOutOfRange)
}

func TestREMSScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkB6Int, _ = REMSScore(70, 65, 115, 26, 88, 13) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkREMSScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkB6Int, _ = REMSScore(70, 65, 115, 26, 88, 13)
	}
}

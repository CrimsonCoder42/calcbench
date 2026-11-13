package critical

import "testing"

// Hoesl 2022 Table 1 / El-Menyar 2017: both sides of every edge. HAND-DERIVED.
func TestLRINECScore_table(t *testing.T) {
	type in struct{ crp, wbc, hb, na, cr, glu float64 }
	base := in{10, 8, 14, 140, 1.0, 100}
	tests := []struct {
		name string
		mod  func(*in)
		want int
	}{
		{"zero", func(*in) {}, 0},
		{"crp 150", func(v *in) { v.crp = 150 }, 0},
		{"crp 150.1", func(v *in) { v.crp = 150.1 }, 4},
		{"wbc 14.9", func(v *in) { v.wbc = 14.9 }, 0},
		{"wbc 15", func(v *in) { v.wbc = 15 }, 1},
		{"wbc 25", func(v *in) { v.wbc = 25 }, 1},
		{"wbc 25.1", func(v *in) { v.wbc = 25.1 }, 2},
		{"hb 13.6", func(v *in) { v.hb = 13.6 }, 0},
		{"hb 13.5", func(v *in) { v.hb = 13.5 }, 1},
		{"hb 11", func(v *in) { v.hb = 11 }, 1},
		{"hb 10.9", func(v *in) { v.hb = 10.9 }, 2},
		{"na 135", func(v *in) { v.na = 135 }, 0},
		{"na 134.9", func(v *in) { v.na = 134.9 }, 2},
		{"cr 1.6", func(v *in) { v.cr = 1.6 }, 0},
		{"cr 1.61", func(v *in) { v.cr = 1.61 }, 2},
		{"glu 180", func(v *in) { v.glu = 180 }, 0},
		{"glu 181", func(v *in) { v.glu = 181 }, 1},
		{"maximum 13", func(v *in) { *v = in{300, 30, 10, 130, 2, 200} }, 13},
	}
	for _, tc := range tests {
		v := base
		tc.mod(&v)
		if got, err := LRINECScore(v.crp, v.wbc, v.hb, v.na, v.cr, v.glu); err != nil || got != tc.want {
			t.Errorf("%s: (%d, %v), want %d", tc.name, got, err, tc.want)
		}
	}
}

// Doc comment example: 4+1+1+2+2+0 = 10 (HAND-DERIVED).
func TestLRINECScore_example(t *testing.T) {
	if got, err := LRINECScore(200, 20, 12, 132, 1.8, 150); err != nil || got != 10 {
		t.Errorf("got (%d, %v), want 10", got, err)
	}
}

// parity: ledger lrinec_score, source case 1. The app scores CRP ≥ 150 as 4; Hoesl prints "≤150 = 0", so CRP
// exactly 150 is 0 here and 4 in the app (a CONFIRM item). Away from 150 they agree: CRP 160, WBC 16, Hb 12,
// Na 136, creatinine 1.0, glucose 200 → 4+1+1+0+0+1 = 7.
func TestLRINECScore_parity(t *testing.T) {
	if got, err := LRINECScore(160, 16, 12, 136, 1.0, 200); err != nil || got != 7 {
		t.Errorf("got (%d, %v), want 7", got, err)
	}
	if got, _ := LRINECScore(150, 8, 14, 140, 1.0, 100); got != 0 {
		t.Errorf("CRP 150: got %d, want 0 (app: 4)", got)
	}
}

func TestLRINECScore_errors(t *testing.T) {
	for i := range 6 {
		for _, v := range nonFiniteB6 {
			a := [6]float64{10, 8, 14, 140, 1, 100}
			a[i] = v
			_, err := LRINECScore(a[0], a[1], a[2], a[3], a[4], a[5])
			wantErrB6(t, "nonfinite", err, ErrNotFinite)
		}
		a := [6]float64{10, 8, 14, 140, 1, 100}
		a[i] = -1
		_, err := LRINECScore(a[0], a[1], a[2], a[3], a[4], a[5])
		wantErrB6(t, "-1", err, ErrOutOfRange)
	}
}

func TestLRINECScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkB6Int, _ = LRINECScore(200, 20, 12, 132, 1.8, 150) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkLRINECScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkB6Int, _ = LRINECScore(200, 20, 12, 132, 1.8, 150)
	}
}

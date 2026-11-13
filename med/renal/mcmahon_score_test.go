package renal

import "testing"

// mcIn is one McMahonScore call.
type mcIn struct {
	age          float64
	female       bool
	cr, ca, cpk  float64
	origin       bool
	phos, bicarb float64
}

func (m mcIn) score() (float64, error) {
	return McMahonScore(m.age, m.female, m.cr, m.ca, m.cpk, m.origin, m.phos, m.bicarb)
}

// zeroMc scores 0 on every item of McMahon 2013 Table 3.
var zeroMc = mcIn{age: 40, cr: 1.0, ca: 9.0, cpk: 10000, phos: 3.0, bicarb: 24}

// McMahon 2013 Table 3 (PMC5152583), Score column: both sides of every printed edge. HAND-DERIVED.
// Points are multiples of 0.5, so equality is exact.
func TestMcMahonScore_table3(t *testing.T) {
	tests := []struct {
		name string
		mod  func(*mcIn)
		want float64
	}{
		{"all zero", func(*mcIn) {}, 0},
		{"age 50", func(m *mcIn) { m.age = 50 }, 0},
		{"age 50.1", func(m *mcIn) { m.age = 50.1 }, 1.5},
		{"age 70", func(m *mcIn) { m.age = 70 }, 1.5},
		{"age 70.1", func(m *mcIn) { m.age = 70.1 }, 2.5},
		{"age 80", func(m *mcIn) { m.age = 80 }, 2.5},
		{"age 80.1", func(m *mcIn) { m.age = 80.1 }, 3},
		{"female", func(m *mcIn) { m.female = true }, 1},
		{"cr 1.39", func(m *mcIn) { m.cr = 1.39 }, 0},
		{"cr 1.4", func(m *mcIn) { m.cr = 1.4 }, 1.5},
		{"cr 2.2", func(m *mcIn) { m.cr = 2.2 }, 1.5},
		{"cr 2.21", func(m *mcIn) { m.cr = 2.21 }, 3},
		{"ca 7.5", func(m *mcIn) { m.ca = 7.5 }, 0},
		{"ca 7.49", func(m *mcIn) { m.ca = 7.49 }, 2},
		{"cpk 40000", func(m *mcIn) { m.cpk = 40000 }, 0},
		{"cpk 40001", func(m *mcIn) { m.cpk = 40001 }, 2},
		{"origin", func(m *mcIn) { m.origin = true }, 3},
		{"phos 3.99", func(m *mcIn) { m.phos = 3.99 }, 0},
		{"phos 4.0", func(m *mcIn) { m.phos = 4.0 }, 1.5},
		{"phos 5.4", func(m *mcIn) { m.phos = 5.4 }, 1.5},
		{"phos 5.41", func(m *mcIn) { m.phos = 5.41 }, 3},
		{"bicarb 19", func(m *mcIn) { m.bicarb = 19 }, 0},
		{"bicarb 18.9", func(m *mcIn) { m.bicarb = 18.9 }, 2},
		{"maximum 19", func(m *mcIn) {
			*m = mcIn{age: 85, female: true, cr: 3, ca: 7, cpk: 50000, origin: true, phos: 6, bicarb: 15}
		}, 19},
	}
	for _, tc := range tests {
		in := zeroMc
		tc.mod(&in)
		if got, err := in.score(); err != nil || got != tc.want {
			t.Errorf("%s: (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// Doc comment example: 2.5+0+3+2+2+3+1.5+2 = 16 (HAND-DERIVED).
func TestMcMahonScore_example(t *testing.T) {
	if got, err := McMahonScore(75, false, 2.5, 7.0, 50000, true, 5.0, 18); err != nil || got != 16 {
		t.Errorf("got (%v, %v), want 16", got, err)
	}
}

// parity: ledger mcmahon_rhabdomyolysis, source case 1 (app points equal Table 3 at integer ages):
// age 65 female, creatinine 1.8, calcium 8.0, CPK 20 000, origin trauma, phosphate 4.5, bicarbonate 22 →
// 1.5+1+1.5+0+0+3+1.5+0 = 8.5. The app's bands (<6 "Low Risk", ≥6, ≥10 "High Risk — 52%") are not in the
// paper (it reports < 5 and > 10) and are not ported.
func TestMcMahonScore_parity(t *testing.T) {
	if got, err := McMahonScore(65, true, 1.8, 8.0, 20000, true, 4.5, 22); err != nil || got != 8.5 {
		t.Errorf("got (%v, %v), want 8.5", got, err)
	}
}

func TestMcMahonScore_errors(t *testing.T) {
	runErrCases(t, []errCase{
		{"age NaN", func() error { _, e := McMahonScore(nan, false, 1, 9, 1, false, 3, 24); return e }, ErrNotFinite},
		{"cr +Inf", func() error { _, e := McMahonScore(40, false, pInf, 9, 1, false, 3, 24); return e }, ErrNotFinite},
		{"ca -Inf", func() error { _, e := McMahonScore(40, false, 1, nInf, 1, false, 3, 24); return e }, ErrNotFinite},
		{"cpk NaN", func() error { _, e := McMahonScore(40, false, 1, 9, nan, false, 3, 24); return e }, ErrNotFinite},
		{"phos NaN", func() error { _, e := McMahonScore(40, false, 1, 9, 1, false, nan, 24); return e }, ErrNotFinite},
		{"bicarb NaN", func() error { _, e := McMahonScore(40, false, 1, 9, 1, false, 3, nan); return e }, ErrNotFinite},
		{"age 17.9", func() error { _, e := McMahonScore(17.9, false, 1, 9, 1, false, 3, 24); return e }, ErrOutOfRange},
		{"age 120.1", func() error { _, e := McMahonScore(120.1, false, 1, 9, 1, false, 3, 24); return e }, ErrOutOfRange},
		{"cr -1", func() error { _, e := McMahonScore(40, false, -1, 9, 1, false, 3, 24); return e }, ErrOutOfRange},
		{"bicarb -1", func() error { _, e := McMahonScore(40, false, 1, 9, 1, false, 3, -1); return e }, ErrOutOfRange},
	})
	for _, age := range []float64{18, 120} {
		if _, err := McMahonScore(age, false, 1, 9, 1, false, 3, 24); err != nil {
			t.Errorf("age %v: %v", age, err)
		}
	}
}

func TestMcMahonScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sink, _ = McMahonScore(75, false, 2.5, 7, 50000, true, 5, 18) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkMcMahonScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sink, _ = McMahonScore(75, false, 2.5, 7, 50000, true, 5, 18)
	}
}

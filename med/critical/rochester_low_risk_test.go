package critical

import "testing"

// Jaskiewicz 1994 abstract: every criterion; both sides of each numeric edge. HAND-DERIVED.
func TestRochesterLowRisk_source(t *testing.T) {
	ok := RochesterInputs{AppearsWell: true, PreviouslyHealthy: true, NoFocalInfection: true, WBCE9PerL: 15.0, BandsE9PerL: 1.5, UrineWBCPerHPF: 10}
	if got, err := RochesterLowRisk(30, ok); err != nil || !got {
		t.Errorf("all met → (%v, %v), want true", got, err)
	}
	mods := []struct {
		name string
		mod  func(*RochesterInputs)
		want bool
	}{
		{"WBC 5.0", func(r *RochesterInputs) { r.WBCE9PerL = 5.0 }, true},
		{"WBC 4.9", func(r *RochesterInputs) { r.WBCE9PerL = 4.9 }, false},
		{"WBC 15.1", func(r *RochesterInputs) { r.WBCE9PerL = 15.1 }, false},
		{"bands 1.6", func(r *RochesterInputs) { r.BandsE9PerL = 1.6 }, false},
		{"urine 11", func(r *RochesterInputs) { r.UrineWBCPerHPF = 11 }, false},
		{"diarrhea, stool 5", func(r *RochesterInputs) { r.Diarrhea, r.StoolWBCPerHPF = true, 5 }, true},
		{"diarrhea, stool 6", func(r *RochesterInputs) { r.Diarrhea, r.StoolWBCPerHPF = true, 6 }, false},
		{"no diarrhea, stool 6", func(r *RochesterInputs) { r.StoolWBCPerHPF = 6 }, true},
		{"not well", func(r *RochesterInputs) { r.AppearsWell = false }, false},
		{"not healthy", func(r *RochesterInputs) { r.PreviouslyHealthy = false }, false},
		{"focal infection", func(r *RochesterInputs) { r.NoFocalInfection = false }, false},
	}
	for _, m := range mods {
		in := ok
		m.mod(&in)
		if got, err := RochesterLowRisk(30, in); err != nil || got != m.want {
			t.Errorf("%s: (%v, %v), want %v", m.name, got, err, m.want)
		}
	}
}

// parity: ledger rochester_criteria, source case 1. The app counts twelve 0/1 items (all met = low risk); the same
// case with every item met is low risk here.
func TestRochesterLowRisk_parity(t *testing.T) {
	in := RochesterInputs{AppearsWell: true, PreviouslyHealthy: true, NoFocalInfection: true, WBCE9PerL: 9, BandsE9PerL: 0.5, UrineWBCPerHPF: 2}
	if got, _ := RochesterLowRisk(45, in); !got {
		t.Error("got false, want true")
	}
}

func TestRochesterLowRisk_errors(t *testing.T) {
	for _, v := range nonFiniteB6 {
		_, err := RochesterLowRisk(v, RochesterInputs{})
		wantErrB6(t, "age", err, ErrNotFinite)
		_, err = RochesterLowRisk(30, RochesterInputs{WBCE9PerL: v})
		wantErrB6(t, "wbc", err, ErrNotFinite)
		_, err = RochesterLowRisk(30, RochesterInputs{StoolWBCPerHPF: v})
		wantErrB6(t, "stool", err, ErrNotFinite)
	}
	for _, a := range []float64{-1, 60.5} {
		_, err := RochesterLowRisk(a, RochesterInputs{})
		wantErrB6(t, "age range", err, ErrOutOfRange)
	}
	_, err := RochesterLowRisk(30, RochesterInputs{BandsE9PerL: -0.1})
	wantErrB6(t, "bands < 0", err, ErrOutOfRange)
}

func TestRochesterLowRisk_noAlloc(t *testing.T) {
	in := RochesterInputs{AppearsWell: true, PreviouslyHealthy: true, NoFocalInfection: true, WBCE9PerL: 9, BandsE9PerL: 0.5, UrineWBCPerHPF: 2}
	if n := testing.AllocsPerRun(100, func() { sinkB6Bool, _ = RochesterLowRisk(30, in) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkRochesterLowRisk(b *testing.B) {
	b.ReportAllocs()
	in := RochesterInputs{AppearsWell: true, PreviouslyHealthy: true, NoFocalInfection: true, WBCE9PerL: 9, BandsE9PerL: 0.5, UrineWBCPerHPF: 2}
	for b.Loop() {
		sinkB6Bool, _ = RochesterLowRisk(30, in)
	}
}

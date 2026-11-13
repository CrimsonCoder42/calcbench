package emergency

import (
	"errors"
	"math"
	"testing"
)

var (
	sinkFloat  float64
	sinkNEDOCS NEDOCSCategory
)

var nedocsBase = NEDOCSInputs{EDCensus: 45, EDBeds: 30, AdmitsWaiting: 10, HospitalBeds: 300, Ventilators: 1, LongestAdmitHours: 6, LongestWaitingRoomHours: 2}

// HAND-DERIVED from "NEDOCS score = 85.8(C/A) + 600(F/B) + 13.4(D) + 0.93(E) + 5.64(G) – 20" (Ilhan 2020
// PMC7360290). Tolerance 1e-9 absolute: a few products of decimal constants.
func TestNEDOCSScore_equation(t *testing.T) {
	cases := []struct {
		name string
		in   NEDOCSInputs
		want float64
	}{
		{"doc example", nedocsBase, 158.96},
		{"empty ED", NEDOCSInputs{EDBeds: 20, HospitalBeds: 200}, -20},
		{"each term", NEDOCSInputs{EDCensus: 10, EDBeds: 10, AdmitsWaiting: 1, HospitalBeds: 100, Ventilators: 2, LongestAdmitHours: 10, LongestWaitingRoomHours: 1},
			85.8 + 6 + 26.8 + 9.3 + 5.64 - 20},
	}
	for _, c := range cases {
		got, err := NEDOCSScore(c.in)
		if err != nil || math.Abs(got-c.want) > 1e-9 {
			t.Errorf("%s: got %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

// parity: ledger nedocs, source case 1: the consumer computes the same equation, then rounds and clamps to
// 0–200; for the doc example it shows 159, this function returns 158.96 and NEDOCSLevel gives level 5.
func TestNEDOCSScore_parity(t *testing.T) {
	got, err := NEDOCSScore(nedocsBase)
	if err != nil || math.Round(got) != 159 {
		t.Fatalf("got %v, %v; want ≈ 159", got, err)
	}
	if lv, _ := NEDOCSLevel(got); lv != NEDOCSSeverelyOvercrowded {
		t.Errorf("level %q, want %q", lv, NEDOCSSeverelyOvercrowded)
	}
}

func TestNEDOCSScore_errors(t *testing.T) {
	cases := []struct {
		name string
		mod  func(*NEDOCSInputs)
		want error
	}{
		{"census NaN", func(in *NEDOCSInputs) { in.EDCensus = math.NaN() }, ErrNotFinite},
		{"beds +Inf", func(in *NEDOCSInputs) { in.EDBeds = math.Inf(1) }, ErrNotFinite},
		{"hours -Inf", func(in *NEDOCSInputs) { in.LongestAdmitHours = math.Inf(-1) }, ErrNotFinite},
		{"ED beds 0", func(in *NEDOCSInputs) { in.EDBeds = 0 }, ErrOutOfRange},
		{"hospital beds 0.5", func(in *NEDOCSInputs) { in.HospitalBeds = 0.5 }, ErrOutOfRange},
		{"ventilators negative", func(in *NEDOCSInputs) { in.Ventilators = -1 }, ErrOutOfRange},
		{"waiting room negative", func(in *NEDOCSInputs) { in.LongestWaitingRoomHours = -1 }, ErrOutOfRange},
	}
	for _, c := range cases {
		in := nedocsBase
		c.mod(&in)
		if _, err := NEDOCSScore(in); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

// Dzelebdzic 2026 bands with Ilhan 2020 names: both sides of each edge.
func TestNEDOCSLevel_edges(t *testing.T) {
	cases := []struct {
		s    float64
		want NEDOCSCategory
	}{
		{-20, NEDOCSNotBusy}, {0, NEDOCSNotBusy}, {20, NEDOCSNotBusy}, {20.5, NEDOCSBusy}, {21, NEDOCSBusy},
		{60, NEDOCSBusy}, {61, NEDOCSExtremelyBusy}, {100, NEDOCSExtremelyBusy}, {101, NEDOCSOvercrowded},
		{140, NEDOCSOvercrowded}, {141, NEDOCSSeverelyOvercrowded}, {180, NEDOCSSeverelyOvercrowded},
		{181, NEDOCSDangerouslyOvercrowded}, {200, NEDOCSDangerouslyOvercrowded}, {250, NEDOCSDangerouslyOvercrowded},
	}
	for _, c := range cases {
		if got, err := NEDOCSLevel(c.s); err != nil || got != c.want {
			t.Errorf("%v: got %q, %v; want %q", c.s, got, err, c.want)
		}
	}
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := NEDOCSLevel(v); !errors.Is(err, ErrNotFinite) {
			t.Errorf("%v: err = %v, want ErrNotFinite", v, err)
		}
	}
}

func TestNEDOCS_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = NEDOCSScore(nedocsBase) }); n != 0 {
		t.Errorf("NEDOCSScore allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sinkNEDOCS, _ = NEDOCSLevel(120) }); n != 0 {
		t.Errorf("NEDOCSLevel allocs = %v, want 0", n)
	}
}

func BenchmarkNEDOCSScore(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkFloat, _ = NEDOCSScore(nedocsBase)
	}
}

func BenchmarkNEDOCSLevel(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkNEDOCS, _ = NEDOCSLevel(float64(i % 210))
	}
}

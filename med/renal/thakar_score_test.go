package renal

import (
	"errors"
	"math"
	"testing"
)

var sinkThakar int

func thakarBase() ThakarInput {
	return ThakarInput{LVEFPercent: 55, PreoperativeCreatinineMgDL: 0.9}
}

// Crosina 2017 Table 1 and Kozioł 2026 Table 1 (both citing Thakar 2005): each item alone, both sides of the
// LVEF and creatinine edges, and the printed maximum 17.
func TestThakarScore_points(t *testing.T) {
	type mod func(*ThakarInput)
	tests := []struct {
		name string
		f    mod
		want int
	}{
		{"none", func(*ThakarInput) {}, 0},
		{"female", func(in *ThakarInput) { in.Female = true }, 1},
		{"chf", func(in *ThakarInput) { in.CongestiveHeartFailure = true }, 1},
		{"lvef 34.9", func(in *ThakarInput) { in.LVEFPercent = 34.9 }, 1},
		{"lvef 35", func(in *ThakarInput) { in.LVEFPercent = 35 }, 0},
		{"iabp", func(in *ThakarInput) { in.PreoperativeIABP = true }, 2},
		{"copd", func(in *ThakarInput) { in.COPD = true }, 1},
		{"insulin dm", func(in *ThakarInput) { in.InsulinRequiringDiabetes = true }, 1},
		{"previous surgery", func(in *ThakarInput) { in.PreviousCardiacSurgery = true }, 1},
		{"emergency", func(in *ThakarInput) { in.EmergencySurgery = true }, 2},
		{"valve only", func(in *ThakarInput) { in.Surgery = ThakarValveOnly }, 1},
		{"cabg + valve", func(in *ThakarInput) { in.Surgery = ThakarCABGPlusValve }, 2},
		{"other cardiac", func(in *ThakarInput) { in.Surgery = ThakarOtherCardiac }, 2},
		{"creatinine 1.19", func(in *ThakarInput) { in.PreoperativeCreatinineMgDL = 1.19 }, 0},
		{"creatinine 1.2", func(in *ThakarInput) { in.PreoperativeCreatinineMgDL = 1.2 }, 2},
		{"creatinine 2.09", func(in *ThakarInput) { in.PreoperativeCreatinineMgDL = 2.09 }, 2},
		{"creatinine 2.1", func(in *ThakarInput) { in.PreoperativeCreatinineMgDL = 2.1 }, 5},
		{"maximum", func(in *ThakarInput) {
			*in = ThakarInput{true, true, 20, true, true, true, true, true, ThakarOtherCardiac, 3}
		}, 17},
	}
	for _, tc := range tests {
		in := thakarBase()
		tc.f(&in)
		if got, err := ThakarScore(in); err != nil || got != tc.want {
			t.Errorf("%s: got %d (%v), want %d", tc.name, got, err, tc.want)
		}
	}
}

// HAND-DERIVED from Table 1: female (1) + LVEF 30 % (1) + valve only (1) + creatinine 1.5 mg/dL (2) = 5.
func TestThakarScore_source(t *testing.T) {
	in := ThakarInput{Female: true, LVEFPercent: 30, Surgery: ThakarValveOnly, PreoperativeCreatinineMgDL: 1.5}
	if got, err := ThakarScore(in); err != nil || got != 5 {
		t.Errorf("got %d (%v), want 5", got, err)
	}
}

// parity: ledger thakar_score, source case 1. No _test.go or JSON example; the entry takes pre-scored items
// (iabp 0/2, emergency 0/2, surgery type 0–2, creatinine 0/2/5). Its all-maximum options sum to 17.
func TestThakarScore_parity(t *testing.T) {
	in := ThakarInput{true, true, 25, true, true, true, true, true, ThakarCABGPlusValve, 2.5}
	if got, _ := ThakarScore(in); got != 17 {
		t.Errorf("got %d, want 17", got)
	}
}

func TestThakarScore_errors(t *testing.T) {
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		in := thakarBase()
		in.LVEFPercent = v
		if _, err := ThakarScore(in); !errors.Is(err, ErrNotFinite) {
			t.Errorf("LVEF %v: err = %v, want ErrNotFinite", v, err)
		}
		in = thakarBase()
		in.PreoperativeCreatinineMgDL = v
		if _, err := ThakarScore(in); !errors.Is(err, ErrNotFinite) {
			t.Errorf("creatinine %v: err = %v, want ErrNotFinite", v, err)
		}
	}
	bad := []func(*ThakarInput){
		func(in *ThakarInput) { in.LVEFPercent = -0.1 },
		func(in *ThakarInput) { in.LVEFPercent = 100.1 },
		func(in *ThakarInput) { in.PreoperativeCreatinineMgDL = 0 },
		func(in *ThakarInput) { in.Surgery = ThakarSurgeryType(4) },
		func(in *ThakarInput) { in.Surgery = ThakarSurgeryType(-1) },
	}
	for i, f := range bad {
		in := thakarBase()
		f(&in)
		if _, err := ThakarScore(in); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("case %d: err = %v, want ErrOutOfRange", i, err)
		}
	}
}

func TestThakarScore_noAlloc(t *testing.T) {
	in := ThakarInput{Female: true, LVEFPercent: 30, Surgery: ThakarValveOnly, PreoperativeCreatinineMgDL: 1.5}
	if n := testing.AllocsPerRun(100, func() { sinkThakar, _ = ThakarScore(in) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkThakarScore(b *testing.B) {
	in := ThakarInput{Female: true, LVEFPercent: 30, Surgery: ThakarValveOnly, PreoperativeCreatinineMgDL: 1.5}
	b.ReportAllocs()
	for b.Loop() {
		sinkThakar, _ = ThakarScore(in)
	}
}

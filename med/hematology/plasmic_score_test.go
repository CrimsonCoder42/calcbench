package hematology

import (
	"errors"
	"math"
	"testing"
)

var plasmicBase = PLASMICFindings{PlateletsE9PerL: 150, MCVfL: 95, INR: 1.6, CreatinineMgDL: 2.5}

// Singh 2026 Table 1 (PMC13185983): each component, both sides of each threshold.
func TestPLASMICScore_components(t *testing.T) {
	cases := []struct {
		name string
		mod  func(*PLASMICFindings)
		want int
	}{
		{"none", func(f *PLASMICFindings) {}, 0},
		{"platelets 29.9", func(f *PLASMICFindings) { f.PlateletsE9PerL = 29.9 }, 1},
		{"platelets 30", func(f *PLASMICFindings) { f.PlateletsE9PerL = 30 }, 0},
		{"hemolysis", func(f *PLASMICFindings) { f.Hemolysis = true }, 1},
		{"no cancer", func(f *PLASMICFindings) { f.NoActiveCancer = true }, 1},
		{"no transplant", func(f *PLASMICFindings) { f.NoTransplant = true }, 1},
		{"MCV 89.9", func(f *PLASMICFindings) { f.MCVfL = 89.9 }, 1},
		{"MCV 90", func(f *PLASMICFindings) { f.MCVfL = 90 }, 0},
		{"INR 1.49", func(f *PLASMICFindings) { f.INR = 1.49 }, 1},
		{"INR 1.5", func(f *PLASMICFindings) { f.INR = 1.5 }, 0},
		{"creatinine 1.99", func(f *PLASMICFindings) { f.CreatinineMgDL = 1.99 }, 1},
		{"creatinine 2.0", func(f *PLASMICFindings) { f.CreatinineMgDL = 2.0 }, 0},
		{"all seven", func(f *PLASMICFindings) { *f = PLASMICFindings{12, true, true, true, 88, 1.2, 1.1} }, 7},
	}
	for _, c := range cases {
		f := plasmicBase
		c.mod(&f)
		got, err := PLASMICScore(f)
		if err != nil || got != c.want {
			t.Errorf("%s: got %d, %v; want %d", c.name, got, err, c.want)
		}
	}
}

// parity: ledger plasmic_ttp, source case 1: the consumer sums seven 0/1 inputs (no cancer and no transplant
// count 1); five met (platelets, hemolysis, no cancer, no transplant, creatinine) → 5 → intermediate risk.
func TestPLASMICScore_parity(t *testing.T) {
	f := PLASMICFindings{PlateletsE9PerL: 10, Hemolysis: true, NoActiveCancer: true, NoTransplant: true, MCVfL: 95, INR: 1.6, CreatinineMgDL: 1.0}
	got, err := PLASMICScore(f)
	if err != nil || got != 5 {
		t.Fatalf("got %d, %v; want 5", got, err)
	}
	if g, _ := PLASMICRiskGroup(got); g != PLASMICIntermediate {
		t.Errorf("group %q, want %q", g, PLASMICIntermediate)
	}
}

func TestPLASMICScore_errors(t *testing.T) {
	cases := []struct {
		name string
		mod  func(*PLASMICFindings)
		want error
	}{
		{"platelets NaN", func(f *PLASMICFindings) { f.PlateletsE9PerL = math.NaN() }, ErrNotFinite},
		{"MCV +Inf", func(f *PLASMICFindings) { f.MCVfL = math.Inf(1) }, ErrNotFinite},
		{"INR -Inf", func(f *PLASMICFindings) { f.INR = math.Inf(-1) }, ErrNotFinite},
		{"creatinine negative", func(f *PLASMICFindings) { f.CreatinineMgDL = -0.1 }, ErrOutOfRange},
		{"platelets negative", func(f *PLASMICFindings) { f.PlateletsE9PerL = -1 }, ErrOutOfRange},
	}
	for _, c := range cases {
		f := plasmicBase
		c.mod(&f)
		if _, err := PLASMICScore(f); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

// Baysal 2023 / Lee 2022: both sides of each group edge.
func TestPLASMICRiskGroup_edges(t *testing.T) {
	cases := []struct {
		s    int
		want PLASMICRisk
	}{{0, PLASMICLow}, {4, PLASMICLow}, {5, PLASMICIntermediate}, {6, PLASMICHigh}, {7, PLASMICHigh}}
	for _, c := range cases {
		if got, err := PLASMICRiskGroup(c.s); err != nil || got != c.want {
			t.Errorf("%d: got %q, %v; want %q", c.s, got, err, c.want)
		}
	}
	for _, s := range []int{-1, 8} {
		if _, err := PLASMICRiskGroup(s); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("%d: err = %v, want ErrOutOfRange", s, err)
		}
	}
}

func TestPLASMIC_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = PLASMICScore(plasmicBase) }); n != 0 {
		t.Errorf("PLASMICScore allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { g, _ := PLASMICRiskGroup(5); sinkString = string(g) }); n != 0 {
		t.Errorf("PLASMICRiskGroup allocs = %v, want 0", n)
	}
}

func BenchmarkPLASMICScore(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkInt, _ = PLASMICScore(plasmicBase)
	}
}

func BenchmarkPLASMICRiskGroup(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		g, _ := PLASMICRiskGroup(i % 8)
		sinkString = string(g)
	}
}

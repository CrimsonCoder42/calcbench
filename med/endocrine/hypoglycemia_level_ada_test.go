package endocrine

import (
	"errors"
	"math"
	"testing"
)

type hypoCase struct {
	g    float64
	want HypoglycemiaLevelADA
}

func checkHypo(t *testing.T, name string, f func(float64) (HypoglycemiaLevelADA, error), cases []hypoCase) {
	t.Helper()
	for _, c := range cases {
		got, err := f(c.g)
		if err != nil {
			t.Errorf("%s(%v): unexpected error %v", name, c.g, err)
			continue
		}
		if got != c.want {
			t.Errorf("%s(%v) = %q, want %q", name, c.g, got, c.want)
		}
	}
}

// ADA 2025 Table 6.4: Level 2 < 54 mg/dL; Level 1 < 70 and ≥ 54 mg/dL; both sides of each edge.
func TestHypoglycemiaLevelADAMgDL(t *testing.T) {
	checkHypo(t, "HypoglycemiaLevelADAMgDL", HypoglycemiaLevelADAMgDL, []hypoCase{
		{20, HypoglycemiaLevel2ADA}, {53.9, HypoglycemiaLevel2ADA}, {54, HypoglycemiaLevel1ADA},
		{69, HypoglycemiaLevel1ADA}, {69.9, HypoglycemiaLevel1ADA}, {70, HypoglycemiaNoLevelADA},
		{400, HypoglycemiaNoLevelADA},
	})
}

// ADA 2025 Table 6.4, mmol/L: Level 2 < 3.0; Level 1 < 3.9 and ≥ 3.0.
func TestHypoglycemiaLevelADAMmolL(t *testing.T) {
	checkHypo(t, "HypoglycemiaLevelADAMmolL", HypoglycemiaLevelADAMmolL, []hypoCase{
		{1.5, HypoglycemiaLevel2ADA}, {2.99, HypoglycemiaLevel2ADA}, {3.0, HypoglycemiaLevel1ADA},
		{3.8, HypoglycemiaLevel1ADA}, {3.89, HypoglycemiaLevel1ADA}, {3.9, HypoglycemiaNoLevelADA},
		{20, HypoglycemiaNoLevelADA},
	})
}

// The source's paired edges (54 ↔ 3.0, 70 ↔ 3.9) give the same level.
func TestHypoglycemiaLevelADA_twins(t *testing.T) {
	for _, p := range [][2]float64{{53, 2.9}, {54, 3.0}, {69, 3.8}, {70, 3.9}} {
		a, _ := HypoglycemiaLevelADAMgDL(p[0])
		b, _ := HypoglycemiaLevelADAMmolL(p[1])
		if a != b {
			t.Errorf("%v mg/dL → %q, %v mmol/L → %q", p[0], a, p[1], b)
		}
	}
}

func TestHypoglycemiaLevelADA_errors(t *testing.T) {
	for name, f := range map[string]func(float64) (HypoglycemiaLevelADA, error){
		"MgDL": HypoglycemiaLevelADAMgDL, "MmolL": HypoglycemiaLevelADAMmolL,
	} {
		for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
			if _, err := f(v); !errors.Is(err, ErrNotFinite) {
				t.Errorf("%s(%v): err = %v, want ErrNotFinite", name, v, err)
			}
		}
		for _, v := range []float64{0, -1} {
			if _, err := f(v); !errors.Is(err, ErrOutOfRange) {
				t.Errorf("%s(%v): err = %v, want ErrOutOfRange", name, v, err)
			}
		}
		if got, err := f(math.SmallestNonzeroFloat64); err != nil || got != HypoglycemiaLevel2ADA {
			t.Errorf("%s(smallest positive) = %q, %v; want Level 2", name, got, err)
		}
	}
}

// parity: ledger ClassifyGlucose, source cases 9–12 (53 "severe_hypoglycemia" = "Level 2 hypoglycemia",
// 54 and 69 "hypoglycemia" = "Level 1 hypoglycemia", 70 not hypoglycemia). The source classes
// hypoglycemia only for fasting and after-meal readings; Table 6.4 has no meal context.
func TestHypoglycemiaLevelADAMgDL_parityClassifier(t *testing.T) {
	checkHypo(t, "HypoglycemiaLevelADAMgDL", HypoglycemiaLevelADAMgDL, []hypoCase{
		{53, HypoglycemiaLevel2ADA}, {54, HypoglycemiaLevel1ADA}, {69, HypoglycemiaLevel1ADA},
		{70, HypoglycemiaNoLevelADA},
	})
}

// parity: ledger bs_ogtt_2h_adult, source case 3 (bands < 54 "Hypoglycemia Level 2 — Clinically
// Significant", 54–69 "Hypoglycemia Level 1 — Alert", read as [54, 70)); the primary's labels are asserted.
// Also ledger bs_fasting_adult / bs_fasting_child / bs_random_adult / bs_random_child, source case 2
// (one band < 70 "Hypoglycemia"): split by the primary into Level 2 and Level 1.
func TestHypoglycemiaLevelADAMgDL_parityBands(t *testing.T) {
	checkHypo(t, "HypoglycemiaLevelADAMgDL", HypoglycemiaLevelADAMgDL, []hypoCase{
		{53.9, HypoglycemiaLevel2ADA}, {54, HypoglycemiaLevel1ADA}, {69.5, HypoglycemiaLevel1ADA},
		{69.9, HypoglycemiaLevel1ADA}, {70, HypoglycemiaNoLevelADA},
	})
}

func TestHypoglycemiaLevelADA_noAlloc(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { _, _ = HypoglycemiaLevelADAMgDL(60) }); a != 0 {
		t.Errorf("MgDL: %v allocs, want 0", a)
	}
	if a := testing.AllocsPerRun(100, func() { _, _ = HypoglycemiaLevelADAMmolL(3.3) }); a != 0 {
		t.Errorf("MmolL: %v allocs, want 0", a)
	}
}

func BenchmarkHypoglycemiaLevelADAMgDL(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = HypoglycemiaLevelADAMgDL(60)
	}
}

func BenchmarkHypoglycemiaLevelADAMmolL(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = HypoglycemiaLevelADAMmolL(3.3)
	}
}

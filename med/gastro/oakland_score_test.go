package gastro

import (
	"math"
	"testing"
)

var oaklandBase = OaklandFindings{AgeYears: 30, HeartRateBpm: 60, SystolicMmHg: 170, HaemoglobinGL: 170}

// Oakland 2017 Table 2 (PMID 28651935): each side of each band edge; the base case (age < 40, SBP and Hb
// above 160) scores 0, below the cohort's observed minimum of 3 but reachable by the table.
func TestOaklandScore_tableBands(t *testing.T) {
	cases := []struct {
		name string
		mod  func(*OaklandFindings)
		want int
	}{
		{"base", func(f *OaklandFindings) {}, 0},
		{"age 39", func(f *OaklandFindings) { f.AgeYears = 39 }, 0},
		{"age 40", func(f *OaklandFindings) { f.AgeYears = 40 }, 1},
		{"age 69", func(f *OaklandFindings) { f.AgeYears = 69 }, 1},
		{"age 70", func(f *OaklandFindings) { f.AgeYears = 70 }, 2},
		{"male", func(f *OaklandFindings) { f.Male = true }, 1},
		{"previous LGIB", func(f *OaklandFindings) { f.PreviousLGIBAdmit = true }, 1},
		{"blood on DRE", func(f *OaklandFindings) { f.BloodOnDRE = true }, 1},
		{"HR 69", func(f *OaklandFindings) { f.HeartRateBpm = 69 }, 0},
		{"HR 70", func(f *OaklandFindings) { f.HeartRateBpm = 70 }, 1},
		{"HR 89", func(f *OaklandFindings) { f.HeartRateBpm = 89 }, 1},
		{"HR 90", func(f *OaklandFindings) { f.HeartRateBpm = 90 }, 2},
		{"HR 109", func(f *OaklandFindings) { f.HeartRateBpm = 109 }, 2},
		{"HR 110", func(f *OaklandFindings) { f.HeartRateBpm = 110 }, 3},
		{"SBP 160", func(f *OaklandFindings) { f.SystolicMmHg = 160 }, 0},
		{"SBP 159", func(f *OaklandFindings) { f.SystolicMmHg = 159 }, 2},
		{"SBP 130", func(f *OaklandFindings) { f.SystolicMmHg = 130 }, 2},
		{"SBP 129", func(f *OaklandFindings) { f.SystolicMmHg = 129 }, 3},
		{"SBP 120", func(f *OaklandFindings) { f.SystolicMmHg = 120 }, 3},
		{"SBP 119", func(f *OaklandFindings) { f.SystolicMmHg = 119 }, 4},
		{"SBP 90", func(f *OaklandFindings) { f.SystolicMmHg = 90 }, 4},
		{"SBP 89", func(f *OaklandFindings) { f.SystolicMmHg = 89 }, 5},
		{"SBP 50", func(f *OaklandFindings) { f.SystolicMmHg = 50 }, 5},
		{"Hb 160", func(f *OaklandFindings) { f.HaemoglobinGL = 160 }, 0},
		{"Hb 159", func(f *OaklandFindings) { f.HaemoglobinGL = 159 }, 4},
		{"Hb 130", func(f *OaklandFindings) { f.HaemoglobinGL = 130 }, 4},
		{"Hb 129", func(f *OaklandFindings) { f.HaemoglobinGL = 129 }, 8},
		{"Hb 110", func(f *OaklandFindings) { f.HaemoglobinGL = 110 }, 8},
		{"Hb 109", func(f *OaklandFindings) { f.HaemoglobinGL = 109 }, 13},
		{"Hb 90", func(f *OaklandFindings) { f.HaemoglobinGL = 90 }, 13},
		{"Hb 89", func(f *OaklandFindings) { f.HaemoglobinGL = 89 }, 17},
		{"Hb 70", func(f *OaklandFindings) { f.HaemoglobinGL = 70 }, 17},
		{"Hb 69", func(f *OaklandFindings) { f.HaemoglobinGL = 69 }, 22},
		{"Hb 36", func(f *OaklandFindings) { f.HaemoglobinGL = 36 }, 22},
		{"doc example", func(f *OaklandFindings) {
			*f = OaklandFindings{AgeYears: 75, Male: true, HeartRateBpm: 95, SystolicMmHg: 125, HaemoglobinGL: 100}
		}, 21},
		{"maximum 35 (every top band of Table 2: 2+1+1+1+3+5+22)", func(f *OaklandFindings) {
			*f = OaklandFindings{AgeYears: 80, Male: true, PreviousLGIBAdmit: true, BloodOnDRE: true, HeartRateBpm: 120, SystolicMmHg: 60, HaemoglobinGL: 50}
		}, 35},
	}
	for _, c := range cases {
		f := oaklandBase
		c.mod(&f)
		got, err := OaklandScore(f)
		if err != nil || got != c.want {
			t.Errorf("%s: got %d, %v; want %d", c.name, got, err, c.want)
		}
	}
}

// Oakland 2017 text: "A patient scoring ≤8 points at presentation has a 95% chance of safe discharge"; a
// woman of 45 with HR 75, SBP 135, Hb 140 g/L → 1 + 1 + 2 + 4 = 8.
func TestOaklandScore_sourceThreshold(t *testing.T) {
	f := OaklandFindings{AgeYears: 45, HeartRateBpm: 75, SystolicMmHg: 135, HaemoglobinGL: 140}
	if got, err := OaklandScore(f); err != nil || got != 8 {
		t.Errorf("got %d, %v; want 8", got, err)
	}
}

// parity: ledger oakland_score, source case 1: the consumer sums picker indexes with SBP points [5,4,3,2,0]
// and Hb points [22,17,13,8,4,0]; age 40–69 (1), male (1), HR 90–109 (2), SBP 90–119 (4), Hb 9–10.9 g/dL
// (13) → 21.
func TestOaklandScore_parity(t *testing.T) {
	f := OaklandFindings{AgeYears: 55, Male: true, HeartRateBpm: 100, SystolicMmHg: 100, HaemoglobinGL: 95}
	if got, err := OaklandScore(f); err != nil || got != 21 {
		t.Errorf("got %d, %v; want 21", got, err)
	}
}

func TestOaklandScore_errors(t *testing.T) {
	bad := func(mod func(*OaklandFindings)) func() error {
		return func() error { f := oaklandBase; mod(&f); _, err := OaklandScore(f); return err }
	}
	runErrCases(t, []errCase{
		{"age NaN", bad(func(f *OaklandFindings) { f.AgeYears = math.NaN() }), ErrNotFinite},
		{"HR +Inf", bad(func(f *OaklandFindings) { f.HeartRateBpm = math.Inf(1) }), ErrNotFinite},
		{"SBP -Inf", bad(func(f *OaklandFindings) { f.SystolicMmHg = math.Inf(-1) }), ErrNotFinite},
		{"Hb NaN", bad(func(f *OaklandFindings) { f.HaemoglobinGL = math.NaN() }), ErrNotFinite},
		{"age 15.9", bad(func(f *OaklandFindings) { f.AgeYears = 15.9 }), ErrOutOfRange},
		{"HR negative", bad(func(f *OaklandFindings) { f.HeartRateBpm = -1 }), ErrOutOfRange},
		{"SBP 49", bad(func(f *OaklandFindings) { f.SystolicMmHg = 49 }), ErrOutOfRange},
		{"Hb 35", bad(func(f *OaklandFindings) { f.HaemoglobinGL = 35 }), ErrOutOfRange},
	})
}

func TestOaklandScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = OaklandScore(oaklandBase) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkOaklandScore(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkInt, _ = OaklandScore(oaklandBase)
	}
}

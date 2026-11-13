package respiratory

import (
	"errors"
	"math"
	"testing"
)

func resectExample() RESECT90Findings {
	return RESECT90Findings{AgeYears: 70, Male: true, PerformanceStatus: 1, DLCOPercentPred: 70, BMIKgM2: 26,
		CreatinineUmolL: 90, RightSided: true, SegmentsResected: 3, Thoracotomy: true, Malignant: true}
}

// HAND-DERIVED from Taylor 2021's equation (PMID 34324664): Logit = −6.036 + 2.87 + 0.493 + 0.183 − 2.03 −
// 1.456 + 0.45 + 0.379 + 0.537 + 0.634 + 0.769 = −3.207; probability 1/(1 + e^3.207). Tolerance 1e-12
// absolute (the same arithmetic in another order).
func TestRESECT90MortalityProbability_equationExample(t *testing.T) {
	got, err := RESECT90MortalityProbability(resectExample())
	want := 1 / (1 + math.Exp(3.207))
	if err != nil || math.Abs(got-want) > 1e-12 {
		t.Errorf("got (%v, %v), want %v", got, err, want)
	}
	if math.Abs(got-0.0389) > 5e-5 {
		t.Errorf("got %v, want ≈ 0.0389", got)
	}
}

// Each binary term moves the logit by its coefficient: odds ratio = exp(β). Tolerance 1e-9 relative.
func TestRESECT90MortalityProbability_eachTerm(t *testing.T) {
	base := RESECT90Findings{AgeYears: 60, DLCOPercentPred: 80, BMIKgM2: 25, CreatinineUmolL: 80, SegmentsResected: 1}
	odds := func(f RESECT90Findings) float64 {
		p, err := RESECT90MortalityProbability(f)
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}
		return p / (1 - p)
	}
	o0 := odds(base)
	terms := []struct {
		name string
		mod  func(*RESECT90Findings)
		beta float64
	}{
		{"male", func(f *RESECT90Findings) { f.Male = true }, 0.493},
		{"anaemia", func(f *RESECT90Findings) { f.Anaemia = true }, 0.242},
		{"arrhythmia", func(f *RESECT90Findings) { f.Arrhythmia = true }, 0.608},
		{"right", func(f *RESECT90Findings) { f.RightSided = true }, 0.379},
		{"thoracotomy", func(f *RESECT90Findings) { f.Thoracotomy = true }, 0.634},
		{"malignant", func(f *RESECT90Findings) { f.Malignant = true }, 0.769},
		{"PS +1", func(f *RESECT90Findings) { f.PerformanceStatus = 1 }, 0.183},
		{"segment +1", func(f *RESECT90Findings) { f.SegmentsResected = 2 }, 0.179},
		{"age +1", func(f *RESECT90Findings) { f.AgeYears = 61 }, 0.041},
		{"DLCO +1", func(f *RESECT90Findings) { f.DLCOPercentPred = 81 }, -0.029},
		{"BMI +1", func(f *RESECT90Findings) { f.BMIKgM2 = 26 }, -0.056},
		{"creatinine +1", func(f *RESECT90Findings) { f.CreatinineUmolL = 81 }, 0.005},
	}
	for _, c := range terms {
		f := base
		c.mod(&f)
		if r := odds(f) / o0; math.Abs(r-math.Exp(c.beta)) > 1e-9*math.Exp(c.beta) {
			t.Errorf("%s: odds ratio %v, want %v", c.name, r, math.Exp(c.beta))
		}
	}
}

// parity: ledger resect_90, source case 1: the source takes creatinine in mg/dL and multiplies by 88.4, then
// uses the paper's coefficients; creatinine 1.0 mg/dL = 88.4 µmol/L, the rest as resectExample →
// Logit −3.215, the same as this function at 88.4 µmol/L. Tolerance 1e-12.
func TestRESECT90MortalityProbability_parity(t *testing.T) {
	f := resectExample()
	f.CreatinineUmolL = 88.4
	got, err := RESECT90MortalityProbability(f)
	want := 1 / (1 + math.Exp(3.207+0.005*1.6))
	if err != nil || math.Abs(got-want) > 1e-12 {
		t.Errorf("got (%v, %v), want %v", got, err, want)
	}
}

func TestRESECT90MortalityProbability_errors(t *testing.T) {
	cases := []struct {
		name string
		mod  func(*RESECT90Findings)
		want error
	}{
		{"age 17", func(f *RESECT90Findings) { f.AgeYears = 17 }, ErrOutOfRange},
		{"age 121", func(f *RESECT90Findings) { f.AgeYears = 121 }, ErrOutOfRange},
		{"PS 5", func(f *RESECT90Findings) { f.PerformanceStatus = 5 }, ErrOutOfRange},
		{"PS -1", func(f *RESECT90Findings) { f.PerformanceStatus = -1 }, ErrOutOfRange},
		{"DLCO 0", func(f *RESECT90Findings) { f.DLCOPercentPred = 0 }, ErrOutOfRange},
		{"BMI 0", func(f *RESECT90Findings) { f.BMIKgM2 = 0 }, ErrOutOfRange},
		{"creatinine 0", func(f *RESECT90Findings) { f.CreatinineUmolL = 0 }, ErrOutOfRange},
		{"segments 20", func(f *RESECT90Findings) { f.SegmentsResected = 20 }, ErrOutOfRange},
		{"segments -1", func(f *RESECT90Findings) { f.SegmentsResected = -1 }, ErrOutOfRange},
		{"age NaN", func(f *RESECT90Findings) { f.AgeYears = math.NaN() }, ErrNotFinite},
		{"DLCO +Inf", func(f *RESECT90Findings) { f.DLCOPercentPred = math.Inf(1) }, ErrNotFinite},
		{"BMI -Inf", func(f *RESECT90Findings) { f.BMIKgM2 = math.Inf(-1) }, ErrNotFinite},
		{"creatinine NaN", func(f *RESECT90Findings) { f.CreatinineUmolL = math.NaN() }, ErrNotFinite},
	}
	for _, c := range cases {
		f := resectExample()
		c.mod(&f)
		if _, err := RESECT90MortalityProbability(f); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

func TestRESECT90MortalityProbability_noAlloc(t *testing.T) {
	f := resectExample()
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = RESECT90MortalityProbability(f) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkRESECT90MortalityProbability(b *testing.B) {
	f := resectExample()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkFloat, _ = RESECT90MortalityProbability(f)
	}
}

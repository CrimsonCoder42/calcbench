package endocrine

import (
	"errors"
	"math"
	"testing"
)

type glycemicCase struct {
	v    float64
	want GlycemicCategoryADA
}

func checkGlycemic(t *testing.T, name string, f func(float64) (GlycemicCategoryADA, error), cases []glycemicCase) {
	t.Helper()
	for _, c := range cases {
		got, err := f(c.v)
		if err != nil {
			t.Errorf("%s(%v): unexpected error %v", name, c.v, err)
			continue
		}
		if got != c.want {
			t.Errorf("%s(%v) = %q, want %q", name, c.v, got, c.want)
		}
	}
}

// checkGlycemicErrors: NaN, ±Inf → ErrNotFinite; 0 and negative → ErrOutOfRange; smallest positive is accepted.
func checkGlycemicErrors(t *testing.T, name string, f func(float64) (GlycemicCategoryADA, error)) {
	t.Helper()
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		if _, err := f(v); !errors.Is(err, ErrNotFinite) {
			t.Errorf("%s(%v): err = %v, want ErrNotFinite", name, v, err)
		}
	}
	for _, v := range []float64{0, -1, math.Copysign(0, -1)} {
		if _, err := f(v); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("%s(%v): err = %v, want ErrOutOfRange", name, v, err)
		}
	}
	if got, err := f(math.SmallestNonzeroFloat64); err != nil || got != GlycemicNoCategoryADA {
		t.Errorf("%s(smallest positive) = %q, %v; want \"\", nil", name, got, err)
	}
	if got, err := f(math.MaxFloat64); err != nil || got != GlycemicDiabetesADA {
		t.Errorf("%s(MaxFloat64) = %q, %v; want Diabetes, nil", name, got, err)
	}
}

func checkGlycemicNoAlloc(t *testing.T, name string, f func(float64) (GlycemicCategoryADA, error), v float64) {
	t.Helper()
	if a := testing.AllocsPerRun(100, func() { _, _ = f(v) }); a != 0 {
		t.Errorf("%s: %v allocs per run, want 0", name, a)
	}
}

// ADA 2026 Tables 2.1 and 2.2: A1C 5.7–6.4 % Prediabetes, ≥ 6.5 % Diabetes; both sides of each edge and
// a value between the printed 6.4 and 6.5.
func TestA1CCategoryADAPercent(t *testing.T) {
	checkGlycemic(t, "A1CCategoryADAPercent", A1CCategoryADAPercent, []glycemicCase{
		{4.0, GlycemicNoCategoryADA}, {5.6, GlycemicNoCategoryADA}, {5.69, GlycemicNoCategoryADA},
		{5.7, GlycemicPrediabetesADA}, {6.4, GlycemicPrediabetesADA}, {6.45, GlycemicPrediabetesADA},
		{6.49, GlycemicPrediabetesADA}, {6.5, GlycemicDiabetesADA}, {12, GlycemicDiabetesADA},
	})
	checkGlycemicErrors(t, "A1CCategoryADAPercent", A1CCategoryADAPercent)
}

// ADA 2026 Tables 2.1 and 2.2, mmol/mol column: 39–47 Prediabetes, ≥ 48 Diabetes.
func TestA1CCategoryADAMmolMol(t *testing.T) {
	checkGlycemic(t, "A1CCategoryADAMmolMol", A1CCategoryADAMmolMol, []glycemicCase{
		{20, GlycemicNoCategoryADA}, {38, GlycemicNoCategoryADA}, {38.9, GlycemicNoCategoryADA},
		{39, GlycemicPrediabetesADA}, {47, GlycemicPrediabetesADA}, {47.5, GlycemicPrediabetesADA},
		{47.9, GlycemicPrediabetesADA}, {48, GlycemicDiabetesADA}, {100, GlycemicDiabetesADA},
	})
	checkGlycemicErrors(t, "A1CCategoryADAMmolMol", A1CCategoryADAMmolMol)
}

// The source's two unit columns agree on the category at the printed paired values (5.7 % ↔ 39, 6.4 ↔ 47,
// 6.5 ↔ 48 mmol/mol).
func TestA1CCategoryADA_twins(t *testing.T) {
	pairs := [][2]float64{{5.6, 38}, {5.7, 39}, {6.4, 47}, {6.5, 48}}
	for _, p := range pairs {
		a, _ := A1CCategoryADAPercent(p[0])
		b, _ := A1CCategoryADAMmolMol(p[1])
		if a != b {
			t.Errorf("%v %% → %q, %v mmol/mol → %q", p[0], a, p[1], b)
		}
	}
}

// parity: ledger bs_a1c_adult, source case 1 (band edges: below 5.7 the source's "Normal A1C", 5.7–6.5
// "Prediabetes (A1C)", ≥ 6.5 "Diabetes (A1C)"); the primary's labels are asserted.
func TestA1CCategoryADAPercent_parity(t *testing.T) {
	checkGlycemic(t, "A1CCategoryADAPercent", A1CCategoryADAPercent, []glycemicCase{
		{5.69, GlycemicNoCategoryADA}, // source "Normal A1C"; the primary names no category
		{5.7, GlycemicPrediabetesADA}, // source "Prediabetes (A1C)"
		{6.49, GlycemicPrediabetesADA},
		{6.5, GlycemicDiabetesADA}, // source "Diabetes (A1C)"
		{25, GlycemicDiabetesADA},  // source input max
	})
}

func TestA1CCategoryADA_noAlloc(t *testing.T) {
	checkGlycemicNoAlloc(t, "A1CCategoryADAPercent", A1CCategoryADAPercent, 6.0)
	checkGlycemicNoAlloc(t, "A1CCategoryADAMmolMol", A1CCategoryADAMmolMol, 42)
}

func BenchmarkA1CCategoryADAPercent(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = A1CCategoryADAPercent(6.0)
	}
}

func BenchmarkA1CCategoryADAMmolMol(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, _ = A1CCategoryADAMmolMol(42)
	}
}

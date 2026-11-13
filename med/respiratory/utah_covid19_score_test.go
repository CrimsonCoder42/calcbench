package respiratory

import (
	"errors"
	"testing"
)

// HAND-DERIVED from Webb 2022 Table 6 (PMID 35239664): 65-year-old man, diabetes, hypertension →
// 3.5 + 1 + 2 + 1 = 7.5; every item at age > 100 → 5.5 + 1 + 2 + 3 × 2 + 9 + 1 = 24.5. Exact comparison
// (multiples of 0.5 are exact in float64).
func TestUtahCOVID19Score_tableExample(t *testing.T) {
	f := UtahCOVID19Findings{AgeYears: 65, Male: true, DiabetesMellitus: true, Hypertension: true}
	if got, err := UtahCOVID19Score(f); err != nil || got != 7.5 {
		t.Errorf("got (%v, %v), want 7.5", got, err)
	}
	all := UtahCOVID19Findings{105, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true}
	if got, err := UtahCOVID19Score(all); err != nil || got != 24.5 {
		t.Errorf("all: got (%v, %v), want 24.5", got, err)
	}
}

// Age bands of Table 6, both sides of each decade edge.
func TestUtahCOVID19Score_ageBands(t *testing.T) {
	cases := []struct {
		age  int
		want float64
	}{{0, 0.5}, {10, 0.5}, {11, 1}, {20, 1}, {21, 1.5}, {30, 1.5}, {31, 2}, {50, 2.5}, {51, 3}, {70, 3.5}, {71, 4},
		{90, 4.5}, {91, 5}, {100, 5}, {101, 5.5}, {120, 5.5}}
	for _, c := range cases {
		if got, err := UtahCOVID19Score(UtahCOVID19Findings{AgeYears: c.age}); err != nil || got != c.want {
			t.Errorf("age %d: got (%v, %v), want %v", c.age, got, err, c.want)
		}
	}
}

// Item weights: 2 for communities of color and the three high-risk comorbidities, 1 for male, each other
// comorbidity and dyspnea.
func TestUtahCOVID19Score_items(t *testing.T) {
	base := UtahCOVID19Findings{AgeYears: 40} // 2
	cases := []struct {
		name string
		mod  func(*UtahCOVID19Findings)
		want float64
	}{
		{"male", func(f *UtahCOVID19Findings) { f.Male = true }, 3},
		{"communities of color", func(f *UtahCOVID19Findings) { f.CommunityOfColor = true }, 4},
		{"immunocompromised", func(f *UtahCOVID19Findings) { f.SeverelyImmunocompromised = true }, 4},
		{"obesity", func(f *UtahCOVID19Findings) { f.ObesityBMIOver30 = true }, 4},
		{"CKD and liver", func(f *UtahCOVID19Findings) { f.ChronicKidneyDisease, f.ChronicLiverDisease = true, true }, 4},
		{"dyspnea", func(f *UtahCOVID19Findings) { f.Dyspnea = true }, 3},
	}
	for _, c := range cases {
		f := base
		c.mod(&f)
		if got, err := UtahCOVID19Score(f); err != nil || got != c.want {
			t.Errorf("%s: got (%v, %v), want %v", c.name, got, err, c.want)
		}
	}
}

// parity: ledger utah_covid_risk, source case 1: 65-year-old woman, diabetes, hypertension → 6.5 in both.
// Source case 2: 45-year-old woman, immunosuppression → the source gives 2.5 + 1 = 3.5 (immunosuppression as an
// "other" comorbidity); Table 6 lists "Severely Immunocompromised" with 2 points. Asserting the PRIMARY: 4.5.
// The source also omits "Male" (1) and "Communities of color" (2), and scores ages 0–10 as 1.0 (Table 6: 0.5).
func TestUtahCOVID19Score_parity(t *testing.T) {
	if got, _ := UtahCOVID19Score(UtahCOVID19Findings{AgeYears: 65, DiabetesMellitus: true, Hypertension: true}); got != 6.5 {
		t.Errorf("case 1: got %v, want 6.5", got)
	}
	if got, _ := UtahCOVID19Score(UtahCOVID19Findings{AgeYears: 45, SeverelyImmunocompromised: true}); got != 4.5 {
		t.Errorf("case 2: got %v, want 4.5 (source 3.5)", got)
	}
}

func TestUtahCOVID19Score_errors(t *testing.T) {
	for _, age := range []int{-1, 121} {
		if _, err := UtahCOVID19Score(UtahCOVID19Findings{AgeYears: age}); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("age %d: err = %v, want ErrOutOfRange", age, err)
		}
	}
}

func TestUtahCOVID19Score_noAlloc(t *testing.T) {
	f := UtahCOVID19Findings{AgeYears: 65, Male: true, DiabetesMellitus: true}
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = UtahCOVID19Score(f) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkUtahCOVID19Score(b *testing.B) {
	f := UtahCOVID19Findings{AgeYears: 65, Male: true, DiabetesMellitus: true}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sinkFloat, _ = UtahCOVID19Score(f)
	}
}

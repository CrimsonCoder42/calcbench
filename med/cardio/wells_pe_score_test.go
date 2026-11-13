package cardio

import "testing"

// Wells 2000 abstract points; maximum 3 + 3 + 1.5 × 3 + 1 + 1 = 12.5. Halves are exact in float64.
func TestWellsPEScore_items(t *testing.T) {
	tests := []struct {
		name string
		f    WellsPEFindings
		want float64
	}{
		{"none", WellsPEFindings{}, 0},
		{"example: HR > 100 + hemoptysis", WellsPEFindings{HeartRateAbove100: true, Hemoptysis: true}, 2.5},
		{"two-level edge 4.0", WellsPEFindings{NoAlternativeDiagnosis: true, Malignancy: true}, 4},
		{"4.5", WellsPEFindings{ClinicalSignsOfDVT: true, PreviousDVTOrPE: true}, 4.5},
		{"maximum", WellsPEFindings{true, true, true, true, true, true, true}, 12.5},
	}
	for _, tc := range tests {
		if got := WellsPEScore(tc.f); got != tc.want {
			t.Errorf("%s: %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestWellsPEScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat = WellsPEScore(WellsPEFindings{Hemoptysis: true}) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkWellsPEScore(b *testing.B) {
	b.ReportAllocs()
	f := WellsPEFindings{HeartRateAbove100: true}
	for b.Loop() {
		sinkFloat = WellsPEScore(f)
	}
}

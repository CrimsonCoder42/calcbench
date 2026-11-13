package cardio

import (
	"math"
	"testing"
)

// Park 2016, Table 5: constant 4, age bands, item points. HAND-DERIVED totals; age edges both sides.
func TestIEMortalityScore_items(t *testing.T) {
	tests := []struct {
		name string
		age  float64
		f    IEMortalityFactors
		want int
	}{
		{"45", 45, IEMortalityFactors{}, 4},
		{"45.5", 45.5, IEMortalityFactors{}, 6},
		{"60", 60, IEMortalityFactors{}, 6},
		{"61", 61, IEMortalityFactors{}, 7},
		{"70", 70, IEMortalityFactors{}, 7},
		{"71", 71, IEMortalityFactors{}, 8},
		{"example", 65, IEMortalityFactors{StaphylococcusAureus: true, MitralVegetation: true, Stroke: true}, 11},
		{"negatives", 30, IEMortalityFactors{SymptomsOver1Month: true, ViridansStreptococci: true, SurgicalTreatment: true}, -1},
	}
	for _, tc := range tests {
		if got, err := IEMortalityScore(tc.age, tc.f); err != nil || got != tc.want {
			t.Errorf("%s: (%d, %v), want %d", tc.name, got, err, tc.want)
		}
	}
	_, err := IEMortalityScore(math.NaN(), IEMortalityFactors{})
	checkErr(t, "NaN", err, ErrNotFinite)
	_, err = IEMortalityScore(17.999, IEMortalityFactors{})
	checkErr(t, "17.999 (below the adult population)", err, ErrOutOfRange)
	// No upper age is stated by the source (A-M1): 101 years scores the > 70 band, 4 + 4 = 8.
	if got, err := IEMortalityScore(101, IEMortalityFactors{}); err != nil || got != 8 {
		t.Errorf("IEMortalityScore(101) = (%d, %v), want (8, nil)", got, err)
	}
}

func TestIEMortalityScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = IEMortalityScore(65, IEMortalityFactors{Stroke: true}) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkIEMortalityScore(b *testing.B) {
	b.ReportAllocs()
	f := IEMortalityFactors{Stroke: true}
	for b.Loop() {
		sinkInt, _ = IEMortalityScore(65, f)
	}
}

package stats

import (
	"errors"
	"fmt"
	"math"
	"testing"
)

// skewedSample has n = 10, s = 2.4967, IQR (type 7) = 1.75, g₁ = 1.6689; the class counts below were
// computed by hand from each paper's formula (Python check, recorded in the work folder).
var skewedSample = []float64{1, 2, 2, 3, 3, 3, 4, 4, 5, 10}

func TestHistogramBinEdges_rules(t *testing.T) {
	for _, tt := range []struct {
		name string
		rule BinRule
		k    int
	}{
		{"Sturges ⌈1 + log₂ 10⌉", BinSturges, 5},
		{"Scott ⌈9/4.0453⌉", BinScott, 3},
		{"Freedman–Diaconis ⌈9/1.6245⌉", BinFreedmanDiaconis, 6},
		{"Doane ⌈6.278⌉", BinDoane, 7},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HistogramBinEdges(skewedSample, tt.rule)
			if err != nil || len(got) != tt.k+1 || got[0] != 1 || got[tt.k] != 10 {
				t.Fatalf("HistogramBinEdges = %v, %v; want %d classes over 1..10", got, err, tt.k)
			}
			for i := range got {
				// Tolerance 1e-14 relative: one lerp per edge.
				if want := 1 + 9*float64(i)/float64(tt.k); !near(got[i], want, 1e-14) {
					t.Errorf("edge %d = %v; want %v", i, got[i], want)
				}
			}
		})
	}
}

// Worked example: 1 … 8 by Sturges → 4 classes, edges 1, 2.75, 4.5, 6.25, 8.
func TestHistogramBinEdges_sturgesWorked(t *testing.T) {
	got, err := HistogramBinEdges([]float64{1, 2, 3, 4, 5, 6, 7, 8}, BinSturges)
	want := []float64{1, 2.75, 4.5, 6.25, 8}
	if err != nil || len(got) != len(want) {
		t.Fatalf("got %v, %v; want %v", got, err, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("edge %d = %v; want %v", i, got[i], want[i])
		}
	}
	if _, err := FrequencyDistribution([]float64{1, 2, 3, 4, 5, 6, 7, 8}, got); err != nil {
		t.Errorf("edges do not cover the data: %v", err)
	}
}

func TestHistogramBinEdges_huge(t *testing.T) {
	m := math.MaxFloat64
	for _, rule := range []BinRule{BinSturges, BinFreedmanDiaconis, BinDoane} {
		got, err := HistogramBinEdges([]float64{-m, 0, 0, m}, rule)
		if err != nil || got[0] != -m || got[len(got)-1] != m {
			t.Errorf("rule %d = %v, %v", rule, got, err)
		}
	}
}

func TestHistogramBinEdges_errors(t *testing.T) {
	m := math.MaxFloat64
	for _, tt := range []struct {
		name   string
		values []float64
		rule   BinRule
		want   error
	}{
		{"bad rule", []float64{1, 2}, 0, ErrOutOfRange},
		{"bad rule high", []float64{1, 2}, 5, ErrOutOfRange},
		{"empty", nil, BinSturges, ErrEmpty},
		{"one", []float64{1}, BinScott, ErrTooFew},
		{"two Doane", []float64{1, 2}, BinDoane, ErrTooFew},
		{"NaN", []float64{1, math.NaN()}, BinSturges, ErrNotFinite},
		{"all equal", []float64{3, 3, 3}, BinSturges, ErrZeroVariance},
		{"zero IQR", []float64{0, 5, 5, 5, 5, 5, 5, 9}, BinFreedmanDiaconis, ErrZeroVariance},
		{"too many classes", outlierSample(), BinFreedmanDiaconis, ErrOutOfRange},
		{"too narrow", []float64{1, math.Nextafter(1, 2), 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}, BinDoane, ErrOutOfRange},
		{"Scott SD overflow", []float64{-m, m}, BinScott, ErrOverflow},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := HistogramBinEdges(tt.values, tt.rule); !errors.Is(err, tt.want) {
				t.Errorf("err = %v; want %v", err, tt.want)
			}
		})
	}
}

func BenchmarkHistogramBinEdges(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		in := benchInputs(n)
		for _, rule := range []BinRule{BinSturges, BinScott, BinFreedmanDiaconis, BinDoane} {
			b.Run(fmt.Sprintf("rule=%d/%d", rule, n), func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					_, _ = HistogramBinEdges(in[0], rule)
				}
			})
		}
	}
}

// outlierSample is 0, 1e-6, …, 999e-6 and one value 1e12: IQR ≈ 5e-4, so Freedman–Diaconis asks for
// about 1e16 classes.
func outlierSample() []float64 {
	v := make([]float64, 1001)
	for i := range 1000 {
		v[i] = float64(i) * 1e-6
	}
	v[1000] = 1e12
	return v
}

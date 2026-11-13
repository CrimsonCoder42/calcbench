package invest

import (
	"fmt"
	"math"
	"testing"
)

func TestInformationRatio_handWorked(t *testing.T) {
	got, err := InformationRatio([]float64{0.02, 0.03, 0.01}, []float64{0.01, 0.01, 0.01})
	if err != nil || !near(got, 1, 1e-14) {
		t.Errorf("InformationRatio = %v, %v; want 1", got, err)
	}
}

// Active returns 1, 2, 3, 4: mean 2.5, s = √(5/3) → 2.5/√(5/3).
func TestInformationRatio_cases(t *testing.T) {
	got, err := InformationRatio([]float64{1, 2, 3, 4}, []float64{0, 0, 0, 0})
	if want := 2.5 / math.Sqrt(5.0/3); err != nil || !near(got, want, 1e-15) {
		t.Errorf("1..4: %v, %v; want %v", got, err, want)
	}
	got, err = InformationRatio([]float64{0, 0, 0, 0}, []float64{1, 2, 3, 4})
	if want := -2.5 / math.Sqrt(5.0/3); err != nil || !near(got, want, 1e-15) {
		t.Errorf("underperforming: %v, %v; want %v", got, err, want)
	}
}

func TestInformationRatio_errors(t *testing.T) {
	_, err := InformationRatio([]float64{0.5, 0.75}, []float64{0.25, 0.5})
	wantErr(t, "zero TE", err, ErrZeroDenominator)
	_, err = InformationRatio([]float64{0.01}, []float64{0.01})
	wantErr(t, "one", err, ErrTooFew)
	_, err = InformationRatio([]float64{0.01, 0.02}, []float64{0.01, 0.02, 0.03})
	wantErr(t, "lengths", err, ErrLengthMismatch)
	_, err = InformationRatio([]float64{math.NaN(), 0.02}, []float64{0.01, 0.02})
	wantErr(t, "NaN", err, ErrNotFinite)
}

func TestInformationRatio_noAlloc(t *testing.T) {
	p, bm := []float64{0.02, 0.03, 0.01}, []float64{0.01, 0.01, 0.01}
	if a := testing.AllocsPerRun(100, func() { sinkFloat, sinkErr = InformationRatio(p, bm) }); a != 0 {
		t.Errorf("allocs = %v, want 0", a)
	}
}

func BenchmarkInformationRatio(b *testing.B) {
	for _, n := range benchSizes {
		p := genReturns(n)
		bm := genReturns(n + 1)[1:]
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, sinkErr = InformationRatio(p, bm)
			}
		})
	}
}

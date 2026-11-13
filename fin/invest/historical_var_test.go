package invest

import (
	"fmt"
	"math"
	"testing"

	"github.com/CrimsonCoder42/calcbench/num/stats"
)

// tailReturns sorted: −0.05, −0.03, −0.01, 0, 0.01, 0.02, 0.02, 0.03, 0.04, 0.05. Confidence 0.8 → p = 0.2,
// np = 2. Type 1: j = 2, g = 0 → X₍₂₎ = −0.03. Type 7: np + 1 − p = 2.8 → 0.2·X₍₂₎ + 0.8·X₍₃₎ = −0.014.
func TestHistoricalVaR_handWorked(t *testing.T) {
	got, err := HistoricalVaR(tailReturns, 0.8, stats.HyndmanFan1)
	if err != nil || !near(got, 0.03, 1e-15) {
		t.Errorf("type 1: %v, %v; want 0.03", got, err)
	}
	got, err = HistoricalVaR(tailReturns, 0.8, stats.HyndmanFan7)
	if err != nil || !near(got, 0.014, 1e-15) {
		t.Errorf("type 7: %v, %v; want 0.014", got, err)
	}
}

func TestHistoricalVaR_cases(t *testing.T) {
	got, err := HistoricalVaR([]float64{0.01}, 0.99, stats.HyndmanFan1)
	if err != nil || got != -0.01 {
		t.Errorf("one gain: %v, %v; want -0.01", got, err)
	}
	got, err = HistoricalVaR([]float64{0.03, -0.02, 0.01}, 0.99, stats.HyndmanFan1)
	if err != nil || got != 0.02 {
		t.Errorf("worst of three: %v, %v; want 0.02", got, err)
	}
}

func TestHistoricalVaR_errors(t *testing.T) {
	_, err := HistoricalVaR(nil, 0.95, stats.HyndmanFan1)
	wantErr(t, "empty", err, ErrTooFew)
	_, err = HistoricalVaR([]float64{math.NaN()}, 0.95, stats.HyndmanFan1)
	wantErr(t, "NaN", err, ErrNotFinite)
	for _, bad := range []float64{0, 1, -0.1, 1.1, math.NaN()} {
		_, err = HistoricalVaR(tailReturns, bad, stats.HyndmanFan1)
		wantErr(t, "bad confidence", err, ErrDomain)
	}
	for _, bad := range []stats.QuantileMethod{0, 10} {
		_, err = HistoricalVaR(tailReturns, 0.95, bad)
		wantErr(t, "bad method", err, ErrDomain)
	}
}

func BenchmarkHistoricalVaR(b *testing.B) {
	for _, n := range benchSizes {
		r := genReturns(n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				sinkFloat, sinkErr = HistoricalVaR(r, 0.99, stats.HyndmanFan1)
			}
		})
	}
}

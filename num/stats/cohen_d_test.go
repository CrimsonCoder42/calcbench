package stats

import (
	"errors"
	"math"
	"testing"
)

// Worked example: Cohen 1988 eq. 2.2.1, d = (m_A − m_B)/σ, applied to the e-Handbook §1.3.5.3 AUTO83B
// means 30.48101 and 20.14458 with the printed pooled s_p = 6.34260 → 1.62968 (hand-checked).
func TestCohenD_worked(t *testing.T) {
	// Tolerance 1e-5 absolute: the inputs are printed to 5 decimals.
	if got, err := CohenD(30.48101, 20.14458, 6.34260); err != nil || math.Abs(got-1.62968) > 1e-5 {
		t.Errorf("CohenD(AUTO83B) = %v, %v; want 1.62968", got, err)
	}
}

func TestCohenD_cases(t *testing.T) {
	m := math.MaxFloat64
	for _, tt := range []struct {
		name             string
		meanA, meanB, sd float64
		want             float64
	}{
		{"equal means", 5, 5, 2, 0},
		{"one SD apart", 12, 10, 2, 1},
		{"negative", 10, 12, 4, -0.5},
		{"difference overflows, d finite", m, -m, m, 2},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := CohenD(tt.meanA, tt.meanB, tt.sd); err != nil || got != tt.want {
				t.Errorf("CohenD(%v, %v, %v) = %v, %v; want %v", tt.meanA, tt.meanB, tt.sd, got, err, tt.want)
			}
		})
	}
}

func TestCohenD_errors(t *testing.T) {
	for _, tt := range []struct {
		name             string
		meanA, meanB, sd float64
		want             error
	}{
		{"NaN meanA", math.NaN(), 0, 1, ErrNotFinite},
		{"Inf meanB", 0, math.Inf(-1), 1, ErrNotFinite},
		{"NaN sd", 0, 0, math.NaN(), ErrNotFinite},
		{"negative sd", 1, 0, -1, ErrNegative},
		{"zero sd", 1, 0, 0, ErrZeroVariance},
		{"overflow", math.MaxFloat64, 0, 0.5, ErrOverflow},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := CohenD(tt.meanA, tt.meanB, tt.sd); !errors.Is(err, tt.want) || got != 0 {
				t.Errorf("CohenD = %v, %v; want 0, %v", got, err, tt.want)
			}
		})
	}
}

func TestCohenD_noAllocation(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { benchSink, _ = CohenD(30.48101, 20.14458, 6.34260) }); a != 0 {
		t.Errorf("CohenD allocates %v times per call, want 0", a)
	}
}

func BenchmarkCohenD(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchSink, _ = CohenD(30.48101, 20.14458, 6.34260)
	}
}

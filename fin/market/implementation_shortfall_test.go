package market

import (
	"math"
	"testing"
)

// Perold 1988 defines shortfall = execution cost + opportunity cost (+ fixed costs, excluded here).
// Hand-checked: buy 100, decision 1000, fills 60 @ 1002 and 20 @ 1005, end 1010, tick 0.01 →
// execution (60·2 + 20·5)·0.01 = 2.2, opportunity (100 − 80)·10·0.01 = 2.0. Tolerance 1e-12.
func TestImplementationShortfall_handChecked(t *testing.T) {
	ex, op, err := ImplementationShortfall(Buy, 1000, []int64{1002, 1005}, []float64{60, 20}, 100, 1010, 0.01)
	if err != nil || !near(ex, 2.2, 1e-12) || !near(op, 2.0, 1e-12) {
		t.Errorf("got %v, %v, %v; want 2.2, 2.0", ex, op, err)
	}
}

func TestImplementationShortfall_cases(t *testing.T) {
	// sell: selling below the decision price is a cost; price fall on the unsold shares is a cost
	ex, op, err := ImplementationShortfall(Sell, 1000, []int64{998}, []float64{50}, 100, 990, 1)
	if err != nil || ex != 100 || op != 500 {
		t.Errorf("sell: %v, %v, %v; want 100, 500", ex, op, err)
	}
	// no fills: all opportunity cost; favourable move is a negative cost
	ex, op, err = ImplementationShortfall(Buy, 1000, nil, nil, 10, 990, 1)
	if err != nil || ex != 0 || op != -100 {
		t.Errorf("no fills: %v, %v, %v; want 0, -100", ex, op, err)
	}
	// fully filled: no opportunity cost
	_, op, err = ImplementationShortfall(Buy, 1000, []int64{1000}, []float64{10}, 10, 2000, 1)
	if err != nil || op != 0 {
		t.Errorf("filled: %v, %v", op, err)
	}
}

func TestImplementationShortfall_errors(t *testing.T) {
	_, _, err := ImplementationShortfall(0, 1, nil, nil, 1, 1, 1)
	wantErr(t, "side", err, ErrOutOfRange)
	_, _, err = ImplementationShortfall(Buy, -1, nil, nil, 1, 1, 1)
	wantErr(t, "decision < 0", err, ErrNegative)
	_, _, err = ImplementationShortfall(Buy, 1, nil, nil, 1, maxTicks+1, 1)
	wantErr(t, "end too large", err, ErrOutOfRange)
	_, _, err = ImplementationShortfall(Buy, 1, nil, nil, math.NaN(), 1, 1)
	wantErr(t, "orderSize NaN", err, ErrNotFinite)
	_, _, err = ImplementationShortfall(Buy, 1, nil, nil, 1, 1, 0)
	wantErr(t, "tick 0", err, ErrOutOfRange)
	_, _, err = ImplementationShortfall(Buy, 1, []int64{1}, nil, 1, 1, 1)
	wantErr(t, "lengths", err, ErrLengthMismatch)
	_, _, err = ImplementationShortfall(Buy, 1, []int64{1}, []float64{-1}, 1, 1, 1)
	wantErr(t, "negative fill", err, ErrNegative)
	_, _, err = ImplementationShortfall(Buy, 1, []int64{1}, []float64{2}, 1, 1, 1)
	wantErr(t, "overfilled", err, ErrDomain)
	_, _, err = ImplementationShortfall(Buy, 0, nil, nil, math.MaxFloat64, maxTicks, 2)
	wantErr(t, "overflow", err, ErrOverflow)
}

func BenchmarkImplementationShortfall(b *testing.B) {
	p, s := make([]int64, 100), make([]float64, 100)
	for i := range p {
		p[i], s[i] = int64(1000+i%7), 1
	}
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkFloat2, sinkErr = ImplementationShortfall(Buy, 1000, p, s, 150, 1010, 0.01)
	}
}

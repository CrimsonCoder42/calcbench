package market

import (
	"math"
	"testing"
)

// Cont, Kukanov, Stoikov 2014 §2.1 define e_n case by case; hand-checked values for each case.
func TestOrderFlowImbalance_cases(t *testing.T) {
	for _, c := range []struct {
		name     string
		bid, ask []int64
		qb, qa   []float64
		want     float64
	}{
		{"one snapshot", []int64{100}, []int64{101}, []float64{5}, []float64{7}, 0},
		// both prices unchanged: q^B_n − q^B_{n−1} − q^A_n + q^A_{n−1} = 8 − 5 − 6 + 7
		{"unchanged prices", []int64{100, 100}, []int64{101, 101}, []float64{5, 8}, []float64{7, 6}, 4},
		// bid up: + q^B_n only; ask up: + q^A_{n−1} only → 3 + 7
		{"both up", []int64{100, 101}, []int64{102, 103}, []float64{5, 3}, []float64{7, 9}, 10},
		// bid down: − q^B_{n−1}; ask down: − q^A_n → −5 − 9
		{"both down", []int64{100, 99}, []int64{102, 101}, []float64{5, 3}, []float64{7, 9}, -14},
		// two steps sum: step 1 = 4 (as above), step 2 bid up (+2), ask same (−6 + 6 = 0) → 6
		{"two steps", []int64{100, 100, 101}, []int64{102, 102, 102}, []float64{5, 8, 2}, []float64{7, 6, 6}, 6},
	} {
		got, err := OrderFlowImbalance(c.bid, c.ask, c.qb, c.qa)
		if err != nil || got != c.want {
			t.Errorf("%s: %v, %v; want %v", c.name, got, err, c.want)
		}
	}
}

func TestOrderFlowImbalance_errors(t *testing.T) {
	_, err := OrderFlowImbalance(nil, nil, nil, nil)
	wantErr(t, "empty", err, ErrEmpty)
	_, err = OrderFlowImbalance([]int64{1, 2}, []int64{3}, []float64{1, 1}, []float64{1, 1})
	wantErr(t, "lengths", err, ErrLengthMismatch)
	_, err = OrderFlowImbalance([]int64{1, 4}, []int64{3, 3}, []float64{1, 1}, []float64{1, 1})
	wantErr(t, "crossed", err, ErrDomain)
	_, err = OrderFlowImbalance([]int64{1}, []int64{3}, []float64{-1}, []float64{1})
	wantErr(t, "negative size", err, ErrNegative)
	_, err = OrderFlowImbalance([]int64{1}, []int64{3}, []float64{1}, []float64{math.NaN()})
	wantErr(t, "NaN size", err, ErrNotFinite)
	_, err = OrderFlowImbalance([]int64{1, 1, 2}, []int64{3, 3, 3}, []float64{0, math.MaxFloat64, math.MaxFloat64}, []float64{0, 0, 0})
	wantErr(t, "overflow", err, ErrOverflow)
}

func BenchmarkOrderFlowImbalance(b *testing.B) {
	const n = 1000
	bid, ask := make([]int64, n), make([]int64, n)
	qb, qa := make([]float64, n), make([]float64, n)
	for i := range n {
		bid[i], ask[i] = int64(100+i%3), int64(103+i%2)
		qb[i], qa[i] = float64(i%7+1), float64(i%5+1)
	}
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, sinkErr = OrderFlowImbalance(bid, ask, qb, qa)
	}
}

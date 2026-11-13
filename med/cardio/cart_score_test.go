package cardio

import (
	"math"
	"testing"
)

// Churpek 2012 Table 3: each band's points, one input varied at a time from a zero-point baseline
// (RR 18, HR 80, DBP 70, age 40).
func TestCARTScore_table3(t *testing.T) {
	tests := []struct {
		rr, hr, dbp, age float64
		want             int
	}{
		{18, 80, 70, 40, 0},
		{20.9, 80, 70, 40, 0}, {21, 80, 70, 40, 8}, {23, 80, 70, 40, 8}, {24, 80, 70, 40, 12}, {25, 80, 70, 40, 12},
		{26, 80, 70, 40, 15}, {29, 80, 70, 40, 15}, {29.1, 80, 70, 40, 22}, {30, 80, 70, 40, 22},
		{18, 109.9, 70, 40, 0}, {18, 110, 70, 40, 4}, {18, 139, 70, 40, 4}, {18, 139.1, 70, 40, 13},
		{18, 80, 50, 40, 0}, {18, 80, 49, 40, 4}, {18, 80, 40, 40, 4}, {18, 80, 39, 40, 6}, {18, 80, 35, 40, 6},
		{18, 80, 34.9, 40, 13},
		{18, 80, 70, 54.9, 0}, {18, 80, 70, 55, 4}, {18, 80, 70, 69, 4}, {18, 80, 70, 69.1, 9},
		{24, 115, 45, 72, 29}, // doc Example (HAND-DERIVED)
		{35, 150, 30, 80, 57}, // every top band: 22 + 13 + 13 + 9
	}
	for _, tc := range tests {
		got, err := CARTScore(tc.rr, tc.hr, tc.dbp, tc.age)
		if err != nil || got != tc.want {
			t.Errorf("CARTScore(%v, %v, %v, %v) = (%d, %v), want %d", tc.rr, tc.hr, tc.dbp, tc.age, got, err, tc.want)
		}
	}
}

// parity: ledger cart_score — the source has no test file and no JSON example; its compute() bands agree with
// Table 3 on the integer inputs below (source cases are these app-band checks, not printed examples).
func TestCARTScore_parity(t *testing.T) {
	tests := []struct {
		rr, hr, dbp, age float64
		want             int
	}{
		{22, 120, 38, 60, 22}, // parity: ledger cart_score, source case 1 (app bands)
		{16, 72, 80, 30, 0},   // parity: ledger cart_score, source case 2 (app bands)
	}
	for i, tc := range tests {
		if got, err := CARTScore(tc.rr, tc.hr, tc.dbp, tc.age); err != nil || got != tc.want {
			t.Errorf("case %d: (%d, %v), want %d", i+1, got, err, tc.want)
		}
	}
}

func TestCARTScore_errors(t *testing.T) {
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		_, err := CARTScore(v, 80, 70, 40)
		checkErr(t, "rr non-finite", err, ErrNotFinite)
		_, err = CARTScore(18, 80, 70, v)
		checkErr(t, "age non-finite", err, ErrNotFinite)
	}
	// a NaN is reported before an out-of-range value elsewhere
	_, err := CARTScore(-1, math.NaN(), 70, 40)
	checkErr(t, "NaN first", err, ErrNotFinite)
	for _, in := range [][4]float64{{-0.1, 80, 70, 40}, {80.1, 80, 70, 40}, {18, 300.1, 70, 40}, {18, 80, -1, 40}, {18, 80, 200.1, 40}, {18, 80, 70, 17.9}, {18, 80, 70, 110.1}} {
		_, err := CARTScore(in[0], in[1], in[2], in[3])
		checkErr(t, "range", err, ErrOutOfRange)
	}
	for _, in := range [][4]float64{{0, 0, 0, 18}, {80, 300, 200, 110}} {
		if _, err := CARTScore(in[0], in[1], in[2], in[3]); err != nil {
			t.Errorf("edge %v: %v", in, err)
		}
	}
}

func TestCARTScore_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = CARTScore(24, 115, 45, 72) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
}

func BenchmarkCARTScore(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkInt, _ = CARTScore(24, 115, 45, 72)
	}
}

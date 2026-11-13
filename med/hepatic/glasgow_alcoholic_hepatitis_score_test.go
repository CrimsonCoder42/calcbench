package hepatic

import (
	"math"
	"testing"
)

// Forrest 2005 Table 4 (PMC1774903): both sides of every cut-off. Integer points, exact comparison.
func TestGlasgowAlcoholicHepatitisScore_example(t *testing.T) {
	type in struct{ age, wcc, urea, ptr, bili float64 }
	tests := []struct {
		name string
		x    in
		want int
	}{
		{"worked", in{55, 16, 7, 1.8, 300}, 11},
		{"minimum", in{30, 8, 3, 1.1, 50}, 5},
		{"maximum", in{70, 30, 20, 3, 400}, 12},
		{"age 49.9", in{49.9, 8, 3, 1.1, 50}, 5},
		{"age 50", in{50, 8, 3, 1.1, 50}, 6},
		{"WCC 14.9", in{30, 14.9, 3, 1.1, 50}, 5},
		{"WCC 15", in{30, 15, 3, 1.1, 50}, 6},
		{"urea 4.9", in{30, 8, 4.9, 1.1, 50}, 5},
		{"urea 5", in{30, 8, 5, 1.1, 50}, 6},
		{"PT ratio 1.49", in{30, 8, 3, 1.49, 50}, 5},
		{"PT ratio 1.5", in{30, 8, 3, 1.5, 50}, 6},
		{"PT ratio 2.0", in{30, 8, 3, 2.0, 50}, 6},
		{"PT ratio 2.01", in{30, 8, 3, 2.01, 50}, 7},
		{"bilirubin 124.9", in{30, 8, 3, 1.1, 124.9}, 5},
		{"bilirubin 125", in{30, 8, 3, 1.1, 125}, 6},
		{"bilirubin 250", in{30, 8, 3, 1.1, 250}, 6},
		{"bilirubin 250.1", in{30, 8, 3, 1.1, 250.1}, 7},
	}
	for _, tc := range tests {
		got, err := GlasgowAlcoholicHepatitisScore(tc.x.age, tc.x.wcc, tc.x.urea, tc.x.ptr, tc.x.bili)
		if err != nil || got != tc.want {
			t.Errorf("%s: got (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

// parity: ledger glasgow_alcoholic_hepatitis. The app takes BUN and bilirubin in mg/dL and converts (BUN/2.8,
// bilirubin × 17.1) and takes the measured PT over the laboratory normal PT; the conversions are the caller's, done here as the app does.
func TestGlasgowAlcoholicHepatitisScore_parity(t *testing.T) {
	cases := []struct {
		age, wbc, bunMgDL, biliMgDL, pt, ptNormal float64
		want                                      int
	}{
		{45, 10, 20, 5, 14, 12, 6},   // parity: ledger glasgow_alcoholic_hepatitis, source case 1: 1+1+2+1+1
		{60, 18, 30, 16, 30, 12, 12}, // source case 2: 2+2+2+3 (bili 273.6)+3 (ratio 2.5)
		{40, 5, 10, 8, 20, 12, 7},    // source case 3: 1+1+1 (urea 3.57)+2 (bili 136.8)+2 (ratio 1.667)
	}
	for i, c := range cases {
		got, err := GlasgowAlcoholicHepatitisScore(c.age, c.wbc, c.bunMgDL/2.8, c.pt/c.ptNormal, c.biliMgDL*17.1)
		if err != nil || got != c.want {
			t.Errorf("case %d: got (%v, %v), want %v", i+1, got, err, c.want)
		}
	}
}

func TestGlasgowAlcoholicHepatitisScore_errors(t *testing.T) {
	call := func(a, w, u, p, b float64) func() error {
		return func() error { _, err := GlasgowAlcoholicHepatitisScore(a, w, u, p, b); return err }
	}
	runErrCases(t, []errCase{
		{"age NaN", call(nan, 10, 5, 1.2, 100), ErrNotFinite},
		{"WCC +Inf", call(50, pInf, 5, 1.2, 100), ErrNotFinite},
		{"urea -Inf", call(50, 10, nInf, 1.2, 100), ErrNotFinite},
		{"PT ratio NaN", call(50, 10, 5, nan, 100), ErrNotFinite},
		{"bilirubin +Inf", call(50, 10, 5, 1.2, math.Inf(1)), ErrNotFinite},
		{"NaN before range", call(0, 10, 5, 1.2, nan), ErrNotFinite},
		{"age 0", call(0, 10, 5, 1.2, 100), ErrOutOfRange},
		{"WCC negative", call(50, -1, 5, 1.2, 100), ErrOutOfRange},
		{"urea 0", call(50, 10, 0, 1.2, 100), ErrOutOfRange},
		{"PT ratio 0", call(50, 10, 5, 0, 100), ErrOutOfRange},
		{"bilirubin 0", call(50, 10, 5, 1.2, 0), ErrOutOfRange},
	})
}

func TestGlasgowAlcoholicHepatitisScore_noAllocation(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkInt, _ = GlasgowAlcoholicHepatitisScore(55, 16, 7, 1.8, 300) }); n != 0 {
		t.Errorf("GlasgowAlcoholicHepatitisScore allocates %v times, want 0", n)
	}
}

func BenchmarkGlasgowAlcoholicHepatitisScore(b *testing.B) {
	b.ReportAllocs()
	age, wcc, urea, ptr, bili := 55.0, 16.0, 7.0, 1.8, 300.0
	for b.Loop() {
		sinkInt, _ = GlasgowAlcoholicHepatitisScore(age, wcc, urea, ptr, bili)
	}
}

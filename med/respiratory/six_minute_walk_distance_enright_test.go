package respiratory

import (
	"errors"
	"math"
	"testing"
)

// Hand-derived from the abstract's men's equation (Enright & Sherrill 1998, PMID 9817683):
// 7.57×170.18 − 5.02×60 − 1.76×77.1107 − 309 = 542.3478 m. Tolerance 1e-9 (exact arithmetic on the inputs).
func TestSixMinuteWalkDistanceEnright_handExample(t *testing.T) {
	kg := 170 * kgPerPound
	want := 7.57*170.18 - 5.02*60 - 1.76*kg - 309
	got, err := SixMinuteWalkDistanceEnright(170.18, kg, 60, false)
	if err != nil || math.Abs(got-want) > 1e-9 {
		t.Errorf("got (%v, %v), want %v", got, err, want)
	}
	// Women: 2.11×162.56 − 2.29×68.0389 − 5.78×55 + 667.
	kgW := 150 * kgPerPound
	wantW := 2.11*162.56 - 2.29*kgW - 5.78*55 + 667
	gotW, err := SixMinuteWalkDistanceEnright(162.56, kgW, 55, true)
	if err != nil || math.Abs(gotW-wantW) > 1e-9 {
		t.Errorf("women: got (%v, %v), want %v", gotW, err, wantW)
	}
}

// parity: ledger six_minute_walk, source cases 1–5. The source shows Math.round of the prediction (and an LLN
// of prediction − 153 m / − 139 m, which was not seen in the primary and is not ported); expected values are the
// source's rounded predictions, so the tolerance is 0.5 m. Inputs in inches and pounds, as the source takes them.
func TestSixMinuteWalkDistanceEnright_parity(t *testing.T) {
	cases := []struct {
		heightIn, ageYears, weightLb float64
		female                       bool
		want                         float64
	}{
		{67, 60, 170, false, 542}, // source case 1
		{64, 55, 150, true, 536},  // source case 2
		{72, 75, 220, false, 523}, // source case 3
		{60, 80, 130, true, 391},  // source case 4
		{70, 40, 180, false, 692}, // source case 5
	}
	for i, c := range cases {
		got, err := SixMinuteWalkDistanceEnrightImperial(c.heightIn, c.weightLb, c.ageYears, c.female)
		if err != nil || math.Abs(got-c.want) > 0.5 {
			t.Errorf("case %d: got (%v, %v), want %v ± 0.5", i+1, got, err, c.want)
		}
	}
}

// Imperial and metric twins agree to 1e-9 m (the conversion is exact).
func TestSixMinuteWalkDistanceEnrightImperial_agreesWithMetric(t *testing.T) {
	a, _ := SixMinuteWalkDistanceEnrightImperial(70, 180, 50, true)
	b, _ := SixMinuteWalkDistanceEnright(70*2.54, 180*0.45359237, 50, true)
	if math.Abs(a-b) > 1e-9 {
		t.Errorf("imperial %v, metric %v", a, b)
	}
}

func TestSixMinuteWalkDistanceEnright_ageEdges(t *testing.T) {
	for _, age := range []float64{40, 80} {
		if _, err := SixMinuteWalkDistanceEnright(170, 70, age, false); err != nil {
			t.Errorf("age %v: err = %v, want nil", age, err)
		}
	}
	for _, age := range []float64{39.999, 80.001} {
		if _, err := SixMinuteWalkDistanceEnright(170, 70, age, false); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("age %v: err = %v, want ErrOutOfRange", age, err)
		}
	}
}

func TestSixMinuteWalkDistanceEnright_errors(t *testing.T) {
	tests := []struct {
		name      string
		h, w, age float64
		imperial  bool
		wantErr   error
	}{
		{"zero height", 0, 70, 50, false, ErrOutOfRange},
		{"negative weight", 170, -1, 50, false, ErrOutOfRange},
		{"NaN height", math.NaN(), 70, 50, false, ErrNotFinite},
		{"+Inf weight", 170, math.Inf(1), 50, false, ErrNotFinite},
		{"-Inf age", 170, 70, math.Inf(-1), false, ErrNotFinite},
		{"NaN age", 170, 70, math.NaN(), false, ErrNotFinite},
		{"overflow", math.MaxFloat64, 70, 50, false, ErrDomain},
		{"imperial zero height", 0, 150, 50, true, ErrOutOfRange},
		{"imperial NaN weight", 67, math.NaN(), 50, true, ErrNotFinite},
		{"imperial age 81", 67, 150, 81, true, ErrOutOfRange},
	}
	for _, tc := range tests {
		var err error
		if tc.imperial {
			_, err = SixMinuteWalkDistanceEnrightImperial(tc.h, tc.w, tc.age, false)
		} else {
			_, err = SixMinuteWalkDistanceEnright(tc.h, tc.w, tc.age, false)
		}
		if !errors.Is(err, tc.wantErr) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.wantErr)
		}
	}
}

func TestSixMinuteWalkDistanceEnright_noAlloc(t *testing.T) {
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = SixMinuteWalkDistanceEnright(170, 70, 60, false) }); n != 0 {
		t.Errorf("allocs = %v, want 0", n)
	}
	if n := testing.AllocsPerRun(100, func() { sinkFloat, _ = SixMinuteWalkDistanceEnrightImperial(67, 170, 60, true) }); n != 0 {
		t.Errorf("imperial allocs = %v, want 0", n)
	}
}

func BenchmarkSixMinuteWalkDistanceEnright(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = SixMinuteWalkDistanceEnright(170, 70, 60, false)
	}
}

func BenchmarkSixMinuteWalkDistanceEnrightImperial(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		sinkFloat, _ = SixMinuteWalkDistanceEnrightImperial(67, 170, 60, true)
	}
}

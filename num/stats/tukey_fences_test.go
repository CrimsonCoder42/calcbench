package stats

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// e-Handbook §7.1.6 worked example (N = 90): inner fences −39.0 and 1211.0, outer −507.75 and 1679.75. Those
// four numbers fix the quartiles: 4·IQ = 1211 − (−39), so IQ = 312.5, Q1 = 429.75, Q3 = 742.25. Exact in
// binary (multiples of 1/4), so compared exactly.
func TestTukeyFences_sourceExample(t *testing.T) {
	f, err := fencesFromQuartiles(429.75, 742.25)
	if err != nil {
		t.Fatal(err)
	}
	want := Fences{Q1: 429.75, Q3: 742.25, IQ: 312.5, InnerLow: -39, InnerHigh: 1211, OuterLow: -507.75, OuterHigh: 1679.75}
	if f != want {
		t.Errorf("fences = %+v, want %+v", f, want)
	}
}

// parity: ledger TukeyFences (robust.go:29), source case S3: resting pulse 62, 64, 63, 61, 65, 118 with
// type-7 quartiles → Q1 62.25, Q3 64.75, inner high 68.5, outer high 72.25. Tolerance 1e-12 absolute.
func TestTukeyFences_parity(t *testing.T) {
	f, err := TukeyFences([]float64{62, 64, 63, 61, 65, 118}, HyndmanFan7)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		name      string
		got, want float64
	}{
		{"Q1", f.Q1, 62.25}, {"Q3", f.Q3, 64.75}, {"IQ", f.IQ, 2.5},
		{"InnerLow", f.InnerLow, 58.5}, {"InnerHigh", f.InnerHigh, 68.5},
		{"OuterLow", f.OuterLow, 54.75}, {"OuterHigh", f.OuterHigh, 72.25},
	} {
		if math.Abs(c.got-c.want) > 1e-12 {
			t.Errorf("%s = %v, want %v", c.name, c.got, c.want)
		}
	}
}

func TestTukeyFences_cases(t *testing.T) {
	tests := []struct {
		name   string
		values []float64
		method QuantileMethod
		q1, q3 float64
	}{
		{"one value", []float64{5}, HyndmanFan7, 5, 5},
		{"all equal", []float64{3, 3, 3, 3}, HyndmanFan7, 3, 3},
		{"1..8 type 7", []float64{8, 1, 7, 2, 6, 3, 5, 4}, HyndmanFan7, 2.75, 6.25},
		{"1..8 type 6", []float64{8, 1, 7, 2, 6, 3, 5, 4}, HyndmanFan6, 2.25, 6.75},
		{"negative", []float64{-4, -3, -2, -1}, HyndmanFan7, -3.25, -1.75},
	}
	for _, tc := range tests {
		f, err := TukeyFences(tc.values, tc.method)
		if err != nil || math.Abs(f.Q1-tc.q1) > 1e-12 || math.Abs(f.Q3-tc.q3) > 1e-12 {
			t.Errorf("%s: Q1, Q3 = %v, %v (%v), want %v, %v", tc.name, f.Q1, f.Q3, err, tc.q1, tc.q3)
			continue
		}
		iq := tc.q3 - tc.q1
		if math.Abs(f.InnerLow-(tc.q1-1.5*iq)) > 1e-12 || math.Abs(f.OuterHigh-(tc.q3+3*iq)) > 1e-12 {
			t.Errorf("%s: fences %+v", tc.name, f)
		}
	}
}

// Oracle: the quartiles equal Percentile's on random samples (fixed seed, 200 samples of up to 300
// values), and the input is unchanged.
func TestTukeyFences_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for range 200 {
		v := randomFloats(rng, 1+rng.IntN(300))
		before := slices.Clone(v)
		f, err := TukeyFences(v, HyndmanFan7)
		if err != nil {
			t.Fatal(err)
		}
		q1, _ := Percentile(v, 0.25, HyndmanFan7)
		q3, _ := Percentile(v, 0.75, HyndmanFan7)
		if f.Q1 != q1 || f.Q3 != q3 {
			t.Fatalf("quartiles %v, %v, want %v, %v", f.Q1, f.Q3, q1, q3)
		}
		if !slices.Equal(v, before) {
			t.Fatal("input modified")
		}
	}
}

func TestTukeyFences_errors(t *testing.T) {
	if _, err := TukeyFences(nil, HyndmanFan7); !errors.Is(err, ErrEmpty) {
		t.Errorf("empty: %v, want ErrEmpty", err)
	}
	if _, err := TukeyFences([]float64{1, math.NaN()}, HyndmanFan7); !errors.Is(err, ErrNotFinite) {
		t.Errorf("NaN: %v, want ErrNotFinite", err)
	}
	if _, err := TukeyFences([]float64{1, math.Inf(1)}, HyndmanFan7); !errors.Is(err, ErrNotFinite) {
		t.Errorf("Inf: %v, want ErrNotFinite", err)
	}
	if _, err := TukeyFences([]float64{1, 2}, QuantileMethod(0)); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("method 0: %v, want ErrOutOfRange", err)
	}
	big := []float64{-math.MaxFloat64, -math.MaxFloat64, math.MaxFloat64, math.MaxFloat64}
	if _, err := TukeyFences(big, HyndmanFan7); !errors.Is(err, ErrOverflow) {
		t.Errorf("huge spread: %v, want ErrOverflow", err)
	}
}

func BenchmarkTukeyFences(b *testing.B) {
	for _, size := range benchSizes {
		v := randomFloats(rand.New(rand.NewPCG(1, 2)), size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = TukeyFences(v, HyndmanFan7)
			}
		})
	}
}

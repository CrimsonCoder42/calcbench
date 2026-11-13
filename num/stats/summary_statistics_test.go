package stats

import (
	"errors"
	"math"
	"math/rand/v2"
	"testing"
)

// Worked example: StRD NumAcc1 → certified mean 10000002 and SD 1; Min 10000001, Max 10000003,
// Median 10000002, Range 2 (hand-checked).
func TestSummaryStatistics_numAcc1(t *testing.T) {
	want := Summary{Count: 3, Mean: 10000002, StandardDeviation: 1, Min: 10000001, Max: 10000003, Median: 10000002, Range: 2}
	if got, err := SummaryStatistics(numAcc1); err != nil || got != want {
		t.Errorf("SummaryStatistics(NumAcc1) = %+v, %v; want %+v", got, err, want)
	}
}

// Every field equals its own function's result exactly (the shared-core guarantee).
func TestSummaryStatistics_agreesWithParts(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	inputs := [][]float64{michelso, mavro, lew, resistivity, numAcc4}
	for range 50 {
		inputs = append(inputs, randomFloats(rng, 2+rng.IntN(300)))
	}
	for i, v := range inputs {
		got, err := SummaryStatistics(v)
		mean, _ := ArithmeticMean(v)
		sd, _ := SampleStandardDeviation(v)
		med, _ := Median(v)
		spread, _ := Range(v)
		lo, _ := Percentile(v, 0, HyndmanFan7)
		hi, _ := Percentile(v, 1, HyndmanFan7)
		want := Summary{Count: len(v), Mean: mean, StandardDeviation: sd, Min: lo, Max: hi, Median: med, Range: spread}
		if err != nil || got != want {
			t.Fatalf("input %d: SummaryStatistics = %+v, %v; want %+v", i, got, err, want)
		}
	}
}

func TestSummaryStatistics_errors(t *testing.T) {
	m := math.MaxFloat64
	for _, tt := range []struct {
		name   string
		values []float64
		want   error
	}{
		{"empty", nil, ErrEmpty},
		{"one element", []float64{1}, ErrTooFew},
		{"NaN", []float64{1, math.NaN()}, ErrNotFinite},
		{"SD overflows", []float64{-m, m}, ErrOverflow},
		{"range overflows", []float64{-m, m, 0, 0, 0, 0, 0, 0}, ErrOverflow},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := SummaryStatistics(tt.values); !errors.Is(err, tt.want) || got != (Summary{}) {
				t.Errorf("SummaryStatistics(%v) = %+v, %v; want zero, %v", tt.values, got, err, tt.want)
			}
		})
	}
}

func TestSummaryStatistics_allocations(t *testing.T) {
	if a := testing.AllocsPerRun(100, func() { s, _ := SummaryStatistics(michelso); benchSink = s.Mean }); a != 1 {
		t.Errorf("SummaryStatistics allocates %v times per call, want 1 (the median's copy)", a)
	}
}

func BenchmarkSummaryStatistics(b *testing.B) {
	for _, size := range benchSizes {
		inputs := benchInputs(size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; b.Loop(); i++ {
				s, _ := SummaryStatistics(inputs[i&7])
				benchSink = s.Median
			}
		})
	}
}

// Tiny spread: the StandardDeviation field is the sample SD, 2^-700 for 2^-700·(1, 2, 3) (hand: s = 1
// for 1, 2, 3, scaled exactly). Before the underflow fix it was 0.
func TestSummaryStatistics_tinySpread(t *testing.T) {
	got, err := SummaryStatistics([]float64{0x1p-700, 0x2p-700, 0x3p-700})
	want := Summary{Count: 3, Mean: 0x2p-700, StandardDeviation: 0x1p-700, Min: 0x1p-700, Max: 0x3p-700, Median: 0x2p-700, Range: 0x2p-700}
	if err != nil || got != want {
		t.Errorf("SummaryStatistics(2^-700·(1 2 3)) = %+v, %v; want %+v", got, err, want)
	}
}

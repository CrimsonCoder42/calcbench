package fitness

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"testing"
	"time"
)

func stream(t0 time.Time, pts ...[2]float64) []DistanceSample {
	out := make([]DistanceSample, len(pts))
	for i, p := range pts {
		out[i] = DistanceSample{At: t0.Add(time.Duration(p[0] * float64(time.Second))), DistanceKmOrMi: p[1]}
	}
	return out
}

// parity: ledger Splits, source cases 1–4 (mile splits 566.667 s and 433.333 s, elapsed 1000 s; unit
// 0.621371 mi gives 3 splits, the first 364.7312 s ±1e-3 as printed; one sample and a backwards distance
// refused).
func TestSplits_parity(t *testing.T) {
	t0 := time.Date(2026, 10, 30, 12, 0, 0, 0, time.UTC)
	s := stream(t0, [2]float64{0, 0}, [2]float64{300, 0.5}, [2]float64{700, 1.25}, [2]float64{1000, 2.0})
	got, err := Splits(s, 1)
	if err != nil {
		t.Fatal(err)
	}
	// Tolerance 1e-6 s: the source's printed values carry 7 decimals.
	if len(got) != 2 || math.Abs(got[0].Seconds-566.6666667) > 1e-6 || math.Abs(got[1].Seconds-433.3333333) > 1e-6 ||
		math.Abs(got[1].ElapsedSeconds-1000) > 1e-9 || got[0].Index != 1 || got[1].Index != 2 {
		t.Fatalf("got %+v", got)
	}
	km, err := Splits(s, 0.621371)
	if err != nil || len(km) != 3 || math.Abs(km[0].Seconds-364.7312) > 1e-3 {
		t.Fatalf("km: %+v, %v", km, err)
	}
	if _, err := Splits(s[:1], 1); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("one sample: %v", err)
	}
	if _, err := Splits(stream(t0, [2]float64{0, 1}, [2]float64{10, 0.5}), 1); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("backwards distance: %v", err)
	}
}

func TestSplits_cases(t *testing.T) {
	t0 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	// No whole unit: empty, non-nil.
	got, err := Splits(stream(t0, [2]float64{0, 0}, [2]float64{60, 0.4}), 1)
	if err != nil || got == nil || len(got) != 0 {
		t.Errorf("partial unit: %v, %v", got, err)
	}
	// A stream that starts part-way: the first boundary after the first sample is 2 km, timed from it.
	got, _ = Splits(stream(t0, [2]float64{0, 1.5}, [2]float64{100, 2.5}), 1)
	if len(got) != 1 || got[0].Index != 2 || got[0].Seconds != 50 {
		t.Errorf("late start: %+v", got)
	}
	// Two boundaries inside one segment, and a boundary exactly on a sample.
	got, _ = Splits(stream(t0, [2]float64{0, 0}, [2]float64{400, 2}, [2]float64{600, 3}), 1)
	if len(got) != 3 || got[0].Seconds != 200 || got[1].Seconds != 200 || got[2].Seconds != 200 {
		t.Errorf("multi: %+v", got)
	}
	// Distance jumps with no time (equal instants).
	got, _ = Splits(stream(t0, [2]float64{0, 0}, [2]float64{0, 1}), 1)
	if len(got) != 1 || got[0].Seconds != 0 {
		t.Errorf("zero time: %+v", got)
	}
}

func TestSplits_errors(t *testing.T) {
	t0 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	ok := stream(t0, [2]float64{0, 0}, [2]float64{60, 1})
	cases := []struct {
		name string
		s    []DistanceSample
		unit float64
		want error
	}{
		{"zero unit", ok, 0, ErrOutOfRange},
		{"negative unit", ok, -1, ErrOutOfRange},
		{"NaN unit", ok, math.NaN(), ErrNotFinite},
		{"empty", nil, 1, ErrOutOfRange},
		{"negative distance", stream(t0, [2]float64{0, -1}, [2]float64{1, 0}), 1, ErrOutOfRange},
		{"NaN distance", stream(t0, [2]float64{0, 0}, [2]float64{1, math.NaN()}), 1, ErrNotFinite},
		{"Inf distance", stream(t0, [2]float64{0, 0}, [2]float64{1, math.Inf(1)}), 1, ErrNotFinite},
		{"backwards time", stream(t0, [2]float64{10, 0}, [2]float64{0, 1}), 1, ErrOutOfRange},
		{"too many splits", ok, 1e-12, ErrDomain},
	}
	for _, c := range cases {
		if _, err := Splits(c.s, c.unit); !errors.Is(err, c.want) {
			t.Errorf("%s: err = %v, want %v", c.name, err, c.want)
		}
	}
}

// parity: ledger Splits, source case 5 (property: the splits' own times sum to the last elapsed time and
// their count is floor(total / unit); 200 random streams, fixed seed).
func TestSplits_property(t *testing.T) {
	r := rand.New(rand.NewPCG(7, 11))
	t0 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for trial := 0; trial < 200; trial++ {
		s := []DistanceSample{{At: t0}}
		for i := 1; i < 40; i++ {
			s = append(s, DistanceSample{At: s[i-1].At.Add(time.Duration(1+r.IntN(120)) * time.Second),
				DistanceKmOrMi: s[i-1].DistanceKmOrMi + r.Float64()*0.3})
		}
		got, err := Splits(s, 1)
		if err != nil {
			t.Fatal(err)
		}
		if want := int(math.Floor(s[len(s)-1].DistanceKmOrMi)); len(got) != want {
			t.Fatalf("trial %d: %d splits, want %d", trial, len(got), want)
		}
		sum := 0.0
		for _, sp := range got {
			sum += sp.Seconds
		}
		// Tolerance 1e-6 s: a sum of up to 12 rounded differences.
		if len(got) > 0 && math.Abs(sum-got[len(got)-1].ElapsedSeconds) > 1e-6 {
			t.Fatalf("trial %d: sum %v, elapsed %v", trial, sum, got[len(got)-1].ElapsedSeconds)
		}
	}
}

func BenchmarkSplits(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000} {
		t0 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
		s := make([]DistanceSample, n)
		for i := range s {
			s[i] = DistanceSample{At: t0.Add(time.Duration(i) * time.Second), DistanceKmOrMi: float64(i) * 0.003}
		}
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := Splits(s, 1); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

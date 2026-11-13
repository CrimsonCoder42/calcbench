package stats

import (
	"errors"
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Worked example: e-Handbook §7.2.6.2, resistivity, p = 0.90, type 6 printed 95.1981.
func TestPercentiles_resistivity(t *testing.T) {
	got, err := Percentiles(resistivity, []float64{0.9}, HyndmanFan6)
	// Tolerance 5e-5 absolute: the handbook prints four decimals.
	if err != nil || math.Abs(got[0]-95.1981) > 5e-5 {
		t.Errorf("Percentiles = %v, %v; want [95.1981]", got, err)
	}
}

// Oracle: every element equals Percentile for every method, on random data with ties.
func TestPercentiles_matchesPercentile(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	ps := []float64{0, 0.01, 0.1, 0.25, 0.5, 0.75, 0.9, 0.99, 1, 1.0 / 3}
	for range 100 {
		n := 1 + r.IntN(40)
		v := make([]float64, n)
		for i := range v {
			v[i] = float64(r.IntN(10))
		}
		for m := HyndmanFan1; m <= HyndmanFan9; m++ {
			got, err := Percentiles(v, ps, m)
			if err != nil {
				t.Fatal(err)
			}
			for i, p := range ps {
				want, _ := Percentile(v, p, m)
				if got[i] != want {
					t.Fatalf("type %d p=%v on %v: got %v want %v", m, p, v, got[i], want)
				}
			}
		}
	}
}

func TestPercentiles_noModify(t *testing.T) {
	v := slices.Clone(resistivity)
	if _, err := Percentiles(v, []float64{0.5}, DefaultQuantileMethod); err != nil || !slices.Equal(v, resistivity) {
		t.Errorf("input modified or err %v", err)
	}
	if got, err := Percentiles(v, nil, DefaultQuantileMethod); err != nil || len(got) != 0 {
		t.Errorf("no ps = %v, %v; want empty", got, err)
	}
}

func TestPercentiles_errors(t *testing.T) {
	for _, tt := range []struct {
		name   string
		values []float64
		ps     []float64
		method QuantileMethod
		want   error
	}{
		{"empty", nil, []float64{0.5}, HyndmanFan7, ErrEmpty},
		{"NaN value", []float64{1, math.NaN()}, []float64{0.5}, HyndmanFan7, ErrNotFinite},
		{"p above 1", []float64{1, 2}, []float64{0.5, 1.5}, HyndmanFan7, ErrOutOfRange},
		{"p NaN", []float64{1, 2}, []float64{math.NaN()}, HyndmanFan7, ErrOutOfRange},
		{"method 0", []float64{1, 2}, []float64{0.5}, 0, ErrOutOfRange},
		{"method 10", []float64{1, 2}, []float64{0.5}, 10, ErrOutOfRange},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Percentiles(tt.values, tt.ps, tt.method); !errors.Is(err, tt.want) {
				t.Errorf("err = %v; want %v", err, tt.want)
			}
		})
	}
}

func BenchmarkPercentiles(b *testing.B) {
	ps := []float64{0.05, 0.25, 0.5, 0.75, 0.95}
	for _, n := range []int{1_000, 10_000, 100_000} {
		in := benchInputs(n)
		b.Run(fmt.Sprint(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = Percentiles(in[0], ps, DefaultQuantileMethod)
			}
		})
	}
}

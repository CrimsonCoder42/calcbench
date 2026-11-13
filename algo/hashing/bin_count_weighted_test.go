package hashing

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"strconv"
	"testing"
)

// Our example (CLRS §11.1 direct addressing, one sum per key): values 0 1 1 2 with weights 0.5 1 2 4 give
// bins 0.5 3 4.
func TestBinCountWeighted_example(t *testing.T) {
	got, err := BinCountWeighted([]int{0, 1, 1, 2}, []float64{0.5, 1, 2, 4}, 3)
	if err != nil || !slices.Equal(got, []float64{0.5, 3, 4}) {
		t.Errorf("BinCountWeighted = (%v, %v)", got, err)
	}
}

// Neumaier 1974: one bin receiving 1, 1e100, 1, -1e100 totals exactly 2.
func TestBinCountWeighted_compensated(t *testing.T) {
	got, err := BinCountWeighted([]int{0, 0, 0, 0}, []float64{1, 1e100, 1, -1e100}, 1)
	if err != nil || got[0] != 2 {
		t.Errorf("BinCountWeighted = (%v, %v), want [2]", got, err)
	}
}

func TestBinCountWeighted_cases(t *testing.T) {
	got, err := BinCountWeighted(nil, nil, 2)
	if err != nil || !slices.Equal(got, []float64{0, 0}) {
		t.Errorf("empty: (%v, %v)", got, err)
	}
	got, err = BinCountWeighted([]int{1}, []float64{-3}, 2)
	if err != nil || !slices.Equal(got, []float64{0, -3}) {
		t.Errorf("negative weight: (%v, %v)", got, err)
	}
}

func TestBinCountWeighted_errors(t *testing.T) {
	tests := []struct {
		name    string
		values  []int
		weights []float64
		bins    int
		want    error
	}{
		{"negative bins", nil, nil, -1, ErrOutOfRange},
		{"value out of range", []int{2}, []float64{1}, 2, ErrOutOfRange},
		{"length", []int{0, 1}, []float64{1}, 2, ErrLength},
		{"NaN", []int{0}, []float64{math.NaN()}, 1, ErrNotFinite},
		{"Inf", []int{0}, []float64{math.Inf(1)}, 1, ErrNotFinite},
		{"overflow", []int{0, 0}, []float64{math.MaxFloat64, math.MaxFloat64}, 1, ErrOverflow},
	}
	for _, tc := range tests {
		if _, err := BinCountWeighted(tc.values, tc.weights, tc.bins); !errors.Is(err, tc.want) {
			t.Errorf("%s: err = %v, want %v", tc.name, err, tc.want)
		}
	}
}

// Oracle: with unit weights it equals BinCount, on 200 random inputs.
func TestBinCountWeighted_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for range 200 {
		values := make([]int, r.IntN(100))
		weights := make([]float64, len(values))
		for i := range values {
			values[i], weights[i] = r.IntN(10), 1
		}
		got, err := BinCountWeighted(values, weights, 10)
		want, _ := BinCount(values, 10)
		if err != nil {
			t.Fatal(err)
		}
		for v := range got {
			if got[v] != float64(want[v]) {
				t.Fatalf("bin %d = %v, want %d", v, got[v], want[v])
			}
		}
	}
}

func BenchmarkBinCountWeighted(b *testing.B) {
	for _, n := range []int{1e3, 1e4, 1e5} {
		r := rand.New(rand.NewPCG(3, 4))
		values := make([]int, n)
		weights := make([]float64, n)
		for i := range values {
			values[i], weights[i] = r.IntN(256), r.Float64()
		}
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = BinCountWeighted(values, weights, 256)
			}
		})
	}
}

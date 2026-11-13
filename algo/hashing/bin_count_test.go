package hashing

import (
	"errors"
	"math/rand/v2"
	"slices"
	"strconv"
	"testing"
)

// Our example (CLRS §8.2, the count array of counting sort): 0 1 1 3 2 1 7 in 8 bins counts 1 3 1 1 0 0 0 1.
func TestBinCount_example(t *testing.T) {
	got, err := BinCount([]int{0, 1, 1, 3, 2, 1, 7}, 8)
	if err != nil || !slices.Equal(got, []int{1, 3, 1, 1, 0, 0, 0, 1}) {
		t.Errorf("BinCount = (%v, %v)", got, err)
	}
}

func TestBinCount_cases(t *testing.T) {
	got, err := BinCount(nil, 0)
	if err != nil || got == nil || len(got) != 0 {
		t.Errorf("empty: (%v, %v)", got, err)
	}
	got, err = BinCount(nil, 3)
	if err != nil || !slices.Equal(got, []int{0, 0, 0}) {
		t.Errorf("no values: (%v, %v)", got, err)
	}
	got, err = BinCount([]int{2, 2, 2}, 3)
	if err != nil || !slices.Equal(got, []int{0, 0, 3}) {
		t.Errorf("all equal: (%v, %v)", got, err)
	}
}

func TestBinCount_errors(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		bins   int
	}{
		{"negative bins", nil, -1},
		{"negative value", []int{0, -1}, 3},
		{"value = bins", []int{3}, 3},
		{"no bins", []int{0}, 0},
	}
	for _, tc := range tests {
		if _, err := BinCount(tc.values, tc.bins); !errors.Is(err, ErrOutOfRange) {
			t.Errorf("%s: err = %v, want ErrOutOfRange", tc.name, err)
		}
	}
}

// Oracle: agrees with FrequencyCount on 300 random inputs.
func TestBinCount_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		values := make([]int, r.IntN(100))
		for i := range values {
			values[i] = r.IntN(20)
		}
		got, err := BinCount(values, 20)
		if err != nil {
			t.Fatal(err)
		}
		freq := FrequencyCount(values)
		for v, c := range got {
			if freq[v] != c {
				t.Fatalf("bin %d = %d, FrequencyCount %d", v, c, freq[v])
			}
		}
	}
}

func BenchmarkBinCount(b *testing.B) {
	for _, n := range []int{1e3, 1e4, 1e5} {
		r := rand.New(rand.NewPCG(3, 4))
		values := make([]int, n)
		for i := range values {
			values[i] = r.IntN(256)
		}
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = BinCount(values, 256)
			}
		})
	}
}

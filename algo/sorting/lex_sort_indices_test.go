package sorting

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"strconv"
	"testing"
)

// Our example (Knuth vol. 3 §5.2, multiple keys): surnames then first names. Records
// (smith, ann) (jones, bob) (smith, al) (jones, bob) sort as jones-bob (1), jones-bob (3, stable), smith-al
// (2), smith-ann (0).
func TestLexSortIndices_example(t *testing.T) {
	surname := []string{"smith", "jones", "smith", "jones"}
	first := []string{"ann", "bob", "al", "bob"}
	got, err := LexSortIndices(nil, [][]string{surname, first})
	if err != nil || !slices.Equal(got, []int{1, 3, 2, 0}) {
		t.Errorf("LexSortIndices = (%v, %v), want [1 3 2 0]", got, err)
	}
}

func TestLexSortIndices_cases(t *testing.T) {
	tests := []struct {
		name string
		keys [][]float64
		want []int
	}{
		{"no keys", nil, []int{}},
		{"empty records", [][]float64{{}, {}}, []int{}},
		{"one key", [][]float64{{3, 1, 2}}, []int{1, 2, 0}},
		{"second key breaks ties", [][]float64{{1, 1, 0}, {9, 8, 7}}, []int{2, 1, 0}},
		{"all equal stable", [][]float64{{1, 1, 1}, {2, 2, 2}}, []int{0, 1, 2}},
		{"NaN first", [][]float64{{1, math.NaN(), 0}}, []int{1, 2, 0}},
		{"signed zeros equal", [][]float64{{0, math.Copysign(0, -1)}, {1, 0}}, []int{1, 0}},
	}
	for _, tc := range tests {
		got, err := LexSortIndices(nil, tc.keys)
		if err != nil || !slices.Equal(got, tc.want) {
			t.Errorf("%s: = (%v, %v), want %v", tc.name, got, err, tc.want)
		}
	}
}

func TestLexSortIndices_errors(t *testing.T) {
	if _, err := LexSortIndices(nil, [][]int{{1, 2}, {1}}); !errors.Is(err, ErrLength) {
		t.Errorf("err = %v, want ErrLength", err)
	}
}

// Oracle: the result equals slices.SortStableFunc over the indexes with the same comparison, on 300 inputs.
func TestLexSortIndices_oracle(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for range 300 {
		n, k := r.IntN(80), 1+r.IntN(3)
		keys := make([][]int, k)
		for j := range keys {
			keys[j] = make([]int, n)
			for i := range keys[j] {
				keys[j][i] = r.IntN(4)
			}
		}
		got, err := LexSortIndices(nil, keys)
		if err != nil {
			t.Fatal(err)
		}
		want := make([]int, n)
		for i := range want {
			want[i] = i
		}
		slices.SortStableFunc(want, func(a, b int) int {
			for _, key := range keys {
				if key[a] != key[b] {
					return key[a] - key[b]
				}
			}
			return 0
		})
		if !slices.Equal(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
}

func BenchmarkLexSortIndices(b *testing.B) {
	for _, n := range []int{1e3, 1e4, 1e5} {
		r := rand.New(rand.NewPCG(3, 4))
		keys := [][]int{make([]int, n), make([]int, n)}
		for i := range n {
			keys[0][i], keys[1][i] = r.IntN(100), r.Int()
		}
		dst := make([]int, n)
		b.Run(strconv.Itoa(n), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = LexSortIndices(dst, keys)
			}
		})
	}
}

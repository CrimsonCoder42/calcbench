package heap

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Worked example written for this repo: the 3×3 matrix with rows 1 4 7 / 2 5 8 / 3 6 9 has every row
// and column ascending and holds 1…9, so its k-th smallest is k (by inspection).
func TestKthSmallestInSortedMatrix_workedExample(t *testing.T) {
	m := [][]int{{1, 4, 7}, {2, 5, 8}, {3, 6, 9}}
	for k := 1; k <= 9; k++ {
		if got, err := KthSmallestInSortedMatrix(m, k); err != nil || got != k {
			t.Errorf("k = %d: (%d, %v), want (%d, nil)", k, got, err, k)
		}
	}
}

func TestKthSmallestInSortedMatrix_cases(t *testing.T) {
	tests := []struct {
		name   string
		matrix [][]int
		k      int
		want   int
	}{
		{"one element", [][]int{{5}}, 1, 5},
		{"one row", [][]int{{1, 3, 5, 7}}, 3, 5},
		{"one column", [][]int{{1}, {3}, {5}}, 2, 3},
		{"two by two, last", [][]int{{1, 2}, {3, 4}}, 4, 4},
		{"duplicates count separately", [][]int{{1, 2, 2}, {2, 2, 3}}, 4, 2},
		{"duplicates, next value", [][]int{{1, 2, 2}, {2, 2, 3}}, 6, 3},
		{"all equal", [][]int{{7, 7}, {7, 7}}, 3, 7},
		{"negative and zero", [][]int{{-9, -1}, {-4, 0}}, 2, -4},
		{"extremes", [][]int{{math.MinInt, 0}, {0, math.MaxInt}}, 4, math.MaxInt},
		{"k-th lies in a late row", [][]int{{1, 100}, {2, 101}, {3, 102}}, 3, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := KthSmallestInSortedMatrix(tt.matrix, tt.k)
			if err != nil || got != tt.want {
				t.Errorf("(%v, %d) = (%d, %v), want (%d, nil)", tt.matrix, tt.k, got, err, tt.want)
			}
		})
	}
}

// Problem-site shape, in our own words: a square matrix with a repeated value straddling rows, asked for
// a rank in its middle; and a 1×1 matrix asked for its only element.
func TestKthSmallestInSortedMatrix_problemSiteShapes(t *testing.T) {
	m := [][]int{{2, 6, 10}, {6, 11, 14}, {12, 13, 16}}
	// Sorted: 2 6 6 10 11 12 13 14 16.
	if got, _ := KthSmallestInSortedMatrix(m, 6); got != 12 {
		t.Errorf("k 6 = %d, want 12", got)
	}
	if got, _ := KthSmallestInSortedMatrix([][]int{{-3}}, 1); got != -3 {
		t.Errorf("1×1 = %d, want -3", got)
	}
}

func TestKthSmallestInSortedMatrix_floatsAndStrings(t *testing.T) {
	f := [][]float64{{math.NaN(), 0.5}, {math.Inf(-1), 1}}
	// cmp order: NaN < -Inf < 0.5 < 1. Columns: NaN < -Inf, 0.5 < 1 (ascending under D-002).
	if got, err := KthSmallestInSortedMatrix(f, 2); err != nil || !math.IsInf(got, -1) {
		t.Errorf("floats k 2 = (%v, %v), want -Inf", got, err)
	}
	s := [][]string{{"ant", "bee"}, {"cat", "dog"}}
	if got, _ := KthSmallestInSortedMatrix(s, 3); got != "cat" {
		t.Errorf("strings k 3 = %q, want cat", got)
	}
}

func TestKthSmallestInSortedMatrix_errors(t *testing.T) {
	m := [][]int{{1, 2}, {3, 4}}
	for _, k := range []int{0, -1, 5, math.MaxInt, math.MinInt} {
		if _, err := KthSmallestInSortedMatrix(m, k); !errors.Is(err, ErrRank) {
			t.Errorf("k = %d: err = %v, want ErrRank", k, err)
		}
	}
	for _, empty := range [][][]int{nil, {}, {{}}, {{}, {}}} {
		if _, err := KthSmallestInSortedMatrix(empty, 1); !errors.Is(err, ErrRank) {
			t.Errorf("empty %v: err = %v, want ErrRank", empty, err)
		}
	}
	if _, err := KthSmallestInSortedMatrix([][]int{{1, 2}, {3}}, 1); !errors.Is(err, ErrRagged) {
		t.Errorf("ragged: err = %v, want ErrRagged", err)
	}
	if _, err := KthSmallestInSortedMatrix([][]int{{1}, {2, 3}}, 1); !errors.Is(err, ErrRagged) {
		t.Errorf("ragged longer: err = %v, want ErrRagged", err)
	}
}

func TestKthSmallestInSortedMatrix_unsortedDoesNotPanic(t *testing.T) {
	m := [][]int{{9, 1}, {0, 5}}
	for k := 1; k <= 4; k++ {
		if _, err := KthSmallestInSortedMatrix(m, k); err != nil {
			t.Errorf("k %d: err = %v", k, err)
		}
	}
}

// sortedMatrix builds an r×c matrix with ascending rows and columns: row base plus column base, both
// ascending, plus a little noise kept monotone by a running max.
func sortedMatrix(rng *rand.Rand, r, c int) [][]int {
	m := make([][]int, r)
	for i := range m {
		m[i] = make([]int, c)
		for j := range m[i] {
			v := rng.IntN(4)
			if i > 0 {
				v = max(v, m[i-1][j]+rng.IntN(3))
			}
			if j > 0 {
				v = max(v, m[i][j-1]+rng.IntN(3))
			}
			m[i][j] = v
		}
	}
	return m
}

func TestKthSmallestInSortedMatrix_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for trial := range 400 {
		r, c := 1+rng.IntN(10), 1+rng.IntN(10)
		m := sortedMatrix(rng, r, c)
		var flat []int
		for _, row := range m {
			flat = append(flat, row...)
		}
		slices.Sort(flat)
		k := 1 + rng.IntN(r*c)
		got, err := KthSmallestInSortedMatrix(m, k)
		if err != nil || got != flat[k-1] {
			t.Fatalf("trial %d: k %d of %v = (%d, %v), want %d", trial, k, m, got, err, flat[k-1])
		}
	}
}

// Square matrices of n = 1e3, 1e4, 1e5 cells (side √n), k = n/2: the median, the costly middle rank.
func BenchmarkKthSmallestInSortedMatrix(b *testing.B) {
	for _, size := range benchSizes {
		side := int(math.Sqrt(float64(size.n)))
		m := sortedMatrix(rand.New(rand.NewPCG(1, 2)), side, side)
		k := side * side / 2
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				_, _ = KthSmallestInSortedMatrix(m, k)
			}
		})
	}
}

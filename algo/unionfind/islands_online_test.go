package unionfind

import (
	"errors"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// Hand-derived example (catalog/unionfind.md, ALG-UNION-008): a 3×3 grid; (1,1) joins the three
// separate groups, and re-adding (0,0) changes nothing.
func TestIslandsOnline_handExample(t *testing.T) {
	adds := []Cell{{0, 0}, {0, 1}, {1, 2}, {2, 1}, {1, 1}, {0, 0}}
	got, err := IslandsOnline(nil, 3, 3, adds)
	if want := []int{1, 1, 2, 3, 1, 1}; err != nil || !slices.Equal(got, want) {
		t.Errorf("IslandsOnline = (%v, %v), want (%v, nil)", got, err, want)
	}
}

func TestIslandsOnline_cases(t *testing.T) {
	tests := []struct {
		name       string
		rows, cols int
		adds       []Cell
		want       []int
	}{
		{"no additions", 4, 4, nil, []int{}},
		{"empty grid, no additions", 0, 0, nil, []int{}},
		{"1×1", 1, 1, []Cell{{0, 0}}, []int{1}},
		{"same cell three times", 2, 2, []Cell{{1, 1}, {1, 1}, {1, 1}}, []int{1, 1, 1}},
		{"a row grows one island", 1, 4, []Cell{{0, 0}, {0, 2}, {0, 1}, {0, 3}}, []int{1, 2, 1, 1}},
		{"a column", 3, 1, []Cell{{2, 0}, {0, 0}, {1, 0}}, []int{1, 2, 1}},
		{"diagonals stay separate", 3, 3, []Cell{{0, 0}, {1, 1}, {2, 2}, {0, 2}, {2, 0}}, []int{1, 2, 3, 4, 5}},
		{"ring then centre", 3, 3, []Cell{{0, 1}, {1, 0}, {1, 2}, {2, 1}, {1, 1}}, []int{1, 2, 3, 4, 1}},
		{"huge grid, few cells", math.MaxInt, math.MaxInt, []Cell{{math.MaxInt - 1, 0}, {math.MaxInt - 1, 1}}, []int{1, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IslandsOnline(nil, tt.rows, tt.cols, tt.adds)
			if err != nil || got == nil || !slices.Equal(got, tt.want) {
				t.Errorf("IslandsOnline(%d, %d, %v) = (%v, %v), want (%v, nil)", tt.rows, tt.cols, tt.adds, got, err, tt.want)
			}
		})
	}
}

func TestIslandsOnline_dstReused(t *testing.T) {
	dst := make([]int, 5, 10)
	got, err := IslandsOnline(dst, 2, 2, []Cell{{0, 0}, {1, 1}})
	if err != nil || !slices.Equal(got, []int{1, 2}) || &got[0] != &dst[0] {
		t.Errorf("IslandsOnline(dst) = (%v, %v), want [1 2] in dst's array", got, err)
	}
}

func TestIslandsOnline_errors(t *testing.T) {
	if _, err := IslandsOnline(nil, -1, 3, nil); !errors.Is(err, ErrNegativeSize) {
		t.Errorf("rows -1: err = %v, want ErrNegativeSize", err)
	}
	if _, err := IslandsOnline(nil, 3, math.MinInt, nil); !errors.Is(err, ErrNegativeSize) {
		t.Errorf("cols MinInt: err = %v, want ErrNegativeSize", err)
	}
	for _, c := range []Cell{{-1, 0}, {0, -1}, {3, 0}, {0, 3}, {math.MinInt, math.MaxInt}} {
		if got, err := IslandsOnline(nil, 3, 3, []Cell{{1, 1}, c}); !errors.Is(err, ErrOutOfRange) || got != nil {
			t.Errorf("addition %v: (%v, %v), want (nil, ErrOutOfRange)", c, got, err)
		}
	}
	if _, err := IslandsOnline(nil, 0, 0, []Cell{{0, 0}}); !errors.Is(err, ErrOutOfRange) {
		t.Errorf("0×0 grid with an addition: err = %v, want ErrOutOfRange", err)
	}
}

// islandsByFlood is the oracle: count 4-connected land groups with a fresh flood fill.
func islandsByFlood(land [][]bool) int {
	rows, cols := len(land), len(land[0])
	seen := make([][]bool, rows)
	for r := range seen {
		seen[r] = make([]bool, cols)
	}
	count := 0
	for r := range rows {
		for c := range cols {
			if !land[r][c] || seen[r][c] {
				continue
			}
			count++
			stack := []Cell{{r, c}}
			seen[r][c] = true
			for len(stack) > 0 {
				p := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				for _, d := range []Cell{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
					q := Cell{p.Row + d.Row, p.Col + d.Col}
					if q.Row >= 0 && q.Row < rows && q.Col >= 0 && q.Col < cols && land[q.Row][q.Col] && !seen[q.Row][q.Col] {
						seen[q.Row][q.Col] = true
						stack = append(stack, q)
					}
				}
			}
		}
	}
	return count
}

func TestIslandsOnline_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for run := 0; run < 200; run++ {
		rows, cols := 1+rng.IntN(8), 1+rng.IntN(8)
		adds := make([]Cell, rng.IntN(rows*cols*2))
		for i := range adds {
			adds[i] = Cell{rng.IntN(rows), rng.IntN(cols)}
		}
		got, err := IslandsOnline(nil, rows, cols, adds)
		if err != nil {
			t.Fatal(err)
		}
		land := make([][]bool, rows)
		for r := range land {
			land[r] = make([]bool, cols)
		}
		for i, c := range adds {
			land[c.Row][c.Col] = true
			if want := islandsByFlood(land); got[i] != want {
				t.Fatalf("run %d step %d: %d islands, want %d", run, i, got[i], want)
			}
		}
	}
}

func BenchmarkIslandsOnline(b *testing.B) {
	const side = 1000
	for _, size := range benchSizes {
		rng := rand.New(rand.NewPCG(1, 2))
		adds := make([]Cell, size.n)
		for i := range adds {
			adds[i] = Cell{rng.IntN(side), rng.IntN(side)}
		}
		dst := make([]int, 0, size.n)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				dst, _ = IslandsOnline(dst, side, side, adds)
			}
		})
	}
}

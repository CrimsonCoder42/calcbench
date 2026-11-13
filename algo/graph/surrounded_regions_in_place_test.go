package graph

import (
	"errors"
	"math/rand/v2"
	"slices"
	"testing"
)

// Hand example (catalog/graph.md): the inner O block is enclosed under 4 neighbours; under 8 neighbours
// (2,2) touches the border O at (3,1) diagonally, so nothing is enclosed.
func TestSurroundedRegionsInPlace_handExample(t *testing.T) {
	board := func() [][]byte { return toGrid("XXXX", "XOOX", "XXOX", "XOXX") }
	grid := board()
	n, err := SurroundedRegionsInPlace(grid, 'O', 'X', FourNeighbors)
	if want := toGrid("XXXX", "XXXX", "XXXX", "XOXX"); n != 3 || err != nil || !slices.EqualFunc(grid, want, slices.Equal) {
		t.Errorf("4 neighbours: %d %v %q, want 3 %q", n, err, grid, want)
	}
	grid = board()
	if n, _ := SurroundedRegionsInPlace(grid, 'O', 'X', EightNeighbors); n != 0 || !slices.EqualFunc(grid, board(), slices.Equal) {
		t.Errorf("8 neighbours: %d %q, want 0 and unchanged", n, grid)
	}
}

func TestSurroundedRegionsInPlace_cases(t *testing.T) {
	tests := []struct {
		name string
		grid [][]byte
		want [][]byte
		n    int
	}{
		{"no rows", [][]byte{}, [][]byte{}, 0},
		{"1×1", toGrid("O"), toGrid("O"), 0},
		{"all border", toGrid("OO", "OO"), toGrid("OO", "OO"), 0},
		{"one enclosed cell", toGrid("XXX", "XOX", "XXX"), toGrid("XXX", "XXX", "XXX"), 1},
		{"other values are walls too", toGrid("ABA", "BOB", "AOA"), toGrid("ABA", "BOB", "AOA"), 0},
		{"two enclosed pockets", toGrid("XXXXX", "XOXOX", "XXXXX"), toGrid("XXXXX", "XXXXX", "XXXXX"), 2},
	}
	for _, tt := range tests {
		n, err := SurroundedRegionsInPlace(tt.grid, 'O', 'X', FourNeighbors)
		if n != tt.n || err != nil || !slices.EqualFunc(tt.grid, tt.want, slices.Equal) {
			t.Errorf("%s: %d %v %q, want %d %q", tt.name, n, err, tt.grid, tt.n, tt.want)
		}
	}
	same := toGrid("XXX", "XOX", "XXX")
	if n, _ := SurroundedRegionsInPlace(same, 'O', 'O', FourNeighbors); n != 0 || same[1][1] != 'O' {
		t.Error("region == fill must change nothing")
	}
}

func TestSurroundedRegionsInPlace_errors(t *testing.T) {
	if _, err := SurroundedRegionsInPlace(toGrid("O"), 'O', 'X', Neighborhood(3)); !errors.Is(err, ErrNeighborhood) {
		t.Errorf("neighborhood: %v", err)
	}
	ragged := toGrid("XXX", "XO")
	if _, err := SurroundedRegionsInPlace(ragged, 'O', 'X', FourNeighbors); !errors.Is(err, ErrNotRectangular) || ragged[1][1] != 'O' {
		t.Errorf("ragged: %v (grid must be unchanged)", err)
	}
}

// Oracle: a region cell is filled exactly when GridBFS from it reaches no border cell.
func TestSurroundedRegionsInPlace_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for run := 0; run < 300; run++ {
		rows, cols := 1+rng.IntN(8), 1+rng.IntN(8)
		grid := randomGrid(rng, rows, cols, 0.6)
		nb := FourNeighbors
		if run%2 == 1 {
			nb = EightNeighbors
		}
		enclosed := make(map[Cell]bool)
		for r := range rows {
			for c := range cols {
				if !isDot(grid[r][c]) {
					continue
				}
				touches := false
				for x := range must(GridBFS(grid, Cell{Row: r, Col: c}, isDot, nb)) {
					touches = touches || x.Row == 0 || x.Col == 0 || x.Row == rows-1 || x.Col == cols-1
				}
				enclosed[Cell{Row: r, Col: c}] = !touches
			}
		}
		before := make([][]byte, rows)
		for r := range grid {
			before[r] = slices.Clone(grid[r])
		}
		n, err := SurroundedRegionsInPlace(grid, '.', '#', nb)
		want := 0
		for r := range rows {
			for c := range cols {
				expect := before[r][c]
				if enclosed[Cell{Row: r, Col: c}] {
					expect = '#'
					want++
				}
				if grid[r][c] != expect {
					t.Fatalf("run %d: cell (%d,%d) = %c, want %c", run, r, c, grid[r][c], expect)
				}
			}
		}
		if n != want || err != nil {
			t.Fatalf("run %d: changed %d (%v), want %d", run, n, err, want)
		}
	}
}

func BenchmarkSurroundedRegionsInPlace(b *testing.B) {
	for _, size := range gridBenchSides {
		src := randomGrid(rand.New(rand.NewPCG(1, 2)), size.side, size.side, 0.5)
		grid := make([][]byte, len(src))
		for r := range grid {
			grid[r] = make([]byte, len(src[r]))
		}
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				for r := range grid {
					copy(grid[r], src[r]) // restore: O(R·C), the same order as the call itself
				}
				_, _ = SurroundedRegionsInPlace(grid, '.', '#', FourNeighbors)
			}
		})
	}
}

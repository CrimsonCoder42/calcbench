package graph

import (
	"errors"
	"math/rand/v2"
	"slices"
	"testing"
)

// Hand example (catalog/graph.md): DFS-VISIT order (CLRS §22.3) on the grid graph, 4 and 8 neighbours.
func TestGridDFS_handExample(t *testing.T) {
	grid := toGrid(
		"..#",
		".##",
		"...",
	)
	if got, want := slices.Collect(must(GridDFS(grid, Cell{}, isDot, FourNeighbors))), cells(0, 0, 0, 1, 1, 0, 2, 0, 2, 1, 2, 2); !slices.Equal(got, want) {
		t.Errorf("GridDFS 4 = %v, want %v", got, want)
	}
	if got, want := slices.Collect(must(GridDFS(grid, Cell{}, isDot, EightNeighbors))), cells(0, 0, 0, 1, 1, 0, 2, 1, 2, 2, 2, 0); !slices.Equal(got, want) {
		t.Errorf("GridDFS 8 = %v, want %v", got, want)
	}
}

func TestGridDFS_cases(t *testing.T) {
	if got := slices.Collect(must(GridDFS(toGrid("#."), Cell{}, isDot, FourNeighbors))); len(got) != 0 {
		t.Errorf("blocked start: %v", got)
	}
	// A 1×200000 corridor: no recursion, so no stack growth with the path length.
	row := make([]byte, 200_000)
	for i := range row {
		row[i] = '.'
	}
	n := 0
	for range must(GridDFS([][]byte{row}, Cell{}, isDot, FourNeighbors)) {
		n++
	}
	if n != len(row) {
		t.Errorf("corridor: %d cells, want %d", n, len(row))
	}
}

func TestGridDFS_errors(t *testing.T) {
	grid := toGrid("..")
	if _, err := GridDFS(grid, Cell{}, isDot, Neighborhood(0)); !errors.Is(err, ErrNeighborhood) {
		t.Errorf("neighborhood 0: %v", err)
	}
	if _, err := GridDFS(toGrid(".", ".."), Cell{}, isDot, FourNeighbors); !errors.Is(err, ErrNotRectangular) {
		t.Errorf("ragged: %v", err)
	}
	if _, err := GridDFS(grid, Cell{Col: 2}, isDot, FourNeighbors); !errors.Is(err, ErrCellOutOfRange) {
		t.Errorf("start out of range: %v", err)
	}
	if _, err := GridDFS(grid, Cell{}, nil, FourNeighbors); !errors.Is(err, ErrNilPredicate) {
		t.Errorf("nil predicate: %v", err)
	}
}

// Oracle: GridDFS order equals the iterative DFS on the explicit grid graph (itself checked against a
// recursive DFS in TestDFS_oracle).
func TestGridDFS_oracle(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 2))
	for run := 0; run < 300; run++ {
		rows, cols := 1+rng.IntN(8), 1+rng.IntN(8)
		grid := randomGrid(rng, rows, cols, 0.6)
		nb := FourNeighbors
		if run%2 == 1 {
			nb = EightNeighbors
		}
		start := Cell{Row: rng.IntN(rows), Col: rng.IntN(cols)}
		var got []int
		for c := range must(GridDFS(grid, start, isDot, nb)) {
			got = append(got, c.Row*cols+c.Col)
		}
		var want []int
		if isDot(grid[start.Row][start.Col]) {
			g := gridAsGraph(grid, nb)
			want = recursiveDFS(g, start.Row*cols+start.Col, make([]bool, rows*cols), nil)
		}
		if !slices.Equal(got, want) {
			t.Fatalf("run %d: GridDFS %v, recursive DFS %v", run, got, want)
		}
	}
}

func BenchmarkGridDFS(b *testing.B) {
	for _, size := range gridBenchSides {
		grid := openGrid(size.side)
		b.Run(size.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				seq, _ := GridDFS(grid, Cell{}, isDot, FourNeighbors)
				for range seq {
				}
			}
		})
	}
}
